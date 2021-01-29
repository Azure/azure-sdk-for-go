package azblob

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"

	"bytes"
	"errors"
	"os"
	//"github.com/Azure/azure-pipeline-go/pipeline"
)

//// CommonResponse returns the headers common to all blob REST API responses.
type CommonResponse interface {
	//// ETag returns the value for header ETag.
	//ETag() *string
	//
	//// LastModified returns the value for header Last-Modified.
	//LastModified() time.Time
	//
	//// RequestID returns the value for header x-ms-request-id.
	//RequestID() string

	//// Date returns the value for header Date.
	//Date() time.Time
	//
	//// Version returns the value for header x-ms-version.
	//Version() string

	// Response returns the raw HTTP response object.
	Response() *http.Response
}

// HighLevelUploadToBlockBlobOption identifies options used by the UploadBufferToBlockBlob and UploadFileToBlockBlob functions.
type HighLevelUploadToBlockBlobOption struct {
	// BlockSize specifies the block size to use; the default (and maximum size) is BlockBlobMaxStageBlockBytes.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are sent to the BlockBlobClient.
	// Note that the progress reporting is not always increasing; it can go down when retrying a request.
	//Progress pipeline.ProgressReceiver

	// BlobHTTPHeaders indicates the HTTP headers to be associated with the blob.
	BlobHTTPHeaders *BlobHttpHeaders

	// Metadata indicates the metadata to be associated with the blob when PutBlockList is called.
	Metadata *map[string]string

	// AccessConditions indicates the access conditions for the block blob.
	ModifiedAccessCondition *ModifiedAccessConditions
	LeaseAccessCondition    *LeaseAccessConditions

	// BlobAccessTier indicates the tier of blob
	BlobAccessTier *AccessTier

	// BlobTagsMap
	BlobTagsMap *map[string]string

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CpkInfo      *CpkInfo
	CpkScopeInfo *CpkScopeInfo

	// Parallelism indicates the maximum number of blocks to upload in parallel (0=default)
	Parallelism uint16
}

func (o HighLevelUploadToBlockBlobOption) getStageBlockOptions() *StageBlockOptions {
	return &StageBlockOptions{
		CpkInfo:               o.CpkInfo,
		CpkScopeInfo:          o.CpkScopeInfo,
		LeaseAccessConditions: o.LeaseAccessCondition,
		//BlockBlobStageBlockOptions: &BlockBlobStageBlockOptions{
		//	RequestId: nil,
		//	Timeout: nil,
		//	TransactionalContentMd5: nil,
		//	TransactionalContentCrc64: nil,
		//},
	}
}

func (o HighLevelUploadToBlockBlobOption) getUploadBlockBlobOptions() *UploadBlockBlobOptions {
	return &UploadBlockBlobOptions{
		BlobTagsMap: o.BlobTagsMap,
		Metadata:    o.Metadata,
		Tier:        o.BlobAccessTier,
		//TransactionalContentMd5:
		BlobHttpHeaders:          o.BlobHTTPHeaders,
		LeaseAccessConditions:    o.LeaseAccessCondition,
		ModifiedAccessConditions: o.ModifiedAccessCondition,
		CpkInfo:                  o.CpkInfo,
		CpkScopeInfo:             o.CpkScopeInfo,
	}
}

func (o *HighLevelUploadToBlockBlobOption) getCommitBlockListOptions() *CommitBlockListOptions {
	return &CommitBlockListOptions{
		BlobTagsMap:     o.BlobTagsMap,
		Metadata:        o.Metadata,
		Tier:            o.BlobAccessTier,
		BlobHTTPHeaders: o.BlobHTTPHeaders,
		CpkInfo:         o.CpkInfo,
		CpkScopeInfo:    o.CpkScopeInfo,
	}
}

// uploadReaderAtToBlockBlob uploads a buffer in blocks to a block blob.
func uploadReaderAtToBlockBlob(ctx context.Context, reader io.ReaderAt, readerSize int64,
	blockBlobClient BlockBlobClient, o HighLevelUploadToBlockBlobOption) (CommonResponse, error) {
	if o.BlockSize == 0 {
		// If bufferSize > (BlockBlobMaxStageBlockBytes * BlockBlobMaxBlocks), then error
		if readerSize > BlockBlobMaxStageBlockBytes*BlockBlobMaxBlocks {
			return BlockBlobUploadResponse{}, errors.New("buffer is too large to upload to a block blob")
		}
		// If bufferSize <= BlockBlobMaxUploadBlobBytes, then Upload should be used with just 1 I/O request
		if readerSize <= BlockBlobMaxUploadBlobBytes {
			o.BlockSize = BlockBlobMaxUploadBlobBytes // Default if unspecified
		} else {
			o.BlockSize = readerSize / BlockBlobMaxBlocks   // buffer / max blocks = block size to use all 50,000 blocks
			if o.BlockSize < BlobDefaultDownloadBlockSize { // If the block size is smaller than 4MB, round up to 4MB
				o.BlockSize = BlobDefaultDownloadBlockSize
			}
			// StageBlock will be called with blockSize blocks and a Parallelism of (BufferSize / BlockSize).
		}
	}

	if readerSize <= BlockBlobMaxUploadBlobBytes {
		// If the size can fit in 1 Upload call, do it this way
		var body io.ReadSeeker = io.NewSectionReader(reader, 0, readerSize)
		//if o.Progress != nil {
		//	body = pipeline.NewRequestBodyProgress(body, o.Progress)
		//}

		uploadBlockBlobOptions := o.getUploadBlockBlobOptions()
		return blockBlobClient.Upload(ctx, body, uploadBlockBlobOptions)
	}

	var numBlocks = uint16(((readerSize - 1) / o.BlockSize) + 1)

	blockIDList := make([]string, numBlocks) // Base-64 encoded block IDs
	//progress := int64(0)
	//progressLock := &sync.Mutex{}

	err := DoBatchTransfer(ctx, BatchTransferOptions{
		OperationName: "uploadReaderAtToBlockBlob",
		TransferSize:  readerSize,
		ChunkSize:     o.BlockSize,
		Parallelism:   o.Parallelism,
		Operation: func(offset int64, count int64, ctx context.Context) error {
			// This function is called once per block.
			// It is passed this block's offset within the buffer and its count of bytes
			// Prepare to read the proper block/section of the buffer
			var body io.ReadSeeker = io.NewSectionReader(reader, offset, count)
			blockNum := offset / o.BlockSize
			//if o.Progress != nil {
			//	blockProgress := int64(0)
			//	body = pipeline.NewRequestBodyProgress(body,
			//		func(bytesTransferred int64) {
			//			diff := bytesTransferred - blockProgress
			//			blockProgress = bytesTransferred
			//			progressLock.Lock() // 1 goroutine at a time gets a progress report
			//			progress += diff
			//			o.Progress(progress)
			//			progressLock.Unlock()
			//		})
			//}

			// Block IDs are unique values to avoid issue if 2+ clients are uploading blocks
			// at the same time causing PutBlockList to get a mix of blocks from all the clients.
			blockIDList[blockNum] = base64.StdEncoding.EncodeToString(newUUID().bytes())
			stageBlockOptions := o.getStageBlockOptions()
			_, err := blockBlobClient.StageBlock(ctx, blockIDList[blockNum], body, stageBlockOptions)
			return err
		},
	})
	if err != nil {
		return BlockBlobUploadResponse{}, err
	}
	// All put blocks were successful, call Put Block List to finalize the blob
	commitBlockListOptions := o.getCommitBlockListOptions()
	return blockBlobClient.CommitBlockList(ctx, blockIDList, commitBlockListOptions)
}

// UploadBufferToBlockBlob uploads a buffer in blocks to a block blob.
func UploadBufferToBlockBlob(ctx context.Context, b []byte,
	blockBlobClient BlockBlobClient, o HighLevelUploadToBlockBlobOption) (CommonResponse, error) {
	return uploadReaderAtToBlockBlob(ctx, bytes.NewReader(b), int64(len(b)), blockBlobClient, o)
}

// UploadFileToBlockBlob uploads a file in blocks to a block blob.
func UploadFileToBlockBlob(ctx context.Context, file *os.File,
	blockBlobClient BlockBlobClient, o HighLevelUploadToBlockBlobOption) (CommonResponse, error) {

	stat, err := file.Stat()
	if err != nil {
		return BlockBlobUploadResponse{}, err
	}
	return uploadReaderAtToBlockBlob(ctx, file, stat.Size(), blockBlobClient, o)
}

///////////////////////////////////////////////////////////////////////////////

const BlobDefaultDownloadBlockSize = int64(4 * 1024 * 1024) // 4MB

// HighLevelDownloadFromBlobOptions identifies options used by the DownloadBlobToBuffer and DownloadBlobToFile functions.
type HighLevelDownloadFromBlobOptions struct {
	// BlockSize specifies the block size to use for each parallel download; the default size is BlobDefaultDownloadBlockSize.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	//Progress pipeline.ProgressReceiver

	// AccessConditions indicates the access conditions used when making HTTP GET requests against the blob.
	ModifiedAccessConditions *ModifiedAccessConditions
	LeaseAccessConditions    *LeaseAccessConditions

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CpkInfo      *CpkInfo
	CpkScopeInfo *CpkScopeInfo

	// Parallelism indicates the maximum number of blocks to download in parallel (0=default)
	Parallelism uint16

	// RetryReaderOptionsPerBlock is used when downloading each block.
	RetryReaderOptionsPerBlock RetryReaderOptions
}

func (o *HighLevelDownloadFromBlobOptions) getBlobPropertiesOptions() *GetBlobPropertiesOptions {
	return &GetBlobPropertiesOptions{
		LeaseAccessConditions:    o.LeaseAccessConditions,
		ModifiedAccessConditions: o.ModifiedAccessConditions,
		CpkInfo:                  o.CpkInfo,
	}
}

func (o *HighLevelDownloadFromBlobOptions) getDownloadBlobOptions(offSet, count int64, rangeGetContentMd5 *bool) *DownloadBlobOptions {
	return &DownloadBlobOptions{
		LeaseAccessConditions:    o.LeaseAccessConditions,
		ModifiedAccessConditions: o.ModifiedAccessConditions,
		CpkInfo:                  o.CpkInfo,
		CpkScopeInfo:             o.CpkScopeInfo,
		Offset:                   &offSet,
		Count:                    &count,
		RangeGetContentMd5:       rangeGetContentMd5,
	}
}

// downloadBlobToWriterAt downloads an Azure blob to a buffer with parallel.
func downloadBlobToWriterAt(ctx context.Context, blobClient BlobClient, offset int64, count int64,
	writer io.WriterAt, o HighLevelDownloadFromBlobOptions, initialDownloadResponse *DownloadResponse) error {
	if o.BlockSize == 0 {
		o.BlockSize = BlobDefaultDownloadBlockSize
	}

	if count == CountToEnd { // If size not specified, calculate it
		if initialDownloadResponse != nil {
			count = *(initialDownloadResponse.ContentLength()) - offset // if we have the length, use it
		} else {
			// If we don't have the length at all, get it
			downloadBlobOptions := o.getDownloadBlobOptions(0, CountToEnd, nil)
			dr, err := blobClient.Download(ctx, downloadBlobOptions)
			if err != nil {
				return err
			}
			count = *(dr.ContentLength()) - offset
		}
	}

	if count <= 0 {
		// The file is empty, there is nothing to download.
		return nil
	}

	// Prepare and do parallel download.
	//progress := int64(0)
	//progressLock := &sync.Mutex{}

	err := DoBatchTransfer(ctx, BatchTransferOptions{
		OperationName: "downloadBlobToWriterAt",
		TransferSize:  count,
		ChunkSize:     o.BlockSize,
		Parallelism:   o.Parallelism,
		Operation: func(chunkStart int64, count int64, ctx context.Context) error {

			downloadBlobOptions := o.getDownloadBlobOptions(chunkStart+offset, count, nil)
			dr, err := blobClient.Download(ctx, downloadBlobOptions)
			if err != nil {
				return err
			}
			body := dr.Body(o.RetryReaderOptionsPerBlock)
			//if o.Progress != nil {
			//	rangeProgress := int64(0)
			//	body = pipeline.NewResponseBodyProgress(
			//		body,
			//		func(bytesTransferred int64) {
			//			diff := bytesTransferred - rangeProgress
			//			rangeProgress = bytesTransferred
			//			progressLock.Lock()
			//			progress += diff
			//			o.Progress(progress)
			//			progressLock.Unlock()
			//		})
			//}
			_, err = io.Copy(newSectionWriter(writer, chunkStart, count), body)
			err = body.Close()
			return err
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// DownloadBlobToBuffer downloads an Azure blob to a buffer with parallel.
// Offset and count are optional, pass 0 for both to download the entire blob.
func DownloadBlobToBuffer(ctx context.Context, blobClient BlobClient, offset int64, count int64,
	b []byte, o HighLevelDownloadFromBlobOptions) error {
	return downloadBlobToWriterAt(ctx, blobClient, offset, count, newBytesWriter(b), o, nil)
}

// DownloadBlobToFile downloads an Azure blob to a local file.
// The file would be truncated if the size doesn't match.
// Offset and count are optional, pass 0 for both to download the entire blob.
func DownloadBlobToFile(ctx context.Context, blobClient BlobClient, offset int64, count int64,
	file *os.File, o HighLevelDownloadFromBlobOptions) error {
	// 1. Calculate the size of the destination file
	var size int64

	if count == CountToEnd {
		// Try to get Azure blob's size
		getBlobPropertiesOptions := o.getBlobPropertiesOptions()
		props, err := blobClient.GetProperties(ctx, getBlobPropertiesOptions)
		if err != nil {
			return err
		}
		size = *props.ContentLength - offset
	} else {
		size = count
	}

	// 2. Compare and try to resize local file's size if it doesn't match Azure blob's size.
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() != size {
		if err = file.Truncate(size); err != nil {
			return err
		}
	}

	if size > 0 {
		return downloadBlobToWriterAt(ctx, blobClient, offset, size, file, o, nil)
	} else { // if the blob's size is 0, there is no need in downloading it
		return nil
	}
}

///////////////////////////////////////////////////////////////////////////////

// BatchTransferOptions identifies options used by DoBatchTransfer.
type BatchTransferOptions struct {
	TransferSize  int64
	ChunkSize     int64
	Parallelism   uint16
	Operation     func(offset int64, chunkSize int64, ctx context.Context) error
	OperationName string
}

// DoBatchTransfer helps to execute operations in a batch manner.
// Can be used by users to customize batch works (for other scenarios that the SDK does not provide)
func DoBatchTransfer(ctx context.Context, o BatchTransferOptions) error {
	if o.ChunkSize == 0 {
		return errors.New("ChunkSize cannot be 0")
	}

	if o.Parallelism == 0 {
		o.Parallelism = 5 // default Parallelism
	}

	// Prepare and do parallel operations.
	numChunks := uint16(((o.TransferSize - 1) / o.ChunkSize) + 1)
	operationChannel := make(chan func() error, o.Parallelism) // Create the channel that release 'Parallelism' goroutines concurrently
	operationResponseChannel := make(chan error, numChunks)    // Holds each response
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create the goroutines that process each operation (in parallel).
	for g := uint16(0); g < o.Parallelism; g++ {
		//grIndex := g
		go func() {
			for f := range operationChannel {
				err := f()
				operationResponseChannel <- err
			}
		}()
	}

	// Add each chunk's operation to the channel.
	for chunkNum := uint16(0); chunkNum < numChunks; chunkNum++ {
		curChunkSize := o.ChunkSize

		if chunkNum == numChunks-1 { // Last chunk
			curChunkSize = o.TransferSize - (int64(chunkNum) * o.ChunkSize) // Remove size of all transferred chunks from total
		}
		offset := int64(chunkNum) * o.ChunkSize

		operationChannel <- func() error {
			return o.Operation(offset, curChunkSize, ctx)
		}
	}
	close(operationChannel)

	// Wait for the operations to complete.
	var firstErr error = nil
	for chunkNum := uint16(0); chunkNum < numChunks; chunkNum++ {
		responseError := <-operationResponseChannel
		// record the first error (the original error which should cause the other chunks to fail with canceled context)
		if responseError != nil && firstErr == nil {
			cancel() // As soon as any operation fails, cancel all remaining operation calls
			firstErr = responseError
		}
	}
	return firstErr
}

////////////////////////////////////////////////////////////////////////////////////////////////

const _1MiB = 1024 * 1024

type UploadStreamToBlockBlobOptions struct {
	// BufferSize sizes the buffer used to read data from source. If < 1 MiB, defaults to 1 MiB.
	BufferSize int
	// MaxBuffers defines the number of simultaneous uploads will be performed to upload the file.
	MaxBuffers       int
	BlobHTTPHeaders  *BlobHttpHeaders
	Metadata         *map[string]string
	AccessConditions *BlobAccessConditions
	BlobAccessTier   *AccessTier
	BlobTagsMap      *map[string]string
	CpkInfo          *CpkInfo
	CpkScopeInfo     *CpkScopeInfo
}

func (u *UploadStreamToBlockBlobOptions) defaults() {
	if u.MaxBuffers == 0 {
		u.MaxBuffers = 1
	}

	if u.BufferSize < _1MiB {
		u.BufferSize = _1MiB
	}
}

func (u *UploadStreamToBlockBlobOptions) getStageBlockOptions() *StageBlockOptions {
	return &StageBlockOptions{}
}

func (u *UploadStreamToBlockBlobOptions) getCommitBlockListOptions() *CommitBlockListOptions {
	return &CommitBlockListOptions{}
}

// UploadStreamToBlockBlob copies the file held in io.Reader to the Blob at blockBlobClient.
// A Context deadline or cancellation will cause this to error.
func UploadStreamToBlockBlob(ctx context.Context, reader io.Reader, blockBlobClient BlockBlobClient,
	o UploadStreamToBlockBlobOptions) (BlockBlobCommitBlockListResponse, error) {
	o.defaults()

	result, err := copyFromReader(ctx, reader, blockBlobClient, o)
	if err != nil {
		return BlockBlobCommitBlockListResponse{}, err
	}

	return result, nil
}

// UploadStreamOptions (defunct) was used internally. This will be removed or made private in a future version.
type UploadStreamOptions struct {
	BufferSize int
	MaxBuffers int
}

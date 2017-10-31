package azblob

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// StreamToBlockBlobOptions identifies options used by the StreamToBlockBlob function. Note that the
// BlockSize field is mandatory and must be set; other fields are optional.
type StreamToBlockBlobOptions struct {
	// BlockSize is mandatory. It specifies the block size to use; the maximum size is BlockBlobMaxPutBlockBytes.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are send in a PutBlock call to the BlockBlobURL.
	Progress pipeline.ProgressReceiver

	// BlobHTTPHeaders indicates the HTTP headers to be associated with the blob when PutBlockList is called.
	BlobHTTPHeaders BlobHTTPHeaders

	// Metadata indicates the metadata to be associated with the blob when PutBlockList is called.
	Metadata Metadata
	// BlobAccessConditions???
}

// StreamToBlockBlob uploads a large stream of data in blocks to a block blob.
func StreamToBlockBlob(ctx context.Context, stream io.ReaderAt, streamSize int64,
	blockBlobURL BlockBlobURL, o StreamToBlockBlobOptions) (*BlockBlobsPutBlockListResponse, error) {

	if o.BlockSize <= 0 || o.BlockSize > BlockBlobMaxPutBlockBytes {
		panic(fmt.Sprintf("BlockSize option must be > 0 and <= %d", BlockBlobMaxPutBlockBytes))
	}

	numBlocks := ((streamSize - int64(1)) / o.BlockSize) + 1
	blockIDList := make([]string, numBlocks) // Base 64 encoded block IDs
	blockSize := o.BlockSize

	for blockNum := int64(0); blockNum < numBlocks; blockNum++ {
		if blockNum == numBlocks-1 { // Last block
			blockSize = streamSize % o.BlockSize
		}

		streamOffset := blockNum * o.BlockSize
		// Prepare to read the proper block/section of the file
		var body io.ReadSeeker = io.NewSectionReader(stream, streamOffset, blockSize)
		if o.Progress != nil {
			body = pipeline.NewRequestBodyProgress(body,
				func(bytesTransferred int64) { o.Progress(streamOffset + bytesTransferred) })
		}

		blockIDList[blockNum] = blockIDUint64ToBase64(uint64(streamOffset)) // The streamOffset is the block ID
		_, err := blockBlobURL.PutBlock(ctx, blockIDList[blockNum], body, LeaseAccessConditions{})
		if err != nil {
			return nil, err
		}
	}
	return blockBlobURL.PutBlockList(ctx, blockIDList, o.Metadata, o.BlobHTTPHeaders, BlobAccessConditions{})
}

// NOTE: The blockID must be <= 64 bytes and ALL blockIDs for the block must be the same length
// These helper functions convert an int64 block ID to a base-64 string
func blockIDUint64ToBase64(blockID uint64) string {
	binaryBlockID := [64 / 8]byte{} // All block IDs are 8 bytes long
	binary.LittleEndian.PutUint64(binaryBlockID[:], blockID)
	return base64.StdEncoding.EncodeToString(binaryBlockID[:])
}

// GetRetryStreamOptions is used to configure a call to NewGetTryStream to download a large stream with intelligent retries.
type GetRetryStreamOptions struct {
	// Range indicates the starting offset and count of bytes within the blob to download.
	Range BlobRange

	// Acc indicates the BlobAccessConditions to use when accessing the blob.
	AC BlobAccessConditions

	// GetBlobResult identifies a function to invoke immediately after GetRetryStream's Read method internally
	// calls GetBlob. This function is invoked after every call to GetBlob. The callback can example GetBlob's
	// response and error information.
	GetBlobResult func(*GetResponse, error)
}

type retryStream struct {
	ctx      context.Context
	blobURL  BlobURL
	o        GetRetryStreamOptions
	response *http.Response
}

// NewGetRetryStream creates a stream over a blob allowing you download the blob's contents.
// When network errors occur, the retry stream internally issues new HTTP GET requests for
// the remaining range of the blob's contents.
func NewGetRetryStream(ctx context.Context, blobURL BlobURL, o GetRetryStreamOptions) io.ReadCloser {
	// BlobAccessConditions may already have an If-Match:etag header
	return &retryStream{ctx: ctx, blobURL: blobURL, o: o, response: nil}
}

func (s *retryStream) Read(p []byte) (n int, err error) {
	for {
		if s.response != nil { // We working with a successful response
			n, err := s.response.Body.Read(p) // Read from the stream
			if err == nil || err == io.EOF {  // We successfully read data or end EOF
				s.o.Range.Offset += int64(n) // Increments the start offset in case we need to make a new HTTP request in the future
				if s.o.Range.Count != 0 {
					s.o.Range.Count -= int64(n) // Decrement the count in case we need to make a new HTTP request in the future
				}
				return n, err // Return the return to the caller
			}
			s.response = nil // Something went wrong; our stream is no longer good
			if nerr, ok := err.(net.Error); ok {
				if !nerr.Timeout() && !nerr.Temporary() {
					return n, err // Not retryable
				}
			} else {
				return n, err // Not retryable, just return
			}
		}

		// We don't have a response stream to read from, try to get one
		response, err := s.blobURL.GetBlob(s.ctx, s.o.Range, s.o.AC, false)
		if s.o.GetBlobResult != nil {
			// If caller desires notification of each GetBlob call, notify them
			s.o.GetBlobResult(response, err)
		}
		if err != nil {
			return 0, err
		}
		// Successful GET; this is the network stream we'll read from
		s.response = response.Response()

		// Ensure that future requests are from the same version of the source
		s.o.AC.IfMatch = response.ETag()

		// Loop around and try to read from this stream
	}
}

func (s *retryStream) Close() error {
	//s.blobURL = BlobURL{} // This blobURL is no longer valid
	return s.response.Body.Close()
}

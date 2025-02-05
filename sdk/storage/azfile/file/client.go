//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Storage file.
type Client base.Client[generated.FileClient]

// NewClient creates an instance of Client with the specified values.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
//
// Note that ClientOptions.FileRequestIntent is currently required for token authentication.
func NewClient(fileURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	audience := base.GetAudience((*base.ClientOptions)(options))
	conOptions := shared.GetClientOptions(options)
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{audience}, &policy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: conOptions.InsecureAllowCredentialWithHTTP,
	})
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewFileClient(fileURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a file or with a shared access signature (SAS) token.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file?<sas token>
//   - options - client options; pass nil to accept the default values
//
// The directoryPath is optional in the fileURL. If omitted, it points to file within the specified share.
func NewClientWithNoCredential(fileURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewFileClient(fileURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file
//   - cred - a SharedKeyCredential created with the matching file's storage account and access key
//   - options - client options; pass nil to accept the default values
//
// The directoryPath is optional in the fileURL. If omitted, it points to file within the specified share.
func NewClientWithSharedKeyCredential(fileURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	return (*Client)(base.NewFileClient(fileURL, azClient, cred, (*base.ClientOptions)(conOptions))), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - shareName - the name of the share within the storage account
//   - filePath - the path of the file within the share
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, shareName string, filePath string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	filePath = strings.ReplaceAll(filePath, "\\", "/")
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, shareName, filePath)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (f *Client) generated() *generated.FileClient {
	return base.InnerClient((*base.Client[generated.FileClient])(f))
}

func (f *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.FileClient])(f))
}

func (f *Client) getClientOptions() *base.ClientOptions {
	return base.GetClientOptions((*base.Client[generated.FileClient])(f))
}

// URL returns the URL endpoint used by the Client object.
func (f *Client) URL() string {
	return f.generated().Endpoint()
}

// Create operation creates a new file or replaces a file. Note it only initializes the file with no content.
//   - fileContentLength: Specifies the maximum size for the file in bytes, up to 4 TB.
//
// ParseNTFSFileAttributes method can be used to convert the file attributes returned in response to NTFSFileAttributes.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-file.
func (f *Client) Create(ctx context.Context, fileContentLength int64, options *CreateOptions) (CreateResponse, error) {
	fileCreateOptions, fileHTTPHeaders, leaseAccessConditions := options.format()
	resp, err := f.generated().Create(ctx, fileContentLength, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions)
	return resp, err
}

// Delete operation removes the file from the storage account.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-file2.
func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().Delete(ctx, opts, leaseAccessConditions)
	return resp, err
}

// Rename operation renames a file, and can optionally set system properties for the file.
//   - destinationPath: the destination path to rename the file to.
//
// For more information, see https://learn.microsoft.com/rest/api/storageservices/rename-file.
func (f *Client) Rename(ctx context.Context, destinationPath string, options *RenameOptions) (RenameResponse, error) {
	destinationPath = strings.Trim(strings.TrimSpace(destinationPath), "/")
	if len(destinationPath) == 0 {
		return RenameResponse{}, errors.New("destination path must not be empty")
	}

	opts, srcLease, destLease, smbInfo, fileHTTPHeaders := options.format()

	urlParts, err := sas.ParseURL(f.URL())
	if err != nil {
		return RenameResponse{}, err
	}

	destParts := strings.Split(destinationPath, "?")
	newDestPath := destParts[0]
	newDestQuery := ""
	if len(destParts) == 2 {
		newDestQuery = destParts[1]
	}

	urlParts.DirectoryOrFilePath = newDestPath
	destURL := urlParts.String()
	// replace the query part if it is present in destination path
	if len(newDestQuery) > 0 {
		destURL = strings.Split(destURL, "?")[0] + "?" + newDestQuery
	}

	destFileClient := (*Client)(base.NewFileClient(destURL, f.generated().InternalClient(), f.sharedKey(), f.getClientOptions()))

	resp, err := destFileClient.generated().Rename(ctx, f.URL(), opts, srcLease, destLease, smbInfo, fileHTTPHeaders)
	return RenameResponse{
		FileClientRenameResponse: resp,
	}, err
}

// GetProperties operation returns all user-defined metadata, standard HTTP properties, and system properties for the file.
// ParseNTFSFileAttributes method can be used to convert the file attributes returned in response to NTFSFileAttributes.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file-properties.
func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().GetProperties(ctx, opts, leaseAccessConditions)
	return resp, err
}

// SetHTTPHeaders operation sets HTTP headers on the file.
// ParseNTFSFileAttributes method can be used to convert the file attributes returned in response to NTFSFileAttributes.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties.
func (f *Client) SetHTTPHeaders(ctx context.Context, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, fileHTTPHeaders, leaseAccessConditions := options.format()
	resp, err := f.generated().SetHTTPHeaders(ctx, opts, fileHTTPHeaders, leaseAccessConditions)
	return resp, err
}

// SetMetadata operation sets user-defined metadata for the specified file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-metadata.
func (f *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().SetMetadata(ctx, opts, leaseAccessConditions)
	return resp, err
}

// StartCopyFromURL operation copies the data at the source URL to a file.
//   - copySource: specifies the URL of the source file or blob, up to 2KiB in length.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/copy-file.
func (f *Client) StartCopyFromURL(ctx context.Context, copySource string, options *StartCopyFromURLOptions) (StartCopyFromURLResponse, error) {
	opts, copyFileSmbInfo, leaseAccessConditions := options.format()
	resp, err := f.generated().StartCopy(ctx, copySource, opts, copyFileSmbInfo, leaseAccessConditions)
	return resp, err
}

// AbortCopy operation cancels a pending Copy File operation, and leaves a destination file with zero length and full metadata.
//   - copyID: the copy identifier provided in the x-ms-copy-id header of the original Copy File operation.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/abort-copy-file.
func (f *Client) AbortCopy(ctx context.Context, copyID string, options *AbortCopyOptions) (AbortCopyResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().AbortCopy(ctx, copyID, opts, leaseAccessConditions)
	return resp, err
}

// Resize operation resizes the file to the specified size.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties.
func (f *Client) Resize(ctx context.Context, size int64, options *ResizeOptions) (ResizeResponse, error) {
	opts, leaseAccessConditions := options.format(size)
	resp, err := f.generated().SetHTTPHeaders(ctx, opts, nil, leaseAccessConditions)
	return resp, err
}

// UploadRange operation uploads a range of bytes to a file.
//   - offset: Specifies the start byte at which the range of bytes is to be written.
//   - body: Specifies the data to be uploaded.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f *Client) UploadRange(ctx context.Context, offset int64, body io.ReadSeekCloser, options *UploadRangeOptions) (UploadRangeResponse, error) {
	rangeParam, contentLength, uploadRangeOptions, leaseAccessConditions, err := options.format(offset, body)
	if err != nil {
		return UploadRangeResponse{}, err
	}

	resp, err := f.generated().UploadRange(ctx, rangeParam, RangeWriteTypeUpdate, contentLength, body, uploadRangeOptions, leaseAccessConditions)
	return resp, err
}

// ClearRange operation clears the specified range and releases the space used in storage for that range.
//   - contentRange: Specifies the range of bytes to be cleared.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f *Client) ClearRange(ctx context.Context, contentRange HTTPRange, options *ClearRangeOptions) (ClearRangeResponse, error) {
	rangeParam, leaseAccessConditions, err := options.format(contentRange)
	if err != nil {
		return ClearRangeResponse{}, err
	}

	resp, err := f.generated().UploadRange(ctx, rangeParam, RangeWriteTypeClear, 0, nil, nil, leaseAccessConditions)
	return resp, err
}

// UploadRangeFromURL operation uploads a range of bytes to a file where the contents are read from a URL.
//   - copySource: Specifies the URL of the source file or blob, up to 2 KB in length.
//   - destinationRange: Specifies the range of bytes in the file to be written.
//   - sourceRange: Bytes of source data in the specified range.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/put-range-from-url.
func (f *Client) UploadRangeFromURL(ctx context.Context, copySource string, sourceOffset int64, destinationOffset int64, count int64, options *UploadRangeFromURLOptions) (UploadRangeFromURLResponse, error) {
	destRange, opts, sourceModifiedAccessConditions, leaseAccessConditions, err := options.format(sourceOffset, destinationOffset, count)
	if err != nil {
		return UploadRangeFromURLResponse{}, err
	}

	resp, err := f.generated().UploadRangeFromURL(ctx, destRange, copySource, 0, opts, sourceModifiedAccessConditions, leaseAccessConditions)
	return resp, err
}

// GetRangeList operation returns the list of valid ranges for a file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-ranges.
func (f *Client) GetRangeList(ctx context.Context, options *GetRangeListOptions) (GetRangeListResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().GetRangeList(ctx, opts, leaseAccessConditions)
	return resp, err
}

// ForceCloseHandles operation closes a handle or handles opened on a file.
//   - handleID - Specifies the handle ID to be closed. Use an asterisk (*) as a wildcard string to specify all handles.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/force-close-handles.
func (f *Client) ForceCloseHandles(ctx context.Context, handleID string, options *ForceCloseHandlesOptions) (ForceCloseHandlesResponse, error) {
	opts := options.format()
	resp, err := f.generated().ForceCloseHandles(ctx, handleID, opts)
	return resp, err
}

// ListHandles operation returns a list of open handles on a file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-handles.
func (f *Client) ListHandles(ctx context.Context, options *ListHandlesOptions) (ListHandlesResponse, error) {
	opts := options.format()
	resp, err := f.generated().ListHandles(ctx, opts)
	return resp, err
}

// CreateHardLink operation returns a list of open handles on a file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-handles.
func (f *Client) CreateHardLink(ctx context.Context, options *CreateHardLinkOptions) (CreateHardLinkResponse, error) {
	if options == nil {
		return CreateHardLinkResponse{}, errors.New("targetFile cannot be nil")
	}
	opts, targetFile, leaseAccessConditions := options.format()
	resp, err := f.generated().CreateHardLink(ctx, targetFile, opts, leaseAccessConditions)
	return resp, err
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at file.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (f *Client) GetSASURL(permissions sas.FilePermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if f.sharedKey() == nil {
		return "", fileerror.MissingSharedKeyCredential
	}
	st := o.format()

	urlParts, err := ParseURL(f.URL())
	if err != nil {
		return "", err
	}

	qps, err := sas.SignatureValues{
		Version:     sas.Version,
		ShareName:   urlParts.ShareName,
		FilePath:    urlParts.DirectoryOrFilePath,
		Permissions: permissions.String(),
		StartTime:   st,
		ExpiryTime:  expiry.UTC(),
	}.SignWithSharedKey(f.sharedKey())
	if err != nil {
		return "", err
	}

	endpoint := f.URL() + "?" + qps.Encode()

	return endpoint, nil
}

// Concurrent Upload Functions -----------------------------------------------------------------------------------------

// uploadFromReader uploads a buffer in chunks to an Azure file.
func (f *Client) uploadFromReader(ctx context.Context, reader io.ReaderAt, actualSize int64, o *uploadFromReaderOptions) error {
	if actualSize > MaxFileSize {
		return errors.New("buffer is too large to upload to a file")
	}
	if o.ChunkSize == 0 {
		o.ChunkSize = MaxUpdateRangeBytes
	}

	if log.Should(exported.EventUpload) {
		urlParts, err := ParseURL(f.URL())
		if err == nil {
			log.Writef(exported.EventUpload, "file name %s actual size %v chunk-size %v chunk-count %v",
				urlParts.DirectoryOrFilePath, actualSize, o.ChunkSize, ((actualSize-1)/o.ChunkSize)+1)
		}
	}

	progress := int64(0)
	progressLock := &sync.Mutex{}

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "uploadFromReader",
		TransferSize:  actualSize,
		ChunkSize:     o.ChunkSize,
		Concurrency:   o.Concurrency,
		Operation: func(ctx context.Context, offset int64, chunkSize int64) error {
			// This function is called once per file range.
			// It is passed this file's offset within the buffer and its count of bytes
			// Prepare to read the proper range/section of the buffer
			if chunkSize < o.ChunkSize {
				// this is the last file range.  Its actual size might be less
				// than the calculated size due to rounding up of the payload
				// size to fit in a whole number of chunks.
				chunkSize = actualSize - offset
			}
			var body io.ReadSeeker = io.NewSectionReader(reader, offset, chunkSize)
			if o.Progress != nil {
				chunkProgress := int64(0)
				body = streaming.NewRequestProgress(streaming.NopCloser(body),
					func(bytesTransferred int64) {
						diff := bytesTransferred - chunkProgress
						chunkProgress = bytesTransferred
						progressLock.Lock() // 1 goroutine at a time gets progress report
						progress += diff
						o.Progress(progress)
						progressLock.Unlock()
					})
			}

			uploadRangeOptions := o.getUploadRangeOptions()
			_, err := f.UploadRange(ctx, offset, streaming.NopCloser(body), uploadRangeOptions)
			return err
		},
	})
	return err
}

// UploadBuffer uploads a buffer in chunks to an Azure file.
func (f *Client) UploadBuffer(ctx context.Context, buffer []byte, options *UploadBufferOptions) error {
	uploadOptions := uploadFromReaderOptions{}
	if options != nil {
		uploadOptions = *options
	}
	return f.uploadFromReader(ctx, bytes.NewReader(buffer), int64(len(buffer)), &uploadOptions)
}

// UploadFile uploads a file in chunks to an Azure file.
func (f *Client) UploadFile(ctx context.Context, file *os.File, options *UploadFileOptions) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	uploadOptions := uploadFromReaderOptions{}
	if options != nil {
		uploadOptions = *options
	}
	return f.uploadFromReader(ctx, file, stat.Size(), &uploadOptions)
}

// UploadStream copies the file held in io.Reader to the file at fileClient.
// A Context deadline or cancellation will cause this to error.
func (f *Client) UploadStream(ctx context.Context, body io.Reader, options *UploadStreamOptions) error {
	if options == nil {
		options = &UploadStreamOptions{}
	}

	err := copyFromReader(ctx, body, f, *options, newMMBPool)
	return err
}

// Concurrent Download Functions -----------------------------------------------------------------------------------------

// download method downloads an Azure file to a WriterAt in parallel.
func (f *Client) download(ctx context.Context, writer io.WriterAt, o downloadOptions) (int64, error) {
	if o.ChunkSize == 0 {
		o.ChunkSize = DefaultDownloadChunkSize
	}

	count := o.Range.Count
	if count == CountToEnd { // If size not specified, calculate it
		// If we don't have the length at all, get it
		getFilePropertiesOptions := o.getFilePropertiesOptions()
		gr, err := f.GetProperties(ctx, getFilePropertiesOptions)
		if err != nil {
			return 0, err
		}
		count = *gr.ContentLength - o.Range.Offset
	}

	if count <= 0 {
		// The file is empty, there is nothing to download.
		return 0, nil
	}

	// Prepare and do parallel download.
	progress := int64(0)
	progressLock := &sync.Mutex{}

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "downloadFileToWriterAt",
		TransferSize:  count,
		ChunkSize:     o.ChunkSize,
		Concurrency:   o.Concurrency,
		Operation: func(ctx context.Context, chunkStart int64, count int64) error {
			downloadFileOptions := o.getDownloadFileOptions(HTTPRange{
				Offset: chunkStart + o.Range.Offset,
				Count:  count,
			})
			dr, err := f.DownloadStream(ctx, downloadFileOptions)
			if err != nil {
				return err
			}
			var body io.ReadCloser = dr.NewRetryReader(ctx, &o.RetryReaderOptionsPerChunk)
			if o.Progress != nil {
				rangeProgress := int64(0)
				body = streaming.NewResponseProgress(
					body,
					func(bytesTransferred int64) {
						diff := bytesTransferred - rangeProgress
						rangeProgress = bytesTransferred
						progressLock.Lock()
						progress += diff
						o.Progress(progress)
						progressLock.Unlock()
					})
			}
			_, err = io.Copy(shared.NewSectionWriter(writer, chunkStart, count), body)
			if err != nil {
				return err
			}
			err = body.Close()
			return err
		},
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DownloadStream operation reads or downloads a file from the system, including its metadata and properties.
// ParseNTFSFileAttributes method can be used to convert the file attributes returned in response to NTFSFileAttributes.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file.
func (f *Client) DownloadStream(ctx context.Context, options *DownloadStreamOptions) (DownloadStreamResponse, error) {
	opts, leaseAccessConditions := options.format()
	if options == nil {
		options = &DownloadStreamOptions{}
	}

	resp, err := f.generated().Download(ctx, opts, leaseAccessConditions)
	if err != nil {
		return DownloadStreamResponse{}, err
	}

	return DownloadStreamResponse{
		DownloadResponse:      resp,
		client:                f,
		getInfo:               httpGetterInfo{Range: options.Range},
		leaseAccessConditions: options.LeaseAccessConditions,
	}, err
}

// DownloadBuffer downloads an Azure file to a buffer with parallel.
func (f *Client) DownloadBuffer(ctx context.Context, buffer []byte, o *DownloadBufferOptions) (int64, error) {
	if o == nil {
		o = &DownloadBufferOptions{}
	}

	return f.download(ctx, shared.NewBytesWriter(buffer), (downloadOptions)(*o))
}

// DownloadFile downloads an Azure file to a local file.
// The file would be truncated if the size doesn't match.
func (f *Client) DownloadFile(ctx context.Context, file *os.File, o *DownloadFileOptions) (int64, error) {
	if o == nil {
		o = &DownloadFileOptions{}
	}
	do := (*downloadOptions)(o)

	// 1. Calculate the size of the destination file
	var size int64

	count := do.Range.Count
	if count == CountToEnd {
		// Try to get Azure file's size
		getFilePropertiesOptions := do.getFilePropertiesOptions()
		props, err := f.GetProperties(ctx, getFilePropertiesOptions)
		if err != nil {
			return 0, err
		}
		size = *props.ContentLength - do.Range.Offset
		do.Range.Count = size
	} else {
		size = count
	}

	// 2. Compare and try to resize local file's size if it doesn't match Azure file's size.
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	if stat.Size() != size {
		if err = file.Truncate(size); err != nil {
			return 0, err
		}
	}

	if size > 0 {
		return f.download(ctx, file, *do)
	} else { // if the file's size is 0, there is no need in downloading it
		return 0, nil
	}
}

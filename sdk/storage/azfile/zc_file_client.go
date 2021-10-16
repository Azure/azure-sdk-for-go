package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
	"net/http"
	"net/url"
)

type FileClient struct {
	client *fileClient
	u      url.URL
	cred   azcore.Credential
}

func NewFileClient(fileURL string, cred azcore.Credential, options *ClientOptions) (FileClient, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return FileClient{}, err
	}
	con := newConnection(fileURL, cred, options.getConnectionOptions())
	return FileClient{client: &fileClient{con: con}, u: *u, cred: cred}, nil
}

func (f FileClient) URL() string {
	return f.u.String()
}

// Create creates a new file or replaces a file. Note that this method only initializes the file.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/create-file.
// Pass default values for SMB properties (ex: "None" for file attributes).
func (f FileClient) Create(ctx context.Context, options *CreateFileOptions) (FileCreateResponse, error) {
	fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileCreateResponse{}, err
	}
	return f.client.Create(ctx, fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions)
}

// StartCopy copies the data at the source URL to a file.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-file.
func (f FileClient) StartCopy(ctx context.Context, source url.URL, options *StartFileCopyOptions) (FileStartCopyResponse, error) {
	fileStartCopyOptions, copyFileSmbInfo, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileStartCopyResponse{}, err
	}
	return f.client.StartCopy(ctx, source.String(), fileStartCopyOptions, copyFileSmbInfo, leaseAccessConditions)
}

// AbortCopy stops a pending copy that was previously started and leaves a destination file with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-file.
func (f FileClient) AbortCopy(ctx context.Context, copyID string, options *AbortFileCopyOptions) (FileAbortCopyResponse, error) {
	fileAbortCopyOptions, leaseAccessConditions := options.format()
	return f.client.AbortCopy(ctx, copyID, fileAbortCopyOptions, leaseAccessConditions)
}

// Download downloads Count bytes of data from the start Offset.
// The response includes all the file’s properties. However, passing true for rangeGetContentMD5 returns the range’s MD5 in the ContentMD5
// response header/property if the range is <= 4 MB;
// The HTTP request fails with 400 (Bad Request) if the requested range is greater than 4 MB.
// Note: Bothoffset must be >=0, Count must be >= 0.
// If Count is CountToEnd (0), then data is read from specified Offset to the end.
// rangeGetContentMD5 only works with partial data downloading.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-file.
func (f FileClient) Download(ctx context.Context, options *DownloadFileOptions) (RetryableDownloadResponse, error) {
	fileDownloadOptions, leaseAccessConditions := options.format()

	dr, err := f.client.Download(ctx, fileDownloadOptions, leaseAccessConditions)
	if err != nil {
		return RetryableDownloadResponse{}, err
	}

	offset := int64(0)
	count := int64(CountToEnd)

	if options != nil && options.Offset != nil {
		offset = *options.Offset
	}

	if options != nil && options.Count != nil {
		count = *options.Count
	}

	return RetryableDownloadResponse{
		f:    f,
		dr:   dr,
		ctx:  ctx,
		info: HTTPGetterInfo{Offset: offset, Count: count, ETag: dr.ETag}, // TODO: Note conditional header is not currently supported in Azure File.
	}, err
}

// Body constructs a stream to read data from with a resilient reader option.
// A zero-value option means to get a raw stream.
func (dr *RetryableDownloadResponse) Body(o RetryReaderOptions) io.ReadCloser {
	if o.MaxRetryRequests == 0 {
		return dr.Response().Body
	}

	return NewRetryReader(
		dr.ctx,
		dr.Response(),
		dr.info,
		o,
		func(ctx context.Context, info HTTPGetterInfo) (*http.Response, error) {
			resp, err := dr.f.Download(ctx, &DownloadFileOptions{
				Offset:             to.Int64Ptr(info.Offset),
				Count:              to.Int64Ptr(info.Count),
				RangeGetContentMD5: to.BoolPtr(false),
			})
			if err != nil {
				return nil, err
			}
			return resp.Response(), err
		})
}

// Delete immediately removes the file from the storage account.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/delete-file2.
func (f FileClient) Delete(ctx context.Context, options *DeleteFileOptions) (FileDeleteResponse, error) {
	fileDeleteOptions, leaseAccessConditions := options.format()
	return f.client.Delete(ctx, fileDeleteOptions, leaseAccessConditions)
}

// GetProperties returns the file's metadata and properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-file-properties.
func (f FileClient) GetProperties(ctx context.Context, options *GetFilePropertiesOptions) (FileGetPropertiesResponse, error) {
	fileGetPropertiesOptions, leaseAccessConditions := options.format()
	return f.client.GetProperties(ctx, fileGetPropertiesOptions, leaseAccessConditions)
}

// SetHTTPHeaders sets file's system properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-file-properties.
func (f FileClient) SetHTTPHeaders(ctx context.Context, options *SetFileHTTPHeadersOptions) (FileSetHTTPHeadersResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileSetHTTPHeadersResponse{}, err
	}

	return f.client.SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions)
}

// SetMetadata sets a file's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-file-metadata.
func (f FileClient) SetMetadata(ctx context.Context, options *SetFileMetadataOptions) (FileSetMetadataResponse, error) {
	fileSetMetadataOptions, leaseAccessConditions := options.format()
	return f.client.SetMetadata(ctx, fileSetMetadataOptions, leaseAccessConditions)
}

// Resize resizes the file to the specified size.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-file-properties.
func (f FileClient) Resize(ctx context.Context, options *ResizeFileOptions) (FileSetHTTPHeadersResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions := options.format()
	return f.client.SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions)
}

// UploadRange writes bytes to a file.
// Offset indicates the Offset at which to begin writing, in bytes.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f FileClient) UploadRange(ctx context.Context, options *UploadFileRangeOptions) (FileUploadRangeResponse, error) {
	rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileUploadRangeResponse{}, err
	}

	return f.client.UploadRange(ctx, rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions)
}

// UploadRangeFromURL Update range with bytes from a specific URL.
// Offset indicates the Offset at which to begin writing, in bytes.
func (f FileClient) UploadRangeFromURL(ctx context.Context, sourceURL url.URL, sourceOffset int64, destinationOffset int64,
	count int64, options *UploadFileRangeFromURLOptions) (FileUploadRangeFromURLResponse, error) {

	rangeParam, copySource, contentLength, fileUploadRangeFromURLOptions, sourceModifiedAccessConditions,
		leaseAccessConditions := options.format(sourceURL, sourceOffset, destinationOffset, count)

	return f.client.UploadRangeFromURL(ctx, rangeParam, copySource, contentLength, fileUploadRangeFromURLOptions,
		sourceModifiedAccessConditions, leaseAccessConditions)
}

// ClearRange clears the specified range and releases the space used in storage for that range.
// Offset means the start Offset of the range to clear.
// Count means Count of bytes to clean, it cannot be CountToEnd (0), and must be explicitly specified.
// If the range specified is not 512-byte aligned, the operation will write zeros to
// the start or end of the range that is not 512-byte aligned and free the rest of the range inside that is 512-byte aligned.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f FileClient) ClearRange(ctx context.Context, options *ClearFileRangeOptions) (FileUploadRangeResponse, error) {
	rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileUploadRangeResponse{}, err
	}

	return f.client.UploadRange(ctx, rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions)
}

// GetRangeList returns the list of valid ranges for a file.
// Use a Count with value CountToEnd (0) to indicate the left part of file start from Offset.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/list-ranges.
func (f FileClient) GetRangeList(ctx context.Context, options *GetFileRangeListOptions) (FileGetRangeListResponse, error) {
	fileGetRangeListOptions, leaseAccessConditions := options.format()
	return f.client.GetRangeList(ctx, fileGetRangeListOptions, leaseAccessConditions)
}

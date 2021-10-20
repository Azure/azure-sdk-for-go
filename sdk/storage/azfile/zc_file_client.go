package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"io"
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

// WithSnapshot creates a new BlobClient object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (f FileClient) WithSnapshot(shareSnapshot string) FileClient {
	fileURLParts := NewFileURLParts(f.URL())
	fileURLParts.ShareSnapshot = shareSnapshot
	u, _ := url.Parse(fileURLParts.URL())

	return FileClient{
		client: &fileClient{
			&connection{u: fileURLParts.URL(), p: f.client.con.p},
		},
		u:    *u,
		cred: f.cred,
	}
}

// Create creates a new file or replaces a file. Note that this method only initializes the file.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/create-file.
// Pass default values for SMB properties (ex: "None" for file attributes).
func (f FileClient) Create(ctx context.Context, options *CreateFileOptions) (FileCreateResponse, error) {
	fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileCreateResponse{}, err
	}
	fileCreateResponse, err := f.client.Create(ctx, fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions)
	return fileCreateResponse, handleError(err)
}

// StartCopy copies the data at the source URL to a file.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-file.
func (f FileClient) StartCopy(ctx context.Context, sourceURL string, options *StartFileCopyOptions) (FileStartCopyResponse, error) {
	fileStartCopyOptions, copyFileSmbInfo, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileStartCopyResponse{}, err
	}
	fileStartCopyResponse, err := f.client.StartCopy(ctx, sourceURL, fileStartCopyOptions, copyFileSmbInfo, leaseAccessConditions)
	return fileStartCopyResponse, handleError(err)
}

// AbortCopy stops a pending copy that was previously started and leaves a destination file with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-file.
func (f FileClient) AbortCopy(ctx context.Context, copyID string, options *AbortFileCopyOptions) (FileAbortCopyResponse, error) {
	fileAbortCopyOptions, leaseAccessConditions := options.format()
	fileAbortCopyResponse, err := f.client.AbortCopy(ctx, copyID, fileAbortCopyOptions, leaseAccessConditions)
	return fileAbortCopyResponse, handleError(err)
}

// Download downloads Count bytes of data from the start Offset.
// The response includes all the file’s properties. However, passing true for rangeGetContentMD5 returns the range’s MD5 in the ContentMD5
// response header/property if the range is <= 4 MB;
// The HTTP request fails with 400 (Bad Request) if the requested range is greater than 4 MB.
// Note: Both offset and count must be >=0.
// If Count is CountToEnd (0), then data is read from specified Offset to the end.
// RangeGetContentMD5 only works with partial data downloading.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-file.
func (f FileClient) Download(ctx context.Context, offset, count int64, options *DownloadFileOptions) (DownloadResponse, error) {

	fileDownloadOptions, leaseAccessConditions := options.format(to.Int64Ptr(offset), to.Int64Ptr(count))

	dr, err := f.client.Download(ctx, fileDownloadOptions, leaseAccessConditions)
	if err != nil {
		return DownloadResponse{}, handleError(err)
	}

	return DownloadResponse{
		FileDownloadResponse: dr,
		f:                    f,
		ctx:                  ctx,
		info:                 HTTPGetterInfo{Offset: offset, Count: count, ETag: dr.ETag},
	}, handleError(err)
}

// Delete immediately removes the file from the storage account.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/delete-file2.
func (f FileClient) Delete(ctx context.Context, options *DeleteFileOptions) (FileDeleteResponse, error) {
	fileDeleteOptions, leaseAccessConditions := options.format()
	fileDeleteResponse, err := f.client.Delete(ctx, fileDeleteOptions, leaseAccessConditions)
	return fileDeleteResponse, handleError(err)
}

// GetProperties returns the file's metadata and properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-file-properties.
func (f FileClient) GetProperties(ctx context.Context, options *GetFilePropertiesOptions) (FileGetPropertiesResponse, error) {
	fileGetPropertiesOptions, leaseAccessConditions := options.format()
	fileGetPropertiesResponse, err := f.client.GetProperties(ctx, fileGetPropertiesOptions, leaseAccessConditions)
	return fileGetPropertiesResponse, handleError(err)
}

// SetHTTPHeaders sets file's system properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-file-properties.
func (f FileClient) SetHTTPHeaders(ctx context.Context, options *SetFileHTTPHeadersOptions) (FileSetHTTPHeadersResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileSetHTTPHeadersResponse{}, err
	}

	fileSetHTTPHeadersResponse, err := f.client.SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions)
	return fileSetHTTPHeadersResponse, handleError(err)
}

// SetMetadata sets a file's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-file-metadata.
func (f FileClient) SetMetadata(ctx context.Context, metadata map[string]string, options *SetFileMetadataOptions) (FileSetMetadataResponse, error) {
	fileSetMetadataOptions, leaseAccessConditions := options.format(metadata)
	fileSetMetadataResponse, err := f.client.SetMetadata(ctx, fileSetMetadataOptions, leaseAccessConditions)
	return fileSetMetadataResponse, handleError(err)
}

// Resize resizes the file to the specified size.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-file-properties.
func (f FileClient) Resize(ctx context.Context, size int64, options *ResizeFileOptions) (FileSetHTTPHeadersResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions := options.format(size)
	fileSetHTTPHeadersResponse, err := f.client.SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, fileSetHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions)
	return fileSetHTTPHeadersResponse, handleError(err)
}

// UploadRange writes bytes to a file.
// Offset indicates the Offset at which to begin writing, in bytes.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f FileClient) UploadRange(ctx context.Context, offset int64, body io.ReadSeekCloser, options *UploadFileRangeOptions) (FileUploadRangeResponse, error) {
	rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions, err := options.format(offset, body)
	if err != nil {
		return FileUploadRangeResponse{}, err
	}

	fileUploadRangeResponse, err := f.client.UploadRange(ctx, rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions)
	return fileUploadRangeResponse, handleError(err)
}

// UploadRangeFromURL Update range with bytes from a specific URL.
// Offset indicates the Offset at which to begin writing, in bytes.
func (f FileClient) UploadRangeFromURL(ctx context.Context, sourceURL string, sourceOffset int64, destinationOffset int64,
	count int64, options *UploadFileRangeFromURLOptions) (FileUploadRangeFromURLResponse, error) {

	rangeParam, copySource, contentLength, fileUploadRangeFromURLOptions, sourceModifiedAccessConditions,
		leaseAccessConditions := options.format(sourceURL, to.Int64Ptr(sourceOffset), to.Int64Ptr(destinationOffset), to.Int64Ptr(count))

	fileUploadRangeFromURLResponse, err := f.client.UploadRangeFromURL(ctx, rangeParam, copySource, contentLength, fileUploadRangeFromURLOptions,
		sourceModifiedAccessConditions, leaseAccessConditions)
	return fileUploadRangeFromURLResponse, handleError(err)
}

// ClearRange clears the specified range and releases the space used in storage for that range.
// Offset means the start Offset of the range to clear.
// Count means Count of bytes to clean, it cannot be CountToEnd (0), and must be explicitly specified.
// If the range specified is not 512-byte aligned, the operation will write zeros to
// the start or end of the range that is not 512-byte aligned and free the rest of the range inside that is 512-byte aligned.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f FileClient) ClearRange(ctx context.Context, offset, count int64, options *ClearFileRangeOptions) (FileUploadRangeResponse, error) {
	rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions, err := options.format(to.Int64Ptr(offset), to.Int64Ptr(count))
	if err != nil {
		return FileUploadRangeResponse{}, err
	}

	fileUploadRangeResponse, err := f.client.UploadRange(ctx, rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions)
	return fileUploadRangeResponse, handleError(err)
}

// GetRangeList returns the list of valid ranges for a file.
// Use a Count with value CountToEnd (0) to indicate the left part of file start from Offset.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/list-ranges.
func (f FileClient) GetRangeList(ctx context.Context, offset, count int64, options *GetFileRangeListOptions) (FileGetRangeListResponse, error) {
	fileGetRangeListOptions, leaseAccessConditions := options.format(to.Int64Ptr(offset), to.Int64Ptr(count))
	fileGetRangeListResponse, err := f.client.GetRangeList(ctx, fileGetRangeListOptions, leaseAccessConditions)
	return fileGetRangeListResponse, handleError(err)
}

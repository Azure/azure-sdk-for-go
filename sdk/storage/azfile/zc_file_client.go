//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"io"
)

type FileClient struct {
	client    *fileClient
	sharedKey *SharedKeyCredential
}

// NewFileClient creates a DirectoryClient object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewFileClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*FileClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(serviceURL, conOptions)

	return &FileClient{
		client: newFileClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewFileClientWithNoCredential creates a DirectoryClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewFileClientWithNoCredential(serviceURL string, options *ClientOptions) (*FileClient, error) {
	conOptions := getConnectionOptions(options)
	conn := newConnection(serviceURL, conOptions)

	return &FileClient{
		client: newFileClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewFileClientWithSharedKey creates a DirectoryClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewFileClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*FileClient, error) {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(serviceURL, conOptions)

	return &FileClient{
		client:    newFileClient(conn.Endpoint(), conn.Pipeline()),
		sharedKey: cred,
	}, nil
}

// NewFileClientFromConnectionString creates a DirectoryClient from the given connection string.
//nolint
func NewFileClientFromConnectionString(connectionString string, options *ClientOptions) (*FileClient, error) {
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewFileClientWithSharedKey(endpoint, credential, options)
}

func (f *FileClient) URL() string {
	return f.client.endpoint
}

// WithSnapshot creates a new BlobClient object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (f *FileClient) WithSnapshot(shareSnapshot string) (*FileClient, error) {
	p, err := NewFileURLParts(f.URL())
	if err != nil {
		return nil, err
	}

	p.ShareSnapshot = shareSnapshot
	endpoint := p.URL()
	fClient := newFileClient(endpoint, f.client.pl)

	return &FileClient{
		client:    fClient,
		sharedKey: f.sharedKey,
	}, nil
}

// Create creates a new file or replaces a file. Note that this method only initializes the file.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/create-file.
// Pass default values for SMB properties (ex: "None" for file attributes).
func (f *FileClient) Create(ctx context.Context, options *FileCreateOptions) (FileCreateResponse, error) {
	fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileCreateResponse{}, err
	}

	fileCreateResponse, err := f.client.Create(ctx, fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime,
		fileCreateOptions, fileHTTPHeaders, leaseAccessConditions)

	return toFileCreateResponse(fileCreateResponse), handleError(err)
}

// StartCopy copies the data at the source URL to a file.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-file.
func (f *FileClient) StartCopy(ctx context.Context, sourceURL string, options *FileStartCopyOptions) (FileStartCopyResponse, error) {
	fileStartCopyOptions, copyFileSmbInfo, leaseAccessConditions, err := options.format()
	if err != nil {
		return FileStartCopyResponse{}, err
	}
	fileStartCopyResponse, err := f.client.StartCopy(ctx, sourceURL, fileStartCopyOptions, copyFileSmbInfo, leaseAccessConditions)
	return toFileStartCopyResponse(fileStartCopyResponse), handleError(err)
}

// AbortCopy stops a pending copy that was previously started and leaves a destination file with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-file.
func (f *FileClient) AbortCopy(ctx context.Context, copyID string, options *FileAbortCopyOptions) (FileAbortCopyResponse, error) {
	fileAbortCopyOptions, leaseAccessConditions := options.format()
	fileAbortCopyResponse, err := f.client.AbortCopy(ctx, copyID, fileAbortCopyOptions, leaseAccessConditions)
	return toFileAbortCopyResponse(fileAbortCopyResponse), handleError(err)
}

// Download downloads Count bytes of data from the start Offset.
// The responseBody includes all the file’s properties. However, passing true for rangeGetContentMD5 returns the range’s MD5 in the ContentMD5
// responseBody header/property if the range is <= 4 MB;
// The HTTP request fails with 400 (Bad Request) if the requested range is greater than 4 MB.
// Note: Both offset and count must be >=0.
// If Count is CountToEnd (0), then data is read from specified Offset to the end.
// RangeGetContentMD5 only works with partial data downloading.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-file.
func (f *FileClient) Download(ctx context.Context, offset, count int64, options *FileDownloadOptions) (FileDownloadResponse, error) {

	fileDownloadOptions, leaseAccessConditions := options.format(offset, count)

	fileDownloadResponse, err := f.client.Download(ctx, fileDownloadOptions, leaseAccessConditions)
	if err != nil {
		return FileDownloadResponse{}, handleError(err)
	}

	return toFileDownloadResponse(ctx, f, fileDownloadResponse, offset, count), handleError(err)
}

// Delete immediately removes the file from the storage account.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/delete-file2.
func (f *FileClient) Delete(ctx context.Context, options *FileDeleteOptions) (FileDeleteResponse, error) {
	fileDeleteOptions, leaseAccessConditions := options.format()
	fileDeleteResponse, err := f.client.Delete(ctx, fileDeleteOptions, leaseAccessConditions)
	return toFileDeleteResponse(fileDeleteResponse), handleError(err)
}

// GetProperties returns the file's metadata and properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-file-properties.
func (f *FileClient) GetProperties(ctx context.Context, options *FileGetPropertiesOptions) (FileGetPropertiesResponse, error) {
	fileGetPropertiesOptions, leaseAccessConditions := options.format()

	fileGetPropertiesResponse, err := f.client.GetProperties(ctx, fileGetPropertiesOptions, leaseAccessConditions)

	return toFileGetPropertiesResponse(fileGetPropertiesResponse), handleError(err)
}

// SetHTTPHeaders sets file's system properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-file-properties.
func (f *FileClient) SetHTTPHeaders(ctx context.Context, options *FileSetHTTPHeadersOptions) (FileSetHTTPHeadersResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, setHTTPHeadersOptions, shareFileHTTPHeaders, leaseAccessConditions, err := options.format()

	if err != nil {
		return FileSetHTTPHeadersResponse{}, err
	}

	fileSetHTTPHeadersResponse, err := f.client.SetHTTPHeaders(ctx, fileAttributes, fileCreationTime,
		fileLastWriteTime, setHTTPHeadersOptions, shareFileHTTPHeaders, leaseAccessConditions)

	return toFileSetHTTPHeadersResponse(fileSetHTTPHeadersResponse), handleError(err)
}

// SetMetadata sets a file's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-file-metadata.
func (f *FileClient) SetMetadata(ctx context.Context, metadata map[string]string, options *FileSetMetadataOptions) (FileSetMetadataResponse, error) {
	fileSetMetadataOptions, leaseAccessConditions := options.format(metadata)

	fileSetMetadataResponse, err := f.client.SetMetadata(ctx, fileSetMetadataOptions, leaseAccessConditions)

	return toFileSetMetadataResponse(fileSetMetadataResponse), handleError(err)
}

// Resize resizes the file to the specified size.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-file-properties.
func (f *FileClient) Resize(ctx context.Context, size int64, options *FileResizeOptions) (FileResizeResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, setHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions := options.format(size)

	setHTTPHeadersResponse, err := f.client.SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime,
		setHTTPHeadersOptions, fileHTTPHeaders, leaseAccessConditions)

	return toFileResizeResponse(setHTTPHeadersResponse), handleError(err)
}

// UploadRange writes bytes to a file.
// Offset indicates the Offset at which to begin writing, in bytes.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f *FileClient) UploadRange(ctx context.Context, offset int64, body io.ReadSeekCloser, options *FileUploadRangeOptions) (FileUploadRangeResponse, error) {
	rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions, err := options.format(offset, body)
	if err != nil {
		return FileUploadRangeResponse{}, err
	}

	uploadRangeResponse, err := f.client.UploadRange(ctx, rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions)
	return toFileUploadRangeResponse(uploadRangeResponse), handleError(err)
}

// UploadRangeFromURL Update range with bytes from a specific URL.
// Offset indicates the Offset at which to begin writing, in bytes.
func (f *FileClient) UploadRangeFromURL(ctx context.Context, sourceURL string, sourceOffset int64, destinationOffset int64,
	count int64, options *FileUploadRangeFromURLOptions) (FileUploadRangeFromURLResponse, error) {

	rangeParam, copySource, contentLength, uploadRangeFromURLOptions,
		sourceModifiedAccessConditions, leaseAccessConditions := options.format(sourceURL, sourceOffset, destinationOffset, count)

	uploadRangeFromURLResponse, err := f.client.UploadRangeFromURL(ctx, rangeParam, copySource, contentLength, uploadRangeFromURLOptions, sourceModifiedAccessConditions, leaseAccessConditions)
	return toFileUploadRangeFromURLResponse(uploadRangeFromURLResponse), handleError(err)
}

// ClearRange clears the specified range and releases the space used in storage for that range.
// Offset means the start Offset of the range to clear.
// Count means Count of bytes to clean, it cannot be CountToEnd (0), and must be explicitly specified.
// If the range specified is not 512-byte aligned, the operation will write zeros to
// the start or end of the range that is not 512-byte aligned and free the rest of the range inside that is 512-byte aligned.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-range.
func (f *FileClient) ClearRange(ctx context.Context, offset, count int64, options *FileClearRangeOptions) (FileClearRangeResponse, error) {
	rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions, err := options.format(offset, count)
	if err != nil {
		return FileClearRangeResponse{}, err
	}

	fileUploadRangeResponse, err := f.client.UploadRange(ctx, rangeParam, fileRangeWrite, contentLength, fileUploadRangeOptions, leaseAccessConditions)

	return toFileClearRangeResponse(fileUploadRangeResponse), handleError(err)
}

// GetRangeList returns the list of valid ranges for a file.
// Use a Count with value CountToEnd (0) to indicate the left part of file start from Offset.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/list-ranges.
func (f *FileClient) GetRangeList(ctx context.Context, offset, count int64, options *FileGetRangeListOptions) (FileGetRangeListResponse, error) {
	fileGetRangeListOptions, leaseAccessConditions := options.format(offset, count)
	getRangeListResponse, err := f.client.GetRangeList(ctx, fileGetRangeListOptions, leaseAccessConditions)

	return toFileGetRangeListResponse(getRangeListResponse), handleError(err)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"io"
	"os"
	"strings"
	"time"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure Storage file.
type Client base.Client[generated.FileClient]

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a file or with a shared access signature (SAS) token.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file?<sas token>
//   - options - client options; pass nil to accept the default values
//
// The directoryPath is optional in the fileURL. If omitted, it points to file within the specified share.
func NewClientWithNoCredential(fileURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewFileClient(fileURL, pl, nil)), nil
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
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewFileClient(fileURL, pl, cred)), nil
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

// URL returns the URL endpoint used by the Client object.
func (f *Client) URL() string {
	return f.generated().Endpoint()
}

// Create operation creates a new file or replaces a file. Note it only initializes the file with no content.
//   - fileContentLength: Specifies the maximum size for the file in bytes, up to 4 TB.
//
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-file.
func (f *Client) Create(ctx context.Context, fileContentLength int64, options *CreateOptions) (CreateResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions := options.format()
	resp, err := f.generated().Create(ctx, fileContentLength, fileAttributes, fileCreationTime, fileLastWriteTime, fileCreateOptions, fileHTTPHeaders, leaseAccessConditions)
	return resp, err
}

// Delete operation removes the file from the storage account.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-file2.
func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().Delete(ctx, opts, leaseAccessConditions)
	return resp, err
}

// GetProperties operation returns all user-defined metadata, standard HTTP properties, and system properties for the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file-properties.
func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := f.generated().GetProperties(ctx, opts, leaseAccessConditions)
	return resp, err
}

// SetHTTPHeaders operation sets HTTP headers on the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties.
func (f *Client) SetHTTPHeaders(ctx context.Context, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, opts, fileHTTPHeaders, leaseAccessConditions := options.format()
	resp, err := f.generated().SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, opts, fileHTTPHeaders, leaseAccessConditions)
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

// DownloadStream operation reads or downloads a file from the system, including its metadata and properties.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file.
func (f *Client) DownloadStream(ctx context.Context, options *DownloadStreamOptions) (DownloadStreamResponse, error) {
	return DownloadStreamResponse{}, nil
}

// DownloadBuffer downloads an Azure file to a buffer with parallel.
func (f *Client) DownloadBuffer(ctx context.Context, buffer []byte, o *DownloadBufferOptions) (int64, error) {
	return 0, nil
}

// DownloadFile downloads an Azure file to a local file.
// The file would be truncated if the size doesn't match.
func (f *Client) DownloadFile(ctx context.Context, file *os.File, o *DownloadFileOptions) (int64, error) {
	return 0, nil
}

// Resize operation resizes the file to the specified size.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties.
func (f *Client) Resize(ctx context.Context, size int64, options *ResizeOptions) (ResizeResponse, error) {
	fileAttributes, fileCreationTime, fileLastWriteTime, opts, leaseAccessConditions := options.format(size)
	resp, err := f.generated().SetHTTPHeaders(ctx, fileAttributes, fileCreationTime, fileLastWriteTime, opts, nil, leaseAccessConditions)
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

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at file.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (f *Client) GetSASURL(permissions sas.FilePermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if f.sharedKey() == nil {
		return "", fileerror.MissingSharedKeyCredential
	}
	st := o.format()

	urlParts, err := sas.ParseURL(f.URL())
	if err != nil {
		return "", err
	}

	qps, err := sas.SignatureValues{
		Version:             sas.Version,
		Protocol:            sas.ProtocolHTTPS,
		ShareName:           urlParts.ShareName,
		DirectoryOrFilePath: urlParts.DirectoryOrFilePath,
		Permissions:         permissions.String(),
		StartTime:           st,
		ExpiryTime:          expiry.UTC(),
	}.SignWithSharedKey(f.sharedKey())
	if err != nil {
		return "", err
	}

	endpoint := f.URL() + "?" + qps.Encode()

	return endpoint, nil
}

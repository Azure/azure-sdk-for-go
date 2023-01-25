//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure Storage file.
type Client base.Client[generated.FileClient]

// NewClient creates an instance of Client with the specified values.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(fileURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a file or with a shared access signature (SAS) token.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(fileURL string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - fileURL - the URL of the file e.g. https://<account>.file.core.windows.net/share/directoryPath/file
//   - cred - a SharedKeyCredential created with the matching file's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(fileURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - shareName - the name of the share within the storage account
//	 - directoryName - the name of the directory within the storage account
//	 - fileName - the name of the file within the storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, shareName string, directoryName string, fileName string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

func (f *Client) generated() *generated.FileClient {
	return base.InnerClient((*base.Client[generated.FileClient])(f))
}

func (f *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.FileClient])(f))
}

// URL returns the URL endpoint used by the Client object.
func (f *Client) URL() string {
	return "s.generated().Endpoint()"
}

// Create operation creates a new file or replaces a file. Note it only initializes the file with no content.
// 	 - fileContentLength - Specifies the maximum size for the file, up to 4 TB.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-file.
func (f *Client) Create(ctx context.Context, fileContentLength int64, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// Delete operation removes the file from the storage account.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-file2.
func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

// GetProperties operation returns all user-defined metadata, standard HTTP properties, and system properties for the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file-properties.
func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil
}

// SetProperties operation sets system properties on the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-properties.
func (f *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	return SetPropertiesResponse{}, nil
}

// SetMetadata operation sets user-defined metadata for the specified file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-metadata.
func (f *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	return SetMetadataResponse{}, nil
}

// StartCopyFromURL operation copies the data at the source URL to a file.
// 	 - copySource - specifies the URL of the source file or blob, up to 2KiB in length.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/copy-file.
func (f *Client) StartCopyFromURL(ctx context.Context, copySource string, options *StartCopyFromURLOptions) (StartCopyFromURLResponse, error) {
	return StartCopyFromURLResponse{}, nil
}

// AbortCopy operation cancels a pending Copy File operation, and leaves a destination file with zero length and full metadata.
// 	 - copyID - the copy identifier provided in the x-ms-copy-id header of the original Copy File operation.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/abort-copy-file.
func (f *Client) AbortCopy(ctx context.Context, copyID string, options *AbortCopyOptions) (AbortCopyResponse, error) {
	return AbortCopyResponse{}, nil
}

// Download operation reads or downloads a file from the system, including its metadata and properties.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file.
func (f *Client) Download(ctx context.Context, options *DownloadOptions) (DownloadResponse, error) {
	return DownloadResponse{}, nil
}

// AcquireLease operation can be used to request a new lease.
// The lease duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *Client) AcquireLease(ctx context.Context, duration int32, options *AcquireLeaseOptions) (AcquireLeaseResponse, error) {
	// TODO: update generated code to make duration as required parameter
	return AcquireLeaseResponse{}, nil
}

// BreakLease operation can be used to break the lease, if the file has an active lease. Once a lease is broken, it cannot be renewed.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *Client) BreakLease(ctx context.Context, options *BreakLeaseOptions) (BreakLeaseResponse, error) {
	return BreakLeaseResponse{}, nil
}

// ChangeLease operation can be used to change the lease ID of an active lease.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *Client) ChangeLease(ctx context.Context, leaseID string, proposedLeaseID string, options *ChangeLeaseOptions) (ChangeLeaseResponse, error) {
	// TODO: update generated code to make proposedLeaseID as required parameter
	return ChangeLeaseResponse{}, nil
}

// ReleaseLease operation can be used to free the lease if it is no longer needed so that another client may immediately acquire a lease against the file.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/lease-file.
func (f *Client) ReleaseLease(ctx context.Context, leaseID string, options *ReleaseLeaseOptions) (ReleaseLeaseResponse, error) {
	return ReleaseLeaseResponse{}, nil
}

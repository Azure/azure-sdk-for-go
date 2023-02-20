//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure Storage directory allowing you to manipulate its directories and files.
type Client base.Client[generated.DirectoryClient]

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a directory or with a shared access signature (SAS) token.
//   - directoryURL - the URL of the directory e.g. https://<account>.file.core.windows.net/share/directory?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(directoryURL string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - directoryURL - the URL of the directory e.g. https://<account>.file.core.windows.net/share/directory
//   - cred - a SharedKeyCredential created with the matching directory's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(directoryURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - shareName - the name of the share within the storage account
//   - directoryName - the name of the directory within the storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, shareName string, directoryName string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

func (d *Client) generated() *generated.DirectoryClient {
	return base.InnerClient((*base.Client[generated.DirectoryClient])(d))
}

func (d *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.DirectoryClient])(d))
}

// URL returns the URL endpoint used by the Client object.
func (d *Client) URL() string {
	return "s.generated().Endpoint()"
}

// NewSubdirectoryClient creates a new Client object by concatenating subDirectoryName to the end of this Client's URL.
// The new subdirectory Client uses the same request policy pipeline as the parent directory Client.
func (d *Client) NewSubdirectoryClient(subDirectoryName string) *Client {
	return nil
}

// NewFileClient creates a new file.Client object by concatenating fileName to the end of this Client's URL.
// The new file.Client uses the same request policy pipeline as the Client.
func (d *Client) NewFileClient(fileName string) *file.Client {
	return nil
}

// Create operation creates a new directory under the specified share or parent directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-directory.
func (d *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// Delete operation removes the specified empty directory. Note that the directory must be empty before it can be deleted.
// Deleting directories that aren't empty returns error 409 (Directory Not Empty).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-directory.
func (d *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

// GetProperties operation returns all system properties for the specified directory, and it can also be used to check the existence of a directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-directory-properties.
func (d *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil
}

// SetProperties operation sets system properties for the specified directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-directory-properties.
func (d *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	return SetPropertiesResponse{}, nil
}

// SetMetadata operation sets user-defined metadata for the specified directory.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-directory-metadata.
func (d *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	return SetMetadataResponse{}, nil
}

// NewListFilesAndDirectoriesPager operation returns a pager for the files and directories starting from the specified Marker.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-directories-and-files.
func (d *Client) NewListFilesAndDirectoriesPager(options *ListFilesAndDirectoriesOptions) *runtime.Pager[ListFilesAndDirectoriesResponse] {
	return nil
}

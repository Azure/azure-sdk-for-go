//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
)

// Client represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
type Client struct {
	path.Client
}
type ClientOptions = path.ClientOptions

// NewClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *exported.SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// Create creates a new file (dfs1).
func (f *Client) Create(options *CreateOptions) (CreateResponse, error) {
	// TODO: format for options should be able to handle the access conditions parameter correctly
	return CreateResponse{}, nil
}

// Delete deletes a file (dfs1).
func (f *Client) Delete(options *DeleteOptions) (DeleteResponse, error) {
	// TODO: recursive set to false when calling generated code
	return DeleteResponse{}, nil
}

// GetProperties gets the properties of a file (blob3)
func (f *Client) GetProperties(options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil
}

// Rename renames a file (dfs1).
func (f *Client) Rename(newName string, options *RenameOptions) error {
	return nil
}

// SetExpiry sets the expiry time of a file (dfs/blob2).
func (f *Client) SetExpiry(expiryOptions ExpiryOptions, options *SetExpiryOptions) (SetExpiryResponse, error) {
	return SetExpiryResponse{}, nil
}

// Upload uploads data to a file.
func (f *Client) Upload() {

}

// Append appends data to a file.
func (f *Client) Append() {

}

// Flush flushes previous uploaded data to a file.
func (f *Client) Flush() {

}

// Download downloads data from a file.
func (f *Client) Download() {

}

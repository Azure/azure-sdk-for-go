//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
)

// FileClient represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
type FileClient struct {
	Client
}

// NewFileClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewFileClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewFileClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewFileClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewFileClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewFileClientWithSharedKeyCredential(serviceURL string, cred *exported.SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewFileClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewFileClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// dfs
func (f *FileClient) Create(options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// dfs
func (f *FileClient) Delete(options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

// blob
func (f *FileClient) GetProperties(options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil

}

func (f *FileClient) SetExpiry() {

}

func (f *FileClient) RenameFile() {

}

func (f *FileClient) UploadData() {

}

func (f *FileClient) AppendData() {

}

func (f *FileClient) FlushData() {

}

func (f *FileClient) Download() {

}

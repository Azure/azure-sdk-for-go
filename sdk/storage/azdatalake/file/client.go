//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// FileClient represents a URL to the Azure Datalake Storage service.
type FileClient base.Client[generated.PathClient]

type Client struct {
	*FileClient
	blobClient                 *blob.Client
	fileClientWithBlobEndpoint *FileClient
}

//TODO: NewClient()

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.ServiceClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	fileClient := base.NewPathClient(serviceURL, azClient, nil, (*base.ClientOptions)(conOptions))
	blobClientOpts := blob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blob.NewClientWithNoCredential(serviceURL, &blobClientOpts)

	return &Client{
		FileClient: (*FileClient)(fileClient),
		blobClient: blobClient,
	}, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.ServiceClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	fileClient := base.NewPathClient(serviceURL, azClient, cred, (*base.ClientOptions)(conOptions))
	blobClientOpts := blob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKeyCredential, _ := blob.NewSharedKeyCredential(cred.AccountName(), cred.AccountKey())
	blobClient, _ := blob.NewClientWithSharedKeyCredential(serviceURL, blobSharedKeyCredential, &blobClientOpts)

	return &Client{
		FileClient: (*FileClient)(fileClient),
		blobClient: blobClient,
	}, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

// Create creates a new file (dfs1).
func (f *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	// TODO: format for options should be able to handle the access conditions parameter correctly
	return CreateResponse{}, nil
}

// Delete deletes a file (dfs1).
func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	// TODO: recursive set to false when calling generated code
	return DeleteResponse{}, nil
}

// GetProperties gets the properties of a file (blob3)
func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	// TODO: format blob response to path response
	return GetPropertiesResponse{}, nil
}

// Rename renames a file (dfs1).
func (f *Client) Rename(ctx context.Context, newName string, options *RenameOptions) (RenameResponse, error) {
	return RenameResponse{}, nil
}

// SetExpiry operation sets an expiry time on an existing file.
func (f *Client) SetExpiry(ctx context.Context, expiryType ExpiryType, o *SetExpiryOptions) (SetExpiryResponse, error) {
	// TODO: consider using the blob client set expiry
	// TODO: call methods in set_expiry.go
	return SetExpiryResponse{}, nil
}

// Upload uploads data to a file.
func (f *Client) Upload(ctx context.Context) {

}

// Append appends data to a file.
func (f *Client) Append(ctx context.Context) {

}

// Flush flushes previous uploaded data to a file.
func (f *Client) Flush(ctx context.Context) {

}

// Download downloads data from a file.
func (f *Client) Download(ctx context.Context) {

}

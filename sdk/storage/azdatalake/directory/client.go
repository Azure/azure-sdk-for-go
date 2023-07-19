//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client]

// NewClient creates an instance of Client with the specified values.
//   - directoryURL - the URL of the directory e.g. https://<account>.dfs.core.windows.net/fs/dir
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(directoryURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	blobURL, directoryURL := shared.GetURLs(directoryURL)

	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FileClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blockblob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blockblob.NewClient(blobURL, cred, &blobClientOpts)
	dirClient := base.NewPathClient(directoryURL, blobURL, blobClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(dirClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - directoryURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/dir?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(directoryURL string, options *ClientOptions) (*Client, error) {
	blobURL, directoryURL := shared.GetURLs(directoryURL)

	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.DirectoryClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blockblob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blockblob.NewClientWithNoCredential(blobURL, &blobClientOpts)
	dirClient := base.NewPathClient(directoryURL, blobURL, blobClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(dirClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - directoryURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/dir
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(directoryURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	blobURL, directoryURL := shared.GetURLs(directoryURL)

	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.DirectoryClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blockblob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKey, err := cred.ConvertToBlobSharedKey()
	if err != nil {
		return nil, err
	}
	blobClient, _ := blockblob.NewClientWithSharedKeyCredential(blobURL, blobSharedKey, &blobClientOpts)
	dirClient := base.NewPathClient(directoryURL, blobURL, blobClient, azClient, cred, (*base.ClientOptions)(conOptions))

	return (*Client)(dirClient), nil
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

func (d *Client) generatedDirClientWithDFS() *generated.PathClient {
	//base.SharedKeyComposite((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	dirClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
	return dirClientWithDFS
}

func (d *Client) generatedDirClientWithBlob() *generated.PathClient {
	_, dirClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
	return dirClientWithBlob
}

func (d *Client) blobClient() *blockblob.Client {
	_, _, blobClient := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
	return blobClient
}

func (d *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.PathClient, generated.PathClient, blockblob.Client])(d))
}

// DFSURL returns the URL endpoint used by the Client object.
func (d *Client) DFSURL() string {
	return d.generatedDirClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (d *Client) BlobURL() string {
	return d.generatedDirClientWithBlob().Endpoint()
}

// Create creates a new directory (dfs1).
func (d *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// Delete removes the directory (dfs1).
func (d *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	//TODO: pass recursive = true
	return DeleteResponse{}, nil
}

// GetProperties returns the properties of the directory (blob3). #TODO: we may just move this method to path client
func (d *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	// TODO: format blob response to path response
	return GetPropertiesResponse{}, nil
}

// Rename renames the directory (dfs1).
func (d *Client) Rename(ctx context.Context, newName string, options *RenameOptions) (RenameResponse, error) {
	return RenameResponse{}, nil
}

// SetAccessControl sets the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) SetAccessControl(ctx context.Context, options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	return SetAccessControlResponse{}, nil
}

// SetAccessControlRecursive sets the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) SetAccessControlRecursive(ctx context.Context, options *SetAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	// TODO explicitly pass SetAccessControlRecursiveMode
	return SetAccessControlRecursiveResponse{}, nil
}

// UpdateAccessControlRecursive updates the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) UpdateAccessControlRecursive(ctx context.Context, options *UpdateAccessControlRecursiveOptions) (UpdateAccessControlRecursiveResponse, error) {
	// TODO explicitly pass SetAccessControlRecursiveMode
	return SetAccessControlRecursiveResponse{}, nil
}

// GetAccessControl gets the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) GetAccessControl(ctx context.Context, options *GetAccessControlOptions) (GetAccessControlResponse, error) {
	return GetAccessControlResponse{}, nil
}

// RemoveAccessControlRecursive removes the owner, owning group, and permissions for a file or directory (dfs1).
func (d *Client) RemoveAccessControlRecursive(ctx context.Context, options *RemoveAccessControlRecursiveOptions) (RemoveAccessControlRecursiveResponse, error) {
	// TODO explicitly pass SetAccessControlRecursiveMode
	return SetAccessControlRecursiveResponse{}, nil
}

// SetMetadata sets the metadata for a file or directory (blob3).
func (d *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	// TODO: call directly into blob
	return SetMetadataResponse{}, nil
}

// SetHTTPHeaders sets the HTTP headers for a file or directory (blob3).
func (d *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	// TODO: call formatBlobHTTPHeaders() since we want to add the blob prefix to our options before calling into blob
	// TODO: call into blob
	return SetHTTPHeadersResponse{}, nil
}

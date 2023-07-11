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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.PathClient, generated.PathClient, blob.Client]

// NewClient creates an instance of Client with the specified values.
//   - fileURL - the URL of the blob e.g. https://<account>.dfs.core.windows.net/fs/file.txt
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(fileURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	blobURL, fileURL := shared.GetURLs(fileURL)
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
	blobClientOpts := blob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blob.NewClient(blobURL, cred, &blobClientOpts)
	fileClient := base.NewPathClient(fileURL, blobURL, blobClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fileClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - fileURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/file.txt?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(fileURL string, options *ClientOptions) (*Client, error) {
	blobURL, fileURL := shared.GetURLs(fileURL)

	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FileClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobClientOpts := blob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobClient, _ := blob.NewClientWithNoCredential(blobURL, &blobClientOpts)
	fileClient := base.NewPathClient(fileURL, blobURL, blobClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fileClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - fileURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs/file.txt
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(fileURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	blobURL, fileURL := shared.GetURLs(fileURL)

	authPolicy := exported.NewSharedKeyCredPolicy(cred)
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
	blobClientOpts := blob.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKey, err := cred.ConvertToBlobSharedKey()
	if err != nil {
		return nil, err
	}
	blobClient, _ := blob.NewClientWithSharedKeyCredential(blobURL, blobSharedKey, &blobClientOpts)
	fileClient := base.NewPathClient(fileURL, blobURL, blobClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fileClient), nil
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

func (f *Client) generatedFileClientWithDFS() *generated.PathClient {
	//base.SharedKeyComposite((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	dirClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blob.Client])(f))
	return dirClientWithDFS
}

func (f *Client) generatedFileClientWithBlob() *generated.PathClient {
	_, dirClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blob.Client])(f))
	return dirClientWithBlob
}

func (f *Client) blobClient() *blob.Client {
	_, _, blobClient := base.InnerClients((*base.CompositeClient[generated.PathClient, generated.PathClient, blob.Client])(f))
	return blobClient
}

func (f *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.PathClient, generated.PathClient, blob.Client])(f))
}

func (f *Client) getClientOptions() *base.ClientOptions {
	return base.GetCompositeClientOptions((*base.CompositeClient[generated.PathClient, generated.PathClient, blob.Client])(f))
}

// DFSURL returns the URL endpoint used by the Client object.
func (f *Client) DFSURL() string {
	return f.generatedFileClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (f *Client) BlobURL() string {
	return f.generatedFileClientWithBlob().Endpoint()
}

// Create creates a new file (dfs1).
func (f *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	lac, mac, httpHeaders, createOpts, cpkOpts := options.format()
	return f.generatedFileClientWithDFS().Create(ctx, createOpts, httpHeaders, lac, mac, nil, cpkOpts)
}

// Delete deletes a file (dfs1).
func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	lac, mac, deleteOpts := options.format()
	return f.generatedFileClientWithDFS().Delete(ctx, deleteOpts, lac, mac)
}

// GetProperties gets the properties of a file (blob3)
func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	// TODO format response
	return f.blobClient().GetProperties(ctx, opts)
}

// Rename renames a file (dfs1).
func (f *Client) Rename(ctx context.Context, newName string, options *RenameOptions) (RenameResponse, error) {
	lac, mac, smac, createOpts := options.format()
	fileURL := runtime.JoinPaths(f.generatedFileClientWithDFS().Endpoint(), newName)
	// TODO: remove new azcore.Client creation after the API for shallow copying with new client name is implemented
	clOpts := f.getClientOptions()
	azClient, err := azcore.NewClient(shared.FileClient, exported.ModuleVersion, *(base.GetPipelineOptions(clOpts)), &(clOpts.ClientOptions))
	if err != nil {
		if log.Should(exported.EventError) {
			log.Writef(exported.EventError, err.Error())
		}
		return RenameResponse{}, err
	}
	fileURL, blobURL := shared.GetURLs(fileURL)
	tempFileClient := (*Client)(base.NewPathClient(fileURL, blobURL, nil, azClient, f.sharedKey(), clOpts))
	// this tempClient does not have a blobClient
	return tempFileClient.generatedFileClientWithDFS().Create(ctx, createOpts, nil, lac, mac, smac, nil)
}

// SetExpiry operation sets an expiry time on an existing file (blob2).
func (f *Client) SetExpiry(ctx context.Context, expiryType ExpiryType, o *SetExpiryOptions) (SetExpiryResponse, error) {
	expMode, opts := expiryType.Format(o)
	return f.generatedFileClientWithBlob().SetExpiry(ctx, expMode, opts)
}

//// Upload uploads data to a file.
//func (f *Client) Upload(ctx context.Context) {
//
//}
//
//// Append appends data to a file.
//func (f *Client) Append(ctx context.Context) {
//
//}
//
//// Flush flushes previous uploaded data to a file.
//func (f *Client) Flush(ctx context.Context) {
//
//}
//
//// Download downloads data from a file.
//func (f *Client) Download(ctx context.Context) {
//
//}

// SetAccessControl sets the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) SetAccessControl(ctx context.Context, options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	opts, lac, mac := options.format()
	return f.generatedFileClientWithDFS().SetAccessControl(ctx, opts, lac, mac)
}

// UpdateAccessControl updates the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) UpdateAccessControl(ctx context.Context, options *UpdateAccessControlOptions) (UpdateAccessControlResponse, error) {
	opts, mode := options.format()
	return f.generatedFileClientWithDFS().SetAccessControlRecursive(ctx, mode, opts)
}

// GetAccessControl gets the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) GetAccessControl(ctx context.Context, options *GetAccessControlOptions) (GetAccessControlResponse, error) {
	opts, lac, mac := options.format()
	return f.generatedFileClientWithDFS().GetProperties(ctx, opts, lac, mac)
}

// RemoveAccessControl removes the owner, owning group, and permissions for a file or directory (dfs1).
func (f *Client) RemoveAccessControl(ctx context.Context, options *RemoveAccessControlOptions) (RemoveAccessControlResponse, error) {
	opts, mode := options.format()
	return f.generatedFileClientWithDFS().SetAccessControlRecursive(ctx, mode, opts)
}

// SetMetadata sets the metadata for a file or directory (blob3).
func (f *Client) SetMetadata(ctx context.Context, metadata map[string]*string, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts := options.format()
	return f.blobClient().SetMetadata(ctx, metadata, opts)
}

// SetHTTPHeaders sets the HTTP headers for a file or directory (blob3).
func (f *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, blobHTTPHeaders := options.format(httpHeaders)
	// TODO: format response since there is a blob sequence header in the response
	return f.blobClient().SetHTTPHeaders(ctx, blobHTTPHeaders, opts)
}

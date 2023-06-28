//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"net/http"
	"strings"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client]

// NewClient creates an instance of Client with the specified values.
//   - filesystemURL - the URL of the blob e.g. https://<account>.dfs.core.windows.net/fs
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(filesystemURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	containerURL := strings.Replace(filesystemURL, ".dfs.", ".blob.", 1)
	filesystemURL = strings.Replace(filesystemURL, ".blob.", ".dfs.", 1)

	authPolicy := shared.NewStorageChallengePolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FilesystemClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	containerClientOpts := container.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobContainerClient, _ := container.NewClient(containerURL, cred, &containerClientOpts)
	fsClient := base.NewFilesystemClient(filesystemURL, containerURL, blobContainerClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fsClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - filesystemURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(filesystemURL string, options *ClientOptions) (*Client, error) {
	containerURL := strings.Replace(filesystemURL, ".dfs.", ".blob.", 1)
	filesystemURL = strings.Replace(filesystemURL, ".blob.", ".dfs.", 1)

	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FilesystemClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	containerClientOpts := container.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobContainerClient, _ := container.NewClientWithNoCredential(containerURL, &containerClientOpts)
	fsClient := base.NewFilesystemClient(filesystemURL, containerURL, blobContainerClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(fsClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - filesystemURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/fs
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(filesystemURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	containerURL := strings.Replace(filesystemURL, ".dfs.", ".blob.", 1)
	filesystemURL = strings.Replace(filesystemURL, ".blob.", ".dfs.", 1)

	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FilesystemClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	containerClientOpts := container.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKeyCredential, _ := blob.NewSharedKeyCredential(cred.AccountName(), cred.AccountKey())
	blobContainerClient, _ := container.NewClientWithSharedKeyCredential(containerURL, blobSharedKeyCredential, &containerClientOpts)
	fsClient := base.NewFilesystemClient(filesystemURL, containerURL, blobContainerClient, azClient, cred, (*base.ClientOptions)(conOptions))

	return (*Client)(fsClient), nil
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

func (fs *Client) generatedFSClientWithDFS() *generated.FileSystemClient {
	//base.SharedKeyComposite((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	fsClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client])(fs))
	return fsClientWithDFS
}

func (fs *Client) generatedFSClientWithBlob() *generated.FileSystemClient {
	_, fsClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client])(fs))
	return fsClientWithBlob
}

func (fs *Client) containerClient() *container.Client {
	_, _, containerClient := base.InnerClients((*base.CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client])(fs))
	return containerClient
}

func (fs *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client])(fs))
}

// URL returns the URL endpoint used by the Client object.
func (fs *Client) URL() string {
	return "s.generated().Endpoint()"
}

// Create creates a new filesystem under the specified account. (blob3).
func (fs *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	opts := options.format()
	return fs.containerClient().Create(ctx, opts)
}

// Delete deletes the specified filesystem and any files or directories it contains. (blob3).
func (fs *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	opts := options.format()
	return fs.containerClient().Delete(ctx, opts)
}

// GetProperties returns all user-defined metadata, standard HTTP properties, and system properties for the filesystem. (blob3).
func (fs *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	newResp := GetPropertiesResponse{}
	resp, err := fs.containerClient().GetProperties(ctx, opts)
	formatFilesystemProperties(&newResp, &resp)
	return newResp, err
}

// SetMetadata sets one or more user-defined name-value pairs for the specified filesystem. (blob3).
func (fs *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	opts := options.format()
	return fs.containerClient().SetMetadata(ctx, opts)
}

// SetAccessPolicy sets the permissions for the specified filesystem or the files and directories under it. (blob3).
func (fs *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	opts := options.format()
	return fs.containerClient().SetAccessPolicy(ctx, opts)
}

// GetAccessPolicy returns the permissions for the specified filesystem or the files and directories under it. (blob3).
func (fs *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	opts := options.format()
	return fs.containerClient().GetAccessPolicy(ctx, opts)
}

// TODO: implement undelete path in fs client as well

// NewListPathsPager operation returns a pager of the shares under the specified account. (dfs1)
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (fs *Client) NewListPathsPager(recursive bool, options *ListPathsOptions) *runtime.Pager[ListPathsSegmentResponse] {
	//TODO: look into possibility of using blob endpoint like list deleted paths is
	//TODO: will use ListPathsCreateRequest
	listOptions := options.format()
	return runtime.NewPager(runtime.PagingHandler[ListPathsSegmentResponse]{
		More: func(page ListPathsSegmentResponse) bool {
			return page.Continuation != nil && len(*page.Continuation) > 0
		},
		Fetcher: func(ctx context.Context, page *ListPathsSegmentResponse) (ListPathsSegmentResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = fs.generatedFSClientWithDFS().ListPathsCreateRequest(ctx, recursive, &listOptions)
			} else {
				listOptions.Continuation = page.Continuation
				req, err = fs.generatedFSClientWithDFS().ListPathsCreateRequest(ctx, recursive, &listOptions)
			}
			if err != nil {
				return ListPathsSegmentResponse{}, err
			}
			resp, err := fs.generatedFSClientWithDFS().InternalClient().Pipeline().Do(req)
			if err != nil {
				return ListPathsSegmentResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListPathsSegmentResponse{}, runtime.NewResponseError(resp)
			}
			return fs.generatedFSClientWithDFS().ListPathsHandleResponse(resp)
		},
	})
}

// NewListDeletedPathsPager operation returns a pager of the shares under the specified account. (dfs op/blob2).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (fs *Client) NewListDeletedPathsPager(options *ListDeletedPathsOptions) *runtime.Pager[ListDeletedPathsSegmentResponse] {
	//TODO: will use ListBlobHierarchySegmentCreateRequest
	return nil
}

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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"net/http"
	"time"
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
	containerURL, filesystemURL := shared.GetURLS(filesystemURL)
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FilesystemClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
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
	containerURL, filesystemURL := shared.GetURLS(filesystemURL)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.FilesystemClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
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
	containerURL, filesystemURL := shared.GetURLS(filesystemURL)
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

	if options == nil {
		options = &ClientOptions{}
	}
	containerClientOpts := container.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKey, err := cred.ConvertToBlobSharedKey()
	if err != nil {
		return nil, err
	}
	blobContainerClient, _ := container.NewClientWithSharedKeyCredential(containerURL, blobSharedKey, &containerClientOpts)
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

// DFSURL returns the URL endpoint used by the Client object.
func (fs *Client) DFSURL() string {
	return fs.generatedFSClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (fs *Client) BlobURL() string {
	return fs.generatedFSClientWithBlob().Endpoint()
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
	// TODO: find a cleaner way to not use lease from blob package
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
	newResp := GetAccessPolicyResponse{}
	resp, err := fs.containerClient().GetAccessPolicy(ctx, opts)
	formatGetAccessPolicyResponse(&newResp, &resp)
	return newResp, err
}

// TODO: implement undelete path in fs client as well

// NewListPathsPager operation returns a pager of the shares under the specified account. (dfs1)
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (fs *Client) NewListPathsPager(recursive bool, options *ListPathsOptions) *runtime.Pager[ListPathsSegmentResponse] {
	//TODO: look into possibility of using blob endpoint like list deleted paths is
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
	listOptions := options.format()
	return runtime.NewPager(runtime.PagingHandler[ListDeletedPathsSegmentResponse]{
		More: func(page ListDeletedPathsSegmentResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListDeletedPathsSegmentResponse) (ListDeletedPathsSegmentResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = fs.generatedFSClientWithDFS().ListBlobHierarchySegmentCreateRequest(ctx, &listOptions)
			} else {
				listOptions.Marker = page.NextMarker
				req, err = fs.generatedFSClientWithDFS().ListBlobHierarchySegmentCreateRequest(ctx, &listOptions)
			}
			if err != nil {
				return ListDeletedPathsSegmentResponse{}, err
			}
			resp, err := fs.generatedFSClientWithDFS().InternalClient().Pipeline().Do(req)
			if err != nil {
				return ListDeletedPathsSegmentResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListDeletedPathsSegmentResponse{}, runtime.NewResponseError(resp)
			}
			return fs.generatedFSClientWithDFS().ListBlobHierarchySegmentHandleResponse(resp)
		},
	})
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at container.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (fs *Client) GetSASURL(permissions sas.FilesystemPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if fs.sharedKey() == nil {
		return "", datalakeerror.MissingSharedKeyCredential
	}
	st := o.format()
	urlParts, err := azdatalake.ParseURL(fs.BlobURL())
	if err != nil {
		return "", err
	}
	qps, err := sas.DatalakeSignatureValues{
		Version:        sas.Version,
		Protocol:       sas.ProtocolHTTPS,
		FilesystemName: urlParts.FilesystemName,
		Permissions:    permissions.String(),
		StartTime:      st,
		ExpiryTime:     expiry.UTC(),
	}.SignWithSharedKey(fs.sharedKey())
	if err != nil {
		return "", err
	}

	endpoint := fs.BlobURL() + "?" + qps.Encode()

	return endpoint, nil
}

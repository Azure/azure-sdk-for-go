//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// Client represents a URL to the Azure Datalake Storage service allowing you to manipulate filesystems.
type Client base.Client[generated.FileSystemClient]

// NewClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(serviceURL string, cred azcore.TokenCredential, options *azdatalake.ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(serviceURL string, options *azdatalake.ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *exported.SharedKeyCredential, options *azdatalake.ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *azdatalake.ClientOptions) (*Client, error) {
	return nil, nil
}

func (fs *Client) generated() *generated.FileSystemClient {
	return base.InnerClient((*base.Client[generated.FileSystemClient])(fs))
}

func (fs *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.FileSystemClient])(fs))
}

// URL returns the URL endpoint used by the Client object.
func (fs *Client) URL() string {
	return "s.generated().Endpoint()"
}

// Create creates a new filesystem under the specified account. (blob3).
func (fs *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// Delete deletes the specified filesystem and any files or directories it contains. (blob3).
func (fs *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

// GetProperties returns all user-defined metadata, standard HTTP properties, and system properties for the filesystem. (blob3).
func (fs *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	// TODO: format blob response to fs response
	return GetPropertiesResponse{}, nil
}

// SetMetadata sets one or more user-defined name-value pairs for the specified filesystem. (blob3).
func (fs *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	return SetMetadataResponse{}, nil
}

// SetAccessPolicy sets the permissions for the specified filesystem or the files and directories under it. (blob3).
func (fs *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	return SetAccessPolicyResponse{}, nil
}

// GetAccessPolicy returns the permissions for the specified filesystem or the files and directories under it. (blob3).
func (fs *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	return GetAccessPolicyResponse{}, nil
}

// UndeletePath restores the specified path that was previously deleted. (dfs op/blob2).
func (fs *Client) UndeletePath(ctx context.Context, path string, options *UndeletePathOptions) (UndeletePathResponse, error) {
	return UndeletePathResponse{}, nil
}

// NewListPathsPager operation returns a pager of the shares under the specified account. (dfs1)
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (fs *Client) NewListPathsPager(recursive bool, options *ListPathsOptions) *runtime.Pager[ListPathsSegmentResponse] {
	//TODO: look into possibility of using blob endpoint like list deleted paths is
	//TODO: will use ListPathsCreateRequest
	return nil
}

// NewListDeletedPathsPager operation returns a pager of the shares under the specified account. (dfs op/blob2).
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (fs *Client) NewListDeletedPathsPager(options *ListDeletedPathsOptions) *runtime.Pager[ListDeletedPathsSegmentResponse] {
	//TODO: will use ListBlobHierarchySegmentCreateRequest
	return nil
}

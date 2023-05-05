//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
type Client base.Client[generated.FileSystemClient]

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

func (f *Client) generated() *generated.FileSystemClient {
	return base.InnerClient((*base.Client[generated.FileSystemClient])(f))
}

func (f *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.FileSystemClient])(f))
}

// URL returns the URL endpoint used by the Client object.
func (f *Client) URL() string {
	return "s.generated().Endpoint()"
}

func (f *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

func (f *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

func (f *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil
}

func (f *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	return SetMetadataResponse{}, nil
}

func (f *Client) ListPaths(ctx context.Context, options *ListPathsOptions) (ListPathsResponse, error) {
	return ListPathsResponse{}, nil
}

func (f *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	return SetAccessPolicyResponse{}, nil
}

func (f *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	return GetAccessPolicyResponse{}, nil
}

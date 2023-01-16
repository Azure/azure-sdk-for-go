//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure File Storage share allowing you to manipulate files.
type Client base.Client[generated.ShareClient]

// NewClient creates an instance of Client with the specified values.
//   - shareURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/share
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(shareURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a storage account or with a shared access signature (SAS) token.
//   - shareURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/share?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(shareURL string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - shareURL - the URL of the storage account e.g. https://<account>.file.core.windows.net/share
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(shareURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

func (s *Client) generated() *generated.ShareClient {
	return base.InnerClient((*base.Client[generated.ShareClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ShareClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return "s.generated().Endpoint()"
}

// WithSnapshot creates a new Client object identical to the source but with the specified share snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base share.
func (s *Client) WithSnapshot(shareSnapshot string) (*Client, error) {
	return nil, nil
}

// Create creates a new share within a storage account. If a share with the same name already exists, the operation fails.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/create-share
func (s *Client) Create(ctx context.Context, options *CreateOptions) (CreateResponse, error) {
	return CreateResponse{}, nil
}

// Delete marks the specified share for deletion. The share and any files contained within it are later deleted during garbage collection.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/delete-share
func (s *Client) Delete(ctx context.Context, options *DeleteOptions) (DeleteResponse, error) {
	return DeleteResponse{}, nil
}

// Restore operation restores a share that had previously been soft-deleted.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/restore-share
func (s *Client) Restore(ctx context.Context, deletedShareVersion string, options *RestoreOptions) (RestoreResponse, error) {
	return RestoreResponse{}, nil
}

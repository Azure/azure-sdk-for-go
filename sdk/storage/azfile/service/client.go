//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a URL to the Azure File Storage service allowing you to manipulate file shares.
type Client base.Client[generated.ServiceClient]

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
func NewClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

func (s *Client) generated() *generated.ServiceClient {
	return base.InnerClient((*base.Client[generated.ServiceClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ServiceClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return "s.generated().Endpoint()"
}

// NewShareClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewShareClient(shareName string) *share.Client {
	return nil
}

// CreateShare is a lifecycle method to creates a new share under the specified account.
// If the share with the same name already exists, a ResourceExistsError will be raised.
// This method returns a client with which to interact with the newly created share.
// For more information see, https://learn.microsoft.com/en-us/rest/api/storageservices/create-share.
func (s *Client) CreateShare(ctx context.Context, shareName string, options *CreateShareOptions) (CreateShareResponse, error) {
	shareClient := s.NewShareClient(shareName)
	createShareResp, err := shareClient.Create(ctx, options)
	return createShareResp, err
}

// DeleteShare is a lifecycle method that marks the specified share for deletion.
// The share and any files contained within it are later deleted during garbage collection.
// If the share is not found, a ResourceNotFoundError will be raised.
// For more information see, https://learn.microsoft.com/en-us/rest/api/storageservices/delete-share.
func (s *Client) DeleteShare(ctx context.Context, shareName string, options *DeleteShareOptions) (DeleteShareResponse, error) {
	shareClient := s.NewShareClient(shareName)
	deleteShareResp, err := shareClient.Delete(ctx, options)
	return deleteShareResp, err
}

// RestoreShare restores soft-deleted share.
// Operation will only be successful if used within the specified number of days set in the delete retention policy.
// For more information see, https://learn.microsoft.com/en-us/rest/api/storageservices/restore-share.
func (s *Client) RestoreShare(ctx context.Context, deletedShareName string, deletedShareVersion string, options *RestoreShareOptions) (RestoreShareResponse, error) {
	shareClient := s.NewShareClient(deletedShareName)
	createShareResp, err := shareClient.Restore(ctx, deletedShareVersion, options)
	return createShareResp, err
}

// GetProperties operation gets the properties of a storage account's File service.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-file-service-properties.
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	resp, err := s.generated().GetProperties(ctx, opts)
	return resp, err
}

// SetProperties operation sets properties for a storage account's File service endpoint.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/set-file-service-properties.
func (s *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	return SetPropertiesResponse{}, nil
}

// NewListSharesPager operation returns a pager of the shares under the specified account.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (s *Client) NewListSharesPager(options *ListSharesOptions) *runtime.Pager[ListSharesSegmentResponse] {
	return nil
}

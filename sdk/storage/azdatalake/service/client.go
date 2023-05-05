//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
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
func NewClientWithSharedKeyCredential(serviceURL string, cred *exported.SharedKeyCredential, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	return nil, nil
}

// NewFilesystemClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewFilesystemClient(filesystemName string) *filesystem.Client {
	return nil
}

// NewDirectoryClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewDirectoryClient(directoryName string) *filesystem.Client {
	return nil
}

// NewFileClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewFileClient(fileName string) *filesystem.Client {
	return nil
}

func (s *Client) generated() *generated.ServiceClient {
	return base.InnerClient((*base.Client[generated.ServiceClient])(s))
}

func (s *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ServiceClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return "s.generated().Endpoint()"
}

func (s *Client) CreateFilesystem(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	createShareResp, err := filesystemClient.Create(ctx, options)
	return createShareResp, err
}

func (s *Client) DeleteFilesystem(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	createShareResp, err := filesystemClient.Create(ctx, options)
	return createShareResp, err
}

func (s *Client) UndeleteFilesystem(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	createShareResp, err := filesystemClient.Create(ctx, options)
	return createShareResp, err
}

func (s *Client) SetServiceProperties(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	createShareResp, err := filesystemClient.Create(ctx, options)
	return createShareResp, err
}

func (s *Client) GetServiceProperties(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	createShareResp, err := filesystemClient.Create(ctx, options)
	return createShareResp, err
}

func (s *Client) ListFilesystems(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	createShareResp, err := filesystemClient.Create(ctx, options)
	return createShareResp, err
}

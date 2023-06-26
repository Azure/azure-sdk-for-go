//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"strings"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client]

// NewClientWithNoCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/
//   - options - client options; pass nil to accept the default values.
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	blobServiceURL := strings.Replace(serviceURL, ".dfs.", ".blob.", 1)
	datalakeServiceURL := strings.Replace(serviceURL, ".blob.", ".dfs.", 1)

	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.ServiceClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSvcClient, _ := service.NewClientWithNoCredential(blobServiceURL, &blobServiceClientOpts)
	svcClient := base.NewServiceClient(datalakeServiceURL, blobServiceURL, blobSvcClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(svcClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	blobServiceURL := strings.Replace(serviceURL, ".dfs.", ".blob.", 1)
	datalakeServiceURL := strings.Replace(serviceURL, ".blob.", ".dfs.", 1)

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

	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKeyCredential, _ := blob.NewSharedKeyCredential(cred.AccountName(), cred.AccountKey())
	blobSvcClient, _ := service.NewClientWithSharedKeyCredential(blobServiceURL, blobSharedKeyCredential, &blobServiceClientOpts)
	svcClient := base.NewServiceClient(datalakeServiceURL, blobServiceURL, blobSvcClient, azClient, cred, (*base.ClientOptions)(conOptions))

	return (*Client)(svcClient), nil
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

func (s *Client) generatedFSClientWithDFS() *generated.ServiceClient {
	svcClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
	return svcClientWithDFS
}

func (s *Client) generatedFSClientWithBlob() *generated.ServiceClient {
	_, svcClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
	return svcClientWithBlob
}

func (s *Client) containerClient() *service.Client {
	_, _, serviceClient := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
	return serviceClient
}

func (s *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return "s.generated().Endpoint()"
}

// CreateFilesystem creates a new filesystem under the specified account. (blob3)
func (s *Client) CreateFilesystem(ctx context.Context, filesystem string, options *CreateFilesystemOptions) (CreateFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	resp, err := filesystemClient.Create(ctx, options)
	return resp, err
}

// DeleteFilesystem deletes the specified filesystem. (blob3)
func (s *Client) DeleteFilesystem(ctx context.Context, filesystem string, options *DeleteFilesystemOptions) (DeleteFilesystemResponse, error) {
	filesystemClient := s.NewFilesystemClient(filesystem)
	resp, err := filesystemClient.Delete(ctx, options)
	return resp, err
}

// SetServiceProperties sets properties for a storage account's File service endpoint. (blob3)
func (s *Client) SetServiceProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	return SetPropertiesResponse{}, nil
}

// GetProperties gets properties for a storage account's File service endpoint. (blob3)
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	return GetPropertiesResponse{}, nil
}

// NewListFilesystemsPager operation returns a pager of the shares under the specified account. (blob3)
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (s *Client) NewListFilesystemsPager(options *ListFilesystemsOptions) *runtime.Pager[ListFilesystemsResponse] {
	return nil
}

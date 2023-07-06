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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
	"time"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client]

// NewClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the blob e.g. https://<account>.dfs.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	blobServiceURL, datalakeServiceURL := shared.GetURLS(serviceURL)
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.ServiceClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSvcClient, _ := service.NewClient(blobServiceURL, cred, &blobServiceClientOpts)
	svcClient := base.NewServiceClient(datalakeServiceURL, blobServiceURL, blobSvcClient, azClient, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(svcClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/
//   - options - client options; pass nil to accept the default values.
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	blobServiceURL, datalakeServiceURL := shared.GetURLS(serviceURL)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(shared.ServiceClient, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
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
	blobServiceURL, datalakeServiceURL := shared.GetURLS(serviceURL)
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

	if options == nil {
		options = &ClientOptions{}
	}
	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKey, err := cred.ConvertToBlobSharedKey()
	if err != nil {
		return nil, err
	}
	blobSvcClient, _ := service.NewClientWithSharedKeyCredential(blobServiceURL, blobSharedKey, &blobServiceClientOpts)
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

func (s *Client) getClientOptions() *base.ClientOptions {
	return base.GetCompositeClientOptions((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
}

// NewFilesystemClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewFilesystemClient(filesystemName string) *filesystem.Client {
	filesystemURL := runtime.JoinPaths(s.generatedServiceClientWithDFS().Endpoint(), filesystemName)
	// TODO: remove new azcore.Client creation after the API for shallow copying with new client name is implemented
	clOpts := s.getClientOptions()
	azClient, err := azcore.NewClient(shared.FilesystemClient, exported.ModuleVersion, *(base.GetPipelineOptions(clOpts)), &(clOpts.ClientOptions))
	if err != nil {
		if log.Should(exported.EventError) {
			log.Writef(exported.EventError, err.Error())
		}
		return nil
	}
	filesystemURL, containerURL := shared.GetURLS(filesystemURL)
	return (*filesystem.Client)(base.NewFilesystemClient(filesystemURL, containerURL, s.serviceClient().NewContainerClient(filesystemName), azClient, s.sharedKey(), clOpts))
}

// NewDirectoryClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewDirectoryClient(directoryName string) *filesystem.Client {
	// TODO: implement once dir client is implemented
	return nil
}

// NewFileClient creates a new share.Client object by concatenating shareName to the end of this Client's URL.
// The new share.Client uses the same request policy pipeline as the Client.
func (s *Client) NewFileClient(fileName string) *filesystem.Client {
	// TODO: implement once file client is implemented
	return nil
}

func (s *Client) generatedServiceClientWithDFS() *generated.ServiceClient {
	svcClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
	return svcClientWithDFS
}

func (s *Client) generatedServiceClientWithBlob() *generated.ServiceClient {
	_, svcClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
	return svcClientWithBlob
}

func (s *Client) serviceClient() *service.Client {
	_, _, serviceClient := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
	return serviceClient
}

func (s *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client])(s))
}

// DFSURL returns the URL endpoint used by the Client object.
func (s *Client) DFSURL() string {
	return s.generatedServiceClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (s *Client) BlobURL() string {
	return s.generatedServiceClientWithBlob().Endpoint()
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

// SetProperties sets properties for a storage account's File service endpoint. (blob3)
func (s *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	opts := options.format()
	return s.serviceClient().SetProperties(ctx, opts)
}

// GetProperties gets properties for a storage account's File service endpoint. (blob3)
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	return s.serviceClient().GetProperties(ctx, opts)

}

// NewListFilesystemsPager operation returns a pager of the shares under the specified account. (blob3)
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (s *Client) NewListFilesystemsPager(options *ListFilesystemsOptions) *runtime.Pager[ListFilesystemsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListFilesystemsResponse]{
		More: func(page ListFilesystemsResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListFilesystemsResponse) (ListFilesystemsResponse, error) {
			if page == nil {
				page = &ListFilesystemsResponse{}
				opts := options.format()
				page.blobPager = s.serviceClient().NewListContainersPager(opts)
			}
			newPage := ListFilesystemsResponse{}
			currPage, err := page.blobPager.NextPage(context.TODO())
			if err != nil {
				return newPage, err
			}
			newPage.Prefix = currPage.Prefix
			newPage.Marker = currPage.Marker
			newPage.MaxResults = currPage.MaxResults
			newPage.NextMarker = currPage.NextMarker
			newPage.Filesystems = convertContainerItemsToFSItems(currPage.ContainerItems)
			newPage.ServiceEndpoint = currPage.ServiceEndpoint
			newPage.blobPager = page.blobPager

			return newPage, err
		},
	})
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (s *Client) GetSASURL(resources sas.AccountResourceTypes, permissions sas.AccountPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	// format all options to blob service options
	res, perms, opts := o.format(resources, permissions)
	return s.serviceClient().GetSASURL(res, perms, expiry, opts)
}

// TODO: Figure out how we can convert from blob delegation key to one defined in datalake
//// GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceURL object.
//// OAuth is required for this call, as well as any role that can delegate access to the storage account.
//func (s *Client) GetUserDelegationCredential(ctx context.Context, info KeyInfo, o *GetUserDelegationCredentialOptions) (*UserDelegationCredential, error) {
//	return s.serviceClient().GetUserDelegationCredential(ctx, info, o)
//}

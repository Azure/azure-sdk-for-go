// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated_blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
)

// FOR SERVICE CLIENT WE STORE THE GENERATED BLOB LAYER IN ORDER TO USE FS LISTING AND THE TRANSFORMS IT HAS

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service.
type Client base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client]

// NewClient creates an instance of Client with the specified values.
//   - serviceURL - the URL of the blob e.g. https://<account>.dfs.core.windows.net/
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	blobServiceURL, datalakeServiceURL := shared.GetURLs(serviceURL)
	audience := base.GetAudience((*base.ClientOptions)(options))
	conOptions := shared.GetClientOptions(options)
	authPolicy := shared.NewStorageChallengePolicy(cred, audience, conOptions.InsecureAllowCredentialWithHTTP)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	perCallPolicies := []policy.Policy{shared.NewIncludeBlobResponsePolicy()}
	if options.PerCallPolicies != nil {
		perCallPolicies = append(perCallPolicies, options.PerCallPolicies...)
	}
	options.PerCallPolicies = perCallPolicies
	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSvcClient, _ := service.NewClient(blobServiceURL, cred, &blobServiceClientOpts)
	svcClient := base.NewServiceClient(datalakeServiceURL, blobServiceURL, blobSvcClient, azClient, nil, &cred, (*base.ClientOptions)(conOptions))

	return (*Client)(svcClient), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/
//   - options - client options; pass nil to accept the default values.
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	blobServiceURL, datalakeServiceURL := shared.GetURLs(serviceURL)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	perCallPolicies := []policy.Policy{shared.NewIncludeBlobResponsePolicy()}
	if options.PerCallPolicies != nil {
		perCallPolicies = append(perCallPolicies, options.PerCallPolicies...)
	}
	options.PerCallPolicies = perCallPolicies
	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSvcClient, _ := service.NewClientWithNoCredential(blobServiceURL, &blobServiceClientOpts)
	svcClient := base.NewServiceClient(datalakeServiceURL, blobServiceURL, blobSvcClient, azClient, nil, nil, (*base.ClientOptions)(conOptions))

	return (*Client)(svcClient), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - serviceURL - the URL of the storage account e.g. https://<account>.dfs.core.windows.net/
//   - cred - a SharedKeyCredential created with the matching storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	blobServiceURL, datalakeServiceURL := shared.GetURLs(serviceURL)
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}
	base.SetPipelineOptions((*base.ClientOptions)(conOptions), &plOpts)

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &ClientOptions{}
	}
	perCallPolicies := []policy.Policy{shared.NewIncludeBlobResponsePolicy()}
	if options.PerCallPolicies != nil {
		perCallPolicies = append(perCallPolicies, options.PerCallPolicies...)
	}
	options.PerCallPolicies = perCallPolicies
	blobServiceClientOpts := service.ClientOptions{
		ClientOptions: options.ClientOptions,
	}
	blobSharedKey, err := exported.ConvertToBlobSharedKey(cred)
	if err != nil {
		return nil, err
	}
	blobSvcClient, _ := service.NewClientWithSharedKeyCredential(blobServiceURL, blobSharedKey, &blobServiceClientOpts)
	svcClient := base.NewServiceClient(datalakeServiceURL, blobServiceURL, blobSvcClient, azClient, cred, nil, (*base.ClientOptions)(conOptions))

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
	return base.GetCompositeClientOptions((*base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client])(s))
}

// NewFileSystemClient creates a new filesystem.Client object by concatenating filesystemName to the end of this Client's URL.
// The new filesystem.Client uses the same request policy pipeline as the Client.
func (s *Client) NewFileSystemClient(filesystemName string) *filesystem.Client {
	filesystemURL := runtime.JoinPaths(s.generatedServiceClientWithDFS().Endpoint(), filesystemName)
	containerURL, filesystemURL := shared.GetURLs(filesystemURL)
	return (*filesystem.Client)(base.NewFileSystemClient(filesystemURL, containerURL, s.serviceClient().NewContainerClient(filesystemName), s.generatedServiceClientWithDFS().InternalClient().WithClientName(exported.ModuleName), s.sharedKey(), s.identityCredential(), s.getClientOptions()))
}

// GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceURL object.
// OAuth is required for this call, as well as any role that can delegate access to the storage account.
func (s *Client) GetUserDelegationCredential(ctx context.Context, info KeyInfo, o *GetUserDelegationCredentialOptions) (*UserDelegationCredential, error) {
	url, err := azdatalake.ParseURL(s.BlobURL())
	if err != nil {
		return nil, err
	}

	getUserDelegationKeyOptions := o.format()
	if o != nil && o.DelegatedUserTenantId != nil {
		info.DelegatedUserTid = o.DelegatedUserTenantId
	}
	udk, err := s.generatedServiceClientWithBlob().GetUserDelegationKey(ctx, info, getUserDelegationKeyOptions)
	if err != nil {
		return nil, exported.ConvertToDFSError(err)
	}

	return exported.NewUserDelegationCredential(strings.Split(url.Host, ".")[0], udk.UserDelegationKey), nil
}

func (s *Client) generatedServiceClientWithDFS() *generated.ServiceClient {
	svcClientWithDFS, _, _ := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client])(s))
	return svcClientWithDFS
}

func (s *Client) generatedServiceClientWithBlob() *generated_blob.ServiceClient {
	_, svcClientWithBlob, _ := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client])(s))
	return svcClientWithBlob
}

func (s *Client) serviceClient() *service.Client {
	_, _, serviceClient := base.InnerClients((*base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client])(s))
	return serviceClient
}

func (s *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client])(s))
}

func (s *Client) identityCredential() *azcore.TokenCredential {
	return base.IdentityCredentialComposite((*base.CompositeClient[generated.ServiceClient, generated_blob.ServiceClient, service.Client])(s))
}

// DFSURL returns the URL endpoint used by the Client object.
func (s *Client) DFSURL() string {
	return s.generatedServiceClientWithDFS().Endpoint()
}

// BlobURL returns the URL endpoint used by the Client object.
func (s *Client) BlobURL() string {
	return s.generatedServiceClientWithBlob().Endpoint()
}

// CreateFileSystem creates a new filesystem under the specified account.
func (s *Client) CreateFileSystem(ctx context.Context, filesystem string, options *CreateFileSystemOptions) (CreateFileSystemResponse, error) {
	filesystemClient := s.NewFileSystemClient(filesystem)
	resp, err := filesystemClient.Create(ctx, options)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// DeleteFileSystem deletes the specified filesystem.
func (s *Client) DeleteFileSystem(ctx context.Context, filesystem string, options *DeleteFileSystemOptions) (DeleteFileSystemResponse, error) {
	filesystemClient := s.NewFileSystemClient(filesystem)
	resp, err := filesystemClient.Delete(ctx, options)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// SetProperties sets properties for a storage account's Datalake service endpoint.
func (s *Client) SetProperties(ctx context.Context, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	opts := options.format()
	resp, err := s.serviceClient().SetProperties(ctx, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

// GetProperties gets properties for a storage account's Datalake service endpoint.
func (s *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts := options.format()
	resp, err := s.serviceClient().GetProperties(ctx, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err

}

// NewListFileSystemsPager operation returns a pager of the shares under the specified account.
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/list-shares
func (s *Client) NewListFileSystemsPager(o *ListFileSystemsOptions) *runtime.Pager[ListFileSystemsResponse] {
	listOptions := generated_blob.ServiceClientListContainersSegmentOptions{}
	defaultInclude := ListFileSystemsInclude{}
	if o != nil {
		if o.Include != defaultInclude && o.Include.Deleted != nil && *o.Include.Deleted {
			listOptions.Include = append(listOptions.Include, generated_blob.ListContainersIncludeTypeDeleted)
		}
		if o.Include != defaultInclude && o.Include.Metadata != nil && *o.Include.Metadata {
			listOptions.Include = append(listOptions.Include, generated_blob.ListContainersIncludeTypeMetadata)
		}
		if o.Include != defaultInclude && o.Include.System != nil && *o.Include.System {
			listOptions.Include = append(listOptions.Include, generated_blob.ListContainersIncludeTypeSystem)
		}
		listOptions.Marker = o.Marker
		listOptions.Maxresults = o.MaxResults
		listOptions.Prefix = o.Prefix
	}
	return runtime.NewPager(runtime.PagingHandler[ListFileSystemsResponse]{
		More: func(page ListFileSystemsResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListFileSystemsResponse) (ListFileSystemsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = s.generatedServiceClientWithBlob().ListContainersSegmentCreateRequest(ctx, &listOptions)
			} else {
				listOptions.Marker = page.NextMarker
				req, err = s.generatedServiceClientWithBlob().ListContainersSegmentCreateRequest(ctx, &listOptions)
			}
			if err != nil {
				return ListFileSystemsResponse{}, exported.ConvertToDFSError(err)
			}
			resp, err := s.generatedServiceClientWithBlob().InternalClient().Pipeline().Do(req)
			if err != nil {
				return ListFileSystemsResponse{}, exported.ConvertToDFSError(err)
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListFileSystemsResponse{}, exported.ConvertToDFSError(runtime.NewResponseError(resp))
			}
			resp1, err := s.generatedServiceClientWithBlob().ListContainersSegmentHandleResponse(resp)
			return resp1, exported.ConvertToDFSError(err)
		},
	})
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (s *Client) GetSASURL(resources sas.AccountResourceTypes, permissions sas.AccountPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	// format all options to blob service options
	res, perms, opts := o.format(resources, permissions)
	resp, err := s.serviceClient().GetSASURL(res, perms, expiry, opts)
	err = exported.ConvertToDFSError(err)
	return resp, err
}

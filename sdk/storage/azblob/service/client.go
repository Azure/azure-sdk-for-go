//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"net/http"
	"strings"
	"time"
)

// Client represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
type Client base.Client[generated.ServiceClient]

// NewClient creates a Client object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := exported.GetConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewServiceClient(serviceURL, pl)), nil
}

// NewClientWithNoCredential creates a Client object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	conOptions := exported.GetConnectionOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewServiceClient(serviceURL, pl)), nil
}

// NewClientWithSharedKey creates a Client object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := exported.GetConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewServiceClient(serviceURL, pl)), nil
}

// NewClientFromConnectionString creates a service client from the given connection string.
//nolint
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
		return NewClientWithSharedKey(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (s *Client) generated() *generated.ServiceClient {
	return base.InnerClient((*base.Client[generated.ServiceClient])(s))
}

func (s *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.ServiceClient])(s))
}

// URL returns the URL endpoint used by the Client object.
func (s *Client) URL() string {
	return s.generated().Endpoint()
}

// NewContainerClient creates a new ContainerClient object by concatenating containerName to the end of
// Client's URL. The new ContainerClient uses the same request policy pipeline as the Client.
// To change the pipeline, create the ContainerClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewContainerClient instead of calling this object's
// NewContainerClient method.
func (s *Client) NewContainerClient(containerName string) *container.Client {
	containerURL := runtime.JoinPaths(s.generated().Endpoint(), containerName)
	return (*container.Client)(base.NewContainerClient(containerURL, s.generated().Pipeline()))
}

// CreateContainer is a lifecycle method to creates a new container under the specified account.
// If the container with the same name already exists, a ResourceExistsError will be raised.
// This method returns a client with which to interact with the newly created container.
func (s *Client) CreateContainer(ctx context.Context, containerName string, options *CreateContainerOptions) (CreateContainerResponse, error) {
	containerClient := s.NewContainerClient(containerName)
	containerCreateResp, err := containerClient.Create(ctx, options)
	return containerCreateResp, err
}

// DeleteContainer is a lifecycle method that marks the specified container for deletion.
// The container and any blobs contained within it are later deleted during garbage collection.
// If the container is not found, a ResourceNotFoundError will be raised.
func (s *Client) DeleteContainer(ctx context.Context, containerName string, options *DeleteContainerOptions) (DeleteContainerResponse, error) {
	containerClient := s.NewContainerClient(containerName)
	containerDeleteResp, err := containerClient.Delete(ctx, options)
	return containerDeleteResp, err
}

// GetAccountInfo provides account level information
func (s *Client) GetAccountInfo(ctx context.Context, o *GetAccountInfoOptions) (GetAccountInfoResponse, error) {
	getAccountInfoOptions := o.format()
	resp, err := s.generated().GetAccountInfo(ctx, getAccountInfoOptions)
	return resp, err
}

// NewListContainersPager operation returns a pager of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (s *Client) NewListContainersPager(o *ListContainersOptions) *runtime.Pager[ListContainersResponse] {
	listOptions := generated.ServiceClientListContainersSegmentOptions{}
	if o != nil {
		if o.Include.Deleted {
			listOptions.Include = append(listOptions.Include, generated.ListContainersIncludeTypeDeleted)
		}
		if o.Include.Metadata {
			listOptions.Include = append(listOptions.Include, generated.ListContainersIncludeTypeMetadata)
		}
		listOptions.Marker = o.Marker
		listOptions.Maxresults = o.MaxResults
		listOptions.Prefix = o.Prefix
	}
	return runtime.NewPager(runtime.PagingHandler[ListContainersResponse]{
		More: func(page ListContainersResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ListContainersResponse) (ListContainersResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = s.generated().ListContainersSegmentCreateRequest(ctx, &listOptions)
			} else {
				listOptions.Marker = page.Marker
				req, err = s.generated().ListContainersSegmentCreateRequest(ctx, &listOptions)
			}
			if err != nil {
				return ListContainersResponse{}, err
			}
			resp, err := s.generated().Pipeline().Do(req)
			if err != nil {
				return ListContainersResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListContainersResponse{}, runtime.NewResponseError(resp)
			}
			return s.generated().ListContainersSegmentHandleResponse(resp)
		},
	})
}

// GetProperties - gets the properties of a storage account's Blob service, including properties for Storage Analytics
// and CORS (Cross-Origin Resource Sharing) rules.
func (s *Client) GetProperties(ctx context.Context, o *GetPropertiesOptions) (GetPropertiesResponse, error) {
	getPropertiesOptions := o.format()
	resp, err := s.generated().GetProperties(ctx, getPropertiesOptions)
	return resp, err
}

// SetProperties Sets the properties of a storage account's Blob service, including Azure Storage Analytics.
// If an element (e.g. analytics_logging) is left as None, the existing settings on the service for that functionality are preserved.
func (s *Client) SetProperties(ctx context.Context, o *SetPropertiesOptions) (SetPropertiesResponse, error) {
	properties, setPropertiesOptions := o.format()
	resp, err := s.generated().SetProperties(ctx, properties, setPropertiesOptions)
	return resp, err
}

// GetStatistics Retrieves statistics related to replication for the Blob service.
// It is only available when read-access geo-redundant replication is enabled for  the storage account.
// With geo-redundant replication, Azure Storage maintains your data durable
// in two locations. In both locations, Azure Storage constantly maintains
// multiple healthy replicas of your data. The location where you read,
// create, update, or delete data is the primary storage account location.
// The primary location exists in the region you choose at the time you
// create an account via the Azure Management Azure classic portal, for
// example, North Central US. The location to which your data is replicated
// is the secondary location. The secondary location is automatically
// determined based on the location of the primary; it is in a second data
// center that resides in the same region as the primary location. Read-only
// access is available from the secondary location, if read-access geo-redundant
// replication is enabled for your storage account.
func (s *Client) GetStatistics(ctx context.Context, o *GetStatisticsOptions) (GetStatisticsResponse, error) {
	getStatisticsOptions := o.format()
	resp, err := s.generated().GetStatistics(ctx, getStatisticsOptions)

	return resp, err
}

// SASResourceTypes type simplifies creating the resource types string for an Azure Storage Account SAS.
// Initialize an instance of this type and then call its String method to set AccountSASSignatureValues's ResourceTypes field.
type SASResourceTypes = exported.AccountSASResourceTypes

// SASServices type simplifies creating the services string for an Azure Storage Account SAS.
// Initialize an instance of this type and then call its String method to set SASServices' Services field.
type SASServices = exported.AccountSASServices

// SASPermissions type simplifies creating the permissions string for an Azure Storage Account SAS.
// Initialize an instance of this type and then call its String method to set AccountSASSignatureValues' Permissions field.
type SASPermissions = exported.AccountSASPermissions

const (
	SASProtocolHTTPS = exported.SASProtocolHTTPS

	SnapshotTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"
)

// SASSignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Storage account.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/constructing-an-account-sas
type SASSignatureValues = exported.AccountSASSignatureValues

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
// This validity can be checked with CanGetAccountSASToken().
func (s *Client) GetSASURL(resources SASResourceTypes, permissions SASPermissions, start time.Time, expiry time.Time) (string, error) {
	if s.sharedKey() == nil {
		return "", errors.New("SAS can only be signed with a SharedKeyCredential")
	}

	qps, err := SASSignatureValues{
		Version:       exported.SASVersion,
		Protocol:      exported.SASProtocolHTTPS,
		Permissions:   permissions.String(),
		Services:      "b",
		ResourceTypes: resources.String(),
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
	}.Sign(s.sharedKey())
	if err != nil {
		return "", err
	}

	endpoint := s.URL()
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	endpoint += "?" + qps.Encode()

	return endpoint, nil
}

//// FindBlobsByTags operation finds all blobs in the storage account whose tags match a given search expression.
//// Filter blobs searches across all containers within a storage account but can be scoped within the expression to a single container.
//// https://docs.microsoft.com/en-us/rest/api/storageservices/find-blobs-by-tags
//// eg. "dog='germanshepherd' and penguin='emperorpenguin'"
//// To specify a container, eg. "@container=’containerName’ and Name = ‘C’"
//func (s *Client) FindBlobsByTags(ctx context.Context, o *ServiceFilterBlobsOptions) (ServiceFilterBlobsResponse, error) {
//	// TODO: Use pager here? Missing support from zz_generated_pagers.go
//	serviceFilterBlobsOptions := o.pointer()
//	resp, err := s.client.FilterBlobs(ctx, serviceFilterBlobsOptions)
//	return toServiceFilterBlobsResponse(resp), err
//}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

//nolint
const (
	// ContainerNameRoot is the special Azure Storage name used to identify a storage account's root container.
	ContainerNameRoot = "$root"

	// ContainerNameLogs is the special Azure Storage name used to identify a storage account's logs container.
	ContainerNameLogs = "$logs"
)

// ServiceClient represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
type ServiceClient struct {
	client    *serviceClient
	sharedKey *SharedKeyCredential
}

// URL returns the URL endpoint used by the ServiceClient object.
func (s ServiceClient) URL() string {
	return s.client.endpoint
}

// NewServiceClient creates a ServiceClient object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewServiceClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) *ServiceClient {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{internal.TokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := internal.NewConnection(serviceURL, conOptions)

	return &ServiceClient{
		client: newServiceClient(conn.Endpoint(), conn.Pipeline()),
	}
}

// NewServiceClientWithNoCredential creates a ServiceClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewServiceClientWithNoCredential(serviceURL string, options *ClientOptions) *ServiceClient {
	conOptions := getConnectionOptions(options)
	conn := internal.NewConnection(serviceURL, conOptions)

	return &ServiceClient{
		client: newServiceClient(conn.Endpoint(), conn.Pipeline()),
	}
}

// NewServiceClientWithSharedKey creates a ServiceClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewServiceClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) *ServiceClient {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := internal.NewConnection(serviceURL, conOptions)

	return &ServiceClient{
		client:    newServiceClient(conn.Endpoint(), conn.Pipeline()),
		sharedKey: cred,
	}
}

// NewServiceClientFromConnectionString creates a service client from the given connection string.
//nolint
func NewServiceClientFromConnectionString(connectionString string, options *ClientOptions) *ServiceClient {
	endpoint, credential, _ := parseConnectionString(connectionString)
	return NewServiceClientWithSharedKey(endpoint, credential, options)
}

// NewContainerClient creates a new ContainerClient object by concatenating containerName to the end of
// ServiceClient's URL. The new ContainerClient uses the same request policy pipeline as the ServiceClient.
// To change the pipeline, create the ContainerClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewContainerClient instead of calling this object's
// NewContainerClient method.
func (s *ServiceClient) NewContainerClient(containerName string) *ContainerClient {
	containerURL := appendToURLPath(s.client.endpoint, containerName)
	return &ContainerClient{
		client:    newContainerClient(containerURL, s.client.pl),
		sharedKey: s.sharedKey,
	}
}

//NewServiceClientWithUserDelegationCredential obtains a UserDelegationKey object using the base ServiceURL object.
//OAuth is required for this call, as well as any role that can delegate access to the storage account.
func (s *ServiceClient) NewServiceClientWithUserDelegationCredential(ctx context.Context, info KeyInfo, timeout *int32, requestID *string) (UserDelegationCredential, error) {
	options := serviceClientGetUserDelegationKeyOptions{
		RequestID: requestID,
		Timeout:   timeout,
	}
	sc := newServiceClient(s.client.endpoint, s.client.pl)
	udk, err := sc.GetUserDelegationKey(ctx, info, &options)
	if err != nil {
		return UserDelegationCredential{}, err
	}
	return *NewUserDelegationCredential(strings.Split(s.client.endpoint, ".")[0], udk.UserDelegationKey), nil
}

// CreateContainer is a lifecycle method to creates a new container under the specified account.
// If the container with the same name already exists, a ResourceExistsError will be raised.
// This method returns a client with which to interact with the newly created container.
func (s *ServiceClient) CreateContainer(ctx context.Context, containerName string, options *ContainerCreateOptions) (ContainerCreateResponse, error) {
	containerClient := s.NewContainerClient(containerName)
	containerCreateResp, err := containerClient.Create(ctx, options)
	return containerCreateResp, err
}

// DeleteContainer is a lifecycle method that marks the specified container for deletion.
// The container and any blobs contained within it are later deleted during garbage collection.
// If the container is not found, a ResourceNotFoundError will be raised.
func (s *ServiceClient) DeleteContainer(ctx context.Context, containerName string, options *ContainerDeleteOptions) (ContainerDeleteResponse, error) {
	containerClient := s.NewContainerClient(containerName)
	containerDeleteResp, err := containerClient.Delete(ctx, options)
	return containerDeleteResp, err
}

// GetAccountInfo provides account level information
func (s *ServiceClient) GetAccountInfo(ctx context.Context, o *ServiceGetAccountInfoOptions) (ServiceGetAccountInfoResponse, error) {
	getAccountInfoOptions := o.format()
	resp, err := s.client.GetAccountInfo(ctx, getAccountInfoOptions)
	return toServiceGetAccountInfoResponse(resp), err
}

// NewListContainersPager operation returns a pager of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (s *ServiceClient) NewListContainersPager(o *ServiceListContainersOptions) *runtime.Pager[ServiceListContainersResponse] {
	listOptions := o.format()
	return runtime.NewPager(runtime.PagingHandler[ServiceListContainersResponse]{
		More: func(page ServiceListContainersResponse) bool {
			if page.Marker == nil || len(*page.Marker) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *ServiceListContainersResponse) (ServiceListContainersResponse, error) {
			var marker *string
			if page != nil {
				if page.NextMarker != nil {
					marker = page.NextMarker
				}
			} else {
				// If provided by the user, then use the one from options bag
				marker = listOptions.Marker
			}

			req, err := s.client.listContainersSegmentCreateRequest(ctx, &listOptions)
			if err != nil {
				return ServiceListContainersResponse{}, err
			}

			if marker != nil {
				queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
				if err != nil {
					return ServiceListContainersResponse{}, err
				}
				queryValues.Set("marker", *marker)
				req.Raw().URL.RawQuery = queryValues.Encode()
			}

			resp, err := s.client.pl.Do(req)
			if err != nil {
				return ServiceListContainersResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ServiceListContainersResponse{}, runtime.NewResponseError(resp)
			}

			generatedResp, err := s.client.listContainersSegmentHandleResponse(resp)
			return toServiceListContainersResponse(generatedResp), err
		},
	})
}

// GetProperties - gets the properties of a storage account's Blob service, including properties for Storage Analytics
// and CORS (Cross-Origin Resource Sharing) rules.
func (s *ServiceClient) GetProperties(ctx context.Context, o *ServiceGetPropertiesOptions) (ServiceGetPropertiesResponse, error) {
	getPropertiesOptions := o.format()
	resp, err := s.client.GetProperties(ctx, getPropertiesOptions)

	return toServiceGetPropertiesResponse(resp), err
}

// SetProperties Sets the properties of a storage account's Blob service, including Azure Storage Analytics.
// If an element (e.g. analytics_logging) is left as None, the existing settings on the service for that functionality are preserved.
func (s *ServiceClient) SetProperties(ctx context.Context, o *ServiceSetPropertiesOptions) (ServiceSetPropertiesResponse, error) {
	properties, setPropertiesOptions := o.format()
	resp, err := s.client.SetProperties(ctx, properties, setPropertiesOptions)

	return toServiceSetPropertiesResponse(resp), err
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
func (s *ServiceClient) GetStatistics(ctx context.Context, o *ServiceGetStatisticsOptions) (ServiceGetStatisticsResponse, error) {
	getStatisticsOptions := o.format()
	resp, err := s.client.GetStatistics(ctx, getStatisticsOptions)

	return toServiceGetStatisticsResponse(resp), err
}

// CanGetAccountSASToken checks if shared key in ServiceClient is nil
func (s *ServiceClient) CanGetAccountSASToken() bool {
	return s.sharedKey != nil
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
// This validity can be checked with CanGetAccountSASToken().
func (s *ServiceClient) GetSASURL(resources AccountSASResourceTypes, permissions AccountSASPermissions, start time.Time, expiry time.Time) (string, error) {
	if s.sharedKey == nil {
		return "", errors.New("SAS can only be signed with a SharedKeyCredential")
	}

	qps, err := AccountSASSignatureValues{
		Version:       internal.SASVersion,
		Protocol:      SASProtocolHTTPS,
		Permissions:   permissions.String(),
		Services:      "b",
		ResourceTypes: resources.String(),
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
	}.Sign(s.sharedKey)
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

// FindBlobsByTags operation finds all blobs in the storage account whose tags match a given search expression.
// Filter blobs searches across all containers within a storage account but can be scoped within the expression to a single container.
// https://docs.microsoft.com/en-us/rest/api/storageservices/find-blobs-by-tags
// eg. "dog='germanshepherd' and penguin='emperorpenguin'"
// To specify a container, eg. "@container=’containerName’ and Name = ‘C’"
func (s *ServiceClient) FindBlobsByTags(ctx context.Context, o *ServiceFilterBlobsOptions) (ServiceFilterBlobsResponse, error) {
	// TODO: Use pager here? Missing support from zz_generated_pagers.go
	serviceFilterBlobsOptions := o.pointer()
	resp, err := s.client.FilterBlobs(ctx, serviceFilterBlobsOptions)
	return toServiceFilterBlobsResponse(resp), err
}

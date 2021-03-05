// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	// ContainerNameRoot is the special Azure Storage name used to identify a storage account's root container.
	ContainerNameRoot = "$root"

	// ContainerNameLogs is the special Azure Storage name used to identify a storage account's logs container.
	ContainerNameLogs = "$logs"
)

// A ServiceClient represents a URL to the Azure Storage Blob service allowing you to manipulate blob containers.
type ServiceClient struct {
	client *serviceClient
	u      url.URL
	cred   StorageAccountCredential
}

// NewServiceClient creates a ServiceClient object using the specified URL, credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewServiceClient(serviceURL string, cred azcore.Credential, options *ClientOptions) (ServiceClient, error) {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return ServiceClient{}, err
	}

	c, _ := cred.(*SharedKeyCredential)

	return ServiceClient{client: &serviceClient{
		con: newConnection(serviceURL, cred, options.getConnectionOptions()),
	}, u: *u, cred: c}, nil
}

// URL returns the URL endpoint used by the ServiceClient object.
func (s ServiceClient) URL() string {
	return s.client.con.u
}

// NewContainerClient creates a new ContainerClient object by concatenating containerName to the end of
// ServiceClient's URL. The new ContainerClient uses the same request policy pipeline as the ServiceClient.
// To change the pipeline, create the ContainerClient and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewContainerClient instead of calling this object's
// NewContainerClient method.
func (s ServiceClient) NewContainerClient(containerName string) ContainerClient {
	containerURL := appendToURLPath(s.client.con.u, containerName)
	containerConnection := newConnectionWithPipeline(containerURL, s.client.con.p)
	return ContainerClient{
		client: &containerClient{
			con: containerConnection,
		},
	}
}

// appendToURLPath appends a string to the end of a URL's path (prefixing the string with a '/' if required)
func appendToURLPath(u string, name string) string {
	// e.g. "https://ms.com/a/b/?k1=v1&k2=v2#f"
	// When you call url.Parse() this is what you'll get:
	//     Scheme: "https"
	//     Opaque: ""
	//       User: nil
	//       Host: "ms.com"
	//       Path: "/a/b/"	This should start with a / and it might or might not have a trailing slash
	//    RawPath: ""
	// ForceQuery: false
	//   RawQuery: "k1=v1&k2=v2"
	//   Fragment: "f"
	uri, _ := url.Parse(u)

	if len(uri.Path) == 0 || uri.Path[len(uri.Path)-1] != '/' {
		uri.Path += "/" // Append "/" to end before appending name
	}
	uri.Path += name
	return uri.String()
}

func (s ServiceClient) GetAccountInfo(ctx context.Context) (ServiceGetAccountInfoResponse, error) {
	resp, err := s.client.GetAccountInfo(ctx, nil)

	return resp, handleError(err)
}

// GetUserDelegationCredential obtains a UserDelegationKey object using the base ServiceClient object.
// OAuth is required for this call, as well as any role that can delegate access to the storage account.
// Strings in KeyInfo should be formatted with SASTimeFormat.
func (s ServiceClient) GetUserDelegationCredential(ctx context.Context, info KeyInfo) (UserDelegationCredential, error) {
	udk, err := s.client.GetUserDelegationKey(ctx, info, nil)
	if err != nil {
		return UserDelegationCredential{}, handleError(err)
	}
	urlParts := NewBlobURLParts(s.URL())
	return NewUserDelegationCredential(strings.Split(urlParts.Host, ".")[0], *udk.UserDelegationKey), nil
}

// The List Containers Segment operation returns a pager of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (s ServiceClient) ListContainersSegment(o *ListContainersSegmentOptions) ListContainersSegmentResponsePager {
	listOptions := o.pointers()
	pager := s.client.ListContainersSegment(listOptions)
	// override the generated advancer, which is incorrect
	if pager.Err() != nil {
		return pager
	}

	p := pager.(*listContainersSegmentResponsePager) // cast to the internal type first
	p.advancer = func(cxt context.Context, response ListContainersSegmentResponseResponse) (*azcore.Request, error) {
		if response.EnumerationResults.NextMarker == nil {
			return nil, errors.New("unexpected missing NextMarker")
		}
		req, err := s.client.listContainersSegmentCreateRequest(cxt, listOptions)
		if err != nil {
			return nil, err
		}
		queryValues, _ := url.ParseQuery(req.URL.RawQuery)
		queryValues.Set("marker", *response.EnumerationResults.NextMarker)

		req.URL.RawQuery = queryValues.Encode()
		return req, nil
	}

	return p
}

func (s ServiceClient) GetProperties(ctx context.Context) (StorageServicePropertiesResponse, error) {
	resp, err := s.client.GetProperties(ctx, nil)

	return resp, handleError(err)
}

func (s ServiceClient) SetProperties(ctx context.Context, properties StorageServiceProperties) (ServiceSetPropertiesResponse, error) {
	resp, err := s.client.SetProperties(ctx, properties, nil)

	return resp, handleError(err)
}

func (s ServiceClient) GetStatistics(ctx context.Context) (StorageServiceStatsResponse, error) {
	resp, err := s.client.GetStatistics(ctx, nil)

	return resp, handleError(err)
}

func (s ServiceClient) CanGetAccountSASToken() bool {
	return s.cred != nil
}

// GetAccountSASToken is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the supplied azcore.Credential during creation was a SharedKeyCredential.
// This validity can be checked with CanGetAccountSASToken().
func (s ServiceClient) GetAccountSASToken(services AccountSASServices, resources AccountSASResourceTypes, permissions AccountSASPermissions, validityTime time.Duration) (SASQueryParameters, error) {
	return AccountSASSignatureValues{
		Version: SASVersion,

		Permissions: permissions.String(),
		Services: services.String(),
		ResourceTypes: resources.String(),

		StartTime: time.Now(),
		ExpiryTime: time.Now().Add(validityTime),
	}.NewSASQueryParameters(s.cred.(*SharedKeyCredential))
}

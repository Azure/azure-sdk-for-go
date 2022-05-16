//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azfile

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// A ServiceClient represents a URL to the Azure Storage File service allowing you to manipulate file share.
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
func NewServiceClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*ServiceClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(serviceURL, conOptions)

	return &ServiceClient{
		client: newServiceClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewServiceClientWithNoCredential creates a ServiceClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewServiceClientWithNoCredential(serviceURL string, options *ClientOptions) (*ServiceClient, error) {
	conOptions := getConnectionOptions(options)
	conn := newConnection(serviceURL, conOptions)

	return &ServiceClient{
		client: newServiceClient(conn.Endpoint(), conn.Pipeline()),
	}, nil
}

// NewServiceClientWithSharedKey creates a ServiceClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewServiceClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*ServiceClient, error) {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(serviceURL, conOptions)

	return &ServiceClient{
		client:    newServiceClient(conn.Endpoint(), conn.Pipeline()),
		sharedKey: cred,
	}, nil
}

// NewServiceClientFromConnectionString creates a service client from the given connection string.
//nolint
func NewServiceClientFromConnectionString(connectionString string, options *ClientOptions) (*ServiceClient, error) {
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewServiceClientWithSharedKey(endpoint, credential, options)
}

// NewShareClient creates a new ShareURL object by concatenating shareName to the end of ServiceURLs' URL.
// The new ShareURL uses the same request policy pipeline as the ServiceURL.
// To change the pipeline, create the ShareURL and then call its WithPipeline method passing in the desired pipeline object.
// Or, call this package's NewShareURL instead of calling this object's NewShareURL method.
func (s *ServiceClient) NewShareClient(shareName string) (*ShareClient, error) {
	shareURL := appendToURLPath(s.client.endpoint, shareName)
	return &ShareClient{
		client:    newShareClient(shareURL, s.client.pl),
		sharedKey: s.sharedKey,
	}, nil
}

// GetProperties returns the properties of the File service.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/get-file-service-properties.
func (s *ServiceClient) GetProperties(ctx context.Context, o *ServiceGetPropertiesOptions) (ServiceGetPropertiesResponse, error) {
	options := o.format()
	getPropertiesResponse, err := s.client.GetProperties(ctx, options)

	return toServiceGetPropertiesResponse(getPropertiesResponse), handleError(err)
}

// SetProperties sets the properties of the File service.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/set-file-service-properties.
func (s *ServiceClient) SetProperties(ctx context.Context, o *ServiceSetPropertiesOptions) (ServiceSetPropertiesResponse, error) {
	properties, serviceSetPropertiesOptions := o.format()
	setPropertiesResponse, err := s.client.SetProperties(ctx, properties, serviceSetPropertiesOptions)

	return toServiceSetPropertiesResponse(setPropertiesResponse), handleError(err)
}

// ListShares operation returns a pager of the containers under the specified account.
// Use an empty Marker to start enumeration from the beginning. Container names are returned in lexicographic order.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/list-containers2.
func (s *ServiceClient) ListShares(o *ServiceListSharesOptions) *runtime.Pager[ServiceListSharesResponse] {
	listOptions := o.format()
	return runtime.NewPager(runtime.PagingHandler[ServiceListSharesResponse]{
		More: func(page ServiceListSharesResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *ServiceListSharesResponse) (ServiceListSharesResponse, error) {
			if page.ListSharesResponse.NextMarker == nil || len(*page.NextMarker) == 0 {
				return ServiceListSharesResponse{}, handleError(errors.New("unexpected missing NextMarker"))
			}
			req, err := s.client.listSharesSegmentCreateRequest(ctx, listOptions)
			if err != nil {
				return ServiceListSharesResponse{}, handleError(err)
			}
			queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
			if err != nil {
				return ServiceListSharesResponse{}, handleError(err)
			}
			queryValues.Set("marker", *page.ListSharesResponse.NextMarker)
			req.Raw().URL.RawQuery = queryValues.Encode()

			resp, err := s.client.pl.Do(req)
			if err != nil {
				return ServiceListSharesResponse{}, handleError(err)
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ServiceListSharesResponse{}, handleError(runtime.NewResponseError(resp))
			}

			generatedResp, err := s.client.listSharesSegmentHandleResponse(resp)
			return toServiceListSharesResponse(generatedResp), handleError(err)
		},
	})
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
		Version:       SASVersion,
		Protocol:      SASProtocolHTTPS,
		Permissions:   permissions.String(),
		Services:      "f",
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

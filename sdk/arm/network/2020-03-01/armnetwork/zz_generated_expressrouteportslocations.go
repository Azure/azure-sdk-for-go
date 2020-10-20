// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// ExpressRoutePortsLocationsOperations contains the methods for the ExpressRoutePortsLocations group.
type ExpressRoutePortsLocationsOperations interface {
	// Get - Retrieves a single ExpressRoutePort peering location, including the list of available bandwidths available at said peering location.
	Get(ctx context.Context, locationName string, options *ExpressRoutePortsLocationsGetOptions) (*ExpressRoutePortsLocationResponse, error)
	// List - Retrieves all ExpressRoutePort peering locations. Does not return available bandwidths for each location. Available bandwidths can only be obtained when retrieving a specific peering location.
	List(options *ExpressRoutePortsLocationsListOptions) ExpressRoutePortsLocationListResultPager
}

// ExpressRoutePortsLocationsClient implements the ExpressRoutePortsLocationsOperations interface.
// Don't use this type directly, use NewExpressRoutePortsLocationsClient() instead.
type ExpressRoutePortsLocationsClient struct {
	*Client
	subscriptionID string
}

// NewExpressRoutePortsLocationsClient creates a new instance of ExpressRoutePortsLocationsClient with the specified values.
func NewExpressRoutePortsLocationsClient(c *Client, subscriptionID string) ExpressRoutePortsLocationsOperations {
	return &ExpressRoutePortsLocationsClient{Client: c, subscriptionID: subscriptionID}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *ExpressRoutePortsLocationsClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

// Get - Retrieves a single ExpressRoutePort peering location, including the list of available bandwidths available at said peering location.
func (client *ExpressRoutePortsLocationsClient) Get(ctx context.Context, locationName string, options *ExpressRoutePortsLocationsGetOptions) (*ExpressRoutePortsLocationResponse, error) {
	req, err := client.GetCreateRequest(ctx, locationName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetHandleError(resp)
	}
	result, err := client.GetHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCreateRequest creates the Get request.
func (client *ExpressRoutePortsLocationsClient) GetCreateRequest(ctx context.Context, locationName string, options *ExpressRoutePortsLocationsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/ExpressRoutePortsLocations/{locationName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetHandleResponse handles the Get response.
func (client *ExpressRoutePortsLocationsClient) GetHandleResponse(resp *azcore.Response) (*ExpressRoutePortsLocationResponse, error) {
	result := ExpressRoutePortsLocationResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ExpressRoutePortsLocation)
}

// GetHandleError handles the Get error response.
func (client *ExpressRoutePortsLocationsClient) GetHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// List - Retrieves all ExpressRoutePort peering locations. Does not return available bandwidths for each location. Available bandwidths can only be obtained when retrieving a specific peering location.
func (client *ExpressRoutePortsLocationsClient) List(options *ExpressRoutePortsLocationsListOptions) ExpressRoutePortsLocationListResultPager {
	return &expressRoutePortsLocationListResultPager{
		pipeline: client.p,
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.ListCreateRequest(ctx, options)
		},
		responder: client.ListHandleResponse,
		errorer:   client.ListHandleError,
		advancer: func(ctx context.Context, resp *ExpressRoutePortsLocationListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ExpressRoutePortsLocationListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// ListCreateRequest creates the List request.
func (client *ExpressRoutePortsLocationsClient) ListCreateRequest(ctx context.Context, options *ExpressRoutePortsLocationsListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/ExpressRoutePortsLocations"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListHandleResponse handles the List response.
func (client *ExpressRoutePortsLocationsClient) ListHandleResponse(resp *azcore.Response) (*ExpressRoutePortsLocationListResultResponse, error) {
	result := ExpressRoutePortsLocationListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ExpressRoutePortsLocationListResult)
}

// ListHandleError handles the List error response.
func (client *ExpressRoutePortsLocationsClient) ListHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// AvailablePrivateEndpointTypesClient contains the methods for the AvailablePrivateEndpointTypes group.
// Don't use this type directly, use NewAvailablePrivateEndpointTypesClient() instead.
type AvailablePrivateEndpointTypesClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewAvailablePrivateEndpointTypesClient creates a new instance of AvailablePrivateEndpointTypesClient with the specified values.
func NewAvailablePrivateEndpointTypesClient(con *armcore.Connection, subscriptionID string) *AvailablePrivateEndpointTypesClient {
	return &AvailablePrivateEndpointTypesClient{con: con, subscriptionID: subscriptionID}
}

// List - Returns all of the resource types that can be linked to a Private Endpoint in this subscription in this region.
// If the operation fails it returns the *CloudError error type.
func (client *AvailablePrivateEndpointTypesClient) List(location string, options *AvailablePrivateEndpointTypesListOptions) AvailablePrivateEndpointTypesResultPager {
	return &availablePrivateEndpointTypesResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, location, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp AvailablePrivateEndpointTypesResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.AvailablePrivateEndpointTypesResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client *AvailablePrivateEndpointTypesClient) listCreateRequest(ctx context.Context, location string, options *AvailablePrivateEndpointTypesListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/locations/{location}/availablePrivateEndpointTypes"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-02-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AvailablePrivateEndpointTypesClient) listHandleResponse(resp *azcore.Response) (AvailablePrivateEndpointTypesResultResponse, error) {
	var val *AvailablePrivateEndpointTypesResult
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return AvailablePrivateEndpointTypesResultResponse{}, err
	}
	return AvailablePrivateEndpointTypesResultResponse{RawResponse: resp.Response, AvailablePrivateEndpointTypesResult: val}, nil
}

// listHandleError handles the List error response.
func (client *AvailablePrivateEndpointTypesClient) listHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// ListByResourceGroup - Returns all of the resource types that can be linked to a Private Endpoint in this subscription in this region.
// If the operation fails it returns the *CloudError error type.
func (client *AvailablePrivateEndpointTypesClient) ListByResourceGroup(location string, resourceGroupName string, options *AvailablePrivateEndpointTypesListByResourceGroupOptions) AvailablePrivateEndpointTypesResultPager {
	return &availablePrivateEndpointTypesResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listByResourceGroupCreateRequest(ctx, location, resourceGroupName, options)
		},
		responder: client.listByResourceGroupHandleResponse,
		errorer:   client.listByResourceGroupHandleError,
		advancer: func(ctx context.Context, resp AvailablePrivateEndpointTypesResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.AvailablePrivateEndpointTypesResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *AvailablePrivateEndpointTypesClient) listByResourceGroupCreateRequest(ctx context.Context, location string, resourceGroupName string, options *AvailablePrivateEndpointTypesListByResourceGroupOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/locations/{location}/availablePrivateEndpointTypes"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-02-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *AvailablePrivateEndpointTypesClient) listByResourceGroupHandleResponse(resp *azcore.Response) (AvailablePrivateEndpointTypesResultResponse, error) {
	var val *AvailablePrivateEndpointTypesResult
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return AvailablePrivateEndpointTypesResultResponse{}, err
	}
	return AvailablePrivateEndpointTypesResultResponse{RawResponse: resp.Response, AvailablePrivateEndpointTypesResult: val}, nil
}

// listByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client *AvailablePrivateEndpointTypesClient) listByResourceGroupHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

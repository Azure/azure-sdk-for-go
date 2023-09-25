//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// TopologyClient contains the methods for the Topology group.
// Don't use this type directly, use NewTopologyClient() instead.
type TopologyClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewTopologyClient creates a new instance of TopologyClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewTopologyClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TopologyClient, error) {
	cl, err := arm.NewClient(moduleName+".TopologyClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &TopologyClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Gets a specific topology component.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-01-01
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - ascLocation - The location where ASC stores the data of the subscription. can be retrieved from Get locations
//   - topologyResourceName - Name of a topology resources collection.
//   - options - TopologyClientGetOptions contains the optional parameters for the TopologyClient.Get method.
func (client *TopologyClient) Get(ctx context.Context, resourceGroupName string, ascLocation string, topologyResourceName string, options *TopologyClientGetOptions) (TopologyClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, ascLocation, topologyResourceName, options)
	if err != nil {
		return TopologyClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TopologyClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TopologyClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *TopologyClient) getCreateRequest(ctx context.Context, resourceGroupName string, ascLocation string, topologyResourceName string, options *TopologyClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/locations/{ascLocation}/topologies/{topologyResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ascLocation == "" {
		return nil, errors.New("parameter ascLocation cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ascLocation}", url.PathEscape(ascLocation))
	if topologyResourceName == "" {
		return nil, errors.New("parameter topologyResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topologyResourceName}", url.PathEscape(topologyResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *TopologyClient) getHandleResponse(resp *http.Response) (TopologyClientGetResponse, error) {
	result := TopologyClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TopologyResource); err != nil {
		return TopologyClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list that allows to build a topology view of a subscription.
//
// Generated from API version 2020-01-01
//   - options - TopologyClientListOptions contains the optional parameters for the TopologyClient.NewListPager method.
func (client *TopologyClient) NewListPager(options *TopologyClientListOptions) (*runtime.Pager[TopologyClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[TopologyClientListResponse]{
		More: func(page TopologyClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TopologyClientListResponse) (TopologyClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return TopologyClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return TopologyClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return TopologyClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *TopologyClient) listCreateRequest(ctx context.Context, options *TopologyClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/topologies"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *TopologyClient) listHandleResponse(resp *http.Response) (TopologyClientListResponse, error) {
	result := TopologyClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TopologyList); err != nil {
		return TopologyClientListResponse{}, err
	}
	return result, nil
}

// NewListByHomeRegionPager - Gets a list that allows to build a topology view of a subscription and location.
//
// Generated from API version 2020-01-01
//   - ascLocation - The location where ASC stores the data of the subscription. can be retrieved from Get locations
//   - options - TopologyClientListByHomeRegionOptions contains the optional parameters for the TopologyClient.NewListByHomeRegionPager
//     method.
func (client *TopologyClient) NewListByHomeRegionPager(ascLocation string, options *TopologyClientListByHomeRegionOptions) (*runtime.Pager[TopologyClientListByHomeRegionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[TopologyClientListByHomeRegionResponse]{
		More: func(page TopologyClientListByHomeRegionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TopologyClientListByHomeRegionResponse) (TopologyClientListByHomeRegionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByHomeRegionCreateRequest(ctx, ascLocation, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return TopologyClientListByHomeRegionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return TopologyClientListByHomeRegionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return TopologyClientListByHomeRegionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByHomeRegionHandleResponse(resp)
		},
	})
}

// listByHomeRegionCreateRequest creates the ListByHomeRegion request.
func (client *TopologyClient) listByHomeRegionCreateRequest(ctx context.Context, ascLocation string, options *TopologyClientListByHomeRegionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/locations/{ascLocation}/topologies"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if ascLocation == "" {
		return nil, errors.New("parameter ascLocation cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ascLocation}", url.PathEscape(ascLocation))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByHomeRegionHandleResponse handles the ListByHomeRegion response.
func (client *TopologyClient) listByHomeRegionHandleResponse(resp *http.Response) (TopologyClientListByHomeRegionResponse, error) {
	result := TopologyClientListByHomeRegionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TopologyList); err != nil {
		return TopologyClientListByHomeRegionResponse{}, err
	}
	return result, nil
}


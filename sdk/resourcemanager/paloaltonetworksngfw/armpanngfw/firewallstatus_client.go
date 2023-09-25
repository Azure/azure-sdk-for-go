//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpanngfw

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

// FirewallStatusClient contains the methods for the FirewallStatus group.
// Don't use this type directly, use NewFirewallStatusClient() instead.
type FirewallStatusClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewFirewallStatusClient creates a new instance of FirewallStatusClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewFirewallStatusClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FirewallStatusClient, error) {
	cl, err := arm.NewClient(moduleName+".FirewallStatusClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &FirewallStatusClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Get a FirewallStatusResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - firewallName - Firewall resource name
//   - options - FirewallStatusClientGetOptions contains the optional parameters for the FirewallStatusClient.Get method.
func (client *FirewallStatusClient) Get(ctx context.Context, resourceGroupName string, firewallName string, options *FirewallStatusClientGetOptions) (FirewallStatusClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, firewallName, options)
	if err != nil {
		return FirewallStatusClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FirewallStatusClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FirewallStatusClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *FirewallStatusClient) getCreateRequest(ctx context.Context, resourceGroupName string, firewallName string, options *FirewallStatusClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/PaloAltoNetworks.Cloudngfw/firewalls/{firewallName}/statuses/default"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if firewallName == "" {
		return nil, errors.New("parameter firewallName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{firewallName}", url.PathEscape(firewallName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *FirewallStatusClient) getHandleResponse(resp *http.Response) (FirewallStatusClientGetResponse, error) {
	result := FirewallStatusClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FirewallStatusResource); err != nil {
		return FirewallStatusClientGetResponse{}, err
	}
	return result, nil
}

// NewListByFirewallsPager - List FirewallStatusResource resources by Firewalls
//
// Generated from API version 2022-08-29
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - firewallName - Firewall resource name
//   - options - FirewallStatusClientListByFirewallsOptions contains the optional parameters for the FirewallStatusClient.NewListByFirewallsPager
//     method.
func (client *FirewallStatusClient) NewListByFirewallsPager(resourceGroupName string, firewallName string, options *FirewallStatusClientListByFirewallsOptions) (*runtime.Pager[FirewallStatusClientListByFirewallsResponse]) {
	return runtime.NewPager(runtime.PagingHandler[FirewallStatusClientListByFirewallsResponse]{
		More: func(page FirewallStatusClientListByFirewallsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *FirewallStatusClientListByFirewallsResponse) (FirewallStatusClientListByFirewallsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByFirewallsCreateRequest(ctx, resourceGroupName, firewallName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return FirewallStatusClientListByFirewallsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return FirewallStatusClientListByFirewallsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return FirewallStatusClientListByFirewallsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByFirewallsHandleResponse(resp)
		},
	})
}

// listByFirewallsCreateRequest creates the ListByFirewalls request.
func (client *FirewallStatusClient) listByFirewallsCreateRequest(ctx context.Context, resourceGroupName string, firewallName string, options *FirewallStatusClientListByFirewallsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/PaloAltoNetworks.Cloudngfw/firewalls/{firewallName}/statuses"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if firewallName == "" {
		return nil, errors.New("parameter firewallName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{firewallName}", url.PathEscape(firewallName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByFirewallsHandleResponse handles the ListByFirewalls response.
func (client *FirewallStatusClient) listByFirewallsHandleResponse(resp *http.Response) (FirewallStatusClientListByFirewallsResponse, error) {
	result := FirewallStatusClientListByFirewallsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FirewallStatusResourceListResult); err != nil {
		return FirewallStatusClientListByFirewallsResponse{}, err
	}
	return result, nil
}


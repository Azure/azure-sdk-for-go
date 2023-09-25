//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagednetworkfabric

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

// IPPrefixesClient contains the methods for the IPPrefixes group.
// Don't use this type directly, use NewIPPrefixesClient() instead.
type IPPrefixesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewIPPrefixesClient creates a new instance of IPPrefixesClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewIPPrefixesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*IPPrefixesClient, error) {
	cl, err := arm.NewClient(moduleName+".IPPrefixesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &IPPrefixesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreate - Implements IP Prefix PUT method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - ipPrefixName - Name of the IP Prefix.
//   - body - Request payload.
//   - options - IPPrefixesClientBeginCreateOptions contains the optional parameters for the IPPrefixesClient.BeginCreate method.
func (client *IPPrefixesClient) BeginCreate(ctx context.Context, resourceGroupName string, ipPrefixName string, body IPPrefix, options *IPPrefixesClientBeginCreateOptions) (*runtime.Poller[IPPrefixesClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, ipPrefixName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[IPPrefixesClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[IPPrefixesClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Implements IP Prefix PUT method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *IPPrefixesClient) create(ctx context.Context, resourceGroupName string, ipPrefixName string, body IPPrefix, options *IPPrefixesClientBeginCreateOptions) (*http.Response, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceGroupName, ipPrefixName, body, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createCreateRequest creates the Create request.
func (client *IPPrefixesClient) createCreateRequest(ctx context.Context, resourceGroupName string, ipPrefixName string, body IPPrefix, options *IPPrefixesClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/ipPrefixes/{ipPrefixName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ipPrefixName == "" {
		return nil, errors.New("parameter ipPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ipPrefixName}", url.PathEscape(ipPrefixName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Implements IP Prefix DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - ipPrefixName - Name of the IP Prefix.
//   - options - IPPrefixesClientBeginDeleteOptions contains the optional parameters for the IPPrefixesClient.BeginDelete method.
func (client *IPPrefixesClient) BeginDelete(ctx context.Context, resourceGroupName string, ipPrefixName string, options *IPPrefixesClientBeginDeleteOptions) (*runtime.Poller[IPPrefixesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, ipPrefixName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[IPPrefixesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[IPPrefixesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Implements IP Prefix DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *IPPrefixesClient) deleteOperation(ctx context.Context, resourceGroupName string, ipPrefixName string, options *IPPrefixesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, ipPrefixName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *IPPrefixesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, ipPrefixName string, options *IPPrefixesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/ipPrefixes/{ipPrefixName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ipPrefixName == "" {
		return nil, errors.New("parameter ipPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ipPrefixName}", url.PathEscape(ipPrefixName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Implements IP Prefix GET method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - ipPrefixName - Name of the IP Prefix.
//   - options - IPPrefixesClientGetOptions contains the optional parameters for the IPPrefixesClient.Get method.
func (client *IPPrefixesClient) Get(ctx context.Context, resourceGroupName string, ipPrefixName string, options *IPPrefixesClientGetOptions) (IPPrefixesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, ipPrefixName, options)
	if err != nil {
		return IPPrefixesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IPPrefixesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return IPPrefixesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *IPPrefixesClient) getCreateRequest(ctx context.Context, resourceGroupName string, ipPrefixName string, options *IPPrefixesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/ipPrefixes/{ipPrefixName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ipPrefixName == "" {
		return nil, errors.New("parameter ipPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ipPrefixName}", url.PathEscape(ipPrefixName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *IPPrefixesClient) getHandleResponse(resp *http.Response) (IPPrefixesClientGetResponse, error) {
	result := IPPrefixesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IPPrefix); err != nil {
		return IPPrefixesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Implements IpPrefixes list by resource group GET method.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - IPPrefixesClientListByResourceGroupOptions contains the optional parameters for the IPPrefixesClient.NewListByResourceGroupPager
//     method.
func (client *IPPrefixesClient) NewListByResourceGroupPager(resourceGroupName string, options *IPPrefixesClientListByResourceGroupOptions) (*runtime.Pager[IPPrefixesClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[IPPrefixesClientListByResourceGroupResponse]{
		More: func(page IPPrefixesClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *IPPrefixesClientListByResourceGroupResponse) (IPPrefixesClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return IPPrefixesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return IPPrefixesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return IPPrefixesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *IPPrefixesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *IPPrefixesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/ipPrefixes"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *IPPrefixesClient) listByResourceGroupHandleResponse(resp *http.Response) (IPPrefixesClientListByResourceGroupResponse, error) {
	result := IPPrefixesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IPPrefixesListResult); err != nil {
		return IPPrefixesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Implements IpPrefixes list by subscription GET method.
//
// Generated from API version 2023-06-15
//   - options - IPPrefixesClientListBySubscriptionOptions contains the optional parameters for the IPPrefixesClient.NewListBySubscriptionPager
//     method.
func (client *IPPrefixesClient) NewListBySubscriptionPager(options *IPPrefixesClientListBySubscriptionOptions) (*runtime.Pager[IPPrefixesClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[IPPrefixesClientListBySubscriptionResponse]{
		More: func(page IPPrefixesClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *IPPrefixesClientListBySubscriptionResponse) (IPPrefixesClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return IPPrefixesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return IPPrefixesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return IPPrefixesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *IPPrefixesClient) listBySubscriptionCreateRequest(ctx context.Context, options *IPPrefixesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ManagedNetworkFabric/ipPrefixes"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *IPPrefixesClient) listBySubscriptionHandleResponse(resp *http.Response) (IPPrefixesClientListBySubscriptionResponse, error) {
	result := IPPrefixesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IPPrefixesListResult); err != nil {
		return IPPrefixesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdate - API to update certain properties of the IP Prefix resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - ipPrefixName - Name of the IP Prefix.
//   - body - IP Prefix properties to update.
//   - options - IPPrefixesClientBeginUpdateOptions contains the optional parameters for the IPPrefixesClient.BeginUpdate method.
func (client *IPPrefixesClient) BeginUpdate(ctx context.Context, resourceGroupName string, ipPrefixName string, body IPPrefixPatch, options *IPPrefixesClientBeginUpdateOptions) (*runtime.Poller[IPPrefixesClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, ipPrefixName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[IPPrefixesClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[IPPrefixesClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - API to update certain properties of the IP Prefix resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *IPPrefixesClient) update(ctx context.Context, resourceGroupName string, ipPrefixName string, body IPPrefixPatch, options *IPPrefixesClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, ipPrefixName, body, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *IPPrefixesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, ipPrefixName string, body IPPrefixPatch, options *IPPrefixesClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/ipPrefixes/{ipPrefixName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ipPrefixName == "" {
		return nil, errors.New("parameter ipPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ipPrefixName}", url.PathEscape(ipPrefixName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
	return nil, err
}
	return req, nil
}


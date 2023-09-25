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

// NeighborGroupsClient contains the methods for the NeighborGroups group.
// Don't use this type directly, use NewNeighborGroupsClient() instead.
type NeighborGroupsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewNeighborGroupsClient creates a new instance of NeighborGroupsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewNeighborGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*NeighborGroupsClient, error) {
	cl, err := arm.NewClient(moduleName+".NeighborGroupsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &NeighborGroupsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreate - Implements the Neighbor Group PUT method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - neighborGroupName - Name of the Neighbor Group.
//   - body - Request payload.
//   - options - NeighborGroupsClientBeginCreateOptions contains the optional parameters for the NeighborGroupsClient.BeginCreate
//     method.
func (client *NeighborGroupsClient) BeginCreate(ctx context.Context, resourceGroupName string, neighborGroupName string, body NeighborGroup, options *NeighborGroupsClientBeginCreateOptions) (*runtime.Poller[NeighborGroupsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, neighborGroupName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NeighborGroupsClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[NeighborGroupsClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Implements the Neighbor Group PUT method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *NeighborGroupsClient) create(ctx context.Context, resourceGroupName string, neighborGroupName string, body NeighborGroup, options *NeighborGroupsClientBeginCreateOptions) (*http.Response, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceGroupName, neighborGroupName, body, options)
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
func (client *NeighborGroupsClient) createCreateRequest(ctx context.Context, resourceGroupName string, neighborGroupName string, body NeighborGroup, options *NeighborGroupsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/neighborGroups/{neighborGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if neighborGroupName == "" {
		return nil, errors.New("parameter neighborGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{neighborGroupName}", url.PathEscape(neighborGroupName))
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

// BeginDelete - Implements Neighbor Group DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - neighborGroupName - Name of the Neighbor Group.
//   - options - NeighborGroupsClientBeginDeleteOptions contains the optional parameters for the NeighborGroupsClient.BeginDelete
//     method.
func (client *NeighborGroupsClient) BeginDelete(ctx context.Context, resourceGroupName string, neighborGroupName string, options *NeighborGroupsClientBeginDeleteOptions) (*runtime.Poller[NeighborGroupsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, neighborGroupName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NeighborGroupsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[NeighborGroupsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Implements Neighbor Group DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *NeighborGroupsClient) deleteOperation(ctx context.Context, resourceGroupName string, neighborGroupName string, options *NeighborGroupsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, neighborGroupName, options)
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
func (client *NeighborGroupsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, neighborGroupName string, options *NeighborGroupsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/neighborGroups/{neighborGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if neighborGroupName == "" {
		return nil, errors.New("parameter neighborGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{neighborGroupName}", url.PathEscape(neighborGroupName))
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

// Get - Gets the Neighbor Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - neighborGroupName - Name of the Neighbor Group.
//   - options - NeighborGroupsClientGetOptions contains the optional parameters for the NeighborGroupsClient.Get method.
func (client *NeighborGroupsClient) Get(ctx context.Context, resourceGroupName string, neighborGroupName string, options *NeighborGroupsClientGetOptions) (NeighborGroupsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, neighborGroupName, options)
	if err != nil {
		return NeighborGroupsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return NeighborGroupsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return NeighborGroupsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *NeighborGroupsClient) getCreateRequest(ctx context.Context, resourceGroupName string, neighborGroupName string, options *NeighborGroupsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/neighborGroups/{neighborGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if neighborGroupName == "" {
		return nil, errors.New("parameter neighborGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{neighborGroupName}", url.PathEscape(neighborGroupName))
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
func (client *NeighborGroupsClient) getHandleResponse(resp *http.Response) (NeighborGroupsClientGetResponse, error) {
	result := NeighborGroupsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NeighborGroup); err != nil {
		return NeighborGroupsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Displays NeighborGroups list by resource group GET method.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - NeighborGroupsClientListByResourceGroupOptions contains the optional parameters for the NeighborGroupsClient.NewListByResourceGroupPager
//     method.
func (client *NeighborGroupsClient) NewListByResourceGroupPager(resourceGroupName string, options *NeighborGroupsClientListByResourceGroupOptions) (*runtime.Pager[NeighborGroupsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[NeighborGroupsClientListByResourceGroupResponse]{
		More: func(page NeighborGroupsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NeighborGroupsClientListByResourceGroupResponse) (NeighborGroupsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NeighborGroupsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return NeighborGroupsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NeighborGroupsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *NeighborGroupsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *NeighborGroupsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/neighborGroups"
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
func (client *NeighborGroupsClient) listByResourceGroupHandleResponse(resp *http.Response) (NeighborGroupsClientListByResourceGroupResponse, error) {
	result := NeighborGroupsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NeighborGroupsListResult); err != nil {
		return NeighborGroupsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Displays NeighborGroups list by subscription GET method.
//
// Generated from API version 2023-06-15
//   - options - NeighborGroupsClientListBySubscriptionOptions contains the optional parameters for the NeighborGroupsClient.NewListBySubscriptionPager
//     method.
func (client *NeighborGroupsClient) NewListBySubscriptionPager(options *NeighborGroupsClientListBySubscriptionOptions) (*runtime.Pager[NeighborGroupsClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[NeighborGroupsClientListBySubscriptionResponse]{
		More: func(page NeighborGroupsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NeighborGroupsClientListBySubscriptionResponse) (NeighborGroupsClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NeighborGroupsClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return NeighborGroupsClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NeighborGroupsClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *NeighborGroupsClient) listBySubscriptionCreateRequest(ctx context.Context, options *NeighborGroupsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ManagedNetworkFabric/neighborGroups"
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
func (client *NeighborGroupsClient) listBySubscriptionHandleResponse(resp *http.Response) (NeighborGroupsClientListBySubscriptionResponse, error) {
	result := NeighborGroupsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NeighborGroupsListResult); err != nil {
		return NeighborGroupsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Updates the Neighbor Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - neighborGroupName - Name of the Neighbor Group.
//   - body - Neighbor Group properties to update. Only annotations are supported.
//   - options - NeighborGroupsClientBeginUpdateOptions contains the optional parameters for the NeighborGroupsClient.BeginUpdate
//     method.
func (client *NeighborGroupsClient) BeginUpdate(ctx context.Context, resourceGroupName string, neighborGroupName string, body NeighborGroupPatch, options *NeighborGroupsClientBeginUpdateOptions) (*runtime.Poller[NeighborGroupsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, neighborGroupName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NeighborGroupsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[NeighborGroupsClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Updates the Neighbor Group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *NeighborGroupsClient) update(ctx context.Context, resourceGroupName string, neighborGroupName string, body NeighborGroupPatch, options *NeighborGroupsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, neighborGroupName, body, options)
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
func (client *NeighborGroupsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, neighborGroupName string, body NeighborGroupPatch, options *NeighborGroupsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/neighborGroups/{neighborGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if neighborGroupName == "" {
		return nil, errors.New("parameter neighborGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{neighborGroupName}", url.PathEscape(neighborGroupName))
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


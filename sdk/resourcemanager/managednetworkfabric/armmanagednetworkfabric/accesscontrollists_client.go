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

// AccessControlListsClient contains the methods for the AccessControlLists group.
// Don't use this type directly, use NewAccessControlListsClient() instead.
type AccessControlListsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewAccessControlListsClient creates a new instance of AccessControlListsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAccessControlListsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AccessControlListsClient, error) {
	cl, err := arm.NewClient(moduleName+".AccessControlListsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AccessControlListsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreate - Implements Access Control List PUT method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - body - Request payload.
//   - options - AccessControlListsClientBeginCreateOptions contains the optional parameters for the AccessControlListsClient.BeginCreate
//     method.
func (client *AccessControlListsClient) BeginCreate(ctx context.Context, resourceGroupName string, accessControlListName string, body AccessControlList, options *AccessControlListsClientBeginCreateOptions) (*runtime.Poller[AccessControlListsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, accessControlListName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AccessControlListsClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AccessControlListsClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Implements Access Control List PUT method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *AccessControlListsClient) create(ctx context.Context, resourceGroupName string, accessControlListName string, body AccessControlList, options *AccessControlListsClientBeginCreateOptions) (*http.Response, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceGroupName, accessControlListName, body, options)
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
func (client *AccessControlListsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, body AccessControlList, options *AccessControlListsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
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

// BeginDelete - Implements Access Control List DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - options - AccessControlListsClientBeginDeleteOptions contains the optional parameters for the AccessControlListsClient.BeginDelete
//     method.
func (client *AccessControlListsClient) BeginDelete(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginDeleteOptions) (*runtime.Poller[AccessControlListsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, accessControlListName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AccessControlListsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AccessControlListsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Implements Access Control List DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *AccessControlListsClient) deleteOperation(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accessControlListName, options)
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
func (client *AccessControlListsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
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

// Get - Implements Access Control List GET method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - options - AccessControlListsClientGetOptions contains the optional parameters for the AccessControlListsClient.Get method.
func (client *AccessControlListsClient) Get(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientGetOptions) (AccessControlListsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, accessControlListName, options)
	if err != nil {
		return AccessControlListsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AccessControlListsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AccessControlListsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AccessControlListsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
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
func (client *AccessControlListsClient) getHandleResponse(resp *http.Response) (AccessControlListsClientGetResponse, error) {
	result := AccessControlListsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessControlList); err != nil {
		return AccessControlListsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Implements AccessControlLists list by resource group GET method.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - AccessControlListsClientListByResourceGroupOptions contains the optional parameters for the AccessControlListsClient.NewListByResourceGroupPager
//     method.
func (client *AccessControlListsClient) NewListByResourceGroupPager(resourceGroupName string, options *AccessControlListsClientListByResourceGroupOptions) (*runtime.Pager[AccessControlListsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AccessControlListsClientListByResourceGroupResponse]{
		More: func(page AccessControlListsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AccessControlListsClientListByResourceGroupResponse) (AccessControlListsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AccessControlListsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AccessControlListsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AccessControlListsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *AccessControlListsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *AccessControlListsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists"
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
func (client *AccessControlListsClient) listByResourceGroupHandleResponse(resp *http.Response) (AccessControlListsClientListByResourceGroupResponse, error) {
	result := AccessControlListsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessControlListsListResult); err != nil {
		return AccessControlListsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Implements AccessControlLists list by subscription GET method.
//
// Generated from API version 2023-06-15
//   - options - AccessControlListsClientListBySubscriptionOptions contains the optional parameters for the AccessControlListsClient.NewListBySubscriptionPager
//     method.
func (client *AccessControlListsClient) NewListBySubscriptionPager(options *AccessControlListsClientListBySubscriptionOptions) (*runtime.Pager[AccessControlListsClientListBySubscriptionResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AccessControlListsClientListBySubscriptionResponse]{
		More: func(page AccessControlListsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AccessControlListsClientListBySubscriptionResponse) (AccessControlListsClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AccessControlListsClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AccessControlListsClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AccessControlListsClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *AccessControlListsClient) listBySubscriptionCreateRequest(ctx context.Context, options *AccessControlListsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ManagedNetworkFabric/accessControlLists"
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
func (client *AccessControlListsClient) listBySubscriptionHandleResponse(resp *http.Response) (AccessControlListsClientListBySubscriptionResponse, error) {
	result := AccessControlListsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessControlListsListResult); err != nil {
		return AccessControlListsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginResync - Implements the operation to the underlying resources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - options - AccessControlListsClientBeginResyncOptions contains the optional parameters for the AccessControlListsClient.BeginResync
//     method.
func (client *AccessControlListsClient) BeginResync(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginResyncOptions) (*runtime.Poller[AccessControlListsClientResyncResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.resync(ctx, resourceGroupName, accessControlListName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AccessControlListsClientResyncResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AccessControlListsClientResyncResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Resync - Implements the operation to the underlying resources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *AccessControlListsClient) resync(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginResyncOptions) (*http.Response, error) {
	var err error
	req, err := client.resyncCreateRequest(ctx, resourceGroupName, accessControlListName, options)
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

// resyncCreateRequest creates the Resync request.
func (client *AccessControlListsClient) resyncCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginResyncOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}/resync"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginUpdate - API to update certain properties of the Access Control List resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - body - Access Control List properties to update.
//   - options - AccessControlListsClientBeginUpdateOptions contains the optional parameters for the AccessControlListsClient.BeginUpdate
//     method.
func (client *AccessControlListsClient) BeginUpdate(ctx context.Context, resourceGroupName string, accessControlListName string, body AccessControlListPatch, options *AccessControlListsClientBeginUpdateOptions) (*runtime.Poller[AccessControlListsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, accessControlListName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AccessControlListsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AccessControlListsClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - API to update certain properties of the Access Control List resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *AccessControlListsClient) update(ctx context.Context, resourceGroupName string, accessControlListName string, body AccessControlListPatch, options *AccessControlListsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, accessControlListName, body, options)
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
func (client *AccessControlListsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, body AccessControlListPatch, options *AccessControlListsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
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

// BeginUpdateAdministrativeState - Implements the operation to the underlying resources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - body - Request payload.
//   - options - AccessControlListsClientBeginUpdateAdministrativeStateOptions contains the optional parameters for the AccessControlListsClient.BeginUpdateAdministrativeState
//     method.
func (client *AccessControlListsClient) BeginUpdateAdministrativeState(ctx context.Context, resourceGroupName string, accessControlListName string, body UpdateAdministrativeState, options *AccessControlListsClientBeginUpdateAdministrativeStateOptions) (*runtime.Poller[AccessControlListsClientUpdateAdministrativeStateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.updateAdministrativeState(ctx, resourceGroupName, accessControlListName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AccessControlListsClientUpdateAdministrativeStateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AccessControlListsClientUpdateAdministrativeStateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// UpdateAdministrativeState - Implements the operation to the underlying resources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *AccessControlListsClient) updateAdministrativeState(ctx context.Context, resourceGroupName string, accessControlListName string, body UpdateAdministrativeState, options *AccessControlListsClientBeginUpdateAdministrativeStateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateAdministrativeStateCreateRequest(ctx, resourceGroupName, accessControlListName, body, options)
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

// updateAdministrativeStateCreateRequest creates the UpdateAdministrativeState request.
func (client *AccessControlListsClient) updateAdministrativeStateCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, body UpdateAdministrativeState, options *AccessControlListsClientBeginUpdateAdministrativeStateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}/updateAdministrativeState"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
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

// BeginValidateConfiguration - Implements the operation to the underlying resources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accessControlListName - Name of the Access Control List.
//   - options - AccessControlListsClientBeginValidateConfigurationOptions contains the optional parameters for the AccessControlListsClient.BeginValidateConfiguration
//     method.
func (client *AccessControlListsClient) BeginValidateConfiguration(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginValidateConfigurationOptions) (*runtime.Poller[AccessControlListsClientValidateConfigurationResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.validateConfiguration(ctx, resourceGroupName, accessControlListName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AccessControlListsClientValidateConfigurationResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[AccessControlListsClientValidateConfigurationResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// ValidateConfiguration - Implements the operation to the underlying resources.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *AccessControlListsClient) validateConfiguration(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginValidateConfigurationOptions) (*http.Response, error) {
	var err error
	req, err := client.validateConfigurationCreateRequest(ctx, resourceGroupName, accessControlListName, options)
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

// validateConfigurationCreateRequest creates the ValidateConfiguration request.
func (client *AccessControlListsClient) validateConfigurationCreateRequest(ctx context.Context, resourceGroupName string, accessControlListName string, options *AccessControlListsClientBeginValidateConfigurationOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/accessControlLists/{accessControlListName}/validateConfiguration"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accessControlListName == "" {
		return nil, errors.New("parameter accessControlListName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accessControlListName}", url.PathEscape(accessControlListName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}


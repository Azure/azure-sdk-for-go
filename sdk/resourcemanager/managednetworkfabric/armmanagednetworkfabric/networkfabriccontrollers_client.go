//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

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

// NetworkFabricControllersClient contains the methods for the NetworkFabricControllers group.
// Don't use this type directly, use NewNetworkFabricControllersClient() instead.
type NetworkFabricControllersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewNetworkFabricControllersClient creates a new instance of NetworkFabricControllersClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewNetworkFabricControllersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*NetworkFabricControllersClient, error) {
	cl, err := arm.NewClient(moduleName+".NetworkFabricControllersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &NetworkFabricControllersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreate - Creates a Network Fabric Controller.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFabricControllerName - Name of the Network Fabric Controller.
//   - body - Request payload.
//   - options - NetworkFabricControllersClientBeginCreateOptions contains the optional parameters for the NetworkFabricControllersClient.BeginCreate
//     method.
func (client *NetworkFabricControllersClient) BeginCreate(ctx context.Context, resourceGroupName string, networkFabricControllerName string, body NetworkFabricController, options *NetworkFabricControllersClientBeginCreateOptions) (*runtime.Poller[NetworkFabricControllersClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, networkFabricControllerName, body, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NetworkFabricControllersClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[NetworkFabricControllersClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Creates a Network Fabric Controller.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *NetworkFabricControllersClient) create(ctx context.Context, resourceGroupName string, networkFabricControllerName string, body NetworkFabricController, options *NetworkFabricControllersClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, networkFabricControllerName, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *NetworkFabricControllersClient) createCreateRequest(ctx context.Context, resourceGroupName string, networkFabricControllerName string, body NetworkFabricController, options *NetworkFabricControllersClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/networkFabricControllers/{networkFabricControllerName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFabricControllerName == "" {
		return nil, errors.New("parameter networkFabricControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFabricControllerName}", url.PathEscape(networkFabricControllerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// BeginDelete - Deletes the Network Fabric Controller resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFabricControllerName - Name of the Network Fabric Controller.
//   - options - NetworkFabricControllersClientBeginDeleteOptions contains the optional parameters for the NetworkFabricControllersClient.BeginDelete
//     method.
func (client *NetworkFabricControllersClient) BeginDelete(ctx context.Context, resourceGroupName string, networkFabricControllerName string, options *NetworkFabricControllersClientBeginDeleteOptions) (*runtime.Poller[NetworkFabricControllersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, networkFabricControllerName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NetworkFabricControllersClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[NetworkFabricControllersClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the Network Fabric Controller resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *NetworkFabricControllersClient) deleteOperation(ctx context.Context, resourceGroupName string, networkFabricControllerName string, options *NetworkFabricControllersClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, networkFabricControllerName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *NetworkFabricControllersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, networkFabricControllerName string, options *NetworkFabricControllersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/networkFabricControllers/{networkFabricControllerName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFabricControllerName == "" {
		return nil, errors.New("parameter networkFabricControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFabricControllerName}", url.PathEscape(networkFabricControllerName))
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

// Get - Shows the provisioning status of Network Fabric Controller.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFabricControllerName - Name of the Network Fabric Controller.
//   - options - NetworkFabricControllersClientGetOptions contains the optional parameters for the NetworkFabricControllersClient.Get
//     method.
func (client *NetworkFabricControllersClient) Get(ctx context.Context, resourceGroupName string, networkFabricControllerName string, options *NetworkFabricControllersClientGetOptions) (NetworkFabricControllersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, networkFabricControllerName, options)
	if err != nil {
		return NetworkFabricControllersClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return NetworkFabricControllersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return NetworkFabricControllersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *NetworkFabricControllersClient) getCreateRequest(ctx context.Context, resourceGroupName string, networkFabricControllerName string, options *NetworkFabricControllersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/networkFabricControllers/{networkFabricControllerName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFabricControllerName == "" {
		return nil, errors.New("parameter networkFabricControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFabricControllerName}", url.PathEscape(networkFabricControllerName))
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
func (client *NetworkFabricControllersClient) getHandleResponse(resp *http.Response) (NetworkFabricControllersClientGetResponse, error) {
	result := NetworkFabricControllersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFabricController); err != nil {
		return NetworkFabricControllersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the NetworkFabricControllers thats available in the resource group.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - NetworkFabricControllersClientListByResourceGroupOptions contains the optional parameters for the NetworkFabricControllersClient.NewListByResourceGroupPager
//     method.
func (client *NetworkFabricControllersClient) NewListByResourceGroupPager(resourceGroupName string, options *NetworkFabricControllersClientListByResourceGroupOptions) *runtime.Pager[NetworkFabricControllersClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[NetworkFabricControllersClientListByResourceGroupResponse]{
		More: func(page NetworkFabricControllersClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NetworkFabricControllersClientListByResourceGroupResponse) (NetworkFabricControllersClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NetworkFabricControllersClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return NetworkFabricControllersClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NetworkFabricControllersClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *NetworkFabricControllersClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *NetworkFabricControllersClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/networkFabricControllers"
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
func (client *NetworkFabricControllersClient) listByResourceGroupHandleResponse(resp *http.Response) (NetworkFabricControllersClientListByResourceGroupResponse, error) {
	result := NetworkFabricControllersClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFabricControllersListResult); err != nil {
		return NetworkFabricControllersClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Lists all the NetworkFabricControllers by subscription.
//
// Generated from API version 2023-06-15
//   - options - NetworkFabricControllersClientListBySubscriptionOptions contains the optional parameters for the NetworkFabricControllersClient.NewListBySubscriptionPager
//     method.
func (client *NetworkFabricControllersClient) NewListBySubscriptionPager(options *NetworkFabricControllersClientListBySubscriptionOptions) *runtime.Pager[NetworkFabricControllersClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[NetworkFabricControllersClientListBySubscriptionResponse]{
		More: func(page NetworkFabricControllersClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NetworkFabricControllersClientListBySubscriptionResponse) (NetworkFabricControllersClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NetworkFabricControllersClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return NetworkFabricControllersClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NetworkFabricControllersClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *NetworkFabricControllersClient) listBySubscriptionCreateRequest(ctx context.Context, options *NetworkFabricControllersClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.ManagedNetworkFabric/networkFabricControllers"
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
func (client *NetworkFabricControllersClient) listBySubscriptionHandleResponse(resp *http.Response) (NetworkFabricControllersClientListBySubscriptionResponse, error) {
	result := NetworkFabricControllersClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFabricControllersListResult); err != nil {
		return NetworkFabricControllersClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Updates are currently not supported for the Network Fabric Controller resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFabricControllerName - Name of the Network Fabric Controller.
//   - body - Network Fabric Controller properties to update.
//   - options - NetworkFabricControllersClientBeginUpdateOptions contains the optional parameters for the NetworkFabricControllersClient.BeginUpdate
//     method.
func (client *NetworkFabricControllersClient) BeginUpdate(ctx context.Context, resourceGroupName string, networkFabricControllerName string, body NetworkFabricControllerPatch, options *NetworkFabricControllersClientBeginUpdateOptions) (*runtime.Poller[NetworkFabricControllersClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, networkFabricControllerName, body, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NetworkFabricControllersClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[NetworkFabricControllersClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Updates are currently not supported for the Network Fabric Controller resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-06-15
func (client *NetworkFabricControllersClient) update(ctx context.Context, resourceGroupName string, networkFabricControllerName string, body NetworkFabricControllerPatch, options *NetworkFabricControllersClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, networkFabricControllerName, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *NetworkFabricControllersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, networkFabricControllerName string, body NetworkFabricControllerPatch, options *NetworkFabricControllersClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedNetworkFabric/networkFabricControllers/{networkFabricControllerName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFabricControllerName == "" {
		return nil, errors.New("parameter networkFabricControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFabricControllerName}", url.PathEscape(networkFabricControllerName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}
//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhybridnetwork

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

// NetworkFunctionsClient contains the methods for the NetworkFunctions group.
// Don't use this type directly, use NewNetworkFunctionsClient() instead.
type NetworkFunctionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewNetworkFunctionsClient creates a new instance of NetworkFunctionsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewNetworkFunctionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*NetworkFunctionsClient, error) {
	cl, err := arm.NewClient(moduleName+".NetworkFunctionsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &NetworkFunctionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a network function resource. This operation can take up to 6 hours to complete.
// This is expected service behavior.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFunctionName - Resource name for the network function resource.
//   - parameters - Parameters supplied in the body to the create or update network function operation.
//   - options - NetworkFunctionsClientBeginCreateOrUpdateOptions contains the optional parameters for the NetworkFunctionsClient.BeginCreateOrUpdate
//     method.
func (client *NetworkFunctionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters NetworkFunction, options *NetworkFunctionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkFunctionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, networkFunctionName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NetworkFunctionsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[NetworkFunctionsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a network function resource. This operation can take up to 6 hours to complete. This
// is expected service behavior.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
func (client *NetworkFunctionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters NetworkFunction, options *NetworkFunctionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, networkFunctionName, parameters, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *NetworkFunctionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters NetworkFunction, options *NetworkFunctionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/networkFunctions/{networkFunctionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFunctionName == "" {
		return nil, errors.New("parameter networkFunctionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFunctionName}", url.PathEscape(networkFunctionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Deletes the specified network function resource. This operation can take up to 1 hour to complete. This is
// expected service behavior.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFunctionName - The name of the network function.
//   - options - NetworkFunctionsClientBeginDeleteOptions contains the optional parameters for the NetworkFunctionsClient.BeginDelete
//     method.
func (client *NetworkFunctionsClient) BeginDelete(ctx context.Context, resourceGroupName string, networkFunctionName string, options *NetworkFunctionsClientBeginDeleteOptions) (*runtime.Poller[NetworkFunctionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, networkFunctionName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NetworkFunctionsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[NetworkFunctionsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified network function resource. This operation can take up to 1 hour to complete. This is expected
// service behavior.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
func (client *NetworkFunctionsClient) deleteOperation(ctx context.Context, resourceGroupName string, networkFunctionName string, options *NetworkFunctionsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, networkFunctionName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *NetworkFunctionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, options *NetworkFunctionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/networkFunctions/{networkFunctionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFunctionName == "" {
		return nil, errors.New("parameter networkFunctionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFunctionName}", url.PathEscape(networkFunctionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginExecuteRequest - Execute a request to services on a network function.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFunctionName - The name of the network function.
//   - parameters - Payload for execute request post call.
//   - options - NetworkFunctionsClientBeginExecuteRequestOptions contains the optional parameters for the NetworkFunctionsClient.BeginExecuteRequest
//     method.
func (client *NetworkFunctionsClient) BeginExecuteRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters ExecuteRequestParameters, options *NetworkFunctionsClientBeginExecuteRequestOptions) (*runtime.Poller[NetworkFunctionsClientExecuteRequestResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.executeRequest(ctx, resourceGroupName, networkFunctionName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[NetworkFunctionsClientExecuteRequestResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[NetworkFunctionsClientExecuteRequestResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// ExecuteRequest - Execute a request to services on a network function.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
func (client *NetworkFunctionsClient) executeRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters ExecuteRequestParameters, options *NetworkFunctionsClientBeginExecuteRequestOptions) (*http.Response, error) {
	req, err := client.executeRequestCreateRequest(ctx, resourceGroupName, networkFunctionName, parameters, options)
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

// executeRequestCreateRequest creates the ExecuteRequest request.
func (client *NetworkFunctionsClient) executeRequestCreateRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters ExecuteRequestParameters, options *NetworkFunctionsClientBeginExecuteRequestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/networkFunctions/{networkFunctionName}/executeRequest"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFunctionName == "" {
		return nil, errors.New("parameter networkFunctionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFunctionName}", url.PathEscape(networkFunctionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// Get - Gets information about the specified network function resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFunctionName - The name of the network function resource.
//   - options - NetworkFunctionsClientGetOptions contains the optional parameters for the NetworkFunctionsClient.Get method.
func (client *NetworkFunctionsClient) Get(ctx context.Context, resourceGroupName string, networkFunctionName string, options *NetworkFunctionsClientGetOptions) (NetworkFunctionsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, networkFunctionName, options)
	if err != nil {
		return NetworkFunctionsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return NetworkFunctionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return NetworkFunctionsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *NetworkFunctionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, options *NetworkFunctionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/networkFunctions/{networkFunctionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFunctionName == "" {
		return nil, errors.New("parameter networkFunctionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFunctionName}", url.PathEscape(networkFunctionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *NetworkFunctionsClient) getHandleResponse(resp *http.Response) (NetworkFunctionsClientGetResponse, error) {
	result := NetworkFunctionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFunction); err != nil {
		return NetworkFunctionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the network function resources in a resource group.
//
// Generated from API version 2022-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - NetworkFunctionsClientListByResourceGroupOptions contains the optional parameters for the NetworkFunctionsClient.NewListByResourceGroupPager
//     method.
func (client *NetworkFunctionsClient) NewListByResourceGroupPager(resourceGroupName string, options *NetworkFunctionsClientListByResourceGroupOptions) *runtime.Pager[NetworkFunctionsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[NetworkFunctionsClientListByResourceGroupResponse]{
		More: func(page NetworkFunctionsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NetworkFunctionsClientListByResourceGroupResponse) (NetworkFunctionsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NetworkFunctionsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return NetworkFunctionsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NetworkFunctionsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *NetworkFunctionsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *NetworkFunctionsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/networkFunctions"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *NetworkFunctionsClient) listByResourceGroupHandleResponse(resp *http.Response) (NetworkFunctionsClientListByResourceGroupResponse, error) {
	result := NetworkFunctionsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFunctionListResult); err != nil {
		return NetworkFunctionsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Lists all the network functions in a subscription.
//
// Generated from API version 2022-01-01-preview
//   - options - NetworkFunctionsClientListBySubscriptionOptions contains the optional parameters for the NetworkFunctionsClient.NewListBySubscriptionPager
//     method.
func (client *NetworkFunctionsClient) NewListBySubscriptionPager(options *NetworkFunctionsClientListBySubscriptionOptions) *runtime.Pager[NetworkFunctionsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[NetworkFunctionsClientListBySubscriptionResponse]{
		More: func(page NetworkFunctionsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *NetworkFunctionsClientListBySubscriptionResponse) (NetworkFunctionsClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return NetworkFunctionsClientListBySubscriptionResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return NetworkFunctionsClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return NetworkFunctionsClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *NetworkFunctionsClient) listBySubscriptionCreateRequest(ctx context.Context, options *NetworkFunctionsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.HybridNetwork/networkFunctions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *NetworkFunctionsClient) listBySubscriptionHandleResponse(resp *http.Response) (NetworkFunctionsClientListBySubscriptionResponse, error) {
	result := NetworkFunctionsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFunctionListResult); err != nil {
		return NetworkFunctionsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// UpdateTags - Updates the tags for the network function resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - networkFunctionName - Resource name for the network function resource.
//   - parameters - Parameters supplied to the update network function tags operation.
//   - options - NetworkFunctionsClientUpdateTagsOptions contains the optional parameters for the NetworkFunctionsClient.UpdateTags
//     method.
func (client *NetworkFunctionsClient) UpdateTags(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters TagsObject, options *NetworkFunctionsClientUpdateTagsOptions) (NetworkFunctionsClientUpdateTagsResponse, error) {
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, networkFunctionName, parameters, options)
	if err != nil {
		return NetworkFunctionsClientUpdateTagsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return NetworkFunctionsClientUpdateTagsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return NetworkFunctionsClientUpdateTagsResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateTagsHandleResponse(resp)
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *NetworkFunctionsClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, networkFunctionName string, parameters TagsObject, options *NetworkFunctionsClientUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridNetwork/networkFunctions/{networkFunctionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkFunctionName == "" {
		return nil, errors.New("parameter networkFunctionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkFunctionName}", url.PathEscape(networkFunctionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *NetworkFunctionsClient) updateTagsHandleResponse(resp *http.Response) (NetworkFunctionsClientUpdateTagsResponse, error) {
	result := NetworkFunctionsClientUpdateTagsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NetworkFunction); err != nil {
		return NetworkFunctionsClientUpdateTagsResponse{}, err
	}
	return result, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetworkcloud

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

// L2NetworksClient contains the methods for the L2Networks group.
// Don't use this type directly, use NewL2NetworksClient() instead.
type L2NetworksClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewL2NetworksClient creates a new instance of L2NetworksClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewL2NetworksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*L2NetworksClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &L2NetworksClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a new layer 2 (L2) network or update the properties of the existing network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - l2NetworkName - The name of the L2 network.
//   - l2NetworkParameters - The request body.
//   - options - L2NetworksClientBeginCreateOrUpdateOptions contains the optional parameters for the L2NetworksClient.BeginCreateOrUpdate
//     method.
func (client *L2NetworksClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, l2NetworkName string, l2NetworkParameters L2Network, options *L2NetworksClientBeginCreateOrUpdateOptions) (*runtime.Poller[L2NetworksClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, l2NetworkName, l2NetworkParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[L2NetworksClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[L2NetworksClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a new layer 2 (L2) network or update the properties of the existing network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
func (client *L2NetworksClient) createOrUpdate(ctx context.Context, resourceGroupName string, l2NetworkName string, l2NetworkParameters L2Network, options *L2NetworksClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "L2NetworksClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, l2NetworkName, l2NetworkParameters, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *L2NetworksClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, l2NetworkName string, l2NetworkParameters L2Network, _ *L2NetworksClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/l2Networks/{l2NetworkName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if l2NetworkName == "" {
		return nil, errors.New("parameter l2NetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{l2NetworkName}", url.PathEscape(l2NetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, l2NetworkParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete the provided layer 2 (L2) network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - l2NetworkName - The name of the L2 network.
//   - options - L2NetworksClientBeginDeleteOptions contains the optional parameters for the L2NetworksClient.BeginDelete method.
func (client *L2NetworksClient) BeginDelete(ctx context.Context, resourceGroupName string, l2NetworkName string, options *L2NetworksClientBeginDeleteOptions) (*runtime.Poller[L2NetworksClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, l2NetworkName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[L2NetworksClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[L2NetworksClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete the provided layer 2 (L2) network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
func (client *L2NetworksClient) deleteOperation(ctx context.Context, resourceGroupName string, l2NetworkName string, options *L2NetworksClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "L2NetworksClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, l2NetworkName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *L2NetworksClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, l2NetworkName string, _ *L2NetworksClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/l2Networks/{l2NetworkName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if l2NetworkName == "" {
		return nil, errors.New("parameter l2NetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{l2NetworkName}", url.PathEscape(l2NetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get properties of the provided layer 2 (L2) network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - l2NetworkName - The name of the L2 network.
//   - options - L2NetworksClientGetOptions contains the optional parameters for the L2NetworksClient.Get method.
func (client *L2NetworksClient) Get(ctx context.Context, resourceGroupName string, l2NetworkName string, options *L2NetworksClientGetOptions) (L2NetworksClientGetResponse, error) {
	var err error
	const operationName = "L2NetworksClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, l2NetworkName, options)
	if err != nil {
		return L2NetworksClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return L2NetworksClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return L2NetworksClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *L2NetworksClient) getCreateRequest(ctx context.Context, resourceGroupName string, l2NetworkName string, _ *L2NetworksClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/l2Networks/{l2NetworkName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if l2NetworkName == "" {
		return nil, errors.New("parameter l2NetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{l2NetworkName}", url.PathEscape(l2NetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *L2NetworksClient) getHandleResponse(resp *http.Response) (L2NetworksClientGetResponse, error) {
	result := L2NetworksClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.L2Network); err != nil {
		return L2NetworksClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Get a list of layer 2 (L2) networks in the provided resource group.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - L2NetworksClientListByResourceGroupOptions contains the optional parameters for the L2NetworksClient.NewListByResourceGroupPager
//     method.
func (client *L2NetworksClient) NewListByResourceGroupPager(resourceGroupName string, options *L2NetworksClientListByResourceGroupOptions) *runtime.Pager[L2NetworksClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[L2NetworksClientListByResourceGroupResponse]{
		More: func(page L2NetworksClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *L2NetworksClientListByResourceGroupResponse) (L2NetworksClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "L2NetworksClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return L2NetworksClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *L2NetworksClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, _ *L2NetworksClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/l2Networks"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
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
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *L2NetworksClient) listByResourceGroupHandleResponse(resp *http.Response) (L2NetworksClientListByResourceGroupResponse, error) {
	result := L2NetworksClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.L2NetworkList); err != nil {
		return L2NetworksClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Get a list of layer 2 (L2) networks in the provided subscription.
//
// Generated from API version 2024-10-01-preview
//   - options - L2NetworksClientListBySubscriptionOptions contains the optional parameters for the L2NetworksClient.NewListBySubscriptionPager
//     method.
func (client *L2NetworksClient) NewListBySubscriptionPager(options *L2NetworksClientListBySubscriptionOptions) *runtime.Pager[L2NetworksClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[L2NetworksClientListBySubscriptionResponse]{
		More: func(page L2NetworksClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *L2NetworksClientListBySubscriptionResponse) (L2NetworksClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "L2NetworksClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return L2NetworksClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *L2NetworksClient) listBySubscriptionCreateRequest(ctx context.Context, _ *L2NetworksClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.NetworkCloud/l2Networks"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *L2NetworksClient) listBySubscriptionHandleResponse(resp *http.Response) (L2NetworksClientListBySubscriptionResponse, error) {
	result := L2NetworksClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.L2NetworkList); err != nil {
		return L2NetworksClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Update tags associated with the provided layer 2 (L2) network.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - l2NetworkName - The name of the L2 network.
//   - l2NetworkUpdateParameters - The request body.
//   - options - L2NetworksClientUpdateOptions contains the optional parameters for the L2NetworksClient.Update method.
func (client *L2NetworksClient) Update(ctx context.Context, resourceGroupName string, l2NetworkName string, l2NetworkUpdateParameters L2NetworkPatchParameters, options *L2NetworksClientUpdateOptions) (L2NetworksClientUpdateResponse, error) {
	var err error
	const operationName = "L2NetworksClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, l2NetworkName, l2NetworkUpdateParameters, options)
	if err != nil {
		return L2NetworksClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return L2NetworksClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return L2NetworksClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *L2NetworksClient) updateCreateRequest(ctx context.Context, resourceGroupName string, l2NetworkName string, l2NetworkUpdateParameters L2NetworkPatchParameters, _ *L2NetworksClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/l2Networks/{l2NetworkName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if l2NetworkName == "" {
		return nil, errors.New("parameter l2NetworkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{l2NetworkName}", url.PathEscape(l2NetworkName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, l2NetworkUpdateParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *L2NetworksClient) updateHandleResponse(resp *http.Response) (L2NetworksClientUpdateResponse, error) {
	result := L2NetworksClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.L2Network); err != nil {
		return L2NetworksClientUpdateResponse{}, err
	}
	return result, nil
}

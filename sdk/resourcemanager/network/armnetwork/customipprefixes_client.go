//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// CustomIPPrefixesClient contains the methods for the CustomIPPrefixes group.
// Don't use this type directly, use NewCustomIPPrefixesClient() instead.
type CustomIPPrefixesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewCustomIPPrefixesClient creates a new instance of CustomIPPrefixesClient with the specified values.
//   - subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
//     ID forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCustomIPPrefixesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CustomIPPrefixesClient, error) {
	cl, err := arm.NewClient(moduleName+".CustomIPPrefixesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CustomIPPrefixesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a custom IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The name of the resource group.
//   - customIPPrefixName - The name of the custom IP prefix.
//   - parameters - Parameters supplied to the create or update custom IP prefix operation.
//   - options - CustomIPPrefixesClientBeginCreateOrUpdateOptions contains the optional parameters for the CustomIPPrefixesClient.BeginCreateOrUpdate
//     method.
func (client *CustomIPPrefixesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, customIPPrefixName string, parameters CustomIPPrefix, options *CustomIPPrefixesClientBeginCreateOrUpdateOptions) (*runtime.Poller[CustomIPPrefixesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, customIPPrefixName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CustomIPPrefixesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CustomIPPrefixesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a custom IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
func (client *CustomIPPrefixesClient) createOrUpdate(ctx context.Context, resourceGroupName string, customIPPrefixName string, parameters CustomIPPrefix, options *CustomIPPrefixesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "CustomIPPrefixesClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, customIPPrefixName, parameters, options)
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
func (client *CustomIPPrefixesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, customIPPrefixName string, parameters CustomIPPrefix, options *CustomIPPrefixesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/customIpPrefixes/{customIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customIPPrefixName == "" {
		return nil, errors.New("parameter customIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customIpPrefixName}", url.PathEscape(customIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the specified custom IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The name of the resource group.
//   - customIPPrefixName - The name of the CustomIpPrefix.
//   - options - CustomIPPrefixesClientBeginDeleteOptions contains the optional parameters for the CustomIPPrefixesClient.BeginDelete
//     method.
func (client *CustomIPPrefixesClient) BeginDelete(ctx context.Context, resourceGroupName string, customIPPrefixName string, options *CustomIPPrefixesClientBeginDeleteOptions) (*runtime.Poller[CustomIPPrefixesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, customIPPrefixName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CustomIPPrefixesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[CustomIPPrefixesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified custom IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
func (client *CustomIPPrefixesClient) deleteOperation(ctx context.Context, resourceGroupName string, customIPPrefixName string, options *CustomIPPrefixesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "CustomIPPrefixesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, customIPPrefixName, options)
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
func (client *CustomIPPrefixesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, customIPPrefixName string, options *CustomIPPrefixesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/customIpPrefixes/{customIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customIPPrefixName == "" {
		return nil, errors.New("parameter customIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customIpPrefixName}", url.PathEscape(customIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified custom IP prefix in a specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The name of the resource group.
//   - customIPPrefixName - The name of the custom IP prefix.
//   - options - CustomIPPrefixesClientGetOptions contains the optional parameters for the CustomIPPrefixesClient.Get method.
func (client *CustomIPPrefixesClient) Get(ctx context.Context, resourceGroupName string, customIPPrefixName string, options *CustomIPPrefixesClientGetOptions) (CustomIPPrefixesClientGetResponse, error) {
	var err error
	const operationName = "CustomIPPrefixesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, customIPPrefixName, options)
	if err != nil {
		return CustomIPPrefixesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CustomIPPrefixesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CustomIPPrefixesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *CustomIPPrefixesClient) getCreateRequest(ctx context.Context, resourceGroupName string, customIPPrefixName string, options *CustomIPPrefixesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/customIpPrefixes/{customIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customIPPrefixName == "" {
		return nil, errors.New("parameter customIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customIpPrefixName}", url.PathEscape(customIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CustomIPPrefixesClient) getHandleResponse(resp *http.Response) (CustomIPPrefixesClientGetResponse, error) {
	result := CustomIPPrefixesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomIPPrefix); err != nil {
		return CustomIPPrefixesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets all custom IP prefixes in a resource group.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The name of the resource group.
//   - options - CustomIPPrefixesClientListOptions contains the optional parameters for the CustomIPPrefixesClient.NewListPager
//     method.
func (client *CustomIPPrefixesClient) NewListPager(resourceGroupName string, options *CustomIPPrefixesClientListOptions) *runtime.Pager[CustomIPPrefixesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[CustomIPPrefixesClientListResponse]{
		More: func(page CustomIPPrefixesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CustomIPPrefixesClientListResponse) (CustomIPPrefixesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "CustomIPPrefixesClient.NewListPager")
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CustomIPPrefixesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CustomIPPrefixesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CustomIPPrefixesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *CustomIPPrefixesClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *CustomIPPrefixesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/customIpPrefixes"
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
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *CustomIPPrefixesClient) listHandleResponse(resp *http.Response) (CustomIPPrefixesClientListResponse, error) {
	result := CustomIPPrefixesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomIPPrefixListResult); err != nil {
		return CustomIPPrefixesClientListResponse{}, err
	}
	return result, nil
}

// NewListAllPager - Gets all the custom IP prefixes in a subscription.
//
// Generated from API version 2023-02-01
//   - options - CustomIPPrefixesClientListAllOptions contains the optional parameters for the CustomIPPrefixesClient.NewListAllPager
//     method.
func (client *CustomIPPrefixesClient) NewListAllPager(options *CustomIPPrefixesClientListAllOptions) *runtime.Pager[CustomIPPrefixesClientListAllResponse] {
	return runtime.NewPager(runtime.PagingHandler[CustomIPPrefixesClientListAllResponse]{
		More: func(page CustomIPPrefixesClientListAllResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CustomIPPrefixesClientListAllResponse) (CustomIPPrefixesClientListAllResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "CustomIPPrefixesClient.NewListAllPager")
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listAllCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CustomIPPrefixesClientListAllResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CustomIPPrefixesClientListAllResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CustomIPPrefixesClientListAllResponse{}, runtime.NewResponseError(resp)
			}
			return client.listAllHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listAllCreateRequest creates the ListAll request.
func (client *CustomIPPrefixesClient) listAllCreateRequest(ctx context.Context, options *CustomIPPrefixesClientListAllOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/customIpPrefixes"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listAllHandleResponse handles the ListAll response.
func (client *CustomIPPrefixesClient) listAllHandleResponse(resp *http.Response) (CustomIPPrefixesClientListAllResponse, error) {
	result := CustomIPPrefixesClientListAllResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomIPPrefixListResult); err != nil {
		return CustomIPPrefixesClientListAllResponse{}, err
	}
	return result, nil
}

// UpdateTags - Updates custom IP prefix tags.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - resourceGroupName - The name of the resource group.
//   - customIPPrefixName - The name of the custom IP prefix.
//   - parameters - Parameters supplied to update custom IP prefix tags.
//   - options - CustomIPPrefixesClientUpdateTagsOptions contains the optional parameters for the CustomIPPrefixesClient.UpdateTags
//     method.
func (client *CustomIPPrefixesClient) UpdateTags(ctx context.Context, resourceGroupName string, customIPPrefixName string, parameters TagsObject, options *CustomIPPrefixesClientUpdateTagsOptions) (CustomIPPrefixesClientUpdateTagsResponse, error) {
	var err error
	const operationName = "CustomIPPrefixesClient.UpdateTags"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, customIPPrefixName, parameters, options)
	if err != nil {
		return CustomIPPrefixesClientUpdateTagsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CustomIPPrefixesClientUpdateTagsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CustomIPPrefixesClientUpdateTagsResponse{}, err
	}
	resp, err := client.updateTagsHandleResponse(httpResp)
	return resp, err
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *CustomIPPrefixesClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, customIPPrefixName string, parameters TagsObject, options *CustomIPPrefixesClientUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/customIpPrefixes/{customIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customIPPrefixName == "" {
		return nil, errors.New("parameter customIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customIpPrefixName}", url.PathEscape(customIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *CustomIPPrefixesClient) updateTagsHandleResponse(resp *http.Response) (CustomIPPrefixesClientUpdateTagsResponse, error) {
	result := CustomIPPrefixesClientUpdateTagsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomIPPrefix); err != nil {
		return CustomIPPrefixesClientUpdateTagsResponse{}, err
	}
	return result, nil
}

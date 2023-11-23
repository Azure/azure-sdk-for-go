//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventgrid

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PermissionBindingsClient contains the methods for the PermissionBindings group.
// Don't use this type directly, use NewPermissionBindingsClient() instead.
type PermissionBindingsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewPermissionBindingsClient creates a new instance of PermissionBindingsClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPermissionBindingsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PermissionBindingsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PermissionBindingsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a permission binding with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-15-preview
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - namespaceName - Name of the namespace.
//   - permissionBindingName - The permission binding name.
//   - permissionBindingInfo - Permission binding information.
//   - options - PermissionBindingsClientBeginCreateOrUpdateOptions contains the optional parameters for the PermissionBindingsClient.BeginCreateOrUpdate
//     method.
func (client *PermissionBindingsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, permissionBindingInfo PermissionBinding, options *PermissionBindingsClientBeginCreateOrUpdateOptions) (*runtime.Poller[PermissionBindingsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, namespaceName, permissionBindingName, permissionBindingInfo, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PermissionBindingsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PermissionBindingsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create or update a permission binding with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-15-preview
func (client *PermissionBindingsClient) createOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, permissionBindingInfo PermissionBinding, options *PermissionBindingsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "PermissionBindingsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, namespaceName, permissionBindingName, permissionBindingInfo, options)
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
func (client *PermissionBindingsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, permissionBindingInfo PermissionBinding, options *PermissionBindingsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/namespaces/{namespaceName}/permissionBindings/{permissionBindingName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if permissionBindingName == "" {
		return nil, errors.New("parameter permissionBindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{permissionBindingName}", url.PathEscape(permissionBindingName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-15-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, permissionBindingInfo); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete an existing permission binding.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-15-preview
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - namespaceName - Name of the namespace.
//   - permissionBindingName - Name of the permission binding.
//   - options - PermissionBindingsClientBeginDeleteOptions contains the optional parameters for the PermissionBindingsClient.BeginDelete
//     method.
func (client *PermissionBindingsClient) BeginDelete(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, options *PermissionBindingsClientBeginDeleteOptions) (*runtime.Poller[PermissionBindingsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, namespaceName, permissionBindingName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PermissionBindingsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PermissionBindingsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete an existing permission binding.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-15-preview
func (client *PermissionBindingsClient) deleteOperation(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, options *PermissionBindingsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "PermissionBindingsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, namespaceName, permissionBindingName, options)
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
func (client *PermissionBindingsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, options *PermissionBindingsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/namespaces/{namespaceName}/permissionBindings/{permissionBindingName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if permissionBindingName == "" {
		return nil, errors.New("parameter permissionBindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{permissionBindingName}", url.PathEscape(permissionBindingName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-15-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get properties of a permission binding.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-15-preview
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - namespaceName - Name of the namespace.
//   - permissionBindingName - Name of the permission binding.
//   - options - PermissionBindingsClientGetOptions contains the optional parameters for the PermissionBindingsClient.Get method.
func (client *PermissionBindingsClient) Get(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, options *PermissionBindingsClientGetOptions) (PermissionBindingsClientGetResponse, error) {
	var err error
	const operationName = "PermissionBindingsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, namespaceName, permissionBindingName, options)
	if err != nil {
		return PermissionBindingsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PermissionBindingsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PermissionBindingsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PermissionBindingsClient) getCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, permissionBindingName string, options *PermissionBindingsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/namespaces/{namespaceName}/permissionBindings/{permissionBindingName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if permissionBindingName == "" {
		return nil, errors.New("parameter permissionBindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{permissionBindingName}", url.PathEscape(permissionBindingName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-15-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PermissionBindingsClient) getHandleResponse(resp *http.Response) (PermissionBindingsClientGetResponse, error) {
	result := PermissionBindingsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PermissionBinding); err != nil {
		return PermissionBindingsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByNamespacePager - Get all the permission bindings under a namespace.
//
// Generated from API version 2023-12-15-preview
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - namespaceName - Name of the namespace.
//   - options - PermissionBindingsClientListByNamespaceOptions contains the optional parameters for the PermissionBindingsClient.NewListByNamespacePager
//     method.
func (client *PermissionBindingsClient) NewListByNamespacePager(resourceGroupName string, namespaceName string, options *PermissionBindingsClientListByNamespaceOptions) *runtime.Pager[PermissionBindingsClientListByNamespaceResponse] {
	return runtime.NewPager(runtime.PagingHandler[PermissionBindingsClientListByNamespaceResponse]{
		More: func(page PermissionBindingsClientListByNamespaceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PermissionBindingsClientListByNamespaceResponse) (PermissionBindingsClientListByNamespaceResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PermissionBindingsClient.NewListByNamespacePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByNamespaceCreateRequest(ctx, resourceGroupName, namespaceName, options)
			}, nil)
			if err != nil {
				return PermissionBindingsClientListByNamespaceResponse{}, err
			}
			return client.listByNamespaceHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByNamespaceCreateRequest creates the ListByNamespace request.
func (client *PermissionBindingsClient) listByNamespaceCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, options *PermissionBindingsClientListByNamespaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/namespaces/{namespaceName}/permissionBindings"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-15-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByNamespaceHandleResponse handles the ListByNamespace response.
func (client *PermissionBindingsClient) listByNamespaceHandleResponse(resp *http.Response) (PermissionBindingsClientListByNamespaceResponse, error) {
	result := PermissionBindingsClientListByNamespaceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PermissionBindingsListResult); err != nil {
		return PermissionBindingsClientListByNamespaceResponse{}, err
	}
	return result, nil
}
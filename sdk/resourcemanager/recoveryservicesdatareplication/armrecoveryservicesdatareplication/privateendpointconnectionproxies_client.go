// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armrecoveryservicesdatareplication

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

// PrivateEndpointConnectionProxiesClient contains the methods for the PrivateEndpointConnectionProxies group.
// Don't use this type directly, use NewPrivateEndpointConnectionProxiesClient() instead.
type PrivateEndpointConnectionProxiesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewPrivateEndpointConnectionProxiesClient creates a new instance of PrivateEndpointConnectionProxiesClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPrivateEndpointConnectionProxiesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PrivateEndpointConnectionProxiesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PrivateEndpointConnectionProxiesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Create - Create a new private endpoint connection proxy which includes both auto and manual approval types. Creating the
// proxy resource will also create a private endpoint connection resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - privateEndpointConnectionProxyName - The private endpoint connection proxy name.
//   - resource - Private endpoint connection creation input.
//   - options - PrivateEndpointConnectionProxiesClientCreateOptions contains the optional parameters for the PrivateEndpointConnectionProxiesClient.Create
//     method.
func (client *PrivateEndpointConnectionProxiesClient) Create(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, resource PrivateEndpointConnectionProxy, options *PrivateEndpointConnectionProxiesClientCreateOptions) (PrivateEndpointConnectionProxiesClientCreateResponse, error) {
	var err error
	const operationName = "PrivateEndpointConnectionProxiesClient.Create"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createCreateRequest(ctx, resourceGroupName, vaultName, privateEndpointConnectionProxyName, resource, options)
	if err != nil {
		return PrivateEndpointConnectionProxiesClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PrivateEndpointConnectionProxiesClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return PrivateEndpointConnectionProxiesClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *PrivateEndpointConnectionProxiesClient) createCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, resource PrivateEndpointConnectionProxy, _ *PrivateEndpointConnectionProxiesClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/privateEndpointConnectionProxies/{privateEndpointConnectionProxyName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if privateEndpointConnectionProxyName == "" {
		return nil, errors.New("parameter privateEndpointConnectionProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateEndpointConnectionProxyName}", url.PathEscape(privateEndpointConnectionProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
		return nil, err
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *PrivateEndpointConnectionProxiesClient) createHandleResponse(resp *http.Response) (PrivateEndpointConnectionProxiesClientCreateResponse, error) {
	result := PrivateEndpointConnectionProxiesClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateEndpointConnectionProxy); err != nil {
		return PrivateEndpointConnectionProxiesClientCreateResponse{}, err
	}
	return result, nil
}

// BeginDelete - Returns the operation to track the deletion of private endpoint connection proxy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - privateEndpointConnectionProxyName - The private endpoint connection proxy name.
//   - options - PrivateEndpointConnectionProxiesClientBeginDeleteOptions contains the optional parameters for the PrivateEndpointConnectionProxiesClient.BeginDelete
//     method.
func (client *PrivateEndpointConnectionProxiesClient) BeginDelete(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, options *PrivateEndpointConnectionProxiesClientBeginDeleteOptions) (*runtime.Poller[PrivateEndpointConnectionProxiesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, vaultName, privateEndpointConnectionProxyName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PrivateEndpointConnectionProxiesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[PrivateEndpointConnectionProxiesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Returns the operation to track the deletion of private endpoint connection proxy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
func (client *PrivateEndpointConnectionProxiesClient) deleteOperation(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, options *PrivateEndpointConnectionProxiesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "PrivateEndpointConnectionProxiesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, vaultName, privateEndpointConnectionProxyName, options)
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
func (client *PrivateEndpointConnectionProxiesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, _ *PrivateEndpointConnectionProxiesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/privateEndpointConnectionProxies/{privateEndpointConnectionProxyName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if privateEndpointConnectionProxyName == "" {
		return nil, errors.New("parameter privateEndpointConnectionProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateEndpointConnectionProxyName}", url.PathEscape(privateEndpointConnectionProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the private endpoint connection proxy details.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - privateEndpointConnectionProxyName - The private endpoint connection proxy name.
//   - options - PrivateEndpointConnectionProxiesClientGetOptions contains the optional parameters for the PrivateEndpointConnectionProxiesClient.Get
//     method.
func (client *PrivateEndpointConnectionProxiesClient) Get(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, options *PrivateEndpointConnectionProxiesClientGetOptions) (PrivateEndpointConnectionProxiesClientGetResponse, error) {
	var err error
	const operationName = "PrivateEndpointConnectionProxiesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, vaultName, privateEndpointConnectionProxyName, options)
	if err != nil {
		return PrivateEndpointConnectionProxiesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PrivateEndpointConnectionProxiesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PrivateEndpointConnectionProxiesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PrivateEndpointConnectionProxiesClient) getCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, _ *PrivateEndpointConnectionProxiesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/privateEndpointConnectionProxies/{privateEndpointConnectionProxyName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if privateEndpointConnectionProxyName == "" {
		return nil, errors.New("parameter privateEndpointConnectionProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateEndpointConnectionProxyName}", url.PathEscape(privateEndpointConnectionProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PrivateEndpointConnectionProxiesClient) getHandleResponse(resp *http.Response) (PrivateEndpointConnectionProxiesClientGetResponse, error) {
	result := PrivateEndpointConnectionProxiesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateEndpointConnectionProxy); err != nil {
		return PrivateEndpointConnectionProxiesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets the all private endpoint connections proxies.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - options - PrivateEndpointConnectionProxiesClientListOptions contains the optional parameters for the PrivateEndpointConnectionProxiesClient.NewListPager
//     method.
func (client *PrivateEndpointConnectionProxiesClient) NewListPager(resourceGroupName string, vaultName string, options *PrivateEndpointConnectionProxiesClientListOptions) *runtime.Pager[PrivateEndpointConnectionProxiesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[PrivateEndpointConnectionProxiesClientListResponse]{
		More: func(page PrivateEndpointConnectionProxiesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PrivateEndpointConnectionProxiesClientListResponse) (PrivateEndpointConnectionProxiesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PrivateEndpointConnectionProxiesClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, vaultName, options)
			}, nil)
			if err != nil {
				return PrivateEndpointConnectionProxiesClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *PrivateEndpointConnectionProxiesClient) listCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, _ *PrivateEndpointConnectionProxiesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/privateEndpointConnectionProxies"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PrivateEndpointConnectionProxiesClient) listHandleResponse(resp *http.Response) (PrivateEndpointConnectionProxiesClientListResponse, error) {
	result := PrivateEndpointConnectionProxiesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateEndpointConnectionProxyListResult); err != nil {
		return PrivateEndpointConnectionProxiesClientListResponse{}, err
	}
	return result, nil
}

// Validate - Returns remote private endpoint connection information after validation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - privateEndpointConnectionProxyName - The private endpoint connection proxy name.
//   - body - The private endpoint connection proxy input.
//   - options - PrivateEndpointConnectionProxiesClientValidateOptions contains the optional parameters for the PrivateEndpointConnectionProxiesClient.Validate
//     method.
func (client *PrivateEndpointConnectionProxiesClient) Validate(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, body PrivateEndpointConnectionProxy, options *PrivateEndpointConnectionProxiesClientValidateOptions) (PrivateEndpointConnectionProxiesClientValidateResponse, error) {
	var err error
	const operationName = "PrivateEndpointConnectionProxiesClient.Validate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.validateCreateRequest(ctx, resourceGroupName, vaultName, privateEndpointConnectionProxyName, body, options)
	if err != nil {
		return PrivateEndpointConnectionProxiesClientValidateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PrivateEndpointConnectionProxiesClientValidateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PrivateEndpointConnectionProxiesClientValidateResponse{}, err
	}
	resp, err := client.validateHandleResponse(httpResp)
	return resp, err
}

// validateCreateRequest creates the Validate request.
func (client *PrivateEndpointConnectionProxiesClient) validateCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, privateEndpointConnectionProxyName string, body PrivateEndpointConnectionProxy, _ *PrivateEndpointConnectionProxiesClientValidateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/privateEndpointConnectionProxies/{privateEndpointConnectionProxyName}/validate"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if privateEndpointConnectionProxyName == "" {
		return nil, errors.New("parameter privateEndpointConnectionProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{privateEndpointConnectionProxyName}", url.PathEscape(privateEndpointConnectionProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return nil, err
	}
	return req, nil
}

// validateHandleResponse handles the Validate response.
func (client *PrivateEndpointConnectionProxiesClient) validateHandleResponse(resp *http.Response) (PrivateEndpointConnectionProxiesClientValidateResponse, error) {
	result := PrivateEndpointConnectionProxiesClientValidateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateEndpointConnectionProxy); err != nil {
		return PrivateEndpointConnectionProxiesClientValidateResponse{}, err
	}
	return result, nil
}

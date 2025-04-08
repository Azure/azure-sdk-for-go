// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armservicenetworking

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

// FrontendsInterfaceClient contains the methods for the FrontendsInterface group.
// Don't use this type directly, use NewFrontendsInterfaceClient() instead.
type FrontendsInterfaceClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewFrontendsInterfaceClient creates a new instance of FrontendsInterfaceClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewFrontendsInterfaceClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FrontendsInterfaceClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &FrontendsInterfaceClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a Frontend
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - trafficControllerName - traffic controller name for path
//   - frontendName - Frontends
//   - resource - Resource create parameters.
//   - options - FrontendsInterfaceClientBeginCreateOrUpdateOptions contains the optional parameters for the FrontendsInterfaceClient.BeginCreateOrUpdate
//     method.
func (client *FrontendsInterfaceClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, resource Frontend, options *FrontendsInterfaceClientBeginCreateOrUpdateOptions) (*runtime.Poller[FrontendsInterfaceClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, trafficControllerName, frontendName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[FrontendsInterfaceClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[FrontendsInterfaceClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a Frontend
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
func (client *FrontendsInterfaceClient) createOrUpdate(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, resource Frontend, options *FrontendsInterfaceClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "FrontendsInterfaceClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, trafficControllerName, frontendName, resource, options)
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
func (client *FrontendsInterfaceClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, resource Frontend, _ *FrontendsInterfaceClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{trafficControllerName}/frontends/{frontendName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if trafficControllerName == "" {
		return nil, errors.New("parameter trafficControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trafficControllerName}", url.PathEscape(trafficControllerName))
	if frontendName == "" {
		return nil, errors.New("parameter frontendName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{frontendName}", url.PathEscape(frontendName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete a Frontend
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - trafficControllerName - traffic controller name for path
//   - frontendName - Frontends
//   - options - FrontendsInterfaceClientBeginDeleteOptions contains the optional parameters for the FrontendsInterfaceClient.BeginDelete
//     method.
func (client *FrontendsInterfaceClient) BeginDelete(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, options *FrontendsInterfaceClientBeginDeleteOptions) (*runtime.Poller[FrontendsInterfaceClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, trafficControllerName, frontendName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[FrontendsInterfaceClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[FrontendsInterfaceClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete a Frontend
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
func (client *FrontendsInterfaceClient) deleteOperation(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, options *FrontendsInterfaceClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "FrontendsInterfaceClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, trafficControllerName, frontendName, options)
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
func (client *FrontendsInterfaceClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, _ *FrontendsInterfaceClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{trafficControllerName}/frontends/{frontendName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if trafficControllerName == "" {
		return nil, errors.New("parameter trafficControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trafficControllerName}", url.PathEscape(trafficControllerName))
	if frontendName == "" {
		return nil, errors.New("parameter frontendName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{frontendName}", url.PathEscape(frontendName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a Frontend
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - trafficControllerName - traffic controller name for path
//   - frontendName - Frontends
//   - options - FrontendsInterfaceClientGetOptions contains the optional parameters for the FrontendsInterfaceClient.Get method.
func (client *FrontendsInterfaceClient) Get(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, options *FrontendsInterfaceClientGetOptions) (FrontendsInterfaceClientGetResponse, error) {
	var err error
	const operationName = "FrontendsInterfaceClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, trafficControllerName, frontendName, options)
	if err != nil {
		return FrontendsInterfaceClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FrontendsInterfaceClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FrontendsInterfaceClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *FrontendsInterfaceClient) getCreateRequest(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, _ *FrontendsInterfaceClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{trafficControllerName}/frontends/{frontendName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if trafficControllerName == "" {
		return nil, errors.New("parameter trafficControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trafficControllerName}", url.PathEscape(trafficControllerName))
	if frontendName == "" {
		return nil, errors.New("parameter frontendName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{frontendName}", url.PathEscape(frontendName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *FrontendsInterfaceClient) getHandleResponse(resp *http.Response) (FrontendsInterfaceClientGetResponse, error) {
	result := FrontendsInterfaceClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Frontend); err != nil {
		return FrontendsInterfaceClientGetResponse{}, err
	}
	return result, nil
}

// NewListByTrafficControllerPager - List Frontend resources by TrafficController
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - trafficControllerName - traffic controller name for path
//   - options - FrontendsInterfaceClientListByTrafficControllerOptions contains the optional parameters for the FrontendsInterfaceClient.NewListByTrafficControllerPager
//     method.
func (client *FrontendsInterfaceClient) NewListByTrafficControllerPager(resourceGroupName string, trafficControllerName string, options *FrontendsInterfaceClientListByTrafficControllerOptions) *runtime.Pager[FrontendsInterfaceClientListByTrafficControllerResponse] {
	return runtime.NewPager(runtime.PagingHandler[FrontendsInterfaceClientListByTrafficControllerResponse]{
		More: func(page FrontendsInterfaceClientListByTrafficControllerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *FrontendsInterfaceClientListByTrafficControllerResponse) (FrontendsInterfaceClientListByTrafficControllerResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "FrontendsInterfaceClient.NewListByTrafficControllerPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByTrafficControllerCreateRequest(ctx, resourceGroupName, trafficControllerName, options)
			}, nil)
			if err != nil {
				return FrontendsInterfaceClientListByTrafficControllerResponse{}, err
			}
			return client.listByTrafficControllerHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByTrafficControllerCreateRequest creates the ListByTrafficController request.
func (client *FrontendsInterfaceClient) listByTrafficControllerCreateRequest(ctx context.Context, resourceGroupName string, trafficControllerName string, _ *FrontendsInterfaceClientListByTrafficControllerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{trafficControllerName}/frontends"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if trafficControllerName == "" {
		return nil, errors.New("parameter trafficControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trafficControllerName}", url.PathEscape(trafficControllerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByTrafficControllerHandleResponse handles the ListByTrafficController response.
func (client *FrontendsInterfaceClient) listByTrafficControllerHandleResponse(resp *http.Response) (FrontendsInterfaceClientListByTrafficControllerResponse, error) {
	result := FrontendsInterfaceClientListByTrafficControllerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FrontendListResult); err != nil {
		return FrontendsInterfaceClientListByTrafficControllerResponse{}, err
	}
	return result, nil
}

// Update - Update a Frontend
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - trafficControllerName - traffic controller name for path
//   - frontendName - Frontends
//   - properties - The resource properties to be updated.
//   - options - FrontendsInterfaceClientUpdateOptions contains the optional parameters for the FrontendsInterfaceClient.Update
//     method.
func (client *FrontendsInterfaceClient) Update(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, properties FrontendUpdate, options *FrontendsInterfaceClientUpdateOptions) (FrontendsInterfaceClientUpdateResponse, error) {
	var err error
	const operationName = "FrontendsInterfaceClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, trafficControllerName, frontendName, properties, options)
	if err != nil {
		return FrontendsInterfaceClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FrontendsInterfaceClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FrontendsInterfaceClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *FrontendsInterfaceClient) updateCreateRequest(ctx context.Context, resourceGroupName string, trafficControllerName string, frontendName string, properties FrontendUpdate, _ *FrontendsInterfaceClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceNetworking/trafficControllers/{trafficControllerName}/frontends/{frontendName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if trafficControllerName == "" {
		return nil, errors.New("parameter trafficControllerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{trafficControllerName}", url.PathEscape(trafficControllerName))
	if frontendName == "" {
		return nil, errors.New("parameter frontendName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{frontendName}", url.PathEscape(frontendName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, properties); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *FrontendsInterfaceClient) updateHandleResponse(resp *http.Response) (FrontendsInterfaceClientUpdateResponse, error) {
	result := FrontendsInterfaceClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Frontend); err != nil {
		return FrontendsInterfaceClientUpdateResponse{}, err
	}
	return result, nil
}

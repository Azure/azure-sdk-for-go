//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappplatform

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

// GatewayRouteConfigsClient contains the methods for the GatewayRouteConfigs group.
// Don't use this type directly, use NewGatewayRouteConfigsClient() instead.
type GatewayRouteConfigsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewGatewayRouteConfigsClient creates a new instance of GatewayRouteConfigsClient with the specified values.
//   - subscriptionID - Gets subscription ID which uniquely identify the Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewGatewayRouteConfigsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GatewayRouteConfigsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &GatewayRouteConfigsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create the default Spring Cloud Gateway route configs or update the existing Spring Cloud Gateway
// route configs.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - gatewayName - The name of Spring Cloud Gateway.
//   - routeConfigName - The name of the Spring Cloud Gateway route config.
//   - gatewayRouteConfigResource - The Spring Cloud Gateway route config for the create or update operation
//   - options - GatewayRouteConfigsClientBeginCreateOrUpdateOptions contains the optional parameters for the GatewayRouteConfigsClient.BeginCreateOrUpdate
//     method.
func (client *GatewayRouteConfigsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, gatewayRouteConfigResource GatewayRouteConfigResource, options *GatewayRouteConfigsClientBeginCreateOrUpdateOptions) (*runtime.Poller[GatewayRouteConfigsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, serviceName, gatewayName, routeConfigName, gatewayRouteConfigResource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GatewayRouteConfigsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[GatewayRouteConfigsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create the default Spring Cloud Gateway route configs or update the existing Spring Cloud Gateway route
// configs.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
func (client *GatewayRouteConfigsClient) createOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, gatewayRouteConfigResource GatewayRouteConfigResource, options *GatewayRouteConfigsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "GatewayRouteConfigsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, gatewayName, routeConfigName, gatewayRouteConfigResource, options)
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
func (client *GatewayRouteConfigsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, gatewayRouteConfigResource GatewayRouteConfigResource, options *GatewayRouteConfigsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/gateways/{gatewayName}/routeConfigs/{routeConfigName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if gatewayName == "" {
		return nil, errors.New("parameter gatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gatewayName}", url.PathEscape(gatewayName))
	if routeConfigName == "" {
		return nil, errors.New("parameter routeConfigName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{routeConfigName}", url.PathEscape(routeConfigName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, gatewayRouteConfigResource); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete the Spring Cloud Gateway route config.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - gatewayName - The name of Spring Cloud Gateway.
//   - routeConfigName - The name of the Spring Cloud Gateway route config.
//   - options - GatewayRouteConfigsClientBeginDeleteOptions contains the optional parameters for the GatewayRouteConfigsClient.BeginDelete
//     method.
func (client *GatewayRouteConfigsClient) BeginDelete(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, options *GatewayRouteConfigsClientBeginDeleteOptions) (*runtime.Poller[GatewayRouteConfigsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, serviceName, gatewayName, routeConfigName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GatewayRouteConfigsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[GatewayRouteConfigsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete the Spring Cloud Gateway route config.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
func (client *GatewayRouteConfigsClient) deleteOperation(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, options *GatewayRouteConfigsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "GatewayRouteConfigsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, gatewayName, routeConfigName, options)
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
func (client *GatewayRouteConfigsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, options *GatewayRouteConfigsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/gateways/{gatewayName}/routeConfigs/{routeConfigName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if gatewayName == "" {
		return nil, errors.New("parameter gatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gatewayName}", url.PathEscape(gatewayName))
	if routeConfigName == "" {
		return nil, errors.New("parameter routeConfigName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{routeConfigName}", url.PathEscape(routeConfigName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the Spring Cloud Gateway route configs.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - gatewayName - The name of Spring Cloud Gateway.
//   - routeConfigName - The name of the Spring Cloud Gateway route config.
//   - options - GatewayRouteConfigsClientGetOptions contains the optional parameters for the GatewayRouteConfigsClient.Get method.
func (client *GatewayRouteConfigsClient) Get(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, options *GatewayRouteConfigsClientGetOptions) (GatewayRouteConfigsClientGetResponse, error) {
	var err error
	const operationName = "GatewayRouteConfigsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, gatewayName, routeConfigName, options)
	if err != nil {
		return GatewayRouteConfigsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GatewayRouteConfigsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GatewayRouteConfigsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *GatewayRouteConfigsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, routeConfigName string, options *GatewayRouteConfigsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/gateways/{gatewayName}/routeConfigs/{routeConfigName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if gatewayName == "" {
		return nil, errors.New("parameter gatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gatewayName}", url.PathEscape(gatewayName))
	if routeConfigName == "" {
		return nil, errors.New("parameter routeConfigName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{routeConfigName}", url.PathEscape(routeConfigName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GatewayRouteConfigsClient) getHandleResponse(resp *http.Response) (GatewayRouteConfigsClientGetResponse, error) {
	result := GatewayRouteConfigsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GatewayRouteConfigResource); err != nil {
		return GatewayRouteConfigsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Handle requests to list all Spring Cloud Gateway route configs.
//
// Generated from API version 2023-12-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serviceName - The name of the Service resource.
//   - gatewayName - The name of Spring Cloud Gateway.
//   - options - GatewayRouteConfigsClientListOptions contains the optional parameters for the GatewayRouteConfigsClient.NewListPager
//     method.
func (client *GatewayRouteConfigsClient) NewListPager(resourceGroupName string, serviceName string, gatewayName string, options *GatewayRouteConfigsClientListOptions) *runtime.Pager[GatewayRouteConfigsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[GatewayRouteConfigsClientListResponse]{
		More: func(page GatewayRouteConfigsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *GatewayRouteConfigsClientListResponse) (GatewayRouteConfigsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "GatewayRouteConfigsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, serviceName, gatewayName, options)
			}, nil)
			if err != nil {
				return GatewayRouteConfigsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *GatewayRouteConfigsClient) listCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, gatewayName string, options *GatewayRouteConfigsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/gateways/{gatewayName}/routeConfigs"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if gatewayName == "" {
		return nil, errors.New("parameter gatewayName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{gatewayName}", url.PathEscape(gatewayName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *GatewayRouteConfigsClient) listHandleResponse(resp *http.Response) (GatewayRouteConfigsClientListResponse, error) {
	result := GatewayRouteConfigsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GatewayRouteConfigResourceCollection); err != nil {
		return GatewayRouteConfigsClientListResponse{}, err
	}
	return result, nil
}

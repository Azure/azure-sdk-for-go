//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// BackendClient contains the methods for the Backend group.
// Don't use this type directly, use NewBackendClient() instead.
type BackendClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewBackendClient creates a new instance of BackendClient with the specified values.
// subscriptionID - Subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID forms
// part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewBackendClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BackendClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &BackendClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or Updates a backend.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// backendID - Identifier of the Backend entity. Must be unique in the current API Management service instance.
// parameters - Create parameters.
// options - BackendClientCreateOrUpdateOptions contains the optional parameters for the BackendClient.CreateOrUpdate method.
func (client *BackendClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, backendID string, parameters BackendContract, options *BackendClientCreateOrUpdateOptions) (BackendClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, backendID, parameters, options)
	if err != nil {
		return BackendClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackendClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return BackendClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *BackendClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, backendID string, parameters BackendContract, options *BackendClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends/{backendId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if backendID == "" {
		return nil, errors.New("parameter backendID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{backendId}", url.PathEscape(backendID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.IfMatch != nil {
		req.Raw().Header.Set("If-Match", *options.IfMatch)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *BackendClient) createOrUpdateHandleResponse(resp *http.Response) (BackendClientCreateOrUpdateResponse, error) {
	result := BackendClientCreateOrUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackendContract); err != nil {
		return BackendClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the specified backend.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// backendID - Identifier of the Backend entity. Must be unique in the current API Management service instance.
// ifMatch - ETag of the Entity. ETag should match the current entity state from the header response of the GET request or
// it should be * for unconditional update.
// options - BackendClientDeleteOptions contains the optional parameters for the BackendClient.Delete method.
func (client *BackendClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, backendID string, ifMatch string, options *BackendClientDeleteOptions) (BackendClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, backendID, ifMatch, options)
	if err != nil {
		return BackendClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackendClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return BackendClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return BackendClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *BackendClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, backendID string, ifMatch string, options *BackendClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends/{backendId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if backendID == "" {
		return nil, errors.New("parameter backendID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{backendId}", url.PathEscape(backendID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("If-Match", ifMatch)
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets the details of the backend specified by its identifier.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// backendID - Identifier of the Backend entity. Must be unique in the current API Management service instance.
// options - BackendClientGetOptions contains the optional parameters for the BackendClient.Get method.
func (client *BackendClient) Get(ctx context.Context, resourceGroupName string, serviceName string, backendID string, options *BackendClientGetOptions) (BackendClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, backendID, options)
	if err != nil {
		return BackendClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackendClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return BackendClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *BackendClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, backendID string, options *BackendClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends/{backendId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if backendID == "" {
		return nil, errors.New("parameter backendID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{backendId}", url.PathEscape(backendID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *BackendClient) getHandleResponse(resp *http.Response) (BackendClientGetResponse, error) {
	result := BackendClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackendContract); err != nil {
		return BackendClientGetResponse{}, err
	}
	return result, nil
}

// GetEntityTag - Gets the entity state (Etag) version of the backend specified by its identifier.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// backendID - Identifier of the Backend entity. Must be unique in the current API Management service instance.
// options - BackendClientGetEntityTagOptions contains the optional parameters for the BackendClient.GetEntityTag method.
func (client *BackendClient) GetEntityTag(ctx context.Context, resourceGroupName string, serviceName string, backendID string, options *BackendClientGetEntityTagOptions) (BackendClientGetEntityTagResponse, error) {
	req, err := client.getEntityTagCreateRequest(ctx, resourceGroupName, serviceName, backendID, options)
	if err != nil {
		return BackendClientGetEntityTagResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackendClientGetEntityTagResponse{}, err
	}
	return client.getEntityTagHandleResponse(resp)
}

// getEntityTagCreateRequest creates the GetEntityTag request.
func (client *BackendClient) getEntityTagCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, backendID string, options *BackendClientGetEntityTagOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends/{backendId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if backendID == "" {
		return nil, errors.New("parameter backendID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{backendId}", url.PathEscape(backendID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodHead, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getEntityTagHandleResponse handles the GetEntityTag response.
func (client *BackendClient) getEntityTagHandleResponse(resp *http.Response) (BackendClientGetEntityTagResponse, error) {
	result := BackendClientGetEntityTagResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Success = true
	}
	return result, nil
}

// NewListByServicePager - Lists a collection of backends in the specified service instance.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// options - BackendClientListByServiceOptions contains the optional parameters for the BackendClient.ListByService method.
func (client *BackendClient) NewListByServicePager(resourceGroupName string, serviceName string, options *BackendClientListByServiceOptions) *runtime.Pager[BackendClientListByServiceResponse] {
	return runtime.NewPager(runtime.PageProcessor[BackendClientListByServiceResponse]{
		More: func(page BackendClientListByServiceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *BackendClientListByServiceResponse) (BackendClientListByServiceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByServiceCreateRequest(ctx, resourceGroupName, serviceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return BackendClientListByServiceResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return BackendClientListByServiceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return BackendClientListByServiceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByServiceHandleResponse(resp)
		},
	})
}

// listByServiceCreateRequest creates the ListByService request.
func (client *BackendClient) listByServiceCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, options *BackendClientListByServiceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByServiceHandleResponse handles the ListByService response.
func (client *BackendClient) listByServiceHandleResponse(resp *http.Response) (BackendClientListByServiceResponse, error) {
	result := BackendClientListByServiceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackendCollection); err != nil {
		return BackendClientListByServiceResponse{}, err
	}
	return result, nil
}

// Reconnect - Notifies the APIM proxy to create a new connection to the backend after the specified timeout. If no timeout
// was specified, timeout of 2 minutes is used.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// backendID - Identifier of the Backend entity. Must be unique in the current API Management service instance.
// options - BackendClientReconnectOptions contains the optional parameters for the BackendClient.Reconnect method.
func (client *BackendClient) Reconnect(ctx context.Context, resourceGroupName string, serviceName string, backendID string, options *BackendClientReconnectOptions) (BackendClientReconnectResponse, error) {
	req, err := client.reconnectCreateRequest(ctx, resourceGroupName, serviceName, backendID, options)
	if err != nil {
		return BackendClientReconnectResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackendClientReconnectResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return BackendClientReconnectResponse{}, runtime.NewResponseError(resp)
	}
	return BackendClientReconnectResponse{}, nil
}

// reconnectCreateRequest creates the Reconnect request.
func (client *BackendClient) reconnectCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, backendID string, options *BackendClientReconnectOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends/{backendId}/reconnect"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if backendID == "" {
		return nil, errors.New("parameter backendID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{backendId}", url.PathEscape(backendID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	if options != nil && options.Parameters != nil {
		return req, runtime.MarshalAsJSON(req, *options.Parameters)
	}
	return req, nil
}

// Update - Updates an existing backend.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// serviceName - The name of the API Management service.
// backendID - Identifier of the Backend entity. Must be unique in the current API Management service instance.
// ifMatch - ETag of the Entity. ETag should match the current entity state from the header response of the GET request or
// it should be * for unconditional update.
// parameters - Update parameters.
// options - BackendClientUpdateOptions contains the optional parameters for the BackendClient.Update method.
func (client *BackendClient) Update(ctx context.Context, resourceGroupName string, serviceName string, backendID string, ifMatch string, parameters BackendUpdateParameters, options *BackendClientUpdateOptions) (BackendClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serviceName, backendID, ifMatch, parameters, options)
	if err != nil {
		return BackendClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BackendClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return BackendClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *BackendClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, backendID string, ifMatch string, parameters BackendUpdateParameters, options *BackendClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/backends/{backendId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if backendID == "" {
		return nil, errors.New("parameter backendID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{backendId}", url.PathEscape(backendID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("If-Match", ifMatch)
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *BackendClient) updateHandleResponse(resp *http.Response) (BackendClientUpdateResponse, error) {
	result := BackendClientUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.BackendContract); err != nil {
		return BackendClientUpdateResponse{}, err
	}
	return result, nil
}

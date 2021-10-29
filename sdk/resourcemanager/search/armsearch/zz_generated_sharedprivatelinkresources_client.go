//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// SharedPrivateLinkResourcesClient contains the methods for the SharedPrivateLinkResources group.
// Don't use this type directly, use NewSharedPrivateLinkResourcesClient() instead.
type SharedPrivateLinkResourcesClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewSharedPrivateLinkResourcesClient creates a new instance of SharedPrivateLinkResourcesClient with the specified values.
func NewSharedPrivateLinkResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *SharedPrivateLinkResourcesClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &SharedPrivateLinkResourcesClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginCreateOrUpdate - Initiates the creation or update of a shared private link resource managed by the search service in the given resource group.
// If the operation fails it returns the *CloudError error type.
func (client *SharedPrivateLinkResourcesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, sharedPrivateLinkResource SharedPrivateLinkResource, options *SearchManagementRequestOptions) (SharedPrivateLinkResourcesCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, searchServiceName, sharedPrivateLinkResourceName, sharedPrivateLinkResource, options)
	if err != nil {
		return SharedPrivateLinkResourcesCreateOrUpdatePollerResponse{}, err
	}
	result := SharedPrivateLinkResourcesCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("SharedPrivateLinkResourcesClient.CreateOrUpdate", "azure-async-operation", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return SharedPrivateLinkResourcesCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &SharedPrivateLinkResourcesCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Initiates the creation or update of a shared private link resource managed by the search service in the given resource group.
// If the operation fails it returns the *CloudError error type.
func (client *SharedPrivateLinkResourcesClient) createOrUpdate(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, sharedPrivateLinkResource SharedPrivateLinkResource, options *SearchManagementRequestOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, searchServiceName, sharedPrivateLinkResourceName, sharedPrivateLinkResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SharedPrivateLinkResourcesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, sharedPrivateLinkResource SharedPrivateLinkResource, options *SearchManagementRequestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Search/searchServices/{searchServiceName}/sharedPrivateLinkResources/{sharedPrivateLinkResourceName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if searchServiceName == "" {
		return nil, errors.New("parameter searchServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{searchServiceName}", url.PathEscape(searchServiceName))
	if sharedPrivateLinkResourceName == "" {
		return nil, errors.New("parameter sharedPrivateLinkResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sharedPrivateLinkResourceName}", url.PathEscape(sharedPrivateLinkResourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.ClientRequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.ClientRequestID)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, sharedPrivateLinkResource)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *SharedPrivateLinkResourcesClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginDelete - Initiates the deletion of the shared private link resource from the search service.
// If the operation fails it returns the *CloudError error type.
func (client *SharedPrivateLinkResourcesClient) BeginDelete(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, options *SearchManagementRequestOptions) (SharedPrivateLinkResourcesDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, searchServiceName, sharedPrivateLinkResourceName, options)
	if err != nil {
		return SharedPrivateLinkResourcesDeletePollerResponse{}, err
	}
	result := SharedPrivateLinkResourcesDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("SharedPrivateLinkResourcesClient.Delete", "azure-async-operation", resp, client.pl, client.deleteHandleError)
	if err != nil {
		return SharedPrivateLinkResourcesDeletePollerResponse{}, err
	}
	result.Poller = &SharedPrivateLinkResourcesDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Initiates the deletion of the shared private link resource from the search service.
// If the operation fails it returns the *CloudError error type.
func (client *SharedPrivateLinkResourcesClient) deleteOperation(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, options *SearchManagementRequestOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, searchServiceName, sharedPrivateLinkResourceName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted, http.StatusNoContent, http.StatusNotFound) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SharedPrivateLinkResourcesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, options *SearchManagementRequestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Search/searchServices/{searchServiceName}/sharedPrivateLinkResources/{sharedPrivateLinkResourceName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if searchServiceName == "" {
		return nil, errors.New("parameter searchServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{searchServiceName}", url.PathEscape(searchServiceName))
	if sharedPrivateLinkResourceName == "" {
		return nil, errors.New("parameter sharedPrivateLinkResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sharedPrivateLinkResourceName}", url.PathEscape(sharedPrivateLinkResourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.ClientRequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.ClientRequestID)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *SharedPrivateLinkResourcesClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Gets the details of the shared private link resource managed by the search service in the given resource group.
// If the operation fails it returns the *CloudError error type.
func (client *SharedPrivateLinkResourcesClient) Get(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, options *SearchManagementRequestOptions) (SharedPrivateLinkResourcesGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, searchServiceName, sharedPrivateLinkResourceName, options)
	if err != nil {
		return SharedPrivateLinkResourcesGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SharedPrivateLinkResourcesGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SharedPrivateLinkResourcesGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SharedPrivateLinkResourcesClient) getCreateRequest(ctx context.Context, resourceGroupName string, searchServiceName string, sharedPrivateLinkResourceName string, options *SearchManagementRequestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Search/searchServices/{searchServiceName}/sharedPrivateLinkResources/{sharedPrivateLinkResourceName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if searchServiceName == "" {
		return nil, errors.New("parameter searchServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{searchServiceName}", url.PathEscape(searchServiceName))
	if sharedPrivateLinkResourceName == "" {
		return nil, errors.New("parameter sharedPrivateLinkResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sharedPrivateLinkResourceName}", url.PathEscape(sharedPrivateLinkResourceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.ClientRequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.ClientRequestID)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SharedPrivateLinkResourcesClient) getHandleResponse(resp *http.Response) (SharedPrivateLinkResourcesGetResponse, error) {
	result := SharedPrivateLinkResourcesGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SharedPrivateLinkResource); err != nil {
		return SharedPrivateLinkResourcesGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *SharedPrivateLinkResourcesClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListByService - Gets a list of all shared private link resources managed by the given service.
// If the operation fails it returns the *CloudError error type.
func (client *SharedPrivateLinkResourcesClient) ListByService(resourceGroupName string, searchServiceName string, options *SearchManagementRequestOptions) *SharedPrivateLinkResourcesListByServicePager {
	return &SharedPrivateLinkResourcesListByServicePager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByServiceCreateRequest(ctx, resourceGroupName, searchServiceName, options)
		},
		advancer: func(ctx context.Context, resp SharedPrivateLinkResourcesListByServiceResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.SharedPrivateLinkResourceListResult.NextLink)
		},
	}
}

// listByServiceCreateRequest creates the ListByService request.
func (client *SharedPrivateLinkResourcesClient) listByServiceCreateRequest(ctx context.Context, resourceGroupName string, searchServiceName string, options *SearchManagementRequestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Search/searchServices/{searchServiceName}/sharedPrivateLinkResources"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if searchServiceName == "" {
		return nil, errors.New("parameter searchServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{searchServiceName}", url.PathEscape(searchServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-08-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.ClientRequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.ClientRequestID)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByServiceHandleResponse handles the ListByService response.
func (client *SharedPrivateLinkResourcesClient) listByServiceHandleResponse(resp *http.Response) (SharedPrivateLinkResourcesListByServiceResponse, error) {
	result := SharedPrivateLinkResourcesListByServiceResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SharedPrivateLinkResourceListResult); err != nil {
		return SharedPrivateLinkResourcesListByServiceResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listByServiceHandleError handles the ListByService error response.
func (client *SharedPrivateLinkResourcesClient) listByServiceHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

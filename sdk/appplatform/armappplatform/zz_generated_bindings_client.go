//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappplatform

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// BindingsClient contains the methods for the Bindings group.
// Don't use this type directly, use NewBindingsClient() instead.
type BindingsClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewBindingsClient creates a new instance of BindingsClient with the specified values.
func NewBindingsClient(con *arm.Connection, subscriptionID string) *BindingsClient {
	return &BindingsClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), subscriptionID: subscriptionID}
}

// BeginCreateOrUpdate - Create a new Binding or update an exiting Binding.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, bindingResource BindingResource, options *BindingsBeginCreateOrUpdateOptions) (BindingsCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, serviceName, appName, bindingName, bindingResource, options)
	if err != nil {
		return BindingsCreateOrUpdatePollerResponse{}, err
	}
	result := BindingsCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("BindingsClient.CreateOrUpdate", "azure-async-operation", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return BindingsCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &BindingsCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Create a new Binding or update an exiting Binding.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) createOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, bindingResource BindingResource, options *BindingsBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, appName, bindingName, bindingResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *BindingsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, bindingResource BindingResource, options *BindingsBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/bindings/{bindingName}"
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
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if bindingName == "" {
		return nil, errors.New("parameter bindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bindingName}", url.PathEscape(bindingName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, bindingResource)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *BindingsClient) createOrUpdateHandleError(resp *http.Response) error {
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

// BeginDelete - Operation to delete a Binding.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) BeginDelete(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, options *BindingsBeginDeleteOptions) (BindingsDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, serviceName, appName, bindingName, options)
	if err != nil {
		return BindingsDeletePollerResponse{}, err
	}
	result := BindingsDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("BindingsClient.Delete", "azure-async-operation", resp, client.pl, client.deleteHandleError)
	if err != nil {
		return BindingsDeletePollerResponse{}, err
	}
	result.Poller = &BindingsDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Operation to delete a Binding.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) deleteOperation(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, options *BindingsBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, appName, bindingName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *BindingsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, options *BindingsBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/bindings/{bindingName}"
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
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if bindingName == "" {
		return nil, errors.New("parameter bindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bindingName}", url.PathEscape(bindingName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *BindingsClient) deleteHandleError(resp *http.Response) error {
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

// Get - Get a Binding and its properties.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) Get(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, options *BindingsGetOptions) (BindingsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, appName, bindingName, options)
	if err != nil {
		return BindingsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return BindingsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return BindingsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *BindingsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, options *BindingsGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/bindings/{bindingName}"
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
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if bindingName == "" {
		return nil, errors.New("parameter bindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bindingName}", url.PathEscape(bindingName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *BindingsClient) getHandleResponse(resp *http.Response) (BindingsGetResponse, error) {
	result := BindingsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BindingResource); err != nil {
		return BindingsGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *BindingsClient) getHandleError(resp *http.Response) error {
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

// List - Handles requests to list all resources in an App.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) List(resourceGroupName string, serviceName string, appName string, options *BindingsListOptions) *BindingsListPager {
	return &BindingsListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, serviceName, appName, options)
		},
		advancer: func(ctx context.Context, resp BindingsListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.BindingResourceCollection.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *BindingsClient) listCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, options *BindingsListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/bindings"
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
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *BindingsClient) listHandleResponse(resp *http.Response) (BindingsListResponse, error) {
	result := BindingsListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BindingResourceCollection); err != nil {
		return BindingsListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *BindingsClient) listHandleError(resp *http.Response) error {
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

// BeginUpdate - Operation to update an exiting Binding.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) BeginUpdate(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, bindingResource BindingResource, options *BindingsBeginUpdateOptions) (BindingsUpdatePollerResponse, error) {
	resp, err := client.update(ctx, resourceGroupName, serviceName, appName, bindingName, bindingResource, options)
	if err != nil {
		return BindingsUpdatePollerResponse{}, err
	}
	result := BindingsUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("BindingsClient.Update", "azure-async-operation", resp, client.pl, client.updateHandleError)
	if err != nil {
		return BindingsUpdatePollerResponse{}, err
	}
	result.Poller = &BindingsUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// Update - Operation to update an exiting Binding.
// If the operation fails it returns the *CloudError error type.
func (client *BindingsClient) update(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, bindingResource BindingResource, options *BindingsBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serviceName, appName, bindingName, bindingResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *BindingsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, bindingName string, bindingResource BindingResource, options *BindingsBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/bindings/{bindingName}"
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
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if bindingName == "" {
		return nil, errors.New("parameter bindingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bindingName}", url.PathEscape(bindingName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, bindingResource)
}

// updateHandleError handles the Update error response.
func (client *BindingsClient) updateHandleError(resp *http.Response) error {
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

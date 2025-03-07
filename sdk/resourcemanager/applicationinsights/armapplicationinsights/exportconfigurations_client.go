// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapplicationinsights

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

// ExportConfigurationsClient contains the methods for the ExportConfigurations group.
// Don't use this type directly, use NewExportConfigurationsClient() instead.
type ExportConfigurationsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewExportConfigurationsClient creates a new instance of ExportConfigurationsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewExportConfigurationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ExportConfigurationsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ExportConfigurationsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Create - Create a Continuous Export configuration of an Application Insights component.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2015-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - exportProperties - Properties that need to be specified to create a Continuous Export configuration of a Application Insights
//     component.
//   - options - ExportConfigurationsClientCreateOptions contains the optional parameters for the ExportConfigurationsClient.Create
//     method.
func (client *ExportConfigurationsClient) Create(ctx context.Context, resourceGroupName string, resourceName string, exportProperties ComponentExportRequest, options *ExportConfigurationsClientCreateOptions) (ExportConfigurationsClientCreateResponse, error) {
	var err error
	const operationName = "ExportConfigurationsClient.Create"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createCreateRequest(ctx, resourceGroupName, resourceName, exportProperties, options)
	if err != nil {
		return ExportConfigurationsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExportConfigurationsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExportConfigurationsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *ExportConfigurationsClient) createCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, exportProperties ComponentExportRequest, _ *ExportConfigurationsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/exportconfiguration"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, exportProperties); err != nil {
		return nil, err
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *ExportConfigurationsClient) createHandleResponse(resp *http.Response) (ExportConfigurationsClientCreateResponse, error) {
	result := ExportConfigurationsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentExportConfigurationArray); err != nil {
		return ExportConfigurationsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a Continuous Export configuration of an Application Insights component.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2015-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - exportID - The Continuous Export configuration ID. This is unique within a Application Insights component.
//   - options - ExportConfigurationsClientDeleteOptions contains the optional parameters for the ExportConfigurationsClient.Delete
//     method.
func (client *ExportConfigurationsClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, exportID string, options *ExportConfigurationsClientDeleteOptions) (ExportConfigurationsClientDeleteResponse, error) {
	var err error
	const operationName = "ExportConfigurationsClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, exportID, options)
	if err != nil {
		return ExportConfigurationsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExportConfigurationsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExportConfigurationsClientDeleteResponse{}, err
	}
	resp, err := client.deleteHandleResponse(httpResp)
	return resp, err
}

// deleteCreateRequest creates the Delete request.
func (client *ExportConfigurationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, exportID string, _ *ExportConfigurationsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/exportconfiguration/{exportId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if exportID == "" {
		return nil, errors.New("parameter exportID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{exportId}", url.PathEscape(exportID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *ExportConfigurationsClient) deleteHandleResponse(resp *http.Response) (ExportConfigurationsClientDeleteResponse, error) {
	result := ExportConfigurationsClientDeleteResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentExportConfiguration); err != nil {
		return ExportConfigurationsClientDeleteResponse{}, err
	}
	return result, nil
}

// Get - Get the Continuous Export configuration for this export id.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2015-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - exportID - The Continuous Export configuration ID. This is unique within a Application Insights component.
//   - options - ExportConfigurationsClientGetOptions contains the optional parameters for the ExportConfigurationsClient.Get
//     method.
func (client *ExportConfigurationsClient) Get(ctx context.Context, resourceGroupName string, resourceName string, exportID string, options *ExportConfigurationsClientGetOptions) (ExportConfigurationsClientGetResponse, error) {
	var err error
	const operationName = "ExportConfigurationsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, exportID, options)
	if err != nil {
		return ExportConfigurationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExportConfigurationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExportConfigurationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ExportConfigurationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, exportID string, _ *ExportConfigurationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/exportconfiguration/{exportId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if exportID == "" {
		return nil, errors.New("parameter exportID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{exportId}", url.PathEscape(exportID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ExportConfigurationsClient) getHandleResponse(resp *http.Response) (ExportConfigurationsClientGetResponse, error) {
	result := ExportConfigurationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentExportConfiguration); err != nil {
		return ExportConfigurationsClientGetResponse{}, err
	}
	return result, nil
}

// List - Gets a list of Continuous Export configuration of an Application Insights component.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2015-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - options - ExportConfigurationsClientListOptions contains the optional parameters for the ExportConfigurationsClient.List
//     method.
func (client *ExportConfigurationsClient) List(ctx context.Context, resourceGroupName string, resourceName string, options *ExportConfigurationsClientListOptions) (ExportConfigurationsClientListResponse, error) {
	var err error
	const operationName = "ExportConfigurationsClient.List"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return ExportConfigurationsClientListResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExportConfigurationsClientListResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExportConfigurationsClientListResponse{}, err
	}
	resp, err := client.listHandleResponse(httpResp)
	return resp, err
}

// listCreateRequest creates the List request.
func (client *ExportConfigurationsClient) listCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, _ *ExportConfigurationsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/exportconfiguration"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ExportConfigurationsClient) listHandleResponse(resp *http.Response) (ExportConfigurationsClientListResponse, error) {
	result := ExportConfigurationsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentExportConfigurationArray); err != nil {
		return ExportConfigurationsClientListResponse{}, err
	}
	return result, nil
}

// Update - Update the Continuous Export configuration for this export id.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2015-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - exportID - The Continuous Export configuration ID. This is unique within a Application Insights component.
//   - exportProperties - Properties that need to be specified to update the Continuous Export configuration.
//   - options - ExportConfigurationsClientUpdateOptions contains the optional parameters for the ExportConfigurationsClient.Update
//     method.
func (client *ExportConfigurationsClient) Update(ctx context.Context, resourceGroupName string, resourceName string, exportID string, exportProperties ComponentExportRequest, options *ExportConfigurationsClientUpdateOptions) (ExportConfigurationsClientUpdateResponse, error) {
	var err error
	const operationName = "ExportConfigurationsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, exportID, exportProperties, options)
	if err != nil {
		return ExportConfigurationsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExportConfigurationsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExportConfigurationsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *ExportConfigurationsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, exportID string, exportProperties ComponentExportRequest, _ *ExportConfigurationsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/exportconfiguration/{exportId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if exportID == "" {
		return nil, errors.New("parameter exportID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{exportId}", url.PathEscape(exportID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, exportProperties); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *ExportConfigurationsClient) updateHandleResponse(resp *http.Response) (ExportConfigurationsClientUpdateResponse, error) {
	result := ExportConfigurationsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentExportConfiguration); err != nil {
		return ExportConfigurationsClientUpdateResponse{}, err
	}
	return result, nil
}

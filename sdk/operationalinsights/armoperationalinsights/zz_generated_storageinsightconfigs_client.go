//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armoperationalinsights

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// StorageInsightConfigsClient contains the methods for the StorageInsightConfigs group.
// Don't use this type directly, use NewStorageInsightConfigsClient() instead.
type StorageInsightConfigsClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewStorageInsightConfigsClient creates a new instance of StorageInsightConfigsClient with the specified values.
func NewStorageInsightConfigsClient(con *arm.Connection, subscriptionID string) *StorageInsightConfigsClient {
	return &StorageInsightConfigsClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), subscriptionID: subscriptionID}
}

// CreateOrUpdate - Create or update a storage insight.
// If the operation fails it returns a generic error.
func (client *StorageInsightConfigsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, storageInsightName string, parameters StorageInsight, options *StorageInsightConfigsCreateOrUpdateOptions) (StorageInsightConfigsCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, storageInsightName, parameters, options)
	if err != nil {
		return StorageInsightConfigsCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StorageInsightConfigsCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return StorageInsightConfigsCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *StorageInsightConfigsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, storageInsightName string, parameters StorageInsight, options *StorageInsightConfigsCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/storageInsightConfigs/{storageInsightName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if storageInsightName == "" {
		return nil, errors.New("parameter storageInsightName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageInsightName}", url.PathEscape(storageInsightName))
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
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *StorageInsightConfigsClient) createOrUpdateHandleResponse(resp *http.Response) (StorageInsightConfigsCreateOrUpdateResponse, error) {
	result := StorageInsightConfigsCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageInsight); err != nil {
		return StorageInsightConfigsCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *StorageInsightConfigsClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	if len(body) == 0 {
		return runtime.NewResponseError(errors.New(resp.Status), resp)
	}
	return runtime.NewResponseError(errors.New(string(body)), resp)
}

// Delete - Deletes a storageInsightsConfigs resource
// If the operation fails it returns a generic error.
func (client *StorageInsightConfigsClient) Delete(ctx context.Context, resourceGroupName string, workspaceName string, storageInsightName string, options *StorageInsightConfigsDeleteOptions) (StorageInsightConfigsDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, storageInsightName, options)
	if err != nil {
		return StorageInsightConfigsDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StorageInsightConfigsDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return StorageInsightConfigsDeleteResponse{}, client.deleteHandleError(resp)
	}
	return StorageInsightConfigsDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *StorageInsightConfigsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, storageInsightName string, options *StorageInsightConfigsDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/storageInsightConfigs/{storageInsightName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if storageInsightName == "" {
		return nil, errors.New("parameter storageInsightName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageInsightName}", url.PathEscape(storageInsightName))
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
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *StorageInsightConfigsClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	if len(body) == 0 {
		return runtime.NewResponseError(errors.New(resp.Status), resp)
	}
	return runtime.NewResponseError(errors.New(string(body)), resp)
}

// Get - Gets a storage insight instance.
// If the operation fails it returns a generic error.
func (client *StorageInsightConfigsClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, storageInsightName string, options *StorageInsightConfigsGetOptions) (StorageInsightConfigsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, storageInsightName, options)
	if err != nil {
		return StorageInsightConfigsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return StorageInsightConfigsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return StorageInsightConfigsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *StorageInsightConfigsClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, storageInsightName string, options *StorageInsightConfigsGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/storageInsightConfigs/{storageInsightName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if storageInsightName == "" {
		return nil, errors.New("parameter storageInsightName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageInsightName}", url.PathEscape(storageInsightName))
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
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *StorageInsightConfigsClient) getHandleResponse(resp *http.Response) (StorageInsightConfigsGetResponse, error) {
	result := StorageInsightConfigsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageInsight); err != nil {
		return StorageInsightConfigsGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *StorageInsightConfigsClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	if len(body) == 0 {
		return runtime.NewResponseError(errors.New(resp.Status), resp)
	}
	return runtime.NewResponseError(errors.New(string(body)), resp)
}

// ListByWorkspace - Lists the storage insight instances within a workspace
// If the operation fails it returns a generic error.
func (client *StorageInsightConfigsClient) ListByWorkspace(resourceGroupName string, workspaceName string, options *StorageInsightConfigsListByWorkspaceOptions) *StorageInsightConfigsListByWorkspacePager {
	return &StorageInsightConfigsListByWorkspacePager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByWorkspaceCreateRequest(ctx, resourceGroupName, workspaceName, options)
		},
		advancer: func(ctx context.Context, resp StorageInsightConfigsListByWorkspaceResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.StorageInsightListResult.ODataNextLink)
		},
	}
}

// listByWorkspaceCreateRequest creates the ListByWorkspace request.
func (client *StorageInsightConfigsClient) listByWorkspaceCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *StorageInsightConfigsListByWorkspaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/storageInsightConfigs"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
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
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByWorkspaceHandleResponse handles the ListByWorkspace response.
func (client *StorageInsightConfigsClient) listByWorkspaceHandleResponse(resp *http.Response) (StorageInsightConfigsListByWorkspaceResponse, error) {
	result := StorageInsightConfigsListByWorkspaceResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageInsightListResult); err != nil {
		return StorageInsightConfigsListByWorkspaceResponse{}, err
	}
	return result, nil
}

// listByWorkspaceHandleError handles the ListByWorkspace error response.
func (client *StorageInsightConfigsClient) listByWorkspaceHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	if len(body) == 0 {
		return runtime.NewResponseError(errors.New(resp.Status), resp)
	}
	return runtime.NewResponseError(errors.New(string(body)), resp)
}

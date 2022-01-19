//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurityinsight

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// SourceControlsClient contains the methods for the SourceControls group.
// Don't use this type directly, use NewSourceControlsClient() instead.
type SourceControlsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewSourceControlsClient creates a new instance of SourceControlsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewSourceControlsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *SourceControlsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &SourceControlsClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// Create - Creates a source control.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// sourceControlID - Source control Id
// sourceControl - The SourceControl
// options - SourceControlsClientCreateOptions contains the optional parameters for the SourceControlsClient.Create method.
func (client *SourceControlsClient) Create(ctx context.Context, resourceGroupName string, workspaceName string, sourceControlID string, sourceControl SourceControl, options *SourceControlsClientCreateOptions) (SourceControlsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, workspaceName, sourceControlID, sourceControl, options)
	if err != nil {
		return SourceControlsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SourceControlsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return SourceControlsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *SourceControlsClient) createCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sourceControlID string, sourceControl SourceControl, options *SourceControlsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/sourcecontrols/{sourceControlId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if sourceControlID == "" {
		return nil, errors.New("parameter sourceControlID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sourceControlId}", url.PathEscape(sourceControlID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, sourceControl)
}

// createHandleResponse handles the Create response.
func (client *SourceControlsClient) createHandleResponse(resp *http.Response) (SourceControlsClientCreateResponse, error) {
	result := SourceControlsClientCreateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SourceControl); err != nil {
		return SourceControlsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a source control.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// sourceControlID - Source control Id
// options - SourceControlsClientDeleteOptions contains the optional parameters for the SourceControlsClient.Delete method.
func (client *SourceControlsClient) Delete(ctx context.Context, resourceGroupName string, workspaceName string, sourceControlID string, options *SourceControlsClientDeleteOptions) (SourceControlsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, sourceControlID, options)
	if err != nil {
		return SourceControlsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SourceControlsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return SourceControlsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return SourceControlsClientDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SourceControlsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sourceControlID string, options *SourceControlsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/sourcecontrols/{sourceControlId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if sourceControlID == "" {
		return nil, errors.New("parameter sourceControlID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sourceControlId}", url.PathEscape(sourceControlID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets a source control byt its identifier.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// sourceControlID - Source control Id
// options - SourceControlsClientGetOptions contains the optional parameters for the SourceControlsClient.Get method.
func (client *SourceControlsClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, sourceControlID string, options *SourceControlsClientGetOptions) (SourceControlsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, sourceControlID, options)
	if err != nil {
		return SourceControlsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SourceControlsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SourceControlsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SourceControlsClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sourceControlID string, options *SourceControlsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/sourcecontrols/{sourceControlId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if sourceControlID == "" {
		return nil, errors.New("parameter sourceControlID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sourceControlId}", url.PathEscape(sourceControlID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SourceControlsClient) getHandleResponse(resp *http.Response) (SourceControlsClientGetResponse, error) {
	result := SourceControlsClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SourceControl); err != nil {
		return SourceControlsClientGetResponse{}, err
	}
	return result, nil
}

// List - Gets all source controls, without source control items.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// options - SourceControlsClientListOptions contains the optional parameters for the SourceControlsClient.List method.
func (client *SourceControlsClient) List(resourceGroupName string, workspaceName string, options *SourceControlsClientListOptions) *SourceControlsClientListPager {
	return &SourceControlsClientListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, workspaceName, options)
		},
		advancer: func(ctx context.Context, resp SourceControlsClientListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.SourceControlList.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *SourceControlsClient) listCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *SourceControlsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/sourcecontrols"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SourceControlsClient) listHandleResponse(resp *http.Response) (SourceControlsClientListResponse, error) {
	result := SourceControlsClientListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SourceControlList); err != nil {
		return SourceControlsClientListResponse{}, err
	}
	return result, nil
}

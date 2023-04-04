//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhealthcareapis

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

// WorkspacePrivateLinkResourcesClient contains the methods for the WorkspacePrivateLinkResources group.
// Don't use this type directly, use NewWorkspacePrivateLinkResourcesClient() instead.
type WorkspacePrivateLinkResourcesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewWorkspacePrivateLinkResourcesClient creates a new instance of WorkspacePrivateLinkResourcesClient with the specified values.
//   - subscriptionID - The subscription identifier.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkspacePrivateLinkResourcesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkspacePrivateLinkResourcesClient, error) {
	cl, err := arm.NewClient(moduleName+".WorkspacePrivateLinkResourcesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkspacePrivateLinkResourcesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Gets a private link resource that need to be created for a workspace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-31-preview
//   - resourceGroupName - The name of the resource group that contains the service instance.
//   - workspaceName - The name of workspace resource.
//   - groupName - The name of the private link resource group.
//   - options - WorkspacePrivateLinkResourcesClientGetOptions contains the optional parameters for the WorkspacePrivateLinkResourcesClient.Get
//     method.
func (client *WorkspacePrivateLinkResourcesClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, groupName string, options *WorkspacePrivateLinkResourcesClientGetOptions) (WorkspacePrivateLinkResourcesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, groupName, options)
	if err != nil {
		return WorkspacePrivateLinkResourcesClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkspacePrivateLinkResourcesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WorkspacePrivateLinkResourcesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *WorkspacePrivateLinkResourcesClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, groupName string, options *WorkspacePrivateLinkResourcesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/privateLinkResources/{groupName}"
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
	if groupName == "" {
		return nil, errors.New("parameter groupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{groupName}", url.PathEscape(groupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WorkspacePrivateLinkResourcesClient) getHandleResponse(resp *http.Response) (WorkspacePrivateLinkResourcesClientGetResponse, error) {
	result := WorkspacePrivateLinkResourcesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateLinkResourceDescription); err != nil {
		return WorkspacePrivateLinkResourcesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByWorkspacePager - Gets the private link resources that need to be created for a workspace.
//
// Generated from API version 2022-01-31-preview
//   - resourceGroupName - The name of the resource group that contains the service instance.
//   - workspaceName - The name of workspace resource.
//   - options - WorkspacePrivateLinkResourcesClientListByWorkspaceOptions contains the optional parameters for the WorkspacePrivateLinkResourcesClient.NewListByWorkspacePager
//     method.
func (client *WorkspacePrivateLinkResourcesClient) NewListByWorkspacePager(resourceGroupName string, workspaceName string, options *WorkspacePrivateLinkResourcesClientListByWorkspaceOptions) *runtime.Pager[WorkspacePrivateLinkResourcesClientListByWorkspaceResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkspacePrivateLinkResourcesClientListByWorkspaceResponse]{
		More: func(page WorkspacePrivateLinkResourcesClientListByWorkspaceResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *WorkspacePrivateLinkResourcesClientListByWorkspaceResponse) (WorkspacePrivateLinkResourcesClientListByWorkspaceResponse, error) {
			req, err := client.listByWorkspaceCreateRequest(ctx, resourceGroupName, workspaceName, options)
			if err != nil {
				return WorkspacePrivateLinkResourcesClientListByWorkspaceResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkspacePrivateLinkResourcesClientListByWorkspaceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkspacePrivateLinkResourcesClientListByWorkspaceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByWorkspaceHandleResponse(resp)
		},
	})
}

// listByWorkspaceCreateRequest creates the ListByWorkspace request.
func (client *WorkspacePrivateLinkResourcesClient) listByWorkspaceCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspacePrivateLinkResourcesClientListByWorkspaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HealthcareApis/workspaces/{workspaceName}/privateLinkResources"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-31-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByWorkspaceHandleResponse handles the ListByWorkspace response.
func (client *WorkspacePrivateLinkResourcesClient) listByWorkspaceHandleResponse(resp *http.Response) (WorkspacePrivateLinkResourcesClientListByWorkspaceResponse, error) {
	result := WorkspacePrivateLinkResourcesClientListByWorkspaceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PrivateLinkResourceListResultDescription); err != nil {
		return WorkspacePrivateLinkResourcesClientListByWorkspaceResponse{}, err
	}
	return result, nil
}

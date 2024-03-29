//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsynapse

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

// WorkspaceManagedSQLServerRecoverableSQLPoolsClient contains the methods for the WorkspaceManagedSQLServerRecoverableSQLPools group.
// Don't use this type directly, use NewWorkspaceManagedSQLServerRecoverableSQLPoolsClient() instead.
type WorkspaceManagedSQLServerRecoverableSQLPoolsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewWorkspaceManagedSQLServerRecoverableSQLPoolsClient creates a new instance of WorkspaceManagedSQLServerRecoverableSQLPoolsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkspaceManagedSQLServerRecoverableSQLPoolsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkspaceManagedSQLServerRecoverableSQLPoolsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkspaceManagedSQLServerRecoverableSQLPoolsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Get recoverable sql pools for workspace managed sql server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - The name of the sql pool
//   - options - WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetOptions contains the optional parameters for the WorkspaceManagedSQLServerRecoverableSQLPoolsClient.Get
//     method.
func (client *WorkspaceManagedSQLServerRecoverableSQLPoolsClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, options *WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetOptions) (WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse, error) {
	var err error
	const operationName = "WorkspaceManagedSQLServerRecoverableSQLPoolsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, options)
	if err != nil {
		return WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *WorkspaceManagedSQLServerRecoverableSQLPoolsClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, options *WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/recoverableSqlPools/{sqlPoolName}"
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
	if sqlPoolName == "" {
		return nil, errors.New("parameter sqlPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlPoolName}", url.PathEscape(sqlPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WorkspaceManagedSQLServerRecoverableSQLPoolsClient) getHandleResponse(resp *http.Response) (WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse, error) {
	result := WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RecoverableSQLPool); err != nil {
		return WorkspaceManagedSQLServerRecoverableSQLPoolsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Get list of recoverable sql pools for workspace managed sql server.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - options - WorkspaceManagedSQLServerRecoverableSQLPoolsClientListOptions contains the optional parameters for the WorkspaceManagedSQLServerRecoverableSQLPoolsClient.NewListPager
//     method.
func (client *WorkspaceManagedSQLServerRecoverableSQLPoolsClient) NewListPager(resourceGroupName string, workspaceName string, options *WorkspaceManagedSQLServerRecoverableSQLPoolsClientListOptions) *runtime.Pager[WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse]{
		More: func(page WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse) (WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "WorkspaceManagedSQLServerRecoverableSQLPoolsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, workspaceName, options)
			}, nil)
			if err != nil {
				return WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *WorkspaceManagedSQLServerRecoverableSQLPoolsClient) listCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *WorkspaceManagedSQLServerRecoverableSQLPoolsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/recoverableSqlPools"
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
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *WorkspaceManagedSQLServerRecoverableSQLPoolsClient) listHandleResponse(resp *http.Response) (WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse, error) {
	result := WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RecoverableSQLPoolListResult); err != nil {
		return WorkspaceManagedSQLServerRecoverableSQLPoolsClientListResponse{}, err
	}
	return result, nil
}

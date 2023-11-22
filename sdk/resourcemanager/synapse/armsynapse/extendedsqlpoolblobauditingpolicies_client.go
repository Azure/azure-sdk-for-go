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

// ExtendedSQLPoolBlobAuditingPoliciesClient contains the methods for the ExtendedSQLPoolBlobAuditingPolicies group.
// Don't use this type directly, use NewExtendedSQLPoolBlobAuditingPoliciesClient() instead.
type ExtendedSQLPoolBlobAuditingPoliciesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewExtendedSQLPoolBlobAuditingPoliciesClient creates a new instance of ExtendedSQLPoolBlobAuditingPoliciesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewExtendedSQLPoolBlobAuditingPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ExtendedSQLPoolBlobAuditingPoliciesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ExtendedSQLPoolBlobAuditingPoliciesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates an extended Sql pool's blob auditing policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - parameters - The extended Sql pool blob auditing policy.
//   - options - ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateOptions contains the optional parameters for the ExtendedSQLPoolBlobAuditingPoliciesClient.CreateOrUpdate
//     method.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, parameters ExtendedSQLPoolBlobAuditingPolicy, options *ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateOptions) (ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "ExtendedSQLPoolBlobAuditingPoliciesClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, parameters, options)
	if err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, parameters ExtendedSQLPoolBlobAuditingPolicy, options *ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/extendedAuditingSettings/{blobAuditingPolicyName}"
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
	urlPath = strings.ReplaceAll(urlPath, "{blobAuditingPolicyName}", url.PathEscape("default"))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) createOrUpdateHandleResponse(resp *http.Response) (ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse, error) {
	result := ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExtendedSQLPoolBlobAuditingPolicy); err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Get - Gets an extended Sql pool's blob auditing policy.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - options - ExtendedSQLPoolBlobAuditingPoliciesClientGetOptions contains the optional parameters for the ExtendedSQLPoolBlobAuditingPoliciesClient.Get
//     method.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, options *ExtendedSQLPoolBlobAuditingPoliciesClientGetOptions) (ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse, error) {
	var err error
	const operationName = "ExtendedSQLPoolBlobAuditingPoliciesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, options)
	if err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, options *ExtendedSQLPoolBlobAuditingPoliciesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/extendedAuditingSettings/{blobAuditingPolicyName}"
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
	urlPath = strings.ReplaceAll(urlPath, "{blobAuditingPolicyName}", url.PathEscape("default"))
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
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) getHandleResponse(resp *http.Response) (ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse, error) {
	result := ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExtendedSQLPoolBlobAuditingPolicy); err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientGetResponse{}, err
	}
	return result, nil
}

// NewListBySQLPoolPager - Lists extended auditing settings of a Sql pool.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - options - ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolOptions contains the optional parameters for the ExtendedSQLPoolBlobAuditingPoliciesClient.NewListBySQLPoolPager
//     method.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) NewListBySQLPoolPager(resourceGroupName string, workspaceName string, sqlPoolName string, options *ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolOptions) *runtime.Pager[ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse] {
	return runtime.NewPager(runtime.PagingHandler[ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse]{
		More: func(page ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse) (ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ExtendedSQLPoolBlobAuditingPoliciesClient.NewListBySQLPoolPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySQLPoolCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, options)
			}, nil)
			if err != nil {
				return ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse{}, err
			}
			return client.listBySQLPoolHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySQLPoolCreateRequest creates the ListBySQLPool request.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) listBySQLPoolCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, options *ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/extendedAuditingSettings"
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

// listBySQLPoolHandleResponse handles the ListBySQLPool response.
func (client *ExtendedSQLPoolBlobAuditingPoliciesClient) listBySQLPoolHandleResponse(resp *http.Response) (ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse, error) {
	result := ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExtendedSQLPoolBlobAuditingPolicyListResult); err != nil {
		return ExtendedSQLPoolBlobAuditingPoliciesClientListBySQLPoolResponse{}, err
	}
	return result, nil
}

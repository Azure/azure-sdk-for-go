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

// SQLPoolWorkloadClassifierClient contains the methods for the SQLPoolWorkloadClassifier group.
// Don't use this type directly, use NewSQLPoolWorkloadClassifierClient() instead.
type SQLPoolWorkloadClassifierClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSQLPoolWorkloadClassifierClient creates a new instance of SQLPoolWorkloadClassifierClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSQLPoolWorkloadClassifierClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLPoolWorkloadClassifierClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SQLPoolWorkloadClassifierClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create Or Update workload classifier for a Sql pool's workload group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - workloadGroupName - The name of the workload group.
//   - workloadClassifierName - The name of the workload classifier.
//   - parameters - The properties of the workload classifier.
//   - options - SQLPoolWorkloadClassifierClientBeginCreateOrUpdateOptions contains the optional parameters for the SQLPoolWorkloadClassifierClient.BeginCreateOrUpdate
//     method.
func (client *SQLPoolWorkloadClassifierClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, parameters WorkloadClassifier, options *SQLPoolWorkloadClassifierClientBeginCreateOrUpdateOptions) (*runtime.Poller[SQLPoolWorkloadClassifierClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, workspaceName, sqlPoolName, workloadGroupName, workloadClassifierName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SQLPoolWorkloadClassifierClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SQLPoolWorkloadClassifierClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create Or Update workload classifier for a Sql pool's workload group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
func (client *SQLPoolWorkloadClassifierClient) createOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, parameters WorkloadClassifier, options *SQLPoolWorkloadClassifierClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "SQLPoolWorkloadClassifierClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, workloadGroupName, workloadClassifierName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SQLPoolWorkloadClassifierClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, parameters WorkloadClassifier, options *SQLPoolWorkloadClassifierClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/workloadGroups/{workloadGroupName}/workloadClassifiers/{workloadClassifierName}"
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
	if workloadGroupName == "" {
		return nil, errors.New("parameter workloadGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadGroupName}", url.PathEscape(workloadGroupName))
	if workloadClassifierName == "" {
		return nil, errors.New("parameter workloadClassifierName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadClassifierName}", url.PathEscape(workloadClassifierName))
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

// BeginDelete - Remove workload classifier of a Sql pool's workload group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - workloadGroupName - The name of the workload group.
//   - workloadClassifierName - The name of the workload classifier.
//   - options - SQLPoolWorkloadClassifierClientBeginDeleteOptions contains the optional parameters for the SQLPoolWorkloadClassifierClient.BeginDelete
//     method.
func (client *SQLPoolWorkloadClassifierClient) BeginDelete(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, options *SQLPoolWorkloadClassifierClientBeginDeleteOptions) (*runtime.Poller[SQLPoolWorkloadClassifierClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, workspaceName, sqlPoolName, workloadGroupName, workloadClassifierName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SQLPoolWorkloadClassifierClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SQLPoolWorkloadClassifierClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Remove workload classifier of a Sql pool's workload group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
func (client *SQLPoolWorkloadClassifierClient) deleteOperation(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, options *SQLPoolWorkloadClassifierClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "SQLPoolWorkloadClassifierClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, workloadGroupName, workloadClassifierName, options)
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
func (client *SQLPoolWorkloadClassifierClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, options *SQLPoolWorkloadClassifierClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/workloadGroups/{workloadGroupName}/workloadClassifiers/{workloadClassifierName}"
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
	if workloadGroupName == "" {
		return nil, errors.New("parameter workloadGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadGroupName}", url.PathEscape(workloadGroupName))
	if workloadClassifierName == "" {
		return nil, errors.New("parameter workloadClassifierName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadClassifierName}", url.PathEscape(workloadClassifierName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Get a workload classifier of Sql pool's workload group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - workloadGroupName - The name of the workload group.
//   - workloadClassifierName - The name of the workload classifier.
//   - options - SQLPoolWorkloadClassifierClientGetOptions contains the optional parameters for the SQLPoolWorkloadClassifierClient.Get
//     method.
func (client *SQLPoolWorkloadClassifierClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, options *SQLPoolWorkloadClassifierClientGetOptions) (SQLPoolWorkloadClassifierClientGetResponse, error) {
	var err error
	const operationName = "SQLPoolWorkloadClassifierClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, workloadGroupName, workloadClassifierName, options)
	if err != nil {
		return SQLPoolWorkloadClassifierClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SQLPoolWorkloadClassifierClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SQLPoolWorkloadClassifierClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *SQLPoolWorkloadClassifierClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, workloadClassifierName string, options *SQLPoolWorkloadClassifierClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/workloadGroups/{workloadGroupName}/workloadClassifiers/{workloadClassifierName}"
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
	if workloadGroupName == "" {
		return nil, errors.New("parameter workloadGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadGroupName}", url.PathEscape(workloadGroupName))
	if workloadClassifierName == "" {
		return nil, errors.New("parameter workloadClassifierName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadClassifierName}", url.PathEscape(workloadClassifierName))
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
func (client *SQLPoolWorkloadClassifierClient) getHandleResponse(resp *http.Response) (SQLPoolWorkloadClassifierClientGetResponse, error) {
	result := SQLPoolWorkloadClassifierClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkloadClassifier); err != nil {
		return SQLPoolWorkloadClassifierClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Get list of Sql pool's workload classifier for workload groups.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - workloadGroupName - The name of the workload group.
//   - options - SQLPoolWorkloadClassifierClientListOptions contains the optional parameters for the SQLPoolWorkloadClassifierClient.NewListPager
//     method.
func (client *SQLPoolWorkloadClassifierClient) NewListPager(resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, options *SQLPoolWorkloadClassifierClientListOptions) *runtime.Pager[SQLPoolWorkloadClassifierClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SQLPoolWorkloadClassifierClientListResponse]{
		More: func(page SQLPoolWorkloadClassifierClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SQLPoolWorkloadClassifierClientListResponse) (SQLPoolWorkloadClassifierClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "SQLPoolWorkloadClassifierClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, workloadGroupName, options)
			}, nil)
			if err != nil {
				return SQLPoolWorkloadClassifierClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *SQLPoolWorkloadClassifierClient) listCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, workloadGroupName string, options *SQLPoolWorkloadClassifierClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/workloadGroups/{workloadGroupName}/workloadClassifiers"
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
	if workloadGroupName == "" {
		return nil, errors.New("parameter workloadGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workloadGroupName}", url.PathEscape(workloadGroupName))
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
func (client *SQLPoolWorkloadClassifierClient) listHandleResponse(resp *http.Response) (SQLPoolWorkloadClassifierClientListResponse, error) {
	result := SQLPoolWorkloadClassifierClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkloadClassifierListResult); err != nil {
		return SQLPoolWorkloadClassifierClientListResponse{}, err
	}
	return result, nil
}

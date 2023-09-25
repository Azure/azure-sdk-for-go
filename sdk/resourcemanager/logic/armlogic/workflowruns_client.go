//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlogic

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// WorkflowRunsClient contains the methods for the WorkflowRuns group.
// Don't use this type directly, use NewWorkflowRunsClient() instead.
type WorkflowRunsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewWorkflowRunsClient creates a new instance of WorkflowRunsClient with the specified values.
//   - subscriptionID - The subscription id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkflowRunsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkflowRunsClient, error) {
	cl, err := arm.NewClient(moduleName+".WorkflowRunsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkflowRunsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Cancel - Cancels a workflow run.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - options - WorkflowRunsClientCancelOptions contains the optional parameters for the WorkflowRunsClient.Cancel method.
func (client *WorkflowRunsClient) Cancel(ctx context.Context, resourceGroupName string, workflowName string, runName string, options *WorkflowRunsClientCancelOptions) (WorkflowRunsClientCancelResponse, error) {
	var err error
	req, err := client.cancelCreateRequest(ctx, resourceGroupName, workflowName, runName, options)
	if err != nil {
		return WorkflowRunsClientCancelResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowRunsClientCancelResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkflowRunsClientCancelResponse{}, err
	}
	return WorkflowRunsClientCancelResponse{}, nil
}

// cancelCreateRequest creates the Cancel request.
func (client *WorkflowRunsClient) cancelCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, options *WorkflowRunsClientCancelOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}/cancel"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if runName == "" {
		return nil, errors.New("parameter runName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runName}", url.PathEscape(runName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a workflow run.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - options - WorkflowRunsClientGetOptions contains the optional parameters for the WorkflowRunsClient.Get method.
func (client *WorkflowRunsClient) Get(ctx context.Context, resourceGroupName string, workflowName string, runName string, options *WorkflowRunsClientGetOptions) (WorkflowRunsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, workflowName, runName, options)
	if err != nil {
		return WorkflowRunsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowRunsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkflowRunsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *WorkflowRunsClient) getCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, options *WorkflowRunsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if runName == "" {
		return nil, errors.New("parameter runName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runName}", url.PathEscape(runName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WorkflowRunsClient) getHandleResponse(resp *http.Response) (WorkflowRunsClientGetResponse, error) {
	result := WorkflowRunsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowRun); err != nil {
		return WorkflowRunsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of workflow runs.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - options - WorkflowRunsClientListOptions contains the optional parameters for the WorkflowRunsClient.NewListPager method.
func (client *WorkflowRunsClient) NewListPager(resourceGroupName string, workflowName string, options *WorkflowRunsClientListOptions) (*runtime.Pager[WorkflowRunsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[WorkflowRunsClientListResponse]{
		More: func(page WorkflowRunsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkflowRunsClientListResponse) (WorkflowRunsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, workflowName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WorkflowRunsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkflowRunsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkflowRunsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *WorkflowRunsClient) listCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, options *WorkflowRunsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *WorkflowRunsClient) listHandleResponse(resp *http.Response) (WorkflowRunsClientListResponse, error) {
	result := WorkflowRunsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowRunListResult); err != nil {
		return WorkflowRunsClientListResponse{}, err
	}
	return result, nil
}


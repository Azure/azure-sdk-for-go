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
	"strings"
)

// WorkflowRunActionRequestHistoriesClient contains the methods for the WorkflowRunActionRequestHistories group.
// Don't use this type directly, use NewWorkflowRunActionRequestHistoriesClient() instead.
type WorkflowRunActionRequestHistoriesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewWorkflowRunActionRequestHistoriesClient creates a new instance of WorkflowRunActionRequestHistoriesClient with the specified values.
//   - subscriptionID - The subscription id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkflowRunActionRequestHistoriesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkflowRunActionRequestHistoriesClient, error) {
	cl, err := arm.NewClient(moduleName+".WorkflowRunActionRequestHistoriesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkflowRunActionRequestHistoriesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Gets a workflow run request history.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - actionName - The workflow action name.
//   - requestHistoryName - The request history name.
//   - options - WorkflowRunActionRequestHistoriesClientGetOptions contains the optional parameters for the WorkflowRunActionRequestHistoriesClient.Get
//     method.
func (client *WorkflowRunActionRequestHistoriesClient) Get(ctx context.Context, resourceGroupName string, workflowName string, runName string, actionName string, requestHistoryName string, options *WorkflowRunActionRequestHistoriesClientGetOptions) (WorkflowRunActionRequestHistoriesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, workflowName, runName, actionName, requestHistoryName, options)
	if err != nil {
		return WorkflowRunActionRequestHistoriesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowRunActionRequestHistoriesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkflowRunActionRequestHistoriesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *WorkflowRunActionRequestHistoriesClient) getCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, actionName string, requestHistoryName string, options *WorkflowRunActionRequestHistoriesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}/actions/{actionName}/requestHistories/{requestHistoryName}"
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
	if actionName == "" {
		return nil, errors.New("parameter actionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionName}", url.PathEscape(actionName))
	if requestHistoryName == "" {
		return nil, errors.New("parameter requestHistoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{requestHistoryName}", url.PathEscape(requestHistoryName))
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
func (client *WorkflowRunActionRequestHistoriesClient) getHandleResponse(resp *http.Response) (WorkflowRunActionRequestHistoriesClientGetResponse, error) {
	result := WorkflowRunActionRequestHistoriesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RequestHistory); err != nil {
		return WorkflowRunActionRequestHistoriesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List a workflow run request history.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - actionName - The workflow action name.
//   - options - WorkflowRunActionRequestHistoriesClientListOptions contains the optional parameters for the WorkflowRunActionRequestHistoriesClient.NewListPager
//     method.
func (client *WorkflowRunActionRequestHistoriesClient) NewListPager(resourceGroupName string, workflowName string, runName string, actionName string, options *WorkflowRunActionRequestHistoriesClientListOptions) (*runtime.Pager[WorkflowRunActionRequestHistoriesClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[WorkflowRunActionRequestHistoriesClientListResponse]{
		More: func(page WorkflowRunActionRequestHistoriesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkflowRunActionRequestHistoriesClientListResponse) (WorkflowRunActionRequestHistoriesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, workflowName, runName, actionName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WorkflowRunActionRequestHistoriesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkflowRunActionRequestHistoriesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkflowRunActionRequestHistoriesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *WorkflowRunActionRequestHistoriesClient) listCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, actionName string, options *WorkflowRunActionRequestHistoriesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}/actions/{actionName}/requestHistories"
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
	if actionName == "" {
		return nil, errors.New("parameter actionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionName}", url.PathEscape(actionName))
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

// listHandleResponse handles the List response.
func (client *WorkflowRunActionRequestHistoriesClient) listHandleResponse(resp *http.Response) (WorkflowRunActionRequestHistoriesClientListResponse, error) {
	result := WorkflowRunActionRequestHistoriesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RequestHistoryListResult); err != nil {
		return WorkflowRunActionRequestHistoriesClientListResponse{}, err
	}
	return result, nil
}


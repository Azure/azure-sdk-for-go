//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armappservice

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

// WorkflowTriggersClient contains the methods for the WorkflowTriggers group.
// Don't use this type directly, use NewWorkflowTriggersClient() instead.
type WorkflowTriggersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewWorkflowTriggersClient creates a new instance of WorkflowTriggersClient with the specified values.
//   - subscriptionID - Your Azure subscription ID. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000).
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkflowTriggersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkflowTriggersClient, error) {
	cl, err := arm.NewClient(moduleName+".WorkflowTriggersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkflowTriggersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Gets a workflow trigger.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01
//   - resourceGroupName - Name of the resource group to which the resource belongs.
//   - name - Site name.
//   - workflowName - The workflow name.
//   - triggerName - The workflow trigger name.
//   - options - WorkflowTriggersClientGetOptions contains the optional parameters for the WorkflowTriggersClient.Get method.
func (client *WorkflowTriggersClient) Get(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientGetOptions) (WorkflowTriggersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, name, workflowName, triggerName, options)
	if err != nil {
		return WorkflowTriggersClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowTriggersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WorkflowTriggersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *WorkflowTriggersClient) getCreateRequest(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/hostruntime/runtime/webhooks/workflow/api/management/workflows/{workflowName}/triggers/{triggerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if triggerName == "" {
		return nil, errors.New("parameter triggerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{triggerName}", url.PathEscape(triggerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WorkflowTriggersClient) getHandleResponse(resp *http.Response) (WorkflowTriggersClientGetResponse, error) {
	result := WorkflowTriggersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowTrigger); err != nil {
		return WorkflowTriggersClientGetResponse{}, err
	}
	return result, nil
}

// GetSchemaJSON - Get the trigger schema as JSON.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01
//   - resourceGroupName - Name of the resource group to which the resource belongs.
//   - name - Site name.
//   - workflowName - The workflow name.
//   - triggerName - The workflow trigger name.
//   - options - WorkflowTriggersClientGetSchemaJSONOptions contains the optional parameters for the WorkflowTriggersClient.GetSchemaJSON
//     method.
func (client *WorkflowTriggersClient) GetSchemaJSON(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientGetSchemaJSONOptions) (WorkflowTriggersClientGetSchemaJSONResponse, error) {
	req, err := client.getSchemaJSONCreateRequest(ctx, resourceGroupName, name, workflowName, triggerName, options)
	if err != nil {
		return WorkflowTriggersClientGetSchemaJSONResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowTriggersClientGetSchemaJSONResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WorkflowTriggersClientGetSchemaJSONResponse{}, runtime.NewResponseError(resp)
	}
	return client.getSchemaJSONHandleResponse(resp)
}

// getSchemaJSONCreateRequest creates the GetSchemaJSON request.
func (client *WorkflowTriggersClient) getSchemaJSONCreateRequest(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientGetSchemaJSONOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/hostruntime/runtime/webhooks/workflow/api/management/workflows/{workflowName}/triggers/{triggerName}/schemas/json"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if triggerName == "" {
		return nil, errors.New("parameter triggerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{triggerName}", url.PathEscape(triggerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getSchemaJSONHandleResponse handles the GetSchemaJSON response.
func (client *WorkflowTriggersClient) getSchemaJSONHandleResponse(resp *http.Response) (WorkflowTriggersClientGetSchemaJSONResponse, error) {
	result := WorkflowTriggersClientGetSchemaJSONResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.JSONSchema); err != nil {
		return WorkflowTriggersClientGetSchemaJSONResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of workflow triggers.
//
// Generated from API version 2022-09-01
//   - resourceGroupName - Name of the resource group to which the resource belongs.
//   - name - Site name.
//   - workflowName - The workflow name.
//   - options - WorkflowTriggersClientListOptions contains the optional parameters for the WorkflowTriggersClient.NewListPager
//     method.
func (client *WorkflowTriggersClient) NewListPager(resourceGroupName string, name string, workflowName string, options *WorkflowTriggersClientListOptions) *runtime.Pager[WorkflowTriggersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkflowTriggersClientListResponse]{
		More: func(page WorkflowTriggersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkflowTriggersClientListResponse) (WorkflowTriggersClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, name, workflowName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WorkflowTriggersClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkflowTriggersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkflowTriggersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *WorkflowTriggersClient) listCreateRequest(ctx context.Context, resourceGroupName string, name string, workflowName string, options *WorkflowTriggersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/hostruntime/runtime/webhooks/workflow/api/management/workflows/{workflowName}/triggers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01")
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
func (client *WorkflowTriggersClient) listHandleResponse(resp *http.Response) (WorkflowTriggersClientListResponse, error) {
	result := WorkflowTriggersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowTriggerListResult); err != nil {
		return WorkflowTriggersClientListResponse{}, err
	}
	return result, nil
}

// ListCallbackURL - Get the callback URL for a workflow trigger.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01
//   - resourceGroupName - Name of the resource group to which the resource belongs.
//   - name - Site name.
//   - workflowName - The workflow name.
//   - triggerName - The workflow trigger name.
//   - options - WorkflowTriggersClientListCallbackURLOptions contains the optional parameters for the WorkflowTriggersClient.ListCallbackURL
//     method.
func (client *WorkflowTriggersClient) ListCallbackURL(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientListCallbackURLOptions) (WorkflowTriggersClientListCallbackURLResponse, error) {
	req, err := client.listCallbackURLCreateRequest(ctx, resourceGroupName, name, workflowName, triggerName, options)
	if err != nil {
		return WorkflowTriggersClientListCallbackURLResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowTriggersClientListCallbackURLResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WorkflowTriggersClientListCallbackURLResponse{}, runtime.NewResponseError(resp)
	}
	return client.listCallbackURLHandleResponse(resp)
}

// listCallbackURLCreateRequest creates the ListCallbackURL request.
func (client *WorkflowTriggersClient) listCallbackURLCreateRequest(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientListCallbackURLOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/hostruntime/runtime/webhooks/workflow/api/management/workflows/{workflowName}/triggers/{triggerName}/listCallbackUrl"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if triggerName == "" {
		return nil, errors.New("parameter triggerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{triggerName}", url.PathEscape(triggerName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listCallbackURLHandleResponse handles the ListCallbackURL response.
func (client *WorkflowTriggersClient) listCallbackURLHandleResponse(resp *http.Response) (WorkflowTriggersClientListCallbackURLResponse, error) {
	result := WorkflowTriggersClientListCallbackURLResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowTriggerCallbackURL); err != nil {
		return WorkflowTriggersClientListCallbackURLResponse{}, err
	}
	return result, nil
}

// BeginRun - Runs a workflow trigger.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01
//   - resourceGroupName - Name of the resource group to which the resource belongs.
//   - name - Site name.
//   - workflowName - The workflow name.
//   - triggerName - The workflow trigger name.
//   - options - WorkflowTriggersClientBeginRunOptions contains the optional parameters for the WorkflowTriggersClient.BeginRun
//     method.
func (client *WorkflowTriggersClient) BeginRun(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientBeginRunOptions) (*runtime.Poller[WorkflowTriggersClientRunResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.run(ctx, resourceGroupName, name, workflowName, triggerName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[WorkflowTriggersClientRunResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[WorkflowTriggersClientRunResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Run - Runs a workflow trigger.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-01
func (client *WorkflowTriggersClient) run(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientBeginRunOptions) (*http.Response, error) {
	req, err := client.runCreateRequest(ctx, resourceGroupName, name, workflowName, triggerName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// runCreateRequest creates the Run request.
func (client *WorkflowTriggersClient) runCreateRequest(ctx context.Context, resourceGroupName string, name string, workflowName string, triggerName string, options *WorkflowTriggersClientBeginRunOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/hostruntime/runtime/webhooks/workflow/api/management/workflows/{workflowName}/triggers/{triggerName}/run"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if triggerName == "" {
		return nil, errors.New("parameter triggerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{triggerName}", url.PathEscape(triggerName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

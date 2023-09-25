//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstreamanalytics

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

// FunctionsClient contains the methods for the Functions group.
// Don't use this type directly, use NewFunctionsClient() instead.
type FunctionsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewFunctionsClient creates a new instance of FunctionsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewFunctionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FunctionsClient, error) {
	cl, err := arm.NewClient(moduleName+".FunctionsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &FunctionsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrReplace - Creates a function or replaces an already existing function under an existing streaming job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - functionName - The name of the function.
//   - function - The definition of the function that will be used to create a new function or replace the existing one under
//     the streaming job.
//   - options - FunctionsClientCreateOrReplaceOptions contains the optional parameters for the FunctionsClient.CreateOrReplace
//     method.
func (client *FunctionsClient) CreateOrReplace(ctx context.Context, resourceGroupName string, jobName string, functionName string, function Function, options *FunctionsClientCreateOrReplaceOptions) (FunctionsClientCreateOrReplaceResponse, error) {
	var err error
	req, err := client.createOrReplaceCreateRequest(ctx, resourceGroupName, jobName, functionName, function, options)
	if err != nil {
		return FunctionsClientCreateOrReplaceResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FunctionsClientCreateOrReplaceResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return FunctionsClientCreateOrReplaceResponse{}, err
	}
	resp, err := client.createOrReplaceHandleResponse(httpResp)
	return resp, err
}

// createOrReplaceCreateRequest creates the CreateOrReplace request.
func (client *FunctionsClient) createOrReplaceCreateRequest(ctx context.Context, resourceGroupName string, jobName string, functionName string, function Function, options *FunctionsClientCreateOrReplaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if functionName == "" {
		return nil, errors.New("parameter functionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{functionName}", url.PathEscape(functionName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*options.IfMatch}
	}
	if options != nil && options.IfNoneMatch != nil {
		req.Raw().Header["If-None-Match"] = []string{*options.IfNoneMatch}
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, function); err != nil {
	return nil, err
}
	return req, nil
}

// createOrReplaceHandleResponse handles the CreateOrReplace response.
func (client *FunctionsClient) createOrReplaceHandleResponse(resp *http.Response) (FunctionsClientCreateOrReplaceResponse, error) {
	result := FunctionsClientCreateOrReplaceResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Function); err != nil {
		return FunctionsClientCreateOrReplaceResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a function from the streaming job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - functionName - The name of the function.
//   - options - FunctionsClientDeleteOptions contains the optional parameters for the FunctionsClient.Delete method.
func (client *FunctionsClient) Delete(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientDeleteOptions) (FunctionsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, jobName, functionName, options)
	if err != nil {
		return FunctionsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FunctionsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return FunctionsClientDeleteResponse{}, err
	}
	return FunctionsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *FunctionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if functionName == "" {
		return nil, errors.New("parameter functionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{functionName}", url.PathEscape(functionName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets details about the specified function.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - functionName - The name of the function.
//   - options - FunctionsClientGetOptions contains the optional parameters for the FunctionsClient.Get method.
func (client *FunctionsClient) Get(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientGetOptions) (FunctionsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, jobName, functionName, options)
	if err != nil {
		return FunctionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FunctionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FunctionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *FunctionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if functionName == "" {
		return nil, errors.New("parameter functionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{functionName}", url.PathEscape(functionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *FunctionsClient) getHandleResponse(resp *http.Response) (FunctionsClientGetResponse, error) {
	result := FunctionsClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Function); err != nil {
		return FunctionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByStreamingJobPager - Lists all of the functions under the specified streaming job.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - options - FunctionsClientListByStreamingJobOptions contains the optional parameters for the FunctionsClient.NewListByStreamingJobPager
//     method.
func (client *FunctionsClient) NewListByStreamingJobPager(resourceGroupName string, jobName string, options *FunctionsClientListByStreamingJobOptions) (*runtime.Pager[FunctionsClientListByStreamingJobResponse]) {
	return runtime.NewPager(runtime.PagingHandler[FunctionsClientListByStreamingJobResponse]{
		More: func(page FunctionsClientListByStreamingJobResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *FunctionsClientListByStreamingJobResponse) (FunctionsClientListByStreamingJobResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByStreamingJobCreateRequest(ctx, resourceGroupName, jobName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return FunctionsClientListByStreamingJobResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return FunctionsClientListByStreamingJobResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return FunctionsClientListByStreamingJobResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByStreamingJobHandleResponse(resp)
		},
	})
}

// listByStreamingJobCreateRequest creates the ListByStreamingJob request.
func (client *FunctionsClient) listByStreamingJobCreateRequest(ctx context.Context, resourceGroupName string, jobName string, options *FunctionsClientListByStreamingJobOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Select != nil {
		reqQP.Set("$select", *options.Select)
	}
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByStreamingJobHandleResponse handles the ListByStreamingJob response.
func (client *FunctionsClient) listByStreamingJobHandleResponse(resp *http.Response) (FunctionsClientListByStreamingJobResponse, error) {
	result := FunctionsClientListByStreamingJobResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FunctionListResult); err != nil {
		return FunctionsClientListByStreamingJobResponse{}, err
	}
	return result, nil
}

// RetrieveDefaultDefinition - Retrieves the default definition of a function based on the parameters specified.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - functionName - The name of the function.
//   - options - FunctionsClientRetrieveDefaultDefinitionOptions contains the optional parameters for the FunctionsClient.RetrieveDefaultDefinition
//     method.
func (client *FunctionsClient) RetrieveDefaultDefinition(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientRetrieveDefaultDefinitionOptions) (FunctionsClientRetrieveDefaultDefinitionResponse, error) {
	var err error
	req, err := client.retrieveDefaultDefinitionCreateRequest(ctx, resourceGroupName, jobName, functionName, options)
	if err != nil {
		return FunctionsClientRetrieveDefaultDefinitionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FunctionsClientRetrieveDefaultDefinitionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FunctionsClientRetrieveDefaultDefinitionResponse{}, err
	}
	resp, err := client.retrieveDefaultDefinitionHandleResponse(httpResp)
	return resp, err
}

// retrieveDefaultDefinitionCreateRequest creates the RetrieveDefaultDefinition request.
func (client *FunctionsClient) retrieveDefaultDefinitionCreateRequest(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientRetrieveDefaultDefinitionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}/retrieveDefaultDefinition"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if functionName == "" {
		return nil, errors.New("parameter functionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{functionName}", url.PathEscape(functionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.FunctionRetrieveDefaultDefinitionParameters != nil {
		if err := runtime.MarshalAsJSON(req, options.FunctionRetrieveDefaultDefinitionParameters); err != nil {
	return nil, err
}
		return req, nil
	}
	return req, nil
}

// retrieveDefaultDefinitionHandleResponse handles the RetrieveDefaultDefinition response.
func (client *FunctionsClient) retrieveDefaultDefinitionHandleResponse(resp *http.Response) (FunctionsClientRetrieveDefaultDefinitionResponse, error) {
	result := FunctionsClientRetrieveDefaultDefinitionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Function); err != nil {
		return FunctionsClientRetrieveDefaultDefinitionResponse{}, err
	}
	return result, nil
}

// BeginTest - Tests if the information provided for a function is valid. This can range from testing the connection to the
// underlying web service behind the function or making sure the function code provided is
// syntactically correct.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - functionName - The name of the function.
//   - options - FunctionsClientBeginTestOptions contains the optional parameters for the FunctionsClient.BeginTest method.
func (client *FunctionsClient) BeginTest(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientBeginTestOptions) (*runtime.Poller[FunctionsClientTestResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.test(ctx, resourceGroupName, jobName, functionName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[FunctionsClientTestResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[FunctionsClientTestResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Test - Tests if the information provided for a function is valid. This can range from testing the connection to the underlying
// web service behind the function or making sure the function code provided is
// syntactically correct.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
func (client *FunctionsClient) test(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientBeginTestOptions) (*http.Response, error) {
	var err error
	req, err := client.testCreateRequest(ctx, resourceGroupName, jobName, functionName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// testCreateRequest creates the Test request.
func (client *FunctionsClient) testCreateRequest(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *FunctionsClientBeginTestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}/test"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if functionName == "" {
		return nil, errors.New("parameter functionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{functionName}", url.PathEscape(functionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Function != nil {
		if err := runtime.MarshalAsJSON(req, *options.Function); err != nil {
	return nil, err
}
		return req, nil
	}
	return req, nil
}

// Update - Updates an existing function under an existing streaming job. This can be used to partially update (ie. update
// one or two properties) a function without affecting the rest the job or function
// definition.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - functionName - The name of the function.
//   - function - A function object. The properties specified here will overwrite the corresponding properties in the existing
//     function (ie. Those properties will be updated). Any properties that are set to null here
//     will mean that the corresponding property in the existing function will remain the same and not change as a result of this
//     PATCH operation.
//   - options - FunctionsClientUpdateOptions contains the optional parameters for the FunctionsClient.Update method.
func (client *FunctionsClient) Update(ctx context.Context, resourceGroupName string, jobName string, functionName string, function Function, options *FunctionsClientUpdateOptions) (FunctionsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, jobName, functionName, function, options)
	if err != nil {
		return FunctionsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FunctionsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FunctionsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *FunctionsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, jobName string, functionName string, function Function, options *FunctionsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/functions/{functionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if functionName == "" {
		return nil, errors.New("parameter functionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{functionName}", url.PathEscape(functionName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*options.IfMatch}
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, function); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *FunctionsClient) updateHandleResponse(resp *http.Response) (FunctionsClientUpdateResponse, error) {
	result := FunctionsClientUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Function); err != nil {
		return FunctionsClientUpdateResponse{}, err
	}
	return result, nil
}


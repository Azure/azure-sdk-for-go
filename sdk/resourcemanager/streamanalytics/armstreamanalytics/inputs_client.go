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

// InputsClient contains the methods for the Inputs group.
// Don't use this type directly, use NewInputsClient() instead.
type InputsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewInputsClient creates a new instance of InputsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewInputsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*InputsClient, error) {
	cl, err := arm.NewClient(moduleName+".InputsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &InputsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrReplace - Creates an input or replaces an already existing input under an existing streaming job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - inputName - The name of the input.
//   - input - The definition of the input that will be used to create a new input or replace the existing one under the streaming
//     job.
//   - options - InputsClientCreateOrReplaceOptions contains the optional parameters for the InputsClient.CreateOrReplace method.
func (client *InputsClient) CreateOrReplace(ctx context.Context, resourceGroupName string, jobName string, inputName string, input Input, options *InputsClientCreateOrReplaceOptions) (InputsClientCreateOrReplaceResponse, error) {
	var err error
	req, err := client.createOrReplaceCreateRequest(ctx, resourceGroupName, jobName, inputName, input, options)
	if err != nil {
		return InputsClientCreateOrReplaceResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InputsClientCreateOrReplaceResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return InputsClientCreateOrReplaceResponse{}, err
	}
	resp, err := client.createOrReplaceHandleResponse(httpResp)
	return resp, err
}

// createOrReplaceCreateRequest creates the CreateOrReplace request.
func (client *InputsClient) createOrReplaceCreateRequest(ctx context.Context, resourceGroupName string, jobName string, inputName string, input Input, options *InputsClientCreateOrReplaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/inputs/{inputName}"
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
	if inputName == "" {
		return nil, errors.New("parameter inputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inputName}", url.PathEscape(inputName))
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
	if err := runtime.MarshalAsJSON(req, input); err != nil {
	return nil, err
}
	return req, nil
}

// createOrReplaceHandleResponse handles the CreateOrReplace response.
func (client *InputsClient) createOrReplaceHandleResponse(resp *http.Response) (InputsClientCreateOrReplaceResponse, error) {
	result := InputsClientCreateOrReplaceResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Input); err != nil {
		return InputsClientCreateOrReplaceResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an input from the streaming job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - inputName - The name of the input.
//   - options - InputsClientDeleteOptions contains the optional parameters for the InputsClient.Delete method.
func (client *InputsClient) Delete(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientDeleteOptions) (InputsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, jobName, inputName, options)
	if err != nil {
		return InputsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InputsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return InputsClientDeleteResponse{}, err
	}
	return InputsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *InputsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/inputs/{inputName}"
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
	if inputName == "" {
		return nil, errors.New("parameter inputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inputName}", url.PathEscape(inputName))
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

// Get - Gets details about the specified input.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - inputName - The name of the input.
//   - options - InputsClientGetOptions contains the optional parameters for the InputsClient.Get method.
func (client *InputsClient) Get(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientGetOptions) (InputsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, jobName, inputName, options)
	if err != nil {
		return InputsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InputsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return InputsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *InputsClient) getCreateRequest(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/inputs/{inputName}"
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
	if inputName == "" {
		return nil, errors.New("parameter inputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inputName}", url.PathEscape(inputName))
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
func (client *InputsClient) getHandleResponse(resp *http.Response) (InputsClientGetResponse, error) {
	result := InputsClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Input); err != nil {
		return InputsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByStreamingJobPager - Lists all of the inputs under the specified streaming job.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - options - InputsClientListByStreamingJobOptions contains the optional parameters for the InputsClient.NewListByStreamingJobPager
//     method.
func (client *InputsClient) NewListByStreamingJobPager(resourceGroupName string, jobName string, options *InputsClientListByStreamingJobOptions) (*runtime.Pager[InputsClientListByStreamingJobResponse]) {
	return runtime.NewPager(runtime.PagingHandler[InputsClientListByStreamingJobResponse]{
		More: func(page InputsClientListByStreamingJobResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *InputsClientListByStreamingJobResponse) (InputsClientListByStreamingJobResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByStreamingJobCreateRequest(ctx, resourceGroupName, jobName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return InputsClientListByStreamingJobResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return InputsClientListByStreamingJobResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return InputsClientListByStreamingJobResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByStreamingJobHandleResponse(resp)
		},
	})
}

// listByStreamingJobCreateRequest creates the ListByStreamingJob request.
func (client *InputsClient) listByStreamingJobCreateRequest(ctx context.Context, resourceGroupName string, jobName string, options *InputsClientListByStreamingJobOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/inputs"
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
func (client *InputsClient) listByStreamingJobHandleResponse(resp *http.Response) (InputsClientListByStreamingJobResponse, error) {
	result := InputsClientListByStreamingJobResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.InputListResult); err != nil {
		return InputsClientListByStreamingJobResponse{}, err
	}
	return result, nil
}

// BeginTest - Tests whether an input’s datasource is reachable and usable by the Azure Stream Analytics service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - inputName - The name of the input.
//   - options - InputsClientBeginTestOptions contains the optional parameters for the InputsClient.BeginTest method.
func (client *InputsClient) BeginTest(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientBeginTestOptions) (*runtime.Poller[InputsClientTestResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.test(ctx, resourceGroupName, jobName, inputName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[InputsClientTestResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[InputsClientTestResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Test - Tests whether an input’s datasource is reachable and usable by the Azure Stream Analytics service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
func (client *InputsClient) test(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientBeginTestOptions) (*http.Response, error) {
	var err error
	req, err := client.testCreateRequest(ctx, resourceGroupName, jobName, inputName, options)
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
func (client *InputsClient) testCreateRequest(ctx context.Context, resourceGroupName string, jobName string, inputName string, options *InputsClientBeginTestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/inputs/{inputName}/test"
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
	if inputName == "" {
		return nil, errors.New("parameter inputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inputName}", url.PathEscape(inputName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Input != nil {
		if err := runtime.MarshalAsJSON(req, *options.Input); err != nil {
	return nil, err
}
		return req, nil
	}
	return req, nil
}

// Update - Updates an existing input under an existing streaming job. This can be used to partially update (ie. update one
// or two properties) an input without affecting the rest the job or input definition.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - inputName - The name of the input.
//   - input - An Input object. The properties specified here will overwrite the corresponding properties in the existing input
//     (ie. Those properties will be updated). Any properties that are set to null here will
//     mean that the corresponding property in the existing input will remain the same and not change as a result of this PATCH
//     operation.
//   - options - InputsClientUpdateOptions contains the optional parameters for the InputsClient.Update method.
func (client *InputsClient) Update(ctx context.Context, resourceGroupName string, jobName string, inputName string, input Input, options *InputsClientUpdateOptions) (InputsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, jobName, inputName, input, options)
	if err != nil {
		return InputsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InputsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return InputsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *InputsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, jobName string, inputName string, input Input, options *InputsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/inputs/{inputName}"
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
	if inputName == "" {
		return nil, errors.New("parameter inputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{inputName}", url.PathEscape(inputName))
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
	if err := runtime.MarshalAsJSON(req, input); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *InputsClient) updateHandleResponse(resp *http.Response) (InputsClientUpdateResponse, error) {
	result := InputsClientUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Input); err != nil {
		return InputsClientUpdateResponse{}, err
	}
	return result, nil
}


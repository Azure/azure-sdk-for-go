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

// OutputsClient contains the methods for the Outputs group.
// Don't use this type directly, use NewOutputsClient() instead.
type OutputsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewOutputsClient creates a new instance of OutputsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewOutputsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*OutputsClient, error) {
	cl, err := arm.NewClient(moduleName+".OutputsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &OutputsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrReplace - Creates an output or replaces an already existing output under an existing streaming job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - outputName - The name of the output.
//   - output - The definition of the output that will be used to create a new output or replace the existing one under the streaming
//     job.
//   - options - OutputsClientCreateOrReplaceOptions contains the optional parameters for the OutputsClient.CreateOrReplace method.
func (client *OutputsClient) CreateOrReplace(ctx context.Context, resourceGroupName string, jobName string, outputName string, output Output, options *OutputsClientCreateOrReplaceOptions) (OutputsClientCreateOrReplaceResponse, error) {
	var err error
	req, err := client.createOrReplaceCreateRequest(ctx, resourceGroupName, jobName, outputName, output, options)
	if err != nil {
		return OutputsClientCreateOrReplaceResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OutputsClientCreateOrReplaceResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return OutputsClientCreateOrReplaceResponse{}, err
	}
	resp, err := client.createOrReplaceHandleResponse(httpResp)
	return resp, err
}

// createOrReplaceCreateRequest creates the CreateOrReplace request.
func (client *OutputsClient) createOrReplaceCreateRequest(ctx context.Context, resourceGroupName string, jobName string, outputName string, output Output, options *OutputsClientCreateOrReplaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs/{outputName}"
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
	if outputName == "" {
		return nil, errors.New("parameter outputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{outputName}", url.PathEscape(outputName))
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
	if err := runtime.MarshalAsJSON(req, output); err != nil {
	return nil, err
}
	return req, nil
}

// createOrReplaceHandleResponse handles the CreateOrReplace response.
func (client *OutputsClient) createOrReplaceHandleResponse(resp *http.Response) (OutputsClientCreateOrReplaceResponse, error) {
	result := OutputsClientCreateOrReplaceResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Output); err != nil {
		return OutputsClientCreateOrReplaceResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an output from the streaming job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - outputName - The name of the output.
//   - options - OutputsClientDeleteOptions contains the optional parameters for the OutputsClient.Delete method.
func (client *OutputsClient) Delete(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientDeleteOptions) (OutputsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, jobName, outputName, options)
	if err != nil {
		return OutputsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OutputsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return OutputsClientDeleteResponse{}, err
	}
	return OutputsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *OutputsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs/{outputName}"
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
	if outputName == "" {
		return nil, errors.New("parameter outputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{outputName}", url.PathEscape(outputName))
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

// Get - Gets details about the specified output.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - outputName - The name of the output.
//   - options - OutputsClientGetOptions contains the optional parameters for the OutputsClient.Get method.
func (client *OutputsClient) Get(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientGetOptions) (OutputsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, jobName, outputName, options)
	if err != nil {
		return OutputsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OutputsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return OutputsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *OutputsClient) getCreateRequest(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs/{outputName}"
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
	if outputName == "" {
		return nil, errors.New("parameter outputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{outputName}", url.PathEscape(outputName))
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
func (client *OutputsClient) getHandleResponse(resp *http.Response) (OutputsClientGetResponse, error) {
	result := OutputsClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Output); err != nil {
		return OutputsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByStreamingJobPager - Lists all of the outputs under the specified streaming job.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - options - OutputsClientListByStreamingJobOptions contains the optional parameters for the OutputsClient.NewListByStreamingJobPager
//     method.
func (client *OutputsClient) NewListByStreamingJobPager(resourceGroupName string, jobName string, options *OutputsClientListByStreamingJobOptions) (*runtime.Pager[OutputsClientListByStreamingJobResponse]) {
	return runtime.NewPager(runtime.PagingHandler[OutputsClientListByStreamingJobResponse]{
		More: func(page OutputsClientListByStreamingJobResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *OutputsClientListByStreamingJobResponse) (OutputsClientListByStreamingJobResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByStreamingJobCreateRequest(ctx, resourceGroupName, jobName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return OutputsClientListByStreamingJobResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return OutputsClientListByStreamingJobResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return OutputsClientListByStreamingJobResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByStreamingJobHandleResponse(resp)
		},
	})
}

// listByStreamingJobCreateRequest creates the ListByStreamingJob request.
func (client *OutputsClient) listByStreamingJobCreateRequest(ctx context.Context, resourceGroupName string, jobName string, options *OutputsClientListByStreamingJobOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs"
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
func (client *OutputsClient) listByStreamingJobHandleResponse(resp *http.Response) (OutputsClientListByStreamingJobResponse, error) {
	result := OutputsClientListByStreamingJobResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OutputListResult); err != nil {
		return OutputsClientListByStreamingJobResponse{}, err
	}
	return result, nil
}

// BeginTest - Tests whether an output’s datasource is reachable and usable by the Azure Stream Analytics service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - outputName - The name of the output.
//   - options - OutputsClientBeginTestOptions contains the optional parameters for the OutputsClient.BeginTest method.
func (client *OutputsClient) BeginTest(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientBeginTestOptions) (*runtime.Poller[OutputsClientTestResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.test(ctx, resourceGroupName, jobName, outputName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[OutputsClientTestResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[OutputsClientTestResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Test - Tests whether an output’s datasource is reachable and usable by the Azure Stream Analytics service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
func (client *OutputsClient) test(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientBeginTestOptions) (*http.Response, error) {
	var err error
	req, err := client.testCreateRequest(ctx, resourceGroupName, jobName, outputName, options)
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
func (client *OutputsClient) testCreateRequest(ctx context.Context, resourceGroupName string, jobName string, outputName string, options *OutputsClientBeginTestOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs/{outputName}/test"
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
	if outputName == "" {
		return nil, errors.New("parameter outputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{outputName}", url.PathEscape(outputName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Output != nil {
		if err := runtime.MarshalAsJSON(req, *options.Output); err != nil {
	return nil, err
}
		return req, nil
	}
	return req, nil
}

// Update - Updates an existing output under an existing streaming job. This can be used to partially update (ie. update one
// or two properties) an output without affecting the rest the job or output definition.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-03-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - jobName - The name of the streaming job.
//   - outputName - The name of the output.
//   - output - An Output object. The properties specified here will overwrite the corresponding properties in the existing output
//     (ie. Those properties will be updated). Any properties that are set to null here will
//     mean that the corresponding property in the existing output will remain the same and not change as a result of this PATCH
//     operation.
//   - options - OutputsClientUpdateOptions contains the optional parameters for the OutputsClient.Update method.
func (client *OutputsClient) Update(ctx context.Context, resourceGroupName string, jobName string, outputName string, output Output, options *OutputsClientUpdateOptions) (OutputsClientUpdateResponse, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, jobName, outputName, output, options)
	if err != nil {
		return OutputsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return OutputsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return OutputsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *OutputsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, jobName string, outputName string, output Output, options *OutputsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StreamAnalytics/streamingjobs/{jobName}/outputs/{outputName}"
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
	if outputName == "" {
		return nil, errors.New("parameter outputName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{outputName}", url.PathEscape(outputName))
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
	if err := runtime.MarshalAsJSON(req, output); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *OutputsClient) updateHandleResponse(resp *http.Response) (OutputsClientUpdateResponse, error) {
	result := OutputsClientUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Output); err != nil {
		return OutputsClientUpdateResponse{}, err
	}
	return result, nil
}


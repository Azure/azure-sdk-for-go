//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatafactory

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// PipelineRunsClient contains the methods for the PipelineRuns group.
// Don't use this type directly, use NewPipelineRunsClient() instead.
type PipelineRunsClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewPipelineRunsClient creates a new instance of PipelineRunsClient with the specified values.
func NewPipelineRunsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *PipelineRunsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &PipelineRunsClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// Cancel - Cancel a pipeline run by its run ID.
// If the operation fails it returns the *CloudError error type.
func (client *PipelineRunsClient) Cancel(ctx context.Context, resourceGroupName string, factoryName string, runID string, options *PipelineRunsCancelOptions) (PipelineRunsCancelResponse, error) {
	req, err := client.cancelCreateRequest(ctx, resourceGroupName, factoryName, runID, options)
	if err != nil {
		return PipelineRunsCancelResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PipelineRunsCancelResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PipelineRunsCancelResponse{}, client.cancelHandleError(resp)
	}
	return PipelineRunsCancelResponse{RawResponse: resp}, nil
}

// cancelCreateRequest creates the Cancel request.
func (client *PipelineRunsClient) cancelCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, runID string, options *PipelineRunsCancelOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/pipelineruns/{runId}/cancel"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	if runID == "" {
		return nil, errors.New("parameter runID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runId}", url.PathEscape(runID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.IsRecursive != nil {
		reqQP.Set("isRecursive", strconv.FormatBool(*options.IsRecursive))
	}
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// cancelHandleError handles the Cancel error response.
func (client *PipelineRunsClient) cancelHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Get a pipeline run by its run ID.
// If the operation fails it returns the *CloudError error type.
func (client *PipelineRunsClient) Get(ctx context.Context, resourceGroupName string, factoryName string, runID string, options *PipelineRunsGetOptions) (PipelineRunsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, factoryName, runID, options)
	if err != nil {
		return PipelineRunsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PipelineRunsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PipelineRunsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *PipelineRunsClient) getCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, runID string, options *PipelineRunsGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/pipelineruns/{runId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	if runID == "" {
		return nil, errors.New("parameter runID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runId}", url.PathEscape(runID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PipelineRunsClient) getHandleResponse(resp *http.Response) (PipelineRunsGetResponse, error) {
	result := PipelineRunsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PipelineRun); err != nil {
		return PipelineRunsGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *PipelineRunsClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// QueryByFactory - Query pipeline runs in the factory based on input filter conditions.
// If the operation fails it returns the *CloudError error type.
func (client *PipelineRunsClient) QueryByFactory(ctx context.Context, resourceGroupName string, factoryName string, filterParameters RunFilterParameters, options *PipelineRunsQueryByFactoryOptions) (PipelineRunsQueryByFactoryResponse, error) {
	req, err := client.queryByFactoryCreateRequest(ctx, resourceGroupName, factoryName, filterParameters, options)
	if err != nil {
		return PipelineRunsQueryByFactoryResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PipelineRunsQueryByFactoryResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PipelineRunsQueryByFactoryResponse{}, client.queryByFactoryHandleError(resp)
	}
	return client.queryByFactoryHandleResponse(resp)
}

// queryByFactoryCreateRequest creates the QueryByFactory request.
func (client *PipelineRunsClient) queryByFactoryCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, filterParameters RunFilterParameters, options *PipelineRunsQueryByFactoryOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/queryPipelineRuns"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, filterParameters)
}

// queryByFactoryHandleResponse handles the QueryByFactory response.
func (client *PipelineRunsClient) queryByFactoryHandleResponse(resp *http.Response) (PipelineRunsQueryByFactoryResponse, error) {
	result := PipelineRunsQueryByFactoryResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PipelineRunsQueryResponse); err != nil {
		return PipelineRunsQueryByFactoryResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// queryByFactoryHandleError handles the QueryByFactory error response.
func (client *PipelineRunsClient) queryByFactoryHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

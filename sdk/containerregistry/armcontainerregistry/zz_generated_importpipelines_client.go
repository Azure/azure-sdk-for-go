// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcontainerregistry

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ImportPipelinesClient contains the methods for the ImportPipelines group.
// Don't use this type directly, use NewImportPipelinesClient() instead.
type ImportPipelinesClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewImportPipelinesClient creates a new instance of ImportPipelinesClient with the specified values.
func NewImportPipelinesClient(con *armcore.Connection, subscriptionID string) *ImportPipelinesClient {
	return &ImportPipelinesClient{con: con, subscriptionID: subscriptionID}
}

// BeginCreate - Creates an import pipeline for a container registry with the specified parameters.
// If the operation fails it returns a generic error.
func (client *ImportPipelinesClient) BeginCreate(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, importPipelineCreateParameters ImportPipeline, options *ImportPipelinesBeginCreateOptions) (ImportPipelinePollerResponse, error) {
	resp, err := client.create(ctx, resourceGroupName, registryName, importPipelineName, importPipelineCreateParameters, options)
	if err != nil {
		return ImportPipelinePollerResponse{}, err
	}
	result := ImportPipelinePollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewLROPoller("ImportPipelinesClient.Create", "", resp, client.con.Pipeline(), client.createHandleError)
	if err != nil {
		return ImportPipelinePollerResponse{}, err
	}
	poller := &importPipelinePoller{
		pt: pt,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (ImportPipelineResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeCreate creates a new ImportPipelinePoller from the specified resume token.
// token - The value must come from a previous call to ImportPipelinePoller.ResumeToken().
func (client *ImportPipelinesClient) ResumeCreate(ctx context.Context, token string) (ImportPipelinePollerResponse, error) {
	pt, err := armcore.NewLROPollerFromResumeToken("ImportPipelinesClient.Create", token, client.con.Pipeline(), client.createHandleError)
	if err != nil {
		return ImportPipelinePollerResponse{}, err
	}
	poller := &importPipelinePoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return ImportPipelinePollerResponse{}, err
	}
	result := ImportPipelinePollerResponse{
		RawResponse: resp,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (ImportPipelineResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// Create - Creates an import pipeline for a container registry with the specified parameters.
// If the operation fails it returns a generic error.
func (client *ImportPipelinesClient) create(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, importPipelineCreateParameters ImportPipeline, options *ImportPipelinesBeginCreateOptions) (*azcore.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, registryName, importPipelineName, importPipelineCreateParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return nil, client.createHandleError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *ImportPipelinesClient) createCreateRequest(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, importPipelineCreateParameters ImportPipeline, options *ImportPipelinesBeginCreateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/importPipelines/{importPipelineName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if importPipelineName == "" {
		return nil, errors.New("parameter importPipelineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{importPipelineName}", url.PathEscape(importPipelineName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(importPipelineCreateParameters)
}

// createHandleError handles the Create error response.
func (client *ImportPipelinesClient) createHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// BeginDelete - Deletes an import pipeline from a container registry.
// If the operation fails it returns a generic error.
func (client *ImportPipelinesClient) BeginDelete(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesBeginDeleteOptions) (HTTPPollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, registryName, importPipelineName, options)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewLROPoller("ImportPipelinesClient.Delete", "", resp, client.con.Pipeline(), client.deleteHandleError)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	poller := &httpPoller{
		pt: pt,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeDelete creates a new HTTPPoller from the specified resume token.
// token - The value must come from a previous call to HTTPPoller.ResumeToken().
func (client *ImportPipelinesClient) ResumeDelete(ctx context.Context, token string) (HTTPPollerResponse, error) {
	pt, err := armcore.NewLROPollerFromResumeToken("ImportPipelinesClient.Delete", token, client.con.Pipeline(), client.deleteHandleError)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	poller := &httpPoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// Delete - Deletes an import pipeline from a container registry.
// If the operation fails it returns a generic error.
func (client *ImportPipelinesClient) deleteOperation(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesBeginDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, registryName, importPipelineName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ImportPipelinesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesBeginDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/importPipelines/{importPipelineName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if importPipelineName == "" {
		return nil, errors.New("parameter importPipelineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{importPipelineName}", url.PathEscape(importPipelineName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *ImportPipelinesClient) deleteHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Get - Gets the properties of the import pipeline.
// If the operation fails it returns a generic error.
func (client *ImportPipelinesClient) Get(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesGetOptions) (ImportPipelineResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, registryName, importPipelineName, options)
	if err != nil {
		return ImportPipelineResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ImportPipelineResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return ImportPipelineResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ImportPipelinesClient) getCreateRequest(ctx context.Context, resourceGroupName string, registryName string, importPipelineName string, options *ImportPipelinesGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/importPipelines/{importPipelineName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if importPipelineName == "" {
		return nil, errors.New("parameter importPipelineName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{importPipelineName}", url.PathEscape(importPipelineName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ImportPipelinesClient) getHandleResponse(resp *azcore.Response) (ImportPipelineResponse, error) {
	var val *ImportPipeline
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return ImportPipelineResponse{}, err
	}
	return ImportPipelineResponse{RawResponse: resp.Response, ImportPipeline: val}, nil
}

// getHandleError handles the Get error response.
func (client *ImportPipelinesClient) getHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// List - Lists all import pipelines for the specified container registry.
// If the operation fails it returns a generic error.
func (client *ImportPipelinesClient) List(resourceGroupName string, registryName string, options *ImportPipelinesListOptions) ImportPipelineListResultPager {
	return &importPipelineListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, registryName, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp ImportPipelineListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ImportPipelineListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client *ImportPipelinesClient) listCreateRequest(ctx context.Context, resourceGroupName string, registryName string, options *ImportPipelinesListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/importPipelines"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ImportPipelinesClient) listHandleResponse(resp *azcore.Response) (ImportPipelineListResultResponse, error) {
	var val *ImportPipelineListResult
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return ImportPipelineListResultResponse{}, err
	}
	return ImportPipelineListResultResponse{RawResponse: resp.Response, ImportPipelineListResult: val}, nil
}

// listHandleError handles the List error response.
func (client *ImportPipelinesClient) listHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RestorePointsClient contains the methods for the RestorePoints group.
// Don't use this type directly, use NewRestorePointsClient() instead.
type RestorePointsClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewRestorePointsClient creates a new instance of RestorePointsClient with the specified values.
func NewRestorePointsClient(con *armcore.Connection, subscriptionID string) *RestorePointsClient {
	return &RestorePointsClient{con: con, subscriptionID: subscriptionID}
}

// BeginCreate - The operation to create the restore point. Updating properties of an existing restore point is not allowed
// If the operation fails it returns the *CloudError error type.
func (client *RestorePointsClient) BeginCreate(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, parameters RestorePoint, options *RestorePointsBeginCreateOptions) (RestorePointPollerResponse, error) {
	resp, err := client.create(ctx, resourceGroupName, restorePointCollectionName, restorePointName, parameters, options)
	if err != nil {
		return RestorePointPollerResponse{}, err
	}
	result := RestorePointPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewLROPoller("RestorePointsClient.Create", "", resp, client.con.Pipeline(), client.createHandleError)
	if err != nil {
		return RestorePointPollerResponse{}, err
	}
	poller := &restorePointPoller{
		pt: pt,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (RestorePointResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeCreate creates a new RestorePointPoller from the specified resume token.
// token - The value must come from a previous call to RestorePointPoller.ResumeToken().
func (client *RestorePointsClient) ResumeCreate(ctx context.Context, token string) (RestorePointPollerResponse, error) {
	pt, err := armcore.NewLROPollerFromResumeToken("RestorePointsClient.Create", token, client.con.Pipeline(), client.createHandleError)
	if err != nil {
		return RestorePointPollerResponse{}, err
	}
	poller := &restorePointPoller{
		pt: pt,
	}
	resp, err := poller.Poll(ctx)
	if err != nil {
		return RestorePointPollerResponse{}, err
	}
	result := RestorePointPollerResponse{
		RawResponse: resp,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (RestorePointResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// Create - The operation to create the restore point. Updating properties of an existing restore point is not allowed
// If the operation fails it returns the *CloudError error type.
func (client *RestorePointsClient) create(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, parameters RestorePoint, options *RestorePointsBeginCreateOptions) (*azcore.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, restorePointCollectionName, restorePointName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, client.createHandleError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *RestorePointsClient) createCreateRequest(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, parameters RestorePoint, options *RestorePointsBeginCreateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/restorePointCollections/{restorePointCollectionName}/restorePoints/{restorePointName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if restorePointCollectionName == "" {
		return nil, errors.New("parameter restorePointCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointCollectionName}", url.PathEscape(restorePointCollectionName))
	if restorePointName == "" {
		return nil, errors.New("parameter restorePointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointName}", url.PathEscape(restorePointName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-03-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// createHandleError handles the Create error response.
func (client *RestorePointsClient) createHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// BeginDelete - The operation to delete the restore point.
// If the operation fails it returns the *CloudError error type.
func (client *RestorePointsClient) BeginDelete(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, options *RestorePointsBeginDeleteOptions) (HTTPPollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, restorePointCollectionName, restorePointName, options)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewLROPoller("RestorePointsClient.Delete", "", resp, client.con.Pipeline(), client.deleteHandleError)
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
func (client *RestorePointsClient) ResumeDelete(ctx context.Context, token string) (HTTPPollerResponse, error) {
	pt, err := armcore.NewLROPollerFromResumeToken("RestorePointsClient.Delete", token, client.con.Pipeline(), client.deleteHandleError)
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

// Delete - The operation to delete the restore point.
// If the operation fails it returns the *CloudError error type.
func (client *RestorePointsClient) deleteOperation(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, options *RestorePointsBeginDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, restorePointCollectionName, restorePointName, options)
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
func (client *RestorePointsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, options *RestorePointsBeginDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/restorePointCollections/{restorePointCollectionName}/restorePoints/{restorePointName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if restorePointCollectionName == "" {
		return nil, errors.New("parameter restorePointCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointCollectionName}", url.PathEscape(restorePointCollectionName))
	if restorePointName == "" {
		return nil, errors.New("parameter restorePointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointName}", url.PathEscape(restorePointName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-03-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *RestorePointsClient) deleteHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// Get - The operation to get the restore point.
// If the operation fails it returns the *CloudError error type.
func (client *RestorePointsClient) Get(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, options *RestorePointsGetOptions) (RestorePointResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, restorePointCollectionName, restorePointName, options)
	if err != nil {
		return RestorePointResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return RestorePointResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return RestorePointResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RestorePointsClient) getCreateRequest(ctx context.Context, resourceGroupName string, restorePointCollectionName string, restorePointName string, options *RestorePointsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/restorePointCollections/{restorePointCollectionName}/restorePoints/{restorePointName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if restorePointCollectionName == "" {
		return nil, errors.New("parameter restorePointCollectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointCollectionName}", url.PathEscape(restorePointCollectionName))
	if restorePointName == "" {
		return nil, errors.New("parameter restorePointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointName}", url.PathEscape(restorePointName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-03-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RestorePointsClient) getHandleResponse(resp *azcore.Response) (RestorePointResponse, error) {
	var val *RestorePoint
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return RestorePointResponse{}, err
	}
	return RestorePointResponse{RawResponse: resp.Response, RestorePoint: val}, nil
}

// getHandleError handles the Get error response.
func (client *RestorePointsClient) getHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

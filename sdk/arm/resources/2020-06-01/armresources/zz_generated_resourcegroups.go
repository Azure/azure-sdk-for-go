// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armresources

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ResourceGroupsClient contains the methods for the ResourceGroups group.
// Don't use this type directly, use NewResourceGroupsClient() instead.
type ResourceGroupsClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewResourceGroupsClient creates a new instance of ResourceGroupsClient with the specified values.
func NewResourceGroupsClient(con *armcore.Connection, subscriptionID string) ResourceGroupsClient {
	return ResourceGroupsClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client ResourceGroupsClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// CheckExistence - Checks whether a resource group exists.
func (client ResourceGroupsClient) CheckExistence(ctx context.Context, resourceGroupName string, options *ResourceGroupsCheckExistenceOptions) (BooleanResponse, error) {
	req, err := client.checkExistenceCreateRequest(ctx, resourceGroupName, options)
	if err != nil {
		return BooleanResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return BooleanResponse{}, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return BooleanResponse{RawResponse: resp.Response, Success: true}, nil
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return BooleanResponse{RawResponse: resp.Response, Success: false}, nil
	} else {
		return BooleanResponse{}, client.checkExistenceHandleError(resp)
	}
}

// checkExistenceCreateRequest creates the CheckExistence request.
func (client ResourceGroupsClient) checkExistenceCreateRequest(ctx context.Context, resourceGroupName string, options *ResourceGroupsCheckExistenceOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodHead, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// checkExistenceHandleError handles the CheckExistence error response.
func (client ResourceGroupsClient) checkExistenceHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// CreateOrUpdate - Creates or updates a resource group.
func (client ResourceGroupsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, parameters ResourceGroup, options *ResourceGroupsCreateOrUpdateOptions) (ResourceGroupResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, parameters, options)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return ResourceGroupResponse{}, client.createOrUpdateHandleError(resp)
	}
	result, err := client.createOrUpdateHandleResponse(resp)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	return result, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client ResourceGroupsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, parameters ResourceGroup, options *ResourceGroupsCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client ResourceGroupsClient) createOrUpdateHandleResponse(resp *azcore.Response) (ResourceGroupResponse, error) {
	result := ResourceGroupResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ResourceGroup)
	return result, err
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client ResourceGroupsClient) createOrUpdateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginDelete - When you delete a resource group, all of its resources are also deleted. Deleting a resource group deletes all of its template deployments
// and currently stored operations.
func (client ResourceGroupsClient) BeginDelete(ctx context.Context, resourceGroupName string, options *ResourceGroupsBeginDeleteOptions) (HTTPPollerResponse, error) {
	resp, err := client.delete(ctx, resourceGroupName, options)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("ResourceGroupsClient.Delete", "", resp, client.deleteHandleError)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeDelete creates a new HTTPPoller from the specified resume token.
// token - The value must come from a previous call to HTTPPoller.ResumeToken().
func (client ResourceGroupsClient) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("ResourceGroupsClient.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Delete - When you delete a resource group, all of its resources are also deleted. Deleting a resource group deletes all of its template deployments and
// currently stored operations.
func (client ResourceGroupsClient) delete(ctx context.Context, resourceGroupName string, options *ResourceGroupsBeginDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client ResourceGroupsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, options *ResourceGroupsBeginDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client ResourceGroupsClient) deleteHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginExportTemplate - Captures the specified resource group as a template.
func (client ResourceGroupsClient) BeginExportTemplate(ctx context.Context, resourceGroupName string, parameters ExportTemplateRequest, options *ResourceGroupsBeginExportTemplateOptions) (ResourceGroupExportResultPollerResponse, error) {
	resp, err := client.exportTemplate(ctx, resourceGroupName, parameters, options)
	if err != nil {
		return ResourceGroupExportResultPollerResponse{}, err
	}
	result := ResourceGroupExportResultPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("ResourceGroupsClient.ExportTemplate", "location", resp, client.exportTemplateHandleError)
	if err != nil {
		return ResourceGroupExportResultPollerResponse{}, err
	}
	poller := &resourceGroupExportResultPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (ResourceGroupExportResultResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeExportTemplate creates a new ResourceGroupExportResultPoller from the specified resume token.
// token - The value must come from a previous call to ResourceGroupExportResultPoller.ResumeToken().
func (client ResourceGroupsClient) ResumeExportTemplate(token string) (ResourceGroupExportResultPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("ResourceGroupsClient.ExportTemplate", token, client.exportTemplateHandleError)
	if err != nil {
		return nil, err
	}
	return &resourceGroupExportResultPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// ExportTemplate - Captures the specified resource group as a template.
func (client ResourceGroupsClient) exportTemplate(ctx context.Context, resourceGroupName string, parameters ExportTemplateRequest, options *ResourceGroupsBeginExportTemplateOptions) (*azcore.Response, error) {
	req, err := client.exportTemplateCreateRequest(ctx, resourceGroupName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.exportTemplateHandleError(resp)
	}
	return resp, nil
}

// exportTemplateCreateRequest creates the ExportTemplate request.
func (client ResourceGroupsClient) exportTemplateCreateRequest(ctx context.Context, resourceGroupName string, parameters ExportTemplateRequest, options *ResourceGroupsBeginExportTemplateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/exportTemplate"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// exportTemplateHandleResponse handles the ExportTemplate response.
func (client ResourceGroupsClient) exportTemplateHandleResponse(resp *azcore.Response) (ResourceGroupExportResultResponse, error) {
	result := ResourceGroupExportResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ResourceGroupExportResult)
	return result, err
}

// exportTemplateHandleError handles the ExportTemplate error response.
func (client ResourceGroupsClient) exportTemplateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - Gets a resource group.
func (client ResourceGroupsClient) Get(ctx context.Context, resourceGroupName string, options *ResourceGroupsGetOptions) (ResourceGroupResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, options)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return ResourceGroupResponse{}, client.getHandleError(resp)
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client ResourceGroupsClient) getCreateRequest(ctx context.Context, resourceGroupName string, options *ResourceGroupsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client ResourceGroupsClient) getHandleResponse(resp *azcore.Response) (ResourceGroupResponse, error) {
	result := ResourceGroupResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ResourceGroup)
	return result, err
}

// getHandleError handles the Get error response.
func (client ResourceGroupsClient) getHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// List - Gets all the resource groups for a subscription.
func (client ResourceGroupsClient) List(options *ResourceGroupsListOptions) ResourceGroupListResultPager {
	return &resourceGroupListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp ResourceGroupListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ResourceGroupListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client ResourceGroupsClient) listCreateRequest(ctx context.Context, options *ResourceGroupsListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	if options != nil && options.Filter != nil {
		query.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		query.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client ResourceGroupsClient) listHandleResponse(resp *azcore.Response) (ResourceGroupListResultResponse, error) {
	result := ResourceGroupListResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ResourceGroupListResult)
	return result, err
}

// listHandleError handles the List error response.
func (client ResourceGroupsClient) listHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Update - Resource groups can be updated through a simple PATCH operation to a group address. The format of the request is the same as that for creating
// a resource group. If a field is unspecified, the current
// value is retained.
func (client ResourceGroupsClient) Update(ctx context.Context, resourceGroupName string, parameters ResourceGroupPatchable, options *ResourceGroupsUpdateOptions) (ResourceGroupResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, parameters, options)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return ResourceGroupResponse{}, client.updateHandleError(resp)
	}
	result, err := client.updateHandleResponse(resp)
	if err != nil {
		return ResourceGroupResponse{}, err
	}
	return result, nil
}

// updateCreateRequest creates the Update request.
func (client ResourceGroupsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, parameters ResourceGroupPatchable, options *ResourceGroupsUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// updateHandleResponse handles the Update response.
func (client ResourceGroupsClient) updateHandleResponse(resp *azcore.Response) (ResourceGroupResponse, error) {
	result := ResourceGroupResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ResourceGroup)
	return result, err
}

// updateHandleError handles the Update error response.
func (client ResourceGroupsClient) updateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

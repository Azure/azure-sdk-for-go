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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DiskAccessesClient contains the methods for the DiskAccesses group.
// Don't use this type directly, use NewDiskAccessesClient() instead.
type DiskAccessesClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewDiskAccessesClient creates a new instance of DiskAccessesClient with the specified values.
func NewDiskAccessesClient(con *armcore.Connection, subscriptionID string) DiskAccessesClient {
	return DiskAccessesClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client DiskAccessesClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// BeginCreateOrUpdate - Creates or updates a disk access resource
func (client DiskAccessesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, diskAccessName string, diskAccess DiskAccess, options *DiskAccessesCreateOrUpdateOptions) (*DiskAccessPollerResponse, error) {
	resp, err := client.CreateOrUpdate(ctx, resourceGroupName, diskAccessName, diskAccess, options)
	if err != nil {
		return nil, err
	}
	result := &DiskAccessPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("DiskAccessesClient.CreateOrUpdate", "", resp, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &diskAccessPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*DiskAccessResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeCreateOrUpdate creates a new DiskAccessPoller from the specified resume token.
// token - The value must come from a previous call to DiskAccessPoller.ResumeToken().
func (client DiskAccessesClient) ResumeCreateOrUpdate(token string) (DiskAccessPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("DiskAccessesClient.CreateOrUpdate", token, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &diskAccessPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// CreateOrUpdate - Creates or updates a disk access resource
func (client DiskAccessesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, diskAccessName string, diskAccess DiskAccess, options *DiskAccessesCreateOrUpdateOptions) (*azcore.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, diskAccessName, diskAccess, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client DiskAccessesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, diskAccessName string, diskAccess DiskAccess, options *DiskAccessesCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/diskAccesses/{diskAccessName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{diskAccessName}", url.PathEscape(diskAccessName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(diskAccess)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client DiskAccessesClient) createOrUpdateHandleResponse(resp *azcore.Response) (*DiskAccessResponse, error) {
	result := DiskAccessResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DiskAccess)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client DiskAccessesClient) createOrUpdateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginDelete - Deletes a disk access resource.
func (client DiskAccessesClient) BeginDelete(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesDeleteOptions) (*HTTPPollerResponse, error) {
	resp, err := client.Delete(ctx, resourceGroupName, diskAccessName, options)
	if err != nil {
		return nil, err
	}
	result := &HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("DiskAccessesClient.Delete", "", resp, client.deleteHandleError)
	if err != nil {
		return nil, err
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
func (client DiskAccessesClient) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("DiskAccessesClient.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Delete - Deletes a disk access resource.
func (client DiskAccessesClient) Delete(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, diskAccessName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client DiskAccessesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/diskAccesses/{diskAccessName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{diskAccessName}", url.PathEscape(diskAccessName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client DiskAccessesClient) deleteHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - Gets information about a disk access resource.
func (client DiskAccessesClient) Get(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesGetOptions) (*DiskAccessResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, diskAccessName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getHandleError(resp)
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client DiskAccessesClient) getCreateRequest(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/diskAccesses/{diskAccessName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{diskAccessName}", url.PathEscape(diskAccessName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client DiskAccessesClient) getHandleResponse(resp *azcore.Response) (*DiskAccessResponse, error) {
	result := DiskAccessResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DiskAccess)
}

// getHandleError handles the Get error response.
func (client DiskAccessesClient) getHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// GetPrivateLinkResources - Gets the private link resources possible under disk access resource
func (client DiskAccessesClient) GetPrivateLinkResources(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesGetPrivateLinkResourcesOptions) (*PrivateLinkResourceListResultResponse, error) {
	req, err := client.getPrivateLinkResourcesCreateRequest(ctx, resourceGroupName, diskAccessName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getPrivateLinkResourcesHandleError(resp)
	}
	result, err := client.getPrivateLinkResourcesHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getPrivateLinkResourcesCreateRequest creates the GetPrivateLinkResources request.
func (client DiskAccessesClient) getPrivateLinkResourcesCreateRequest(ctx context.Context, resourceGroupName string, diskAccessName string, options *DiskAccessesGetPrivateLinkResourcesOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/diskAccesses/{diskAccessName}/privateLinkResources"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{diskAccessName}", url.PathEscape(diskAccessName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getPrivateLinkResourcesHandleResponse handles the GetPrivateLinkResources response.
func (client DiskAccessesClient) getPrivateLinkResourcesHandleResponse(resp *azcore.Response) (*PrivateLinkResourceListResultResponse, error) {
	result := PrivateLinkResourceListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PrivateLinkResourceListResult)
}

// getPrivateLinkResourcesHandleError handles the GetPrivateLinkResources error response.
func (client DiskAccessesClient) getPrivateLinkResourcesHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// List - Lists all the disk access resources under a subscription.
func (client DiskAccessesClient) List(options *DiskAccessesListOptions) DiskAccessListPager {
	return &diskAccessListPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp *DiskAccessListResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.DiskAccessList.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client DiskAccessesClient) listCreateRequest(ctx context.Context, options *DiskAccessesListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/diskAccesses"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client DiskAccessesClient) listHandleResponse(resp *azcore.Response) (*DiskAccessListResponse, error) {
	result := DiskAccessListResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DiskAccessList)
}

// listHandleError handles the List error response.
func (client DiskAccessesClient) listHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// ListByResourceGroup - Lists all the disk access resources under a resource group.
func (client DiskAccessesClient) ListByResourceGroup(resourceGroupName string, options *DiskAccessesListByResourceGroupOptions) DiskAccessListPager {
	return &diskAccessListPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
		},
		responder: client.listByResourceGroupHandleResponse,
		errorer:   client.listByResourceGroupHandleError,
		advancer: func(ctx context.Context, resp *DiskAccessListResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.DiskAccessList.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client DiskAccessesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *DiskAccessesListByResourceGroupOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/diskAccesses"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client DiskAccessesClient) listByResourceGroupHandleResponse(resp *azcore.Response) (*DiskAccessListResponse, error) {
	result := DiskAccessListResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DiskAccessList)
}

// listByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client DiskAccessesClient) listByResourceGroupHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginUpdate - Updates (patches) a disk access resource.
func (client DiskAccessesClient) BeginUpdate(ctx context.Context, resourceGroupName string, diskAccessName string, diskAccess DiskAccessUpdate, options *DiskAccessesUpdateOptions) (*DiskAccessPollerResponse, error) {
	resp, err := client.Update(ctx, resourceGroupName, diskAccessName, diskAccess, options)
	if err != nil {
		return nil, err
	}
	result := &DiskAccessPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("DiskAccessesClient.Update", "", resp, client.updateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &diskAccessPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*DiskAccessResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeUpdate creates a new DiskAccessPoller from the specified resume token.
// token - The value must come from a previous call to DiskAccessPoller.ResumeToken().
func (client DiskAccessesClient) ResumeUpdate(token string) (DiskAccessPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("DiskAccessesClient.Update", token, client.updateHandleError)
	if err != nil {
		return nil, err
	}
	return &diskAccessPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Update - Updates (patches) a disk access resource.
func (client DiskAccessesClient) Update(ctx context.Context, resourceGroupName string, diskAccessName string, diskAccess DiskAccessUpdate, options *DiskAccessesUpdateOptions) (*azcore.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, diskAccessName, diskAccess, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client DiskAccessesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, diskAccessName string, diskAccess DiskAccessUpdate, options *DiskAccessesUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/diskAccesses/{diskAccessName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{diskAccessName}", url.PathEscape(diskAccessName))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(diskAccess)
}

// updateHandleResponse handles the Update response.
func (client DiskAccessesClient) updateHandleResponse(resp *azcore.Response) (*DiskAccessResponse, error) {
	result := DiskAccessResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.DiskAccess)
}

// updateHandleError handles the Update error response.
func (client DiskAccessesClient) updateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

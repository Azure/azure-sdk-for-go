// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VirtualWansClient contains the methods for the VirtualWans group.
// Don't use this type directly, use NewVirtualWansClient() instead.
type VirtualWansClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewVirtualWansClient creates a new instance of VirtualWansClient with the specified values.
func NewVirtualWansClient(con *armcore.Connection, subscriptionID string) VirtualWansClient {
	return VirtualWansClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client VirtualWansClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// BeginCreateOrUpdate - Creates a VirtualWAN resource if it doesn't exist else updates the existing VirtualWAN.
func (client VirtualWansClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, virtualWanName string, wanParameters VirtualWan, options *VirtualWansBeginCreateOrUpdateOptions) (VirtualWanPollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, virtualWanName, wanParameters, options)
	if err != nil {
		return VirtualWanPollerResponse{}, err
	}
	result := VirtualWanPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("VirtualWansClient.CreateOrUpdate", "azure-async-operation", resp, client.createOrUpdateHandleError)
	if err != nil {
		return VirtualWanPollerResponse{}, err
	}
	poller := &virtualWanPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (VirtualWanResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeCreateOrUpdate creates a new VirtualWanPoller from the specified resume token.
// token - The value must come from a previous call to VirtualWanPoller.ResumeToken().
func (client VirtualWansClient) ResumeCreateOrUpdate(token string) (VirtualWanPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("VirtualWansClient.CreateOrUpdate", token, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &virtualWanPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// CreateOrUpdate - Creates a VirtualWAN resource if it doesn't exist else updates the existing VirtualWAN.
func (client VirtualWansClient) createOrUpdate(ctx context.Context, resourceGroupName string, virtualWanName string, wanParameters VirtualWan, options *VirtualWansBeginCreateOrUpdateOptions) (*azcore.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, virtualWanName, wanParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client VirtualWansClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, virtualWanName string, wanParameters VirtualWan, options *VirtualWansBeginCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{VirtualWANName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{VirtualWANName}", url.PathEscape(virtualWanName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(wanParameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client VirtualWansClient) createOrUpdateHandleResponse(resp *azcore.Response) (VirtualWanResponse, error) {
	result := VirtualWanResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.VirtualWan)
	return result, err
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client VirtualWansClient) createOrUpdateHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// BeginDelete - Deletes a VirtualWAN.
func (client VirtualWansClient) BeginDelete(ctx context.Context, resourceGroupName string, virtualWanName string, options *VirtualWansBeginDeleteOptions) (HTTPPollerResponse, error) {
	resp, err := client.delete(ctx, resourceGroupName, virtualWanName, options)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("VirtualWansClient.Delete", "location", resp, client.deleteHandleError)
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
func (client VirtualWansClient) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("VirtualWansClient.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Delete - Deletes a VirtualWAN.
func (client VirtualWansClient) delete(ctx context.Context, resourceGroupName string, virtualWanName string, options *VirtualWansBeginDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, virtualWanName, options)
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
func (client VirtualWansClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, virtualWanName string, options *VirtualWansBeginDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{VirtualWANName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{VirtualWANName}", url.PathEscape(virtualWanName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client VirtualWansClient) deleteHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - Retrieves the details of a VirtualWAN.
func (client VirtualWansClient) Get(ctx context.Context, resourceGroupName string, virtualWanName string, options *VirtualWansGetOptions) (VirtualWanResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, virtualWanName, options)
	if err != nil {
		return VirtualWanResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return VirtualWanResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return VirtualWanResponse{}, client.getHandleError(resp)
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return VirtualWanResponse{}, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client VirtualWansClient) getCreateRequest(ctx context.Context, resourceGroupName string, virtualWanName string, options *VirtualWansGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{VirtualWANName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{VirtualWANName}", url.PathEscape(virtualWanName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client VirtualWansClient) getHandleResponse(resp *azcore.Response) (VirtualWanResponse, error) {
	result := VirtualWanResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.VirtualWan)
	return result, err
}

// getHandleError handles the Get error response.
func (client VirtualWansClient) getHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// List - Lists all the VirtualWANs in a subscription.
func (client VirtualWansClient) List(options *VirtualWansListOptions) ListVirtualWaNsResultPager {
	return &listVirtualWaNsResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp ListVirtualWaNsResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ListVirtualWaNsResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client VirtualWansClient) listCreateRequest(ctx context.Context, options *VirtualWansListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/virtualWans"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client VirtualWansClient) listHandleResponse(resp *azcore.Response) (ListVirtualWaNsResultResponse, error) {
	result := ListVirtualWaNsResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ListVirtualWaNsResult)
	return result, err
}

// listHandleError handles the List error response.
func (client VirtualWansClient) listHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// ListByResourceGroup - Lists all the VirtualWANs in a resource group.
func (client VirtualWansClient) ListByResourceGroup(resourceGroupName string, options *VirtualWansListByResourceGroupOptions) ListVirtualWaNsResultPager {
	return &listVirtualWaNsResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
		},
		responder: client.listByResourceGroupHandleResponse,
		errorer:   client.listByResourceGroupHandleError,
		advancer: func(ctx context.Context, resp ListVirtualWaNsResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ListVirtualWaNsResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client VirtualWansClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *VirtualWansListByResourceGroupOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client VirtualWansClient) listByResourceGroupHandleResponse(resp *azcore.Response) (ListVirtualWaNsResultResponse, error) {
	result := ListVirtualWaNsResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.ListVirtualWaNsResult)
	return result, err
}

// listByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client VirtualWansClient) listByResourceGroupHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// UpdateTags - Updates a VirtualWAN tags.
func (client VirtualWansClient) UpdateTags(ctx context.Context, resourceGroupName string, virtualWanName string, wanParameters TagsObject, options *VirtualWansUpdateTagsOptions) (VirtualWanResponse, error) {
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, virtualWanName, wanParameters, options)
	if err != nil {
		return VirtualWanResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return VirtualWanResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return VirtualWanResponse{}, client.updateTagsHandleError(resp)
	}
	result, err := client.updateTagsHandleResponse(resp)
	if err != nil {
		return VirtualWanResponse{}, err
	}
	return result, nil
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client VirtualWansClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, virtualWanName string, wanParameters TagsObject, options *VirtualWansUpdateTagsOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{VirtualWANName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{VirtualWANName}", url.PathEscape(virtualWanName))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-07-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(wanParameters)
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client VirtualWansClient) updateTagsHandleResponse(resp *azcore.Response) (VirtualWanResponse, error) {
	result := VirtualWanResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.VirtualWan)
	return result, err
}

// updateTagsHandleError handles the UpdateTags error response.
func (client VirtualWansClient) updateTagsHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

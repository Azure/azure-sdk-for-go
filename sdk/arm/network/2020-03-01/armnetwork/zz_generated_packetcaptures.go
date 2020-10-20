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

// PacketCapturesOperations contains the methods for the PacketCaptures group.
type PacketCapturesOperations interface {
	// BeginCreate - Create and start a packet capture on the specified VM.
	BeginCreate(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, parameters PacketCapture, options *PacketCapturesCreateOptions) (*PacketCaptureResultPollerResponse, error)
	// ResumeCreate - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeCreate(token string) (PacketCaptureResultPoller, error)
	// BeginDelete - Deletes the specified packet capture session.
	BeginDelete(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesDeleteOptions) (*HTTPPollerResponse, error)
	// ResumeDelete - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeDelete(token string) (HTTPPoller, error)
	// Get - Gets a packet capture session by name.
	Get(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetOptions) (*PacketCaptureResultResponse, error)
	// BeginGetStatus - Query the status of a running packet capture session.
	BeginGetStatus(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetStatusOptions) (*PacketCaptureQueryStatusResultPollerResponse, error)
	// ResumeGetStatus - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeGetStatus(token string) (PacketCaptureQueryStatusResultPoller, error)
	// List - Lists all packet capture sessions within the specified resource group.
	List(ctx context.Context, resourceGroupName string, networkWatcherName string, options *PacketCapturesListOptions) (*PacketCaptureListResultResponse, error)
	// BeginStop - Stops a specified packet capture session.
	BeginStop(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesStopOptions) (*HTTPPollerResponse, error)
	// ResumeStop - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeStop(token string) (HTTPPoller, error)
}

// PacketCapturesClient implements the PacketCapturesOperations interface.
// Don't use this type directly, use NewPacketCapturesClient() instead.
type PacketCapturesClient struct {
	*Client
	subscriptionID string
}

// NewPacketCapturesClient creates a new instance of PacketCapturesClient with the specified values.
func NewPacketCapturesClient(c *Client, subscriptionID string) PacketCapturesOperations {
	return &PacketCapturesClient{Client: c, subscriptionID: subscriptionID}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *PacketCapturesClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

func (client *PacketCapturesClient) BeginCreate(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, parameters PacketCapture, options *PacketCapturesCreateOptions) (*PacketCaptureResultPollerResponse, error) {
	resp, err := client.Create(ctx, resourceGroupName, networkWatcherName, packetCaptureName, parameters, options)
	if err != nil {
		return nil, err
	}
	result := &PacketCaptureResultPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("PacketCapturesClient.Create", "azure-async-operation", resp, client.CreateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &packetCaptureResultPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*PacketCaptureResultResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *PacketCapturesClient) ResumeCreate(token string) (PacketCaptureResultPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("PacketCapturesClient.Create", token, client.CreateHandleError)
	if err != nil {
		return nil, err
	}
	return &packetCaptureResultPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// Create - Create and start a packet capture on the specified VM.
func (client *PacketCapturesClient) Create(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, parameters PacketCapture, options *PacketCapturesCreateOptions) (*azcore.Response, error) {
	req, err := client.CreateCreateRequest(ctx, resourceGroupName, networkWatcherName, packetCaptureName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, client.CreateHandleError(resp)
	}
	return resp, nil
}

// CreateCreateRequest creates the Create request.
func (client *PacketCapturesClient) CreateCreateRequest(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, parameters PacketCapture, options *PacketCapturesCreateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/packetCaptures/{packetCaptureName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{networkWatcherName}", url.PathEscape(networkWatcherName))
	urlPath = strings.ReplaceAll(urlPath, "{packetCaptureName}", url.PathEscape(packetCaptureName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// CreateHandleResponse handles the Create response.
func (client *PacketCapturesClient) CreateHandleResponse(resp *azcore.Response) (*PacketCaptureResultResponse, error) {
	result := PacketCaptureResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PacketCaptureResult)
}

// CreateHandleError handles the Create error response.
func (client *PacketCapturesClient) CreateHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

func (client *PacketCapturesClient) BeginDelete(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesDeleteOptions) (*HTTPPollerResponse, error) {
	resp, err := client.Delete(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	result := &HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("PacketCapturesClient.Delete", "location", resp, client.DeleteHandleError)
	if err != nil {
		return nil, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *PacketCapturesClient) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("PacketCapturesClient.Delete", token, client.DeleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// Delete - Deletes the specified packet capture session.
func (client *PacketCapturesClient) Delete(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesDeleteOptions) (*azcore.Response, error) {
	req, err := client.DeleteCreateRequest(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusAccepted, http.StatusNoContent) {
		return nil, client.DeleteHandleError(resp)
	}
	return resp, nil
}

// DeleteCreateRequest creates the Delete request.
func (client *PacketCapturesClient) DeleteCreateRequest(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/packetCaptures/{packetCaptureName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{networkWatcherName}", url.PathEscape(networkWatcherName))
	urlPath = strings.ReplaceAll(urlPath, "{packetCaptureName}", url.PathEscape(packetCaptureName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// DeleteHandleError handles the Delete error response.
func (client *PacketCapturesClient) DeleteHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - Gets a packet capture session by name.
func (client *PacketCapturesClient) Get(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetOptions) (*PacketCaptureResultResponse, error) {
	req, err := client.GetCreateRequest(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetHandleError(resp)
	}
	result, err := client.GetHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCreateRequest creates the Get request.
func (client *PacketCapturesClient) GetCreateRequest(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/packetCaptures/{packetCaptureName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{networkWatcherName}", url.PathEscape(networkWatcherName))
	urlPath = strings.ReplaceAll(urlPath, "{packetCaptureName}", url.PathEscape(packetCaptureName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetHandleResponse handles the Get response.
func (client *PacketCapturesClient) GetHandleResponse(resp *azcore.Response) (*PacketCaptureResultResponse, error) {
	result := PacketCaptureResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PacketCaptureResult)
}

// GetHandleError handles the Get error response.
func (client *PacketCapturesClient) GetHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

func (client *PacketCapturesClient) BeginGetStatus(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetStatusOptions) (*PacketCaptureQueryStatusResultPollerResponse, error) {
	resp, err := client.GetStatus(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	result := &PacketCaptureQueryStatusResultPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("PacketCapturesClient.GetStatus", "location", resp, client.GetStatusHandleError)
	if err != nil {
		return nil, err
	}
	poller := &packetCaptureQueryStatusResultPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*PacketCaptureQueryStatusResultResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *PacketCapturesClient) ResumeGetStatus(token string) (PacketCaptureQueryStatusResultPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("PacketCapturesClient.GetStatus", token, client.GetStatusHandleError)
	if err != nil {
		return nil, err
	}
	return &packetCaptureQueryStatusResultPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// GetStatus - Query the status of a running packet capture session.
func (client *PacketCapturesClient) GetStatus(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetStatusOptions) (*azcore.Response, error) {
	req, err := client.GetStatusCreateRequest(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.GetStatusHandleError(resp)
	}
	return resp, nil
}

// GetStatusCreateRequest creates the GetStatus request.
func (client *PacketCapturesClient) GetStatusCreateRequest(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesGetStatusOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/packetCaptures/{packetCaptureName}/queryStatus"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{networkWatcherName}", url.PathEscape(networkWatcherName))
	urlPath = strings.ReplaceAll(urlPath, "{packetCaptureName}", url.PathEscape(packetCaptureName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetStatusHandleResponse handles the GetStatus response.
func (client *PacketCapturesClient) GetStatusHandleResponse(resp *azcore.Response) (*PacketCaptureQueryStatusResultResponse, error) {
	result := PacketCaptureQueryStatusResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PacketCaptureQueryStatusResult)
}

// GetStatusHandleError handles the GetStatus error response.
func (client *PacketCapturesClient) GetStatusHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// List - Lists all packet capture sessions within the specified resource group.
func (client *PacketCapturesClient) List(ctx context.Context, resourceGroupName string, networkWatcherName string, options *PacketCapturesListOptions) (*PacketCaptureListResultResponse, error) {
	req, err := client.ListCreateRequest(ctx, resourceGroupName, networkWatcherName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListHandleError(resp)
	}
	result, err := client.ListHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListCreateRequest creates the List request.
func (client *PacketCapturesClient) ListCreateRequest(ctx context.Context, resourceGroupName string, networkWatcherName string, options *PacketCapturesListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/packetCaptures"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{networkWatcherName}", url.PathEscape(networkWatcherName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListHandleResponse handles the List response.
func (client *PacketCapturesClient) ListHandleResponse(resp *azcore.Response) (*PacketCaptureListResultResponse, error) {
	result := PacketCaptureListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.PacketCaptureListResult)
}

// ListHandleError handles the List error response.
func (client *PacketCapturesClient) ListHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

func (client *PacketCapturesClient) BeginStop(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesStopOptions) (*HTTPPollerResponse, error) {
	resp, err := client.Stop(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	result := &HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("PacketCapturesClient.Stop", "location", resp, client.StopHandleError)
	if err != nil {
		return nil, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *PacketCapturesClient) ResumeStop(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("PacketCapturesClient.Stop", token, client.StopHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// Stop - Stops a specified packet capture session.
func (client *PacketCapturesClient) Stop(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesStopOptions) (*azcore.Response, error) {
	req, err := client.StopCreateRequest(ctx, resourceGroupName, networkWatcherName, packetCaptureName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.StopHandleError(resp)
	}
	return resp, nil
}

// StopCreateRequest creates the Stop request.
func (client *PacketCapturesClient) StopCreateRequest(ctx context.Context, resourceGroupName string, networkWatcherName string, packetCaptureName string, options *PacketCapturesStopOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkWatchers/{networkWatcherName}/packetCaptures/{packetCaptureName}/stop"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{networkWatcherName}", url.PathEscape(networkWatcherName))
	urlPath = strings.ReplaceAll(urlPath, "{packetCaptureName}", url.PathEscape(packetCaptureName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// StopHandleError handles the Stop error response.
func (client *PacketCapturesClient) StopHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

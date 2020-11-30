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

// SnapshotsClient contains the methods for the Snapshots group.
// Don't use this type directly, use NewSnapshotsClient() instead.
type SnapshotsClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewSnapshotsClient creates a new instance of SnapshotsClient with the specified values.
func NewSnapshotsClient(con *armcore.Connection, subscriptionID string) SnapshotsClient {
	return SnapshotsClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client SnapshotsClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// BeginCreateOrUpdate - Creates or updates a snapshot.
func (client SnapshotsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, snapshotName string, snapshot Snapshot, options *SnapshotsCreateOrUpdateOptions) (SnapshotPollerResponse, error) {
	resp, err := client.CreateOrUpdate(ctx, resourceGroupName, snapshotName, snapshot, options)
	if err != nil {
		return SnapshotPollerResponse{}, err
	}
	result := SnapshotPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("SnapshotsClient.CreateOrUpdate", "", resp, client.createOrUpdateHandleError)
	if err != nil {
		return SnapshotPollerResponse{}, err
	}
	poller := &snapshotPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (SnapshotResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeCreateOrUpdate creates a new SnapshotPoller from the specified resume token.
// token - The value must come from a previous call to SnapshotPoller.ResumeToken().
func (client SnapshotsClient) ResumeCreateOrUpdate(token string) (SnapshotPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("SnapshotsClient.CreateOrUpdate", token, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &snapshotPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// CreateOrUpdate - Creates or updates a snapshot.
func (client SnapshotsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, snapshotName string, snapshot Snapshot, options *SnapshotsCreateOrUpdateOptions) (*azcore.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, snapshotName, snapshot, options)
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
func (client SnapshotsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, snapshotName string, snapshot Snapshot, options *SnapshotsCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{snapshotName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{snapshotName}", url.PathEscape(snapshotName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(snapshot)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client SnapshotsClient) createOrUpdateHandleResponse(resp *azcore.Response) (SnapshotResponse, error) {
	result := SnapshotResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.Snapshot)
	return result, err
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client SnapshotsClient) createOrUpdateHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// BeginDelete - Deletes a snapshot.
func (client SnapshotsClient) BeginDelete(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsDeleteOptions) (HTTPPollerResponse, error) {
	resp, err := client.Delete(ctx, resourceGroupName, snapshotName, options)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("SnapshotsClient.Delete", "", resp, client.deleteHandleError)
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
func (client SnapshotsClient) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("SnapshotsClient.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Delete - Deletes a snapshot.
func (client SnapshotsClient) Delete(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsDeleteOptions) (*azcore.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, snapshotName, options)
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
func (client SnapshotsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{snapshotName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{snapshotName}", url.PathEscape(snapshotName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client SnapshotsClient) deleteHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Get - Gets information about a snapshot.
func (client SnapshotsClient) Get(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsGetOptions) (SnapshotResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, snapshotName, options)
	if err != nil {
		return SnapshotResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return SnapshotResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return SnapshotResponse{}, client.getHandleError(resp)
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return SnapshotResponse{}, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client SnapshotsClient) getCreateRequest(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{snapshotName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{snapshotName}", url.PathEscape(snapshotName))
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
func (client SnapshotsClient) getHandleResponse(resp *azcore.Response) (SnapshotResponse, error) {
	result := SnapshotResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.Snapshot)
	return result, err
}

// getHandleError handles the Get error response.
func (client SnapshotsClient) getHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// BeginGrantAccess - Grants access to a snapshot.
func (client SnapshotsClient) BeginGrantAccess(ctx context.Context, resourceGroupName string, snapshotName string, grantAccessData GrantAccessData, options *SnapshotsGrantAccessOptions) (AccessURIPollerResponse, error) {
	resp, err := client.GrantAccess(ctx, resourceGroupName, snapshotName, grantAccessData, options)
	if err != nil {
		return AccessURIPollerResponse{}, err
	}
	result := AccessURIPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("SnapshotsClient.GrantAccess", "location", resp, client.grantAccessHandleError)
	if err != nil {
		return AccessURIPollerResponse{}, err
	}
	poller := &accessUriPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (AccessURIResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeGrantAccess creates a new AccessURIPoller from the specified resume token.
// token - The value must come from a previous call to AccessURIPoller.ResumeToken().
func (client SnapshotsClient) ResumeGrantAccess(token string) (AccessURIPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("SnapshotsClient.GrantAccess", token, client.grantAccessHandleError)
	if err != nil {
		return nil, err
	}
	return &accessUriPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// GrantAccess - Grants access to a snapshot.
func (client SnapshotsClient) GrantAccess(ctx context.Context, resourceGroupName string, snapshotName string, grantAccessData GrantAccessData, options *SnapshotsGrantAccessOptions) (*azcore.Response, error) {
	req, err := client.grantAccessCreateRequest(ctx, resourceGroupName, snapshotName, grantAccessData, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.grantAccessHandleError(resp)
	}
	return resp, nil
}

// grantAccessCreateRequest creates the GrantAccess request.
func (client SnapshotsClient) grantAccessCreateRequest(ctx context.Context, resourceGroupName string, snapshotName string, grantAccessData GrantAccessData, options *SnapshotsGrantAccessOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{snapshotName}/beginGetAccess"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{snapshotName}", url.PathEscape(snapshotName))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(grantAccessData)
}

// grantAccessHandleResponse handles the GrantAccess response.
func (client SnapshotsClient) grantAccessHandleResponse(resp *azcore.Response) (AccessURIResponse, error) {
	result := AccessURIResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AccessURI)
	return result, err
}

// grantAccessHandleError handles the GrantAccess error response.
func (client SnapshotsClient) grantAccessHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// List - Lists snapshots under a subscription.
func (client SnapshotsClient) List(options *SnapshotsListOptions) SnapshotListPager {
	return &snapshotListPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp SnapshotListResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.SnapshotList.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client SnapshotsClient) listCreateRequest(ctx context.Context, options *SnapshotsListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/snapshots"
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
func (client SnapshotsClient) listHandleResponse(resp *azcore.Response) (SnapshotListResponse, error) {
	result := SnapshotListResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.SnapshotList)
	return result, err
}

// listHandleError handles the List error response.
func (client SnapshotsClient) listHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListByResourceGroup - Lists snapshots under a resource group.
func (client SnapshotsClient) ListByResourceGroup(resourceGroupName string, options *SnapshotsListByResourceGroupOptions) SnapshotListPager {
	return &snapshotListPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
		},
		responder: client.listByResourceGroupHandleResponse,
		errorer:   client.listByResourceGroupHandleError,
		advancer: func(ctx context.Context, resp SnapshotListResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.SnapshotList.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client SnapshotsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *SnapshotsListByResourceGroupOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots"
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
func (client SnapshotsClient) listByResourceGroupHandleResponse(resp *azcore.Response) (SnapshotListResponse, error) {
	result := SnapshotListResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.SnapshotList)
	return result, err
}

// listByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client SnapshotsClient) listByResourceGroupHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// BeginRevokeAccess - Revokes access to a snapshot.
func (client SnapshotsClient) BeginRevokeAccess(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsRevokeAccessOptions) (HTTPPollerResponse, error) {
	resp, err := client.RevokeAccess(ctx, resourceGroupName, snapshotName, options)
	if err != nil {
		return HTTPPollerResponse{}, err
	}
	result := HTTPPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("SnapshotsClient.RevokeAccess", "location", resp, client.revokeAccessHandleError)
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

// ResumeRevokeAccess creates a new HTTPPoller from the specified resume token.
// token - The value must come from a previous call to HTTPPoller.ResumeToken().
func (client SnapshotsClient) ResumeRevokeAccess(token string) (HTTPPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("SnapshotsClient.RevokeAccess", token, client.revokeAccessHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// RevokeAccess - Revokes access to a snapshot.
func (client SnapshotsClient) RevokeAccess(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsRevokeAccessOptions) (*azcore.Response, error) {
	req, err := client.revokeAccessCreateRequest(ctx, resourceGroupName, snapshotName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.revokeAccessHandleError(resp)
	}
	return resp, nil
}

// revokeAccessCreateRequest creates the RevokeAccess request.
func (client SnapshotsClient) revokeAccessCreateRequest(ctx context.Context, resourceGroupName string, snapshotName string, options *SnapshotsRevokeAccessOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{snapshotName}/endGetAccess"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{snapshotName}", url.PathEscape(snapshotName))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// revokeAccessHandleError handles the RevokeAccess error response.
func (client SnapshotsClient) revokeAccessHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// BeginUpdate - Updates (patches) a snapshot.
func (client SnapshotsClient) BeginUpdate(ctx context.Context, resourceGroupName string, snapshotName string, snapshot SnapshotUpdate, options *SnapshotsUpdateOptions) (SnapshotPollerResponse, error) {
	resp, err := client.Update(ctx, resourceGroupName, snapshotName, snapshot, options)
	if err != nil {
		return SnapshotPollerResponse{}, err
	}
	result := SnapshotPollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("SnapshotsClient.Update", "", resp, client.updateHandleError)
	if err != nil {
		return SnapshotPollerResponse{}, err
	}
	poller := &snapshotPoller{
		pt:       pt,
		pipeline: client.con.Pipeline(),
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (SnapshotResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

// ResumeUpdate creates a new SnapshotPoller from the specified resume token.
// token - The value must come from a previous call to SnapshotPoller.ResumeToken().
func (client SnapshotsClient) ResumeUpdate(token string) (SnapshotPoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("SnapshotsClient.Update", token, client.updateHandleError)
	if err != nil {
		return nil, err
	}
	return &snapshotPoller{
		pipeline: client.con.Pipeline(),
		pt:       pt,
	}, nil
}

// Update - Updates (patches) a snapshot.
func (client SnapshotsClient) Update(ctx context.Context, resourceGroupName string, snapshotName string, snapshot SnapshotUpdate, options *SnapshotsUpdateOptions) (*azcore.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, snapshotName, snapshot, options)
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
func (client SnapshotsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, snapshotName string, snapshot SnapshotUpdate, options *SnapshotsUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/snapshots/{snapshotName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{snapshotName}", url.PathEscape(snapshotName))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-30")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(snapshot)
}

// updateHandleResponse handles the Update response.
func (client SnapshotsClient) updateHandleResponse(resp *azcore.Response) (SnapshotResponse, error) {
	result := SnapshotResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.Snapshot)
	return result, err
}

// updateHandleError handles the Update error response.
func (client SnapshotsClient) updateHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

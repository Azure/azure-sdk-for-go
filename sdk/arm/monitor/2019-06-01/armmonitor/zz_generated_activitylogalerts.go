// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// ActivityLogAlertsOperations contains the methods for the ActivityLogAlerts group.
type ActivityLogAlertsOperations interface {
	// CreateOrUpdate - Create a new activity log alert or update an existing one.
	CreateOrUpdate(ctx context.Context, resourceGroupName string, activityLogAlertName string, activityLogAlert ActivityLogAlertResource, options *ActivityLogAlertsCreateOrUpdateOptions) (*ActivityLogAlertResourceResponse, error)
	// Delete - Delete an activity log alert.
	Delete(ctx context.Context, resourceGroupName string, activityLogAlertName string, options *ActivityLogAlertsDeleteOptions) (*http.Response, error)
	// Get - Get an activity log alert.
	Get(ctx context.Context, resourceGroupName string, activityLogAlertName string, options *ActivityLogAlertsGetOptions) (*ActivityLogAlertResourceResponse, error)
	// ListByResourceGroup - Get a list of all activity log alerts in a resource group.
	ListByResourceGroup(ctx context.Context, resourceGroupName string, options *ActivityLogAlertsListByResourceGroupOptions) (*ActivityLogAlertListResponse, error)
	// ListBySubscriptionID - Get a list of all activity log alerts in a subscription.
	ListBySubscriptionID(ctx context.Context, options *ActivityLogAlertsListBySubscriptionIDOptions) (*ActivityLogAlertListResponse, error)
	// Update - Updates an existing ActivityLogAlertResource's tags. To update other fields use the CreateOrUpdate method.
	Update(ctx context.Context, resourceGroupName string, activityLogAlertName string, activityLogAlertPatch ActivityLogAlertPatchBody, options *ActivityLogAlertsUpdateOptions) (*ActivityLogAlertResourceResponse, error)
}

// ActivityLogAlertsClient implements the ActivityLogAlertsOperations interface.
// Don't use this type directly, use NewActivityLogAlertsClient() instead.
type ActivityLogAlertsClient struct {
	*Client
	subscriptionID string
}

// NewActivityLogAlertsClient creates a new instance of ActivityLogAlertsClient with the specified values.
func NewActivityLogAlertsClient(c *Client, subscriptionID string) ActivityLogAlertsOperations {
	return &ActivityLogAlertsClient{Client: c, subscriptionID: subscriptionID}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *ActivityLogAlertsClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

// CreateOrUpdate - Create a new activity log alert or update an existing one.
func (client *ActivityLogAlertsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, activityLogAlertName string, activityLogAlert ActivityLogAlertResource, options *ActivityLogAlertsCreateOrUpdateOptions) (*ActivityLogAlertResourceResponse, error) {
	req, err := client.CreateOrUpdateCreateRequest(ctx, resourceGroupName, activityLogAlertName, activityLogAlert, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated) {
		return nil, client.CreateOrUpdateHandleError(resp)
	}
	result, err := client.CreateOrUpdateHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ActivityLogAlertsClient) CreateOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, activityLogAlertName string, activityLogAlert ActivityLogAlertResource, options *ActivityLogAlertsCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/activityLogAlerts/{activityLogAlertName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{activityLogAlertName}", url.PathEscape(activityLogAlertName))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-04-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(activityLogAlert)
}

// CreateOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ActivityLogAlertsClient) CreateOrUpdateHandleResponse(resp *azcore.Response) (*ActivityLogAlertResourceResponse, error) {
	result := ActivityLogAlertResourceResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ActivityLogAlertResource)
}

// CreateOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *ActivityLogAlertsClient) CreateOrUpdateHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Delete - Delete an activity log alert.
func (client *ActivityLogAlertsClient) Delete(ctx context.Context, resourceGroupName string, activityLogAlertName string, options *ActivityLogAlertsDeleteOptions) (*http.Response, error) {
	req, err := client.DeleteCreateRequest(ctx, resourceGroupName, activityLogAlertName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusNoContent) {
		return nil, client.DeleteHandleError(resp)
	}
	return resp.Response, nil
}

// DeleteCreateRequest creates the Delete request.
func (client *ActivityLogAlertsClient) DeleteCreateRequest(ctx context.Context, resourceGroupName string, activityLogAlertName string, options *ActivityLogAlertsDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/activityLogAlerts/{activityLogAlertName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{activityLogAlertName}", url.PathEscape(activityLogAlertName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-04-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// DeleteHandleError handles the Delete error response.
func (client *ActivityLogAlertsClient) DeleteHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Get - Get an activity log alert.
func (client *ActivityLogAlertsClient) Get(ctx context.Context, resourceGroupName string, activityLogAlertName string, options *ActivityLogAlertsGetOptions) (*ActivityLogAlertResourceResponse, error) {
	req, err := client.GetCreateRequest(ctx, resourceGroupName, activityLogAlertName, options)
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
func (client *ActivityLogAlertsClient) GetCreateRequest(ctx context.Context, resourceGroupName string, activityLogAlertName string, options *ActivityLogAlertsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/activityLogAlerts/{activityLogAlertName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{activityLogAlertName}", url.PathEscape(activityLogAlertName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-04-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// GetHandleResponse handles the Get response.
func (client *ActivityLogAlertsClient) GetHandleResponse(resp *azcore.Response) (*ActivityLogAlertResourceResponse, error) {
	result := ActivityLogAlertResourceResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ActivityLogAlertResource)
}

// GetHandleError handles the Get error response.
func (client *ActivityLogAlertsClient) GetHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// ListByResourceGroup - Get a list of all activity log alerts in a resource group.
func (client *ActivityLogAlertsClient) ListByResourceGroup(ctx context.Context, resourceGroupName string, options *ActivityLogAlertsListByResourceGroupOptions) (*ActivityLogAlertListResponse, error) {
	req, err := client.ListByResourceGroupCreateRequest(ctx, resourceGroupName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListByResourceGroupHandleError(resp)
	}
	result, err := client.ListByResourceGroupHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ActivityLogAlertsClient) ListByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ActivityLogAlertsListByResourceGroupOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/activityLogAlerts"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-04-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ActivityLogAlertsClient) ListByResourceGroupHandleResponse(resp *azcore.Response) (*ActivityLogAlertListResponse, error) {
	result := ActivityLogAlertListResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ActivityLogAlertList)
}

// ListByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client *ActivityLogAlertsClient) ListByResourceGroupHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// ListBySubscriptionID - Get a list of all activity log alerts in a subscription.
func (client *ActivityLogAlertsClient) ListBySubscriptionID(ctx context.Context, options *ActivityLogAlertsListBySubscriptionIDOptions) (*ActivityLogAlertListResponse, error) {
	req, err := client.ListBySubscriptionIDCreateRequest(ctx, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.ListBySubscriptionIDHandleError(resp)
	}
	result, err := client.ListBySubscriptionIDHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListBySubscriptionIDCreateRequest creates the ListBySubscriptionID request.
func (client *ActivityLogAlertsClient) ListBySubscriptionIDCreateRequest(ctx context.Context, options *ActivityLogAlertsListBySubscriptionIDOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/microsoft.insights/activityLogAlerts"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-04-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListBySubscriptionIDHandleResponse handles the ListBySubscriptionID response.
func (client *ActivityLogAlertsClient) ListBySubscriptionIDHandleResponse(resp *azcore.Response) (*ActivityLogAlertListResponse, error) {
	result := ActivityLogAlertListResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ActivityLogAlertList)
}

// ListBySubscriptionIDHandleError handles the ListBySubscriptionID error response.
func (client *ActivityLogAlertsClient) ListBySubscriptionIDHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Update - Updates an existing ActivityLogAlertResource's tags. To update other fields use the CreateOrUpdate method.
func (client *ActivityLogAlertsClient) Update(ctx context.Context, resourceGroupName string, activityLogAlertName string, activityLogAlertPatch ActivityLogAlertPatchBody, options *ActivityLogAlertsUpdateOptions) (*ActivityLogAlertResourceResponse, error) {
	req, err := client.UpdateCreateRequest(ctx, resourceGroupName, activityLogAlertName, activityLogAlertPatch, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.UpdateHandleError(resp)
	}
	result, err := client.UpdateHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateCreateRequest creates the Update request.
func (client *ActivityLogAlertsClient) UpdateCreateRequest(ctx context.Context, resourceGroupName string, activityLogAlertName string, activityLogAlertPatch ActivityLogAlertPatchBody, options *ActivityLogAlertsUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/activityLogAlerts/{activityLogAlertName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{activityLogAlertName}", url.PathEscape(activityLogAlertName))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2017-04-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(activityLogAlertPatch)
}

// UpdateHandleResponse handles the Update response.
func (client *ActivityLogAlertsClient) UpdateHandleResponse(resp *azcore.Response) (*ActivityLogAlertResourceResponse, error) {
	result := ActivityLogAlertResourceResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.ActivityLogAlertResource)
}

// UpdateHandleError handles the Update error response.
func (client *ActivityLogAlertsClient) UpdateHandleError(resp *azcore.Response) error {
	var err ErrorResponse
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

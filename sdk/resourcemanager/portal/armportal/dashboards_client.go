//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armportal

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

// DashboardsClient contains the methods for the Dashboards group.
// Don't use this type directly, use NewDashboardsClient() instead.
type DashboardsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewDashboardsClient creates a new instance of DashboardsClient with the specified values.
//   - subscriptionID - The Azure subscription ID. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000)
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewDashboardsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DashboardsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &DashboardsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a Dashboard.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01-preview
//   - resourceGroupName - The name of the resource group.
//   - dashboardName - The name of the dashboard.
//   - dashboard - The parameters required to create or update a dashboard.
//   - options - DashboardsClientCreateOrUpdateOptions contains the optional parameters for the DashboardsClient.CreateOrUpdate
//     method.
func (client *DashboardsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, dashboardName string, dashboard Dashboard, options *DashboardsClientCreateOrUpdateOptions) (DashboardsClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "DashboardsClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, dashboardName, dashboard, options)
	if err != nil {
		return DashboardsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DashboardsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return DashboardsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *DashboardsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, dashboardName string, dashboard Dashboard, options *DashboardsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Portal/dashboards/{dashboardName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dashboardName == "" {
		return nil, errors.New("parameter dashboardName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dashboardName}", url.PathEscape(dashboardName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, dashboard); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *DashboardsClient) createOrUpdateHandleResponse(resp *http.Response) (DashboardsClientCreateOrUpdateResponse, error) {
	result := DashboardsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dashboard); err != nil {
		return DashboardsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the Dashboard.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01-preview
//   - resourceGroupName - The name of the resource group.
//   - dashboardName - The name of the dashboard.
//   - options - DashboardsClientDeleteOptions contains the optional parameters for the DashboardsClient.Delete method.
func (client *DashboardsClient) Delete(ctx context.Context, resourceGroupName string, dashboardName string, options *DashboardsClientDeleteOptions) (DashboardsClientDeleteResponse, error) {
	var err error
	const operationName = "DashboardsClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, dashboardName, options)
	if err != nil {
		return DashboardsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DashboardsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return DashboardsClientDeleteResponse{}, err
	}
	return DashboardsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DashboardsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, dashboardName string, options *DashboardsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Portal/dashboards/{dashboardName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dashboardName == "" {
		return nil, errors.New("parameter dashboardName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dashboardName}", url.PathEscape(dashboardName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the Dashboard.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01-preview
//   - resourceGroupName - The name of the resource group.
//   - dashboardName - The name of the dashboard.
//   - options - DashboardsClientGetOptions contains the optional parameters for the DashboardsClient.Get method.
func (client *DashboardsClient) Get(ctx context.Context, resourceGroupName string, dashboardName string, options *DashboardsClientGetOptions) (DashboardsClientGetResponse, error) {
	var err error
	const operationName = "DashboardsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, dashboardName, options)
	if err != nil {
		return DashboardsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DashboardsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNotFound) {
		err = runtime.NewResponseError(httpResp)
		return DashboardsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *DashboardsClient) getCreateRequest(ctx context.Context, resourceGroupName string, dashboardName string, options *DashboardsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Portal/dashboards/{dashboardName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dashboardName == "" {
		return nil, errors.New("parameter dashboardName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dashboardName}", url.PathEscape(dashboardName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DashboardsClient) getHandleResponse(resp *http.Response) (DashboardsClientGetResponse, error) {
	result := DashboardsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dashboard); err != nil {
		return DashboardsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Gets all the Dashboards within a resource group.
//
// Generated from API version 2020-09-01-preview
//   - resourceGroupName - The name of the resource group.
//   - options - DashboardsClientListByResourceGroupOptions contains the optional parameters for the DashboardsClient.NewListByResourceGroupPager
//     method.
func (client *DashboardsClient) NewListByResourceGroupPager(resourceGroupName string, options *DashboardsClientListByResourceGroupOptions) *runtime.Pager[DashboardsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[DashboardsClientListByResourceGroupResponse]{
		More: func(page DashboardsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DashboardsClientListByResourceGroupResponse) (DashboardsClientListByResourceGroupResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "DashboardsClient.NewListByResourceGroupPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			}, nil)
			if err != nil {
				return DashboardsClientListByResourceGroupResponse{}, err
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *DashboardsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *DashboardsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Portal/dashboards"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *DashboardsClient) listByResourceGroupHandleResponse(resp *http.Response) (DashboardsClientListByResourceGroupResponse, error) {
	result := DashboardsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DashboardListResult); err != nil {
		return DashboardsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - Gets all the dashboards within a subscription.
//
// Generated from API version 2020-09-01-preview
//   - options - DashboardsClientListBySubscriptionOptions contains the optional parameters for the DashboardsClient.NewListBySubscriptionPager
//     method.
func (client *DashboardsClient) NewListBySubscriptionPager(options *DashboardsClientListBySubscriptionOptions) *runtime.Pager[DashboardsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[DashboardsClientListBySubscriptionResponse]{
		More: func(page DashboardsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DashboardsClientListBySubscriptionResponse) (DashboardsClientListBySubscriptionResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "DashboardsClient.NewListBySubscriptionPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listBySubscriptionCreateRequest(ctx, options)
			}, nil)
			if err != nil {
				return DashboardsClientListBySubscriptionResponse{}, err
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *DashboardsClient) listBySubscriptionCreateRequest(ctx context.Context, options *DashboardsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Portal/dashboards"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *DashboardsClient) listBySubscriptionHandleResponse(resp *http.Response) (DashboardsClientListBySubscriptionResponse, error) {
	result := DashboardsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DashboardListResult); err != nil {
		return DashboardsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Updates an existing Dashboard.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01-preview
//   - resourceGroupName - The name of the resource group.
//   - dashboardName - The name of the dashboard.
//   - dashboard - The updatable fields of a Dashboard.
//   - options - DashboardsClientUpdateOptions contains the optional parameters for the DashboardsClient.Update method.
func (client *DashboardsClient) Update(ctx context.Context, resourceGroupName string, dashboardName string, dashboard PatchableDashboard, options *DashboardsClientUpdateOptions) (DashboardsClientUpdateResponse, error) {
	var err error
	const operationName = "DashboardsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, dashboardName, dashboard, options)
	if err != nil {
		return DashboardsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DashboardsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNotFound) {
		err = runtime.NewResponseError(httpResp)
		return DashboardsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *DashboardsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, dashboardName string, dashboard PatchableDashboard, options *DashboardsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Portal/dashboards/{dashboardName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dashboardName == "" {
		return nil, errors.New("parameter dashboardName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dashboardName}", url.PathEscape(dashboardName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, dashboard); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *DashboardsClient) updateHandleResponse(resp *http.Response) (DashboardsClientUpdateResponse, error) {
	result := DashboardsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Dashboard); err != nil {
		return DashboardsClientUpdateResponse{}, err
	}
	return result, nil
}

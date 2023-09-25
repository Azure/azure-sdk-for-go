//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurityinsights

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

// WatchlistsClient contains the methods for the Watchlists group.
// Don't use this type directly, use NewWatchlistsClient() instead.
type WatchlistsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewWatchlistsClient creates a new instance of WatchlistsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWatchlistsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WatchlistsClient, error) {
	cl, err := arm.NewClient(moduleName+".WatchlistsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WatchlistsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a Watchlist and its Watchlist Items (bulk creation, e.g. through text/csv content type).
// To create a Watchlist and its Items, we should call this endpoint with rawContent and
// contentType properties.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-10-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - watchlistAlias - The watchlist alias
//   - watchlist - The watchlist
//   - options - WatchlistsClientCreateOrUpdateOptions contains the optional parameters for the WatchlistsClient.CreateOrUpdate
//     method.
func (client *WatchlistsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, watchlist Watchlist, options *WatchlistsClientCreateOrUpdateOptions) (WatchlistsClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, watchlistAlias, watchlist, options)
	if err != nil {
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *WatchlistsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, watchlist Watchlist, options *WatchlistsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/watchlists/{watchlistAlias}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if watchlistAlias == "" {
		return nil, errors.New("parameter watchlistAlias cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{watchlistAlias}", url.PathEscape(watchlistAlias))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, watchlist); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *WatchlistsClient) createOrUpdateHandleResponse(resp *http.Response) (WatchlistsClientCreateOrUpdateResponse, error) {
	result := WatchlistsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Watchlist); err != nil {
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a watchlist.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-10-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - watchlistAlias - The watchlist alias
//   - options - WatchlistsClientDeleteOptions contains the optional parameters for the WatchlistsClient.Delete method.
func (client *WatchlistsClient) Delete(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, options *WatchlistsClientDeleteOptions) (WatchlistsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, watchlistAlias, options)
	if err != nil {
		return WatchlistsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WatchlistsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return WatchlistsClientDeleteResponse{}, err
	}
	return WatchlistsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *WatchlistsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, options *WatchlistsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/watchlists/{watchlistAlias}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if watchlistAlias == "" {
		return nil, errors.New("parameter watchlistAlias cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{watchlistAlias}", url.PathEscape(watchlistAlias))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a watchlist, without its watchlist items.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-10-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - watchlistAlias - The watchlist alias
//   - options - WatchlistsClientGetOptions contains the optional parameters for the WatchlistsClient.Get method.
func (client *WatchlistsClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, options *WatchlistsClientGetOptions) (WatchlistsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, watchlistAlias, options)
	if err != nil {
		return WatchlistsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WatchlistsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WatchlistsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *WatchlistsClient) getCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, options *WatchlistsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/watchlists/{watchlistAlias}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if watchlistAlias == "" {
		return nil, errors.New("parameter watchlistAlias cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{watchlistAlias}", url.PathEscape(watchlistAlias))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WatchlistsClient) getHandleResponse(resp *http.Response) (WatchlistsClientGetResponse, error) {
	result := WatchlistsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Watchlist); err != nil {
		return WatchlistsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Get all watchlists, without watchlist items.
//
// Generated from API version 2021-10-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - options - WatchlistsClientListOptions contains the optional parameters for the WatchlistsClient.NewListPager method.
func (client *WatchlistsClient) NewListPager(resourceGroupName string, workspaceName string, options *WatchlistsClientListOptions) (*runtime.Pager[WatchlistsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[WatchlistsClientListResponse]{
		More: func(page WatchlistsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WatchlistsClientListResponse) (WatchlistsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, workspaceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return WatchlistsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WatchlistsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WatchlistsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *WatchlistsClient) listCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, options *WatchlistsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{workspaceName}/providers/Microsoft.SecurityInsights/watchlists"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *WatchlistsClient) listHandleResponse(resp *http.Response) (WatchlistsClientListResponse, error) {
	result := WatchlistsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WatchlistList); err != nil {
		return WatchlistsClientListResponse{}, err
	}
	return result, nil
}


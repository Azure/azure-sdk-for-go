//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurityinsights

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// WatchlistsClient contains the methods for the Watchlists group.
// Don't use this type directly, use NewWatchlistsClient() instead.
type WatchlistsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewWatchlistsClient creates a new instance of WatchlistsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewWatchlistsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WatchlistsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &WatchlistsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a Watchlist and its Watchlist Items (bulk creation, e.g. through text/csv content type).
// To create a Watchlist and its Items, we should call this endpoint with either rawContent or a
// valid SAR URI and contentType properties. The rawContent is mainly used for small watchlist (content size below 3.8 MB).
// The SAS URI enables the creation of large watchlist, where the content size can
// go up to 500 MB. The status of processing such large file can be polled through the URL returned in Azure-AsyncOperation
// header.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// watchlistAlias - Watchlist Alias
// watchlist - The watchlist
// options - WatchlistsClientCreateOrUpdateOptions contains the optional parameters for the WatchlistsClient.CreateOrUpdate
// method.
func (client *WatchlistsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, watchlist Watchlist, options *WatchlistsClientCreateOrUpdateOptions) (WatchlistsClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, workspaceName, watchlistAlias, watchlist, options)
	if err != nil {
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return WatchlistsClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, watchlist)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *WatchlistsClient) createOrUpdateHandleResponse(resp *http.Response) (WatchlistsClientCreateOrUpdateResponse, error) {
	result := WatchlistsClientCreateOrUpdateResponse{}
	if val := resp.Header.Get("Azure-AsyncOperation"); val != "" {
		result.AzureAsyncOperation = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Watchlist); err != nil {
		return WatchlistsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a watchlist.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// watchlistAlias - Watchlist Alias
// options - WatchlistsClientDeleteOptions contains the optional parameters for the WatchlistsClient.Delete method.
func (client *WatchlistsClient) Delete(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, options *WatchlistsClientDeleteOptions) (WatchlistsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, workspaceName, watchlistAlias, options)
	if err != nil {
		return WatchlistsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return WatchlistsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return WatchlistsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return client.deleteHandleResponse(resp)
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
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *WatchlistsClient) deleteHandleResponse(resp *http.Response) (WatchlistsClientDeleteResponse, error) {
	result := WatchlistsClientDeleteResponse{}
	if val := resp.Header.Get("Azure-AsyncOperation"); val != "" {
		result.AzureAsyncOperation = &val
	}
	return result, nil
}

// Get - Gets a watchlist, without its watchlist items.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// watchlistAlias - Watchlist Alias
// options - WatchlistsClientGetOptions contains the optional parameters for the WatchlistsClient.Get method.
func (client *WatchlistsClient) Get(ctx context.Context, resourceGroupName string, workspaceName string, watchlistAlias string, options *WatchlistsClientGetOptions) (WatchlistsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, workspaceName, watchlistAlias, options)
	if err != nil {
		return WatchlistsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return WatchlistsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return WatchlistsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
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

// List - Gets all watchlists, without watchlist items.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// workspaceName - The name of the workspace.
// options - WatchlistsClientListOptions contains the optional parameters for the WatchlistsClient.List method.
func (client *WatchlistsClient) List(resourceGroupName string, workspaceName string, options *WatchlistsClientListOptions) *runtime.Pager[WatchlistsClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[WatchlistsClientListResponse]{
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
			resp, err := client.pl.Do(req)
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
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

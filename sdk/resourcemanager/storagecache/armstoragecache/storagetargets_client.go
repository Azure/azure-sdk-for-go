//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstoragecache

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

// StorageTargetsClient contains the methods for the StorageTargets group.
// Don't use this type directly, use NewStorageTargetsClient() instead.
type StorageTargetsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewStorageTargetsClient creates a new instance of StorageTargetsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewStorageTargetsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*StorageTargetsClient, error) {
	cl, err := arm.NewClient(moduleName+".StorageTargetsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &StorageTargetsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create or update a Storage Target. This operation is allowed at any time, but if the cache is down
// or unhealthy, the actual creation/modification of the Storage Target may be delayed until the cache
// is healthy again.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cacheName - Name of cache. Length of name must not be greater than 80 and chars must be from the [-0-9a-zA-Z_] char class.
//   - storageTargetName - Name of Storage Target.
//   - storagetarget - Object containing the definition of a Storage Target.
//   - options - StorageTargetsClientBeginCreateOrUpdateOptions contains the optional parameters for the StorageTargetsClient.BeginCreateOrUpdate
//     method.
func (client *StorageTargetsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, storagetarget StorageTarget, options *StorageTargetsClientBeginCreateOrUpdateOptions) (*runtime.Poller[StorageTargetsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, cacheName, storageTargetName, storagetarget, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[StorageTargetsClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[StorageTargetsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create or update a Storage Target. This operation is allowed at any time, but if the cache is down or
// unhealthy, the actual creation/modification of the Storage Target may be delayed until the cache
// is healthy again.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *StorageTargetsClient) createOrUpdate(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, storagetarget StorageTarget, options *StorageTargetsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, cacheName, storageTargetName, storagetarget, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *StorageTargetsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, storagetarget StorageTarget, options *StorageTargetsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StorageCache/caches/{cacheName}/storageTargets/{storageTargetName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if cacheName == "" {
		return nil, errors.New("parameter cacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cacheName}", url.PathEscape(cacheName))
	if storageTargetName == "" {
		return nil, errors.New("parameter storageTargetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageTargetName}", url.PathEscape(storageTargetName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, storagetarget); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDNSRefresh - Tells a storage target to refresh its DNS information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cacheName - Name of cache. Length of name must not be greater than 80 and chars must be from the [-0-9a-zA-Z_] char class.
//   - storageTargetName - Name of Storage Target.
//   - options - StorageTargetsClientBeginDNSRefreshOptions contains the optional parameters for the StorageTargetsClient.BeginDNSRefresh
//     method.
func (client *StorageTargetsClient) BeginDNSRefresh(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginDNSRefreshOptions) (*runtime.Poller[StorageTargetsClientDNSRefreshResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.dNSRefresh(ctx, resourceGroupName, cacheName, storageTargetName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[StorageTargetsClientDNSRefreshResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[StorageTargetsClientDNSRefreshResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// DNSRefresh - Tells a storage target to refresh its DNS information.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *StorageTargetsClient) dNSRefresh(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginDNSRefreshOptions) (*http.Response, error) {
	var err error
	req, err := client.dnsRefreshCreateRequest(ctx, resourceGroupName, cacheName, storageTargetName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// dnsRefreshCreateRequest creates the DNSRefresh request.
func (client *StorageTargetsClient) dnsRefreshCreateRequest(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginDNSRefreshOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StorageCache/caches/{cacheName}/storageTargets/{storageTargetName}/dnsRefresh"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if cacheName == "" {
		return nil, errors.New("parameter cacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cacheName}", url.PathEscape(cacheName))
	if storageTargetName == "" {
		return nil, errors.New("parameter storageTargetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageTargetName}", url.PathEscape(storageTargetName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginDelete - Removes a Storage Target from a cache. This operation is allowed at any time, but if the cache is down or
// unhealthy, the actual removal of the Storage Target may be delayed until the cache is healthy
// again. Note that if the cache has data to flush to the Storage Target, the data will be flushed before the Storage Target
// will be deleted.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cacheName - Name of cache. Length of name must not be greater than 80 and chars must be from the [-0-9a-zA-Z_] char class.
//   - storageTargetName - Name of Storage Target.
//   - options - StorageTargetsClientBeginDeleteOptions contains the optional parameters for the StorageTargetsClient.BeginDelete
//     method.
func (client *StorageTargetsClient) BeginDelete(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginDeleteOptions) (*runtime.Poller[StorageTargetsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, cacheName, storageTargetName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[StorageTargetsClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[StorageTargetsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Removes a Storage Target from a cache. This operation is allowed at any time, but if the cache is down or unhealthy,
// the actual removal of the Storage Target may be delayed until the cache is healthy
// again. Note that if the cache has data to flush to the Storage Target, the data will be flushed before the Storage Target
// will be deleted.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *StorageTargetsClient) deleteOperation(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, cacheName, storageTargetName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *StorageTargetsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StorageCache/caches/{cacheName}/storageTargets/{storageTargetName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if cacheName == "" {
		return nil, errors.New("parameter cacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cacheName}", url.PathEscape(cacheName))
	if storageTargetName == "" {
		return nil, errors.New("parameter storageTargetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageTargetName}", url.PathEscape(storageTargetName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	if options != nil && options.Force != nil {
		reqQP.Set("force", *options.Force)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Returns a Storage Target from a cache.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cacheName - Name of cache. Length of name must not be greater than 80 and chars must be from the [-0-9a-zA-Z_] char class.
//   - storageTargetName - Name of Storage Target.
//   - options - StorageTargetsClientGetOptions contains the optional parameters for the StorageTargetsClient.Get method.
func (client *StorageTargetsClient) Get(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientGetOptions) (StorageTargetsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, cacheName, storageTargetName, options)
	if err != nil {
		return StorageTargetsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StorageTargetsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return StorageTargetsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *StorageTargetsClient) getCreateRequest(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StorageCache/caches/{cacheName}/storageTargets/{storageTargetName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if cacheName == "" {
		return nil, errors.New("parameter cacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cacheName}", url.PathEscape(cacheName))
	if storageTargetName == "" {
		return nil, errors.New("parameter storageTargetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageTargetName}", url.PathEscape(storageTargetName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *StorageTargetsClient) getHandleResponse(resp *http.Response) (StorageTargetsClientGetResponse, error) {
	result := StorageTargetsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageTarget); err != nil {
		return StorageTargetsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByCachePager - Returns a list of Storage Targets for the specified cache.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cacheName - Name of cache. Length of name must not be greater than 80 and chars must be from the [-0-9a-zA-Z_] char class.
//   - options - StorageTargetsClientListByCacheOptions contains the optional parameters for the StorageTargetsClient.NewListByCachePager
//     method.
func (client *StorageTargetsClient) NewListByCachePager(resourceGroupName string, cacheName string, options *StorageTargetsClientListByCacheOptions) (*runtime.Pager[StorageTargetsClientListByCacheResponse]) {
	return runtime.NewPager(runtime.PagingHandler[StorageTargetsClientListByCacheResponse]{
		More: func(page StorageTargetsClientListByCacheResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *StorageTargetsClientListByCacheResponse) (StorageTargetsClientListByCacheResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByCacheCreateRequest(ctx, resourceGroupName, cacheName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return StorageTargetsClientListByCacheResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return StorageTargetsClientListByCacheResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return StorageTargetsClientListByCacheResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByCacheHandleResponse(resp)
		},
	})
}

// listByCacheCreateRequest creates the ListByCache request.
func (client *StorageTargetsClient) listByCacheCreateRequest(ctx context.Context, resourceGroupName string, cacheName string, options *StorageTargetsClientListByCacheOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.StorageCache/caches/{cacheName}/storageTargets"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if cacheName == "" {
		return nil, errors.New("parameter cacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cacheName}", url.PathEscape(cacheName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByCacheHandleResponse handles the ListByCache response.
func (client *StorageTargetsClient) listByCacheHandleResponse(resp *http.Response) (StorageTargetsClientListByCacheResponse, error) {
	result := StorageTargetsClientListByCacheResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StorageTargetsResult); err != nil {
		return StorageTargetsClientListByCacheResponse{}, err
	}
	return result, nil
}

// BeginRestoreDefaults - Tells a storage target to restore its settings to their default values.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - cacheName - Name of cache. Length of name must not be greater than 80 and chars must be from the [-0-9a-zA-Z_] char class.
//   - storageTargetName - Name of Storage Target.
//   - options - StorageTargetsClientBeginRestoreDefaultsOptions contains the optional parameters for the StorageTargetsClient.BeginRestoreDefaults
//     method.
func (client *StorageTargetsClient) BeginRestoreDefaults(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginRestoreDefaultsOptions) (*runtime.Poller[StorageTargetsClientRestoreDefaultsResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.restoreDefaults(ctx, resourceGroupName, cacheName, storageTargetName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[StorageTargetsClientRestoreDefaultsResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[StorageTargetsClientRestoreDefaultsResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// RestoreDefaults - Tells a storage target to restore its settings to their default values.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-05-01
func (client *StorageTargetsClient) restoreDefaults(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginRestoreDefaultsOptions) (*http.Response, error) {
	var err error
	req, err := client.restoreDefaultsCreateRequest(ctx, resourceGroupName, cacheName, storageTargetName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// restoreDefaultsCreateRequest creates the RestoreDefaults request.
func (client *StorageTargetsClient) restoreDefaultsCreateRequest(ctx context.Context, resourceGroupName string, cacheName string, storageTargetName string, options *StorageTargetsClientBeginRestoreDefaultsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageCache/caches/{cacheName}/storageTargets/{storageTargetName}/restoreDefaults"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if cacheName == "" {
		return nil, errors.New("parameter cacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{cacheName}", url.PathEscape(cacheName))
	if storageTargetName == "" {
		return nil, errors.New("parameter storageTargetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageTargetName}", url.PathEscape(storageTargetName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}


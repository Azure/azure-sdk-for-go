//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetworkcloud

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

// BmcKeySetsClient contains the methods for the BmcKeySets group.
// Don't use this type directly, use NewBmcKeySetsClient() instead.
type BmcKeySetsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewBmcKeySetsClient creates a new instance of BmcKeySetsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewBmcKeySetsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BmcKeySetsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &BmcKeySetsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a new baseboard management controller key set or update the existing one for the provided
// cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bmcKeySetName - The name of the baseboard management controller key set.
//   - bmcKeySetParameters - The request body.
//   - options - BmcKeySetsClientBeginCreateOrUpdateOptions contains the optional parameters for the BmcKeySetsClient.BeginCreateOrUpdate
//     method.
func (client *BmcKeySetsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, bmcKeySetParameters BmcKeySet, options *BmcKeySetsClientBeginCreateOrUpdateOptions) (*runtime.Poller[BmcKeySetsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, bmcKeySetName, bmcKeySetParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[BmcKeySetsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[BmcKeySetsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a new baseboard management controller key set or update the existing one for the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
func (client *BmcKeySetsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, bmcKeySetParameters BmcKeySet, options *BmcKeySetsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "BmcKeySetsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, bmcKeySetName, bmcKeySetParameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *BmcKeySetsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, bmcKeySetParameters BmcKeySet, options *BmcKeySetsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bmcKeySets/{bmcKeySetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if bmcKeySetName == "" {
		return nil, errors.New("parameter bmcKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bmcKeySetName}", url.PathEscape(bmcKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, bmcKeySetParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete the baseboard management controller key set of the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bmcKeySetName - The name of the baseboard management controller key set.
//   - options - BmcKeySetsClientBeginDeleteOptions contains the optional parameters for the BmcKeySetsClient.BeginDelete method.
func (client *BmcKeySetsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, options *BmcKeySetsClientBeginDeleteOptions) (*runtime.Poller[BmcKeySetsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, bmcKeySetName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[BmcKeySetsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[BmcKeySetsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete the baseboard management controller key set of the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
func (client *BmcKeySetsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, options *BmcKeySetsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "BmcKeySetsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, bmcKeySetName, options)
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
func (client *BmcKeySetsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, options *BmcKeySetsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bmcKeySets/{bmcKeySetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if bmcKeySetName == "" {
		return nil, errors.New("parameter bmcKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bmcKeySetName}", url.PathEscape(bmcKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get baseboard management controller key set of the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bmcKeySetName - The name of the baseboard management controller key set.
//   - options - BmcKeySetsClientGetOptions contains the optional parameters for the BmcKeySetsClient.Get method.
func (client *BmcKeySetsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, options *BmcKeySetsClientGetOptions) (BmcKeySetsClientGetResponse, error) {
	var err error
	const operationName = "BmcKeySetsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, bmcKeySetName, options)
	if err != nil {
		return BmcKeySetsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return BmcKeySetsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return BmcKeySetsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *BmcKeySetsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, options *BmcKeySetsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bmcKeySets/{bmcKeySetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if bmcKeySetName == "" {
		return nil, errors.New("parameter bmcKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bmcKeySetName}", url.PathEscape(bmcKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *BmcKeySetsClient) getHandleResponse(resp *http.Response) (BmcKeySetsClientGetResponse, error) {
	result := BmcKeySetsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BmcKeySet); err != nil {
		return BmcKeySetsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByClusterPager - Get a list of baseboard management controller key sets for the provided cluster.
//
// Generated from API version 2024-06-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - options - BmcKeySetsClientListByClusterOptions contains the optional parameters for the BmcKeySetsClient.NewListByClusterPager
//     method.
func (client *BmcKeySetsClient) NewListByClusterPager(resourceGroupName string, clusterName string, options *BmcKeySetsClientListByClusterOptions) *runtime.Pager[BmcKeySetsClientListByClusterResponse] {
	return runtime.NewPager(runtime.PagingHandler[BmcKeySetsClientListByClusterResponse]{
		More: func(page BmcKeySetsClientListByClusterResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *BmcKeySetsClientListByClusterResponse) (BmcKeySetsClientListByClusterResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "BmcKeySetsClient.NewListByClusterPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByClusterCreateRequest(ctx, resourceGroupName, clusterName, options)
			}, nil)
			if err != nil {
				return BmcKeySetsClientListByClusterResponse{}, err
			}
			return client.listByClusterHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByClusterCreateRequest creates the ListByCluster request.
func (client *BmcKeySetsClient) listByClusterCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, options *BmcKeySetsClientListByClusterOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bmcKeySets"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByClusterHandleResponse handles the ListByCluster response.
func (client *BmcKeySetsClient) listByClusterHandleResponse(resp *http.Response) (BmcKeySetsClientListByClusterResponse, error) {
	result := BmcKeySetsClientListByClusterResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BmcKeySetList); err != nil {
		return BmcKeySetsClientListByClusterResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Patch properties of baseboard management controller key set for the provided cluster, or update the tags
// associated with it. Properties and tag updates can be done independently.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bmcKeySetName - The name of the baseboard management controller key set.
//   - bmcKeySetUpdateParameters - The request body.
//   - options - BmcKeySetsClientBeginUpdateOptions contains the optional parameters for the BmcKeySetsClient.BeginUpdate method.
func (client *BmcKeySetsClient) BeginUpdate(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, bmcKeySetUpdateParameters BmcKeySetPatchParameters, options *BmcKeySetsClientBeginUpdateOptions) (*runtime.Poller[BmcKeySetsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, clusterName, bmcKeySetName, bmcKeySetUpdateParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[BmcKeySetsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[BmcKeySetsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Patch properties of baseboard management controller key set for the provided cluster, or update the tags associated
// with it. Properties and tag updates can be done independently.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-06-01-preview
func (client *BmcKeySetsClient) update(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, bmcKeySetUpdateParameters BmcKeySetPatchParameters, options *BmcKeySetsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "BmcKeySetsClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, clusterName, bmcKeySetName, bmcKeySetUpdateParameters, options)
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

// updateCreateRequest creates the Update request.
func (client *BmcKeySetsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bmcKeySetName string, bmcKeySetUpdateParameters BmcKeySetPatchParameters, options *BmcKeySetsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bmcKeySets/{bmcKeySetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if clusterName == "" {
		return nil, errors.New("parameter clusterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{clusterName}", url.PathEscape(clusterName))
	if bmcKeySetName == "" {
		return nil, errors.New("parameter bmcKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bmcKeySetName}", url.PathEscape(bmcKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, bmcKeySetUpdateParameters); err != nil {
		return nil, err
	}
	return req, nil
}

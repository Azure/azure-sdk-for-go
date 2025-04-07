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

// BareMetalMachineKeySetsClient contains the methods for the BareMetalMachineKeySets group.
// Don't use this type directly, use NewBareMetalMachineKeySetsClient() instead.
type BareMetalMachineKeySetsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewBareMetalMachineKeySetsClient creates a new instance of BareMetalMachineKeySetsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewBareMetalMachineKeySetsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*BareMetalMachineKeySetsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &BareMetalMachineKeySetsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a new bare metal machine key set or update the existing one for the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bareMetalMachineKeySetName - The name of the bare metal machine key set.
//   - bareMetalMachineKeySetParameters - The request body.
//   - options - BareMetalMachineKeySetsClientBeginCreateOrUpdateOptions contains the optional parameters for the BareMetalMachineKeySetsClient.BeginCreateOrUpdate
//     method.
func (client *BareMetalMachineKeySetsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, bareMetalMachineKeySetParameters BareMetalMachineKeySet, options *BareMetalMachineKeySetsClientBeginCreateOrUpdateOptions) (*runtime.Poller[BareMetalMachineKeySetsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, bareMetalMachineKeySetParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[BareMetalMachineKeySetsClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[BareMetalMachineKeySetsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Create a new bare metal machine key set or update the existing one for the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
func (client *BareMetalMachineKeySetsClient) createOrUpdate(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, bareMetalMachineKeySetParameters BareMetalMachineKeySet, options *BareMetalMachineKeySetsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "BareMetalMachineKeySetsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, bareMetalMachineKeySetParameters, options)
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
func (client *BareMetalMachineKeySetsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, bareMetalMachineKeySetParameters BareMetalMachineKeySet, _ *BareMetalMachineKeySetsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bareMetalMachineKeySets/{bareMetalMachineKeySetName}"
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
	if bareMetalMachineKeySetName == "" {
		return nil, errors.New("parameter bareMetalMachineKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bareMetalMachineKeySetName}", url.PathEscape(bareMetalMachineKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, bareMetalMachineKeySetParameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete the bare metal machine key set of the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bareMetalMachineKeySetName - The name of the bare metal machine key set.
//   - options - BareMetalMachineKeySetsClientBeginDeleteOptions contains the optional parameters for the BareMetalMachineKeySetsClient.BeginDelete
//     method.
func (client *BareMetalMachineKeySetsClient) BeginDelete(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, options *BareMetalMachineKeySetsClientBeginDeleteOptions) (*runtime.Poller[BareMetalMachineKeySetsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[BareMetalMachineKeySetsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[BareMetalMachineKeySetsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete the bare metal machine key set of the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
func (client *BareMetalMachineKeySetsClient) deleteOperation(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, options *BareMetalMachineKeySetsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "BareMetalMachineKeySetsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, options)
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
func (client *BareMetalMachineKeySetsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, _ *BareMetalMachineKeySetsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bareMetalMachineKeySets/{bareMetalMachineKeySetName}"
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
	if bareMetalMachineKeySetName == "" {
		return nil, errors.New("parameter bareMetalMachineKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bareMetalMachineKeySetName}", url.PathEscape(bareMetalMachineKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get bare metal machine key set of the provided cluster.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bareMetalMachineKeySetName - The name of the bare metal machine key set.
//   - options - BareMetalMachineKeySetsClientGetOptions contains the optional parameters for the BareMetalMachineKeySetsClient.Get
//     method.
func (client *BareMetalMachineKeySetsClient) Get(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, options *BareMetalMachineKeySetsClientGetOptions) (BareMetalMachineKeySetsClientGetResponse, error) {
	var err error
	const operationName = "BareMetalMachineKeySetsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, options)
	if err != nil {
		return BareMetalMachineKeySetsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return BareMetalMachineKeySetsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return BareMetalMachineKeySetsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *BareMetalMachineKeySetsClient) getCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, _ *BareMetalMachineKeySetsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bareMetalMachineKeySets/{bareMetalMachineKeySetName}"
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
	if bareMetalMachineKeySetName == "" {
		return nil, errors.New("parameter bareMetalMachineKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bareMetalMachineKeySetName}", url.PathEscape(bareMetalMachineKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *BareMetalMachineKeySetsClient) getHandleResponse(resp *http.Response) (BareMetalMachineKeySetsClientGetResponse, error) {
	result := BareMetalMachineKeySetsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BareMetalMachineKeySet); err != nil {
		return BareMetalMachineKeySetsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByClusterPager - Get a list of bare metal machine key sets for the provided cluster.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - options - BareMetalMachineKeySetsClientListByClusterOptions contains the optional parameters for the BareMetalMachineKeySetsClient.NewListByClusterPager
//     method.
func (client *BareMetalMachineKeySetsClient) NewListByClusterPager(resourceGroupName string, clusterName string, options *BareMetalMachineKeySetsClientListByClusterOptions) *runtime.Pager[BareMetalMachineKeySetsClientListByClusterResponse] {
	return runtime.NewPager(runtime.PagingHandler[BareMetalMachineKeySetsClientListByClusterResponse]{
		More: func(page BareMetalMachineKeySetsClientListByClusterResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *BareMetalMachineKeySetsClientListByClusterResponse) (BareMetalMachineKeySetsClientListByClusterResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "BareMetalMachineKeySetsClient.NewListByClusterPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByClusterCreateRequest(ctx, resourceGroupName, clusterName, options)
			}, nil)
			if err != nil {
				return BareMetalMachineKeySetsClientListByClusterResponse{}, err
			}
			return client.listByClusterHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByClusterCreateRequest creates the ListByCluster request.
func (client *BareMetalMachineKeySetsClient) listByClusterCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, _ *BareMetalMachineKeySetsClientListByClusterOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bareMetalMachineKeySets"
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
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByClusterHandleResponse handles the ListByCluster response.
func (client *BareMetalMachineKeySetsClient) listByClusterHandleResponse(resp *http.Response) (BareMetalMachineKeySetsClientListByClusterResponse, error) {
	result := BareMetalMachineKeySetsClientListByClusterResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.BareMetalMachineKeySetList); err != nil {
		return BareMetalMachineKeySetsClientListByClusterResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Patch properties of bare metal machine key set for the provided cluster, or update the tags associated with
// it. Properties and tag updates can be done independently.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - clusterName - The name of the cluster.
//   - bareMetalMachineKeySetName - The name of the bare metal machine key set.
//   - bareMetalMachineKeySetUpdateParameters - The request body.
//   - options - BareMetalMachineKeySetsClientBeginUpdateOptions contains the optional parameters for the BareMetalMachineKeySetsClient.BeginUpdate
//     method.
func (client *BareMetalMachineKeySetsClient) BeginUpdate(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, bareMetalMachineKeySetUpdateParameters BareMetalMachineKeySetPatchParameters, options *BareMetalMachineKeySetsClientBeginUpdateOptions) (*runtime.Poller[BareMetalMachineKeySetsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, bareMetalMachineKeySetUpdateParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[BareMetalMachineKeySetsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[BareMetalMachineKeySetsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Patch properties of bare metal machine key set for the provided cluster, or update the tags associated with it.
// Properties and tag updates can be done independently.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-10-01-preview
func (client *BareMetalMachineKeySetsClient) update(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, bareMetalMachineKeySetUpdateParameters BareMetalMachineKeySetPatchParameters, options *BareMetalMachineKeySetsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "BareMetalMachineKeySetsClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, clusterName, bareMetalMachineKeySetName, bareMetalMachineKeySetUpdateParameters, options)
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
func (client *BareMetalMachineKeySetsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, clusterName string, bareMetalMachineKeySetName string, bareMetalMachineKeySetUpdateParameters BareMetalMachineKeySetPatchParameters, _ *BareMetalMachineKeySetsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetworkCloud/clusters/{clusterName}/bareMetalMachineKeySets/{bareMetalMachineKeySetName}"
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
	if bareMetalMachineKeySetName == "" {
		return nil, errors.New("parameter bareMetalMachineKeySetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bareMetalMachineKeySetName}", url.PathEscape(bareMetalMachineKeySetName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, bareMetalMachineKeySetUpdateParameters); err != nil {
		return nil, err
	}
	return req, nil
}

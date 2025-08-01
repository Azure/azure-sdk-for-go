// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package armrecoveryservicesdatareplication

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

// ReplicationExtensionClient contains the methods for the ReplicationExtension group.
// Don't use this type directly, use NewReplicationExtensionClient() instead.
type ReplicationExtensionClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewReplicationExtensionClient creates a new instance of ReplicationExtensionClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewReplicationExtensionClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ReplicationExtensionClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ReplicationExtensionClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreate - Creates the replication extension in the given vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - replicationExtensionName - The replication extension name.
//   - resource - Replication extension model.
//   - options - ReplicationExtensionClientBeginCreateOptions contains the optional parameters for the ReplicationExtensionClient.BeginCreate
//     method.
func (client *ReplicationExtensionClient) BeginCreate(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, resource ReplicationExtensionModel, options *ReplicationExtensionClientBeginCreateOptions) (*runtime.Poller[ReplicationExtensionClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, vaultName, replicationExtensionName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationExtensionClientCreateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationExtensionClientCreateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Create - Creates the replication extension in the given vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
func (client *ReplicationExtensionClient) create(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, resource ReplicationExtensionModel, options *ReplicationExtensionClientBeginCreateOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationExtensionClient.BeginCreate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createCreateRequest(ctx, resourceGroupName, vaultName, replicationExtensionName, resource, options)
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

// createCreateRequest creates the Create request.
func (client *ReplicationExtensionClient) createCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, resource ReplicationExtensionModel, _ *ReplicationExtensionClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/replicationExtensions/{replicationExtensionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if replicationExtensionName == "" {
		return nil, errors.New("parameter replicationExtensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{replicationExtensionName}", url.PathEscape(replicationExtensionName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["Content-Type"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the replication extension in the given vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - replicationExtensionName - The replication extension name.
//   - options - ReplicationExtensionClientBeginDeleteOptions contains the optional parameters for the ReplicationExtensionClient.BeginDelete
//     method.
func (client *ReplicationExtensionClient) BeginDelete(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, options *ReplicationExtensionClientBeginDeleteOptions) (*runtime.Poller[ReplicationExtensionClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, vaultName, replicationExtensionName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationExtensionClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationExtensionClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes the replication extension in the given vault.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
func (client *ReplicationExtensionClient) deleteOperation(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, options *ReplicationExtensionClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationExtensionClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, vaultName, replicationExtensionName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ReplicationExtensionClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, _ *ReplicationExtensionClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/replicationExtensions/{replicationExtensionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if replicationExtensionName == "" {
		return nil, errors.New("parameter replicationExtensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{replicationExtensionName}", url.PathEscape(replicationExtensionName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the details of the replication extension.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - replicationExtensionName - The replication extension name.
//   - options - ReplicationExtensionClientGetOptions contains the optional parameters for the ReplicationExtensionClient.Get
//     method.
func (client *ReplicationExtensionClient) Get(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, options *ReplicationExtensionClientGetOptions) (ReplicationExtensionClientGetResponse, error) {
	var err error
	const operationName = "ReplicationExtensionClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, vaultName, replicationExtensionName, options)
	if err != nil {
		return ReplicationExtensionClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ReplicationExtensionClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ReplicationExtensionClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ReplicationExtensionClient) getCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, replicationExtensionName string, _ *ReplicationExtensionClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/replicationExtensions/{replicationExtensionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if replicationExtensionName == "" {
		return nil, errors.New("parameter replicationExtensionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{replicationExtensionName}", url.PathEscape(replicationExtensionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReplicationExtensionClient) getHandleResponse(resp *http.Response) (ReplicationExtensionClientGetResponse, error) {
	result := ReplicationExtensionClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReplicationExtensionModel); err != nil {
		return ReplicationExtensionClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets the list of replication extensions in the given vault.
//
// Generated from API version 2024-09-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - vaultName - The vault name.
//   - options - ReplicationExtensionClientListOptions contains the optional parameters for the ReplicationExtensionClient.NewListPager
//     method.
func (client *ReplicationExtensionClient) NewListPager(resourceGroupName string, vaultName string, options *ReplicationExtensionClientListOptions) *runtime.Pager[ReplicationExtensionClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicationExtensionClientListResponse]{
		More: func(page ReplicationExtensionClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicationExtensionClientListResponse) (ReplicationExtensionClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReplicationExtensionClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, vaultName, options)
			}, nil)
			if err != nil {
				return ReplicationExtensionClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *ReplicationExtensionClient) listCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, _ *ReplicationExtensionClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataReplication/replicationVaults/{vaultName}/replicationExtensions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ReplicationExtensionClient) listHandleResponse(resp *http.Response) (ReplicationExtensionClientListResponse, error) {
	result := ReplicationExtensionClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReplicationExtensionModelListResult); err != nil {
		return ReplicationExtensionClientListResponse{}, err
	}
	return result, nil
}

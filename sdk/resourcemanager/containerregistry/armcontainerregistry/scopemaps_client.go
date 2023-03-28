//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armcontainerregistry

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

// ScopeMapsClient contains the methods for the ScopeMaps group.
// Don't use this type directly, use NewScopeMapsClient() instead.
type ScopeMapsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewScopeMapsClient creates a new instance of ScopeMapsClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewScopeMapsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScopeMapsClient, error) {
	cl, err := arm.NewClient(moduleName+".ScopeMapsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ScopeMapsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreate - Creates a scope map for a container registry with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - registryName - The name of the container registry.
//   - scopeMapName - The name of the scope map.
//   - scopeMapCreateParameters - The parameters for creating a scope map.
//   - options - ScopeMapsClientBeginCreateOptions contains the optional parameters for the ScopeMapsClient.BeginCreate method.
func (client *ScopeMapsClient) BeginCreate(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, scopeMapCreateParameters ScopeMap, options *ScopeMapsClientBeginCreateOptions) (*runtime.Poller[ScopeMapsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, registryName, scopeMapName, scopeMapCreateParameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ScopeMapsClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ScopeMapsClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Creates a scope map for a container registry with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
func (client *ScopeMapsClient) create(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, scopeMapCreateParameters ScopeMap, options *ScopeMapsClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, registryName, scopeMapName, scopeMapCreateParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *ScopeMapsClient) createCreateRequest(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, scopeMapCreateParameters ScopeMap, options *ScopeMapsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/scopeMaps/{scopeMapName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if scopeMapName == "" {
		return nil, errors.New("parameter scopeMapName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeMapName}", url.PathEscape(scopeMapName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, scopeMapCreateParameters)
}

// BeginDelete - Deletes a scope map from a container registry.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - registryName - The name of the container registry.
//   - scopeMapName - The name of the scope map.
//   - options - ScopeMapsClientBeginDeleteOptions contains the optional parameters for the ScopeMapsClient.BeginDelete method.
func (client *ScopeMapsClient) BeginDelete(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, options *ScopeMapsClientBeginDeleteOptions) (*runtime.Poller[ScopeMapsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, registryName, scopeMapName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ScopeMapsClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ScopeMapsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a scope map from a container registry.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
func (client *ScopeMapsClient) deleteOperation(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, options *ScopeMapsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, registryName, scopeMapName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ScopeMapsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, options *ScopeMapsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/scopeMaps/{scopeMapName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if scopeMapName == "" {
		return nil, errors.New("parameter scopeMapName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeMapName}", url.PathEscape(scopeMapName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the properties of the specified scope map.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - registryName - The name of the container registry.
//   - scopeMapName - The name of the scope map.
//   - options - ScopeMapsClientGetOptions contains the optional parameters for the ScopeMapsClient.Get method.
func (client *ScopeMapsClient) Get(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, options *ScopeMapsClientGetOptions) (ScopeMapsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, registryName, scopeMapName, options)
	if err != nil {
		return ScopeMapsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScopeMapsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ScopeMapsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ScopeMapsClient) getCreateRequest(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, options *ScopeMapsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/scopeMaps/{scopeMapName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if scopeMapName == "" {
		return nil, errors.New("parameter scopeMapName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeMapName}", url.PathEscape(scopeMapName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ScopeMapsClient) getHandleResponse(resp *http.Response) (ScopeMapsClientGetResponse, error) {
	result := ScopeMapsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScopeMap); err != nil {
		return ScopeMapsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the scope maps for the specified container registry.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - registryName - The name of the container registry.
//   - options - ScopeMapsClientListOptions contains the optional parameters for the ScopeMapsClient.NewListPager method.
func (client *ScopeMapsClient) NewListPager(resourceGroupName string, registryName string, options *ScopeMapsClientListOptions) *runtime.Pager[ScopeMapsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ScopeMapsClientListResponse]{
		More: func(page ScopeMapsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ScopeMapsClientListResponse) (ScopeMapsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, registryName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ScopeMapsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ScopeMapsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ScopeMapsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ScopeMapsClient) listCreateRequest(ctx context.Context, resourceGroupName string, registryName string, options *ScopeMapsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/scopeMaps"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ScopeMapsClient) listHandleResponse(resp *http.Response) (ScopeMapsClientListResponse, error) {
	result := ScopeMapsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScopeMapListResult); err != nil {
		return ScopeMapsClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Updates a scope map with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - registryName - The name of the container registry.
//   - scopeMapName - The name of the scope map.
//   - scopeMapUpdateParameters - The parameters for updating a scope map.
//   - options - ScopeMapsClientBeginUpdateOptions contains the optional parameters for the ScopeMapsClient.BeginUpdate method.
func (client *ScopeMapsClient) BeginUpdate(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, scopeMapUpdateParameters ScopeMapUpdateParameters, options *ScopeMapsClientBeginUpdateOptions) (*runtime.Poller[ScopeMapsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, registryName, scopeMapName, scopeMapUpdateParameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ScopeMapsClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ScopeMapsClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Updates a scope map with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01-preview
func (client *ScopeMapsClient) update(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, scopeMapUpdateParameters ScopeMapUpdateParameters, options *ScopeMapsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, registryName, scopeMapName, scopeMapUpdateParameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *ScopeMapsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, registryName string, scopeMapName string, scopeMapUpdateParameters ScopeMapUpdateParameters, options *ScopeMapsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/scopeMaps/{scopeMapName}"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if registryName == "" {
		return nil, errors.New("parameter registryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registryName}", url.PathEscape(registryName))
	if scopeMapName == "" {
		return nil, errors.New("parameter scopeMapName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeMapName}", url.PathEscape(scopeMapName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, scopeMapUpdateParameters)
}

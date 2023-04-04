//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armredhatopenshift

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

// SyncIdentityProvidersClient contains the methods for the SyncIdentityProviders group.
// Don't use this type directly, use NewSyncIdentityProvidersClient() instead.
type SyncIdentityProvidersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSyncIdentityProvidersClient creates a new instance of SyncIdentityProvidersClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSyncIdentityProvidersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SyncIdentityProvidersClient, error) {
	cl, err := arm.NewClient(moduleName+".SyncIdentityProvidersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SyncIdentityProvidersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - The operation returns properties of a SyncIdentityProvider.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the OpenShift cluster resource.
//   - childResourceName - The name of the SyncIdentityProvider resource.
//   - parameters - The SyncIdentityProvider resource.
//   - options - SyncIdentityProvidersClientCreateOrUpdateOptions contains the optional parameters for the SyncIdentityProvidersClient.CreateOrUpdate
//     method.
func (client *SyncIdentityProvidersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, parameters SyncIdentityProvider, options *SyncIdentityProvidersClientCreateOrUpdateOptions) (SyncIdentityProvidersClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, resourceName, childResourceName, parameters, options)
	if err != nil {
		return SyncIdentityProvidersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SyncIdentityProvidersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return SyncIdentityProvidersClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SyncIdentityProvidersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, parameters SyncIdentityProvider, options *SyncIdentityProvidersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openshiftclusters/{resourceName}/syncIdentityProvider/{childResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if childResourceName == "" {
		return nil, errors.New("parameter childResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{childResourceName}", url.PathEscape(childResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *SyncIdentityProvidersClient) createOrUpdateHandleResponse(resp *http.Response) (SyncIdentityProvidersClientCreateOrUpdateResponse, error) {
	result := SyncIdentityProvidersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SyncIdentityProvider); err != nil {
		return SyncIdentityProvidersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - The operation returns nothing.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the OpenShift cluster resource.
//   - childResourceName - The name of the SyncIdentityProvider resource.
//   - options - SyncIdentityProvidersClientDeleteOptions contains the optional parameters for the SyncIdentityProvidersClient.Delete
//     method.
func (client *SyncIdentityProvidersClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, options *SyncIdentityProvidersClientDeleteOptions) (SyncIdentityProvidersClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, resourceName, childResourceName, options)
	if err != nil {
		return SyncIdentityProvidersClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SyncIdentityProvidersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return SyncIdentityProvidersClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return SyncIdentityProvidersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SyncIdentityProvidersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, options *SyncIdentityProvidersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openshiftclusters/{resourceName}/syncIdentityProvider/{childResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if childResourceName == "" {
		return nil, errors.New("parameter childResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{childResourceName}", url.PathEscape(childResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - The operation returns properties of a SyncIdentityProvider.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the OpenShift cluster resource.
//   - childResourceName - The name of the SyncIdentityProvider resource.
//   - options - SyncIdentityProvidersClientGetOptions contains the optional parameters for the SyncIdentityProvidersClient.Get
//     method.
func (client *SyncIdentityProvidersClient) Get(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, options *SyncIdentityProvidersClientGetOptions) (SyncIdentityProvidersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, childResourceName, options)
	if err != nil {
		return SyncIdentityProvidersClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SyncIdentityProvidersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SyncIdentityProvidersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SyncIdentityProvidersClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, options *SyncIdentityProvidersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openshiftclusters/{resourceName}/syncIdentityProvider/{childResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if childResourceName == "" {
		return nil, errors.New("parameter childResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{childResourceName}", url.PathEscape(childResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SyncIdentityProvidersClient) getHandleResponse(resp *http.Response) (SyncIdentityProvidersClientGetResponse, error) {
	result := SyncIdentityProvidersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SyncIdentityProvider); err != nil {
		return SyncIdentityProvidersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - The operation returns properties of each SyncIdentityProvider.
//
// Generated from API version 2022-09-04
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the OpenShift cluster resource.
//   - options - SyncIdentityProvidersClientListOptions contains the optional parameters for the SyncIdentityProvidersClient.NewListPager
//     method.
func (client *SyncIdentityProvidersClient) NewListPager(resourceGroupName string, resourceName string, options *SyncIdentityProvidersClientListOptions) *runtime.Pager[SyncIdentityProvidersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SyncIdentityProvidersClientListResponse]{
		More: func(page SyncIdentityProvidersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SyncIdentityProvidersClientListResponse) (SyncIdentityProvidersClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, resourceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SyncIdentityProvidersClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return SyncIdentityProvidersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SyncIdentityProvidersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *SyncIdentityProvidersClient) listCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *SyncIdentityProvidersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openShiftCluster/{resourceName}/syncIdentityProviders"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SyncIdentityProvidersClient) listHandleResponse(resp *http.Response) (SyncIdentityProvidersClientListResponse, error) {
	result := SyncIdentityProvidersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SyncIdentityProviderList); err != nil {
		return SyncIdentityProvidersClientListResponse{}, err
	}
	return result, nil
}

// Update - The operation returns properties of a SyncIdentityProvider.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-04
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the OpenShift cluster resource.
//   - childResourceName - The name of the SyncIdentityProvider resource.
//   - parameters - The SyncIdentityProvider resource.
//   - options - SyncIdentityProvidersClientUpdateOptions contains the optional parameters for the SyncIdentityProvidersClient.Update
//     method.
func (client *SyncIdentityProvidersClient) Update(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, parameters SyncIdentityProviderUpdate, options *SyncIdentityProvidersClientUpdateOptions) (SyncIdentityProvidersClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, childResourceName, parameters, options)
	if err != nil {
		return SyncIdentityProvidersClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SyncIdentityProvidersClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SyncIdentityProvidersClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *SyncIdentityProvidersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, childResourceName string, parameters SyncIdentityProviderUpdate, options *SyncIdentityProvidersClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RedHatOpenShift/openshiftclusters/{resourceName}/syncIdentityProvider/{childResourceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	if childResourceName == "" {
		return nil, errors.New("parameter childResourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{childResourceName}", url.PathEscape(childResourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-04")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *SyncIdentityProvidersClient) updateHandleResponse(resp *http.Response) (SyncIdentityProvidersClientUpdateResponse, error) {
	result := SyncIdentityProvidersClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SyncIdentityProvider); err != nil {
		return SyncIdentityProvidersClientUpdateResponse{}, err
	}
	return result, nil
}

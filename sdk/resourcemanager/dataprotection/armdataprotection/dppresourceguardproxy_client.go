//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdataprotection

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

// DppResourceGuardProxyClient contains the methods for the DppResourceGuardProxy group.
// Don't use this type directly, use NewDppResourceGuardProxyClient() instead.
type DppResourceGuardProxyClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDppResourceGuardProxyClient creates a new instance of DppResourceGuardProxyClient with the specified values.
// subscriptionID - The subscription Id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDppResourceGuardProxyClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DppResourceGuardProxyClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &DppResourceGuardProxyClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Delete -
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-11-01-preview
// resourceGroupName - The name of the resource group where the backup vault is present.
// vaultName - The name of the backup vault.
// options - DppResourceGuardProxyClientDeleteOptions contains the optional parameters for the DppResourceGuardProxyClient.Delete
// method.
func (client *DppResourceGuardProxyClient) Delete(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, options *DppResourceGuardProxyClientDeleteOptions) (DppResourceGuardProxyClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, vaultName, resourceGuardProxyName, options)
	if err != nil {
		return DppResourceGuardProxyClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DppResourceGuardProxyClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return DppResourceGuardProxyClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return DppResourceGuardProxyClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DppResourceGuardProxyClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, options *DppResourceGuardProxyClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupResourceGuardProxies/{resourceGuardProxyName}"
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
	if resourceGuardProxyName == "" {
		return nil, errors.New("parameter resourceGuardProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGuardProxyName}", url.PathEscape(resourceGuardProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get -
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-11-01-preview
// resourceGroupName - The name of the resource group where the backup vault is present.
// vaultName - The name of the backup vault.
// options - DppResourceGuardProxyClientGetOptions contains the optional parameters for the DppResourceGuardProxyClient.Get
// method.
func (client *DppResourceGuardProxyClient) Get(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, options *DppResourceGuardProxyClientGetOptions) (DppResourceGuardProxyClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, vaultName, resourceGuardProxyName, options)
	if err != nil {
		return DppResourceGuardProxyClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DppResourceGuardProxyClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DppResourceGuardProxyClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DppResourceGuardProxyClient) getCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, options *DppResourceGuardProxyClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupResourceGuardProxies/{resourceGuardProxyName}"
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
	if resourceGuardProxyName == "" {
		return nil, errors.New("parameter resourceGuardProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGuardProxyName}", url.PathEscape(resourceGuardProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DppResourceGuardProxyClient) getHandleResponse(resp *http.Response) (DppResourceGuardProxyClientGetResponse, error) {
	result := DppResourceGuardProxyClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ResourceGuardProxyBaseResource); err != nil {
		return DppResourceGuardProxyClientGetResponse{}, err
	}
	return result, nil
}

// resourceGroupName - The name of the resource group where the backup vault is present.
// vaultName - The name of the backup vault.
// options - DppResourceGuardProxyClientListOptions contains the optional parameters for the DppResourceGuardProxyClient.List
// method.
func (client *DppResourceGuardProxyClient) NewListPager(resourceGroupName string, vaultName string, options *DppResourceGuardProxyClientListOptions) *runtime.Pager[DppResourceGuardProxyClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[DppResourceGuardProxyClientListResponse]{
		More: func(page DppResourceGuardProxyClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DppResourceGuardProxyClientListResponse) (DppResourceGuardProxyClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, vaultName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return DppResourceGuardProxyClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return DppResourceGuardProxyClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DppResourceGuardProxyClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *DppResourceGuardProxyClient) listCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, options *DppResourceGuardProxyClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupResourceGuardProxies"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *DppResourceGuardProxyClient) listHandleResponse(resp *http.Response) (DppResourceGuardProxyClientListResponse, error) {
	result := DppResourceGuardProxyClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ResourceGuardProxyBaseResourceList); err != nil {
		return DppResourceGuardProxyClientListResponse{}, err
	}
	return result, nil
}

// Put -
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-11-01-preview
// resourceGroupName - The name of the resource group where the backup vault is present.
// vaultName - The name of the backup vault.
// parameters - Request body for operation
// options - DppResourceGuardProxyClientPutOptions contains the optional parameters for the DppResourceGuardProxyClient.Put
// method.
func (client *DppResourceGuardProxyClient) Put(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, parameters ResourceGuardProxyBaseResource, options *DppResourceGuardProxyClientPutOptions) (DppResourceGuardProxyClientPutResponse, error) {
	req, err := client.putCreateRequest(ctx, resourceGroupName, vaultName, resourceGuardProxyName, parameters, options)
	if err != nil {
		return DppResourceGuardProxyClientPutResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DppResourceGuardProxyClientPutResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DppResourceGuardProxyClientPutResponse{}, runtime.NewResponseError(resp)
	}
	return client.putHandleResponse(resp)
}

// putCreateRequest creates the Put request.
func (client *DppResourceGuardProxyClient) putCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, parameters ResourceGuardProxyBaseResource, options *DppResourceGuardProxyClientPutOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupResourceGuardProxies/{resourceGuardProxyName}"
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
	if resourceGuardProxyName == "" {
		return nil, errors.New("parameter resourceGuardProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGuardProxyName}", url.PathEscape(resourceGuardProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// putHandleResponse handles the Put response.
func (client *DppResourceGuardProxyClient) putHandleResponse(resp *http.Response) (DppResourceGuardProxyClientPutResponse, error) {
	result := DppResourceGuardProxyClientPutResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ResourceGuardProxyBaseResource); err != nil {
		return DppResourceGuardProxyClientPutResponse{}, err
	}
	return result, nil
}

// UnlockDelete -
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-11-01-preview
// resourceGroupName - The name of the resource group where the backup vault is present.
// vaultName - The name of the backup vault.
// parameters - Request body for operation
// options - DppResourceGuardProxyClientUnlockDeleteOptions contains the optional parameters for the DppResourceGuardProxyClient.UnlockDelete
// method.
func (client *DppResourceGuardProxyClient) UnlockDelete(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, parameters UnlockDeleteRequest, options *DppResourceGuardProxyClientUnlockDeleteOptions) (DppResourceGuardProxyClientUnlockDeleteResponse, error) {
	req, err := client.unlockDeleteCreateRequest(ctx, resourceGroupName, vaultName, resourceGuardProxyName, parameters, options)
	if err != nil {
		return DppResourceGuardProxyClientUnlockDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DppResourceGuardProxyClientUnlockDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DppResourceGuardProxyClientUnlockDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return client.unlockDeleteHandleResponse(resp)
}

// unlockDeleteCreateRequest creates the UnlockDelete request.
func (client *DppResourceGuardProxyClient) unlockDeleteCreateRequest(ctx context.Context, resourceGroupName string, vaultName string, resourceGuardProxyName string, parameters UnlockDeleteRequest, options *DppResourceGuardProxyClientUnlockDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupResourceGuardProxies/{resourceGuardProxyName}/unlockDelete"
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
	if resourceGuardProxyName == "" {
		return nil, errors.New("parameter resourceGuardProxyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGuardProxyName}", url.PathEscape(resourceGuardProxyName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// unlockDeleteHandleResponse handles the UnlockDelete response.
func (client *DppResourceGuardProxyClient) unlockDeleteHandleResponse(resp *http.Response) (DppResourceGuardProxyClientUnlockDeleteResponse, error) {
	result := DppResourceGuardProxyClientUnlockDeleteResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UnlockDeleteResponse); err != nil {
		return DppResourceGuardProxyClientUnlockDeleteResponse{}, err
	}
	return result, nil
}

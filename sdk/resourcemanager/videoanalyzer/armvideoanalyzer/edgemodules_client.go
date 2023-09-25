//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armvideoanalyzer

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// EdgeModulesClient contains the methods for the EdgeModules group.
// Don't use this type directly, use NewEdgeModulesClient() instead.
type EdgeModulesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewEdgeModulesClient creates a new instance of EdgeModulesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewEdgeModulesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EdgeModulesClient, error) {
	cl, err := arm.NewClient(moduleName+".EdgeModulesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &EdgeModulesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates a new edge module or updates an existing one. An edge module resource enables a single instance
// of an Azure Video Analyzer IoT edge module to interact with the Video Analyzer Account. This is
// used for authorization and also to make sure that the particular edge module instance only has access to the data it requires
// from the Azure Video Analyzer service. A new edge module resource should
// be created for every new instance of an Azure Video Analyzer edge module deployed to you Azure IoT edge environment. Edge
// module resources can be deleted if the specific module is not in use anymore.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accountName - The Azure Video Analyzer account name.
//   - edgeModuleName - The Edge Module name.
//   - parameters - The request parameters
//   - options - EdgeModulesClientCreateOrUpdateOptions contains the optional parameters for the EdgeModulesClient.CreateOrUpdate
//     method.
func (client *EdgeModulesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, parameters EdgeModuleEntity, options *EdgeModulesClientCreateOrUpdateOptions) (EdgeModulesClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, accountName, edgeModuleName, parameters, options)
	if err != nil {
		return EdgeModulesClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EdgeModulesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return EdgeModulesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *EdgeModulesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, parameters EdgeModuleEntity, options *EdgeModulesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if edgeModuleName == "" {
		return nil, errors.New("parameter edgeModuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeModuleName}", url.PathEscape(edgeModuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *EdgeModulesClient) createOrUpdateHandleResponse(resp *http.Response) (EdgeModulesClientCreateOrUpdateResponse, error) {
	result := EdgeModulesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EdgeModuleEntity); err != nil {
		return EdgeModulesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an existing edge module resource. Deleting the edge module resource will prevent an Azure Video Analyzer
// IoT edge module which was previously initiated with the module provisioning token from
// communicating with the cloud.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accountName - The Azure Video Analyzer account name.
//   - edgeModuleName - The Edge Module name.
//   - options - EdgeModulesClientDeleteOptions contains the optional parameters for the EdgeModulesClient.Delete method.
func (client *EdgeModulesClient) Delete(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, options *EdgeModulesClientDeleteOptions) (EdgeModulesClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, edgeModuleName, options)
	if err != nil {
		return EdgeModulesClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EdgeModulesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return EdgeModulesClientDeleteResponse{}, err
	}
	return EdgeModulesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *EdgeModulesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, options *EdgeModulesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if edgeModuleName == "" {
		return nil, errors.New("parameter edgeModuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeModuleName}", url.PathEscape(edgeModuleName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieves an existing edge module resource with the given name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accountName - The Azure Video Analyzer account name.
//   - edgeModuleName - The Edge Module name.
//   - options - EdgeModulesClientGetOptions contains the optional parameters for the EdgeModulesClient.Get method.
func (client *EdgeModulesClient) Get(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, options *EdgeModulesClientGetOptions) (EdgeModulesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, edgeModuleName, options)
	if err != nil {
		return EdgeModulesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EdgeModulesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EdgeModulesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *EdgeModulesClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, options *EdgeModulesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if edgeModuleName == "" {
		return nil, errors.New("parameter edgeModuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeModuleName}", url.PathEscape(edgeModuleName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EdgeModulesClient) getHandleResponse(resp *http.Response) (EdgeModulesClientGetResponse, error) {
	result := EdgeModulesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EdgeModuleEntity); err != nil {
		return EdgeModulesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List all existing edge module resources, along with their JSON representations.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accountName - The Azure Video Analyzer account name.
//   - options - EdgeModulesClientListOptions contains the optional parameters for the EdgeModulesClient.NewListPager method.
func (client *EdgeModulesClient) NewListPager(resourceGroupName string, accountName string, options *EdgeModulesClientListOptions) (*runtime.Pager[EdgeModulesClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[EdgeModulesClientListResponse]{
		More: func(page EdgeModulesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EdgeModulesClientListResponse) (EdgeModulesClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, accountName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return EdgeModulesClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return EdgeModulesClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return EdgeModulesClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *EdgeModulesClient) listCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *EdgeModulesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *EdgeModulesClient) listHandleResponse(resp *http.Response) (EdgeModulesClientListResponse, error) {
	result := EdgeModulesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EdgeModuleEntityCollection); err != nil {
		return EdgeModulesClientListResponse{}, err
	}
	return result, nil
}

// ListProvisioningToken - Creates a new provisioning token. A provisioning token allows for a single instance of Azure Video
// analyzer IoT edge module to be initialized and authorized to the cloud account. The provisioning
// token itself is short lived and it is only used for the initial handshake between IoT edge module and the cloud. After
// the initial handshake, the IoT edge module will agree on a set of authentication
// keys which will be auto-rotated as long as the module is able to periodically connect to the cloud. A new provisioning
// token can be generated for the same IoT edge module in case the module state lost
// or reset.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accountName - The Azure Video Analyzer account name.
//   - edgeModuleName - The Edge Module name.
//   - parameters - The request parameters
//   - options - EdgeModulesClientListProvisioningTokenOptions contains the optional parameters for the EdgeModulesClient.ListProvisioningToken
//     method.
func (client *EdgeModulesClient) ListProvisioningToken(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, parameters ListProvisioningTokenInput, options *EdgeModulesClientListProvisioningTokenOptions) (EdgeModulesClientListProvisioningTokenResponse, error) {
	var err error
	req, err := client.listProvisioningTokenCreateRequest(ctx, resourceGroupName, accountName, edgeModuleName, parameters, options)
	if err != nil {
		return EdgeModulesClientListProvisioningTokenResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EdgeModulesClientListProvisioningTokenResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EdgeModulesClientListProvisioningTokenResponse{}, err
	}
	resp, err := client.listProvisioningTokenHandleResponse(httpResp)
	return resp, err
}

// listProvisioningTokenCreateRequest creates the ListProvisioningToken request.
func (client *EdgeModulesClient) listProvisioningTokenCreateRequest(ctx context.Context, resourceGroupName string, accountName string, edgeModuleName string, parameters ListProvisioningTokenInput, options *EdgeModulesClientListProvisioningTokenOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}/listProvisioningToken"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if edgeModuleName == "" {
		return nil, errors.New("parameter edgeModuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{edgeModuleName}", url.PathEscape(edgeModuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// listProvisioningTokenHandleResponse handles the ListProvisioningToken response.
func (client *EdgeModulesClient) listProvisioningTokenHandleResponse(resp *http.Response) (EdgeModulesClientListProvisioningTokenResponse, error) {
	result := EdgeModulesClientListProvisioningTokenResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EdgeModuleProvisioningToken); err != nil {
		return EdgeModulesClientListProvisioningTokenResponse{}, err
	}
	return result, nil
}


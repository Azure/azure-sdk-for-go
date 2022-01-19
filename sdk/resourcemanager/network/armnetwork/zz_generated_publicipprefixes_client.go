//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// PublicIPPrefixesClient contains the methods for the PublicIPPrefixes group.
// Don't use this type directly, use NewPublicIPPrefixesClient() instead.
type PublicIPPrefixesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewPublicIPPrefixesClient creates a new instance of PublicIPPrefixesClient with the specified values.
// subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
// ID forms part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewPublicIPPrefixesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *PublicIPPrefixesClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &PublicIPPrefixesClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// BeginCreateOrUpdate - Creates or updates a static or dynamic public IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// publicIPPrefixName - The name of the public IP prefix.
// parameters - Parameters supplied to the create or update public IP prefix operation.
// options - PublicIPPrefixesClientBeginCreateOrUpdateOptions contains the optional parameters for the PublicIPPrefixesClient.BeginCreateOrUpdate
// method.
func (client *PublicIPPrefixesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, publicIPPrefixName string, parameters PublicIPPrefix, options *PublicIPPrefixesClientBeginCreateOrUpdateOptions) (PublicIPPrefixesClientCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, publicIPPrefixName, parameters, options)
	if err != nil {
		return PublicIPPrefixesClientCreateOrUpdatePollerResponse{}, err
	}
	result := PublicIPPrefixesClientCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("PublicIPPrefixesClient.CreateOrUpdate", "location", resp, client.pl)
	if err != nil {
		return PublicIPPrefixesClientCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &PublicIPPrefixesClientCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates or updates a static or dynamic public IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *PublicIPPrefixesClient) createOrUpdate(ctx context.Context, resourceGroupName string, publicIPPrefixName string, parameters PublicIPPrefix, options *PublicIPPrefixesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, publicIPPrefixName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PublicIPPrefixesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, publicIPPrefixName string, parameters PublicIPPrefix, options *PublicIPPrefixesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPPrefixes/{publicIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publicIPPrefixName == "" {
		return nil, errors.New("parameter publicIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publicIpPrefixName}", url.PathEscape(publicIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Deletes the specified public IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// publicIPPrefixName - The name of the PublicIpPrefix.
// options - PublicIPPrefixesClientBeginDeleteOptions contains the optional parameters for the PublicIPPrefixesClient.BeginDelete
// method.
func (client *PublicIPPrefixesClient) BeginDelete(ctx context.Context, resourceGroupName string, publicIPPrefixName string, options *PublicIPPrefixesClientBeginDeleteOptions) (PublicIPPrefixesClientDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, publicIPPrefixName, options)
	if err != nil {
		return PublicIPPrefixesClientDeletePollerResponse{}, err
	}
	result := PublicIPPrefixesClientDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("PublicIPPrefixesClient.Delete", "location", resp, client.pl)
	if err != nil {
		return PublicIPPrefixesClientDeletePollerResponse{}, err
	}
	result.Poller = &PublicIPPrefixesClientDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Deletes the specified public IP prefix.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *PublicIPPrefixesClient) deleteOperation(ctx context.Context, resourceGroupName string, publicIPPrefixName string, options *PublicIPPrefixesClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, publicIPPrefixName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PublicIPPrefixesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, publicIPPrefixName string, options *PublicIPPrefixesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPPrefixes/{publicIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publicIPPrefixName == "" {
		return nil, errors.New("parameter publicIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publicIpPrefixName}", url.PathEscape(publicIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets the specified public IP prefix in a specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// publicIPPrefixName - The name of the public IP prefix.
// options - PublicIPPrefixesClientGetOptions contains the optional parameters for the PublicIPPrefixesClient.Get method.
func (client *PublicIPPrefixesClient) Get(ctx context.Context, resourceGroupName string, publicIPPrefixName string, options *PublicIPPrefixesClientGetOptions) (PublicIPPrefixesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, publicIPPrefixName, options)
	if err != nil {
		return PublicIPPrefixesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PublicIPPrefixesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PublicIPPrefixesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *PublicIPPrefixesClient) getCreateRequest(ctx context.Context, resourceGroupName string, publicIPPrefixName string, options *PublicIPPrefixesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPPrefixes/{publicIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publicIPPrefixName == "" {
		return nil, errors.New("parameter publicIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publicIpPrefixName}", url.PathEscape(publicIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PublicIPPrefixesClient) getHandleResponse(resp *http.Response) (PublicIPPrefixesClientGetResponse, error) {
	result := PublicIPPrefixesClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublicIPPrefix); err != nil {
		return PublicIPPrefixesClientGetResponse{}, err
	}
	return result, nil
}

// List - Gets all public IP prefixes in a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// options - PublicIPPrefixesClientListOptions contains the optional parameters for the PublicIPPrefixesClient.List method.
func (client *PublicIPPrefixesClient) List(resourceGroupName string, options *PublicIPPrefixesClientListOptions) *PublicIPPrefixesClientListPager {
	return &PublicIPPrefixesClientListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, options)
		},
		advancer: func(ctx context.Context, resp PublicIPPrefixesClientListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.PublicIPPrefixListResult.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *PublicIPPrefixesClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *PublicIPPrefixesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPPrefixes"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *PublicIPPrefixesClient) listHandleResponse(resp *http.Response) (PublicIPPrefixesClientListResponse, error) {
	result := PublicIPPrefixesClientListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublicIPPrefixListResult); err != nil {
		return PublicIPPrefixesClientListResponse{}, err
	}
	return result, nil
}

// ListAll - Gets all the public IP prefixes in a subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// options - PublicIPPrefixesClientListAllOptions contains the optional parameters for the PublicIPPrefixesClient.ListAll
// method.
func (client *PublicIPPrefixesClient) ListAll(options *PublicIPPrefixesClientListAllOptions) *PublicIPPrefixesClientListAllPager {
	return &PublicIPPrefixesClientListAllPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listAllCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp PublicIPPrefixesClientListAllResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.PublicIPPrefixListResult.NextLink)
		},
	}
}

// listAllCreateRequest creates the ListAll request.
func (client *PublicIPPrefixesClient) listAllCreateRequest(ctx context.Context, options *PublicIPPrefixesClientListAllOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/publicIPPrefixes"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listAllHandleResponse handles the ListAll response.
func (client *PublicIPPrefixesClient) listAllHandleResponse(resp *http.Response) (PublicIPPrefixesClientListAllResponse, error) {
	result := PublicIPPrefixesClientListAllResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublicIPPrefixListResult); err != nil {
		return PublicIPPrefixesClientListAllResponse{}, err
	}
	return result, nil
}

// UpdateTags - Updates public IP prefix tags.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group.
// publicIPPrefixName - The name of the public IP prefix.
// parameters - Parameters supplied to update public IP prefix tags.
// options - PublicIPPrefixesClientUpdateTagsOptions contains the optional parameters for the PublicIPPrefixesClient.UpdateTags
// method.
func (client *PublicIPPrefixesClient) UpdateTags(ctx context.Context, resourceGroupName string, publicIPPrefixName string, parameters TagsObject, options *PublicIPPrefixesClientUpdateTagsOptions) (PublicIPPrefixesClientUpdateTagsResponse, error) {
	req, err := client.updateTagsCreateRequest(ctx, resourceGroupName, publicIPPrefixName, parameters, options)
	if err != nil {
		return PublicIPPrefixesClientUpdateTagsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PublicIPPrefixesClientUpdateTagsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PublicIPPrefixesClientUpdateTagsResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateTagsHandleResponse(resp)
}

// updateTagsCreateRequest creates the UpdateTags request.
func (client *PublicIPPrefixesClient) updateTagsCreateRequest(ctx context.Context, resourceGroupName string, publicIPPrefixName string, parameters TagsObject, options *PublicIPPrefixesClientUpdateTagsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/publicIPPrefixes/{publicIpPrefixName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if publicIPPrefixName == "" {
		return nil, errors.New("parameter publicIPPrefixName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{publicIpPrefixName}", url.PathEscape(publicIPPrefixName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateTagsHandleResponse handles the UpdateTags response.
func (client *PublicIPPrefixesClient) updateTagsHandleResponse(resp *http.Response) (PublicIPPrefixesClientUpdateTagsResponse, error) {
	result := PublicIPPrefixesClientUpdateTagsResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PublicIPPrefix); err != nil {
		return PublicIPPrefixesClientUpdateTagsResponse{}, err
	}
	return result, nil
}

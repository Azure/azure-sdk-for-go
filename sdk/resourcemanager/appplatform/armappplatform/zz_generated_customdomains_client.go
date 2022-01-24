//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armappplatform

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

// CustomDomainsClient contains the methods for the CustomDomains group.
// Don't use this type directly, use NewCustomDomainsClient() instead.
type CustomDomainsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewCustomDomainsClient creates a new instance of CustomDomainsClient with the specified values.
// subscriptionID - Gets subscription ID which uniquely identify the Microsoft Azure subscription. The subscription ID forms
// part of the URI for every service call.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewCustomDomainsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *CustomDomainsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &CustomDomainsClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// BeginCreateOrUpdate - Create or update custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// appName - The name of the App resource.
// domainName - The name of the custom domain resource.
// domainResource - Parameters for the create or update operation
// options - CustomDomainsClientBeginCreateOrUpdateOptions contains the optional parameters for the CustomDomainsClient.BeginCreateOrUpdate
// method.
func (client *CustomDomainsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, domainResource CustomDomainResource, options *CustomDomainsClientBeginCreateOrUpdateOptions) (CustomDomainsClientCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, serviceName, appName, domainName, domainResource, options)
	if err != nil {
		return CustomDomainsClientCreateOrUpdatePollerResponse{}, err
	}
	result := CustomDomainsClientCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("CustomDomainsClient.CreateOrUpdate", "azure-async-operation", resp, client.pl)
	if err != nil {
		return CustomDomainsClientCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &CustomDomainsClientCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Create or update custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *CustomDomainsClient) createOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, domainResource CustomDomainResource, options *CustomDomainsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, appName, domainName, domainResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *CustomDomainsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, domainResource CustomDomainResource, options *CustomDomainsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/domains/{domainName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if domainName == "" {
		return nil, errors.New("parameter domainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{domainName}", url.PathEscape(domainName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, domainResource)
}

// BeginDelete - Delete the custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// appName - The name of the App resource.
// domainName - The name of the custom domain resource.
// options - CustomDomainsClientBeginDeleteOptions contains the optional parameters for the CustomDomainsClient.BeginDelete
// method.
func (client *CustomDomainsClient) BeginDelete(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, options *CustomDomainsClientBeginDeleteOptions) (CustomDomainsClientDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceGroupName, serviceName, appName, domainName, options)
	if err != nil {
		return CustomDomainsClientDeletePollerResponse{}, err
	}
	result := CustomDomainsClientDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("CustomDomainsClient.Delete", "azure-async-operation", resp, client.pl)
	if err != nil {
		return CustomDomainsClientDeletePollerResponse{}, err
	}
	result.Poller = &CustomDomainsClientDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Delete the custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *CustomDomainsClient) deleteOperation(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, options *CustomDomainsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, appName, domainName, options)
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
func (client *CustomDomainsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, options *CustomDomainsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/domains/{domainName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if domainName == "" {
		return nil, errors.New("parameter domainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{domainName}", url.PathEscape(domainName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Get the custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// appName - The name of the App resource.
// domainName - The name of the custom domain resource.
// options - CustomDomainsClientGetOptions contains the optional parameters for the CustomDomainsClient.Get method.
func (client *CustomDomainsClient) Get(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, options *CustomDomainsClientGetOptions) (CustomDomainsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, appName, domainName, options)
	if err != nil {
		return CustomDomainsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return CustomDomainsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CustomDomainsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *CustomDomainsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, options *CustomDomainsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/domains/{domainName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if domainName == "" {
		return nil, errors.New("parameter domainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{domainName}", url.PathEscape(domainName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CustomDomainsClient) getHandleResponse(resp *http.Response) (CustomDomainsClientGetResponse, error) {
	result := CustomDomainsClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomDomainResource); err != nil {
		return CustomDomainsClientGetResponse{}, err
	}
	return result, nil
}

// List - List the custom domains of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// appName - The name of the App resource.
// options - CustomDomainsClientListOptions contains the optional parameters for the CustomDomainsClient.List method.
func (client *CustomDomainsClient) List(resourceGroupName string, serviceName string, appName string, options *CustomDomainsClientListOptions) *CustomDomainsClientListPager {
	return &CustomDomainsClientListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, serviceName, appName, options)
		},
		advancer: func(ctx context.Context, resp CustomDomainsClientListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.CustomDomainResourceCollection.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *CustomDomainsClient) listCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, options *CustomDomainsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/domains"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *CustomDomainsClient) listHandleResponse(resp *http.Response) (CustomDomainsClientListResponse, error) {
	result := CustomDomainsClientListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomDomainResourceCollection); err != nil {
		return CustomDomainsClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serviceName - The name of the Service resource.
// appName - The name of the App resource.
// domainName - The name of the custom domain resource.
// domainResource - Parameters for the create or update operation
// options - CustomDomainsClientBeginUpdateOptions contains the optional parameters for the CustomDomainsClient.BeginUpdate
// method.
func (client *CustomDomainsClient) BeginUpdate(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, domainResource CustomDomainResource, options *CustomDomainsClientBeginUpdateOptions) (CustomDomainsClientUpdatePollerResponse, error) {
	resp, err := client.update(ctx, resourceGroupName, serviceName, appName, domainName, domainResource, options)
	if err != nil {
		return CustomDomainsClientUpdatePollerResponse{}, err
	}
	result := CustomDomainsClientUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("CustomDomainsClient.Update", "azure-async-operation", resp, client.pl)
	if err != nil {
		return CustomDomainsClientUpdatePollerResponse{}, err
	}
	result.Poller = &CustomDomainsClientUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// Update - Update custom domain of one lifecycle application.
// If the operation fails it returns an *azcore.ResponseError type.
func (client *CustomDomainsClient) update(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, domainResource CustomDomainResource, options *CustomDomainsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serviceName, appName, domainName, domainResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *CustomDomainsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, appName string, domainName string, domainResource CustomDomainResource, options *CustomDomainsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppPlatform/Spring/{serviceName}/apps/{appName}/domains/{domainName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if appName == "" {
		return nil, errors.New("parameter appName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{appName}", url.PathEscape(appName))
	if domainName == "" {
		return nil, errors.New("parameter domainName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{domainName}", url.PathEscape(domainName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, domainResource)
}

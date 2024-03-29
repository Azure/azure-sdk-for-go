//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatalakeanalytics

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

// FirewallRulesClient contains the methods for the FirewallRules group.
// Don't use this type directly, use NewFirewallRulesClient() instead.
type FirewallRulesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewFirewallRulesClient creates a new instance of FirewallRulesClient with the specified values.
//   - subscriptionID - Get subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
//     forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewFirewallRulesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FirewallRulesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &FirewallRulesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates the specified firewall rule. During update, the firewall rule with the specified name
// will be replaced with this new firewall rule.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-11-01-preview
//   - resourceGroupName - The name of the Azure resource group.
//   - accountName - The name of the Data Lake Analytics account.
//   - firewallRuleName - The name of the firewall rule to create or update.
//   - parameters - Parameters supplied to create or update the firewall rule.
//   - options - FirewallRulesClientCreateOrUpdateOptions contains the optional parameters for the FirewallRulesClient.CreateOrUpdate
//     method.
func (client *FirewallRulesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, parameters CreateOrUpdateFirewallRuleParameters, options *FirewallRulesClientCreateOrUpdateOptions) (FirewallRulesClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "FirewallRulesClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, accountName, firewallRuleName, parameters, options)
	if err != nil {
		return FirewallRulesClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FirewallRulesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FirewallRulesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *FirewallRulesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, parameters CreateOrUpdateFirewallRuleParameters, options *FirewallRulesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/firewallRules/{firewallRuleName}"
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
	if firewallRuleName == "" {
		return nil, errors.New("parameter firewallRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{firewallRuleName}", url.PathEscape(firewallRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *FirewallRulesClient) createOrUpdateHandleResponse(resp *http.Response) (FirewallRulesClientCreateOrUpdateResponse, error) {
	result := FirewallRulesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FirewallRule); err != nil {
		return FirewallRulesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the specified firewall rule from the specified Data Lake Analytics account
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-11-01-preview
//   - resourceGroupName - The name of the Azure resource group.
//   - accountName - The name of the Data Lake Analytics account.
//   - firewallRuleName - The name of the firewall rule to delete.
//   - options - FirewallRulesClientDeleteOptions contains the optional parameters for the FirewallRulesClient.Delete method.
func (client *FirewallRulesClient) Delete(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, options *FirewallRulesClientDeleteOptions) (FirewallRulesClientDeleteResponse, error) {
	var err error
	const operationName = "FirewallRulesClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, firewallRuleName, options)
	if err != nil {
		return FirewallRulesClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FirewallRulesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return FirewallRulesClientDeleteResponse{}, err
	}
	return FirewallRulesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *FirewallRulesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, options *FirewallRulesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/firewallRules/{firewallRuleName}"
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
	if firewallRuleName == "" {
		return nil, errors.New("parameter firewallRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{firewallRuleName}", url.PathEscape(firewallRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified Data Lake Analytics firewall rule.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-11-01-preview
//   - resourceGroupName - The name of the Azure resource group.
//   - accountName - The name of the Data Lake Analytics account.
//   - firewallRuleName - The name of the firewall rule to retrieve.
//   - options - FirewallRulesClientGetOptions contains the optional parameters for the FirewallRulesClient.Get method.
func (client *FirewallRulesClient) Get(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, options *FirewallRulesClientGetOptions) (FirewallRulesClientGetResponse, error) {
	var err error
	const operationName = "FirewallRulesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, firewallRuleName, options)
	if err != nil {
		return FirewallRulesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FirewallRulesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FirewallRulesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *FirewallRulesClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, options *FirewallRulesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/firewallRules/{firewallRuleName}"
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
	if firewallRuleName == "" {
		return nil, errors.New("parameter firewallRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{firewallRuleName}", url.PathEscape(firewallRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *FirewallRulesClient) getHandleResponse(resp *http.Response) (FirewallRulesClientGetResponse, error) {
	result := FirewallRulesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FirewallRule); err != nil {
		return FirewallRulesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByAccountPager - Lists the Data Lake Analytics firewall rules within the specified Data Lake Analytics account.
//
// Generated from API version 2019-11-01-preview
//   - resourceGroupName - The name of the Azure resource group.
//   - accountName - The name of the Data Lake Analytics account.
//   - options - FirewallRulesClientListByAccountOptions contains the optional parameters for the FirewallRulesClient.NewListByAccountPager
//     method.
func (client *FirewallRulesClient) NewListByAccountPager(resourceGroupName string, accountName string, options *FirewallRulesClientListByAccountOptions) *runtime.Pager[FirewallRulesClientListByAccountResponse] {
	return runtime.NewPager(runtime.PagingHandler[FirewallRulesClientListByAccountResponse]{
		More: func(page FirewallRulesClientListByAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *FirewallRulesClientListByAccountResponse) (FirewallRulesClientListByAccountResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "FirewallRulesClient.NewListByAccountPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByAccountCreateRequest(ctx, resourceGroupName, accountName, options)
			}, nil)
			if err != nil {
				return FirewallRulesClientListByAccountResponse{}, err
			}
			return client.listByAccountHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByAccountCreateRequest creates the ListByAccount request.
func (client *FirewallRulesClient) listByAccountCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *FirewallRulesClientListByAccountOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/firewallRules"
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
	reqQP.Set("api-version", "2019-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByAccountHandleResponse handles the ListByAccount response.
func (client *FirewallRulesClient) listByAccountHandleResponse(resp *http.Response) (FirewallRulesClientListByAccountResponse, error) {
	result := FirewallRulesClientListByAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FirewallRuleListResult); err != nil {
		return FirewallRulesClientListByAccountResponse{}, err
	}
	return result, nil
}

// Update - Updates the specified firewall rule.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-11-01-preview
//   - resourceGroupName - The name of the Azure resource group.
//   - accountName - The name of the Data Lake Analytics account.
//   - firewallRuleName - The name of the firewall rule to update.
//   - options - FirewallRulesClientUpdateOptions contains the optional parameters for the FirewallRulesClient.Update method.
func (client *FirewallRulesClient) Update(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, options *FirewallRulesClientUpdateOptions) (FirewallRulesClientUpdateResponse, error) {
	var err error
	const operationName = "FirewallRulesClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, accountName, firewallRuleName, options)
	if err != nil {
		return FirewallRulesClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return FirewallRulesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return FirewallRulesClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *FirewallRulesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, options *FirewallRulesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/firewallRules/{firewallRuleName}"
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
	if firewallRuleName == "" {
		return nil, errors.New("parameter firewallRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{firewallRuleName}", url.PathEscape(firewallRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Parameters != nil {
		if err := runtime.MarshalAsJSON(req, *options.Parameters); err != nil {
			return nil, err
		}
		return req, nil
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *FirewallRulesClient) updateHandleResponse(resp *http.Response) (FirewallRulesClientUpdateResponse, error) {
	result := FirewallRulesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.FirewallRule); err != nil {
		return FirewallRulesClientUpdateResponse{}, err
	}
	return result, nil
}

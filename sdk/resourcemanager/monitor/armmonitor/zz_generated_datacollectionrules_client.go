//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

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

// DataCollectionRulesClient contains the methods for the DataCollectionRules group.
// Don't use this type directly, use NewDataCollectionRulesClient() instead.
type DataCollectionRulesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDataCollectionRulesClient creates a new instance of DataCollectionRulesClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDataCollectionRulesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *DataCollectionRulesClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &DataCollectionRulesClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// Create - Creates or updates a data collection rule.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// dataCollectionRuleName - The name of the data collection rule. The name is case insensitive.
// options - DataCollectionRulesClientCreateOptions contains the optional parameters for the DataCollectionRulesClient.Create
// method.
func (client *DataCollectionRulesClient) Create(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientCreateOptions) (DataCollectionRulesClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, dataCollectionRuleName, options)
	if err != nil {
		return DataCollectionRulesClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataCollectionRulesClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return DataCollectionRulesClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *DataCollectionRulesClient) createCreateRequest(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataCollectionRuleName == "" {
		return nil, errors.New("parameter dataCollectionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataCollectionRuleName}", url.PathEscape(dataCollectionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	if options != nil && options.Body != nil {
		return req, runtime.MarshalAsJSON(req, *options.Body)
	}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *DataCollectionRulesClient) createHandleResponse(resp *http.Response) (DataCollectionRulesClientCreateResponse, error) {
	result := DataCollectionRulesClientCreateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataCollectionRuleResource); err != nil {
		return DataCollectionRulesClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a data collection rule.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// dataCollectionRuleName - The name of the data collection rule. The name is case insensitive.
// options - DataCollectionRulesClientDeleteOptions contains the optional parameters for the DataCollectionRulesClient.Delete
// method.
func (client *DataCollectionRulesClient) Delete(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientDeleteOptions) (DataCollectionRulesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, dataCollectionRuleName, options)
	if err != nil {
		return DataCollectionRulesClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataCollectionRulesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return DataCollectionRulesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return DataCollectionRulesClientDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DataCollectionRulesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataCollectionRuleName == "" {
		return nil, errors.New("parameter dataCollectionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataCollectionRuleName}", url.PathEscape(dataCollectionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Returns the specified data collection rule.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// dataCollectionRuleName - The name of the data collection rule. The name is case insensitive.
// options - DataCollectionRulesClientGetOptions contains the optional parameters for the DataCollectionRulesClient.Get method.
func (client *DataCollectionRulesClient) Get(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientGetOptions) (DataCollectionRulesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, dataCollectionRuleName, options)
	if err != nil {
		return DataCollectionRulesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataCollectionRulesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataCollectionRulesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DataCollectionRulesClient) getCreateRequest(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataCollectionRuleName == "" {
		return nil, errors.New("parameter dataCollectionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataCollectionRuleName}", url.PathEscape(dataCollectionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DataCollectionRulesClient) getHandleResponse(resp *http.Response) (DataCollectionRulesClientGetResponse, error) {
	result := DataCollectionRulesClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataCollectionRuleResource); err != nil {
		return DataCollectionRulesClientGetResponse{}, err
	}
	return result, nil
}

// ListByResourceGroup - Lists all data collection rules in the specified resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// options - DataCollectionRulesClientListByResourceGroupOptions contains the optional parameters for the DataCollectionRulesClient.ListByResourceGroup
// method.
func (client *DataCollectionRulesClient) ListByResourceGroup(resourceGroupName string, options *DataCollectionRulesClientListByResourceGroupOptions) *DataCollectionRulesClientListByResourceGroupPager {
	return &DataCollectionRulesClientListByResourceGroupPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
		},
		advancer: func(ctx context.Context, resp DataCollectionRulesClientListByResourceGroupResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.DataCollectionRuleResourceListResult.NextLink)
		},
	}
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *DataCollectionRulesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *DataCollectionRulesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *DataCollectionRulesClient) listByResourceGroupHandleResponse(resp *http.Response) (DataCollectionRulesClientListByResourceGroupResponse, error) {
	result := DataCollectionRulesClientListByResourceGroupResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataCollectionRuleResourceListResult); err != nil {
		return DataCollectionRulesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// ListBySubscription - Lists all data collection rules in the specified subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// options - DataCollectionRulesClientListBySubscriptionOptions contains the optional parameters for the DataCollectionRulesClient.ListBySubscription
// method.
func (client *DataCollectionRulesClient) ListBySubscription(options *DataCollectionRulesClientListBySubscriptionOptions) *DataCollectionRulesClientListBySubscriptionPager {
	return &DataCollectionRulesClientListBySubscriptionPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listBySubscriptionCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp DataCollectionRulesClientListBySubscriptionResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.DataCollectionRuleResourceListResult.NextLink)
		},
	}
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *DataCollectionRulesClient) listBySubscriptionCreateRequest(ctx context.Context, options *DataCollectionRulesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Insights/dataCollectionRules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *DataCollectionRulesClient) listBySubscriptionHandleResponse(resp *http.Response) (DataCollectionRulesClientListBySubscriptionResponse, error) {
	result := DataCollectionRulesClientListBySubscriptionResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataCollectionRuleResourceListResult); err != nil {
		return DataCollectionRulesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Updates part of a data collection rule.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// dataCollectionRuleName - The name of the data collection rule. The name is case insensitive.
// options - DataCollectionRulesClientUpdateOptions contains the optional parameters for the DataCollectionRulesClient.Update
// method.
func (client *DataCollectionRulesClient) Update(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientUpdateOptions) (DataCollectionRulesClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, dataCollectionRuleName, options)
	if err != nil {
		return DataCollectionRulesClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataCollectionRulesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataCollectionRulesClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *DataCollectionRulesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, dataCollectionRuleName string, options *DataCollectionRulesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataCollectionRuleName == "" {
		return nil, errors.New("parameter dataCollectionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataCollectionRuleName}", url.PathEscape(dataCollectionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	if options != nil && options.Body != nil {
		return req, runtime.MarshalAsJSON(req, *options.Body)
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *DataCollectionRulesClient) updateHandleResponse(resp *http.Response) (DataCollectionRulesClientUpdateResponse, error) {
	result := DataCollectionRulesClientUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataCollectionRuleResource); err != nil {
		return DataCollectionRulesClientUpdateResponse{}, err
	}
	return result, nil
}

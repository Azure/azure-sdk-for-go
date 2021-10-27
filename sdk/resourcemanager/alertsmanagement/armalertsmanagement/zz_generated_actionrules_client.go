//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armalertsmanagement

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ActionRulesClient contains the methods for the ActionRules group.
// Don't use this type directly, use NewActionRulesClient() instead.
type ActionRulesClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewActionRulesClient creates a new instance of ActionRulesClient with the specified values.
func NewActionRulesClient(con *arm.Connection, subscriptionID string) *ActionRulesClient {
	return &ActionRulesClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), subscriptionID: subscriptionID}
}

// CreateUpdate - Creates/Updates a specific action rule
// If the operation fails it returns the *ErrorResponse error type.
func (client *ActionRulesClient) CreateUpdate(ctx context.Context, resourceGroupName string, actionRuleName string, actionRule ActionRule, options *ActionRulesCreateUpdateOptions) (ActionRulesCreateUpdateResponse, error) {
	req, err := client.createUpdateCreateRequest(ctx, resourceGroupName, actionRuleName, actionRule, options)
	if err != nil {
		return ActionRulesCreateUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ActionRulesCreateUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ActionRulesCreateUpdateResponse{}, client.createUpdateHandleError(resp)
	}
	return client.createUpdateHandleResponse(resp)
}

// createUpdateCreateRequest creates the CreateUpdate request.
func (client *ActionRulesClient) createUpdateCreateRequest(ctx context.Context, resourceGroupName string, actionRuleName string, actionRule ActionRule, options *ActionRulesCreateUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules/{actionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionRuleName == "" {
		return nil, errors.New("parameter actionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionRuleName}", url.PathEscape(actionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, actionRule)
}

// createUpdateHandleResponse handles the CreateUpdate response.
func (client *ActionRulesClient) createUpdateHandleResponse(resp *http.Response) (ActionRulesCreateUpdateResponse, error) {
	result := ActionRulesCreateUpdateResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionRule); err != nil {
		return ActionRulesCreateUpdateResponse{}, err
	}
	return result, nil
}

// createUpdateHandleError handles the CreateUpdate error response.
func (client *ActionRulesClient) createUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Delete - Deletes a given action rule
// If the operation fails it returns the *ErrorResponse error type.
func (client *ActionRulesClient) Delete(ctx context.Context, resourceGroupName string, actionRuleName string, options *ActionRulesDeleteOptions) (ActionRulesDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, actionRuleName, options)
	if err != nil {
		return ActionRulesDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ActionRulesDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ActionRulesDeleteResponse{}, client.deleteHandleError(resp)
	}
	return client.deleteHandleResponse(resp)
}

// deleteCreateRequest creates the Delete request.
func (client *ActionRulesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, actionRuleName string, options *ActionRulesDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules/{actionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionRuleName == "" {
		return nil, errors.New("parameter actionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionRuleName}", url.PathEscape(actionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *ActionRulesClient) deleteHandleResponse(resp *http.Response) (ActionRulesDeleteResponse, error) {
	result := ActionRulesDeleteResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.Value); err != nil {
		return ActionRulesDeleteResponse{}, err
	}
	return result, nil
}

// deleteHandleError handles the Delete error response.
func (client *ActionRulesClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetByName - Get a specific action rule
// If the operation fails it returns the *ErrorResponse error type.
func (client *ActionRulesClient) GetByName(ctx context.Context, resourceGroupName string, actionRuleName string, options *ActionRulesGetByNameOptions) (ActionRulesGetByNameResponse, error) {
	req, err := client.getByNameCreateRequest(ctx, resourceGroupName, actionRuleName, options)
	if err != nil {
		return ActionRulesGetByNameResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ActionRulesGetByNameResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ActionRulesGetByNameResponse{}, client.getByNameHandleError(resp)
	}
	return client.getByNameHandleResponse(resp)
}

// getByNameCreateRequest creates the GetByName request.
func (client *ActionRulesClient) getByNameCreateRequest(ctx context.Context, resourceGroupName string, actionRuleName string, options *ActionRulesGetByNameOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules/{actionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionRuleName == "" {
		return nil, errors.New("parameter actionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionRuleName}", url.PathEscape(actionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getByNameHandleResponse handles the GetByName response.
func (client *ActionRulesClient) getByNameHandleResponse(resp *http.Response) (ActionRulesGetByNameResponse, error) {
	result := ActionRulesGetByNameResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionRule); err != nil {
		return ActionRulesGetByNameResponse{}, err
	}
	return result, nil
}

// getByNameHandleError handles the GetByName error response.
func (client *ActionRulesClient) getByNameHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListByResourceGroup - List all action rules of the subscription, created in given resource group and given input filters
// If the operation fails it returns the *ErrorResponse error type.
func (client *ActionRulesClient) ListByResourceGroup(resourceGroupName string, options *ActionRulesListByResourceGroupOptions) *ActionRulesListByResourceGroupPager {
	return &ActionRulesListByResourceGroupPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
		},
		advancer: func(ctx context.Context, resp ActionRulesListByResourceGroupResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.ActionRulesList.NextLink)
		},
	}
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ActionRulesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ActionRulesListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.TargetResourceGroup != nil {
		reqQP.Set("targetResourceGroup", *options.TargetResourceGroup)
	}
	if options != nil && options.TargetResourceType != nil {
		reqQP.Set("targetResourceType", *options.TargetResourceType)
	}
	if options != nil && options.TargetResource != nil {
		reqQP.Set("targetResource", *options.TargetResource)
	}
	if options != nil && options.Severity != nil {
		reqQP.Set("severity", string(*options.Severity))
	}
	if options != nil && options.MonitorService != nil {
		reqQP.Set("monitorService", string(*options.MonitorService))
	}
	if options != nil && options.ImpactedScope != nil {
		reqQP.Set("impactedScope", *options.ImpactedScope)
	}
	if options != nil && options.Description != nil {
		reqQP.Set("description", *options.Description)
	}
	if options != nil && options.AlertRuleID != nil {
		reqQP.Set("alertRuleId", *options.AlertRuleID)
	}
	if options != nil && options.ActionGroup != nil {
		reqQP.Set("actionGroup", *options.ActionGroup)
	}
	if options != nil && options.Name != nil {
		reqQP.Set("name", *options.Name)
	}
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ActionRulesClient) listByResourceGroupHandleResponse(resp *http.Response) (ActionRulesListByResourceGroupResponse, error) {
	result := ActionRulesListByResourceGroupResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionRulesList); err != nil {
		return ActionRulesListByResourceGroupResponse{}, err
	}
	return result, nil
}

// listByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client *ActionRulesClient) listByResourceGroupHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListBySubscription - List all action rules of the subscription and given input filters
// If the operation fails it returns the *ErrorResponse error type.
func (client *ActionRulesClient) ListBySubscription(options *ActionRulesListBySubscriptionOptions) *ActionRulesListBySubscriptionPager {
	return &ActionRulesListBySubscriptionPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listBySubscriptionCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp ActionRulesListBySubscriptionResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.ActionRulesList.NextLink)
		},
	}
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *ActionRulesClient) listBySubscriptionCreateRequest(ctx context.Context, options *ActionRulesListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.AlertsManagement/actionRules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.TargetResourceGroup != nil {
		reqQP.Set("targetResourceGroup", *options.TargetResourceGroup)
	}
	if options != nil && options.TargetResourceType != nil {
		reqQP.Set("targetResourceType", *options.TargetResourceType)
	}
	if options != nil && options.TargetResource != nil {
		reqQP.Set("targetResource", *options.TargetResource)
	}
	if options != nil && options.Severity != nil {
		reqQP.Set("severity", string(*options.Severity))
	}
	if options != nil && options.MonitorService != nil {
		reqQP.Set("monitorService", string(*options.MonitorService))
	}
	if options != nil && options.ImpactedScope != nil {
		reqQP.Set("impactedScope", *options.ImpactedScope)
	}
	if options != nil && options.Description != nil {
		reqQP.Set("description", *options.Description)
	}
	if options != nil && options.AlertRuleID != nil {
		reqQP.Set("alertRuleId", *options.AlertRuleID)
	}
	if options != nil && options.ActionGroup != nil {
		reqQP.Set("actionGroup", *options.ActionGroup)
	}
	if options != nil && options.Name != nil {
		reqQP.Set("name", *options.Name)
	}
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *ActionRulesClient) listBySubscriptionHandleResponse(resp *http.Response) (ActionRulesListBySubscriptionResponse, error) {
	result := ActionRulesListBySubscriptionResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionRulesList); err != nil {
		return ActionRulesListBySubscriptionResponse{}, err
	}
	return result, nil
}

// listBySubscriptionHandleError handles the ListBySubscription error response.
func (client *ActionRulesClient) listBySubscriptionHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Update - Update enabled flag and/or tags for the given action rule
// If the operation fails it returns the *ErrorResponse error type.
func (client *ActionRulesClient) Update(ctx context.Context, resourceGroupName string, actionRuleName string, actionRulePatch PatchObject, options *ActionRulesUpdateOptions) (ActionRulesUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, actionRuleName, actionRulePatch, options)
	if err != nil {
		return ActionRulesUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ActionRulesUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ActionRulesUpdateResponse{}, client.updateHandleError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *ActionRulesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, actionRuleName string, actionRulePatch PatchObject, options *ActionRulesUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AlertsManagement/actionRules/{actionRuleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionRuleName == "" {
		return nil, errors.New("parameter actionRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionRuleName}", url.PathEscape(actionRuleName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-05-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, actionRulePatch)
}

// updateHandleResponse handles the Update response.
func (client *ActionRulesClient) updateHandleResponse(resp *http.Response) (ActionRulesUpdateResponse, error) {
	result := ActionRulesUpdateResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionRule); err != nil {
		return ActionRulesUpdateResponse{}, err
	}
	return result, nil
}

// updateHandleError handles the Update error response.
func (client *ActionRulesClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

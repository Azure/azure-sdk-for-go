//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventhub

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

// EventHubsClient contains the methods for the EventHubs group.
// Don't use this type directly, use NewEventHubsClient() instead.
type EventHubsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewEventHubsClient creates a new instance of EventHubsClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewEventHubsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*EventHubsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &EventHubsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a new Event Hub as a nested resource within a Namespace.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - parameters - Parameters supplied to create an Event Hub resource.
//   - options - EventHubsClientCreateOrUpdateOptions contains the optional parameters for the EventHubsClient.CreateOrUpdate
//     method.
func (client *EventHubsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, parameters Eventhub, options *EventHubsClientCreateOrUpdateOptions) (EventHubsClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "EventHubsClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, parameters, options)
	if err != nil {
		return EventHubsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *EventHubsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, parameters Eventhub, options *EventHubsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *EventHubsClient) createOrUpdateHandleResponse(resp *http.Response) (EventHubsClientCreateOrUpdateResponse, error) {
	result := EventHubsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Eventhub); err != nil {
		return EventHubsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// CreateOrUpdateAuthorizationRule - Creates or updates an AuthorizationRule for the specified Event Hub. Creation/update
// of the AuthorizationRule will take a few seconds to take effect.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - authorizationRuleName - The authorization rule name.
//   - parameters - The shared access AuthorizationRule.
//   - options - EventHubsClientCreateOrUpdateAuthorizationRuleOptions contains the optional parameters for the EventHubsClient.CreateOrUpdateAuthorizationRule
//     method.
func (client *EventHubsClient) CreateOrUpdateAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters AuthorizationRule, options *EventHubsClientCreateOrUpdateAuthorizationRuleOptions) (EventHubsClientCreateOrUpdateAuthorizationRuleResponse, error) {
	var err error
	const operationName = "EventHubsClient.CreateOrUpdateAuthorizationRule"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateAuthorizationRuleCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, authorizationRuleName, parameters, options)
	if err != nil {
		return EventHubsClientCreateOrUpdateAuthorizationRuleResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientCreateOrUpdateAuthorizationRuleResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientCreateOrUpdateAuthorizationRuleResponse{}, err
	}
	resp, err := client.createOrUpdateAuthorizationRuleHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateAuthorizationRuleCreateRequest creates the CreateOrUpdateAuthorizationRule request.
func (client *EventHubsClient) createOrUpdateAuthorizationRuleCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters AuthorizationRule, options *EventHubsClientCreateOrUpdateAuthorizationRuleOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if authorizationRuleName == "" {
		return nil, errors.New("parameter authorizationRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationRuleName}", url.PathEscape(authorizationRuleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateAuthorizationRuleHandleResponse handles the CreateOrUpdateAuthorizationRule response.
func (client *EventHubsClient) createOrUpdateAuthorizationRuleHandleResponse(resp *http.Response) (EventHubsClientCreateOrUpdateAuthorizationRuleResponse, error) {
	result := EventHubsClientCreateOrUpdateAuthorizationRuleResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AuthorizationRule); err != nil {
		return EventHubsClientCreateOrUpdateAuthorizationRuleResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an Event Hub from the specified Namespace and resource group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - options - EventHubsClientDeleteOptions contains the optional parameters for the EventHubsClient.Delete method.
func (client *EventHubsClient) Delete(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, options *EventHubsClientDeleteOptions) (EventHubsClientDeleteResponse, error) {
	var err error
	const operationName = "EventHubsClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, options)
	if err != nil {
		return EventHubsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientDeleteResponse{}, err
	}
	return EventHubsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *EventHubsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, options *EventHubsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// DeleteAuthorizationRule - Deletes an Event Hub AuthorizationRule.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - authorizationRuleName - The authorization rule name.
//   - options - EventHubsClientDeleteAuthorizationRuleOptions contains the optional parameters for the EventHubsClient.DeleteAuthorizationRule
//     method.
func (client *EventHubsClient) DeleteAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, options *EventHubsClientDeleteAuthorizationRuleOptions) (EventHubsClientDeleteAuthorizationRuleResponse, error) {
	var err error
	const operationName = "EventHubsClient.DeleteAuthorizationRule"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteAuthorizationRuleCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, authorizationRuleName, options)
	if err != nil {
		return EventHubsClientDeleteAuthorizationRuleResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientDeleteAuthorizationRuleResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientDeleteAuthorizationRuleResponse{}, err
	}
	return EventHubsClientDeleteAuthorizationRuleResponse{}, nil
}

// deleteAuthorizationRuleCreateRequest creates the DeleteAuthorizationRule request.
func (client *EventHubsClient) deleteAuthorizationRuleCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, options *EventHubsClientDeleteAuthorizationRuleOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if authorizationRuleName == "" {
		return nil, errors.New("parameter authorizationRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationRuleName}", url.PathEscape(authorizationRuleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets an Event Hubs description for the specified Event Hub.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - options - EventHubsClientGetOptions contains the optional parameters for the EventHubsClient.Get method.
func (client *EventHubsClient) Get(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, options *EventHubsClientGetOptions) (EventHubsClientGetResponse, error) {
	var err error
	const operationName = "EventHubsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, options)
	if err != nil {
		return EventHubsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *EventHubsClient) getCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, options *EventHubsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EventHubsClient) getHandleResponse(resp *http.Response) (EventHubsClientGetResponse, error) {
	result := EventHubsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Eventhub); err != nil {
		return EventHubsClientGetResponse{}, err
	}
	return result, nil
}

// GetAuthorizationRule - Gets an AuthorizationRule for an Event Hub by rule name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - authorizationRuleName - The authorization rule name.
//   - options - EventHubsClientGetAuthorizationRuleOptions contains the optional parameters for the EventHubsClient.GetAuthorizationRule
//     method.
func (client *EventHubsClient) GetAuthorizationRule(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, options *EventHubsClientGetAuthorizationRuleOptions) (EventHubsClientGetAuthorizationRuleResponse, error) {
	var err error
	const operationName = "EventHubsClient.GetAuthorizationRule"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getAuthorizationRuleCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, authorizationRuleName, options)
	if err != nil {
		return EventHubsClientGetAuthorizationRuleResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientGetAuthorizationRuleResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientGetAuthorizationRuleResponse{}, err
	}
	resp, err := client.getAuthorizationRuleHandleResponse(httpResp)
	return resp, err
}

// getAuthorizationRuleCreateRequest creates the GetAuthorizationRule request.
func (client *EventHubsClient) getAuthorizationRuleCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, options *EventHubsClientGetAuthorizationRuleOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if authorizationRuleName == "" {
		return nil, errors.New("parameter authorizationRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationRuleName}", url.PathEscape(authorizationRuleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getAuthorizationRuleHandleResponse handles the GetAuthorizationRule response.
func (client *EventHubsClient) getAuthorizationRuleHandleResponse(resp *http.Response) (EventHubsClientGetAuthorizationRuleResponse, error) {
	result := EventHubsClientGetAuthorizationRuleResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AuthorizationRule); err != nil {
		return EventHubsClientGetAuthorizationRuleResponse{}, err
	}
	return result, nil
}

// NewListAuthorizationRulesPager - Gets the authorization rules for an Event Hub.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - options - EventHubsClientListAuthorizationRulesOptions contains the optional parameters for the EventHubsClient.NewListAuthorizationRulesPager
//     method.
func (client *EventHubsClient) NewListAuthorizationRulesPager(resourceGroupName string, namespaceName string, eventHubName string, options *EventHubsClientListAuthorizationRulesOptions) *runtime.Pager[EventHubsClientListAuthorizationRulesResponse] {
	return runtime.NewPager(runtime.PagingHandler[EventHubsClientListAuthorizationRulesResponse]{
		More: func(page EventHubsClientListAuthorizationRulesResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EventHubsClientListAuthorizationRulesResponse) (EventHubsClientListAuthorizationRulesResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "EventHubsClient.NewListAuthorizationRulesPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listAuthorizationRulesCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, options)
			}, nil)
			if err != nil {
				return EventHubsClientListAuthorizationRulesResponse{}, err
			}
			return client.listAuthorizationRulesHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listAuthorizationRulesCreateRequest creates the ListAuthorizationRules request.
func (client *EventHubsClient) listAuthorizationRulesCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, options *EventHubsClientListAuthorizationRulesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listAuthorizationRulesHandleResponse handles the ListAuthorizationRules response.
func (client *EventHubsClient) listAuthorizationRulesHandleResponse(resp *http.Response) (EventHubsClientListAuthorizationRulesResponse, error) {
	result := EventHubsClientListAuthorizationRulesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AuthorizationRuleListResult); err != nil {
		return EventHubsClientListAuthorizationRulesResponse{}, err
	}
	return result, nil
}

// NewListByNamespacePager - Gets all the Event Hubs in a Namespace.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - options - EventHubsClientListByNamespaceOptions contains the optional parameters for the EventHubsClient.NewListByNamespacePager
//     method.
func (client *EventHubsClient) NewListByNamespacePager(resourceGroupName string, namespaceName string, options *EventHubsClientListByNamespaceOptions) *runtime.Pager[EventHubsClientListByNamespaceResponse] {
	return runtime.NewPager(runtime.PagingHandler[EventHubsClientListByNamespaceResponse]{
		More: func(page EventHubsClientListByNamespaceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EventHubsClientListByNamespaceResponse) (EventHubsClientListByNamespaceResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "EventHubsClient.NewListByNamespacePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByNamespaceCreateRequest(ctx, resourceGroupName, namespaceName, options)
			}, nil)
			if err != nil {
				return EventHubsClientListByNamespaceResponse{}, err
			}
			return client.listByNamespaceHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByNamespaceCreateRequest creates the ListByNamespace request.
func (client *EventHubsClient) listByNamespaceCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, options *EventHubsClientListByNamespaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByNamespaceHandleResponse handles the ListByNamespace response.
func (client *EventHubsClient) listByNamespaceHandleResponse(resp *http.Response) (EventHubsClientListByNamespaceResponse, error) {
	result := EventHubsClientListByNamespaceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ListResult); err != nil {
		return EventHubsClientListByNamespaceResponse{}, err
	}
	return result, nil
}

// ListKeys - Gets the ACS and SAS connection strings for the Event Hub.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - authorizationRuleName - The authorization rule name.
//   - options - EventHubsClientListKeysOptions contains the optional parameters for the EventHubsClient.ListKeys method.
func (client *EventHubsClient) ListKeys(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, options *EventHubsClientListKeysOptions) (EventHubsClientListKeysResponse, error) {
	var err error
	const operationName = "EventHubsClient.ListKeys"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listKeysCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, authorizationRuleName, options)
	if err != nil {
		return EventHubsClientListKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientListKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientListKeysResponse{}, err
	}
	resp, err := client.listKeysHandleResponse(httpResp)
	return resp, err
}

// listKeysCreateRequest creates the ListKeys request.
func (client *EventHubsClient) listKeysCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, options *EventHubsClientListKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}/listKeys"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if authorizationRuleName == "" {
		return nil, errors.New("parameter authorizationRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationRuleName}", url.PathEscape(authorizationRuleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listKeysHandleResponse handles the ListKeys response.
func (client *EventHubsClient) listKeysHandleResponse(resp *http.Response) (EventHubsClientListKeysResponse, error) {
	result := EventHubsClientListKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessKeys); err != nil {
		return EventHubsClientListKeysResponse{}, err
	}
	return result, nil
}

// RegenerateKeys - Regenerates the ACS and SAS connection strings for the Event Hub.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - eventHubName - The Event Hub name
//   - authorizationRuleName - The authorization rule name.
//   - parameters - Parameters supplied to regenerate the AuthorizationRule Keys (PrimaryKey/SecondaryKey).
//   - options - EventHubsClientRegenerateKeysOptions contains the optional parameters for the EventHubsClient.RegenerateKeys
//     method.
func (client *EventHubsClient) RegenerateKeys(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters RegenerateAccessKeyParameters, options *EventHubsClientRegenerateKeysOptions) (EventHubsClientRegenerateKeysResponse, error) {
	var err error
	const operationName = "EventHubsClient.RegenerateKeys"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.regenerateKeysCreateRequest(ctx, resourceGroupName, namespaceName, eventHubName, authorizationRuleName, parameters, options)
	if err != nil {
		return EventHubsClientRegenerateKeysResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EventHubsClientRegenerateKeysResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return EventHubsClientRegenerateKeysResponse{}, err
	}
	resp, err := client.regenerateKeysHandleResponse(httpResp)
	return resp, err
}

// regenerateKeysCreateRequest creates the RegenerateKeys request.
func (client *EventHubsClient) regenerateKeysCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, eventHubName string, authorizationRuleName string, parameters RegenerateAccessKeyParameters, options *EventHubsClientRegenerateKeysOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/eventhubs/{eventHubName}/authorizationRules/{authorizationRuleName}/regenerateKeys"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if eventHubName == "" {
		return nil, errors.New("parameter eventHubName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventHubName}", url.PathEscape(eventHubName))
	if authorizationRuleName == "" {
		return nil, errors.New("parameter authorizationRuleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationRuleName}", url.PathEscape(authorizationRuleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// regenerateKeysHandleResponse handles the RegenerateKeys response.
func (client *EventHubsClient) regenerateKeysHandleResponse(resp *http.Response) (EventHubsClientRegenerateKeysResponse, error) {
	result := EventHubsClientRegenerateKeysResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessKeys); err != nil {
		return EventHubsClientRegenerateKeysResponse{}, err
	}
	return result, nil
}

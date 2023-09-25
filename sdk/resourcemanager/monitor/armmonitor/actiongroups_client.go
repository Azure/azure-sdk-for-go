//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

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

// ActionGroupsClient contains the methods for the ActionGroups group.
// Don't use this type directly, use NewActionGroupsClient() instead.
type ActionGroupsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewActionGroupsClient creates a new instance of ActionGroupsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewActionGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ActionGroupsClient, error) {
	cl, err := arm.NewClient(moduleName+".ActionGroupsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ActionGroupsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateNotificationsAtActionGroupResourceLevel - Send test notifications to a set of provided receivers
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - notificationRequest - The notification request body which includes the contact details
//   - options - ActionGroupsClientBeginCreateNotificationsAtActionGroupResourceLevelOptions contains the optional parameters
//     for the ActionGroupsClient.BeginCreateNotificationsAtActionGroupResourceLevel method.
func (client *ActionGroupsClient) BeginCreateNotificationsAtActionGroupResourceLevel(ctx context.Context, resourceGroupName string, actionGroupName string, notificationRequest NotificationRequestBody, options *ActionGroupsClientBeginCreateNotificationsAtActionGroupResourceLevelOptions) (*runtime.Poller[ActionGroupsClientCreateNotificationsAtActionGroupResourceLevelResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createNotificationsAtActionGroupResourceLevel(ctx, resourceGroupName, actionGroupName, notificationRequest, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ActionGroupsClientCreateNotificationsAtActionGroupResourceLevelResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ActionGroupsClientCreateNotificationsAtActionGroupResourceLevelResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateNotificationsAtActionGroupResourceLevel - Send test notifications to a set of provided receivers
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
func (client *ActionGroupsClient) createNotificationsAtActionGroupResourceLevel(ctx context.Context, resourceGroupName string, actionGroupName string, notificationRequest NotificationRequestBody, options *ActionGroupsClientBeginCreateNotificationsAtActionGroupResourceLevelOptions) (*http.Response, error) {
	var err error
	const operationName = "ActionGroupsClient.BeginCreateNotificationsAtActionGroupResourceLevel"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createNotificationsAtActionGroupResourceLevelCreateRequest(ctx, resourceGroupName, actionGroupName, notificationRequest, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createNotificationsAtActionGroupResourceLevelCreateRequest creates the CreateNotificationsAtActionGroupResourceLevel request.
func (client *ActionGroupsClient) createNotificationsAtActionGroupResourceLevelCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, notificationRequest NotificationRequestBody, options *ActionGroupsClientBeginCreateNotificationsAtActionGroupResourceLevelOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}/createNotifications"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, notificationRequest); err != nil {
	return nil, err
}
	return req, nil
}

// CreateOrUpdate - Create a new action group or update an existing one.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - actionGroup - The action group to create or use for the update.
//   - options - ActionGroupsClientCreateOrUpdateOptions contains the optional parameters for the ActionGroupsClient.CreateOrUpdate
//     method.
func (client *ActionGroupsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, actionGroupName string, actionGroup ActionGroupResource, options *ActionGroupsClientCreateOrUpdateOptions) (ActionGroupsClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "ActionGroupsClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, actionGroupName, actionGroup, options)
	if err != nil {
		return ActionGroupsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ActionGroupsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ActionGroupsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ActionGroupsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, actionGroup ActionGroupResource, options *ActionGroupsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, actionGroup); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ActionGroupsClient) createOrUpdateHandleResponse(resp *http.Response) (ActionGroupsClientCreateOrUpdateResponse, error) {
	result := ActionGroupsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionGroupResource); err != nil {
		return ActionGroupsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete an action group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - options - ActionGroupsClientDeleteOptions contains the optional parameters for the ActionGroupsClient.Delete method.
func (client *ActionGroupsClient) Delete(ctx context.Context, resourceGroupName string, actionGroupName string, options *ActionGroupsClientDeleteOptions) (ActionGroupsClientDeleteResponse, error) {
	var err error
	const operationName = "ActionGroupsClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, actionGroupName, options)
	if err != nil {
		return ActionGroupsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ActionGroupsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return ActionGroupsClientDeleteResponse{}, err
	}
	return ActionGroupsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ActionGroupsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, options *ActionGroupsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// EnableReceiver - Enable a receiver in an action group. This changes the receiver's status from Disabled to Enabled. This
// operation is only supported for Email or SMS receivers.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - enableRequest - The receiver to re-enable.
//   - options - ActionGroupsClientEnableReceiverOptions contains the optional parameters for the ActionGroupsClient.EnableReceiver
//     method.
func (client *ActionGroupsClient) EnableReceiver(ctx context.Context, resourceGroupName string, actionGroupName string, enableRequest EnableRequest, options *ActionGroupsClientEnableReceiverOptions) (ActionGroupsClientEnableReceiverResponse, error) {
	var err error
	const operationName = "ActionGroupsClient.EnableReceiver"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.enableReceiverCreateRequest(ctx, resourceGroupName, actionGroupName, enableRequest, options)
	if err != nil {
		return ActionGroupsClientEnableReceiverResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ActionGroupsClientEnableReceiverResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ActionGroupsClientEnableReceiverResponse{}, err
	}
	return ActionGroupsClientEnableReceiverResponse{}, nil
}

// enableReceiverCreateRequest creates the EnableReceiver request.
func (client *ActionGroupsClient) enableReceiverCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, enableRequest EnableRequest, options *ActionGroupsClientEnableReceiverOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}/subscribe"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, enableRequest); err != nil {
	return nil, err
}
	return req, nil
}

// Get - Get an action group.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - options - ActionGroupsClientGetOptions contains the optional parameters for the ActionGroupsClient.Get method.
func (client *ActionGroupsClient) Get(ctx context.Context, resourceGroupName string, actionGroupName string, options *ActionGroupsClientGetOptions) (ActionGroupsClientGetResponse, error) {
	var err error
	const operationName = "ActionGroupsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, actionGroupName, options)
	if err != nil {
		return ActionGroupsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ActionGroupsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ActionGroupsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ActionGroupsClient) getCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, options *ActionGroupsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ActionGroupsClient) getHandleResponse(resp *http.Response) (ActionGroupsClientGetResponse, error) {
	result := ActionGroupsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionGroupResource); err != nil {
		return ActionGroupsClientGetResponse{}, err
	}
	return result, nil
}

// GetTestNotificationsAtActionGroupResourceLevel - Get the test notifications by the notification id
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - notificationID - The notification id
//   - options - ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelOptions contains the optional parameters for
//     the ActionGroupsClient.GetTestNotificationsAtActionGroupResourceLevel method.
func (client *ActionGroupsClient) GetTestNotificationsAtActionGroupResourceLevel(ctx context.Context, resourceGroupName string, actionGroupName string, notificationID string, options *ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelOptions) (ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse, error) {
	var err error
	const operationName = "ActionGroupsClient.GetTestNotificationsAtActionGroupResourceLevel"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getTestNotificationsAtActionGroupResourceLevelCreateRequest(ctx, resourceGroupName, actionGroupName, notificationID, options)
	if err != nil {
		return ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse{}, err
	}
	resp, err := client.getTestNotificationsAtActionGroupResourceLevelHandleResponse(httpResp)
	return resp, err
}

// getTestNotificationsAtActionGroupResourceLevelCreateRequest creates the GetTestNotificationsAtActionGroupResourceLevel request.
func (client *ActionGroupsClient) getTestNotificationsAtActionGroupResourceLevelCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, notificationID string, options *ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}/notificationStatus/{notificationId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	if notificationID == "" {
		return nil, errors.New("parameter notificationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{notificationId}", url.PathEscape(notificationID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getTestNotificationsAtActionGroupResourceLevelHandleResponse handles the GetTestNotificationsAtActionGroupResourceLevel response.
func (client *ActionGroupsClient) getTestNotificationsAtActionGroupResourceLevelHandleResponse(resp *http.Response) (ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse, error) {
	result := ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TestNotificationDetailsResponse); err != nil {
		return ActionGroupsClientGetTestNotificationsAtActionGroupResourceLevelResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Get a list of all action groups in a resource group.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - ActionGroupsClientListByResourceGroupOptions contains the optional parameters for the ActionGroupsClient.NewListByResourceGroupPager
//     method.
func (client *ActionGroupsClient) NewListByResourceGroupPager(resourceGroupName string, options *ActionGroupsClientListByResourceGroupOptions) (*runtime.Pager[ActionGroupsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ActionGroupsClientListByResourceGroupResponse]{
		More: func(page ActionGroupsClientListByResourceGroupResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ActionGroupsClientListByResourceGroupResponse) (ActionGroupsClientListByResourceGroupResponse, error) {
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ActionGroupsClient.NewListByResourceGroupPager")
			req, err := client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			if err != nil {
				return ActionGroupsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ActionGroupsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ActionGroupsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ActionGroupsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ActionGroupsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ActionGroupsClient) listByResourceGroupHandleResponse(resp *http.Response) (ActionGroupsClientListByResourceGroupResponse, error) {
	result := ActionGroupsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionGroupList); err != nil {
		return ActionGroupsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionIDPager - Get a list of all action groups in a subscription.
//
// Generated from API version 2023-01-01
//   - options - ActionGroupsClientListBySubscriptionIDOptions contains the optional parameters for the ActionGroupsClient.NewListBySubscriptionIDPager
//     method.
func (client *ActionGroupsClient) NewListBySubscriptionIDPager(options *ActionGroupsClientListBySubscriptionIDOptions) (*runtime.Pager[ActionGroupsClientListBySubscriptionIDResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ActionGroupsClientListBySubscriptionIDResponse]{
		More: func(page ActionGroupsClientListBySubscriptionIDResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ActionGroupsClientListBySubscriptionIDResponse) (ActionGroupsClientListBySubscriptionIDResponse, error) {
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ActionGroupsClient.NewListBySubscriptionIDPager")
			req, err := client.listBySubscriptionIDCreateRequest(ctx, options)
			if err != nil {
				return ActionGroupsClientListBySubscriptionIDResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ActionGroupsClientListBySubscriptionIDResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ActionGroupsClientListBySubscriptionIDResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionIDHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listBySubscriptionIDCreateRequest creates the ListBySubscriptionID request.
func (client *ActionGroupsClient) listBySubscriptionIDCreateRequest(ctx context.Context, options *ActionGroupsClientListBySubscriptionIDOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Insights/actionGroups"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionIDHandleResponse handles the ListBySubscriptionID response.
func (client *ActionGroupsClient) listBySubscriptionIDHandleResponse(resp *http.Response) (ActionGroupsClientListBySubscriptionIDResponse, error) {
	result := ActionGroupsClientListBySubscriptionIDResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionGroupList); err != nil {
		return ActionGroupsClientListBySubscriptionIDResponse{}, err
	}
	return result, nil
}

// Update - Updates an existing action group's tags. To update other fields use the CreateOrUpdate method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-01-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - actionGroupName - The name of the action group.
//   - actionGroupPatch - Parameters supplied to the operation.
//   - options - ActionGroupsClientUpdateOptions contains the optional parameters for the ActionGroupsClient.Update method.
func (client *ActionGroupsClient) Update(ctx context.Context, resourceGroupName string, actionGroupName string, actionGroupPatch ActionGroupPatchBody, options *ActionGroupsClientUpdateOptions) (ActionGroupsClientUpdateResponse, error) {
	var err error
	const operationName = "ActionGroupsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, actionGroupName, actionGroupPatch, options)
	if err != nil {
		return ActionGroupsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ActionGroupsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ActionGroupsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *ActionGroupsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, actionGroupName string, actionGroupPatch ActionGroupPatchBody, options *ActionGroupsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/actionGroups/{actionGroupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if actionGroupName == "" {
		return nil, errors.New("parameter actionGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionGroupName}", url.PathEscape(actionGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-01-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, actionGroupPatch); err != nil {
	return nil, err
}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *ActionGroupsClient) updateHandleResponse(resp *http.Response) (ActionGroupsClientUpdateResponse, error) {
	result := ActionGroupsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ActionGroupResource); err != nil {
		return ActionGroupsClientUpdateResponse{}, err
	}
	return result, nil
}


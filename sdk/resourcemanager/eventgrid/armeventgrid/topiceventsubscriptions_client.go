// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventgrid

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

// TopicEventSubscriptionsClient contains the methods for the TopicEventSubscriptions group.
// Don't use this type directly, use NewTopicEventSubscriptionsClient() instead.
type TopicEventSubscriptionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewTopicEventSubscriptionsClient creates a new instance of TopicEventSubscriptionsClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewTopicEventSubscriptionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TopicEventSubscriptionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &TopicEventSubscriptionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Asynchronously creates a new event subscription or updates an existing event subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the domain topic.
//   - eventSubscriptionName - Name of the event subscription to be created. Event subscription names must be between 3 and 64
//     characters in length and use alphanumeric letters only.
//   - eventSubscriptionInfo - Event subscription properties containing the destination and filter information.
//   - options - TopicEventSubscriptionsClientBeginCreateOrUpdateOptions contains the optional parameters for the TopicEventSubscriptionsClient.BeginCreateOrUpdate
//     method.
func (client *TopicEventSubscriptionsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, eventSubscriptionInfo EventSubscription, options *TopicEventSubscriptionsClientBeginCreateOrUpdateOptions) (*runtime.Poller[TopicEventSubscriptionsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, topicName, eventSubscriptionName, eventSubscriptionInfo, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[TopicEventSubscriptionsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[TopicEventSubscriptionsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Asynchronously creates a new event subscription or updates an existing event subscription.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
func (client *TopicEventSubscriptionsClient) createOrUpdate(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, eventSubscriptionInfo EventSubscription, options *TopicEventSubscriptionsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "TopicEventSubscriptionsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, topicName, eventSubscriptionName, eventSubscriptionInfo, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *TopicEventSubscriptionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, eventSubscriptionInfo EventSubscription, _ *TopicEventSubscriptionsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions/{eventSubscriptionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if eventSubscriptionName == "" {
		return nil, errors.New("parameter eventSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventSubscriptionName}", url.PathEscape(eventSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, eventSubscriptionInfo); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Delete an existing event subscription for a topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the topic.
//   - eventSubscriptionName - Name of the event subscription to be deleted.
//   - options - TopicEventSubscriptionsClientBeginDeleteOptions contains the optional parameters for the TopicEventSubscriptionsClient.BeginDelete
//     method.
func (client *TopicEventSubscriptionsClient) BeginDelete(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, options *TopicEventSubscriptionsClientBeginDeleteOptions) (*runtime.Poller[TopicEventSubscriptionsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, topicName, eventSubscriptionName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[TopicEventSubscriptionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[TopicEventSubscriptionsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete an existing event subscription for a topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
func (client *TopicEventSubscriptionsClient) deleteOperation(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, options *TopicEventSubscriptionsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "TopicEventSubscriptionsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, topicName, eventSubscriptionName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *TopicEventSubscriptionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, _ *TopicEventSubscriptionsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions/{eventSubscriptionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if eventSubscriptionName == "" {
		return nil, errors.New("parameter eventSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventSubscriptionName}", url.PathEscape(eventSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Get properties of an event subscription of a topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the topic.
//   - eventSubscriptionName - Name of the event subscription to be found.
//   - options - TopicEventSubscriptionsClientGetOptions contains the optional parameters for the TopicEventSubscriptionsClient.Get
//     method.
func (client *TopicEventSubscriptionsClient) Get(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, options *TopicEventSubscriptionsClientGetOptions) (TopicEventSubscriptionsClientGetResponse, error) {
	var err error
	const operationName = "TopicEventSubscriptionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, topicName, eventSubscriptionName, options)
	if err != nil {
		return TopicEventSubscriptionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TopicEventSubscriptionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TopicEventSubscriptionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *TopicEventSubscriptionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, _ *TopicEventSubscriptionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions/{eventSubscriptionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if eventSubscriptionName == "" {
		return nil, errors.New("parameter eventSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventSubscriptionName}", url.PathEscape(eventSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *TopicEventSubscriptionsClient) getHandleResponse(resp *http.Response) (TopicEventSubscriptionsClientGetResponse, error) {
	result := TopicEventSubscriptionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventSubscription); err != nil {
		return TopicEventSubscriptionsClientGetResponse{}, err
	}
	return result, nil
}

// GetDeliveryAttributes - Get all delivery attributes for an event subscription for topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the topic.
//   - eventSubscriptionName - Name of the event subscription.
//   - options - TopicEventSubscriptionsClientGetDeliveryAttributesOptions contains the optional parameters for the TopicEventSubscriptionsClient.GetDeliveryAttributes
//     method.
func (client *TopicEventSubscriptionsClient) GetDeliveryAttributes(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, options *TopicEventSubscriptionsClientGetDeliveryAttributesOptions) (TopicEventSubscriptionsClientGetDeliveryAttributesResponse, error) {
	var err error
	const operationName = "TopicEventSubscriptionsClient.GetDeliveryAttributes"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getDeliveryAttributesCreateRequest(ctx, resourceGroupName, topicName, eventSubscriptionName, options)
	if err != nil {
		return TopicEventSubscriptionsClientGetDeliveryAttributesResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TopicEventSubscriptionsClientGetDeliveryAttributesResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TopicEventSubscriptionsClientGetDeliveryAttributesResponse{}, err
	}
	resp, err := client.getDeliveryAttributesHandleResponse(httpResp)
	return resp, err
}

// getDeliveryAttributesCreateRequest creates the GetDeliveryAttributes request.
func (client *TopicEventSubscriptionsClient) getDeliveryAttributesCreateRequest(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, _ *TopicEventSubscriptionsClientGetDeliveryAttributesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions/{eventSubscriptionName}/getDeliveryAttributes"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if eventSubscriptionName == "" {
		return nil, errors.New("parameter eventSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventSubscriptionName}", url.PathEscape(eventSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getDeliveryAttributesHandleResponse handles the GetDeliveryAttributes response.
func (client *TopicEventSubscriptionsClient) getDeliveryAttributesHandleResponse(resp *http.Response) (TopicEventSubscriptionsClientGetDeliveryAttributesResponse, error) {
	result := TopicEventSubscriptionsClientGetDeliveryAttributesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeliveryAttributeListResult); err != nil {
		return TopicEventSubscriptionsClientGetDeliveryAttributesResponse{}, err
	}
	return result, nil
}

// GetFullURL - Get the full endpoint URL for an event subscription for topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the domain topic.
//   - eventSubscriptionName - Name of the event subscription.
//   - options - TopicEventSubscriptionsClientGetFullURLOptions contains the optional parameters for the TopicEventSubscriptionsClient.GetFullURL
//     method.
func (client *TopicEventSubscriptionsClient) GetFullURL(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, options *TopicEventSubscriptionsClientGetFullURLOptions) (TopicEventSubscriptionsClientGetFullURLResponse, error) {
	var err error
	const operationName = "TopicEventSubscriptionsClient.GetFullURL"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getFullURLCreateRequest(ctx, resourceGroupName, topicName, eventSubscriptionName, options)
	if err != nil {
		return TopicEventSubscriptionsClientGetFullURLResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return TopicEventSubscriptionsClientGetFullURLResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return TopicEventSubscriptionsClientGetFullURLResponse{}, err
	}
	resp, err := client.getFullURLHandleResponse(httpResp)
	return resp, err
}

// getFullURLCreateRequest creates the GetFullURL request.
func (client *TopicEventSubscriptionsClient) getFullURLCreateRequest(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, _ *TopicEventSubscriptionsClientGetFullURLOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions/{eventSubscriptionName}/getFullUrl"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if eventSubscriptionName == "" {
		return nil, errors.New("parameter eventSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventSubscriptionName}", url.PathEscape(eventSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getFullURLHandleResponse handles the GetFullURL response.
func (client *TopicEventSubscriptionsClient) getFullURLHandleResponse(resp *http.Response) (TopicEventSubscriptionsClientGetFullURLResponse, error) {
	result := TopicEventSubscriptionsClientGetFullURLResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventSubscriptionFullURL); err != nil {
		return TopicEventSubscriptionsClientGetFullURLResponse{}, err
	}
	return result, nil
}

// NewListPager - List all event subscriptions that have been created for a specific topic.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the topic.
//   - options - TopicEventSubscriptionsClientListOptions contains the optional parameters for the TopicEventSubscriptionsClient.NewListPager
//     method.
func (client *TopicEventSubscriptionsClient) NewListPager(resourceGroupName string, topicName string, options *TopicEventSubscriptionsClientListOptions) *runtime.Pager[TopicEventSubscriptionsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[TopicEventSubscriptionsClientListResponse]{
		More: func(page TopicEventSubscriptionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TopicEventSubscriptionsClientListResponse) (TopicEventSubscriptionsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "TopicEventSubscriptionsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, topicName, options)
			}, nil)
			if err != nil {
				return TopicEventSubscriptionsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *TopicEventSubscriptionsClient) listCreateRequest(ctx context.Context, resourceGroupName string, topicName string, options *TopicEventSubscriptionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *TopicEventSubscriptionsClient) listHandleResponse(resp *http.Response) (TopicEventSubscriptionsClientListResponse, error) {
	result := TopicEventSubscriptionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventSubscriptionsListResult); err != nil {
		return TopicEventSubscriptionsClientListResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update an existing event subscription for a topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
//   - resourceGroupName - The name of the resource group within the user's subscription.
//   - topicName - Name of the domain.
//   - eventSubscriptionName - Name of the event subscription to be updated.
//   - eventSubscriptionUpdateParameters - Updated event subscription information.
//   - options - TopicEventSubscriptionsClientBeginUpdateOptions contains the optional parameters for the TopicEventSubscriptionsClient.BeginUpdate
//     method.
func (client *TopicEventSubscriptionsClient) BeginUpdate(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, eventSubscriptionUpdateParameters EventSubscriptionUpdateParameters, options *TopicEventSubscriptionsClientBeginUpdateOptions) (*runtime.Poller[TopicEventSubscriptionsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, topicName, eventSubscriptionName, eventSubscriptionUpdateParameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[TopicEventSubscriptionsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[TopicEventSubscriptionsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Update an existing event subscription for a topic.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-02-15
func (client *TopicEventSubscriptionsClient) update(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, eventSubscriptionUpdateParameters EventSubscriptionUpdateParameters, options *TopicEventSubscriptionsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "TopicEventSubscriptionsClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, topicName, eventSubscriptionName, eventSubscriptionUpdateParameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *TopicEventSubscriptionsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, topicName string, eventSubscriptionName string, eventSubscriptionUpdateParameters EventSubscriptionUpdateParameters, _ *TopicEventSubscriptionsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/topics/{topicName}/eventSubscriptions/{eventSubscriptionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if eventSubscriptionName == "" {
		return nil, errors.New("parameter eventSubscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{eventSubscriptionName}", url.PathEscape(eventSubscriptionName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, eventSubscriptionUpdateParameters); err != nil {
		return nil, err
	}
	return req, nil
}

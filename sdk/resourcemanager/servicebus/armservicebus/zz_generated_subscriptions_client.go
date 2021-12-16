//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicebus

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// SubscriptionsClient contains the methods for the Subscriptions group.
// Don't use this type directly, use NewSubscriptionsClient() instead.
type SubscriptionsClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewSubscriptionsClient creates a new instance of SubscriptionsClient with the specified values.
func NewSubscriptionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *SubscriptionsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &SubscriptionsClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// CreateOrUpdate - Creates a topic subscription.
// If the operation fails it returns the *ErrorResponse error type.
func (client *SubscriptionsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, parameters SBSubscription, options *SubscriptionsCreateOrUpdateOptions) (SubscriptionsCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, namespaceName, topicName, subscriptionName, parameters, options)
	if err != nil {
		return SubscriptionsCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SubscriptionsCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SubscriptionsCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SubscriptionsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, parameters SBSubscription, options *SubscriptionsCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/subscriptions/{subscriptionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if subscriptionName == "" {
		return nil, errors.New("parameter subscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionName}", url.PathEscape(subscriptionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *SubscriptionsClient) createOrUpdateHandleResponse(resp *http.Response) (SubscriptionsCreateOrUpdateResponse, error) {
	result := SubscriptionsCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SBSubscription); err != nil {
		return SubscriptionsCreateOrUpdateResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *SubscriptionsClient) createOrUpdateHandleError(resp *http.Response) error {
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

// Delete - Deletes a subscription from the specified topic.
// If the operation fails it returns the *ErrorResponse error type.
func (client *SubscriptionsClient) Delete(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, options *SubscriptionsDeleteOptions) (SubscriptionsDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, namespaceName, topicName, subscriptionName, options)
	if err != nil {
		return SubscriptionsDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SubscriptionsDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return SubscriptionsDeleteResponse{}, client.deleteHandleError(resp)
	}
	return SubscriptionsDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SubscriptionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, options *SubscriptionsDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/subscriptions/{subscriptionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if subscriptionName == "" {
		return nil, errors.New("parameter subscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionName}", url.PathEscape(subscriptionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *SubscriptionsClient) deleteHandleError(resp *http.Response) error {
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

// Get - Returns a subscription description for the specified topic.
// If the operation fails it returns the *ErrorResponse error type.
func (client *SubscriptionsClient) Get(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, options *SubscriptionsGetOptions) (SubscriptionsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, namespaceName, topicName, subscriptionName, options)
	if err != nil {
		return SubscriptionsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return SubscriptionsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SubscriptionsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SubscriptionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, subscriptionName string, options *SubscriptionsGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/subscriptions/{subscriptionName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if subscriptionName == "" {
		return nil, errors.New("parameter subscriptionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionName}", url.PathEscape(subscriptionName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SubscriptionsClient) getHandleResponse(resp *http.Response) (SubscriptionsGetResponse, error) {
	result := SubscriptionsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SBSubscription); err != nil {
		return SubscriptionsGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *SubscriptionsClient) getHandleError(resp *http.Response) error {
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

// ListByTopic - List all the subscriptions under a specified topic.
// If the operation fails it returns the *ErrorResponse error type.
func (client *SubscriptionsClient) ListByTopic(resourceGroupName string, namespaceName string, topicName string, options *SubscriptionsListByTopicOptions) *SubscriptionsListByTopicPager {
	return &SubscriptionsListByTopicPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByTopicCreateRequest(ctx, resourceGroupName, namespaceName, topicName, options)
		},
		advancer: func(ctx context.Context, resp SubscriptionsListByTopicResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.SBSubscriptionListResult.NextLink)
		},
	}
}

// listByTopicCreateRequest creates the ListByTopic request.
func (client *SubscriptionsClient) listByTopicCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, topicName string, options *SubscriptionsListByTopicOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{namespaceName}/topics/{topicName}/subscriptions"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if topicName == "" {
		return nil, errors.New("parameter topicName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{topicName}", url.PathEscape(topicName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01-preview")
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByTopicHandleResponse handles the ListByTopic response.
func (client *SubscriptionsClient) listByTopicHandleResponse(resp *http.Response) (SubscriptionsListByTopicResponse, error) {
	result := SubscriptionsListByTopicResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SBSubscriptionListResult); err != nil {
		return SubscriptionsListByTopicResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listByTopicHandleError handles the ListByTopic error response.
func (client *SubscriptionsClient) listByTopicHandleError(resp *http.Response) error {
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

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

// ChannelsClient contains the methods for the Channels group.
// Don't use this type directly, use NewChannelsClient() instead.
type ChannelsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewChannelsClient creates a new instance of ChannelsClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewChannelsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ChannelsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ChannelsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Synchronously creates or updates a new channel with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-04-01-preview
//   - resourceGroupName - The name of the resource group within the partners subscription.
//   - partnerNamespaceName - Name of the partner namespace.
//   - channelName - Name of the channel.
//   - channelInfo - Channel information.
//   - options - ChannelsClientCreateOrUpdateOptions contains the optional parameters for the ChannelsClient.CreateOrUpdate method.
func (client *ChannelsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, channelInfo Channel, options *ChannelsClientCreateOrUpdateOptions) (ChannelsClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "ChannelsClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, partnerNamespaceName, channelName, channelInfo, options)
	if err != nil {
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ChannelsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, channelInfo Channel, _ *ChannelsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/partnerNamespaces/{partnerNamespaceName}/channels/{channelName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if partnerNamespaceName == "" {
		return nil, errors.New("parameter partnerNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerNamespaceName}", url.PathEscape(partnerNamespaceName))
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, channelInfo); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ChannelsClient) createOrUpdateHandleResponse(resp *http.Response) (ChannelsClientCreateOrUpdateResponse, error) {
	result := ChannelsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Channel); err != nil {
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// BeginDelete - Delete an existing channel.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-04-01-preview
//   - resourceGroupName - The name of the resource group within the partners subscription.
//   - partnerNamespaceName - Name of the partner namespace.
//   - channelName - Name of the channel.
//   - options - ChannelsClientBeginDeleteOptions contains the optional parameters for the ChannelsClient.BeginDelete method.
func (client *ChannelsClient) BeginDelete(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, options *ChannelsClientBeginDeleteOptions) (*runtime.Poller[ChannelsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, partnerNamespaceName, channelName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ChannelsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ChannelsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Delete an existing channel.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-04-01-preview
func (client *ChannelsClient) deleteOperation(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, options *ChannelsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ChannelsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, partnerNamespaceName, channelName, options)
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
func (client *ChannelsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, _ *ChannelsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/partnerNamespaces/{partnerNamespaceName}/channels/{channelName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if partnerNamespaceName == "" {
		return nil, errors.New("parameter partnerNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerNamespaceName}", url.PathEscape(partnerNamespaceName))
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Get properties of a channel.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-04-01-preview
//   - resourceGroupName - The name of the resource group within the partners subscription.
//   - partnerNamespaceName - Name of the partner namespace.
//   - channelName - Name of the channel.
//   - options - ChannelsClientGetOptions contains the optional parameters for the ChannelsClient.Get method.
func (client *ChannelsClient) Get(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, options *ChannelsClientGetOptions) (ChannelsClientGetResponse, error) {
	var err error
	const operationName = "ChannelsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, partnerNamespaceName, channelName, options)
	if err != nil {
		return ChannelsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ChannelsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ChannelsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ChannelsClient) getCreateRequest(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, _ *ChannelsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/partnerNamespaces/{partnerNamespaceName}/channels/{channelName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if partnerNamespaceName == "" {
		return nil, errors.New("parameter partnerNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerNamespaceName}", url.PathEscape(partnerNamespaceName))
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ChannelsClient) getHandleResponse(resp *http.Response) (ChannelsClientGetResponse, error) {
	result := ChannelsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Channel); err != nil {
		return ChannelsClientGetResponse{}, err
	}
	return result, nil
}

// GetFullURL - Get the full endpoint URL of a partner destination channel.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-04-01-preview
//   - resourceGroupName - The name of the resource group within the partners subscription.
//   - partnerNamespaceName - Name of the partner namespace.
//   - channelName - Name of the Channel.
//   - options - ChannelsClientGetFullURLOptions contains the optional parameters for the ChannelsClient.GetFullURL method.
func (client *ChannelsClient) GetFullURL(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, options *ChannelsClientGetFullURLOptions) (ChannelsClientGetFullURLResponse, error) {
	var err error
	const operationName = "ChannelsClient.GetFullURL"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getFullURLCreateRequest(ctx, resourceGroupName, partnerNamespaceName, channelName, options)
	if err != nil {
		return ChannelsClientGetFullURLResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ChannelsClientGetFullURLResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ChannelsClientGetFullURLResponse{}, err
	}
	resp, err := client.getFullURLHandleResponse(httpResp)
	return resp, err
}

// getFullURLCreateRequest creates the GetFullURL request.
func (client *ChannelsClient) getFullURLCreateRequest(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, _ *ChannelsClientGetFullURLOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/partnerNamespaces/{partnerNamespaceName}/channels/{channelName}/getFullUrl"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if partnerNamespaceName == "" {
		return nil, errors.New("parameter partnerNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerNamespaceName}", url.PathEscape(partnerNamespaceName))
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getFullURLHandleResponse handles the GetFullURL response.
func (client *ChannelsClient) getFullURLHandleResponse(resp *http.Response) (ChannelsClientGetFullURLResponse, error) {
	result := ChannelsClientGetFullURLResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventSubscriptionFullURL); err != nil {
		return ChannelsClientGetFullURLResponse{}, err
	}
	return result, nil
}

// NewListByPartnerNamespacePager - List all the channels in a partner namespace.
//
// Generated from API version 2025-04-01-preview
//   - resourceGroupName - The name of the resource group within the partners subscription.
//   - partnerNamespaceName - Name of the partner namespace.
//   - options - ChannelsClientListByPartnerNamespaceOptions contains the optional parameters for the ChannelsClient.NewListByPartnerNamespacePager
//     method.
func (client *ChannelsClient) NewListByPartnerNamespacePager(resourceGroupName string, partnerNamespaceName string, options *ChannelsClientListByPartnerNamespaceOptions) *runtime.Pager[ChannelsClientListByPartnerNamespaceResponse] {
	return runtime.NewPager(runtime.PagingHandler[ChannelsClientListByPartnerNamespaceResponse]{
		More: func(page ChannelsClientListByPartnerNamespaceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ChannelsClientListByPartnerNamespaceResponse) (ChannelsClientListByPartnerNamespaceResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ChannelsClient.NewListByPartnerNamespacePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByPartnerNamespaceCreateRequest(ctx, resourceGroupName, partnerNamespaceName, options)
			}, nil)
			if err != nil {
				return ChannelsClientListByPartnerNamespaceResponse{}, err
			}
			return client.listByPartnerNamespaceHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByPartnerNamespaceCreateRequest creates the ListByPartnerNamespace request.
func (client *ChannelsClient) listByPartnerNamespaceCreateRequest(ctx context.Context, resourceGroupName string, partnerNamespaceName string, options *ChannelsClientListByPartnerNamespaceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/partnerNamespaces/{partnerNamespaceName}/channels"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if partnerNamespaceName == "" {
		return nil, errors.New("parameter partnerNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerNamespaceName}", url.PathEscape(partnerNamespaceName))
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
	reqQP.Set("api-version", "2025-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByPartnerNamespaceHandleResponse handles the ListByPartnerNamespace response.
func (client *ChannelsClient) listByPartnerNamespaceHandleResponse(resp *http.Response) (ChannelsClientListByPartnerNamespaceResponse, error) {
	result := ChannelsClientListByPartnerNamespaceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ChannelsListResult); err != nil {
		return ChannelsClientListByPartnerNamespaceResponse{}, err
	}
	return result, nil
}

// Update - Synchronously updates a channel with the specified parameters.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2025-04-01-preview
//   - resourceGroupName - The name of the resource group within the partners subscription.
//   - partnerNamespaceName - Name of the partner namespace.
//   - channelName - Name of the channel.
//   - channelUpdateParameters - Channel update information.
//   - options - ChannelsClientUpdateOptions contains the optional parameters for the ChannelsClient.Update method.
func (client *ChannelsClient) Update(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, channelUpdateParameters ChannelUpdateParameters, options *ChannelsClientUpdateOptions) (ChannelsClientUpdateResponse, error) {
	var err error
	const operationName = "ChannelsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, partnerNamespaceName, channelName, channelUpdateParameters, options)
	if err != nil {
		return ChannelsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ChannelsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ChannelsClientUpdateResponse{}, err
	}
	return ChannelsClientUpdateResponse{}, nil
}

// updateCreateRequest creates the Update request.
func (client *ChannelsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, partnerNamespaceName string, channelName string, channelUpdateParameters ChannelUpdateParameters, _ *ChannelsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/partnerNamespaces/{partnerNamespaceName}/channels/{channelName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if partnerNamespaceName == "" {
		return nil, errors.New("parameter partnerNamespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerNamespaceName}", url.PathEscape(partnerNamespaceName))
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if err := runtime.MarshalAsJSON(req, channelUpdateParameters); err != nil {
		return nil, err
	}
	return req, nil
}

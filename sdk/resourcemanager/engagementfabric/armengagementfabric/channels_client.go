//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armengagementfabric

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

// ChannelsClient contains the methods for the Channels group.
// Don't use this type directly, use NewChannelsClient() instead.
type ChannelsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewChannelsClient creates a new instance of ChannelsClient with the specified values.
//   - subscriptionID - Subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewChannelsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ChannelsClient, error) {
	cl, err := arm.NewClient(moduleName+".ChannelsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ChannelsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or Update the EngagementFabric channel
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-09-01-preview
//   - resourceGroupName - Resource Group Name
//   - accountName - Account Name
//   - channelName - Channel Name
//   - channel - The EngagementFabric channel description
//   - options - ChannelsClientCreateOrUpdateOptions contains the optional parameters for the ChannelsClient.CreateOrUpdate method.
func (client *ChannelsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, channelName string, channel Channel, options *ChannelsClientCreateOrUpdateOptions) (ChannelsClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, accountName, channelName, channel, options)
	if err != nil {
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ChannelsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ChannelsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, channelName string, channel Channel, options *ChannelsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EngagementFabric/Accounts/{accountName}/Channels/{channelName}"
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
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, channel); err != nil {
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

// Delete - Delete the EngagementFabric channel
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-09-01-preview
//   - resourceGroupName - Resource Group Name
//   - accountName - Account Name
//   - channelName - The EngagementFabric channel name
//   - options - ChannelsClientDeleteOptions contains the optional parameters for the ChannelsClient.Delete method.
func (client *ChannelsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, channelName string, options *ChannelsClientDeleteOptions) (ChannelsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, channelName, options)
	if err != nil {
		return ChannelsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ChannelsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return ChannelsClientDeleteResponse{}, err
	}
	return ChannelsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ChannelsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, channelName string, options *ChannelsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EngagementFabric/Accounts/{accountName}/Channels/{channelName}"
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
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get the EngagementFabric channel
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-09-01-preview
//   - resourceGroupName - Resource Group Name
//   - accountName - Account Name
//   - channelName - Channel Name
//   - options - ChannelsClientGetOptions contains the optional parameters for the ChannelsClient.Get method.
func (client *ChannelsClient) Get(ctx context.Context, resourceGroupName string, accountName string, channelName string, options *ChannelsClientGetOptions) (ChannelsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, channelName, options)
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
func (client *ChannelsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, channelName string, options *ChannelsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EngagementFabric/Accounts/{accountName}/Channels/{channelName}"
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
	if channelName == "" {
		return nil, errors.New("parameter channelName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{channelName}", url.PathEscape(channelName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-09-01-preview")
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

// NewListByAccountPager - List the EngagementFabric channels
//
// Generated from API version 2018-09-01-preview
//   - resourceGroupName - Resource Group Name
//   - accountName - Account Name
//   - options - ChannelsClientListByAccountOptions contains the optional parameters for the ChannelsClient.NewListByAccountPager
//     method.
func (client *ChannelsClient) NewListByAccountPager(resourceGroupName string, accountName string, options *ChannelsClientListByAccountOptions) (*runtime.Pager[ChannelsClientListByAccountResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ChannelsClientListByAccountResponse]{
		More: func(page ChannelsClientListByAccountResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ChannelsClientListByAccountResponse) (ChannelsClientListByAccountResponse, error) {
			req, err := client.listByAccountCreateRequest(ctx, resourceGroupName, accountName, options)
			if err != nil {
				return ChannelsClientListByAccountResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ChannelsClientListByAccountResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ChannelsClientListByAccountResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByAccountHandleResponse(resp)
		},
	})
}

// listByAccountCreateRequest creates the ListByAccount request.
func (client *ChannelsClient) listByAccountCreateRequest(ctx context.Context, resourceGroupName string, accountName string, options *ChannelsClientListByAccountOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EngagementFabric/Accounts/{accountName}/Channels"
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
	reqQP.Set("api-version", "2018-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByAccountHandleResponse handles the ListByAccount response.
func (client *ChannelsClient) listByAccountHandleResponse(resp *http.Response) (ChannelsClientListByAccountResponse, error) {
	result := ChannelsClientListByAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ChannelList); err != nil {
		return ChannelsClientListByAccountResponse{}, err
	}
	return result, nil
}


//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// StaticMembersClient contains the methods for the StaticMembers group.
// Don't use this type directly, use NewStaticMembersClient() instead.
type StaticMembersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewStaticMembersClient creates a new instance of StaticMembersClient with the specified values.
//   - subscriptionID - The subscription credentials which uniquely identify the Microsoft Azure subscription. The subscription
//     ID forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewStaticMembersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*StaticMembersClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &StaticMembersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a static member.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - networkManagerName - The name of the network manager.
//   - networkGroupName - The name of the network group.
//   - staticMemberName - The name of the static member.
//   - parameters - Parameters supplied to the specify the static member to create
//   - options - StaticMembersClientCreateOrUpdateOptions contains the optional parameters for the StaticMembersClient.CreateOrUpdate
//     method.
func (client *StaticMembersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string, parameters StaticMember, options *StaticMembersClientCreateOrUpdateOptions) (StaticMembersClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "StaticMembersClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, networkManagerName, networkGroupName, staticMemberName, parameters, options)
	if err != nil {
		return StaticMembersClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StaticMembersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return StaticMembersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *StaticMembersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string, parameters StaticMember, options *StaticMembersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{networkManagerName}/networkGroups/{networkGroupName}/staticMembers/{staticMemberName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkManagerName == "" {
		return nil, errors.New("parameter networkManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkManagerName}", url.PathEscape(networkManagerName))
	if networkGroupName == "" {
		return nil, errors.New("parameter networkGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkGroupName}", url.PathEscape(networkGroupName))
	if staticMemberName == "" {
		return nil, errors.New("parameter staticMemberName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{staticMemberName}", url.PathEscape(staticMemberName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *StaticMembersClient) createOrUpdateHandleResponse(resp *http.Response) (StaticMembersClientCreateOrUpdateResponse, error) {
	result := StaticMembersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StaticMember); err != nil {
		return StaticMembersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a static member.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - networkManagerName - The name of the network manager.
//   - networkGroupName - The name of the network group.
//   - staticMemberName - The name of the static member.
//   - options - StaticMembersClientDeleteOptions contains the optional parameters for the StaticMembersClient.Delete method.
func (client *StaticMembersClient) Delete(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string, options *StaticMembersClientDeleteOptions) (StaticMembersClientDeleteResponse, error) {
	var err error
	const operationName = "StaticMembersClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, networkManagerName, networkGroupName, staticMemberName, options)
	if err != nil {
		return StaticMembersClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StaticMembersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return StaticMembersClientDeleteResponse{}, err
	}
	return StaticMembersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *StaticMembersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string, options *StaticMembersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{networkManagerName}/networkGroups/{networkGroupName}/staticMembers/{staticMemberName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkManagerName == "" {
		return nil, errors.New("parameter networkManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkManagerName}", url.PathEscape(networkManagerName))
	if networkGroupName == "" {
		return nil, errors.New("parameter networkGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkGroupName}", url.PathEscape(networkGroupName))
	if staticMemberName == "" {
		return nil, errors.New("parameter staticMemberName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{staticMemberName}", url.PathEscape(staticMemberName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the specified static member.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - networkManagerName - The name of the network manager.
//   - networkGroupName - The name of the network group.
//   - staticMemberName - The name of the static member.
//   - options - StaticMembersClientGetOptions contains the optional parameters for the StaticMembersClient.Get method.
func (client *StaticMembersClient) Get(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string, options *StaticMembersClientGetOptions) (StaticMembersClientGetResponse, error) {
	var err error
	const operationName = "StaticMembersClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, networkManagerName, networkGroupName, staticMemberName, options)
	if err != nil {
		return StaticMembersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return StaticMembersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return StaticMembersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *StaticMembersClient) getCreateRequest(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string, options *StaticMembersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{networkManagerName}/networkGroups/{networkGroupName}/staticMembers/{staticMemberName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkManagerName == "" {
		return nil, errors.New("parameter networkManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkManagerName}", url.PathEscape(networkManagerName))
	if networkGroupName == "" {
		return nil, errors.New("parameter networkGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkGroupName}", url.PathEscape(networkGroupName))
	if staticMemberName == "" {
		return nil, errors.New("parameter staticMemberName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{staticMemberName}", url.PathEscape(staticMemberName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *StaticMembersClient) getHandleResponse(resp *http.Response) (StaticMembersClientGetResponse, error) {
	result := StaticMembersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StaticMember); err != nil {
		return StaticMembersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists the specified static member.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group.
//   - networkManagerName - The name of the network manager.
//   - networkGroupName - The name of the network group.
//   - options - StaticMembersClientListOptions contains the optional parameters for the StaticMembersClient.NewListPager method.
func (client *StaticMembersClient) NewListPager(resourceGroupName string, networkManagerName string, networkGroupName string, options *StaticMembersClientListOptions) *runtime.Pager[StaticMembersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[StaticMembersClientListResponse]{
		More: func(page StaticMembersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *StaticMembersClientListResponse) (StaticMembersClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "StaticMembersClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, networkManagerName, networkGroupName, options)
			}, nil)
			if err != nil {
				return StaticMembersClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *StaticMembersClient) listCreateRequest(ctx context.Context, resourceGroupName string, networkManagerName string, networkGroupName string, options *StaticMembersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/networkManagers/{networkManagerName}/networkGroups/{networkGroupName}/staticMembers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if networkManagerName == "" {
		return nil, errors.New("parameter networkManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkManagerName}", url.PathEscape(networkManagerName))
	if networkGroupName == "" {
		return nil, errors.New("parameter networkGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{networkGroupName}", url.PathEscape(networkGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *StaticMembersClient) listHandleResponse(resp *http.Response) (StaticMembersClientListResponse, error) {
	result := StaticMembersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.StaticMemberListResult); err != nil {
		return StaticMembersClientListResponse{}, err
	}
	return result, nil
}

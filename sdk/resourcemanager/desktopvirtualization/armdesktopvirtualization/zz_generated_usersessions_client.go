//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdesktopvirtualization

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// UserSessionsClient contains the methods for the UserSessions group.
// Don't use this type directly, use NewUserSessionsClient() instead.
type UserSessionsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewUserSessionsClient creates a new instance of UserSessionsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewUserSessionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*UserSessionsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &UserSessionsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Delete - Remove a userSession.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-10-preview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// hostPoolName - The name of the host pool within the specified resource group
// sessionHostName - The name of the session host within the specified host pool
// userSessionID - The name of the user session within the specified session host
// options - UserSessionsClientDeleteOptions contains the optional parameters for the UserSessionsClient.Delete method.
func (client *UserSessionsClient) Delete(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientDeleteOptions) (UserSessionsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, userSessionID, options)
	if err != nil {
		return UserSessionsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserSessionsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return UserSessionsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return UserSessionsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *UserSessionsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}/userSessions/{userSessionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if hostPoolName == "" {
		return nil, errors.New("parameter hostPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{hostPoolName}", url.PathEscape(hostPoolName))
	if sessionHostName == "" {
		return nil, errors.New("parameter sessionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sessionHostName}", url.PathEscape(sessionHostName))
	if userSessionID == "" {
		return nil, errors.New("parameter userSessionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{userSessionId}", url.PathEscape(userSessionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-10-preview")
	if options != nil && options.Force != nil {
		reqQP.Set("force", strconv.FormatBool(*options.Force))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Disconnect - Disconnect a userSession.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-10-preview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// hostPoolName - The name of the host pool within the specified resource group
// sessionHostName - The name of the session host within the specified host pool
// userSessionID - The name of the user session within the specified session host
// options - UserSessionsClientDisconnectOptions contains the optional parameters for the UserSessionsClient.Disconnect method.
func (client *UserSessionsClient) Disconnect(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientDisconnectOptions) (UserSessionsClientDisconnectResponse, error) {
	req, err := client.disconnectCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, userSessionID, options)
	if err != nil {
		return UserSessionsClientDisconnectResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserSessionsClientDisconnectResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserSessionsClientDisconnectResponse{}, runtime.NewResponseError(resp)
	}
	return UserSessionsClientDisconnectResponse{}, nil
}

// disconnectCreateRequest creates the Disconnect request.
func (client *UserSessionsClient) disconnectCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientDisconnectOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}/userSessions/{userSessionId}/disconnect"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if hostPoolName == "" {
		return nil, errors.New("parameter hostPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{hostPoolName}", url.PathEscape(hostPoolName))
	if sessionHostName == "" {
		return nil, errors.New("parameter sessionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sessionHostName}", url.PathEscape(sessionHostName))
	if userSessionID == "" {
		return nil, errors.New("parameter userSessionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{userSessionId}", url.PathEscape(userSessionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a userSession.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-10-preview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// hostPoolName - The name of the host pool within the specified resource group
// sessionHostName - The name of the session host within the specified host pool
// userSessionID - The name of the user session within the specified session host
// options - UserSessionsClientGetOptions contains the optional parameters for the UserSessionsClient.Get method.
func (client *UserSessionsClient) Get(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientGetOptions) (UserSessionsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, userSessionID, options)
	if err != nil {
		return UserSessionsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserSessionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserSessionsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *UserSessionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}/userSessions/{userSessionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if hostPoolName == "" {
		return nil, errors.New("parameter hostPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{hostPoolName}", url.PathEscape(hostPoolName))
	if sessionHostName == "" {
		return nil, errors.New("parameter sessionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sessionHostName}", url.PathEscape(sessionHostName))
	if userSessionID == "" {
		return nil, errors.New("parameter userSessionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{userSessionId}", url.PathEscape(userSessionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *UserSessionsClient) getHandleResponse(resp *http.Response) (UserSessionsClientGetResponse, error) {
	result := UserSessionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserSession); err != nil {
		return UserSessionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List userSessions.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-10-preview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// hostPoolName - The name of the host pool within the specified resource group
// sessionHostName - The name of the session host within the specified host pool
// options - UserSessionsClientListOptions contains the optional parameters for the UserSessionsClient.List method.
func (client *UserSessionsClient) NewListPager(resourceGroupName string, hostPoolName string, sessionHostName string, options *UserSessionsClientListOptions) *runtime.Pager[UserSessionsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[UserSessionsClientListResponse]{
		More: func(page UserSessionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *UserSessionsClientListResponse) (UserSessionsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return UserSessionsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return UserSessionsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return UserSessionsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *UserSessionsClient) listCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *UserSessionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}/userSessions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if hostPoolName == "" {
		return nil, errors.New("parameter hostPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{hostPoolName}", url.PathEscape(hostPoolName))
	if sessionHostName == "" {
		return nil, errors.New("parameter sessionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sessionHostName}", url.PathEscape(sessionHostName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *UserSessionsClient) listHandleResponse(resp *http.Response) (UserSessionsClientListResponse, error) {
	result := UserSessionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserSessionList); err != nil {
		return UserSessionsClientListResponse{}, err
	}
	return result, nil
}

// NewListByHostPoolPager - List userSessions.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-10-preview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// hostPoolName - The name of the host pool within the specified resource group
// options - UserSessionsClientListByHostPoolOptions contains the optional parameters for the UserSessionsClient.ListByHostPool
// method.
func (client *UserSessionsClient) NewListByHostPoolPager(resourceGroupName string, hostPoolName string, options *UserSessionsClientListByHostPoolOptions) *runtime.Pager[UserSessionsClientListByHostPoolResponse] {
	return runtime.NewPager(runtime.PagingHandler[UserSessionsClientListByHostPoolResponse]{
		More: func(page UserSessionsClientListByHostPoolResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *UserSessionsClientListByHostPoolResponse) (UserSessionsClientListByHostPoolResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByHostPoolCreateRequest(ctx, resourceGroupName, hostPoolName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return UserSessionsClientListByHostPoolResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return UserSessionsClientListByHostPoolResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return UserSessionsClientListByHostPoolResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByHostPoolHandleResponse(resp)
		},
	})
}

// listByHostPoolCreateRequest creates the ListByHostPool request.
func (client *UserSessionsClient) listByHostPoolCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, options *UserSessionsClientListByHostPoolOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/userSessions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if hostPoolName == "" {
		return nil, errors.New("parameter hostPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{hostPoolName}", url.PathEscape(hostPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-10-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByHostPoolHandleResponse handles the ListByHostPool response.
func (client *UserSessionsClient) listByHostPoolHandleResponse(resp *http.Response) (UserSessionsClientListByHostPoolResponse, error) {
	result := UserSessionsClientListByHostPoolResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserSessionList); err != nil {
		return UserSessionsClientListByHostPoolResponse{}, err
	}
	return result, nil
}

// SendMessage - Send a message to a user.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-10-preview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// hostPoolName - The name of the host pool within the specified resource group
// sessionHostName - The name of the session host within the specified host pool
// userSessionID - The name of the user session within the specified session host
// options - UserSessionsClientSendMessageOptions contains the optional parameters for the UserSessionsClient.SendMessage
// method.
func (client *UserSessionsClient) SendMessage(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientSendMessageOptions) (UserSessionsClientSendMessageResponse, error) {
	req, err := client.sendMessageCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, userSessionID, options)
	if err != nil {
		return UserSessionsClientSendMessageResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return UserSessionsClientSendMessageResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserSessionsClientSendMessageResponse{}, runtime.NewResponseError(resp)
	}
	return UserSessionsClientSendMessageResponse{}, nil
}

// sendMessageCreateRequest creates the SendMessage request.
func (client *UserSessionsClient) sendMessageCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, userSessionID string, options *UserSessionsClientSendMessageOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}/userSessions/{userSessionId}/sendMessage"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if hostPoolName == "" {
		return nil, errors.New("parameter hostPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{hostPoolName}", url.PathEscape(hostPoolName))
	if sessionHostName == "" {
		return nil, errors.New("parameter sessionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sessionHostName}", url.PathEscape(sessionHostName))
	if userSessionID == "" {
		return nil, errors.New("parameter userSessionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{userSessionId}", url.PathEscape(userSessionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-10-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.SendMessage != nil {
		return req, runtime.MarshalAsJSON(req, *options.SendMessage)
	}
	return req, nil
}

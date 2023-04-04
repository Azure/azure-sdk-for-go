//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdesktopvirtualization

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

// SessionHostsClient contains the methods for the SessionHosts group.
// Don't use this type directly, use NewSessionHostsClient() instead.
type SessionHostsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSessionHostsClient creates a new instance of SessionHostsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSessionHostsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SessionHostsClient, error) {
	cl, err := arm.NewClient(moduleName+".SessionHostsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SessionHostsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Delete - Remove a SessionHost.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-09
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - hostPoolName - The name of the host pool within the specified resource group
//   - sessionHostName - The name of the session host within the specified host pool
//   - options - SessionHostsClientDeleteOptions contains the optional parameters for the SessionHostsClient.Delete method.
func (client *SessionHostsClient) Delete(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *SessionHostsClientDeleteOptions) (SessionHostsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, options)
	if err != nil {
		return SessionHostsClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SessionHostsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return SessionHostsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return SessionHostsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SessionHostsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *SessionHostsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-09")
	if options != nil && options.Force != nil {
		reqQP.Set("force", strconv.FormatBool(*options.Force))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a session host.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-09
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - hostPoolName - The name of the host pool within the specified resource group
//   - sessionHostName - The name of the session host within the specified host pool
//   - options - SessionHostsClientGetOptions contains the optional parameters for the SessionHostsClient.Get method.
func (client *SessionHostsClient) Get(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *SessionHostsClientGetOptions) (SessionHostsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, options)
	if err != nil {
		return SessionHostsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SessionHostsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SessionHostsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SessionHostsClient) getCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *SessionHostsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-09")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SessionHostsClient) getHandleResponse(resp *http.Response) (SessionHostsClientGetResponse, error) {
	result := SessionHostsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SessionHost); err != nil {
		return SessionHostsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List sessionHosts.
//
// Generated from API version 2022-09-09
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - hostPoolName - The name of the host pool within the specified resource group
//   - options - SessionHostsClientListOptions contains the optional parameters for the SessionHostsClient.NewListPager method.
func (client *SessionHostsClient) NewListPager(resourceGroupName string, hostPoolName string, options *SessionHostsClientListOptions) *runtime.Pager[SessionHostsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SessionHostsClientListResponse]{
		More: func(page SessionHostsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SessionHostsClientListResponse) (SessionHostsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, hostPoolName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SessionHostsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return SessionHostsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SessionHostsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *SessionHostsClient) listCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, options *SessionHostsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts"
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-09")
	if options != nil && options.PageSize != nil {
		reqQP.Set("pageSize", strconv.FormatInt(int64(*options.PageSize), 10))
	}
	if options != nil && options.IsDescending != nil {
		reqQP.Set("isDescending", strconv.FormatBool(*options.IsDescending))
	}
	if options != nil && options.InitialSkip != nil {
		reqQP.Set("initialSkip", strconv.FormatInt(int64(*options.InitialSkip), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SessionHostsClient) listHandleResponse(resp *http.Response) (SessionHostsClientListResponse, error) {
	result := SessionHostsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SessionHostList); err != nil {
		return SessionHostsClientListResponse{}, err
	}
	return result, nil
}

// Update - Update a session host.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-09-09
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - hostPoolName - The name of the host pool within the specified resource group
//   - sessionHostName - The name of the session host within the specified host pool
//   - options - SessionHostsClientUpdateOptions contains the optional parameters for the SessionHostsClient.Update method.
func (client *SessionHostsClient) Update(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *SessionHostsClientUpdateOptions) (SessionHostsClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, hostPoolName, sessionHostName, options)
	if err != nil {
		return SessionHostsClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SessionHostsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SessionHostsClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *SessionHostsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, hostPoolName string, sessionHostName string, options *SessionHostsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DesktopVirtualization/hostPools/{hostPoolName}/sessionHosts/{sessionHostName}"
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
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-09-09")
	if options != nil && options.Force != nil {
		reqQP.Set("force", strconv.FormatBool(*options.Force))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.SessionHost != nil {
		return req, runtime.MarshalAsJSON(req, *options.SessionHost)
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *SessionHostsClient) updateHandleResponse(resp *http.Response) (SessionHostsClientUpdateResponse, error) {
	result := SessionHostsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SessionHost); err != nil {
		return SessionHostsClientUpdateResponse{}, err
	}
	return result, nil
}

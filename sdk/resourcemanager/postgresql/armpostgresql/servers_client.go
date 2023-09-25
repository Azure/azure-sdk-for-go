//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpostgresql

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

// ServersClient contains the methods for the Servers group.
// Don't use this type directly, use NewServersClient() instead.
type ServersClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewServersClient creates a new instance of ServersClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewServersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ServersClient, error) {
	cl, err := arm.NewClient(moduleName+".ServersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ServersClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreate - Creates a new server, or will overwrite an existing server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serverName - The name of the server.
//   - parameters - The required parameters for creating or updating a server.
//   - options - ServersClientBeginCreateOptions contains the optional parameters for the ServersClient.BeginCreate method.
func (client *ServersClient) BeginCreate(ctx context.Context, resourceGroupName string, serverName string, parameters ServerForCreate, options *ServersClientBeginCreateOptions) (*runtime.Poller[ServersClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, serverName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ServersClientCreateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ServersClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Creates a new server, or will overwrite an existing server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
func (client *ServersClient) create(ctx context.Context, resourceGroupName string, serverName string, parameters ServerForCreate, options *ServersClientBeginCreateOptions) (*http.Response, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceGroupName, serverName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createCreateRequest creates the Create request.
func (client *ServersClient) createCreateRequest(ctx context.Context, resourceGroupName string, serverName string, parameters ServerForCreate, options *ServersClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes a server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serverName - The name of the server.
//   - options - ServersClientBeginDeleteOptions contains the optional parameters for the ServersClient.BeginDelete method.
func (client *ServersClient) BeginDelete(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientBeginDeleteOptions) (*runtime.Poller[ServersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, serverName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ServersClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ServersClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes a server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
func (client *ServersClient) deleteOperation(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serverName, options)
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
func (client *ServersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets information about a server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serverName - The name of the server.
//   - options - ServersClientGetOptions contains the optional parameters for the ServersClient.Get method.
func (client *ServersClient) Get(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientGetOptions) (ServersClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverName, options)
	if err != nil {
		return ServersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ServersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ServersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ServersClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ServersClient) getHandleResponse(resp *http.Response) (ServersClientGetResponse, error) {
	result := ServersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Server); err != nil {
		return ServersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List all the servers in a given subscription.
//
// Generated from API version 2017-12-01
//   - options - ServersClientListOptions contains the optional parameters for the ServersClient.NewListPager method.
func (client *ServersClient) NewListPager(options *ServersClientListOptions) (*runtime.Pager[ServersClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ServersClientListResponse]{
		More: func(page ServersClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ServersClientListResponse) (ServersClientListResponse, error) {
			req, err := client.listCreateRequest(ctx, options)
			if err != nil {
				return ServersClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ServersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ServersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ServersClient) listCreateRequest(ctx context.Context, options *ServersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.DBforPostgreSQL/servers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ServersClient) listHandleResponse(resp *http.Response) (ServersClientListResponse, error) {
	result := ServersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ServerListResult); err != nil {
		return ServersClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List all the servers in a given resource group.
//
// Generated from API version 2017-12-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - options - ServersClientListByResourceGroupOptions contains the optional parameters for the ServersClient.NewListByResourceGroupPager
//     method.
func (client *ServersClient) NewListByResourceGroupPager(resourceGroupName string, options *ServersClientListByResourceGroupOptions) (*runtime.Pager[ServersClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ServersClientListByResourceGroupResponse]{
		More: func(page ServersClientListByResourceGroupResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ServersClientListByResourceGroupResponse) (ServersClientListByResourceGroupResponse, error) {
			req, err := client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			if err != nil {
				return ServersClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ServersClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ServersClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ServersClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ServersClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ServersClient) listByResourceGroupHandleResponse(resp *http.Response) (ServersClientListByResourceGroupResponse, error) {
	result := ServersClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ServerListResult); err != nil {
		return ServersClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// BeginRestart - Restarts a server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serverName - The name of the server.
//   - options - ServersClientBeginRestartOptions contains the optional parameters for the ServersClient.BeginRestart method.
func (client *ServersClient) BeginRestart(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientBeginRestartOptions) (*runtime.Poller[ServersClientRestartResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.restart(ctx, resourceGroupName, serverName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ServersClientRestartResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ServersClientRestartResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Restart - Restarts a server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
func (client *ServersClient) restart(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientBeginRestartOptions) (*http.Response, error) {
	var err error
	req, err := client.restartCreateRequest(ctx, resourceGroupName, serverName, options)
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

// restartCreateRequest creates the Restart request.
func (client *ServersClient) restartCreateRequest(ctx context.Context, resourceGroupName string, serverName string, options *ServersClientBeginRestartOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}/restart"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginUpdate - Updates an existing server. The request body can contain one to many of the properties present in the normal
// server definition.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serverName - The name of the server.
//   - parameters - The required parameters for updating a server.
//   - options - ServersClientBeginUpdateOptions contains the optional parameters for the ServersClient.BeginUpdate method.
func (client *ServersClient) BeginUpdate(ctx context.Context, resourceGroupName string, serverName string, parameters ServerUpdateParameters, options *ServersClientBeginUpdateOptions) (*runtime.Poller[ServersClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, serverName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[ServersClientUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[ServersClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Updates an existing server. The request body can contain one to many of the properties present in the normal server
// definition.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2017-12-01
func (client *ServersClient) update(ctx context.Context, resourceGroupName string, serverName string, parameters ServerUpdateParameters, options *ServersClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serverName, parameters, options)
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

// updateCreateRequest creates the Update request.
func (client *ServersClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serverName string, parameters ServerUpdateParameters, options *ServersClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/servers/{serverName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-12-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}


//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armstoragesync

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
	"strings"
)

// RegisteredServersClient contains the methods for the RegisteredServers group.
// Don't use this type directly, use NewRegisteredServersClient() instead.
type RegisteredServersClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewRegisteredServersClient creates a new instance of RegisteredServersClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewRegisteredServersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RegisteredServersClient, error) {
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
	client := &RegisteredServersClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreate - Add a new registered server.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// storageSyncServiceName - Name of Storage Sync Service resource.
// serverID - GUID identifying the on-premises server.
// parameters - Body of Registered Server object.
// options - RegisteredServersClientBeginCreateOptions contains the optional parameters for the RegisteredServersClient.BeginCreate
// method.
func (client *RegisteredServersClient) BeginCreate(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters RegisteredServerCreateParameters, options *RegisteredServersClientBeginCreateOptions) (*runtime.Poller[RegisteredServersClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, storageSyncServiceName, serverID, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[RegisteredServersClientCreateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[RegisteredServersClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Add a new registered server.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
func (client *RegisteredServersClient) create(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters RegisteredServerCreateParameters, options *RegisteredServersClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, storageSyncServiceName, serverID, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *RegisteredServersClient) createCreateRequest(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters RegisteredServerCreateParameters, options *RegisteredServersClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/registeredServers/{serverId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if storageSyncServiceName == "" {
		return nil, errors.New("parameter storageSyncServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageSyncServiceName}", url.PathEscape(storageSyncServiceName))
	if serverID == "" {
		return nil, errors.New("parameter serverID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverId}", url.PathEscape(serverID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Delete the given registered server.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// storageSyncServiceName - Name of Storage Sync Service resource.
// serverID - GUID identifying the on-premises server.
// options - RegisteredServersClientBeginDeleteOptions contains the optional parameters for the RegisteredServersClient.BeginDelete
// method.
func (client *RegisteredServersClient) BeginDelete(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, options *RegisteredServersClientBeginDeleteOptions) (*runtime.Poller[RegisteredServersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, storageSyncServiceName, serverID, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[RegisteredServersClientDeleteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[RegisteredServersClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Delete the given registered server.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
func (client *RegisteredServersClient) deleteOperation(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, options *RegisteredServersClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, storageSyncServiceName, serverID, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *RegisteredServersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, options *RegisteredServersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/registeredServers/{serverId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if storageSyncServiceName == "" {
		return nil, errors.New("parameter storageSyncServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageSyncServiceName}", url.PathEscape(storageSyncServiceName))
	if serverID == "" {
		return nil, errors.New("parameter serverID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverId}", url.PathEscape(serverID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a given registered server.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// storageSyncServiceName - Name of Storage Sync Service resource.
// serverID - GUID identifying the on-premises server.
// options - RegisteredServersClientGetOptions contains the optional parameters for the RegisteredServersClient.Get method.
func (client *RegisteredServersClient) Get(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, options *RegisteredServersClientGetOptions) (RegisteredServersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, storageSyncServiceName, serverID, options)
	if err != nil {
		return RegisteredServersClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RegisteredServersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RegisteredServersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RegisteredServersClient) getCreateRequest(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, options *RegisteredServersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/registeredServers/{serverId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if storageSyncServiceName == "" {
		return nil, errors.New("parameter storageSyncServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageSyncServiceName}", url.PathEscape(storageSyncServiceName))
	if serverID == "" {
		return nil, errors.New("parameter serverID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverId}", url.PathEscape(serverID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RegisteredServersClient) getHandleResponse(resp *http.Response) (RegisteredServersClientGetResponse, error) {
	result := RegisteredServersClientGetResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if val := resp.Header.Get("x-ms-correlation-request-id"); val != "" {
		result.XMSCorrelationRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.RegisteredServer); err != nil {
		return RegisteredServersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByStorageSyncServicePager - Get a given registered server list.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// storageSyncServiceName - Name of Storage Sync Service resource.
// options - RegisteredServersClientListByStorageSyncServiceOptions contains the optional parameters for the RegisteredServersClient.ListByStorageSyncService
// method.
func (client *RegisteredServersClient) NewListByStorageSyncServicePager(resourceGroupName string, storageSyncServiceName string, options *RegisteredServersClientListByStorageSyncServiceOptions) *runtime.Pager[RegisteredServersClientListByStorageSyncServiceResponse] {
	return runtime.NewPager(runtime.PagingHandler[RegisteredServersClientListByStorageSyncServiceResponse]{
		More: func(page RegisteredServersClientListByStorageSyncServiceResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *RegisteredServersClientListByStorageSyncServiceResponse) (RegisteredServersClientListByStorageSyncServiceResponse, error) {
			req, err := client.listByStorageSyncServiceCreateRequest(ctx, resourceGroupName, storageSyncServiceName, options)
			if err != nil {
				return RegisteredServersClientListByStorageSyncServiceResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return RegisteredServersClientListByStorageSyncServiceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RegisteredServersClientListByStorageSyncServiceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByStorageSyncServiceHandleResponse(resp)
		},
	})
}

// listByStorageSyncServiceCreateRequest creates the ListByStorageSyncService request.
func (client *RegisteredServersClient) listByStorageSyncServiceCreateRequest(ctx context.Context, resourceGroupName string, storageSyncServiceName string, options *RegisteredServersClientListByStorageSyncServiceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/registeredServers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if storageSyncServiceName == "" {
		return nil, errors.New("parameter storageSyncServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageSyncServiceName}", url.PathEscape(storageSyncServiceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByStorageSyncServiceHandleResponse handles the ListByStorageSyncService response.
func (client *RegisteredServersClient) listByStorageSyncServiceHandleResponse(resp *http.Response) (RegisteredServersClientListByStorageSyncServiceResponse, error) {
	result := RegisteredServersClientListByStorageSyncServiceResponse{}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.XMSRequestID = &val
	}
	if val := resp.Header.Get("x-ms-correlation-request-id"); val != "" {
		result.XMSCorrelationRequestID = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.RegisteredServerArray); err != nil {
		return RegisteredServersClientListByStorageSyncServiceResponse{}, err
	}
	return result, nil
}

// BeginTriggerRollover - Triggers Server certificate rollover.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// storageSyncServiceName - Name of Storage Sync Service resource.
// serverID - Server Id
// parameters - Body of Trigger Rollover request.
// options - RegisteredServersClientBeginTriggerRolloverOptions contains the optional parameters for the RegisteredServersClient.BeginTriggerRollover
// method.
func (client *RegisteredServersClient) BeginTriggerRollover(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters TriggerRolloverRequest, options *RegisteredServersClientBeginTriggerRolloverOptions) (*runtime.Poller[RegisteredServersClientTriggerRolloverResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.triggerRollover(ctx, resourceGroupName, storageSyncServiceName, serverID, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[RegisteredServersClientTriggerRolloverResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[RegisteredServersClientTriggerRolloverResponse](options.ResumeToken, client.pl, nil)
	}
}

// TriggerRollover - Triggers Server certificate rollover.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
func (client *RegisteredServersClient) triggerRollover(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters TriggerRolloverRequest, options *RegisteredServersClientBeginTriggerRolloverOptions) (*http.Response, error) {
	req, err := client.triggerRolloverCreateRequest(ctx, resourceGroupName, storageSyncServiceName, serverID, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// triggerRolloverCreateRequest creates the TriggerRollover request.
func (client *RegisteredServersClient) triggerRolloverCreateRequest(ctx context.Context, resourceGroupName string, storageSyncServiceName string, serverID string, parameters TriggerRolloverRequest, options *RegisteredServersClientBeginTriggerRolloverOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.StorageSync/storageSyncServices/{storageSyncServiceName}/registeredServers/{serverId}/triggerRollover"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if storageSyncServiceName == "" {
		return nil, errors.New("parameter storageSyncServiceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{storageSyncServiceName}", url.PathEscape(storageSyncServiceName))
	if serverID == "" {
		return nil, errors.New("parameter serverID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverId}", url.PathEscape(serverID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

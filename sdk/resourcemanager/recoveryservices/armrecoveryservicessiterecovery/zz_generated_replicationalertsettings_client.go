//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armrecoveryservicessiterecovery

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

// ReplicationAlertSettingsClient contains the methods for the ReplicationAlertSettings group.
// Don't use this type directly, use NewReplicationAlertSettingsClient() instead.
type ReplicationAlertSettingsClient struct {
	host              string
	resourceName      string
	resourceGroupName string
	subscriptionID    string
	pl                runtime.Pipeline
}

// NewReplicationAlertSettingsClient creates a new instance of ReplicationAlertSettingsClient with the specified values.
// resourceName - The name of the recovery services vault.
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// subscriptionID - The subscription Id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewReplicationAlertSettingsClient(resourceName string, resourceGroupName string, subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ReplicationAlertSettingsClient, error) {
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
	client := &ReplicationAlertSettingsClient{
		resourceName:      resourceName,
		resourceGroupName: resourceGroupName,
		subscriptionID:    subscriptionID,
		host:              ep,
		pl:                pl,
	}
	return client, nil
}

// Create - Create or update an email notification(alert) configuration.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// alertSettingName - The name of the email notification(alert) configuration.
// request - The input to configure the email notification(alert).
// options - ReplicationAlertSettingsClientCreateOptions contains the optional parameters for the ReplicationAlertSettingsClient.Create
// method.
func (client *ReplicationAlertSettingsClient) Create(ctx context.Context, alertSettingName string, request ConfigureAlertRequest, options *ReplicationAlertSettingsClientCreateOptions) (ReplicationAlertSettingsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, alertSettingName, request, options)
	if err != nil {
		return ReplicationAlertSettingsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ReplicationAlertSettingsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReplicationAlertSettingsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *ReplicationAlertSettingsClient) createCreateRequest(ctx context.Context, alertSettingName string, request ConfigureAlertRequest, options *ReplicationAlertSettingsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationAlertSettings/{alertSettingName}"
	if client.resourceName == "" {
		return nil, errors.New("parameter client.resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(client.resourceName))
	if client.resourceGroupName == "" {
		return nil, errors.New("parameter client.resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(client.resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if alertSettingName == "" {
		return nil, errors.New("parameter alertSettingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{alertSettingName}", url.PathEscape(alertSettingName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, request)
}

// createHandleResponse handles the Create response.
func (client *ReplicationAlertSettingsClient) createHandleResponse(resp *http.Response) (ReplicationAlertSettingsClientCreateResponse, error) {
	result := ReplicationAlertSettingsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Alert); err != nil {
		return ReplicationAlertSettingsClientCreateResponse{}, err
	}
	return result, nil
}

// Get - Gets the details of the specified email notification(alert) configuration.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// alertSettingName - The name of the email notification configuration.
// options - ReplicationAlertSettingsClientGetOptions contains the optional parameters for the ReplicationAlertSettingsClient.Get
// method.
func (client *ReplicationAlertSettingsClient) Get(ctx context.Context, alertSettingName string, options *ReplicationAlertSettingsClientGetOptions) (ReplicationAlertSettingsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, alertSettingName, options)
	if err != nil {
		return ReplicationAlertSettingsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ReplicationAlertSettingsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReplicationAlertSettingsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ReplicationAlertSettingsClient) getCreateRequest(ctx context.Context, alertSettingName string, options *ReplicationAlertSettingsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationAlertSettings/{alertSettingName}"
	if client.resourceName == "" {
		return nil, errors.New("parameter client.resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(client.resourceName))
	if client.resourceGroupName == "" {
		return nil, errors.New("parameter client.resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(client.resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if alertSettingName == "" {
		return nil, errors.New("parameter alertSettingName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{alertSettingName}", url.PathEscape(alertSettingName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReplicationAlertSettingsClient) getHandleResponse(resp *http.Response) (ReplicationAlertSettingsClientGetResponse, error) {
	result := ReplicationAlertSettingsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Alert); err != nil {
		return ReplicationAlertSettingsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets the list of email notification(alert) configurations for the vault.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// options - ReplicationAlertSettingsClientListOptions contains the optional parameters for the ReplicationAlertSettingsClient.List
// method.
func (client *ReplicationAlertSettingsClient) NewListPager(options *ReplicationAlertSettingsClientListOptions) *runtime.Pager[ReplicationAlertSettingsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicationAlertSettingsClientListResponse]{
		More: func(page ReplicationAlertSettingsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicationAlertSettingsClientListResponse) (ReplicationAlertSettingsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ReplicationAlertSettingsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ReplicationAlertSettingsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ReplicationAlertSettingsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ReplicationAlertSettingsClient) listCreateRequest(ctx context.Context, options *ReplicationAlertSettingsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationAlertSettings"
	if client.resourceName == "" {
		return nil, errors.New("parameter client.resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(client.resourceName))
	if client.resourceGroupName == "" {
		return nil, errors.New("parameter client.resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(client.resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ReplicationAlertSettingsClient) listHandleResponse(resp *http.Response) (ReplicationAlertSettingsClientListResponse, error) {
	result := ReplicationAlertSettingsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AlertCollection); err != nil {
		return ReplicationAlertSettingsClientListResponse{}, err
	}
	return result, nil
}

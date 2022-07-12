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

// RecoveryPointsClient contains the methods for the RecoveryPoints group.
// Don't use this type directly, use NewRecoveryPointsClient() instead.
type RecoveryPointsClient struct {
	host              string
	resourceName      string
	resourceGroupName string
	subscriptionID    string
	pl                runtime.Pipeline
}

// NewRecoveryPointsClient creates a new instance of RecoveryPointsClient with the specified values.
// resourceName - The name of the recovery services vault.
// resourceGroupName - The name of the resource group where the recovery services vault is present.
// subscriptionID - The subscription Id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewRecoveryPointsClient(resourceName string, resourceGroupName string, subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RecoveryPointsClient, error) {
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
	client := &RecoveryPointsClient{
		resourceName:      resourceName,
		resourceGroupName: resourceGroupName,
		subscriptionID:    subscriptionID,
		host:              ep,
		pl:                pl,
	}
	return client, nil
}

// Get - Get the details of specified recovery point.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// fabricName - The fabric name.
// protectionContainerName - The protection container name.
// replicatedProtectedItemName - The replication protected item name.
// recoveryPointName - The recovery point name.
// options - RecoveryPointsClientGetOptions contains the optional parameters for the RecoveryPointsClient.Get method.
func (client *RecoveryPointsClient) Get(ctx context.Context, fabricName string, protectionContainerName string, replicatedProtectedItemName string, recoveryPointName string, options *RecoveryPointsClientGetOptions) (RecoveryPointsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, fabricName, protectionContainerName, replicatedProtectedItemName, recoveryPointName, options)
	if err != nil {
		return RecoveryPointsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RecoveryPointsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RecoveryPointsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RecoveryPointsClient) getCreateRequest(ctx context.Context, fabricName string, protectionContainerName string, replicatedProtectedItemName string, recoveryPointName string, options *RecoveryPointsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{fabricName}/replicationProtectionContainers/{protectionContainerName}/replicationProtectedItems/{replicatedProtectedItemName}/recoveryPoints/{recoveryPointName}"
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
	if fabricName == "" {
		return nil, errors.New("parameter fabricName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{fabricName}", url.PathEscape(fabricName))
	if protectionContainerName == "" {
		return nil, errors.New("parameter protectionContainerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{protectionContainerName}", url.PathEscape(protectionContainerName))
	if replicatedProtectedItemName == "" {
		return nil, errors.New("parameter replicatedProtectedItemName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{replicatedProtectedItemName}", url.PathEscape(replicatedProtectedItemName))
	if recoveryPointName == "" {
		return nil, errors.New("parameter recoveryPointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{recoveryPointName}", url.PathEscape(recoveryPointName))
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
func (client *RecoveryPointsClient) getHandleResponse(resp *http.Response) (RecoveryPointsClientGetResponse, error) {
	result := RecoveryPointsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RecoveryPoint); err != nil {
		return RecoveryPointsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByReplicationProtectedItemsPager - Lists the available recovery points for a replication protected item.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// fabricName - The fabric name.
// protectionContainerName - The protection container name.
// replicatedProtectedItemName - The replication protected item name.
// options - RecoveryPointsClientListByReplicationProtectedItemsOptions contains the optional parameters for the RecoveryPointsClient.ListByReplicationProtectedItems
// method.
func (client *RecoveryPointsClient) NewListByReplicationProtectedItemsPager(fabricName string, protectionContainerName string, replicatedProtectedItemName string, options *RecoveryPointsClientListByReplicationProtectedItemsOptions) *runtime.Pager[RecoveryPointsClientListByReplicationProtectedItemsResponse] {
	return runtime.NewPager(runtime.PagingHandler[RecoveryPointsClientListByReplicationProtectedItemsResponse]{
		More: func(page RecoveryPointsClientListByReplicationProtectedItemsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RecoveryPointsClientListByReplicationProtectedItemsResponse) (RecoveryPointsClientListByReplicationProtectedItemsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByReplicationProtectedItemsCreateRequest(ctx, fabricName, protectionContainerName, replicatedProtectedItemName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RecoveryPointsClientListByReplicationProtectedItemsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return RecoveryPointsClientListByReplicationProtectedItemsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RecoveryPointsClientListByReplicationProtectedItemsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByReplicationProtectedItemsHandleResponse(resp)
		},
	})
}

// listByReplicationProtectedItemsCreateRequest creates the ListByReplicationProtectedItems request.
func (client *RecoveryPointsClient) listByReplicationProtectedItemsCreateRequest(ctx context.Context, fabricName string, protectionContainerName string, replicatedProtectedItemName string, options *RecoveryPointsClientListByReplicationProtectedItemsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{fabricName}/replicationProtectionContainers/{protectionContainerName}/replicationProtectedItems/{replicatedProtectedItemName}/recoveryPoints"
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
	if fabricName == "" {
		return nil, errors.New("parameter fabricName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{fabricName}", url.PathEscape(fabricName))
	if protectionContainerName == "" {
		return nil, errors.New("parameter protectionContainerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{protectionContainerName}", url.PathEscape(protectionContainerName))
	if replicatedProtectedItemName == "" {
		return nil, errors.New("parameter replicatedProtectedItemName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{replicatedProtectedItemName}", url.PathEscape(replicatedProtectedItemName))
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

// listByReplicationProtectedItemsHandleResponse handles the ListByReplicationProtectedItems response.
func (client *RecoveryPointsClient) listByReplicationProtectedItemsHandleResponse(resp *http.Response) (RecoveryPointsClientListByReplicationProtectedItemsResponse, error) {
	result := RecoveryPointsClientListByReplicationProtectedItemsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RecoveryPointCollection); err != nil {
		return RecoveryPointsClientListByReplicationProtectedItemsResponse{}, err
	}
	return result, nil
}

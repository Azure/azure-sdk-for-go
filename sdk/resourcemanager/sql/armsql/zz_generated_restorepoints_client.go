//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsql

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

// RestorePointsClient contains the methods for the RestorePoints group.
// Don't use this type directly, use NewRestorePointsClient() instead.
type RestorePointsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewRestorePointsClient creates a new instance of RestorePointsClient with the specified values.
// subscriptionID - The subscription ID that identifies an Azure subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewRestorePointsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RestorePointsClient, error) {
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
	client := &RestorePointsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreate - Creates a restore point for a data warehouse.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// parameters - The definition for creating the restore point of this database.
// options - RestorePointsClientBeginCreateOptions contains the optional parameters for the RestorePointsClient.BeginCreate
// method.
func (client *RestorePointsClient) BeginCreate(ctx context.Context, resourceGroupName string, serverName string, databaseName string, parameters CreateDatabaseRestorePointDefinition, options *RestorePointsClientBeginCreateOptions) (*runtime.Poller[RestorePointsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, serverName, databaseName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[RestorePointsClientCreateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[RestorePointsClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Creates a restore point for a data warehouse.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
func (client *RestorePointsClient) create(ctx context.Context, resourceGroupName string, serverName string, databaseName string, parameters CreateDatabaseRestorePointDefinition, options *RestorePointsClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, serverName, databaseName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *RestorePointsClient) createCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, parameters CreateDatabaseRestorePointDefinition, options *RestorePointsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/restorePoints"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// Delete - Deletes a restore point.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// restorePointName - The name of the restore point.
// options - RestorePointsClientDeleteOptions contains the optional parameters for the RestorePointsClient.Delete method.
func (client *RestorePointsClient) Delete(ctx context.Context, resourceGroupName string, serverName string, databaseName string, restorePointName string, options *RestorePointsClientDeleteOptions) (RestorePointsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serverName, databaseName, restorePointName, options)
	if err != nil {
		return RestorePointsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RestorePointsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RestorePointsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return RestorePointsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *RestorePointsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, restorePointName string, options *RestorePointsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/restorePoints/{restorePointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if restorePointName == "" {
		return nil, errors.New("parameter restorePointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointName}", url.PathEscape(restorePointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Gets a restore point.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// restorePointName - The name of the restore point.
// options - RestorePointsClientGetOptions contains the optional parameters for the RestorePointsClient.Get method.
func (client *RestorePointsClient) Get(ctx context.Context, resourceGroupName string, serverName string, databaseName string, restorePointName string, options *RestorePointsClientGetOptions) (RestorePointsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverName, databaseName, restorePointName, options)
	if err != nil {
		return RestorePointsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RestorePointsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RestorePointsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RestorePointsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, restorePointName string, options *RestorePointsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/restorePoints/{restorePointName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if restorePointName == "" {
		return nil, errors.New("parameter restorePointName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{restorePointName}", url.PathEscape(restorePointName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RestorePointsClient) getHandleResponse(resp *http.Response) (RestorePointsClientGetResponse, error) {
	result := RestorePointsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RestorePoint); err != nil {
		return RestorePointsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByDatabasePager - Gets a list of database restore points.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// options - RestorePointsClientListByDatabaseOptions contains the optional parameters for the RestorePointsClient.ListByDatabase
// method.
func (client *RestorePointsClient) NewListByDatabasePager(resourceGroupName string, serverName string, databaseName string, options *RestorePointsClientListByDatabaseOptions) *runtime.Pager[RestorePointsClientListByDatabaseResponse] {
	return runtime.NewPager(runtime.PagingHandler[RestorePointsClientListByDatabaseResponse]{
		More: func(page RestorePointsClientListByDatabaseResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RestorePointsClientListByDatabaseResponse) (RestorePointsClientListByDatabaseResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByDatabaseCreateRequest(ctx, resourceGroupName, serverName, databaseName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RestorePointsClientListByDatabaseResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return RestorePointsClientListByDatabaseResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RestorePointsClientListByDatabaseResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByDatabaseHandleResponse(resp)
		},
	})
}

// listByDatabaseCreateRequest creates the ListByDatabase request.
func (client *RestorePointsClient) listByDatabaseCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, options *RestorePointsClientListByDatabaseOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/restorePoints"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByDatabaseHandleResponse handles the ListByDatabase response.
func (client *RestorePointsClient) listByDatabaseHandleResponse(resp *http.Response) (RestorePointsClientListByDatabaseResponse, error) {
	result := RestorePointsClientListByDatabaseResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RestorePointListResult); err != nil {
		return RestorePointsClientListByDatabaseResponse{}, err
	}
	return result, nil
}

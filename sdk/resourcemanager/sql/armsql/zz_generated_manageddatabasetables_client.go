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

// ManagedDatabaseTablesClient contains the methods for the ManagedDatabaseTables group.
// Don't use this type directly, use NewManagedDatabaseTablesClient() instead.
type ManagedDatabaseTablesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewManagedDatabaseTablesClient creates a new instance of ManagedDatabaseTablesClient with the specified values.
// subscriptionID - The subscription ID that identifies an Azure subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewManagedDatabaseTablesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ManagedDatabaseTablesClient, error) {
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
	client := &ManagedDatabaseTablesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Get managed database table
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// managedInstanceName - The name of the managed instance.
// databaseName - The name of the database.
// schemaName - The name of the schema.
// tableName - The name of the table.
// options - ManagedDatabaseTablesClientGetOptions contains the optional parameters for the ManagedDatabaseTablesClient.Get
// method.
func (client *ManagedDatabaseTablesClient) Get(ctx context.Context, resourceGroupName string, managedInstanceName string, databaseName string, schemaName string, tableName string, options *ManagedDatabaseTablesClientGetOptions) (ManagedDatabaseTablesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, managedInstanceName, databaseName, schemaName, tableName, options)
	if err != nil {
		return ManagedDatabaseTablesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ManagedDatabaseTablesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ManagedDatabaseTablesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ManagedDatabaseTablesClient) getCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, databaseName string, schemaName string, tableName string, options *ManagedDatabaseTablesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/databases/{databaseName}/schemas/{schemaName}/tables/{tableName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if schemaName == "" {
		return nil, errors.New("parameter schemaName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{schemaName}", url.PathEscape(schemaName))
	if tableName == "" {
		return nil, errors.New("parameter tableName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tableName}", url.PathEscape(tableName))
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
func (client *ManagedDatabaseTablesClient) getHandleResponse(resp *http.Response) (ManagedDatabaseTablesClientGetResponse, error) {
	result := ManagedDatabaseTablesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DatabaseTable); err != nil {
		return ManagedDatabaseTablesClientGetResponse{}, err
	}
	return result, nil
}

// NewListBySchemaPager - List managed database tables
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// managedInstanceName - The name of the managed instance.
// databaseName - The name of the database.
// schemaName - The name of the schema.
// options - ManagedDatabaseTablesClientListBySchemaOptions contains the optional parameters for the ManagedDatabaseTablesClient.ListBySchema
// method.
func (client *ManagedDatabaseTablesClient) NewListBySchemaPager(resourceGroupName string, managedInstanceName string, databaseName string, schemaName string, options *ManagedDatabaseTablesClientListBySchemaOptions) *runtime.Pager[ManagedDatabaseTablesClientListBySchemaResponse] {
	return runtime.NewPager(runtime.PagingHandler[ManagedDatabaseTablesClientListBySchemaResponse]{
		More: func(page ManagedDatabaseTablesClientListBySchemaResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ManagedDatabaseTablesClientListBySchemaResponse) (ManagedDatabaseTablesClientListBySchemaResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySchemaCreateRequest(ctx, resourceGroupName, managedInstanceName, databaseName, schemaName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ManagedDatabaseTablesClientListBySchemaResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ManagedDatabaseTablesClientListBySchemaResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ManagedDatabaseTablesClientListBySchemaResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySchemaHandleResponse(resp)
		},
	})
}

// listBySchemaCreateRequest creates the ListBySchema request.
func (client *ManagedDatabaseTablesClient) listBySchemaCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, databaseName string, schemaName string, options *ManagedDatabaseTablesClientListBySchemaOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/databases/{databaseName}/schemas/{schemaName}/tables"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if schemaName == "" {
		return nil, errors.New("parameter schemaName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{schemaName}", url.PathEscape(schemaName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySchemaHandleResponse handles the ListBySchema response.
func (client *ManagedDatabaseTablesClient) listBySchemaHandleResponse(resp *http.Response) (ManagedDatabaseTablesClientListBySchemaResponse, error) {
	result := ManagedDatabaseTablesClientListBySchemaResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DatabaseTableListResult); err != nil {
		return ManagedDatabaseTablesClientListBySchemaResponse{}, err
	}
	return result, nil
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsynapse

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

// SQLPoolTableColumnsClient contains the methods for the SQLPoolTableColumns group.
// Don't use this type directly, use NewSQLPoolTableColumnsClient() instead.
type SQLPoolTableColumnsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewSQLPoolTableColumnsClient creates a new instance of SQLPoolTableColumnsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSQLPoolTableColumnsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLPoolTableColumnsClient, error) {
	cl, err := arm.NewClient(moduleName+".SQLPoolTableColumnsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SQLPoolTableColumnsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// NewListByTableNamePager - Gets columns in a given table in a SQL pool.
//
// Generated from API version 2021-06-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - workspaceName - The name of the workspace.
//   - sqlPoolName - SQL pool name
//   - schemaName - The name of the schema.
//   - tableName - The name of the table.
//   - options - SQLPoolTableColumnsClientListByTableNameOptions contains the optional parameters for the SQLPoolTableColumnsClient.NewListByTableNamePager
//     method.
func (client *SQLPoolTableColumnsClient) NewListByTableNamePager(resourceGroupName string, workspaceName string, sqlPoolName string, schemaName string, tableName string, options *SQLPoolTableColumnsClientListByTableNameOptions) (*runtime.Pager[SQLPoolTableColumnsClientListByTableNameResponse]) {
	return runtime.NewPager(runtime.PagingHandler[SQLPoolTableColumnsClientListByTableNameResponse]{
		More: func(page SQLPoolTableColumnsClientListByTableNameResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SQLPoolTableColumnsClientListByTableNameResponse) (SQLPoolTableColumnsClientListByTableNameResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByTableNameCreateRequest(ctx, resourceGroupName, workspaceName, sqlPoolName, schemaName, tableName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SQLPoolTableColumnsClientListByTableNameResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return SQLPoolTableColumnsClientListByTableNameResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SQLPoolTableColumnsClientListByTableNameResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByTableNameHandleResponse(resp)
		},
	})
}

// listByTableNameCreateRequest creates the ListByTableName request.
func (client *SQLPoolTableColumnsClient) listByTableNameCreateRequest(ctx context.Context, resourceGroupName string, workspaceName string, sqlPoolName string, schemaName string, tableName string, options *SQLPoolTableColumnsClientListByTableNameOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Synapse/workspaces/{workspaceName}/sqlPools/{sqlPoolName}/schemas/{schemaName}/tables/{tableName}/columns"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workspaceName == "" {
		return nil, errors.New("parameter workspaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workspaceName}", url.PathEscape(workspaceName))
	if sqlPoolName == "" {
		return nil, errors.New("parameter sqlPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlPoolName}", url.PathEscape(sqlPoolName))
	if schemaName == "" {
		return nil, errors.New("parameter schemaName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{schemaName}", url.PathEscape(schemaName))
	if tableName == "" {
		return nil, errors.New("parameter tableName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{tableName}", url.PathEscape(tableName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByTableNameHandleResponse handles the ListByTableName response.
func (client *SQLPoolTableColumnsClient) listByTableNameHandleResponse(resp *http.Response) (SQLPoolTableColumnsClientListByTableNameResponse, error) {
	result := SQLPoolTableColumnsClientListByTableNameResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SQLPoolColumnListResult); err != nil {
		return SQLPoolTableColumnsClientListByTableNameResponse{}, err
	}
	return result, nil
}


//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armazuredata

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

// SQLServersClient contains the methods for the SQLServers group.
// Don't use this type directly, use NewSQLServersClient() instead.
type SQLServersClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSQLServersClient creates a new instance of SQLServersClient with the specified values.
//   - subscriptionID - Subscription ID that identifies an Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSQLServersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLServersClient, error) {
	cl, err := arm.NewClient(moduleName+".SQLServersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SQLServersClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a SQL Server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-07-24-preview
//   - resourceGroupName - Name of the resource group that contains the resource. You can obtain this value from the Azure Resource
//     Manager API or the portal.
//   - sqlServerRegistrationName - Name of the SQL Server registration.
//   - sqlServerName - Name of the SQL Server.
//   - parameters - The SQL Server to be created or updated.
//   - options - SQLServersClientCreateOrUpdateOptions contains the optional parameters for the SQLServersClient.CreateOrUpdate
//     method.
func (client *SQLServersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, sqlServerName string, parameters SQLServer, options *SQLServersClientCreateOrUpdateOptions) (SQLServersClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, sqlServerRegistrationName, sqlServerName, parameters, options)
	if err != nil {
		return SQLServersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SQLServersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return SQLServersClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *SQLServersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, sqlServerName string, parameters SQLServer, options *SQLServersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureData/sqlServerRegistrations/{sqlServerRegistrationName}/sqlServers/{sqlServerName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sqlServerRegistrationName == "" {
		return nil, errors.New("parameter sqlServerRegistrationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerRegistrationName}", url.PathEscape(sqlServerRegistrationName))
	if sqlServerName == "" {
		return nil, errors.New("parameter sqlServerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerName}", url.PathEscape(sqlServerName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-07-24-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *SQLServersClient) createOrUpdateHandleResponse(resp *http.Response) (SQLServersClientCreateOrUpdateResponse, error) {
	result := SQLServersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SQLServer); err != nil {
		return SQLServersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a SQL Server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-07-24-preview
//   - resourceGroupName - Name of the resource group that contains the resource. You can obtain this value from the Azure Resource
//     Manager API or the portal.
//   - sqlServerRegistrationName - Name of the SQL Server registration.
//   - sqlServerName - Name of the SQL Server.
//   - options - SQLServersClientDeleteOptions contains the optional parameters for the SQLServersClient.Delete method.
func (client *SQLServersClient) Delete(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, sqlServerName string, options *SQLServersClientDeleteOptions) (SQLServersClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, sqlServerRegistrationName, sqlServerName, options)
	if err != nil {
		return SQLServersClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SQLServersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return SQLServersClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return SQLServersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SQLServersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, sqlServerName string, options *SQLServersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureData/sqlServerRegistrations/{sqlServerRegistrationName}/sqlServers/{sqlServerName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sqlServerRegistrationName == "" {
		return nil, errors.New("parameter sqlServerRegistrationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerRegistrationName}", url.PathEscape(sqlServerRegistrationName))
	if sqlServerName == "" {
		return nil, errors.New("parameter sqlServerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerName}", url.PathEscape(sqlServerName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-07-24-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a SQL Server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-07-24-preview
//   - resourceGroupName - Name of the resource group that contains the resource. You can obtain this value from the Azure Resource
//     Manager API or the portal.
//   - sqlServerRegistrationName - Name of the SQL Server registration.
//   - sqlServerName - Name of the SQL Server.
//   - options - SQLServersClientGetOptions contains the optional parameters for the SQLServersClient.Get method.
func (client *SQLServersClient) Get(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, sqlServerName string, options *SQLServersClientGetOptions) (SQLServersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, sqlServerRegistrationName, sqlServerName, options)
	if err != nil {
		return SQLServersClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SQLServersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return SQLServersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *SQLServersClient) getCreateRequest(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, sqlServerName string, options *SQLServersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureData/sqlServerRegistrations/{sqlServerRegistrationName}/sqlServers/{sqlServerName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sqlServerRegistrationName == "" {
		return nil, errors.New("parameter sqlServerRegistrationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerRegistrationName}", url.PathEscape(sqlServerRegistrationName))
	if sqlServerName == "" {
		return nil, errors.New("parameter sqlServerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerName}", url.PathEscape(sqlServerName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	reqQP.Set("api-version", "2019-07-24-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SQLServersClient) getHandleResponse(resp *http.Response) (SQLServersClientGetResponse, error) {
	result := SQLServersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SQLServer); err != nil {
		return SQLServersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Gets all SQL Servers in a SQL Server Registration.
//
// Generated from API version 2019-07-24-preview
//   - resourceGroupName - Name of the resource group that contains the resource. You can obtain this value from the Azure Resource
//     Manager API or the portal.
//   - sqlServerRegistrationName - Name of the SQL Server registration.
//   - options - SQLServersClientListByResourceGroupOptions contains the optional parameters for the SQLServersClient.NewListByResourceGroupPager
//     method.
func (client *SQLServersClient) NewListByResourceGroupPager(resourceGroupName string, sqlServerRegistrationName string, options *SQLServersClientListByResourceGroupOptions) *runtime.Pager[SQLServersClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[SQLServersClientListByResourceGroupResponse]{
		More: func(page SQLServersClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SQLServersClientListByResourceGroupResponse) (SQLServersClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, sqlServerRegistrationName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return SQLServersClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return SQLServersClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return SQLServersClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *SQLServersClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, sqlServerRegistrationName string, options *SQLServersClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AzureData/sqlServerRegistrations/{sqlServerRegistrationName}/sqlServers"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sqlServerRegistrationName == "" {
		return nil, errors.New("parameter sqlServerRegistrationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sqlServerRegistrationName}", url.PathEscape(sqlServerRegistrationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	reqQP.Set("api-version", "2019-07-24-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *SQLServersClient) listByResourceGroupHandleResponse(resp *http.Response) (SQLServersClientListByResourceGroupResponse, error) {
	result := SQLServersClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SQLServerListResult); err != nil {
		return SQLServersClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

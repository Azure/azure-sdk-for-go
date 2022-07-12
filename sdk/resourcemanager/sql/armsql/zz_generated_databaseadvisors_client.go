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

// DatabaseAdvisorsClient contains the methods for the DatabaseAdvisors group.
// Don't use this type directly, use NewDatabaseAdvisorsClient() instead.
type DatabaseAdvisorsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDatabaseAdvisorsClient creates a new instance of DatabaseAdvisorsClient with the specified values.
// subscriptionID - The subscription ID that identifies an Azure subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDatabaseAdvisorsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DatabaseAdvisorsClient, error) {
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
	client := &DatabaseAdvisorsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Gets a database advisor.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// advisorName - The name of the Database Advisor.
// options - DatabaseAdvisorsClientGetOptions contains the optional parameters for the DatabaseAdvisorsClient.Get method.
func (client *DatabaseAdvisorsClient) Get(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, options *DatabaseAdvisorsClientGetOptions) (DatabaseAdvisorsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverName, databaseName, advisorName, options)
	if err != nil {
		return DatabaseAdvisorsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DatabaseAdvisorsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DatabaseAdvisorsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DatabaseAdvisorsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, options *DatabaseAdvisorsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/advisors/{advisorName}"
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
	if advisorName == "" {
		return nil, errors.New("parameter advisorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{advisorName}", url.PathEscape(advisorName))
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
func (client *DatabaseAdvisorsClient) getHandleResponse(resp *http.Response) (DatabaseAdvisorsClientGetResponse, error) {
	result := DatabaseAdvisorsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Advisor); err != nil {
		return DatabaseAdvisorsClientGetResponse{}, err
	}
	return result, nil
}

// ListByDatabase - Gets a list of database advisors.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// options - DatabaseAdvisorsClientListByDatabaseOptions contains the optional parameters for the DatabaseAdvisorsClient.ListByDatabase
// method.
func (client *DatabaseAdvisorsClient) ListByDatabase(ctx context.Context, resourceGroupName string, serverName string, databaseName string, options *DatabaseAdvisorsClientListByDatabaseOptions) (DatabaseAdvisorsClientListByDatabaseResponse, error) {
	req, err := client.listByDatabaseCreateRequest(ctx, resourceGroupName, serverName, databaseName, options)
	if err != nil {
		return DatabaseAdvisorsClientListByDatabaseResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DatabaseAdvisorsClientListByDatabaseResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DatabaseAdvisorsClientListByDatabaseResponse{}, runtime.NewResponseError(resp)
	}
	return client.listByDatabaseHandleResponse(resp)
}

// listByDatabaseCreateRequest creates the ListByDatabase request.
func (client *DatabaseAdvisorsClient) listByDatabaseCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, options *DatabaseAdvisorsClientListByDatabaseOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/advisors"
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
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByDatabaseHandleResponse handles the ListByDatabase response.
func (client *DatabaseAdvisorsClient) listByDatabaseHandleResponse(resp *http.Response) (DatabaseAdvisorsClientListByDatabaseResponse, error) {
	result := DatabaseAdvisorsClientListByDatabaseResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AdvisorArray); err != nil {
		return DatabaseAdvisorsClientListByDatabaseResponse{}, err
	}
	return result, nil
}

// Update - Updates a database advisor.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// serverName - The name of the server.
// databaseName - The name of the database.
// advisorName - The name of the Database Advisor.
// parameters - The requested advisor resource state.
// options - DatabaseAdvisorsClientUpdateOptions contains the optional parameters for the DatabaseAdvisorsClient.Update method.
func (client *DatabaseAdvisorsClient) Update(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, parameters Advisor, options *DatabaseAdvisorsClientUpdateOptions) (DatabaseAdvisorsClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serverName, databaseName, advisorName, parameters, options)
	if err != nil {
		return DatabaseAdvisorsClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DatabaseAdvisorsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DatabaseAdvisorsClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *DatabaseAdvisorsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, parameters Advisor, options *DatabaseAdvisorsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/advisors/{advisorName}"
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
	if advisorName == "" {
		return nil, errors.New("parameter advisorName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{advisorName}", url.PathEscape(advisorName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *DatabaseAdvisorsClient) updateHandleResponse(resp *http.Response) (DatabaseAdvisorsClientUpdateResponse, error) {
	result := DatabaseAdvisorsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Advisor); err != nil {
		return DatabaseAdvisorsClientUpdateResponse{}, err
	}
	return result, nil
}

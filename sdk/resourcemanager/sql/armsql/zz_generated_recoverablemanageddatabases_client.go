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

// RecoverableManagedDatabasesClient contains the methods for the RecoverableManagedDatabases group.
// Don't use this type directly, use NewRecoverableManagedDatabasesClient() instead.
type RecoverableManagedDatabasesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewRecoverableManagedDatabasesClient creates a new instance of RecoverableManagedDatabasesClient with the specified values.
// subscriptionID - The subscription ID that identifies an Azure subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewRecoverableManagedDatabasesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RecoverableManagedDatabasesClient, error) {
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
	client := &RecoverableManagedDatabasesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Gets a recoverable managed database.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// managedInstanceName - The name of the managed instance.
// options - RecoverableManagedDatabasesClientGetOptions contains the optional parameters for the RecoverableManagedDatabasesClient.Get
// method.
func (client *RecoverableManagedDatabasesClient) Get(ctx context.Context, resourceGroupName string, managedInstanceName string, recoverableDatabaseName string, options *RecoverableManagedDatabasesClientGetOptions) (RecoverableManagedDatabasesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, managedInstanceName, recoverableDatabaseName, options)
	if err != nil {
		return RecoverableManagedDatabasesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return RecoverableManagedDatabasesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return RecoverableManagedDatabasesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RecoverableManagedDatabasesClient) getCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, recoverableDatabaseName string, options *RecoverableManagedDatabasesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/recoverableDatabases/{recoverableDatabaseName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
	if recoverableDatabaseName == "" {
		return nil, errors.New("parameter recoverableDatabaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{recoverableDatabaseName}", url.PathEscape(recoverableDatabaseName))
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
func (client *RecoverableManagedDatabasesClient) getHandleResponse(resp *http.Response) (RecoverableManagedDatabasesClientGetResponse, error) {
	result := RecoverableManagedDatabasesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RecoverableManagedDatabase); err != nil {
		return RecoverableManagedDatabasesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByInstancePager - Gets a list of recoverable managed databases.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-11-01-preview
// resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
// Resource Manager API or the portal.
// managedInstanceName - The name of the managed instance.
// options - RecoverableManagedDatabasesClientListByInstanceOptions contains the optional parameters for the RecoverableManagedDatabasesClient.ListByInstance
// method.
func (client *RecoverableManagedDatabasesClient) NewListByInstancePager(resourceGroupName string, managedInstanceName string, options *RecoverableManagedDatabasesClientListByInstanceOptions) *runtime.Pager[RecoverableManagedDatabasesClientListByInstanceResponse] {
	return runtime.NewPager(runtime.PagingHandler[RecoverableManagedDatabasesClientListByInstanceResponse]{
		More: func(page RecoverableManagedDatabasesClientListByInstanceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RecoverableManagedDatabasesClientListByInstanceResponse) (RecoverableManagedDatabasesClientListByInstanceResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByInstanceCreateRequest(ctx, resourceGroupName, managedInstanceName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RecoverableManagedDatabasesClientListByInstanceResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return RecoverableManagedDatabasesClientListByInstanceResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RecoverableManagedDatabasesClientListByInstanceResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByInstanceHandleResponse(resp)
		},
	})
}

// listByInstanceCreateRequest creates the ListByInstance request.
func (client *RecoverableManagedDatabasesClient) listByInstanceCreateRequest(ctx context.Context, resourceGroupName string, managedInstanceName string, options *RecoverableManagedDatabasesClientListByInstanceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/managedInstances/{managedInstanceName}/recoverableDatabases"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if managedInstanceName == "" {
		return nil, errors.New("parameter managedInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{managedInstanceName}", url.PathEscape(managedInstanceName))
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

// listByInstanceHandleResponse handles the ListByInstance response.
func (client *RecoverableManagedDatabasesClient) listByInstanceHandleResponse(resp *http.Response) (RecoverableManagedDatabasesClientListByInstanceResponse, error) {
	result := RecoverableManagedDatabasesClientListByInstanceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RecoverableManagedDatabaseListResult); err != nil {
		return RecoverableManagedDatabasesClientListByInstanceResponse{}, err
	}
	return result, nil
}

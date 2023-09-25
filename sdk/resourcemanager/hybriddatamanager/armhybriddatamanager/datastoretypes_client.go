//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armhybriddatamanager

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

// DataStoreTypesClient contains the methods for the DataStoreTypes group.
// Don't use this type directly, use NewDataStoreTypesClient() instead.
type DataStoreTypesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewDataStoreTypesClient creates a new instance of DataStoreTypesClient with the specified values.
//   - subscriptionID - The Subscription Id
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewDataStoreTypesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DataStoreTypesClient, error) {
	cl, err := arm.NewClient(moduleName+".DataStoreTypesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &DataStoreTypesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Get - Gets the data store/repository type given its name.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-06-01
//   - dataStoreTypeName - The data store/repository type name for which details are needed.
//   - resourceGroupName - The Resource Group Name
//   - dataManagerName - The name of the DataManager Resource within the specified resource group. DataManager names must be between
//     3 and 24 characters in length and use any alphanumeric and underscore only
//   - options - DataStoreTypesClientGetOptions contains the optional parameters for the DataStoreTypesClient.Get method.
func (client *DataStoreTypesClient) Get(ctx context.Context, dataStoreTypeName string, resourceGroupName string, dataManagerName string, options *DataStoreTypesClientGetOptions) (DataStoreTypesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, dataStoreTypeName, resourceGroupName, dataManagerName, options)
	if err != nil {
		return DataStoreTypesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DataStoreTypesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DataStoreTypesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *DataStoreTypesClient) getCreateRequest(ctx context.Context, dataStoreTypeName string, resourceGroupName string, dataManagerName string, options *DataStoreTypesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridData/dataManagers/{dataManagerName}/dataStoreTypes/{dataStoreTypeName}"
	if dataStoreTypeName == "" {
		return nil, errors.New("parameter dataStoreTypeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataStoreTypeName}", url.PathEscape(dataStoreTypeName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataManagerName == "" {
		return nil, errors.New("parameter dataManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataManagerName}", url.PathEscape(dataManagerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DataStoreTypesClient) getHandleResponse(resp *http.Response) (DataStoreTypesClientGetResponse, error) {
	result := DataStoreTypesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataStoreType); err != nil {
		return DataStoreTypesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByDataManagerPager - Gets all the data store/repository types that the resource supports.
//
// Generated from API version 2019-06-01
//   - resourceGroupName - The Resource Group Name
//   - dataManagerName - The name of the DataManager Resource within the specified resource group. DataManager names must be between
//     3 and 24 characters in length and use any alphanumeric and underscore only
//   - options - DataStoreTypesClientListByDataManagerOptions contains the optional parameters for the DataStoreTypesClient.NewListByDataManagerPager
//     method.
func (client *DataStoreTypesClient) NewListByDataManagerPager(resourceGroupName string, dataManagerName string, options *DataStoreTypesClientListByDataManagerOptions) (*runtime.Pager[DataStoreTypesClientListByDataManagerResponse]) {
	return runtime.NewPager(runtime.PagingHandler[DataStoreTypesClientListByDataManagerResponse]{
		More: func(page DataStoreTypesClientListByDataManagerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DataStoreTypesClientListByDataManagerResponse) (DataStoreTypesClientListByDataManagerResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByDataManagerCreateRequest(ctx, resourceGroupName, dataManagerName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return DataStoreTypesClientListByDataManagerResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return DataStoreTypesClientListByDataManagerResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DataStoreTypesClientListByDataManagerResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByDataManagerHandleResponse(resp)
		},
	})
}

// listByDataManagerCreateRequest creates the ListByDataManager request.
func (client *DataStoreTypesClient) listByDataManagerCreateRequest(ctx context.Context, resourceGroupName string, dataManagerName string, options *DataStoreTypesClientListByDataManagerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.HybridData/dataManagers/{dataManagerName}/dataStoreTypes"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if dataManagerName == "" {
		return nil, errors.New("parameter dataManagerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataManagerName}", url.PathEscape(dataManagerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByDataManagerHandleResponse handles the ListByDataManager response.
func (client *DataStoreTypesClient) listByDataManagerHandleResponse(resp *http.Response) (DataStoreTypesClientListByDataManagerResponse, error) {
	result := DataStoreTypesClientListByDataManagerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataStoreTypeList); err != nil {
		return DataStoreTypesClientListByDataManagerResponse{}, err
	}
	return result, nil
}


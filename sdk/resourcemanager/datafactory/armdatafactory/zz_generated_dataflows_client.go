//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatafactory

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

// DataFlowsClient contains the methods for the DataFlows group.
// Don't use this type directly, use NewDataFlowsClient() instead.
type DataFlowsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDataFlowsClient creates a new instance of DataFlowsClient with the specified values.
// subscriptionID - The subscription identifier.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDataFlowsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DataFlowsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &DataFlowsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a data flow.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// factoryName - The factory name.
// dataFlowName - The data flow name.
// dataFlow - Data flow resource definition.
// options - DataFlowsClientCreateOrUpdateOptions contains the optional parameters for the DataFlowsClient.CreateOrUpdate
// method.
func (client *DataFlowsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, factoryName string, dataFlowName string, dataFlow DataFlowResource, options *DataFlowsClientCreateOrUpdateOptions) (DataFlowsClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, factoryName, dataFlowName, dataFlow, options)
	if err != nil {
		return DataFlowsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataFlowsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataFlowsClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *DataFlowsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, dataFlowName string, dataFlow DataFlowResource, options *DataFlowsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/dataflows/{dataFlowName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	if dataFlowName == "" {
		return nil, errors.New("parameter dataFlowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataFlowName}", url.PathEscape(dataFlowName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.IfMatch != nil {
		req.Raw().Header.Set("If-Match", *options.IfMatch)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, dataFlow)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *DataFlowsClient) createOrUpdateHandleResponse(resp *http.Response) (DataFlowsClientCreateOrUpdateResponse, error) {
	result := DataFlowsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataFlowResource); err != nil {
		return DataFlowsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a data flow.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// factoryName - The factory name.
// dataFlowName - The data flow name.
// options - DataFlowsClientDeleteOptions contains the optional parameters for the DataFlowsClient.Delete method.
func (client *DataFlowsClient) Delete(ctx context.Context, resourceGroupName string, factoryName string, dataFlowName string, options *DataFlowsClientDeleteOptions) (DataFlowsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, factoryName, dataFlowName, options)
	if err != nil {
		return DataFlowsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataFlowsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return DataFlowsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return DataFlowsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DataFlowsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, dataFlowName string, options *DataFlowsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/dataflows/{dataFlowName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	if dataFlowName == "" {
		return nil, errors.New("parameter dataFlowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataFlowName}", url.PathEscape(dataFlowName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets a data flow.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// factoryName - The factory name.
// dataFlowName - The data flow name.
// options - DataFlowsClientGetOptions contains the optional parameters for the DataFlowsClient.Get method.
func (client *DataFlowsClient) Get(ctx context.Context, resourceGroupName string, factoryName string, dataFlowName string, options *DataFlowsClientGetOptions) (DataFlowsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, factoryName, dataFlowName, options)
	if err != nil {
		return DataFlowsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataFlowsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataFlowsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DataFlowsClient) getCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, dataFlowName string, options *DataFlowsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/dataflows/{dataFlowName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	if dataFlowName == "" {
		return nil, errors.New("parameter dataFlowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataFlowName}", url.PathEscape(dataFlowName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.IfNoneMatch != nil {
		req.Raw().Header.Set("If-None-Match", *options.IfNoneMatch)
	}
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DataFlowsClient) getHandleResponse(resp *http.Response) (DataFlowsClientGetResponse, error) {
	result := DataFlowsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataFlowResource); err != nil {
		return DataFlowsClientGetResponse{}, err
	}
	return result, nil
}

// ListByFactory - Lists data flows.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// factoryName - The factory name.
// options - DataFlowsClientListByFactoryOptions contains the optional parameters for the DataFlowsClient.ListByFactory method.
func (client *DataFlowsClient) ListByFactory(resourceGroupName string, factoryName string, options *DataFlowsClientListByFactoryOptions) *runtime.Pager[DataFlowsClientListByFactoryResponse] {
	return runtime.NewPager(runtime.PageProcessor[DataFlowsClientListByFactoryResponse]{
		More: func(page DataFlowsClientListByFactoryResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DataFlowsClientListByFactoryResponse) (DataFlowsClientListByFactoryResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByFactoryCreateRequest(ctx, resourceGroupName, factoryName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return DataFlowsClientListByFactoryResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return DataFlowsClientListByFactoryResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DataFlowsClientListByFactoryResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByFactoryHandleResponse(resp)
		},
	})
}

// listByFactoryCreateRequest creates the ListByFactory request.
func (client *DataFlowsClient) listByFactoryCreateRequest(ctx context.Context, resourceGroupName string, factoryName string, options *DataFlowsClientListByFactoryOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataFactory/factories/{factoryName}/dataflows"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if factoryName == "" {
		return nil, errors.New("parameter factoryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{factoryName}", url.PathEscape(factoryName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByFactoryHandleResponse handles the ListByFactory response.
func (client *DataFlowsClient) listByFactoryHandleResponse(resp *http.Response) (DataFlowsClientListByFactoryResponse, error) {
	result := DataFlowsClientListByFactoryResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataFlowListResponse); err != nil {
		return DataFlowsClientListByFactoryResponse{}, err
	}
	return result, nil
}

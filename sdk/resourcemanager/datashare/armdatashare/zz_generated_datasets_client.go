//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatashare

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

// DataSetsClient contains the methods for the DataSets group.
// Don't use this type directly, use NewDataSetsClient() instead.
type DataSetsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewDataSetsClient creates a new instance of DataSetsClient with the specified values.
// subscriptionID - The subscription identifier
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewDataSetsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DataSetsClient, error) {
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
	client := &DataSetsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Create a DataSet
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share to add the data set to.
// dataSetName - The name of the dataSet.
// dataSet - The new data set information.
// options - DataSetsClientCreateOptions contains the optional parameters for the DataSetsClient.Create method.
func (client *DataSetsClient) Create(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, dataSet DataSetClassification, options *DataSetsClientCreateOptions) (DataSetsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, shareName, dataSetName, dataSet, options)
	if err != nil {
		return DataSetsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataSetsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return DataSetsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *DataSetsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, dataSet DataSetClassification, options *DataSetsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/dataSets/{dataSetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	if dataSetName == "" {
		return nil, errors.New("parameter dataSetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataSetName}", url.PathEscape(dataSetName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, dataSet)
}

// createHandleResponse handles the Create response.
func (client *DataSetsClient) createHandleResponse(resp *http.Response) (DataSetsClientCreateResponse, error) {
	result := DataSetsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return DataSetsClientCreateResponse{}, err
	}
	return result, nil
}

// BeginDelete - Delete a DataSet in a share
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// dataSetName - The name of the dataSet.
// options - DataSetsClientBeginDeleteOptions contains the optional parameters for the DataSetsClient.BeginDelete method.
func (client *DataSetsClient) BeginDelete(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, options *DataSetsClientBeginDeleteOptions) (*runtime.Poller[DataSetsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, accountName, shareName, dataSetName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[DataSetsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[DataSetsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Delete a DataSet in a share
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
func (client *DataSetsClient) deleteOperation(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, options *DataSetsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, shareName, dataSetName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DataSetsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, options *DataSetsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/dataSets/{dataSetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	if dataSetName == "" {
		return nil, errors.New("parameter dataSetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataSetName}", url.PathEscape(dataSetName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a DataSet in a share
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// dataSetName - The name of the dataSet.
// options - DataSetsClientGetOptions contains the optional parameters for the DataSetsClient.Get method.
func (client *DataSetsClient) Get(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, options *DataSetsClientGetOptions) (DataSetsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, shareName, dataSetName, options)
	if err != nil {
		return DataSetsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return DataSetsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DataSetsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DataSetsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, dataSetName string, options *DataSetsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/dataSets/{dataSetName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	if dataSetName == "" {
		return nil, errors.New("parameter dataSetName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{dataSetName}", url.PathEscape(dataSetName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *DataSetsClient) getHandleResponse(resp *http.Response) (DataSetsClientGetResponse, error) {
	result := DataSetsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return DataSetsClientGetResponse{}, err
	}
	return result, nil
}

// NewListBySharePager - List DataSets in a share
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-09-01
// resourceGroupName - The resource group name.
// accountName - The name of the share account.
// shareName - The name of the share.
// options - DataSetsClientListByShareOptions contains the optional parameters for the DataSetsClient.ListByShare method.
func (client *DataSetsClient) NewListBySharePager(resourceGroupName string, accountName string, shareName string, options *DataSetsClientListByShareOptions) *runtime.Pager[DataSetsClientListByShareResponse] {
	return runtime.NewPager(runtime.PagingHandler[DataSetsClientListByShareResponse]{
		More: func(page DataSetsClientListByShareResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *DataSetsClientListByShareResponse) (DataSetsClientListByShareResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByShareCreateRequest(ctx, resourceGroupName, accountName, shareName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return DataSetsClientListByShareResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return DataSetsClientListByShareResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return DataSetsClientListByShareResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByShareHandleResponse(resp)
		},
	})
}

// listByShareCreateRequest creates the ListByShare request.
func (client *DataSetsClient) listByShareCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *DataSetsClientListByShareOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/dataSets"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByShareHandleResponse handles the ListByShare response.
func (client *DataSetsClient) listByShareHandleResponse(resp *http.Response) (DataSetsClientListByShareResponse, error) {
	result := DataSetsClientListByShareResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DataSetList); err != nil {
		return DataSetsClientListByShareResponse{}, err
	}
	return result, nil
}

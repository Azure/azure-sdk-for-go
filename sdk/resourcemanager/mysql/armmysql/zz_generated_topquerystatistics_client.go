//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmysql

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

// TopQueryStatisticsClient contains the methods for the TopQueryStatistics group.
// Don't use this type directly, use NewTopQueryStatisticsClient() instead.
type TopQueryStatisticsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewTopQueryStatisticsClient creates a new instance of TopQueryStatisticsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewTopQueryStatisticsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*TopQueryStatisticsClient, error) {
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
	client := &TopQueryStatisticsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Retrieve the query statistic for specified identifier.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverName - The name of the server.
// queryStatisticID - The Query Statistic identifier.
// options - TopQueryStatisticsClientGetOptions contains the optional parameters for the TopQueryStatisticsClient.Get method.
func (client *TopQueryStatisticsClient) Get(ctx context.Context, resourceGroupName string, serverName string, queryStatisticID string, options *TopQueryStatisticsClientGetOptions) (TopQueryStatisticsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverName, queryStatisticID, options)
	if err != nil {
		return TopQueryStatisticsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return TopQueryStatisticsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return TopQueryStatisticsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *TopQueryStatisticsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverName string, queryStatisticID string, options *TopQueryStatisticsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/servers/{serverName}/topQueryStatistics/{queryStatisticId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if queryStatisticID == "" {
		return nil, errors.New("parameter queryStatisticID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{queryStatisticId}", url.PathEscape(queryStatisticID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *TopQueryStatisticsClient) getHandleResponse(resp *http.Response) (TopQueryStatisticsClientGetResponse, error) {
	result := TopQueryStatisticsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.QueryStatistic); err != nil {
		return TopQueryStatisticsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByServerPager - Retrieve the Query-Store top queries for specified metric and aggregation.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-06-01
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverName - The name of the server.
// parameters - The required parameters for retrieving top query statistics.
// options - TopQueryStatisticsClientListByServerOptions contains the optional parameters for the TopQueryStatisticsClient.ListByServer
// method.
func (client *TopQueryStatisticsClient) NewListByServerPager(resourceGroupName string, serverName string, parameters TopQueryStatisticsInput, options *TopQueryStatisticsClientListByServerOptions) *runtime.Pager[TopQueryStatisticsClientListByServerResponse] {
	return runtime.NewPager(runtime.PagingHandler[TopQueryStatisticsClientListByServerResponse]{
		More: func(page TopQueryStatisticsClientListByServerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TopQueryStatisticsClientListByServerResponse) (TopQueryStatisticsClientListByServerResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByServerCreateRequest(ctx, resourceGroupName, serverName, parameters, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return TopQueryStatisticsClientListByServerResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return TopQueryStatisticsClientListByServerResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return TopQueryStatisticsClientListByServerResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByServerHandleResponse(resp)
		},
	})
}

// listByServerCreateRequest creates the ListByServer request.
func (client *TopQueryStatisticsClient) listByServerCreateRequest(ctx context.Context, resourceGroupName string, serverName string, parameters TopQueryStatisticsInput, options *TopQueryStatisticsClientListByServerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforMySQL/servers/{serverName}/topQueryStatistics"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// listByServerHandleResponse handles the ListByServer response.
func (client *TopQueryStatisticsClient) listByServerHandleResponse(resp *http.Response) (TopQueryStatisticsClientListByServerResponse, error) {
	result := TopQueryStatisticsClientListByServerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.TopQueryStatisticsResultList); err != nil {
		return TopQueryStatisticsClientListByServerResponse{}, err
	}
	return result, nil
}

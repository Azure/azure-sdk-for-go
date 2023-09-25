//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsql

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

// ElasticPoolActivitiesClient contains the methods for the ElasticPoolActivities group.
// Don't use this type directly, use NewElasticPoolActivitiesClient() instead.
type ElasticPoolActivitiesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewElasticPoolActivitiesClient creates a new instance of ElasticPoolActivitiesClient with the specified values.
//   - subscriptionID - The subscription ID that identifies an Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewElasticPoolActivitiesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticPoolActivitiesClient, error) {
	cl, err := arm.NewClient(moduleName+".ElasticPoolActivitiesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ElasticPoolActivitiesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// NewListByElasticPoolPager - Returns elastic pool activities.
//
// Generated from API version 2014-04-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool for which to get the current activity.
//   - options - ElasticPoolActivitiesClientListByElasticPoolOptions contains the optional parameters for the ElasticPoolActivitiesClient.NewListByElasticPoolPager
//     method.
func (client *ElasticPoolActivitiesClient) NewListByElasticPoolPager(resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolActivitiesClientListByElasticPoolOptions) (*runtime.Pager[ElasticPoolActivitiesClientListByElasticPoolResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ElasticPoolActivitiesClientListByElasticPoolResponse]{
		More: func(page ElasticPoolActivitiesClientListByElasticPoolResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ElasticPoolActivitiesClientListByElasticPoolResponse) (ElasticPoolActivitiesClientListByElasticPoolResponse, error) {
			req, err := client.listByElasticPoolCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, options)
			if err != nil {
				return ElasticPoolActivitiesClientListByElasticPoolResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ElasticPoolActivitiesClientListByElasticPoolResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ElasticPoolActivitiesClientListByElasticPoolResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByElasticPoolHandleResponse(resp)
		},
	})
}

// listByElasticPoolCreateRequest creates the ListByElasticPool request.
func (client *ElasticPoolActivitiesClient) listByElasticPoolCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolActivitiesClientListByElasticPoolOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}/elasticPoolActivity"
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
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2014-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByElasticPoolHandleResponse handles the ListByElasticPool response.
func (client *ElasticPoolActivitiesClient) listByElasticPoolHandleResponse(resp *http.Response) (ElasticPoolActivitiesClientListByElasticPoolResponse, error) {
	result := ElasticPoolActivitiesClientListByElasticPoolResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ElasticPoolActivityListResult); err != nil {
		return ElasticPoolActivitiesClientListByElasticPoolResponse{}, err
	}
	return result, nil
}


//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcosmos

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

// PartitionKeyRangeIDRegionClient contains the methods for the PartitionKeyRangeIDRegion group.
// Don't use this type directly, use NewPartitionKeyRangeIDRegionClient() instead.
type PartitionKeyRangeIDRegionClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewPartitionKeyRangeIDRegionClient creates a new instance of PartitionKeyRangeIDRegionClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPartitionKeyRangeIDRegionClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PartitionKeyRangeIDRegionClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PartitionKeyRangeIDRegionClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// NewListMetricsPager - Retrieves the metrics determined by the given filter for the given partition key range id and region.
//
// Generated from API version 2023-03-15-preview
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - accountName - Cosmos DB database account name.
//   - region - Cosmos DB region, with spaces between words and each word capitalized.
//   - databaseRid - Cosmos DB database rid.
//   - collectionRid - Cosmos DB collection rid.
//   - partitionKeyRangeID - Partition Key Range Id for which to get data.
//   - filter - An OData filter expression that describes a subset of metrics to return. The parameters that can be filtered are
//     name.value (name of the metric, can have an or of multiple names), startTime, endTime,
//     and timeGrain. The supported operator is eq.
//   - options - PartitionKeyRangeIDRegionClientListMetricsOptions contains the optional parameters for the PartitionKeyRangeIDRegionClient.NewListMetricsPager
//     method.
func (client *PartitionKeyRangeIDRegionClient) NewListMetricsPager(resourceGroupName string, accountName string, region string, databaseRid string, collectionRid string, partitionKeyRangeID string, filter string, options *PartitionKeyRangeIDRegionClientListMetricsOptions) *runtime.Pager[PartitionKeyRangeIDRegionClientListMetricsResponse] {
	return runtime.NewPager(runtime.PagingHandler[PartitionKeyRangeIDRegionClientListMetricsResponse]{
		More: func(page PartitionKeyRangeIDRegionClientListMetricsResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *PartitionKeyRangeIDRegionClientListMetricsResponse) (PartitionKeyRangeIDRegionClientListMetricsResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PartitionKeyRangeIDRegionClient.NewListMetricsPager")
			req, err := client.listMetricsCreateRequest(ctx, resourceGroupName, accountName, region, databaseRid, collectionRid, partitionKeyRangeID, filter, options)
			if err != nil {
				return PartitionKeyRangeIDRegionClientListMetricsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PartitionKeyRangeIDRegionClientListMetricsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PartitionKeyRangeIDRegionClientListMetricsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listMetricsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listMetricsCreateRequest creates the ListMetrics request.
func (client *PartitionKeyRangeIDRegionClient) listMetricsCreateRequest(ctx context.Context, resourceGroupName string, accountName string, region string, databaseRid string, collectionRid string, partitionKeyRangeID string, filter string, options *PartitionKeyRangeIDRegionClientListMetricsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DocumentDB/databaseAccounts/{accountName}/region/{region}/databases/{databaseRid}/collections/{collectionRid}/partitionKeyRangeId/{partitionKeyRangeId}/metrics"
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
	if region == "" {
		return nil, errors.New("parameter region cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{region}", url.PathEscape(region))
	if databaseRid == "" {
		return nil, errors.New("parameter databaseRid cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseRid}", url.PathEscape(databaseRid))
	if collectionRid == "" {
		return nil, errors.New("parameter collectionRid cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{collectionRid}", url.PathEscape(collectionRid))
	if partitionKeyRangeID == "" {
		return nil, errors.New("parameter partitionKeyRangeID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partitionKeyRangeId}", url.PathEscape(partitionKeyRangeID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-03-15-preview")
	reqQP.Set("$filter", filter)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listMetricsHandleResponse handles the ListMetrics response.
func (client *PartitionKeyRangeIDRegionClient) listMetricsHandleResponse(resp *http.Response) (PartitionKeyRangeIDRegionClientListMetricsResponse, error) {
	result := PartitionKeyRangeIDRegionClientListMetricsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PartitionMetricListResult); err != nil {
		return PartitionKeyRangeIDRegionClientListMetricsResponse{}, err
	}
	return result, nil
}

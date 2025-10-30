// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// QueryItemsResponse contains response from the item query operation.
type QueryItemsResponse struct {
	Response
	// ContinuationToken contains the value of the x-ms-continuation header in the response.
	// It can be used to stop a query and resume it later.
	ContinuationToken *string
	// Contains the query metrics related to the query execution
	QueryMetrics *string
	// IndexMetrics contains the index utilization metrics if QueryOptions.PopulateIndexMetrics = true
	IndexMetrics *string
	// List of items.
	Items [][]byte
}

func newQueryResponse(resp *http.Response) (QueryItemsResponse, error) {
	response := QueryItemsResponse{
		Response: newResponse(resp),
	}

	continuationToken := resp.Header.Get(cosmosHeaderContinuationToken)
	if continuationToken != "" {
		response.ContinuationToken = &continuationToken
	}
	queryMetrics := resp.Header.Get(cosmosHeaderQueryMetrics)
	if queryMetrics != "" {
		response.QueryMetrics = &queryMetrics
	}
	queryIndexUtilization := resp.Header.Get(cosmosHeaderIndexUtilization)
	if queryIndexUtilization != "" {
		response.IndexMetrics = &queryIndexUtilization
	}

	result := queryServiceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return QueryItemsResponse{}, err
	}

	marshalledValue := make([][]byte, 0)
	for _, e := range result.Documents {
		m, err := json.Marshal(e)
		if err != nil {
			return QueryItemsResponse{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}
	response.Items = marshalledValue

	return response, nil
}

type queryServiceResponse struct {
	Documents []json.RawMessage `json:"Documents,omitempty"`
}

// QueryContainersResponse contains response from the container query operation.
type QueryContainersResponse struct {
	Response
	// ContinuationToken contains the value of the x-ms-continuation header in the response.
	// It can be used to stop a query and resume it later.
	ContinuationToken *string
	// List of containers.
	Containers []ContainerProperties
}

func newContainersQueryResponse(resp *http.Response) (QueryContainersResponse, error) {
	response := QueryContainersResponse{
		Response: newResponse(resp),
	}

	continuationToken := resp.Header.Get(cosmosHeaderContinuationToken)
	if continuationToken != "" {
		response.ContinuationToken = &continuationToken
	}
	result := queryContainersServiceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return QueryContainersResponse{}, err
	}

	response.Containers = result.Containers

	return response, nil
}

type queryContainersServiceResponse struct {
	Containers []ContainerProperties `json:"DocumentCollections,omitempty"`
}

// QueryDatabasesResponse contains response from the database query operation.
type QueryDatabasesResponse struct {
	Response
	// ContinuationToken contains the value of the x-ms-continuation header in the response.
	// It can be used to stop a query and resume it later.
	ContinuationToken *string
	// List of databases.
	Databases []DatabaseProperties
}

func newDatabasesQueryResponse(resp *http.Response) (QueryDatabasesResponse, error) {
	response := QueryDatabasesResponse{
		Response: newResponse(resp),
	}

	continuationToken := resp.Header.Get(cosmosHeaderContinuationToken)
	if continuationToken != "" {
		response.ContinuationToken = &continuationToken
	}

	result := queryDatabasesServiceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result); err != nil {
		return QueryDatabasesResponse{}, err
	}

	response.Databases = result.Databases

	return response, nil
}

type queryDatabasesServiceResponse struct {
	Databases []DatabaseProperties `json:"Databases,omitempty"`
}

// ReadManyItemsResponse contains response from the item query operation.
type ReadManyItemsResponse struct {
	RequestCharge float32
	// List of items.
	Items [][]byte
}

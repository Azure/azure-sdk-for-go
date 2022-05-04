// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// QueryItemsResponse contains response from the query operation.
type QueryItemsResponse struct {
	Response
	// ContinuationToken contains the value of the x-ms-continuation header in the response.
	// It can be used to stop a query and resume it later.
	ContinuationToken string
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

	response.ContinuationToken = resp.Header.Get(cosmosHeaderContinuationToken)
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
	Documents []map[string]interface{} `json:"Documents,omitempty"`
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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
	defer resp.Body.Close()
	body, err := azruntime.Payload(resp)
	if err != nil {
		return response, err
	}
	response.Value = body
	return response, nil
}

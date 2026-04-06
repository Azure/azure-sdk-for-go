// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// partitionKeyRangeResponse represents the response from GET /pkranges endpoint
// Contains the list of partition key ranges
// Rid is for the high level resource id
// count is for the number of partition key ranges returned
type partitionKeyRangeResponse struct {
	// ResourceID is the resource id of the partition key ranges
	ResourceID string `json:"_rid"`
	// PartitionKeyRanges contains the list of partition key ranges
	PartitionKeyRanges []partitionKeyRange `json:"PartitionKeyRanges"`
	// Count is the number of partition key ranges returned in the response
	Count int `json:"_count"`
	Response
}

// newPartitionKeyRangeResponse creates a new partitionKeyRangeResponse from an HTTP response
// It will parse the HTTP response and return a list of PartitionKeyRangeProperty objects
func newPartitionKeyRangeResponse(resp *http.Response) (partitionKeyRangeResponse, error) {
	response := partitionKeyRangeResponse{
		Response: newResponse(resp),
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := azruntime.Payload(resp)
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}

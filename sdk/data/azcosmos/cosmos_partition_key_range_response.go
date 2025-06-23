// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// PartitionKeyRange represents a single partition key range from GET /pkranges endpoint
// Includes all the needed fields to represent a partition key range
type PartitionKeyRange struct {
	ID                      string   `json:"id"`
	Rid                     string   `json:"_rid"`
	ETag                    string   `json:"_etag"`
	MinInclusive            string   `json:"minInclusive"`
	MaxExclusive            string   `json:"maxExclusive"`
	RidPrefix               int      `json:"ridPrefix"`
	Self                    string   `json:"_self"`
	ThroughputFraction      float64  `json:"throughputFraction"`
	Status                  string   `json:"status"`
	Parents                 []string `json:"parents"`
	OwnedArchivalPKRangeIds []string `json:"ownedArchivalPKRangeIds"`
	Timestamp               int64    `json:"_ts"`
	LSN                     int64    `json:"_lsn"`
}

// PartitionKeyRangeResponse represents the response from GET /pkranges endpoint
// Contains the list of partition key ranges
// Rid is for the high level resource id
// count is for the number of partition key ranges returned
type PartitionKeyRangeResponse struct {
	Response
	Rid                string              `json:"_rid"`
	PartitionKeyRanges []PartitionKeyRange `json:"PartitionKeyRanges"`
	Count              int                 `json:"_count"`
}

// newPartitionKeyRangeResponse creates a new PartitionKeyRangeResponse from an HTTP response
// It will parse the HTTP response and return a list of PartitionKeyRange objects
func newPartitionKeyRangeResponse(resp *http.Response) (PartitionKeyRangeResponse, error) {
	response := PartitionKeyRangeResponse{
		Response: newResponse(resp),
	}

	defer resp.Body.Close() // Standard practice to close the response body

	// Read the response body and slice with azruntime.Payload
	body, err := azruntime.Payload(resp)
	if err != nil {
		return response, err
	}
	// Unmarshal the JSON response into the PartitionKeysResponse struct
	// Acts like a deserializer that will map the JSON body to the struct createed
	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}

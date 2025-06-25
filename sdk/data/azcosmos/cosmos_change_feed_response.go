// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos


import (
	"encoding/json"
	"net/http"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ChangeFeedResponse contains the result of a change feed request.
type ChangeFeedResponse struct {
	// Response contains base information common to all responses.
	Response

	// ResourceID is the resource ID from the response
	ResourceID string `json:"_rid,omitempty"`

	// // ContinuationToken can be used in a subsequent request to continue reading changes.
	// ContinuationToken *string

	// Items contains the documents that have changed.
	Items [][]byte

	// Count is the number of items in the response.
	Count int
}

// newChangeFeedResponse creates a new ChangeFeedResponse from an HTTP response.
func newChangeFeedResponse(response *azResponse) (ChangeFeedResponse, error) {
	items, err := response.readItems()
	if err != nil {
		return ChangeFeedResponse{}, err
	}

	// continuationToken := response.readContinuationToken()
	
	// Extract the resource ID from the response
	var responseBody struct {
		ResourceID string `json:"_rid,omitempty"`
	}
	
	body, err := azruntime.Payload(response.Response)
	if err == nil {
		_ = json.Unmarshal(body, &responseBody)
	}
	
	result := ChangeFeedResponse{
		Response:   newResponse(response.Response),
		ResourceID: responseBody.ResourceID,
		Items:      items,
		Count:      len(items),
	}

	// if continuationToken != "" {
	// 	result.ContinuationToken = &continuationToken
	// }

	return result, nil
}

// newChangeFeedResponse creates a new ChangeFeedResponse from an HTTP response.
func newChangeFeedResponse(response *azResponse) (ChangeFeedResponse, error) {
	items, err := response.readItems()
	if err != nil {
		return ChangeFeedResponse{}, err
	}

	// continuationToken := response.readContinuationToken()
	
	result := ChangeFeedResponse{
		Response: newResponse(response.Response),
		Items:    items,
		Count:    len(items),
	}

	// if continuationToken != "" {
	// 	result.ContinuationToken = &continuationToken
	// }

	return result, nil
}

// // newPartitionKeyRangeResponse creates a new PartitionKeyRangeResponse from an HTTP response
// // It will parse the HTTP response and return a list of PartitionKeyRangeProperty objects
// func newPartitionKeyRangeResponse(resp *http.Response) (PartitionKeyRangeResponse, error) {
// 	response := PartitionKeyRangeResponse{
// 		Response: newResponse(resp),
// 	}

// 	defer resp.Body.Close()

// 	body, err := azruntime.Payload(resp)
// 	if err != nil {
// 		return response, err
// 	}

// 	if err := json.Unmarshal(body, &response); err != nil {
// 		return response, err
// 	}

// 	return response, nil
// }

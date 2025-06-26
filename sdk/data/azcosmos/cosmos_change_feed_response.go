// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ChangeFeedResponse contains the result of a change feed request.
type ChangeFeedResponse struct {
	// ResourceID is the unique identifier for the resource.
	ResourceID string `json:"_rid"`
	// Documents is a list of changed documents returned in the change feed.
	Documents []json.RawMessage `json:"Documents"`
	// Count is the number of documents returned in this page.
	Count int `json:"_count"`
	Response
}

// newChangeFeedResponse creates a new ChangeFeedResponse from an HTTP response.
func newChangeFeedResponse(response *http.Response) (ChangeFeedResponse, error) {
	if response == nil {
		return ChangeFeedResponse{}, errNilResponse
	}

	if response.StatusCode == http.StatusNotModified {
		// Handle 304 Not Modified response (no changes since the specified ETag)
		result := ChangeFeedResponse{
			Documents:       []json.RawMessage{}, // Empty array for documents
			Count:           0,
			Response:        response,
		}

		return result, nil
	}

	// For non-304 responses, read the body and parse the feed
	documents, err := readDocumentsFromResponse(response)
	if err != nil {
		return ChangeFeedResponse{}, err
	}

	// Create the response with parsed documents
	result := ChangeFeedResponse{
		ResourceID:  	response.Header.Get(cosmosHeaderResourceId),
		Documents:   	documents,
		Count:       	len(documents),
		Response: 		newResponse(response),
	}

	return result, nil
}

// readDocumentsFromResponse reads JSON documents from an HTTP response body.
func readDocumentsFromResponse(response *http.Response) ([]json.RawMessage, error) {
	if response == nil {
		return nil, errNilResponse
	}

	// Parse response body as a JSON array
	var documents []json.RawMessage
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&documents); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			// Empty response body, return empty array
			return []json.RawMessage{}, nil
		}
		return nil, err
	}

	return documents, nil
}

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
	// ResourceID is the unique identifier for the resource.
	ResourceID string `json:"_rid"`
	// Documents is a list of changed documents returned in the change feed.
	Documents []json.RawMessage `json:"Documents"`
	// Count is the number of documents returned in this page.
	Count int `json:"_count"`

	// Selected HTTP headers that we're retrieving from the response
	ETag              string
	ContinuationToken string
	LSN               string

	Response
}

// newChangeFeedResponse creates a new ChangeFeedResponse from an HTTP response.
func newChangeFeedResponse(resp *http.Response) (ChangeFeedResponse, error) {
	response := ChangeFeedResponse{
		Response:          newResponse(resp),
		ETag:              resp.Header.Get("etag"),
		ContinuationToken: resp.Header.Get("x-ms-continuation"),
		LSN:               resp.Header.Get("lsn"),
	}

	if resp.StatusCode == http.StatusNotModified {
		// Handle 304 Not Modified response (no changes since the specified ETag)
		response.Documents = []json.RawMessage{}
		response.Count = 0
		return response, nil
	}

	// For non-304 responses, unmarshal the response body
	defer resp.Body.Close()
	body, err := azruntime.Payload(resp)
	if err != nil {
		return response, err
	}
	// Parse the response into our response structure
	if err := json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}

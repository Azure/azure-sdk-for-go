// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

// GetContinuation from ChangeFeedResponse
func (c ChangeFeedResponse) GetContinuation() string {
	return c.ETag
}

// GetContRanges extracts the continuation token range from the ChangeFeedResponse.
func (c ChangeFeedResponse) GetContRanges() (min string, max string, ok bool) {
	if c.ContinuationToken == "" {
		return "", "", false
	}

	// Parse the continuation token JSON
	var contToken struct {
		Token *string
		Range struct {
			Min string
			Max string
		}
	}

	if err := json.Unmarshal([]byte(c.ContinuationToken), &contToken); err != nil {
		// Not a valid JSON continuation token
		return "", "", false
	}

	// Check if range values exist
	if contToken.Range.Min == "" || contToken.Range.Max == "" {
		return "", "", false
	}

	return contToken.Range.Min, contToken.Range.Max, true
}

// getCompositeContinuationToken creates a composite continuation token from the response.
// This token combines the feed range information with the ETag for use in subsequent requests.
func (c ChangeFeedResponse) getCompositeContinuationToken() (string, error) {
	// Extract the range from the continuation token
	min, max, ok := c.GetContRanges()
	if !ok {
		// No valid range in continuation token
		return "", nil
	}

	// Get the ETag
	etag := c.GetContinuation()
	fmt.Printf("ETag is this: %s\n", etag)
	if etag == "" {
		// No ETag available
		return "", nil
	}

	// Create the change feed range with continuation
	etagValue := azcore.ETag(etag)
	cfRange := newChangeFeedRange(min, max, &ChangeFeedRangeOptions{
		ContinuationToken: &etagValue,
	})

	// Create composite token
	compositeToken := newCompositeContinuationToken(c.ResourceID, []changeFeedRange{cfRange})

	// Marshal to JSON
	tokenBytes, err := json.Marshal(compositeToken)
	if err != nil {
		return "", err
	}

	return string(tokenBytes), nil
}

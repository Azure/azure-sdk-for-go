// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
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

	// ContinuationToken is the token used to continue reading the change feed.
	ContinuationToken string

	// Store the feed range if it was used in the request.
	FeedRange *FeedRange

	Response
}

// newChangeFeedResponse creates a new ChangeFeedResponse from an HTTP response.
func newChangeFeedResponse(resp *http.Response) (ChangeFeedResponse, error) {
	response := ChangeFeedResponse{
		Response: newResponse(resp),
	}

	if resp.StatusCode == http.StatusNotModified {
		response.Documents = []json.RawMessage{}
		response.Count = 0
		return response, nil
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

// PopulateCompositeContinuationToken generates and sets the composite continuation token if a feed range was used
func (response *ChangeFeedResponse) PopulateCompositeContinuationToken() {
	if response.FeedRange != nil && response.ETag != "" {
		compositeToken, err := response.GetCompositeContinuationToken()
		if err == nil && compositeToken != "" {
			response.ContinuationToken = compositeToken
		}
	}
}

// GetContinuation from ChangeFeedResponse
func (c ChangeFeedResponse) GetContinuation() string {
	return string(c.ETag)
}

// GetContRanges extracts the continuation token range from the ChangeFeedResponse.
func (c ChangeFeedResponse) GetContRanges() (min string, max string, ok bool) {
	if c.FeedRange != nil {
		return c.FeedRange.MinInclusive, c.FeedRange.MaxExclusive, true
	}

	if c.ContinuationToken == "" {
		return "", "", false
	}

	return "", "", false
}

// GetCompositeContinuationToken creates a composite continuation token from the response.
// This token combines the feed range information with the ETag for use in subsequent requests.
func (c ChangeFeedResponse) GetCompositeContinuationToken() (string, error) {
	min, max, ok := c.GetContRanges()
	if !ok {
		return "", nil
	}

	etag := c.GetContinuation()
	if etag == "" {
		return "", nil
	}

	etagValue := azcore.ETag(etag)
	cfRange := newChangeFeedRange(min, max, &ChangeFeedRangeOptions{
		ContinuationToken: &etagValue,
	})

	compositeToken := newCompositeContinuationToken(c.ResourceID, []changeFeedRange{cfRange})

	tokenBytes, err := json.Marshal(compositeToken)
	if err != nil {
		return "", err
	}

	return string(tokenBytes), nil
}

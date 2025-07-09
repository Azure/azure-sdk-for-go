// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestNewCompositeContinuationToken(t *testing.T) {
	// Creating the ResourceID, ContinuationToken to insert into the changeFeedRange
	resourceID := "testResource"
	continuationToken := azcore.ETag("14")
	changeFeedRangeInstance := newChangeFeedRange(
		"",
		"FF",
		&ChangeFeedRangeOptions{
			ContinuationToken: &continuationToken,
		},
	)

	compositeContinuationToken := newCompositeContinuationToken(resourceID, []changeFeedRange{changeFeedRangeInstance})
	t.Logf("ResourceID: %s\nContinuation: %+v", compositeContinuationToken.ResourceID, compositeContinuationToken.Continuation)

	// Marshal the compositeContinuationToken to JSON
	data, err := json.Marshal(compositeContinuationToken)
	if err != nil {
		t.Fatalf("Failed to marshal composite token: %v", err)
	}

	// Assertting the JSON output
	expectedJSON := `{"resourceId":"testResource","continuation":[{"minInclusive":"","maxExclusive":"FF","continuationToken":"14"}]}`
	if string(data) != expectedJSON {
		t.Errorf("Unexpected JSON output.\nExpected: %s\nActual:   %s", expectedJSON, string(data))
	}
}

func TestEmptyCompositeContinuationToken(t *testing.T) {
	// Test case with no FeedRange - should return empty token
	response := ChangeFeedResponse{
		ResourceID: "testResource",
		ETag:       "14",
	}

	token, err := response.getCompositeContinuationToken()
	if err != nil {
		t.Fatalf("Failed to get composite token: %v", err)
	}

	if token != "" {
		t.Errorf("Expected empty token but got: %s", token)
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestNewCompositeContinuationToken(t *testing.T) {
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

	data, err := json.Marshal(compositeContinuationToken)
	if err != nil {
		t.Fatalf("Failed to marshal composite token: %v", err)
	}

	expectedJSON := `{"version":1,"resourceId":"testResource","continuation":[{"minInclusive":"","maxExclusive":"FF","continuationToken":"14"}]}`
	if string(data) != expectedJSON {
		t.Errorf("Unexpected JSON output.\nExpected: %s\nActual:   %s", expectedJSON, string(data))
	}

	if compositeContinuationToken.Version != cosmosCompositeContinuationTokenVersion {
		t.Errorf("Unexpected version. Expected: %d, Actual: %d", cosmosCompositeContinuationTokenVersion, compositeContinuationToken.Version)
	}
}

func TestCompositeContinuationTokenWithMultipleRanges(t *testing.T) {
	resourceID := "testResource"
	continuation1 := azcore.ETag("14")
	continuation2 := azcore.ETag("15")

	changeFeedRange1 := newChangeFeedRange("00", "7F", &ChangeFeedRangeOptions{ContinuationToken: &continuation1})
	changeFeedRange2 := newChangeFeedRange("80", "FF", &ChangeFeedRangeOptions{ContinuationToken: &continuation2})

	compositeContinuationToken := newCompositeContinuationToken(resourceID, []changeFeedRange{changeFeedRange1, changeFeedRange2})

	if compositeContinuationToken.ResourceID != resourceID {
		t.Errorf("Resource ID mismatch. Expected: %s, Got: %s", resourceID, compositeContinuationToken.ResourceID)
	}

	if len(compositeContinuationToken.Continuation) != 2 {
		t.Errorf("Expected 2 continuation entries, got: %d", len(compositeContinuationToken.Continuation))
	}

	data, err := json.Marshal(compositeContinuationToken)
	if err != nil {
		t.Fatalf("Failed to marshal multi-range token: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"minInclusive":"00"`) {
		t.Errorf("First range minInclusive not found in JSON: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"maxExclusive":"7F"`) {
		t.Errorf("First range maxExclusive not found in JSON: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"minInclusive":"80"`) {
		t.Errorf("Second range minInclusive not found in JSON: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"maxExclusive":"FF"`) {
		t.Errorf("Second range maxExclusive not found in JSON: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"continuationToken":"14"`) {
		t.Errorf("First continuation token not found in JSON: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"continuationToken":"15"`) {
		t.Errorf("Second continuation token not found in JSON: %s", jsonStr)
	}
}

func TestEmptyCompositeContinuationToken(t *testing.T) {
	response := ChangeFeedResponse{
		ResourceID: "testResource",
	}
	token, err := response.GetCompositeContinuationToken()
	if err != nil {
		t.Fatalf("Failed to get composite token: %v", err)
	}

	if token != "" {
		t.Errorf("Expected empty token but got: %s", token)
	}
}

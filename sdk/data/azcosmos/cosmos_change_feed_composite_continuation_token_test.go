// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
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

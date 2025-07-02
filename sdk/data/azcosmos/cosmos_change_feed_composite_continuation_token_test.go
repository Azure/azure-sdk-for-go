// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
)

func TestNewCompositeContinuationToken(t *testing.T) {
	// Creating the ResourceID, ContinuationToken to insert into the ChangeFeedRange
	resourceID := "testResource"
	continuationToken := "14"
	changeFeedRange := newChangeFeedRange(
		"",
		"FF",
		&ChangeFeedRangeOptions{
			ContinuationToken: &continuationToken,
		},
	)

	compositeContinuationToken := newCompositeContinuationToken(resourceID, []ChangeFeedRange{changeFeedRange})
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

	// Unmarshal back to struct and check fields
	var unmarshaled CompositeContinuationToken
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if unmarshaled.ResourceID != resourceID {
		t.Errorf("ResourceID mismatch: got %s, want %s", unmarshaled.ResourceID, resourceID)
	}
	if len(unmarshaled.Continuation) != 1 {
		t.Fatalf("Expected 1 continuation, got %d", len(unmarshaled.Continuation))
	}
	if unmarshaled.Continuation[0].MinInclusive != "" ||
		unmarshaled.Continuation[0].MaxExclusive != "FF" ||
		unmarshaled.Continuation[0].ContinuationToken == nil ||
		string(*unmarshaled.Continuation[0].ContinuationToken) != "14" {
		t.Errorf("Continuation fields mismatch: %+v", unmarshaled.Continuation[0])
	}
}

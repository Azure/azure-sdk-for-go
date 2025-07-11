// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestNewContinuationTokenForPartitionKey(t *testing.T) {
	// Creating the ResourceID, PartitionKey, and Continuation token
	resourceID := "testResource"
	partitionKey := NewPartitionKeyString("testPartitionKey")
	continuationToken := azcore.ETag("14")

	// Create the continuation token for partition key
	tokenForPartitionKey := newContinuationTokenForPartitionKey(resourceID, &partitionKey, &continuationToken)
	t.Logf("ResourceID: %s\nPartitionKey: %+v\nContinuation: %v",
		tokenForPartitionKey.ResourceID,
		tokenForPartitionKey.PartitionKey,
		tokenForPartitionKey.Continuation)

	// Marshal the token to JSON
	data, err := json.Marshal(tokenForPartitionKey)
	require.NoError(t, err, "Failed to marshal token for partition key")

	// Verify the token serializes with the expected format
	expectedJSON := `{"resourceId":"testResource","partitionKey":["testPartitionKey"],"continuation":"14"}`
	require.Equal(t, expectedJSON, string(data), "Unexpected JSON output")

	// Test deserialization
	var unmarshaled continuationTokenForPartitionKey
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err, "Failed to unmarshal")

	require.Equal(t, resourceID, unmarshaled.ResourceID, "ResourceID mismatch")
	require.NotNil(t, unmarshaled.PartitionKey, "PartitionKey should not be nil")
	require.NotNil(t, unmarshaled.Continuation, "Continuation token should not be nil")
	require.Equal(t, string(continuationToken), string(*unmarshaled.Continuation), "Continuation token mismatch")
}

func TestContinuationTokenForPartitionKeyWithComplexKey(t *testing.T) {
	// Testing with a hierarchical partition key
	resourceID := "complexResource"
	partitionKey := NewPartitionKey().
		AppendString("level1").
		AppendNumber(42).
		AppendBool(true)
	continuationToken := azcore.ETag("complex-token")

	// Create the continuation token for partition key
	tokenForPartitionKey := newContinuationTokenForPartitionKey(resourceID, &partitionKey, &continuationToken)

	// Marshal the token to JSON
	data, err := json.Marshal(tokenForPartitionKey)
	require.NoError(t, err, "Failed to marshal complex token")

	// Verify serialization result with the expected format
	expectedJSON := `{"resourceId":"complexResource","partitionKey":["level1",42,true],"continuation":"complex-token"}`
	require.Equal(t, expectedJSON, string(data), "Unexpected JSON for complex token")

	// Test deserialization
	var unmarshaled continuationTokenForPartitionKey
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err, "Failed to unmarshal complex token")

	require.Equal(t, resourceID, unmarshaled.ResourceID, "ResourceID mismatch")
	require.NotNil(t, unmarshaled.PartitionKey, "PartitionKey should not be nil")
	require.NotNil(t, unmarshaled.Continuation, "Continuation token should not be nil")
	require.Equal(t, string(continuationToken), string(*unmarshaled.Continuation), "Continuation token mismatch")
}

func TestEmptyContinuationTokenForPartitionKey(t *testing.T) {
	// Test case with no FeedRange - should return empty token
	response := ChangeFeedResponse{
		ResourceID: "testResource",
		ETag:       "14",
	}

	token, err := response.GetContinuationTokenForPartitionKey()
	if err != nil {
		t.Fatalf("Failed to get composite token: %v", err)
	}

	if token != "" {
		t.Errorf("Expected empty token but got: %s", token)
	}
}

// Add this test to cosmos_change_feed_response_test.go if that file exists
func TestContinuationTokenForPartitionKeyWithNoETag(t *testing.T) {
	// Test case with PartitionKey but no ETag - should return empty token
	pk := NewPartitionKeyString("testPK")
	response := ChangeFeedResponse{
		ResourceID:   "testResource",
		PartitionKey: &pk,
		// No ETag set
	}

	token, err := response.GetContinuationTokenForPartitionKey()
	require.NoError(t, err, "Failed to get partition key token")

	require.Empty(t, token, "Expected empty token but got: %s", token)
}

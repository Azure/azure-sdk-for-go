// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

// TestPartitionKeyHeaderNotDoubleEncoded verifies that partition key headers
// are not double-encoded (e.g., [["value"]] instead of ["value"])
func TestPartitionKeyHeaderNotDoubleEncoded(t *testing.T) {
	pk := NewPartitionKeyString("test-value")

	// Get the JSON string representation
	pkJSON, err := pk.toJsonString()
	require.NoError(t, err)

	// Should be ["test-value"], NOT [["test-value"]]
	require.Equal(t, `["test-value"]`, pkJSON, "Partition key should be single-encoded JSON array")

	// Verify it doesn't start with [[ (double encoding)
	require.NotEqual(t, '[', pkJSON[1], "Partition key should not be double-encoded")
}

// TestPartitionKeyHeaderSetCorrectlyInPipeline verifies the header policy sets
// the partition key header with correct single encoding through the full pipeline
func TestPartitionKeyHeaderSetCorrectlyInPipeline(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	// Capture the actual header value
	var capturedPKHeader string
	capturer := &headerCapturer{onCapture: func(req *http.Request) {
		capturedPKHeader = req.Header.Get(cosmosHeaderPartitionKey)
	}}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, capturer}},
		&policy.ClientOptions{Transport: srv})

	req, err := azruntime.NewRequest(context.Background(), http.MethodPost, srv.URL())
	require.NoError(t, err)

	partitionKey := NewPartitionKeyString("my-partition-value")
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		headerOptionsOverride: &headerOptionsOverride{
			partitionKey: &partitionKey,
		},
	})

	_, err = pl.Do(req)
	require.NoError(t, err)

	// CRITICAL: The header MUST be single-encoded as ["value"], NOT double-encoded as [["value"]]
	require.Equal(t, `["my-partition-value"]`, capturedPKHeader,
		"Partition key header must be single-encoded JSON array")

	// Explicit check for double-encoding bug
	require.NotContains(t, capturedPKHeader, `[["`,
		"Partition key header must NOT be double-encoded")
	require.NotContains(t, capturedPKHeader, `"]]`,
		"Partition key header must NOT be double-encoded")
}

// TestPartitionKeyWithSpecialCharacters ensures special characters are properly escaped
func TestPartitionKeyWithSpecialCharacters(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected string
	}{
		{"simple string", "test", `["test"]`},
		{"with spaces", "hello world", `["hello world"]`},
		{"with quotes", `say "hello"`, `["say \"hello\""]`},
		{"with unicode", "café", `["caf\u00e9"]`},
		{"with backslash", `path\to\file`, `["path\\to\\file"]`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pk := NewPartitionKeyString(tc.value)
			pkJSON, err := pk.toJsonString()
			require.NoError(t, err)
			require.Equal(t, tc.expected, pkJSON)

			// Ensure no double encoding
			require.NotContains(t, pkJSON, `[["`)
		})
	}
}

// headerCapturer is a policy that captures HTTP headers for testing
type headerCapturer struct {
	onCapture func(req *http.Request)
}

func (p *headerCapturer) Do(req *policy.Request) (*http.Response, error) {
	if p.onCapture != nil {
		p.onCapture(req.Raw())
	}
	return req.Next()
}

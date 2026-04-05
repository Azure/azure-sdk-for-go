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

// TestGatewayModeNoDirectRoutingHeaders verifies that when Direct Mode transport
// is NOT enabled, no partition key range ID or effective partition key headers
// are added to requests (which would cause Gateway to reject them).
func TestGatewayModeNoDirectRoutingHeaders(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	var capturedHeaders http.Header
	capturer := &headerCapturer{onCapture: func(req *http.Request) {
		capturedHeaders = req.Header.Clone()
	}}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, capturer}},
		&policy.ClientOptions{Transport: srv})

	req, err := azruntime.NewRequest(context.Background(), http.MethodPost, srv.URL())
	require.NoError(t, err)

	partitionKey := NewPartitionKeyString("test-pk")
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		headerOptionsOverride: &headerOptionsOverride{
			partitionKey: &partitionKey,
		},
	})

	_, err = pl.Do(req)
	require.NoError(t, err)

	pkHeader := capturedHeaders.Get(cosmosHeaderPartitionKey)
	require.Equal(t, `["test-pk"]`, pkHeader, "Partition key header should be set")

	pkRangeID := capturedHeaders.Get(cosmosHeaderPartitionKeyRangeId)
	require.Empty(t, pkRangeID, "Partition key range ID must NOT be set in Gateway mode")

	epk := capturedHeaders.Get(cosmosHeaderEffectivePartitionKey)
	require.Empty(t, epk, "Effective partition key must NOT be set in Gateway mode")

	collRid := capturedHeaders.Get(cosmosHeaderCollectionRid)
	require.Empty(t, collRid, "Collection RID must NOT be set in Gateway mode")
}

// TestDirectModeHeadersSetWhenDirectTransportEnabled verifies that routing
// headers ARE set when Direct Mode transport is active
func TestDirectModeHeadersSetWhenDirectTransportEnabled(t *testing.T) {
	headerPolicy := &headerPolicies{}
	srv, close := mock.NewTLSServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusOK))

	var capturedHeaders http.Header
	capturer := &headerCapturer{onCapture: func(req *http.Request) {
		capturedHeaders = req.Header.Clone()
	}}

	pl := azruntime.NewPipeline("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{headerPolicy, capturer}},
		&policy.ClientOptions{Transport: srv})

	req, err := azruntime.NewRequest(context.Background(), http.MethodPost, srv.URL())
	require.NoError(t, err)

	partitionKey := NewPartitionKeyString("test-pk")
	req.SetOperationValue(pipelineRequestOptions{
		isWriteOperation: true,
		headerOptionsOverride: &headerOptionsOverride{
			partitionKey:          &partitionKey,
			partitionKeyRangeID:   "0",
			effectivePartitionKey: "ABC123",
			collectionRID:         "dOlzAKWJ04c=",
		},
	})

	_, err = pl.Do(req)
	require.NoError(t, err)

	pkHeader := capturedHeaders.Get(cosmosHeaderPartitionKey)
	require.Equal(t, `["test-pk"]`, pkHeader)

	pkRangeID := capturedHeaders.Get(cosmosHeaderPartitionKeyRangeId)
	require.Equal(t, "0", pkRangeID, "Partition key range ID should be set when explicitly provided")

	epk := capturedHeaders.Get(cosmosHeaderEffectivePartitionKey)
	require.Equal(t, "ABC123", epk, "Effective partition key should be set when explicitly provided")

	collRid := capturedHeaders.Get(cosmosHeaderCollectionRid)
	require.Equal(t, "dOlzAKWJ04c=", collRid, "Collection RID should be set when explicitly provided")
}

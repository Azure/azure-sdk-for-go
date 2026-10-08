// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestTransactionalBatchNativeRequest(t *testing.T) {
	disabled, enabled := false, true
	for _, contentResponse := range []*bool{nil, &disabled, &enabled} {
		batch := NewTransactionalBatch(NewPartitionKeyString("tenant").AppendNull().AppendNumber(3))
		require.NoError(t, batch.CreateItem([]byte(`{"id":"item","pk":"tenant","child":null,"leaf":3}`), nil))
		req, err := newTransactionalBatchRequest(batch, &TransactionalBatchOptions{
			Operation: OperationOptions{EnableContentResponseOnWrite: contentResponse}, SessionToken: "1:2",
		})
		require.NoError(t, err)
		native, release := inspectNativeItemRequest(req)
		release()
		require.Equal(t, int32(18), native.kind)
		require.Equal(t, 3, native.partitionKeyLen)
		require.Empty(t, native.itemID)
		require.Equal(t, "1:2", native.sessionToken)
		require.Equal(t, int32(-1), native.maxItemCount)
		require.Zero(t, native.preconditionKind, "conditions belong to individual JSON operations")
		require.Equal(t, req.body, native.body, "the singleton request carries the documented JSON operation array")
		expected := int32(0)
		if contentResponse != nil {
			expected = 1
			if *contentResponse {
				expected = 2
			}
		}
		require.Equal(t, expected, native.contentResponseWrite)
	}
}

func TestTransactionalBatchNativeCompletionOwnsEnvelopeAndHeaders(t *testing.T) {
	body := []byte(`[{"statusCode":201,"requestCharge":1.25,"eTag":"\"item\"","resourceBody":{"id":"item"}}]`)
	result := syntheticBatchCompletion(body, http.StatusOK)
	response, err := result.batchResponse(Response{RequestCharge: 2.5}, 1)
	require.NoError(t, err)
	require.True(t, response.Success)
	require.Equal(t, 7.0, response.RequestCharge)
	require.Equal(t, 3200, response.SubStatus)
	require.Equal(t, "activity-123", response.ActivityID)
	require.Equal(t, SessionToken("1:2"), response.SessionToken)
	require.Equal(t, azcore.ETag(`"etag"`), response.ETag)
	require.Equal(t, 125*time.Millisecond, response.RetryAfter)
	require.True(t, result.fromWire)
	require.Equal(t, body, response.Body)
	require.Equal(t, azcore.ETag(`"item"`), response.OperationResults[0].ETag)
	require.Equal(t, `{"id":"item"}`, string(response.OperationResults[0].ResourceBody))
}

func TestTransactionalBatchAuthoritativeCompletionSurvivesCancellation(t *testing.T) {
	for _, test := range []struct {
		status  int
		body    string
		success bool
	}{
		{200, `[{"statusCode":201},{"statusCode":204}]`, true},
		{207, `[{"statusCode":424},{"statusCode":409}]`, false},
	} {
		for range 100 {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			pending := &pendingOperation{result: make(chan completionResult, 1)}
			pending.result <- syntheticBatchCompletion([]byte(test.body), test.status)
			result, err := awaitOperationResult(ctx, pending, true)
			require.NoError(t, err, "received native write outcomes outrank context cancellation")
			response, err := result.batchResponse(Response{RequestCharge: 2}, 2)
			require.NoError(t, err)
			require.Equal(t, test.success, response.Success)
			require.Equal(t, test.body, string(response.Body))
			require.Equal(t, 6.5, response.RequestCharge)
			require.Equal(t, SessionToken("1:2"), response.SessionToken)
			require.Empty(t, pending.result)
		}
	}
}

func TestTransactionalBatchCancellationRetainsAllAvailableMetadata(t *testing.T) {
	for _, original := range []completionResult{
		{
			response: ItemResponse{Response: Response{StatusCode: 207, SubStatus: 3, RequestCharge: 5, ActivityID: "response", AttemptCount: 2},
				SessionToken: "session", ETag: `"etag"`},
			body: []byte("response-body"), retryAfter: time.Second, fromWire: true,
		},
		{err: &Error{
			Code: CodeClientOperationTimeout, StatusCode: 408, SubStatus: 20008, RequestCharge: 6, ActivityID: "error",
			AttemptCount: 3, SessionToken: "error-session", ETag: `"error-etag"`, RetryAfter: time.Second,
			FromWire: true, Body: []byte("error-body"),
		}},
	} {
		err := completionCancellationError(context.DeadlineExceeded, original)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		var cancelled *Error
		require.ErrorAs(t, err, &cancelled)
		require.Equal(t, CodeOperationCancelled, cancelled.Code)
		require.True(t, cancelled.FromWire)
		require.Equal(t, time.Second, cancelled.RetryAfter)
		if original.err == nil {
			require.Equal(t, original.response.StatusCode, cancelled.StatusCode)
			require.Equal(t, original.response.SubStatus, cancelled.SubStatus)
			require.Equal(t, original.response.RequestCharge, cancelled.RequestCharge)
			require.Equal(t, original.response.AttemptCount, cancelled.AttemptCount)
			require.Equal(t, original.response.SessionToken, cancelled.SessionToken)
			require.Equal(t, original.response.ETag, cancelled.ETag)
			require.Equal(t, original.body, cancelled.Body)
			original.body[0] = 'x'
			require.Equal(t, "response-body", string(cancelled.Body))
		} else {
			completionErr := original.err.(*Error)
			require.Equal(t, completionErr.StatusCode, cancelled.StatusCode)
			require.Equal(t, completionErr.SubStatus, cancelled.SubStatus)
			require.Equal(t, completionErr.RequestCharge, cancelled.RequestCharge)
			require.Equal(t, completionErr.AttemptCount, cancelled.AttemptCount)
			require.Equal(t, completionErr.SessionToken, cancelled.SessionToken)
			require.Equal(t, completionErr.ETag, cancelled.ETag)
			require.Equal(t, completionErr.Body, cancelled.Body)
			completionErr.Body[0] = 'x'
			require.Equal(t, "error-body", string(cancelled.Body))
		}
	}
}

func TestTransactionalBatchMetadataValidationErrorPreservesEvidence(t *testing.T) {
	result := syntheticBatchCompletion([]byte(`{"partitionKey":{"paths":["/pk","/child"],"kind":"MultiHash","version":2}}`), 200)
	err := validateCompleteBatchPartitionKey(result.body, NewPartitionKeyString("tenant"))
	require.Error(t, err)
	var enriched *Error
	require.ErrorAs(t, batchPartitionKeyError(err, result), &enriched)
	require.Equal(t, CodeBadRequest, enriched.Code)
	require.Equal(t, 4.5, enriched.RequestCharge)
	require.Equal(t, 200, enriched.StatusCode)
	require.Equal(t, SessionToken("1:2"), enriched.SessionToken)
	require.Equal(t, azcore.ETag(`"etag"`), enriched.ETag)
	require.Equal(t, 125*time.Millisecond, enriched.RetryAfter)
	require.Equal(t, result.body, enriched.Body)
	require.False(t, enriched.FromWire, "a successful metadata fetch must not classify a local validation failure as a wire failure")
	result.body[0] = 'x'
	require.NotEqual(t, result.body, enriched.Body)
}

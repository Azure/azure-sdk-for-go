// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
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
		require.NoError(t, batch.CreateItem(json.RawMessage(`{"id":"item","pk":"tenant","child":null,"leaf":3}`), nil))
		etag := azcore.ETag(`"unchanged"`)
		require.NoError(t, batch.ReadItem("item", &TransactionalBatchReadItemOptions{IfNoneMatchETag: &etag}))
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
		require.JSONEq(t, `[
			{"operationType":"Create","resourceBody":{"id":"item","pk":"tenant","child":null,"leaf":3}},
			{"operationType":"Read","id":"item","ifNoneMatch":"\"unchanged\""}
		]`, string(native.body))
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
	require.Equal(t, http.StatusOK, response.StatusCode)
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

func TestTransactionalBatchCancellationPreservesAvailableCompletion(t *testing.T) {
	for _, test := range []struct {
		status   int
		body     string
		statuses []int
	}{
		{200, `[{"statusCode":201},{"statusCode":204}]`, []int{201, 204}},
		{207, `[{"statusCode":424},{"statusCode":409}]`, []int{424, 409}},
		{207, `[{"statusCode":201},{"statusCode":304,"eTag":"\"unchanged\""}]`, []int{201, 304}},
	} {
		for range 100 {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			pending := &pendingOperation{result: make(chan completionResult, 1)}
			pending.result <- syntheticBatchCompletion([]byte(test.body), test.status)
			result, err := awaitOperationResult(ctx, pending)
			require.NoError(t, err, "already available native outcomes survive context cancellation")
			response, err := result.batchResponse(Response{RequestCharge: 2}, 2)
			require.NoError(t, err)
			require.Equal(t, test.status, response.StatusCode)
			require.Len(t, response.OperationResults, len(test.statuses))
			for i, status := range test.statuses {
				require.Equal(t, status, response.OperationResults[i].StatusCode)
			}
			require.Equal(t, test.body, string(response.Body))
			require.Equal(t, 6.5, response.RequestCharge)
			require.Equal(t, SessionToken("1:2"), response.SessionToken)
			require.Empty(t, pending.result)
		}
	}
}

func TestTransactionalBatchCancellationReleasesLateCompletion(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(cause.Error(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			if cause == context.DeadlineExceeded {
				cancel()
				ctx, cancel = context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
			}
			defer cancel()
			pending := &pendingOperation{result: make(chan completionResult, 1), delivered: make(chan struct{})}
			type outcome struct {
				result completionResult
				err    error
			}
			done := make(chan outcome, 1)
			go func() {
				result, err := awaitOperationResult(ctx, pending)
				done <- outcome{result, err}
			}()
			cancel()
			select {
			case result := <-done:
				require.Zero(t, result.result)
				require.ErrorIs(t, result.err, cause)
				var cosmosErr *Error
				require.ErrorAs(t, result.err, &cosmosErr)
				require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
			case <-time.After(time.Second):
				t.Fatal("batch did not stop waiting for an undelivered completion")
			}
			pending.deliver(syntheticBatchCompletion([]byte(`[{"statusCode":201},{"statusCode":304}]`), 207))
			require.Empty(t, pending.result, "late batch results must be released rather than retained")
			select {
			case <-pending.delivered:
			default:
				t.Fatal("late native completion was not drained")
			}
		})
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

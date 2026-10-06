// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResultAfterCancellationReturnsTerminalCompletion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result := completionResult{body: []byte("created")}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.NoError(t, err)
		require.Equal(t, result.body, got.body)
	})

	t.Run("failure", func(t *testing.T) {
		serviceErr := &Error{Code: CodeConflict}
		result := completionResult{err: serviceErr}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.NoError(t, err)
		require.Same(t, serviceErr, got.err)
	})

	t.Run("cancelled", func(t *testing.T) {
		result := completionResult{
			cancelled: true,
			err: &Error{
				Code:          CodeOperationCancelled,
				RequestCharge: 1.5,
				ActivityID:    "activity-id",
			},
		}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.Empty(t, got)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.NotErrorIs(t, err, context.Canceled)
		var cosmosErr *Error
		require.True(t, errors.As(err, &cosmosErr))
		require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
		require.Equal(t, 1.5, cosmosErr.RequestCharge)
		require.Equal(t, "activity-id", cosmosErr.ActivityID)
	})
}

func TestAwaitOperationResultCancelledWithDeliveredMetadata(t *testing.T) {
	for _, failed := range []bool{false, true} {
		for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
			for range 100 {
				ctx, cancel := context.WithCancel(t.Context())
				if cause == context.DeadlineExceeded {
					cancel()
					ctx, cancel = context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
				}
				cancel()
				result := completionResult{
					response: ItemResponse{Response: Response{RequestCharge: 3.5, ActivityID: "success"}},
				}
				charge, activity := 3.5, "success"
				if failed {
					result.err = &Error{Code: CodeBadRequest, RequestCharge: 4.75, ActivityID: "failure"}
					charge, activity = 4.75, "failure"
				}
				results := make(chan completionResult, 1)
				results <- result
				got, err := awaitOperationResult(ctx, results, false, func() {
					t.Fatal("an already-delivered completion must not be abandoned")
				})
				require.Zero(t, got)
				require.Empty(t, results)
				require.ErrorIs(t, err, cause)
				var cosmosErr *Error
				require.ErrorAs(t, err, &cosmosErr)
				require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
				require.Equal(t, charge, cosmosErr.RequestCharge)
				require.Equal(t, activity, cosmosErr.ActivityID)
			}
		}
	}
}

func TestAwaitOperationResultAbandonsOnlyUndelivered(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	abandoned := false
	got, err := awaitOperationResult(ctx, make(chan completionResult, 1), false, func() { abandoned = true })
	require.True(t, abandoned)
	require.Zero(t, got)
	require.ErrorIs(t, err, context.Canceled)
}

func TestAwaitOperationResultPreservesAuthoritativeOutcome(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, result := range []completionResult{
		{body: []byte("committed")},
		{err: &Error{Code: CodeConflict, RequestCharge: 2}},
	} {
		results := make(chan completionResult, 1)
		results <- result
		got, err := awaitOperationResult(ctx, results, true, func() { t.Fatal("authoritative operation abandoned") })
		require.NoError(t, err)
		require.Equal(t, result, got)
	}
}

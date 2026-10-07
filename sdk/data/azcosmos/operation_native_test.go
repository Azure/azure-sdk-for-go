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

	t.Run("client operation timeout", func(t *testing.T) {
		// The native driver's own end-to-end timeout can fire a ClientOperationTimeout completion
		// in a race with the caller's context deadline. An authoritative wait always awaits the
		// real completion, so without this it would surface that raw error instead of one callers
		// can detect with errors.Is(err, context.DeadlineExceeded), even though the context is why
		// the operation stopped mattering to the caller.
		diagnostics := &Diagnostics{AttemptCount: 1, StatusCode: 408, SubStatus: 20008}
		result := completionResult{
			err: &Error{
				Code:          CodeClientOperationTimeout,
				RequestCharge: 0.5,
				ActivityID:    "timeout-activity-id",
				Diagnostics:   diagnostics,
				AttemptCount:  1,
				StatusCode:    408,
				SubStatus:     20008,
			},
		}

		got, err := resultAfterCancellation(context.DeadlineExceeded, result)

		require.Empty(t, got)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		var cosmosErr *Error
		require.True(t, errors.As(err, &cosmosErr))
		require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
		require.Equal(t, 0.5, cosmosErr.RequestCharge)
		require.Equal(t, "timeout-activity-id", cosmosErr.ActivityID)
	})

	t.Run("cancelled", func(t *testing.T) {
		diagnostics := &Diagnostics{AttemptCount: 2, StatusCode: 429, SubStatus: 3200}
		result := completionResult{
			cancelled: true,
			err: &Error{
				Code:          CodeOperationCancelled,
				RequestCharge: 1.5,
				ActivityID:    "activity-id",
				Diagnostics:   diagnostics,
				AttemptCount:  2,
				StatusCode:    429,
				SubStatus:     3200,
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
		require.Same(t, diagnostics, cosmosErr.Diagnostics)
		require.Equal(t, uint32(2), cosmosErr.AttemptCount)
		require.Equal(t, 429, cosmosErr.StatusCode)
		require.Equal(t, 3200, cosmosErr.SubStatus)
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
				pending := &pendingOperation{result: make(chan completionResult, 1)}
				pending.result <- result
				got, err := awaitOperationResult(ctx, pending, false)
				require.Zero(t, got)
				require.Empty(t, pending.result)
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
	pending := &pendingOperation{result: make(chan completionResult, 1)}
	got, err := awaitOperationResult(ctx, pending, false)
	require.Zero(t, got)
	require.ErrorIs(t, err, context.Canceled)
	// awaitOperationResult must have abandoned the pending operation rather than leaving it open:
	// abandoning it again finds nothing buffered, since nothing was ever delivered.
	_, ok := pending.abandon()
	require.False(t, ok)
}

func TestAwaitOperationResultPreservesAuthoritativeOutcome(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	for _, result := range []completionResult{
		{body: []byte("committed")},
		{err: &Error{Code: CodeConflict, RequestCharge: 2}},
	} {
		pending := &pendingOperation{result: make(chan completionResult, 1)}
		pending.result <- result
		got, err := awaitOperationResult(ctx, pending, true)
		require.NoError(t, err)
		require.Equal(t, result, got)
	}
}

func TestAwaitOperationResultWrapsClientOperationTimeoutAfterCancellation(t *testing.T) {
	// The native driver's own end-to-end timeout races the caller's context deadline on an
	// independent clock. An authoritative wait always awaits the real completion rather than
	// abandoning early, so when the native clock wins that race, this must still surface as a
	// cancellation error rather than the raw ClientOperationTimeout, matching the non-authoritative
	// path's errors.Is(err, context.DeadlineExceeded) contract.
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	pending := &pendingOperation{result: make(chan completionResult, 1)}
	pending.result <- completionResult{
		err: &Error{Code: CodeClientOperationTimeout, RequestCharge: 0.5, ActivityID: "timeout-activity-id"},
	}

	got, err := awaitOperationResult(ctx, pending, true)

	require.Zero(t, got)
	require.ErrorIs(t, err, context.Canceled)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
	require.Equal(t, 0.5, cosmosErr.RequestCharge)
	require.Equal(t, "timeout-activity-id", cosmosErr.ActivityID)
}

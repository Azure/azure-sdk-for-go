// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAwaitOperationResultWithoutCancellation(t *testing.T) {
	for _, result := range []completionResult{
		{body: []byte("committed")},
		{err: &Error{Code: CodeConflict, RequestCharge: 2}},
		{err: &Error{Code: CodeClientOperationTimeout}},
	} {
		pending := &pendingOperation{result: make(chan completionResult, 1)}
		pending.deliver(result)
		got, err := awaitOperationResult(context.Background(), pending)
		require.NoError(t, err)
		require.Equal(t, result, got)
	}
}

func TestAwaitOperationResultCancelledWithDeliveredMetadata(t *testing.T) {
	const trackingID PatchTrackingID = "00112233-4455-6677-8899-aabbccddeeff"
	for _, code := range []Code{"", CodeBadRequest, CodeClientOperationTimeout} {
		for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
			for range 100 {
				ctx, cancel := context.WithCancel(context.Background())
				if cause == context.DeadlineExceeded {
					cancel()
					ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				}
				cancel()
				diagnostics := &Diagnostics{AttemptCount: 2}
				result := completionResult{response: ItemResponse{
					Response: Response{RequestCharge: 3.5, ActivityID: "activity", AttemptCount: 2,
						StatusCode: 200, SubStatus: 42, Diagnostics: diagnostics},
					PatchTrackingID: trackingID,
				}}
				if code != "" {
					result.err = &Error{Code: code, RequestCharge: 3.5, ActivityID: "activity",
						AttemptCount: 2, StatusCode: 400, SubStatus: 42, Diagnostics: diagnostics,
						PatchTrackingID: trackingID}
				}
				pending := &pendingOperation{result: make(chan completionResult, 1)}
				pending.deliver(result)
				got, err := awaitOperationResult(ctx, pending)
				require.Zero(t, got)
				require.ErrorIs(t, err, cause)
				require.Empty(t, pending.result)
				var cancelled *Error
				require.ErrorAs(t, err, &cancelled)
				require.Equal(t, CodeOperationCancelled, cancelled.Code)
				require.Equal(t, 3.5, cancelled.RequestCharge)
				require.Equal(t, "activity", cancelled.ActivityID)
				require.Equal(t, uint32(2), cancelled.AttemptCount)
				require.Equal(t, 42, cancelled.SubStatus)
				require.Same(t, diagnostics, cancelled.Diagnostics)
				require.Equal(t, trackingID, cancelled.PatchTrackingID)
				require.Same(t, err, patchCancellationError(err, ""))
				require.Equal(t, trackingID, cancelled.PatchTrackingID)
			}
		}
	}
}

func TestAwaitOperationResultAbandonsOnlyUndelivered(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	pending := &pendingOperation{result: make(chan completionResult, 1)}
	got, err := awaitOperationResult(ctx, pending)
	require.Zero(t, got)
	require.ErrorIs(t, err, context.Canceled)
	pending.deliver(completionResult{body: []byte("late committed write")})
	require.Empty(t, pending.result, "late results must be released rather than retained")
	const id PatchTrackingID = "00112233-4455-6677-8899-aabbccddeeff"
	var cancelled *Error
	require.ErrorAs(t, patchCancellationError(err, id), &cancelled)
	require.Equal(t, id, cancelled.PatchTrackingID)
}

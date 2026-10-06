// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/stretchr/testify/require"
)

func TestThrottledCompletionCopiesAndClassifiesEveryField(t *testing.T) {
	for _, tt := range []struct {
		name            string
		packedSubStatus int
	}{
		{"header refines absent packed sub-status", 0},
		{"packed sub-status is authoritative", 3200},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := syntheticThrottledCompletion(tt.packedSubStatus)

			require.Equal(t, CodeThrottled, err.Code)
			require.Equal(t, http.StatusTooManyRequests, err.StatusCode)
			require.Equal(t, 3200, err.SubStatus)
			require.Equal(t, 4.5, err.RequestCharge)
			require.Equal(t, "activity-123", err.ActivityID)
			require.Equal(t, SessionToken("1:2"), err.SessionToken)
			require.Equal(t, azcore.ETag("\"etag\""), err.ETag)
			require.Equal(t, 125*time.Millisecond, err.RetryAfter)
			require.Equal(t, "Request rate is large.", err.Message)
			require.JSONEq(t, `{"code":"TooManyRequests","message":"slow down"}`, string(err.Body))
			require.True(t, err.FromWire)
		})
	}

}

func TestNativeResponseCopiesUpdatedSessionToken(t *testing.T) {
	response := syntheticSessionCompletion()
	require.Equal(t, SessionToken("0:-1#43"), response.SessionToken)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Zero(t, response.SubStatus)
}

func TestAbsentNativeDiagnostics(t *testing.T) {
	diagnostics, err := copyDiagnostics(nil)
	require.NoError(t, err)
	require.Nil(t, diagnostics)
}

func TestDiagnosticDurationsSaturate(t *testing.T) {
	require.Equal(t, time.Duration(1<<63-1), durationMicros(^uint64(0)))
	require.Equal(t, time.Duration(1<<63-1), durationMillis(^uint64(0)))
	require.Equal(t, 3*time.Microsecond, durationMicros(3))
	require.Equal(t, 3*time.Millisecond, durationMillis(3))
}

// stubDiagnosticsCopy substitutes copyDiagnosticsForCompletion for the calling test, restoring the
// original on cleanup. cosmos_diagnostics_t is an opaque native handle only the driver can produce
// and a real copy failure can only be forced by actual OOM, so this is the seam tests use to
// simulate a diagnostics-copy outcome instead.
func stubDiagnosticsCopy(t *testing.T, diagnostics *Diagnostics, err error) {
	t.Helper()
	t.Cleanup(stubCopyDiagnosticsForCompletion(diagnostics, err))
}

// captureDiagnosticsLog installs a log listener for the calling test and returns the messages
// logged under EventDiagnostics, restoring the default listener on cleanup.
func captureDiagnosticsLog(t *testing.T) *[]string {
	t.Helper()
	messages := &[]string{}
	log.SetListener(func(event log.Event, message string) {
		if event == EventDiagnostics {
			*messages = append(*messages, message)
		}
	})
	t.Cleanup(func() { log.SetListener(nil) })
	return messages
}

func TestDiagnosticsCopyFailureDoesNotMaskSuccess(t *testing.T) {
	fakeDiagnostics := &Diagnostics{AttemptCount: 3, TotalRequestCharge: 7}
	stubDiagnosticsCopy(t, fakeDiagnostics, errors.New("boom"))
	messages := captureDiagnosticsLog(t)

	result := syntheticCompletionResult(syntheticOutcomeOK, http.StatusOK, false)

	require.NoError(t, result.err)
	require.Equal(t, http.StatusOK, result.response.StatusCode)
	require.Same(t, fakeDiagnostics, result.response.Diagnostics)
	require.Equal(t, uint32(3), result.response.AttemptCount)
	require.Len(t, *messages, 1)
	require.Contains(t, (*messages)[0], "boom")
}

func TestDiagnosticsCopyFailureDoesNotMaskCancellation(t *testing.T) {
	fakeDiagnostics := &Diagnostics{AttemptCount: 2}
	stubDiagnosticsCopy(t, fakeDiagnostics, errors.New("boom"))
	messages := captureDiagnosticsLog(t)

	result := syntheticCompletionResult(syntheticOutcomeCancelled, 0, false)

	require.True(t, result.cancelled)
	var cosmosErr *Error
	require.ErrorAs(t, result.err, &cosmosErr)
	require.Equal(t, CodeOperationCancelled, cosmosErr.Code)
	require.Same(t, fakeDiagnostics, cosmosErr.Diagnostics)
	require.Equal(t, uint32(2), cosmosErr.AttemptCount)
	// The cancellation cause must still fall back to context.Canceled; a diagnostics copy
	// failure must never hijack Unwrap()'s cause.
	require.ErrorIs(t, cosmosErr, context.Canceled)
	require.Len(t, *messages, 1)
}

func TestDiagnosticsCopyFailureAppendsToOperationErrorWithoutBecomingItsCause(t *testing.T) {
	stubDiagnosticsCopy(t, nil, errors.New("boom"))
	messages := captureDiagnosticsLog(t)

	result := syntheticCompletionResult(syntheticOutcomeError, http.StatusNotFound, true)

	var cosmosErr *Error
	require.ErrorAs(t, result.err, &cosmosErr)
	require.Equal(t, CodeNotFound, cosmosErr.Code)
	require.Contains(t, cosmosErr.Message, "boom")
	require.Nil(t, cosmosErr.Unwrap(), "a diagnostics copy failure must not become the operation error's cause")
	require.Len(t, *messages, 1)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/stretchr/testify/require"
)

func TestNativePackedStatusRepresentation(t *testing.T) {
	for _, tt := range []struct {
		bits       uint32
		signed     int64
		httpStatus int
		subStatus  int
	}{
		{0x00000000, 0, 0, 0},
		{0x019403ea, 26477546, 404, 1002},
		{0x01ad0c80, 28118144, 429, 3200},
		{0x01f74e2a, 32984618, 503, 20010},
		{0x0000ffff, 65535, 0, 65535},
		{0x7fffffff, 2147483647, 32767, 65535},
		{0x80000000, -2147483648, 32768, 0},
		{0xffff0000, -65536, 65535, 0},
		{0xffffffff, -1, 65535, 65535},
	} {
		t.Run(fmt.Sprintf("%08x", tt.bits), func(t *testing.T) {
			bits, signed, width, status, subStatus := nativeStatusBits(tt.bits)
			require.Equal(t, tt.bits, bits)
			require.Equal(t, tt.signed, signed)
			require.Equal(t, 4, width, "cosmos_status_code_t must remain signed int32")
			require.Equal(t, tt.httpStatus, status)
			require.Equal(t, tt.subStatus, subStatus)
			macroStatus, macroSubStatus := nativeStatusMacroParts(tt.bits)
			require.Equal(t, macroStatus, status, "Go decoding must match COSMOS_STATUS_HTTP")
			require.Equal(t, macroSubStatus, subStatus, "Go decoding must match COSMOS_STATUS_SUB")
		})
	}
}

func TestNativeCompletionErrorMappings(t *testing.T) {
	for _, tt := range []struct {
		name       string
		httpStatus int
		packedSub  int
		headerSub  int
		wantSub    int
		fromWire   bool
		cancelled  bool
		want       Code
	}{
		{"wire bad request", 400, 0, -1, 0, true, false, CodeBadRequest},
		{"wire unauthorized", 401, 0, 0, 0, true, false, CodeUnauthorized},
		{"wire forbidden", 403, 0, 0, 0, true, false, CodeForbidden},
		{"wire missing", 404, 0, 0, 0, true, false, CodeNotFound},
		{"wire session", 404, 1002, 0, 1002, true, false, CodeSessionUnavailable},
		{"header session fallback", 404, 0, 1002, 1002, true, false, CodeSessionUnavailable},
		{"wire timeout", 408, 0, 0, 0, true, false, CodeRequestTimeout},
		{"wire conflict ignores 1002", 409, 1002, 0, 1002, true, false, CodeConflict},
		{"wire gone", 410, 0, 0, 0, true, false, CodeGone},
		{"wire split", 410, 1002, 3200, 1002, true, false, CodePartitionKeyRangeGone},
		{"wire precondition", 412, 0, 0, 0, true, false, CodePreconditionFailed},
		{"wire throttle", 429, 3200, 0, 3200, true, false, CodeThrottled},
		{"wire unavailable ignores synthetic", 503, 20010, 0, 20010, true, false, CodeServiceUnavailable},
		{"unknown wire status and substatus", 599, 65535, 0, 65535, true, false, CodeServiceError},
		{"client transport", 503, 20010, 0, 20010, false, false, CodeTransportFailure},
		{"client unavailable after service failure", 503, 20310, 0, 20310, false, false, CodeServiceUnavailable},
		{"client unavailable after region failures", 503, 20311, 0, 20311, false, false, CodeServiceUnavailable},
		{"client unavailable after partition failures", 503, 20312, 0, 20312, false, false, CodeServiceUnavailable},
		{"client timeout", 408, 20008, 0, 20008, false, false, CodeClientOperationTimeout},
		{"client response serialization", 500, 20020, 0, 20020, false, false, CodeSerializationFailed},
		{"client request serialization", 400, 20021, 0, 20021, false, false, CodeSerializationFailed},
		{"client authentication", 401, 20402, 0, 20402, false, false, CodeAuthenticationFailed},
		{"client generated 401", 401, 20401, 0, 20401, false, false, CodeAuthenticationFailed},
		{"client validation", 400, 20353, 0, 20353, false, false, CodeClientError},
		{"unknown synthetic", 503, 22000, 0, 22000, false, false, CodeClientError},
		{"no service status", 0, 0, -1, 0, false, false, CodeClientError},
		{"reserved cancellation substatus remains generic", 500, 20360, 0, 20360, false, false, CodeClientError},
		{"reserved cancellation outcome", 0, 20360, 0, 20360, false, true, CodeOperationCancelled},
		{"negative packed sentinel retains unsigned halves", 65535, 65535, 0, 65535, false, false, CodeClientError},
		{"high packed sign bit", 32768, 1002, 0, 1002, false, false, CodeClientError},
	} {
		t.Run(tt.name, func(t *testing.T) {
			diagnostics := &Diagnostics{
				AttemptCount: 2, RegionsContacted: []string{"synthetic-region"},
				Attempts: []DiagnosticAttempt{{StatusCode: tt.httpStatus, SubStatus: tt.wantSub}},
			}
			stubDiagnosticsCopy(t, diagnostics, nil)
			var body []byte
			if tt.fromWire {
				body = []byte(`{"code":"synthetic","message":"boundary failure"}`)
			}
			result := nativeMappingCompletion(uint32(tt.httpStatus)<<16|uint32(tt.packedSub),
				tt.httpStatus, tt.headerSub, tt.fromWire, tt.cancelled, body)
			var failure *Error
			require.ErrorAs(t, result.err, &failure)
			require.Equal(t, ItemResponse{}, result.response)
			require.Equal(t, tt.cancelled, result.cancelled)
			require.Equal(t, tt.want, failure.Code)
			require.Equal(t, tt.httpStatus, failure.StatusCode)
			require.Equal(t, tt.wantSub, failure.SubStatus)
			require.Equal(t, tt.fromWire, failure.FromWire)
			require.Equal(t, 4.5, failure.RequestCharge)
			require.Equal(t, "activity-123", failure.ActivityID)
			require.Equal(t, SessionToken("1:2"), failure.SessionToken)
			require.Equal(t, azcore.ETag(`"etag"`), failure.ETag)
			require.Equal(t, 125*time.Millisecond, failure.RetryAfter)
			require.Equal(t, body, failure.Body)
			if tt.cancelled {
				require.ErrorIs(t, failure, context.Canceled)
				require.NotEmpty(t, failure.Message)
			} else {
				require.Equal(t, "Request rate is large.", failure.Message)
				require.Nil(t, failure.Unwrap(), "native message/body are not Go causes")
			}
			require.Same(t, diagnostics, failure.Diagnostics)
			require.Equal(t, uint32(2), failure.AttemptCount)
			require.Equal(t, tt.httpStatus, failure.Diagnostics.StatusCode)
			require.Equal(t, tt.wantSub, failure.Diagnostics.SubStatus)
		})
	}
}

func TestNativeWireErrorDiagnosticsOwnTheirSnapshot(t *testing.T) {
	const body = `{"code":"BadRequest","message":"synthetic mapping failure"}`
	server := newMetadataServer(t, func(w http.ResponseWriter, r *http.Request) bool {
		if !strings.Contains(r.URL.Path, "/docs/") {
			return false
		}
		w.Header().Set("x-ms-substatus", "4242")
		w.Header().Set("x-ms-activity-id", "synthetic-native-error")
		w.Header().Set("x-ms-request-charge", "3.5")
		w.Header().Set("x-ms-retry-after-ms", "17")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(body))
		if err != nil {
			t.Errorf("writing synthetic response: %v", err)
		}
		return true
	})
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	client, err := NewClientWithKey(ctx, server.URL, mustKeyCredential(t), nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(ctx, "db", "items")
	require.NoError(t, err)
	_, err = container.ReadItem(ctx, NewPartitionKeyString("synthetic"), "synthetic", nil)
	var failure *Error
	require.ErrorAs(t, err, &failure)
	require.NoError(t, client.Close())

	require.Equal(t, CodeBadRequest, failure.Code)
	require.Equal(t, 400, failure.StatusCode)
	require.Equal(t, 4242, failure.SubStatus)
	require.True(t, failure.FromWire)
	require.Equal(t, "synthetic-native-error", failure.ActivityID)
	require.Equal(t, 3.5, failure.RequestCharge)
	require.Equal(t, 17*time.Millisecond, failure.RetryAfter)
	require.Equal(t, []byte(body), failure.Body)
	require.NotEmpty(t, failure.Message)
	require.Nil(t, failure.Unwrap())
	requireEmulatorDiagnostics(t, failure.Diagnostics, 400, 4242, true)
	require.Equal(t, failure.AttemptCount, failure.Diagnostics.AttemptCount)
	require.Equal(t, failure.RequestCharge, failure.Diagnostics.TotalRequestCharge)
	cloned := cloneError(failure).(*Error)
	cloned.Diagnostics.Attempts[0].SubStatus = 0
	cloned.Diagnostics.RegionsContacted[0] = "changed"
	cloned.Body[0] = 'x'
	require.Equal(t, 4242, failure.Diagnostics.Attempts[0].SubStatus)
	require.NotEqual(t, "changed", failure.Diagnostics.RegionsContacted[0])
	require.Equal(t, []byte(body), failure.Body)
}

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
	diagnostics, err := copyDiagnostics(nil, DiagnosticsVerbosityDefault)
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

// TestTranslateCompletionOutcomePropagatesEveryDiagnosticsVerbosity walks every DiagnosticsVerbosity
// constant, verifying translateCompletionOutcome passes the exact value a caller gave it through
// to copyDiagnosticsForCompletion unchanged, regardless of the completion's outcome. The seam is
// stubbed because cosmos_diagnostics_t is an opaque native handle tests cannot construct; this
// proves the plumbing between the two functions; lazy rendering is covered separately below.
func TestTranslateCompletionOutcomePropagatesEveryDiagnosticsVerbosity(t *testing.T) {
	for _, outcome := range []syntheticOutcome{syntheticOutcomeOK, syntheticOutcomeCancelled, syntheticOutcomeError} {
		for _, verbosity := range []DiagnosticsVerbosity{
			DiagnosticsVerbosityDefault,
			DiagnosticsVerbositySummary,
			DiagnosticsVerbosityDetailed,
		} {
			t.Run(fmt.Sprintf("%d/%s", outcome, verbosity), func(t *testing.T) {
				restore, seen := stubCopyDiagnosticsForCompletionCapturing(&Diagnostics{}, nil)
				t.Cleanup(restore)

				syntheticCompletionResultWithVerbosity(outcome, http.StatusOK, false, verbosity)

				require.Equal(t, []DiagnosticsVerbosity{verbosity}, *seen)
			})
		}
	}
}

func TestTranslateCompletionOutcomeDefersDiagnosticsJSON(t *testing.T) {
	original := marshalDiagnosticsJSON
	var calls atomic.Int32
	marshalDiagnosticsJSON = func(value any) ([]byte, error) {
		calls.Add(1)
		return original(value)
	}
	t.Cleanup(func() { marshalDiagnosticsJSON = original })

	for _, verbosity := range []DiagnosticsVerbosity{DiagnosticsVerbositySummary, DiagnosticsVerbosityDetailed} {
		t.Run(verbosity.String(), func(t *testing.T) {
			fakeDiagnostics := sampleDiagnostics()
			stubDiagnosticsCopy(t, fakeDiagnostics, nil)

			result := syntheticCompletionResultWithVerbosity(syntheticOutcomeOK, http.StatusOK, false, verbosity)
			require.Zero(t, calls.Load(), "completion translation must not render JSON")

			rendered, err := result.response.Diagnostics.JSON()
			require.NoError(t, err)
			require.Equal(t, int32(1), calls.Load())
			require.Equal(t, verbosity == DiagnosticsVerbosityDetailed, strings.Contains(rendered, `"attempts"`))
			calls.Store(0)
		})
	}
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

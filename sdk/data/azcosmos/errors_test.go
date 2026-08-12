// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCodeForDriverStatus(t *testing.T) {
	for _, tt := range []struct {
		name      string
		status    driverStatus
		subStatus int
		want      Code
	}{
		{"bad request", driverStatusBadRequest, -1, CodeBadRequest},
		{"unauthorized", driverStatusUnauthorized, -1, CodeUnauthorized},
		{"forbidden", driverStatusForbidden, -1, CodeForbidden},
		{"not found", driverStatusNotFound, -1, CodeNotFound},
		{"request timeout", driverStatusTimeout, -1, CodeRequestTimeout},
		{"conflict", driverStatusConflict, -1, CodeConflict},
		{"gone", driverStatusGone, -1, CodeGone},
		{"precondition failed", driverStatusPreconditionFailed, -1, CodePreconditionFailed},
		{"throttled", driverStatusThrottled, -1, CodeThrottled},
		{"service unavailable", driverStatusServiceUnavailable, -1, CodeServiceUnavailable},
		{"service error", driverStatusServiceError, -1, CodeServiceError},

		// Failures the client produced carry no HTTP status, so classifying on one would have
		// collapsed every one of these into CodeUnknown.
		{"client error", driverStatusClientError, -1, CodeClientError},
		{"transport failure", driverStatusTransportFailure, -1, CodeTransportFailure},
		{"serialization failed", driverStatusSerializationFailed, -1, CodeSerializationFailed},
		{"authentication failed", driverStatusAuthenticationFailed, -1, CodeAuthenticationFailed},
		{"client operation timeout", driverStatusClientOperationTimeout, -1, CodeClientOperationTimeout},
		{"cancelled", driverStatusOperationCancelled, -1, CodeOperationCancelled},

		// Cosmos DB overloads 404 and 410 with sub-status 1002; classifying either of those on
		// status alone would tell the caller something categorically wrong.
		{"session token not satisfied is not a missing resource", driverStatusNotFound, 1002, CodeSessionUnavailable},
		{"partition key range gone is not a bare gone", driverStatusGone, 1002, CodePartitionKeyRangeGone},

		// An unrelated sub-status must not disturb the status-based classification.
		{"not found with unrelated sub-status", driverStatusNotFound, 1003, CodeNotFound},
		{"gone with unrelated sub-status", driverStatusGone, 1007, CodeGone},

		// 1002 only means something on the statuses that overload it.
		{"1002 on an unrelated status", driverStatusConflict, 1002, CodeConflict},

		// Codes the driver adds later must still classify by band rather than falling to
		// CodeUnknown, which is what the header requires of consumers.
		{"argument validation is client-side", driverStatus(1), -1, CodeClientError},
		{"invalid utf-8 is client-side", driverStatus(2), -1, CodeClientError},
		{"reserved auth and conversion band", driverStatus(1001), -1, CodeClientError},
		{"unmapped wire code is a service error", driverStatus(2451), -1, CodeServiceError},
		{"unmapped plumbing code is client-side", driverStatus(3999), -1, CodeClientError},
		{"invalid partition key", driverStatus(4004), -1, CodeClientError},
		{"queue shutdown", driverStatus(4011), -1, CodeClientError},

		// The warning band still delivers its result, so it must not be reported as a client
		// failure. 5001 is options ignored on a cache hit, where the driver is still returned.
		{"warning band is not a failure", driverStatus(5001), -1, CodeUnknown},

		{"success is not a failure", driverStatusSuccess, -1, CodeUnknown},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, codeForDriverStatus(tt.status, tt.subStatus))
		})
	}
}

// Each mapped status must produce its own code and no other, which catches both a mis-mapping and
// two constants having been given the same value.
func TestCodeForDriverStatusIsInjective(t *testing.T) {
	statuses := []driverStatus{
		driverStatusBadRequest, driverStatusUnauthorized, driverStatusForbidden,
		driverStatusNotFound, driverStatusTimeout, driverStatusConflict, driverStatusGone,
		driverStatusPreconditionFailed, driverStatusThrottled, driverStatusServiceUnavailable,
		driverStatusServiceError, driverStatusClientError, driverStatusTransportFailure,
		driverStatusSerializationFailed, driverStatusAuthenticationFailed,
		driverStatusClientOperationTimeout, driverStatusOperationCancelled,
	}

	seen := make(map[Code]driverStatus, len(statuses))
	for _, status := range statuses {
		code := codeForDriverStatus(status, -1)
		require.NotEqual(t, CodeUnknown, code, "status %d should be classified", status)

		previous, duplicated := seen[code]
		require.False(t, duplicated, "statuses %d and %d both classify as %q", previous, status, code)
		seen[code] = status
	}
}

// The driver reports an absent sub-status as -1, while the API documents zero, so the boundary has
// to translate. A raw -1 would otherwise read as a real sub-status in Error.
func TestNormalizeSubStatus(t *testing.T) {
	for _, tt := range []struct {
		name      string
		subStatus int32
		want      int
	}{
		{"absent", -1, 0},
		{"any other negative sentinel", -42, 0},
		{"already zero", 0, 0},
		{"present", 1002, 1002},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, normalizeSubStatus(tt.subStatus))
		})
	}
}

// An absent sub-status must not reach the message, whether it was normalized or copied verbatim.
func TestErrorOmitsAbsentSubStatus(t *testing.T) {
	normalized := (&Error{Code: CodeClientError, SubStatus: normalizeSubStatus(-1)}).Error()
	require.NotContains(t, normalized, "sub-status")

	unnormalized := (&Error{Code: CodeClientError, SubStatus: -1}).Error()
	require.NotContains(t, unnormalized, "sub-status")
	require.NotContains(t, unnormalized, "-1")

	withStatus := (&Error{Code: CodeNotFound, StatusCode: http.StatusNotFound, SubStatus: -1}).Error()
	require.Contains(t, withStatus, "status 404")
	require.NotContains(t, withStatus, "404/")
}

// A cancelled operation has to answer to the vocabulary Go callers already use, so that code
// written against context cancellation works without knowing about Cosmos codes.
func TestErrorUnwrapsCancellation(t *testing.T) {
	cancelled := fmt.Errorf("reading item: %w", &Error{Code: CodeOperationCancelled})
	require.ErrorIs(t, cancelled, context.Canceled)

	var cosmosErr *Error
	require.True(t, errors.As(cancelled, &cosmosErr))
	require.Equal(t, CodeOperationCancelled, cosmosErr.Code)

	// Every other failure is a Cosmos failure and nothing else.
	require.NotErrorIs(t, &Error{Code: CodeThrottled}, context.Canceled)
	require.NoError(t, (&Error{Code: CodeThrottled}).Unwrap())
}

func TestErrorAsRetrievesFields(t *testing.T) {
	// Wrapped, because callers see errors that operations have annotated on the way out.
	err := fmt.Errorf("reading item: %w", &Error{
		Code:          CodeThrottled,
		StatusCode:    http.StatusTooManyRequests,
		SubStatus:     3200,
		Message:       "Request rate is large",
		RequestCharge: 4.5,
		ActivityID:    "8fd3d1d1-5cbb-4a2a-9c4b-6a2b1f6c9d55",
		SessionToken:  "0:-1#42",
		ETag:          `"00000000-0000-0000-0000-000000000000"`,
		RetryAfter:    150 * time.Millisecond,
		Body:          []byte(`{"code":"TooManyRequests"}`),
		FromWire:      true,
	})

	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeThrottled, cosmosErr.Code)
	require.Equal(t, http.StatusTooManyRequests, cosmosErr.StatusCode)
	require.Equal(t, 3200, cosmosErr.SubStatus)
	require.Equal(t, "Request rate is large", cosmosErr.Message)
	require.Equal(t, "8fd3d1d1-5cbb-4a2a-9c4b-6a2b1f6c9d55", cosmosErr.ActivityID)
	require.Equal(t, SessionToken("0:-1#42"), cosmosErr.SessionToken)
	require.Equal(t, 150*time.Millisecond, cosmosErr.RetryAfter)
	require.True(t, cosmosErr.FromWire)

	// A throttled request is still billed, so the charge has to survive onto the error: it is the
	// only place a caller can account for the RUs a failed operation consumed.
	require.Equal(t, float64(4.5), cosmosErr.RequestCharge)
	require.Equal(t, []byte(`{"code":"TooManyRequests"}`), cosmosErr.Body)
}

// The preview returns errNotImplemented from every operation, so it has to satisfy the errors.As
// idiom the package documents; a bare errors.New would not.
func TestNotImplementedIsRetrievableAsError(t *testing.T) {
	err := fmt.Errorf("reading item: %w", error(errNotImplemented))

	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientError, cosmosErr.Code)
	require.False(t, cosmosErr.FromWire)
	require.Zero(t, cosmosErr.StatusCode)
}

func TestErrorMessage(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "code and status",
			err:  &Error{Code: CodeNotFound, StatusCode: http.StatusNotFound},
			want: "azcosmos: operation failed: NotFound: status 404",
		},
		{
			name: "sub-status is reported alongside status",
			err:  &Error{Code: CodeSessionUnavailable, StatusCode: http.StatusNotFound, SubStatus: 1002},
			want: "azcosmos: operation failed: SessionUnavailable: status 404/1002",
		},
		{
			name: "sub-status survives a missing status",
			err:  &Error{SubStatus: 1002},
			want: "azcosmos: operation failed: sub-status 1002",
		},
		{
			name: "retry-after is reported so throttling is diagnosable from logs",
			err: &Error{
				Code:       CodeThrottled,
				StatusCode: http.StatusTooManyRequests,
				RetryAfter: 150 * time.Millisecond,
				ActivityID: "abc",
			},
			want: "azcosmos: operation failed: Throttled: status 429 (retry after 150ms) (activity id: abc)",
		},
		{
			name: "message and activity id",
			err:  &Error{Code: CodeConflict, StatusCode: http.StatusConflict, Message: "Resource already exists", ActivityID: "abc"},
			want: "azcosmos: operation failed: Conflict: status 409: Resource already exists (activity id: abc)",
		},
		{
			name: "client-side failure carries no status",
			err:  &Error{Code: CodeClientError, Message: "connection refused"},
			want: "azcosmos: operation failed: ClientError: connection refused",
		},
		{
			name: "unclassified failure",
			err:  &Error{StatusCode: http.StatusInternalServerError},
			want: "azcosmos: operation failed: status 500",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.err.Error())
		})
	}
}

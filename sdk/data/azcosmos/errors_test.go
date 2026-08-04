// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCodeForStatus(t *testing.T) {
	for _, tt := range []struct {
		name       string
		statusCode int
		subStatus  int
		want       Code
	}{
		{"bad request", http.StatusBadRequest, 0, CodeBadRequest},
		{"unauthorized", http.StatusUnauthorized, 0, CodeUnauthorized},
		{"forbidden", http.StatusForbidden, 0, CodeForbidden},
		{"not found", http.StatusNotFound, 0, CodeNotFound},
		{"request timeout", http.StatusRequestTimeout, 0, CodeRequestTimeout},
		{"conflict", http.StatusConflict, 0, CodeConflict},
		{"gone", http.StatusGone, 0, CodeGone},
		{"precondition failed", http.StatusPreconditionFailed, 0, CodePreconditionFailed},
		{"throttled", http.StatusTooManyRequests, 0, CodeThrottled},
		{"service unavailable", http.StatusServiceUnavailable, 0, CodeServiceUnavailable},

		// Cosmos DB overloads 404 and 410 with sub-status 1002; classifying either of those on
		// status alone would tell the caller something categorically wrong.
		{"session token not satisfied is not a missing resource", http.StatusNotFound, 1002, CodeSessionUnavailable},
		{"partition key range gone is not a bare gone", http.StatusGone, 1002, CodePartitionKeyRangeGone},

		// An unrelated sub-status must not disturb the status-based classification.
		{"not found with unrelated sub-status", http.StatusNotFound, 1003, CodeNotFound},
		{"gone with unrelated sub-status", http.StatusGone, 1007, CodeGone},

		// 1002 only means something on the statuses that overload it.
		{"1002 on an unrelated status", http.StatusConflict, 1002, CodeConflict},

		{"unmapped status", http.StatusInternalServerError, 0, CodeUnknown},
		{"no status at all", 0, 0, CodeUnknown},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, codeForStatus(tt.statusCode, tt.subStatus))
		})
	}
}

// Each mapped status must produce its own code and no other, which catches both a mis-mapping and
// two constants having been given the same value.
func TestCodeForStatusIsInjective(t *testing.T) {
	statuses := []int{
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound,
		http.StatusRequestTimeout, http.StatusConflict, http.StatusGone,
		http.StatusPreconditionFailed, http.StatusTooManyRequests, http.StatusServiceUnavailable,
	}

	seen := make(map[Code]int, len(statuses))
	for _, status := range statuses {
		code := codeForStatus(status, 0)
		require.NotEqual(t, CodeUnknown, code, "status %d should be classified", status)

		previous, duplicated := seen[code]
		require.False(t, duplicated, "statuses %d and %d both classify as %q", previous, status, code)
		seen[code] = status
	}
}

func TestErrorAsRetrievesFields(t *testing.T) {
	// Wrapped, because callers see errors that operations have annotated on the way out.
	err := fmt.Errorf("reading item: %w", &Error{
		Code:         CodeThrottled,
		StatusCode:   http.StatusTooManyRequests,
		SubStatus:    3200,
		Message:      "Request rate is large",
		ActivityID:   "8fd3d1d1-5cbb-4a2a-9c4b-6a2b1f6c9d55",
		SessionToken: "0:-1#42",
		ETag:         `"00000000-0000-0000-0000-000000000000"`,
		RetryAfter:   150 * time.Millisecond,
		FromWire:     true,
	})

	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeThrottled, cosmosErr.Code)
	require.Equal(t, http.StatusTooManyRequests, cosmosErr.StatusCode)
	require.Equal(t, 3200, cosmosErr.SubStatus)
	require.Equal(t, "Request rate is large", cosmosErr.Message)
	require.Equal(t, "8fd3d1d1-5cbb-4a2a-9c4b-6a2b1f6c9d55", cosmosErr.ActivityID)
	require.Equal(t, "0:-1#42", cosmosErr.SessionToken)
	require.Equal(t, 150*time.Millisecond, cosmosErr.RetryAfter)
	require.True(t, cosmosErr.FromWire)
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

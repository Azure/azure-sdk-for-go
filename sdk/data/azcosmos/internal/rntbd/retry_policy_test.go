// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetryPolicy_ShouldRetryWithGoneException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{RntbdError: RntbdError{StatusCode: StatusGone}}

	result1 := policy.ShouldRetry(goneErr, retryCtx)
	require.True(t, result1.ShouldRetry)
	require.True(t, result1.ForceRefreshAddressCache)
	require.Equal(t, 1, result1.AttemptNumber)
	require.Equal(t, time.Duration(0), result1.BackoffTime)

	result2 := policy.ShouldRetry(goneErr, retryCtx)
	require.True(t, result2.ShouldRetry)
	require.True(t, result2.ForceRefreshAddressCache)
	require.Equal(t, 2, result2.AttemptNumber)
	require.Equal(t, 1*time.Second, result2.BackoffTime)

	result3 := policy.ShouldRetry(goneErr, retryCtx)
	require.True(t, result3.ShouldRetry)
	require.True(t, result3.ForceRefreshAddressCache)
	require.Equal(t, 3, result3.AttemptNumber)
	require.Equal(t, 2*time.Second, result3.BackoffTime)

	result4 := policy.ShouldRetry(goneErr, retryCtx)
	require.True(t, result4.ShouldRetry)
	require.True(t, result4.ForceRefreshAddressCache)
	require.Equal(t, 4, result4.AttemptNumber)
	require.Equal(t, 4*time.Second, result4.BackoffTime)
}

func TestRetryPolicy_ShouldRetryWithPartitionIsMigratingException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	migratingErr := &PartitionIsMigratingError{RntbdError: RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusCompletingPartitionMigrate}}

	result := policy.ShouldRetry(migratingErr, retryCtx)

	require.True(t, result.ShouldRetry)
	require.True(t, retryCtx.ForceCollectionRoutingMapRefresh)
	require.True(t, result.ForceRefreshAddressCache)
}

func TestRetryPolicy_ShouldRetryWithInvalidPartitionException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()

	invalidErr := &InvalidPartitionError{
		RntbdError:                      RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusNameCacheIsStale},
		IsBasedOn410ResponseFromService: true,
	}

	result := policy.ShouldRetry(invalidErr, retryCtx)
	require.True(t, result.ShouldRetry)
	require.True(t, result.ForceRefreshAddressCache)
}

func TestRetryPolicy_ShouldRetryWithPartitionKeyRangeIsSplittingException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	retryCtx.ResolvedPartitionKeyRange = "some-range"
	retryCtx.QuorumSelectedLSN = 100

	splittingErr := &PartitionKeyRangeIsSplittingError{RntbdError: RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusCompletingSplit}}

	result := policy.ShouldRetry(splittingErr, retryCtx)

	require.True(t, result.ShouldRetry)
	require.True(t, retryCtx.ForcePartitionKeyRangeRefresh)
	require.Nil(t, retryCtx.ResolvedPartitionKeyRange)
	require.Equal(t, int64(-1), retryCtx.QuorumSelectedLSN)
	require.False(t, result.ForceRefreshAddressCache)
}

func TestRetryPolicy_ShouldRetryWithGenericException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	badRequestErr := &BadRequestError{RntbdError: RntbdError{StatusCode: StatusBadRequest}}

	result := policy.ShouldRetry(badRequestErr, retryCtx)

	require.False(t, result.ShouldRetry)
}

func TestRetryPolicy_BackoffProgression(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	goneErr := &GoneError{RntbdError: RntbdError{StatusCode: StatusGone}}

	expectedBackoffs := []time.Duration{0, 1 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 15 * time.Second, 15 * time.Second}

	for i, expected := range expectedBackoffs {
		result := policy.ShouldRetry(goneErr, NewRetryContext())
		require.True(t, result.ShouldRetry, "attempt %d should allow retry", i+1)
		require.Equal(t, expected, result.BackoffTime, "attempt %d backoff mismatch", i+1)
	}
}

func TestRetryPolicy_RetryWithExceptionStored(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryWithErr := &RetryWithError{RntbdError: RntbdError{StatusCode: StatusRetryWith}}

	result := policy.ShouldRetry(retryWithErr, NewRetryContext())
	require.True(t, result.ShouldRetry)
	// RetryWith uses ms-scale backoff, not address cache refresh
	require.False(t, result.ForceRefreshAddressCache)

	require.Equal(t, retryWithErr, policy.lastRetryWithException)
}

func TestRetryPolicy_RetryWithBackoffProgression(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryWithErr := &RetryWithError{RntbdError: RntbdError{StatusCode: StatusRetryWith}}

	// First attempt: no backoff
	result1 := policy.ShouldRetry(retryWithErr, NewRetryContext())
	require.True(t, result1.ShouldRetry)
	require.Equal(t, time.Duration(0), result1.BackoffTime)

	// Second attempt: ~10ms (+random salt 0-4ms)
	result2 := policy.ShouldRetry(retryWithErr, NewRetryContext())
	require.True(t, result2.ShouldRetry)
	require.GreaterOrEqual(t, result2.BackoffTime, 10*time.Millisecond)
	require.LessOrEqual(t, result2.BackoffTime, 15*time.Millisecond)

	// Third attempt: ~20ms (+random salt)
	result3 := policy.ShouldRetry(retryWithErr, NewRetryContext())
	require.True(t, result3.ShouldRetry)
	require.GreaterOrEqual(t, result3.BackoffTime, 20*time.Millisecond)
	require.LessOrEqual(t, result3.BackoffTime, 25*time.Millisecond)
}

func TestRetryPolicy_LeaseNotFoundException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	leaseErr := &LeaseNotFoundError{RntbdError: RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusLeaseNotFound}}

	result := policy.ShouldRetry(leaseErr, NewRetryContext())
	require.False(t, result.ShouldRetry)
	require.NotNil(t, result.Exception)

	var serviceErr *ServiceUnavailableError
	require.ErrorAs(t, result.Exception, &serviceErr)
	require.Equal(t, int64(SubStatusLeaseNotFound), serviceErr.SubStatusCode)
}

func TestRetryPolicy_PartitionKeyRangeGoneException(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)
	pkRangeGoneErr := &PartitionKeyRangeGoneError{RntbdError: RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusPartitionKeyRangeGone}}

	result := policy.ShouldRetry(pkRangeGoneErr, NewRetryContext())

	require.True(t, result.ShouldRetry)
	require.True(t, result.ForceRefreshAddressCache)
}

func TestRetryPolicy_DefaultWaitTime(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(0)
	require.Equal(t, DefaultWaitTimeInSeconds, policy.waitTimeInSeconds)
}

func TestRetryPolicy_NoRetryForUnknownError(t *testing.T) {
	policy := NewGoneAndRetryWithRetryPolicy(30)

	unknownErr := &ConflictError{RntbdError: RntbdError{StatusCode: StatusConflict}}
	result := policy.ShouldRetry(unknownErr, NewRetryContext())
	require.False(t, result.ShouldRetry)

	forbiddenErr := &ForbiddenError{RntbdError: RntbdError{StatusCode: StatusForbidden}}
	result = policy.ShouldRetry(forbiddenErr, NewRetryContext())
	require.False(t, result.ShouldRetry)

	notFoundErr := &NotFoundError{RntbdError: RntbdError{StatusCode: StatusNotFound}}
	result = policy.ShouldRetry(notFoundErr, NewRetryContext())
	require.False(t, result.ShouldRetry)
}

// =============================================================================
// Non-Idempotent Write Safety Tests (Stage 9)
// =============================================================================

func TestRetryPolicy_ReadOnlyOperationsAlwaysRetry(t *testing.T) {
	// Read-only operations should always be allowed to retry regardless of send state
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: false, // NOT from service - transport error
	}

	readOnlyOps := []OperationType{
		OperationRead,
		OperationReadFeed,
		OperationQuery,
		OperationSQLQuery,
		OperationHead,
		OperationHeadFeed,
	}

	for _, opType := range readOnlyOps {
		t.Run(opType.String(), func(t *testing.T) {
			req := &ServiceRequest{OperationType: opType}
			// Simulate that sending has started
			req.MarkSendingRequestStarted(time.Now().UnixNano())

			result := policy.ShouldRetryWithRequest(goneErr, retryCtx, req)
			require.True(t, result.ShouldRetry, "read-only operation %s should always retry", opType.String())
		})
	}
}

func TestRetryPolicy_WriteBeforeSendStartedRetries(t *testing.T) {
	// Write operations BEFORE send started should retry (no network write yet)
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: false, // NOT from service - transport error
	}

	writeOps := []OperationType{
		OperationCreate,
		OperationReplace,
		OperationDelete,
		OperationUpsert,
		OperationPatch,
		OperationBatch,
	}

	for _, opType := range writeOps {
		t.Run(opType.String(), func(t *testing.T) {
			req := &ServiceRequest{OperationType: opType}
			// Do NOT mark sending started - send hasn't begun

			result := policy.ShouldRetryWithRequest(goneErr, retryCtx, req)
			require.True(t, result.ShouldRetry, "write operation %s should retry before send started", opType.String())
		})
	}
}

func TestRetryPolicy_WriteAfterSendStartedWithService410Retries(t *testing.T) {
	// Write operations AFTER send started WITH service 410 should retry
	// (server confirmed it didn't process the request)
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: true, // From service - safe to retry
	}

	writeOps := []OperationType{
		OperationCreate,
		OperationReplace,
		OperationDelete,
		OperationUpsert,
		OperationPatch,
		OperationBatch,
	}

	for _, opType := range writeOps {
		t.Run(opType.String(), func(t *testing.T) {
			req := &ServiceRequest{OperationType: opType}
			req.MarkSendingRequestStarted(time.Now().UnixNano())

			result := policy.ShouldRetryWithRequest(goneErr, retryCtx, req)
			require.True(t, result.ShouldRetry, "write operation %s should retry with service 410", opType.String())
		})
	}
}

func TestRetryPolicy_WriteAfterSendStartedWithoutService410ReturnsTimeout(t *testing.T) {
	// Write operations AFTER send started WITHOUT service 410 should NOT retry
	// (could cause duplicate writes)
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: false, // Transport error - NOT from service
	}

	writeOps := []OperationType{
		OperationCreate,
		OperationReplace,
		OperationDelete,
		OperationUpsert,
		OperationPatch,
		OperationBatch,
	}

	for _, opType := range writeOps {
		t.Run(opType.String(), func(t *testing.T) {
			req := &ServiceRequest{OperationType: opType}
			req.MarkSendingRequestStarted(time.Now().UnixNano())

			result := policy.ShouldRetryWithRequest(goneErr, retryCtx, req)
			require.False(t, result.ShouldRetry, "write operation %s should NOT retry without service 410", opType.String())
			require.NotNil(t, result.Exception)

			var timeoutErr *RequestTimeoutError
			require.ErrorAs(t, result.Exception, &timeoutErr)
			require.Equal(t, int32(StatusRequestTimeout), timeoutErr.StatusCode)
		})
	}
}

func TestRetryPolicy_NonIdempotentWriteRetriesEnabledFlag(t *testing.T) {
	// When nonIdempotentWriteRetriesEnabled is true, even unsafe writes should retry
	policy := NewGoneAndRetryWithRetryPolicyWithWriteRetries(30, true)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: false, // Transport error - normally unsafe
	}

	req := &ServiceRequest{OperationType: OperationCreate}
	req.MarkSendingRequestStarted(time.Now().UnixNano())

	result := policy.ShouldRetryWithRequest(goneErr, retryCtx, req)
	require.True(t, result.ShouldRetry, "should retry when nonIdempotentWriteRetriesEnabled is true")
}

func TestRetryPolicy_NilRequestAllowsRetry(t *testing.T) {
	// Nil request should allow retry (backward compatibility)
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	goneErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: false,
	}

	// Using nil request via ShouldRetry (not ShouldRetryWithRequest)
	result := policy.ShouldRetry(goneErr, retryCtx)
	require.True(t, result.ShouldRetry, "nil request should allow retry for backward compatibility")
}

func TestRetryPolicy_AllGoneErrorTypesCheckWriteSafety(t *testing.T) {
	// All Gone error types should check write safety
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()

	// Create a write request that has started sending
	writeReq := &ServiceRequest{OperationType: OperationCreate}
	writeReq.MarkSendingRequestStarted(time.Now().UnixNano())

	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "GoneError",
			err:  &GoneError{RntbdError: RntbdError{StatusCode: StatusGone}, IsBasedOn410ResponseFromService: false},
		},
		{
			name: "InvalidPartitionError",
			err:  &InvalidPartitionError{RntbdError: RntbdError{StatusCode: StatusGone}, IsBasedOn410ResponseFromService: false},
		},
		{
			name: "PartitionIsMigratingError",
			err:  &PartitionIsMigratingError{RntbdError: RntbdError{StatusCode: StatusGone}, IsBasedOn410ResponseFromService: false},
		},
		{
			name: "PartitionKeyRangeIsSplittingError",
			err:  &PartitionKeyRangeIsSplittingError{RntbdError: RntbdError{StatusCode: StatusGone}, IsBasedOn410ResponseFromService: false},
		},
		{
			name: "PartitionKeyRangeGoneError",
			err:  &PartitionKeyRangeGoneError{RntbdError: RntbdError{StatusCode: StatusGone}, IsBasedOn410ResponseFromService: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := policy.ShouldRetryWithRequest(tc.err, retryCtx, writeReq)
			require.False(t, result.ShouldRetry, "%s should not retry unsafe write", tc.name)
			require.NotNil(t, result.Exception)

			var timeoutErr *RequestTimeoutError
			require.ErrorAs(t, result.Exception, &timeoutErr)
		})
	}
}

func TestRetryPolicy_ServiceRequestSendTracking(t *testing.T) {
	// Test ServiceRequest send tracking methods
	req := &ServiceRequest{OperationType: OperationCreate}

	require.False(t, req.HasSendingRequestStarted(), "should not have started initially")
	require.Equal(t, int64(0), req.GetSentTime(), "sent time should be 0 initially")

	sentTime := time.Now().UnixNano()
	req.MarkSendingRequestStarted(sentTime)

	require.True(t, req.HasSendingRequestStarted(), "should have started after marking")
	require.Equal(t, sentTime, req.GetSentTime(), "sent time should match")
}

func TestRetryPolicy_ServiceRequestIsReadOnly(t *testing.T) {
	// Test IsReadOnly classification
	readOps := []OperationType{
		OperationRead,
		OperationReadFeed,
		OperationQuery,
		OperationSQLQuery,
		OperationHead,
		OperationHeadFeed,
	}

	writeOps := []OperationType{
		OperationCreate,
		OperationReplace,
		OperationDelete,
		OperationUpsert,
		OperationPatch,
		OperationBatch,
		OperationExecuteJavaScript,
	}

	for _, opType := range readOps {
		req := &ServiceRequest{OperationType: opType}
		require.True(t, req.IsReadOnly(), "operation %s should be read-only", opType.String())
	}

	for _, opType := range writeOps {
		req := &ServiceRequest{OperationType: opType}
		require.False(t, req.IsReadOnly(), "operation %s should NOT be read-only", opType.String())
	}
}

func TestRetryPolicy_PartitionIsMigratingPreservesWriteSafety(t *testing.T) {
	// PartitionIsMigratingError should also check write safety
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	migratingErr := &PartitionIsMigratingError{
		RntbdError:                      RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusCompletingPartitionMigrate},
		IsBasedOn410ResponseFromService: false,
	}

	// Read operation - should retry
	readReq := &ServiceRequest{OperationType: OperationRead}
	readReq.MarkSendingRequestStarted(time.Now().UnixNano())
	result := policy.ShouldRetryWithRequest(migratingErr, retryCtx, readReq)
	require.True(t, result.ShouldRetry)
	require.True(t, result.ForceRefreshAddressCache)
	require.True(t, retryCtx.ForceCollectionRoutingMapRefresh)

	// Write operation after send without service 410 - should NOT retry
	retryCtx2 := NewRetryContext()
	writeReq := &ServiceRequest{OperationType: OperationCreate}
	writeReq.MarkSendingRequestStarted(time.Now().UnixNano())
	result = policy.ShouldRetryWithRequest(migratingErr, retryCtx2, writeReq)
	require.False(t, result.ShouldRetry)
}

func TestRetryPolicy_PartitionKeyRangeIsSplittingPreservesWriteSafety(t *testing.T) {
	// PartitionKeyRangeIsSplittingError should also check write safety
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()
	retryCtx.ResolvedPartitionKeyRange = "some-range"
	retryCtx.QuorumSelectedLSN = 100

	splittingErr := &PartitionKeyRangeIsSplittingError{
		RntbdError:                      RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusCompletingSplit},
		IsBasedOn410ResponseFromService: false,
	}

	// Read operation - should retry
	readReq := &ServiceRequest{OperationType: OperationRead}
	readReq.MarkSendingRequestStarted(time.Now().UnixNano())
	result := policy.ShouldRetryWithRequest(splittingErr, retryCtx, readReq)
	require.True(t, result.ShouldRetry)
	require.True(t, retryCtx.ForcePartitionKeyRangeRefresh)
	require.Nil(t, retryCtx.ResolvedPartitionKeyRange)
	require.Equal(t, int64(-1), retryCtx.QuorumSelectedLSN)

	// Write operation after send without service 410 - should NOT retry
	retryCtx2 := NewRetryContext()
	writeReq := &ServiceRequest{OperationType: OperationCreate}
	writeReq.MarkSendingRequestStarted(time.Now().UnixNano())
	result = policy.ShouldRetryWithRequest(splittingErr, retryCtx2, writeReq)
	require.False(t, result.ShouldRetry)
}

func TestRetryPolicy_ErrorFromServiceVsTransport(t *testing.T) {
	// Test the distinction between service 410 and transport errors
	policy := NewGoneAndRetryWithRetryPolicy(30)
	writeReq := &ServiceRequest{OperationType: OperationCreate}
	writeReq.MarkSendingRequestStarted(time.Now().UnixNano())

	// Service 410 - safe to retry
	serviceErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: true,
	}
	result := policy.ShouldRetryWithRequest(serviceErr, NewRetryContext(), writeReq)
	require.True(t, result.ShouldRetry, "service 410 should allow retry")

	// Transport error - NOT safe to retry
	transportErr := &GoneError{
		RntbdError:                      RntbdError{StatusCode: StatusGone},
		IsBasedOn410ResponseFromService: false,
	}
	result = policy.ShouldRetryWithRequest(transportErr, NewRetryContext(), writeReq)
	require.False(t, result.ShouldRetry, "transport error should NOT allow retry")
}

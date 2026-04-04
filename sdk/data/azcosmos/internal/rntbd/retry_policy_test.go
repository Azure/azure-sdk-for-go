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
	// New SDK: InvalidPartitionError is no longer retried by GoneRetryPolicy
	policy := NewGoneAndRetryWithRetryPolicy(30)
	retryCtx := NewRetryContext()

	invalidErr := &InvalidPartitionError{RntbdError: RntbdError{StatusCode: StatusGone, SubStatusCode: SubStatusNameCacheIsStale}}

	result := policy.ShouldRetry(invalidErr, retryCtx)
	require.False(t, result.ShouldRetry)
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

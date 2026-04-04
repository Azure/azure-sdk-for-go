// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"math/rand"
	"time"
)

// GoneRetryPolicy constants (second-scale backoff)
const (
	DefaultWaitTimeInSeconds    = 30
	MaximumBackoffTimeInSeconds = 15
	InitialBackoffTimeSeconds   = 1
	BackoffMultiplier           = 2
)

// RetryWithRetryPolicy constants (millisecond-scale backoff)
const (
	retryWithDefaultWaitTimeInSeconds = 30
	retryWithMaximumBackoffTimeInMs   = 1000
	retryWithInitialBackoffTimeMs     = 10
	retryWithBackoffMultiplier        = 2
	retryWithRandomSaltInMs           = 5
)

type ShouldRetryResult struct {
	ShouldRetry              bool
	BackoffTime              time.Duration
	Timeout                  time.Duration
	AttemptNumber            int
	ForceRefreshAddressCache bool
	Exception                error
}

func NoRetry() ShouldRetryResult {
	return ShouldRetryResult{ShouldRetry: false}
}

func NoRetryWithError(err error) ShouldRetryResult {
	return ShouldRetryResult{ShouldRetry: false, Exception: err}
}

type RetryContext struct {
	QuorumSelectedLSN                int64
	GlobalCommittedSelectedLSN       int64
	ResolvedPartitionKeyRange        interface{}
	QuorumSelectedStoreResponse      interface{}
	ForceCollectionRoutingMapRefresh bool
	ForcePartitionKeyRangeRefresh    bool
	ForceNameCacheRefresh            bool
}

func NewRetryContext() *RetryContext {
	return &RetryContext{
		QuorumSelectedLSN:          -1,
		GlobalCommittedSelectedLSN: -1,
	}
}

// -----------------------------------------------------------------------------
// RetryWithRetryPolicy — millisecond-scale backoff for RetryWith (449) errors
// -----------------------------------------------------------------------------

type retryWithRetryPolicy struct {
	waitTimeInSeconds int
	attemptCount      int
	currentBackoffMs  int
	startTime         time.Time
}

func newRetryWithRetryPolicy(waitTimeInSeconds int) *retryWithRetryPolicy {
	if waitTimeInSeconds <= 0 {
		waitTimeInSeconds = retryWithDefaultWaitTimeInSeconds
	}
	return &retryWithRetryPolicy{
		waitTimeInSeconds: waitTimeInSeconds,
		attemptCount:      1,
		currentBackoffMs:  retryWithInitialBackoffTimeMs,
		startTime:         time.Now(),
	}
}

func (p *retryWithRetryPolicy) shouldRetry(err error) (ShouldRetryResult, bool) {
	if _, ok := err.(*RetryWithError); !ok {
		// Not a RetryWith error — signal caller to try the next policy
		return ShouldRetryResult{}, false
	}

	currentAttempt := p.attemptCount

	var backoffTime time.Duration
	if p.attemptCount > 1 {
		elapsedMs := int(time.Since(p.startTime).Milliseconds())
		remainingMs := p.waitTimeInSeconds*1000 - elapsedMs

		if remainingMs <= 0 {
			return NoRetryWithError(&ServiceUnavailableError{
				RntbdError: RntbdError{
					StatusCode: StatusServiceUnavailable,
					Message:    "operation timed out (RetryWith)",
				},
			}), true
		}

		backoffMs := p.currentBackoffMs + rand.Intn(retryWithRandomSaltInMs)
		if backoffMs > remainingMs {
			backoffMs = remainingMs
		}
		if backoffMs > retryWithMaximumBackoffTimeInMs {
			backoffMs = retryWithMaximumBackoffTimeInMs
		}
		backoffTime = time.Duration(backoffMs) * time.Millisecond

		newBackoff := p.currentBackoffMs * retryWithBackoffMultiplier
		if newBackoff < retryWithInitialBackoffTimeMs {
			newBackoff = retryWithInitialBackoffTimeMs
		}
		if newBackoff > retryWithMaximumBackoffTimeInMs {
			newBackoff = retryWithMaximumBackoffTimeInMs
		}
		p.currentBackoffMs = newBackoff
	}
	p.attemptCount++

	elapsedMs := int(time.Since(p.startTime).Milliseconds())
	remainingMs := p.waitTimeInSeconds*1000 - elapsedMs - int(backoffTime.Milliseconds())
	var timeout time.Duration
	if remainingMs > 0 {
		timeout = time.Duration(remainingMs) * time.Millisecond
	} else {
		timeout = time.Duration(retryWithMaximumBackoffTimeInMs) * time.Millisecond
	}

	return ShouldRetryResult{
		ShouldRetry:              true,
		BackoffTime:              backoffTime,
		Timeout:                  timeout,
		AttemptNumber:            currentAttempt,
		ForceRefreshAddressCache: false,
	}, true
}

// -----------------------------------------------------------------------------
// GoneAndRetryWithRetryPolicy — composite: delegates to retryWithPolicy first,
// then gonePolicy. Matches Java GoneAndRetryWithRetryPolicy.
// -----------------------------------------------------------------------------

type GoneAndRetryWithRetryPolicy struct {
	waitTimeInSeconds      int
	attemptCount           int
	currentBackoffSeconds  int
	startTime              time.Time
	lastRetryWithException *RetryWithError
	retryWithPolicy        *retryWithRetryPolicy
}

func NewGoneAndRetryWithRetryPolicy(waitTimeInSeconds int) *GoneAndRetryWithRetryPolicy {
	if waitTimeInSeconds <= 0 {
		waitTimeInSeconds = DefaultWaitTimeInSeconds
	}
	return &GoneAndRetryWithRetryPolicy{
		waitTimeInSeconds:     waitTimeInSeconds,
		attemptCount:          1,
		currentBackoffSeconds: InitialBackoffTimeSeconds,
		startTime:             time.Now(),
		retryWithPolicy:       newRetryWithRetryPolicy(waitTimeInSeconds),
	}
}

func (p *GoneAndRetryWithRetryPolicy) ShouldRetry(err error, retryCtx *RetryContext) ShouldRetryResult {
	if retryCtx == nil {
		retryCtx = NewRetryContext()
	}

	// First, let RetryWithRetryPolicy handle RetryWith errors (ms-scale backoff)
	if retryWithErr, ok := err.(*RetryWithError); ok {
		p.lastRetryWithException = retryWithErr
		result, handled := p.retryWithPolicy.shouldRetry(err)
		if handled {
			return result
		}
	}

	// GoneRetryPolicy handles the rest
	switch err.(type) {
	case *GoneError:
		return p.handleGoneException(retryCtx)
	case *PartitionIsMigratingError:
		return p.handlePartitionIsMigratingException(retryCtx)
	case *PartitionKeyRangeIsSplittingError:
		return p.handlePartitionKeyRangeIsSplittingException(retryCtx)
	case *PartitionKeyRangeGoneError:
		return p.handleGoneException(retryCtx)
	case *LeaseNotFoundError:
		return p.handleLeaseNotFoundException()
	default:
		return NoRetry()
	}
}

func (p *GoneAndRetryWithRetryPolicy) handleGoneException(retryCtx *RetryContext) ShouldRetryResult {
	currentRetryAttemptCount := p.attemptCount

	var backoffTime time.Duration
	if p.attemptCount > 1 {
		elapsedSeconds := int(time.Since(p.startTime).Seconds())
		remainingSeconds := p.waitTimeInSeconds - elapsedSeconds

		if remainingSeconds <= 0 {
			if p.lastRetryWithException != nil {
				return NoRetryWithError(p.lastRetryWithException)
			}
			return NoRetryWithError(&ServiceUnavailableError{
				RntbdError: RntbdError{
					StatusCode: StatusServiceUnavailable,
					Message:    "operation timed out",
				},
			})
		}

		backoffSeconds := p.currentBackoffSeconds
		if backoffSeconds > remainingSeconds {
			backoffSeconds = remainingSeconds
		}
		if backoffSeconds > MaximumBackoffTimeInSeconds {
			backoffSeconds = MaximumBackoffTimeInSeconds
		}
		backoffTime = time.Duration(backoffSeconds) * time.Second
		p.currentBackoffSeconds *= BackoffMultiplier
	}
	p.attemptCount++

	elapsedSeconds := int(time.Since(p.startTime).Seconds())
	remainingSeconds := p.waitTimeInSeconds - elapsedSeconds - int(backoffTime.Seconds())
	var timeout time.Duration
	if remainingSeconds > 0 {
		timeout = time.Duration(remainingSeconds) * time.Second
	} else {
		timeout = time.Duration(MaximumBackoffTimeInSeconds) * time.Second
	}

	return ShouldRetryResult{
		ShouldRetry:              true,
		BackoffTime:              backoffTime,
		Timeout:                  timeout,
		AttemptNumber:            currentRetryAttemptCount,
		ForceRefreshAddressCache: true,
	}
}

func (p *GoneAndRetryWithRetryPolicy) handlePartitionIsMigratingException(retryCtx *RetryContext) ShouldRetryResult {
	retryCtx.ForceCollectionRoutingMapRefresh = true

	result := p.handleGoneException(retryCtx)
	result.ForceRefreshAddressCache = true
	return result
}

func (p *GoneAndRetryWithRetryPolicy) handlePartitionKeyRangeIsSplittingException(retryCtx *RetryContext) ShouldRetryResult {
	retryCtx.ResolvedPartitionKeyRange = nil
	retryCtx.QuorumSelectedLSN = -1
	retryCtx.ForcePartitionKeyRangeRefresh = true

	result := p.handleGoneException(retryCtx)
	result.ForceRefreshAddressCache = false
	return result
}

func (p *GoneAndRetryWithRetryPolicy) handleLeaseNotFoundException() ShouldRetryResult {
	return NoRetryWithError(&ServiceUnavailableError{
		RntbdError: RntbdError{
			StatusCode:    StatusServiceUnavailable,
			SubStatusCode: SubStatusLeaseNotFound,
			Message:       "lease not found",
		},
	})
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"time"
)

// TimeoutHelper tracks elapsed time and remaining time for request timeouts.
// It is used to enforce request-level timeouts across retry attempts.
//
// This is a direct port of the Java SDK's TimeoutHelper class from:
// azure-cosmosdb-java/commons/src/main/java/com/microsoft/azure/cosmosdb/internal/directconnectivity/TimeoutHelper.java
type TimeoutHelper struct {
	startTime time.Time
	timeout   time.Duration
}

// NewTimeoutHelper creates a new TimeoutHelper with the specified timeout duration.
// The start time is recorded at creation.
func NewTimeoutHelper(timeout time.Duration) *TimeoutHelper {
	return &TimeoutHelper{
		startTime: time.Now(),
		timeout:   timeout,
	}
}

// IsElapsed returns true if the timeout duration has elapsed since creation.
func (t *TimeoutHelper) IsElapsed() bool {
	return time.Since(t.startTime) >= t.timeout
}

// GetRemainingTime returns the remaining time until timeout.
// This can return a negative duration if the timeout has already elapsed.
func (t *TimeoutHelper) GetRemainingTime() time.Duration {
	return t.timeout - time.Since(t.startTime)
}

// GetElapsedTime returns the time elapsed since creation.
func (t *TimeoutHelper) GetElapsedTime() time.Duration {
	return time.Since(t.startTime)
}

// ThrowTimeoutIfElapsed returns a RequestTimeoutError if the timeout has elapsed.
// Returns nil if the timeout has not elapsed.
func (t *TimeoutHelper) ThrowTimeoutIfElapsed() error {
	if t.IsElapsed() {
		return &RequestTimeoutError{
			RntbdError: RntbdError{
				StatusCode: StatusRequestTimeout,
				Message:    "request timed out",
			},
		}
	}
	return nil
}

// ThrowGoneIfElapsed returns a GoneError if the timeout has elapsed.
// Returns nil if the timeout has not elapsed.
func (t *TimeoutHelper) ThrowGoneIfElapsed() error {
	if t.IsElapsed() {
		return &GoneError{
			RntbdError: RntbdError{
				StatusCode:    StatusGone,
				SubStatusCode: SubStatusTimeoutGenerated410,
				Message:       "request timed out (gone)",
			},
		}
	}
	return nil
}

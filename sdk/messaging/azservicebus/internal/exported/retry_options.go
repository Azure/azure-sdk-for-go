// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import "time"

// NOTE: this is exposed via type-aliasing in azservicebus/client.go

// RetryOptions represent the options for retries.
type RetryOptions struct {
	// MaxRetries specifies the maximum number of attempts a failed operation will be retried
	// before producing an error.
	// The default value is three.  A value less than zero means one try and no retries.
	MaxRetries int32

	// RetryDelay specifies the initial amount of delay to use before retrying an operation.
	// The delay increases exponentially with each retry up to the maximum specified by MaxRetryDelay.
	// The default value is four seconds.  A value less than zero means no delay between retries.
	RetryDelay time.Duration

	// MaxRetryDelay specifies the maximum delay allowed before retrying an operation.
	// Typically the value is greater than or equal to the value specified in RetryDelay.
	// The default Value is 120 seconds.  A value less than zero means there is no cap.
	MaxRetryDelay time.Duration

	// TryTimeout specifies the maximum duration of a single attempt, not the whole
	// operation. When set to a positive value, each retry attempt is bounded by a
	// fresh TryTimeout, recomputed per attempt (so a slow attempt does not shrink
	// later attempts' budgets) but still capped by the caller's remaining deadline,
	// so a single slow attempt cannot consume the entire operation. An attempt that
	// exceeds TryTimeout is treated as a retryable failure (subject to MaxRetries);
	// the overall bound remains the caller's context.
	//
	// The default value is zero, which means no per-attempt timeout: an attempt is
	// bounded only by the caller's context, preserving the historical behavior. Set
	// a positive value (for example 60 seconds, matching the .NET SDK's default per
	// attempt) to opt in to per-attempt bounding. A negative value is also treated
	// as no per-attempt timeout.
	TryTimeout time.Duration
}

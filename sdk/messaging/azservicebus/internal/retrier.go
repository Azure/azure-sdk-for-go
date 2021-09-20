// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"time"

	"github.com/jpillora/backoff"
)

// TODO: we should discuss what a common policy would be.
var DefaultRetryPolicy Retrier = &BackoffRetrier{
	Backoff: backoff.Backoff{
		Factor: 1,
		Min:    time.Second * 5,
	},
	MaxRetries: 5,
}

// A retrier that allows you to do a basic for loop and get backoff
// and retry limits. See `Try` for more details on how to use it.
type Retrier interface {
	// Copies the retrier. Retriers are stateful and must be copied
	// before starting a set of retries.
	Copy() Retrier

	// Returns true if the retries were exhausted.
	Exhausted() bool

	// Try marks an attempt to call (first call to Try() does not sleep).
	// Will return false if the `ctx` is cancelled or if we exhaust our retries.
	//
	//    rp := RetryPolicy{Backoff:defaultBackoffPolicy, MaxRetries:5}
	//
	//    for rp.Try(ctx) {
	//       <your code>
	//    }
	//
	//    if rp.Cancelled() || rp.Exhausted() {
	//       // no more retries needed
	//    }
	//
	Try(ctx context.Context) bool
}

// Encapsulates a backoff policy, which allows you to configure the amount of
// time in between retries as well as the maximum retries allowed (via MaxRetries)
// NOTE: this should be copied by the caller as it is stateful.
type BackoffRetrier struct {
	Backoff    backoff.Backoff
	MaxRetries int

	tries int
}

// Copies the backoff retrier since it's stateful.
func (rp *BackoffRetrier) Copy() Retrier {
	copy := *rp
	return &copy
}

// Exhausted is true if all the retries have been used.
func (rp *BackoffRetrier) Exhausted() bool {
	return rp.tries > rp.MaxRetries
}

// Try marks an attempt to call (first call to Try() does not sleep).
// Will return false if the `ctx` is cancelled or if we exhaust our retries.
//
//    rp := RetryPolicy{Backoff:defaultBackoffPolicy, MaxRetries:5}
//
//    for rp.Try(ctx) {
//       <your code>
//    }
//
//    if rp.Cancelled() || rp.Exhausted() {
//       // no more retries needed
//    }
//
func (rp *BackoffRetrier) Try(ctx context.Context) bool {
	defer func() { rp.tries++ }()

	if rp.tries == 0 {
		// first 'try' is always free
		return true
	}

	if rp.Exhausted() {
		return false
	}

	select {
	case <-time.After(rp.Backoff.Duration()):
		return true
	case <-ctx.Done():
		return false
	}
}

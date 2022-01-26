// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package utils

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// EventRetry is the name for retry events
const EventRetry = "azsb.Retry"

type RetryFnArgs struct {
	I int32
	// LastErr is the returned error from the previous loop.
	// If you have potentially expensive
	LastErr error

	resetAttempts bool
}

// ResetAttempts causes the current retry attempt number to be reset
// in time for the next recovery (should we fail).
// NOTE: Use of this should be pretty rare, it's really only needed when you have
// a situation like Receiver.ReceiveMessages() that can recovery but intentionally
// does not return.
func (rf *RetryFnArgs) ResetAttempts() {
	rf.resetAttempts = true
}

// Retry runs a standard retry loop. It executes your passed in fn as the body of the loop.
// It returns if it exceeds the number of configured retry options or if 'isFatal' returns true.
func Retry(ctx context.Context, name string, fn func(ctx context.Context, args *RetryFnArgs) error, isFatalFn func(err error) bool, o RetryOptions) error {
	if isFatalFn == nil {
		panic("isFatalFn is nil, errors would panic")
	}

	var ro RetryOptions = o
	setDefaults(&ro)

	var err error

	for i := int32(0); i <= ro.MaxRetries; i++ {
		if i > 0 {
			sleep := calcDelay(ro, i)
			log.Writef(EventRetry, "(%s) Attempt %d sleeping for %s", name, i, sleep)
			time.Sleep(sleep)
		}

		args := RetryFnArgs{
			I:       i,
			LastErr: err,
		}
		err = fn(ctx, &args)

		if args.resetAttempts {
			log.Writef(EventRetry, "(%s) Resetting attempts", name)
			i = int32(0)
		}

		if err != nil {
			if isFatalFn(err) {
				log.Writef(EventRetry, "(%s) Attempt %d returned non-retryable error: %s", name, i, err.Error())
				return err
			} else {
				log.Writef(EventRetry, "(%s) Attempt %d returned retryable error: %s", name, i, err.Error())
			}

			continue
		}

		return nil
	}

	return err
}

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
}

func setDefaults(o *RetryOptions) {
	if o.MaxRetries == 0 {
		o.MaxRetries = 3
	} else if o.MaxRetries < 0 {
		o.MaxRetries = 0
	}
	if o.MaxRetryDelay == 0 {
		o.MaxRetryDelay = 120 * time.Second
	} else if o.MaxRetryDelay < 0 {
		// not really an unlimited cap, but sufficiently large enough to be considered as such
		o.MaxRetryDelay = math.MaxInt64
	}
	if o.RetryDelay == 0 {
		o.RetryDelay = 4 * time.Second
	} else if o.RetryDelay < 0 {
		o.RetryDelay = 0
	}
}

// (adapted from from azcore/policy_retry)
func calcDelay(o RetryOptions, try int32) time.Duration {
	if try == 0 {
		return 0
	}

	pow := func(number int64, exponent int32) int64 { // pow is nested helper function
		var result int64 = 1
		for n := int32(0); n < exponent; n++ {
			result *= number
		}
		return result
	}

	delay := time.Duration(pow(2, try)-1) * o.RetryDelay

	// Introduce some jitter:  [0.0, 1.0) / 2 = [0.0, 0.5) + 0.8 = [0.8, 1.3)
	delay = time.Duration(delay.Seconds() * (rand.Float64()/2 + 0.8) * float64(time.Second)) // NOTE: We want math/rand; not crypto/rand
	if delay > o.MaxRetryDelay {
		delay = o.MaxRetryDelay
	}
	return delay
}

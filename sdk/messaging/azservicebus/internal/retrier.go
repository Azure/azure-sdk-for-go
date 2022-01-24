// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/jpillora/backoff"
)

// A retrier that allows you to do a basic for loop and get backoff
// and retry limits. See `Try` for more details on how to use it.
type Retrier interface {
	// Copies the retrier. Retriers are stateful and must be copied
	// before starting a set of retries.
	Copy() Retrier

	// Exhausted is true if the retries were exhausted.
	Exhausted() bool

	// CurrentTry is the current try (0 for the first run before retries)
	CurrentTry() int

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
type backoffRetrier struct {
	backoff    backoff.Backoff
	MaxRetries int

	tries int
}

// BackoffRetrierParams are parameters for NewBackoffRetrier.
type BackoffRetrierParams struct {
	// MaxRetries is the maximum number of tries (after the first attempt)
	// that are allowed.
	MaxRetries int
	// Factor is the multiplying factor for each increment step
	Factor float64
	// Jitter eases contention by randomizing backoff steps
	Jitter bool
	// Min and Max are the minimum and maximum values of the counter
	Min, Max time.Duration
}

// NewBackoffRetrier creates a retrier that allows for configurable
// min/max times, jitter and maximum retries.
func NewBackoffRetrier(params BackoffRetrierParams) Retrier {
	return &backoffRetrier{
		backoff: backoff.Backoff{
			Factor: params.Factor,
			Jitter: params.Jitter,
			Min:    params.Min,
			Max:    params.Max,
		},
		MaxRetries: params.MaxRetries,
	}
}

// Copies the backoff retrier since it's stateful.
func (rp *backoffRetrier) Copy() Retrier {
	copy := *rp
	return &copy
}

// Exhausted is true if all the retries have been used.
func (rp *backoffRetrier) Exhausted() bool {
	return rp.tries > rp.MaxRetries
}

// CurrentTry is the current try number (0 for the first run before retries)
func (rp *backoffRetrier) CurrentTry() int {
	return rp.tries
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
func (rp *backoffRetrier) Try(ctx context.Context) bool {
	defer func() { rp.tries++ }()

	select {
	case <-ctx.Done():
		return false
	default:
	}

	if rp.tries == 0 {
		// first 'try' is always free
		return true
	}

	if rp.Exhausted() {
		return false
	}

	select {
	case <-time.After(rp.backoff.Duration()):
		return true
	case <-ctx.Done():
		return false
	}
}

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
// 'isFatal' can be nil, and defaults to just checking that ServiceBusError(err).recoveryKind != recoveryKindNonRetriable.
// It returns if it exceeds the number of configured retry options or if 'isFatal' returns true.
func Retry(ctx context.Context, name string, fn func(ctx context.Context, args *RetryFnArgs) error, isFatalFn func(err error) bool, o RetryOptions) error {
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
			if isFatalFn != nil {
				if isFatalFn(err) {
					log.Writef(EventRetry, "(%s) Attempt %d returned non-retryable error: %s", name, i, err.Error())
					return err
				} else {
					log.Writef(EventRetry, "(%s) Attempt %d returned retryable error: %s", name, i, err.Error())
				}
			} else {
				recoveryKind := ToSBE(ctx, err).RecoveryKind
				if recoveryKind == RecoveryKindFatal {
					log.Writef(EventRetry, "(%s) Attempt %d returned non-retryable error: %s", name, i, err.Error())
					return err
				} else {
					log.Writef(EventRetry, "(%s) Attempt %d returned retryable error with recovery kind %s: %s", name, i, recoveryKind, err.Error())
				}
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

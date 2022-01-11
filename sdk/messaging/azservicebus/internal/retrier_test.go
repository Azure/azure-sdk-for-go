// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestRetrier(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		retrier := NewBackoffRetrier(struct {
			MaxRetries int
			Factor     float64
			Jitter     bool
			Min        time.Duration
			Max        time.Duration
		}{
			MaxRetries: 5,
			Factor:     0,
		})

		require := require.New(t)

		// first iteration is always free (ie, that's not the
		// retry part)
		require.True(retrier.Try(context.Background()))

		// now we're doing retries
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))
		require.True(retrier.Try(context.Background()))

		// and it's the 6th retry that fails since we've exhausted
		// the retries we're allotted.
		require.False(retrier.Try(context.Background()))
		require.True(retrier.Exhausted())
	})

	t.Run("Cancellation", func(t *testing.T) {
		retrier := NewBackoffRetrier(struct {
			MaxRetries int
			Factor     float64
			Jitter     bool
			Min        time.Duration
			Max        time.Duration
		}{
			MaxRetries: 5,
			Factor:     0,
		})

		// first iteration is always free (ie, that's not the
		// retry part)
		cancelledContext, cancel := context.WithCancel(context.Background())
		cancel()
		require.False(t, retrier.Try(cancelledContext))
	})
}

var fastRetryOptions = RetryOptions{
	// note: omitting MaxRetries just to give a sanity check that
	// we do setDefaults() before we run.
	RetryDelay:    time.Millisecond,
	MaxRetryDelay: time.Millisecond,
}

func TestRetryBasic(t *testing.T) {
	called := 0

	err := Retry(context.Background(), "retrytest", func(ctx context.Context, args *RetryFnArgs) error {
		require.NotNil(t, args)
		require.NotNil(t, ctx)

		called++

		return &amqp.DetachError{}
	}, nil, fastRetryOptions)

	var de *amqp.DetachError
	require.ErrorAs(t, err, &de)
	require.EqualValues(t, 4, called)
}

func TestRetryWithFatalError(t *testing.T) {
	called := 0

	err := Retry(context.Background(), "retrytest", func(ctx context.Context, args *RetryFnArgs) error {
		require.NotNil(t, args)
		require.NotNil(t, ctx)

		called++

		return &amqp.Error{
			// this is just a basic non-recoverable situation - typically happens if the
			// lock period expires.
			Condition: amqp.ErrorCondition("com.microsoft:message-lock-lost"),
		}
	}, nil, fastRetryOptions)

	// fatal error so we only called the function once
	require.EqualValues(t, 1, called)

	var testErr *amqp.Error

	require.ErrorAs(t, err, &testErr)
	require.EqualValues(t, "com.microsoft:message-lock-lost", testErr.Condition)
}

func TestRetryCustomIsFatal(t *testing.T) {
	called := 0
	var totallyHarmlessErrorAsFatal = errors.New("I'm supposed to be harmless but the custom error handler is going to make me fatal")
	var isFatalErr error

	err := Retry(context.Background(), "retrytest", func(ctx context.Context, args *RetryFnArgs) error {
		require.NotNil(t, args)
		require.NotNil(t, ctx)

		called++

		return totallyHarmlessErrorAsFatal
	}, func(err error) bool {
		require.Nil(t, isFatalErr, "should only get called once")
		isFatalErr = err
		return true
	}, fastRetryOptions)

	// fatal error so we only called the function once
	require.EqualValues(t, 1, called)

	require.ErrorIs(t, err, totallyHarmlessErrorAsFatal)
	require.ErrorIs(t, isFatalErr, totallyHarmlessErrorAsFatal)
}

func TestRetryDefaults(t *testing.T) {
	ro := RetryOptions{}
	setDefaults(&ro)

	require.EqualValues(t, 3, ro.MaxRetries)
	require.EqualValues(t, 4*time.Second, ro.RetryDelay)
	require.EqualValues(t, 2*time.Minute, ro.MaxRetryDelay)

	// this is an interesting default. Anything < 0 basically
	// causes the max delay to be "infinite"
	ro.MaxRetryDelay = -1
	// whereas this just normalizes to '0'
	ro.RetryDelay = -1
	ro.MaxRetries = -1
	setDefaults(&ro)
	require.EqualValues(t, time.Duration(math.MaxInt64), ro.MaxRetryDelay)
	require.EqualValues(t, 0, ro.MaxRetries)
	require.EqualValues(t, time.Duration(0), ro.RetryDelay)
}

func TestCalcDelay(t *testing.T) {
	// calcDelay introduces some jitter, automatically.
	ro := RetryOptions{}
	setDefaults(&ro)
	d := calcDelay(ro, 0)
	require.EqualValues(t, 0, d)

	// by default the first calc is 2^attempt
	d = calcDelay(ro, 1)
	require.LessOrEqual(t, d, 6*time.Second)
	require.GreaterOrEqual(t, d, time.Second)
}

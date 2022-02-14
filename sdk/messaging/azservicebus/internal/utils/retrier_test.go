// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetrier(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		ctx := context.Background()

		called := 0

		err := Retry(ctx, "Retrier", func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return nil
		}, func(err error) bool {
			panic("won't get called")
		}, RetryOptions{})

		require.Nil(t, err)
		require.EqualValues(t, 1, called)
	})

	t.Run("FailsThenSucceeds", func(t *testing.T) {
		ctx := context.Background()

		called := 0
		isFatalCalled := 0

		isFatalFn := func(err error) bool {
			require.NotNil(t, err)
			// we'll just keep saying the errors aren't fatal.
			isFatalCalled++
			return false
		}

		err := Retry(ctx, "FailsThenSucceeds", func(ctx context.Context, args *RetryFnArgs) error {
			called++

			if args.I == 3 {
				// we're on the last iteration, succeed
				return nil
			}

			return fmt.Errorf("Error, iteration %d", args.I)
		}, isFatalFn, fastRetryOptions)

		require.EqualValues(t, 4, called)
		require.EqualValues(t, 3, isFatalCalled)

		// if an attempt succeeds then there's no error (despite previous failed tries)
		require.NoError(t, err)
	})

	t.Run("FatalFailure", func(t *testing.T) {
		ctx := context.Background()
		called := 0

		isFatalFn := func(err error) bool {
			require.EqualValues(t, "isFatalFn says this is a fatal error", err.Error())
			return true
		}

		err := Retry(ctx, "FatalFailure", func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return errors.New("isFatalFn says this is a fatal error")
		}, isFatalFn, RetryOptions{})

		require.EqualValues(t, "isFatalFn says this is a fatal error", err.Error())
		require.EqualValues(t, 1, called)
	})

	t.Run("Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		isFatalFn := func(err error) bool {
			return errors.Is(err, context.Canceled)
		}

		// it's up to
		err := Retry(ctx, "Cancellation", func(ctx context.Context, args *RetryFnArgs) error {
			// NOTE: it's up to the underlying function to handle cancellation. `Retry` doesn't
			// do anything but propagate it.
			select {
			case <-ctx.Done():
			default:
				require.Fail(t, "Context should have been cancelled")
			}

			return context.Canceled
		}, isFatalFn, RetryOptions{})

		require.ErrorIs(t, context.Canceled, err)
	})

	t.Run("ResetAttempts", func(t *testing.T) {
		isFatalFn := func(err error) bool {
			return errors.Is(err, context.Canceled)
		}

		customRetryOptions := fastRetryOptions
		customRetryOptions.MaxRetries = 1

		var actualAttempts []int32

		err := Retry(context.Background(), "ResetAttempts", func(ctx context.Context, args *RetryFnArgs) error {
			actualAttempts = append(actualAttempts, args.I)

			if len(actualAttempts) == 3 {
				args.ResetAttempts()
			}

			return errors.New("whatever")
		}, isFatalFn, RetryOptions{
			MaxRetries:    2,
			RetryDelay:    time.Millisecond,
			MaxRetryDelay: time.Millisecond,
		})

		expectedAttempts := []int32{
			0, 1, 2, // we resetted attempts here.
			1, 2, // and we start at the first retry attempt again.
		}

		require.EqualValues(t, "whatever", err.Error())
		require.EqualValues(t, expectedAttempts, actualAttempts)
	})

	t.Run("DisableRetries", func(t *testing.T) {
		isFatalFn := func(err error) bool {
			return errors.Is(err, context.Canceled)
		}

		customRetryOptions := fastRetryOptions
		customRetryOptions.MaxRetries = -1

		called := 0

		err := Retry(context.Background(), "ResetAttempts", func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return errors.New("whatever")
		}, isFatalFn, customRetryOptions)

		require.EqualValues(t, 1, called)
		require.EqualValues(t, "whatever", err.Error())
	})
}

func Test_calcDelay(t *testing.T) {
	t.Run("can't exceed max retry delay", func(t *testing.T) {
		duration := calcDelay(RetryOptions{
			RetryDelay:    time.Hour,
			MaxRetryDelay: time.Minute,
		}, 1)

		require.EqualValues(t, time.Minute, duration)
	})

	t.Run("increases with jitter", func(t *testing.T) {
		duration := calcDelay(RetryOptions{
			RetryDelay:    time.Minute,
			MaxRetryDelay: time.Hour,
		}, 1)

		require.GreaterOrEqual(t, duration, time.Duration((2-1)*time.Minute.Seconds()*0.8*float64(time.Second)))
		require.LessOrEqual(t, duration, time.Duration((2-1)*time.Minute.Seconds()*1.3*float64(time.Second)))

		duration = calcDelay(RetryOptions{
			RetryDelay:    time.Minute,
			MaxRetryDelay: time.Hour,
		}, 2)

		require.GreaterOrEqual(t, duration, time.Duration((2*2-1)*time.Minute.Seconds()*0.8*float64(time.Second)))
		require.LessOrEqual(t, duration, time.Duration((2*2-1)*time.Minute.Seconds()*1.3*float64(time.Second)))

		duration = calcDelay(RetryOptions{
			RetryDelay:    time.Minute,
			MaxRetryDelay: time.Hour,
		}, 3)

		require.GreaterOrEqual(t, duration, time.Duration((2*2*2-1)*time.Minute.Seconds()*0.8*float64(time.Second)))
		require.LessOrEqual(t, duration, time.Duration((2*2*2-1)*time.Minute.Seconds()*1.3*float64(time.Second)))
	})
}

var fastRetryOptions = RetryOptions{
	// note: omitting MaxRetries just to give a sanity check that
	// we do setDefaults() before we run.
	RetryDelay:    time.Millisecond,
	MaxRetryDelay: time.Millisecond,
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

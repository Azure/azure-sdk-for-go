// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

const testLogEvent = azlog.Event("testLogEvent")

func TestRetrier(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		ctx := context.Background()

		called := 0

		err := Retry(ctx, testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return nil
		}, func(err error) bool {
			panic("won't get called")
		}, exported.RetryOptions{}, nil)

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

		err := Retry(ctx, testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
			called++

			if args.I == 3 {
				// we're on the last iteration, succeed
				return nil
			}

			return fmt.Errorf("Error, iteration %d", args.I)
		}, isFatalFn, fastRetryOptions, nil)

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

		err := Retry(ctx, testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return errors.New("isFatalFn says this is a fatal error")
		}, isFatalFn, exported.RetryOptions{}, nil)

		require.EqualValues(t, "isFatalFn says this is a fatal error", err.Error())
		require.EqualValues(t, 1, called)
	})

	t.Run("ResetAttempts", func(t *testing.T) {
		isFatalFn := func(err error) bool {
			return errors.Is(err, context.Canceled)
		}

		var actualAttempts []int32

		maxRetries := int32(2)

		err := Retry(context.Background(), testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
			actualAttempts = append(actualAttempts, args.I)

			if len(actualAttempts) == int(maxRetries+1) {
				args.ResetAttempts()
			}

			return errors.New("whatever")
		}, isFatalFn, exported.RetryOptions{
			MaxRetries:    maxRetries,
			RetryDelay:    time.Millisecond,
			MaxRetryDelay: time.Millisecond,
		}, nil)

		expectedAttempts := []int32{
			0, 1, 2, // we resetted attempts here.
			0, 1, 2, // and we start at the first retry attempt again.
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

		err := Retry(context.Background(), testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return errors.New("whatever")
		}, isFatalFn, customRetryOptions, nil)

		require.EqualValues(t, 1, called)
		require.EqualValues(t, "whatever", err.Error())
	})
}

func TestCancellationCancelsSleep(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	isFatalFn := func(err error) bool {
		return errors.Is(err, context.Canceled)
	}

	called := 0

	err := Retry(ctx, testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
		called++
		return errors.New("try again")
	}, isFatalFn, exported.RetryOptions{
		RetryDelay: time.Hour,
	}, nil)

	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, called, 1)
}

func TestCancellationFromUserFunc(t *testing.T) {
	alreadyCancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	canceledfromFunc := errors.New("the user func got the cancellation signal")

	isFatalFn := func(err error) bool {
		return errors.Is(err, canceledfromFunc)
	}

	called := 0

	err := Retry(alreadyCancelledCtx, testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
		called++

		select {
		case <-ctx.Done():
			return canceledfromFunc
		default:
			panic("Context should have been cancelled")
		}
	}, isFatalFn, exported.RetryOptions{}, nil)

	require.Error(t, err)
	require.ErrorIs(t, err, canceledfromFunc)
}

func TestCancellationTimeoutsArentPropagatedToUser(t *testing.T) {
	isFatalFn := func(err error) bool {
		// we want to exhaust all retries and run through the "sleep between retries" logic.
		return false
	}

	tryAgainErr := errors.New("try again")
	called := 0

	err := Retry(context.Background(), testLogEvent, "notused", func(ctx context.Context, args *RetryFnArgs) error {
		called++
		require.NoError(t, ctx.Err(), "our sleep/timeout doesn't show up for users")
		return tryAgainErr
	}, isFatalFn, exported.RetryOptions{
		RetryDelay: time.Millisecond,
	}, nil)

	require.Error(t, err)
	require.ErrorIs(t, err, tryAgainErr, "error should be propagated from user callback")
	require.Equal(t, called, 1+3, "all attempts exhausted since we never returned a fatal error")
}

func Test_calcDelay(t *testing.T) {
	t.Run("can't exceed max retry delay", func(t *testing.T) {
		duration := calcDelay(exported.RetryOptions{
			RetryDelay:    time.Hour,
			MaxRetryDelay: time.Minute,
		}, 1)

		require.EqualValues(t, time.Minute, duration)
	})

	t.Run("increases with jitter", func(t *testing.T) {
		duration := calcDelay(exported.RetryOptions{
			RetryDelay:    time.Minute,
			MaxRetryDelay: time.Hour,
		}, 1)

		require.GreaterOrEqual(t, duration, time.Duration((2-1)*time.Minute.Seconds()*0.8*float64(time.Second)))
		require.LessOrEqual(t, duration, time.Duration((2-1)*time.Minute.Seconds()*1.3*float64(time.Second)))

		duration = calcDelay(exported.RetryOptions{
			RetryDelay:    time.Minute,
			MaxRetryDelay: time.Hour,
		}, 2)

		require.GreaterOrEqual(t, duration, time.Duration((2*2-1)*time.Minute.Seconds()*0.8*float64(time.Second)))
		require.LessOrEqual(t, duration, time.Duration((2*2-1)*time.Minute.Seconds()*1.3*float64(time.Second)))

		duration = calcDelay(exported.RetryOptions{
			RetryDelay:    time.Minute,
			MaxRetryDelay: time.Hour,
		}, 3)

		require.GreaterOrEqual(t, duration, time.Duration((2*2*2-1)*time.Minute.Seconds()*0.8*float64(time.Second)))
		require.LessOrEqual(t, duration, time.Duration((2*2*2-1)*time.Minute.Seconds()*1.3*float64(time.Second)))
	})
}

var fastRetryOptions = exported.RetryOptions{
	// note: omitting MaxRetries just to give a sanity check that
	// we do setDefaults() before we run.
	RetryDelay:    time.Millisecond,
	MaxRetryDelay: time.Millisecond,
}

func TestRetryDefaults(t *testing.T) {
	ro := exported.RetryOptions{}
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
	ro := exported.RetryOptions{}
	setDefaults(&ro)
	d := calcDelay(ro, 0)
	require.EqualValues(t, 0, d)

	// by default the first calc is 2^attempt
	d = calcDelay(ro, 1)
	require.LessOrEqual(t, d, 6*time.Second)
	require.GreaterOrEqual(t, d, time.Second)
}

func TestRetryLogging(t *testing.T) {
	t.Run("normal error", func(t *testing.T) {
		logsFn := test.CaptureLogsForTest(false)

		err := Retry(context.Background(), testLogEvent, "(my_operation)", func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc", "Attempt %d, within test func, returning error hello", args.I)
			return errors.New("hello")
		}, func(err error) bool {
			return false
		}, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, nil)
		require.EqualError(t, err, "hello")

		require.Equal(t, []string{
			"[TestFunc] Attempt 0, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 0 returned retryable error: hello",

			"[testLogEvent] (my_operation) Retry attempt 1 sleeping for <time elided>",
			"[TestFunc] Attempt 1, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 1 returned retryable error: hello",

			"[testLogEvent] (my_operation) Retry attempt 2 sleeping for <time elided>",
			"[TestFunc] Attempt 2, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 2 returned retryable error: hello",

			"[testLogEvent] (my_operation) Retry attempt 3 sleeping for <time elided>",
			"[TestFunc] Attempt 3, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 3 returned retryable error: hello",
		}, normalizeRetryLogLines(logsFn()))
	})

	t.Run("normal error2", func(t *testing.T) {
		test.EnableStdoutLogging(t)

		err := Retry(context.Background(), testLogEvent, "(my_operation)", func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc", "Attempt %d, within test func, returning error hello", args.I)
			return errors.New("hello")
		}, func(err error) bool {
			return false
		}, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, nil)
		require.EqualError(t, err, "hello")
	})

	t.Run("cancellation error", func(t *testing.T) {
		logsFn := test.CaptureLogsForTest(false)

		err := Retry(context.Background(), testLogEvent, "(test_operation)", func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc",
				"Attempt %d, within test func", args.I)
			return context.Canceled
		}, func(err error) bool {
			return errors.Is(err, context.Canceled)
		}, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, nil)
		require.ErrorIs(t, err, context.Canceled)

		require.Equal(t, []string{
			"[TestFunc] Attempt 0, within test func",
			"[testLogEvent] (test_operation) Retry attempt 0 was cancelled, stopping: context canceled",
		}, normalizeRetryLogLines(logsFn()))
	})

	t.Run("custom fatal error", func(t *testing.T) {
		logsFn := test.CaptureLogsForTest(false)

		err := Retry(context.Background(), testLogEvent, "(test_operation)", func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc",
				"Attempt %d, within test func", args.I)
			return errors.New("custom fatal error")
		}, func(err error) bool {
			return true
		}, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, nil)
		require.EqualError(t, err, "custom fatal error")

		require.Equal(t, []string{
			"[TestFunc] Attempt 0, within test func",
			"[testLogEvent] (test_operation) Retry attempt 0 returned non-retryable error: custom fatal error",
		}, normalizeRetryLogLines(logsFn()))
	})

	t.Run("with reset attempts", func(t *testing.T) {
		logsFn := test.CaptureLogsForTest(false)
		reset := false

		err := Retry(context.Background(), testLogEvent, "(test_operation)", func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc", "Attempt %d, within test func", args.I)

			if !reset {
				azlog.Writef("TestFunc", "Attempt %d, resetting", args.I)
				args.ResetAttempts()
				reset = true
				return &amqp.LinkError{}
			}

			if reset {
				azlog.Writef("TestFunc", "Attempt %d, return nil", args.I)
				return nil
			}

			return errors.New("custom fatal error")
		}, func(err error) bool {
			var de amqp.LinkError
			return errors.Is(err, &de)
		}, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, nil)
		require.Nil(t, err)

		require.Equal(t, []string{
			"[TestFunc] Attempt 0, within test func",
			"[TestFunc] Attempt 0, resetting",
			"[testLogEvent] (test_operation) Resetting retry attempts",
			"[testLogEvent] (test_operation) Retry attempt -1 returned retryable error: amqp: link closed",
			"[TestFunc] Attempt 0, within test func",
			"[TestFunc] Attempt 0, return nil",
		}, normalizeRetryLogLines(logsFn()))
	})
}

// retryRE is used to replace the 'retry time' with a consistent string to make
// unit tests against logging simpler
// A typical string: "[azsb.Retry] (retry) Attempt 1 sleeping for 1.10233ms"
var retryRE = regexp.MustCompile(`[\d.]+(µs|ms|ns)`)

func normalizeRetryLogLines(msgs []string) []string {
	var newLogs []string

	for i := 0; i < len(msgs); i++ {
		newLogs = append(newLogs, retryRE.ReplaceAllString(msgs[i], "<time elided>"))
	}

	return newLogs
}

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
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

const testLogEvent = azlog.Event("testLogEvent")

func TestRetrier(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		ctx := context.Background()

		called := 0

		err := Retry(ctx, testLogEvent, func() string { return "notused" }, exported.RetryOptions{}, func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return nil
		}, func(err error) bool {
			panic("won't get called")
		})

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

		err := Retry(ctx, testLogEvent, func() string { return "notused" }, fastRetryOptions, func(ctx context.Context, args *RetryFnArgs) error {
			called++

			if args.I == 3 {
				// we're on the last iteration, succeed
				return nil
			}

			return fmt.Errorf("Error, iteration %d", args.I)
		}, isFatalFn)

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

		err := Retry(ctx, testLogEvent, func() string { return "notused" }, exported.RetryOptions{}, func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return errors.New("isFatalFn says this is a fatal error")
		}, isFatalFn)

		require.EqualValues(t, "isFatalFn says this is a fatal error", err.Error())
		require.EqualValues(t, 1, called)
	})

	t.Run("ResetAttempts", func(t *testing.T) {
		isFatalFn := func(err error) bool {
			return errors.Is(err, context.Canceled)
		}

		var actualAttempts []int32

		maxRetries := int32(2)

		err := Retry(context.Background(), testLogEvent, func() string { return "notused" }, exported.RetryOptions{
			MaxRetries:    maxRetries,
			RetryDelay:    time.Millisecond,
			MaxRetryDelay: time.Millisecond,
		}, func(ctx context.Context, args *RetryFnArgs) error {
			actualAttempts = append(actualAttempts, args.I)

			if len(actualAttempts) == int(maxRetries+1) {
				args.ResetAttempts()
			}

			return errors.New("whatever")
		}, isFatalFn)

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

		err := Retry(context.Background(), testLogEvent, func() string { return "notused" }, customRetryOptions, func(ctx context.Context, args *RetryFnArgs) error {
			called++
			return errors.New("whatever")
		}, isFatalFn)

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

	err := Retry(ctx, testLogEvent, func() string { return "notused" }, exported.RetryOptions{
		RetryDelay: time.Hour,
	}, func(ctx context.Context, args *RetryFnArgs) error {
		called++
		return errors.New("try again")
	}, isFatalFn)

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

	err := Retry(alreadyCancelledCtx, testLogEvent, func() string { return "notused" }, exported.RetryOptions{}, func(ctx context.Context, args *RetryFnArgs) error {
		called++

		select {
		case <-ctx.Done():
			return canceledfromFunc
		default:
			panic("Context should have been cancelled")
		}
	}, isFatalFn)

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

	err := Retry(context.Background(), testLogEvent, func() string { return "notused" }, exported.RetryOptions{
		RetryDelay: time.Millisecond,
	}, func(ctx context.Context, args *RetryFnArgs) error {
		called++
		require.NoError(t, ctx.Err(), "our sleep/timeout doesn't show up for users")
		return tryAgainErr
	}, isFatalFn)

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

func TestCalcDelay2(t *testing.T) {
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
	var logs []string

	azlog.SetListener(func(e azlog.Event, s string) {
		logs = append(logs, fmt.Sprintf("[%-10s] %s", e, s))
	})

	defer azlog.SetListener(nil)

	t.Run("normal error", func(t *testing.T) {
		logs = nil

		err := Retry(context.Background(), testLogEvent, func() string { return "my_operation" }, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc", "Attempt %d, within test func, returning error hello", args.I)
			return errors.New("hello")
		}, func(err error) bool {
			return false
		})
		require.EqualError(t, err, "hello")

		require.Equal(t, []string{
			"[TestFunc  ] Attempt 0, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 0 returned retryable error: hello",

			"[testLogEvent] (my_operation) Retry attempt 1 sleeping for <time elided>",
			"[TestFunc  ] Attempt 1, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 1 returned retryable error: hello",

			"[testLogEvent] (my_operation) Retry attempt 2 sleeping for <time elided>",
			"[TestFunc  ] Attempt 2, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 2 returned retryable error: hello",

			"[testLogEvent] (my_operation) Retry attempt 3 sleeping for <time elided>",
			"[TestFunc  ] Attempt 3, within test func, returning error hello",
			"[testLogEvent] (my_operation) Retry attempt 3 returned retryable error: hello",
		}, normalizeRetryLogLines(logs))
	})

	t.Run("cancellation error", func(t *testing.T) {
		logs = nil

		err := Retry(context.Background(), testLogEvent, func() string { return "test_operation" }, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc",
				"Attempt %d, within test func", args.I)
			return context.Canceled
		}, func(err error) bool {
			return errors.Is(err, context.Canceled)
		})
		require.ErrorIs(t, err, context.Canceled)

		require.Equal(t, []string{
			"[TestFunc  ] Attempt 0, within test func",
			"[testLogEvent] (test_operation) Retry attempt 0 was cancelled, stopping: context canceled",
		}, normalizeRetryLogLines(logs))
	})

	t.Run("custom fatal error", func(t *testing.T) {
		logs = nil

		err := Retry(context.Background(), testLogEvent, func() string { return "test_operation" }, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, func(ctx context.Context, args *RetryFnArgs) error {
			azlog.Writef("TestFunc",
				"Attempt %d, within test func", args.I)
			return errors.New("custom fatal error")
		}, func(err error) bool {
			return true
		})
		require.EqualError(t, err, "custom fatal error")

		require.Equal(t, []string{
			"[TestFunc  ] Attempt 0, within test func",
			"[testLogEvent] (test_operation) Retry attempt 0 returned non-retryable error: custom fatal error",
		}, normalizeRetryLogLines(logs))
	})

	t.Run("with reset attempts", func(t *testing.T) {
		logs = nil

		reset := false

		err := Retry(context.Background(), testLogEvent, func() string { return "test_operation" }, exported.RetryOptions{
			RetryDelay: time.Microsecond,
		}, func(ctx context.Context, args *RetryFnArgs) error {
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
			var de *amqp.LinkError
			return !errors.As(err, &de)
		})
		require.Nil(t, err)

		require.Equal(t, []string{
			"[TestFunc  ] Attempt 0, within test func",
			"[TestFunc  ] Attempt 0, resetting",
			"[testLogEvent] (test_operation) Resetting retry attempts",
			"[testLogEvent] (test_operation) Retry attempt -1 returned retryable error: amqp: link closed",
			"[TestFunc  ] Attempt 0, within test func",
			"[TestFunc  ] Attempt 0, return nil",
		}, normalizeRetryLogLines(logs))
	})
}

func BenchmarkCalcDelay_defaultSettings(b *testing.B) {
	retryOptions := exported.RetryOptions{}
	setDefaults(&retryOptions)

	for i := 0; i < b.N; i++ {
		calcDelay(retryOptions, 32)
	}
}

func BenchmarkCalcDelay_overflow(b *testing.B) {
	retryOptions := exported.RetryOptions{
		RetryDelay:    1,
		MaxRetryDelay: math.MaxInt64,
	}
	setDefaults(&retryOptions)

	for i := 0; i < b.N; i++ {
		calcDelay(retryOptions, 100)
	}
}

func TestCalcDelay(t *testing.T) {
	requireWithinJitter := func(t testing.TB, expected, actual time.Duration) {
		lower, upper := float64(expected)*0.8, float64(expected)*1.3
		require.Truef(
			t, float64(actual) >= lower && float64(actual) <= upper,
			"%.2f not within jitter of %.2f", actual.Seconds(), expected.Seconds(),
		)
	}

	t.Run("basic cases", func(t *testing.T) {
		retryOptions := exported.RetryOptions{
			RetryDelay:    1 * time.Second,
			MaxRetryDelay: 30 * time.Second,
		}
		setDefaults(&retryOptions)

		for i := int32(1); i <= 5; i++ {
			delay := float64(calcDelay(retryOptions, i))
			expected := float64((1<<i - 1) * int64(retryOptions.RetryDelay))
			requireWithinJitter(
				t, time.Duration(expected), time.Duration(delay),
			)
		}
		for i := int32(6); i < 100; i++ {
			require.Equal(
				t,
				calcDelay(retryOptions, i),
				retryOptions.MaxRetryDelay,
			)
		}
	})

	t.Run("overflow", func(t *testing.T) {
		retryOptions := exported.RetryOptions{
			RetryDelay:    1,
			MaxRetryDelay: math.MaxInt64,
		}
		setDefaults(&retryOptions)

		for i := int32(63); i < 100000; i++ {
			requireWithinJitter(
				t, math.MaxInt64, calcDelay(retryOptions, i),
			)
		}
	})
}

func TestLinkRecoveryDelay(t *testing.T) {
	isFatalFn := func(err error) bool {
		return false // All errors are retryable
	}

	t.Run("LinkRecoveryDelay with specific value", func(t *testing.T) {
		var actualDelays []time.Duration
		lastTime := time.Now()

		err := Retry(context.Background(), testLogEvent, func() string { return "test" }, exported.RetryOptions{
			MaxRetries:        2,
			RetryDelay:        100 * time.Millisecond, // Normal delay (not used)
			MaxRetryDelay:     200 * time.Millisecond,
			LinkRecoveryDelay: 10 * time.Millisecond, // Special link recovery delay
		}, func(ctx context.Context, args *RetryFnArgs) error {
			now := time.Now()
			if args.I > 0 {
				actualDelays = append(actualDelays, now.Sub(lastTime))
			}
			lastTime = now

			// Simulate link recovery
			args.UseLinkRecoveryDelay()

			return errors.New("link error")
		}, isFatalFn)

		require.Error(t, err)
		require.Len(t, actualDelays, 2) // 2 retries after initial attempt

		// All delays should be approximately LinkRecoveryDelay (10ms), not exponential backoff
		for i, delay := range actualDelays {
			require.GreaterOrEqual(t, delay, 8*time.Millisecond, "delay %d too short", i)
			require.LessOrEqual(t, delay, 15*time.Millisecond, "delay %d too long", i)
		}
	})

	t.Run("LinkRecoveryDelay with negative value (no delay)", func(t *testing.T) {
		var actualDelays []time.Duration
		lastTime := time.Now()

		err := Retry(context.Background(), testLogEvent, func() string { return "test" }, exported.RetryOptions{
			MaxRetries:        2,
			RetryDelay:        100 * time.Millisecond, // Normal delay (not used)
			MaxRetryDelay:     200 * time.Millisecond,
			LinkRecoveryDelay: -1, // No delay
		}, func(ctx context.Context, args *RetryFnArgs) error {
			now := time.Now()
			if args.I > 0 {
				actualDelays = append(actualDelays, now.Sub(lastTime))
			}
			lastTime = now

			// Simulate link recovery
			args.UseLinkRecoveryDelay()

			return errors.New("link error")
		}, isFatalFn)

		require.Error(t, err)
		require.Len(t, actualDelays, 2)

		// All delays should be essentially zero (immediate retry)
		for i, delay := range actualDelays {
			require.LessOrEqual(t, delay, 5*time.Millisecond, "delay %d should be near zero", i)
		}
	})

	t.Run("LinkRecoveryDelay zero uses normal exponential backoff", func(t *testing.T) {
		var actualDelays []time.Duration
		lastTime := time.Now()

		err := Retry(context.Background(), testLogEvent, func() string { return "test" }, exported.RetryOptions{
			MaxRetries:        2,
			RetryDelay:        10 * time.Millisecond,
			MaxRetryDelay:     100 * time.Millisecond,
			LinkRecoveryDelay: 0, // Use normal exponential backoff
		}, func(ctx context.Context, args *RetryFnArgs) error {
			now := time.Now()
			if args.I > 0 {
				actualDelays = append(actualDelays, now.Sub(lastTime))
			}
			lastTime = now

			// Simulate link recovery
			args.UseLinkRecoveryDelay()

			return errors.New("link error")
		}, isFatalFn)

		require.Error(t, err)
		require.Len(t, actualDelays, 2)

		// First retry should be ~10ms * (2^1 - 1) = ~10ms with jitter
		require.GreaterOrEqual(t, actualDelays[0], 8*time.Millisecond)
		require.LessOrEqual(t, actualDelays[0], 15*time.Millisecond)

		// Second retry should be ~10ms * (2^2 - 1) = ~30ms with jitter
		require.GreaterOrEqual(t, actualDelays[1], 24*time.Millisecond)
		require.LessOrEqual(t, actualDelays[1], 40*time.Millisecond)
	})

	t.Run("Normal error without UseLinkRecoveryDelay uses exponential backoff", func(t *testing.T) {
		var actualDelays []time.Duration
		lastTime := time.Now()

		err := Retry(context.Background(), testLogEvent, func() string { return "test" }, exported.RetryOptions{
			MaxRetries:        2,
			RetryDelay:        10 * time.Millisecond,
			MaxRetryDelay:     100 * time.Millisecond,
			LinkRecoveryDelay: 5 * time.Millisecond, // Should NOT be used
		}, func(ctx context.Context, args *RetryFnArgs) error {
			now := time.Now()
			if args.I > 0 {
				actualDelays = append(actualDelays, now.Sub(lastTime))
			}
			lastTime = now

			// Do NOT call UseLinkRecoveryDelay() - this is a normal error
			return errors.New("normal error")
		}, isFatalFn)

		require.Error(t, err)
		require.Len(t, actualDelays, 2)

		// Should use exponential backoff, not LinkRecoveryDelay
		// First retry: ~10ms with jitter
		require.GreaterOrEqual(t, actualDelays[0], 8*time.Millisecond)
		require.LessOrEqual(t, actualDelays[0], 15*time.Millisecond)

		// Second retry: ~30ms with jitter
		require.GreaterOrEqual(t, actualDelays[1], 24*time.Millisecond)
		require.LessOrEqual(t, actualDelays[1], 40*time.Millisecond)
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

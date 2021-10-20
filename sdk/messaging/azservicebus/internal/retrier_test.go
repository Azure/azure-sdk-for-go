// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"
	"time"

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

func TestCyclingRetrier(t *testing.T) {
	ctx := context.Background()
	var durations []time.Duration

	timeAfter := func(duration time.Duration) <-chan time.Time {
		durations = append(durations, duration)

		ch := make(chan time.Time)
		close(ch)
		return ch
	}

	retrier := NewCyclingRetrier(CyclingRetrierOptions{
		Factor:    2,
		Min:       1,
		Max:       4,
		timeAfter: timeAfter,
	})

	// first try is free
	require.EqualValues(t, 0, retrier.CurrentTry())
	require.True(t, retrier.Try(ctx))
	require.EqualValues(t, 1, retrier.CurrentTry())
	require.Empty(t, durations)

	// retries
	require.True(t, retrier.Try(ctx))
	require.True(t, retrier.Try(ctx))
	require.True(t, retrier.Try(ctx))

	require.EqualValues(t, []time.Duration{
		1, 2, 4,
	}, durations)

	// the retrier will automatically reset since we hit our 'max' time limit.
	durations = nil

	// first try again
	require.EqualValues(t, 0, retrier.CurrentTry())
	require.True(t, retrier.Try(ctx))
	require.EqualValues(t, 1, retrier.CurrentTry())
	require.Empty(t, durations)

	// retries
	require.True(t, retrier.Try(ctx))
	require.True(t, retrier.Try(ctx))
	require.True(t, retrier.Try(ctx))

	// and...
	require.EqualValues(t, []time.Duration{
		1, 2, 4,
	}, durations)

	// and we respect cancellations
	retrier = retrier.Copy()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.False(t, retrier.Try(ctx))

	retrier = retrier.Copy()
	ctx, cancel = context.WithCancel(context.Background())

	require.True(t, retrier.Try(ctx))

	// subsequent attempts can be cancelled as well.
	cyclingRetrier := retrier.(*cyclingRetrier)
	cyclingRetrier.timeAfter = func(duration time.Duration) <-chan time.Time {
		cancel()
		// ensure that the ctx is what stopped the try.
		ch := make(chan time.Time)
		return ch
	}

	require.False(t, retrier.Try(ctx))
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLocalIdleTracker(t *testing.T) {
	t.Run("is hierarchical", func(t *testing.T) {
		idleTracker := &LocalIdleTracker{MaxDuration: time.Hour}

		parentCtx, cancelParent := context.WithCancel(context.Background())

		ctx, cancel := idleTracker.NewContextWithDeadline(parentCtx)
		defer cancel()

		require.NoError(t, ctx.Err())

		cancelParent()

		require.ErrorIs(t, ctx.Err(), context.Canceled)
	})

	t.Run("expires in MaxDuration if no previous cancel time exists", func(t *testing.T) {
		maxDuration := 2 * time.Second

		idleTracker := &LocalIdleTracker{
			MaxDuration: maxDuration,
		}

		ctx, cancel := idleTracker.NewContextWithDeadline(context.Background())
		defer cancel()
		require.Nil(t, ctx.Err())

		deadline, ok := ctx.Deadline()
		require.True(t, ok)
		require.GreaterOrEqual(t, deadline.Add(-maxDuration).UnixNano(), int64(0), "our deadline was set appropriately into the future")
	})

	t.Run("with a previous cancel", func(t *testing.T) {
		maxWait := 2 * time.Second
		idleStartTime := time.Now().Add(time.Hour)

		idleTracker := &LocalIdleTracker{
			MaxDuration: maxWait,
			IdleStart:   idleStartTime,
		}

		ctx, cancel := idleTracker.NewContextWithDeadline(context.Background())
		defer cancel()
		require.Nil(t, ctx.Err())

		deadline, ok := ctx.Deadline()
		require.True(t, ok)
		require.Equal(t, deadline.Add(-maxWait), idleStartTime, "deadline used our idle start time as the base, not time.Now()")
	})

	t.Run("user cancels", func(t *testing.T) {
		idleTracker := &LocalIdleTracker{
			MaxDuration: 30 * time.Minute,
		}

		parentCtx, cancelParent := context.WithCancel(context.Background())
		cancelParent()

		// The user cancelled here - since that cancellation is specifically for the _first_
		// message it means they didn't receive anything. If they do this for long enough the
		// link will be considered idle.

		twoHoursFromNow := time.Now().Add(2 * time.Hour)
		require.Zero(t, idleTracker.IdleStart)
		err := idleTracker.Check(parentCtx, twoHoursFromNow, context.Canceled)

		require.ErrorIs(t, err, context.Canceled)
		require.Equal(t, idleTracker.IdleStart, twoHoursFromNow, "time of first cancel is recorded (gets used as the base for future idle calculations)")

		// now we have a successful call, and it resets the idle time back to zero (ie, we're no longer in danger of being idle)
		err = idleTracker.Check(parentCtx, twoHoursFromNow, nil)
		require.NoError(t, err)
		require.Zero(t, idleTracker.IdleStart, "a successful call resets our idle tracking")

		// we also reset the idle time back to zero if there's any other error since
		// those errors will be dealt with by the recovery code.
		idleTracker.IdleStart = time.Now()

		err = idleTracker.Check(parentCtx, twoHoursFromNow, errors.New("some other error"))
		require.EqualError(t, err, "some other error")
		require.Zero(t, idleTracker.IdleStart, "an error is also considered as proof that the link is alive")
	})

	t.Run("idle deadline expires", func(t *testing.T) {
		twoHoursAndOneMinuteAgo := time.Now().Add(-time.Hour - time.Minute)
		idleTracker := &LocalIdleTracker{
			MaxDuration: time.Hour,
			IdleStart:   twoHoursAndOneMinuteAgo,
		}

		parentCtx, cancelParent := context.WithCancel(context.Background())
		defer cancelParent()

		err := idleTracker.Check(parentCtx, time.Now(), context.DeadlineExceeded)
		require.ErrorIs(t, err, localIdleError)
		require.True(t, IsLocalIdleError(err))
	})
}

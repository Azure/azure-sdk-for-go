// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorefresh

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// fakeClock is a manually advanced clock shared by a test and the cache under test.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// entry returns a value valid for 5 minutes from now, due for refresh 30 seconds before that.
func (c *fakeClock) entry(v int) Entry[int] {
	expires := c.Now().Add(5 * time.Minute)
	return Entry[int]{Value: v, ExpiresOn: expires, RefreshOn: expires.Add(-30 * time.Second)}
}

// counter hands out increasing values and counts acquisitions.
type counter struct{ n atomic.Int32 }

func (c *counter) next() int  { return int(c.n.Add(1)) }
func (c *counter) calls() int { return int(c.n.Load()) }

func newCache(clock *fakeClock, acquire AcquireFunc[int]) *Cache[int] {
	return New(acquire, &Options{Now: clock.Now})
}

// waitFor polls cond until it holds; background refreshes complete on their own goroutine.
func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	require.Eventually(t, cond, 5*time.Second, time.Millisecond)
}

func TestGetAcquiresOnceAndCaches(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	c := newCache(clock, func(context.Context) (Entry[int], error) { return clock.entry(calls.next()), nil })

	for range 3 {
		v, err := c.Get(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, v)
	}
	require.Equal(t, 1, calls.calls(), "a valid value is served from the cache")
}

func TestConcurrentColdCallersShareOneAcquisition(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	release := make(chan struct{})
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		v := calls.next()
		<-release
		return clock.entry(v), nil
	})

	const callers = 50
	results := make(chan int, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := c.Get(context.Background())
			require.NoError(t, err)
			results <- v
		}()
	}
	waitFor(t, func() bool { return calls.calls() == 1 })
	close(release)
	wg.Wait()
	close(results)

	require.Equal(t, 1, calls.calls(), "concurrent cold callers must not each acquire")
	for v := range results {
		require.Equal(t, 1, v)
	}
}

func TestConcurrentColdCallersShareOneFailure(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	release := make(chan struct{})
	boom := errors.New("boom")
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		calls.next()
		<-release
		return Entry[int]{}, boom
	})
	const callers = 20
	// release the acquisition only once every other caller is waiting on it; a caller that
	// arrives after the failure is supposed to start a new acquisition
	var waiting atomic.Int32
	c.testHookWaiting = func() { waiting.Add(1) }

	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := c.Get(context.Background())
			require.ErrorIs(t, err, boom)
		}()
	}
	waitFor(t, func() bool { return calls.calls() == 1 && waiting.Load() == callers-1 })
	close(release)
	wg.Wait()
	require.Equal(t, 1, calls.calls(), "waiters share the failure rather than each retrying in turn")

	// A failure isn't cached: the next call acquires again.
	c.testHookWaiting = nil
	_, err := c.Get(context.Background())
	require.ErrorIs(t, err, boom)
	require.Equal(t, 2, calls.calls())
}

func TestExpiryIsAuthoritative(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		v := calls.next()
		// no early refresh, so only expiry can cause a re-acquire
		e := clock.entry(v)
		e.RefreshOn = e.ExpiresOn
		return e, nil
	})

	v, err := c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)

	clock.Advance(5*time.Minute - time.Nanosecond)
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v, "still valid just before expiry")

	clock.Advance(time.Nanosecond)
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, v, "an expired value is never returned")
}

func TestNoEarlyRefreshWhenRefreshOnEqualsExpiresOn(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		calls.next()
		expires := clock.Now().Add(5 * time.Minute)
		return Entry[int]{Value: -1, ExpiresOn: expires, RefreshOn: expires}, nil
	})

	_, err := c.Get(context.Background())
	require.NoError(t, err)
	clock.Advance(4*time.Minute + 59*time.Second)
	for range 5 {
		_, err = c.Get(context.Background())
		require.NoError(t, err)
	}
	require.Equal(t, 1, calls.calls(), "a sentinel with RefreshOn == ExpiresOn is honored for its whole lifetime")

	clock.Advance(time.Second)
	_, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, calls.calls(), "a single foreground re-acquire happens at expiry")
}

func TestBackgroundRefreshDoesNotBlockCallers(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	release := make(chan struct{})
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		v := calls.next()
		if v == 2 {
			<-release
		}
		return clock.entry(v), nil
	})

	v, err := c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)

	// Enter the 30 second refresh window.
	clock.Advance(4*time.Minute + 31*time.Second)

	// The caller that triggers the refresh returns the current value immediately, even though
	// the refresh is blocked.
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)
	waitFor(t, func() bool { return calls.calls() == 2 })

	// Callers during the refresh keep getting the current value and don't start another.
	for range 10 {
		v, err = c.Get(context.Background())
		require.NoError(t, err)
		require.Equal(t, 1, v)
	}
	require.Equal(t, 2, calls.calls())

	close(release)
	waitFor(t, func() bool {
		v, err := c.Get(context.Background())
		return err == nil && v == 2
	})
	require.Equal(t, 2, calls.calls(), "the refreshed value is promoted without another acquisition")
}

func TestBackgroundRefreshFailureKeepsValueAndThrottles(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	var fail atomic.Bool
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		v := calls.next()
		if fail.Load() {
			return Entry[int]{}, errors.New("refresh failed")
		}
		return clock.entry(v), nil
	})

	v, err := c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)

	fail.Store(true)
	clock.Advance(4*time.Minute + 31*time.Second)
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)
	waitFor(t, func() bool { return calls.calls() == 2 })

	// Wait for the failed refresh to be recorded, then confirm callers keep the old value and
	// the retry is throttled.
	waitFor(t, func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.st.background != nil && c.st.background.completed()
	})
	for range 5 {
		v, err = c.Get(context.Background())
		require.NoError(t, err, "a failed background refresh never fails a caller")
		require.Equal(t, 1, v)
	}
	require.Equal(t, 2, calls.calls(), "the retry waits out the throttle")

	// After the throttle (30 seconds), but still before expiry, another refresh is attempted.
	clock.Advance(20 * time.Second)
	_, _ = c.Get(context.Background())
	require.Equal(t, 2, calls.calls())
	clock.Advance(10 * time.Second)
	_, err = c.Get(context.Background())
	require.Error(t, err, "the value has now expired, so the failure surfaces in the foreground")
}

func TestBackgroundRefreshTimeoutRetriesOnNextCall(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	var hang atomic.Bool
	c := New(func(ctx context.Context) (Entry[int], error) {
		v := calls.next()
		if hang.Load() {
			<-ctx.Done()
			return Entry[int]{}, ctx.Err()
		}
		return clock.entry(v), nil
	}, &Options{Now: clock.Now, BackgroundAcquireTimeout: 20 * time.Millisecond})

	_, err := c.Get(context.Background())
	require.NoError(t, err)

	hang.Store(true)
	clock.Advance(4*time.Minute + 31*time.Second)
	v, err := c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)
	waitFor(t, func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.st.background != nil && c.st.background.completed()
	})

	// A timed-out refresh is retried on the very next call rather than throttled.
	hang.Store(false)
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)
	waitFor(t, func() bool { return calls.calls() == 3 })
	waitFor(t, func() bool {
		v, err := c.Get(context.Background())
		return err == nil && v == 3
	})
}

func TestBackgroundRefreshIsDetachedFromCallerCancellation(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	c := newCache(clock, func(ctx context.Context) (Entry[int], error) {
		v := calls.next()
		if v == 2 {
			// the caller that triggered this refresh has already gone away
			time.Sleep(10 * time.Millisecond)
			if ctx.Err() != nil {
				return Entry[int]{}, ctx.Err()
			}
		}
		return clock.entry(v), nil
	})
	_, err := c.Get(context.Background())
	require.NoError(t, err)

	clock.Advance(4*time.Minute + 31*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	_, err = c.Get(ctx)
	require.NoError(t, err)
	cancel()

	waitFor(t, func() bool {
		v, err := c.Get(context.Background())
		return err == nil && v == 2
	})
}

func TestInvalidateIfOnlyClearsMatchingValue(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	c := newCache(clock, func(context.Context) (Entry[int], error) { return clock.entry(calls.next()), nil })

	v, err := c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v)

	c.InvalidateIf(func(cur int) bool { return cur == 99 })
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v, "a non-matching invalidation is a no-op")

	c.InvalidateIf(func(cur int) bool { return cur == 1 })
	require.Equal(t, 1, calls.calls(), "invalidation never acquires")
	v, err = c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, v, "the next call acquires a replacement")
}

func TestConcurrentInvalidationsCauseOneReplacement(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	c := newCache(clock, func(context.Context) (Entry[int], error) { return clock.entry(calls.next()), nil })

	rejected, err := c.Get(context.Background())
	require.NoError(t, err)

	// Many requests were in flight with the same value and are all rejected. Each invalidates
	// and fetches again; only the first invalidation may clear the cache.
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.InvalidateIf(func(cur int) bool { return cur == rejected })
			v, err := c.Get(context.Background())
			require.NoError(t, err)
			require.NotEqual(t, rejected, v)
		}()
	}
	wg.Wait()
	require.Equal(t, 2, calls.calls(), "one replacement, not one per rejected request")
}

func TestInvalidateIfIgnoresInProgressAcquisition(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	release := make(chan struct{})
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		v := calls.next()
		<-release
		return clock.entry(v), nil
	})

	done := make(chan int)
	go func() {
		v, _ := c.Get(context.Background())
		done <- v
	}()
	waitFor(t, func() bool { return calls.calls() == 1 })
	c.InvalidateIf(func(int) bool { return true })
	close(release)
	require.Equal(t, 1, <-done)

	v, err := c.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, v, "an invalidation made before the value existed doesn't discard it")
}

func TestWaiterHonorsItsOwnCancellation(t *testing.T) {
	clock := newFakeClock()
	release := make(chan struct{})
	defer close(release)
	c := newCache(clock, func(context.Context) (Entry[int], error) {
		<-release
		return clock.entry(1), nil
	})
	go func() { _, _ = c.Get(context.Background()) }()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	waitFor(t, func() bool {
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.st != nil
	})
	_, err := c.Get(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestCanceledAcquirerDoesNotFailOtherWaiters(t *testing.T) {
	clock, calls := newFakeClock(), &counter{}
	started := make(chan struct{})
	c := newCache(clock, func(ctx context.Context) (Entry[int], error) {
		v := calls.next()
		if v == 1 {
			close(started)
			<-ctx.Done()
			return Entry[int]{}, ctx.Err()
		}
		return clock.entry(v), nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	firstErr := make(chan error)
	go func() {
		_, err := c.Get(ctx)
		firstErr <- err
	}()
	<-started

	var waiting atomic.Int32
	c.testHookWaiting = func() { waiting.Add(1) }
	waiterResult := make(chan int)
	go func() {
		v, err := c.Get(context.Background())
		require.NoError(t, err)
		waiterResult <- v
	}()
	waitFor(t, func() bool { return waiting.Load() == 1 })
	cancel()

	require.ErrorIs(t, <-firstErr, context.Canceled)
	require.Equal(t, 2, <-waiterResult, "the waiter retries instead of inheriting the cancellation")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Package autorefresh provides a cache for a single expiring value that is refreshed in the
// background shortly before it expires.
package autorefresh

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Entry is a cached value together with the times that govern its lifetime.
type Entry[T any] struct {
	// Value is the cached value.
	Value T

	// ExpiresOn is when Value stops being valid. Once it has passed, Value is never returned
	// again and the next call to [Cache.Get] acquires a replacement in the foreground.
	ExpiresOn time.Time

	// RefreshOn is when a background refresh should start so that a replacement is ready
	// before ExpiresOn. Set it equal to ExpiresOn for a value that must not be refreshed early.
	RefreshOn time.Time
}

// AcquireFunc acquires a new value for a [Cache].
type AcquireFunc[T any] func(ctx context.Context) (Entry[T], error)

// Options contains the optional configuration for a [Cache].
type Options struct {
	// BackgroundAcquireTimeout bounds a background refresh. When it elapses, the cache keeps
	// its current, still-valid value and starts another refresh on the next call.
	// The default is 30 seconds.
	BackgroundAcquireTimeout time.Duration

	// Now returns the current time. The default is [time.Now]. It exists for tests.
	Now func() time.Time
}

const defaultBackgroundAcquireTimeout = 30 * time.Second

// maxCancellationRetries is how many times a caller whose context can't be canceled retries
// after an acquisition started by another caller was canceled.
const maxCancellationRetries = 3

// Cache is a thread-safe cache for a single expiring value.
//
// Only one acquisition runs at a time. Concurrent callers that find the cache empty or expired
// wait for that acquisition and share its result, including its error, rather than each
// acquiring their own. When the value reaches its RefreshOn time, the caller that notices
// starts a refresh in the background and returns the current value without waiting, as do
// the callers after it while the refresh runs. A background refresh that fails or times out
// leaves the current value in place until it expires.
//
// The state machine is a port of AutoRefreshingCache from the .NET Azure Storage SDK, which is
// itself modeled on the access token cache in Azure.Core.
type Cache[T any] struct {
	acquire                  AcquireFunc[T]
	backgroundAcquireTimeout time.Duration
	now                      func() time.Time

	mu sync.Mutex
	// st is nil when the cache is empty, either before the first call or after invalidation.
	st *state[T]

	// testHookWaiting, when set, is called just before a caller waits on an acquisition.
	testHookWaiting func()
}

// New creates a Cache that calls acquire to obtain its value.
func New[T any](acquire AcquireFunc[T], opts *Options) *Cache[T] {
	c := &Cache[T]{
		acquire:                  acquire,
		backgroundAcquireTimeout: defaultBackgroundAcquireTimeout,
		now:                      time.Now,
	}
	if opts != nil {
		if opts.BackgroundAcquireTimeout > 0 {
			c.backgroundAcquireTimeout = opts.BackgroundAcquireTimeout
		}
		if opts.Now != nil {
			c.now = opts.Now
		}
	}
	return c
}

// Get returns the cached value, acquiring it first when the cache is empty, the value has
// expired or its last acquisition failed. ctx bounds how long the caller waits; it is also
// passed to an acquisition this caller runs in the foreground.
func (c *Cache[T]) Get(ctx context.Context) (T, error) {
	retries := maxCancellationRetries
	for {
		shouldAcquire, st := c.evaluate()
		if shouldAcquire {
			if st.background != nil {
				// The current value is still valid but due for a refresh: return it, and refresh
				// in the background so that neither this caller nor the ones after it wait.
				current := st.current.entry
				go c.acquireInBackground(context.WithoutCancel(ctx), st.background, current)
				return current.Value, nil
			}
			// This caller runs the acquisition; everyone else waits on st.current.
			c.acquireInForeground(ctx, st.current)
		} else if c.testHookWaiting != nil && !st.current.completed() {
			c.testHookWaiting()
		}

		select {
		case <-st.current.done:
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}

		if st.current.canceled && ctx.Err() == nil {
			// The acquisition belonged to a caller whose context was canceled, which says
			// nothing about this caller, so try again with a fresh acquisition.
			retries--
			if ctx.Done() == nil && retries <= 0 {
				var zero T
				return zero, st.current.err
			}
			continue
		}
		if st.current.err != nil {
			var zero T
			return zero, st.current.err
		}
		return st.current.entry.Value, nil
	}
}

// InvalidateIf empties the cache when its current value satisfies isCurrent, so that the next
// call to Get acquires a replacement. It is a no-op when the cache is empty, an acquisition is
// in progress, or the current value has already been replaced. Checking the value the caller
// used, rather than invalidating unconditionally, keeps many callers that were rejected with the
// same value from discarding each other's replacements.
//
// InvalidateIf never acquires a value.
func (c *Cache[T]) InvalidateIf(isCurrent func(T) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.st != nil && c.st.current.succeeded() && isCurrent(c.st.current.entry.Value) {
		c.st = nil
	}
}

// evaluate decides whether the caller should acquire a value. When it returns true and the
// returned state has a background flight, the caller starts a background refresh; otherwise the
// caller acquires into the state's current flight.
func (c *Cache[T]) evaluate() (bool, *state[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// First call, or the first after an invalidation.
	if c.st == nil {
		c.st = &state[T]{current: newFlight[T]()}
		return true, c.st
	}

	// An acquisition is in progress; wait for it.
	if !c.st.current.completed() {
		if c.st.background != nil {
			// A background result can't be promoted over an acquisition still in flight.
			c.st = &state[T]{current: c.st.current}
		}
		return false, c.st
	}

	now := c.now()

	// A background refresh finished with a value that is still valid; promote it.
	if c.st.background != nil && c.st.background.completed() && c.st.background.entry.ExpiresOn.After(now) {
		c.st = &state[T]{current: c.st.background}
	}

	// The current value failed or expired; it is re-acquired in the foreground.
	if !c.st.current.succeeded() || !now.Before(c.st.current.entry.ExpiresOn) {
		c.st = &state[T]{current: newFlight[T]()}
		return true, c.st
	}

	// The current value is valid but due for a refresh, and none is running yet.
	if !now.Before(c.st.current.entry.RefreshOn) && c.st.background == nil {
		c.st = &state[T]{current: c.st.current, background: newFlight[T]()}
		return true, c.st
	}

	return false, c.st
}

// acquireInForeground runs the acquisition for f on the caller's goroutine.
func (c *Cache[T]) acquireInForeground(ctx context.Context, f *flight[T]) {
	entry, err := c.acquire(ctx)
	if err != nil {
		f.fail(err, ctx.Err() != nil)
		return
	}
	f.complete(entry)
}

// acquireInBackground refreshes into f with a bounded timeout. The background flight always
// completes with a usable value: on failure it holds the current value, rescheduled for another
// refresh attempt. A timeout retries on the next call; any other failure waits out the timeout
// again before retrying so a failing dependency isn't hit on every call.
func (c *Cache[T]) acquireInBackground(ctx context.Context, f *flight[T], current Entry[T]) {
	ctx, cancel := context.WithTimeout(ctx, c.backgroundAcquireTimeout)
	defer cancel()

	entry, err := c.acquire(ctx)
	switch {
	case err == nil:
		f.complete(entry)
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		current.RefreshOn = c.now()
		f.complete(current)
	default:
		current.RefreshOn = c.now().Add(c.backgroundAcquireTimeout)
		f.complete(current)
	}
}

// state is an immutable snapshot of the cache. Every transition replaces it.
type state[T any] struct {
	current    *flight[T]
	background *flight[T]
}

// flight is one acquisition. done is closed once it has completed, after which its other
// fields are read-only.
type flight[T any] struct {
	done     chan struct{}
	entry    Entry[T]
	err      error
	canceled bool
}

func newFlight[T any]() *flight[T] {
	return &flight[T]{done: make(chan struct{})}
}

func (f *flight[T]) complete(entry Entry[T]) {
	f.entry = entry
	close(f.done)
}

func (f *flight[T]) fail(err error, canceled bool) {
	f.err, f.canceled = err, canceled
	close(f.done)
}

func (f *flight[T]) completed() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

func (f *flight[T]) succeeded() bool {
	return f.completed() && f.err == nil
}

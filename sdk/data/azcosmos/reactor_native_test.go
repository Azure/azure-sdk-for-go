// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The reactor owns a completion queue and a goroutine blocked in C. These cover its lifetime
// without a service, which the emulator tests cannot: they need an account to talk to, and a leak
// or a hang here would show up there as a timeout rather than as itself.

func TestReactorStartsAndStops(t *testing.T) {
	d := &nativeDriver{}
	require.NoError(t, d.buildRuntime())
	t.Cleanup(func() { _ = d.close() })

	r, err := newReactor(d.runtime)
	require.NoError(t, err)
	require.NotNil(t, r.queue)

	r.close()
	require.Nil(t, r.queue, "close must not leave a dangling queue")
}

// close has to be safe to call more than once, because the driver's close path calls it and a
// caller may already have.
func TestReactorCloseIsIdempotent(t *testing.T) {
	d := &nativeDriver{}
	require.NoError(t, d.buildRuntime())
	t.Cleanup(func() { _ = d.close() })

	r, err := newReactor(d.runtime)
	require.NoError(t, err)

	r.close()
	require.NotPanics(t, r.close)
	require.NotPanics(t, (*reactor)(nil).close, "a nil reactor is what an unopened driver has")
}

// Repeated cycles catch a queue or goroutine that close forgets to release, which a single cycle
// would not. Each iteration also outlives at least one wait timeout, so the loop is exercised
// rather than only started and stopped.
func TestReactorLifecycleIsRepeatable(t *testing.T) {
	d := &nativeDriver{}
	require.NoError(t, d.buildRuntime())
	t.Cleanup(func() { _ = d.close() })

	for range 10 {
		r, err := newReactor(d.runtime)
		require.NoError(t, err)
		r.close()
	}
}

// close must return rather than block, even while the loop is parked in the C wait call. Shutting
// the queue down is what wakes it; without that this would take a full wait timeout per close.
func TestReactorCloseDoesNotWaitOutTheTimeout(t *testing.T) {
	d := &nativeDriver{}
	require.NoError(t, d.buildRuntime())
	t.Cleanup(func() { _ = d.close() })

	waitStarted := make(chan struct{})
	var waitStartedOnce sync.Once
	r, err := newReactorWithWait(d.runtime, 10_000, func() {
		waitStartedOnce.Do(func() { close(waitStarted) })
	})
	require.NoError(t, err)
	<-waitStarted

	done := make(chan struct{})
	go func() {
		r.close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("close blocked; the queue shutdown did not wake the wait call")
	}
}

// A queue cannot be created without a runtime, and reporting that as an error rather than
// returning a nil reactor is what keeps the failure at construction.
func TestReactorRequiresARuntime(t *testing.T) {
	_, err := newReactor(nil)
	require.Error(t, err)

	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientError, cosmosErr.Code)
}

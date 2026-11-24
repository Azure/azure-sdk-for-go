// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCursorReactorCorrelatesConcurrentDelivery(t *testing.T) {
	r := &cursorReactor{pending: make(map[uintptr]chan cursorDelivery)}
	const count = completionBatch * 3
	cookies := make([]uintptr, count)
	waiters := make([]chan cursorDelivery, count)
	expected := make([]error, count)
	for i := range cookies {
		cookie, waiter, err := r.register()
		require.NoError(t, err)
		require.NotZero(t, cookie)
		cookies[i], waiters[i] = cookie, waiter
		expected[i] = errors.New("unique delivery")
	}
	var wg sync.WaitGroup
	for i := count - 1; i >= 0; i-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.deliver(cursorDelivery{err: expected[i]}, cookies[i])
		}()
	}
	wg.Wait()
	for i, waiter := range waiters {
		select {
		case got := <-waiter:
			require.Same(t, expected[i], got.err)
		case <-time.After(time.Second):
			t.Fatal("completion was lost or sent to a different waiter")
		}
	}
	require.Empty(t, r.pending)
}

func TestCursorReactorRejectionAndFailure(t *testing.T) {
	r := &cursorReactor{pending: make(map[uintptr]chan cursorDelivery)}
	cookie, rejected, err := r.register()
	require.NoError(t, err)
	r.forget(cookie)
	r.deliver(cursorDelivery{err: errors.New("late result")}, cookie)
	require.Empty(t, rejected)

	_, first, err := r.register()
	require.NoError(t, err)
	_, second, err := r.register()
	require.NoError(t, err)
	failure := &Error{Code: CodeClientError, Message: "queue failure"}
	r.fail(failure)
	gotFirst, gotSecond := <-first, <-second
	require.Equal(t, failure, gotFirst.err)
	require.Equal(t, failure, gotSecond.err)
	require.NotSame(t, gotFirst.err, gotSecond.err, "callers must not share mutable errors")
	require.Empty(t, r.pending)
	_, waiter, err := r.register()
	require.Equal(t, failure, err)
	require.Nil(t, waiter)
}

func TestCursorReactorRefusesCookieReuse(t *testing.T) {
	r := &cursorReactor{pending: make(map[uintptr]chan cursorDelivery), next: ^uintptr(0)}
	cookie, waiter, err := r.register()
	require.ErrorContains(t, err, "correlation space exhausted")
	require.Zero(t, cookie)
	require.Nil(t, waiter)
}

func TestCursorReactorSharedLifetime(t *testing.T) {
	client := newTestContainer(t).database.client
	driver := client.driver
	first, err := driver.testIdleCursor()
	require.NoError(t, err)
	second, err := driver.testIdleCursor()
	require.NoError(t, err)
	require.Same(t, first.reactor, second.reactor)
	reactor := first.reactor
	first.close()
	require.NotNil(t, reactor.queue, "closing a cursor must not free its sibling's queue")
	require.NotNil(t, second.queue)
	require.NoError(t, client.Close())
	require.Nil(t, second.queue)
	require.Nil(t, reactor.queue)
	require.Empty(t, driver.cursors)
	require.NotPanics(t, reactor.close)
	_, _, err = reactor.register()
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
	_, err = driver.queryReactor()
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

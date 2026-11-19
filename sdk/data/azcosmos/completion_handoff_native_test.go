// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"runtime/cgo"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func handoffPending(t *testing.T) (*pendingOperation, uintptr) {
	t.Helper()
	p := &pendingOperation{result: make(chan completionResult, 1), delivered: make(chan struct{})}
	handle := cgo.NewHandle(p)
	t.Cleanup(handle.Delete)
	return p, uintptr(handle)
}

// This test also compiles on the original implementation: it fails because translation blocks
// the collector, not because it assumes the new ownership representation.
func TestBlockedCallerTranslationDoesNotBlockCollector(t *testing.T) {
	started, resume := make(chan struct{}), make(chan struct{})
	restore := stubBlockingDiagnosticsCopy(started, resume)
	defer restore()
	first, cookieA := handoffPending(t)
	second, cookieB := handoffPending(t)
	var workers sync.WaitGroup
	var once sync.Once
	defer func() {
		once.Do(func() { close(resume) })
		workers.Wait()
	}()
	firstDone := make(chan error, 1)
	workers.Add(2)
	go func() {
		defer workers.Done()
		_, err := awaitOperationResult(context.Background(), first)
		firstDone <- err
	}()
	go func() {
		defer workers.Done()
		for _, cookie := range []uintptr{cookieA, cookieB} {
			_, release := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "")
			release()
		}
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("first result did not reach preparation")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := awaitOperationResult(ctx, second)
	require.NoError(t, err, "blocked preparation must not prevent another completion's delivery")
	require.Equal(t, 200, result.response.StatusCode)
	select {
	case <-firstDone:
		t.Fatal("first caller returned before preparation was released")
	default:
	}
	once.Do(func() { close(resume) })
	select {
	case err := <-firstDone:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("first caller did not finish preparation")
	}
}

func TestCompletionHandoffSurvivesBatchReuse(t *testing.T) {
	for _, verbosity := range []DiagnosticsVerbosity{DiagnosticsVerbositySummary, DiagnosticsVerbosityDetailed} {
		t.Run(verbosity.String(), func(t *testing.T) {
			diagnostics := &Diagnostics{AttemptCount: 3}
			restore, seen := stubCopyDiagnosticsForCompletionCapturing(diagnostics, nil)
			defer restore()
			pending, cookie := handoffPending(t)
			cleared, release := dispatchSyntheticCompletion(&reactor{verbosity: verbosity}, cookie, 201, `{"id":"one"}`)
			defer func() { release() }()
			require.True(t, cleared, "the poll slot must lose ownership before publication")
			require.Empty(t, *seen, "the collector must not prepare diagnostics")
			moved := (<-pending.result).native
			require.NotNil(t, moved)
			pending.result <- completionResult{native: moved}
			result, err := awaitOperationResult(context.Background(), pending)
			require.NoError(t, err)
			require.Nil(t, result.native)
			require.Nil(t, moved.record.body)
			require.Nil(t, moved.record.backing)
			require.Equal(t, []DiagnosticsVerbosity{verbosity}, *seen)
			require.Same(t, diagnostics, result.response.Diagnostics)
			require.Equal(t, uint32(3), result.response.AttemptCount)
			require.Equal(t, 201, result.response.StatusCode)
			// Release borrowed fixture bytes before inspecting the Go copy.
			release()
			release = func() {}
			require.JSONEq(t, `{"id":"one"}`, string(result.body))
		})
	}
}

func TestCompletionHandoffCancellationAndLateDiscard(t *testing.T) {
	for range 100 {
		pending, cookie := handoffPending(t)
		ctx, cancel := context.WithCancel(context.Background())
		cleared, release := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "")
		release()
		require.True(t, cleared)
		cancel()
		result, err := awaitOperationResult(ctx, pending)
		require.NoError(t, err, "a delivered outcome wins over cancellation")
		require.Equal(t, 200, result.response.StatusCode)
		require.Nil(t, result.native)
	}
	restore, seen := stubCopyDiagnosticsForCompletionCapturing(nil, nil)
	defer restore()
	pending, cookie := handoffPending(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := awaitOperationResult(ctx, pending)
	require.ErrorIs(t, err, context.Canceled)
	cleared, release := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "")
	release()
	require.True(t, cleared)
	require.Empty(t, *seen, "late abandoned results must be freed without translation")
	require.Empty(t, pending.result)
	select {
	case <-pending.delivered:
	default:
		t.Fatal("late discard must finish before signalling cleanup")
	}
}

func TestCompletionHandoffUnknownCookiesAndMultipleBatches(t *testing.T) {
	wrongType := cgo.NewHandle("not an operation")
	defer wrongType.Delete()
	deleted := cgo.NewHandle(&pendingOperation{})
	deleted.Delete()
	for _, cookie := range []uintptr{0, uintptr(deleted), uintptr(wrongType)} {
		cleared, release := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "unclaimed")
		release()
		require.False(t, cleared, "unclaimed slots retain ownership for batch cleanup")
	}
	for range completionBatch * 3 {
		pending, cookie := handoffPending(t)
		cleared, release := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "discard")
		require.True(t, cleared)
		result := <-pending.result
		owned := result.native
		require.NotNil(t, owned)
		result.release()
		require.Nil(t, owned.record.body)
		require.Nil(t, owned.record.backing)
		release()
	}
}

func TestCompletionHandoffConcurrentAbandonment(t *testing.T) {
	for range 100 {
		pending, cookie := handoffPending(t)
		var workers sync.WaitGroup
		workers.Add(2)
		go func() {
			defer workers.Done()
			_, release := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "")
			release()
		}()
		go func() {
			defer workers.Done()
			if result, ok := pending.abandon(); ok {
				result.release()
			}
		}()
		workers.Wait()
		require.Empty(t, pending.result)
		select {
		case <-pending.delivered:
		default:
			t.Fatal("delivery did not finish")
		}
	}
}

func TestCloseWaitsForCallerPreparation(t *testing.T) {
	client := sharedTestClient(t, nil, OperationOptions{})
	started, resume := make(chan struct{}), make(chan struct{})
	restore := stubBlockingDiagnosticsCopy(started, resume)
	defer restore()
	var once sync.Once
	defer once.Do(func() { close(resume) })
	pending, cookie := handoffPending(t)
	cleared, releaseFixture := dispatchSyntheticCompletion(&reactor{}, cookie, 200, "")
	releaseFixture()
	require.True(t, cleared)
	done := make(chan error, 1)
	go func() {
		release, err := client.acquire()
		if err == nil {
			defer release()
			_, err = awaitOperationResult(context.Background(), pending)
		}
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("caller did not start preparation")
	}
	nativeClose := make(chan struct{})
	client.driver.beforePendingWait = func() { close(nativeClose) }
	closing := make(chan struct{}, 2)
	closed := make(chan error, 2)
	for range 2 {
		go func() {
			closing <- struct{}{}
			closed <- client.Close()
		}()
	}
	<-closing
	<-closing
	select {
	case <-nativeClose:
		t.Fatal("native teardown began before caller preparation finished")
	case <-time.After(25 * time.Millisecond):
	}
	once.Do(func() { close(resume) })
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("caller did not finish")
	}
	for range 2 {
		select {
		case err := <-closed:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("concurrent Close did not finish")
		}
	}
}

func BenchmarkCompletionHandoff(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		pending := &pendingOperation{result: make(chan completionResult, 1)}
		handle := cgo.NewHandle(pending)
		_, release := dispatchSyntheticCompletion(&reactor{}, uintptr(handle), 200, "")
		result, err := awaitOperationResult(context.Background(), pending)
		if err != nil || result.response.StatusCode != 200 {
			b.Fatal(result, err)
		}
		release()
		handle.Delete()
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func sharedTestClient(t *testing.T, runtime *Runtime, options OperationOptions) *Client {
	t.Helper()
	key, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", key, &ClientOptions{
		Runtime: runtime, Operation: options,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	return client
}

func TestRuntimeSiblingIsolationAndOwnership(t *testing.T) {
	runtime, err := NewRuntime(&RuntimeOptions{ApplicationID: "shared-test"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	first := sharedTestClient(t, runtime, OperationOptions{})
	second := sharedTestClient(t, runtime, OperationOptions{})
	require.Same(t, first.runtime, second.runtime)
	require.NotSame(t, first.driver.account, second.driver.account)
	require.NotSame(t, first.driver.reactor.queue, second.driver.reactor.queue)
	require.NoError(t, first.Close())
	release, err := second.acquire()
	require.NoError(t, err)
	release()
	require.NoError(t, runtime.Close())
	_, err = second.acquire()
	require.Error(t, err)
	require.Error(t, runtime.SetOperationOptions(OperationOptions{}))
	_, err = NewClientWithKey("https://myaccount.documents.azure.com", KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: runtime})
	require.Error(t, err)
}

func TestRuntimeCloseDrainsAndRejectsNewWork(t *testing.T) {
	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	client := sharedTestClient(t, runtime, OperationOptions{})
	release, err := client.acquire()
	require.NoError(t, err)
	done := make(chan error, 1)
	go func() { done <- runtime.Close() }()
	require.Eventually(t, func() bool {
		runtime.mu.Lock()
		defer runtime.mu.Unlock()
		return runtime.closing
	}, time.Second, time.Millisecond)
	_, err = runtime.acquire()
	require.Error(t, err)
	select {
	case <-done:
		t.Fatal("Close freed an admitted operation's resources")
	default:
	}
	release()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("Close did not finish after operation drain")
	}
}

func TestRuntimeSnapshotTimeoutHierarchyAndReplacement(t *testing.T) {
	runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: 8 * time.Second}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	client := sharedTestClient(t, runtime, OperationOptions{})
	oldCtx, _, releaseOld, err := client.driver.snapshot(context.Background(), OperationOptions{})
	require.NoError(t, err)
	defer releaseOld()
	oldDeadline, ok := oldCtx.Deadline()
	require.True(t, ok)
	require.NoError(t, runtime.SetOperationOptions(OperationOptions{EndToEndTimeout: 2 * time.Second}))
	newCtx, _, releaseNew, err := client.driver.snapshot(context.Background(), OperationOptions{})
	require.NoError(t, err)
	defer releaseNew()
	newDeadline, ok := newCtx.Deadline()
	require.True(t, ok)
	require.Greater(t, oldDeadline.Sub(newDeadline), 5*time.Second)
	require.Equal(t, oldDeadline, func() time.Time { d, _ := oldCtx.Deadline(); return d }())
	account := sharedTestClient(t, runtime, OperationOptions{EndToEndTimeout: 4 * time.Second})
	accountCtx, _, releaseAccount, err := account.driver.snapshot(context.Background(), OperationOptions{})
	require.NoError(t, err)
	defer releaseAccount()
	deadline, ok := accountCtx.Deadline()
	require.True(t, ok)
	require.InDelta(t, 4*time.Second, time.Until(deadline), float64(200*time.Millisecond))
	requestCtx, _, releaseRequest, err := account.driver.snapshot(context.Background(), OperationOptions{EndToEndTimeout: 50 * time.Millisecond})
	require.NoError(t, err)
	defer releaseRequest()
	deadline, _ = requestCtx.Deadline()
	require.LessOrEqual(t, time.Until(deadline), 50*time.Millisecond)
}

func TestRuntimeConcurrentAttachmentUpdateAndClose(t *testing.T) {
	for range 10 {
		runtime, err := NewRuntime(nil)
		require.NoError(t, err)
		client := sharedTestClient(t, runtime, OperationOptions{})
		var wg sync.WaitGroup
		for i := range 6 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				switch i % 3 {
				case 0:
					_ = runtime.SetOperationOptions(OperationOptions{MaxFailoverRetryCount: to(uint32(i))})
				case 1:
					_ = client.Close()
				case 2:
					attached, err := NewClientWithKey("https://myaccount.documents.azure.com", KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: runtime})
					if err == nil {
						_ = attached.Close()
					}
				}
			}()
		}
		require.NoError(t, runtime.Close())
		wg.Wait()
		require.Empty(t, runtime.clients)
	}
}

func TestClientAliasIdentityAndOptionOwnership(t *testing.T) {
	runtime, err := NewRuntime(&RuntimeOptions{ApplicationID: "shared-test"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	key := KeyCredential{accountKey: emulatorKey}
	_, err = NewClientWithKey("https://myaccount.documents.azure.com", key, &ClientOptions{Runtime: runtime, ApplicationID: "different"})
	require.ErrorContains(t, err, "ApplicationID")
	_, err = NewClientWithKey("https://myaccount.documents.azure.com", key, &ClientOptions{
		Runtime: runtime, EnableContentResponseOnWrite: to(true),
		Operation: OperationOptions{EnableContentResponseOnWrite: to(false)},
	})
	require.ErrorContains(t, err, "conflicting")
	headers := map[string]string{"x-custom": "original"}
	client := sharedTestClient(t, runtime, OperationOptions{CustomHeaders: headers})
	headers["x-custom"] = "changed"
	require.Equal(t, "original", client.options.Operation.CustomHeaders["x-custom"])
}

func TestRuntimeUpdateCopiesAndRejectsInvalidReplacement(t *testing.T) {
	defaults := OperationOptions{EndToEndTimeout: 5 * time.Second, CustomHeaders: map[string]string{"x-test": "original"}}
	runtime, err := NewRuntime(&RuntimeOptions{Operation: defaults})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	client := sharedTestClient(t, runtime, OperationOptions{})
	require.NoError(t, runtime.SetOperationOptions(defaults))
	defaults.EndToEndTimeout = time.Second
	defaults.CustomHeaders["x-test"] = "changed"
	require.Error(t, runtime.SetOperationOptions(OperationOptions{EndToEndTimeout: -1}))
	ctx, _, release, err := client.driver.snapshot(context.Background(), OperationOptions{})
	require.NoError(t, err)
	defer release()
	deadline, ok := ctx.Deadline()
	require.True(t, ok)
	require.InDelta(t, 5*time.Second, time.Until(deadline), float64(200*time.Millisecond))
}

func TestRuntimeUpdateDoesNotRestartInitializationBudget(t *testing.T) {
	runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: time.Second}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	client, err := NewClientWithKey(localBlackholeEndpoint(t), KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: runtime})
	require.NoError(t, err)
	container, err := client.NewContainer("db", "items")
	require.NoError(t, err)
	done := make(chan error, 1)
	start := time.Now()
	go func() {
		_, err := container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "id", nil)
		done <- err
	}()
	require.Eventually(t, func() bool {
		client.driver.mu.Lock()
		defer client.driver.mu.Unlock()
		return client.driver.creating != nil
	}, time.Second, time.Millisecond)
	require.NoError(t, runtime.SetOperationOptions(OperationOptions{EndToEndTimeout: 10 * time.Second}))
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Less(t, time.Since(start), 3*time.Second, "update replaced the admitted timeout")
	case <-time.After(3 * time.Second):
		t.Fatal("runtime update restarted the initialization budget")
	}
}

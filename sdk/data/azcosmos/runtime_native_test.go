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
	client, err := NewClientWithKey("https://myaccount.documents.azure.com", KeyCredential{accountKey: emulatorKey},
		&ClientOptions{Runtime: runtime, Operation: options})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	return client
}

func TestDefaultRuntimeSharedAcrossConcurrentClients(t *testing.T) {
	const count = 12
	results := make(chan *Client, count)
	failures := make(chan error, count)
	var wg sync.WaitGroup
	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client, err := NewClientWithKey("https://default.documents.azure.com", KeyCredential{accountKey: emulatorKey}, nil)
			if err != nil {
				failures <- err
				return
			}
			results <- client
		}()
	}
	wg.Wait()
	close(results)
	close(failures)
	for err := range failures {
		require.NoError(t, err)
	}
	var runtime *Runtime
	for client := range results {
		if runtime == nil {
			runtime = client.runtime
		}
		require.Same(t, runtime, client.runtime)
		require.NoError(t, client.Close())
	}
	require.NotNil(t, runtime.native)
	client := sharedTestClient(t, nil, OperationOptions{})
	require.Same(t, runtime, client.runtime)
}

func TestRuntimeOwnerClosePreservesAttachedClients(t *testing.T) {
	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	first := sharedTestClient(t, runtime, OperationOptions{})
	second := sharedTestClient(t, runtime, OperationOptions{})
	require.NoError(t, runtime.Close())
	require.NotNil(t, runtime.native)
	_, err = NewClientWithKey("https://another.documents.azure.com", KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: runtime})
	require.Error(t, err)
	require.NoError(t, first.Close())
	release, err := second.acquire()
	require.NoError(t, err)
	release()
	require.NotNil(t, runtime.native)
	require.NoError(t, second.Close())
	require.Nil(t, runtime.native)
	require.NoError(t, runtime.Close())
}

func TestRuntimeConstructionCopiesDefaults(t *testing.T) {
	timeout := 5 * time.Second
	options := RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: &timeout}}
	runtime, err := NewRuntime(&options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	timeout = time.Second
	client := sharedTestClient(t, runtime, OperationOptions{})
	ctx, _, release, err := client.driver.snapshot(context.Background(), OperationOptions{})
	require.NoError(t, err)
	defer release()
	deadline, ok := ctx.Deadline()
	require.True(t, ok)
	require.InDelta(t, float64(5*time.Second), float64(time.Until(deadline)), float64(200*time.Millisecond))
}

func TestRuntimeTimeoutPresenceAndClamp(t *testing.T) {
	runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: to(8 * time.Second)}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	client := sharedTestClient(t, runtime, OperationOptions{EndToEndTimeout: to(4 * time.Second)})
	for _, tt := range []struct {
		timeout *time.Duration
		want    time.Duration
	}{
		{nil, 4 * time.Second}, {to(time.Duration(0)), time.Second}, {to(time.Millisecond), time.Second}, {to(2 * time.Second), 2 * time.Second},
	} {
		ctx, _, release, err := client.driver.snapshot(context.Background(), OperationOptions{EndToEndTimeout: tt.timeout})
		require.NoError(t, err)
		deadline, ok := ctx.Deadline()
		require.True(t, ok)
		require.InDelta(t, float64(tt.want), float64(time.Until(deadline)), float64(200*time.Millisecond))
		release()
	}
	parent, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	ctx, _, release, err := client.driver.snapshot(parent, OperationOptions{EndToEndTimeout: to(time.Duration(0))})
	require.NoError(t, err)
	defer release()
	want, _ := parent.Deadline()
	got, _ := ctx.Deadline()
	require.Equal(t, want, got)
}

func TestRuntimeContextDeadlineFallback(t *testing.T) {
	for _, tt := range []struct {
		name     string
		runtime  *time.Duration
		client   *time.Duration
		request  *time.Duration
		deadline time.Duration
		wantMS   int64
	}{
		{name: "no deadline", wantMS: -1},
		{name: "context fallback", deadline: 3 * time.Second, wantMS: 3000},
		{name: "subsecond clamp", deadline: 100 * time.Millisecond, wantMS: 1000},
		{name: "shorter runtime timeout", runtime: to(time.Second), deadline: 3 * time.Second, wantMS: 1000},
		{name: "explicit runtime preserved", runtime: to(5 * time.Second), deadline: 3 * time.Second, wantMS: 5000},
		{name: "explicit client preserved", client: to(5 * time.Second), deadline: 3 * time.Second, wantMS: 5000},
		{name: "explicit request preserved", request: to(5 * time.Second), deadline: 3 * time.Second, wantMS: 5000},
	} {
		t.Run(tt.name, func(t *testing.T) {
			runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: tt.runtime}})
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, runtime.Close()) })
			client := sharedTestClient(t, runtime, OperationOptions{EndToEndTimeout: tt.client})
			ctx := t.Context()
			if tt.deadline != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.deadline)
				defer cancel()
			}
			_, timeout, release, err := client.driver.captureSnapshot(ctx, OperationOptions{EndToEndTimeout: tt.request})
			require.NoError(t, err)
			defer release()
			if tt.name == "context fallback" {
				require.InDelta(t, tt.wantMS, timeout, 200, "assert the native snapshot timeout, not just the Go deadline")
				require.LessOrEqual(t, timeout, tt.wantMS)
			} else {
				require.Equal(t, tt.wantMS, timeout)
			}
		})
	}
}

func TestClientNativeDefaultsReusedAcrossSnapshots(t *testing.T) {
	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	client := sharedTestClient(t, runtime, OperationOptions{
		EndToEndTimeout: to(4 * time.Second), CustomHeaders: map[string]string{"x-test": "original"},
	})
	driver := client.driver
	options := driver.clientOptions
	require.NotNil(t, options)
	// Mutating the internal Go copy proves snapshots use the construction-time native allocation.
	driver.cfg.options.Operation.EndToEndTimeout = to(8 * time.Second)
	for range 4 {
		t.Run("snapshot", func(t *testing.T) {
			t.Parallel()
			_, timeout, release, err := driver.captureSnapshot(t.Context(), OperationOptions{})
			require.NoError(t, err)
			defer release()
			require.Equal(t, int64(4000), timeout)
			require.Same(t, options, driver.clientOptions)
		})
	}
	t.Cleanup(func() {
		driver.beforePendingWait = func() {
			require.Same(t, options, driver.clientOptions, "defaults must remain alive until work drains")
		}
		require.NoError(t, client.Close())
		require.Nil(t, driver.clientOptions)
		require.Nil(t, driver.releaseClientOptions)
	})
}

func TestRuntimeSamplingValidation(t *testing.T) {
	for _, duration := range []time.Duration{0, time.Second - 1, time.Minute + 1} {
		_, err := NewRuntime(&RuntimeOptions{CPURefreshInterval: &duration})
		require.ErrorContains(t, err, "CPURefreshInterval")
	}
	runtime, err := NewRuntime(&RuntimeOptions{CPURefreshInterval: to(time.Second)})
	require.NoError(t, err)
	require.NoError(t, runtime.Close())
}

func TestClientIdentityAndOptionOwnership(t *testing.T) {
	runtime, err := NewRuntime(&RuntimeOptions{ApplicationID: "shared-test"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	headers := map[string]string{"x-custom": "original"}
	client := sharedTestClient(t, runtime, OperationOptions{CustomHeaders: headers})
	headers["x-custom"] = "changed"
	require.Equal(t, "original", client.options.Operation.CustomHeaders["x-custom"])
}

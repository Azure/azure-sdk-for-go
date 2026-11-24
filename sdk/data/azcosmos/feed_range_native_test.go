// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build azcosmos_native_preview && cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNativeFeedRangeValueContracts(t *testing.T) {
	left, err := NewFeedRange("", "80")
	require.NoError(t, err)
	right, err := NewFeedRange("80", "FF")
	require.NoError(t, err)
	full, err := NewFeedRange("", "FF")
	require.NoError(t, err)
	parsed, err := ParseFeedRange(left.String())
	require.NoError(t, err)
	require.Equal(t, left, parsed, "native handles and strings have been freed")
	require.Equal(t, "", left.MinInclusive())
	require.Equal(t, "80", left.MaxExclusive())
	require.False(t, left.IsLogicalPartition())
	point, err := NewFeedRange("80", "80")
	require.NoError(t, err, "equal native bounds denote a single EPK")
	require.Equal(t, "80", point.MinInclusive())
	require.Equal(t, "80", point.MaxExclusive())
	parsedPoint, err := ParseFeedRange(point.String())
	require.NoError(t, err)
	require.Equal(t, point, parsedPoint)
	for _, tt := range []struct {
		left, right     FeedRange
		subset, overlap bool
	}{
		{left, right, false, false}, {right, left, false, false},
		{left, full, true, true}, {full, left, false, true},
		{parsed, left, true, true},
	} {
		subset, err := tt.left.IsSubsetOf(tt.right)
		require.NoError(t, err)
		require.Equal(t, tt.subset, subset)
		overlap, err := tt.left.Overlaps(tt.right)
		require.NoError(t, err)
		require.Equal(t, tt.overlap, overlap)
	}
	for _, bounds := range [][2]string{{"FF", "80"}, {"not-hex", "FF"}} {
		value, err := NewFeedRange(bounds[0], bounds[1])
		require.Zero(t, value)
		require.Error(t, err)
	}
	for _, serialized := range []string{"not base64", "e30=", left.String() + "trailing"} {
		value, err := ParseFeedRange(serialized)
		require.Zero(t, value)
		require.Error(t, err)
	}
}

func TestNativeFeedRangeDiscoveryAndPartitionKeyOwnership(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	container := service.container(t, "hierarchical")
	discovered, err := container.GetFeedRanges(t.Context(), nil)
	require.NoError(t, err)
	require.Len(t, discovered.FeedRanges, 2)
	require.Equal(t, "", discovered.FeedRanges[0].MinInclusive())
	require.Equal(t, "80", discovered.FeedRanges[0].MaxExclusive())
	require.Equal(t, "80", discovered.FeedRanges[1].MinInclusive())
	require.Equal(t, "FF", discovered.FeedRanges[1].MaxExclusive())
	refreshed, err := container.GetFeedRanges(t.Context(), &GetFeedRangesOptions{ForceRefresh: true})
	require.NoError(t, err)
	require.Equal(t, discovered.FeedRanges, refreshed.FeedRanges)
	reads := service.routingReads.Load()
	_, err = container.GetFeedRanges(t.Context(), nil)
	require.NoError(t, err)
	require.Equal(t, reads, service.routingReads.Load(), "cached discovery must not force a routing-map request")
	_, err = container.GetFeedRanges(t.Context(), &GetFeedRangesOptions{ForceRefresh: true})
	require.NoError(t, err)
	require.Greater(t, service.routingReads.Load(), reads, "force refresh must reach native routing-map discovery")

	key := NewPartitionKeyString("synthetic").AppendString("category").AppendString("item")
	logical, err := container.GetFeedRangeForPartitionKey(t.Context(), key)
	require.NoError(t, err)
	require.False(t, logical.FeedRange.IsLogicalPartition(), "public values normalize logical keys to persistable EPK bounds")
	require.Equal(t, 3.5, logical.RequestCharge)
	parsed, err := ParseFeedRange(logical.FeedRange.String())
	require.NoError(t, err)
	require.False(t, parsed.IsLogicalPartition(), "serialization intentionally preserves only bounds")
	require.Equal(t, logical.FeedRange.MinInclusive(), parsed.MinInclusive())
	require.Equal(t, logical.FeedRange.MaxExclusive(), parsed.MaxExclusive())
	prefix, err := container.GetFeedRangeForPartitionKey(t.Context(), NewPartitionKeyString("synthetic"))
	require.NoError(t, err)
	require.False(t, prefix.FeedRange.IsLogicalPartition())
	require.NoError(t, container.database.client.Close())
	for _, value := range append(discovered.FeedRanges, logical.FeedRange, prefix.FeedRange) {
		parsed, err := ParseFeedRange(value.String())
		require.NoError(t, err, "values must outlive the discovery client")
		require.Equal(t, value.MinInclusive(), parsed.MinInclusive())
	}
}

func TestNativeFeedRangeListCopyEmpty(t *testing.T) {
	discovered, err := copyFeedRangeEntries(0, nil)
	require.NoError(t, err)
	require.NotNil(t, discovered)
	require.Empty(t, discovered)
}

func TestNativeFeedRangeDiscoveryCancellationDrainsLateList(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	container := service.container(t, "items")
	entered, proceed := make(chan struct{}), make(chan struct{})
	var enteredOnce, proceedOnce sync.Once
	unblock := func() { proceedOnce.Do(func() { close(proceed) }) }
	defer unblock()
	service.setRoutingHook(func(http.ResponseWriter, *http.Request) bool {
		enteredOnce.Do(func() { close(entered) })
		<-proceed
		return false
	})
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	finished := make(chan error, 1)
	go func() {
		response, err := container.GetFeedRanges(ctx, &GetFeedRangesOptions{ForceRefresh: true})
		if response.Response != (Response{}) || response.FeedRanges != nil {
			finished <- fmt.Errorf("canceled discovery returned ranges")
			return
		}
		finished <- err
	}()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("discovery did not reach native routing request")
	}
	cancel()
	select {
	case err := <-finished:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("discovery did not abandon its Go wait")
	}
	waiting := make(chan struct{})
	container.database.client.driver.beforePendingWait = func() { close(waiting) }
	closed := make(chan error, 1)
	go func() { closed <- container.database.client.Close() }()
	select {
	case <-waiting:
	case <-time.After(5 * time.Second):
		t.Fatal("client did not reach pending discovery cleanup")
	}
	select {
	case err := <-closed:
		t.Fatalf("client freed native resources before discovery delivery: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	unblock()
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("late discovery list was not drained")
	}
}

func TestNativeFeedRangeDiscoveryHonorsClientTimeout(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	timeout := time.Second
	credential, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := NewClientWithKey(t.Context(), service.url, credential,
		&ClientOptions{Operation: OperationOptions{EndToEndTimeout: &timeout}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(t.Context(), "db", "items")
	require.NoError(t, err)

	entered, proceed := make(chan struct{}), make(chan struct{})
	var enteredOnce, proceedOnce sync.Once
	unblock := func() { proceedOnce.Do(func() { close(proceed) }) }
	defer unblock()
	service.setRoutingHook(func(http.ResponseWriter, *http.Request) bool {
		enteredOnce.Do(func() { close(entered) })
		<-proceed
		return false
	})
	finished := make(chan error, 1)
	go func() {
		_, err := container.GetFeedRanges(t.Context(), &GetFeedRangesOptions{ForceRefresh: true})
		finished <- err
	}()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("discovery did not reach native routing request")
	}
	select {
	case err := <-finished:
		require.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(5 * time.Second):
		t.Fatal("discovery did not honor the client timeout")
	}
}

func TestNativeFeedRangeDiscoveryFailedRefreshUsesCachedRanges(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	container := service.container(t, "items")
	cached, err := container.GetFeedRanges(t.Context(), nil)
	require.NoError(t, err)
	reads := service.routingReads.Load()
	service.setRoutingHook(func(w http.ResponseWriter, _ *http.Request) bool {
		w.Header().Set("x-ms-substatus", "1001")
		w.Header().Set("x-ms-retry-after-ms", "25")
		w.Header().Set("ETag", "routing-failure")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, `{"code":"BadRequest","message":"synthetic routing failure"}`)
		return true
	})
	response, err := container.GetFeedRanges(t.Context(), &GetFeedRangesOptions{ForceRefresh: true})
	require.NoError(t, err, "native cached-success limitation: https://github.com/Azure/azure-sdk-for-rust/issues/5637")
	require.Greater(t, service.routingReads.Load(), reads)
	require.Equal(t, cached.FeedRanges, response.FeedRanges)
}

func TestNativeFeedRangeColdTopologyFailure(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	service.setRoutingHook(func(w http.ResponseWriter, _ *http.Request) bool {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, `{"code":"BadRequest","message":"synthetic cold routing failure"}`)
		return true
	})
	credential, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := NewClientWithKey(t.Context(), service.url, credential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(t.Context(), "db", "items")
	require.Nil(t, container)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientError, cosmosErr.Code)
	require.Equal(t, http.StatusServiceUnavailable, cosmosErr.StatusCode)
	require.Equal(t, 20305, cosmosErr.SubStatus)
	require.Greater(t, service.routingReads.Load(), int64(0))
}

func TestNativeFeedRangeDiscoveryConcurrentClose(t *testing.T) {
	service := newChangeFeedContractService(t, 2, false)
	container := service.container(t, "items")
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := container.GetFeedRanges(t.Context(), &GetFeedRangesOptions{ForceRefresh: true})
			if err == nil {
				require.Len(t, result.FeedRanges, 2)
			} else {
				var cosmosErr *Error
				require.ErrorAs(t, err, &cosmosErr)
				require.Equal(t, CodeClientClosed, cosmosErr.Code)
				require.Zero(t, result)
			}
		}()
	}
	require.NoError(t, container.database.client.Close())
	wg.Wait()
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	result, err := container.GetFeedRanges(canceled, nil)
	require.Zero(t, result)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

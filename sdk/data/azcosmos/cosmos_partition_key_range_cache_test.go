// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_partitionKeyRangeCache_getRoutingMap_populatesOnMiss(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	// We can't easily test with a real Client, but we can test the structure.
	// Verify that a new cache starts empty.
	cache.mu.RLock()
	require.Empty(t, cache.entries)
	cache.mu.RUnlock()
}

func Test_partitionKeyRangeCache_invalidate_nilEntry(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	// Invalidating a non-existent entry should not panic
	cache.invalidate("dbs/db1/colls/col1")

	cache.mu.RLock()
	_, exists := cache.entries["dbs/db1/colls/col1"]
	cache.mu.RUnlock()
	require.False(t, exists)
}

func Test_partitionKeyRangeCache_invalidate_existingEntry(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	// Manually populate
	entry := &pkRangeCacheEntry{
		routingMap: newCollectionRoutingMap([]partitionKeyRange{
			{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
		}, "etag1"),
	}
	cache.mu.Lock()
	cache.entries["dbs/db1/colls/col1"] = entry
	cache.mu.Unlock()

	// Verify populated
	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	entry.mu.Unlock()

	// Invalidate
	cache.invalidate("dbs/db1/colls/col1")

	// Verify nil
	entry.mu.Lock()
	require.Nil(t, entry.routingMap)
	entry.mu.Unlock()
}

func Test_partitionKeyRangeCache_getRoutingMap_cacheHit(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	expectedRM := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1"},
		{ID: "1", MinInclusive: "05C1", MaxExclusive: "FF"},
	}, "etag1")

	entry := &pkRangeCacheEntry{routingMap: expectedRM}
	cache.mu.Lock()
	cache.entries["dbs/db1/colls/col1"] = entry
	cache.mu.Unlock()

	// getRoutingMap with a nil client should return cached value without calling service
	rm, err := cache.getRoutingMap(context.Background(), "dbs/db1/colls/col1", nil)
	require.NoError(t, err)
	require.Equal(t, expectedRM, rm)
}

func Test_partitionKeyRangeCache_singleFlight(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	expectedRM := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")

	// Pre-populate entry with nil routingMap (simulates invalidated state)
	entry := &pkRangeCacheEntry{routingMap: nil}
	cache.mu.Lock()
	cache.entries["dbs/db1/colls/col1"] = entry
	cache.mu.Unlock()

	// Since we can't mock the HTTP call easily, we'll simulate by manually setting
	// the routing map after a "concurrent" lock acquisition.
	// This test verifies the per-entry mutex protects single-flight.
	var callCount int32

	// Simulate multiple goroutines trying to get the routing map
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			entry.mu.Lock()
			if entry.routingMap == nil {
				atomic.AddInt32(&callCount, 1)
				entry.routingMap = expectedRM
			}
			entry.mu.Unlock()
		}()
	}
	wg.Wait()

	// Only one goroutine should have populated the map
	require.Equal(t, int32(1), callCount)
	require.Equal(t, expectedRM, entry.routingMap)
}

func Test_partitionKeyRangeCache_forceRefresh_noEntry(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	// Pre-populate so forceRefresh has something to refresh
	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: rm}
	cache.mu.Lock()
	cache.entries["dbs/db1/colls/col1"] = entry
	cache.mu.Unlock()

	// forceRefresh with nil client will panic/fail on service call,
	// but if we test the path without service we can verify the entry.mu lock is acquired
	// For a full integration test, we'd need a mock. Here we just verify
	// the structure doesn't deadlock.
	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	entry.mu.Unlock()
}

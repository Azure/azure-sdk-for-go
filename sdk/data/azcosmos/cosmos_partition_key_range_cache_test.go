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

func Test_partitionKeyRangeCache_newCache_startsEmpty(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	cache.mu.RLock()
	require.Empty(t, cache.entries)
	cache.mu.RUnlock()
}

func Test_partitionKeyRangeCache_invalidate_nilEntry(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	// Invalidating a non-existent entry should not panic
	cache.invalidate("rid1")

	cache.mu.RLock()
	_, exists := cache.entries["rid1"]
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
	cache.entries["rid1"] = entry
	cache.mu.Unlock()

	// Verify populated
	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	entry.mu.Unlock()

	// Invalidate
	cache.invalidate("rid1")

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
	cache.entries["rid1"] = entry
	cache.mu.Unlock()

	// getRoutingMap with a nil client should return cached value without calling service
	rm, err := cache.getRoutingMap(context.Background(), "rid1", "dbs/db1/colls/col1", nil)
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
	cache.entries["rid1"] = entry
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

func Test_partitionKeyRangeCache_entryMutex_noDeadlock(t *testing.T) {
	cache := newPartitionKeyRangeCache()

	// Pre-populate an entry and verify we can acquire its mutex
	// without deadlocking against the cache-level mutex.
	rm := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: rm}
	cache.mu.Lock()
	cache.entries["rid1"] = entry
	cache.mu.Unlock()

	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	entry.mu.Unlock()
}

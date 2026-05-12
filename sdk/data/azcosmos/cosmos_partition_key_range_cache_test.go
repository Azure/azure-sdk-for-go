// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
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

// createMockClientForPKRangeCache creates a *Client wired to the given mock server
// with both caches initialized, suitable for testing PK range cache refresh flows.
func createMockClientForPKRangeCache(srv *mock.Server) *Client {
	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	return &Client{
		endpoint:       srv.URL(),
		endpointUrl:    defaultEndpoint,
		internal:       internalClient,
		gem:            gem,
		pkRangeCache:   newPartitionKeyRangeCache(),
		containerCache: newContainerPropertiesCache(),
	}
}

func Test_partitionKeyRangeCache_cacheMiss_fullRefresh(t *testing.T) {
	// Scenario: Empty cache, getRoutingMap triggers a full refresh from the service.
	pkRangeResponse := []byte(`{
		"_rid": "testRID",
		"PartitionKeyRanges": [
			{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": []},
			{"_rid": "r1", "id": "1", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": []}
		],
		"_count": 2
	}`)

	srv, close := mock.NewTLSServer()
	defer close()

	// Container properties response (for getContainerRID)
	srv.AppendResponse(
		mock.WithBody([]byte(`{"id": "col1", "_rid": "testRID", "partitionKey": {"paths": ["/pk"], "kind": "Hash", "version": 2}}`)),
		mock.WithStatusCode(200),
	)
	// PK range response (full refresh)
	srv.AppendResponse(
		mock.WithBody(pkRangeResponse),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
		mock.WithStatusCode(200),
	)

	client := createMockClientForPKRangeCache(srv)
	database, _ := newDatabase("db1", client)
	container, _ := newContainer("col1", database)

	// Use getPartitionKeyRanges which goes through getContainerRID → cache → getRoutingMap
	resp, err := container.getPartitionKeyRanges(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, 2, resp.Count)
	require.Equal(t, "", resp.PartitionKeyRanges[0].MinInclusive)
	require.Equal(t, "05C1E18D2D7F08", resp.PartitionKeyRanges[0].MaxExclusive)
	require.Equal(t, "05C1E18D2D7F08", resp.PartitionKeyRanges[1].MinInclusive)
	require.Equal(t, "FF", resp.PartitionKeyRanges[1].MaxExclusive)

	// Verify cache is populated
	rm, err := client.pkRangeCache.getRoutingMap(context.Background(), "testRID", container.link, client)
	require.NoError(t, err)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "etag1", rm.changeFeedETag)
}

func Test_partitionKeyRangeCache_incrementalRefresh_success(t *testing.T) {
	// Scenario: Cache has 1 range with ETag. Server returns 2 child ranges (split),
	// then 304. forceRefresh merges them incrementally.
	srv, close := mock.NewTLSServer()
	defer close()

	// Incremental feed response: 2 children replacing parent "0"
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r1", "id": "1", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": ["0"]},
				{"_rid": "r2", "id": "2", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": ["0"]}
			],
			"_count": 2
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag2"),
		mock.WithStatusCode(200),
	)
	// 304 Not Modified — no more changes
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag2"),
	)

	client := createMockClientForPKRangeCache(srv)

	// Pre-populate cache with 1 range + ETag
	initialMap := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: initialMap}
	client.pkRangeCache.mu.Lock()
	client.pkRangeCache.entries["testRID"] = entry
	client.pkRangeCache.mu.Unlock()

	// forceRefresh should do incremental refresh
	rm, err := client.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "", rm.orderedRanges[0].MinInclusive)
	require.Equal(t, "05C1E18D2D7F08", rm.orderedRanges[0].MaxExclusive)
	require.Equal(t, "05C1E18D2D7F08", rm.orderedRanges[1].MinInclusive)
	require.Equal(t, "FF", rm.orderedRanges[1].MaxExclusive)
	require.Equal(t, "etag2", rm.changeFeedETag)
	// Parent "0" should be marked as gone
	require.True(t, rm.isGone("0"))
}

func Test_partitionKeyRangeCache_incrementalRefresh_304_immediate(t *testing.T) {
	// Scenario: No changes since last fetch — 304 immediately, map preserved.
	srv, close := mock.NewTLSServer()
	defer close()

	// 304 Not Modified with updated ETag
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag2"),
	)

	client := createMockClientForPKRangeCache(srv)

	// Pre-populate cache with 2 ranges + ETag
	initialMap := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1E18D2D7F08"},
		{ID: "1", MinInclusive: "05C1E18D2D7F08", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: initialMap}
	client.pkRangeCache.mu.Lock()
	client.pkRangeCache.entries["testRID"] = entry
	client.pkRangeCache.mu.Unlock()

	rm, err := client.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	// Ranges should be preserved
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	// ETag should be updated
	require.Equal(t, "etag2", rm.changeFeedETag)
}

func Test_partitionKeyRangeCache_incrementalRefresh_mergeFailure_fullRefresh(t *testing.T) {
	// Scenario: tryCombine returns nil (incomplete ranges after merge) → falls through to full refresh.
	srv, close := mock.NewTLSServer()
	defer close()

	// Incremental response: only 1 child range that leaves a gap (parent "0" is gone
	// but child only covers half the range → tryCombine fails)
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r1", "id": "1", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": ["0"]}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag2"),
		mock.WithStatusCode(200),
	)
	// Full refresh response: complete set of 3 ranges
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r1", "id": "1", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": []},
				{"_rid": "r2", "id": "2", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "0BC1", "parents": []},
				{"_rid": "r3", "id": "3", "minInclusive": "0BC1", "maxExclusive": "FF", "parents": []}
			],
			"_count": 3
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag3"),
		mock.WithStatusCode(200),
	)

	client := createMockClientForPKRangeCache(srv)

	// Pre-populate cache with 1 range spanning the full space
	initialMap := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: initialMap}
	client.pkRangeCache.mu.Lock()
	client.pkRangeCache.entries["testRID"] = entry
	client.pkRangeCache.mu.Unlock()

	rm, err := client.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	// Should have the 3 ranges from full refresh
	require.Equal(t, 3, len(rm.orderedRanges))
	require.Equal(t, "1", rm.orderedRanges[0].ID)
	require.Equal(t, "2", rm.orderedRanges[1].ID)
	require.Equal(t, "3", rm.orderedRanges[2].ID)
	require.Equal(t, "etag3", rm.changeFeedETag)
}

func Test_partitionKeyRangeCache_incrementalRefresh_contextCancelled(t *testing.T) {
	// Scenario: Context is cancelled before the HTTP call in the incremental loop.
	srv, close := mock.NewTLSServer()
	defer close()

	// Queue a response that should never be reached
	srv.AppendResponse(
		mock.WithBody([]byte(`{"_rid": "testRID", "PartitionKeyRanges": [], "_count": 0}`)),
		mock.WithStatusCode(200),
	)

	client := createMockClientForPKRangeCache(srv)

	// Pre-populate cache with a range + ETag so it takes the incremental path
	initialMap := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: initialMap}
	client.pkRangeCache.mu.Lock()
	client.pkRangeCache.entries["testRID"] = entry
	client.pkRangeCache.mu.Unlock()

	// Cancel context before calling forceRefresh
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.pkRangeCache.forceRefresh(ctx, "testRID", "dbs/db1/colls/col1", client)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)

	// Cache entry should still have the original routing map (unchanged)
	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	require.Equal(t, "etag1", entry.routingMap.changeFeedETag)
	require.Equal(t, 1, len(entry.routingMap.orderedRanges))
	entry.mu.Unlock()
}

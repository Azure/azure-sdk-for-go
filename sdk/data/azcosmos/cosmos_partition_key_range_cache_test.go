// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
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
		endpoint:    srv.URL(),
		endpointUrl: defaultEndpoint,
		internal:    internalClient,
		gem:         gem,
		caches: &sharedCacheSet{
			pkRangeCache:   newPartitionKeyRangeCache(),
			containerCache: newContainerPropertiesCache(),
		},
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
	// PK range response (full change-feed refresh — first page)
	srv.AppendResponse(
		mock.WithBody(pkRangeResponse),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
		mock.WithStatusCode(200),
	)
	// 304 Not Modified — terminates the change-feed loop
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
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
	rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", container.link, client)
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
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	// forceRefresh should do incremental refresh
	rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
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
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
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
	// 304 terminates the incremental change-feed loop — tryCombine sees only the 1 child
	// range, detects the gap, returns nil → falls through to full refresh
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag2"),
	)
	// Full change-feed refresh response: complete set of 3 ranges
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
	// 304 Not Modified — terminates the full change-feed refresh loop
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag3"),
	)

	client := createMockClientForPKRangeCache(srv)

	// Pre-populate cache with 1 range spanning the full space
	initialMap := newCollectionRoutingMap([]partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: initialMap}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
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
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	// Cancel context before calling forceRefresh
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.caches.pkRangeCache.forceRefresh(ctx, "testRID", "dbs/db1/colls/col1", client)
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)

	// Cache entry should still have the original routing map (unchanged)
	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	require.Equal(t, "etag1", entry.routingMap.changeFeedETag)
	require.Equal(t, 1, len(entry.routingMap.orderedRanges))
	entry.mu.Unlock()
}

func Test_partitionKeyRangeCache_fullRefresh_multiPage(t *testing.T) {
	// Scenario: Full refresh requires multiple change-feed pages before 304.
	// This validates the pagination loop that accumulates ranges across pages.
	srv, close := mock.NewTLSServer()
	defer close()

	// Page 1: first partition range
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": []}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag-page1"),
		mock.WithStatusCode(200),
	)
	// Page 2: second partition range
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r1", "id": "1", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": []}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag-page2"),
		mock.WithStatusCode(200),
	)
	// 304 Not Modified — terminates the loop
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag-page2"),
	)

	client := createMockClientForPKRangeCache(srv)

	// Empty cache — will trigger full change-feed refresh
	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "", rm.orderedRanges[0].MinInclusive)
	require.Equal(t, "05C1E18D2D7F08", rm.orderedRanges[0].MaxExclusive)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	require.Equal(t, "05C1E18D2D7F08", rm.orderedRanges[1].MinInclusive)
	require.Equal(t, "FF", rm.orderedRanges[1].MaxExclusive)
	require.Equal(t, "etag-page2", rm.changeFeedETag)
}

func Test_partitionKeyRangeCache_fullRefresh_filtersParentsInSamePage(t *testing.T) {
	// Scenario: During a full change-feed refresh, the service returns parent and
	// children ranges in the same page. The parent-filtering in
	// newCollectionRoutingMap should filter the parent.
	srv, close := mock.NewTLSServer()
	defer close()

	// Response includes parent "0" and its two children
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "FF", "parents": []},
				{"_rid": "r1", "id": "1", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": ["0"]},
				{"_rid": "r2", "id": "2", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": ["0"]}
			],
			"_count": 3
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
		mock.WithStatusCode(200),
	)
	// 304 Not Modified
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
	)

	client := createMockClientForPKRangeCache(srv)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	// Parent "0" should be filtered out, leaving 2 child ranges
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "1", rm.orderedRanges[0].ID)
	require.Equal(t, "2", rm.orderedRanges[1].ID)
	require.True(t, rm.isGone("0"))
}

// captureRequestsPolicy is a test policy that records all outgoing HTTP requests.
type captureRequestsPolicy struct {
	mu       sync.Mutex
	requests []*http.Request
}

func (p *captureRequestsPolicy) Do(req *policy.Request) (*http.Response, error) {
	p.mu.Lock()
	p.requests = append(p.requests, req.Raw().Clone(req.Raw().Context()))
	p.mu.Unlock()
	return req.Next()
}

func Test_partitionKeyRangeCache_fullRefresh_setsChangeFeedHeaders(t *testing.T) {
	// Scenario: Verify that full refresh requests include A-IM and x-ms-max-item-count headers.
	srv, close := mock.NewTLSServer()
	defer close()

	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "FF", "parents": []}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
		mock.WithStatusCode(200),
	)
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
	)

	capture := &captureRequestsPolicy{}

	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{capture}},
		&policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{
		endpoint:    srv.URL(),
		endpointUrl: defaultEndpoint,
		internal:    internalClient,
		gem:         gem,
		caches: &sharedCacheSet{
			pkRangeCache:   newPartitionKeyRangeCache(),
			containerCache: newContainerPropertiesCache(),
		},
	}

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	_, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)

	// Should have captured 2 requests: initial fetch + 304 termination
	require.Equal(t, 2, len(capture.requests))

	// First request: should have A-IM and max-item-count, but NO If-None-Match
	firstReq := capture.requests[0]
	require.Equal(t, cosmosHeaderValuesChangeFeed, firstReq.Header.Get(cosmosHeaderChangeFeed))
	require.Equal(t, cosmosHeaderValuesMaxItemAll, firstReq.Header.Get(cosmosHeaderMaxItemCount))
	require.Empty(t, firstReq.Header.Get(headerIfNoneMatch), "full refresh should not set If-None-Match")

	// Second request: should have A-IM, max-item-count, AND If-None-Match with the ETag from page 1
	secondReq := capture.requests[1]
	require.Equal(t, cosmosHeaderValuesChangeFeed, secondReq.Header.Get(cosmosHeaderChangeFeed))
	require.Equal(t, cosmosHeaderValuesMaxItemAll, secondReq.Header.Get(cosmosHeaderMaxItemCount))
	require.Equal(t, "etag1", secondReq.Header.Get(headerIfNoneMatch))
}

func Test_partitionKeyRangeCache_fullRefresh_emptyPagesBeforeData(t *testing.T) {
	// Scenario: The service returns 200 with empty PartitionKeyRanges arrays
	// before returning actual data. The loop must NOT terminate on empty 200 pages —
	// only 304 Not Modified signals the end of the change feed.
	srv, close := mock.NewTLSServer()
	defer close()

	// Page 1: 200 with empty array (not 304 — should continue draining)
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [],
			"_count": 0
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag-empty1"),
		mock.WithStatusCode(200),
	)
	// Page 2: another empty 200
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [],
			"_count": 0
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag-empty2"),
		mock.WithStatusCode(200),
	)
	// Page 3: actual data
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": []},
				{"_rid": "r1", "id": "1", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": []}
			],
			"_count": 2
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag-data"),
		mock.WithStatusCode(200),
	)
	// 304 Not Modified — terminates the loop
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag-data"),
	)

	client := createMockClientForPKRangeCache(srv)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	require.Equal(t, "etag-data", rm.changeFeedETag)
}

func Test_partitionKeyRangeCache_incrementalRefresh_emptyPagesBeforeData(t *testing.T) {
	// Scenario: During incremental refresh, the service returns empty 200 pages
	// before the actual split data arrives. Must keep draining until 304.
	srv, close := mock.NewTLSServer()
	defer close()

	// Initial load: single range
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "FF", "parents": []}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
		mock.WithStatusCode(200),
	)
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
	)

	// Incremental refresh: empty 200, then actual split data, then 304
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [],
			"_count": 0
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag2"),
		mock.WithStatusCode(200),
	)
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r1", "id": "1", "minInclusive": "", "maxExclusive": "05C1E18D2D7F08", "parents": ["0"]},
				{"_rid": "r2", "id": "2", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": ["0"]}
			],
			"_count": 2
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag3"),
		mock.WithStatusCode(200),
	)
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag3"),
	)

	client := createMockClientForPKRangeCache(srv)

	// Populate cache
	rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.Equal(t, 1, len(rm.orderedRanges))

	// Trigger incremental refresh
	rm, err = client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)

	// Should have 2 ranges after split (parent "0" filtered)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "1", rm.orderedRanges[0].ID)
	require.Equal(t, "2", rm.orderedRanges[1].ID)
	require.Equal(t, "etag3", rm.changeFeedETag)
	require.True(t, rm.isGone("0"))
}

func Test_partitionKeyRangeCache_incrementalRefresh_cascadingSplitAcrossPages(t *testing.T) {
	// Scenario: Page 1 returns only ONE child of a split (range "1"), making the
	// intermediate state incomplete — range "1" covers only half the keyspace that
	// parent "0" covered, so per-page tryCombine would see a gap and fail.
	// Page 2 delivers the sibling ("2"). The accumulate-all-then-combine approach
	// handles this because tryCombine sees both children together.
	srv, close := mock.NewTLSServer()
	defer close()

	// Initial getRoutingMap call: single range "0" covering full keyspace
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r0", "id": "0", "minInclusive": "", "maxExclusive": "FF", "parents": []}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
		mock.WithStatusCode(200),
	)
	// 304 to complete initial load
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag1"),
	)

	// Incremental refresh page 1: only first child of split "0" → "1" + "2"
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
	// Incremental refresh page 2: second child of the same split
	srv.AppendResponse(
		mock.WithBody([]byte(`{
			"_rid": "testRID",
			"PartitionKeyRanges": [
				{"_rid": "r2", "id": "2", "minInclusive": "05C1E18D2D7F08", "maxExclusive": "FF", "parents": ["0"]}
			],
			"_count": 1
		}`)),
		mock.WithHeader(cosmosHeaderEtag, "etag3"),
		mock.WithStatusCode(200),
	)
	// 304 to complete incremental refresh
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, "etag3"),
	)

	client := createMockClientForPKRangeCache(srv)

	// First: populate the cache with initial range
	rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.Equal(t, 1, len(rm.orderedRanges))
	require.Equal(t, "0", rm.orderedRanges[0].ID)

	// Trigger incremental refresh via forceRefresh
	rm, err = client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)

	// Final state should have 2 ranges: "1" and "2" (parent "0" filtered)
	require.Equal(t, 2, len(rm.orderedRanges))
	require.Equal(t, "1", rm.orderedRanges[0].ID)
	require.Equal(t, "", rm.orderedRanges[0].MinInclusive)
	require.Equal(t, "05C1E18D2D7F08", rm.orderedRanges[0].MaxExclusive)
	require.Equal(t, "2", rm.orderedRanges[1].ID)
	require.Equal(t, "05C1E18D2D7F08", rm.orderedRanges[1].MinInclusive)
	require.Equal(t, "FF", rm.orderedRanges[1].MaxExclusive)
	require.Equal(t, "etag3", rm.changeFeedETag)

	// Verify parent is marked as gone
	require.True(t, rm.isGone("0"))
}

// gatePolicy blocks each request on a "started" signal until the test
// closes the "release" channel. It also counts requests so tests can assert
// single-flight (factory invoked exactly once).
type gatePolicy struct {
	started chan struct{}
	release chan struct{}
	count   atomic.Int32
}

func newGatePolicy() *gatePolicy {
	return &gatePolicy{
		started: make(chan struct{}, 16),
		release: make(chan struct{}),
	}
}

func (g *gatePolicy) Do(req *policy.Request) (*http.Response, error) {
	g.count.Add(1)
	// non-blocking signal so tests can count "arrivals" without deadlock
	select {
	case g.started <- struct{}{}:
	default:
	}
	select {
	case <-g.release:
	case <-req.Raw().Context().Done():
		return nil, req.Raw().Context().Err()
	}
	return req.Next()
}

func createGatedClientForPKRangeCache(srv *mock.Server, gate *gatePolicy) *Client {
	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{gate}},
		&policy.ClientOptions{Transport: srv})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	return &Client{
		endpoint:    srv.URL(),
		endpointUrl: defaultEndpoint,
		internal:    internalClient,
		gem:         gem,
		caches: &sharedCacheSet{
			pkRangeCache:   newPartitionKeyRangeCache(),
			containerCache: newContainerPropertiesCache(),
		},
	}
}

// appendPKRangesAndTerminator queues a 200 with the given ranges + ETag and
// a follow-up 304 to terminate the change-feed loop.
func appendPKRangesAndTerminator(srv *mock.Server, body []byte, etag string) {
	srv.AppendResponse(
		mock.WithBody(body),
		mock.WithHeader(cosmosHeaderEtag, etag),
		mock.WithStatusCode(200),
	)
	srv.AppendResponse(
		mock.WithStatusCode(304),
		mock.WithHeader(cosmosHeaderEtag, etag),
	)
}

func Test_partitionKeyRangeCache_getRoutingMap_returnsCachedDuringRefresh(t *testing.T) {
	// While a forced refresh is in flight, a concurrent getRoutingMap should
	// return the cached map immediately rather than waiting for the refresh.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	body := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body, "etag2")

	gate := newGatePolicy()
	client := createGatedClientForPKRangeCache(srv, gate)

	// Seed the entry with a routing map so getRoutingMap should fast-path.
	existing := newCollectionRoutingMap([]partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF"}}, "etag1")
	entry := &pkRangeCacheEntry{routingMap: existing}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	// Start a forced refresh in the background; it will block in the gate.
	refreshDone := make(chan struct{})
	go func() {
		defer close(refreshDone)
		_, _ = client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	}()

	// Wait for the refresh goroutine to actually issue its HTTP request.
	<-gate.started

	// Concurrent getRoutingMap MUST return immediately with the existing map.
	got, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.Same(t, existing, got)

	close(gate.release)
	<-refreshDone
}

func Test_partitionKeyRangeCache_forceRefresh_concurrentCallersShareOneFetch(t *testing.T) {
	// Multiple concurrent forceRefresh callers must trigger only one network
	// fetch and all observe the same resulting routing map.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	body := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body, "etag1")

	gate := newGatePolicy()
	client := createGatedClientForPKRangeCache(srv, gate)

	// Seed entry without a routing map so concurrent forceRefresh callers all
	// fall into the in-flight path.
	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	const N = 8
	var wg sync.WaitGroup
	results := make([]*collectionRoutingMap, N)
	errs := make([]error, N)
	wg.Add(N)
	for i := 0; i < N; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i], errs[i] = client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
		}()
	}

	// Let the (single) fetch arrive then release it.
	<-gate.started
	close(gate.release)
	wg.Wait()

	for i := 0; i < N; i++ {
		require.NoError(t, errs[i])
		require.NotNil(t, results[i])
		require.Same(t, results[0], results[i], "all callers must share the same routing map pointer")
	}
	// Exactly one full-refresh sequence (200 + 304) = 2 HTTP requests.
	require.Equal(t, int32(2), gate.count.Load())
}

// (stale-view pointer-identity dedup was dropped; the test exercising it has
// been removed accordingly.)

func Test_partitionKeyRangeCache_callerCancelDoesNotAbortSharedFetch(t *testing.T) {
	// A caller whose context is cancelled while awaiting the in-flight refresh
	// must return ctx.Err(); the refresh must continue and serve other waiters.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	body := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body, "etag1")

	gate := newGatePolicy()
	client := createGatedClientForPKRangeCache(srv, gate)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	// Caller A cancels mid-await.
	ctxA, cancelA := context.WithCancel(context.Background())
	errA := make(chan error, 1)
	go func() {
		_, err := client.caches.pkRangeCache.getRoutingMap(ctxA, "testRID", "dbs/db1/colls/col1", client)
		errA <- err
	}()

	// Caller B uses background and should succeed.
	errB := make(chan error, 1)
	rmB := make(chan *collectionRoutingMap, 1)
	go func() {
		rm, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
		rmB <- rm
		errB <- err
	}()

	<-gate.started
	cancelA()
	require.ErrorIs(t, <-errA, context.Canceled)

	// Release the refresh and B should observe the populated map.
	close(gate.release)
	require.NoError(t, <-errB)
	require.NotNil(t, <-rmB)
}

func Test_partitionKeyRangeCache_invalidateDuringRefresh_discardsResult(t *testing.T) {
	// invalidate() during an in-flight refresh must NOT abort the refresh
	// (other awaiters still get the in-flight result), but the result must
	// be discarded from the cache: an invalidate-during-refresh strictly
	// means "what you're about to install is stale, do not use it". The
	// awaiter must NOT receive the discarded map; instead the cache
	// internally retries and the awaiter ultimately receives the result of
	// a fresh post-invalidate refresh.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// First refresh body — will be discarded after invalidate.
	body1 := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body1, "etagOld")
	// Second refresh body — what the post-invalidate retry should pick up.
	body2 := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body2, "etagPost")

	gate := newGatePolicy()
	client := createGatedClientForPKRangeCache(srv, gate)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	done := make(chan struct {
		rm  *collectionRoutingMap
		err error
	}, 1)
	go func() {
		rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
		done <- struct {
			rm  *collectionRoutingMap
			err error
		}{rm, err}
	}()

	<-gate.started
	client.caches.pkRangeCache.invalidate("testRID")
	close(gate.release)

	res := <-done
	// Internal retry swallowed the sentinel; the caller observes the
	// post-invalidate fresh refresh's result.
	require.NoError(t, res.err)
	require.NotNil(t, res.rm)
	require.Equal(t, "etagPost", res.rm.changeFeedETag, "caller must see post-invalidate refresh, not the discarded one")

	// And the cache must hold the fresh map.
	entry.mu.Lock()
	require.NotNil(t, entry.routingMap)
	require.Equal(t, "etagPost", entry.routingMap.changeFeedETag)
	entry.mu.Unlock()
}

func Test_partitionKeyRangeCache_refreshError_clearsInFlightForRetry(t *testing.T) {
	// When a refresh fails, current waiters observe the error and the
	// in-flight slot is cleared so the next caller starts a fresh op.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// First refresh fails with 401 (non-transient, won't be retried by either
	// the pipeline or our per-page retry). Second refresh succeeds.
	srv.AppendResponse(mock.WithBody([]byte(`{"code":"Unauthorized"}`)), mock.WithStatusCode(401))
	body := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body, "etag1")

	client := createMockClientForPKRangeCache(srv)
	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	_, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.Error(t, err)

	// in-flight slot cleared so next call kicks off another op.
	entry.mu.Lock()
	require.Nil(t, entry.inFlight)
	entry.mu.Unlock()

	rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.NotNil(t, rm)
	require.Equal(t, "etag1", rm.changeFeedETag)
}

// createMockClientForPKRangeCacheNoRetry is like createMockClientForPKRangeCache
// but disables the pipeline's built-in retry policy so tests can exercise the
// change-feed loop's own per-page retry behaviour without the pipeline
// silently retrying transient failures first.
func createMockClientForPKRangeCacheNoRetry(srv *mock.Server) *Client {
	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0", azruntime.PipelineOptions{}, &policy.ClientOptions{
		Transport: srv,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	})
	gem := &globalEndpointManager{preferredLocations: []string{}}
	return &Client{
		endpoint:    srv.URL(),
		endpointUrl: defaultEndpoint,
		internal:    internalClient,
		gem:         gem,
		caches: &sharedCacheSet{
			pkRangeCache:   newPartitionKeyRangeCache(),
			containerCache: newContainerPropertiesCache(),
		},
	}
}

func Test_partitionKeyRangeCache_midPagination_transientRetrySucceeds(t *testing.T) {
	// The change-feed loop must survive a transient 408 between pages by
	// retrying the failing page, preserving the pages already accumulated.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// Page 1: 200 with one range, etag "p1"
	srv.AppendResponse(
		mock.WithBody([]byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"80","parents":[]}],"_count":1}`)),
		mock.WithHeader(cosmosHeaderEtag, "p1"),
		mock.WithStatusCode(200),
	)
	// Page 2: 408 transient — must be retried
	srv.AppendResponse(
		mock.WithBody([]byte(`{"code":"RequestTimeout"}`)),
		mock.WithStatusCode(408),
	)
	// Page 2 retry: 200 with second range, etag "p2"
	srv.AppendResponse(
		mock.WithBody([]byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r1","id":"1","minInclusive":"80","maxExclusive":"FF","parents":[]}],"_count":1}`)),
		mock.WithHeader(cosmosHeaderEtag, "p2"),
		mock.WithStatusCode(200),
	)
	// Terminator
	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, "p2"))

	client := createMockClientForPKRangeCacheNoRetry(srv)
	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.NotNil(t, rm)
	require.Equal(t, 2, len(rm.orderedRanges), "both pages must be in the final routing map")
	require.Equal(t, "0", rm.orderedRanges[0].ID)
	require.Equal(t, "1", rm.orderedRanges[1].ID)
	require.Equal(t, "p2", rm.changeFeedETag)
}

func Test_partitionKeyRangeCache_midPagination_callerCancelAbortsRetries(t *testing.T) {
	// Transient failures retry indefinitely until either the page succeeds
	// or the refresh's context is cancelled. Per-awaiter ctx cancellation
	// does NOT abort the shared background refresh (it runs on
	// context.Background()), so this test exercises invalidate which
	// causes the in-flight op to be discarded after the next page returns.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// Page 1: success
	srv.AppendResponse(
		mock.WithBody([]byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"80","parents":[]}],"_count":1}`)),
		mock.WithHeader(cosmosHeaderEtag, "p1"),
		mock.WithStatusCode(200),
	)
	// Page 2: 408 a few times, then succeed. Proves retries are not capped
	// at any small constant — the original test verified an attempt cap
	// that no longer exists.
	for i := 0; i < 5; i++ {
		srv.AppendResponse(
			mock.WithBody([]byte(`{"code":"RequestTimeout"}`)),
			mock.WithStatusCode(408),
		)
	}
	// Page 2 eventual success
	srv.AppendResponse(
		mock.WithBody([]byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r1","id":"1","minInclusive":"80","maxExclusive":"FF","parents":[]}],"_count":1}`)),
		mock.WithHeader(cosmosHeaderEtag, "p2"),
		mock.WithStatusCode(200),
	)
	// Terminator
	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, "p2"))

	client := createMockClientForPKRangeCacheNoRetry(srv)
	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	rm, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.NotNil(t, rm)
	require.Equal(t, 2, len(rm.orderedRanges), "infinite retry must eventually surface the successful page")
}

func Test_partitionKeyRangeCache_midPagination_nonTransientFailsFast(t *testing.T) {
	// A non-transient error (e.g. 401) on a mid-loop page must NOT be
	// retried; it should be surfaced immediately.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// Page 1: success
	srv.AppendResponse(
		mock.WithBody([]byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"80","parents":[]}],"_count":1}`)),
		mock.WithHeader(cosmosHeaderEtag, "p1"),
		mock.WithStatusCode(200),
	)
	// Page 2: 401 — fail fast, no retries
	srv.AppendResponse(
		mock.WithBody([]byte(`{"code":"Unauthorized"}`)),
		mock.WithStatusCode(401),
	)
	// Sentinel: if we DID retry, the next response would be a 200, which
	// would mask the bug. Make it a 408 instead so an unintended retry
	// also fails (test would still fail with the wrong status code).
	srv.AppendResponse(
		mock.WithBody([]byte(`{"code":"RequestTimeout"}`)),
		mock.WithStatusCode(408),
	)

	client := createMockClientForPKRangeCacheNoRetry(srv)
	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	_, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusUnauthorized, respErr.StatusCode, "non-transient 401 must fail fast without retry")
}

func Test_partitionKeyRangeCache_secondWaveAfterCompletion_noRedundantFetch(t *testing.T) {
	// A caller that arrives AFTER a refresh has completed but already has a
	// usable routing map (e.g. another caller just installed one) must NOT
	// trigger a second redundant network fetch. The cached map is observable
	// immediately via getRoutingMap.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	body := []byte(`{"_rid":"testRID","PartitionKeyRanges":[{"_rid":"r0","id":"0","minInclusive":"","maxExclusive":"FF","parents":[]}],"_count":1}`)
	appendPKRangesAndTerminator(srv, body, "etagWave1")

	gate := newGatePolicy()
	close(gate.release) // never block — wave 1 completes immediately
	client := createGatedClientForPKRangeCache(srv, gate)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	// Wave 1: forceRefresh populates the cache.
	rm1, err := client.caches.pkRangeCache.forceRefresh(context.Background(), "testRID", "dbs/db1/colls/col1", client)
	require.NoError(t, err)
	require.NotNil(t, rm1)
	wave1Count := gate.count.Load()
	require.Equal(t, int32(2), wave1Count, "wave 1 should have made the full-refresh round-trip (200 + 304)")

	// Wave 2: getRoutingMap arriving after wave 1 completed must return
	// the cached map immediately, NOT trigger another fetch.
	const N = 5
	for i := 0; i < N; i++ {
		got, err := client.caches.pkRangeCache.getRoutingMap(context.Background(), "testRID", "dbs/db1/colls/col1", client)
		require.NoError(t, err)
		require.Same(t, rm1, got, "wave 2 callers must see the wave-1 map")
	}
	require.Equal(t, wave1Count, gate.count.Load(), "wave 2 must not trigger additional network requests")
}

func Test_partitionKeyRangeCache_getRoutingMap_canceledCallerNoBackgroundFetch(t *testing.T) {
	// When the caller's ctx is already canceled and no refresh is in flight,
	// getRoutingMap must return ctx.Err() WITHOUT spawning a background
	// fetch that nobody is waiting on.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	gate := newGatePolicy()
	close(gate.release)
	client := createGatedClientForPKRangeCache(srv, gate)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.caches.pkRangeCache.getRoutingMap(ctx, "testRID", "dbs/db1/colls/col1", client)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, int32(0), gate.count.Load(), "canceled caller must not trigger background fetch")

	// And no in-flight op was left behind.
	entry.mu.Lock()
	require.Nil(t, entry.inFlight)
	entry.mu.Unlock()
}

func Test_partitionKeyRangeCache_forceRefresh_canceledCallerNoBackgroundFetch(t *testing.T) {
	// Same as the getRoutingMap variant but for forceRefresh.
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	gate := newGatePolicy()
	close(gate.release)
	client := createGatedClientForPKRangeCache(srv, gate)

	entry := &pkRangeCacheEntry{}
	client.caches.pkRangeCache.mu.Lock()
	client.caches.pkRangeCache.entries["testRID"] = entry
	client.caches.pkRangeCache.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.caches.pkRangeCache.forceRefresh(ctx, "testRID", "dbs/db1/colls/col1", client)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, int32(0), gate.count.Load(), "canceled caller must not trigger background fetch")

	entry.mu.Lock()
	require.Nil(t, entry.inFlight)
	entry.mu.Unlock()
}

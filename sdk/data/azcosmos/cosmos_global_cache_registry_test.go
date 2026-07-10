// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"weak"

	"github.com/stretchr/testify/require"
)

// registryRefCount returns the current reference count recorded in the registry
// entry for the given endpoint, or -1 if no entry exists. Test helper.
func registryRefCount(endpoint string) int {
	val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	if !ok {
		return -1
	}
	entry := val.(*cacheRegistryEntry)
	entry.mu.Lock()
	defer entry.mu.Unlock()
	return entry.refCount
}

func TestNormalizeEndpoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://MyAccount.documents.azure.com/", "https://myaccount.documents.azure.com"},
		{"https://MyAccount.documents.azure.com", "https://myaccount.documents.azure.com"},
		{"https://MYACCOUNT.DOCUMENTS.AZURE.COM///", "https://myaccount.documents.azure.com"},
		{"https://myaccount.documents.azure.com:443/", "https://myaccount.documents.azure.com"},
		{"https://myaccount.documents.azure.com:443", "https://myaccount.documents.azure.com"},
		{"http://myaccount.documents.azure.com:80", "http://myaccount.documents.azure.com"},
		{"https://localhost:8081/", "https://localhost"},
		{"https://localhost:8081", "https://localhost"},
		{"  https://myaccount.documents.azure.com  ", "https://myaccount.documents.azure.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			require.Equal(t, tt.expected, normalizeEndpoint(tt.input))
		})
	}
}

func TestAcquireCaches_SameEndpoint_ReturnsSameInstance(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	set1 := acquireCaches("https://account1.documents.azure.com/")
	set2 := acquireCaches("https://account1.documents.azure.com/")
	set3 := acquireCaches("https://Account1.Documents.Azure.Com") // different case, no trailing slash

	require.Same(t, set1, set2, "same endpoint should return same cache set")
	require.Same(t, set1, set3, "normalization should make these equivalent")
	require.Equal(t, 3, registryRefCount("https://account1.documents.azure.com/"))
}

func TestAcquireCaches_DifferentEndpoints_ReturnDifferentInstances(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	set1 := acquireCaches("https://account1.documents.azure.com")
	set2 := acquireCaches("https://account2.documents.azure.com")

	require.NotSame(t, set1, set2, "different endpoints should return different cache sets")
	require.NotSame(t, set1.containerCache, set2.containerCache)
	require.NotSame(t, set1.pkRangeCache, set2.pkRangeCache)
}

func TestReleaseCaches_RemovesEntryWhenZeroRefs(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	set1 := acquireCaches(endpoint)
	_ = acquireCaches(endpoint) // refCount = 2

	releaseCaches(endpoint) // refCount = 1
	require.Equal(t, 1, registryRefCount(endpoint))

	// Entry should still be in the registry
	val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.True(t, ok)
	require.Same(t, set1, val.(*cacheRegistryEntry).weakRef.Value())

	releaseCaches(endpoint) // refCount = 0 → removed
	_, ok = globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.False(t, ok, "entry should be removed when refCount reaches 0")
}

func TestReleaseCaches_NewAcquireAfterFullRelease_CreatesNew(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	set1 := acquireCaches(endpoint)
	releaseCaches(endpoint) // refCount = 0, removed

	set2 := acquireCaches(endpoint) // new instance
	require.NotSame(t, set1, set2, "should create a fresh cache set after full release")
}

func TestAcquireCaches_CollectedWhenClientDiscardedWithoutClose(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"

	// Acquire a cache set the way a Client would, but never call releaseCaches
	// (i.e. never call Close) and drop the only strong reference — simulating a
	// discarded Client. The set escapes only through the weak registry entry.
	weakRef := func() weak.Pointer[sharedCacheSet] {
		set := acquireCaches(endpoint)
		require.NotNil(t, set)
		val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
		require.True(t, ok)
		w := val.(*cacheRegistryEntry).weakRef
		require.NotNil(t, w.Value())
		return w
	}()

	// Because the registry holds only a weak reference, the cache set becomes
	// eligible for collection once the discarded Client's strong reference is
	// gone — even though Close was never called.
	collected := false
	for i := 0; i < 100; i++ {
		runtime.GC()
		if weakRef.Value() == nil {
			collected = true
			break
		}
	}
	require.True(t, collected, "cache set should be reclaimed once the discarded client's strong ref is gone")

	// The small registry entry remains, but reacquiring must build a fresh set.
	set2 := acquireCaches(endpoint)
	require.NotNil(t, set2)
	val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.True(t, ok)
	require.Same(t, set2, val.(*cacheRegistryEntry).weakRef.Value())
	releaseCaches(endpoint)
}

func TestAcquireCaches_ReacquireAfterCollectionResetsRefCount(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"

	// Two acquires, no releases, then both strong refs discarded (leaked clients).
	weakRef := func() weak.Pointer[sharedCacheSet] {
		_ = acquireCaches(endpoint)
		_ = acquireCaches(endpoint) // refCount = 2; both refs dropped on return
		val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
		require.True(t, ok)
		return val.(*cacheRegistryEntry).weakRef
	}()
	require.Equal(t, 2, registryRefCount(endpoint))

	collected := false
	for i := 0; i < 100; i++ {
		runtime.GC()
		if weakRef.Value() == nil {
			collected = true
			break
		}
	}
	require.True(t, collected, "cache set should be reclaimed after both leaked refs are gone")

	// A fresh acquire must reset the stale refCount to 1, not resume from 2.
	set := acquireCaches(endpoint)
	require.NotNil(t, set)
	require.Equal(t, 1, registryRefCount(endpoint))
	releaseCaches(endpoint)
}

func TestSharedCaches_CrossClientCacheHit(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"

	// Simulate two clients targeting the same endpoint
	clientA := &Client{endpoint: endpoint, caches: acquireCaches(endpoint)}
	clientB := &Client{endpoint: endpoint, caches: acquireCaches(endpoint)}
	defer clientA.Close()
	defer clientB.Close()

	// Verify they share the exact same cache instances
	require.Same(t, clientA.getContainerCache(), clientB.getContainerCache())
	require.Same(t, clientA.getPKRangeCache(), clientB.getPKRangeCache())

	// Client A populates the container cache
	containerLink := "dbs/db1/colls/col1"
	props := &ContainerProperties{
		ID:         "col1",
		ResourceID: "rid-abc",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Version: 2,
		},
	}
	clientA.getContainerCache().set(containerLink, props)

	// Client B reads from cache — gets the same value with zero HTTP calls
	containerB := &ContainerClient{link: containerLink}
	result, err := clientB.getContainerCache().getProperties(context.Background(), containerB)
	require.NoError(t, err)
	require.Equal(t, props, result)
	require.Equal(t, "rid-abc", result.ResourceID)
}

func TestSharedCaches_CrossClientInvalidationVisibility(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"

	clientA := &Client{endpoint: endpoint, caches: acquireCaches(endpoint)}
	clientB := &Client{endpoint: endpoint, caches: acquireCaches(endpoint)}
	defer clientA.Close()
	defer clientB.Close()

	containerLink := "dbs/db1/colls/col1"
	props := &ContainerProperties{
		ID:         "col1",
		ResourceID: "rid-abc",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Version: 2,
		},
	}
	clientA.getContainerCache().set(containerLink, props)

	// Client B invalidates (e.g., got a 410)
	clientB.getContainerCache().invalidate(containerLink)

	// Client A should also see the invalidation since they share the cache
	containerA := &ContainerClient{link: containerLink}
	// getProperties will try to refresh — but no pipeline is wired, so
	// we check the entry directly to confirm the nil state is visible
	clientA.getContainerCache().mu.RLock()
	entry := clientA.getContainerCache().entries[containerLink]
	clientA.getContainerCache().mu.RUnlock()
	require.NotNil(t, entry)
	entry.mu.Lock()
	require.Nil(t, entry.props, "invalidation by Client B should be visible to Client A")
	entry.mu.Unlock()
	_ = containerA
}

func TestSharedCaches_ConcurrentClientsRefreshSingleFlight(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	caches := acquireCaches(endpoint)
	defer releaseCaches(endpoint)

	containerLink := "dbs/db1/colls/col1"
	props := &ContainerProperties{
		ID:         "col1",
		ResourceID: "rid-abc",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Version: 2,
		},
	}

	// Pre-populate cache, then invalidate to simulate a 410
	caches.containerCache.set(containerLink, props)
	caches.containerCache.invalidate(containerLink)

	// Now manually set the entry.props to simulate a refresh completing
	// while many goroutines are waiting on the entry lock.
	// This tests that all waiters get the value once ONE sets it.
	const numClients = 50
	var wg sync.WaitGroup
	var fetchCount atomic.Int64
	results := make([]*ContainerProperties, numClients)

	// Grab the entry and manually control the lock to simulate the race
	caches.containerCache.mu.RLock()
	entry := caches.containerCache.entries[containerLink]
	caches.containerCache.mu.RUnlock()
	require.NotNil(t, entry)

	// Lock the entry before spawning goroutines — simulates a refresh in progress
	entry.mu.Lock()

	wg.Add(numClients)
	for i := 0; i < numClients; i++ {
		go func(idx int) {
			defer wg.Done()
			// Each goroutine tries to acquire entry lock (simulating getProperties)
			entry.mu.Lock()
			if entry.props == nil {
				// "I'm the one refreshing" — simulate HTTP call
				fetchCount.Add(1)
				entry.props = &ContainerProperties{
					ID:         "col1",
					ResourceID: "rid-new",
				}
			}
			results[idx] = entry.props
			entry.mu.Unlock()
		}(i)
	}

	// Release the lock — one goroutine will "win" and set the value,
	// all others will see it already set.
	entry.mu.Unlock()
	wg.Wait()

	// Exactly one goroutine should have done the "fetch"
	require.Equal(t, int64(1), fetchCount.Load(),
		"only one goroutine should refresh; others should see the populated value")

	// All goroutines should have gotten the same result
	for i := 0; i < numClients; i++ {
		require.NotNil(t, results[i])
		require.Equal(t, "rid-new", results[i].ResourceID)
	}
}

func TestAcquireCaches_ConcurrentSafe(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	const goroutines = 100

	var wg sync.WaitGroup
	results := make([]*sharedCacheSet, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			results[idx] = acquireCaches(endpoint)
		}(i)
	}
	wg.Wait()

	// All goroutines should have gotten the same instance
	for i := 1; i < goroutines; i++ {
		require.Same(t, results[0], results[i], "concurrent acquires should return same instance")
	}
	require.Equal(t, goroutines, registryRefCount(endpoint))
}

func TestAcquireCaches_ResilientToInterleavedRelease(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"

	// Simulate the TOCTOU scenario: acquire, then concurrent release+acquire.
	// After the race, every acquire must still return the same registered instance.
	set1 := acquireCaches(endpoint) // refCount = 1

	// Release fully (refCount -> 0, entry removed)
	releaseCaches(endpoint)

	// Now a new acquire should create a fresh set (the old one is gone)
	set2 := acquireCaches(endpoint)
	require.NotSame(t, set1, set2, "after full release, a new set should be created")
	require.Equal(t, 1, registryRefCount(endpoint))

	// Stress test: interleave acquires and releases concurrently
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			s := acquireCaches(endpoint)
			// Verify the returned set is actually in the registry
			val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
			require.True(t, ok, "acquired set must be in registry")
			require.Same(t, s, val.(*cacheRegistryEntry).weakRef.Value(), "acquired set must match registry entry")
			releaseCaches(endpoint)
		}()
	}
	wg.Wait()

	// set2 should still be valid (we haven't released it)
	val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.True(t, ok)
	require.Same(t, set2, val.(*cacheRegistryEntry).weakRef.Value())
	require.Equal(t, 1, registryRefCount(endpoint))
	releaseCaches(endpoint)
}

func TestClientClose_Idempotent(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	client := &Client{endpoint: endpoint, caches: acquireCaches(endpoint)}

	// Also acquire a second reference so we can observe the refCount
	_ = acquireCaches(endpoint) // refCount = 2

	client.Close() // refCount -> 1
	client.Close() // idempotent — should NOT decrement again
	client.Close() // still idempotent

	// RefCount should be 1 (the second acquire), not -1 or 0
	_, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.True(t, ok, "entry should still exist — second reference is alive")
	require.Equal(t, 1, registryRefCount(endpoint))

	// Clean up the second reference
	releaseCaches(endpoint)
}

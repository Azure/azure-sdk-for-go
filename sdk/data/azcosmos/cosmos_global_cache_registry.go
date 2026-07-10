// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/url"
	"strings"
	"sync"
	"weak"
)

// sharedCacheSet groups the metadata caches that are shared across all Client
// instances targeting the same Cosmos DB account endpoint. This ensures that
// partition key range and container property metadata is fetched once per
// account regardless of how many Client instances exist.
type sharedCacheSet struct {
	containerCache *containerPropertiesCache
	pkRangeCache   *partitionKeyRangeCache
}

// cacheRegistryEntry is the small, stable value kept strongly in the registry
// map for a given endpoint. It holds only a weak reference to the potentially
// large sharedCacheSet, so the caches become eligible for garbage collection
// once the last Client holding a strong reference is discarded — even if the
// user never calls Close. The entry itself (a mutex, a weak pointer, and an
// int) is the only per-endpoint memory that leaks in that case.
//
// The mutex makes the refCount adjustment and the weak-pointer check atomic as
// a unit. Client creation is not a hot path, so a per-endpoint lock is an
// acceptable simplification over lock-free reference counting.
type cacheRegistryEntry struct {
	mu       sync.Mutex
	weakRef  weak.Pointer[sharedCacheSet]
	refCount int // guarded by mu
}

// globalCacheRegistry is a process-level registry of shared cache sets keyed
// by normalized account endpoint. It ensures singleton caches per account.
var globalCacheRegistry sync.Map // map[string]*cacheRegistryEntry

// normalizeEndpoint returns a canonical form of the endpoint for use as a
// registry key. It lowercases the host and strips paths. The port is included
// only when it is not the scheme's default (443 for https, 80 for http); a
// port that matches the scheme default is removed so that
// "https://account.documents.azure.com:443/" and
// "https://account.documents.azure.com" resolve to the same key.
//
// Non-default ports are deliberately preserved so that distinct targets such
// as multiple emulator instances or proxy/forwarding configurations are never
// collapsed onto a single cache entry. Confusing two account endpoints for one
// another would be a major failure, so we err on the side of keeping keys
// distinct.
func normalizeEndpoint(endpoint string) string {
	u, err := url.Parse(strings.TrimSpace(endpoint))
	if err != nil || u.Host == "" {
		// Fallback for malformed input
		return strings.TrimRight(strings.ToLower(endpoint), "/")
	}

	scheme := strings.ToLower(u.Scheme)
	port := u.Port()
	switch {
	case scheme == "https" && port == "443":
		port = ""
	case scheme == "http" && port == "80":
		port = ""
	}

	host := strings.ToLower(u.Hostname())
	if port == "" {
		return scheme + "://" + host
	}
	return scheme + "://" + host + ":" + port
}

// acquireCaches returns the shared cache set for the given endpoint, creating
// one if it doesn't exist. The returned pointer is a strong reference that the
// caller (the Client) must retain; the registry keeps only a weak reference, so
// the caches are reclaimed once every Client for the endpoint is discarded.
// The caller should call releaseCaches when the Client is closed to
// deterministically remove the (small) registry entry.
//
// A fresh cache set is created when the endpoint has no entry, when a previous
// set was garbage collected (because all its Clients were discarded without
// Close), or when a previous set was deterministically released. This covers
// callers that frequently create and discard Clients for the same endpoint.
func acquireCaches(endpoint string) *sharedCacheSet {
	key := normalizeEndpoint(endpoint)

	for {
		val, _ := globalCacheRegistry.LoadOrStore(key, &cacheRegistryEntry{})
		entry := val.(*cacheRegistryEntry)

		entry.mu.Lock()
		// A concurrent releaseCaches may have removed this exact entry from the
		// registry after our LoadOrStore. If so, retry so we don't attach to an
		// orphaned entry that other callers can no longer find.
		if cur, ok := globalCacheRegistry.Load(key); !ok || cur != val {
			entry.mu.Unlock()
			continue
		}

		if set := entry.weakRef.Value(); set != nil {
			entry.refCount++
			entry.mu.Unlock()
			return set
		}

		// The endpoint was never populated, or its cache set was collected or
		// released. Create a fresh set and reset the count.
		set := &sharedCacheSet{
			containerCache: newContainerPropertiesCache(),
			pkRangeCache:   newPartitionKeyRangeCache(),
		}
		entry.weakRef = weak.Make(set)
		entry.refCount = 1
		entry.mu.Unlock()
		return set
	}
}

// releaseCaches decrements the reference count for the given endpoint's cache
// set and removes the entry from the registry when no clients remain.
func releaseCaches(endpoint string) {
	key := normalizeEndpoint(endpoint)
	val, ok := globalCacheRegistry.Load(key)
	if !ok {
		return
	}
	entry := val.(*cacheRegistryEntry)
	entry.mu.Lock()
	entry.refCount--
	if entry.refCount <= 0 {
		entry.refCount = 0
		entry.weakRef = weak.Pointer[sharedCacheSet]{}
		globalCacheRegistry.CompareAndDelete(key, val)
	}
	entry.mu.Unlock()
}

// resetGlobalCacheRegistry clears the global cache registry.
// This is intended for test isolation only.
func resetGlobalCacheRegistry() {
	globalCacheRegistry.Range(func(key, _ any) bool {
		globalCacheRegistry.Delete(key)
		return true
	})
}

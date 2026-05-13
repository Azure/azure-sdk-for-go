// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strings"
	"sync"
	"sync/atomic"
)

// sharedCacheSet groups the metadata caches that are shared across all Client
// instances targeting the same Cosmos DB account endpoint. This ensures that
// partition key range and container property metadata is fetched once per
// account regardless of how many Client instances exist.
type sharedCacheSet struct {
	containerCache *containerPropertiesCache
	pkRangeCache   *partitionKeyRangeCache
	refCount       atomic.Int64
}

// globalCacheRegistry is a process-level registry of shared cache sets keyed
// by normalized account endpoint. It ensures singleton caches per account.
var globalCacheRegistry sync.Map // map[string]*sharedCacheSet

// normalizeEndpoint returns a canonical form of the endpoint for use as a
// registry key. It lowercases and strips the trailing slash.
func normalizeEndpoint(endpoint string) string {
	return strings.TrimRight(strings.ToLower(endpoint), "/")
}

// acquireCaches returns the shared cache set for the given endpoint, creating
// one if it doesn't exist. The caller must call releaseCaches when the Client
// is closed to allow cleanup.
//
// The implementation uses a verify-after-increment pattern to prevent a TOCTOU
// race with releaseCaches: after incrementing the refCount, it confirms the
// entry is still in the registry. If a concurrent release evicted it, the
// increment is undone and the loop retries.
func acquireCaches(endpoint string) *sharedCacheSet {
	key := normalizeEndpoint(endpoint)

	for {
		if val, ok := globalCacheRegistry.Load(key); ok {
			set := val.(*sharedCacheSet)
			set.refCount.Add(1)
			// Verify the entry is still registered after we incremented.
			// A concurrent releaseCaches could have evicted it between our
			// Load and Add, leaving us with an orphaned cache set.
			if val2, ok2 := globalCacheRegistry.Load(key); ok2 && val2 == val {
				return set
			}
			// Entry was evicted — undo our increment and retry.
			set.refCount.Add(-1)
			continue
		}

		// Slow path: create new and use LoadOrStore to avoid races
		newSet := &sharedCacheSet{
			containerCache: newContainerPropertiesCache(),
			pkRangeCache:   newPartitionKeyRangeCache(),
		}
		newSet.refCount.Store(1)

		_, loaded := globalCacheRegistry.LoadOrStore(key, newSet)
		if loaded {
			// Another goroutine created it first — use theirs via the retry loop
			// to get the verify-after-increment safety.
			continue
		}
		return newSet
	}
}

// releaseCaches decrements the reference count for the given endpoint's cache
// set and removes it from the registry when no clients remain.
func releaseCaches(endpoint string) {
	key := normalizeEndpoint(endpoint)
	if val, ok := globalCacheRegistry.Load(key); ok {
		set := val.(*sharedCacheSet)
		if set.refCount.Add(-1) <= 0 {
			globalCacheRegistry.CompareAndDelete(key, val)
		}
	}
}

// resetGlobalCacheRegistry clears the global cache registry.
// This is intended for test isolation only.
func resetGlobalCacheRegistry() {
	globalCacheRegistry.Range(func(key, _ any) bool {
		globalCacheRegistry.Delete(key)
		return true
	})
}

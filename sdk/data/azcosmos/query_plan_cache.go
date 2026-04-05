// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// QueryPlanCache caches query plans to reduce Gateway roundtrips.
// Similar to Java SDK's QueryPlanCache (5000 entries, v4.20.0+).
type QueryPlanCache struct {
	mu       sync.RWMutex
	entries  map[string]*queryPlanCacheEntry
	maxSize  int
	ttl      time.Duration
	evictSeq uint64 // Monotonic counter for LRU eviction
}

type queryPlanCacheEntry struct {
	plan       []byte
	createdAt  time.Time
	accessedAt time.Time
	accessSeq  uint64
}

// QueryPlanCacheOptions configures the query plan cache.
type QueryPlanCacheOptions struct {
	// MaxSize is the maximum number of query plans to cache.
	// Default is 5000 (matching Java SDK).
	MaxSize int

	// TTL is the time-to-live for cached query plans.
	// Default is 5 minutes.
	TTL time.Duration
}

// DefaultQueryPlanCacheOptions returns the default query plan cache configuration.
func DefaultQueryPlanCacheOptions() *QueryPlanCacheOptions {
	return &QueryPlanCacheOptions{
		MaxSize: 5000,
		TTL:     5 * time.Minute,
	}
}

// newQueryPlanCache creates a new query plan cache.
func newQueryPlanCache(opts *QueryPlanCacheOptions) *QueryPlanCache {
	if opts == nil {
		opts = DefaultQueryPlanCacheOptions()
	}

	maxSize := opts.MaxSize
	if maxSize <= 0 {
		maxSize = 5000
	}

	ttl := opts.TTL
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &QueryPlanCache{
		entries: make(map[string]*queryPlanCacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// cacheKey generates a cache key from the query text and container link.
// The key is based on the SHA256 hash of the query to handle long queries.
func (c *QueryPlanCache) cacheKey(containerLink, query string) string {
	hash := sha256.Sum256([]byte(containerLink + ":" + query))
	return hex.EncodeToString(hash[:])
}

// Get retrieves a cached query plan if available and not expired.
func (c *QueryPlanCache) Get(containerLink, query string) ([]byte, bool) {
	c.mu.RLock()
	key := c.cacheKey(containerLink, query)
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	// Check TTL
	if time.Since(entry.createdAt) > c.ttl {
		c.Remove(containerLink, query)
		return nil, false
	}

	// Update access time for LRU (requires write lock)
	c.mu.Lock()
	entry.accessedAt = time.Now()
	c.evictSeq++
	entry.accessSeq = c.evictSeq
	c.mu.Unlock()

	return entry.plan, true
}

// Set stores a query plan in the cache.
func (c *QueryPlanCache) Set(containerLink, query string, plan []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.cacheKey(containerLink, query)

	// Evict if at capacity (LRU eviction)
	if len(c.entries) >= c.maxSize {
		c.evictOldestLocked()
	}

	now := time.Now()
	c.evictSeq++
	c.entries[key] = &queryPlanCacheEntry{
		plan:       plan,
		createdAt:  now,
		accessedAt: now,
		accessSeq:  c.evictSeq,
	}
}

// Remove removes a query plan from the cache.
func (c *QueryPlanCache) Remove(containerLink, query string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.cacheKey(containerLink, query)
	delete(c.entries, key)
}

// Clear removes all entries from the cache.
func (c *QueryPlanCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*queryPlanCacheEntry)
}

// Size returns the current number of cached entries.
func (c *QueryPlanCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// evictOldestLocked removes the least recently used entry.
// Must be called with the write lock held.
func (c *QueryPlanCache) evictOldestLocked() {
	if len(c.entries) == 0 {
		return
	}

	var oldestKey string
	var oldestSeq uint64 = ^uint64(0) // Max uint64

	for key, entry := range c.entries {
		if entry.accessSeq < oldestSeq {
			oldestSeq = entry.accessSeq
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
	}
}

// evictExpiredLocked removes all expired entries.
// Must be called with the write lock held.
func (c *QueryPlanCache) evictExpiredLocked() {
	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > c.ttl {
			delete(c.entries, key)
		}
	}
}

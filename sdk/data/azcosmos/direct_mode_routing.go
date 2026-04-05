// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"strings"
	"sync"
	"time"
)

const (
	minEffectivePartitionKey = ""
	maxEffectivePartitionKey = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
)

type directModeRoutingCache struct {
	mu      sync.RWMutex
	entries map[string]*routingCacheEntry
	ttl     time.Duration
}

type routingCacheEntry struct {
	collectionRID string
	pkRanges      []partitionKeyRange
	pkDefinition  *PartitionKeyDefinition
	createdAt     time.Time
}

func newDirectModeRoutingCache(ttl time.Duration) *directModeRoutingCache {
	return &directModeRoutingCache{
		entries: make(map[string]*routingCacheEntry),
		ttl:     ttl,
	}
}

func routingCacheKey(databaseName, containerName string) string {
	return databaseName + "/" + containerName
}

func (c *directModeRoutingCache) get(databaseName, containerName string) *routingCacheEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := routingCacheKey(databaseName, containerName)
	entry, ok := c.entries[key]
	if !ok {
		return nil
	}

	if c.ttl > 0 && time.Since(entry.createdAt) > c.ttl {
		return nil
	}

	return entry
}

func (c *directModeRoutingCache) set(databaseName, containerName string, entry *routingCacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := routingCacheKey(databaseName, containerName)
	entry.createdAt = time.Now()
	c.entries[key] = entry
}

func (c *directModeRoutingCache) invalidate(databaseName, containerName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := routingCacheKey(databaseName, containerName)
	delete(c.entries, key)
}

func (c *directModeRoutingCache) invalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*routingCacheEntry)
}

type directModeRouter struct {
	cache *directModeRoutingCache
}

func newDirectModeRouter(cacheTTL time.Duration) *directModeRouter {
	return &directModeRouter{
		cache: newDirectModeRoutingCache(cacheTTL),
	}
}

type routingInfo struct {
	collectionRID         string
	partitionKeyRangeID   string
	effectivePartitionKey string
}

func (r *directModeRouter) resolve(
	ctx context.Context,
	container *ContainerClient,
	partitionKey *PartitionKey,
) (*routingInfo, error) {
	databaseName := container.database.id
	containerName := container.id

	entry := r.cache.get(databaseName, containerName)
	if entry == nil {
		var err error
		entry, err = r.fetchAndCache(ctx, container)
		if err != nil {
			return nil, err
		}
	}

	pkRangeID, epkStr := r.findPKRangeIDAndEPK(partitionKey, entry)

	return &routingInfo{
		collectionRID:         entry.collectionRID,
		partitionKeyRangeID:   pkRangeID,
		effectivePartitionKey: epkStr,
	}, nil
}

func (r *directModeRouter) fetchAndCache(ctx context.Context, container *ContainerClient) (*routingCacheEntry, error) {
	containerResp, err := container.Read(ctx, nil)
	if err != nil {
		return nil, err
	}

	pkRangesResp, err := container.getPartitionKeyRanges(ctx, nil)
	if err != nil {
		return nil, err
	}

	entry := &routingCacheEntry{
		collectionRID: containerResp.ContainerProperties.ResourceID,
		pkRanges:      pkRangesResp.PartitionKeyRanges,
		pkDefinition:  &containerResp.ContainerProperties.PartitionKeyDefinition,
	}

	r.cache.set(container.database.id, container.id, entry)
	return entry, nil
}

func (r *directModeRouter) findPKRangeIDAndEPK(partitionKey *PartitionKey, entry *routingCacheEntry) (string, string) {
	if entry.pkDefinition == nil || len(entry.pkRanges) == 0 {
		if len(entry.pkRanges) > 0 {
			return entry.pkRanges[0].ID, ""
		}
		return "0", ""
	}

	effectivePK := partitionKey.computeEffectivePartitionKey(entry.pkDefinition.Kind, entry.pkDefinition.Version)
	epkStr := effectivePK.EPK

	for _, pkRange := range entry.pkRanges {
		if isEPKInRange(epkStr, pkRange.MinInclusive, pkRange.MaxExclusive) {
			return pkRange.ID, epkStr
		}
	}

	if len(entry.pkRanges) > 0 {
		return entry.pkRanges[0].ID, epkStr
	}
	return "0", epkStr
}

func isEPKInRange(epk, minInclusive, maxExclusive string) bool {
	epk = strings.ToUpper(epk)
	minInclusive = strings.ToUpper(minInclusive)
	maxExclusive = strings.ToUpper(maxExclusive)

	if minInclusive == "" {
		minInclusive = minEffectivePartitionKey
	}
	if maxExclusive == "" {
		maxExclusive = maxEffectivePartitionKey
	}

	return epk >= minInclusive && epk < maxExclusive
}

func (r *directModeRouter) invalidateCache(databaseName, containerName string) {
	r.cache.invalidate(databaseName, containerName)
}

func (r *directModeRouter) invalidateAllCaches() {
	r.cache.invalidateAll()
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"sync"
	"time"
)

// partitionKeyRangeCache provides a cache for partition key ranges of a Cosmos container.
// It allows for efficient retrieval of partition key ranges without making repeated calls to the server.
type partitionKeyRangeCache struct {
	partitionKeyRangeCache map[string]PartitionKeyRangeProperties
	resourceID             string
	mu                     sync.RWMutex
	lastRefreshed          time.Time
}

// newPartitionKeyRangeCache creates a new empty cache for partition key ranges.
func newPartitionKeyRangeCache(resourceID string) *partitionKeyRangeCache {
	return &partitionKeyRangeCache{
		partitionKeyRangeCache: make(map[string]PartitionKeyRangeProperties),
		resourceID:             resourceID,
		lastRefreshed:          time.Time{}, // Zero time indicates never refreshed
	}
}

// refresh updates the cache with the latest partition key ranges from the container.
// It acquires a write lock on the cache during the update.
func (c *partitionKeyRangeCache) refresh(ctx context.Context, container *ContainerClient) error {
	response, err := container.GetPartitionKeyRange(ctx, nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear existing cache and populate with new data
	c.partitionKeyRangeCache = make(map[string]PartitionKeyRangeProperties)
	for _, pkRange := range response.PartitionKeyRanges {
		c.partitionKeyRangeCache[pkRange.ID] = pkRange
	}

	c.resourceID = response.ResourceID
	c.lastRefreshed = time.Now()

	return nil
}

// getByID retrieves a partition key range by its ID.
// It acquires a read lock on the cache during retrieval.
func (c *partitionKeyRangeCache) getByID(id string) (PartitionKeyRangeProperties, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	pkRange, exists := c.partitionKeyRangeCache[id]
	return pkRange, exists
}

// getAll returns all partition key ranges in the cache.
// It acquires a read lock on the cache during retrieval.
func (c *partitionKeyRangeCache) getAll() []PartitionKeyRangeProperties {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ranges := make([]PartitionKeyRangeProperties, 0, len(c.partitionKeyRangeCache))
	for _, pkRange := range c.partitionKeyRangeCache {
		ranges = append(ranges, pkRange)
	}

	return ranges
}

// getByMinMax finds all partition key ranges that overlap with the specified min-max range.
// It acquires a read lock on the cache during retrieval.
func (c *partitionKeyRangeCache) getByMinMax(minInclusive, maxExclusive string) []PartitionKeyRangeProperties {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var overlappingRanges []PartitionKeyRangeProperties

	for _, pkRange := range c.partitionKeyRangeCache {
		// A range [a,b) overlaps with [c,d) if and only if a < d and c < b
		if pkRange.MinInclusive < maxExclusive && minInclusive < pkRange.MaxExclusive {
			overlappingRanges = append(overlappingRanges, pkRange)
		}
	}

	return overlappingRanges
}

// getFeedRanges converts the cached partition key ranges to FeedRange objects.
// It acquires a read lock on the cache during conversion.
func (c *partitionKeyRangeCache) getFeedRanges() []FeedRange {
	c.mu.RLock()
	defer c.mu.RUnlock()

	feedRanges := make([]FeedRange, 0, len(c.partitionKeyRangeCache))
	for _, pkRange := range c.partitionKeyRangeCache {
		feedRanges = append(feedRanges, NewFeedRange(pkRange.MinInclusive, pkRange.MaxExclusive))
	}

	return feedRanges
}

// getLastRefreshed returns the time when the cache was last refreshed.
// It acquires a read lock on the cache during retrieval.
func (c *partitionKeyRangeCache) getLastRefreshed() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.lastRefreshed
}

// needsRefresh determines if the cache needs to be refreshed based on the specified max age.
// It acquires a read lock on the cache during evaluation.
func (c *partitionKeyRangeCache) needsRefresh(maxAge time.Duration) bool {
	maxAgeAdjusted := time.Now().Add(-maxAge)
	c.mu.RLock()
	defer c.mu.RUnlock()

	// If cache has never been refreshed, it needs a refresh
	if c.lastRefreshed.IsZero() {
		return true
	}

	// Check if the cache is older than the specified max age
	return c.lastRefreshed.Before(maxAgeAdjusted)
}

// refreshIfNeeded refreshes the cache if it's older than the specified max age.
func (c *partitionKeyRangeCache) refreshIfNeeded(ctx context.Context, container *ContainerClient, maxAge time.Duration) error {
	if c.needsRefresh(maxAge) {
		return c.refresh(ctx, container)
	}
	return nil
}

// BuildPartitionKeyRangeCache creates and populates a partition key range cache for the specified container.
// It returns the populated cache and any error that occurred during cache population.
func BuildPartitionKeyRangeCache(ctx context.Context, container *ContainerClient) (*partitionKeyRangeCache, error) {
	cache := newPartitionKeyRangeCache("")
	err := cache.refresh(ctx, container)
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"sync"
)

// containerPropertiesCache provides a client-level cache of container properties
// (specifically PartitionKeyDefinition). Keyed by container link with event-driven
// invalidation (no TTL).
type containerPropertiesCache struct {
	mu      sync.RWMutex
	entries map[string]*containerPropsCacheEntry
}

type containerPropsCacheEntry struct {
	mu    sync.Mutex // single-flights refresh for this container
	props *ContainerProperties
}

func newContainerPropertiesCache() *containerPropertiesCache {
	return &containerPropertiesCache{
		entries: make(map[string]*containerPropsCacheEntry),
	}
}

// getProperties returns the cached container properties for the given container.
// If the cache is empty, it fetches from the service using the provided ContainerClient.
func (c *containerPropertiesCache) getProperties(
	ctx context.Context,
	container *ContainerClient,
) (*ContainerProperties, error) {
	containerLink := container.link

	// Fast path: read lock check
	c.mu.RLock()
	entry, exists := c.entries[containerLink]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		if entry.props != nil {
			props := entry.props
			entry.mu.Unlock()
			return props, nil
		}
		// Entry exists but props is nil (invalidated) — refresh under lock
		props, err := c.refreshEntry(ctx, container, entry)
		entry.mu.Unlock()
		return props, err
	}

	// Slow path: create entry
	c.mu.Lock()
	entry, exists = c.entries[containerLink]
	if !exists {
		entry = &containerPropsCacheEntry{}
		c.entries[containerLink] = entry
	}
	c.mu.Unlock()

	entry.mu.Lock()
	if entry.props != nil {
		props := entry.props
		entry.mu.Unlock()
		return props, nil
	}
	props, err := c.refreshEntry(ctx, container, entry)
	entry.mu.Unlock()
	return props, err
}

// invalidate removes the cached properties for the given container link,
// forcing the next access to fetch fresh data.
func (c *containerPropertiesCache) invalidate(containerLink string) {
	c.mu.RLock()
	entry, exists := c.entries[containerLink]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		entry.props = nil
		entry.mu.Unlock()
	}
}

// set directly populates the cache with the given container properties.
// This is used when a Read() call already fetched the properties.
func (c *containerPropertiesCache) set(containerLink string, props *ContainerProperties) {
	c.mu.RLock()
	entry, exists := c.entries[containerLink]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		entry.props = props
		entry.mu.Unlock()
		return
	}

	c.mu.Lock()
	entry, exists = c.entries[containerLink]
	if !exists {
		entry = &containerPropsCacheEntry{}
		c.entries[containerLink] = entry
	}
	c.mu.Unlock()

	entry.mu.Lock()
	entry.props = props
	entry.mu.Unlock()
}

// refreshEntry fetches container properties from the service.
// Caller must hold entry.mu.
func (c *containerPropertiesCache) refreshEntry(
	ctx context.Context,
	container *ContainerClient,
	entry *containerPropsCacheEntry,
) (*ContainerProperties, error) {
	resp, err := container.Read(ctx, nil)
	if err != nil {
		return nil, err
	}

	entry.props = resp.ContainerProperties
	return entry.props, nil
}

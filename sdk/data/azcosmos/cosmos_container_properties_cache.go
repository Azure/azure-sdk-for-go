// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"sync"
)

// containerPropertiesCache provides a client-level cache of container properties
// (specifically PartitionKeyDefinition). It maintains a dual index — by container
// link (name-based path) and by ResourceID (RID) — so that lookups succeed
// regardless of which identifier the caller has. When a reference is fetched or
// inserted, both indices are cross-populated.
type containerPropertiesCache struct {
	mu          sync.RWMutex
	entries     map[string]*containerPropsCacheEntry // keyed by container link
	entriesByID map[string]*containerPropsCacheEntry // keyed by ResourceID
}

type containerPropsCacheEntry struct {
	mu    sync.Mutex // single-flights refresh for this container
	props *ContainerProperties
}

func newContainerPropertiesCache() *containerPropertiesCache {
	return &containerPropertiesCache{
		entries:     make(map[string]*containerPropsCacheEntry),
		entriesByID: make(map[string]*containerPropsCacheEntry),
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
		if err == nil {
			c.updateRIDIndex(entry, props)
		}
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
	if err == nil {
		c.updateRIDIndex(entry, props)
	}
	return props, err
}

// getPropertiesByRID looks up cached container properties by ResourceID.
// Returns nil if the RID is not in the cache.
func (c *containerPropertiesCache) getPropertiesByRID(resourceID string) *ContainerProperties {
	if resourceID == "" {
		return nil
	}
	c.mu.RLock()
	entry, exists := c.entriesByID[resourceID]
	c.mu.RUnlock()

	if !exists {
		return nil
	}

	entry.mu.Lock()
	props := entry.props
	entry.mu.Unlock()
	return props
}

// invalidate removes the cached properties for the given container link,
// forcing the next access to fetch fresh data. Also removes the RID index entry.
// This is called during 410/Gone retry paths to handle scenarios where a container
// has been deleted and recreated with the same name but a new ResourceID.
func (c *containerPropertiesCache) invalidate(containerLink string) {
	c.mu.RLock()
	entry, exists := c.entries[containerLink]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		props := entry.props
		entry.props = nil
		entry.mu.Unlock()

		// Remove the RID index entry if we had cached props
		if props != nil && props.ResourceID != "" {
			c.mu.Lock()
			delete(c.entriesByID, props.ResourceID)
			c.mu.Unlock()
		}
	}
}

// set directly populates the cache with the given container properties.
// This is used when a Read() or Replace() call already fetched the properties.
// Cross-populates both the link-based and RID-based indices.
func (c *containerPropertiesCache) set(containerLink string, props *ContainerProperties) {
	c.mu.RLock()
	entry, exists := c.entries[containerLink]
	c.mu.RUnlock()

	if exists {
		entry.mu.Lock()
		entry.props = props
		entry.mu.Unlock()
	} else {
		// Slow path: upgrade to write lock and double-check.
		// We release c.mu before touching entry.mu to maintain
		// consistent lock order (c.mu → entry.mu, never reversed).
		c.mu.Lock()
		entry, exists = c.entries[containerLink]
		if !exists {
			entry = &containerPropsCacheEntry{props: props}
			c.entries[containerLink] = entry
		}
		c.mu.Unlock()

		if exists {
			entry.mu.Lock()
			entry.props = props
			entry.mu.Unlock()
		}
	}

	c.updateRIDIndex(entry, props)
}

// updateRIDIndex cross-populates the RID-based index.
// Must be called WITHOUT entry.mu held to maintain lock order (c.mu → entry.mu).
func (c *containerPropertiesCache) updateRIDIndex(entry *containerPropsCacheEntry, props *ContainerProperties) {
	if props != nil && props.ResourceID != "" {
		c.mu.Lock()
		c.entriesByID[props.ResourceID] = entry
		c.mu.Unlock()
	}
}

// refreshEntry fetches container properties directly from the service.
// This bypasses container.Read() to avoid deadlock — the caller already holds entry.mu,
// and Read() calls cache.set() which would try to re-acquire the same lock.
// It uses readContainerRaw() to share the HTTP call logic with Read().
// Caller must hold entry.mu.
// NOTE: This method must NOT acquire c.mu — callers update the RID index
// after releasing entry.mu via updateRIDIndex() to prevent lock-order inversion.
func (c *containerPropertiesCache) refreshEntry(
	ctx context.Context,
	container *ContainerClient,
	entry *containerPropsCacheEntry,
) (*ContainerProperties, error) {
	response, err := container.readContainerRaw(ctx, nil)
	if err != nil {
		return nil, err
	}

	if response.ContainerProperties == nil {
		return nil, errors.New("container properties response contained no properties")
	}

	entry.props = response.ContainerProperties

	return entry.props, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_containerPropertiesCache_newCache(t *testing.T) {
	cache := newContainerPropertiesCache()
	require.NotNil(t, cache)
	require.Empty(t, cache.entries)
}

func Test_containerPropertiesCache_invalidate_noEntry(t *testing.T) {
	cache := newContainerPropertiesCache()

	// Should not panic
	cache.invalidate("dbs/db1/colls/col1")

	cache.mu.RLock()
	_, exists := cache.entries["dbs/db1/colls/col1"]
	cache.mu.RUnlock()
	require.False(t, exists)
}

func Test_containerPropertiesCache_invalidate_existingEntry(t *testing.T) {
	cache := newContainerPropertiesCache()

	entry := &containerPropsCacheEntry{
		props: &ContainerProperties{
			ID: "col1",
			PartitionKeyDefinition: PartitionKeyDefinition{
				Paths:   []string{"/pk"},
				Version: 2,
			},
		},
	}
	cache.mu.Lock()
	cache.entries["dbs/db1/colls/col1"] = entry
	cache.mu.Unlock()

	// Verify populated
	entry.mu.Lock()
	require.NotNil(t, entry.props)
	entry.mu.Unlock()

	// Invalidate
	cache.invalidate("dbs/db1/colls/col1")

	// Verify nil
	entry.mu.Lock()
	require.Nil(t, entry.props)
	entry.mu.Unlock()
}

func Test_containerPropertiesCache_getProperties_cacheHit(t *testing.T) {
	cache := newContainerPropertiesCache()

	expectedProps := &ContainerProperties{
		ID: "col1",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Version: 2,
		},
	}
	entry := &containerPropsCacheEntry{props: expectedProps}
	cache.mu.Lock()
	cache.entries["dbs/db1/colls/col1"] = entry
	cache.mu.Unlock()

	// Create a minimal ContainerClient with link matching the cache key
	container := &ContainerClient{
		link: "dbs/db1/colls/col1",
	}

	props, err := cache.getProperties(nil, container) //nolint:staticcheck // nil context is fine for cache hit
	require.NoError(t, err)
	require.Equal(t, expectedProps, props)
}

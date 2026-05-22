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

func Test_containerPropertiesCache_set_crossPopulatesRID(t *testing.T) {
	cache := newContainerPropertiesCache()

	props := &ContainerProperties{
		ID:         "col1",
		ResourceID: "rid123",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Version: 2,
		},
	}

	cache.set("dbs/db1/colls/col1", props)

	// Should be retrievable by link
	cache.mu.RLock()
	entry, exists := cache.entries["dbs/db1/colls/col1"]
	cache.mu.RUnlock()
	require.True(t, exists)
	entry.mu.Lock()
	require.Equal(t, props, entry.props)
	entry.mu.Unlock()

	// Should also be retrievable by RID
	result := cache.getPropertiesByRID("rid123")
	require.NotNil(t, result)
	require.Equal(t, "col1", result.ID)
}

func Test_containerPropertiesCache_getPropertiesByRID_miss(t *testing.T) {
	cache := newContainerPropertiesCache()

	result := cache.getPropertiesByRID("nonexistent")
	require.Nil(t, result)
}

func Test_containerPropertiesCache_getPropertiesByRID_emptyRID(t *testing.T) {
	cache := newContainerPropertiesCache()

	result := cache.getPropertiesByRID("")
	require.Nil(t, result)
}

func Test_containerPropertiesCache_invalidate_removesRIDIndex(t *testing.T) {
	cache := newContainerPropertiesCache()

	props := &ContainerProperties{
		ID:         "col1",
		ResourceID: "rid456",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Version: 2,
		},
	}

	cache.set("dbs/db1/colls/col1", props)

	// Verify RID lookup works
	require.NotNil(t, cache.getPropertiesByRID("rid456"))

	// Invalidate by link
	cache.invalidate("dbs/db1/colls/col1")

	// RID index should be removed
	require.Nil(t, cache.getPropertiesByRID("rid456"))
}

func Test_containerPropertiesCache_set_multipleContainers(t *testing.T) {
	cache := newContainerPropertiesCache()

	props1 := &ContainerProperties{
		ID:                     "col1",
		ResourceID:             "rid1",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk1"}},
	}
	props2 := &ContainerProperties{
		ID:                     "col2",
		ResourceID:             "rid2",
		PartitionKeyDefinition: PartitionKeyDefinition{Paths: []string{"/pk2"}},
	}

	cache.set("dbs/db1/colls/col1", props1)
	cache.set("dbs/db1/colls/col2", props2)

	// Both should be retrievable by RID
	r1 := cache.getPropertiesByRID("rid1")
	r2 := cache.getPropertiesByRID("rid2")
	require.NotNil(t, r1)
	require.NotNil(t, r2)
	require.Equal(t, "col1", r1.ID)
	require.Equal(t, "col2", r2.ID)
}

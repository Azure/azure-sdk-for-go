// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryPlanCache_SetAndGet(t *testing.T) {
	cache := newQueryPlanCache(nil)

	containerLink := "dbs/testdb/colls/testcol"
	query := "SELECT * FROM c WHERE c.id = @id"
	plan := []byte(`{"queryInfo": {"hasAggregates": false}}`)

	cache.Set(containerLink, query, plan)

	cached, ok := cache.Get(containerLink, query)
	require.True(t, ok)
	require.Equal(t, plan, cached)
}

func TestQueryPlanCache_MissOnDifferentQuery(t *testing.T) {
	cache := newQueryPlanCache(nil)

	containerLink := "dbs/testdb/colls/testcol"
	query1 := "SELECT * FROM c WHERE c.id = @id"
	query2 := "SELECT * FROM c WHERE c.name = @name"
	plan := []byte(`{"queryInfo": {"hasAggregates": false}}`)

	cache.Set(containerLink, query1, plan)

	_, ok := cache.Get(containerLink, query2)
	require.False(t, ok)
}

func TestQueryPlanCache_MissOnDifferentContainer(t *testing.T) {
	cache := newQueryPlanCache(nil)

	containerLink1 := "dbs/testdb/colls/testcol1"
	containerLink2 := "dbs/testdb/colls/testcol2"
	query := "SELECT * FROM c WHERE c.id = @id"
	plan := []byte(`{"queryInfo": {"hasAggregates": false}}`)

	cache.Set(containerLink1, query, plan)

	_, ok := cache.Get(containerLink2, query)
	require.False(t, ok)
}

func TestQueryPlanCache_TTLExpiration(t *testing.T) {
	opts := &QueryPlanCacheOptions{
		MaxSize: 100,
		TTL:     50 * time.Millisecond,
	}
	cache := newQueryPlanCache(opts)

	containerLink := "dbs/testdb/colls/testcol"
	query := "SELECT * FROM c"
	plan := []byte(`{"queryInfo": {}}`)

	cache.Set(containerLink, query, plan)

	cached, ok := cache.Get(containerLink, query)
	require.True(t, ok)
	require.Equal(t, plan, cached)

	time.Sleep(60 * time.Millisecond)

	_, ok = cache.Get(containerLink, query)
	require.False(t, ok)
}

func TestQueryPlanCache_LRUEviction(t *testing.T) {
	opts := &QueryPlanCacheOptions{
		MaxSize: 2,
		TTL:     5 * time.Minute,
	}
	cache := newQueryPlanCache(opts)

	containerLink := "dbs/testdb/colls/testcol"
	query1 := "SELECT * FROM c WHERE c.id = 1"
	query2 := "SELECT * FROM c WHERE c.id = 2"
	query3 := "SELECT * FROM c WHERE c.id = 3"
	plan := []byte(`{"queryInfo": {}}`)

	cache.Set(containerLink, query1, plan)
	cache.Set(containerLink, query2, plan)

	cache.Get(containerLink, query1)

	cache.Set(containerLink, query3, plan)

	_, ok := cache.Get(containerLink, query1)
	require.True(t, ok, "query1 should still be cached (recently accessed)")

	_, ok = cache.Get(containerLink, query2)
	require.False(t, ok, "query2 should be evicted (LRU)")

	_, ok = cache.Get(containerLink, query3)
	require.True(t, ok, "query3 should be cached")
}

func TestQueryPlanCache_Remove(t *testing.T) {
	cache := newQueryPlanCache(nil)

	containerLink := "dbs/testdb/colls/testcol"
	query := "SELECT * FROM c"
	plan := []byte(`{"queryInfo": {}}`)

	cache.Set(containerLink, query, plan)

	_, ok := cache.Get(containerLink, query)
	require.True(t, ok)

	cache.Remove(containerLink, query)

	_, ok = cache.Get(containerLink, query)
	require.False(t, ok)
}

func TestQueryPlanCache_Clear(t *testing.T) {
	cache := newQueryPlanCache(nil)

	containerLink := "dbs/testdb/colls/testcol"
	query1 := "SELECT * FROM c WHERE c.id = 1"
	query2 := "SELECT * FROM c WHERE c.id = 2"
	plan := []byte(`{"queryInfo": {}}`)

	cache.Set(containerLink, query1, plan)
	cache.Set(containerLink, query2, plan)

	require.Equal(t, 2, cache.Size())

	cache.Clear()

	require.Equal(t, 0, cache.Size())

	_, ok := cache.Get(containerLink, query1)
	require.False(t, ok)
}

func TestQueryPlanCache_Size(t *testing.T) {
	cache := newQueryPlanCache(nil)

	require.Equal(t, 0, cache.Size())

	containerLink := "dbs/testdb/colls/testcol"
	plan := []byte(`{"queryInfo": {}}`)

	cache.Set(containerLink, "query1", plan)
	require.Equal(t, 1, cache.Size())

	cache.Set(containerLink, "query2", plan)
	require.Equal(t, 2, cache.Size())

	cache.Remove(containerLink, "query1")
	require.Equal(t, 1, cache.Size())
}

func TestQueryPlanCache_DefaultOptions(t *testing.T) {
	opts := DefaultQueryPlanCacheOptions()

	require.Equal(t, 5000, opts.MaxSize)
	require.Equal(t, 5*time.Minute, opts.TTL)
}

func TestQueryPlanCache_CacheKeyConsistency(t *testing.T) {
	cache := newQueryPlanCache(nil)

	containerLink := "dbs/testdb/colls/testcol"
	query := "SELECT * FROM c WHERE c.id = @id AND c.name = @name"

	key1 := cache.cacheKey(containerLink, query)
	key2 := cache.cacheKey(containerLink, query)

	require.Equal(t, key1, key2, "same inputs should produce same cache key")

	key3 := cache.cacheKey(containerLink, query+" ")
	require.NotEqual(t, key1, key3, "different queries should produce different keys")
}

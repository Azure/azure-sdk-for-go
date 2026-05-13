// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeEndpoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://MyAccount.documents.azure.com/", "https://myaccount.documents.azure.com"},
		{"https://MyAccount.documents.azure.com", "https://myaccount.documents.azure.com"},
		{"https://MYACCOUNT.DOCUMENTS.AZURE.COM///", "https://myaccount.documents.azure.com"},
		{"https://localhost:8081/", "https://localhost:8081"},
		{"https://localhost:8081", "https://localhost:8081"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			require.Equal(t, tt.expected, normalizeEndpoint(tt.input))
		})
	}
}

func TestAcquireCaches_SameEndpoint_ReturnsSameInstance(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	set1 := acquireCaches("https://account1.documents.azure.com/")
	set2 := acquireCaches("https://account1.documents.azure.com/")
	set3 := acquireCaches("https://Account1.Documents.Azure.Com") // different case, no trailing slash

	require.Same(t, set1, set2, "same endpoint should return same cache set")
	require.Same(t, set1, set3, "normalization should make these equivalent")
	require.Equal(t, int64(3), set1.refCount.Load())
}

func TestAcquireCaches_DifferentEndpoints_ReturnDifferentInstances(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	set1 := acquireCaches("https://account1.documents.azure.com")
	set2 := acquireCaches("https://account2.documents.azure.com")

	require.NotSame(t, set1, set2, "different endpoints should return different cache sets")
	require.NotSame(t, set1.containerCache, set2.containerCache)
	require.NotSame(t, set1.pkRangeCache, set2.pkRangeCache)
}

func TestReleaseCaches_RemovesEntryWhenZeroRefs(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	set1 := acquireCaches(endpoint)
	_ = acquireCaches(endpoint) // refCount = 2

	releaseCaches(endpoint) // refCount = 1
	require.Equal(t, int64(1), set1.refCount.Load())

	// Entry should still be in the registry
	val, ok := globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.True(t, ok)
	require.Same(t, set1, val.(*sharedCacheSet))

	releaseCaches(endpoint) // refCount = 0 → removed
	_, ok = globalCacheRegistry.Load(normalizeEndpoint(endpoint))
	require.False(t, ok, "entry should be removed when refCount reaches 0")
}

func TestReleaseCaches_NewAcquireAfterFullRelease_CreatesNew(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	set1 := acquireCaches(endpoint)
	releaseCaches(endpoint) // refCount = 0, removed

	set2 := acquireCaches(endpoint) // new instance
	require.NotSame(t, set1, set2, "should create a fresh cache set after full release")
}

func TestAcquireCaches_ConcurrentSafe(t *testing.T) {
	resetGlobalCacheRegistry()
	defer resetGlobalCacheRegistry()

	endpoint := "https://account1.documents.azure.com"
	const goroutines = 100

	var wg sync.WaitGroup
	results := make([]*sharedCacheSet, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			results[idx] = acquireCaches(endpoint)
		}(i)
	}
	wg.Wait()

	// All goroutines should have gotten the same instance
	for i := 1; i < goroutines; i++ {
		require.Same(t, results[0], results[i], "concurrent acquires should return same instance")
	}
	require.Equal(t, int64(goroutines), results[0].refCount.Load())
}

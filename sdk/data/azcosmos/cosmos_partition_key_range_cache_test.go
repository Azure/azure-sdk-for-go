// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
	"time"
)

func TestNewPartitionKeyRangeCache(t *testing.T) {
	resourceID := "test-resource-id"
	cache := newPartitionKeyRangeCache(resourceID)

	if cache.resourceID != resourceID {
		t.Errorf("Expected resourceID to be %s, got %s", resourceID, cache.resourceID)
	}

	if !cache.lastRefreshed.IsZero() {
		t.Errorf("Expected lastRefreshed to be zero time, got %v", cache.lastRefreshed)
	}

	if len(cache.partitionKeyRangeCache) != 0 {
		t.Errorf("Expected cache to be empty, got %d entries", len(cache.partitionKeyRangeCache))
	}
}

func TestPartitionKeyRangeCacheNeedsRefresh(t *testing.T) {
	cache := newPartitionKeyRangeCache("test-resource-id")

	// Never refreshed cache should need refresh
	if !cache.needsRefresh(1 * time.Hour) {
		t.Error("Expected needsRefresh to return true for never refreshed cache")
	}

	// Set last refreshed to now
	cache.lastRefreshed = time.Now()

	// Cache refreshed just now shouldn't need refresh with 2 hour max age
	if cache.needsRefresh(2 * time.Hour) {
		t.Error("Expected needsRefresh to return false for recently refreshed cache")
	}

	// Cache should need refresh with 0 max age
	if !cache.needsRefresh(0) {
		t.Error("Expected needsRefresh to return true with 0 max age")
	}

	// Set last refreshed to 2 hours ago
	cache.lastRefreshed = time.Now().Add(-2 * time.Hour)

	// Cache refreshed 2 hours ago should need refresh with 1 hour max age
	if !cache.needsRefresh(1 * time.Hour) {
		t.Error("Expected needsRefresh to return true for old cache")
	}
}

func TestGetByID(t *testing.T) {
	cache := newPartitionKeyRangeCache("test-resource-id")

	// Create test range
	testRange := PartitionKeyRangeProperties{
		ID:           "test-id",
		MinInclusive: "00",
		MaxExclusive: "FF",
	}

	// Add to cache
	cache.partitionKeyRangeCache[testRange.ID] = testRange

	// Test retrieval
	retrievedRange, exists := cache.getByID("test-id")
	if !exists {
		t.Error("Expected to find partition key range with ID 'test-id'")
	}

	if retrievedRange.ID != testRange.ID {
		t.Errorf("Expected ID to be %s, got %s", testRange.ID, retrievedRange.ID)
	}

	// Test non-existent ID
	_, exists = cache.getByID("non-existent-id")
	if exists {
		t.Error("Expected not to find partition key range with ID 'non-existent-id'")
	}
}

func TestGetByMinMax(t *testing.T) {
	cache := newPartitionKeyRangeCache("test-resource-id")

	// Create test ranges
	ranges := []PartitionKeyRangeProperties{
		{ID: "1", MinInclusive: "00", MaxExclusive: "20"},
		{ID: "2", MinInclusive: "20", MaxExclusive: "40"},
		{ID: "3", MinInclusive: "40", MaxExclusive: "60"},
		{ID: "4", MinInclusive: "60", MaxExclusive: "80"},
		{ID: "5", MinInclusive: "80", MaxExclusive: "FF"},
	}

	// Add to cache
	for _, r := range ranges {
		cache.partitionKeyRangeCache[r.ID] = r
	}

	// Test retrievals
	// Case 1: Exact match
	overlapping := cache.getByMinMax("20", "40")
	if len(overlapping) != 1 || overlapping[0].ID != "2" {
		t.Errorf("Expected 1 overlapping range with ID '2', got %d ranges", len(overlapping))
	}

	// Case 2: Overlapping multiple ranges
	overlapping = cache.getByMinMax("10", "50")
	if len(overlapping) != 3 {
		t.Errorf("Expected 3 overlapping ranges, got %d ranges", len(overlapping))
	}

	// Case 3: No overlap
	overlapping = cache.getByMinMax("FF", "FF1")
	if len(overlapping) != 0 {
		t.Errorf("Expected 0 overlapping ranges, got %d ranges", len(overlapping))
	}

	// Case 4: Full range
	overlapping = cache.getByMinMax("00", "FF")
	if len(overlapping) != 5 {
		t.Errorf("Expected 5 overlapping ranges, got %d ranges", len(overlapping))
	}
}

func TestGetFeedRanges(t *testing.T) {
	cache := newPartitionKeyRangeCache("test-resource-id")

	// Create test ranges
	ranges := []PartitionKeyRangeProperties{
		{ID: "1", MinInclusive: "00", MaxExclusive: "20"},
		{ID: "2", MinInclusive: "20", MaxExclusive: "40"},
	}

	// Add to cache
	for _, r := range ranges {
		cache.partitionKeyRangeCache[r.ID] = r
	}

	// Get feed ranges
	feedRanges := cache.getFeedRanges()

	if len(feedRanges) != 2 {
		t.Errorf("Expected 2 feed ranges, got %d ranges", len(feedRanges))
	}

	// Verify feed ranges match partition key ranges
	for i, fr := range feedRanges {
		found := false
		for _, pkr := range ranges {
			if fr.MinInclusive == pkr.MinInclusive && fr.MaxExclusive == pkr.MaxExclusive {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Feed range %d (%s-%s) doesn't match any partition key range",
				i, fr.MinInclusive, fr.MaxExclusive)
		}
	}
}

// The following test requires a real Cosmos DB container
// Uncomment and modify to run against a real container
/*
func TestRefresh(t *testing.T) {
	// Skip in CI/short test runs
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Create test client
	client, err := NewClientWithKey("<endpoint>", "<key>")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Get database and container
	database := client.NewDatabase("<database>")
	container := database.NewContainer("<container>")

	// Create cache
	cache := newPartitionKeyRangeCache("")

	// Refresh cache
	err = cache.refresh(context.Background(), container)
	if err != nil {
		t.Fatalf("Failed to refresh cache: %v", err)
	}

	// Verify cache was populated
	if len(cache.partitionKeyRangeCache) == 0 {
		t.Error("Expected cache to be populated after refresh")
	}

	// Verify last refreshed time was updated
	if cache.lastRefreshed.IsZero() {
		t.Error("Expected lastRefreshed to be updated after refresh")
	}
}
*/

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/mock"
)

const (
	// These keys have been tested to ensure they end up in separate PK ranges for a 40000 RU container with the test data inserted.
	// Conveniently, they also have descriptive names.

	partition1Key string = "partition1"
	partition2Key string = "partition2"
	partition3Key string = "partition3"

	ruCountForMultiplePartitions int32 = 40000

	partitionCount    int = 3
	itemsPerPartition int = 10
)

var partitionKeys = [...]string{partition1Key, partition2Key, partition3Key}

func generateMockItem(partitionIndex int, itemIndex int) mock.MockItem {
	// Reuse the partitionKeys defined above so generated items match the test partition names.
	pk := partitionKeys[partitionIndex]
	return mock.MockItem{
		// make sure id and merge order are not the same
		ID:           strconv.Itoa(partitionIndex*itemsPerPartition + itemIndex + 1),
		PartitionKey: pk,
		// The merge order should alternate between partitions
		MergeOrder: partitionIndex + itemIndex*partitionCount,
	}
}

func generateMockItems(partitions int, itemsPerPartition int) []mock.MockItem {
	items := make([]mock.MockItem, 0, partitions*itemsPerPartition)
	for i := 0; i < partitions; i++ {
		for j := 0; j < itemsPerPartition; j++ {
			items = append(items, generateMockItem(i, j))
		}
	}
	return items
}

func createTestItems(t *testing.T, database *DatabaseClient, items []mock.MockItem) (*ContainerClient, error) {
	properties := ContainerProperties{
		ID: "TestContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/partitionKey"},
		},
	}

	// Force the creation of a container with multiple physical partitions
	throughput := NewManualThroughputProperties(ruCountForMultiplePartitions)
	_, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{
		ThroughputProperties: &throughput,
	})
	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}

	container, err := database.NewContainer("TestContainer")
	if err != nil {
		t.Fatalf("failed to create container client: %v", err)
	}
	for _, item := range items {
		serializedItem, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		_, err = container.UpsertItem(context.TODO(), NewPartitionKeyString(item.PartitionKey), serializedItem, nil)
		if err != nil {
			return nil, err
		}
	}

	return container, nil
}

func TestQueryViaQueryEngine(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{}))
	database := emulatorTests.createDatabase(t, context.TODO(), client, "TestQueryViaQueryEngine")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	// generate items and create the container with them
	items := generateMockItems(3, 10)
	container, err := createTestItems(t, database, items)
	if err != nil {
		t.Fatalf("Failed to create test items: %v", err)
	}

	options := &QueryOptions{
		QueryEngine: mock.NewMockQueryEngine(),
	}
	pager := container.NewQueryItemsPager("SELECT * FROM c ORDER BY c.mergeOrder", NewPartitionKey(), options)

	expectedPartitionId := 0
	expectedMergeOrder := 0
	itemCount := 0
	for pager.More() {
		response, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to get next page: %v", err)
		}
		for i, item := range response.Items {
			itemCount++
			var testItem mock.MockItem
			if err := json.Unmarshal(item, &testItem); err != nil {
				t.Fatalf("Failed to unmarshal item: %v", err)
			}

			if testItem.PartitionKey != partitionKeys[expectedPartitionId] {
				t.Fatalf("Expected partition key of item #%d with ID %s to be %s, got %s", i, testItem.ID, partitionKeys[expectedPartitionId], testItem.PartitionKey)
			}

			if testItem.MergeOrder != expectedMergeOrder {
				t.Fatalf("Expected merge order of item #%d with ID %s to be %d, got %d", i, testItem.ID, expectedMergeOrder, testItem.MergeOrder)
			}

			expectedPartitionId = (expectedPartitionId + 1) % partitionCount
			expectedMergeOrder++
		}
	}

	if itemCount != partitionCount*itemsPerPartition {
		t.Fatalf("Expected %d items, got %d", partitionCount*itemsPerPartition, itemCount)
	}
}

func TestQueryOverrideWithoutParameters(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{}))
	database := emulatorTests.createDatabase(t, context.TODO(), client, "TestQueryOverrideWithoutParameters")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	items := generateMockItems(3, 5)
	container, err := createTestItems(t, database, items)
	if err != nil {
		t.Fatalf("Failed to create test items: %v", err)
	}

	override := "SELECT * FROM c WHERE c.id = 'override'"
	cfg := &mock.QueryRequestConfig{Query: &override, IncludeParameters: false}
	engine := mock.WithQueryRequestConfig(cfg)

	options := &QueryOptions{QueryEngine: engine}
	pager := container.NewQueryItemsPager("SELECT * FROM c WHERE c.id = @param1", NewPartitionKey(), options)

	resultItems := make([]mock.MockItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("failed to get next page: %v", err)
		}
		for _, it := range resp.Items {
			var mi mock.MockItem
			if err := json.Unmarshal(it, &mi); err != nil {
				t.Fatalf("failed to unmarshal item: %v", err)
			}
			resultItems = append(resultItems, mi)
		}
	}

	if len(resultItems) != 0 {
		t.Fatalf("expected 0 results for override query without parameters, got %d", len(resultItems))
	}
}

func TestQueryOverrideWithParameters(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{}))
	database := emulatorTests.createDatabase(t, context.TODO(), client, "TestQueryOverrideWithParameters")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	items := generateMockItems(3, 5)
	container, err := createTestItems(t, database, items)
	if err != nil {
		t.Fatalf("Failed to create test items: %v", err)
	}

	override := "SELECT * FROM c WHERE c.mergeOrder = @targetOrder"
	cfg := &mock.QueryRequestConfig{Query: &override, IncludeParameters: true}
	engine := mock.WithQueryRequestConfig(cfg)

	// choose a target merge order present in the test data: use the first item's merge order (0)
	if strconv.Itoa(items[0].MergeOrder) == items[0].ID {
		t.Fatalf("Test data generation error: item ID and MergeOrder should not match")
	}
	target := items[0].MergeOrder

	// Build original query that uses a parameter which should be forwarded to the override when includeParameters=true
	options := &QueryOptions{QueryEngine: engine, QueryParameters: []QueryParameter{{Name: "@targetOrder", Value: target}}}
	pager := container.NewQueryItemsPager("SELECT * FROM c WHERE c.id = @targetOrder", NewPartitionKey(), options)

	resultItems := make([]mock.MockItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("failed to get next page: %v", err)
		}
		for _, it := range resp.Items {
			var mi mock.MockItem
			if err := json.Unmarshal(it, &mi); err != nil {
				t.Fatalf("failed to unmarshal item: %v", err)
			}
			resultItems = append(resultItems, mi)
		}
	}

	// Expect items whose MergeOrder == target
	expected := 0
	for _, it := range resultItems {
		if it.MergeOrder == target {
			expected++
		}
	}
	if expected == 0 {
		t.Fatalf("expected at least one matching item for target merge order %d", target)
	}
}

func TestNoQueryOverrideUsesOriginal(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{}))
	database := emulatorTests.createDatabase(t, context.TODO(), client, "TestNoQueryOverrideUsesOriginal")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	items := generateMockItems(3, 5)
	container, err := createTestItems(t, database, items)
	if err != nil {
		t.Fatalf("Failed to create test items: %v", err)
	}

	// No override: Query = nil
	cfg := &mock.QueryRequestConfig{Query: nil, IncludeParameters: false}
	engine := mock.WithQueryRequestConfig(cfg)

	// We will query by mergeOrder using a parameter
	target := 0
	options := &QueryOptions{QueryEngine: engine, QueryParameters: []QueryParameter{{Name: "@targetOrder", Value: target}}}
	pager := container.NewQueryItemsPager("SELECT * FROM c WHERE c.mergeOrder = @targetOrder", NewPartitionKey(), options)

	resultItems := make([]mock.MockItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("failed to get next page: %v", err)
		}
		for _, it := range resp.Items {
			var mi mock.MockItem
			if err := json.Unmarshal(it, &mi); err != nil {
				t.Fatalf("failed to unmarshal item: %v", err)
			}
			resultItems = append(resultItems, mi)
		}
	}

	// Expect items whose MergeOrder == target
	expected := 0
	for _, it := range resultItems {
		if it.MergeOrder == target {
			expected++
		}
	}
	if expected == 0 {
		t.Fatalf("expected at least one matching item for target merge order %d", target)
	}
}

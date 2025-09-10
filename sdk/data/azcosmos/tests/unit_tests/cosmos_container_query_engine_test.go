// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	azcosmosinternal "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal"
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

func createTestItems(t *testing.T, database *DatabaseClient) (*ContainerClient, error) {
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

	for i := 0; i < partitionCount; i++ {
		for j := 0; j < itemsPerPartition; j++ {
			item := azcosmosinternal.MockItem{
				ID:           strconv.Itoa(i*itemsPerPartition + j),
				PartitionKey: partitionKeys[i],

				// The merge order should alternate between partitions
				MergeOrder: i + j*partitionCount,
			}
			serializedItem, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			_, err = container.UpsertItem(context.TODO(), NewPartitionKeyString(item.PartitionKey), serializedItem, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return container, nil
}

func TestQueryViaQueryEngine(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{}))
	database := emulatorTests.createDatabase(t, context.TODO(), client, "TestQueryViaQueryEngine")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	container, err := createTestItems(t, database)
	if err != nil {
		t.Fatalf("Failed to create test items: %v", err)
	}

	options := &QueryOptions{
		QueryEngine: azcosmosinternal.NewMockQueryEngine(),
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
			var testItem azcosmosinternal.MockItem
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

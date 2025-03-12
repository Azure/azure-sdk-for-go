package azcosmos

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	cosmosmock "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/mock"
)

const (
	// These keys have been tested to ensure they end up in separate PK ranges for an otherwise-empty 40000 RU container, using https://github.com/analogrelay/cosmos-pkhash
	// Conveniently, they also have descriptive names
	// These are the values, their EPK, and the range they are in:
	//   partition1 => 069AFC9298E63A9D08429B60A594626C in 0
	//   partition2 => 1721B45D14F4AF59263B278F53573BFF in 2
	//   partition3 => 2DF2EB69837D3D4E551791CE80D2636F in 5

	partition1Key string = "partition1"
	partition2Key string = "partition2"
	partition3Key string = "partition3"

	ruCount int32 = 40000

	partitionCount    int = 3
	itemsPerPartition int = 10
)

var partitionKeys = [...]string{partition1Key, partition2Key, partition3Key}

func createTestItems(t *testing.T, database *DatabaseClient, emulatorTests *emulatorTests, client *Client) (*ContainerClient, error) {
	properties := ContainerProperties{
		ID: "TestContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/partitionKey"},
		},
	}

	// Force the creation of a container with multiple physical partitions
	throughput := NewManualThroughputProperties(40000)
	_, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{
		ThroughputProperties: &throughput,
	})

	container, err := database.NewContainer("TestContainer")

	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}

	for i := 0; i < partitionCount; i++ {
		for j := 0; j < itemsPerPartition; j++ {
			item := cosmosmock.MockItem{
				ID:           strconv.Itoa(i*itemsPerPartition + j),
				PartitionKey: partitionKeys[i],

				// The merge order should alternate between partitions
				MergeOrder: i + j*partitionCount,
			}
			serializedItem, err := json.Marshal(item)
			if err != nil {
				return nil, err
			}
			container.UpsertItem(context.TODO(), NewPartitionKeyString(item.PartitionKey), serializedItem, nil)
		}
	}

	return container, nil
}

func TestQueryViaQueryEngine(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{}))
	database := emulatorTests.createDatabase(t, context.TODO(), client, "TestQueryViaQueryEngine")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	container, err := createTestItems(t, database, emulatorTests, client)
	if err != nil {
		t.Fatalf("Failed to create test items: %v", err)
	}

	options := &QueryOptions{
		UnstablePreviewQueryEngine: cosmosmock.NewMockQueryEngine(),
	}
	pager := container.NewQueryItemsPager("SELECT * FROM c ORDER BY c.mergeOrder", NewPartitionKey(), options)

	expectedPartitionId := 0
	expectedMergeOrder := 0
	for pager.More() {
		response, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to get next page: %v", err)
		}
		for i, item := range response.Items {
			var testItem cosmosmock.MockItem
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
			expectedMergeOrder = expectedMergeOrder + 1
		}
	}
}

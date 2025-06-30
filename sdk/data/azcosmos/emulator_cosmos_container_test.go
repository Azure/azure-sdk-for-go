// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestContainerCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_container aContainer", "read_container aContainer", "query_containers containerCRUD", "replace_container aContainer", "read_container_throughput aContainer", "replace_container_throughput aContainer", "delete_container aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
		IndexingPolicy: &IndexingPolicy{
			IncludedPaths: []IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []ExcludedPath{
				{Path: "/\"_etag\"/?"},
			},
			Automatic:    true,
			IndexingMode: IndexingModeConsistent,
		},
	}

	throughput := NewManualThroughputProperties(400)

	resp, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if resp.ContainerProperties.ID != properties.ID {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container, _ := database.NewContainer("aContainer")
	resp, err = container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	receivedIds := []string{}
	opt := QueryContainersOptions{
		QueryParameters: []QueryParameter{
			{"@id", "aContainer"},
		},
	}
	queryPager := database.NewQueryContainersPager("SELECT * FROM root r WHERE r.id = @id", &opt)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("Failed to query databases: %v", err)
		}

		for _, db := range queryResponse.Containers {
			receivedIds = append(receivedIds, db.ID)
		}
	}

	if len(receivedIds) != 1 {
		t.Fatalf("Expected 1 container, got %d", len(receivedIds))
	}

	updatedProperties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
		IndexingPolicy: &IndexingPolicy{
			IncludedPaths: []IncludedPath{},
			ExcludedPaths: []ExcludedPath{},
			Automatic:     false,
			IndexingMode:  IndexingModeNone,
		},
	}

	resp, err = container.Replace(context.TODO(), updatedProperties, nil)
	if err != nil {
		t.Fatalf("Failed to update container: %v", err)
	}

	throughputResponse, err := container.ReadThroughput(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read throughput: %v", err)
	}

	mt, hasManualThroughput := throughputResponse.ThroughputProperties.ManualThroughput()
	if !hasManualThroughput {
		t.Fatalf("Expected manual throughput to be available")
	}

	if mt != 400 {
		t.Errorf("Unexpected throughput: %v", mt)
	}

	newScale := NewManualThroughputProperties(500)
	throughputResponse, err = container.ReplaceThroughput(context.TODO(), newScale, nil)
	if err != nil {
		t.Fatalf("Failed to replace throughput: %v", err)
	}

	mt, hasManualThroughput = throughputResponse.ThroughputProperties.ManualThroughput()
	if !hasManualThroughput {
		t.Fatalf("Expected manual throughput to be available")
	}

	if mt != 500 {
		t.Errorf("Unexpected throughput: %v", mt)
	}

	resp, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

func TestContainerAutoscaleCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_container aContainer", "read_container aContainer", "read_container_throughput aContainer", "replace_container_throughput aContainer", "delete_container aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
		IndexingPolicy: &IndexingPolicy{
			IncludedPaths: []IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []ExcludedPath{
				{Path: "/\"_etag\"/?"},
			},
			Automatic:    true,
			IndexingMode: IndexingModeConsistent,
		},
	}

	throughput := NewAutoscaleThroughputProperties(5000)

	resp, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if resp.ContainerProperties.ID != properties.ID {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container, _ := database.NewContainer("aContainer")
	resp, err = container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	throughputResponse, err := container.ReadThroughput(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read throughput: %v", err)
	}

	maxru, hasAutoscale := throughputResponse.ThroughputProperties.AutoscaleMaxThroughput()
	if !hasAutoscale {
		t.Fatalf("Expected autoscale throughput to be available")
	}

	if maxru != 5000 {
		t.Errorf("Unexpected throughput: %v", maxru)
	}

	newScale := NewAutoscaleThroughputProperties(10000)
	_, err = container.ReplaceThroughput(context.TODO(), newScale, nil)
	if err != nil {
		t.Errorf("Failed to read throughput: %v", err)
	}

	resp, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

func TestContainerVectorSearch(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_container vectorContainer", "read_container vectorContainer", "delete_container vectorContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "vectorSearch")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	// Create container with vector embedding and indexing policies
	properties := ContainerProperties{
		ID: "vectorContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
		VectorEmbeddingPolicy: &VectorEmbeddingPolicy{
			VectorEmbeddings: []VectorEmbedding{
				{
					Path:             "/embedding",
					DataType:         VectorDataTypeFloat32,
					DistanceFunction: VectorDistanceFunctionCosine,
					Dimensions:       3,
				},
				{
					Path:             "/textEmbedding",
					DataType:         VectorDataTypeFloat32,
					DistanceFunction: VectorDistanceFunctionDotProduct,
					Dimensions:       384, // Use smaller dimension for flat index compatibility
				},
			},
		},
		IndexingPolicy: &IndexingPolicy{
			Automatic:    true,
			IndexingMode: IndexingModeConsistent,
			IncludedPaths: []IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []ExcludedPath{
				{Path: "/\"_etag\"/?"},
				{Path: "/embedding/*"},     // Exclude vector path from standard indexing
				{Path: "/textEmbedding/*"}, // Exclude vector path from standard indexing
			},
			VectorIndexes: []VectorIndex{
				{
					Path: "/embedding",
					Type: VectorIndexTypeFlat,
				},
				{
					Path: "/textEmbedding",
					Type: VectorIndexTypeFlat, // Use flat instead of diskANN for emulator compatibility
				},
			},
		},
	}

	throughput := NewManualThroughputProperties(400)
	_, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		t.Fatalf("Failed to create vector container: %v", err)
	}

	container, _ := database.NewContainer("vectorContainer")

	// Read the container back to validate properties were set correctly
	resp, err := container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	readProperties := resp.ContainerProperties

	// Validate basic properties
	if readProperties.ID != properties.ID {
		t.Errorf("Expected container ID %s, got %s", properties.ID, readProperties.ID)
	}

	// Validate vector embedding policy
	if readProperties.VectorEmbeddingPolicy == nil {
		t.Fatalf("Expected VectorEmbeddingPolicy to be set, but it was nil")
	}

	if len(readProperties.VectorEmbeddingPolicy.VectorEmbeddings) != 2 {
		t.Fatalf("Expected 2 vector embeddings, got %d", len(readProperties.VectorEmbeddingPolicy.VectorEmbeddings))
	}

	// Validate first vector embedding
	embedding1 := readProperties.VectorEmbeddingPolicy.VectorEmbeddings[0]
	if embedding1.Path != "/embedding" {
		t.Errorf("Expected first embedding path /embedding, got %s", embedding1.Path)
	}
	if embedding1.DataType != VectorDataTypeFloat32 {
		t.Errorf("Expected first embedding data type float32, got %s", embedding1.DataType)
	}
	if embedding1.DistanceFunction != VectorDistanceFunctionCosine {
		t.Errorf("Expected first embedding distance function cosine, got %s", embedding1.DistanceFunction)
	}
	if embedding1.Dimensions != 3 {
		t.Errorf("Expected first embedding dimensions 3, got %d", embedding1.Dimensions)
	}

	// Validate second vector embedding
	embedding2 := readProperties.VectorEmbeddingPolicy.VectorEmbeddings[1]
	if embedding2.Path != "/textEmbedding" {
		t.Errorf("Expected second embedding path /textEmbedding, got %s", embedding2.Path)
	}
	if embedding2.DataType != VectorDataTypeFloat32 {
		t.Errorf("Expected second embedding data type float32, got %s", embedding2.DataType)
	}
	if embedding2.DistanceFunction != VectorDistanceFunctionDotProduct {
		t.Errorf("Expected second embedding distance function dotproduct, got %s", embedding2.DistanceFunction)
	}
	if embedding2.Dimensions != 384 {
		t.Errorf("Expected second embedding dimensions 384, got %d", embedding2.Dimensions)
	}

	// Validate vector indexing policy
	if readProperties.IndexingPolicy == nil {
		t.Fatalf("Expected IndexingPolicy to be set, but it was nil")
	}

	if len(readProperties.IndexingPolicy.VectorIndexes) != 2 {
		t.Fatalf("Expected 2 vector indexes, got %d", len(readProperties.IndexingPolicy.VectorIndexes))
	}

	// Validate first vector index
	index1 := readProperties.IndexingPolicy.VectorIndexes[0]
	if index1.Path != "/embedding" {
		t.Errorf("Expected first vector index path /embedding, got %s", index1.Path)
	}
	if index1.Type != VectorIndexTypeFlat {
		t.Errorf("Expected first vector index type flat, got %s", index1.Type)
	}

	// Validate second vector index
	index2 := readProperties.IndexingPolicy.VectorIndexes[1]
	if index2.Path != "/textEmbedding" {
		t.Errorf("Expected second vector index path /textEmbedding, got %s", index2.Path)
	}
	if index2.Type != VectorIndexTypeFlat {
		t.Errorf("Expected second vector index type flat, got %s", index2.Type)
	}

	// Clean up
	_, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

func TestContainerFullTextSearch(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_container fullTextContainer", "read_container fullTextContainer", "delete_container fullTextContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "fullTextSearch")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	// Create container with full-text policy and indexing
	properties := ContainerProperties{
		ID: "fullTextContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
		FullTextPolicy: &FullTextPolicy{
			DefaultLanguage: "en-US",
			FullTextPaths: []FullTextPath{
				{
					Path:     "/title",
					Language: "en-US",
				},
				{
					Path:     "/description",
					Language: "en-US",
				},
			},
		},
		IndexingPolicy: &IndexingPolicy{
			Automatic:    true,
			IndexingMode: IndexingModeConsistent,
			IncludedPaths: []IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []ExcludedPath{
				{Path: "/\"_etag\"/?"},
			},
			FullTextIndexes: []FullTextIndex{
				{
					Path: "/title",
				},
				{
					Path: "/description",
				},
			},
		},
	}

	throughput := NewManualThroughputProperties(400)
	_, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		t.Fatalf("Failed to create full-text container: %v", err)
	}

	container, _ := database.NewContainer("fullTextContainer")

	// Read the container back to validate properties were set correctly
	resp, err := container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	readProperties := resp.ContainerProperties

	// Validate basic properties
	if readProperties.ID != properties.ID {
		t.Errorf("Expected container ID %s, got %s", properties.ID, readProperties.ID)
	}

	// Validate full-text policy
	if readProperties.FullTextPolicy == nil {
		t.Fatalf("Expected FullTextPolicy to be set, but it was nil")
	}

	if readProperties.FullTextPolicy.DefaultLanguage != "en-US" {
		t.Errorf("Expected default language en-US, got %s", readProperties.FullTextPolicy.DefaultLanguage)
	}

	if len(readProperties.FullTextPolicy.FullTextPaths) != 2 {
		t.Fatalf("Expected 2 full text paths, got %d", len(readProperties.FullTextPolicy.FullTextPaths))
	}

	// Validate first full text path
	path1 := readProperties.FullTextPolicy.FullTextPaths[0]
	if path1.Path != "/title" {
		t.Errorf("Expected first path /title, got %s", path1.Path)
	}
	if path1.Language != "en-US" {
		t.Errorf("Expected first path language en-US, got %s", path1.Language)
	}

	// Validate second full text path
	path2 := readProperties.FullTextPolicy.FullTextPaths[1]
	if path2.Path != "/description" {
		t.Errorf("Expected second path /description, got %s", path2.Path)
	}
	if path2.Language != "en-US" {
		t.Errorf("Expected second path language en-US, got %s", path2.Language)
	}

	// Validate full-text indexing policy
	if readProperties.IndexingPolicy == nil {
		t.Fatalf("Expected IndexingPolicy to be set, but it was nil")
	}

	if len(readProperties.IndexingPolicy.FullTextIndexes) != 2 {
		t.Fatalf("Expected 2 full text indexes, got %d", len(readProperties.IndexingPolicy.FullTextIndexes))
	}

	// Validate first full text index
	index1 := readProperties.IndexingPolicy.FullTextIndexes[0]
	if index1.Path != "/title" {
		t.Errorf("Expected first full text index path /title, got %s", index1.Path)
	}

	// Validate second full text index
	index2 := readProperties.IndexingPolicy.FullTextIndexes[1]
	if index2.Path != "/description" {
		t.Errorf("Expected second full text index path /description, got %s", index2.Path)
	}
	// Try to insert some sample data for full-text search testing
	sampleItems := []map[string]interface{}{
		{
			"id":          "1",
			"pk":          "test",
			"title":       "Azure Cosmos DB Full Text Search",
			"description": "Learn about the powerful full-text search capabilities in Azure Cosmos DB",
		},
		{
			// An item that should not match the full-text search query.
			// This means it should not contain the word "search" in the title or description.
			"id":          "2",
			"pk":          "test",
			"title":       "Not related",
			"description": "An unrelated item that should not match the query",
		},
	}

	partitionKey := NewPartitionKeyString("test")
	for _, item := range sampleItems {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			t.Fatalf("Failed to marshal sample item: %v", err)
		}
		_, err = container.CreateItem(context.TODO(), partitionKey, itemBytes, nil)
		if err != nil {
			t.Fatalf("Failed to create sample item: %v", err)
		}
	}

	// Try to execute a full-text search query (this may fail if the SDK doesn't support it yet)
	queryText := `SELECT * FROM c WHERE FullTextContains(c.title, "search") OR FullTextContains(c.description, "search")`
	queryPager := container.NewQueryItemsPager(queryText, partitionKey, nil)

	if !queryPager.More() {
		t.Errorf("Expected results from full-text search query, but got none")
	}

	page, err := queryPager.NextPage(context.TODO())
	if err != nil {
		t.Errorf("Failed to execute full-text search query: %v", err)
	}
	if len(page.Items) != 1 {
		t.Errorf("Expected 1 result from full-text search query, but got %d", len(page.Items))
	}

	var resultItem map[string]interface{}
	err = json.Unmarshal(page.Items[0], &resultItem)
	if err != nil {
		t.Errorf("Failed to unmarshal full-text search result: %v", err)
	} else {
		if resultItem["id"] != "1" {
			t.Errorf("Expected result item ID '1', got '%s'", resultItem["id"])
		}
		if resultItem["title"] != "Azure Cosmos DB Full Text Search" {
			t.Errorf("Expected result item title 'Azure Cosmos DB Full Text Search', got '%s'", resultItem["title"])
		}
	}

	// Clean up
	_, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

func TestEmulatorContainerReadPartitionKeyRanges(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_container aContainer", "read_partition_key_ranges aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerGETPKR")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}

	throughput := NewManualThroughputProperties(30000)

	resp, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if resp.ContainerProperties.ID != properties.ID {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container, _ := database.NewContainer("aContainer")

	item := map[string]interface{}{
		"id": "testitem1",
	}
	itemBytes, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal item: %v", err)
	}
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("testitem1"), itemBytes, nil)
	if err != nil {
		t.Fatalf("Failed to insert item: %v", err)
	}
	time.Sleep(2 * time.Second)

	pkRangesResponse, err := container.getPartitionKeyRanges(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read partition key ranges: %v", err)
	}
	t.Logf("PK Ranges Response: %+v", pkRangesResponse)

	if len(pkRangesResponse.PartitionKeyRanges) == 0 {
		t.Fatalf("Expected at least one partition key range, got none")
	}

	if pkRangesResponse.PartitionKeyRanges[1].ID == "" {
		t.Errorf("Expected partition key range ID to be set, but got empty string")
	}

	if pkRangesResponse.PartitionKeyRanges[1].MinInclusive == "" {
		t.Errorf("Expected partition key range MinInclusive to be set, but got empty string")
	}

	if pkRangesResponse.PartitionKeyRanges[1].MaxExclusive == "" {
		t.Errorf("Expected partition key range MaxExclusive to be set, but got empty string")
	}
}

func TestEmulatorContainerGetFeedRanges(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{"create_container aContainer", "read_partition_key_ranges aContainer"},
	}))

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerGFR")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		ID: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}

	throughput := NewManualThroughputProperties(30000)

	resp, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: &throughput})
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if resp.ContainerProperties.ID != properties.ID {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container, _ := database.NewContainer("aContainer")

	item := map[string]interface{}{
		"id": "testitem1",
	}
	itemBytes, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal item: %v", err)
	}
	_, err = container.CreateItem(context.TODO(), NewPartitionKeyString("testitem1"), itemBytes, nil)
	if err != nil {
		t.Fatalf("Failed to insert item: %v", err)
	}
	time.Sleep(2 * time.Second)

	feedRanges, err := container.GetFeedRanges(context.TODO())
	if err != nil {
		t.Fatalf("Failed to get feed ranges: %v", err)
	}
	t.Logf("Feed Ranges: %+v", feedRanges)

	if len(feedRanges) == 0 {
		t.Fatalf("Expected at least one feed range, got none")
	}

	for i, feedRange := range feedRanges {
		if feedRange.MinInclusive == "" {
			t.Errorf("Feed range %d MinInclusive is empty", i)
		}
		if feedRange.MaxExclusive == "" {
			t.Errorf("Feed range %d MaxExclusive is empty", i)
		}
	}
}

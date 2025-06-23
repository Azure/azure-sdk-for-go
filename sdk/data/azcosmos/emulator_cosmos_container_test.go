// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"
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

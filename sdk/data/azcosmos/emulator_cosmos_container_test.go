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

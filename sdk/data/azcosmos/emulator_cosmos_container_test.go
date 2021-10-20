// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"
)

func TestContainerCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		Id: "aContainer",
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

	resp, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: throughput})
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if resp.ContainerProperties.Id != properties.Id {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container := resp.ContainerProperties.Container
	resp, err = container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	updatedProperties := ContainerProperties{
		Id: "aContainer",
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

	mt, err := throughputResponse.ThroughputProperties.ManualThroughput()
	if err != nil {
		t.Errorf("Failed to read throughput: %v", err)
	}

	if mt != 400 {
		t.Errorf("Unexpected throughput: %v", mt)
	}

	newScale := NewManualThroughputProperties(500)
	_, err = container.ReplaceThroughput(context.TODO(), *newScale, nil)
	if err != nil {
		t.Errorf("Failed to read throughput: %v", err)
	}

	resp, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

func TestContainerAutoscaleCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := ContainerProperties{
		Id: "aContainer",
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

	resp, err := database.CreateContainer(context.TODO(), properties, &CreateContainerOptions{ThroughputProperties: throughput})
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if resp.ContainerProperties.Id != properties.Id {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container := resp.ContainerProperties.Container
	resp, err = container.Read(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	throughputResponse, err := container.ReadThroughput(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read throughput: %v", err)
	}

	maxru, err := throughputResponse.ThroughputProperties.AutoscaleMaxThroughput()
	if err != nil {
		t.Errorf("Failed to read throughput: %v", err)
	}

	if maxru != 5000 {
		t.Errorf("Unexpected throughput: %v", maxru)
	}

	newScale := NewAutoscaleThroughputProperties(10000)
	_, err = container.ReplaceThroughput(context.TODO(), *newScale, nil)
	if err != nil {
		t.Errorf("Failed to read throughput: %v", err)
	}

	resp, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

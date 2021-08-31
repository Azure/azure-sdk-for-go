// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos_emulator_tests

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/cosmos/azcosmos"
)

func TestContainerCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests()
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := azcosmos.CosmosContainerProperties{
		Id: "aContainer",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
		IndexingPolicy: &azcosmos.IndexingPolicy{
			IncludedPaths: []azcosmos.IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []azcosmos.ExcludedPath{
				{Path: "/\"_etag\"/?"},
			},
			Automatic:    true,
			IndexingMode: azcosmos.IndexingModeConsistent,
		},
	}

	throughput := azcosmos.NewManualThroughputProperties(400)

	resp, err := database.CreateContainer(context.TODO(), properties, throughput, nil)
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

	updatedProperties := azcosmos.CosmosContainerProperties{
		Id: "aContainer",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
		IndexingPolicy: &azcosmos.IndexingPolicy{
			IncludedPaths: []azcosmos.IncludedPath{},
			ExcludedPaths: []azcosmos.ExcludedPath{},
			Automatic:     false,
			IndexingMode:  azcosmos.IndexingModeNone,
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

	newScale := azcosmos.NewManualThroughputProperties(500)
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
	emulatorTests := newEmulatorTests()
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := azcosmos.CosmosContainerProperties{
		Id: "aContainer",
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
		IndexingPolicy: &azcosmos.IndexingPolicy{
			IncludedPaths: []azcosmos.IncludedPath{
				{Path: "/*"},
			},
			ExcludedPaths: []azcosmos.ExcludedPath{
				{Path: "/\"_etag\"/?"},
			},
			Automatic:    true,
			IndexingMode: azcosmos.IndexingModeConsistent,
		},
	}

	throughput := azcosmos.NewAutoscaleThroughputProperties(5000)

	resp, err := database.CreateContainer(context.TODO(), properties, throughput, nil)
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

	newScale := azcosmos.NewAutoscaleThroughputProperties(10000)
	_, err = container.ReplaceThroughput(context.TODO(), *newScale, nil)
	if err != nil {
		t.Errorf("Failed to read throughput: %v", err)
	}

	resp, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

// +build emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"
)

func TestContainerCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests()
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "containerCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := CosmosContainerProperties{
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

	resp, err := database.AddContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container2: %v", err)
	}

	if resp.RawResponse.StatusCode != 201 {
		t.Fatal(emulatorTests.parseErrorResponse(resp.RawResponse))
	}

	if resp.ContainerProperties.Id != properties.Id {
		t.Errorf("Unexpected id match: %v", resp.ContainerProperties)
	}

	if resp.ContainerProperties.PartitionKeyDefinition.Paths[0] != properties.PartitionKeyDefinition.Paths[0] {
		t.Errorf("Unexpected path match: %v", resp.ContainerProperties)
	}

	container := resp.ContainerProperties.Container
	resp, err = container.Get(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to read container: %v", err)
	}

	if resp.RawResponse.StatusCode != 200 {
		t.Fatal(emulatorTests.parseErrorResponse(resp.RawResponse))
	}

	updatedProperties := CosmosContainerProperties{
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

	resp, err = container.Update(context.TODO(), updatedProperties, nil)
	if err != nil {
		t.Fatalf("Failed to update container: %v", err)
	}

	if resp.RawResponse.StatusCode != 200 {
		t.Fatal(emulatorTests.parseErrorResponse(resp.RawResponse))
	}

	resp, err = container.Delete(context.TODO(), nil)
	if err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}

	if resp.RawResponse.StatusCode != 204 {
		t.Fatal(emulatorTests.parseErrorResponse(resp.RawResponse))
	}
}

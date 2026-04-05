// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"testing"
)

// TestDirectModeItemCRUD tests direct mode (RNTBD) document operations against the emulator
func TestDirectModeItemCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getDirectModeClient(t, newSpanValidator(t, &spanMatcher{
		ExpectedSpans: []string{},
	}))

	dbName := "directModeTest"

	// Clean up stale database from previous failed test runs
	existingDb, _ := client.NewDatabase(dbName)
	_, _ = existingDb.Delete(context.TODO(), nil)

	database := emulatorTests.createDatabase(t, context.TODO(), client, dbName)
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)

	properties := ContainerProperties{
		ID: "testContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/pk"},
		},
	}

	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	item := map[string]interface{}{
		"id":    "item1",
		"pk":    "partition1",
		"value": "test data",
	}

	container, _ := database.NewContainer("testContainer")
	pk := NewPartitionKeyString("partition1")

	marshalled, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Attempting to create item via Direct Mode (RNTBD)...")

	// This will use RNTBD direct mode for the document operation
	itemResponse, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
	if err != nil {
		t.Fatalf("Failed to create item via direct mode: %v", err)
	}

	t.Logf("Item created successfully! SessionToken: %v", *itemResponse.SessionToken)

	// Read the item back
	itemResponse, err = container.ReadItem(context.TODO(), pk, "item1", nil)
	if err != nil {
		t.Fatalf("Failed to read item via direct mode: %v", err)
	}

	var readItem map[string]interface{}
	err = json.Unmarshal(itemResponse.Value, &readItem)
	if err != nil {
		t.Fatalf("Failed to unmarshal item: %v", err)
	}

	if readItem["id"] != "item1" {
		t.Fatalf("Expected id 'item1', got '%v'", readItem["id"])
	}

	t.Log("Direct mode item read successful!")

	// Delete the item
	_, err = container.DeleteItem(context.TODO(), pk, "item1", nil)
	if err != nil {
		t.Fatalf("Failed to delete item via direct mode: %v", err)
	}

	t.Log("Direct mode CRUD test passed!")
}

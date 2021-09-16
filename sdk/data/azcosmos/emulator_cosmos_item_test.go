// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"testing"
)

func TestItemCRUD(t *testing.T) {
	emulatorTests := newEmulatorTests(t)
	client := emulatorTests.getClient(t)

	database := emulatorTests.createDatabase(t, context.TODO(), client, "itemCRUD")
	defer emulatorTests.deleteDatabase(t, context.TODO(), database)
	properties := CosmosContainerProperties{
		Id: "aContainer",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}

	resp, err := database.CreateContainer(context.TODO(), properties, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	item := map[string]string{
		"id":    "1",
		"value": "2",
	}

	container := resp.ContainerProperties.Container
	pk, err := NewPartitionKey("1")
	if err != nil {
		t.Fatalf("Failed to create pk: %v", err)
	}

	itemResponse, err := container.CreateItem(context.TODO(), pk, item, nil)
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	// No content on write by default
	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}

	itemResponse, err = container.ReadItem(context.TODO(), pk, "1", nil)
	if err != nil {
		t.Fatalf("Failed to read item: %v", err)
	}

	if len(itemResponse.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", itemResponse.Value)
	}

	var itemResponseBody map[string]interface{}
	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}
	if itemResponseBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", itemResponseBody["id"])
	}
	if itemResponseBody["value"] != "2" {
		t.Fatalf("Expected value to be 2, got %v", itemResponseBody["value"])
	}

	item["value"] = "3"
	itemResponse, err = container.ReplaceItem(context.TODO(), pk, "1", item, &CosmosItemRequestOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		t.Fatalf("Failed to replace item: %v", err)
	}

	// Explicitly requesting body on write
	if len(itemResponse.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", itemResponse.Value)
	}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}
	if itemResponseBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", itemResponseBody["id"])
	}
	if itemResponseBody["value"] != "3" {
		t.Fatalf("Expected value to be 3, got %v", itemResponseBody["value"])
	}

	item["value"] = "4"
	itemResponse, err = container.UpsertItem(context.TODO(), pk, item, &CosmosItemRequestOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		t.Fatalf("Failed to upsert item: %v", err)
	}

	// Explicitly requesting body on write
	if len(itemResponse.Value) == 0 {
		t.Fatalf("Expected non-empty response, got %v", itemResponse.Value)
	}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal item response: %v", err)
	}
	if itemResponseBody["id"] != "1" {
		t.Fatalf("Expected id to be 1, got %v", itemResponseBody["id"])
	}
	if itemResponseBody["value"] != "4" {
		t.Fatalf("Expected value to be 4, got %v", itemResponseBody["value"])
	}

	itemResponse, err = container.DeleteItem(context.TODO(), pk, "1", nil)
	if err != nil {
		t.Fatalf("Failed to replace item: %v", err)
	}

	if len(itemResponse.Value) != 0 {
		t.Fatalf("Expected empty response, got %v", itemResponse.Value)
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

// Helper to create a container with simple string id PK and return container + cleanup func
func setupContainerForReadMany(t *testing.T, e *emulatorTests, client *Client, dbName string, containerName string) *ContainerClient {
	database := e.createDatabase(t, context.TODO(), client, dbName)
	// create container
	properties := ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}
	_, err := database.CreateContainer(context.TODO(), properties, nil)
	if err != nil {
		t.Fatalf("failed to create container: %v", err)
	}
	c, _ := database.NewContainer(containerName)
	return c
}

func TestReadMany_NilItemsSlice(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_nilitems_db", "rmnil")
	defer e.deleteDatabase(t, context.TODO(), container.database)

	// Pass nil items slice; should return empty response and no error
	resp, err := container.ReadManyItems(context.TODO(), nil, nil)
	if err != nil {
		t.Fatalf("expected no error for nil items slice, got: %v", err)
	}
	if len(resp.Items) != 0 {
		t.Fatalf("expected zero items, got: %d", len(resp.Items))
	}
}

func TestReadMany_ReadSeveralItems(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_many_db", "rmmany")
	defer e.deleteDatabase(t, context.TODO(), container.database)

	// create 3 items
	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("%d", i)
		item := map[string]string{"id": id, "pk": id}
		marshalled, _ := json.Marshal(item)
		pk := NewPartitionKeyString(item["id"]) // partition is id
		_, err := container.CreateItem(context.TODO(), pk, marshalled, nil)
		if err != nil {
			t.Fatalf("failed to create item %d: %v", i, err)
		}
	}

	// prepare identities
	idents := make([]ItemIdentity, 0, 3)
	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("%d", i)
		idents = append(idents, ItemIdentity{ID: id, PartitionKey: NewPartitionKeyString(id)})
	}

	resp, err := container.ReadManyItems(context.TODO(), idents, nil)
	if err != nil {
		t.Fatalf("ReadManyItems failed: %v", err)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("expected 3 items, got: %d", len(resp.Items))
	}
	if resp.RequestCharge <= 0 {
		t.Fatalf("expected positive request charge, got: %f", resp.RequestCharge)
	}
	// verify items ids are as expected as the items created before
	for i := 0; i < 3; i++ {
		var returnedItem map[string]interface{}
		err := json.Unmarshal(resp.Items[i], &returnedItem)
		if err != nil {
			t.Fatalf("failed to unmarshal returned item %d: %v", i, err)
		}
		expectedID := fmt.Sprintf("%d", i)
		// id in the returned JSON might be a string or a number; stringify for comparison
		idVal := returnedItem["id"]
		gotID := fmt.Sprintf("%v", idVal)
		if gotID != expectedID {
			t.Fatalf("expected item id %s, got: %s", expectedID, gotID)
		}
	}

}

func TestReadMany_NilIDReturnsError(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_nils_db", "rmnils")
	defer e.deleteDatabase(t, context.TODO(), container.database)

	// create one item
	item := map[string]string{"id": "x", "pk": "x"}
	marshalled, _ := json.Marshal(item)
	_, err := container.CreateItem(context.TODO(), NewPartitionKeyString("x"), marshalled, nil)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	// pass an identity with empty id
	idents := []ItemIdentity{{ID: "", PartitionKey: NewPartitionKeyString("x")}}
	_, err = container.ReadManyItems(context.TODO(), idents, nil)
	if err == nil {
		t.Fatalf("expected error for empty id in identity, got nil")
	}
}

// Additional test: partial failure - one identity valid, one invalid -> expect error
func TestReadMany_PartialFailure(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_partial_db", "rmpartial")
	defer e.deleteDatabase(t, context.TODO(), container.database)

	// create a valid item
	item := map[string]string{"id": "good", "pk": "good"}
	marshalled, _ := json.Marshal(item)
	_, err := container.CreateItem(context.TODO(), NewPartitionKeyString("good"), marshalled, nil)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	idents := []ItemIdentity{
		{ID: "good", PartitionKey: NewPartitionKeyString("good")},
		{ID: "missing", PartitionKey: NewPartitionKeyString("missing")},
	}

	_, err = container.ReadManyItems(context.TODO(), idents, nil)
	if err == nil {
		t.Fatalf("expected an error when one of the reads fails, got nil")
	}
}

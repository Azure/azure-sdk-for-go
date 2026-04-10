// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	azcosmosinternal "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal"
	"github.com/stretchr/testify/require"
)

// Helper to create a container with simple string id PK and return container + cleanup func
func setupContainerForReadMany(t *testing.T, e *emulatorTests, client *Client, dbName string, containerName string) *ContainerClient {
	database := e.createDatabase(t, context.Background(), client, dbName)
	// create container
	properties := ContainerProperties{
		ID: containerName,
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths: []string{"/id"},
		},
	}
	_, err := database.CreateContainer(context.Background(), properties, nil)
	require.NoError(t, err, "failed to create container")
	c, _ := database.NewContainer(containerName)
	return c
}

func TestReadMany_NilItemsSlice(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_nilitems_db", "rmnil")
	defer e.deleteDatabase(t, context.Background(), container.database)

	// Pass nil items slice; should return empty response and no error
	resp, err := container.ReadManyItems(context.Background(), nil, nil)
	require.NoError(t, err)
	require.Empty(t, resp.Items)
}

func TestReadMany_ReadSeveralItems(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_many_db", "rmmany")
	defer e.deleteDatabase(t, context.Background(), container.database)

	// create 3 items
	for i := 0; i < 3; i++ {
		item := map[string]string{"id": fmt.Sprintf("%d", i), "pk": fmt.Sprintf("pk_%d", i)}
		marshalled, err := json.Marshal(item)
		require.NoError(t, err)
		pk := NewPartitionKeyString(item["id"]) // partition is id
		_, err = container.CreateItem(context.Background(), pk, marshalled, nil)
		require.NoError(t, err)
	}

	// prepare identities
	idents := make([]ItemIdentity, 0, 3)
	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("%d", i)
		idents = append(idents, ItemIdentity{ID: id, PartitionKey: NewPartitionKeyString(id)})
	}

	resp, err := container.ReadManyItems(context.Background(), idents, nil)
	require.NoError(t, err)
	require.Equal(t, 3, len(resp.Items))
	require.Positive(t, resp.RequestCharge, "expected positive request charge")
	// verify items ids are as expected as the items created before
	for i := 0; i < 3; i++ {
		var returnedItem map[string]interface{}
		err := json.Unmarshal(resp.Items[i], &returnedItem)
		require.NoError(t, err, "failed to unmarshal returned item %d", i)
		expectedID := fmt.Sprintf("%d", i)
		// id in the returned JSON might be a string or a number; stringify for comparison
		idVal := returnedItem["id"]
		gotID := fmt.Sprintf("%v", idVal)
		require.Equal(t, expectedID, gotID)
	}

}

func TestReadMany_NilIDReturnsError(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_nils_db", "rmnils")
	defer e.deleteDatabase(t, context.Background(), container.database)

	// create one item
	item := map[string]string{"id": "x", "pk": "x"}
	marshalled, err := json.Marshal(item)
	require.NoError(t, err)
	_, err = container.CreateItem(context.Background(), NewPartitionKeyString("x"), marshalled, nil)
	require.NoError(t, err)

	// pass an identity with empty id
	idents := []ItemIdentity{{ID: "", PartitionKey: NewPartitionKeyString("x")}}
	_, err = container.ReadManyItems(context.Background(), idents, nil)
	require.Error(t, err, "expected error for empty id in identity")
}

// Additional test: partial failure - one identity valid, one missing -> expect success with only found items returned
func TestReadMany_PartialFailure(t *testing.T) {
	e := newEmulatorTests(t)
	client := e.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	container := setupContainerForReadMany(t, e, client, "readmany_partial_db", "rmpartial")
	defer e.deleteDatabase(t, context.Background(), container.database)

	// create a valid item
	item := map[string]string{"id": "good", "pk": "good"}
	item2 := map[string]string{"id": "good2", "pk": "good2"}
	items := []map[string]string{item, item2}
	for _, item := range items {
		marshalled, err := json.Marshal(item)
		require.NoError(t, err)
		_, err = container.CreateItem(context.Background(), NewPartitionKeyString(item["id"]), marshalled, nil)
		require.NoError(t, err, "failed to create item")
	}

	idents := []ItemIdentity{
		{ID: "good", PartitionKey: NewPartitionKeyString("good")},
		{ID: "missing", PartitionKey: NewPartitionKeyString("missing")},
		{ID: "good2", PartitionKey: NewPartitionKeyString("good2")},
	}

	resp, err := container.ReadManyItems(context.Background(), idents, nil)
	require.NoError(t, err)
	require.Equal(t, 2, len(resp.Items))

	var returnedItem map[string]interface{}
	err = json.Unmarshal(resp.Items[0], &returnedItem)
	require.NoError(t, err, "failed to unmarshal returned item")
	idVal := returnedItem["id"]
	gotID := fmt.Sprintf("%v", idVal)
	require.Equal(t, "good", gotID)

	returnedItem = map[string]interface{}{}
	err = json.Unmarshal(resp.Items[1], &returnedItem)
	require.NoError(t, err, "failed to unmarshal returned item")
	idVal = returnedItem["id"]
	gotID = fmt.Sprintf("%v", idVal)
	require.Equal(t, "good2", gotID)

}

func TestReadMany_WithQueryEngine_EmptyItems(t *testing.T) {
	emulator := newEmulatorTests(t)
	client := emulator.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	db := emulator.createDatabase(t, context.Background(), client, "rm_qeng_empty_db")
	defer emulator.deleteDatabase(t, context.Background(), db)

	container, err := db.NewContainer("c")
	require.NoError(t, err)

	// call ReadMany with empty list and a mock engine
	options := &ReadManyOptions{QueryEngine: azcosmosinternal.NewMockQueryEngine()}
	resp, err := container.ReadManyItems(context.Background(), []ItemIdentity{}, options)
	require.NoError(t, err)
	require.Empty(t, resp.Items)
}

func TestReadMany_WithQueryEngine_ReturnsItems(t *testing.T) {
	emulator := newEmulatorTests(t)
	client := emulator.getClient(t, newSpanValidator(t, &spanMatcher{ExpectedSpans: []string{}}))
	db := emulator.createDatabase(t, context.Background(), client, "rm_qeng_db")
	defer emulator.deleteDatabase(t, context.Background(), db)

	// create container and some items
	_, err := db.CreateContainer(context.Background(), ContainerProperties{ID: "c", PartitionKeyDefinition: PartitionKeyDefinition{
		Paths: []string{"/pk"},
	}}, nil)
	require.NoError(t, err)
	container, err := db.NewContainer("c")
	require.NoError(t, err)

	// insert two items
	for i := 0; i < 2; i++ {
		itm := map[string]string{"id": fmt.Sprintf("%d", i), "pk": fmt.Sprintf("pk_%d", i)}
		b, err := json.Marshal(itm)
		require.NoError(t, err)
		_, err = container.CreateItem(context.Background(), NewPartitionKeyString(itm["pk"]), b, nil)
		require.NoError(t, err)
	}

	// Build item identities to ask for
	idents := []ItemIdentity{{ID: "0", PartitionKey: NewPartitionKeyString("pk_0")}, {ID: "1", PartitionKey: NewPartitionKeyString("pk_1")}}

	// Use the mock query engine which will echo these identities as documents
	options := &ReadManyOptions{QueryEngine: azcosmosinternal.NewMockQueryEngine()}
	resp, err := container.ReadManyItems(context.Background(), idents, options)
	require.NoError(t, err)
	// Expect two items per engine's behavior
	require.Equal(t, 2, len(resp.Items))
}

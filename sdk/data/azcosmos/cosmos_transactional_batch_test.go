// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

func TestTransactionalBatchIDsJSONEscaped(t *testing.T) {
	etag := azcore.ETag(`etag"value`)
	for _, id := range []string{
		"safe", `item"quoted`, `allowed","id":"victim`,
		`allowed","operationType":"Delete","id":"victim`,
		`allowed","ifMatch":"*`, `item\name`, "item\n\t\r", "item\x00", "\u4e16\u754c",
	} {
		for _, options := range []*TransactionalBatchItemOptions{nil, {IfMatchETag: &etag}} {
			batch := &TransactionalBatch{}
			batch.DeleteItem(id, options)
			batch.ReadItem(id, options)
			batch.ReplaceItem(id, []byte(`{"value":1}`), options)
			patch := PatchOperations{}
			patch.AppendAdd("/value", 1)
			batch.PatchItem(id, patch, options)
			for _, op := range batch.operations {
				payload, err := json.Marshal(op)
				require.NoError(t, err)
				var fields map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(payload, &fields))
				var actualID string
				require.NoError(t, json.Unmarshal(fields["id"], &actualID))
				require.Equal(t, id, actualID, "payload: %s", payload)
				require.Equal(t, 1, bytes.Count(payload, []byte(`"id":`)), "payload: %s", payload)
				var operationType string
				require.NoError(t, json.Unmarshal(fields["operationType"], &operationType))
				expectedFields := map[string]any{"id": id}
				switch op.(type) {
				case batchOperationRead:
					require.Equal(t, "Read", operationType)
				case batchOperationDelete:
					require.Equal(t, "Delete", operationType)
				case batchOperationReplace:
					require.Equal(t, "Replace", operationType)
					expectedFields["resourceBody"] = json.RawMessage(`{"value":1}`)
				case batchOperationPatch:
					require.Equal(t, "Patch", operationType)
					expectedFields["resourceBody"] = patch
				}
				expectedFields["operationType"] = operationType
				if options != nil && operationType != "Read" {
					expectedFields["ifMatch"] = etag
				}
				expected, err := json.Marshal(expectedFields)
				require.NoError(t, err)
				require.JSONEq(t, string(expected), string(payload))
			}
		}
	}
}

func TestTransactionalBatchResourceNameValidation(t *testing.T) {
	for _, id := range []string{"", ".", "..", "other/item", `other\item`, "item?", "item#", "item\t", "item\n", "item\r"} {
		body, err := json.Marshal(map[string]string{"id": id})
		require.NoError(t, err)
		batch := &TransactionalBatch{}
		batch.CreateItem(body, nil)
		batch.UpsertItem(body, nil)
		batch.ReplaceItem(id, []byte(`{}`), nil)
		batch.ReplaceItem("safe", body, nil)
		batch.ReadItem(id, nil)
		batch.DeleteItem(id, nil)
		batch.PatchItem(id, PatchOperations{}, nil)
		container := &ContainerClient{database: &DatabaseClient{client: &Client{enableResourceNameValidation: true}}}
		for _, op := range batch.operations {
			invalid := TransactionalBatch{operations: []batchOperation{batchOperationRead{operationType: "Read", id: "safe"}, op}}
			require.ErrorContains(t, invalid.validateResourceNames(&Client{enableResourceNameValidation: true}), "batch operation at index 1: invalid item ID")
			require.NoError(t, invalid.validateResourceNames(&Client{}))
			_, err := container.ExecuteTransactionalBatch(context.Background(), invalid, nil)
			require.ErrorContains(t, err, "batch operation at index 1: invalid item ID")
		}
	}
}

func TestTransactionalBatchInvalidUTF8IDs(t *testing.T) {
	client := &Client{enableResourceNameValidation: true}
	container := &ContainerClient{database: &DatabaseClient{client: client}}
	id := "item" + string([]byte{0xff})
	body := []byte(`{"id":"` + id + `"}`)
	batch := &TransactionalBatch{}
	batch.CreateItem(body, nil)
	batch.UpsertItem(body, nil)
	batch.ReplaceItem("safe", body, nil)
	batch.ReplaceItem(id, []byte(`{}`), nil)
	batch.ReadItem(id, nil)
	batch.DeleteItem(id, nil)
	batch.PatchItem(id, PatchOperations{}, nil)
	for _, op := range batch.operations {
		b := TransactionalBatch{operations: []batchOperation{op}}
		_, err := container.ExecuteTransactionalBatch(context.Background(), b, nil)
		require.ErrorContains(t, err, "invalid item ID")
		require.ErrorContains(t, err, "valid UTF-8")
		require.NoError(t, b.validateResourceNames(&Client{}))
	}
}

func TestTransactionalBatchRejectsMalformedItemBodies(t *testing.T) {
	for _, body := range []string{
		``, `not JSON`, `{"id":"safe"`,
		`{"id":"safe"},"id":"victim"`,
		`{"id":"safe"},"operationType":"Delete","id":"victim"`,
	} {
		batch := &TransactionalBatch{}
		batch.CreateItem([]byte(body), nil)
		batch.UpsertItem([]byte(body), nil)
		batch.ReplaceItem("safe", []byte(body), nil)
		for _, op := range batch.operations {
			_, err := json.Marshal(op)
			require.ErrorContains(t, err, "batch item body must be valid JSON")
		}
	}
}

func TestTransactionalBatchEscapedIDReachesPipeline(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(map[bool]string{false: "default", true: "enabled"}[enabled], func(t *testing.T) {
			id := `allowed","id":"victim`
			if !enabled {
				id = "allowed\\name" + id
			}
			reachedPipeline := errors.New("request reached pipeline")
			var body []byte
			internal, err := azcore.NewClient("azcosmostest", "v1.0.0", runtime.PipelineOptions{
				PerCall: []policy.Policy{policyFunc(func(req *policy.Request) (*http.Response, error) {
					var err error
					body, err = io.ReadAll(req.Raw().Body)
					require.NoError(t, err)
					return nil, reachedPipeline
				})},
			}, nil)
			require.NoError(t, err)
			client := &Client{
				endpoint: "https://localhost", internal: internal, gem: &globalEndpointManager{},
				endpointUrl:                  &url.URL{Scheme: "https", Host: "localhost"},
				enableResourceNameValidation: enabled,
			}
			container, err := client.NewContainer("database", "container")
			require.NoError(t, err)
			batch := container.NewTransactionalBatch(NewPartitionKeyString("pk"))
			batch.DeleteItem(id, nil)
			_, err = container.ExecuteTransactionalBatch(context.Background(), batch, nil)
			require.ErrorIs(t, err, reachedPipeline)
			expected, err := json.Marshal([]map[string]string{{"operationType": "Delete", "id": id}})
			require.NoError(t, err)
			require.JSONEq(t, string(expected), string(body))
			require.Equal(t, 1, bytes.Count(body, []byte(`"id":`)))
		})
	}
}

func TestTransactionalBatchCreateItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	body := map[string]string{
		"foo": "bar",
	}

	itemMarshall, _ := json.Marshal(body)
	batch.CreateItem(itemMarshall, nil)

	if len(batch.operations) != 1 {
		t.Errorf("Expected 1 operation, but got %v", len(batch.operations))
	}

	if batch.operations[0].getOperationType() != operationTypeCreate {
		t.Errorf("Expected operation type %v, but got %v", operationTypeCreate, batch.operations[0].getOperationType())
	}

	asCreate := batch.operations[0].(batchOperationCreate)

	if asCreate.operationType != "Create" {
		t.Errorf("Expected operation type %v, but got %v", "Create", asCreate.operationType)
	}

	if string(asCreate.resourceBody) != string(itemMarshall) {
		t.Errorf("Expected body %v, but got %v", string(itemMarshall), string(asCreate.resourceBody))
	}
}

func TestTransactionalBatchReadItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	itemId := "bar"
	batch.ReadItem(itemId, nil)

	if len(batch.operations) != 1 {
		t.Errorf("Expected 1 operation, but got %v", len(batch.operations))
	}

	if batch.operations[0].getOperationType() != operationTypeRead {
		t.Errorf("Expected operation type %v, but got %v", operationTypeRead, batch.operations[0].getOperationType())
	}

	asRead := batch.operations[0].(batchOperationRead)

	if asRead.operationType != "Read" {
		t.Errorf("Expected operation type %v, but got %v", "Read", asRead.operationType)
	}

	if asRead.id != itemId {
		t.Errorf("Expected id %v, but got %v", itemId, asRead.id)
	}
}

func TestTransactionalBatchUpsertItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	body := map[string]string{
		"foo": "bar",
	}

	itemMarshall, _ := json.Marshal(body)

	options := &TransactionalBatchItemOptions{}
	etag := azcore.ETag("someEtag")
	options.IfMatchETag = &etag
	batch.UpsertItem(itemMarshall, options)

	if len(batch.operations) != 1 {
		t.Errorf("Expected 1 operation, but got %v", len(batch.operations))
	}

	if batch.operations[0].getOperationType() != operationTypeUpsert {
		t.Errorf("Expected operation type %v, but got %v", operationTypeUpsert, batch.operations[0].getOperationType())
	}

	asUpsert := batch.operations[0].(batchOperationUpsert)

	if asUpsert.operationType != "Upsert" {
		t.Errorf("Expected operation type %v, but got %v", "Upsert", asUpsert.operationType)
	}

	if asUpsert.ifMatch != options.IfMatchETag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asUpsert.ifMatch)
	}

	if string(asUpsert.resourceBody) != string(itemMarshall) {
		t.Errorf("Expected body %v, but got %v", string(itemMarshall), string(asUpsert.resourceBody))
	}
}

func TestTransactionalBatchReplaceItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	body := map[string]string{
		"foo": "bar",
	}

	itemMarshall, _ := json.Marshal(body)

	options := &TransactionalBatchItemOptions{}
	etag := azcore.ETag("someEtag")
	options.IfMatchETag = &etag
	itemId := "bar"
	batch.ReplaceItem(itemId, itemMarshall, options)

	if len(batch.operations) != 1 {
		t.Errorf("Expected 1 operation, but got %v", len(batch.operations))
	}

	if batch.operations[0].getOperationType() != operationTypeReplace {
		t.Errorf("Expected operation type %v, but got %v", operationTypeReplace, batch.operations[0].getOperationType())
	}

	asReplace := batch.operations[0].(batchOperationReplace)

	if asReplace.operationType != "Replace" {
		t.Errorf("Expected operation type %v, but got %v", "Replace", asReplace.operationType)
	}

	if asReplace.id != itemId {
		t.Errorf("Expected id %v, but got %v", itemId, asReplace.id)
	}

	if asReplace.ifMatch != options.IfMatchETag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asReplace.ifMatch)
	}

	if string(asReplace.resourceBody) != string(itemMarshall) {
		t.Errorf("Expected body %v, but got %v", string(itemMarshall), string(asReplace.resourceBody))
	}
}

func TestTransactionalBatchDeleteItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	options := &TransactionalBatchItemOptions{}
	etag := azcore.ETag("someEtag")
	options.IfMatchETag = &etag
	itemId := "bar"
	batch.DeleteItem(itemId, options)

	if len(batch.operations) != 1 {
		t.Errorf("Expected 1 operation, but got %v", len(batch.operations))
	}

	if batch.operations[0].getOperationType() != operationTypeDelete {
		t.Errorf("Expected operation type %v, but got %v", operationTypeDelete, batch.operations[0].getOperationType())
	}

	asDelete := batch.operations[0].(batchOperationDelete)

	if asDelete.operationType != "Delete" {
		t.Errorf("Expected operation type %v, but got %v", "Delete", asDelete.operationType)
	}

	if asDelete.ifMatch != options.IfMatchETag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asDelete.ifMatch)
	}

	if asDelete.id != itemId {
		t.Errorf("Expected id %v, but got %v", itemId, asDelete.id)
	}
}

func TestTransactionalBatchPatchItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	patchOperations := PatchOperations{}

	options := &TransactionalBatchItemOptions{}
	etag := azcore.ETag("someEtag")
	options.IfMatchETag = &etag
	itemId := "bar"

	patchOperations.AppendAdd("/foo", "bar")

	batch.PatchItem(itemId, patchOperations, options)

	if len(batch.operations) != 1 {
		t.Errorf("Expected 1 operation, but got %v", len(batch.operations))
	}

	if batch.operations[0].getOperationType() != operationTypePatch {
		t.Errorf("Expected operation type %v, but got %v", operationTypePatch, batch.operations[0].getOperationType())
	}

	asPatch := batch.operations[0].(batchOperationPatch)

	if asPatch.operationType != "Patch" {
		t.Errorf("Expected operation type %v, but got %v", "Patch", asPatch.operationType)
	}

	if asPatch.ifMatch != options.IfMatchETag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asPatch.ifMatch)
	}

	if len(asPatch.patchOperations.operations) != len(patchOperations.operations) {
		t.Errorf("Expected patch operations %v, but got %v", patchOperations, asPatch.patchOperations)
	}
}

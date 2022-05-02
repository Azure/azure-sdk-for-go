// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

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

	pkValue, _ := asCreate.partitionKey.toJsonString()
	batchPkValue, _ := batch.partitionKey.toJsonString()

	if pkValue != batchPkValue {
		t.Errorf("Expected partition key %v, but got %v", pkValue, batchPkValue)
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

	pkValue, _ := asRead.partitionKey.toJsonString()
	batchPkValue, _ := batch.partitionKey.toJsonString()

	if pkValue != batchPkValue {
		t.Errorf("Expected partition key %v, but got %v", pkValue, batchPkValue)
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
	options.IfMatchEtag = &etag
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

	if asUpsert.ifMatch != options.IfMatchEtag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asUpsert.ifMatch)
	}

	if string(asUpsert.resourceBody) != string(itemMarshall) {
		t.Errorf("Expected body %v, but got %v", string(itemMarshall), string(asUpsert.resourceBody))
	}

	pkValue, _ := asUpsert.partitionKey.toJsonString()
	batchPkValue, _ := batch.partitionKey.toJsonString()

	if pkValue != batchPkValue {
		t.Errorf("Expected partition key %v, but got %v", pkValue, batchPkValue)
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
	options.IfMatchEtag = &etag
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

	if asReplace.ifMatch != options.IfMatchEtag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asReplace.ifMatch)
	}

	if string(asReplace.resourceBody) != string(itemMarshall) {
		t.Errorf("Expected body %v, but got %v", string(itemMarshall), string(asReplace.resourceBody))
	}

	pkValue, _ := asReplace.partitionKey.toJsonString()
	batchPkValue, _ := batch.partitionKey.toJsonString()

	if pkValue != batchPkValue {
		t.Errorf("Expected partition key %v, but got %v", pkValue, batchPkValue)
	}
}

func TestTransactionalBatchDeleteItem(t *testing.T) {
	batch := &TransactionalBatch{}
	batch.partitionKey = NewPartitionKeyString("foo")
	options := &TransactionalBatchItemOptions{}
	etag := azcore.ETag("someEtag")
	options.IfMatchEtag = &etag
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

	pkValue, _ := asDelete.partitionKey.toJsonString()
	batchPkValue, _ := batch.partitionKey.toJsonString()

	if pkValue != batchPkValue {
		t.Errorf("Expected partition key %v, but got %v", pkValue, batchPkValue)
	}

	if asDelete.ifMatch != options.IfMatchEtag {
		t.Errorf("Expected ifMatch %v, but got %v", etag, asDelete.ifMatch)
	}

	if asDelete.id != itemId {
		t.Errorf("Expected id %v, but got %v", itemId, asDelete.id)
	}
}

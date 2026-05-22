// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// TransactionalBatch is a batch of operations to be executed in a single transaction.
// See https://docs.microsoft.com/azure/cosmos-db/sql/transactional-batch
type TransactionalBatch struct {
	partitionKey PartitionKey
	operations   []batchOperation
}

// CreateItem adds a create operation to the batch.
func (b *TransactionalBatch) CreateItem(item []byte, o *TransactionalBatchItemOptions) {
	b.operations = append(b.operations,
		batchOperationCreate{
			operationType: "Create",
			resourceBody:  item})
}

// DeleteItem adds a delete operation to the batch.
func (b *TransactionalBatch) DeleteItem(itemID string, o *TransactionalBatchItemOptions) {
	if o == nil {
		o = &TransactionalBatchItemOptions{}
	}
	b.operations = append(b.operations,
		batchOperationDelete{
			operationType: "Delete",
			id:            itemID,
			ifMatch:       o.IfMatchETag})
}

// ReplaceItem adds a replace operation to the batch.
func (b *TransactionalBatch) ReplaceItem(itemID string, item []byte, o *TransactionalBatchItemOptions) {
	if o == nil {
		o = &TransactionalBatchItemOptions{}
	}
	b.operations = append(b.operations,
		batchOperationReplace{
			operationType: "Replace",
			id:            itemID,
			resourceBody:  item,
			ifMatch:       o.IfMatchETag})
}

// UpsertItem adds an upsert operation to the batch.
func (b *TransactionalBatch) UpsertItem(item []byte, o *TransactionalBatchItemOptions) {
	if o == nil {
		o = &TransactionalBatchItemOptions{}
	}
	b.operations = append(b.operations,
		batchOperationUpsert{
			operationType: "Upsert",
			resourceBody:  item,
			ifMatch:       o.IfMatchETag})
}

// ReadItem adds a read operation to the batch.
func (b *TransactionalBatch) ReadItem(itemID string, o *TransactionalBatchItemOptions) {
	b.operations = append(b.operations,
		batchOperationRead{
			operationType: "Read",
			id:            itemID})
}

// PatchItem adds a patch operation to the batch
func (b *TransactionalBatch) PatchItem(itemID string, p PatchOperations, o *TransactionalBatchItemOptions) {
	if o == nil {
		o = &TransactionalBatchItemOptions{}
	}
	b.operations = append(b.operations,
		batchOperationPatch{
			operationType:   "Patch",
			id:              itemID,
			patchOperations: p,
			ifMatch:         o.IfMatchETag,
		})
}

type batchOperation interface {
	getOperationType() operationType
}

type batchOperationCreate struct {
	operationType string
	resourceBody  []byte
}

func (b batchOperationCreate) getOperationType() operationType {
	return operationTypeCreate
}

// MarshalJSON implements the json.Marshaler interface
func (b batchOperationCreate) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	fmt.Fprint(buffer, ",\"resourceBody\":")
	buffer.Write(b.resourceBody)
	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

type batchOperationDelete struct {
	operationType string
	ifMatch       *azcore.ETag
	id            string
}

func (b batchOperationDelete) getOperationType() operationType {
	return operationTypeDelete
}

// MarshalJSON implements the json.Marshaler interface
func (b batchOperationDelete) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	fmt.Fprintf(buffer, ",\"id\":\"%s\"", b.id)
	if b.ifMatch != nil {
		fmt.Fprint(buffer, ",\"ifMatch\":")
		etag, err := json.Marshal(b.ifMatch)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

type batchOperationReplace struct {
	operationType string
	ifMatch       *azcore.ETag
	id            string
	resourceBody  []byte
}

func (b batchOperationReplace) getOperationType() operationType {
	return operationTypeReplace
}

// MarshalJSON implements the json.Marshaler interface
func (b batchOperationReplace) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	if b.ifMatch != nil {
		fmt.Fprint(buffer, ",\"ifMatch\":")
		etag, err := json.Marshal(b.ifMatch)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	fmt.Fprintf(buffer, ",\"id\":\"%s\"", b.id)
	fmt.Fprint(buffer, ",\"resourceBody\":")
	buffer.Write(b.resourceBody)
	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

type batchOperationUpsert struct {
	operationType string
	ifMatch       *azcore.ETag
	resourceBody  []byte
}

func (b batchOperationUpsert) getOperationType() operationType {
	return operationTypeUpsert
}

// MarshalJSON implements the json.Marshaler interface
func (b batchOperationUpsert) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	if b.ifMatch != nil {
		fmt.Fprint(buffer, ",\"ifMatch\":")
		etag, err := json.Marshal(b.ifMatch)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	fmt.Fprint(buffer, ",\"resourceBody\":")
	buffer.Write(b.resourceBody)
	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

type batchOperationPatch struct {
	operationType   string
	id              string
	ifMatch         *azcore.ETag
	patchOperations PatchOperations
}

func (b batchOperationPatch) getOperationType() operationType {
	return operationTypePatch
}

// MarshalJSON implements the json.Marshaler interface
func (b batchOperationPatch) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)

	if b.ifMatch != nil {
		fmt.Fprint(buffer, ",\"ifMatch\":")
		etag, err := json.Marshal(b.ifMatch)
		if err != nil {
			return nil, err
		}
		buffer.Write(etag)
	}

	fmt.Fprintf(buffer, ",\"id\":\"%s\"", b.id)
	fmt.Fprint(buffer, ",\"resourceBody\":")
	p, err := json.Marshal(b.patchOperations)
	if err != nil {
		return nil, err
	}
	buffer.Write(p)
	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

type batchOperationRead struct {
	operationType string
	id            string
}

func (b batchOperationRead) getOperationType() operationType {
	return operationTypeRead
}

// MarshalJSON implements the json.Marshaler interface
func (b batchOperationRead) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	fmt.Fprintf(buffer, ",\"id\":\"%s\"", b.id)
	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

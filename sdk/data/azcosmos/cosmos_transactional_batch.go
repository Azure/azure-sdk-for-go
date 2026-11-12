// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (b TransactionalBatch) validateResourceNames(client *Client) error {
	if !client.enableResourceNameValidation {
		return nil
	}
	for i, operation := range b.operations {
		var err error
		switch op := operation.(type) {
		case batchOperationCreate:
			err = client.validateItemBodyID(op.resourceBody)
		case batchOperationUpsert:
			err = client.validateItemBodyID(op.resourceBody)
		case batchOperationReplace:
			err = client.validateResourceName("item", op.id)
			if err == nil {
				err = client.validateItemBodyID(op.resourceBody)
			}
		case batchOperationRead:
			err = client.validateResourceName("item", op.id)
		case batchOperationDelete:
			err = client.validateResourceName("item", op.id)
		case batchOperationPatch:
			err = client.validateResourceName("item", op.id)
		}
		if err != nil {
			return fmt.Errorf("batch operation at index %d: %w", i, err)
		}
	}
	return nil
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
	if !json.Valid(b.resourceBody) {
		return nil, errors.New("batch item body must be valid JSON")
	}
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
	id, err := json.Marshal(b.id)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	fmt.Fprint(buffer, ",\"id\":")
	buffer.Write(id)
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
	if !json.Valid(b.resourceBody) {
		return nil, errors.New("batch item body must be valid JSON")
	}
	id, err := json.Marshal(b.id)
	if err != nil {
		return nil, err
	}
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

	fmt.Fprint(buffer, ",\"id\":")
	buffer.Write(id)
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
	if !json.Valid(b.resourceBody) {
		return nil, errors.New("batch item body must be valid JSON")
	}
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
	id, err := json.Marshal(b.id)
	if err != nil {
		return nil, err
	}
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

	fmt.Fprint(buffer, ",\"id\":")
	buffer.Write(id)
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
	id, err := json.Marshal(b.id)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBufferString("{")
	fmt.Fprintf(buffer, "\"operationType\":\"%s\"", b.operationType)
	fmt.Fprint(buffer, ",\"id\":")
	buffer.Write(id)
	fmt.Fprint(buffer, "}")
	return buffer.Bytes(), nil
}

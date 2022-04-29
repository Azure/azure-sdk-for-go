// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
)

// TransactionalBatch is a batch of operations to be executed in a single transaction.
// See https://docs.microsoft.com/azure/cosmos-db/sql/transactional-batch
type TransactionalBatch struct {
	partitionKey PartitionKey
	container    *ContainerClient
	operations   []batchOperation
}

// CreateItem adds a create operation to the batch.
func (b *TransactionalBatch) CreateItem(item []byte, o *TransactionalBatchOptions) {
	b.operations = append(b.operations, batchOperationCreate{item, o})
}

// DeleteItem adds a delete operation to the batch.
func (b *TransactionalBatch) DeleteItem(itemId string, o *TransactionalBatchOptions) {
	b.operations = append(b.operations, batchOperationDelete{itemId, o})
}

// DeleteItem adds a delete operation to the batch.
func (b *TransactionalBatch) ReplaceItem(itemId string, item []byte, o *TransactionalBatchOptions) {
	b.operations = append(b.operations, batchOperationReplace{itemId, item, o})
}

// DeleteItem adds a delete operation to the batch.
func (b *TransactionalBatch) UpsertItem(item []byte, o *TransactionalBatchOptions) {
	b.operations = append(b.operations, batchOperationUpsert{item, o})
}

// DeleteItem adds a delete operation to the batch.
func (b *TransactionalBatch) ReadItem(itemId string, o *TransactionalBatchOptions) {
	b.operations = append(b.operations, batchOperationRead{itemId, o})
}

// Execute executes the transactional batch.
func (b *TransactionalBatch) Execute(ctx context.Context) (TransactionalBatchResponse, error) {
	if len(b.operations) == 0 {
		return TransactionalBatchResponse{}, errors.New("no operations in batch")
	}

	return TransactionalBatchResponse{}, nil
}

type batchOperation interface {
	getOperationType() operationType
	getOptions() *TransactionalBatchOptions
}

type batchOperationCreate struct {
	item    []byte
	options *TransactionalBatchOptions
}

func (b batchOperationCreate) getOperationType() operationType {
	return operationTypeCreate
}

func (b batchOperationCreate) getOptions() *TransactionalBatchOptions {
	return b.options
}

type batchOperationDelete struct {
	itemId  string
	options *TransactionalBatchOptions
}

func (b batchOperationDelete) getOperationType() operationType {
	return operationTypeDelete
}

func (b batchOperationDelete) getOptions() *TransactionalBatchOptions {
	return b.options
}

type batchOperationReplace struct {
	itemId  string
	item    []byte
	options *TransactionalBatchOptions
}

func (b batchOperationReplace) getOperationType() operationType {
	return operationTypeReplace
}

func (b batchOperationReplace) getOptions() *TransactionalBatchOptions {
	return b.options
}

type batchOperationUpsert struct {
	item    []byte
	options *TransactionalBatchOptions
}

func (b batchOperationUpsert) getOperationType() operationType {
	return operationTypeUpsert
}

func (b batchOperationUpsert) getOptions() *TransactionalBatchOptions {
	return b.options
}

type batchOperationRead struct {
	itemId  string
	options *TransactionalBatchOptions
}

func (b batchOperationRead) getOperationType() operationType {
	return operationTypeRead
}

func (b batchOperationRead) getOptions() *TransactionalBatchOptions {
	return b.options
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	maxTransactionalBatchOperations = 100
	maxTransactionalBatchBytes      = 2 * 1024 * 1024
)

// TransactionalBatch is an ordered set of operations for one complete logical partition.
//
// Construct it with [NewTransactionalBatch]. Item bodies and conditional options are copied
// when appended. Copies of a batch can be extended independently. A completed builder can be
// executed repeatedly or concurrently, but must not be mutated concurrently with other access.
// Reexecuting a batch is a new transaction, not an idempotent replay.
//
// An append error makes the builder unusable, even when the caller ignores that error.
// The zero value has no partition key and cannot be executed. A batch supports create, read,
// upsert, replace, and delete operations; transactional PATCH is not supported.
type TransactionalBatch struct {
	partitionKey PartitionKey
	operations   []transactionalBatchOperation
	err          error
}

// NewTransactionalBatch constructs an empty batch for partitionKey.
//
// The key must have exactly one component per path in the container's partition key definition,
// including every level of a hierarchical partition key. Prefix and cross-partition batches
// are not supported. Construction performs no network I/O; execution validates completeness.
func NewTransactionalBatch(partitionKey PartitionKey) TransactionalBatch {
	return TransactionalBatch{partitionKey: partitionKey}
}

// TransactionalBatchCreateItemOptions configures [TransactionalBatch.CreateItem].
// There are currently no create-specific settings. A nil options pointer selects defaults.
type TransactionalBatchCreateItemOptions struct{}

// TransactionalBatchReadItemOptions configures [TransactionalBatch.ReadItem].
// A nil options pointer selects defaults.
type TransactionalBatchReadItemOptions struct {
	// IfMatchETag makes the read conditional on the item having this ETag.
	// A mismatch fails the entire transaction.
	IfMatchETag *azcore.ETag
}

// TransactionalBatchUpsertItemOptions configures [TransactionalBatch.UpsertItem].
// A nil options pointer selects defaults.
type TransactionalBatchUpsertItemOptions struct {
	// IfMatchETag requires an existing item to have this ETag. A missing item is created.
	// It must not be set together with IfNoneMatchETag.
	IfMatchETag *azcore.ETag

	// IfNoneMatchETag applies an If-None-Match precondition. Use "*" to require creation.
	// It must not be set together with IfMatchETag.
	IfNoneMatchETag *azcore.ETag
}

// TransactionalBatchReplaceItemOptions configures [TransactionalBatch.ReplaceItem].
// A nil options pointer selects defaults.
type TransactionalBatchReplaceItemOptions struct {
	// IfMatchETag makes the replacement conditional on the item having this ETag.
	IfMatchETag *azcore.ETag
}

// TransactionalBatchDeleteItemOptions configures [TransactionalBatch.DeleteItem].
// A nil options pointer selects defaults.
type TransactionalBatchDeleteItemOptions struct {
	// IfMatchETag makes the deletion conditional on the item having this ETag.
	IfMatchETag *azcore.ETag
}

type transactionalBatchOperation struct {
	OperationType string          `json:"operationType"`
	ID            string          `json:"id,omitempty"`
	ResourceBody  json.RawMessage `json:"resourceBody,omitempty"`
	IfMatch       string          `json:"ifMatch,omitempty"`
	IfNoneMatch   string          `json:"ifNoneMatch,omitempty"`
}

// CreateItem appends a create operation. item is the item's encoded JSON.
//
// The SDK does not inspect application-defined fields, including id and partition key properties.
// All item bodies must belong to the batch's logical partition. options may be nil.
func (b *TransactionalBatch) CreateItem(item []byte, options *TransactionalBatchCreateItemOptions) error {
	return b.appendOperation(transactionalBatchOperation{OperationType: "Create"}, item, true, nil, nil)
}

// ReadItem appends a read operation. itemID must not be empty. options may be nil.
//
// Reads return item bodies even when write content responses are disabled. Conditional
// If-None-Match reads are not exposed until native 304/transaction outcome semantics are verified.
func (b *TransactionalBatch) ReadItem(itemID string, options *TransactionalBatchReadItemOptions) error {
	var ifMatch *azcore.ETag
	if options != nil {
		ifMatch = options.IfMatchETag
	}
	return b.appendOperation(transactionalBatchOperation{OperationType: "Read", ID: itemID}, nil, false, ifMatch, nil)
}

// UpsertItem appends an operation that creates or replaces the item encoded in item.
// The SDK leaves application id and partition key consistency checks to the driver/service.
// options may be nil.
func (b *TransactionalBatch) UpsertItem(item []byte, options *TransactionalBatchUpsertItemOptions) error {
	var ifMatch, ifNoneMatch *azcore.ETag
	if options != nil {
		ifMatch, ifNoneMatch = options.IfMatchETag, options.IfNoneMatchETag
	}
	return b.appendOperation(transactionalBatchOperation{OperationType: "Upsert"}, item, true, ifMatch, ifNoneMatch)
}

// ReplaceItem appends a replacement for itemID. item is the replacement's encoded JSON.
// itemID must match the item's id property; the SDK does not parse the body to compare them.
// options may be nil.
func (b *TransactionalBatch) ReplaceItem(itemID string, item []byte, options *TransactionalBatchReplaceItemOptions) error {
	var ifMatch *azcore.ETag
	if options != nil {
		ifMatch = options.IfMatchETag
	}
	return b.appendOperation(transactionalBatchOperation{OperationType: "Replace", ID: itemID}, item, true, ifMatch, nil)
}

// DeleteItem appends a deletion of itemID. itemID must not be empty. options may be nil.
func (b *TransactionalBatch) DeleteItem(itemID string, options *TransactionalBatchDeleteItemOptions) error {
	var ifMatch *azcore.ETag
	if options != nil {
		ifMatch = options.IfMatchETag
	}
	return b.appendOperation(transactionalBatchOperation{OperationType: "Delete", ID: itemID}, nil, false, ifMatch, nil)
}

func (b *TransactionalBatch) appendOperation(operation transactionalBatchOperation, item []byte, requireItem bool, ifMatch, ifNoneMatch *azcore.ETag) error {
	if b == nil {
		return errors.New("azcosmos: TransactionalBatch must not be nil")
	}
	if b.err != nil {
		return b.err
	}
	fail := func(err error) error {
		b.err = err
		return err
	}
	if len(b.operations) >= maxTransactionalBatchOperations {
		return fail(fmt.Errorf("azcosmos: transactional batch must not exceed %d operations", maxTransactionalBatchOperations))
	}
	if operation.OperationType == "Read" || operation.OperationType == "Replace" || operation.OperationType == "Delete" {
		if operation.ID == "" {
			return fail(errors.New("azcosmos: batch item ID must not be empty"))
		}
		if !utf8.ValidString(operation.ID) || strings.IndexByte(operation.ID, 0) >= 0 {
			return fail(errors.New("azcosmos: batch item ID must be valid UTF-8 without NUL bytes"))
		}
	}
	if requireItem {
		if len(item) == 0 || !utf8.Valid(item) || !json.Valid(item) {
			return fail(errors.New("azcosmos: batch item must contain valid, non-empty UTF-8 JSON"))
		}
		operation.ResourceBody = append(json.RawMessage(nil), item...)
	}
	if ifMatch != nil && ifNoneMatch != nil {
		return fail(errors.New("azcosmos: IfMatchETag and IfNoneMatchETag must not both be set"))
	}
	for _, condition := range []struct {
		name  string
		value *azcore.ETag
	}{
		{"IfMatchETag", ifMatch},
		{"IfNoneMatchETag", ifNoneMatch},
	} {
		if err := validateETag(condition.name, condition.value); err != nil {
			return fail(err)
		}
		if condition.value != nil && !utf8.ValidString(string(*condition.value)) {
			return fail(fmt.Errorf("azcosmos: %s must contain valid UTF-8", condition.name))
		}
	}
	if ifMatch != nil {
		operation.IfMatch = string(*ifMatch)
	}
	if ifNoneMatch != nil {
		operation.IfNoneMatch = string(*ifNoneMatch)
	}
	operations := make([]transactionalBatchOperation, len(b.operations)+1)
	copy(operations, b.operations)
	operations[len(b.operations)] = operation
	b.operations = operations
	return nil
}

func (b TransactionalBatch) marshal() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}
	if err := validateBatchPartitionKey(b.partitionKey); err != nil {
		return nil, err
	}
	if len(b.operations) == 0 || len(b.operations) > maxTransactionalBatchOperations {
		return nil, fmt.Errorf("azcosmos: transactional batch must contain between 1 and %d operations", maxTransactionalBatchOperations)
	}
	body, err := json.Marshal(b.operations)
	if err != nil {
		return nil, &Error{Code: CodeSerializationFailed, Message: "encoding transactional batch", cause: err}
	}
	if len(body) > maxTransactionalBatchBytes {
		return nil, fmt.Errorf("azcosmos: encoded transactional batch must not exceed %d bytes", maxTransactionalBatchBytes)
	}
	return body, nil
}

func validateBatchPartitionKey(partitionKey PartitionKey) error {
	if err := validateItemArguments(partitionKey); err != nil {
		return err
	}
	for _, component := range partitionKey.components {
		switch component.kind {
		case partitionKeyKindString:
			if !utf8.ValidString(component.stringValue) {
				return errors.New("azcosmos: batch partition key strings must contain valid UTF-8")
			}
		case partitionKeyKindNumber:
			if math.IsNaN(component.numberValue) || math.IsInf(component.numberValue, 0) {
				return errors.New("azcosmos: batch partition key numbers must be finite")
			}
		}
	}
	return nil
}

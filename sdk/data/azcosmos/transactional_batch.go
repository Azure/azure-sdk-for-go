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
// Construct it with [NewTransactionalBatch] or [ContainerClient.NewTransactionalBatch].
// Items are serialized and copied when appended, along with conditional options.
// Copies of a batch can be extended independently. A completed builder can be
// executed repeatedly or concurrently, but must not be mutated concurrently with other access.
// Executing a batch again starts a new transaction, not an idempotent replay.
//
// An append error makes the builder unusable, even when the caller ignores that error.
// The zero value has no partition key and cannot be executed. A batch supports create, read,
// upsert, replace, and delete operations; transactional PATCH is not supported.
type TransactionalBatch struct {
	partitionKey PartitionKey
	operations   []transactionalBatchOperation
	err          error
	serializer   Serializer
}

// NewTransactionalBatch constructs an empty batch for partitionKey using encoding/json.
// Use [ContainerClient.NewTransactionalBatch] to use a client's configured serializer.
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
	// IfMatchETag forwards an If-Match read condition to the native driver.
	// Live checks with native driver v0.2.0 returned HTTP 304 for a matching ETag
	// and HTTP 200 for a mismatch. Do not use this condition to guard writes:
	// accompanying writes can still commit.
	// It must not be set together with IfNoneMatchETag.
	IfMatchETag *azcore.ETag

	// IfNoneMatchETag forwards an If-None-Match read condition to the native driver.
	// Live checks with native driver v0.2.0 returned HTTP 200 for both matching
	// and different ETags. Do not rely on this option to produce HTTP 304.
	// It must not be set together with IfMatchETag.
	IfNoneMatchETag *azcore.ETag
}

// TransactionalBatchUpsertItemOptions configures [TransactionalBatch.UpsertItem].
// A nil options pointer selects defaults.
type TransactionalBatchUpsertItemOptions struct {
	// IfMatchETag requires an existing item to have this ETag. A missing item is created.
	// It must not be set together with IfNoneMatchETag.
	IfMatchETag *azcore.ETag

	// IfNoneMatchETag forwards an If-None-Match upsert condition to the native driver.
	// Live checks with native driver v0.2.0 replaced existing items even with "*".
	// This is not a creation-only guard; use [TransactionalBatch.CreateItem] instead.
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

// CreateItem serializes item with the batch's serializer and appends a create operation.
//
// The SDK does not inspect application-defined fields, including id and partition key properties.
// All item bodies must belong to the batch's logical partition. options may be nil.
func (b *TransactionalBatch) CreateItem[T any](item T, options *TransactionalBatchCreateItemOptions) error {
	return b.appendOperation(transactionalBatchOperation{OperationType: "Create"}, item, true, nil, nil)
}

// ReadItem appends a read operation. itemID must not be empty. options may be nil.
//
// With native driver v0.2.0, disabling write content responses can also suppress HTTP 200 read
// bodies. Explicitly enable TransactionalBatchOptions.Operation.EnableContentResponseOnWrite
// when read bodies are required. HTTP 304 has no item body. Statuses do not establish batch commitment.
func (b *TransactionalBatch) ReadItem(itemID string, options *TransactionalBatchReadItemOptions) error {
	var ifMatch, ifNoneMatch *azcore.ETag
	if options != nil {
		ifMatch, ifNoneMatch = options.IfMatchETag, options.IfNoneMatchETag
	}
	return b.appendOperation(transactionalBatchOperation{OperationType: "Read", ID: itemID}, nil, false, ifMatch, ifNoneMatch)
}

// UpsertItem serializes item with the batch's serializer and appends a create-or-replace operation.
// The SDK leaves application id and partition key consistency checks to the driver/service.
// options may be nil.
func (b *TransactionalBatch) UpsertItem[T any](item T, options *TransactionalBatchUpsertItemOptions) error {
	var ifMatch, ifNoneMatch *azcore.ETag
	if options != nil {
		ifMatch, ifNoneMatch = options.IfMatchETag, options.IfNoneMatchETag
	}
	return b.appendOperation(transactionalBatchOperation{OperationType: "Upsert"}, item, true, ifMatch, ifNoneMatch)
}

// ReplaceItem serializes item with the batch's serializer and appends a replacement for itemID.
// itemID must match the item's id property; the SDK does not parse the body to compare them.
// options may be nil.
func (b *TransactionalBatch) ReplaceItem[T any](itemID string, item T, options *TransactionalBatchReplaceItemOptions) error {
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

func (b *TransactionalBatch) appendOperation(operation transactionalBatchOperation, item any, requireItem bool, ifMatch, ifNoneMatch *azcore.ETag) error {
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
		body, err := marshalItem(b.serializer, item)
		if err != nil {
			return fail(err)
		}
		if !utf8.Valid(body) {
			return fail(&Error{Code: CodeSerializationFailed, Message: "encoding batch item: serializer must produce valid UTF-8"})
		}
		operation.ResourceBody = append(json.RawMessage(nil), body...)
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

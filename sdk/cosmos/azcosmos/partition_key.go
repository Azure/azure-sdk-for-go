// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "encoding/json"

// PartitionKey represents a logical partition key value.
type PartitionKey struct {
	isNone               bool
	partitionKeyInternal *partitionKeyInternal
}

// NewPartitionKeyNone creates a partition key value for non-partitioned containers.
func NewPartitionKeyNone() *PartitionKey {
	return &PartitionKey{
		isNone:               true,
		partitionKeyInternal: nonePartitionKey,
	}
}

// NewPartitionKey creates a new partition key.
// value - the partition key value.
func NewPartitionKey(value interface{}) (*PartitionKey, error) {
	pkInternal, err := newPartitionKeyInternal([]interface{}{value})
	if err != nil {
		return nil, err
	}
	return &PartitionKey{
		partitionKeyInternal: pkInternal,
		isNone:               false,
	}, nil
}

func (pk *PartitionKey) toJsonString() (string, error) {
	if pk.isNone {
		return "", nil
	}

	res, err := json.Marshal(pk.partitionKeyInternal)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

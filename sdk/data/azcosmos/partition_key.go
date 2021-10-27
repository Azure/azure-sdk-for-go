// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
)

// PartitionKey represents a logical partition key value.
type PartitionKey struct {
	values []interface{}
}

func NewPartitionKeyString(value string) PartitionKey {
	components := []interface{}{value}
	return PartitionKey{
		values: components,
	}
}

func NewPartitionKeyBool(value bool) PartitionKey {
	components := []interface{}{value}
	return PartitionKey{
		values: components,
	}
}

func NewPartitionKeyNumber(value float64) PartitionKey {
	components := []interface{}{value}
	return PartitionKey{
		values: components,
	}
}

func (pk *PartitionKey) toJsonString() (string, error) {
	res, err := json.Marshal(pk.values)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

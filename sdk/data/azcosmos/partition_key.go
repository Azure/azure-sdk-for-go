// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
)

// PartitionKey represents a logical partition key value.
type PartitionKey struct {
	values []interface{}
}

func (p PartitionKey) MarshalJSON() ([]byte, error) {
	// TODO: support multicomponent partition keys
	switch val := p.values[0].(type) {
	case nil:
		return []byte("null"), nil
	case bool, string, float64:
		return json.Marshal(val)
	default:
		return nil, fmt.Errorf("PartitionKey can only be a string, bool, or a number: '%T'", p.values[0])
	}
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

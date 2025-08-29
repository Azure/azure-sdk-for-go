// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"strconv"
	"strings"
)

// PartitionKey represents a logical partition key value.
type PartitionKey struct {
	values []interface{}
}

// NullPartitionKey represents a partition key with a null value.
var NullPartitionKey PartitionKey = PartitionKey{
	values: []interface{}{nil},
}

// NewPartitionKey creates a new partition key.
func NewPartitionKey() PartitionKey {
	return PartitionKey{
		values: []interface{}{},
	}
}

// NewPartitionKeyString creates a partition key with a string value.
func NewPartitionKeyString(value string) PartitionKey {
	components := []interface{}{value}
	return PartitionKey{
		values: components,
	}
}

// NewPartitionKeyBool creates a partition key with a boolean value.
func NewPartitionKeyBool(value bool) PartitionKey {
	components := []interface{}{value}
	return PartitionKey{
		values: components,
	}
}

// NewPartitionKeyNumber creates a partition key with a numeric value.
func NewPartitionKeyNumber(value float64) PartitionKey {
	components := []interface{}{value}
	return PartitionKey{
		values: components,
	}
}

// AppendString appends a string value to the partition key.
func (pk PartitionKey) AppendString(value string) PartitionKey {
	pk.values = append(pk.values, value)
	return pk
}

// AppendBool appends a boolean value to the partition key.
func (pk PartitionKey) AppendBool(value bool) PartitionKey {
	pk.values = append(pk.values, value)
	return pk
}

// AppendNumber appends a numeric value to the partition key.
func (pk PartitionKey) AppendNumber(value float64) PartitionKey {
	pk.values = append(pk.values, value)
	return pk
}

// AppendNull appends a null value to the partition key.
func (pk PartitionKey) AppendNull() PartitionKey {
	pk.values = append(pk.values, nil)
	return pk
}

func (pk *PartitionKey) toJsonString() (string, error) {
	var completeJson strings.Builder
	completeJson.Grow(256)
	completeJson.WriteString("[")
	for index, i := range pk.values {
		switch v := i.(type) {
		case string:
			// json marshall does not support escaping ASCII as an option
			escaped := strconv.QuoteToASCII(v)
			completeJson.WriteString(escaped)
		default:
			res, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			completeJson.WriteString(string(res))
		}

		if index < len(pk.values)-1 {
			completeJson.WriteString(",")
		}
	}

	completeJson.WriteString("]")
	return completeJson.String(), nil
}

// MarshalJSON implements the json.Marshaler interface for PartitionKey.
func (pk PartitionKey) MarshalJSON() ([]byte, error) {
	if len(pk.values) == 0 {
		return []byte("[]"), nil
	}

	jsonString, err := pk.toJsonString()
	if err != nil {
		return nil, err
	}

	return []byte(jsonString), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for PartitionKey.
func (pk *PartitionKey) UnmarshalJSON(data []byte) error {
	var values []interface{}
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	pk.values = values
	return nil
}

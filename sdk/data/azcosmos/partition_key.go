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

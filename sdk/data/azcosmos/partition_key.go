// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/epk"
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
	fmt.Fprint(&completeJson, "[")
	for index, i := range pk.values {
		switch v := i.(type) {
		case string:
			// json marshall does not support escaping ASCII as an option
			escaped := strconv.QuoteToASCII(v)
			fmt.Fprint(&completeJson, escaped)
		default:
			res, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			fmt.Fprint(&completeJson, string(res))
		}

		if index < len(pk.values)-1 {
			fmt.Fprint(&completeJson, ",")
		}
	}

	fmt.Fprint(&completeJson, "]")
	return completeJson.String(), nil
}

// computeEffectivePartitionKey computes the effective partition key hash for
// this partition key value.
func (pk *PartitionKey) computeEffectivePartitionKey(kind PartitionKeyKind, version int) epk.EffectivePartitionKey {
	values := make([]interface{}, len(pk.values))
	for i, v := range pk.values {
		values[i] = v
	}

	// Empty values → undefined partition key
	if len(values) == 0 {
		values = []interface{}{epk.UndefinedMarker{}}
	}

	var epkStr string
	switch {
	case version == 1:
		epkStr = epk.ComputeV1(values)
	case kind == PartitionKeyKindMultiHash:
		epkStr = epk.ComputeV2MultiHash(values)
	default:
		epkStr = epk.ComputeV2Hash(values)
	}

	return epk.EffectivePartitionKey{EPK: epkStr}
}

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
	copy(values, pk.values)

	// Empty values → undefined partition key
	if len(values) == 0 {
		values = []interface{}{epk.UndefinedMarker{}}
	}

	var epkStr string
	switch {
	case version == 1:
		epkStr = epk.ComputeV1(values)
	case kind == PartitionKeyKindMultiHash:
		epkStr = epk.ComputeV2MultiHashForRouting(values)
	default:
		epkStr = epk.ComputeV2HashForRouting(values)
	}

	return epk.EffectivePartitionKey{EPK: epkStr}
}

// epkRange represents an effective partition key range for routing.
// For a point key, Min == Max (both equal to the EPK).
// For a prefix key on a MultiHash container, Min is the prefix EPK and
// Max is prefix EPK + "FF" (an exclusive upper bound).
type epkRange struct {
	Min string // inclusive
	Max string // exclusive (empty means same as Min, i.e., point)
}

// isRange returns true if this represents a range (prefix key) rather than a point.
func (r epkRange) isRange() bool {
	return r.Max != "" && r.Min != r.Max
}

// computeEPKRange computes the EPK range for a partition key given the container's
// partition key definition. For full keys it returns a point range. For prefix keys
// on MultiHash containers it returns a range [prefix_epk, prefix_epk+"FF").
// Non-MultiHash containers require exactly the right number of components.
func computeEPKRange(pk *PartitionKey, pkDef PartitionKeyDefinition) (epkRange, error) {
	pkVersion := pkDef.Version
	if pkVersion == 0 {
		pkVersion = 1
	}

	// Undefined PK (no components) is a concrete value, not a prefix.
	// It hashes to a single deterministic EPK regardless of the number of
	// definition paths, so always return a point range.
	if len(pk.values) == 0 {
		epkVal := pk.computeEffectivePartitionKey(pkDef.Kind, pkVersion)
		return epkRange{Min: epkVal.EPK, Max: epkVal.EPK}, nil
	}

	componentCount := len(pk.values)
	pathCount := len(pkDef.Paths)

	if componentCount > pathCount {
		return epkRange{}, fmt.Errorf("more partition key components (%d) than definition paths (%d)", componentCount, pathCount)
	}

	if pkDef.Kind != PartitionKeyKindMultiHash && componentCount != pathCount {
		return epkRange{}, fmt.Errorf("non-MultiHash containers require exactly %d components, got %d", pathCount, componentCount)
	}

	epkVal := pk.computeEffectivePartitionKey(pkDef.Kind, pkVersion)

	isPrefix := pkDef.Kind == PartitionKeyKindMultiHash && componentCount < pathCount
	if isPrefix {
		// "FF" is safe as an upper-bound sentinel because maskTopBitsForRouting
		// ensures every EPK component's first hex digit is in [0-3].
		return epkRange{
			Min: epkVal.EPK,
			Max: epkVal.EPK + "FF",
		}, nil
	}

	return epkRange{
		Min: epkVal.EPK,
		Max: epkVal.EPK,
	}, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package epk

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PartitionKeyInternal represents an internal partition key with multiple components.
// This is used for routing and comparison operations.
type PartitionKeyInternal struct {
	components []interface{}
}

// PartitionKeyComponentType represents the type ordering for partition key components.
// The ordering is: Undefined < Null < False < True < Number < String < Infinity
type PartitionKeyComponentType int

const (
	ComponentTypeUndefined PartitionKeyComponentType = iota
	ComponentTypeNull
	ComponentTypeFalse
	ComponentTypeTrue
	ComponentTypeNumber
	ComponentTypeString
	ComponentTypeInfinity
)

// Special partition key values
var (
	// Empty represents an empty partition key (minimum inclusive value).
	Empty = PartitionKeyInternal{components: []interface{}{}}

	// InclusiveMinimum is an alias for Empty.
	InclusiveMinimum = Empty

	// ExclusiveMaximum represents the maximum partition key value (Infinity).
	ExclusiveMaximum = PartitionKeyInternal{components: []interface{}{infinityMarker{}}}
)

// infinityMarker is a sentinel type representing the "Infinity" partition key value.
type infinityMarker struct{}

// FromJsonString parses a JSON string into a PartitionKeyInternal.
// Returns an error if the JSON is invalid or cannot be parsed.
func FromJsonString(jsonStr string) (PartitionKeyInternal, error) {
	if jsonStr == "" {
		return PartitionKeyInternal{}, fmt.Errorf("unable to deserialize partition key value: %s", jsonStr)
	}

	// Special case: "Infinity" represents ExclusiveMaximum
	if jsonStr == `"Infinity"` {
		return ExclusiveMaximum, nil
	}

	var raw interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return PartitionKeyInternal{}, fmt.Errorf("unable to deserialize partition key value: %s", jsonStr)
	}

	// Must be an array
	arr, ok := raw.([]interface{})
	if !ok {
		return PartitionKeyInternal{}, fmt.Errorf("unable to deserialize partition key value: %s", jsonStr)
	}

	components := make([]interface{}, len(arr))
	for i, v := range arr {
		switch val := v.(type) {
		case nil, bool, float64, string:
			components[i] = val
		case map[string]interface{}:
			// Empty object {} represents Undefined
			if len(val) == 0 {
				components[i] = UndefinedMarker{}
			} else {
				return PartitionKeyInternal{}, fmt.Errorf("unable to deserialize partition key value: %s", jsonStr)
			}
		default:
			return PartitionKeyInternal{}, fmt.Errorf("unable to deserialize partition key value: %s", jsonStr)
		}
	}

	return PartitionKeyInternal{components: components}, nil
}

// FromObjectArray creates a PartitionKeyInternal from an array of objects.
// If strict is true, invalid types will cause an error.
// If strict is false, invalid types are converted to Undefined.
func FromObjectArray(objects []interface{}, strict bool) (PartitionKeyInternal, error) {
	components := make([]interface{}, len(objects))
	for i, obj := range objects {
		switch v := obj.(type) {
		case nil, bool, float64, string:
			components[i] = v
		case int:
			components[i] = float64(v)
		case int64:
			components[i] = float64(v)
		case UndefinedMarker:
			components[i] = v
		default:
			if strict {
				return PartitionKeyInternal{}, fmt.Errorf("invalid partition key value type: %T", obj)
			}
			// Non-strict mode: treat unknown types as Undefined
			components[i] = UndefinedMarker{}
		}
	}
	return PartitionKeyInternal{components: components}, nil
}

// ToJson returns the JSON representation of the partition key.
func (pk PartitionKeyInternal) ToJson() string {
	// Special case: ExclusiveMaximum
	if len(pk.components) == 1 {
		if _, ok := pk.components[0].(infinityMarker); ok {
			return `"Infinity"`
		}
	}

	var sb strings.Builder
	sb.WriteByte('[')
	for i, comp := range pk.components {
		if i > 0 {
			sb.WriteByte(',')
		}
		switch v := comp.(type) {
		case nil:
			sb.WriteString("null")
		case bool:
			if v {
				sb.WriteString("true")
			} else {
				sb.WriteString("false")
			}
		case float64:
			// Format number to match Java behavior
			sb.WriteString(formatNumber(v))
		case string:
			data, _ := json.Marshal(v)
			sb.Write(data)
		case UndefinedMarker:
			sb.WriteString("{}")
		case infinityMarker:
			// Should not happen in normal arrays
			sb.WriteString("{}")
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

func formatNumber(f float64) string {
	// Format to match Java: 5 -> "5.0", 5.5 -> "5.5"
	s := fmt.Sprintf("%v", f)
	if !strings.Contains(s, ".") && !strings.Contains(s, "e") && !strings.Contains(s, "E") {
		s += ".0"
	}
	return s
}

// Components returns the partition key components.
func (pk PartitionKeyInternal) Components() []interface{} {
	return pk.components
}

// IsEmpty returns true if the partition key has no components.
func (pk PartitionKeyInternal) IsEmpty() bool {
	return len(pk.components) == 0
}

// getComponentType returns the type of a partition key component for ordering.
func getComponentType(comp interface{}) PartitionKeyComponentType {
	switch v := comp.(type) {
	case UndefinedMarker:
		return ComponentTypeUndefined
	case nil:
		return ComponentTypeNull
	case bool:
		if v {
			return ComponentTypeTrue
		}
		return ComponentTypeFalse
	case float64:
		return ComponentTypeNumber
	case string:
		return ComponentTypeString
	case infinityMarker:
		return ComponentTypeInfinity
	default:
		return ComponentTypeUndefined
	}
}

// CompareTo compares two partition keys.
// Returns -1 if pk < other, 0 if pk == other, 1 if pk > other.
// The type ordering is: Undefined < Null < False < True < Number < String < Infinity
func (pk PartitionKeyInternal) CompareTo(other PartitionKeyInternal) int {
	// Compare component by component
	minLen := len(pk.components)
	if len(other.components) < minLen {
		minLen = len(other.components)
	}

	for i := 0; i < minLen; i++ {
		cmp := compareComponents(pk.components[i], other.components[i])
		if cmp != 0 {
			return cmp
		}
	}

	// If all compared components are equal, shorter key is smaller
	if len(pk.components) < len(other.components) {
		return -1
	}
	if len(pk.components) > len(other.components) {
		return 1
	}
	return 0
}

// compareComponents compares two partition key components.
func compareComponents(a, b interface{}) int {
	typeA := getComponentType(a)
	typeB := getComponentType(b)

	// Different types: compare by type ordering
	if typeA != typeB {
		if typeA < typeB {
			return -1
		}
		return 1
	}

	// Same type: compare values
	switch typeA {
	case ComponentTypeUndefined, ComponentTypeNull, ComponentTypeFalse, ComponentTypeTrue, ComponentTypeInfinity:
		// These types have no value comparison - equal if same type
		return 0
	case ComponentTypeNumber:
		numA := a.(float64)
		numB := b.(float64)
		if numA < numB {
			return -1
		}
		if numA > numB {
			return 1
		}
		return 0
	case ComponentTypeString:
		strA := a.(string)
		strB := b.(string)
		if strA < strB {
			return -1
		}
		if strA > strB {
			return 1
		}
		return 0
	}

	return 0
}

// Contains returns true if this partition key is a prefix of (or equal to) the other partition key.
// An empty partition key contains all other partition keys.
func (pk PartitionKeyInternal) Contains(other PartitionKeyInternal) bool {
	// Empty partition key contains everything
	if len(pk.components) == 0 {
		return true
	}

	// If pk has more components than other, it cannot be a prefix
	if len(pk.components) > len(other.components) {
		return false
	}

	// Check if pk is a prefix of other
	for i := 0; i < len(pk.components); i++ {
		if compareComponents(pk.components[i], other.components[i]) != 0 {
			return false
		}
	}

	return true
}

// Equals returns true if two partition keys are equal.
func (pk PartitionKeyInternal) Equals(other PartitionKeyInternal) bool {
	return pk.CompareTo(other) == 0
}

// Error messages matching Java SDK
var (
	ErrTooFewPartitionKeyComponents  = fmt.Errorf("PartitionKey has fewer components than defined in the PartitionKeyDefinition")
	ErrTooManyPartitionKeyComponents = fmt.Errorf("PartitionKey has more components than defined in the PartitionKeyDefinition")
)

// ValidateComponentCount validates that the partition key has the expected number of components.
// Returns an error if the partition key has too few or too many components.
func (pk PartitionKeyInternal) ValidateComponentCount(expectedCount int) error {
	// Special cases: Empty and ExclusiveMaximum are always valid
	if pk.IsEmpty() {
		return nil
	}
	if len(pk.components) == 1 {
		if _, ok := pk.components[0].(infinityMarker); ok {
			return nil
		}
	}

	if len(pk.components) < expectedCount {
		return ErrTooFewPartitionKeyComponents
	}
	if len(pk.components) > expectedCount {
		return ErrTooManyPartitionKeyComponents
	}
	return nil
}

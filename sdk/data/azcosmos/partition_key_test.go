// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	validTypes := map[string]PartitionKey{
		"[10.5]":            NewPartitionKeyNumber(float64(10.5)),
		"[10]":              NewPartitionKeyNumber(float64(10)),
		"[\"some string\"]": NewPartitionKeyString("some string"),
		"[true]":            NewPartitionKeyBool(true),
		"[false]":           NewPartitionKeyBool(false),
		"[null]":            NullPartitionKey,
	}

	for expectedSerialization, pk := range validTypes {
		if len(pk.values) != 1 {
			t.Errorf("Expected partition key to have 1 component, but it has %v", len(pk.values))
		}

		serialization, err := pk.toJsonString()
		if err != nil {
			t.Errorf("Failed to serialize PK for %v, got %v", pk, err)
		}

		if serialization != expectedSerialization {
			t.Errorf("Expected serialization %v, but got %v", expectedSerialization, serialization)
		}
	}
}

func TestPartitionKeyAppends(t *testing.T) {
	validTypes := map[string]PartitionKey{
		"[\"key0\"]":           NewPartitionKey().AppendString("key0"),
		"[true]":               NewPartitionKey().AppendBool(true),
		"[false]":              NewPartitionKey().AppendBool(false),
		"[10.5]":               NewPartitionKey().AppendNumber(10.5),
		"[10]":                 NewPartitionKey().AppendNumber(10),
		"[null]":               NewPartitionKey().AppendNull(),
		"[\"key0\",true,10.5]": NewPartitionKey().AppendString("key0").AppendBool(true).AppendNumber(10.5),
		"[null,null,null]":     NewPartitionKey().AppendNull().AppendNull().AppendNull(),
	}

	for expectedSerialization, pk := range validTypes {
		if len(pk.values) < 1 {
			t.Errorf("Expected partition key to have at least 1 component, but it has %v", len(pk.values))
		}

		serialization, err := pk.toJsonString()
		if err != nil {
			t.Errorf("Failed to serialize PK for %v, got %v", pk, err)
		}

		if serialization != expectedSerialization {
			t.Errorf("Expected serialization %v, but got %v", expectedSerialization, serialization)
		}
	}
}

func TestPartitionKeyEquality(t *testing.T) {
	pk := NewPartitionKeyNumber(float64(10.5))
	pk2 := NewPartitionKeyNumber(float64(10.5))

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyNumber(float64(50))
	pk2 = NewPartitionKeyNumber(float64(50))

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyBool(true)
	pk2 = NewPartitionKeyBool(true)

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyBool(false)
	pk2 = NewPartitionKeyBool(false)

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyString("some string")
	pk2 = NewPartitionKeyString("some string")

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NullPartitionKey
	pk2 = NullPartitionKey
	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}
}

func TestPartitionKeyMarshalJSON(t *testing.T) {
	// Test simple string partition key
	pk := NewPartitionKeyString("testPartitionKey")
	data, err := json.Marshal(pk)
	require.NoError(t, err, "Failed to marshal partition key")
	require.Equal(t, `["testPartitionKey"]`, string(data), "Unexpected JSON output for string partition key")

	// Test number partition key
	pk = NewPartitionKeyNumber(42)
	data, err = json.Marshal(pk)
	require.NoError(t, err, "Failed to marshal partition key")
	require.Equal(t, `[42]`, string(data), "Unexpected JSON output for number partition key")

	// Test boolean partition key
	pk = NewPartitionKeyBool(true)
	data, err = json.Marshal(pk)
	require.NoError(t, err, "Failed to marshal partition key")
	require.Equal(t, `[true]`, string(data), "Unexpected JSON output for boolean partition key")

	// Test null partition key
	data, err = json.Marshal(NullPartitionKey)
	require.NoError(t, err, "Failed to marshal null partition key")
	require.Equal(t, `[null]`, string(data), "Unexpected JSON output for null partition key")

	// Test empty partition key
	pk = NewPartitionKey()
	data, err = json.Marshal(pk)
	require.NoError(t, err, "Failed to marshal empty partition key")
	require.Equal(t, `[]`, string(data), "Unexpected JSON output for empty partition key")

	// Test complex partition key
	pk = NewPartitionKey().
		AppendString("level1").
		AppendNumber(42).
		AppendBool(true)
	data, err = json.Marshal(pk)
	require.NoError(t, err, "Failed to marshal complex partition key")
	require.Equal(t, `["level1",42,true]`, string(data), "Unexpected JSON output for complex partition key")
}

func TestPartitionKeyUnmarshalJSON(t *testing.T) {
	// Test unmarshaling string partition key
	var pk PartitionKey
	err := json.Unmarshal([]byte(`["testPartitionKey"]`), &pk)
	require.NoError(t, err, "Failed to unmarshal string partition key")
	require.Len(t, pk.values, 1, "Expected 1 value in partition key")
	require.Equal(t, "testPartitionKey", pk.values[0], "Unexpected value in partition key")

	// Test unmarshaling number partition key
	err = json.Unmarshal([]byte(`[42]`), &pk)
	require.NoError(t, err, "Failed to unmarshal number partition key")
	require.Len(t, pk.values, 1, "Expected 1 value in partition key")
	// JSON numbers are unmarshaled as float64
	require.Equal(t, float64(42), pk.values[0], "Unexpected value in partition key")

	// Test unmarshaling boolean partition key
	err = json.Unmarshal([]byte(`[true]`), &pk)
	require.NoError(t, err, "Failed to unmarshal boolean partition key")
	require.Len(t, pk.values, 1, "Expected 1 value in partition key")
	require.Equal(t, true, pk.values[0], "Unexpected value in partition key")

	// Test unmarshaling null partition key
	err = json.Unmarshal([]byte(`[null]`), &pk)
	require.NoError(t, err, "Failed to unmarshal null partition key")
	require.Len(t, pk.values, 1, "Expected 1 value in partition key")
	require.Nil(t, pk.values[0], "Expected nil value in partition key")

	// Test unmarshaling empty partition key
	err = json.Unmarshal([]byte(`[]`), &pk)
	require.NoError(t, err, "Failed to unmarshal empty partition key")
	require.Len(t, pk.values, 0, "Expected 0 values in partition key")

	// Test unmarshaling complex partition key
	err = json.Unmarshal([]byte(`["level1",42,true]`), &pk)
	require.NoError(t, err, "Failed to unmarshal complex partition key")
	require.Len(t, pk.values, 3, "Expected 3 values in partition key")
	require.Equal(t, "level1", pk.values[0], "Unexpected first value in partition key")
	require.Equal(t, float64(42), pk.values[1], "Unexpected second value in partition key")
	require.Equal(t, true, pk.values[2], "Unexpected third value in partition key")
}

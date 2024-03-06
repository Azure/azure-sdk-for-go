// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"reflect"
	"testing"
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

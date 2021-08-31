// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"reflect"
	"testing"
)

func TestInvalidPartitionKeyValues(t *testing.T) {
	invalidTypes := []interface{}{
		complex64(0),
		complex128(0),
		[]byte(nil),
		[]byte{},
		// whatever other type of struct
		cosmosOperationContext{},
	}

	for _, invalidType := range invalidTypes {
		_, err := NewPartitionKey(invalidType)
		if err == nil {
			t.Errorf("Expected error for partition key type %v", invalidType)
		}
	}
}

func TestValidPartitionKeyValues(t *testing.T) {
	validTypes := map[interface{}]string{
		nil:           "[null]",
		true:          "[true]",
		false:         "[false]",
		"some string": "[\"some string\"]",
		int(10):       "[10]",
		int8(10):      "[10]",
		int16(10):     "[10]",
		int32(10):     "[10]",
		int64(10):     "[10]",
		uint(10):      "[10]",
		uint8(10):     "[10]",
		uint16(10):    "[10]",
		uint32(10):    "[10]",
		uint64(10):    "[10]",
		float32(10.5): "[10.5]",
		float64(10.5): "[10.5]",
	}

	for validType, expectedSerialization := range validTypes {
		pk, err := NewPartitionKey(validType)
		if err != nil {
			t.Errorf("Expected success for partition key type %v and got %v", validType, err)
		}

		if len(pk.partitionKeyInternal.components) != 1 {
			t.Errorf("Expected partition key to have 1 component, but it has %v", len(pk.partitionKeyInternal.components))
		}

		serialization, err := pk.toJsonString()
		if err != nil {
			t.Errorf("Failed to serialize PK for %v, got %v", validType, err)
		}

		if serialization != expectedSerialization {
			t.Errorf("Expected serialization %v, but got %v", expectedSerialization, serialization)
		}
	}
}

func TestPartitionKeyEmpty(t *testing.T) {
	pk := &PartitionKey{
		partitionKeyInternal: emptyPartitionKey,
	}

	serialization, err := pk.toJsonString()
	if err != nil {
		t.Errorf("Failed to serialize PK, %v", err)
	}

	if serialization != "[]" {
		t.Errorf("Expected serialization [], but got %v", serialization)
	}
}

func TestPartitionKeyUndefined(t *testing.T) {
	pk := &PartitionKey{
		partitionKeyInternal: undefinedPartitionKey,
	}

	serialization, err := pk.toJsonString()
	if err != nil {
		t.Errorf("Failed to serialize PK, %v", err)
	}

	if serialization != "[{}]" {
		t.Errorf("Expected serialization [{}], but got %v", serialization)
	}
}

func TestPartitionKeyEquality(t *testing.T) {
	validTypes := []interface{}{
		nil,
		true,
		false,
		"some string",
		int(10),
		int8(10),
		int16(10),
		int32(10),
		int64(10),
		uint(10),
		uint8(10),
		uint16(10),
		uint32(10),
		uint64(10),
		float32(10.5),
		float64(10.5),
	}

	for _, validType := range validTypes {
		pk, err := NewPartitionKey(validType)
		if err != nil {
			t.Errorf("Expected success for partition key type %v and got %v", validType, err)
		}

		pk2, _ := NewPartitionKey(validType)
		if !reflect.DeepEqual(pk, pk2) {
			t.Errorf("Expected %v to equal %v", pk, pk2)
		}
	}
}

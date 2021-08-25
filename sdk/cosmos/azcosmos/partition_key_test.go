// +build !emulator
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
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
	validTypes := []interface{}{
		nil,
		true,
		false,
		"some string",
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		float32(0),
		float64(0),
	}

	for _, validType := range validTypes {
		pk, err := NewPartitionKey(validType)
		if err != nil {
			t.Errorf("Expected success for partition key type %v and got %v", validType, err)
		}

		if len(pk.partitionKeyInternal.components) != 1 {
			t.Errorf("Expected partition key to have 1 component, but it has %v", len(pk.partitionKeyInternal.components))
		}
	}
}

package azcosmos

import (
	"testing"
)

func TestPartitionKeyBuilder(t *testing.T) {
	validTypes := map[string]PartitionKey{
		"[\"key0\"]":           NewPartitionKeyBuilder().AppendString("key0").Build(),
		"[true]":               NewPartitionKeyBuilder().AppendBool(true).Build(),
		"[false]":              NewPartitionKeyBuilder().AppendBool(false).Build(),
		"[10.5]":               NewPartitionKeyBuilder().AppendNumber(10.5).Build(),
		"[10]":                 NewPartitionKeyBuilder().AppendNumber(10).Build(),
		"[null]":               NewPartitionKeyBuilder().AppendNull().Build(),
		"[\"key0\",true,10.5]": NewPartitionKeyBuilder().AppendString("key0").AppendBool(true).AppendNumber(10.5).Build(),
		"[null,null,null]":     NewPartitionKeyBuilder().AppendNull().AppendNull().AppendNull().Build(),
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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestPartitionKeyDefinitionSerialization(t *testing.T) {
	pkd_kind_unset_len_one := PartitionKeyDefinition{
		Paths:   []string{"somePath"},
		Version: 2,
	}

	jsonString, err := pkd_kind_unset_len_one.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"kind":"Hash","paths":["somePath"],"version":2}`
	if string(jsonString) != expected {
		t.Errorf("Expected serialization %v, but got %v", expected, string(jsonString))
	}

	pkd_kind_unset_len_two := PartitionKeyDefinition{
		Paths:   []string{"somePath", "someOtherPath"},
		Version: 2,
	}

	jsonString, err = pkd_kind_unset_len_two.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected = `{"kind":"MultiHash","paths":["somePath","someOtherPath"],"version":2}`
	if string(jsonString) != expected {
		t.Errorf("Expected serialization %v, but got %v", expected, string(jsonString))
	}

	pkd_kind_set := PartitionKeyDefinition{
		Kind:    PartitionKeyKindMultiHash,
		Paths:   []string{"somePath"},
		Version: 2,
	}

	jsonString, err = pkd_kind_set.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	expected = `{"kind":"MultiHash","paths":["somePath"],"version":2}`
	if string(jsonString) != expected {
		t.Errorf("Expected serialization %v, but got %v", expected, string(jsonString))
	}
}

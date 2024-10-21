// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"reflect"
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

func TestPartitionKeyDefinitionDeserialization(t *testing.T) {
	def := PartitionKeyDefinition{
		Paths:   []string{"somePath"},
		Version: 2,
	}

	jsonString, err := def.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var otherDef PartitionKeyDefinition
	err = otherDef.UnmarshalJSON(jsonString)

	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	// Kind is inferred based on the number of paths
	if otherDef.Kind != PartitionKeyKindHash {
		t.Errorf("Expected Kind to be %v, but got %v", def.Kind, otherDef.Kind)
	}

	if !reflect.DeepEqual(def.Paths, otherDef.Paths) {
		t.Errorf("Expected Paths to be %v, but got %v", def.Paths, otherDef.Paths)
	}

	if def.Version != otherDef.Version {
		t.Errorf("Expected Version to be %v, but got %v", def.Version, otherDef.Version)
	}

	def = PartitionKeyDefinition{
		Paths:   []string{"somePath", "someOtherPath"},
		Version: 2,
	}

	jsonString, err = def.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	err = otherDef.UnmarshalJSON(jsonString)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	// Kind is inferred based on the number of paths
	if otherDef.Kind != PartitionKeyKindMultiHash {
		t.Errorf("Expected Kind to be %v, but got %v", PartitionKeyKindMultiHash, otherDef.Kind)
	}

	if !reflect.DeepEqual(def.Paths, otherDef.Paths) {
		t.Errorf("Expected Paths to be %v, but got %v", def.Paths, otherDef.Paths)
	}

	if def.Version != otherDef.Version {
		t.Errorf("Expected Version to be %v, but got %v", def.Version, otherDef.Version)
	}

	def = PartitionKeyDefinition{
		Kind:    PartitionKeyKindMultiHash,
		Paths:   []string{"somePath"},
		Version: 2,
	}

	jsonString, err = def.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	err = otherDef.UnmarshalJSON(jsonString)
	if err != nil {
		t.Fatal(err, string(jsonString))
	}

	if def.Kind != otherDef.Kind {
		t.Errorf("Expected Kind to be %v, but got %v", def.Kind, otherDef.Kind)
	}

	if !reflect.DeepEqual(def.Paths, otherDef.Paths) {
		t.Errorf("Expected Paths to be %v, but got %v", def.Paths, otherDef.Paths)
	}

	if def.Version != otherDef.Version {
		t.Errorf("Expected Version to be %v, but got %v", def.Version, otherDef.Version)
	}
}

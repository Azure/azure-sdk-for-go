// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"testing"
)

func TestPatchSetCondition(t *testing.T) {
	patch := PatchOperations{}
	query := "from c where c.taskNum = 3"
	patch.SetCondition(query)

	if patch.condition == nil {
		t.Fatalf("Expected condition to be set")
	}

	if *patch.condition != query {
		t.Fatalf("Expected condition to be %v, but got %v", query, *patch.condition)
	}
}

func TestPatchAppendAdd(t *testing.T) {
	patch := PatchOperations{}
	patch.AppendAdd("/foo", "bar")

	if len(patch.operations) != 1 {
		t.Fatalf("Expected 1 operation, but got %v", len(patch.operations))
	}

	if patch.operations[0].Op != patchOperationTypeAdd {
		t.Fatalf("Expected operation type %v, but got %v", patchOperationTypeAdd, patch.operations[0].Op)
	}

	if patch.operations[0].Path != "/foo" {
		t.Fatalf("Expected path %v, but got %v", "/foo", patch.operations[0].Path)
	}

	if patch.operations[0].Value != "bar" {
		t.Fatalf("Expected value %v, but got %v", "bar", patch.operations[0].Value)
	}

	jsonString, err := json.Marshal(patch)
	if err != nil {
		t.Fatal(err)
	}

	expectedSerialization := `{"operations":[{"op":"add","path":"/foo","value":"bar"}]}`

	if string(jsonString) != expectedSerialization {
		t.Fatalf("Expected serialization %v, but got %v", expectedSerialization, string(jsonString))
	}
}

func TestPatchAppendReplace(t *testing.T) {
	patch := PatchOperations{}
	patch.AppendReplace("/foo", "bar")

	if len(patch.operations) != 1 {
		t.Fatalf("Expected 1 operation, but got %v", len(patch.operations))
	}

	if patch.operations[0].Op != patchOperationTypeReplace {
		t.Fatalf("Expected operation type %v, but got %v", patchOperationTypeReplace, patch.operations[0].Op)
	}

	if patch.operations[0].Path != "/foo" {
		t.Fatalf("Expected path %v, but got %v", "/foo", patch.operations[0].Path)
	}

	if patch.operations[0].Value != "bar" {
		t.Fatalf("Expected value %v, but got %v", "bar", patch.operations[0].Value)
	}

	jsonString, err := json.Marshal(patch)
	if err != nil {
		t.Fatal(err)
	}

	expectedSerialization := `{"operations":[{"op":"replace","path":"/foo","value":"bar"}]}`

	if string(jsonString) != expectedSerialization {
		t.Fatalf("Expected serialization %v, but got %v", expectedSerialization, string(jsonString))
	}
}

func TestPatchAppendRemove(t *testing.T) {
	patch := PatchOperations{}
	patch.AppendRemove("/foo")

	if len(patch.operations) != 1 {
		t.Fatalf("Expected 1 operation, but got %v", len(patch.operations))
	}

	if patch.operations[0].Op != patchOperationTypeRemove {
		t.Fatalf("Expected operation type %v, but got %v", patchOperationTypeRemove, patch.operations[0].Op)
	}

	if patch.operations[0].Path != "/foo" {
		t.Fatalf("Expected path %v, but got %v", "/foo", patch.operations[0].Path)
	}

	if patch.operations[0].Value != nil {
		t.Fatalf("Expected value to be nil, but got %v", patch.operations[0].Value)
	}

	jsonString, err := json.Marshal(patch)
	if err != nil {
		t.Fatal(err)
	}

	expectedSerialization := `{"operations":[{"op":"remove","path":"/foo"}]}`

	if string(jsonString) != expectedSerialization {
		t.Fatalf("Expected serialization %v, but got %v", expectedSerialization, string(jsonString))
	}
}

func TestPatchAppendIncrement(t *testing.T) {
	patch := PatchOperations{}
	value := int64(5)
	patch.AppendIncrement("/foo", value)

	if len(patch.operations) != 1 {
		t.Fatalf("Expected 1 operation, but got %v", len(patch.operations))
	}

	if patch.operations[0].Op != patchOperationTypeIncrement {
		t.Fatalf("Expected operation type %v, but got %v", patchOperationTypeIncrement, patch.operations[0].Op)
	}

	if patch.operations[0].Path != "/foo" {
		t.Fatalf("Expected path %v, but got %v", "/foo", patch.operations[0].Path)
	}

	if patch.operations[0].Value != value {
		t.Fatalf("Expected value to be %v, but got %v", value, patch.operations[0].Value)
	}

	jsonString, err := json.Marshal(patch)
	if err != nil {
		t.Fatal(err)
	}

	expectedSerialization := `{"operations":[{"op":"incr","path":"/foo","value":5}]}`

	if string(jsonString) != expectedSerialization {
		t.Fatalf("Expected serialization %v, but got %v", expectedSerialization, string(jsonString))
	}
}

func TestPatchAppendSet(t *testing.T) {
	patch := PatchOperations{}
	patch.AppendSet("/foo", "bar")

	if len(patch.operations) != 1 {
		t.Fatalf("Expected 1 operation, but got %v", len(patch.operations))
	}

	if patch.operations[0].Op != patchOperationTypeSet {
		t.Fatalf("Expected operation type %v, but got %v", patchOperationTypeSet, patch.operations[0].Op)
	}

	if patch.operations[0].Path != "/foo" {
		t.Fatalf("Expected path %v, but got %v", "/foo", patch.operations[0].Path)
	}

	if patch.operations[0].Value != "bar" {
		t.Fatalf("Expected value to be bar, but got %v", patch.operations[0].Value)
	}

	jsonString, err := json.Marshal(patch)
	if err != nil {
		t.Fatal(err)
	}

	expectedSerialization := `{"operations":[{"op":"set","path":"/foo","value":"bar"}]}`

	if string(jsonString) != expectedSerialization {
		t.Fatalf("Expected serialization %v, but got %v", expectedSerialization, string(jsonString))
	}
}

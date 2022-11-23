// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// PatchOperationType defines supported values for operation types in Patch Document.
type patchOperationType string

const (
	// Represents a patch operation Add.
	patchOperationTypeAdd patchOperationType = "add"
	// Represents a patch operation Replace.
	patchOperationTypeReplace patchOperationType = "replace"
	// Represents a patch operation Remove.
	patchOperationTypeRemove patchOperationType = "remove"
	// Represents a patch operation Set.
	patchOperationTypeSet patchOperationType = "set"
	// Represents a patch operation Increment.
	patchOperationTypeIncrement patchOperationType = "incr"
)

// PatchOperation represents individual patch operation.
type patchOperation struct {
	Op    patchOperationType `json:"op"`
	Path  string             `json:"path"`
	Value interface{}        `json:"value,omitempty"`
}

// PatchOperations represents the patch request.
// See https://learn.microsoft.com/azure/cosmos-db/partial-document-update
type PatchOperations struct {
	condition  *string          `json:"condition,omitempty"`
	operations []patchOperation `json:"operations"`
}

// SetCondition sets condition for the patch request.
func (p PatchOperations) SetCondition(condition string) {
	p.condition = &condition
}

// AppendReplace appends a replace operation to the patch request.
func (p PatchOperations) AppendReplace(path string, value interface{}) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeReplace,
		Path:  path,
		Value: value,
	})
}

// AppendAdd appends an add operation to the patch request.
func (p PatchOperations) AppendAdd(path string, value interface{}) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeAdd,
		Path:  path,
		Value: value,
	})
}

// AppendSet appends a set operation to the patch request.
func (p PatchOperations) AppendSet(path string, value interface{}) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeSet,
		Path:  path,
		Value: value,
	})
}

// AppendRemove appends a remove operation to the patch request.
func (p PatchOperations) AppendRemove(path string) {
	p.operations = append(p.operations, patchOperation{
		Op:   patchOperationTypeRemove,
		Path: path,
	})
}

// AppendIncrement appends an increment operation to the patch request.
func (p PatchOperations) AppendIncrement(path string, value int32) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeIncrement,
		Path:  path,
		Value: value,
	})
}
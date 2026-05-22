// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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
	Value any                `json:"value,omitempty"`
}

// PatchOperations represents the patch request.
// See https://learn.microsoft.com/azure/cosmos-db/partial-document-update
type PatchOperations struct {
	condition  *string
	operations []patchOperation
}

// MarshalJSON implements the json.Marshaler interface
func (o PatchOperations) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	if o.condition != nil {
		fmt.Fprintf(buffer, "\"condition\":\"%s\",", *o.condition)
	}
	fmt.Fprint(buffer, "\"operations\":[")
	for i, operation := range o.operations {
		if i > 0 {
			fmt.Fprint(buffer, ",")
		}
		operationBytes, err := json.Marshal(operation)
		if err != nil {
			return nil, err
		}
		buffer.Write(operationBytes)
	}
	fmt.Fprint(buffer, "]}")
	return buffer.Bytes(), nil
}

// SetCondition sets condition for the patch request.
func (p *PatchOperations) SetCondition(condition string) {
	p.condition = &condition
}

// AppendReplace appends a replace operation to the patch request.
func (p *PatchOperations) AppendReplace(path string, value any) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeReplace,
		Path:  path,
		Value: value,
	})
}

// AppendAdd appends an add operation to the patch request.
func (p *PatchOperations) AppendAdd(path string, value any) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeAdd,
		Path:  path,
		Value: value,
	})
}

// AppendSet appends a set operation to the patch request.
func (p *PatchOperations) AppendSet(path string, value any) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeSet,
		Path:  path,
		Value: value,
	})
}

// AppendRemove appends a remove operation to the patch request.
func (p *PatchOperations) AppendRemove(path string) {
	p.operations = append(p.operations, patchOperation{
		Op:   patchOperationTypeRemove,
		Path: path,
	})
}

// AppendIncrement appends an increment operation to the patch request.
func (p *PatchOperations) AppendIncrement(path string, value int64) {
	p.operations = append(p.operations, patchOperation{
		Op:    patchOperationTypeIncrement,
		Path:  path,
		Value: value,
	})
}

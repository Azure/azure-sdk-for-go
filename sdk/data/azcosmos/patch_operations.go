// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore azsdk

package azcosmos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type patchOperation struct {
	Op    string          `json:"op"`
	Path  string          `json:"path"`
	From  string          `json:"from,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`
}

// PatchOperations is an ordered set of changes for [ContainerClient.PatchItem].
//
// Its zero value is ready to use. Values are encoded and copied when appended, so later changes
// to caller-owned maps, slices, or byte buffers do not change the patch. Copying PatchOperations
// also produces an independent builder: appending to either copy does not change the other.
//
// Patch execution uses the driver's automatic strategy. Operations that are not intrinsically
// retry-safe can use a client-side read-modify-write path, which permanently adds the internal
// _azsdkPatchTracking property to the item. A stable tracking ID across separate PatchItem calls
// is not exposed yet, so deduplication applies only within one call.
type PatchOperations struct {
	operations []patchOperation
	err        error
}

// AppendAdd appends an operation that adds a value at path.
func (p *PatchOperations) AppendAdd(path string, value any) error {
	return p.appendValue("add", path, value, false)
}

// AppendSet appends an operation that sets a value at path, adding it when it does not exist.
func (p *PatchOperations) AppendSet(path string, value any) error {
	return p.appendValue("set", path, value, false)
}

// AppendReplace appends an operation that replaces the existing value at path.
func (p *PatchOperations) AppendReplace(path string, value any) error {
	return p.appendValue("replace", path, value, false)
}

// AppendRemove appends an operation that removes the value at path.
func (p *PatchOperations) AppendRemove(path string) error {
	if err := p.prepareAppend("remove", path); err != nil {
		return err
	}
	p.appendOwned(patchOperation{Op: "remove", Path: path})
	return nil
}

// AppendIncrement appends an operation that increments the number at path by value.
//
// value must encode as one JSON number. NaN, infinities, strings, objects, arrays, booleans, and
// null are rejected.
func (p *PatchOperations) AppendIncrement(path string, value any) error {
	return p.appendValue("incr", path, value, true)
}

// AppendMove appends an operation that moves a value from fromPath to path.
func (p *PatchOperations) AppendMove(fromPath string, path string) error {
	if err := p.prepareAppend("move", path); err != nil {
		return err
	}
	if err := validatePatchPath(fromPath); err != nil {
		return p.setError(fmt.Errorf("azcosmos: patch move source path: %w", err))
	}
	p.appendOwned(patchOperation{Op: "move", Path: path, From: fromPath})
	return nil
}

func (p *PatchOperations) appendValue(op string, path string, value any, numberOnly bool) error {
	if err := p.prepareAppend(op, path); err != nil {
		return err
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return p.setError(fmt.Errorf("azcosmos: encoding patch %s value: %w", op, err))
	}
	if numberOnly && !isJSONNumber(encoded) {
		return p.setError(errors.New("azcosmos: patch increment value must be a JSON number"))
	}

	p.appendOwned(patchOperation{
		Op:    op,
		Path:  path,
		Value: append(json.RawMessage(nil), encoded...),
	})
	return nil
}

func (p *PatchOperations) prepareAppend(op string, path string) error {
	if p == nil {
		return errors.New("azcosmos: PatchOperations must not be nil")
	}
	if p.err != nil {
		return p.err
	}
	if err := validatePatchPath(path); err != nil {
		return p.setError(fmt.Errorf("azcosmos: patch %s path: %w", op, err))
	}
	return nil
}

func (p *PatchOperations) appendOwned(operation patchOperation) {
	operations := make([]patchOperation, len(p.operations)+1)
	copy(operations, p.operations)
	operations[len(p.operations)] = operation
	p.operations = operations
}

func (p *PatchOperations) setError(err error) error {
	if p != nil && p.err == nil {
		p.err = err
	}
	return err
}

func (p PatchOperations) marshal() ([]byte, error) {
	if p.err != nil {
		return nil, p.err
	}
	if len(p.operations) == 0 {
		return nil, errors.New("azcosmos: patch must contain at least one operation")
	}
	return json.Marshal(struct {
		Operations []patchOperation `json:"operations"`
	}{
		Operations: p.operations,
	})
}

func validatePatchPath(path string) error {
	if path == "" || path[0] != '/' {
		return errors.New("path must be a non-empty RFC 6901 JSON pointer beginning with '/'")
	}
	for i := 0; i < len(path); i++ {
		if path[i] != '~' {
			continue
		}
		if i+1 >= len(path) || (path[i+1] != '0' && path[i+1] != '1') {
			return errors.New("path contains an invalid RFC 6901 escape; '~' must be followed by '0' or '1'")
		}
		i++
	}
	return nil
}

func isJSONNumber(encoded []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return false
	}
	if _, ok := value.(json.Number); !ok {
		return false
	}
	return decoder.Decode(&struct{}{}) == io.EOF
}

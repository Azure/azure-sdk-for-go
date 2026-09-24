// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatchOperationsMarshalEveryOperation(t *testing.T) {
	var operations PatchOperations
	require.NoError(t, operations.AppendAdd("/added", map[string]any{"x": 1}))
	require.NoError(t, operations.AppendSet("/set", []string{"a", "b"}))
	require.NoError(t, operations.AppendReplace("/replaced", true))
	require.NoError(t, operations.AppendRemove("/removed"))
	require.NoError(t, operations.AppendIncrement("/count", json.Number("2.5")))
	require.NoError(t, operations.AppendMove("/source", "/destination"))

	body, err := operations.marshal()
	require.NoError(t, err)
	require.JSONEq(t, `{
		"operations": [
			{"op":"add","path":"/added","value":{"x":1}},
			{"op":"set","path":"/set","value":["a","b"]},
			{"op":"replace","path":"/replaced","value":true},
			{"op":"remove","path":"/removed"},
			{"op":"incr","path":"/count","value":2.5},
			{"op":"move","path":"/destination","from":"/source"}
		]
	}`, string(body))
}

func TestPatchOperationsOwnValuesAtAppendTime(t *testing.T) {
	value := map[string]any{
		"slice": []string{"original"},
	}
	raw := json.RawMessage(`{"state":"original"}`)

	var operations PatchOperations
	require.NoError(t, operations.AppendAdd("/map", value))
	require.NoError(t, operations.AppendSet("/raw", raw))

	value["slice"].([]string)[0] = "changed"
	value["added"] = true
	copy(raw, `{"state":"changed!"}`)

	body, err := operations.marshal()
	require.NoError(t, err)
	require.JSONEq(t, `{
		"operations": [
			{"op":"add","path":"/map","value":{"slice":["original"]}},
			{"op":"set","path":"/raw","value":{"state":"original"}}
		]
	}`, string(body))
}

func TestPatchOperationsCopiesDoNotAliasAppends(t *testing.T) {
	var original PatchOperations
	require.NoError(t, original.AppendSet("/first", 1))

	copied := original
	require.NoError(t, original.AppendSet("/original", 2))
	require.NoError(t, copied.AppendSet("/copied", 3))

	originalBody, err := original.marshal()
	require.NoError(t, err)
	copiedBody, err := copied.marshal()
	require.NoError(t, err)

	require.JSONEq(t, `{"operations":[
		{"op":"set","path":"/first","value":1},
		{"op":"set","path":"/original","value":2}
	]}`, string(originalBody))
	require.JSONEq(t, `{"operations":[
		{"op":"set","path":"/first","value":1},
		{"op":"set","path":"/copied","value":3}
	]}`, string(copiedBody))
}

func TestPatchOperationsValidateRFC6901Paths(t *testing.T) {
	for _, path := range []string{"/", "/property", "/a~0b~1c", "/array/0", "/snowman-\u2603"} {
		t.Run("valid "+path, func(t *testing.T) {
			var operations PatchOperations
			require.NoError(t, operations.AppendRemove(path))
		})
	}

	for _, path := range []string{"", "property", "~0", "/dangling~", "/bad~2escape"} {
		t.Run("invalid "+path, func(t *testing.T) {
			var operations PatchOperations
			require.ErrorContains(t, operations.AppendRemove(path), "RFC 6901")
		})
	}
}

func TestPatchOperationsValidateBothMovePaths(t *testing.T) {
	var invalidSource PatchOperations
	require.ErrorContains(t, invalidSource.AppendMove("source", "/destination"), "source path")

	var invalidDestination PatchOperations
	require.ErrorContains(t, invalidDestination.AppendMove("/source", "destination"), "move path")
}

func TestPatchOperationsRejectInvalidUTF8Paths(t *testing.T) {
	invalidPath := string([]byte{'/', 0xff})

	var destination PatchOperations
	require.ErrorContains(t, destination.AppendSet(invalidPath, 1), "valid UTF-8")

	var source PatchOperations
	err := source.AppendMove(invalidPath, "/destination")
	require.ErrorContains(t, err, "source path")
	require.ErrorContains(t, err, "valid UTF-8")
}

func TestPatchOperationsIncrementRequiresJSONNumber(t *testing.T) {
	for _, value := range []any{int64(-2), uint64(3), float64(4.5), json.Number("6e2")} {
		t.Run("valid", func(t *testing.T) {
			var operations PatchOperations
			require.NoError(t, operations.AppendIncrement("/count", value))
		})
	}

	for _, value := range []any{nil, "1", true, []int{1}, map[string]int{"n": 1}, math.NaN(), math.Inf(1)} {
		t.Run("invalid", func(t *testing.T) {
			var operations PatchOperations
			require.Error(t, operations.AppendIncrement("/count", value))
		})
	}
}

func TestPatchOperationsIncrementPreservesNumericCategory(t *testing.T) {
	type namedFloat64 float64
	wholeFloat := 5.0

	for _, tt := range []struct {
		name  string
		value any
		want  string
	}{
		{"integer", int64(1), "1"},
		{"whole float32", float32(1), "1.0"},
		{"whole float64", float64(2), "2.0"},
		{"named whole float", namedFloat64(3), "3.0"},
		{"fractional float", float64(1.5), "1.5"},
		{"scientific float", float64(1e21), "1e+21"},
		{"explicit JSON float", json.Number("4.0"), "4.0"},
		{"explicit JSON integer", json.Number("4"), "4"},
		{"whole float pointer", &wholeFloat, "5.0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var operations PatchOperations
			require.NoError(t, operations.AppendIncrement("/amount", tt.value))
			require.Equal(t, tt.want, string(operations.operations[0].Value))
		})
	}
}

func TestPatchOperationsRetainFirstError(t *testing.T) {
	var operations PatchOperations
	first := operations.AppendRemove("not-a-pointer")
	require.Error(t, first)

	second := operations.AppendSet("/valid", 1)
	require.ErrorIs(t, second, first)

	_, err := operations.marshal()
	require.ErrorIs(t, err, first)
}

func TestPatchOperationsRejectValueThatCannotEncode(t *testing.T) {
	for _, value := range []any{make(chan int), failingJSONMarshaler{}} {
		var operations PatchOperations
		err := operations.AppendAdd("/invalid", value)
		require.ErrorContains(t, err, "encoding patch add value")

		_, marshalErr := operations.marshal()
		require.ErrorIs(t, marshalErr, err)
	}
}

func TestPatchOperationsRequireAtLeastOneOperation(t *testing.T) {
	_, err := (PatchOperations{}).marshal()
	require.ErrorContains(t, err, "at least one operation")
}

func TestPatchOperationsDoNotImposeServerOperationLimit(t *testing.T) {
	var operations PatchOperations
	for i := range 11 {
		require.NoError(t, operations.AppendSet("/value", i))
	}

	body, err := operations.marshal()
	require.NoError(t, err)

	var document struct {
		Operations []json.RawMessage `json:"operations"`
	}
	require.NoError(t, json.Unmarshal(body, &document))
	require.Len(t, document.Operations, 11)
}

func TestPatchOperationsNilReceiverReturnsError(t *testing.T) {
	var operations *PatchOperations
	require.Error(t, operations.AppendSet("/value", 1))
	require.Error(t, operations.AppendRemove("/value"))
	require.Error(t, operations.AppendMove("/source", "/destination"))
}

type failingJSONMarshaler struct{}

func (failingJSONMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cannot marshal")
}

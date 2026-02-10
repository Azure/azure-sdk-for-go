// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialization(t *testing.T) {
	validTypes := map[string]PartitionKey{
		"[10.5]":            NewPartitionKeyNumber(float64(10.5)),
		"[10]":              NewPartitionKeyNumber(float64(10)),
		"[\"some string\"]": NewPartitionKeyString("some string"),
		"[true]":            NewPartitionKeyBool(true),
		"[false]":           NewPartitionKeyBool(false),
		"[null]":            NullPartitionKey,
	}

	for expectedSerialization, pk := range validTypes {
		if len(pk.values) != 1 {
			t.Errorf("Expected partition key to have 1 component, but it has %v", len(pk.values))
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

func TestPartitionKeyAppends(t *testing.T) {
	validTypes := map[string]PartitionKey{
		"[\"key0\"]":           NewPartitionKey().AppendString("key0"),
		"[true]":               NewPartitionKey().AppendBool(true),
		"[false]":              NewPartitionKey().AppendBool(false),
		"[10.5]":               NewPartitionKey().AppendNumber(10.5),
		"[10]":                 NewPartitionKey().AppendNumber(10),
		"[null]":               NewPartitionKey().AppendNull(),
		"[\"key0\",true,10.5]": NewPartitionKey().AppendString("key0").AppendBool(true).AppendNumber(10.5),
		"[null,null,null]":     NewPartitionKey().AppendNull().AppendNull().AppendNull(),
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

func TestPartitionKeyEquality(t *testing.T) {
	pk := NewPartitionKeyNumber(float64(10.5))
	pk2 := NewPartitionKeyNumber(float64(10.5))

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyNumber(float64(50))
	pk2 = NewPartitionKeyNumber(float64(50))

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyBool(true)
	pk2 = NewPartitionKeyBool(true)

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyBool(false)
	pk2 = NewPartitionKeyBool(false)

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NewPartitionKeyString("some string")
	pk2 = NewPartitionKeyString("some string")

	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}

	pk = NullPartitionKey
	pk2 = NullPartitionKey
	if !reflect.DeepEqual(pk, pk2) {
		t.Errorf("Expected %v to equal %v", pk, pk2)
	}
}

// XML structures for parsing the baseline test data.
type epkBaselineResults struct {
	Results []epkBaselineResult `xml:"Result"`
}

type epkBaselineResult struct {
	Input  epkBaselineInput  `xml:"Input"`
	Output epkBaselineOutput `xml:"Output"`
}

type epkBaselineInput struct {
	Description       string `xml:"Description"`
	PartitionKeyValue string `xml:"PartitionKeyValue"`
}

type epkBaselineOutput struct {
	PartitionKeyHashV1 string `xml:"PartitionKeyHashV1"`
	PartitionKeyHashV2 string `xml:"PartitionKeyHashV2"`
}

// parsePartitionKeyValue converts the XML PartitionKeyValue string into a
// PartitionKey. The input is either the special string "UNDEFINED" or a valid
// JSON value (null, bool, number, string, or array of strings for hierarchical PKs).
// Post-unmarshal, the magic strings "NaN", "-Infinity", and "Infinity" are
// converted to their float64 equivalents.
func parsePartitionKeyValue(raw string) PartitionKey {
	if raw == "UNDEFINED" {
		return NewPartitionKey()
	}

	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		panic("failed to parse partition key value: " + raw + ": " + err.Error())
	}

	switch val := v.(type) {
	case nil:
		return NullPartitionKey
	case bool:
		return NewPartitionKeyBool(val)
	case float64:
		return NewPartitionKeyNumber(val)
	case string:
		// Convert magic number strings to actual float64 values
		switch val {
		case "NaN":
			return NewPartitionKeyNumber(math.NaN())
		case "-Infinity":
			return NewPartitionKeyNumber(math.Inf(-1))
		case "Infinity":
			return NewPartitionKeyNumber(math.Inf(1))
		default:
			return NewPartitionKeyString(val)
		}
	case []interface{}:
		pk := NewPartitionKey()
		for _, elem := range val {
			s, ok := elem.(string)
			if !ok {
				panic("non-string element in list partition key value: " + raw)
			}
			pk = pk.AppendString(s)
		}
		return pk
	default:
		panic("unexpected JSON type in partition key value: " + raw)
	}
}

func partitionKeyKindForCategory(category string) PartitionKeyKind {
	if category == "Lists" {
		return PartitionKeyKindMultiHash
	}
	return PartitionKeyKindHash
}

func loadBaselineResults(t *testing.T, filename string) epkBaselineResults {
	t.Helper()
	data, err := os.ReadFile(filename)
	require.NoError(t, err, "failed to read baseline file %s", filename)
	var results epkBaselineResults
	require.NoError(t, xml.Unmarshal(data, &results), "failed to parse baseline file %s", filename)
	return results
}

// TestComputeEffectivePartitionKey_Baseline enumerates every case in the XML
// baseline files and verifies that computeEffectivePartitionKey produces the
// expected V1 and V2 hash values.
func TestComputeEffectivePartitionKey_Baseline(t *testing.T) {
	files := map[string]string{
		"Singletons": "testdata/PartitionKeyHashBaselineTest.Singletons.xml",
		"Numbers":    "testdata/PartitionKeyHashBaselineTest.Numbers.xml",
		"Strings":    "testdata/PartitionKeyHashBaselineTest.Strings.xml",
		"Lists":      "testdata/PartitionKeyHashBaselineTest.Lists.xml",
	}

	for category, filename := range files {
		baseline := loadBaselineResults(t, filename)
		kind := partitionKeyKindForCategory(category)

		for _, tc := range baseline.Results {
			pk := parsePartitionKeyValue(tc.Input.PartitionKeyValue)

			// Test V1 hash
			t.Run(category+"/V1/"+tc.Input.Description, func(t *testing.T) {
				epk := pk.computeEffectivePartitionKey(kind, 1)
				require.Equal(t, tc.Output.PartitionKeyHashV1, epk.epk,
					"V1 hash mismatch for %s (value: %s)", tc.Input.Description, tc.Input.PartitionKeyValue)
			})

			// Test V2 hash
			t.Run(category+"/V2/"+tc.Input.Description, func(t *testing.T) {
				epk := pk.computeEffectivePartitionKey(kind, 2)
				require.Equal(t, tc.Output.PartitionKeyHashV2, epk.epk,
					"V2 hash mismatch for %s (value: %s)", tc.Input.Description, tc.Input.PartitionKeyValue)
			})
		}
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
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
// PartitionKey. The format varies by test file category:
//   - Singletons: "UNDEFINED", "null", "true", "false"
//   - Numbers: numeric literals (e.g. "0", "-1", "5E-324")
//   - Strings: JSON-quoted strings (e.g. `"asdf"`)
//   - Lists: JSON arrays of strings (e.g. `["/path1","/path2"]`)
func parsePartitionKeyValue(raw string, category string) PartitionKey {
	switch category {
	case "Singletons":
		switch raw {
		case "UNDEFINED":
			return NewPartitionKey()
		case "null":
			return NullPartitionKey
		case "true":
			return NewPartitionKeyBool(true)
		case "false":
			return NewPartitionKeyBool(false)
		default:
			panic("unknown singleton value: " + raw)
		}

	case "Numbers":
		// Some "number" test cases are actually quoted strings (NaN, Infinity)
		if strings.HasPrefix(raw, "\"") {
			var s string
			if err := json.Unmarshal([]byte(raw), &s); err != nil {
				panic("failed to parse quoted number value: " + raw)
			}
			return NewPartitionKeyString(s)
		}
		if raw == "-0" {
			return NewPartitionKeyNumber(math.Copysign(0, -1))
		}
		n, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			panic("failed to parse number: " + raw + ": " + err.Error())
		}
		return NewPartitionKeyNumber(n)

	case "Strings":
		var s string
		if err := json.Unmarshal([]byte(raw), &s); err != nil {
			panic("failed to parse string value: " + raw)
		}
		return NewPartitionKeyString(s)

	case "Lists":
		var paths []string
		if err := json.Unmarshal([]byte(raw), &paths); err != nil {
			panic("failed to parse list value: " + raw)
		}
		pk := NewPartitionKey()
		for _, p := range paths {
			pk = pk.AppendString(p)
		}
		return pk

	default:
		panic("unknown category: " + category)
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
			pk := parsePartitionKeyValue(tc.Input.PartitionKeyValue, category)

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

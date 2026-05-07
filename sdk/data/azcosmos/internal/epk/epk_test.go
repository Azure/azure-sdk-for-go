// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package epk

import (
	"encoding/json"
	"encoding/xml"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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

// parseValues converts the XML PartitionKeyValue string into a slice of
// interface{} values suitable for passing to ComputeV1/ComputeV2Hash/etc.
// The input is either the special string "UNDEFINED", a valid JSON scalar, or
// a JSON array of strings (for hierarchical/multi-hash partition keys).
func parseValues(raw string) []interface{} {
	if raw == "UNDEFINED" {
		return []interface{}{UndefinedMarker{}}
	}

	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		panic("failed to parse partition key value: " + raw + ": " + err.Error())
	}

	switch val := v.(type) {
	case nil:
		return []interface{}{nil}
	case bool:
		return []interface{}{val}
	case float64:
		return []interface{}{val}
	case string:
		switch val {
		case "NaN":
			// Use .NET's NaN bit pattern (0xFFF8000000000000) so the hash
			// matches the baseline. Go's math.NaN() uses a different bit
			// pattern (0x7FF8000000000001).
			return []interface{}{math.Float64frombits(0xFFF8000000000000)}
		case "-Infinity":
			return []interface{}{math.Inf(-1)}
		case "Infinity":
			return []interface{}{math.Inf(1)}
		default:
			return []interface{}{val}
		}
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, elem := range val {
			s, ok := elem.(string)
			if !ok {
				panic("non-string element in list partition key value: " + raw)
			}
			result[i] = s
		}
		return result
	default:
		panic("unexpected JSON type in partition key value: " + raw)
	}
}

func loadBaselineResults(t *testing.T, filename string) epkBaselineResults {
	t.Helper()
	data, err := os.ReadFile(filename)
	require.NoError(t, err, "failed to read baseline file %s", filename)
	var results epkBaselineResults
	require.NoError(t, xml.Unmarshal(data, &results), "failed to parse baseline file %s", filename)
	return results
}

// TestComputeEPK_Baseline enumerates every case in the XML baseline files and
// verifies that ComputeV1 and ComputeV2Hash/ComputeV2MultiHash produce the
// expected hash values.
func TestComputeEPK_Baseline(t *testing.T) {
	type fileSpec struct {
		filename  string
		multiHash bool
	}

	files := map[string]fileSpec{
		"Singletons": {filename: "testdata/PartitionKeyHashBaselineTest.Singletons.xml"},
		"Numbers":    {filename: "testdata/PartitionKeyHashBaselineTest.Numbers.xml"},
		"Strings":    {filename: "testdata/PartitionKeyHashBaselineTest.Strings.xml"},
		"Lists":      {filename: "testdata/PartitionKeyHashBaselineTest.Lists.xml", multiHash: true},
	}

	for category, spec := range files {
		baseline := loadBaselineResults(t, spec.filename)

		for _, tc := range baseline.Results {
			values := parseValues(tc.Input.PartitionKeyValue)
			expectedV1 := tc.Output.PartitionKeyHashV1
			expectedV2 := tc.Output.PartitionKeyHashV2

			t.Run(category+"/V1/"+tc.Input.Description, func(t *testing.T) {
				actual := ComputeV1(values)
				require.Equal(t, expectedV1, actual,
					"V1 hash mismatch for %s (value: %s)", tc.Input.Description, tc.Input.PartitionKeyValue)
			})

			t.Run(category+"/V2/"+tc.Input.Description, func(t *testing.T) {
				var actual string
				if spec.multiHash {
					actual = computeV2MultiHash(values)
				} else {
					actual = computeV2Hash(values)
				}
				require.Equal(t, expectedV2, actual,
					"V2 hash mismatch for %s (value: %s)", tc.Input.Description, tc.Input.PartitionKeyValue)
			})
		}
	}
}

func TestCompareEPK_EqualSameLength(t *testing.T) {
	require.Equal(t, 0, CompareEPK("06AB34CFE4E48223", "06AB34CFE4E48223"))
}

func TestCompareEPK_LessThan(t *testing.T) {
	require.Equal(t, -1, CompareEPK("06AB34CFE4E48223", "06AB34CFE4E48224"))
}

func TestCompareEPK_GreaterThan(t *testing.T) {
	require.Equal(t, 1, CompareEPK("06AB34CFE4E48224", "06AB34CFE4E48223"))
}

func TestCompareEPK_ZeroPaddedTailEqual(t *testing.T) {
	// A 32-char partial EPK should equal its 64-char zero-padded equivalent
	partial := "06AB34CFE4E482236BCACBBF50E234AB"
	full := "06AB34CFE4E482236BCACBBF50E234AB00000000000000000000000000000000"
	require.Equal(t, 0, CompareEPK(partial, full))
	require.Equal(t, 0, CompareEPK(full, partial))
}

func TestCompareEPK_NonZeroTailNotEqual(t *testing.T) {
	partial := "06AB34CFE4E482236BCACBBF50E234AB"
	full := "06AB34CFE4E482236BCACBBF50E234AB00000000000000000000000000000001"
	require.Equal(t, -1, CompareEPK(partial, full))
	require.Equal(t, 1, CompareEPK(full, partial))
}

func TestCompareEPK_EmptyStrings(t *testing.T) {
	require.Equal(t, 0, CompareEPK("", ""))
	require.Equal(t, 0, CompareEPK("", "00000"))
	require.Equal(t, -1, CompareEPK("", "00001"))
}

func TestCompareEPK_FFSentinel(t *testing.T) {
	// "FF" should be greater than any masked EPK (first hex digit in [0-3])
	require.Equal(t, 1, CompareEPK("FF", "3FFFFFFFFFFFFFFF"))
	require.Equal(t, -1, CompareEPK("3FFFFFFFFFFFFFFF", "FF"))
}

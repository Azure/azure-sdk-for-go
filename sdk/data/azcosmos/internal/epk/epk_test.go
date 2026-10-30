// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package epk

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"strings"
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

func TestComputeV1Routing(t *testing.T) {
	// Routing vectors from the Java and .NET SDKs, not raw Murmur32 hash baselines.
	// https://github.com/Azure/azure-sdk-for-rust/blob/e259b26317ebe1b4e3171c89869f595e459455ac/sdk/cosmos/azure_data_cosmos_driver/testdata/EffectivePartitionKeyConformance.json
	tests := []struct {
		name     string
		values   []interface{}
		expected string
	}{
		{"empty string", []interface{}{""}, "05C1CF33970FF80800"},
		{"string", []interface{}{"partitionKey"}, "05C1E1B3D9CD2608716273756A756A706F4C667A00"},
		{"null", []interface{}{nil}, "05C1ED45D7475601"},
		{"undefined", []interface{}{UndefinedMarker{}}, "05C1D529E345DC00"},
		{"true", []interface{}{true}, "05C1D7C5A903D803"},
		{"false", []interface{}{false}, "05C1DB857D857C02"},
		{"negative number", []interface{}{float64(-128)}, "05C1D73349F54C053FA0"},
		{"positive number", []interface{}{float64(127)}, "05C1DD539DDFCC05C05FE0"},
		{"smallest positive number", []interface{}{math.SmallestNonzeroFloat64}, "05C1E5C91F4D3005800101010101010102"},
		{"largest number", []interface{}{math.MaxFloat64}, "05C1CBE367C53005FFEFFFFFFFFFFFFFFE"},
		{"large negative number", []interface{}{float64(math.MinInt64)}, "05C1DB35F33D1C053C20"},
		{"large positive number", []interface{}{float64(math.MaxInt64)}, "05C1B799AB2DD005C3E0"},
		{
			"99 characters",
			[]interface{}{strings.Repeat("abcdefghijklmnopqrstuvwxyz", 4)[:99]},
			"05C19B2DC38FC00862636465666768696A6B6C6D6E6F707172737475767778797A7B62636465666768696A6B6C6D6E6F707172737475767778797A7B62636465666768696A6B6C6D6E6F707172737475767778797A7B62636465666768696A6B6C6D6E6F7071727374757600",
		},
		{"100 characters", []interface{}{strings.Repeat("a", 100)}, "05C1EB5921F70608" + strings.Repeat("62", 100) + "00"},
		{"101 characters", []interface{}{strings.Repeat("a", 101)}, "05C1EB5921F70608" + strings.Repeat("62", 100) + "00"},
		{"1024 characters", []interface{}{strings.Repeat("a", 1024)}, "05C1EB5921F70608" + strings.Repeat("62", 100) + "00"},
		{
			"combined components",
			[]interface{}{float64(2), true, false, nil, UndefinedMarker{}, "\u041f\u0440\u0438\u0432\u0435\u0442!"},
			"05C1D19581B37C05C0000302010008D1A0D281D1B9D1B3D1B6D2832200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, ComputeV1(tt.values))
		})
	}
}

func TestComputeV1StringBoundaries(t *testing.T) {
	// Service-compatible UTF-16 truncation and 101-byte binary encoding vectors.
	// https://github.com/Azure/azure-sdk-for-rust/blob/e259b26317ebe1b4e3171c89869f595e459455ac/sdk/cosmos/azure_data_cosmos_driver/src/models/effective_partition_key.rs
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"100 bytes", "tenant\x00" + strings.Repeat("a", 91) + "\u00e9", "05C1EFBD31E1740875666F626F7501" + strings.Repeat("62", 91) + "C4AA00"},
		{"101 bytes", "tenant\x00" + strings.Repeat("a", 92) + "\u00e9", "05C1C7CDE1BBB00875666F626F7501" + strings.Repeat("62", 92) + "C4AA"},
		{"truncate before hashing", "tenant\x00" + strings.Repeat("a", 92) + "\u00e9z", "05C1C7CDE1BBB00875666F626F7501" + strings.Repeat("62", 92) + "C4AA"},
		{"99 UTF-16 units", strings.Repeat("a", 97) + "\U0001f600", "05C1C973BF796808" + strings.Repeat("62", 97) + "F1A09981"},
		{"100 UTF-16 units", strings.Repeat("a", 98) + "\U0001f600", "05C1DD87B7E74C08" + strings.Repeat("62", 98) + "F1A099"},
		{"split surrogate pair", strings.Repeat("a", 99) + "\U0001f600", "05C1E153DB455808" + strings.Repeat("62", 99) + "F0C0"},
		{"truncate surrogate pair", strings.Repeat("a", 100) + "\U0001f600", "05C1EB5921F70608" + strings.Repeat("62", 100) + "00"},
		{"split UTF-8 bytes", strings.Repeat("a", 99) + "\u20ac", "05C1D747E54D9008" + strings.Repeat("62", 99) + "E383"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, ComputeV1([]interface{}{tt.value}))
		})
	}
}

func TestWriteNumberForBinaryEncodingV1(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"zero", 0, "058000"},
		{"negative zero", math.Copysign(0, -1), "058000"},
		{"integer", 5, "05C014"},
		{"fraction", 5.12312419050912359123, "05C0153F858949735550"},
		{"negative integer", -128, "053FA0"},
		{"smallest positive number", math.SmallestNonzeroFloat64, "05800101010101010102"},
		{"largest number", math.MaxFloat64, "05FFEFFFFFFFFFFFFFFE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var encoded []byte
			writeNumberForBinaryEncodingV1(tt.value, &encoded)
			require.Equal(t, tt.expected, toHexUpper(encoded))
		})
	}
}

func TestComputeV1DoesNotMutateComponents(t *testing.T) {
	value := strings.Repeat("a", 101)
	values := []interface{}{value}
	ComputeV1(values)
	require.Equal(t, []interface{}{value}, values)
	require.Empty(t, ComputeV1(nil))
}

func TestComputeV2UnicodeStringsNotTruncated(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"tenant\x00" + strings.Repeat("a", 92) + "\u00e9z", "3C0D04FB07AECB6B45A1592465FDFAA7"},
		{strings.Repeat("a", 99) + "\U0001f600", "3DCDC7ABFC5BF49867B328ECDB14F0BF"},
	}
	for _, tt := range tests {
		values := []interface{}{tt.value}
		require.Equal(t, tt.expected, ComputeV2HashForRouting(values))
		require.Equal(t, tt.expected, ComputeV2MultiHashForRouting(values))
	}
}

// The XML V1 fixtures describe raw per-component Murmur32 hashes, not routing
// EPKs. V2 fixtures likewise contain raw, unmasked hashes.
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

			t.Run(category+"/RawV1Hash/"+tc.Input.Description, func(t *testing.T) {
				var hashes strings.Builder
				for _, value := range values {
					var input []byte
					writeForHashing(value, &input)
					fmt.Fprintf(&hashes, "%032X", murmurhash3_32(input, 0))
				}
				require.Equal(t, expectedV1, hashes.String(),
					"raw hash mismatch for %s (value: %s)", tc.Input.Description, tc.Input.PartitionKeyValue)
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

func TestMaskTopBitsForRouting(t *testing.T) {
	// Already valid (first byte ≤ 0x3F) — unchanged
	require.Equal(t, "3FAABBCC", maskTopBitsForRouting("3FAABBCC"))
	// 0xFF & 0x3F = 0x3F
	require.Equal(t, "3FAABBCC", maskTopBitsForRouting("FFAABBCC"))
	// 0xC0 & 0x3F = 0x00
	require.Equal(t, "00AABBCC", maskTopBitsForRouting("C0AABBCC"))
	// 0x80 & 0x3F = 0x00
	require.Equal(t, "00112233", maskTopBitsForRouting("80112233"))
	// 0x40 & 0x3F = 0x00
	require.Equal(t, "00112233", maskTopBitsForRouting("40112233"))
	// Edge: empty string
	require.Equal(t, "", maskTopBitsForRouting(""))
	// Edge: single char
	require.Equal(t, "A", maskTopBitsForRouting("A"))
}

func TestComputeV2HashForRouting_MaskingApplied(t *testing.T) {
	// null produces a raw V2 hash with first byte 0x77, which should be masked to 0x37
	result := ComputeV2HashForRouting([]interface{}{nil})
	require.True(t, len(result) >= 2, "result should be at least 2 hex chars")
	// After masking, the first hex digit must be in [0-3]
	firstDigit := result[0]
	require.True(t, firstDigit >= '0' && firstDigit <= '3',
		"first hex digit should be in [0-3] after masking, got %c", firstDigit)
}

func TestComputeV2MultiHashForRouting_MaskingApplied(t *testing.T) {
	// Each 32-char component should have its first byte masked independently
	result := ComputeV2MultiHashForRouting([]interface{}{"hello", "world"})
	require.Equal(t, 64, len(result), "two components should produce 64 hex chars")
	// Check first byte of each 32-char component is in [0x00, 0x3F]
	for i := 0; i < 2; i++ {
		firstDigit := result[i*32]
		require.True(t, firstDigit >= '0' && firstDigit <= '3',
			"component %d first hex digit should be in [0-3] after masking, got %c", i, firstDigit)
	}
}

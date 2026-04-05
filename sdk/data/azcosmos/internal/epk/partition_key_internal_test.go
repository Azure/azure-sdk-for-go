// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package epk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromJsonString_EmptyPartitionKey(t *testing.T) {
	pk, err := FromJsonString("[]")
	require.NoError(t, err)
	require.True(t, pk.Equals(Empty))
	require.Equal(t, "[]", pk.ToJson())
}

func TestFromJsonString_VariousTypes(t *testing.T) {
	pk, err := FromJsonString(`["aa", null, true, false, {}, 5, 5.5]`)
	require.NoError(t, err)

	expected, err := FromObjectArray([]interface{}{"aa", nil, true, false, UndefinedMarker{}, float64(5), 5.5}, true)
	require.NoError(t, err)
	require.True(t, pk.Equals(expected))
	require.Equal(t, `["aa",null,true,false,{},5.0,5.5]`, pk.ToJson())
}

func TestFromJsonString_DeserializeEmptyString(t *testing.T) {
	_, err := FromJsonString("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to deserialize partition key value")
}

func TestFromJsonString_DeserializeNull(t *testing.T) {
	_, err := FromJsonString("null")
	require.Error(t, err)
}

func TestFromJsonString_InvalidString(t *testing.T) {
	_, err := FromJsonString("[aa]")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to deserialize partition key value")
}

func TestFromJsonString_InvalidNumber(t *testing.T) {
	_, err := FromJsonString("[1.a]")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to deserialize partition key value")
}

func TestFromJsonString_MissingBraces(t *testing.T) {
	_, err := FromJsonString("[{]")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to deserialize partition key value")
}

func TestFromJsonString_MissingValue(t *testing.T) {
	_, err := FromJsonString("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to deserialize partition key value: ")
}

func TestFromJsonString_MaxValue(t *testing.T) {
	pk, err := FromJsonString(`"Infinity"`)
	require.NoError(t, err)
	require.True(t, pk.Equals(ExclusiveMaximum))
}

func TestFromJsonString_MinValue(t *testing.T) {
	pk, err := FromJsonString("[]")
	require.NoError(t, err)
	require.True(t, pk.Equals(InclusiveMinimum))
}

func TestFromJsonString_UndefinedValue(t *testing.T) {
	pk, err := FromJsonString("[]")
	require.NoError(t, err)
	require.True(t, pk.Equals(Empty))
}

func TestFromJsonString_JsonConvertDefaultSettings(t *testing.T) {
	pk, err := FromJsonString("[123.0]")
	require.NoError(t, err)
	require.Equal(t, "[123.0]", pk.ToJson())
}

func TestFromJsonString_UnicodeCharacters(t *testing.T) {
	pk, err := FromJsonString(`["电脑"]`)
	require.NoError(t, err)
	require.Equal(t, `["电脑"]`, pk.ToJson())
}

func TestCompareTo(t *testing.T) {
	testCases := []struct {
		left     string
		right    string
		expected int
	}{
		{"[]", "[]", 0},
		{"[]", "[{}]", -1},
		{"[]", "[false]", -1},
		{"[]", "[true]", -1},
		{"[]", "[null]", -1},
		{"[]", "[2]", -1},
		{"[]", `["aa"]`, -1},
		{"[]", `"Infinity"`, -1},

		{"[{}]", "[]", 1},
		{"[{}]", "[{}]", 0},
		{"[{}]", "[false]", -1},
		{"[{}]", "[true]", -1},
		{"[{}]", "[null]", -1},
		{"[{}]", "[2]", -1},
		{"[{}]", `["aa"]`, -1},
		{"[{}]", `"Infinity"`, -1},

		{"[false]", "[]", 1},
		{"[false]", "[{}]", 1},
		{"[false]", "[null]", 1},
		{"[false]", "[false]", 0},
		{"[false]", "[true]", -1},
		{"[false]", "[2]", -1},
		{"[false]", `["aa"]`, -1},
		{"[false]", `"Infinity"`, -1},

		{"[true]", "[]", 1},
		{"[true]", "[{}]", 1},
		{"[true]", "[null]", 1},
		{"[true]", "[false]", 1},
		{"[true]", "[true]", 0},
		{"[true]", "[2]", -1},
		{"[true]", `["aa"]`, -1},
		{"[true]", `"Infinity"`, -1},

		{"[null]", "[]", 1},
		{"[null]", "[{}]", 1},
		{"[null]", "[null]", 0},
		{"[null]", "[false]", -1},
		{"[null]", "[true]", -1},
		{"[null]", "[2]", -1},
		{"[null]", `["aa"]`, -1},
		{"[null]", `"Infinity"`, -1},

		{"[2]", "[]", 1},
		{"[2]", "[{}]", 1},
		{"[2]", "[null]", 1},
		{"[2]", "[false]", 1},
		{"[2]", "[true]", 1},
		{"[1]", "[2]", -1},
		{"[2]", "[2]", 0},
		{"[3]", "[2]", 1},
		{"[2.1234344]", "[2]", 1},
		{"[2]", `["aa"]`, -1},
		{"[2]", `"Infinity"`, -1},

		{`["aa"]`, "[]", 1},
		{`["aa"]`, "[{}]", 1},
		{`["aa"]`, "[null]", 1},
		{`["aa"]`, "[false]", 1},
		{`["aa"]`, "[true]", 1},
		{`["aa"]`, "[2]", 1},
		{`[""]`, `["aa"]`, -1},
		{`["aa"]`, `["aa"]`, 0},
		{`["b"]`, `["aa"]`, 1},
		{`["aa"]`, `"Infinity"`, -1},

		{`"Infinity"`, "[]", 1},
		{`"Infinity"`, "[{}]", 1},
		{`"Infinity"`, "[null]", 1},
		{`"Infinity"`, "[false]", 1},
		{`"Infinity"`, "[true]", 1},
		{`"Infinity"`, "[2]", 1},
		{`"Infinity"`, `["aa"]`, 1},
		{`"Infinity"`, `"Infinity"`, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.left+"_vs_"+tc.right, func(t *testing.T) {
			left, err := FromJsonString(tc.left)
			require.NoError(t, err)
			right, err := FromJsonString(tc.right)
			require.NoError(t, err)
			require.Equal(t, tc.expected, left.CompareTo(right))
		})
	}
}

func TestFromObjectArray_InvalidPartitionKeyValue(t *testing.T) {
	type customType struct{}
	_, err := FromObjectArray([]interface{}{2, true, customType{}}, true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid partition key value type")
}

func TestContains(t *testing.T) {
	testCases := []struct {
		parent   string
		child    string
		expected bool
	}{
		{"[]", "[]", true},
		{"[]", "[{}]", true},
		{"[]", "[null]", true},
		{"[]", "[true]", true},
		{"[]", "[false]", true},
		{"[]", "[2]", true},
		{"[]", `["fdfd"]`, true},

		{"[2]", "[]", false},
		{"[2]", "[2]", true},
		{`[2]`, `[2, "USA"]`, true},
		{`[1]`, `[2, "USA"]`, false},
	}

	for _, tc := range testCases {
		t.Run(tc.parent+"_contains_"+tc.child, func(t *testing.T) {
			parent, err := FromJsonString(tc.parent)
			require.NoError(t, err)
			child, err := FromJsonString(tc.child)
			require.NoError(t, err)
			require.Equal(t, tc.expected, parent.Contains(child))
		})
	}
}

func TestFromObjectArray_InvalidPartitionKeyValueNonStrict(t *testing.T) {
	type customType struct{}
	strictResult, err := FromObjectArray([]interface{}{2, true, UndefinedMarker{}}, true)
	require.NoError(t, err)

	nonStrictResult, err := FromObjectArray([]interface{}{2, true, customType{}}, false)
	require.NoError(t, err)

	require.True(t, strictResult.Equals(nonStrictResult))
}

func TestToJson_ExclusiveMaximum(t *testing.T) {
	require.Equal(t, `"Infinity"`, ExclusiveMaximum.ToJson())
}

func TestToJson_Numbers(t *testing.T) {
	pk, _ := FromJsonString("[5]")
	require.Equal(t, "[5.0]", pk.ToJson())

	pk, _ = FromJsonString("[5.5]")
	require.Equal(t, "[5.5]", pk.ToJson())
}

func TestComponents(t *testing.T) {
	pk, _ := FromJsonString(`["a", null, true]`)
	comps := pk.Components()
	require.Len(t, comps, 3)
	require.Equal(t, "a", comps[0])
	require.Nil(t, comps[1])
	require.Equal(t, true, comps[2])
}

func TestIsEmpty(t *testing.T) {
	require.True(t, Empty.IsEmpty())
	require.True(t, InclusiveMinimum.IsEmpty())

	pk, _ := FromJsonString("[1]")
	require.False(t, pk.IsEmpty())
}

func TestValidateComponentCount_TooFewPartitionKeyComponents(t *testing.T) {
	pk, err := FromJsonString(`["PartitionKeyValue"]`)
	require.NoError(t, err)

	err = pk.ValidateComponentCount(2)
	require.Error(t, err)
	require.Equal(t, ErrTooFewPartitionKeyComponents, err)
}

func TestValidateComponentCount_TooManyPartitionKeyComponents(t *testing.T) {
	pk, err := FromJsonString("[true, false]")
	require.NoError(t, err)

	err = pk.ValidateComponentCount(1)
	require.Error(t, err)
	require.Equal(t, ErrTooManyPartitionKeyComponents, err)
}

func TestValidateComponentCount_ExactMatch(t *testing.T) {
	pk, err := FromJsonString(`["value1", "value2"]`)
	require.NoError(t, err)

	err = pk.ValidateComponentCount(2)
	require.NoError(t, err)
}

func TestValidateComponentCount_SpecialValues(t *testing.T) {
	require.NoError(t, Empty.ValidateComponentCount(1))
	require.NoError(t, InclusiveMinimum.ValidateComponentCount(2))
	require.NoError(t, ExclusiveMaximum.ValidateComponentCount(3))
}

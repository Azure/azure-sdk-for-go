// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegionId_Canonicalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple lowercase",
			input:    "westus",
			expected: "westus",
		},
		{
			name:     "Mixed case",
			input:    "WestUS",
			expected: "westus",
		},
		{
			name:     "With spaces",
			input:    "West US",
			expected: "westus",
		},
		{
			name:     "With multiple spaces",
			input:    " West   US ",
			expected: "westus",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := newRegionId(tt.input)
			assert.Equal(t, tt.expected, id.String())
			assert.Equal(t, regionId(tt.expected), id)
		})
	}
}

func TestRegionId_Equal(t *testing.T) {
	id1 := newRegionId("West US")
	id2 := newRegionId("westus")
	id3 := newRegionId("East US")

	assert.True(t, id1.Equal(id2))
	assert.False(t, id1.Equal(id3))
}

func TestRegionId_UnmarshalJSON(t *testing.T) {
	jsonStr := `" West US "`
	var id regionId
	err := id.UnmarshalJSON([]byte(jsonStr))
	assert.NoError(t, err)
	assert.Equal(t, "westus", id.String())
}

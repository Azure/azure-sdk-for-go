// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAzdErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single JSON message",
			input:    `{"type":"consoleMessage","timestamp":"2024-01-01T00:00:00Z","data":{"message":"\nERROR: fetching token: authentication failed\n"}}`,
			expected: "ERROR: fetching token: authentication failed",
		},
		{
			name:     "plain text error (not JSON)",
			input:    "ERROR: plain text error message",
			expected: "ERROR: plain text error message",
		},
		{
			name:     "empty message",
			input:    "",
			expected: "",
		},
		{
			name:     "JSON with empty data.message",
			input:    `{"type":"consoleMessage","timestamp":"2024-01-01T00:00:00Z","data":{"message":""}}`,
			expected: `{"type":"consoleMessage","timestamp":"2024-01-01T00:00:00Z","data":{"message":""}}`,
		},
		{
			name:     "whitespace trimming",
			input:    `{"type":"consoleMessage","timestamp":"2024-01-01T00:00:00Z","data":{"message":"   \n\nERROR: token failed   \n\n"}}`,
			expected: "ERROR: token failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAzdErrorMessage(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

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
			name: "multiple JSON messages",
			input: `{"type":"consoleMessage","timestamp":"2024-01-01T00:00:00Z","data":{"message":"\nERROR: fetching token: AADSTS50079\n"}}
{"type":"consoleMessage","timestamp":"2024-01-01T00:00:01Z","data":{"message":"Suggestion: run azd auth login\n"}}`,
			expected: "ERROR: fetching token: AADSTS50079 Suggestion: run azd auth login",
		},
		{
			name:     "plain text error (not JSON)",
			input:    "ERROR: plain text error message",
			expected: "ERROR: plain text error message",
		},
		{
			name:     "mixed JSON and plain text",
			input:    `{"type":"consoleMessage","timestamp":"2024-01-01T00:00:00Z","data":{"message":"ERROR: token expired\n"}}` + "\nsome plain text",
			expected: "ERROR: token expired",
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
		{
			name: "real world multi-line example",
			input: `{"type":"consoleMessage","timestamp":"...","data":{"message":"\nERROR: fetching token: AADSTS50079: Due to a configuration change made by your administrator, or because you moved to a new location, you must enroll in multi-factor authentication'. Trace ID: ... Correlation ID: ... Timestamp: ...\n"}}
{"type":"consoleMessage","timestamp":"...","data":{"message":"Suggestion: reauthentication required, run azd auth login --scope ... to acquire a new token.\n"}}`,
			expected: "ERROR: fetching token: AADSTS50079: Due to a configuration change made by your administrator, or because you moved to a new location, you must enroll in multi-factor authentication'. Trace ID: ... Correlation ID: ... Timestamp: ... Suggestion: reauthentication required, run azd auth login --scope ... to acquire a new token.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAzdErrorMessage(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

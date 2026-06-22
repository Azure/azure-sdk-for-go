// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractAzdError(t *testing.T) {
	aadError := "AADSTS90002: Tenant 'test' not found. Check to make sure you have the correct tenant ID"

	tests := []struct {
		name     string
		stderr   string
		expected string
	}{
		{
			name:     "v1.24.0+ structured error only",
			stderr:   `{"error":"` + aadError + `","message":"Authentication with Azure failed.","suggestion":"Run 'azd auth login' to sign in again."}`,
			expected: aadError,
		},
		{
			name: "v1.23.7 structured error preceded by empty consoleMessage",
			stderr: `{"type":"consoleMessage","timestamp":"2026-04-14T22:02:36.154700776Z","data":{"message":"\n"}}
{"error":"` + aadError + `","message":"Authentication with Azure failed.","suggestion":"Run 'azd auth login' to sign in again."}`,
			expected: aadError,
		},
		{
			name:     "pre-v1.23.7 legacy consoleMessage",
			stderr:   `{"type":"consoleMessage","timestamp":"2026-04-14T22:03:37.687535934Z","data":{"message":"\nERROR: ` + aadError + `\n\n"}}`,
			expected: "ERROR: " + aadError,
		},
		{
			name: "pre-v1.23.7 multiple consoleMessage lines",
			stderr: `{"type":"consoleMessage","timestamp":"...","data":{"message":"\nERROR: ` + aadError + `\n"}}
{"type":"consoleMessage","timestamp":"...","data":{"message":"Suggestion: reauthentication required\n"}}`,
			expected: "ERROR: " + aadError,
		},
		{
			name:     "non-JSON plain text preserved",
			stderr:   "some plain error text",
			expected: "some plain error text",
		},
		{
			name:     "empty consoleMessage not extracted",
			stderr:   `{"type":"consoleMessage","data":{"message":"\n"}}`,
			expected: `{"type":"consoleMessage","data":{"message":"\n"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := extractAzdError(tt.stderr)
			require.Equal(t, tt.expected, actual)
		})
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package sanitizer

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSanitizedHeader(t *testing.T) {
	s := New(Options{AllowedHeaders: []string{"allowed"}})
	header := http.Header{
		"allowed":      []string{"first", "second"},
		"Content-Type": []string{"application/json"},
		"secret":       []string{"hidden"},
	}
	require.Equal(t, "first", s.GetSanitizedHeader("allowed", header))
	require.Equal(t, "application/json", s.GetSanitizedHeader("content-type", header))
	require.Equal(t, redactedValue, s.GetSanitizedHeader("secret", header))
	require.Empty(t, s.GetSanitizedHeader("allowed", nil))
}

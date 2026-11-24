// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratedFileUsesCurrentSourcePin(t *testing.T) {
	contents, err := os.ReadFile("../../../region_proximity.go")
	require.NoError(t, err)

	provenance := fmt.Sprintf("at %s (SHA-256 %s); DO NOT EDIT.", rustRevision, rustSourceHash)
	require.True(t, strings.Contains(string(contents), provenance),
		"generated file does not contain current provenance %q", provenance)
}

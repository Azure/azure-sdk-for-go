// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"math"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAvailableSystemMemoryBytes asserts the per-platform memory query returns
// a positive reading on supported platforms and the "unknown" sentinel (0)
// elsewhere.
func TestAvailableSystemMemoryBytes(t *testing.T) {
	got := availableSystemMemoryBytes()
	switch runtime.GOOS {
	case "linux", "darwin", "windows":
		require.Greater(t, got, uint64(0), "expected a positive memory reading on %s", runtime.GOOS)
	default:
		require.Equal(t, uint64(0), got, "unsupported platforms report 0 (unknown)")
	}
}

// TestCheckBufferMemoryBudget covers the buffer-method memory guard: tiny and
// non-positive requests always pass, while an astronomically large request
// fails only when available memory is known (>0). When the platform can't
// determine available memory the check is skipped rather than guessed.
func TestCheckBufferMemoryBudget(t *testing.T) {
	for _, size := range []int64{0, -1, 1, 1024} {
		require.NoError(t, checkBufferMemoryBudget("--upload-method buffer", size),
			"size %d must pass the budget check", size)
	}

	err := checkBufferMemoryBudget("--upload-method buffer", math.MaxInt64)
	if availableSystemMemoryBytes() == 0 {
		require.NoError(t, err, "budget check must be skipped when available memory is unknown")
		return
	}
	require.Error(t, err, "an astronomically large buffer must exceed the memory budget")
	require.Contains(t, err.Error(), "--upload-method buffer", "error must echo the triggering flag label")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConcurrencyValue_InBounds(t *testing.T) {
	val := DefaultConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))
}

func TestDefaultConcurrencyValue_Deterministic(t *testing.T) {
	val1 := DefaultConcurrencyValue()
	val2 := DefaultConcurrencyValue()
	require.Equal(t, val1, val2)
}

func TestDefaultConcurrencyValue_MatchesCPU(t *testing.T) {
	cpus := runtime.NumCPU()
	val := DefaultConcurrencyValue()
	if cpus < 8 {
		require.Equal(t, uint16(8), val)
	} else if cpus > 96 {
		require.Equal(t, uint16(96), val)
	} else {
		require.Equal(t, uint16(cpus), val)
	}
}

func TestDefaultConcurrencyValue_LegacyEnvVar(t *testing.T) {
	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "true")
	require.Equal(t, uint16(DefaultConcurrency), DefaultConcurrencyValue())

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "TRUE")
	require.Equal(t, uint16(DefaultConcurrency), DefaultConcurrencyValue())

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "false")
	val := DefaultConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))

	t.Setenv("AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY", "")
	val = DefaultConcurrencyValue()
	require.GreaterOrEqual(t, val, uint16(8))
	require.LessOrEqual(t, val, uint16(96))
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStartCPUProfileDisabled(t *testing.T) {
	stop, err := startCPUProfile(false, "")
	require.NoError(t, err)
	require.Nil(t, stop)
}

func TestStartCPUProfileRequiresPath(t *testing.T) {
	stop, err := startCPUProfile(true, "")
	require.EqualError(t, err, "CPU profile path cannot be empty")
	require.Nil(t, stop)
}

func TestStartCPUProfileWritesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles", "cpu.pprof")
	stop, err := startCPUProfile(true, path)
	require.NoError(t, err)
	require.NotNil(t, stop)

	deadline := time.Now().Add(25 * time.Millisecond)
	result := uint64(1)
	for time.Now().Before(deadline) {
		result = result*1664525 + 1013904223
	}
	require.NotZero(t, result)
	require.NoError(t, stop())

	stat, err := os.Stat(path)
	require.NoError(t, err)
	require.Positive(t, stat.Size())
}

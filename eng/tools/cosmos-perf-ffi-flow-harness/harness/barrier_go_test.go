//go:build linux

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMeasurementBarrierAndMonotonicBounds(t *testing.T) {
	directory := t.TempDir()
	t.Setenv("COSMOS_PERF_BARRIER_DIR", directory)
	require.NoError(t, os.WriteFile(filepath.Join(directory, "ready-start"), nil, 0600))
	require.NoError(t, measurementBarrier("start"))
	_, err := os.Stat(filepath.Join(directory, "request-start"))
	require.NoError(t, err)
	begin, err := monotonicNS()
	require.NoError(t, err)
	end, err := monotonicNS()
	require.NoError(t, err)
	require.GreaterOrEqual(t, end, begin)
	require.NoError(t, measurementBounds(begin, end))
	data, err := os.ReadFile(filepath.Join(directory, "measurement-bounds.json"))
	require.NoError(t, err)
	var value map[string]uint64
	require.NoError(t, json.Unmarshal(data, &value))
	require.Equal(t, begin, value["monotonicStartNs"])
	require.Equal(t, end, value["monotonicEndNs"])
	require.Error(t, measurementBounds(begin, end), "existing bounds must not be overwritten")
}

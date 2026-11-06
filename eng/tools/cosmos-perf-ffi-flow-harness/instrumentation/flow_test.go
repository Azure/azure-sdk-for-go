// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func flowNativeEpoch() uint64 {
	return uint64(time.Now().Unix()) - 3
}

func flowNativePhase(uint32) uint64 {
	return uint64(time.Now().UnixNano())
}

func flowNativeWrite(string) int {
	return 0
}

func TestFlowPhasesConcurrentCountersAndHistogram(t *testing.T) {
	require.Error(t, ExperimentalFlowPhase(3))
	require.NoError(t, ExperimentalFlowPhase(1))
	flowEvent(0, 900)
	require.NoError(t, ExperimentalFlowPhase(2))
	group := sync.WaitGroup{}
	for range 10 {
		group.Add(1)
		go func() {
			defer group.Done()
			for range 100 {
				flowEvent(0, 1)
				flowEvent(1, 1)
				flowWait(1)
				flowEvent(6, 1)
				flowResumed(time.Now())
			}
		}()
	}
	group.Wait()
	flowWait(3)
	flowEvent(6, 3)
	for range 3 {
		flowResumed(time.Now())
	}
	flowWait(0)
	require.Error(t, ExperimentalFlowWrite(t.TempDir()), "measurement must end before snapshot")
	require.NoError(t, ExperimentalFlowPhase(0))
	directory := t.TempDir()
	require.NoError(t, ExperimentalFlowWrite(directory))
	data, err := os.ReadFile(filepath.Join(directory, "go-flow.json"))
	require.NoError(t, err)
	var value struct {
		Overflow uint64 `json:"overflow"`
		Rows     []struct {
			Phase     int               `json:"phase"`
			Counts    map[string]uint64 `json:"counts"`
			Histogram map[string]uint64 `json:"batchHistogram"`
		} `json:"rows"`
	}
	require.NoError(t, json.Unmarshal(data, &value))
	require.Zero(t, value.Overflow)
	totals := map[string]uint64{}
	histogram := map[string]uint64{}
	for _, row := range value.Rows {
		if row.Phase != 2 {
			continue
		}
		for name, count := range row.Counts {
			totals[name] += count
		}
		for size, count := range row.Histogram {
			histogram[size] += count
		}
	}
	require.EqualValues(t, 1000, totals["submitAttempts"])
	require.EqualValues(t, 1000, totals["submitAccepted"])
	require.EqualValues(t, 1002, totals["waitReturns"])
	require.EqualValues(t, 1, totals["emptyWaits"])
	require.EqualValues(t, 1003, totals["drained"])
	require.EqualValues(t, 1003, totals["delivered"])
	require.EqualValues(t, 1003, totals["resumed"])
	require.Equal(t, map[string]uint64{"1": 1000, "3": 1}, histogram)
	require.Error(t, ExperimentalFlowWrite(directory), "existing ledger must not be overwritten")
}

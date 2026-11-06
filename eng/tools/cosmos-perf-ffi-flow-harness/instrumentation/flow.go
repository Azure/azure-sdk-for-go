// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

const flowSeconds = 1024

const (
	flowSubmitAttempts = iota
	flowSubmitAccepted
	flowSubmitErrors
	flowWaitReturns
	flowEmptyWaits
	flowDrained
	flowDelivered
	flowResumedCount
	flowCount
)

var flowNames = [flowCount]string{
	"submitAttempts", "submitAccepted", "submitErrors", "waitReturns",
	"emptyWaits", "drained", "delivered", "resumed",
}

type flowBucket struct {
	counts          [flowCount]atomic.Uint64
	batch           [65]atomic.Uint64
	resumeNanos     atomic.Uint64
	resumeMaxNanos  atomic.Uint64
	resumeHistogram [65]atomic.Uint64
}

var flowCounters struct {
	gate     sync.RWMutex
	phase    uint32
	epoch    uint64
	init     sync.Once
	start    [3]uint64
	end      [3]uint64
	buckets  [3][flowSeconds]flowBucket
	overflow atomic.Uint64
}

func flowBucketLocked() *flowBucket {
	flowCounters.init.Do(func() { flowCounters.epoch = flowNativeEpoch() })
	now := time.Now().Unix()
	if now < 0 || uint64(now) < flowCounters.epoch || uint64(now)-flowCounters.epoch >= flowSeconds {
		flowCounters.overflow.Add(1)
		return nil
	}
	return &flowCounters.buckets[flowCounters.phase][uint64(now)-flowCounters.epoch]
}

func flowEvent(index int, count uint64) {
	flowCounters.gate.RLock()
	defer flowCounters.gate.RUnlock()
	if bucket := flowBucketLocked(); bucket != nil {
		bucket.counts[index].Add(count)
	}
}

func flowWait(count int) {
	flowCounters.gate.RLock()
	defer flowCounters.gate.RUnlock()
	bucket := flowBucketLocked()
	if bucket == nil {
		return
	}
	bucket.counts[flowWaitReturns].Add(1)
	if count == 0 {
		bucket.counts[flowEmptyWaits].Add(1)
		return
	}
	if count < 0 || count >= len(bucket.batch) {
		flowCounters.overflow.Add(1)
		return
	}
	bucket.counts[flowDrained].Add(uint64(count))
	bucket.batch[count].Add(1)
}

func flowResumed(deliveredAt time.Time) {
	if deliveredAt.IsZero() {
		flowCounters.overflow.Add(1)
		return
	}
	delay := time.Since(deliveredAt)
	if delay < 0 {
		flowCounters.overflow.Add(1)
		return
	}
	flowCounters.gate.RLock()
	defer flowCounters.gate.RUnlock()
	if bucket := flowBucketLocked(); bucket != nil {
		value := uint64(delay)
		bucket.counts[flowResumedCount].Add(1)
		bucket.resumeNanos.Add(value)
		bucket.resumeHistogram[bits.Len64(value)].Add(1)
		for old := bucket.resumeMaxNanos.Load(); value > old; old = bucket.resumeMaxNanos.Load() {
			if bucket.resumeMaxNanos.CompareAndSwap(old, value) {
				break
			}
		}
	}
}

// ExperimentalFlowPhase changes phase only at a drained benchmark phase boundary.
// Available only in the isolated diagnostic SDK copy, never in the clean binary.
func ExperimentalFlowPhase(phase uint32) error {
	if phase >= 3 {
		return errors.New("flow phase must be 0, 1, or 2")
	}
	flowCounters.gate.Lock()
	defer flowCounters.gate.Unlock()
	flowCounters.init.Do(func() { flowCounters.epoch = flowNativeEpoch() })
	now := flowNativePhase(phase)
	flowCounters.end[flowCounters.phase] = now
	flowCounters.start[phase] = now
	flowCounters.phase = phase
	return nil
}

// ExperimentalFlowWrite preserves the independent Go and Rust ledgers after client shutdown.
func ExperimentalFlowWrite(directory string) error {
	flowCounters.gate.Lock()
	defer flowCounters.gate.Unlock()
	if flowCounters.phase != 0 {
		return errors.New("flow output requires the measurement phase to have ended")
	}
	rows := make([]map[string]any, 0)
	for phase := range flowCounters.buckets {
		for second := range flowCounters.buckets[phase] {
			bucket := &flowCounters.buckets[phase][second]
			counts := make(map[string]uint64, len(flowNames))
			var total uint64
			for index, name := range flowNames {
				counts[name] = bucket.counts[index].Load()
				total += counts[name]
			}
			if total == 0 {
				continue
			}
			histogram := make(map[int]uint64)
			for size := 1; size < len(bucket.batch); size++ {
				if count := bucket.batch[size].Load(); count != 0 {
					histogram[size] = count
				}
			}
			rows = append(rows, map[string]any{
				"phase": phase, "unixSecond": flowCounters.epoch + uint64(second),
				"counts": counts, "batchHistogram": histogram,
				"deliveryToResumeTotalNs": bucket.resumeNanos.Load(),
				"deliveryToResumeMaxNs":   bucket.resumeMaxNanos.Load(),
				"deliveryToResumeLog2Ns": func() map[int]uint64 {
					values := make(map[int]uint64)
					for bin := range bucket.resumeHistogram {
						if count := bucket.resumeHistogram[bin].Load(); count != 0 {
							values[bin] = count
						}
					}
					return values
				}(),
			})
		}
	}
	value := map[string]any{
		"schemaVersion": 1, "rows": rows, "overflow": flowCounters.overflow.Load(),
		"phaseStartUnixNs": flowCounters.start, "phaseEndUnixNs": flowCounters.end,
		"deliveryDefinition":   "completion resolved to its Go operation, counted immediately before publishing its result",
		"resumptionDefinition": "waiting goroutine received its completion, before SDK response preparation; delay includes channel publication and scheduling, not pure scheduler time",
		"phaseDefinition":      "0=setup/teardown, 1=excluded warmup, 2=measurement including final closed-loop drain",
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal Go flow counters: %w", err)
	}
	file, err := os.OpenFile(filepath.Join(directory, "go-flow.json"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open Go flow output: %w", err)
	}
	_, writeErr := file.Write(append(data, '\n'))
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("write Go flow counters: %w", err)
	}
	if code := flowNativeWrite(filepath.Join(directory, "native-flow.json")); code != 0 {
		return fmt.Errorf("write native flow counters: status %d", code)
	}
	return nil
}

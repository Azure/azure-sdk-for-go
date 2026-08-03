// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"fmt"
	"runtime"
	"runtime/metrics"
	"sync"
	"time"
)

const (
	// cpuTotalMetric is the cumulative CPU budget for the Go runtime, i.e.
	// GOMAXPROCS * wallTime. It INCLUDES /cpu/classes/idle (time the scheduler
	// had a P available but no goroutine was runnable), so it must not be used
	// directly to compute CPU utilization. We subtract cpuIdleMetric from it to
	// obtain CPU-seconds actually consumed by user code, the GC, and the
	// scavenger.
	cpuTotalMetric    = "/cpu/classes/total:cpu-seconds"
	cpuIdleMetric     = "/cpu/classes/idle:cpu-seconds"
	memoryTotalMetric = "/memory/classes/total:bytes"
)

// processStatsSampler periodically samples the process CPU usage (in percent of all
// available cores) and total memory obtained from the OS. It is intended to mirror
// the CPU/memory tracking added to the .NET/Python perf-automation runners.
type processStatsSampler struct {
	interval      time.Duration
	done          chan struct{}
	wg            sync.WaitGroup
	mu            sync.Mutex
	cpuSamples    []float64
	memorySamples []uint64
	startOnce     sync.Once
	stopOnce      sync.Once
	started       bool
}

func newProcessStatsSampler(interval time.Duration) *processStatsSampler {
	return &processStatsSampler{
		interval: interval,
		done:     make(chan struct{}),
	}
}

// Start begins periodic sampling in a background goroutine. Safe to call
// multiple times; only the first call has effect.
func (s *processStatsSampler) Start() {
	s.startOnce.Do(func() {
		s.started = true
		s.wg.Add(1)
		go s.run()
	})
}

// Stop halts sampling and waits for the sampler goroutine to exit. Safe to call
// multiple times; only the first call closes the done channel.
func (s *processStatsSampler) Stop() {
	if !s.started {
		return
	}
	s.stopOnce.Do(func() {
		close(s.done)
	})
	s.wg.Wait()
}

func (s *processStatsSampler) run() {
	defer s.wg.Done()

	samples := []metrics.Sample{
		{Name: cpuTotalMetric},
		{Name: cpuIdleMetric},
		{Name: memoryTotalMetric},
	}
	metrics.Read(samples)
	lastCPU, cpuOK := readBusyCPUSeconds(samples[0], samples[1])
	lastTime := time.Now()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case now := <-ticker.C:
			metrics.Read(samples)
			cpu, ok := readBusyCPUSeconds(samples[0], samples[1])
			mem, memOK := readMemoryBytes(samples[2])
			elapsed := now.Sub(lastTime).Seconds()
			if cpuOK && ok && elapsed > 0 {
				cpuPct := (cpu - lastCPU) / (elapsed * float64(runtime.NumCPU())) * 100.0
				s.mu.Lock()
				s.cpuSamples = append(s.cpuSamples, cpuPct)
				s.mu.Unlock()
			}
			if memOK {
				s.mu.Lock()
				s.memorySamples = append(s.memorySamples, mem)
				s.mu.Unlock()
			}
			lastCPU, cpuOK = cpu, ok
			lastTime = now
		}
	}
}

// AverageCPU returns the average sampled CPU percent across all available cores,
// or -1 if no samples were collected.
func (s *processStatsSampler) AverageCPU() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.cpuSamples) == 0 {
		return -1
	}
	var sum float64
	for _, v := range s.cpuSamples {
		sum += v
	}
	return sum / float64(len(s.cpuSamples))
}

// AverageMemory returns the average sampled memory in bytes, or -1 if no samples
// were collected.
func (s *processStatsSampler) AverageMemory() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.memorySamples) == 0 {
		return -1
	}
	var sum uint64
	for _, v := range s.memorySamples {
		sum += v
	}
	return int64(sum / uint64(len(s.memorySamples)))
}

// SampleCount returns the number of samples collected so far.
func (s *processStatsSampler) SampleCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cpuSamples)
}

// LastSample returns the most recently observed CPU% and memory in bytes
// suitable for emitting on the live status line. Returns (-1, 0) when no
// samples have been collected yet.
func (s *processStatsSampler) LastSample() (cpuPercent float64, memoryBytes uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.cpuSamples) == 0 {
		return -1, 0
	}
	cpuPercent = s.cpuSamples[len(s.cpuSamples)-1]
	if len(s.memorySamples) > 0 {
		memoryBytes = s.memorySamples[len(s.memorySamples)-1]
	}
	return cpuPercent, memoryBytes
}

// Summary returns a single-line human-readable summary of the collected stats.
func (s *processStatsSampler) Summary() string {
	cpu := s.AverageCPU()
	mem := s.AverageMemory()
	memMiB := -1.0
	if mem >= 0 {
		memMiB = float64(mem) / (1024.0 * 1024.0)
	}
	return fmt.Sprintf(
		"Process stats: averageCpu=%.2f%% averageMemory(MiB)=%.2f samples=%d",
		cpu, memMiB, s.SampleCount(),
	)
}

func readCPUSeconds(s metrics.Sample) (float64, bool) {
	switch s.Value.Kind() {
	case metrics.KindFloat64:
		return s.Value.Float64(), true
	case metrics.KindUint64:
		return float64(s.Value.Uint64()), true
	}
	return 0, false
}

// readBusyCPUSeconds returns total - idle, i.e. the CPU-seconds actually
// consumed by Go code (user goroutines + GC + scavenger), excluding the
// scheduler-idle time that /cpu/classes/total includes. This matches what
// the .NET runner reports via Process.TotalProcessorTime, which counts only
// busy CPU time.
func readBusyCPUSeconds(total, idle metrics.Sample) (float64, bool) {
	totalSec, ok := readCPUSeconds(total)
	if !ok {
		return 0, false
	}
	idleSec, ok := readCPUSeconds(idle)
	if !ok {
		// If the idle class isn't available on this Go version, fall back to
		// the total so the metric at least remains monotonic (even though the
		// resulting percentage will be inflated).
		return totalSec, true
	}
	return totalSec - idleSec, true
}

func readMemoryBytes(s metrics.Sample) (uint64, bool) {
	switch s.Value.Kind() {
	case metrics.KindUint64:
		return s.Value.Uint64(), true
	case metrics.KindFloat64:
		return uint64(s.Value.Float64()), true
	}
	return 0, false
}

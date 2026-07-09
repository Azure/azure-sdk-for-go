// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package stats tracks latency histograms and writes perf result documents.
package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

const maxLatencyUS int64 = 3_600_000_000

// Stats collects per-operation latency and error counts.
type Stats struct {
	shards map[string]*operationStatsShard
}

type operationStatsShard struct {
	mu sync.Mutex
	operationStats
}

type operationStats struct {
	histogram *hdrhistogram.Histogram
	count     uint64
	min       time.Duration
	max       time.Duration
	sum       time.Duration
	errors    uint64

	backendHistogram *hdrhistogram.Histogram
	backendCount     uint64
	backendMin       time.Duration
	backendMax       time.Duration
	backendSum       time.Duration
}

// Summary is a drained interval summary for one operation.
type Summary struct {
	Name        string
	Count       uint64
	Errors      uint64
	Min         time.Duration
	Max         time.Duration
	Mean        time.Duration
	P50         time.Duration
	P90         time.Duration
	P99         time.Duration
	BackendMin  *time.Duration
	BackendMax  *time.Duration
	BackendMean *time.Duration
	BackendP50  *time.Duration
	BackendP90  *time.Duration
	BackendP99  *time.Duration
}

// ProcessMetrics contains process and system CPU/memory values.
type ProcessMetrics struct {
	CPUPercent             float64
	MemoryBytes            uint64
	SystemCPUPercent       float64
	SystemTotalMemoryBytes uint64
	SystemUsedMemoryBytes  uint64
	CgroupCPUPercent       *float64
}

// ConfigSnapshot is emitted with each PerfResult document.
type ConfigSnapshot struct {
	Concurrency       uint64
	ApplicationRegion string
	PreferredRegions  string
	ExcludedRegions   string
	GOMAXPROCS        uint64
	PPCBEnabled       bool
	PyroscopeEnabled  bool
}

// PerfResult is the summary document schema stored in Cosmos DB.
type PerfResult struct {
	ID           string  `json:"id"`
	PartitionKey string  `json:"partition_key"`
	WorkloadID   string  `json:"workload_id"`
	CommitSHA    string  `json:"commit_sha"`
	Hostname     string  `json:"hostname"`
	Timestamp    string  `json:"TIMESTAMP"`
	Operation    string  `json:"operation"`
	Count        uint64  `json:"count"`
	Errors       uint64  `json:"errors"`
	MinMS        float64 `json:"min_ms"`
	MaxMS        float64 `json:"max_ms"`
	MeanMS       float64 `json:"mean_ms"`
	P50MS        float64 `json:"p50_ms"`
	P90MS        float64 `json:"p90_ms"`
	P99MS        float64 `json:"p99_ms"`

	BackendMinMS  *float64 `json:"backend_min_ms,omitempty"`
	BackendMaxMS  *float64 `json:"backend_max_ms,omitempty"`
	BackendMeanMS *float64 `json:"backend_mean_ms,omitempty"`
	BackendP50MS  *float64 `json:"backend_p50_ms,omitempty"`
	BackendP90MS  *float64 `json:"backend_p90_ms,omitempty"`
	BackendP99MS  *float64 `json:"backend_p99_ms,omitempty"`

	CPUPercent             float64  `json:"cpu_percent"`
	MemoryBytes            uint64   `json:"memory_bytes"`
	SystemCPUPercent       float64  `json:"system_cpu_percent"`
	SystemTotalMemoryBytes uint64   `json:"system_total_memory_bytes"`
	SystemUsedMemoryBytes  uint64   `json:"system_used_memory_bytes"`
	CgroupCPUPercent       *float64 `json:"cgroup_cpu_percent,omitempty"`

	ConfigConcurrency       *uint64 `json:"config_concurrency,omitempty"`
	ConfigApplicationRegion *string `json:"config_application_region,omitempty"`
	ConfigPreferredRegions  *string `json:"config_preferred_regions,omitempty"`
	ConfigExcludedRegions   *string `json:"config_excluded_regions,omitempty"`
	ConfigPPCBEnabled       *bool   `json:"config_ppcb_enabled,omitempty"`
	ConfigPyroscopeEnabled  *bool   `json:"config_pyroscope_enabled,omitempty"`
	ConfigGOMAXPROCS        *uint64 `json:"config_gomaxprocs,omitempty"`
	RuntimeGoroutines       uint64  `json:"runtime_goroutines"`
}

// ErrorResult is the per-operation failure document schema.
type ErrorResult struct {
	ID            string  `json:"id"`
	PartitionKey  string  `json:"partition_key"`
	WorkloadID    string  `json:"workload_id"`
	CommitSHA     string  `json:"commit_sha"`
	Hostname      string  `json:"hostname"`
	Timestamp     string  `json:"TIMESTAMP"`
	Operation     string  `json:"operation"`
	ErrorMessage  string  `json:"error_message"`
	SourceMessage *string `json:"source_message"`
}

// New creates a new stats collector.
func New(operationNames []string) *Stats {
	shards := make(map[string]*operationStatsShard, len(operationNames))
	for _, name := range operationNames {
		shards[name] = &operationStatsShard{operationStats: newOperationStats()}
	}
	return &Stats{shards: shards}
}

func newOperationStats() operationStats {
	return operationStats{
		histogram:        hdrhistogram.New(1, maxLatencyUS, 3),
		min:              time.Duration(1<<63 - 1),
		backendHistogram: hdrhistogram.New(1, maxLatencyUS, 3),
		backendMin:       time.Duration(1<<63 - 1),
	}
}

// RecordLatency records a successful operation latency.
func (s *Stats) RecordLatency(operation string, latency time.Duration, backend *time.Duration) {
	shard := s.shards[operation]
	if shard == nil {
		return
	}
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.record(latency, backend)
}

// RecordError records an operation error.
func (s *Stats) RecordError(operation string) {
	shard := s.shards[operation]
	if shard == nil {
		return
	}
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.errors++
}

// DrainSummaries drains and resets interval statistics.
func (s *Stats) DrainSummaries() []Summary {
	summaries := make([]Summary, 0, len(s.shards))
	for name, shard := range s.shards {
		shard.mu.Lock()
		drained := shard.operationStats
		shard.operationStats = newOperationStats()
		shard.mu.Unlock()
		if drained.count == 0 && drained.errors == 0 {
			continue
		}
		summaries = append(summaries, computeSummary(name, drained))
	}
	sortSummaries(summaries)
	return summaries
}

func (o *operationStats) record(latency time.Duration, backend *time.Duration) {
	o.count++
	o.sum += latency
	if latency < o.min {
		o.min = latency
	}
	if latency > o.max {
		o.max = latency
	}
	_ = o.histogram.RecordValue(clampMicros(latency))

	if backend != nil {
		d := *backend
		o.backendCount++
		o.backendSum += d
		if d < o.backendMin {
			o.backendMin = d
		}
		if d > o.backendMax {
			o.backendMax = d
		}
		_ = o.backendHistogram.RecordValue(clampMicros(d))
	}
}

func computeSummary(name string, s operationStats) Summary {
	if s.count == 0 {
		return Summary{Name: name, Errors: s.errors}
	}
	mean := time.Duration(int64(s.sum) / int64(s.count))
	summary := Summary{
		Name:   name,
		Count:  s.count,
		Errors: s.errors,
		Min:    s.min,
		Max:    s.max,
		Mean:   mean,
		P50:    time.Duration(s.histogram.ValueAtQuantile(50)) * time.Microsecond,
		P90:    time.Duration(s.histogram.ValueAtQuantile(90)) * time.Microsecond,
		P99:    time.Duration(s.histogram.ValueAtQuantile(99)) * time.Microsecond,
	}
	if s.backendCount > 0 {
		backendMean := time.Duration(int64(s.backendSum) / int64(s.backendCount))
		summary.BackendMin = durationPtr(s.backendMin)
		summary.BackendMax = durationPtr(s.backendMax)
		summary.BackendMean = durationPtr(backendMean)
		summary.BackendP50 = durationPtr(time.Duration(s.backendHistogram.ValueAtQuantile(50)) * time.Microsecond)
		summary.BackendP90 = durationPtr(time.Duration(s.backendHistogram.ValueAtQuantile(90)) * time.Microsecond)
		summary.BackendP99 = durationPtr(time.Duration(s.backendHistogram.ValueAtQuantile(99)) * time.Microsecond)
	}
	return summary
}

func clampMicros(d time.Duration) int64 {
	us := d.Microseconds()
	if us < 1 {
		return 1
	}
	if us > maxLatencyUS {
		return maxLatencyUS
	}
	return us
}

func durationPtr(d time.Duration) *time.Duration { return &d }

func sortSummaries(summaries []Summary) {
	for i := 1; i < len(summaries); i++ {
		for j := i; j > 0 && summaries[j-1].Name > summaries[j].Name; j-- {
			summaries[j-1], summaries[j] = summaries[j], summaries[j-1]
		}
	}
}

// PrintReport prints operation latency summaries.
func PrintReport(summaries []Summary) {
	if len(summaries) == 0 {
		fmt.Println("  (no operations recorded)")
		return
	}
	fmt.Printf("  %-15s %8s %8s %10s %10s %10s %10s %10s %10s %10s\n", "Operation", "Count", "Errors", "Min", "Max", "Mean", "P50", "P90", "P99", "BackendP99")
	fmt.Printf("  %s\n", strings.Repeat("-", 114))
	for _, s := range summaries {
		backend := "—"
		if s.BackendP99 != nil {
			backend = formatDuration(*s.BackendP99)
		}
		fmt.Printf("  %-15s %8d %8d %10s %10s %10s %10s %10s %10s %10s\n",
			s.Name, s.Count, s.Errors, formatDuration(s.Min), formatDuration(s.Max), formatDuration(s.Mean), formatDuration(s.P50), formatDuration(s.P90), formatDuration(s.P99), backend)
	}
}

// RefreshProcessMetrics captures process and system metrics.
func RefreshProcessMetrics() *ProcessMetrics {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil
	}
	procCPU, _ := proc.CPUPercent()
	memInfo, _ := proc.MemoryInfo()
	systemCPU := 0.0
	if values, err := cpu.Percent(0, false); err == nil && len(values) > 0 {
		systemCPU = values[0]
	}
	vm, _ := mem.VirtualMemory()
	metrics := &ProcessMetrics{CPUPercent: procCPU, SystemCPUPercent: systemCPU, CgroupCPUPercent: readCgroupCPUPercent()}
	if memInfo != nil {
		metrics.MemoryBytes = memInfo.RSS
	}
	if vm != nil {
		metrics.SystemTotalMemoryBytes = vm.Total
		metrics.SystemUsedMemoryBytes = vm.Used
	}
	return metrics
}

// PrintProcessMetrics prints process and system metrics.
func PrintProcessMetrics(metrics *ProcessMetrics) {
	if metrics == nil {
		return
	}
	fmt.Printf("  Process: CPU %.1f%%, Memory %s\n", metrics.CPUPercent, formatBytes(metrics.MemoryBytes))
	fmt.Printf("  System:  CPU %.1f%%, Memory %s/%s\n", metrics.SystemCPUPercent, formatBytes(metrics.SystemUsedMemoryBytes), formatBytes(metrics.SystemTotalMemoryBytes))
	if metrics.CgroupCPUPercent != nil {
		fmt.Printf("  Cgroup:  CPU %.1f%% (kubectl-equivalent)\n", *metrics.CgroupCPUPercent)
	}
}

var cgroupMu sync.Mutex
var previousCgroupUsage *struct {
	usageUsec uint64
	time      time.Time
}

func readCgroupCPUPercent() *float64 {
	stat, err := os.ReadFile("/sys/fs/cgroup/cpu.stat")
	if err != nil {
		return nil
	}
	var usageUsec uint64
	for _, line := range strings.Split(string(stat), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "usage_usec" {
			_, _ = fmt.Sscanf(fields[1], "%d", &usageUsec)
			break
		}
	}
	if usageUsec == 0 {
		return nil
	}
	maxContent, err := os.ReadFile("/sys/fs/cgroup/cpu.max")
	if err != nil {
		return nil
	}
	parts := strings.Fields(string(maxContent))
	if len(parts) == 0 || parts[0] == "max" {
		return nil
	}
	var quotaUsec, periodUsec uint64
	_, _ = fmt.Sscanf(parts[0], "%d", &quotaUsec)
	periodUsec = 100000
	if len(parts) > 1 {
		_, _ = fmt.Sscanf(parts[1], "%d", &periodUsec)
	}
	if quotaUsec == 0 || periodUsec == 0 {
		return nil
	}
	cores := float64(quotaUsec) / float64(periodUsec)
	if cores <= 0 {
		return nil
	}

	now := time.Now()
	cgroupMu.Lock()
	defer cgroupMu.Unlock()
	if previousCgroupUsage == nil {
		previousCgroupUsage = &struct {
			usageUsec uint64
			time      time.Time
		}{usageUsec: usageUsec, time: now}
		return nil
	}
	deltaUsage := usageUsec - previousCgroupUsage.usageUsec
	deltaWall := now.Sub(previousCgroupUsage.time).Seconds() * 1_000_000
	previousCgroupUsage.usageUsec = usageUsec
	previousCgroupUsage.time = now
	if deltaWall <= 0 {
		return nil
	}
	pct := (float64(deltaUsage) / deltaWall / cores) * 100
	return &pct
}

// UpsertResults writes summary documents to Cosmos DB.
func UpsertResults(ctx context.Context, container *azcosmos.ContainerClient, summaries []Summary, metrics *ProcessMetrics, cfg ConfigSnapshot, workloadID, commitSHA, hostname string) {
	now := time.Now().UTC().Format(time.RFC3339)
	cpuPercent, memoryBytes, systemCPUPercent, systemTotalMemoryBytes, systemUsedMemoryBytes := 0.0, uint64(0), 0.0, uint64(0), uint64(0)
	var cgroupCPUPercent *float64
	if metrics != nil {
		cpuPercent = metrics.CPUPercent
		memoryBytes = metrics.MemoryBytes
		systemCPUPercent = metrics.SystemCPUPercent
		systemTotalMemoryBytes = metrics.SystemTotalMemoryBytes
		systemUsedMemoryBytes = metrics.SystemUsedMemoryBytes
		cgroupCPUPercent = metrics.CgroupCPUPercent
	}
	for _, s := range summaries {
		result := PerfResult{
			ID:                      uuid.NewString(),
			PartitionKey:            uuid.NewString(),
			WorkloadID:              workloadID,
			CommitSHA:               commitSHA,
			Hostname:                hostname,
			Timestamp:               now,
			Operation:               s.Name,
			Count:                   s.Count,
			Errors:                  s.Errors,
			MinMS:                   durationMS(s.Min),
			MaxMS:                   durationMS(s.Max),
			MeanMS:                  durationMS(s.Mean),
			P50MS:                   durationMS(s.P50),
			P90MS:                   durationMS(s.P90),
			P99MS:                   durationMS(s.P99),
			BackendMinMS:            durationMSPtr(s.BackendMin),
			BackendMaxMS:            durationMSPtr(s.BackendMax),
			BackendMeanMS:           durationMSPtr(s.BackendMean),
			BackendP50MS:            durationMSPtr(s.BackendP50),
			BackendP90MS:            durationMSPtr(s.BackendP90),
			BackendP99MS:            durationMSPtr(s.BackendP99),
			CPUPercent:              cpuPercent,
			MemoryBytes:             memoryBytes,
			SystemCPUPercent:        systemCPUPercent,
			SystemTotalMemoryBytes:  systemTotalMemoryBytes,
			SystemUsedMemoryBytes:   systemUsedMemoryBytes,
			CgroupCPUPercent:        cgroupCPUPercent,
			ConfigConcurrency:       uint64Ptr(cfg.Concurrency),
			ConfigApplicationRegion: stringPtr(cfg.ApplicationRegion),
			ConfigPreferredRegions:  nonEmptyStringPtr(cfg.PreferredRegions),
			ConfigExcludedRegions:   nonEmptyStringPtr(cfg.ExcludedRegions),
			ConfigPPCBEnabled:       boolPtr(cfg.PPCBEnabled),
			ConfigPyroscopeEnabled:  boolPtr(cfg.PyroscopeEnabled),
			ConfigGOMAXPROCS:        uint64Ptr(cfg.GOMAXPROCS),
			RuntimeGoroutines:       uint64(runtime.NumGoroutine()),
		}
		body, err := json.Marshal(result)
		if err == nil {
			_, err = container.UpsertItem(ctx, azcosmos.NewPartitionKeyString(result.PartitionKey), body, nil)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to upsert perf result: %v\n", err)
		}
	}
}

// UpsertError writes a single error document to Cosmos DB.
func UpsertError(ctx context.Context, container *azcosmos.ContainerClient, operation string, err error, workloadID, commitSHA, hostname string) {
	UpsertErrorWithSource(ctx, container, operation, err, "", workloadID, commitSHA, hostname)
}

// UpsertErrorWithSource writes an error document and lets the caller override the source_message
// (used by the panic recover handler to persist the full goroutine stack trace alongside the panic message).
func UpsertErrorWithSource(ctx context.Context, container *azcosmos.ContainerClient, operation string, err error, sourceOverride string, workloadID, commitSHA, hostname string) {
	src := sourceMessage(err)
	if sourceOverride != "" {
		src = &sourceOverride
	}
	doc := ErrorResult{
		ID:            uuid.NewString(),
		PartitionKey:  uuid.NewString(),
		WorkloadID:    workloadID,
		CommitSHA:     commitSHA,
		Hostname:      hostname,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Operation:     operation,
		ErrorMessage:  err.Error(),
		SourceMessage: src,
	}
	body, marshalErr := json.Marshal(doc)
	if marshalErr == nil {
		_, marshalErr = container.UpsertItem(ctx, azcosmos.NewPartitionKeyString(doc.PartitionKey), body, nil)
	}
	if marshalErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to upsert error result: %v\n", marshalErr)
	}
}

func sourceMessage(err error) *string {
	var parts []string
	for current := errors.Unwrap(err); current != nil; current = errors.Unwrap(current) {
		parts = append(parts, current.Error())
	}
	if len(parts) == 0 {
		return nil
	}
	joined := strings.Join(parts, " → ")
	return &joined
}

func durationMS(d time.Duration) float64 { return d.Seconds() * 1000 }
func durationMSPtr(d *time.Duration) *float64 {
	if d == nil {
		return nil
	}
	ms := durationMS(*d)
	return &ms
}
func uint64Ptr(v uint64) *uint64 { return &v }
func boolPtr(v bool) *bool       { return &v }
func stringPtr(v string) *string { return &v }
func nonEmptyStringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func formatDuration(d time.Duration) string {
	ms := d.Seconds() * 1000
	if ms < 1000 {
		return fmt.Sprintf("%.1fms", ms)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func formatBytes(bytes uint64) string {
	const kb = 1024
	const mb = kb * 1024
	const gb = mb * 1024
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/gb)
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/mb)
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/kb)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

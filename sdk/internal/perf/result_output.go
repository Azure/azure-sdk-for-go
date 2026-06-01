// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type runSummary struct {
	TestName         string  `json:"testName"`
	DurationSeconds  int     `json:"durationSeconds"`
	WarmupSeconds    int     `json:"warmupSeconds"`
	Parallel         int     `json:"parallel"`
	TotalOperations  int64   `json:"totalOperations"`
	OpsPerSecond     float64 `json:"opsPerSecond"`
	SecondsPerOp     float64 `json:"secondsPerOp"`
	WeightedAvgSec   float64 `json:"weightedAverageSeconds"`
	TimestampUTC     string  `json:"timestampUtc"`
	LatencySummary   string  `json:"latencySummary,omitempty"`
	CallTypeSummary  string  `json:"callTypeSummary,omitempty"`
	ResourceSummary  string  `json:"resourceSummary,omitempty"`
	WorkloadConfig   string  `json:"workloadConfig,omitempty"`
	SelectedWorkload string  `json:"selectedWorkload,omitempty"`

	AverageCPUPercent   float64 `json:"averageCpuPercent"`
	AverageMemoryBytes  int64   `json:"averageMemoryBytes"`
	ProcessStatsSummary string  `json:"processStatsSummary,omitempty"`
}

func writeRunArtifacts(prefix string, summary runSummary) error {
	if prefix == "" {
		return nil
	}

	if err := writeJSON(prefix+".json", summary); err != nil {
		return err
	}
	if err := writeCSV(prefix+".csv", summary); err != nil {
		return err
	}
	if err := writeText(prefix+".txt", summary); err != nil {
		return err
	}
	if err := writeMarkdown(prefix+".md", summary); err != nil {
		return err
	}

	return nil
}

func writeJSON(path string, summary runSummary) error {
	b, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON artifact: %w", err)
	}
	return os.WriteFile(path, b, 0o600)
}

func writeCSV(path string, summary runSummary) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create CSV artifact %s: %w", path, err)
	}
	defer func() {
		_ = f.Close()
	}()

	w := csv.NewWriter(f)
	rows := [][]string{
		{"testName", summary.TestName},
		{"durationSeconds", fmt.Sprintf("%d", summary.DurationSeconds)},
		{"warmupSeconds", fmt.Sprintf("%d", summary.WarmupSeconds)},
		{"parallel", fmt.Sprintf("%d", summary.Parallel)},
		{"totalOperations", fmt.Sprintf("%d", summary.TotalOperations)},
		{"opsPerSecond", fmt.Sprintf("%.6f", summary.OpsPerSecond)},
		{"secondsPerOp", fmt.Sprintf("%.6f", summary.SecondsPerOp)},
		{"weightedAverageSeconds", fmt.Sprintf("%.6f", summary.WeightedAvgSec)},
		{"timestampUtc", summary.TimestampUTC},
		{"latencySummary", summary.LatencySummary},
		{"callTypeSummary", summary.CallTypeSummary},
		{"resourceSummary", summary.ResourceSummary},
		{"workloadConfig", summary.WorkloadConfig},
		{"selectedWorkload", summary.SelectedWorkload},
		{"averageCpuPercent", fmt.Sprintf("%.6f", summary.AverageCPUPercent)},
		{"averageMemoryBytes", fmt.Sprintf("%d", summary.AverageMemoryBytes)},
		{"processStatsSummary", summary.ProcessStatsSummary},
	}
	for _, row := range rows {
		if err = w.Write(row); err != nil {
			return fmt.Errorf("failed writing CSV artifact %s: %w", path, err)
		}
	}
	w.Flush()
	if err = w.Error(); err != nil {
		return fmt.Errorf("failed flushing CSV artifact %s: %w", path, err)
	}
	return nil
}

func writeText(path string, summary runSummary) error {
	content := renderText(summary)
	return os.WriteFile(path, []byte(content), 0o600)
}

func writeMarkdown(path string, summary runSummary) error {
	content := renderMarkdown(summary)
	return os.WriteFile(path, []byte(content), 0o600)
}

func renderText(summary runSummary) string {
	lines := []string{
		fmt.Sprintf("Test: %s", summary.TestName),
		fmt.Sprintf("Duration(s): %d", summary.DurationSeconds),
		fmt.Sprintf("Warmup(s): %d", summary.WarmupSeconds),
		fmt.Sprintf("Parallel: %d", summary.Parallel),
		fmt.Sprintf("TotalOperations: %d", summary.TotalOperations),
		fmt.Sprintf("OpsPerSecond: %.6f", summary.OpsPerSecond),
		fmt.Sprintf("SecondsPerOp: %.6f", summary.SecondsPerOp),
		fmt.Sprintf("WeightedAverageSeconds: %.6f", summary.WeightedAvgSec),
		fmt.Sprintf("TimestampUTC: %s", summary.TimestampUTC),
	}
	if summary.LatencySummary != "" {
		lines = append(lines, summary.LatencySummary)
	}
	if summary.CallTypeSummary != "" {
		lines = append(lines, summary.CallTypeSummary)
	}
	if summary.ResourceSummary != "" {
		lines = append(lines, summary.ResourceSummary)
	}
	if summary.WorkloadConfig != "" {
		lines = append(lines, fmt.Sprintf("WorkloadConfig: %s", summary.WorkloadConfig))
	}
	if summary.SelectedWorkload != "" {
		lines = append(lines, fmt.Sprintf("SelectedWorkload: %s", summary.SelectedWorkload))
	}
	lines = append(lines,
		fmt.Sprintf("AverageCpuPercent: %.2f", summary.AverageCPUPercent),
		fmt.Sprintf("AverageMemoryBytes: %d", summary.AverageMemoryBytes),
	)
	if summary.ProcessStatsSummary != "" {
		lines = append(lines, summary.ProcessStatsSummary)
	}
	return strings.Join(lines, "\n") + "\n"
}

func renderMarkdown(summary runSummary) string {
	lines := []string{
		"# Perf Run Summary",
		"",
		"| Field | Value |",
		"| --- | --- |",
		fmt.Sprintf("| Test | %s |", summary.TestName),
		fmt.Sprintf("| Duration (s) | %d |", summary.DurationSeconds),
		fmt.Sprintf("| Warmup (s) | %d |", summary.WarmupSeconds),
		fmt.Sprintf("| Parallel | %d |", summary.Parallel),
		fmt.Sprintf("| Total Operations | %d |", summary.TotalOperations),
		fmt.Sprintf("| Ops/s | %.6f |", summary.OpsPerSecond),
		fmt.Sprintf("| Seconds/op | %.6f |", summary.SecondsPerOp),
		fmt.Sprintf("| Weighted Avg (s) | %.6f |", summary.WeightedAvgSec),
		fmt.Sprintf("| Timestamp (UTC) | %s |", summary.TimestampUTC),
		fmt.Sprintf("| Average CPU (%%) | %.2f |", summary.AverageCPUPercent),
		fmt.Sprintf("| Average Memory (bytes) | %d |", summary.AverageMemoryBytes),
	}
	if summary.WorkloadConfig != "" {
		lines = append(lines, fmt.Sprintf("| Workload Config | %s |", summary.WorkloadConfig))
	}
	if summary.SelectedWorkload != "" {
		lines = append(lines, fmt.Sprintf("| Selected Workload | %s |", summary.SelectedWorkload))
	}
	if summary.LatencySummary != "" {
		lines = append(lines, "", "## Latency", "", summary.LatencySummary)
	}
	if summary.CallTypeSummary != "" {
		lines = append(lines, "", "## By Call Type", "", summary.CallTypeSummary)
	}
	if summary.ResourceSummary != "" {
		lines = append(lines, "", "## Resource Telemetry", "", summary.ResourceSummary)
	}
	if summary.ProcessStatsSummary != "" {
		lines = append(lines, "", "## Process Stats", "", summary.ProcessStatsSummary)
	}
	return strings.Join(lines, "\n") + "\n"
}

func newRunSummary(testName string, totalOperations int64, opsPerSecond float64, latencySummary string, callTypeSummary string, resourceSummary string) runSummary {
	secondsPerOp := 0.0
	weightedAvg := 0.0
	if opsPerSecond > 0 {
		secondsPerOp = 1.0 / opsPerSecond
		weightedAvg = float64(totalOperations) / opsPerSecond
	}

	return runSummary{
		TestName:         testName,
		DurationSeconds:  duration,
		WarmupSeconds:    warmUpDuration,
		Parallel:         parallelInstances,
		TotalOperations:  totalOperations,
		OpsPerSecond:     opsPerSecond,
		SecondsPerOp:     secondsPerOp,
		WeightedAvgSec:   weightedAvg,
		TimestampUTC:     time.Now().UTC().Format(time.RFC3339),
		LatencySummary:   latencySummary,
		CallTypeSummary:  callTypeSummary,
		ResourceSummary:  resourceSummary,
		WorkloadConfig:   workloadConfigPath,
		SelectedWorkload: workloadName,

		AverageCPUPercent:  -1,
		AverageMemoryBytes: -1,
	}
}

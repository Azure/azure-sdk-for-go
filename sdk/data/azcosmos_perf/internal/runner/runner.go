// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package runner executes perf operations with a fixed worker pool.
package runner

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/operations"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/stats"
)

// ConfigSnapshot is a runtime config snapshot emitted with perf results.
type ConfigSnapshot = stats.ConfigSnapshot

// RunConfig configures a perf run.
type RunConfig struct {
	Container        *azcosmos.ContainerClient
	Operations       []operations.Operation
	Stats            *stats.Stats
	Concurrency      int
	Duration         time.Duration
	ReportInterval   time.Duration
	ResultsContainer *azcosmos.ContainerClient
	WorkloadID       string
	CommitSHA        string
	Hostname         string
	ConfigSnapshot   ConfigSnapshot
}

// Run executes operations until ctx is cancelled or Duration elapses.
func Run(ctx context.Context, cfg RunConfig) {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if cfg.Duration > 0 {
		timer := time.AfterFunc(cfg.Duration, func() {
			fmt.Println("\nDuration elapsed, shutting down...")
			cancel()
		})
		defer timer.Stop()
	}

	reporterDone := make(chan struct{})
	go runReporter(runCtx, reporterDone, cfg)

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-runCtx.Done():
					return
				default:
				}

				op := cfg.Operations[rand.Intn(len(cfg.Operations))]
				runIteration(runCtx, op, cfg)
			}
		}()
	}

	wg.Wait()
	cancel()
	<-reporterDone

	totalElapsed := time.Since(start)
	fmt.Printf("\n=== Final Report (total: %.1fs) ===\n", totalElapsed.Seconds())
	metrics := stats.RefreshProcessMetrics()
	stats.PrintProcessMetrics(metrics)
	summaries := cfg.Stats.DrainSummaries()
	stats.PrintReport(summaries)
	upsertCtx, cancelUpsert := context.WithTimeout(context.Background(), 60*time.Second)
	stats.UpsertResults(upsertCtx, cfg.ResultsContainer, summaries, metrics, cfg.ConfigSnapshot, cfg.WorkloadID, cfg.CommitSHA, cfg.Hostname)
	cancelUpsert()
}

func runIteration(ctx context.Context, op operations.Operation, cfg RunConfig) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			fmt.Fprintf(os.Stderr, "panic in operation %q: %v\n%s\n", op.Name(), r, stack)
			cfg.Stats.RecordError(op.Name())
			panicErr := fmt.Errorf("panic: %v", r)
			errCtx, cancelErr := context.WithTimeout(context.Background(), 30*time.Second)
			stats.UpsertErrorWithSource(errCtx, cfg.ResultsContainer, op.Name(), panicErr, string(stack), cfg.WorkloadID, cfg.CommitSHA, cfg.Hostname)
			cancelErr()
		}
	}()

	opStart := time.Now()
	backend, err := op.Execute(ctx, cfg.Container)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		cfg.Stats.RecordError(op.Name())
		errCtx, cancelErr := context.WithTimeout(context.Background(), 30*time.Second)
		stats.UpsertError(errCtx, cfg.ResultsContainer, op.Name(), err, cfg.WorkloadID, cfg.CommitSHA, cfg.Hostname)
		cancelErr()
		return
	}
	cfg.Stats.RecordLatency(op.Name(), time.Since(opStart), backend)
}

func runReporter(ctx context.Context, done chan<- struct{}, cfg RunConfig) {
	defer close(done)
	ticker := time.NewTicker(cfg.ReportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("\n--- Interval Report ---")
			metrics := stats.RefreshProcessMetrics()
			stats.PrintProcessMetrics(metrics)
			summaries := cfg.Stats.DrainSummaries()
			stats.PrintReport(summaries)
			upsertCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			stats.UpsertResults(upsertCtx, cfg.ResultsContainer, summaries, metrics, cfg.ConfigSnapshot, cfg.WorkloadID, cfg.CommitSHA, cfg.Hostname)
			cancel()
		}
	}
}

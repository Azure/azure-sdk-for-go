// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/operations"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/runner"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/setup"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/stats"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		if errors.Is(err, config.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := config.Parse(args)
	if err != nil {
		return err
	}

	if os.Getenv("PYROSCOPE_SERVER_URL") != "" {
		fmt.Fprintln(os.Stderr, "Pyroscope server configured — profiles collected via eBPF auto-instrumentation")
	}

	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	ctx, cancel := context.WithCancel(signalCtx)
	defer cancel()
	normalDone := make(chan struct{})
	defer close(normalDone)
	go func() {
		select {
		case <-signalCtx.Done():
			select {
			case <-normalDone:
				return
			default:
			}
			fmt.Println("\nShutting down…")
			cancel()
		case <-normalDone:
		}
	}()

	client, err := buildClient(cfg.Endpoint, cfg.Auth, cfg.Key, cfg.PreferredRegionList())
	if err != nil {
		return err
	}
	db, err := client.NewDatabase(cfg.Database)
	if err != nil {
		return err
	}
	if err := setup.EnsureDatabase(ctx, client, cfg.Database); err != nil {
		return err
	}

	var defaultTTL *int32
	if cfg.DefaultTTL > 0 {
		ttl := cfg.DefaultTTL
		defaultTTL = &ttl
	}
	container, err := setup.EnsureContainer(ctx, db, cfg.Container, cfg.Throughput, defaultTTL)
	if err != nil {
		return err
	}

	seeded, err := seed.SeedContainer(ctx, container, cfg.SeedCount, cfg.Concurrency)
	if err != nil {
		return err
	}
	sharedItems := seed.NewSharedItems(seeded)

	ops := operations.CreateOperations(cfg, sharedItems)
	fmt.Printf("\nStarting perf test: %d operation(s), concurrency=%d\n", len(ops), cfg.Concurrency)
	for _, op := range ops {
		fmt.Printf("  - %s\n", op.Name())
	}
	if cfg.Duration > 0 {
		fmt.Printf("  Duration: %ds\n", cfg.Duration)
	} else {
		fmt.Println("  Duration: indefinite (Ctrl+C to stop)")
	}
	fmt.Println()

	commitSHA := cfg.CommitSHA
	if commitSHA == "" {
		commitSHA = resolveCommitSHA()
	}
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "unknown"
	}

	resultsContainer, err := buildResultsContainer(ctx, cfg, db)
	if err != nil {
		return err
	}

	configSnapshot := runner.ConfigSnapshot{
		Concurrency:       uint64(cfg.Concurrency),
		ApplicationRegion: cfg.ApplicationRegion,
		PreferredRegions:  strings.Join(cfg.PreferredRegionList(), ", "),
		ExcludedRegions:   strings.Join(cfg.ExcludedRegions, ", "),
		GOMAXPROCS:        uint64(runtime.GOMAXPROCS(0)),
		PPCBEnabled:       envBoolDefault("AZURE_COSMOS_PER_PARTITION_CIRCUIT_BREAKER_ENABLED", true),
		PyroscopeEnabled:  os.Getenv("PYROSCOPE_SERVER_URL") != "",
	}

	opNames := make([]string, len(ops))
	for i, op := range ops {
		opNames[i] = op.Name()
	}
	runStats := stats.New(opNames)
	duration := time.Duration(0)
	if cfg.Duration > 0 {
		duration = time.Duration(cfg.Duration) * time.Second
	}

	runner.Run(ctx, runner.RunConfig{
		Container:        container,
		Operations:       ops,
		Stats:            runStats,
		Concurrency:      cfg.Concurrency,
		Duration:         duration,
		ReportInterval:   time.Duration(cfg.ReportInterval) * time.Second,
		ResultsContainer: resultsContainer,
		WorkloadID:       cfg.WorkloadID,
		CommitSHA:        commitSHA,
		Hostname:         hostname,
		ConfigSnapshot:   configSnapshot,
	})
	return nil
}

func buildResultsContainer(ctx context.Context, cfg config.Config, workloadDB *azcosmos.DatabaseClient) (*azcosmos.ContainerClient, error) {
	if cfg.ResultsEndpoint == "" {
		container, err := setup.EnsureContainer(ctx, workloadDB, cfg.ResultsContainer, 10000, int32Ptr(86400))
		if err != nil {
			return nil, err
		}
		fmt.Printf("Perf results will be stored in container '%s'. Workload ID: %s\n", cfg.ResultsContainer, cfg.WorkloadID)
		return container, nil
	}

	resultsAuth := cfg.Auth
	if cfg.ResultsAuth != nil {
		resultsAuth = *cfg.ResultsAuth
	}
	resultsClient, err := buildClient(cfg.ResultsEndpoint, resultsAuth, cfg.ResultsKey, cfg.PreferredRegionList())
	if err != nil {
		return nil, err
	}
	if err := setup.EnsureDatabase(ctx, resultsClient, cfg.ResultsDatabase); err != nil {
		return nil, err
	}
	resultsDB, err := resultsClient.NewDatabase(cfg.ResultsDatabase)
	if err != nil {
		return nil, err
	}
	container, err := setup.EnsureContainer(ctx, resultsDB, cfg.ResultsContainer, 10000, int32Ptr(86400))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Perf results will be stored on separate account '%s' in '%s/%s'. Workload ID: %s\n", cfg.ResultsEndpoint, cfg.ResultsDatabase, cfg.ResultsContainer, cfg.WorkloadID)
	return container, nil
}

func buildClient(endpoint string, auth config.AuthMethod, key string, preferredRegions []string) (*azcosmos.Client, error) {
	options := &azcosmos.ClientOptions{PreferredRegions: preferredRegions}
	options.PerCallPolicies = append(options.PerCallPolicies, operations.NewPipelinePolicy())
	switch auth {
	case config.AuthKey:
		cred, err := azcosmos.NewKeyCredential(key)
		if err != nil {
			return nil, err
		}
		return azcosmos.NewClientWithKey(endpoint, cred, options)
	case config.AuthAAD:
		cred, err := createAADCredential()
		if err != nil {
			return nil, err
		}
		return azcosmos.NewClient(endpoint, cred, options)
	default:
		return nil, fmt.Errorf("unsupported auth method %q", auth)
	}
}

func createAADCredential() (azcore.TokenCredential, error) {
	workload, workloadErr := azidentity.NewWorkloadIdentityCredential(nil)
	if workloadErr == nil {
		return workload, nil
	}
	managed, managedErr := azidentity.NewManagedIdentityCredential(nil)
	if managedErr == nil {
		return managed, nil
	}
	return nil, fmt.Errorf("failed to create AAD credential. WorkloadIdentityCredential: %v, ManagedIdentityCredential: %v", workloadErr, managedErr)
}

func resolveCommitSHA() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	sha := strings.TrimSpace(string(out))
	if sha == "" {
		return "unknown"
	}
	return sha
}

func envBoolDefault(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func int32Ptr(v int32) *int32 { return &v }

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package config parses and validates cosmos-perf CLI configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// AuthMethod identifies how the perf tool authenticates to Cosmos DB.
type AuthMethod string

const (
	AuthKey AuthMethod = "key"
	AuthAAD AuthMethod = "aad"
)

// ExcludeRegionsScope controls which operation types receive excluded regions.
type ExcludeRegionsScope string

const (
	ExcludeRegionsReads  ExcludeRegionsScope = "reads"
	ExcludeRegionsWrites ExcludeRegionsScope = "writes"
	ExcludeRegionsBoth   ExcludeRegionsScope = "both"
)

// ErrHelp indicates cobra handled --help and no config was parsed.
var ErrHelp = errors.New("help requested")

// Config is the parsed command-line configuration.
type Config struct {
	Endpoint          string
	Database          string
	Container         string
	Auth              AuthMethod
	Key               string
	ApplicationRegion string
	PreferredRegions  []string
	ExcludedRegions   []string
	ExcludeRegionsFor ExcludeRegionsScope

	NoReads      bool
	NoQueries    bool
	NoUpserts    bool
	NoCreates    bool
	NoReadMany   bool
	NoChangeFeed bool

	ReadManyBatchSize  int
	ChangeFeedMaxItems int32

	Concurrency    int
	Duration       int64
	SeedCount      int
	ReportInterval int64
	Throughput     int32
	DefaultTTL     int32

	ResultsContainer string
	ResultsEndpoint  string
	ResultsDatabase  string
	ResultsAuth      *AuthMethod
	ResultsKey       string

	WorkloadID string
	CommitSHA  string
}

// Parse parses CLI arguments using cobra.
func Parse(args []string) (Config, error) {
	cfg := Config{
		Database:           "perfdb",
		Container:          "perfcontainer",
		ExcludeRegionsFor:  ExcludeRegionsBoth,
		Concurrency:        50,
		SeedCount:          1000,
		ReportInterval:     300,
		Throughput:         100000,
		DefaultTTL:         3600,
		ResultsContainer:   "perfresults",
		ResultsDatabase:    "perfdb",
		WorkloadID:         uuid.NewString(),
		ReadManyBatchSize:  20,
		ChangeFeedMaxItems: 100,
	}
	var auth string
	var resultsAuth string

	executed := false
	cmd := &cobra.Command{
		Use:           "cosmos-perf",
		Short:         "Azure Cosmos DB Go SDK performance testing tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			executed = true
			cfg.Key = firstNonEmpty(cfg.Key, os.Getenv("AZURE_COSMOS_KEY"))
			cfg.ResultsKey = firstNonEmpty(cfg.ResultsKey, os.Getenv("AZURE_COSMOS_RESULTS_KEY"))

			if auth != "" {
				parsed, err := parseAuth(auth, "--auth")
				if err != nil {
					return err
				}
				cfg.Auth = parsed
			}
			if resultsAuth != "" {
				parsed, err := parseAuth(resultsAuth, "--results-auth")
				if err != nil {
					return err
				}
				cfg.ResultsAuth = &parsed
			}
			return cfg.Validate()
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&cfg.Endpoint, "endpoint", "", "Cosmos DB account endpoint URL")
	flags.StringVar(&cfg.Database, "database", cfg.Database, "Database name")
	flags.StringVar(&cfg.Container, "container", cfg.Container, "Container name")
	flags.StringVar(&auth, "auth", "", "Authentication method: key or aad")
	flags.StringVar(&cfg.Key, "key", "", "Account key (or AZURE_COSMOS_KEY)")
	flags.StringVar(&cfg.ApplicationRegion, "application-region", "", "Azure region where the application is running")
	flags.StringSliceVar(&cfg.PreferredRegions, "preferred-regions", nil, "Additional comma-separated preferred regions")
	flags.StringSliceVar(&cfg.ExcludedRegions, "excluded-regions", nil, "Comma-separated excluded regions")
	flags.Var((*excludeRegionsScopeValue)(&cfg.ExcludeRegionsFor), "exclude-regions-for", "Scope for excluded regions: reads, writes, or both")

	flags.BoolVar(&cfg.NoReads, "no-reads", false, "Disable point read operations")
	flags.BoolVar(&cfg.NoQueries, "no-queries", false, "Disable query operations")
	flags.BoolVar(&cfg.NoUpserts, "no-upserts", false, "Disable upsert operations")
	flags.BoolVar(&cfg.NoCreates, "no-creates", false, "Disable create operations")
	flags.BoolVar(&cfg.NoReadMany, "no-read-many", false, "Disable read-many operations")
	flags.BoolVar(&cfg.NoChangeFeed, "no-change-feed", false, "Disable change feed operations")
	flags.IntVar(&cfg.ReadManyBatchSize, "read-many-batch-size", cfg.ReadManyBatchSize, "Items per ReadManyItems call")
	flags.Int32Var(&cfg.ChangeFeedMaxItems, "change-feed-max-items", cfg.ChangeFeedMaxItems, "MaxItemCount for change feed")

	flags.IntVar(&cfg.Concurrency, "concurrency", cfg.Concurrency, "Number of concurrent operations")
	flags.Int64Var(&cfg.Duration, "duration", 0, "Run duration in seconds; omit for indefinite")
	flags.IntVar(&cfg.SeedCount, "seed-count", cfg.SeedCount, "Number of items to pre-seed")
	flags.Int64Var(&cfg.ReportInterval, "report-interval", cfg.ReportInterval, "Stats reporting interval in seconds")
	flags.Int32Var(&cfg.Throughput, "throughput", cfg.Throughput, "Throughput (RU/s) when creating the container")
	flags.Int32Var(&cfg.DefaultTTL, "default-ttl", cfg.DefaultTTL, "Default TTL in seconds for items; 0 disables TTL")

	flags.StringVar(&cfg.ResultsContainer, "results-container", cfg.ResultsContainer, "Container for storing perf results")
	flags.StringVar(&cfg.ResultsEndpoint, "results-endpoint", "", "Cosmos DB endpoint for results account")
	flags.StringVar(&cfg.ResultsDatabase, "results-database", cfg.ResultsDatabase, "Database name on the results account")
	flags.StringVar(&resultsAuth, "results-auth", "", "Authentication method for results account: key or aad")
	flags.StringVar(&cfg.ResultsKey, "results-key", "", "Results account key (or AZURE_COSMOS_RESULTS_KEY)")
	flags.StringVar(&cfg.WorkloadID, "workload-id", cfg.WorkloadID, "Unique identifier for this workload instance")
	flags.StringVar(&cfg.CommitSHA, "commit-sha", "", "Git commit SHA stamped on result documents")

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		return Config{}, err
	}
	if !executed {
		return Config{}, ErrHelp
	}
	return cfg, nil
}

// Validate validates semantic configuration constraints.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Endpoint) == "" {
		return fmt.Errorf("--endpoint is required")
	}
	if c.Auth == "" {
		return fmt.Errorf("--auth is required")
	}
	if strings.TrimSpace(c.ApplicationRegion) == "" {
		return fmt.Errorf("--application-region is required")
	}
	if c.Auth == AuthKey && strings.TrimSpace(c.Key) == "" {
		return fmt.Errorf("account key is required for key auth. Use --key or set AZURE_COSMOS_KEY env var")
	}
	if c.NoReads && c.NoQueries && c.NoUpserts && c.NoCreates && c.NoReadMany && c.NoChangeFeed {
		return fmt.Errorf("all operations are disabled. Enable at least one")
	}
	if c.Concurrency < 1 {
		return fmt.Errorf("--concurrency must be at least 1")
	}
	if c.SeedCount < 1 {
		return fmt.Errorf("--seed-count must be at least 1")
	}
	if c.ReadManyBatchSize < 1 {
		return fmt.Errorf("--read-many-batch-size must be at least 1")
	}
	if c.ChangeFeedMaxItems < 1 {
		return fmt.Errorf("--change-feed-max-items must be at least 1")
	}
	if c.Duration < 0 {
		return fmt.Errorf("--duration cannot be negative")
	}
	if c.ReportInterval < 1 {
		return fmt.Errorf("--report-interval must be at least 1")
	}
	if c.Throughput < 1 {
		return fmt.Errorf("--throughput must be at least 1")
	}
	if c.DefaultTTL < 0 {
		return fmt.Errorf("--default-ttl cannot be negative")
	}
	if c.ResultsEndpoint != "" {
		resultsAuth := c.Auth
		if c.ResultsAuth != nil {
			resultsAuth = *c.ResultsAuth
		}
		if resultsAuth == AuthKey && strings.TrimSpace(c.ResultsKey) == "" {
			return fmt.Errorf("results account key is required. Use --results-key or set AZURE_COSMOS_RESULTS_KEY")
		}
	}
	return nil
}

// PreferredRegionList returns application region followed by additional preferred regions.
func (c Config) PreferredRegionList() []string {
	seen := map[string]struct{}{}
	regions := make([]string, 0, 1+len(c.PreferredRegions))
	for _, region := range append([]string{c.ApplicationRegion}, c.PreferredRegions...) {
		region = strings.TrimSpace(region)
		if region == "" {
			continue
		}
		key := strings.ToLower(region)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		regions = append(regions, region)
	}
	return regions
}

func parseAuth(value, flag string) (AuthMethod, error) {
	switch AuthMethod(strings.ToLower(strings.TrimSpace(value))) {
	case AuthKey:
		return AuthKey, nil
	case AuthAAD:
		return AuthAAD, nil
	default:
		return "", fmt.Errorf("%s must be key or aad", flag)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

type excludeRegionsScopeValue ExcludeRegionsScope

func (v *excludeRegionsScopeValue) String() string { return string(*v) }
func (v *excludeRegionsScopeValue) Type() string   { return "scope" }
func (v *excludeRegionsScopeValue) Set(s string) error {
	switch ExcludeRegionsScope(strings.ToLower(strings.TrimSpace(s))) {
	case ExcludeRegionsReads:
		*v = excludeRegionsScopeValue(ExcludeRegionsReads)
	case ExcludeRegionsWrites:
		*v = excludeRegionsScopeValue(ExcludeRegionsWrites)
	case ExcludeRegionsBoth:
		*v = excludeRegionsScopeValue(ExcludeRegionsBoth)
	default:
		return fmt.Errorf("must be reads, writes, or both")
	}
	return nil
}

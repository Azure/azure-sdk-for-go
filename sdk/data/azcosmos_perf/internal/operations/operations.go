// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package operations contains executable perf operations.
package operations

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos_perf/internal/seed"
)

const requestDurationHeader = "x-ms-request-duration-ms"
const excludedRegionsHeader = "x-ms-excluded-regions"

// PerfItem is the item type used for perf operations.
type PerfItem = seed.PerfItem

// Operation is a single executable perf test operation.
type Operation interface {
	Name() string
	Execute(ctx context.Context, c *azcosmos.ContainerClient) (*time.Duration, error)
}

// CreateOperations creates enabled operations from CLI configuration.
func CreateOperations(cfg config.Config, items *seed.SharedItems) []Operation {
	readRegions, writeRegions := regionScopes(cfg)
	ops := make([]Operation, 0, 6)
	if !cfg.NoReads {
		ops = append(ops, NewReadItemOperation(items, readRegions))
	}
	if !cfg.NoQueries {
		ops = append(ops, NewQueryItemsOperation(items, readRegions))
	}
	if !cfg.NoReadMany {
		ops = append(ops, NewReadManyItemsOperation(items, cfg.ReadManyBatchSize, readRegions))
	}
	if !cfg.NoUpserts {
		ops = append(ops, NewUpsertItemOperation(items, writeRegions))
	}
	if !cfg.NoCreates {
		ops = append(ops, NewCreateItemOperation(items, writeRegions))
	}
	if !cfg.NoChangeFeed {
		ops = append(ops, NewChangeFeedItemsOperation(cfg.ChangeFeedMaxItems, readRegions))
	}
	return ops
}

func regionScopes(cfg config.Config) (readRegions, writeRegions []string) {
	if len(cfg.ExcludedRegions) == 0 {
		return nil, nil
	}
	regions := cleanRegions(cfg.ExcludedRegions)
	switch cfg.ExcludeRegionsFor {
	case config.ExcludeRegionsReads:
		readRegions = regions
	case config.ExcludeRegionsWrites:
		writeRegions = regions
	default:
		readRegions = regions
		writeRegions = regions
	}
	return readRegions, writeRegions
}

func cleanRegions(regions []string) []string {
	result := make([]string, 0, len(regions))
	for _, region := range regions {
		region = strings.TrimSpace(region)
		if region != "" {
			result = append(result, region)
		}
	}
	return result
}

type backendCollectorKey struct{}
type excludedRegionsKey struct{}

type backendCollector struct {
	mu    sync.Mutex
	total time.Duration
	count int
}

func newBackendCollector() *backendCollector { return &backendCollector{} }
func (c *backendCollector) record(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.total += d
	c.count++
}
func (c *backendCollector) duration() *time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.count == 0 {
		return nil
	}
	// Match the Rust perf tool: multi-page operations report total backend processing time.
	d := c.total
	return &d
}

func withBackendCollector(ctx context.Context, collector *backendCollector) context.Context {
	return context.WithValue(ctx, backendCollectorKey{}, collector)
}

func withExcludedRegions(ctx context.Context, regions []string) context.Context {
	if len(regions) == 0 {
		return ctx
	}
	return context.WithValue(ctx, excludedRegionsKey{}, regions)
}

func prepareContext(ctx context.Context, regions []string) (context.Context, *backendCollector) {
	collector := newBackendCollector()
	ctx = withBackendCollector(ctx, collector)
	ctx = withExcludedRegions(ctx, regions)
	return ctx, collector
}

// PipelinePolicy captures backend durations and applies per-operation excluded regions.
type PipelinePolicy struct{}

// NewPipelinePolicy returns the policy required by this perf package.
func NewPipelinePolicy() policy.Policy { return PipelinePolicy{} }

// Do implements policy.Policy.
func (PipelinePolicy) Do(req *policy.Request) (*http.Response, error) {
	if regions, ok := req.Raw().Context().Value(excludedRegionsKey{}).([]string); ok && len(regions) > 0 {
		req.Raw().Header.Set(excludedRegionsHeader, strings.Join(regions, ","))
	}

	resp, err := req.Next()
	if resp != nil {
		recordBackendDuration(req.Raw().Context(), resp)
	}
	if err != nil {
		var responseErr *azcore.ResponseError
		if errors.As(err, &responseErr) && responseErr.RawResponse != nil {
			recordBackendDuration(req.Raw().Context(), responseErr.RawResponse)
		}
	}
	return resp, err
}

func recordBackendDuration(ctx context.Context, resp *http.Response) {
	collector, ok := ctx.Value(backendCollectorKey{}).(*backendCollector)
	if !ok || collector == nil {
		return
	}
	if d := ExtractBackendDuration(resp); d != nil {
		collector.record(*d)
	}
}

// ExtractBackendDuration parses x-ms-request-duration-ms from a raw response.
func ExtractBackendDuration(resp *http.Response) *time.Duration {
	if resp == nil {
		return nil
	}
	value := resp.Header.Get(requestDurationHeader)
	if value == "" {
		return nil
	}
	ms, err := strconv.ParseFloat(value, 64)
	if err != nil || !isFiniteNonNegative(ms) {
		return nil
	}
	d := time.Duration(ms * float64(time.Millisecond))
	return &d
}

func isFiniteNonNegative(v float64) bool {
	return v >= 0 && !math.IsInf(v, 0) && !math.IsNaN(v)
}

func marshalItem(item PerfItem) ([]byte, error) {
	return json.Marshal(item)
}

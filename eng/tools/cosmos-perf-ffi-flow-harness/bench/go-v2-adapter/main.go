// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

// cSpell:words colls documentdb

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"runtime/pprof"
	runtimeTrace "runtime/trace"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type options struct {
	endpoint  string
	database  string
	container string
	owner     string
	key       string
	keyFD     string
	workers   int
	warmup    time.Duration
	measure   time.Duration
}

type adapterResult struct {
	Workload            map[string]any `json:"workload"`
	TransportProof      map[string]any `json:"transportProof"`
	Counts              map[string]any `json:"counts"`
	Timing              map[string]any `json:"timing"`
	LatencyMilliseconds map[string]any `json:"latencyMilliseconds"`
	QueueDelay          any            `json:"queueDelayMilliseconds"`
	Gauges              map[string]any `json:"gauges"`
}

func main() {
	if err := runMain(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runMain(arguments []string) error {
	if err := rejectSecretArguments(arguments); err != nil {
		return err
	}
	flags := flag.NewFlagSet("cosmos-perf-go", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var value options
	flags.StringVar(&value.endpoint, "endpoint", "", "Cosmos account endpoint")
	flags.StringVar(&value.database, "database", "", "database name")
	flags.StringVar(&value.container, "container", "", "container name")
	flags.StringVar(&value.owner, "owner", "", "campaign owner")
	flags.StringVar(&value.keyFD, "key-fd", os.Getenv("COSMOS_PERF_KEY_FD"), "inherited account-key descriptor")
	flags.IntVar(&value.workers, "workers", envInt("COSMOS_PERF_WORKERS", 0), "worker count")
	warmupSeconds := flags.Int("warmup-seconds", 30, "excluded warmup duration")
	measurementSeconds := flags.Int("measurement-seconds", 300, "measurement duration")
	if err := flags.Parse(arguments); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("exactly one action is required: seed, verify, preflight, run, or cleanup")
	}
	if value.endpoint == "" || value.database == "" || value.container == "" {
		return errors.New("--endpoint, --database, and --container are required")
	}
	if err := validateOwner(value.owner); err != nil {
		return err
	}
	value.warmup = time.Duration(*warmupSeconds) * time.Second
	value.measure = time.Duration(*measurementSeconds) * time.Second

	keyBytes, err := readSecretFD(value.keyFD)
	if err != nil {
		return err
	}
	value.key = string(keyBytes)
	clearBytes(keyBytes)
	defer func() { value.key = "" }()

	ctx := context.Background()
	client, err := newSDKClient(ctx, value)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}
	defer func() { _ = client.Close() }()

	switch flags.Arg(0) {
	case "seed":
		return seed(ctx, client, value)
	case "verify":
		return verify(ctx, client, value)
	case "preflight":
		proof, err := preflight(ctx, client, value)
		if err != nil {
			return err
		}
		return writeJSON(os.Stdout, proof)
	case "run":
		return runBenchmark(ctx, client, value)
	case "cleanup":
		return cleanup(ctx, client, value)
	default:
		return fmt.Errorf("unknown action %q", flags.Arg(0))
	}
}

func seed(ctx context.Context, client sdkClient, value options) error {
	for index := 0; index < seedCount; index++ {
		body, err := makeDocument(value.owner, index)
		if err != nil {
			return err
		}
		existing, err := client.Read(ctx, itemPK(value.owner, index), itemID(value.owner, index))
		switch {
		case err == nil:
			owned, ownerErr := ownedBy(existing.Body, value.owner)
			if ownerErr != nil || !owned {
				return fmt.Errorf("refusing to overwrite non-owned item %q", itemID(value.owner, index))
			}
			_, err = client.Upsert(ctx, itemPK(value.owner, index), itemID(value.owner, index), body)
		case isNotFound(err):
			_, err = client.Create(ctx, itemPK(value.owner, index), itemID(value.owner, index), body)
		}
		if err != nil {
			return fmt.Errorf("seed item %d: %w", index, err)
		}
	}
	return writeJSON(os.Stdout, map[string]any{"action": "seed", "owner": value.owner, "count": seedCount})
}

func verify(ctx context.Context, client sdkClient, value options) error {
	for index := 0; index < seedCount; index++ {
		response, err := client.Read(ctx, itemPK(value.owner, index), itemID(value.owner, index))
		if err != nil {
			return fmt.Errorf("verify item %d: %w", index, err)
		}
		if err := validateDocument(response.Body, value.owner, index); err != nil {
			return fmt.Errorf("verify item %d: %w", index, err)
		}
	}
	return writeJSON(os.Stdout, map[string]any{"action": "verify", "owner": value.owner, "count": seedCount})
}

func cleanup(ctx context.Context, client sdkClient, value options) error {
	deleted := 0
	skipped := 0
	for index := 0; index < seedCount; index++ {
		response, err := client.Read(ctx, itemPK(value.owner, index), itemID(value.owner, index))
		if isNotFound(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read before cleanup item %d: %w", index, err)
		}
		owned, err := ownedBy(response.Body, value.owner)
		if err != nil || !owned {
			skipped++
			continue
		}
		if _, err := client.Delete(ctx, itemPK(value.owner, index), itemID(value.owner, index)); err != nil {
			return fmt.Errorf("delete item %d: %w", index, err)
		}
		deleted++
	}
	return writeJSON(os.Stdout, map[string]any{
		"action": "cleanup", "owner": value.owner, "deleted": deleted, "nonOwnedSkipped": skipped,
	})
}

func preflight(ctx context.Context, client sdkClient, value options) (map[string]any, error) {
	_, err := client.Read(ctx, itemPK(value.owner, 0), itemID(value.owner, 0))
	if err != nil {
		return nil, fmt.Errorf("SDK preflight read: %w", err)
	}
	probe, err := signedHTTPProbe(ctx, value)
	if err != nil {
		return nil, err
	}
	if probe.TransportKind != "gateway-v1" || probe.HTTPVersion != "HTTP/2" {
		return nil, fmt.Errorf("transport proof failed: kind=%q HTTP=%q", probe.TransportKind, probe.HTTPVersion)
	}
	return map[string]any{
		"transportKind":         probe.TransportKind,
		"negotiatedHTTPVersion": probe.HTTPVersion,
		"binaryEncoding":        false,
	}, nil
}

func signedHTTPProbe(ctx context.Context, value options) (operationResponse, error) {
	resourceLink := fmt.Sprintf("dbs/%s/colls/%s/docs/%s", value.database, value.container, itemID(value.owner, 0))
	date := time.Now().UTC().Format(http.TimeFormat)
	payload := strings.ToLower("get") + "\n" + "docs\n" + resourceLink + "\n" + strings.ToLower(date) + "\n\n"
	decodedKey, err := base64.StdEncoding.DecodeString(value.key)
	if err != nil {
		return operationResponse{}, errors.New("account key is not valid base64")
	}
	defer clearBytes(decodedKey)
	mac := hmac.New(sha256.New, decodedKey)
	_, _ = mac.Write([]byte(payload))
	token := url.QueryEscape("type=master&ver=1.0&sig=" + base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	requestURL := strings.TrimRight(value.endpoint, "/") + "/" + resourceLink
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return operationResponse{}, err
	}
	request.Header.Set("Authorization", token)
	request.Header.Set("x-ms-date", date)
	request.Header.Set("x-ms-version", "2020-11-05")
	request.Header.Set("x-ms-documentdb-partitionkey", fmt.Sprintf("[\"%s\"]", itemPK(value.owner, 0)))
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ForceAttemptHTTP2 = true
	response, err := (&http.Client{Transport: transport}).Do(request)
	if err != nil {
		return operationResponse{}, fmt.Errorf("transport preflight request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	_, _ = io.Copy(io.Discard, response.Body)
	if response.StatusCode != http.StatusOK {
		return operationResponse{}, fmt.Errorf("transport preflight returned HTTP %d", response.StatusCode)
	}
	kind := ""
	if response.Header.Get("x-ms-gatewayversion") != "" {
		kind = "gateway-v1"
	}
	return operationResponse{
		StatusCode:    response.StatusCode,
		TransportKind: kind,
		HTTPVersion:   normalizeHTTPVersion(response.Proto),
	}, nil
}

func runBenchmark(ctx context.Context, client sdkClient, value options) error {
	if value.workers != 8 && value.workers != 16 && value.workers != 32 && value.workers != 64 {
		return errors.New("--workers must be 8, 16, 32, or 64")
	}
	proof, err := preflight(ctx, client, value)
	if err != nil {
		return err
	}
	if err := startDiagnostics(); err != nil {
		return err
	}
	defer stopDiagnostics()

	if err := runPhase(ctx, client, value, value.warmup, false, nil); err != nil {
		return fmt.Errorf("warmup: %w", err)
	}
	cpuStart := processCPUSeconds()
	started := time.Now()
	var state benchmarkState
	if err := runPhase(ctx, client, value, value.measure, true, &state); err != nil {
		return fmt.Errorf("measurement: %w", err)
	}
	wallSeconds := time.Since(started).Seconds()
	cpuSeconds := processCPUSeconds() - cpuStart
	if state.errors.Load() != 0 {
		return fmt.Errorf("measurement completed with %d errors", state.errors.Load())
	}
	latencies := state.mergedLatencies()
	result := adapterResult{
		Workload: map[string]any{
			"partitionKeyPath": "/pk", "seedCount": seedCount, "documentBytes": documentBytes,
			"operation": "point-read", "workers": value.workers,
			"warmupSeconds": int(value.warmup.Seconds()), "measurementSeconds": int(value.measure.Seconds()),
			"closedLoop": true, "keySequence": "round-robin-0-through-1023",
			"rawResponseMode": true, "validationOutsideTimedPath": true,
			"gatewayProtocol": "GatewayV1HTTP2", "binaryEncoding": false, "authentication": "account-key",
		},
		TransportProof: proof,
		Counts: map[string]any{
			"success": state.success.Load(), "error": state.errors.Load(),
			"retry": state.retries.Load(), "throttle": state.throttles.Load(),
		},
		Timing:              map[string]any{"wallSeconds": wallSeconds, "processCPUSeconds": cpuSeconds},
		LatencyMilliseconds: percentiles(latencies),
		QueueDelay:          nil,
		Gauges: map[string]any{
			"queueHighWater": nil, "inFlightHighWater": value.workers, "completionHighWater": nil,
		},
	}
	output := os.Getenv("COSMOS_PERF_ADAPTER_RESULT")
	if output == "" {
		return errors.New("COSMOS_PERF_ADAPTER_RESULT is required for run")
	}
	file, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	return writeJSON(file, result)
}

type benchmarkState struct {
	sequence  atomic.Uint64
	success   atomic.Uint64
	errors    atomic.Uint64
	retries   atomic.Uint64
	throttles atomic.Uint64
	latencies [][]int64
}

func (state *benchmarkState) mergedLatencies() []int64 {
	var count int
	for _, values := range state.latencies {
		count += len(values)
	}
	result := make([]int64, 0, count)
	for _, values := range state.latencies {
		result = append(result, values...)
	}
	return result
}

func runPhase(ctx context.Context, client sdkClient, value options, duration time.Duration, record bool, state *benchmarkState) error {
	if duration <= 0 {
		return errors.New("phase duration must be positive")
	}
	if state == nil {
		state = &benchmarkState{}
	}
	state.latencies = make([][]int64, value.workers)
	deadline := time.Now().Add(duration)
	var wait sync.WaitGroup
	var firstErr atomic.Pointer[error]
	for worker := 0; worker < value.workers; worker++ {
		wait.Add(1)
		go func(worker int) {
			defer wait.Done()
			latencies := make([]int64, 0, 4096)
			for time.Now().Before(deadline) && firstErr.Load() == nil {
				index := keyIndex(state.sequence.Add(1) - 1)
				started := time.Now()
				response, err := client.Read(ctx, itemPK(value.owner, index), itemID(value.owner, index))
				elapsed := time.Since(started)
				if err != nil {
					state.errors.Add(1)
					captured := error(err)
					firstErr.CompareAndSwap(nil, &captured)
					break
				}
				if record {
					latencies = append(latencies, elapsed.Nanoseconds())
				}
				if response.Attempts > 1 {
					state.retries.Add(uint64(response.Attempts - 1))
				}
				if response.StatusCode == http.StatusTooManyRequests {
					state.throttles.Add(1)
				}
				state.success.Add(1)
			}
			state.latencies[worker] = latencies
		}(worker)
	}
	wait.Wait()
	if captured := firstErr.Load(); captured != nil {
		return *captured
	}
	return nil
}

func percentiles(values []int64) map[string]any {
	sort.Slice(values, func(left, right int) bool { return values[left] < values[right] })
	result := map[string]any{}
	for name, quantile := range map[string]float64{"p50": .5, "p90": .9, "p95": .95, "p99": .99, "p999": .999} {
		if len(values) == 0 {
			result[name] = float64(0)
			continue
		}
		index := int(float64(len(values)-1) * quantile)
		result[name] = float64(values[index]) / float64(time.Millisecond)
	}
	return result
}

var (
	cpuProfileFile   *os.File
	traceProfileFile *os.File
)

func startDiagnostics() error {
	if path := os.Getenv("COSMOS_PERF_BLOCK_PROFILE"); path != "" {
		runtime.SetBlockProfileRate(1)
	}
	if path := os.Getenv("COSMOS_PERF_MUTEX_PROFILE"); path != "" {
		runtime.SetMutexProfileFraction(1)
	}
	if path := os.Getenv("COSMOS_PERF_CPU_PROFILE"); path != "" {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		cpuProfileFile = file
		if err := pprof.StartCPUProfile(file); err != nil {
			return err
		}
	}
	if path := os.Getenv("COSMOS_PERF_RUNTIME_TRACE"); path != "" {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		traceProfileFile = file
		if err := runtimeTrace.Start(file); err != nil {
			return err
		}
		time.AfterFunc(10*time.Second, func() {
			runtimeTrace.Stop()
			_ = traceProfileFile.Close()
			traceProfileFile = nil
		})
	}
	return nil
}

func stopDiagnostics() {
	if cpuProfileFile != nil {
		pprof.StopCPUProfile()
		_ = cpuProfileFile.Close()
	}
	writeProfile("heap", os.Getenv("COSMOS_PERF_HEAP_PROFILE"))
	writeProfile("block", os.Getenv("COSMOS_PERF_BLOCK_PROFILE"))
	writeProfile("mutex", os.Getenv("COSMOS_PERF_MUTEX_PROFILE"))
	if traceProfileFile != nil {
		runtimeTrace.Stop()
		_ = traceProfileFile.Close()
	}
}

func writeProfile(name string, path string) {
	if path == "" {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	_ = pprof.Lookup(name).WriteTo(file, 0)
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func envInt(name string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return fallback
	}
	return value
}

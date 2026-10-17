// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore Failf PPAF PPCB ppaf ppcb

package azcosmos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// The fault-injection ABI uses its own operation discriminants, distinct from operationKind.
	nativeFaultOperationCreateItemTest = int32(3)

	emulatorEastUS = string(RegionEastUS)
	emulatorWestUS = string(RegionWestUS)
)

type availabilityReadItem struct {
	id string
	pk PartitionKey
}

func nativeAvailabilityContainer(
	t *testing.T,
	endpoint string,
	databaseID string,
	containerID string,
	options ClientOptions,
	rules []nativeFaultRule,
) (*Client, *ContainerClient) {
	t.Helper()

	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	options.Runtime = runtime

	cred, err := NewKeyCredential(emulatorKey)
	require.NoError(t, err)
	client, err := NewClientWithKey(endpoint, cred, &options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	client.driver.faultRules = rules
	require.NoError(t, client.Initialize(t.Context()))

	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	return client, container
}

func seedAvailabilityReadItems(t *testing.T, count int) []availabilityReadItem {
	t.Helper()

	writer := emulatorContainerWithOptions(t, &ClientOptions{
		Routing: PreferredRegions(RegionEastUS),
	})
	westReader := emulatorContainerWithOptions(t, &ClientOptions{
		Routing: PreferredRegions(RegionWestUS),
		Operation: OperationOptions{
			ConsistencyStrategy: ReadConsistencyStrategyEventual,
		},
	})

	baseID := uniqueItemID(t)
	items := make([]availabilityReadItem, count)
	for i := range items {
		id := fmt.Sprintf("%s-%02d", baseID, i)
		pkValue := fmt.Sprintf("availability-native-partition-%02d", i)
		pk := NewPartitionKeyString(pkValue)
		body, err := json.Marshal(map[string]string{"id": id, "pk": pkValue})
		require.NoError(t, err)
		trackEmulatorItem(t, writer, pk, id)
		_, err = writer.CreateItem(t.Context(), pk, id, body, nil)
		require.NoError(t, err)
		items[i] = availabilityReadItem{id: id, pk: pk}
	}

	require.Eventually(t, func() bool {
		for _, item := range items {
			if _, err := westReader.ReadItem(t.Context(), item.pk, item.id, nil); err != nil {
				return false
			}
		}
		return true
	}, 5*time.Second, 50*time.Millisecond, "seeded items did not replicate to West US")
	return items
}

func requireDiagnosticAttempt(
	t *testing.T,
	diagnostics *Diagnostics,
	region string,
	status int,
	subStatus int,
) DiagnosticAttempt {
	t.Helper()
	require.NotNil(t, diagnostics)
	for _, attempt := range diagnostics.Attempts {
		if attempt.Region == region && attempt.StatusCode == status && attempt.SubStatus == subStatus {
			require.NotEmpty(t, attempt.Endpoint)
			return attempt
		}
	}
	require.Failf(t, "missing diagnostic attempt",
		"region=%q status=%d substatus=%d attempts=%+v", region, status, subStatus, diagnostics.Attempts)
	return DiagnosticAttempt{}
}

func requireSingleRegionAttempt(t *testing.T, diagnostics *Diagnostics, region string, status int) {
	t.Helper()
	require.NotNil(t, diagnostics)
	require.Equal(t, []string{region}, diagnostics.RegionsContacted)
	require.Equal(t, uint32(1), diagnostics.AttemptCount)
	require.Len(t, diagnostics.Attempts, 1)
	require.Equal(t, region, diagnostics.Attempts[0].Region)
	require.Equal(t, status, diagnostics.Attempts[0].StatusCode)
	require.NotEmpty(t, diagnostics.Attempts[0].Endpoint)
}

func requireAvailabilityStatusError(t *testing.T, err error, code Code, status int) *Error {
	t.Helper()
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, code, cosmosErr.Code)
	require.Equal(t, status, cosmosErr.StatusCode)
	require.True(t, cosmosErr.FromWire)
	return cosmosErr
}

func configurePPCBTestEnvironment(t *testing.T, enabled bool) {
	t.Helper()
	// A successful alternate hedge would hide the failures that should trip PPCB.
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "false")
	t.Setenv("AZURE_COSMOS_PPCB_ENABLED", "true")
	t.Setenv("AZURE_COSMOS_PPCB_ENABLED_OVERRIDE", fmt.Sprintf("%t", enabled))
	t.Setenv("AZURE_COSMOS_PPCB_READ_FAILURE_THRESHOLD", "2")
	t.Setenv("AZURE_COSMOS_PPCB_COUNTER_RESET_WINDOW_MS", "1000")
	t.Setenv("AZURE_COSMOS_PPCB_PARTITION_UNAVAILABILITY_DURATION_MS", "1000")
	t.Setenv("AZURE_COSMOS_PPCB_FAILBACK_SWEEP_INTERVAL_MS", "1000")
}

func ppafEmulatorConfiguration(t *testing.T) (endpoint, databaseID, containerID string) {
	t.Helper()
	_, databaseID, containerID = emulatorConfiguration(t)
	endpoint = os.Getenv("AZCOSMOS_PPAF_ENDPOINT")
	if endpoint == "" {
		t.Skip("set AZCOSMOS_PPAF_ENDPOINT to test single-master per-partition failover")
	}
	return
}

// The pinned emulator cannot move partition writers. Model that service behavior
// at regional proxies, but persist successful writes in the real emulator.
func ppafWriteFailoverFixture(
	t *testing.T,
	endpoint, databaseID, containerID string,
	partitionKeys []string,
) (fixtureEndpoint, samePartitionKey, otherPartitionKey string) {
	t.Helper()
	require.NotEmpty(t, partitionKeys)
	client := &http.Client{Timeout: 5 * time.Second}
	t.Cleanup(client.CloseIdleConnections)
	get := func(path, partitionKey string) (int, http.Header, []byte) {
		t.Helper()
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, endpoint+path, nil)
		require.NoError(t, err)
		if partitionKey != "" {
			value, err := json.Marshal([]string{partitionKey})
			require.NoError(t, err)
			request.Header.Set("x-ms-documentdb-partitionkey", string(value))
		}
		response, err := client.Do(request)
		require.NoError(t, err)
		body, err := io.ReadAll(response.Body)
		require.NoError(t, response.Body.Close())
		require.NoError(t, err)
		return response.StatusCode, response.Header, body
	}

	status, _, metadata := get("", "")
	require.Equal(t, http.StatusOK, status)
	var account struct {
		ReadableLocations []struct {
			Name                    string `json:"name"`
			DatabaseAccountEndpoint string `json:"databaseAccountEndpoint"`
		} `json:"readableLocations"`
		EnableMultipleWriteLocations bool `json:"enableMultipleWriteLocations"`
	}
	require.NoError(t, json.Unmarshal(metadata, &account))
	require.False(t, account.EnableMultipleWriteLocations)
	require.Len(t, account.ReadableLocations, 2)
	require.Equal(t, "East US", account.ReadableLocations[0].Name)
	require.Equal(t, "West US", account.ReadableLocations[1].Name)

	docsPath := "dbs/" + url.PathEscape(databaseID) + "/colls/" + url.PathEscape(containerID) + "/docs"
	rangeIDs := make(map[string]string, len(partitionKeys))
	for _, pk := range partitionKeys {
		// A missing-item read identifies the physical partition without modifying it.
		status, headers, _ := get(docsPath+"/"+url.PathEscape(uniqueItemID(t)), pk)
		require.Equal(t, http.StatusNotFound, status)
		rangeIDs[pk] = headers.Get("x-ms-documentdb-partitionkeyrangeid")
		require.NotEmpty(t, rangeIDs[pk])
	}
	movedRange := rangeIDs[partitionKeys[0]]
	for _, pk := range partitionKeys[1:] {
		if rangeIDs[pk] == movedRange {
			samePartitionKey = pk
		} else {
			otherPartitionKey = pk
		}
	}
	require.NotEmpty(t, samePartitionKey, "candidate keys must include another key in the moved physical partition")
	require.NotEmpty(t, otherPartitionKey, "candidate keys must include an unmoved physical partition")

	servers := make([]*httptest.Server, len(account.ReadableLocations))
	proxies := make([]*httputil.ReverseProxy, len(servers))
	for i, location := range account.ReadableLocations {
		upstream, err := url.Parse(location.DatabaseAccountEndpoint)
		require.NoError(t, err)
		proxies[i] = httputil.NewSingleHostReverseProxy(upstream)
		proxies[i].ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
			t.Errorf("PPAF emulator proxy: %v", err)
			http.Error(w, "emulator proxy failed", http.StatusBadGateway)
		}
		servers[i] = httptest.NewUnstartedServer(nil)
		t.Cleanup(servers[i].Close)
		proxyEndpoint := "http://" + servers[i].Listener.Addr().String() + "/"
		original, err := json.Marshal(location.DatabaseAccountEndpoint)
		require.NoError(t, err)
		replacement, err := json.Marshal(proxyEndpoint)
		require.NoError(t, err)
		require.Contains(t, string(metadata), string(original))
		metadata = bytes.ReplaceAll(metadata, original, replacement)
	}
	for i, server := range servers {
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet && request.URL.Path == "/" {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write(metadata); err != nil {
					t.Errorf("writing PPAF account metadata: %v", err)
				}
				return
			}
			if request.Method != http.MethodGet &&
				(request.URL.Path == "/"+docsPath || strings.HasPrefix(request.URL.Path, "/"+docsPath+"/")) {
				var pk []string
				err := json.Unmarshal([]byte(request.Header.Get("x-ms-documentdb-partitionkey")), &pk)
				if err != nil || len(pk) != 1 || rangeIDs[pk[0]] == "" {
					t.Errorf("unexpected partition key at PPAF fixture: %v", err)
					http.Error(w, "unknown fixture partition key", http.StatusBadRequest)
					return
				}
				writer := 0
				if rangeIDs[pk[0]] == movedRange {
					writer = 1
				}
				if i != writer {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("x-ms-substatus", "3")
					w.WriteHeader(http.StatusForbidden)
					if _, err := io.WriteString(w, `{"code":"Forbidden","message":"partition writer is in another region"}`); err != nil {
						t.Errorf("writing PPAF wrong-writer response: %v", err)
					}
					return
				}
				// The proxy models writer ownership; the emulator supplies the actual CRUD response.
				proxies[0].ServeHTTP(w, request)
				return
			}
			proxies[i].ServeHTTP(w, request)
		})
		server.Start()
	}
	return servers[0].URL + "/", samePartitionKey, otherPartitionKey
}

func TestNativeHedgingAlternateRegionWins(t *testing.T) {
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "true")
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	delay := 250 * time.Millisecond
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			HedgingEnabled:       to(true),
			AvailabilityStrategy: HedgingAvailability(10 * time.Millisecond),
		},
	}, []nativeFaultRule{{
		id: "delay-preferred-read", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(1)), delay: &delay,
	}})

	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	requireItemFields(t, response.Value, map[string]any{"id": item.id})
	// Rust omits the losing leg when it is canceled before observing a terminal outcome.
	requireSingleRegionAttempt(t, response.Diagnostics, string(RegionWestUS), http.StatusOK)
	require.Equal(t, response.AttemptCount, response.Diagnostics.AttemptCount)

	// The one-hit delay must not pin routing to the alternate. A generous threshold
	// keeps scheduler latency from making this recovery check another hedge race.
	next, err := container.ReadItem(t.Context(), item.pk, item.id, &ReadItemOptions{
		Operation: OperationOptions{AvailabilityStrategy: HedgingAvailability(5 * time.Second)},
	})
	require.NoError(t, err)
	requireSingleRegionAttempt(t, next.Diagnostics, string(RegionEastUS), http.StatusOK)
}

func TestNativeHedgingPrimaryWinsBeforeThreshold(t *testing.T) {
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "true")
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			HedgingEnabled:       to(true),
			AvailabilityStrategy: HedgingAvailability(5 * time.Second),
		},
	}, nil)

	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	requireItemFields(t, response.Value, map[string]any{"id": item.id})
	requireSingleRegionAttempt(t, response.Diagnostics, string(RegionEastUS), http.StatusOK)
	require.Equal(t, response.AttemptCount, response.Diagnostics.AttemptCount)
}

func TestNativeHedgingPreservesObservedLosingResponse(t *testing.T) {
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "true")
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			HedgingEnabled:       to(true),
			AvailabilityStrategy: HedgingAvailability(50 * time.Millisecond),
		},
	}, []nativeFaultRule{{
		id: "throttle-primary-before-hedge", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		errorType: nativeFaultErrorTooManyRequests, hitLimit: to(int64(1)),
		retryAfter: to(time.Second),
	}})

	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	requireItemFields(t, response.Value, map[string]any{"id": item.id})
	require.NotNil(t, response.Diagnostics)
	require.Equal(t, uint32(2), response.AttemptCount)
	require.Equal(t, response.AttemptCount, response.Diagnostics.AttemptCount)
	require.Len(t, response.Diagnostics.Attempts, 2)
	require.Equal(t, []string{string(RegionEastUS), string(RegionWestUS)}, response.Diagnostics.RegionsContacted)

	// Unlike an unfinished losing request, the primary's observed 429 survives cancellation.
	primary := requireDiagnosticAttempt(t, response.Diagnostics, string(RegionEastUS),
		http.StatusTooManyRequests, int(nativeFaultThrottleSubStatus))
	alternate := requireDiagnosticAttempt(t, response.Diagnostics, string(RegionWestUS), http.StatusOK, 0)
	require.Equal(t, primary, response.Diagnostics.Attempts[0])
	require.Equal(t, alternate, response.Diagnostics.Attempts[1])
	require.Equal(t, primary.RequestCharge+alternate.RequestCharge, response.Diagnostics.TotalRequestCharge)
}

func TestNativeHedgingDisabledDoesNotContactAlternateRegion(t *testing.T) {
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "false")
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	delay := 50 * time.Millisecond
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			HedgingEnabled:       to(true),
			AvailabilityStrategy: HedgingAvailability(time.Millisecond),
		},
	}, []nativeFaultRule{{
		id: "delay-with-hedging-disabled", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(1)), delay: &delay,
	}})

	start := time.Now()
	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	elapsed := time.Since(start)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.GreaterOrEqual(t, elapsed, delay)
	requireSingleRegionAttempt(t, response.Diagnostics, string(RegionEastUS), http.StatusOK)
}

func TestNativeHedgingDisabledPerOperation(t *testing.T) {
	// Register restoration before removing environment switches that override the strategy.
	for _, name := range []string{"AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "AZURE_COSMOS_HEDGING_ENABLED"} {
		t.Setenv(name, "")
		require.NoError(t, os.Unsetenv(name))
	}
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	delay := 100 * time.Millisecond
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			AvailabilityStrategy: HedgingAvailability(10 * time.Millisecond),
		},
	}, []nativeFaultRule{{
		id: "delay-with-operation-hedging-disabled", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(1)), delay: &delay,
	}})

	start := time.Now()
	response, err := container.ReadItem(t.Context(), item.pk, item.id, &ReadItemOptions{
		Operation: OperationOptions{AvailabilityStrategy: DisabledAvailability()},
	})
	elapsed := time.Since(start)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.GreaterOrEqual(t, elapsed, delay)
	requireSingleRegionAttempt(t, response.Diagnostics, string(RegionEastUS), http.StatusOK)
}

func TestNativeHedgingCancellationDrainsLateCompletions(t *testing.T) {
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "true")
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	delay := 500 * time.Millisecond
	client, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			HedgingEnabled:       to(true),
			AvailabilityStrategy: HedgingAvailability(10 * time.Millisecond),
		},
	}, []nativeFaultRule{
		{
			id: "delay-cancelled-primary", kind: nativeFaultOperationReadItem,
			region: emulatorEastUS, containerID: containerID,
			hitLimit: to(int64(1)), delay: &delay,
		},
		{
			id: "delay-cancelled-hedge", kind: nativeFaultOperationReadItem,
			region: emulatorWestUS, containerID: containerID,
			hitLimit: to(int64(1)), delay: &delay,
		},
	})

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	response, err := container.ReadItem(ctx, item.pk, item.id, nil)
	require.Equal(t, ItemResponse{}, response)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), time.Second)

	closed := make(chan error, 1)
	go func() { closed <- client.Close() }()
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Client.Close did not drain the late native hedge completions")
	}
}

func TestNativePPCBTripsPerPartitionAndFailsBack(t *testing.T) {
	configurePPCBTestEnvironment(t, true)
	items := seedAvailabilityReadItems(t, 16)
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			// Native PPCB records 503 failures on the failover-retry path.
			MaxFailoverRetryCount: to(uint32(1)),
		},
	}, []nativeFaultRule{{
		id: "trip-east-us-read-partition", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(2)), customStatusCode: to(int32(http.StatusServiceUnavailable)),
	}})

	partitionA := items[0]
	for range 2 {
		response, err := container.ReadItem(t.Context(), partitionA.pk, partitionA.id, nil)
		require.NoError(t, err)
		require.NotNil(t, response.Diagnostics)
		require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, response.Diagnostics.RegionsContacted)
		require.Equal(t, uint32(2), response.Diagnostics.AttemptCount)
		requireDiagnosticAttempt(t, response.Diagnostics, emulatorEastUS, http.StatusServiceUnavailable, 0)
		requireDiagnosticAttempt(t, response.Diagnostics, emulatorWestUS, http.StatusOK, 0)
	}

	routedAway, err := container.ReadItem(t.Context(), partitionA.pk, partitionA.id, nil)
	require.NoError(t, err)
	requireSingleRegionAttempt(t, routedAway.Diagnostics, emulatorWestUS, http.StatusOK)

	var partitionB *availabilityReadItem
	for i := 1; i < len(items); i++ {
		response, err := container.ReadItem(t.Context(), items[i].pk, items[i].id, nil)
		require.NoError(t, err)
		require.NotNil(t, response.Diagnostics)
		require.NotEmpty(t, response.Diagnostics.RegionsContacted)
		if response.Diagnostics.RegionsContacted[0] == emulatorEastUS {
			partitionB = &items[i]
			requireSingleRegionAttempt(t, response.Diagnostics, emulatorEastUS, http.StatusOK)
			break
		}
	}
	require.NotNil(t, partitionB, "candidate keys must cover a partition distinct from partition A")

	var failedBack ItemResponse
	require.Eventually(t, func() bool {
		response, err := container.ReadItem(t.Context(), partitionA.pk, partitionA.id, nil)
		if err != nil || response.Diagnostics == nil || len(response.Diagnostics.RegionsContacted) == 0 {
			return false
		}
		if response.Diagnostics.RegionsContacted[0] != emulatorEastUS {
			return false
		}
		failedBack = response
		return true
	}, 8*time.Second, 100*time.Millisecond, "PPCB did not probe and fail back to East US")
	requireSingleRegionAttempt(t, failedBack.Diagnostics, emulatorEastUS, http.StatusOK)
}

func TestNativePPCBKillSwitchPreventsTrip(t *testing.T) {
	configurePPCBTestEnvironment(t, false)
	item := seedAvailabilityReadItems(t, 1)[0]
	endpoint, databaseID, containerID := emulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			MaxFailoverRetryCount: to(uint32(1)),
		},
	}, []nativeFaultRule{{
		id: "disabled-ppcb-east-us-failures", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(2)), customStatusCode: to(int32(http.StatusServiceUnavailable)),
	}})

	first, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	require.NotNil(t, first.Diagnostics)
	require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, first.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, first.Diagnostics, emulatorEastUS, http.StatusServiceUnavailable, 0)
	requireDiagnosticAttempt(t, first.Diagnostics, emulatorWestUS, http.StatusOK, 0)

	// Without PPCB, the first 503 marks the endpoint unavailable account-wide.
	// With PPCB enabled, a second East US failure would be needed to trip it.
	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	requireSingleRegionAttempt(t, response.Diagnostics, emulatorWestUS, http.StatusOK)
}

func TestNativePPAFDiscoversWriteRegionPerPartition(t *testing.T) {
	t.Setenv("AZURE_COSMOS_HEDGING_ENABLED_OVERRIDE", "false")
	endpoint, databaseID, containerID := ppafEmulatorConfiguration(t)
	partitionA := "availability-ppaf-partition-a"
	partitionKeys := []string{partitionA}
	for i := range 16 {
		partitionKeys = append(partitionKeys, fmt.Sprintf("availability-ppaf-partition-b-%02d", i))
	}
	fixtureEndpoint, samePartitionKey, otherPartitionKey := ppafWriteFailoverFixture(t, endpoint, databaseID, containerID, partitionKeys)
	_, container := nativeAvailabilityContainer(t, fixtureEndpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			MaxFailoverRetryCount: to(uint32(1)),
		},
	}, nil)

	type storedItem struct {
		id string
		pk string
	}
	var items []storedItem
	create := func(pkValue string) ItemResponse {
		t.Helper()
		pk := NewPartitionKeyString(pkValue)
		id := uniqueItemID(t)
		body, err := json.Marshal(map[string]string{"id": id, "pk": pkValue})
		require.NoError(t, err)
		trackEmulatorItem(t, container, pk, id)
		response, err := container.CreateItem(t.Context(), pk, id, body, nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, response.StatusCode)
		items = append(items, storedItem{id: id, pk: pkValue})
		return response
	}

	first := create(partitionA)
	require.NotNil(t, first.Diagnostics)
	require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, first.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, first.Diagnostics, emulatorEastUS, http.StatusForbidden, 3)
	requireDiagnosticAttempt(t, first.Diagnostics, emulatorWestUS, http.StatusCreated, 0)

	for _, pk := range []string{partitionA, samePartitionKey} {
		response := create(pk)
		requireSingleRegionAttempt(t, response.Diagnostics, emulatorWestUS, http.StatusCreated)
	}

	// The initial 403/3 also marks East US unavailable account-wide. A cold
	// partition must discover its own writer before testing warm-cache isolation.
	other := create(otherPartitionKey)
	require.NotNil(t, other.Diagnostics)
	require.Equal(t, []string{emulatorWestUS, emulatorEastUS}, other.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, other.Diagnostics, emulatorWestUS, http.StatusForbidden, 3)
	requireDiagnosticAttempt(t, other.Diagnostics, emulatorEastUS, http.StatusCreated, 0)

	for _, expected := range []struct {
		pk     string
		region string
	}{
		{partitionA, emulatorWestUS},
		{samePartitionKey, emulatorWestUS},
		{otherPartitionKey, emulatorEastUS},
	} {
		response := create(expected.pk)
		requireSingleRegionAttempt(t, response.Diagnostics, expected.region, http.StatusCreated)
	}

	_, reader := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionWestUS),
		Operation: OperationOptions{
			ConsistencyStrategy:   ReadConsistencyStrategyEventual,
			MaxFailoverRetryCount: to(uint32(0)),
		},
	}, nil)
	for _, item := range items {
		var stored ItemResponse
		require.Eventually(t, func() bool {
			response, err := reader.ReadItem(t.Context(), NewPartitionKeyString(item.pk), item.id, nil)
			stored = response
			return err == nil
		}, 5*time.Second, 50*time.Millisecond, "PPAF write did not persist and replicate")
		requireSingleRegionAttempt(t, stored.Diagnostics, emulatorWestUS, http.StatusOK)
		requireItemFields(t, stored.Value, map[string]any{"id": item.id, "pk": item.pk})
	}
}

func TestNativePPAFDiscardsUnprovenWriteRegion(t *testing.T) {
	endpoint, databaseID, containerID := ppafEmulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			// The 403/3 discovery uses its own retry budget; stop at West US's 503.
			MaxFailoverRetryCount: to(uint32(0)),
		},
	}, []nativeFaultRule{
		{
			id: "unproven-east-us-write-forbidden", kind: nativeFaultOperationCreateItemTest,
			region: emulatorEastUS, containerID: containerID,
			hitLimit:         to(int64(1)),
			customStatusCode: to(int32(http.StatusForbidden)), customSubStatus: to(int32(3)),
		},
		{
			id: "unproven-west-us-response", kind: nativeFaultOperationCreateItemTest,
			region: emulatorWestUS, containerID: containerID,
			hitLimit: to(int64(1)), customStatusCode: to(int32(http.StatusServiceUnavailable)),
		},
	})

	pkValue := "availability-ppaf-unproven-partition"
	pk := NewPartitionKeyString(pkValue)
	failedID := uniqueItemID(t)
	failedBody, err := json.Marshal(map[string]string{"id": failedID, "pk": pkValue})
	require.NoError(t, err)
	response, err := container.CreateItem(t.Context(), pk, failedID, failedBody, nil)
	require.Equal(t, ItemResponse{}, response)
	cosmosErr := requireAvailabilityStatusError(t, err, CodeServiceUnavailable, http.StatusServiceUnavailable)
	require.NotNil(t, cosmosErr.Diagnostics)
	require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, cosmosErr.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, cosmosErr.Diagnostics, emulatorEastUS, http.StatusForbidden, 3)
	requireDiagnosticAttempt(t, cosmosErr.Diagnostics, emulatorWestUS, http.StatusServiceUnavailable, 0)

	nextID := uniqueItemID(t)
	nextBody, err := json.Marshal(map[string]string{"id": nextID, "pk": pkValue})
	require.NoError(t, err)
	trackEmulatorItem(t, container, pk, nextID)
	next, err := container.CreateItem(t.Context(), pk, nextID, nextBody, nil)
	require.NoError(t, err)
	requireSingleRegionAttempt(t, next.Diagnostics, emulatorEastUS, http.StatusCreated)
}

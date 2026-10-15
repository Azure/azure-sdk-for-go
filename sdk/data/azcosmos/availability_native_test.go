// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// The fault-injection ABI uses its own operation discriminants, distinct from operationKind.
	nativeFaultOperationCreateItemTest = int32(3)

	emulatorEastUS = "East US"
	emulatorWestUS = "West US"
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
	require.NotEmpty(t, response.Value)
	require.NotNil(t, response.Diagnostics)
	require.Equal(t, response.AttemptCount, response.Diagnostics.AttemptCount)
	require.GreaterOrEqual(t, response.AttemptCount, uint32(2))
	require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, response.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, response.Diagnostics, emulatorWestUS, http.StatusOK, 0)

	var retainedEast bool
	for _, attempt := range response.Diagnostics.Attempts {
		require.NotEmpty(t, attempt.Endpoint)
		require.NotEmpty(t, attempt.Region)
		if attempt.Region == emulatorEastUS {
			retainedEast = true
			require.Contains(t, []int{0, http.StatusOK}, attempt.StatusCode)
		}
	}
	require.True(t, retainedEast, "the delayed preferred attempt must remain visible in diagnostics")
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

	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	requireSingleRegionAttempt(t, response.Diagnostics, emulatorEastUS, http.StatusOK)
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
			MaxFailoverRetryCount: to(uint32(0)),
		},
	}, []nativeFaultRule{{
		id: "trip-east-us-read-partition", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(2)), customStatusCode: to(int32(http.StatusServiceUnavailable)),
	}})

	partitionA := items[0]
	for range 2 {
		response, err := container.ReadItem(t.Context(), partitionA.pk, partitionA.id, nil)
		require.Equal(t, ItemResponse{}, response)
		cosmosErr := requireAvailabilityStatusError(t, err, CodeServiceUnavailable, http.StatusServiceUnavailable)
		requireSingleRegionAttempt(t, cosmosErr.Diagnostics, emulatorEastUS, http.StatusServiceUnavailable)
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
			MaxFailoverRetryCount: to(uint32(0)),
		},
	}, []nativeFaultRule{{
		id: "disabled-ppcb-east-us-failures", kind: nativeFaultOperationReadItem,
		region: emulatorEastUS, containerID: containerID,
		hitLimit: to(int64(2)), customStatusCode: to(int32(http.StatusServiceUnavailable)),
	}})

	for range 2 {
		_, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
		cosmosErr := requireAvailabilityStatusError(t, err, CodeServiceUnavailable, http.StatusServiceUnavailable)
		requireSingleRegionAttempt(t, cosmosErr.Diagnostics, emulatorEastUS, http.StatusServiceUnavailable)
	}
	response, err := container.ReadItem(t.Context(), item.pk, item.id, nil)
	require.NoError(t, err)
	requireSingleRegionAttempt(t, response.Diagnostics, emulatorEastUS, http.StatusOK)
}

func TestNativePPAFDiscoversWriteRegionPerPartition(t *testing.T) {
	endpoint, databaseID, containerID := ppafEmulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			MaxFailoverRetryCount: to(uint32(1)),
		},
	}, []nativeFaultRule{{
		id: "east-us-write-forbidden", kind: nativeFaultOperationCreateItemTest,
		region: emulatorEastUS, containerID: containerID,
		hitLimit:         to(int64(1)),
		customStatusCode: to(int32(http.StatusForbidden)), customSubStatus: to(int32(3)),
	}})

	partitionA := "availability-ppaf-partition-a"
	firstID := uniqueItemID(t)
	firstPK := NewPartitionKeyString(partitionA)
	firstBody, err := json.Marshal(map[string]string{"id": firstID, "pk": partitionA})
	require.NoError(t, err)
	trackEmulatorItem(t, container, firstPK, firstID)
	first, err := container.CreateItem(t.Context(), firstPK, firstID, firstBody, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, first.StatusCode)
	require.NotNil(t, first.Diagnostics)
	require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, first.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, first.Diagnostics, emulatorEastUS, http.StatusForbidden, 3)
	requireDiagnosticAttempt(t, first.Diagnostics, emulatorWestUS, http.StatusCreated, 0)

	secondID := uniqueItemID(t)
	secondBody, err := json.Marshal(map[string]string{"id": secondID, "pk": partitionA})
	require.NoError(t, err)
	trackEmulatorItem(t, container, firstPK, secondID)
	second, err := container.CreateItem(t.Context(), firstPK, secondID, secondBody, nil)
	require.NoError(t, err)
	requireSingleRegionAttempt(t, second.Diagnostics, emulatorWestUS, http.StatusCreated)

	var isolated bool
	for i := 0; i < 16; i++ {
		pkValue := fmt.Sprintf("availability-ppaf-partition-b-%02d", i)
		pk := NewPartitionKeyString(pkValue)
		id := uniqueItemID(t)
		body, err := json.Marshal(map[string]string{"id": id, "pk": pkValue})
		require.NoError(t, err)
		trackEmulatorItem(t, container, pk, id)
		response, err := container.CreateItem(t.Context(), pk, id, body, nil)
		require.NoError(t, err)
		require.NotNil(t, response.Diagnostics)
		require.NotEmpty(t, response.Diagnostics.RegionsContacted)
		if response.Diagnostics.RegionsContacted[0] == emulatorEastUS {
			isolated = true
			requireSingleRegionAttempt(t, response.Diagnostics, emulatorEastUS, http.StatusCreated)
			break
		}
	}
	require.True(t, isolated, "candidate keys must cover a partition distinct from partition A")
}

func TestNativePPAFDiscardsUnprovenWriteRegion(t *testing.T) {
	endpoint, databaseID, containerID := ppafEmulatorConfiguration(t)
	_, container := nativeAvailabilityContainer(t, endpoint, databaseID, containerID, ClientOptions{
		Routing: PreferredRegions(RegionEastUS, RegionWestUS),
		Operation: OperationOptions{
			MaxFailoverRetryCount: to(uint32(1)),
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
			hitLimit: to(int64(1)), customStatusCode: to(int32(http.StatusTeapot)),
		},
	})

	pkValue := "availability-ppaf-unproven-partition"
	pk := NewPartitionKeyString(pkValue)
	failedID := uniqueItemID(t)
	failedBody, err := json.Marshal(map[string]string{"id": failedID, "pk": pkValue})
	require.NoError(t, err)
	response, err := container.CreateItem(t.Context(), pk, failedID, failedBody, nil)
	require.Equal(t, ItemResponse{}, response)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, http.StatusTeapot, cosmosErr.StatusCode)
	require.NotNil(t, cosmosErr.Diagnostics)
	require.Equal(t, []string{emulatorEastUS, emulatorWestUS}, cosmosErr.Diagnostics.RegionsContacted)
	requireDiagnosticAttempt(t, cosmosErr.Diagnostics, emulatorEastUS, http.StatusForbidden, 3)
	requireDiagnosticAttempt(t, cosmosErr.Diagnostics, emulatorWestUS, http.StatusTeapot, 0)

	nextID := uniqueItemID(t)
	nextBody, err := json.Marshal(map[string]string{"id": nextID, "pk": pkValue})
	require.NoError(t, err)
	trackEmulatorItem(t, container, pk, nextID)
	next, err := container.CreateItem(t.Context(), pk, nextID, nextBody, nil)
	require.NoError(t, err)
	requireSingleRegionAttempt(t, next.Diagnostics, emulatorEastUS, http.StatusCreated)
}

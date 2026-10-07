// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// These rules are installed in the native driver before initialization. A single Go ReadItem
// submits once; the driver's diagnostics count the attempts made inside that submission.
func faultInjectedContainer(t *testing.T, rule nativeFaultRule) *ContainerClient {
	t.Helper()
	client, databaseID, containerID := emulatorClientConfigured(t, nil, false)
	client.driver.faultRules = []nativeFaultRule{rule}
	require.NoError(t, client.Initialize(t.Context()))
	container, err := client.NewContainer(databaseID, containerID)
	require.NoError(t, err)
	return container
}

func seedFaultInjectedItem(t *testing.T, container *ContainerClient) (string, PartitionKey, SessionToken) {
	t.Helper()
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	trackEmulatorItem(t, container, pk, id)
	created, err := container.CreateItem(t.Context(), pk, id, []byte(`{"id":"`+id+`","pk":"`+id+`"}`), nil)
	require.NoError(t, err)
	require.NotEmpty(t, created.SessionToken)
	return id, pk, created.SessionToken
}

func TestNativeThrottleRetryReturnsSingleFinalResponse(t *testing.T) {
	container := faultInjectedContainer(t, nativeFaultRule{
		id: "throttle-twice", kind: 1, errorType: 2, hitLimit: 2,
		delayMS: -1, retryAfter: 1,
	})
	id, pk, token := seedFaultInjectedItem(t, container)
	response, err := container.ReadItem(t.Context(), pk, id, &ReadItemOptions{SessionToken: token})
	require.NoError(t, err)
	requireItemFields(t, response.Value, map[string]any{"id": id, "pk": id})
	require.Equal(t, uint32(3), response.AttemptCount, "two retries belong to one native submission")
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Zero(t, response.SubStatus)
	require.NotNil(t, response.Diagnostics)
	require.Equal(t, response.AttemptCount, response.Diagnostics.AttemptCount)
	require.Equal(t, response.StatusCode, response.Diagnostics.StatusCode)
	require.Zero(t, response.Diagnostics.SubStatus)
	require.True(t, response.Diagnostics.Completed)
	require.False(t, response.Diagnostics.Failed)
	require.False(t, response.Diagnostics.Compacted)
	require.Positive(t, response.Diagnostics.Elapsed)
	require.NotEmpty(t, response.Diagnostics.RegionsContacted)
	require.Len(t, response.Diagnostics.Attempts, 3)
	require.GreaterOrEqual(t, response.Diagnostics.TotalRequestCharge, response.RequestCharge)
	require.True(t, json.Valid([]byte(response.Diagnostics.JSON)))
	for _, attempt := range response.Diagnostics.Attempts[:2] {
		require.Equal(t, http.StatusTooManyRequests, attempt.StatusCode)
		require.Equal(t, 3200, attempt.SubStatus)
		require.NotEmpty(t, attempt.Endpoint)
	}
	require.Equal(t, http.StatusOK, response.Diagnostics.Attempts[2].StatusCode)
	require.NotEmpty(t, response.SessionToken)

	next, err := container.ReadItem(t.Context(), pk, id, &ReadItemOptions{SessionToken: response.SessionToken})
	require.NoError(t, err)
	require.Equal(t, uint32(1), next.AttemptCount)
	require.Len(t, response.Diagnostics.Attempts, 3, "the snapshot survives later native completions")
	require.True(t, json.Valid([]byte(response.Diagnostics.JSON)))
}

func TestNativeThrottleRetryExhaustion(t *testing.T) {
	container := faultInjectedContainer(t, nativeFaultRule{
		id: "always-throttle", kind: 1, errorType: 2, hitLimit: 100,
		delayMS: -1, retryAfter: 1,
	})
	id, pk, _ := seedFaultInjectedItem(t, container)
	response, err := container.ReadItem(t.Context(), pk, id, nil)
	require.Equal(t, ItemResponse{}, response)
	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeThrottled, cosmosErr.Code)
	require.Equal(t, http.StatusTooManyRequests, cosmosErr.StatusCode)
	require.Equal(t, 3200, cosmosErr.SubStatus)
	require.NotNil(t, cosmosErr.Diagnostics)
	require.Equal(t, cosmosErr.AttemptCount, cosmosErr.Diagnostics.AttemptCount)
	require.Equal(t, cosmosErr.StatusCode, cosmosErr.Diagnostics.StatusCode)
	require.Equal(t, cosmosErr.SubStatus, cosmosErr.Diagnostics.SubStatus)
	require.True(t, cosmosErr.Diagnostics.Completed)
	require.True(t, cosmosErr.Diagnostics.Failed)
	require.Positive(t, cosmosErr.Diagnostics.Elapsed)
	require.NotEmpty(t, cosmosErr.Diagnostics.RegionsContacted)
	require.True(t, json.Valid([]byte(cosmosErr.Diagnostics.JSON)))
	require.Equal(t, int(cosmosErr.AttemptCount), len(cosmosErr.Diagnostics.Attempts))
	for _, attempt := range cosmosErr.Diagnostics.Attempts {
		require.Equal(t, http.StatusTooManyRequests, attempt.StatusCode)
		require.Equal(t, 3200, attempt.SubStatus)
	}
	require.True(t, cosmosErr.FromWire)
	require.Positive(t, cosmosErr.RetryAfter)
}

func TestNativeStaleSessionRecovery(t *testing.T) {
	container := faultInjectedContainer(t, nativeFaultRule{
		id: "stale-session-once", kind: 1, errorType: 4, hitLimit: 1,
		delayMS: -1, retryAfter: -1,
	})
	id, pk, token := seedFaultInjectedItem(t, container)
	response, err := container.ReadItem(t.Context(), pk, id, &ReadItemOptions{
		SessionToken: token,
		Operation:    OperationOptions{ConsistencyStrategy: ReadConsistencyStrategySession},
	})
	require.NoError(t, err)
	require.Equal(t, uint32(2), response.AttemptCount)
	require.NotNil(t, response.Diagnostics)
	require.Len(t, response.Diagnostics.Attempts, 2)
	require.Equal(t, http.StatusNotFound, response.Diagnostics.Attempts[0].StatusCode)
	require.Equal(t, subStatusReadSessionNotAvailable, response.Diagnostics.Attempts[0].SubStatus)
	require.Equal(t, http.StatusOK, response.Diagnostics.Attempts[1].StatusCode)
	require.NotEmpty(t, response.SessionToken)

	next, err := container.ReadItem(t.Context(), pk, id, &ReadItemOptions{SessionToken: response.SessionToken})
	require.NoError(t, err)
	require.Equal(t, uint32(1), next.AttemptCount)
}

func TestNativeThrottleCancellationReturnsPromptly(t *testing.T) {
	container := faultInjectedContainer(t, nativeFaultRule{
		id: "delayed-throttle", kind: 1, errorType: 2, hitLimit: 100,
		delayMS: 2000, retryAfter: 1,
	})
	id, pk, _ := seedFaultInjectedItem(t, container)
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	response, err := container.ReadItem(ctx, pk, id, nil)
	require.Equal(t, ItemResponse{}, response)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, time.Since(start), time.Second)
	// The v0.2 ABI has no native cancellation. This checks the Go wait/handle lifecycle,
	// not termination of the native retry task.
}

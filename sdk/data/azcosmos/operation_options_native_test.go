// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFullOperationOptionsNativePresence(t *testing.T) {
	options := OperationOptions{
		ConsistencyStrategy: ReadConsistencyStrategyDefault, EnableContentResponseOnWrite: to(false),
		ExcludedRegions: []Region{}, EndToEndTimeout: 3 * time.Second, PatchStrategy: PatchStrategyClientSide,
		SessionCapturingDisabled: to(false), MaxFailoverRetryCount: to(uint32(math.MaxUint32)),
		MaxSessionRetryCount: to(uint32(0)), EndpointUnavailabilityTTL: to(time.Duration(0)),
		CustomHeaders: map[string]string{"X-Test": "value"}, BinaryEncoding: &BinaryEncodingOptions{Enabled: true, RequestTextResponse: true},
		ThroughputControl: ThroughputControlOptions{ThroughputBucket: to(uint32(math.MaxUint32)), PriorityLevel: PriorityLevelLow},
		ThrottlingRetry:   ThrottlingRetryOptions{MaxRetryCount: to(uint32(0)), MaxRetryWaitTime: to(time.Duration(0))},
		HedgingEnabled:    to(false), AvailabilityStrategy: HedgingAvailability(time.Nanosecond),
	}
	native, release := inspectNativeOperationOptions(options)
	defer release()
	require.Equal(t, int64(math.MaxUint32), native.maxFailoverRetryCount)
	require.Zero(t, native.maxSessionRetryCount)
	require.Zero(t, native.endpointTTLMillis)
	require.NotNil(t, native.excludedRegions)
	require.Empty(t, native.excludedRegions)
	require.Equal(t, map[string]string{"x-test": "value"}, native.customHeaders)
	require.Equal(t, int8(2), native.binaryEncodingEnabled)
	require.Equal(t, int8(2), native.binaryTextResponse)
	require.Equal(t, int8(1), native.sessionCapturingDisabled)
	require.Equal(t, int64(math.MaxUint32), native.throughputBucket)
	require.Equal(t, int32(2), native.priorityLevel)
	require.Zero(t, native.maxThrottleRetryCount)
	require.Zero(t, native.maxThrottleRetryWaitMillis)
	require.Equal(t, int8(1), native.hedgingEnabled)
	require.Equal(t, int32(2), native.availabilityStrategy)
	require.Equal(t, int64(1), native.hedgeThresholdMillis)
	client, releaseClient, err := inspectNativeClientOptions(ClientOptions{Operation: options})
	require.NoError(t, err)
	defer releaseClient()
	require.Equal(t, native, client.operationOptions)
	empty, releaseEmpty := inspectNativeOperationOptions(OperationOptions{CustomHeaders: map[string]string{}})
	defer releaseEmpty()
	require.NotNil(t, empty.customHeaders)
	require.Empty(t, empty.customHeaders)
}

func TestPatchTrackingNativeFieldsAndPrecedence(t *testing.T) {
	for _, retention := range []time.Duration{0, time.Nanosecond, 1900 * time.Millisecond, 2 * time.Second} {
		request, release := inspectNativeItemRequest(itemRequest{
			kind:          operationKindPatchItem,
			options:       OperationOptions{PatchStrategy: PatchStrategyServerSide},
			patchStrategy: PatchStrategyClientSide, patchMaxAttempts: to(uint8(255)),
			patchTrackingID:       "00112233-4455-6677-8899-aabbccddeeff",
			patchTrackingCapacity: to(uint16(65535)), patchTrackingRetention: &retention,
		})
		want, _ := nativePatchStrategy(PatchStrategyClientSide)
		require.Equal(t, want, request.patchStrategy)
		require.Equal(t, uint8(255), request.patchMaxAttempts)
		require.Equal(t, "00112233-4455-6677-8899-aabbccddeeff", request.patchTrackingID)
		require.Equal(t, uint16(65535), request.patchTrackingCapacity)
		require.Equal(t, uint32(max(1, retention/time.Second)), request.patchTrackingRetention)
		release()
	}
}

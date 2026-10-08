// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func to[T any](value T) *T { return &value }

func TestOperationOptionsValidation(t *testing.T) {
	for name, options := range map[string]OperationOptions{
		"consistency":      {ConsistencyStrategy: "invalid"},
		"patch":            {PatchStrategy: "invalid"},
		"timeout":          {EndToEndTimeout: to(time.Duration(-1))},
		"ttl":              {EndpointUnavailabilityTTL: to(time.Duration(-1))},
		"throttle wait":    {ThrottlingRetry: ThrottlingRetryOptions{MaxRetryWaitTime: to(time.Duration(-1))}},
		"priority":         {ThroughputControl: ThroughputControlOptions{PriorityLevel: "invalid"}},
		"hedging zero":     {AvailabilityStrategy: HedgingAvailability(0)},
		"hedging negative": {AvailabilityStrategy: HedgingAvailability(-1)},
		"header name":      {CustomHeaders: map[string]string{"bad name": "value"}},
		"header empty":     {CustomHeaders: map[string]string{"": "value"}},
		"header collision": {CustomHeaders: map[string]string{"X-Custom": "a", "x-custom": "b"}},
		"header injection": {CustomHeaders: map[string]string{"x-custom": "a\r\nb"}},
		"header nul":       {CustomHeaders: map[string]string{"x-custom": "a\x00b"}},
		"region":           {ExcludedRegions: []Region{""}},
	} {
		t.Run(name, func(t *testing.T) { require.Error(t, options.validate()) })
	}
	require.NoError(t, (OperationOptions{
		MaxFailoverRetryCount: to(uint32(math.MaxUint32)), MaxSessionRetryCount: to(uint32(0)),
		EndToEndTimeout: to(time.Duration(time.Nanosecond)), EndpointUnavailabilityTTL: to(time.Duration(0)),
		ExcludedRegions: []Region{}, CustomHeaders: map[string]string{"X-Custom": ""},
		AvailabilityStrategy: HedgingAvailability(time.Nanosecond),
	}).validate())
}

func TestOperationOptionsCloneOwnsAllMutableValues(t *testing.T) {
	original := OperationOptions{
		EnableContentResponseOnWrite: to(true), ExcludedRegions: []Region{RegionEastUS},
		SessionCapturingDisabled: to(true), MaxFailoverRetryCount: to(uint32(1)), MaxSessionRetryCount: to(uint32(2)),
		EndpointUnavailabilityTTL: to(time.Second), CustomHeaders: map[string]string{"x-custom": "first"},
		BinaryEncoding: &BinaryEncodingOptions{Enabled: to(true)}, HedgingEnabled: to(true),
		ThroughputControl: ThroughputControlOptions{ThroughputBucket: to(uint32(3))},
		ThrottlingRetry:   ThrottlingRetryOptions{MaxRetryCount: to(uint32(4)), MaxRetryWaitTime: to(time.Second)},
	}
	snapshot := original.clone()
	want := original.clone()
	*original.EnableContentResponseOnWrite = false
	original.ExcludedRegions[0] = RegionWestUS
	*original.SessionCapturingDisabled = false
	*original.MaxFailoverRetryCount = 8
	*original.MaxSessionRetryCount = 9
	*original.EndpointUnavailabilityTTL = 2 * time.Second
	original.CustomHeaders["x-custom"] = "second"
	original.BinaryEncoding.Enabled = to(false)
	*original.HedgingEnabled = false
	*original.ThroughputControl.ThroughputBucket = 8
	*original.ThrottlingRetry.MaxRetryCount = 8
	*original.ThrottlingRetry.MaxRetryWaitTime = 2 * time.Second
	require.Equal(t, want, snapshot)
	empty := (OperationOptions{ExcludedRegions: []Region{}, CustomHeaders: map[string]string{}}).clone()
	require.NotNil(t, empty.ExcludedRegions)
	require.NotNil(t, empty.CustomHeaders)
	require.Nil(t, (OperationOptions{}).clone().ExcludedRegions)
	require.Nil(t, (OperationOptions{}).clone().CustomHeaders)
}

func TestCommonValidationAcrossAllItemOperations(t *testing.T) {
	client := &Client{closed: true}
	container, err := client.NewContainer("db", "container")
	require.NoError(t, err)
	options := OperationOptions{MaxSessionRetryCount: to(uint32(0)), EndToEndTimeout: to(time.Duration(-1))}
	pk := NewPartitionKeyString("pk")
	ctx := context.Background()
	patch := PatchOperations{}
	require.NoError(t, patch.AppendSet("/name", "value"))
	calls := []func() (ItemResponse, error){
		func() (ItemResponse, error) {
			return container.ReadItem(ctx, pk, "id", &ReadItemOptions{Operation: options})
		},
		func() (ItemResponse, error) {
			return container.CreateItem(ctx, pk, "id", []byte("{}"), &CreateItemOptions{Operation: options})
		},
		func() (ItemResponse, error) {
			return container.ReplaceItem(ctx, pk, "id", []byte("{}"), &ReplaceItemOptions{Operation: options})
		},
		func() (ItemResponse, error) {
			return container.UpsertItem(ctx, pk, "id", []byte("{}"), &UpsertItemOptions{Operation: options})
		},
		func() (ItemResponse, error) {
			return container.DeleteItem(ctx, pk, "id", &DeleteItemOptions{Operation: options})
		},
		func() (ItemResponse, error) {
			return container.PatchItem(ctx, pk, "id", patch, &PatchItemOptions{Operation: options})
		},
	}
	for _, call := range calls {
		response, err := call()
		require.ErrorContains(t, err, "EndToEndTimeout")
		require.Empty(t, response)
	}
}

func TestTrackingValidation(t *testing.T) {
	for _, options := range []PatchItemOptions{
		{MaxAttempts: to(uint8(0))}, {TrackingCapacity: to(uint16(0))},
		{TrackingRetention: to(time.Duration(-1))},
		{TrackingID: "not-a-uuid"},
	} {
		require.Error(t, options.validateTracking())
	}
	require.NoError(t, (PatchItemOptions{
		MaxAttempts: to(uint8(255)), TrackingCapacity: to(uint16(65535)),
		TrackingRetention: to(time.Duration(0)), TrackingID: "00112233-4455-6677-8899-aabbccddeeff",
	}).validateTracking())
}

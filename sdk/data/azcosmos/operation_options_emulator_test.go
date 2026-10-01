// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore azsdk

package azcosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/require"
)

func runtimeEmulatorContainer(t *testing.T, runtime *Runtime, defaults OperationOptions) *ContainerClient {
	t.Helper()
	endpoint, database, container := emulatorConfiguration(t)
	client, err := NewClientWithKey(endpoint, KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: runtime, Operation: defaults})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	result, err := client.NewContainer(database, container)
	require.NoError(t, err)
	return result
}

func TestEmulatorBinaryDefaultAndHierarchy(t *testing.T) {
	emulatorConfiguration(t)
	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	container := runtimeEmulatorContainer(t, runtime, OperationOptions{})
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	_, err = container.CreateItem(t.Context(), pk, id, []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"count":1}`, id, id)), nil)
	require.NoError(t, err)
	trackEmulatorItem(t, container, pk, id)
	response, err := container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.NotEmpty(t, response.Value)
	require.Equal(t, byte(0x80), response.Value[0], "unset options must negotiate Cosmos binary JSON")

	text := OperationOptions{BinaryEncoding: &BinaryEncodingOptions{Enabled: true, RequestTextResponse: true}}
	require.NoError(t, runtime.SetOperationOptions(text))
	text.BinaryEncoding.RequestTextResponse = false
	response, err = container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	var item map[string]any
	require.NoError(t, json.Unmarshal(response.Value, &item))
	require.Equal(t, id, item["id"])
	require.Equal(t, float64(1), item["count"])
	text.BinaryEncoding.RequestTextResponse = true

	clientRuntime, err := NewRuntime(&RuntimeOptions{Operation: text})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, clientRuntime.Close()) })
	binaryClient := runtimeEmulatorContainer(t, clientRuntime, OperationOptions{BinaryEncoding: &BinaryEncodingOptions{Enabled: true}})
	response, err = binaryClient.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, byte(0x80), response.Value[0], "client group must replace runtime text-response group")
	response, err = binaryClient.ReadItem(t.Context(), pk, id, &ReadItemOptions{Operation: text})
	require.NoError(t, err)
	require.True(t, json.Valid(response.Value), "request must override client")
	response, err = binaryClient.ReadItem(t.Context(), pk, id, &ReadItemOptions{Operation: OperationOptions{BinaryEncoding: &BinaryEncodingOptions{}}})
	require.NoError(t, err)
	require.True(t, json.Valid(response.Value), "explicit zero group disables binary")
	require.NoError(t, runtime.SetOperationOptions(OperationOptions{}))
	response, err = container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, byte(0x80), response.Value[0], "replacement must clear prior runtime text preference")
}

func TestEmulatorSnapshotSurvivesUpdateDuringInitialization(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	setup := runtimeEmulatorContainer(t, nil, OperationOptions{})
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	_, err := setup.CreateItem(t.Context(), pk, id, []byte(fmt.Sprintf(`{"id":%q,"pk":%q}`, id, id)), nil)
	require.NoError(t, err)
	trackEmulatorItem(t, setup, pk, id)
	runtime, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{
		EndToEndTimeout: 10 * time.Second,
		BinaryEncoding:  &BinaryEncodingOptions{Enabled: true, RequestTextResponse: true},
	}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	credential := &delayedTokenCredential{started: make(chan struct{}), release: make(chan struct{})}
	client, err := NewClient(endpoint, credential, &ClientOptions{Runtime: runtime})
	require.NoError(t, err)
	container, err := client.NewContainer(database, containerID)
	require.NoError(t, err)
	type result struct {
		response ItemResponse
		err      error
	}
	done := make(chan result, 1)
	go func() {
		response, err := container.ReadItem(t.Context(), pk, id, nil)
		done <- result{response, err}
	}()
	select {
	case <-credential.started:
	case <-time.After(5 * time.Second):
		t.Fatal("lazy initialization did not start")
	}
	require.NoError(t, runtime.SetOperationOptions(OperationOptions{
		EndToEndTimeout: time.Second, BinaryEncoding: &BinaryEncodingOptions{Enabled: true},
	}))
	close(credential.release)
	select {
	case got := <-done:
		require.NoError(t, got.err)
		require.True(t, json.Valid(got.response.Value), "admitted operation must retain pre-update text encoding")
	case <-time.After(15 * time.Second):
		t.Fatal("admitted operation did not complete")
	}
	response, err := container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, byte(0x80), response.Value[0])
}

func TestEmulatorPatchTrackingAcrossApplicationRetries(t *testing.T) {
	container := emulatorContainer(t)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	_, err := container.CreateItem(t.Context(), pk, id, []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"count":0}`, id, id)), nil)
	require.NoError(t, err)
	trackEmulatorItem(t, container, pk, id)
	var patch PatchOperations
	require.NoError(t, patch.AppendIncrement("/count", 1))
	options := &PatchItemOptions{
		Operation: OperationOptions{PatchStrategy: PatchStrategyServerSide},
		Strategy:  PatchStrategyClientSide, TrackingID: "00112233-4455-6677-8899-aabbccddeeff",
		MaxAttempts: to(uint8(3)), TrackingCapacity: to(uint16(8)), TrackingRetention: to(time.Hour),
	}
	for range 2 {
		response, err := container.PatchItem(t.Context(), pk, id, patch, options)
		require.NoError(t, err)
		require.Equal(t, options.TrackingID, response.PatchTrackingID)
	}
	response, err := container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	var item map[string]any
	require.NoError(t, json.Unmarshal(response.Value, &item))
	require.Equal(t, float64(1), item["count"], "same tracking identity suppresses repeated application patch")
	require.Contains(t, item, "_azsdkPatchTracking")
	_, err = container.PatchItem(context.Background(), pk, "missing-"+id, patch, options)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, options.TrackingID, cosmosErr.PatchTrackingID)
}

func TestEmulatorRuntimeRejectsCachedCredentialReuse(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	runtime, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, runtime.Close()) })
	firstCredential := &shortLivedIdentityCredential{identity: "emulator-client-a"}
	secondCredential := &shortLivedIdentityCredential{identity: "emulator-client-b"}
	first, err := NewClient(endpoint, firstCredential, &ClientOptions{Runtime: runtime})
	require.NoError(t, err)
	require.NoError(t, first.Initialize(t.Context()))
	firstContainer, err := first.NewContainer(database, containerID)
	require.NoError(t, err)
	readMissing := func(container *ContainerClient) {
		_, err := container.ReadItem(t.Context(), NewPartitionKeyString("missing"), uniqueItemID(t), nil)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeNotFound, cosmosErr.Code)
	}
	readMissing(firstContainer)
	second, err := NewClient(endpoint, secondCredential, &ClientOptions{Runtime: runtime})
	require.Nil(t, second)
	require.ErrorContains(t, err, "already been attached")
	require.Zero(t, secondCredential.calls.Load())
	require.NoError(t, first.Close())
	second, err = NewClient(endpoint, secondCredential, &ClientOptions{Runtime: runtime})
	require.Nil(t, second)
	require.ErrorContains(t, err, "already been attached")
	require.Zero(t, secondCredential.calls.Load())

	second, err = NewClient(endpoint, secondCredential, nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, second.Close()) })
	container, err := second.NewContainer(database, containerID)
	require.NoError(t, err)
	time.Sleep(1100 * time.Millisecond)
	readMissing(container)
	require.Positive(t, secondCredential.calls.Load())
}

type shortLivedIdentityCredential struct {
	identity string
	calls    atomic.Int32
}

func (c *shortLivedIdentityCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	return azcore.AccessToken{Token: c.identity, ExpiresOn: time.Now().Add(time.Second)}, nil
}

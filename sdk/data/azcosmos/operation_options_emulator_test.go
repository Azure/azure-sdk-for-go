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

type shortLivedIdentityCredential struct {
	identity string
	calls    atomic.Int32
}

func (c *shortLivedIdentityCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	c.calls.Add(1)
	return azcore.AccessToken{Token: c.identity, ExpiresOn: time.Now().Add(time.Second)}, nil
}

func TestEmulatorBinaryClientDefaultAndRequestOverride(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	shared, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{
		BinaryEncoding: &BinaryEncodingOptions{Enabled: to(false)},
	}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	client, err := NewClientWithKey(endpoint, KeyCredential{accountKey: emulatorKey}, &ClientOptions{Runtime: shared})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	container, err := client.NewContainer(database, containerID)
	require.NoError(t, err)
	id := uniqueItemID(t)
	pk := NewPartitionKeyString(id)
	_, err = container.CreateItem(t.Context(), pk, id, []byte(fmt.Sprintf(`{"id":%q,"pk":%q,"count":1}`, id, id)), nil)
	require.NoError(t, err)
	trackEmulatorItem(t, container, pk, id)
	response, err := container.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, byte(0x80), response.Value[0], "SDK client default wins over runtime encoding for reads")
	for _, encoding := range []*BinaryEncodingOptions{
		{RequestTextResponse: true}, {Enabled: to(false)},
	} {
		response, err = container.ReadItem(t.Context(), pk, id, &ReadItemOptions{Operation: OperationOptions{BinaryEncoding: encoding}})
		require.NoError(t, err)
		require.True(t, json.Valid(response.Value))
	}
	textClient, err := NewClientWithKey(endpoint, KeyCredential{accountKey: emulatorKey},
		&ClientOptions{Runtime: shared, BinaryEncoding: &BinaryEncodingOptions{RequestTextResponse: true}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, textClient.Close()) })
	textContainer, err := textClient.NewContainer(database, containerID)
	require.NoError(t, err)
	response, err = textContainer.ReadItem(t.Context(), pk, id, nil)
	require.NoError(t, err)
	require.True(t, json.Valid(response.Value))
	response, err = textContainer.ReadItem(t.Context(), pk, id, &ReadItemOptions{Operation: OperationOptions{BinaryEncoding: &BinaryEncodingOptions{}}})
	require.NoError(t, err)
	require.Equal(t, byte(0x80), response.Value[0], "zero encoding group enables binary")
}

func TestEmulatorSameAccountClientsShareCachedCredentials(t *testing.T) {
	endpoint, database, containerID := emulatorConfiguration(t)
	shared, err := NewRuntime(nil)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	firstCredential := &shortLivedIdentityCredential{identity: "emulator-client-a"}
	secondCredential := &shortLivedIdentityCredential{identity: "emulator-client-b"}
	first, err := NewClient(endpoint, firstCredential, &ClientOptions{Runtime: shared})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, first.Close()) })
	second, err := NewClient(endpoint, secondCredential, &ClientOptions{Runtime: shared})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, second.Close()) })
	require.NoError(t, first.Initialize(t.Context()))
	require.NoError(t, second.Initialize(t.Context()))
	readMissing := func(client *Client) {
		container, err := client.NewContainer(database, containerID)
		require.NoError(t, err)
		_, err = container.ReadItem(t.Context(), NewPartitionKeyString("missing"), uniqueItemID(t), nil)
		var serviceErr *Error
		require.ErrorAs(t, err, &serviceErr)
		require.Equal(t, CodeNotFound, serviceErr.Code)
	}
	readMissing(first)
	calls := firstCredential.calls.Load()
	require.NoError(t, first.Close())
	require.NoError(t, shared.Close())
	time.Sleep(1100 * time.Millisecond)
	readMissing(second)
	require.Greater(t, firstCredential.calls.Load(), calls, "native cache keeps the first resolved account credential")
}

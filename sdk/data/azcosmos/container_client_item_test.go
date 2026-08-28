// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestContainer(t *testing.T) *ContainerClient {
	t.Helper()

	client, err := NewClientWithKey("https://myaccount.documents.azure.com", mustKeyCredential(t), nil)
	require.NoError(t, err)
	container, err := client.NewContainer("db", "items")
	require.NoError(t, err)
	return container
}

func TestReadItemRejectsEmptyID(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.ReadItem(context.Background(), NewPartitionKeyString("pk"), "", nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, errDriverUnavailable, "argument validation should run before the operation is attempted")
}

func TestCreateItemRejectsEmptyItem(t *testing.T) {
	container := newTestContainer(t)

	for _, item := range [][]byte{nil, {}} {
		_, err := container.CreateItem(context.Background(), NewPartitionKeyString("pk"), "item-1", item, nil)
		require.Error(t, err)
		require.NotErrorIs(t, err, errDriverUnavailable, "argument validation should run before the operation is attempted")
	}
}

// The id addresses the item being created, so an empty one is a caller mistake that has to be
// reported as itself rather than reaching the driver as a null argument.
func TestCreateItemRejectsEmptyID(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.CreateItem(context.Background(), NewPartitionKeyString("pk"), "", []byte(`{"id":"item-1"}`), nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, errDriverUnavailable, "argument validation should run before the operation is attempted")
}

// Argument validation runs before the context is consulted, so a caller's deterministic mistake is
// reported as itself rather than being masked by a deadline that happened to fire first.
func TestItemOperationsValidateArgumentsBeforeContext(t *testing.T) {
	container := newTestContainer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "", nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.CreateItem(ctx, NewPartitionKeyString("pk"), "item-1", nil, nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, context.Canceled)
}

// A partition key with no components is always a caller mistake: no container has a definition
// with zero paths, and PartitionKey has a usable zero value, so "declared but never appended to"
// is an easy thing to get wrong.
func TestItemOperationsRejectEmptyPartitionKey(t *testing.T) {
	container := newTestContainer(t)

	_, err := container.ReadItem(context.Background(), PartitionKey{}, "item-1", nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, errDriverUnavailable)

	_, err = container.CreateItem(context.Background(), PartitionKey{}, "item-1", []byte(`{"id":"item-1"}`), nil)
	require.Error(t, err)
	require.NotErrorIs(t, err, errDriverUnavailable)
}

// An already-cancelled context must be honored rather than starting work that is bound to fail,
// and the caller must get the context's own error so errors.Is against context.Canceled works.
func TestItemOperationsHonorCancelledContext(t *testing.T) {
	container := newTestContainer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "item-1", nil)
	require.ErrorIs(t, err, context.Canceled)

	_, err = container.CreateItem(ctx, NewPartitionKeyString("pk"), "item-1", []byte(`{"id":"item-1"}`), nil)
	require.ErrorIs(t, err, context.Canceled)
}

// A well-formed call reaches the driver in a driver-backed build and reports that it is not
// implemented otherwise. Either way it fails as an *Error, so the documented errors.As idiom works
// whichever build a caller has.
func TestItemOperationsReportErrorsAsCosmosErrors(t *testing.T) {
	if driverAvailable {
		// The driver would try to reach the endpoint, which is what the emulator tests cover.
		t.Skip("driver-backed build: covered by the emulator tests")
	}

	container := newTestContainer(t)
	pk := NewPartitionKeyString("pk")

	_, readErr := container.ReadItem(context.Background(), pk, "item-1", nil)
	_, createErr := container.CreateItem(context.Background(), pk, "item-1", []byte(`{"id":"item-1"}`), nil)

	for _, err := range []error{readErr, createErr} {
		var cosmosErr *Error
		require.True(t, errors.As(err, &cosmosErr))
		require.Equal(t, CodeClientError, cosmosErr.Code)
	}
}

// The point of factoring OperationOptions out is that every operation takes the same driver-level
// settings without restating them, so a knob added there reaches all of them at once. This pins
// that both item operations carry it, which a per-type copy would not guarantee.
func TestItemOptionsShareOperationOptions(t *testing.T) {
	shared := OperationOptions{
		ConsistencyStrategy: ReadConsistencyStrategySession,
		ExcludedRegions:     []Region{RegionEastUS},
		EndToEndTimeout:     5 * time.Second,
	}

	read := ReadItemOptions{Operation: shared}
	create := CreateItemOptions{Operation: shared}

	require.Equal(t, shared, read.Operation)
	require.Equal(t, shared, create.Operation)
	require.Equal(t, []Region{RegionEastUS}, read.Operation.ExcludedRegions,
		"excluded regions are typed, not free strings")
}

// Close guarantees that later operations fail with CodeClientClosed. Shutdown is exactly when a
// caller's context is also likely to be cancelled, so the closed client has to win: checking the
// context first would report context.Canceled and hide the real problem.
func TestClosedClientReportedAheadOfCancelledContext(t *testing.T) {
	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := container.ReadItem(ctx, NewPartitionKeyString("pk"), "item-1", nil)
	var cosmosErr *Error
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
	require.NotErrorIs(t, err, context.Canceled)

	_, err = container.CreateItem(ctx, NewPartitionKeyString("pk"), "x", []byte(`{"id":"x"}`), nil)
	require.True(t, errors.As(err, &cosmosErr))
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

// The ABI distinguishes "inherit whatever is configured" from "use the default for this account's
// consistency level". Collapsing both onto the zero value would make one of them unreachable, and
// silently so: the read would just use the wrong strategy.
func TestReadConsistencyStrategyUnsetIsNotDefault(t *testing.T) {
	require.NotEqual(t, ReadConsistencyStrategyUnset, ReadConsistencyStrategyDefault)

	var zero ReadConsistencyStrategy
	require.Equal(t, ReadConsistencyStrategyUnset, zero, "the zero value must mean inherit")
}

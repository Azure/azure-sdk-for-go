// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

// TestBootstrapPrimarySuccess ports Rust scenario bootstrap.primary-success: a client built
// against the account's primary endpoint initializes and serves a trivial operation.
func TestBootstrapPrimarySuccess(t *testing.T) {
	fx := newFixture(t, "bootstrap.primary-success", "")
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString(id)
	body := mustMarshalItem(t, item(id, id, 1))

	created, err := fx.Container.CreateItem(context.Background(), pk, id, body, nil)
	require.NoError(t, err)
	trackItem(t, fx.Container, pk, id)
	require.Equal(t, 201, created.StatusCode)
}

// TestItemLifecycleSmoke ports the smokeTests-profile tier of Rust scenario item.lifecycle: one
// account-default create/read/replace/delete round trip. The full lifecycleConsistencyMatrix
// (five read-consistency-strategy cases) and readConsistencyOverrideMatrix are deferred; see
// internal/e2e/implementations.go's reasons for consistency.feed-read-strategies.
func TestItemLifecycleSmoke(t *testing.T) {
	fx := newFixture(t, "item.lifecycle", "")
	ctx := context.Background()
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString("A")

	created, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 1)), nil)
	require.NoError(t, err)
	trackItem(t, fx.Container, pk, id)
	require.Equal(t, 201, created.StatusCode)
	requireCriticalDiagnostics(t, created.Diagnostics, 201)

	read, err := fx.Container.ReadItem(ctx, pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, 200, read.StatusCode)
	require.Equal(t, item(id, "A", 1), mustUnmarshalItem(t, read.Value))

	replaced, err := fx.Container.ReplaceItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 2)), &azcosmos.ReplaceItemOptions{
		Operation: azcosmos.OperationOptions{EnableContentResponseOnWrite: to(true)},
	})
	require.NoError(t, err)
	require.Equal(t, 200, replaced.StatusCode)
	requireCriticalDiagnostics(t, replaced.Diagnostics, 200)
	require.Equal(t, item(id, "A", 2), mustUnmarshalItem(t, replaced.Value))

	deleted, err := fx.Container.DeleteItem(ctx, pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, 204, deleted.StatusCode)
	requireCriticalDiagnostics(t, deleted.Diagnostics, 204)

	_, err = fx.Container.ReadItem(ctx, pk, id, nil)
	requireCode(t, err, azcosmos.CodeNotFound)
}

// TestItemUpsertCreateThenUpdate ports Rust scenario item.upsert-create-update: UpsertItem
// creates an absent item, then updates it in place on a second call.
func TestItemUpsertCreateThenUpdate(t *testing.T) {
	fx := newFixture(t, "item.upsert-create-update", "")
	ctx := context.Background()
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString("A")

	created, err := fx.Container.UpsertItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 1)), nil)
	require.NoError(t, err)
	trackItem(t, fx.Container, pk, id)
	require.Equal(t, 201, created.StatusCode)

	updated, err := fx.Container.UpsertItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 2)), nil)
	require.NoError(t, err)
	require.Equal(t, 200, updated.StatusCode)

	read, err := fx.Container.ReadItem(ctx, pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, item(id, "A", 2), mustUnmarshalItem(t, read.Value))
}

// TestItemCreateConflictPreservesOriginal ports Rust scenario item.create-conflict: creating a
// duplicate id/partition-key leaves the original item untouched and reports CodeConflict with no
// sub-status. It exercises both a scalar (hash) partition key and a hierarchical one, reusing the
// pre-provisioned "items" and "query-hierarchical" containers respectively.
func TestItemCreateConflictPreservesOriginal(t *testing.T) {
	ctx := context.Background()

	t.Run("hash", func(t *testing.T) {
		fx := newFixture(t, "item.create-conflict", "")
		id := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString("A")

		_, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 1)), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, pk, id)

		_, err = fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 2)), nil)
		cosmosErr := requireCode(t, err, azcosmos.CodeConflict)
		require.Zero(t, cosmosErr.SubStatus)
		requireCriticalDiagnostics(t, cosmosErr.Diagnostics, 409)

		read, err := fx.Container.ReadItem(ctx, pk, id, nil)
		require.NoError(t, err)
		require.Equal(t, item(id, "A", 1), mustUnmarshalItem(t, read.Value))
	})

	t.Run("hierarchical", func(t *testing.T) {
		fx := newFixture(t, "item.create-conflict", "query-hierarchical")
		id := uniqueID(t)
		pk := azcosmos.NewPartitionKeyString("tenant-a").AppendString("user-1")
		body := func(value int64) []byte {
			v, err := jsonMarshalTenantItem(id, "tenant-a", "user-1", value)
			require.NoError(t, err)
			return v
		}

		_, err := fx.Container.CreateItem(ctx, pk, id, body(1), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, pk, id)

		_, err = fx.Container.CreateItem(ctx, pk, id, body(2), nil)
		cosmosErr := requireCode(t, err, azcosmos.CodeConflict)
		require.Zero(t, cosmosErr.SubStatus)
	})
}

// TestItemNotFoundDoesNotCrossPartitionKeys ports Rust scenario
// item.not-found-wrong-partition-key: reading an item under a partition key it was not created
// under reports CodeNotFound, the same as an id that was never created at all; the item under its
// real partition key is unaffected.
func TestItemNotFoundDoesNotCrossPartitionKeys(t *testing.T) {
	fx := newFixture(t, "item.not-found-wrong-partition-key", "")
	ctx := context.Background()
	id := uniqueID(t)
	ownerPK := azcosmos.NewPartitionKeyString("A")
	wrongPK := azcosmos.NewPartitionKeyString("B")

	_, err := fx.Container.CreateItem(ctx, ownerPK, id, mustMarshalItem(t, item(id, "A", 1)), nil)
	require.NoError(t, err)
	trackItem(t, fx.Container, ownerPK, id)

	_, err = fx.Container.ReadItem(ctx, wrongPK, id, nil)
	requireCode(t, err, azcosmos.CodeNotFound)

	_, err = fx.Container.ReadItem(ctx, ownerPK, uniqueID(t), nil)
	requireCode(t, err, azcosmos.CodeNotFound)

	read, err := fx.Container.ReadItem(ctx, ownerPK, id, nil)
	require.NoError(t, err)
	require.Equal(t, item(id, "A", 1), mustUnmarshalItem(t, read.Value))
}

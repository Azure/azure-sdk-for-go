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

// TestItemOptimisticConcurrencyStaleETag ports Rust scenario item.optimistic-concurrency: an
// IfMatchETag precondition built from the current ETag permits a replace, but reusing the
// now-stale ETag fails with CodePreconditionFailed and leaves the successful update intact.
func TestItemOptimisticConcurrencyStaleETag(t *testing.T) {
	fx := newFixture(t, "item.optimistic-concurrency", "")
	ctx := context.Background()
	id := uniqueID(t)
	pk := azcosmos.NewPartitionKeyString("A")

	created, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 1)), nil)
	require.NoError(t, err)
	trackItem(t, fx.Container, pk, id)
	initialETag := created.ETag

	replaced, err := fx.Container.ReplaceItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 2)), &azcosmos.ReplaceItemOptions{
		IfMatchETag: &initialETag,
	})
	require.NoError(t, err)
	require.Equal(t, 200, replaced.StatusCode)

	staleETag := initialETag
	_, err = fx.Container.ReplaceItem(ctx, pk, id, mustMarshalItem(t, item(id, "A", 3)), &azcosmos.ReplaceItemOptions{
		IfMatchETag: &staleETag,
	})
	cosmosErr := requireCode(t, err, azcosmos.CodePreconditionFailed)
	require.Zero(t, cosmosErr.SubStatus)
	requireCriticalDiagnostics(t, cosmosErr.Diagnostics, 412)

	read, err := fx.Container.ReadItem(ctx, pk, id, nil)
	require.NoError(t, err)
	require.Equal(t, item(id, "A", 2), mustUnmarshalItem(t, read.Value))
}

// TestItemScalarPartitionKeysRemainDistinct ports Rust scenario item.scalar-partition-keys: every
// scalar partition-key kind round-trips independently under the same shared item id without the
// kinds colliding with each other.
//
// Rust's version also covers "undefined" (a header value of {}, meaning "no partition key
// property"). That case is omitted here: the pinned in-memory emulator's own header parser
// (azure_data_cosmos_driver::in_memory_emulator::epk::json_to_pk_component) rejects any JSON
// object/array component - including the driver's own {} wire encoding for Undefined
// (partition_key.rs's Serialize impl emits "{}" for InnerPartitionKeyValue::Undefined) - with
// "partition key components must be scalar". This is a gap in the pinned emulator build, not a Go
// SDK defect; bumping the emulator pin is a separate decision (see catalog/README.md).
func TestItemScalarPartitionKeysRemainDistinct(t *testing.T) {
	fx := newFixture(t, "item.scalar-partition-keys", "")
	ctx := context.Background()
	const sharedID = "same-id"

	cases := []struct {
		name string
		pk   azcosmos.PartitionKey
		body []byte
	}{
		{"string", azcosmos.NewPartitionKeyString("A"), []byte(`{"id":"same-id","pk":"A","value":1}`)},
		{"number", azcosmos.NewPartitionKeyNumber(42), []byte(`{"id":"same-id","pk":42,"value":2}`)},
		{"bool", azcosmos.NewPartitionKeyBool(true), []byte(`{"id":"same-id","pk":true,"value":3}`)},
		{"null", azcosmos.NewPartitionKeyNull(), []byte(`{"id":"same-id","pk":null,"value":4}`)},
	}

	for _, c := range cases {
		created, err := fx.Container.CreateItem(ctx, c.pk, sharedID, c.body, nil)
		require.NoErrorf(t, err, "case %s", c.name)
		require.Equalf(t, 201, created.StatusCode, "case %s", c.name)
		t.Cleanup(func(pk azcosmos.PartitionKey, name string) func() {
			return func() {
				_, err := fx.Container.DeleteItem(context.Background(), pk, sharedID, nil)
				if err == nil {
					return
				}
				var cosmosErr *azcosmos.Error
				require.ErrorAsf(t, err, &cosmosErr, "case %s cleanup", name)
				require.Equalf(t, azcosmos.CodeNotFound, cosmosErr.Code, "case %s cleanup", name)
			}
		}(c.pk, c.name))
	}

	for _, c := range cases {
		read, err := fx.Container.ReadItem(ctx, c.pk, sharedID, nil)
		require.NoErrorf(t, err, "case %s", c.name)
		var stored map[string]any
		mustUnmarshalInto(t, read.Value, &stored)
		require.Equalf(t, sharedID, stored["id"], "case %s", c.name)
	}
}

// TestItemHierarchicalPartitionKeyPointAndPrefix ports Rust scenario
// item.hierarchical-partition-key: point reads/replaces/deletes address exactly one logical
// partition under a two-level key, and a partition-key prefix query returns every logical
// partition sharing the leading component without crossing into a sibling top-level value.
func TestItemHierarchicalPartitionKeyPointAndPrefix(t *testing.T) {
	fx := newFixture(t, "item.hierarchical-partition-key", "query-hierarchical")
	ctx := context.Background()
	tenant := uniqueID(t)

	firstID, secondID := tenant+"-1", tenant+"-2"
	firstPK := azcosmos.NewPartitionKeyString(tenant).AppendString("user-1")
	secondPK := azcosmos.NewPartitionKeyString(tenant).AppendString("user-2")

	for _, c := range []struct {
		pk    azcosmos.PartitionKey
		id    string
		child string
		value int64
	}{{firstPK, firstID, "user-1", 1}, {secondPK, secondID, "user-2", 2}} {
		body, err := jsonMarshalTenantItem(c.id, tenant, c.child, c.value)
		require.NoError(t, err)
		_, err = fx.Container.CreateItem(ctx, c.pk, c.id, body, nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, c.pk, c.id)
	}

	replaced, err := fx.Container.ReplaceItem(ctx, firstPK, firstID, func() []byte {
		b, err := jsonMarshalTenantItem(firstID, tenant, "user-1", 10)
		require.NoError(t, err)
		return b
	}(), nil)
	require.NoError(t, err)
	require.Equal(t, 200, replaced.StatusCode)

	read, err := fx.Container.ReadItem(ctx, firstPK, firstID, nil)
	require.NoError(t, err)
	require.Equal(t, int64(10), mustUnmarshalTenantItem(t, read.Value).Value)

	prefix := azcosmos.NewPartitionKeyString(tenant)
	pager := fx.Container.NewQueryItemsPager(
		azcosmos.NewQuery("SELECT VALUE c.id FROM c ORDER BY c.id ASC"),
		azcosmos.NewFeedScopeForPartitionKey(prefix), nil)
	defer func() { require.NoError(t, pager.Close()) }()
	var ids []string
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(t, err)
		for _, raw := range page.Items {
			var id string
			mustUnmarshalInto(t, raw, &id)
			ids = append(ids, id)
		}
	}
	require.ElementsMatch(t, []string{firstID, secondID}, ids)

	deleted, err := fx.Container.DeleteItem(ctx, firstPK, firstID, nil)
	require.NoError(t, err)
	require.Equal(t, 204, deleted.StatusCode)

	_, err = fx.Container.ReadItem(ctx, firstPK, firstID, nil)
	_ = requireCode(t, err, azcosmos.CodeNotFound)
}

func mustUnmarshalTenantItem(t *testing.T, body []byte) TenantItem {
	t.Helper()
	var value TenantItem
	mustUnmarshalInto(t, body, &value)
	return value
}

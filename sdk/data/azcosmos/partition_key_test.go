// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPartitionKeyConstructorsProduceOneComponent(t *testing.T) {
	for _, tt := range []struct {
		name string
		pk   PartitionKey
		want partitionKeyComponent
	}{
		{"string", NewPartitionKeyString("Contoso"), partitionKeyComponent{kind: partitionKeyKindString, stringValue: "Contoso"}},
		{"number", NewPartitionKeyNumber(98052), partitionKeyComponent{kind: partitionKeyKindNumber, numberValue: 98052}},
		{"bool", NewPartitionKeyBool(true), partitionKeyComponent{kind: partitionKeyKindBool, boolValue: true}},
		{"null", NewPartitionKeyNull(), partitionKeyComponent{kind: partitionKeyKindNull}},
		{"undefined", NewPartitionKeyUndefined(), partitionKeyComponent{kind: partitionKeyKindUndefined}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, 1, tt.pk.Len())
			require.Equal(t, tt.want, tt.pk.components[0])
		})
	}
}

// Null and undefined route differently in Cosmos DB, so the two must never collapse into the same
// component.
func TestPartitionKeyNullAndUndefinedAreDistinct(t *testing.T) {
	require.NotEqual(t, NewPartitionKeyNull().components[0], NewPartitionKeyUndefined().components[0])
}

func TestPartitionKeyEmpty(t *testing.T) {
	require.Equal(t, 0, NewPartitionKey().Len())
	require.Equal(t, 0, PartitionKey{}.Len(), "the zero value should be a usable empty key")
}

func TestPartitionKeyHierarchical(t *testing.T) {
	pk := NewPartitionKeyString("Contoso").
		AppendString("Redmond").
		AppendNumber(98052)

	require.Equal(t, 3, pk.Len())
	require.Equal(t, []partitionKeyComponent{
		{kind: partitionKeyKindString, stringValue: "Contoso"},
		{kind: partitionKeyKindString, stringValue: "Redmond"},
		{kind: partitionKeyKindNumber, numberValue: 98052},
	}, pk.components)
}

// Appending must not mutate the receiver, or a caller who holds a partially built key and extends
// it twice gets a key they never asked for. v1 appended in place, so two keys built from a shared
// prefix could overwrite each other's components.
func TestPartitionKeyAppendDoesNotMutateReceiver(t *testing.T) {
	prefix := NewPartitionKeyString("Contoso")

	redmond := prefix.AppendString("Redmond")
	seattle := prefix.AppendString("Seattle")

	require.Equal(t, 1, prefix.Len(), "the shared prefix must be unchanged")
	require.Equal(t, "Contoso", prefix.components[0].stringValue)

	require.Equal(t, "Redmond", redmond.components[1].stringValue)
	require.Equal(t, "Seattle", seattle.components[1].stringValue,
		"branches built from a shared prefix must not share a backing array")
}

// The same aliasing trap, reached through a longer prefix where the backing array has spare
// capacity from an earlier append.
func TestPartitionKeyBranchesFromGrownPrefixAreIndependent(t *testing.T) {
	prefix := NewPartitionKeyString("a").AppendString("b").AppendString("c")

	first := prefix.AppendString("first")
	second := prefix.AppendString("second")

	require.Equal(t, "first", first.components[3].stringValue)
	require.Equal(t, "second", second.components[3].stringValue)
	require.Equal(t, 3, prefix.Len())
}

func TestPartitionKeyAppendOnZeroValue(t *testing.T) {
	var pk PartitionKey

	appended := pk.AppendString("Contoso")

	require.Equal(t, 0, pk.Len(), "appending must not mutate the zero value")
	require.Equal(t, 1, appended.Len())
}

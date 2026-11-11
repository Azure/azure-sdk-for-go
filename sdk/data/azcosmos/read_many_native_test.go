// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadManyNativeRequestMapping(t *testing.T) {
	hierarchical := NewPartitionKeyString("tenant").AppendString("user")
	items, err := NewReadManyItemsSelection(
		NewItemIdentity(NewPartitionKeyString("pk"), "a"),
		NewItemIdentity(hierarchical, "b"),
	).WithFilter("c.v = @v").WithParameter("@v", map[string]int{"n": 1})
	require.NoError(t, err)
	req, err := newReadManyRequest(items, &ReadManyOptions{SessionToken: "0:1#2", PageSizeHint: 7, MaxFanOut: 3})
	require.NoError(t, err)
	native := inspectNativeReadMany(&req)
	require.Equal(t, uint32(readManySelectionItems), native.selectionKind)
	require.Equal(t, 2, native.identities)
	require.Equal(t, []string{"a", "b"}, native.itemIDs)
	require.Equal(t, []int{1, 2}, native.pkLens)
	require.Equal(t, "c.v = @v", native.filter)
	require.Equal(t, map[string]string{"@v": `{"n":1}`}, native.parameters)
	require.Equal(t, int32(7), native.pageSize)
	require.Equal(t, uint32(3), native.fanOut)
	require.Equal(t, "0:1#2", native.sessionToken)
	require.NotZero(t, native.abiVersion)
	require.NotZero(t, native.structSize)

	partitions, err := newReadManyRequest(NewReadManyPartitionsSelection(hierarchical), nil)
	require.NoError(t, err)
	native = inspectNativeReadMany(&partitions)
	require.Equal(t, uint32(readManySelectionPartitions), native.selectionKind)
	require.Equal(t, []string{""}, native.itemIDs, "partition selections carry no item id")
	require.Equal(t, []int{2}, native.pkLens)
	require.Empty(t, native.filter)
	require.Nil(t, native.parameters)
	require.Equal(t, int32(unsetMaxItemCount), native.pageSize)
	require.Zero(t, native.fanOut)
}

func TestNativeResourceNameValidation(t *testing.T) {
	require.Equal(t, int8(0), int8(nativeResourceNameValidation(nil)))
	require.Equal(t, int8(1), int8(nativeResourceNameValidation(to(false))))
	require.Equal(t, int8(2), int8(nativeResourceNameValidation(to(true))))
}

func TestNativeDriverOptionsBuildWithResourceNameValidation(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		d := &nativeDriver{cfg: driverConfig{options: ClientOptions{ResourceNameValidation: to(enabled)}}}
		require.NoError(t, d.buildRuntime())
		require.NoError(t, d.buildAccount(driverConfig{
			endpoint: "https://myaccount.documents.azure.com", accountKey: emulatorKey,
		}))
		options, err := d.buildDriverOptions()
		require.NoError(t, err)
		require.NotNil(t, options)
		freeNativeDriverOptions(options)
		require.NoError(t, d.close())
	}
}

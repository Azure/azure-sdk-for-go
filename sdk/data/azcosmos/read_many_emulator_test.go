// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type readManyEmulatorItem struct {
	ID     string `json:"id"`
	PK     string `json:"pk"`
	Status string `json:"status"`
}

// seedReadManyItems creates count items in each partition and returns their identities by
// partition. Partition key values are unique to the test, so partition selections see only
// these items.
func seedReadManyItems(t *testing.T, container *ContainerClient, partitions, count int) ([]string, map[string][]string) {
	t.Helper()
	base := uniqueItemID(t)
	pks := make([]string, partitions)
	ids := make(map[string][]string, partitions)
	for p := range pks {
		pks[p] = fmt.Sprintf("%s-pk%d", base, p)
		for i := range count {
			id := fmt.Sprintf("%s-item%d", pks[p], i)
			status := "active"
			if i%2 == 1 {
				status = "archived"
			}
			body, err := json.Marshal(readManyEmulatorItem{ID: id, PK: pks[p], Status: status})
			require.NoError(t, err)
			pk := NewPartitionKeyString(pks[p])
			_, err = container.CreateItem(t.Context(), pk, id, json.RawMessage(body), nil)
			require.NoError(t, err)
			trackEmulatorItem(t, container, pk, id)
			ids[pks[p]] = append(ids[pks[p]], id)
		}
	}
	return pks, ids
}

// drainReadMany reads every page, checks each page's metadata, and returns the sorted item IDs.
func drainReadMany(t *testing.T, pager *ReadManyItemsPager[readManyEmulatorItem]) ([]string, int) {
	t.Helper()
	t.Cleanup(func() { require.NoError(t, pager.Close()) })
	var ids []string
	pages := 0
	for pager.More() {
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		pages++
		require.NotNil(t, page.Diagnostics, "every page carries diagnostics")
		for _, item := range page.Items {
			ids = append(ids, item.ID)
		}
	}
	sort.Strings(ids)
	return ids, pages
}

func sortedCopy(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}

func TestEmulatorReadManyItems(t *testing.T) {
	container := emulatorContainer(t)
	pks, ids := seedReadManyItems(t, container, 2, 3)

	t.Run("skips missing items and deduplicates", func(t *testing.T) {
		pk0, pk1 := NewPartitionKeyString(pks[0]), NewPartitionKeyString(pks[1])
		selection := NewReadManyItemsSelection(
			NewItemIdentity(pk0, ids[pks[0]][0]),
			NewItemIdentity(pk0, ids[pks[0]][0]),
			NewItemIdentity(pk1, ids[pks[1]][2]),
			NewItemIdentity(pk0, "missing-"+ids[pks[0]][1]),
			// The right ID under the wrong partition key is a different, missing item.
			NewItemIdentity(pk1, ids[pks[0]][1]),
		)
		got, _ := drainReadMany(t, container.NewReadManyItemsPager[readManyEmulatorItem](selection, nil))
		require.Equal(t, sortedCopy([]string{ids[pks[0]][0], ids[pks[1]][2]}), got)
	})

	t.Run("empty selection reads nothing", func(t *testing.T) {
		got, _ := drainReadMany(t, container.NewReadManyItemsPager[readManyEmulatorItem](NewReadManyItemsSelection(), nil))
		require.Empty(t, got)
	})

	t.Run("response metadata", func(t *testing.T) {
		selection := NewReadManyItemsSelection(NewItemIdentity(NewPartitionKeyString(pks[0]), ids[pks[0]][0]))
		pager := container.NewReadManyItemsPager[readManyEmulatorItem](selection, nil)
		t.Cleanup(func() { require.NoError(t, pager.Close()) })
		var charge float64
		var items []readManyEmulatorItem
		for pager.More() {
			page, err := pager.NextPage(t.Context())
			require.NoError(t, err)
			require.NotNil(t, page.Diagnostics)
			charge += page.RequestCharge
			items = append(items, page.Items...)
		}
		require.Len(t, items, 1)
		require.Equal(t, readManyEmulatorItem{ID: ids[pks[0]][0], PK: pks[0], Status: "active"}, items[0])
		require.Positive(t, charge)
	})
}

func TestEmulatorReadManyPartitions(t *testing.T) {
	container := emulatorContainer(t)
	pks, ids := seedReadManyItems(t, container, 3, 4)
	all := append(append([]string(nil), ids[pks[0]]...), ids[pks[1]]...)

	selection := NewReadManyPartitionsSelection(NewPartitionKeyString(pks[0]), NewPartitionKeyString(pks[1]))

	t.Run("multiple partitions", func(t *testing.T) {
		got, _ := drainReadMany(t, container.NewReadManyItemsPager[readManyEmulatorItem](selection, nil))
		require.Equal(t, sortedCopy(all), got, "the unselected partition is not read")
	})

	t.Run("filter with parameter", func(t *testing.T) {
		filtered, err := selection.WithFilter("c.status = @status").WithParameter("@status", "archived")
		require.NoError(t, err)
		got, _ := drainReadMany(t, container.NewReadManyItemsPager[readManyEmulatorItem](filtered, nil))
		var want []string
		for _, pk := range pks[:2] {
			want = append(want, ids[pk][1], ids[pk][3])
		}
		require.Equal(t, sortedCopy(want), got)
	})

	t.Run("page size hint returns the same items", func(t *testing.T) {
		got, pages := drainReadMany(t, container.NewReadManyItemsPager[readManyEmulatorItem](selection,
			&ReadManyOptions{PageSizeHint: 1, MaxFanOut: 1}))
		require.Equal(t, sortedCopy(all), got)
		require.Greater(t, pages, 1, "a page size hint of one splits the results")
	})

	t.Run("concurrent pagers", func(t *testing.T) {
		const readers = 4
		results := make([][]string, readers)
		var wg sync.WaitGroup
		for i := range results {
			wg.Add(1)
			go func() {
				defer wg.Done()
				pager := container.NewReadManyItemsPager[readManyEmulatorItem](selection, nil)
				defer func() { _ = pager.Close() }()
				var got []string
				for pager.More() {
					page, err := pager.NextPage(t.Context())
					if err != nil {
						t.Errorf("reader %d: %v", i, err)
						return
					}
					for _, item := range page.Items {
						got = append(got, item.ID)
					}
				}
				sort.Strings(got)
				results[i] = got
			}()
		}
		wg.Wait()
		for i, got := range results {
			require.Equal(t, sortedCopy(all), got, "reader %d", i)
		}
	})
}

// Resource name validation is on by default, so an unsafe ID is rejected before a request is
// sent. Disabling it lets the request reach the service, which reports its own error instead.
func TestEmulatorResourceNameValidation(t *testing.T) {
	id := uniqueItemID(t) + "/unsafe"
	pk := NewPartitionKeyString("resource-name-validation")

	read := func(t *testing.T, options *ClientOptions) *Error {
		t.Helper()
		container := emulatorContainerWithOptions(t, options)
		_, err := container.ReadItem(t.Context(), pk, id, nil)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		return cosmosErr
	}

	t.Run("enabled by default", func(t *testing.T) {
		err := read(t, nil)
		require.False(t, err.FromWire, "rejected before the request is sent: %v", err)
		require.Equal(t, 400, err.StatusCode)
	})

	t.Run("disabled", func(t *testing.T) {
		disabled := false
		err := read(t, &ClientOptions{ResourceNameValidation: &disabled})
		require.True(t, err.FromWire, "the service answers when validation is disabled: %v", err)
	})
}

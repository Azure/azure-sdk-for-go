// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCursorPageCopiesAllItemsAndMetrics(t *testing.T) {
	for _, items := range [][][]byte{
		{},
		{[]byte(`{"id":"a"}`)},
		{[]byte("null"), {}, []byte(`9007199254740993`), []byte(`[1,true]`), []byte(`"text"`)},
	} {
		page, end, err := syntheticCursorPage(nil, items, 2, 2)
		require.NoError(t, err)
		require.False(t, end)
		require.Equal(t, items, page.Items)
		require.Equal(t, `{"UtilizedIndexes":[]}`, page.IndexMetrics)
		require.Equal(t, "retrievedDocumentCount=2", page.QueryMetrics)
	}
}

func TestCursorPageRepresentationsAndEnd(t *testing.T) {
	page, end, err := syntheticCursorPage([]byte(`{"Documents":[1,2]}`), nil, 1, 2)
	require.NoError(t, err)
	require.False(t, end)
	require.Equal(t, [][]byte{[]byte("1"), []byte("2")}, page.Items)
	_, end, err = syntheticCursorPage(nil, nil, 0, 4)
	require.NoError(t, err)
	require.True(t, end)
	for _, tc := range []struct {
		body         []byte
		items        [][]byte
		kind, result uint32
	}{
		{[]byte(`{`), nil, 1, 2},
		{nil, nil, 0, 2},
		{nil, nil, 3, 2},
		{nil, nil, 2, 1},
		{nil, [][]byte{[]byte("1")}, 2, 4},
	} {
		page, end, err := syntheticCursorPage(tc.body, tc.items, tc.kind, tc.result)
		require.Error(t, err)
		require.Zero(t, page)
		require.False(t, end)
	}
}

func TestCursorClientCloseAndConcurrentCleanup(t *testing.T) {
	container := newTestContainer(t)
	client := container.database.client
	driver := client.driver
	cursor, err := driver.testIdleCursor()
	require.NoError(t, err)
	require.NotNil(t, cursor.queue)
	pager := &QueryItemsPager{client: client, cursor: cursor}
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			require.NoError(t, pager.Close())
		}()
		go func() {
			defer wg.Done()
			require.NoError(t, client.Close())
		}()
	}
	wg.Wait()
	require.Nil(t, cursor.queue)
	require.Empty(t, driver.cursors)
	require.False(t, pager.More())
}

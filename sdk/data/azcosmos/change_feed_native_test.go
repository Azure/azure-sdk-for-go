// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNativeChangeFeedSelectors(t *testing.T) {
	instant := time.Date(2026, 1, 2, 3, 4, 5, 123, time.UTC)
	for _, mode := range []struct {
		value  ChangeFeedMode
		native uint32
	}{
		{ChangeFeedModeLatestVersion, 1}, {ChangeFeedModeAllVersionsAndDeletes, 2},
	} {
		for _, start := range []struct {
			value     ChangeFeedStartFrom
			native    uint32
			timestamp string
		}{
			{NewChangeFeedStartFromBeginning(), 1, ""},
			{NewChangeFeedStartFromNow(), 2, ""},
			{NewChangeFeedStartFromPointInTime(instant), 3, "2026-01-02T03:04:05.000000123Z"},
		} {
			req, err := newChangeFeedRequest(NewFeedScopeForFullContainer(), start.value, &ChangeFeedOptions{
				Mode:         mode.value,
				Feed:         FeedOptions{PageSizeHint: 7, MaxFanOut: 3, ContinuationToken: "opaque checkpoint"},
				Operation:    OperationOptions{EndToEndTimeout: 2500 * time.Millisecond},
				SessionToken: "0:3",
			})
			require.NoError(t, err)
			native, err := inspectNativeFullChangeFeed(&req)
			require.NoError(t, err)
			require.True(t, native.full)
			require.Zero(t, native.kind, "change feed must not be submitted as a query")
			require.Equal(t, mode.native, native.mode)
			require.Equal(t, start.native, native.start)
			require.Equal(t, start.timestamp, native.timestamp)
			require.EqualValues(t, 7, native.pageSize)
			require.EqualValues(t, 3, native.fanOut)
			require.EqualValues(t, 2500, native.timeoutMillis)
			require.Equal(t, "opaque checkpoint", native.token)
			require.Equal(t, "0:3", native.session)
			require.Zero(t, native.indexMetrics)
			require.Zero(t, native.queryMetrics)
		}
	}
	req, err := newChangeFeedRequest(NewFeedScopeForFullContainer(), NewChangeFeedStartFromNow(), nil)
	require.NoError(t, err)
	native, err := inspectNativeFullChangeFeed(&req)
	require.NoError(t, err)
	require.EqualValues(t, 1, native.mode)
	require.EqualValues(t, unsetMaxItemCount, native.pageSize)
	require.Zero(t, native.fanOut)
	require.EqualValues(t, -1, native.timeoutMillis, "unset leaves the native timeout policy unchanged")
	require.Empty(t, native.token)
	require.Empty(t, native.session)
}

func TestNativeChangeFeedPageRepresentationsAndOwnership(t *testing.T) {
	for _, tt := range []struct {
		name   string
		body   []byte
		items  [][]byte
		kind   uint32
		status int
		want   [][]byte
	}{
		{"idle no payload", nil, nil, 0, 304, nil},
		{"idle raw bytes", nil, nil, 1, 304, nil},
		{"empty envelope", []byte(`{"Documents":[]}`), nil, 1, 200, [][]byte{}},
		{"raw envelopes", []byte(`{"Documents":[{"current":{"id":"a"},"previous":null},{"metadata":{"operationType":"delete"},"previous":{"id":"b"}}]}`), nil, 1, 200,
			[][]byte{[]byte(`{"current":{"id":"a"},"previous":null}`), []byte(`{"metadata":{"operationType":"delete"},"previous":{"id":"b"}}`)}},
		{"pre-split envelopes", nil, [][]byte{[]byte(`{"current":{"n":9007199254740993},"future":true}`), {}, []byte(`null`)}, 2, 200,
			[][]byte{[]byte(`{"current":{"n":9007199254740993},"future":true}`), {}, []byte(`null`)}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			page, err := syntheticChangeFeedCursorPage(tt.body, tt.items, tt.kind, 2, tt.status)
			require.NoError(t, err)
			require.Equal(t, tt.want, page.Items, "the fixture freed all native input before this assertion")
			require.Equal(t, tt.status, page.StatusCode)
			require.Equal(t, 3.5, page.RequestCharge)
			require.Equal(t, "native-change-activity", page.ActivityID)
			require.EqualValues(t, "native-position", page.ETag)
			require.Equal(t, SessionToken("0:2"), page.SessionToken)
		})
	}
}

func TestNativeChangeFeedRejectsInvalidPages(t *testing.T) {
	for _, tt := range []struct {
		body         []byte
		items        [][]byte
		kind, result uint32
	}{
		{[]byte(`{`), nil, 1, 2},
		{[]byte(`{"Documents":null}`), nil, 1, 2},
		{nil, nil, 3, 2},
		{[]byte(`{"Documents":[]}`), nil, 0, 2},
		{nil, [][]byte{[]byte(`1`)}, 0, 2},
		{nil, nil, 0, 4},
		{nil, nil, 2, 1},
	} {
		page, err := syntheticChangeFeedCursorPage(tt.body, tt.items, tt.kind, tt.result, 304)
		require.Zero(t, page)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, 304, cosmosErr.StatusCode)
		require.Equal(t, 3.5, cosmosErr.RequestCharge)
		require.Equal(t, "native-change-activity", cosmosErr.ActivityID)
		require.EqualValues(t, "native-position", cosmosErr.ETag)
		require.Equal(t, SessionToken("0:2"), cosmosErr.SessionToken)
	}
	for kind := range uint32(4) {
		page, err := syntheticInvalidChangeFeedItems(kind)
		require.Zero(t, page)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
	}
}

func TestNativeChangeFeedCancellationPreservesCompletionMetadata(t *testing.T) {
	for _, cause := range []error{context.Canceled, context.DeadlineExceeded} {
		err := syntheticChangeFeedCancellation(cause)
		require.ErrorIs(t, err, cause)
		require.Equal(t, CodeOperationCancelled, err.Code)
		require.Equal(t, 304, err.StatusCode)
		require.Equal(t, 3.5, err.RequestCharge)
		require.Equal(t, "native-change-activity", err.ActivityID)
		require.EqualValues(t, "native-position", err.ETag)
		require.Equal(t, SessionToken("0:2"), err.SessionToken)
	}
}

func TestNativeChangeFeedClientCloseAndConcurrentCleanup(t *testing.T) {
	client := newTestContainer(t).database.client
	driver := client.driver
	cursor, err := driver.testIdleCursor()
	require.NoError(t, err)
	pager := &ChangeFeedPager{client: client, cursor: &nativeChangeFeedCursor{cursor}}
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
	require.Nil(t, cursor.handle)
	require.Empty(t, driver.cursors)
	require.False(t, pager.More())
}

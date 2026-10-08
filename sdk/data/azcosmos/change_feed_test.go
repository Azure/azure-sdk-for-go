// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChangeFeedRequestDefaultsAndSnapshots(t *testing.T) {
	scope := NewFeedScopeForPartitionKey(NewPartitionKeyString("partition"))
	start := NewChangeFeedStartFromNow()
	req, err := newChangeFeedRequest(scope, start, nil)
	require.NoError(t, err)
	require.Equal(t, ChangeFeedModeLatestVersion, req.options.Mode)
	require.Equal(t, start, req.startFrom)
	require.Equal(t, scope.partitionKey, req.partitionKey)

	enabled := true
	options := &ChangeFeedOptions{
		Mode: ChangeFeedModeAllVersionsAndDeletes,
		Feed: FeedOptions{PageSizeHint: 5, MaxFanOut: 2, ContinuationToken: "opaque checkpoint"},
		Operation: OperationOptions{
			EnableContentResponseOnWrite: &enabled, ExcludedRegions: []Region{RegionEastUS},
		},
		SessionToken: "0:1",
	}
	req, err = newChangeFeedRequest(scope, start, options)
	require.NoError(t, err)
	enabled = false
	options.Operation.ExcludedRegions[0] = RegionWestUS
	options.Mode, options.Feed, options.SessionToken = "", FeedOptions{}, ""
	require.True(t, *req.options.Operation.EnableContentResponseOnWrite)
	require.Equal(t, []Region{RegionEastUS}, req.options.Operation.ExcludedRegions)
	require.Equal(t, ChangeFeedModeAllVersionsAndDeletes, req.options.Mode)
	require.Equal(t, FeedOptions{PageSizeHint: 5, MaxFanOut: 2, ContinuationToken: "opaque checkpoint"}, req.options.Feed)
	require.Equal(t, SessionToken("0:1"), req.options.SessionToken)
}

func TestChangeFeedStartAndScopeContracts(t *testing.T) {
	instant := time.Date(2026, 1, 2, 3, 4, 5, 123, time.FixedZone("offset", 3600))
	starts := []ChangeFeedStartFrom{
		NewChangeFeedStartFromBeginning(), NewChangeFeedStartFromNow(),
		NewChangeFeedStartFromPointInTime(instant),
	}
	require.Equal(t, time.UTC, starts[2].time.Location())
	require.True(t, instant.Equal(starts[2].time))
	prefix := NewPartitionKeyString("tenant").AppendNull()
	for _, scope := range []FeedScope{
		NewFeedScopeForPartitionKey(NewPartitionKeyString("partition")),
		NewFeedScopeForPartitionKey(prefix), NewFeedScopeForFullContainer(),
	} {
		for _, mode := range []ChangeFeedMode{ChangeFeedModeLatestVersion, ChangeFeedModeAllVersionsAndDeletes} {
			for _, start := range starts {
				options := &ChangeFeedOptions{Mode: mode, Feed: FeedOptions{ContinuationToken: "not parsed by Go"}}
				req, err := newChangeFeedRequest(scope, start, options)
				require.NoError(t, err, "mode/start service restrictions must not become Go routing policy")
				require.Equal(t, scope.fullContainer, req.fullContainer)
				require.Equal(t, start, req.startFrom)
				require.Equal(t, options.Feed.ContinuationToken, req.options.Feed.ContinuationToken)
			}
		}
	}
}

func TestChangeFeedPagerArgumentAndLifetimeOrdering(t *testing.T) {
	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	scope, now := NewFeedScopeForFullContainer(), NewChangeFeedStartFromNow()
	for _, tt := range []struct {
		name    string
		scope   FeedScope
		start   ChangeFeedStartFrom
		options *ChangeFeedOptions
		message string
	}{
		{"empty scope", FeedScope{}, now, nil, "partition key"},
		{"unset start", scope, ChangeFeedStartFrom{}, nil, "explicit start"},
		{"unset resume start", scope, ChangeFeedStartFrom{}, &ChangeFeedOptions{Feed: FeedOptions{ContinuationToken: "resume"}}, "explicit start"},
		{"invalid start", scope, ChangeFeedStartFrom{kind: 4}, nil, "explicit start"},
		{"negative year", scope, NewChangeFeedStartFromPointInTime(time.Date(-1, 1, 1, 0, 0, 0, 0, time.UTC)), nil, "RFC 3339"},
		{"large year", scope, NewChangeFeedStartFromPointInTime(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)), nil, "RFC 3339"},
		{"invalid mode", scope, now, &ChangeFeedOptions{Mode: "invalid"}, "mode"},
		{"negative hint", scope, now, &ChangeFeedOptions{Feed: FeedOptions{PageSizeHint: -1}}, "page size hint"},
		{"invalid token", scope, now, &ChangeFeedOptions{Feed: FeedOptions{ContinuationToken: "a\x00b"}}, "continuation token"},
		{"invalid consistency", scope, now, &ChangeFeedOptions{Operation: OperationOptions{ConsistencyStrategy: "invalid"}}, "consistency"},
		{"invalid session", scope, now, &ChangeFeedOptions{SessionToken: "a\x00b"}, "session"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pager := container.NewChangeFeedPager(tt.scope, tt.start, tt.options)
			page, err := pager.NextPage(ctx)
			require.Zero(t, page)
			require.ErrorContains(t, err, tt.message)
			require.NotErrorIs(t, err, context.Canceled)
			requireNotDriverUnavailable(t, err)
			token, err := pager.ContinuationToken(ctx)
			require.Empty(t, token)
			require.ErrorContains(t, err, tt.message)
			require.NoError(t, pager.Close())
		})
	}
	page, err := container.NewChangeFeedPager(scope, now, nil).NextPage(ctx)
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
	page, err = newTestContainer(t).NewChangeFeedPager(scope, now, nil).NextPage(ctx)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)
}

type testChangeFeedCursor struct {
	calls         int
	checkpoints   int
	closed        int
	nextErr       error
	checkpointErr error
	entered       chan struct{}
	proceed       chan struct{}
}

func (c *testChangeFeedCursor) next(context.Context) (ChangeFeedResponse, error) {
	c.calls++
	if c.entered != nil {
		close(c.entered)
		<-c.proceed
	}
	if c.nextErr != nil {
		return ChangeFeedResponse{}, c.nextErr
	}
	if c.calls <= 2 {
		return ChangeFeedResponse{Response: Response{StatusCode: 304}, ETag: "idle-position"}, nil
	}
	return ChangeFeedResponse{Items: [][]byte{[]byte(`{"current":{"id":"change"}}`)}}, nil
}

func (c *testChangeFeedCursor) checkpoint(context.Context) (string, error) {
	c.checkpoints++
	if c.checkpointErr != nil {
		return "", c.checkpointErr
	}
	return fmt.Sprintf("checkpoint-%d", c.calls), nil
}

func (c *testChangeFeedCursor) close() { c.closed++ }

func TestChangeFeedPagerIdleAndCheckpointContracts(t *testing.T) {
	cursor := &testChangeFeedCursor{}
	pager := &ChangeFeedPager{client: newTestContainer(t).database.client, cursor: cursor}
	for range 2 {
		page, err := pager.NextPage(t.Context())
		require.NoError(t, err)
		require.Equal(t, 304, page.StatusCode)
		require.EqualValues(t, "idle-position", page.ETag)
		require.Empty(t, page.Items)
		require.True(t, pager.More(), "idle pages must remain pollable")
	}
	for range 2 {
		token, err := pager.ContinuationToken(t.Context())
		require.NoError(t, err)
		require.Equal(t, "checkpoint-2", token)
		require.Equal(t, 2, cursor.calls, "checkpoint must not fetch changes")
	}
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 1)
	require.True(t, pager.More())
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, "checkpoint-3", token)
	require.NoError(t, pager.Close())
	require.NoError(t, pager.Close())
	require.False(t, pager.More())
	require.Equal(t, 1, cursor.closed)
	page, err = pager.NextPage(t.Context())
	require.Error(t, err)
	require.Zero(t, page)
	token, err = pager.ContinuationToken(t.Context())
	require.Error(t, err)
	require.Empty(t, token)
	require.Equal(t, 3, cursor.calls)
}

func TestChangeFeedPagerCancellationAndTerminalFailure(t *testing.T) {
	for _, failure := range []error{context.Canceled, context.DeadlineExceeded, &Error{Code: CodeBadRequest, FromWire: true}} {
		for _, checkpoint := range []bool{false, true} {
			cursor := &testChangeFeedCursor{}
			pager := &ChangeFeedPager{client: newTestContainer(t).database.client, cursor: cursor}
			if checkpoint {
				cursor.checkpointErr = failure
				token, err := pager.ContinuationToken(t.Context())
				require.ErrorIs(t, err, failure)
				require.Empty(t, token)
			} else {
				cursor.nextErr = failure
				page, err := pager.NextPage(t.Context())
				require.ErrorIs(t, err, failure)
				require.Zero(t, page)
			}
			require.False(t, pager.More())
			require.Equal(t, 1, cursor.closed)
			require.NoError(t, pager.Close())
			require.Equal(t, 1, cursor.closed)
		}
	}
	cursor := &testChangeFeedCursor{}
	pager := &ChangeFeedPager{client: newTestContainer(t).database.client, cursor: cursor}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	page, err := pager.NextPage(ctx)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)
	token, err := pager.ContinuationToken(ctx)
	require.Empty(t, token)
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, cursor.calls)
	require.Zero(t, cursor.checkpoints)
	require.True(t, pager.More(), "pre-admission cancellation must not lose progress")
	require.NoError(t, pager.Close())

	unopened := newTestContainer(t).NewChangeFeedPager(NewFeedScopeForFullContainer(), NewChangeFeedStartFromNow(), nil)
	token, err = unopened.ContinuationToken(t.Context())
	require.Empty(t, token)
	require.Error(t, err)
	require.True(t, unopened.More())
	require.NoError(t, unopened.Close())
}

func TestChangeFeedCloseWaitsForActiveCall(t *testing.T) {
	client := newTestContainer(t).database.client
	cursor := &testChangeFeedCursor{entered: make(chan struct{}), proceed: make(chan struct{})}
	pager := &ChangeFeedPager{client: client, cursor: cursor}
	next := make(chan error, 1)
	go func() {
		_, err := pager.NextPage(t.Context())
		next <- err
	}()
	<-cursor.entered
	closedPager, closedClient := make(chan error, 1), make(chan error, 1)
	go func() { closedPager <- pager.Close() }()
	go func() { closedClient <- client.Close() }()
	select {
	case err := <-closedPager:
		t.Fatalf("pager closed during an active call: %v", err)
	case err := <-closedClient:
		t.Fatalf("client closed during an active call: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(cursor.proceed)
	require.NoError(t, <-next)
	require.NoError(t, <-closedPager)
	require.NoError(t, <-closedClient)
	require.Equal(t, 1, cursor.closed)
	require.False(t, pager.More())
}

func TestChangeFeedPreservesRawEnvelopesAndMetadata(t *testing.T) {
	items := []string{
		`{"current":{"id":"one","large":9007199254740993},"metadata":{"operationType":"create"},"future":[1,true]}`,
		`{"current":{"id":"two"},"previous":{"id":"two","value":1},"metadata":{"operationType":"replace","lsn":9007199254740993}}`,
		`{"previous":{"id":"three"},"metadata":{"operationType":"delete"},"current":null}`,
	}
	body := []byte(`{"Documents":[` + items[0] + `,` + items[1] + `,` + items[2] + `],"_count":3}`)
	response := ChangeFeedResponse{
		Response: Response{StatusCode: 200, SubStatus: 1, AttemptCount: 2, RequestCharge: 3.5,
			ActivityID: "change-activity", Diagnostics: &Diagnostics{StatusCode: 200, AttemptCount: 2}},
		ETag: "position", SessionToken: "0:2",
	}
	page, err := decodeChangeFeedPage(body, response)
	require.NoError(t, err)
	require.Equal(t, response.Response, page.Response)
	require.Equal(t, response.ETag, page.ETag)
	require.Equal(t, response.SessionToken, page.SessionToken)
	require.Len(t, page.Items, 3)
	clear(body)
	for i, item := range page.Items {
		require.Equal(t, items[i], string(item), "raw envelopes must not borrow or unwrap the page")
	}
	for _, body := range []string{"", `{"Documents":[]}`} {
		page, err := decodeChangeFeedPage([]byte(body), response)
		require.NoError(t, err)
		require.Empty(t, page.Items)
		require.Equal(t, response.Response, page.Response)
	}
	for _, body := range []string{`{`, `null`, `{}`, `{"Documents":null}`, `{"Documents":{}}`, `[]`} {
		page, err := decodeChangeFeedPage([]byte(body), response)
		require.Zero(t, page)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
		require.Equal(t, response.RequestCharge, cosmosErr.RequestCharge)
		require.Equal(t, response.ActivityID, cosmosErr.ActivityID)
		require.Equal(t, response.StatusCode, cosmosErr.StatusCode)
		require.Equal(t, response.SubStatus, cosmosErr.SubStatus)
		require.Equal(t, response.AttemptCount, cosmosErr.AttemptCount)
		require.Equal(t, response.ETag, cosmosErr.ETag)
		require.Equal(t, response.SessionToken, cosmosErr.SessionToken)
		require.Same(t, response.Diagnostics, cosmosErr.Diagnostics)
		require.NotNil(t, errors.Unwrap(err))
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryOwnsParameterSnapshots(t *testing.T) {
	value := map[string]any{"values": []int{1, 2}}
	base := NewQuery("SELECT VALUE @value")
	first, err := base.WithParameter("@value", value)
	require.NoError(t, err)
	value["values"].([]int)[0] = 99
	value["extra"] = true
	second, err := first.WithParameter("@value", nil)
	require.NoError(t, err)
	third, err := first.WithParameter("@other", json.RawMessage(`9007199254740993`))
	require.NoError(t, err)

	for _, tt := range []struct {
		query Query
		body  string
	}{
		{base, `{"query":"SELECT VALUE @value"}`},
		{first, `{"query":"SELECT VALUE @value","parameters":[{"name":"@value","value":{"values":[1,2]}}]}`},
		{second, `{"query":"SELECT VALUE @value","parameters":[{"name":"@value","value":null}]}`},
		{third, `{"query":"SELECT VALUE @value","parameters":[{"name":"@value","value":{"values":[1,2]}},{"name":"@other","value":9007199254740993}]}`},
	} {
		body, err := tt.query.body()
		require.NoError(t, err)
		require.Equal(t, tt.body, string(body))
	}
}

func TestQueryParameterValidation(t *testing.T) {
	for _, name := range []string{"", "@", "value", "@with space", "@with\nnewline", "@nul\x00"} {
		_, err := NewQuery("SELECT * FROM c").WithParameter(name, 1)
		require.Error(t, err, name)
	}
	for _, value := range []any{math.NaN(), math.Inf(1), make(chan int), json.RawMessage(`{`)} {
		_, err := NewQuery("SELECT VALUE @v").WithParameter("@v", value)
		var cosmosErr *Error
		require.ErrorAs(t, err, &cosmosErr)
		require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
		require.NotNil(t, errors.Unwrap(err))
	}
}

func TestQueryRequestOwnsOptions(t *testing.T) {
	enabled := true
	options := &QueryOptions{
		Operation: OperationOptions{
			EnableContentResponseOnWrite: &enabled,
			ExcludedRegions:              []Region{RegionEastUS},
			CustomHeaders:                map[string]string{"x-test": "original"},
			BinaryEncoding:               &BinaryEncodingOptions{Enabled: to(true), RequestTextResponse: true},
			MaxFailoverRetryCount:        to(uint32(3)),
			ThrottlingRetry:              ThrottlingRetryOptions{MaxRetryCount: to(uint32(1))},
			ThroughputControl:            ThroughputControlOptions{ThroughputBucket: to(uint32(2))},
		},
		Feed:                 FeedOptions{PageSizeHint: 5, ContinuationToken: "resume"},
		SessionToken:         SessionToken("0:1"),
		PopulateIndexMetrics: &enabled,
		PopulateQueryMetrics: &enabled,
	}
	req, err := newQueryRequest(NewQuery("SELECT * FROM c"), NewFeedScopeForPartitionKey(NewPartitionKeyNull()), options)
	require.NoError(t, err)
	enabled = false
	options.Operation.ExcludedRegions[0] = RegionWestUS
	options.Operation.CustomHeaders["x-test"] = "changed"
	options.Operation.BinaryEncoding.Enabled = to(false)
	*options.Operation.MaxFailoverRetryCount = 9
	*options.Operation.ThrottlingRetry.MaxRetryCount = 9
	*options.Operation.ThroughputControl.ThroughputBucket = 9
	options.Feed = FeedOptions{}
	options.SessionToken = ""
	require.True(t, *req.options.Operation.EnableContentResponseOnWrite)
	require.True(t, *req.options.PopulateIndexMetrics)
	require.True(t, *req.options.PopulateQueryMetrics)
	require.Equal(t, []Region{RegionEastUS}, req.options.Operation.ExcludedRegions)
	require.Equal(t, "original", req.options.Operation.CustomHeaders["x-test"])
	require.True(t, req.options.Operation.BinaryEncoding.enabled())
	require.Equal(t, uint32(3), *req.options.Operation.MaxFailoverRetryCount)
	require.Equal(t, uint32(1), *req.options.Operation.ThrottlingRetry.MaxRetryCount)
	require.Equal(t, uint32(2), *req.options.Operation.ThroughputControl.ThroughputBucket)
	require.Equal(t, FeedOptions{PageSizeHint: 5, ContinuationToken: "resume"}, req.options.Feed)
	require.Equal(t, SessionToken("0:1"), req.options.SessionToken)
}

func TestQueryPagerArgumentAndLifetimeOrdering(t *testing.T) {
	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	scope := NewFeedScopeForPartitionKey(NewPartitionKeyString("partition"))
	for _, tt := range []struct {
		name    string
		query   Query
		scope   FeedScope
		options *QueryOptions
		message string
	}{
		{"empty query", Query{}, scope, nil, "query text"},
		{"blank query", NewQuery(" \t"), scope, nil, "query text"},
		{"empty scope", NewQuery("SELECT * FROM c"), FeedScope{}, nil, "partition key"},
		{"negative hint", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Feed: FeedOptions{PageSizeHint: -1}}, "page size hint"},
		{"invalid token", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Feed: FeedOptions{ContinuationToken: "a\x00b"}}, "continuation token"},
		{"invalid consistency", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Operation: OperationOptions{ConsistencyStrategy: "invalid"}}, "consistency"},
		{"invalid timeout", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Operation: OperationOptions{EndToEndTimeout: to(time.Duration(-1))}}, "EndToEndTimeout"},
		{"invalid header", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Operation: OperationOptions{CustomHeaders: map[string]string{"x-test": "\r\n"}}}, "custom header"},
		{"invalid availability", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Operation: OperationOptions{AvailabilityStrategy: HedgingAvailability(0)}}, "positive threshold"},
		{"raw binary", NewQuery("SELECT * FROM c"), scope, &QueryOptions{Operation: OperationOptions{BinaryEncoding: &BinaryEncodingOptions{Enabled: to(true)}}}, "queries require text JSON"},
		{"invalid session", NewQuery("SELECT * FROM c"), scope, &QueryOptions{SessionToken: "a\x00b"}, "session"},
		{"invalid plan mode", NewQuery("SELECT * FROM c"), scope, &QueryOptions{QueryPlanMode: "invalid"}, "query plan mode"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pager := container.NewQueryItemsPager(tt.query, tt.scope, tt.options)
			page, err := pager.NextPage(ctx)
			require.ErrorContains(t, err, tt.message)
			require.Zero(t, page)
			require.NotErrorIs(t, err, context.Canceled)
			requireNotDriverUnavailable(t, err)
		})
	}
	page, err := container.NewQueryItemsPager(NewQuery("SELECT * FROM c"), scope, nil).NextPage(ctx)
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	open := newTestContainer(t)
	page, err = open.NewQueryItemsPager(NewQuery("SELECT * FROM c"), scope, nil).NextPage(ctx)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)
}

type testQueryCursor struct {
	calls         int
	closed        int
	checkpointErr error
	nextErr       error
}

func (c *testQueryCursor) next(context.Context) (QueryItemsResponse, bool, error) {
	c.calls++
	if c.nextErr != nil {
		return QueryItemsResponse{}, false, c.nextErr
	}
	switch c.calls {
	case 1:
		return QueryItemsResponse{}, false, nil
	case 2:
		return QueryItemsResponse{Items: [][]byte{[]byte("1"), []byte("2")}}, false, nil
	default:
		return QueryItemsResponse{}, true, nil
	}
}

func (c *testQueryCursor) checkpoint(context.Context) (string, error) {
	return "snapshot", c.checkpointErr
}

func (c *testQueryCursor) close() { c.closed++ }

func TestQueryPagerRetainsProgressWithoutCheckpoint(t *testing.T) {
	cursor := &testQueryCursor{checkpointErr: &Error{Code: CodeBadRequest, StatusCode: 400, SubStatus: subStatusBufferedQueryContinuationUnsupported}}
	pager := &QueryItemsPager{client: newTestContainer(t).database.client, cursor: cursor}
	require.True(t, pager.More())
	empty, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Empty(t, empty.Items)
	require.True(t, pager.More(), "an empty page must not truncate the query")
	_, err = pager.ContinuationToken(t.Context())
	require.ErrorIs(t, err, cursor.checkpointErr)
	require.True(t, pager.More())
	page, err := pager.NextPage(t.Context())
	require.NoError(t, err)
	require.Len(t, page.Items, 2)
	_, err = pager.NextPage(t.Context())
	require.NoError(t, err)
	require.False(t, pager.More())
	_, err = pager.NextPage(t.Context())
	require.Error(t, err)
	require.Equal(t, 3, cursor.calls)
	require.Equal(t, 1, cursor.closed)
	require.NoError(t, pager.Close())
	require.Equal(t, 1, cursor.closed)
}

func TestUnsupportedQueryCheckpoint(t *testing.T) {
	require.Equal(t, 20124, subStatusBufferedQueryContinuationUnsupported)
	require.Equal(t, 20117, subStatusContinuationTokenNonQueryOperation)
	for _, tt := range []struct {
		name string
		err  error
		want bool
	}{
		{"buffered query", &Error{StatusCode: 400, SubStatus: subStatusBufferedQueryContinuationUnsupported}, true},
		{"non-query operation", &Error{StatusCode: 400, SubStatus: subStatusContinuationTokenNonQueryOperation}, true},
		{"wrapped", fmt.Errorf("checkpoint: %w", &Error{StatusCode: 400, SubStatus: subStatusBufferedQueryContinuationUnsupported}), true},
		{"buffered wire error", &Error{StatusCode: 400, SubStatus: subStatusBufferedQueryContinuationUnsupported, FromWire: true}, false},
		{"non-query wire error", &Error{StatusCode: 400, SubStatus: subStatusContinuationTokenNonQueryOperation, FromWire: true}, false},
		{"buffered wrong status", &Error{StatusCode: 503, SubStatus: subStatusBufferedQueryContinuationUnsupported}, false},
		{"non-query wrong status", &Error{StatusCode: 503, SubStatus: subStatusContinuationTokenNonQueryOperation}, false},
		{"unknown substatus", &Error{StatusCode: 400, SubStatus: 9999}, false},
		{"other error", errors.New("checkpoint failed"), false},
		{"nil", nil, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, unsupportedQueryCheckpoint(tt.err))
		})
	}
}

func TestQueryPagerTerminalFailureClosesCursor(t *testing.T) {
	cursor := &testQueryCursor{nextErr: &Error{Code: CodeBadRequest}}
	pager := &QueryItemsPager{client: newTestContainer(t).database.client, cursor: cursor}
	page, err := pager.NextPage(t.Context())
	require.ErrorIs(t, err, cursor.nextErr)
	require.Zero(t, page)
	require.False(t, pager.More())
	require.Equal(t, 1, cursor.closed)
}

func TestQueryPagerCheckpointAndCloseContracts(t *testing.T) {
	for _, failure := range []error{context.Canceled, &Error{Code: CodeClientError}} {
		cursor := &testQueryCursor{checkpointErr: failure}
		pager := &QueryItemsPager{client: newTestContainer(t).database.client, cursor: cursor}
		_, err := pager.ContinuationToken(t.Context())
		require.ErrorIs(t, err, failure)
		require.False(t, pager.More())
		require.Equal(t, 1, cursor.closed)
	}
	cursor := &testQueryCursor{}
	pager := &QueryItemsPager{client: newTestContainer(t).database.client, cursor: cursor}
	token, err := pager.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.Equal(t, "snapshot", token)
	require.Zero(t, cursor.calls, "checkpoint must not advance the cursor")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = pager.NextPage(ctx)
	require.ErrorIs(t, err, context.Canceled)
	require.True(t, pager.More(), "pre-admission cancellation does not lose a page")
	require.Zero(t, cursor.calls)
	require.NoError(t, pager.Close())
	require.False(t, pager.More())
	_, err = pager.ContinuationToken(t.Context())
	require.Error(t, err)
}

func TestQueryPagerRejectsNilContextsWithoutAdvancing(t *testing.T) {
	cursor := &testQueryCursor{}
	pager := &QueryItemsPager{client: newTestClient(t), cursor: cursor}
	defer func() { require.NoError(t, pager.Close()) }()
	page, err := pager.NextPage(nil) //nolint:staticcheck // verifies the boundary guard
	require.ErrorContains(t, err, "context must not be nil")
	require.Zero(t, page)
	token, err := pager.ContinuationToken(nil) //nolint:staticcheck // verifies the boundary guard
	require.ErrorContains(t, err, "context must not be nil")
	require.Empty(t, token)
	require.True(t, pager.More())
	require.Zero(t, cursor.calls)
}

func TestDecodeQueryPage(t *testing.T) {
	response := Response{RequestCharge: 3.5, ActivityID: "query-activity"}
	body := []byte(`{"Documents":[{"id":"first"},9007199254740993,null,[1,true],"text"],"_count":5}`)
	page, err := decodeQueryPage(body, response, "0:2", "planner-token", false)
	require.NoError(t, err)
	require.Equal(t, response, page.Response)
	require.Equal(t, SessionToken("0:2"), page.SessionToken)
	require.Equal(t, [][]byte{[]byte(`{"id":"first"}`), []byte(`9007199254740993`), []byte(`null`), []byte(`[1,true]`), []byte(`"text"`)}, page.Items)
	clear(body)
	require.Equal(t, `{"id":"first"}`, string(page.Items[0]), "items must not borrow the envelope")

	for _, body := range []string{`{`, `null`, `{}`, `{"Documents":null}`, `{"Documents":{}}`, `{"id":"first"}`, `[]`, ``} {
		t.Run(body, func(t *testing.T) {
			page, err := decodeQueryPage([]byte(body), response, "0:2", "next", false)
			require.Zero(t, page)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
			require.Equal(t, 3.5, cosmosErr.RequestCharge)
			require.Equal(t, "query-activity", cosmosErr.ActivityID)
			require.Equal(t, SessionToken("0:2"), cosmosErr.SessionToken)
			require.False(t, cosmosErr.FromWire)
		})
	}
	page, err = decodeQueryPage([]byte(`{"Documents":[]}`), response, "", "more", false)
	require.NoError(t, err)
	require.Empty(t, page.Items)
	page, err = decodeQueryPage(nil, Response{}, "", "", true)
	require.NoError(t, err)
	require.Empty(t, page.Items)
	_, err = decodeQueryPage([]byte(`{"Documents":[]}`), response, "", "", true)
	require.Error(t, err)
	_, err = decodeQueryPage(nil, response, "", "more", true)
	require.Error(t, err)
}

func TestQueryPartitionScopeValidation(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		key  PartitionKey
		ok   bool
	}{
		{"hash", `{"partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`, NewPartitionKeyNull(), true},
		{"legacy hash", `{"partitionKey":{"paths":["/pk"],"kind":"Hash"}}`, NewPartitionKeyUndefined(), true},
		{"hierarchical", `{"partitionKey":{"paths":["/a","/b"],"kind":"MultiHash","version":2}}`, NewPartitionKeyString("a").AppendNumber(2), true},
		{"prefix", `{"partitionKey":{"paths":["/a","/b"],"kind":"MultiHash","version":2}}`, NewPartitionKeyString("a"), true},
		{"too long", `{"partitionKey":{"paths":["/pk"],"kind":"Hash","version":2}}`, NewPartitionKeyString("a").AppendNull(), false},
		{"range", `{"partitionKey":{"paths":["/pk"],"kind":"Range","version":1}}`, NewPartitionKeyString("a"), false},
		{"invalid multihash", `{"partitionKey":{"paths":["/a","/b"],"kind":"MultiHash","version":1}}`, NewPartitionKeyString("a").AppendNull(), false},
		{"unknown version", `{"partitionKey":{"paths":["/pk"],"kind":"Hash","version":3}}`, NewPartitionKeyString("a"), false},
		{"missing definition", `{}`, NewPartitionKeyString("a"), false},
		{"invalid JSON", `{`, NewPartitionKeyString("a"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := validateQueryPartitionKey([]byte(tt.body), tt.key)
			if tt.ok {
				require.NoError(t, err)
			} else {
				var cosmosErr *Error
				require.ErrorAs(t, err, &cosmosErr)
			}
		})
	}
}

func TestQuerySetupChargeOnErrors(t *testing.T) {
	original := &Error{Code: CodeThrottled, RequestCharge: 2, ActivityID: "query", StatusCode: 429, FromWire: true}
	var got *Error
	require.ErrorAs(t, addQuerySetupCharge(original, Response{RequestCharge: 3}), &got)
	require.Equal(t, 5.0, got.RequestCharge)
	require.Equal(t, 2.0, original.RequestCharge)
	require.Equal(t, 429, got.StatusCode)
	require.True(t, got.FromWire)
	require.Equal(t, "query", got.ActivityID)

	err := addQuerySetupCharge(context.DeadlineExceeded, Response{RequestCharge: 3, ActivityID: "metadata"})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.ErrorAs(t, err, &got)
	require.Equal(t, CodeOperationCancelled, got.Code)
	require.Equal(t, 3.0, got.RequestCharge)
}

func TestQueryDiagnosticBuild(t *testing.T) {
	if driverAvailable {
		t.Skip("requires the diagnostic build")
	}
	page, err := newTestContainer(t).NewQueryItemsPager(NewQuery("SELECT * FROM c"),
		NewFeedScopeForPartitionKey(NewPartitionKeyString("pk")), nil).NextPage(t.Context())
	require.Zero(t, page)
	require.ErrorContains(t, err, "cannot reach the Cosmos driver")
}

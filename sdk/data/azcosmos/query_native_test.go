// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQuerySnapshotTimeoutIncludesLazyInitialization(t *testing.T) {
	shared, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: to(time.Duration(time.Second))}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	credential := &blockingTokenCredential{started: make(chan struct{}), stopped: make(chan struct{})}
	client, err := NewClient("https://myaccount.documents.azure.com", credential, &ClientOptions{Runtime: shared})
	require.NoError(t, err)
	t.Cleanup(func() { client.driver.tokenProvider.cancel(); require.NoError(t, client.Close()) })
	container, err := client.NewContainer("db", "items")
	require.NoError(t, err)
	pager := container.NewQueryItemsPager(NewQuery("SELECT * FROM c"),
		NewFeedScopeForPartitionKey(NewPartitionKeyString("pk")), nil)
	done := make(chan error, 1)
	start := time.Now()
	go func() {
		page, err := pager.NextPage(context.Background())
		if len(page.Items) != 0 {
			done <- errors.New("failed query returned items")
			return
		}
		done <- err
	}()
	select {
	case <-credential.started:
	case <-time.After(3 * time.Second):
		t.Fatal("query did not reach lazy initialization")
	}
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Less(t, time.Since(start), 3*time.Second)
	case <-time.After(3 * time.Second):
		t.Fatal("query did not honor its runtime budget")
	}
}

func TestQueryNativeOptions(t *testing.T) {
	yes, no := true, false
	req, err := newQueryRequest(NewQuery("SELECT * FROM c"), NewFeedScopeForFullContainer(), &QueryOptions{
		Operation:     OperationOptions{EndToEndTimeout: to(5 * time.Second)},
		Feed:          FeedOptions{MaxFanOut: 250, PageSizeHint: 17},
		QueryPlanMode: QueryPlanModeGatewayOnly, PopulateIndexMetrics: &yes, PopulateQueryMetrics: &no,
	})
	require.NoError(t, err)
	got, err := inspectNativeFullQuery(&req)
	require.NoError(t, err)
	require.Equal(t, nativeQueryOptions{full: true, fanOut: 250, pageSize: 17, mode: 2, indexMetrics: 2, queryMetrics: 1, timeoutMillis: 5000}, got)

	req, err = newQueryRequest(NewQuery("SELECT * FROM c"), NewFeedScopeForFullContainer(), nil)
	require.NoError(t, err)
	got, err = inspectNativeFullQuery(&req)
	require.NoError(t, err)
	require.True(t, got.full)
	require.Zero(t, got.fanOut)
	require.EqualValues(t, -1, got.pageSize)
	require.Zero(t, got.mode)
	require.Zero(t, got.indexMetrics)
	require.Zero(t, got.queryMetrics)
}

func TestQueryContextInheritsTimeoutAndPreservesCallerDeadline(t *testing.T) {
	shared, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: to(3 * time.Second)}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	client := sharedTestClient(t, shared, OperationOptions{EndToEndTimeout: to(2 * time.Second)})
	for _, options := range []OperationOptions{{}, {EndToEndTimeout: to(time.Duration(0))}} {
		ctx, release, err := client.queryContext(context.Background(), options)
		require.NoError(t, err)
		deadline, bounded := ctx.Deadline()
		require.True(t, bounded)
		want := 2 * time.Second
		if options.EndToEndTimeout != nil {
			want = time.Second
		}
		require.InDelta(t, float64(want), float64(time.Until(deadline)), float64(200*time.Millisecond))
		release()
	}
	parent, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	ctx, release, err := client.queryContext(parent, OperationOptions{})
	require.NoError(t, err)
	defer release()
	want, _ := parent.Deadline()
	got, _ := ctx.Deadline()
	require.Equal(t, want, got)
}

func TestQueryCompletionCopiesPageAndPlannerToken(t *testing.T) {
	page, err := syntheticQueryCompletion([]byte(`{"Documents":[{"id":"first"},{"id":"second"}]}`), "planner-token", 200)
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte(`{"id":"first"}`), []byte(`{"id":"second"}`)}, page.Items)
	end, err := syntheticQueryCompletion(nil, "", 0)
	require.NoError(t, err)
	require.Zero(t, end)
}

func TestQueryPageSetupActivityID(t *testing.T) {
	for _, tt := range []struct {
		name     string
		result   completionResult
		setup    Response
		expected Response
	}{
		{
			name:     "exhausted first fetch",
			setup:    Response{RequestCharge: 3, ActivityID: "metadata"},
			expected: Response{RequestCharge: 3, ActivityID: "metadata"},
		},
		{
			name: "preserve query activity",
			result: completionResult{
				response: ItemResponse{Response: Response{RequestCharge: 2, ActivityID: "query"}},
				body:     []byte(`{"Documents":[]}`), httpStatus: 200,
			},
			setup:    Response{RequestCharge: 3, ActivityID: "metadata"},
			expected: Response{RequestCharge: 5, ActivityID: "query"},
		},
		{
			name:  "uncharged setup",
			setup: Response{ActivityID: "metadata"},
		},
		{
			name: "exhausted later fetch",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			page, err := tt.result.queryPage(tt.setup)
			require.NoError(t, err)
			require.Equal(t, tt.expected, page.Response)
			require.Empty(t, page.Items)
		})
	}
}

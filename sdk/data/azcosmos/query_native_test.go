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
	shared, err := NewRuntime(&RuntimeOptions{Operation: OperationOptions{EndToEndTimeout: time.Second}})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, shared.Close()) })
	credential := &blockingTokenCredential{started: make(chan struct{}), stopped: make(chan struct{})}
	client, err := NewClient("https://myaccount.documents.azure.com", credential, &ClientOptions{Runtime: shared})
	require.NoError(t, err)
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
	require.NoError(t, shared.SetOperationOptions(OperationOptions{EndToEndTimeout: 10 * time.Second}))
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Less(t, time.Since(start), 3*time.Second)
	case <-time.After(3 * time.Second):
		t.Fatal("runtime update replaced the admitted query budget")
	}
}

func TestQueryCompletionCopiesPageAndPlannerToken(t *testing.T) {
	page, err := syntheticQueryCompletion([]byte(`{"Documents":[{"id":"first"},{"id":"second"}]}`), "planner-token", 200)
	require.NoError(t, err)
	require.Equal(t, "planner-token", page.ContinuationToken)
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
			require.Empty(t, page.ContinuationToken)
		})
	}
}

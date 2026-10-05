// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryNativeOptions(t *testing.T) {
	yes, no := true, false
	req, err := newQueryRequest(NewQuery("SELECT * FROM c"), NewFeedScopeForFullContainer(), &QueryOptions{
		Operation:     OperationOptions{EndToEndTimeout: 5 * time.Second},
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

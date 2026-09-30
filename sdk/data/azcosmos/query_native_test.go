// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

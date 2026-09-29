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

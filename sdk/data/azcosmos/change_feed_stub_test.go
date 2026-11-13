// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeFeedUnavailableDriverClosesPager(t *testing.T) {
	pager := newTestContainer(t).NewChangeFeedPager(NewFeedScopeForFullContainer(), NewChangeFeedStartFromNow(), nil)
	require.True(t, pager.More())
	page, err := pager.NextPage(t.Context())
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientError, cosmosErr.Code)
	require.ErrorContains(t, err, "driver")
	require.False(t, pager.More())
	require.NoError(t, pager.Close())
	token, err := pager.ContinuationToken(t.Context())
	require.Empty(t, token)
	require.Error(t, err)
}

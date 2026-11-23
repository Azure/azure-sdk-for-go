// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFeedRangeInvalidValues(t *testing.T) {
	for _, serialized := range []string{"", "a\x00b"} {
		value, err := ParseFeedRange(serialized)
		require.Zero(t, value)
		require.Error(t, err)
	}
	for _, bounds := range [][2]string{{"\x00", "FF"}, {"", "F\x00F"}} {
		value, err := NewFeedRange(bounds[0], bounds[1])
		require.Zero(t, value)
		require.ErrorContains(t, err, "NUL")
	}
	var zero FeedRange
	ok, err := zero.IsSubsetOf(zero)
	require.False(t, ok)
	require.ErrorContains(t, err, "zero value")
	ok, err = zero.Overlaps(zero)
	require.False(t, ok)
	require.ErrorContains(t, err, "zero value")

	container := newTestContainer(t)
	require.NoError(t, container.database.client.Close())
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	scope := NewFeedScopeForFeedRange(zero)
	query := container.NewQueryItemsPager[map[string]any](NewQuery("SELECT * FROM c"), scope, nil)
	defer func() { require.NoError(t, query.Close()) }()
	page, err := query.NextPage(ctx)
	require.Zero(t, page)
	require.ErrorContains(t, err, "feed range")
	feed := container.NewChangeFeedPager(scope, NewChangeFeedStartFromBeginning(), nil)
	defer func() { require.NoError(t, feed.Close()) }()
	changes, err := feed.NextPage(ctx)
	require.Zero(t, changes)
	require.ErrorContains(t, err, "feed range")
}

func TestFeedRangeDiscoveryLifetimeOrdering(t *testing.T) {
	container := newTestContainer(t)
	page, err := container.GetFeedRanges(nil, nil) //nolint:staticcheck // verifies the boundary guard
	require.Zero(t, page)
	require.ErrorContains(t, err, "context must not be nil")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	page, err = container.GetFeedRanges(ctx, nil)
	require.Zero(t, page)
	require.ErrorIs(t, err, context.Canceled)
	require.NoError(t, container.database.client.Close())
	page, err = container.GetFeedRanges(ctx, nil)
	require.Zero(t, page)
	var cosmosErr *Error
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)

	response, err := container.GetFeedRangeForPartitionKey(ctx, PartitionKey{})
	require.Zero(t, response)
	require.ErrorContains(t, err, "partition key")
	response, err = container.GetFeedRangeForPartitionKey(ctx, NewPartitionKeyString("synthetic"))
	require.Zero(t, response)
	require.ErrorAs(t, err, &cosmosErr)
	require.Equal(t, CodeClientClosed, cosmosErr.Code)
}

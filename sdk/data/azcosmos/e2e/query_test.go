// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"fmt"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

// TestQueryParameterizedFilterAndOrder ports Rust scenario query.parameterized-filter: bound
// @-parameters both filter and order results, and reversing ORDER BY reverses the result order,
// proving the parameters are not merely accepted but actually applied server-side.
func TestQueryParameterizedFilterAndOrder(t *testing.T) {
	fx := newFixture(t, "query.parameterized-filter", "")
	ctx := context.Background()
	pk := uniqueID(t)
	partitionKey := azcosmos.NewPartitionKeyString(pk)

	for id, score := range map[string]int64{"a": 3, "b": 1, "c": 2} {
		_, err := fx.Container.CreateItem(ctx, partitionKey, pk+"-"+id, mustMarshalItem(t, itemWithScore(pk+"-"+id, pk, score)), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, partitionKey, pk+"-"+id)
	}

	ascending, err := azcosmos.NewQuery("SELECT * FROM c WHERE c.pk = @pk AND c.score >= @min ORDER BY c.score ASC").WithParameter("@pk", pk)
	require.NoError(t, err)
	ascending, err = ascending.WithParameter("@min", 2)
	require.NoError(t, err)
	require.Equal(t, []string{"c", "a"}, queryOrderedSuffixes(t, fx, ascending, partitionKey, pk))

	descending, err := azcosmos.NewQuery("SELECT * FROM c WHERE c.pk = @pk AND c.score >= @min ORDER BY c.score DESC").WithParameter("@pk", pk)
	require.NoError(t, err)
	descending, err = descending.WithParameter("@min", 2)
	require.NoError(t, err)
	require.Equal(t, []string{"a", "c"}, queryOrderedSuffixes(t, fx, descending, partitionKey, pk))
}

// queryOrderedSuffixes drains query into a slice of ids with the pk+"-" prefix stripped, so
// assertions read as the bare case labels ("a", "b", "c") the scenario uses.
func queryOrderedSuffixes(t *testing.T, fx fixture, query azcosmos.Query, scope azcosmos.PartitionKey, prefix string) []string {
	t.Helper()
	pager := fx.Container.NewQueryItemsPager(query, azcosmos.NewFeedScopeForPartitionKey(scope), nil)
	defer func() { require.NoError(t, pager.Close()) }()
	var ids []string
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, raw := range page.Items {
			value := mustUnmarshalItem(t, raw)
			ids = append(ids, value.ID[len(prefix)+1:])
		}
	}
	return ids
}

// TestQueryInvalidSyntaxIsNotAnEmptyFeed ports Rust scenario query.invalid-syntax: malformed
// query text is reported as a typed CodeBadRequest failure, whether it fails while opening the
// pager or on its first page, and must never look like a successful empty feed.
func TestQueryInvalidSyntaxIsNotAnEmptyFeed(t *testing.T) {
	fx := newFixture(t, "query.invalid-syntax", "")
	scope := azcosmos.NewFeedScopeForPartitionKey(azcosmos.NewPartitionKeyString("A"))

	pager := fx.Container.NewQueryItemsPager(azcosmos.NewQuery("SELECT FROM"), scope, nil)
	defer func() { require.NoError(t, pager.Close()) }()
	_, err := pager.NextPage(context.Background())
	_ = requireCode(t, err, azcosmos.CodeBadRequest)
}

// TestQueryPaginationResumesWithoutLossOrDuplication ports Rust scenario
// query.pagination-resume: a continuation snapshot taken after the first page lets a fresh pager
// resume and, combined with the items already delivered, deliver every item exactly once with no
// loss or duplication.
func TestQueryPaginationResumesWithoutLossOrDuplication(t *testing.T) {
	fx := newFixture(t, "query.pagination-resume", "")
	ctx := context.Background()
	const itemCount = 24
	run := uniqueID(t)
	partitionKeys := make([]string, 6)
	for i := range partitionKeys {
		partitionKeys[i] = fmt.Sprintf("%s-partition-%d", run, i)
	}

	for i := 0; i < itemCount; i++ {
		id := fmt.Sprintf("%s-item-%d", run, i)
		pkValue := partitionKeys[i%len(partitionKeys)]
		pk := azcosmos.NewPartitionKeyString(pkValue)
		_, err := fx.Container.CreateItem(ctx, pk, id, mustMarshalItem(t, item(id, pkValue, int64(i))), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, pk, id)
	}

	query := azcosmos.NewQuery("SELECT * FROM c WHERE STARTSWITH(c.pk, @run) ORDER BY c.value ASC")
	query, err := query.WithParameter("@run", run)
	require.NoError(t, err)
	options := &azcosmos.QueryOptions{Feed: azcosmos.FeedOptions{PageSizeHint: 2}}

	first := fx.Container.NewQueryItemsPager(query, azcosmos.NewFeedScopeForFullContainer(), options)
	firstPage, err := first.NextPage(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, firstPage.Items)
	require.Less(t, len(firstPage.Items), itemCount)
	token, err := first.ContinuationToken(ctx)
	require.NoError(t, err)
	require.NoError(t, first.Close())

	var values []int64
	for _, raw := range firstPage.Items {
		values = append(values, mustUnmarshalItem(t, raw).Value)
	}

	resumed := fx.Container.NewQueryItemsPager(query, azcosmos.NewFeedScopeForFullContainer(), &azcosmos.QueryOptions{
		Feed: azcosmos.FeedOptions{PageSizeHint: 2, ContinuationToken: token},
	})
	defer func() { require.NoError(t, resumed.Close()) }()
	for resumed.More() {
		page, err := resumed.NextPage(ctx)
		require.NoError(t, err)
		for _, raw := range page.Items {
			values = append(values, mustUnmarshalItem(t, raw).Value)
		}
	}

	require.Len(t, values, itemCount, "resumed query must deliver every item exactly once")
	seen := make(map[int64]bool, itemCount)
	for _, v := range values {
		require.Falsef(t, seen[v], "value %d delivered more than once", v)
		seen[v] = true
	}
}

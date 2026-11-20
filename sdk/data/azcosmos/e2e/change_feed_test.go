// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"encoding/json"
	"fmt"
	"testing"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

func TestChangeFeedPaginationResumesWithoutLossOrDuplication(t *testing.T) {
	const scenarioID = "changefeed.pagination-resume"
	fx := newFixture(t, scenarioID, "")
	run := uniqueID(t)
	partitionKey := run + "-partition"
	pk := azcosmos.NewPartitionKeyString(partitionKey)
	const itemCount = 6
	wanted := make(map[string]bool, itemCount)
	for i := range itemCount {
		id := fmt.Sprintf("%s-item-%d", run, i)
		_, err := fx.Container.CreateItem(t.Context(), pk, id, mustMarshalItem(t, item(id, partitionKey, int64(i))), nil)
		require.NoError(t, err)
		trackItem(t, fx.Container, pk, id)
		wanted[id] = false
	}

	scope := azcosmos.NewFeedScopeForPartitionKey(pk)
	start := azcosmos.NewChangeFeedStartFromBeginning()
	first := fx.Container.NewChangeFeedPager(scope, start, &azcosmos.ChangeFeedOptions{
		Mode: azcosmos.ChangeFeedModeLatestVersion,
		Feed: azcosmos.FeedOptions{PageSizeHint: 1},
	})
	t.Cleanup(func() { require.NoError(t, first.Close()) })

	var firstPage azcosmos.ChangeFeedResponse
	for range 8 {
		var err error
		firstPage, err = first.NextPage(t.Context())
		require.NoError(t, err)
		if len(firstPage.Items) > 0 {
			break
		}
	}
	require.NotEmpty(t, firstPage.Items, "the initial snapshot must contain the seeded items")
	require.Less(t, len(firstPage.Items), itemCount, "the checkpoint must interrupt a multi-page snapshot")
	accept := func(items [][]byte) {
		for _, raw := range items {
			var image struct {
				ID      string `json:"id"`
				Current *struct {
					ID string `json:"id"`
				} `json:"current"`
			}
			require.NoError(t, json.Unmarshal(raw, &image))
			if image.Current != nil {
				image.ID = image.Current.ID
			}
			seen, expected := wanted[image.ID]
			require.Truef(t, expected, "unexpected change-feed item %q", image.ID)
			require.Falsef(t, seen, "change-feed item %q was replayed", image.ID)
			wanted[image.ID] = true
		}
	}
	accept(firstPage.Items)
	token, err := first.ContinuationToken(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NoError(t, first.Close())

	resumeFixture := newFixture(t, scenarioID, "")
	resumed := resumeFixture.Container.NewChangeFeedPager(scope, start, &azcosmos.ChangeFeedOptions{
		Mode: azcosmos.ChangeFeedModeLatestVersion,
		Feed: azcosmos.FeedOptions{PageSizeHint: 1, ContinuationToken: token},
	})
	t.Cleanup(func() { require.NoError(t, resumed.Close()) })

	allSeen := func() bool {
		for _, seen := range wanted {
			if !seen {
				return false
			}
		}
		return true
	}
	for range itemCount + 8 {
		if allSeen() {
			break
		}
		page, err := resumed.NextPage(t.Context())
		require.NoError(t, err)
		accept(page.Items)
	}
	require.True(t, allSeen(), "resume must deliver every remaining item from the stable snapshot")
	page, err := resumed.NextPage(t.Context())
	require.NoError(t, err)
	require.Empty(t, page.Items, "the completed snapshot must not replay changes after resume")
}

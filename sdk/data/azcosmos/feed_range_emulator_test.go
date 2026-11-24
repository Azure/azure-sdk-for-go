// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build azcosmos_native_preview && cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEmulatorFeedRangeDisjointReads(t *testing.T) {
	for _, hierarchical := range []bool{false, true} {
		t.Run(fmt.Sprintf("hierarchical-%t", hierarchical), func(t *testing.T) {
			client, databaseID, containerID := emulatorClient(t)
			if hierarchical {
				containerID = "query-hierarchical"
			}
			container, err := client.NewContainer(t.Context(), databaseID, containerID)
			require.NoError(t, err)
			discovered, err := container.GetFeedRanges(t.Context(), &GetFeedRangesOptions{ForceRefresh: true})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(discovered.FeedRanges), 2, "provision a multi-range emulator container for this contract")
			prefix := fmt.Sprintf("feed-range-%d-", time.Now().UnixNano())
			wanted := make(map[string]bool)
			for i := range 16 {
				id := fmt.Sprintf("%s%d", prefix, i)
				key := NewPartitionKeyString(id)
				if hierarchical {
					key = key.AppendString("child")
				}
				seedChangeFeedEmulatorItem(t, container, key, id, id, "child")
				wanted[id] = true
			}
			query, err := NewQuery("SELECT * FROM c WHERE STARTSWITH(c.id, @prefix)").WithParameter("@prefix", prefix)
			require.NoError(t, err)
			queries, changes := make(map[string]bool), make(map[string]bool)
			for _, value := range discovered.FeedRanges {
				parsed, err := ParseFeedRange(value.String())
				require.NoError(t, err)
				scope := NewFeedScopeForFeedRange(parsed)
				pager := container.NewQueryItemsPager[json.RawMessage](query, scope, &QueryOptions{Feed: FeedOptions{PageSizeHint: 2}})
				for pager.More() {
					page, err := pager.NextPage(t.Context())
					require.NoError(t, err)
					for _, body := range page.Items {
						id := changeFeedEmulatorItemID(t, body)
						require.True(t, wanted[id])
						require.False(t, queries[id], "disjoint range queries must not duplicate data")
						queries[id] = true
					}
				}
				require.NoError(t, pager.Close())
				feed := container.NewChangeFeedPager(scope, NewChangeFeedStartFromBeginning(), &ChangeFeedOptions{Feed: FeedOptions{PageSizeHint: 1000}})
				for range 32 {
					page, err := feed.NextPage(t.Context())
					require.NoError(t, err)
					for _, body := range page.Items {
						id := changeFeedEmulatorItemID(t, body)
						if wanted[id] {
							require.False(t, changes[id], "disjoint range change feeds must not duplicate data")
							changes[id] = true
						}
					}
					if page.StatusCode == 304 {
						break
					}
				}
				token, err := feed.ContinuationToken(t.Context())
				require.NoError(t, err)
				require.NoError(t, feed.Close())
				resumed := container.NewChangeFeedPager(scope, NewChangeFeedStartFromNow(), &ChangeFeedOptions{Feed: FeedOptions{ContinuationToken: token}})
				page, err := resumed.NextPage(t.Context())
				require.NoError(t, err)
				for _, body := range page.Items {
					require.False(t, wanted[changeFeedEmulatorItemID(t, body)], "resume must not replay processed seed items")
				}
				require.NoError(t, resumed.Close())
			}
			require.Equal(t, wanted, queries, "range queries must have no gaps")
			require.Equal(t, wanted, changes, "range change feeds must have no gaps")
		})
	}
}

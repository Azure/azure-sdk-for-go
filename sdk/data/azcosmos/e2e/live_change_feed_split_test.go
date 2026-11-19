// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/stretchr/testify/require"
)

func TestAzureLiveChangeFeedSplit(t *testing.T) {
	if os.Getenv("AZCOSMOS_LIVE_SPLIT") != "1" {
		t.Skip("set AZCOSMOS_LIVE_SPLIT=1 for the manual physical-split test")
	}
	require.Equal(t, "azureLive", os.Getenv("AZURE_COSMOS_E2E_BACKEND"))
	require.Empty(t, os.Getenv("EMULATOR"))
	endpoint, databaseID, containerID := os.Getenv("AZCOSMOS_ENDPOINT"),
		os.Getenv("AZCOSMOS_DATABASE"), os.Getenv("AZCOSMOS_SPLIT_CONTAINER")
	target, err := url.Parse(endpoint)
	liveNoError(t, "parse split account endpoint", err)
	require.True(t, target.Scheme == "https" && strings.HasSuffix(target.Hostname(), ".documents.azure.com") &&
		target.User == nil && target.Port() == "" && target.RawQuery == "" && target.Fragment == "" &&
		(target.Path == "" || target.Path == "/"), "use a public Azure Cosmos DB account endpoint")
	require.NotEmpty(t, databaseID)
	require.NotEmpty(t, containerID, "use a fresh disposable split container, not the smoke container")
	require.NotEqual(t, os.Getenv("AZCOSMOS_CONTAINER"), containerID)
	require.NotEqual(t, os.Getenv("AZCOSMOS_HPK_CONTAINER"), containerID)

	f := newLiveFixture(t, endpoint, databaseID, containerID)
	ranges := liveReadRanges(t, f, f.ctx)
	require.Len(t, ranges, 1, "start with an unsplit container")
	parent := ranges[0]
	require.NotEmpty(t, parent.ID)
	require.Empty(t, parent.Min)
	require.Equal(t, "FF", parent.Max)

	scope, start := azcosmos.NewFeedScopeForFullContainer(), azcosmos.NewChangeFeedStartFromBeginning()
	baseline := f.prefix + "-before"
	f.create(t, "-before", 0)
	initial := liveChangeFeedPager(t, f.container, scope, start, "")
	liveCollectChangeFeed(t, f.ctx, initial, newLiveChangeFeedProgress([]string{baseline}))
	token := liveChangeFeedCheckpoint(t, f, initial)

	t.Log("checkpoint saved; now raise the dedicated container's manual throughput to 13000 RU/s")
	children := liveWaitForSplit(t, f, parent)
	t.Logf("verified physical split: one parent replaced by %d child ranges", len(children))

	// Give post-split operations a fresh budget, independent of topology polling.
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	f.ctx = ctx
	idle := liveChangeFeedPager(t, liveFreshContainer(t, f), scope, start, token)
	liveRequireSplitIdle(t, ctx, idle, len(children))
	liveNoError(t, "close idle post-split pager", idle.Close())

	writer := liveFreshContainer(t, f)
	var expected []string
	for i := range 24 {
		seed := *f
		seed.prefix = fmt.Sprintf("%s-%02d", f.prefix, i)
		seed.pk = azcosmos.NewPartitionKeyString(seed.prefix)
		seed.container = writer
		seed.create(t, "-after", int64(i))
		expected = append(expected, seed.prefix+"-after")
	}
	resumed := liveChangeFeedPager(t, liveFreshContainer(t, f), scope, start, token)
	liveCollectChangeFeed(t, ctx, resumed, newLiveChangeFeedProgress(expected))
	liveRequireSplitIdle(t, ctx, resumed, len(children))
	next := liveChangeFeedCheckpoint(t, f, resumed)
	noReplay := liveChangeFeedPager(t, liveFreshContainer(t, f), scope, start, next)
	liveRequireSplitIdle(t, ctx, noReplay, len(children))
	liveRequireChildWrites(t, f, children, baseline, expected)
}

func liveWaitForSplit(t *testing.T, f *liveFixture, parent livePhysicalRange) []livePhysicalRange {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Minute)
	defer cancel()
	for {
		ranges := liveReadRanges(t, f, ctx)
		if isLiveSplit(parent, ranges) {
			return ranges
		}
		select {
		case <-ctx.Done():
			t.Fatal("Azure did not report complete child ranges within 15 minutes; recovery was not exercised")
		case <-time.After(15 * time.Second):
		}
	}
}

func isLiveSplit(parent livePhysicalRange, ranges []livePhysicalRange) bool {
	if len(ranges) < 2 {
		return false
	}
	slices.SortFunc(ranges, func(a, b livePhysicalRange) int { return strings.Compare(a.Min, b.Min) })
	ids := map[string]bool{parent.ID: true, "": true}
	next := parent.Min
	for _, child := range ranges {
		if ids[child.ID] || child.Min != next || child.Max <= child.Min {
			return false
		}
		ids[child.ID] = true
		next = child.Max
	}
	return next == parent.Max
}

func liveRequireSplitIdle(t *testing.T, ctx context.Context, pager *azcosmos.ChangeFeedPager, ranges int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for range ranges * 2 {
		page, err := pager.NextPage(ctx)
		liveNoError(t, "poll idle split feed", err)
		liveCheckChangeFeedPage(t, pager, page)
		require.Equal(t, http.StatusNotModified, page.StatusCode, "idle split recovery must return 304 without replay")
	}
}

// Prove that the resume checks more than one child, independently of the SDK
// pager under test. All writes are tiny and fit in one bounded page per range.
func liveRequireChildWrites(t *testing.T, f *liveFixture, children []livePhysicalRange, baseline string, expected []string) {
	t.Helper()
	progress := newLiveChangeFeedProgress(expected)
	for _, child := range children {
		before := progress.remaining
		for _, raw := range liveReadChildFeed(t, f, child.ID) {
			var envelope struct {
				ID      string `json:"id"`
				Current struct {
					ID string `json:"id"`
				} `json:"current"`
			}
			liveNoError(t, "decode child change-feed probe entry", json.Unmarshal(raw, &envelope))
			id := envelope.Current.ID
			if id == "" {
				id = envelope.ID
			}
			if id != baseline {
				liveNoError(t, "verify child change-feed IDs", progress.acceptID(id))
			}
		}
		require.Less(t, progress.remaining, before, "each child must contain at least one post-split write")
	}
	require.Zero(t, progress.remaining, "all expected writes must be visible across the child ranges")
}

func liveReadChildFeed(t *testing.T, f *liveFixture, rangeID string) []json.RawMessage {
	t.Helper()
	headers := http.Header{}
	headers.Set("A-IM", "Incremental feed")
	headers.Set("x-ms-documentdb-partitionkeyrangeid", rangeID)
	headers.Set("x-ms-max-item-count", "100")
	response := liveReadContainerFeed(t, f, f.ctx, "docs", headers)
	defer response.Body.Close()
	var page struct {
		Documents []json.RawMessage `json:"Documents"`
	}
	liveNoError(t, "decode child change-feed probe", json.NewDecoder(io.LimitReader(response.Body, 64*1024)).Decode(&page))
	return page.Documents
}

func TestLiveSplitTopology(t *testing.T) {
	parent := livePhysicalRange{ID: "0", Max: "FF"}
	for _, tc := range []struct {
		name   string
		ranges []livePhysicalRange
		valid  bool
	}{
		{"unsplit", []livePhysicalRange{parent}, false},
		{"disordered children", []livePhysicalRange{{"2", "80", "FF"}, {"1", "", "80"}}, true},
		{"stale parent", []livePhysicalRange{parent, {"1", "", "80"}, {"2", "80", "FF"}}, false},
		{"gap", []livePhysicalRange{{"1", "", "70"}, {"2", "80", "FF"}}, false},
		{"overlap", []livePhysicalRange{{"1", "", "90"}, {"2", "80", "FF"}}, false},
		{"duplicate ID", []livePhysicalRange{{"1", "", "80"}, {"1", "80", "FF"}}, false},
		{"partial coverage", []livePhysicalRange{{"1", "", "80"}, {"2", "80", "F0"}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.valid, isLiveSplit(parent, tc.ranges))
		})
	}
}

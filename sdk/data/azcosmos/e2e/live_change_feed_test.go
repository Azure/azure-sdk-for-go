// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

// These local cases refine the pinned catalog's pagination-resume scenario; they
// are not new shared scenario IDs. Every case must run for the scenario to pass.
func liveChangeFeed(t *testing.T, f *liveFixture) {
	cases := []struct {
		name string
		run  func(*testing.T)
	}{
		{"idle-checkpoint-fresh-resume", func(t *testing.T) {
			liveChangeFeedIdleResume(t, f)
		}},
		{"point-in-time-server-timestamps", func(t *testing.T) {
			liveChangeFeedPointInTime(t, newLiveFixture(t, f.endpoint, f.databaseID, f.containerID))
		}},
		{"hpk-single-range-prefix-partial-resume", func(t *testing.T) {
			id := os.Getenv("AZCOSMOS_HPK_CONTAINER")
			require.NotEmpty(t, id, "the live resource deployment must provide AZCOSMOS_HPK_CONTAINER")
			liveChangeFeedPrefix(t, newLiveFixtureWithPaths(t, f.endpoint, f.databaseID, id,
				[]string{"/tenant", "/user"}))
		}},
	}
	executed := 0
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			executed++
			test.run(t)
		})
	}
	require.Equal(t, len(cases), executed, "all local change-feed cases are required; do not filter out subtests")
}

func liveChangeFeedPager(t *testing.T, container *azcosmos.ContainerClient, scope azcosmos.FeedScope,
	start azcosmos.ChangeFeedStartFrom, token string) *azcosmos.ChangeFeedPager {
	t.Helper()
	pager := container.NewChangeFeedPager(scope, start, &azcosmos.ChangeFeedOptions{
		Mode: azcosmos.ChangeFeedModeLatestVersion,
		Feed: azcosmos.FeedOptions{PageSizeHint: 1, ContinuationToken: token},
	})
	t.Cleanup(func() { liveNoError(t, "close change-feed pager", pager.Close()) })
	return pager
}

func liveFreshContainer(t *testing.T, f *liveFixture) *azcosmos.ContainerClient {
	t.Helper()
	client := newLiveClient(t, f.ctx, f.endpoint)
	container, err := client.NewContainer(f.ctx, f.databaseID, f.containerID)
	liveNoError(t, "open fresh-client container", err)
	return container
}

func liveChangeFeedCheckpoint(t *testing.T, f *liveFixture, pager *azcosmos.ChangeFeedPager) string {
	t.Helper()
	token, err := pager.ContinuationToken(f.ctx)
	liveNoError(t, "checkpoint change feed", err)
	require.True(t, token != "", "checkpoint must not be empty")
	liveNoError(t, "close change-feed pager before fresh-client resume", pager.Close())
	return token
}

func liveChangeFeedIdleResume(t *testing.T, f *liveFixture) {
	// A unique complete partition key makes this idle assertion local, not a claim
	// that the first 304 from a full-container fan-out means every range is idle.
	scope := azcosmos.NewFeedScopeForPartitionKey(f.pk)
	start := azcosmos.NewChangeFeedStartFromBeginning()
	f.create(t, "-old", 0)
	initial := liveChangeFeedPager(t, f.container, scope, start, "")
	liveCollectChangeFeed(t, f.ctx, initial, newLiveChangeFeedProgress([]string{f.prefix + "-old"}))
	token := liveChangeFeedCheckpoint(t, f, initial)

	ids := []string{f.prefix + "-a", f.prefix + "-b", f.prefix + "-c"}
	for i, suffix := range []string{"-a", "-b", "-c"} {
		f.create(t, suffix, int64(i))
	}
	resumed := liveChangeFeedPager(t, liveFreshContainer(t, f), scope, start, token)
	liveCollectChangeFeed(t, f.ctx, resumed, newLiveChangeFeedProgress(ids))
	next := liveChangeFeedCheckpoint(t, f, resumed)

	noReplay := liveChangeFeedPager(t, liveFreshContainer(t, f), scope, start, next)
	for range 3 {
		liveCollectChangeFeed(t, f.ctx, noReplay, newLiveChangeFeedProgress(nil))
	}
}

func liveChangeFeedPointInTime(t *testing.T, f *liveFixture) {
	f.create(t, "-before", 1)
	before := liveReadServerTimestamp(t, f, f.prefix+"-before")
	startTime := time.Unix(before+1, 0).UTC()
	// The service has second precision. Separate the images using elapsed time,
	// then verify both sides against server _ts rather than trusting local UTC.
	liveChangeFeedWait(t, f.ctx, 3*time.Second)
	f.create(t, "-after", 2)
	after := liveReadServerTimestamp(t, f, f.prefix+"-after")
	require.Less(t, before, startTime.Unix())
	require.Greater(t, after, startTime.Unix(), "the after image must be strictly beyond the service-second boundary")

	scope := azcosmos.NewFeedScopeForPartitionKey(f.pk)
	start := azcosmos.NewChangeFeedStartFromPointInTime(startTime)
	pager := liveChangeFeedPager(t, f.container, scope, start, "")
	liveCollectChangeFeed(t, f.ctx, pager, newLiveChangeFeedProgress([]string{f.prefix + "-after"}))
}

func liveReadServerTimestamp(t *testing.T, f *liveFixture, id string) int64 {
	t.Helper()
	response, err := f.container.ReadItem(f.ctx, f.pk, id, nil)
	liveNoError(t, "read seeded server timestamp", err)
	var document struct {
		ID string `json:"id"`
		TS int64  `json:"_ts"`
	}
	liveNoError(t, "decode seeded server timestamp", json.Unmarshal(response.Value, &document))
	require.Equal(t, id, document.ID)
	require.Positive(t, document.TS)
	return document.TS
}

func liveChangeFeedPrefix(t *testing.T, f *liveFixture) {
	liveRequireSingleRange(t, f)
	ids := []string{f.prefix + "-a1", f.prefix + "-b1", f.prefix + "-a2", f.prefix + "-b2"}
	// The foreign tenant shares the verified physical range. Prefix isolation
	// must therefore filter actual out-of-scope changes, not merely route elsewhere.
	liveSeedHierarchicalItem(t, f, f.prefix+"-foreign", f.prefix+"-other", "a")
	for i, id := range ids {
		child := "a"
		if i%2 != 0 {
			child = "b"
		}
		liveSeedHierarchicalItem(t, f, id, f.prefix, child)
	}
	scope := azcosmos.NewFeedScopeForPartitionKey(f.pk)
	start := azcosmos.NewChangeFeedStartFromBeginning()
	pager := liveChangeFeedPager(t, f.container, scope, start, "")
	progress := newLiveChangeFeedProgress(ids)
	for range 64 {
		page, err := pager.NextPage(f.ctx)
		liveNoError(t, "fetch first hierarchical prefix page", err)
		liveCheckChangeFeedPage(t, pager, page)
		liveNoError(t, "validate prefix page IDs", progress.accept(page.Items))
		if len(page.Items) != 0 {
			break
		}
		liveChangeFeedWait(t, f.ctx, 500*time.Millisecond)
	}
	require.Less(t, progress.remaining, len(ids), "a nonempty first prefix page is required")
	require.Positive(t, progress.remaining, "checkpoint must interrupt delivery, not merely follow the entire snapshot")
	token := liveChangeFeedCheckpoint(t, f, pager)
	resumed := liveChangeFeedPager(t, liveFreshContainer(t, f), scope, start, token)
	liveCollectChangeFeed(t, f.ctx, resumed, progress)
	liveRequireSingleRange(t, f)
}

func liveSeedHierarchicalItem(t *testing.T, f *liveFixture, id, tenant, user string) {
	t.Helper()
	key := azcosmos.NewPartitionKeyString(tenant).AppendString(user)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_, err := f.container.DeleteItem(ctx, key, id, nil)
		var serviceErr *azcosmos.Error
		if errors.As(err, &serviceErr) && serviceErr.Code == azcosmos.CodeNotFound {
			return
		}
		liveNoError(t, "delete owned hierarchical item", err)
	})
	body, err := json.Marshal(map[string]string{"id": id, "tenant": tenant, "user": user})
	liveNoError(t, "encode hierarchical item", err)
	_, err = f.container.CreateItem(f.ctx, key, id, json.RawMessage(body), nil)
	liveNoError(t, "create hierarchical item", err)
}

// No public SDK range-enumeration API exists. This bounded, read-only metadata
// request verifies topology without parsing opaque checkpoints or assuming RU implies ranges.
func liveReadRanges(t *testing.T, f *liveFixture, ctx context.Context) []livePhysicalRange {
	t.Helper()
	response := liveReadContainerFeed(t, f, ctx, "pkranges", nil)
	defer response.Body.Close()
	require.Empty(t, response.Header.Get("x-ms-continuation"), "all routing ranges must fit in this bounded response")
	var topology struct {
		Ranges []livePhysicalRange `json:"PartitionKeyRanges"`
	}
	liveNoError(t, "decode physical ranges", json.NewDecoder(io.LimitReader(response.Body, 64*1024)).Decode(&topology))
	return topology.Ranges
}

type livePhysicalRange struct {
	ID  string `json:"id"`
	Min string `json:"minInclusive"`
	Max string `json:"maxExclusive"`
}

func liveRequireSingleRange(t *testing.T, f *liveFixture) {
	t.Helper()
	ranges := liveReadRanges(t, f, f.ctx)
	require.Len(t, ranges, 1, "fixture must have one physical range")
	require.Empty(t, ranges[0].Min)
	require.Equal(t, "FF", ranges[0].Max)
}

func liveReadContainerFeed(t *testing.T, f *liveFixture, ctx context.Context, resource string, headers http.Header) *http.Response {
	t.Helper()
	cred, err := credential.New(nil)
	liveNoError(t, "create feed-probe credential", err)
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{"https://cosmos.azure.com/.default"}})
	liveNoError(t, "acquire feed-probe token", err)
	path := strings.TrimRight(f.endpoint, "/") + "/dbs/" + url.PathEscape(f.databaseID) +
		"/colls/" + url.PathEscape(f.containerID) + "/" + resource
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	liveNoError(t, "construct feed probe", err)
	if headers != nil {
		request.Header = headers.Clone()
	}
	request.Header.Set("Authorization", url.QueryEscape("type=aad&ver=1.0&sig="+token.Token))
	request.Header.Set("x-ms-version", "2020-07-15")
	request.Header.Set("x-ms-date", time.Now().UTC().Format(http.TimeFormat))
	client := &http.Client{Timeout: 20 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	response, err := client.Do(request)
	liveNoError(t, "read Cosmos feed probe", err)
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		t.Fatalf("read Cosmos %s feed: HTTP %d", resource, response.StatusCode)
	}
	return response
}

type liveChangeFeedProgress struct {
	seen      map[string]bool
	remaining int
}

func newLiveChangeFeedProgress(ids []string) *liveChangeFeedProgress {
	p := &liveChangeFeedProgress{seen: make(map[string]bool), remaining: len(ids)}
	for _, id := range ids {
		p.seen[id] = false
	}
	return p
}

func (p *liveChangeFeedProgress) accept(items [][]byte) error {
	for _, raw := range items {
		var envelope struct {
			Current *struct {
				ID string `json:"id"`
			} `json:"current"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			return errors.New("invalid change-feed envelope")
		}
		if envelope.Current == nil || envelope.Current.ID == "" {
			return errors.New("change-feed envelope has no current ID")
		}
		if err := p.acceptID(envelope.Current.ID); err != nil {
			return err
		}
	}
	return nil
}

func (p *liveChangeFeedProgress) acceptID(id string) error {
	seen, expected := p.seen[id]
	if !expected || seen {
		return errors.New("unexpected or replayed change-feed ID")
	}
	p.seen[id] = true
	p.remaining--
	return nil
}

func liveCheckChangeFeedPage(t *testing.T, pager *azcosmos.ChangeFeedPager, page azcosmos.ChangeFeedResponse) {
	t.Helper()
	require.Contains(t, []int{http.StatusOK, http.StatusNotModified}, page.StatusCode)
	require.True(t, pager.More(), "idle and nonempty change-feed pages must remain pollable")
	if page.StatusCode == http.StatusNotModified {
		require.Empty(t, page.Items)
	}
}

func liveCollectChangeFeed(t *testing.T, ctx context.Context, pager *azcosmos.ChangeFeedPager, progress *liveChangeFeedProgress) {
	t.Helper()
	for range 64 {
		page, err := pager.NextPage(ctx)
		liveNoError(t, "poll bounded change feed", err)
		liveCheckChangeFeedPage(t, pager, page)
		liveNoError(t, "validate exact seeded change-feed IDs", progress.accept(page.Items))
		if page.StatusCode == http.StatusNotModified && progress.remaining == 0 {
			return
		}
		if len(page.Items) == 0 {
			liveChangeFeedWait(t, ctx, 500*time.Millisecond)
		}
	}
	t.Fatal("bounded polling did not deliver the exact seeded IDs followed by an idle page")
}

func liveChangeFeedWait(t *testing.T, ctx context.Context, delay time.Duration) {
	t.Helper()
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		liveNoError(t, "wait for change-feed visibility", ctx.Err())
	case <-timer.C:
	}
}

func TestLiveChangeFeedProgress(t *testing.T) {
	progress := newLiveChangeFeedProgress([]string{"first", "second"})
	require.NoError(t, progress.accept(nil), "idle pages must not consume expected IDs")
	require.NoError(t, progress.accept([][]byte{[]byte(`{"current":{"id":"first"}}`)}))
	require.Equal(t, 1, progress.remaining)
	require.NoError(t, progress.accept([][]byte{[]byte(`{"current":{"id":"second"}}`)}))
	require.Zero(t, progress.remaining)
	require.ErrorContains(t, progress.accept([][]byte{[]byte(`{"current":{"id":"first"}}`)}), "replayed")
	for _, raw := range []string{
		`{"current":{"id":"foreign"}}`, `{"current":{"id":"first"}}`,
		`{"id":"second"}`, `{"current":null}`, `{"current":{"id":""}}`, `not-json`,
	} {
		t.Run(raw, func(t *testing.T) {
			require.Error(t, progress.accept([][]byte{[]byte(raw)}), "foreign, replayed and malformed images must fail")
		})
	}
	require.NoError(t, newLiveChangeFeedProgress(nil).accept(nil))
	require.Error(t, newLiveChangeFeedProgress(nil).accept([][]byte{[]byte(`{"current":{"id":"replay"}}`)}))
}

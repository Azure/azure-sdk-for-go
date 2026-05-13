// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

// returnGoneOnChangeFeedPolicy is the change-feed-side analogue of
// returnGoneOnQueryPolicy in cosmos_container_read_many_test.go: it returns 410/Gone
// with a configurable PK-range substatus on change-feed requests (identified via
// the A-IM = "Incremental feed" header) until maxGone such responses have been
// emitted, after which subsequent calls pass through.
type returnGoneOnChangeFeedPolicy struct {
	maxGone   int32
	substatus string
	count     atomic.Int32
}

func (p *returnGoneOnChangeFeedPolicy) Do(req *policy.Request) (*http.Response, error) {
	// Match only true change-feed reads (against /docs), not the PK-range cache's
	// incremental refresh (against /pkranges) which also carries the A-IM header.
	if req.Raw().Header.Get(cosmosHeaderChangeFeed) != cosmosHeaderValuesChangeFeed ||
		!strings.HasSuffix(req.Raw().URL.Path, "/docs") {
		return req.Next()
	}
	n := p.count.Add(1)
	if n <= p.maxGone {
		headers := http.Header{}
		headers.Set(cosmosHeaderSubstatus, p.substatus)
		return &http.Response{
			StatusCode: http.StatusGone,
			Status:     "410 Gone",
			Header:     headers,
			Body:       io.NopCloser(strings.NewReader(`{"message":"Gone"}`)),
			Request:    req.Raw(),
		}, nil
	}
	return req.Next()
}

// createChangeFeedTestClient mirrors createReadManyTestClient but pre-seeds the PK
// range cache with the supplied physical ranges. Container cache is pre-populated
// with ResourceID=testRID so cross-container validation can be exercised.
func createChangeFeedTestClient(t *testing.T, srv *mock.Server, policies []policy.Policy, ranges []partitionKeyRange) *Client {
	t.Helper()
	defaultEndpoint, err := url.Parse(srv.URL())
	require.NoError(t, err)

	internalClient, err := azcore.NewClient("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: policies},
		&policy.ClientOptions{Transport: srv})
	require.NoError(t, err)

	containerCache := newContainerPropertiesCache()
	pkRangeCache := newPartitionKeyRangeCache()
	gem := &globalEndpointManager{preferredLocations: []string{}}

	client := &Client{
		endpoint:    srv.URL(),
		endpointUrl: defaultEndpoint,
		internal:    internalClient,
		gem:         gem,
		caches: &sharedCacheSet{
			containerCache: containerCache,
			pkRangeCache:   pkRangeCache,
		},
	}

	containerLink := "dbs/databaseId/colls/containerId"
	containerCache.set(containerLink, &ContainerProperties{
		ID:         "containerId",
		ResourceID: "testRID",
		PartitionKeyDefinition: PartitionKeyDefinition{
			Paths:   []string{"/pk"},
			Kind:    PartitionKeyKindHash,
			Version: 2,
		},
	})

	pkRangeCache.entries["testRID"] = &pkRangeCacheEntry{
		routingMap: newCollectionRoutingMap(ranges, "etag1"),
	}

	return client
}

// TestGetChangeFeed_410Gone_TriggersCacheRefreshAndRetry — Phase 1's headline win:
// when the gateway returns 410/Gone with a PK-range substatus, the change-feed
// request must (a) refresh the PK-range cache, (b) retry on the freshly-resolved
// range, (c) succeed transparently to the caller without losing the response.
func TestGetChangeFeed_410Gone_TriggersCacheRefreshAndRetry(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	gonePolicy := &returnGoneOnChangeFeedPolicy{maxGone: 1, substatus: subStatusPartitionKeyRangeGone}
	ranges := []partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF", ResourceID: "testRID"}}
	client := createChangeFeedTestClient(t, srv, []policy.Policy{gonePolicy}, ranges)

	containerPropsResp := []byte(`{
		"id": "containerId",
		"_rid": "testRID",
		"_self": "dbs/db1/colls/containerId/",
		"partitionKey": {"paths": ["/pk"], "kind": "Hash", "version": 2}
	}`)
	pkRangeResp := []byte(`{
		"_rid": "testRID",
		"PartitionKeyRanges": [{"_rid": "testRID", "id": "0", "minInclusive": "", "maxExclusive": "FF"}],
		"_count": 1
	}`)
	cfBody := []byte(`{"_rid":"testRID","Documents":[{"id":"doc1"}],"_count":1}`)

	// Sequence after the policy returns 410 once:
	//   container props re-fetch → PK ranges refresh → 304 (incremental loop end)
	//   → change-feed retry succeeds (passes through gonePolicy because maxGone=1)
	srv.AppendResponse(mock.WithBody(containerPropsResp), mock.WithStatusCode(200))
	srv.AppendResponse(mock.WithBody(pkRangeResp), mock.WithStatusCode(200),
		mock.WithHeader(cosmosHeaderEtag, "etag2"))
	srv.AppendResponse(mock.WithStatusCode(304))
	srv.AppendResponse(mock.WithBody(cfBody), mock.WithStatusCode(200),
		mock.WithHeader(cosmosHeaderEtag, "\"new-etag\""),
		mock.WithHeader(cosmosHeaderRequestCharge, "1.0"))

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		FeedRange: &FeedRange{MinInclusive: "00", MaxExclusive: "FF"},
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Count, "retry must surface the success page from the second attempt")
	require.Equal(t, int32(2), gonePolicy.count.Load(), "expected initial 410 + 1 retry on change-feed path")
	require.NotEmpty(t, resp.ContinuationToken, "response must carry a continuation token after a successful retry")
}

// TestGetChangeFeed_RetryCapAt3_ReturnsLastErrorOnRepeated410 — when the cache
// keeps returning the same range and 410s keep coming, the loop must surface the
// final 410 to the caller after exactly maxPKRangeGoneRetries attempts.
func TestGetChangeFeed_RetryCapAt3_ReturnsLastErrorOnRepeated410(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	gonePolicy := &returnGoneOnChangeFeedPolicy{maxGone: 100, substatus: subStatusPartitionKeyRangeGone}
	ranges := []partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF", ResourceID: "testRID"}}
	client := createChangeFeedTestClient(t, srv, []policy.Policy{gonePolicy}, ranges)

	containerPropsResp := []byte(`{
		"id": "containerId",
		"_rid": "testRID",
		"_self": "dbs/db1/colls/containerId/",
		"partitionKey": {"paths": ["/pk"], "kind": "Hash", "version": 2}
	}`)
	pkRangeResp := []byte(`{
		"_rid": "testRID",
		"PartitionKeyRanges": [{"_rid": "testRID", "id": "0", "minInclusive": "", "maxExclusive": "FF"}],
		"_count": 1
	}`)

	// Each retry needs container props + PK ranges (with ETag) + 304 to terminate the
	// incremental refresh loop. After maxPKRangeGoneRetries refreshes, the next 410 is
	// surfaced to the caller without further retries.
	for i := 0; i < maxPKRangeGoneRetries; i++ {
		srv.AppendResponse(mock.WithBody(containerPropsResp), mock.WithStatusCode(200))
		srv.AppendResponse(mock.WithBody(pkRangeResp), mock.WithStatusCode(200),
			mock.WithHeader(cosmosHeaderEtag, "etag-refresh"))
		srv.AppendResponse(mock.WithStatusCode(304))
	}

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	_, err := container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		FeedRange: &FeedRange{MinInclusive: "00", MaxExclusive: "FF"},
	})
	require.Error(t, err)
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusGone, respErr.StatusCode)
	require.Equal(t, int32(maxPKRangeGoneRetries+1), gonePolicy.count.Load(),
		"expected initial attempt + maxPKRangeGoneRetries retries")
}

// TestGetChangeFeed_TokenResourceIDMismatch_Rejected — the cross-container token
// reuse guard. A token whose ResourceID doesn't match the current container's
// must be rejected loudly; otherwise the EPK boundaries in the token would be
// misinterpreted against the wrong routing map and silent wrong data could leak.
func TestGetChangeFeed_TokenResourceIDMismatch_Rejected(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	ranges := []partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF", ResourceID: "testRID"}}
	client := createChangeFeedTestClient(t, srv, nil, ranges)

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	mismatchToken := &compositeContinuationToken{
		Version:    cosmosCompositeContinuationTokenVersion,
		ResourceID: "differentContainerRID",
		Continuation: []changeFeedRange{
			{MinInclusive: "00", MaxExclusive: "FF"},
		},
	}
	tokenJSON, err := json.Marshal(mismatchToken)
	require.NoError(t, err)
	tokenStr := string(tokenJSON)

	_, err = container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		Continuation: &tokenStr,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ResourceID",
		"error message must call out ResourceID mismatch so customers can diagnose token-reuse bugs")
	require.Contains(t, err.Error(), "differentContainerRID")
	require.Contains(t, err.Error(), "testRID")
}

// TestGetChangeFeed_NoOverlapAfterRefresh_ReturnsErrFeedRangeUnresolved — when
// the customer's FeedRange genuinely doesn't overlap any current physical range
// (e.g., wrong container, malformed range, container recreated under same name
// with different boundaries), even a forced cache refresh can't help. The error
// MUST be wrapped as ErrFeedRangeUnresolved so callers can detect-and-recover
// rather than seeing a generic opaque failure.
func TestGetChangeFeed_NoOverlapAfterRefresh_ReturnsErrFeedRangeUnresolved(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// The cached routing map is complete (covers ["", "FF")), so the routing-map's
	// own completeness check is satisfied. The customer's FeedRange ["GG", "HH") is
	// entirely above the routed key space — a malformed/foreign range that no refresh
	// can rescue. The cache stays fresh; we never hit HTTP.
	ranges := []partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF", ResourceID: "testRID"}}
	client := createChangeFeedTestClient(t, srv, nil, ranges)

	// The forced refresh hits these; we keep the same complete-but-non-overlapping range.
	containerPropsResp := []byte(`{
		"id": "containerId",
		"_rid": "testRID",
		"_self": "dbs/db1/colls/containerId/",
		"partitionKey": {"paths": ["/pk"], "kind": "Hash", "version": 2}
	}`)
	pkRangeResp := []byte(`{
		"_rid": "testRID",
		"PartitionKeyRanges": [{"_rid": "testRID", "id": "0", "minInclusive": "", "maxExclusive": "FF"}],
		"_count": 1
	}`)
	srv.AppendResponse(mock.WithBody(containerPropsResp), mock.WithStatusCode(200))
	srv.AppendResponse(mock.WithBody(pkRangeResp), mock.WithStatusCode(200),
		mock.WithHeader(cosmosHeaderEtag, "etag-refresh"))
	srv.AppendResponse(mock.WithStatusCode(304))

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	_, err := container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		FeedRange: &FeedRange{MinInclusive: "GG", MaxExclusive: "HH"},
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrFeedRangeUnresolved),
		"unresolvable FeedRange MUST be wrapped as ErrFeedRangeUnresolved; got %T: %v", err, err)
}

// TestGetChangeFeed_SplitExpansion_RoutesToFirstChild — when the customer's
// FeedRange overlaps multiple physical ranges (the post-split case), the loop
// must replace the head with one queue entry per child and route the actual
// HTTP call against the first child. The returned token must encode the full
// child queue so the next call drains the remaining children.
func TestGetChangeFeed_SplitExpansion_RoutesToFirstChild(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// Cache has two children where the customer's parent range used to be.
	ranges := []partitionKeyRange{
		{ID: "1a", MinInclusive: "00", MaxExclusive: "55", ResourceID: "testRID"},
		{ID: "1b", MinInclusive: "55", MaxExclusive: "FF", ResourceID: "testRID"},
	}
	client := createChangeFeedTestClient(t, srv, nil, ranges)

	cfBody := []byte(`{"_rid":"testRID","Documents":[{"id":"docFromFirstChild"}],"_count":1}`)
	srv.AppendResponse(mock.WithBody(cfBody), mock.WithStatusCode(200),
		mock.WithHeader(cosmosHeaderEtag, "\"etag-1a\""))

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	resp, err := container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		FeedRange: &FeedRange{MinInclusive: "00", MaxExclusive: "FF"},
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Count)
	require.NotEmpty(t, resp.ContinuationToken)

	// The persisted token must enumerate BOTH children so a subsequent call
	// continues with 1b. After a successful 200 from 1a, advance() rotates 1a
	// to the tail with its new ETag — so the queue is [1b, 1a-with-etag].
	var token compositeContinuationToken
	require.NoError(t, json.Unmarshal([]byte(resp.ContinuationToken), &token))
	require.Len(t, token.Continuation, 2, "split expansion must surface both children in the continuation queue")
	require.Equal(t, "55", token.Continuation[0].MinInclusive,
		"after one successful 200, the new head must be the unread sibling (1b at [55, FF))")
	require.Equal(t, "00", token.Continuation[1].MinInclusive,
		"completed sub-range must be rotated to the tail with its updated ETag")
	require.NotNil(t, token.Continuation[1].ContinuationToken, "tail entry must carry the freshly-issued ETag")
	require.Equal(t, azcore.ETag("\"etag-1a\""), *token.Continuation[1].ContinuationToken)
}

// TestGetChangeFeed_MultiRangeContinuation_RoundTrips — a token issued by an
// earlier call (multi-element queue) drives the next call against the queue's
// head, not against options.FeedRange. This is the resume-after-split pattern.
func TestGetChangeFeed_MultiRangeContinuation_RoundTrips(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	ranges := []partitionKeyRange{
		{ID: "1a", MinInclusive: "00", MaxExclusive: "55", ResourceID: "testRID"},
		{ID: "1b", MinInclusive: "55", MaxExclusive: "FF", ResourceID: "testRID"},
	}
	client := createChangeFeedTestClient(t, srv, nil, ranges)

	cfBody := []byte(`{"_rid":"testRID","Documents":[{"id":"docFromHead"}],"_count":1}`)
	srv.AppendResponse(mock.WithBody(cfBody), mock.WithStatusCode(200),
		mock.WithHeader(cosmosHeaderEtag, "\"new-etag\""))

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	// Customer-supplied resume token: queue head is the second child (the unfinished one).
	headETag := azcore.ETag("\"prev-etag-1b\"")
	tailETag := azcore.ETag("\"prev-etag-1a\"")
	multiRangeToken := compositeContinuationToken{
		Version:    cosmosCompositeContinuationTokenVersion,
		ResourceID: "testRID",
		Continuation: []changeFeedRange{
			{MinInclusive: "55", MaxExclusive: "FF", ContinuationToken: &headETag},
			{MinInclusive: "00", MaxExclusive: "55", ContinuationToken: &tailETag},
		},
	}
	tokenJSON, err := json.Marshal(multiRangeToken)
	require.NoError(t, err)
	tokenStr := string(tokenJSON)

	resp, err := container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		Continuation: &tokenStr,
	})
	require.NoError(t, err)
	require.Equal(t, 1, resp.Count)

	// After this call, the queue must rotate so the previously-tail entry is now the head.
	var rt compositeContinuationToken
	require.NoError(t, json.Unmarshal([]byte(resp.ContinuationToken), &rt))
	require.Len(t, rt.Continuation, 2, "queue length must be preserved across calls")
	require.Equal(t, "00", rt.Continuation[0].MinInclusive, "head should rotate to the unfinished sibling")
	require.Equal(t, "55", rt.Continuation[1].MinInclusive, "completed entry should rotate to the tail")
	require.NotNil(t, rt.Continuation[1].ContinuationToken)
	require.Equal(t, azcore.ETag("\"new-etag\""), *rt.Continuation[1].ContinuationToken,
		"tail entry must carry the freshly-issued ETag from this call's response")
}

// TestGetChangeFeed_SplitDuringDrain_QueriesEveryNewlyInsertedChild — regression
// guard for a subtle accounting bug: when one of the queued sub-ranges splits
// MID-drain (i.e., the cache learns about new children in the middle of a single
// GetChangeFeed call's drain rotation), the rotation budget must be expanded to
// account for the inserted siblings. Otherwise the loop bails early and silently
// skips queryable ranges this call.
//
// Setup: two queued entries (A=[00,55), B=[55,FF)). The cache reflects a split
// of B into B1=[55,AA) and B2=[AA,FF). The drain sequence MUST be:
//
//	A → 304, rotate. Queue: [B, A]. budget consumed=1/2.
//	B → split detected, replaced with [B1, B2]. Queue: [B1, B2, A]. budget bumped to 3.
//	B1 → 304, rotate. Queue: [B2, A, B1]. budget consumed=2/3.
//	B2 → 304, rotate. Queue: [A, B1, B2]. budget consumed=3/3 → exit.
//
// Without the budget bump, the loop would break after B1 (rotations==2 ==
// originalQueueLen==2) and never query B2 in this call. The test therefore
// asserts that THREE distinct CF reads happened (one for A, one for B1,
// one for B2), not two.
func TestGetChangeFeed_SplitDuringDrain_QueriesEveryNewlyInsertedChild(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	// The cache reflects the post-split topology: A is unchanged, B has split into B1+B2.
	ranges := []partitionKeyRange{
		{ID: "A", MinInclusive: "00", MaxExclusive: "55", ResourceID: "testRID"},
		{ID: "B1", MinInclusive: "55", MaxExclusive: "AA", ResourceID: "testRID"},
		{ID: "B2", MinInclusive: "AA", MaxExclusive: "FF", ResourceID: "testRID"},
	}

	// Count CF requests AND remember which PK-range-id was queried each time so we can
	// assert B2 is actually visited in this single call (the regression target).
	var cfCount atomic.Int32
	queriedPKRangeIDs := make([]string, 0, 8)
	var mu sync.Mutex
	tracker := policyFunc(func(req *policy.Request) (*http.Response, error) {
		if req.Raw().Header.Get(cosmosHeaderChangeFeed) == cosmosHeaderValuesChangeFeed &&
			strings.HasSuffix(req.Raw().URL.Path, "/docs") {
			cfCount.Add(1)
			mu.Lock()
			queriedPKRangeIDs = append(queriedPKRangeIDs, req.Raw().Header.Get(cosmosHeaderPartitionKeyRangeId))
			mu.Unlock()
		}
		return req.Next()
	})
	client := createChangeFeedTestClient(t, srv, []policy.Policy{tracker}, ranges)

	// All three sub-range reads return 304 (no new changes). The drain must hit every one.
	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, "\"etag-A\""))
	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, "\"etag-B1\""))
	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, "\"etag-B2\""))

	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	// Resume token preserves the original two-element queue (A, B). The drain must
	// expand B → [B1, B2] mid-rotation and still query B2 within this call.
	prevA := azcore.ETag("\"prev-A\"")
	prevB := azcore.ETag("\"prev-B\"")
	resumeToken := compositeContinuationToken{
		Version:    cosmosCompositeContinuationTokenVersion,
		ResourceID: "testRID",
		Continuation: []changeFeedRange{
			{MinInclusive: "00", MaxExclusive: "55", ContinuationToken: &prevA},
			{MinInclusive: "55", MaxExclusive: "FF", ContinuationToken: &prevB},
		},
	}
	tokenJSON, err := json.Marshal(resumeToken)
	require.NoError(t, err)
	tokenStr := string(tokenJSON)

	_, err = container.GetChangeFeed(context.Background(), &ChangeFeedOptions{
		Continuation: &tokenStr,
	})
	require.NoError(t, err)

	require.Equal(t, int32(3), cfCount.Load(),
		"split-during-drain must query every newly-inserted child in the same call (got %d, want 3); queried IDs=%v",
		cfCount.Load(), queriedPKRangeIDs)

	// The actual order of visitation should be A, B1, B2 (FIFO with split-expand-on-head).
	mu.Lock()
	got := append([]string(nil), queriedPKRangeIDs...)
	mu.Unlock()
	require.Equal(t, []string{"A", "B1", "B2"}, got,
		"drain order must be A, B1, B2 (FIFO with split expansion on the head)")
}

// policyFunc adapts a function to policy.Policy so we can use closures inline.
type policyFunc func(*policy.Request) (*http.Response, error)

func (f policyFunc) Do(req *policy.Request) (*http.Response, error) { return f(req) }

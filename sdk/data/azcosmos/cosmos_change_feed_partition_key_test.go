// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// changeFeedRequest is what changeFeedRequestRecorder captures from one change
// feed request.
type changeFeedRequest struct {
	pkRangeID    string
	partitionKey string
	ifNoneMatch  string
}

// changeFeedRequestRecorder records the routing headers of every change feed
// request (not partition key range reads) that passes through it.
type changeFeedRequestRecorder struct {
	mu       sync.Mutex
	requests []changeFeedRequest
}

func (r *changeFeedRequestRecorder) Do(req *policy.Request) (*http.Response, error) {
	h := req.Raw().Header
	if h.Get(cosmosHeaderChangeFeed) == cosmosHeaderValuesChangeFeed && strings.HasSuffix(req.Raw().URL.Path, "/docs") {
		r.mu.Lock()
		r.requests = append(r.requests, changeFeedRequest{
			pkRangeID:    h.Get(cosmosHeaderPartitionKeyRangeId),
			partitionKey: h.Get(cosmosHeaderPartitionKey),
			ifNoneMatch:  h.Get(headerIfNoneMatch),
		})
		r.mu.Unlock()
	}
	return req.Next()
}

func (r *changeFeedRequestRecorder) recorded() []changeFeedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]changeFeedRequest(nil), r.requests...)
}

// testPartitionKeyEPK returns the EPK of pk in the container that
// createChangeFeedTestClient sets up (Hash v2 on /pk).
func testPartitionKeyEPK(t *testing.T, pk PartitionKey) string {
	t.Helper()
	r, err := computeEPKRange(&pk, PartitionKeyDefinition{Paths: []string{"/pk"}, Kind: PartitionKeyKindHash, Version: 2})
	require.NoError(t, err)
	require.Equal(t, r.Min, r.Max)
	require.NotEmpty(t, r.Min)
	return r.Min
}

func decodeChangeFeedToken(t *testing.T, s string) compositeContinuationToken {
	t.Helper()
	var token compositeContinuationToken
	require.NoError(t, json.Unmarshal([]byte(s), &token))
	return token
}

// TestReadChangeFeed_PartitionKey_LetsServiceRoute is the regression test for
// partition key reads on containers with several physical partitions. Naming
// a range that doesn't hold the key gets 410 Gone from the service, so the
// requests must carry the partition key and no partition key range ID.
func TestReadChangeFeed_PartitionKey_LetsServiceRoute(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	pk := NewPartitionKeyString("customer-42")
	pkEPK := testPartitionKeyEPK(t, pk)
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: pkEPK, ResourceID: "testRID"},
		{ID: "1", MinInclusive: pkEPK, MaxExclusive: "FF", ResourceID: "testRID"},
	}
	recorder := &changeFeedRequestRecorder{}
	client := createChangeFeedTestClient(t, srv, []policy.Policy{recorder}, ranges)
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	srv.AppendResponse(mock.WithStatusCode(200),
		mock.WithBody([]byte(`{"_rid":"testRID","Documents":[{"id":"order-1","pk":"customer-42"}],"_count":1}`)),
		mock.WithHeader(cosmosHeaderEtag, `"etag-1"`))
	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, `"etag-1"`))

	// No FeedRange: the partition key alone scopes the read.
	resp, err := container.ReadChangeFeed(context.Background(), &ChangeFeedOptions{PartitionKey: &pk})
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Count)

	token := decodeChangeFeedToken(t, resp.ContinuationToken)
	require.Len(t, token.Continuation, 1, "a partition key read must track one queue entry")
	assert.Equal(t, pkEPK, token.Continuation[0].MinInclusive)
	assert.Equal(t, pkEPK, token.Continuation[0].MaxExclusive)
	etag1 := azcore.ETag(`"etag-1"`)
	assert.Equal(t, &etag1, token.Continuation[0].ContinuationToken)

	resp, err = container.ReadChangeFeed(context.Background(), &ChangeFeedOptions{
		PartitionKey: &pk,
		Continuation: &resp.ContinuationToken,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotModified, resp.RawResponse.StatusCode)

	pkJSON, err := pk.toJsonString()
	require.NoError(t, err)
	assert.Equal(t, []changeFeedRequest{
		{partitionKey: string(pkJSON)},
		{partitionKey: string(pkJSON), ifNoneMatch: `"etag-1"`},
	}, recorder.recorded(), "requests must carry the partition key and no partition key range ID")
}

// TestReadChangeFeed_PartitionKey_ResumesFromFeedRangeToken checks that a
// token issued for a whole FeedRange (one entry per physical partition) can
// resume a partition key read: only the entry that covers the key is kept,
// with its ETag.
func TestReadChangeFeed_PartitionKey_ResumesFromFeedRangeToken(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	pk := NewPartitionKeyString("customer-42")
	pkEPK := testPartitionKeyEPK(t, pk)
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: pkEPK, ResourceID: "testRID"},
		{ID: "1", MinInclusive: pkEPK, MaxExclusive: "FF", ResourceID: "testRID"},
	}
	recorder := &changeFeedRequestRecorder{}
	client := createChangeFeedTestClient(t, srv, []policy.Policy{recorder}, ranges)
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	srv.AppendResponse(mock.WithStatusCode(304), mock.WithHeader(cosmosHeaderEtag, `"etag-1"`))

	etag0, etag1 := azcore.ETag(`"etag-0"`), azcore.ETag(`"etag-1"`)
	tokenJSON, err := json.Marshal(compositeContinuationToken{
		Version:    cosmosCompositeContinuationTokenVersion,
		ResourceID: "testRID",
		Continuation: []changeFeedRange{
			{MinInclusive: "", MaxExclusive: pkEPK, ContinuationToken: &etag0},
			{MinInclusive: pkEPK, MaxExclusive: "FF", ContinuationToken: &etag1},
		},
	})
	require.NoError(t, err)
	tokenStr := string(tokenJSON)

	resp, err := container.ReadChangeFeed(context.Background(), &ChangeFeedOptions{
		PartitionKey: &pk,
		FeedRange:    &FeedRange{MinInclusive: "", MaxExclusive: "FF"},
		Continuation: &tokenStr,
	})
	require.NoError(t, err)

	requests := recorder.recorded()
	require.Len(t, requests, 1)
	assert.Empty(t, requests[0].pkRangeID)
	assert.Equal(t, `"etag-1"`, requests[0].ifNoneMatch, "the ETag of the entry covering the key must be used")

	token := decodeChangeFeedToken(t, resp.ContinuationToken)
	require.Len(t, token.Continuation, 1)
	assert.Equal(t, pkEPK, token.Continuation[0].MinInclusive)
	assert.Equal(t, pkEPK, token.Continuation[0].MaxExclusive)
}

// TestReadChangeFeed_PartitionKey_RejectsTokenForOtherKey checks that a token
// issued for one partition key can't silently resume a different one.
func TestReadChangeFeed_PartitionKey_RejectsTokenForOtherKey(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	ranges := []partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF", ResourceID: "testRID"}}
	client := createChangeFeedTestClient(t, srv, nil, ranges)
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	otherEPK := testPartitionKeyEPK(t, NewPartitionKeyString("customer-7"))
	etag := azcore.ETag(`"etag-7"`)
	tokenJSON, err := json.Marshal(compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
		ResourceID:   "testRID",
		Continuation: []changeFeedRange{{MinInclusive: otherEPK, MaxExclusive: otherEPK, ContinuationToken: &etag}},
	})
	require.NoError(t, err)
	tokenStr := string(tokenJSON)

	pk := NewPartitionKeyString("customer-42")
	_, err = container.ReadChangeFeed(context.Background(), &ChangeFeedOptions{
		PartitionKey: &pk,
		Continuation: &tokenStr,
	})
	assert.ErrorContains(t, err, "doesn't cover ChangeFeedOptions.PartitionKey")
}

// TestReadChangeFeed_PartitionKeyToken_RequiresPartitionKey checks that a
// token from a partition key read can't resume without the partition key.
// Without it the request would carry no key, so the error must say what's
// missing instead of reporting an unresolved FeedRange.
func TestReadChangeFeed_PartitionKeyToken_RequiresPartitionKey(t *testing.T) {
	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()

	recorder := &changeFeedRequestRecorder{}
	ranges := []partitionKeyRange{{ID: "0", MinInclusive: "", MaxExclusive: "FF", ResourceID: "testRID"}}
	client := createChangeFeedTestClient(t, srv, []policy.Policy{recorder}, ranges)
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	pkEPK := testPartitionKeyEPK(t, NewPartitionKeyString("customer-42"))
	etag := azcore.ETag(`"etag-42"`)
	tokenJSON, err := json.Marshal(compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
		ResourceID:   "testRID",
		Continuation: []changeFeedRange{{MinInclusive: pkEPK, MaxExclusive: pkEPK, ContinuationToken: &etag}},
	})
	require.NoError(t, err)
	tokenStr := string(tokenJSON)

	_, err = container.ReadChangeFeed(context.Background(), &ChangeFeedOptions{Continuation: &tokenStr})
	require.ErrorContains(t, err, "set ChangeFeedOptions.PartitionKey to the same partition key")
	assert.NotErrorIs(t, err, ErrFeedRangeUnresolved)
	assert.Empty(t, recorder.recorded(), "no change feed request may be sent without the partition key")
}

func TestScopeChangeFeedQueueToEPK(t *testing.T) {
	etagA, etagB := azcore.ETag(`"a"`), azcore.ETag(`"b"`)
	queue := []changeFeedRange{
		{MinInclusive: "", MaxExclusive: "40", ContinuationToken: &etagA},
		{MinInclusive: "40", MaxExclusive: "FF", ContinuationToken: &etagB},
	}

	scoped, err := scopeChangeFeedQueueToEPK(queue, "40")
	require.NoError(t, err)
	assert.Equal(t, []changeFeedRange{{MinInclusive: "40", MaxExclusive: "40", ContinuationToken: &etagB}}, scoped,
		"MinInclusive is inclusive")

	scoped, err = scopeChangeFeedQueueToEPK(queue, "3F")
	require.NoError(t, err)
	assert.Equal(t, &etagA, scoped[0].ContinuationToken, "MaxExclusive is exclusive")

	scoped, err = scopeChangeFeedQueueToEPK(scoped, "3F")
	require.NoError(t, err, "a token scoped to a key resumes the same key")
	assert.Equal(t, "3F", scoped[0].MinInclusive)

	_, err = scopeChangeFeedQueueToEPK(scoped, "3E")
	assert.Error(t, err, "a token scoped to one key must not resume another")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestFindPhysicalRangeForEPK_SingleRange(t *testing.T) {
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}
	id, ok := findPhysicalRangeForEPK("05C1DFECD6D5C72619A4ACA689462EB6", ranges)
	require.True(t, ok)
	require.Equal(t, "0", id)
}

func TestFindPhysicalRangeForEPK_MultipleRanges(t *testing.T) {
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "05C1E0"},
		{ID: "1", MinInclusive: "05C1E0", MaxExclusive: "0BED02"},
		{ID: "2", MinInclusive: "0BED02", MaxExclusive: "FF"},
	}

	// Value that falls in range 0
	id, ok := findPhysicalRangeForEPK("05C1DF", ranges)
	require.True(t, ok)
	require.Equal(t, "0", id)

	// Value at boundary of range 1
	id, ok = findPhysicalRangeForEPK("05C1E0", ranges)
	require.True(t, ok)
	require.Equal(t, "1", id)

	// Value in range 2
	id, ok = findPhysicalRangeForEPK("0BED03", ranges)
	require.True(t, ok)
	require.Equal(t, "2", id)
}

func TestFindPhysicalRangeForEPK_NotFound(t *testing.T) {
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "05C1E0", MaxExclusive: "0BED02"},
	}

	// Value below all ranges
	_, ok := findPhysicalRangeForEPK("00000000", ranges)
	require.False(t, ok)

	// Value above all ranges
	_, ok = findPhysicalRangeForEPK("FF", ranges)
	require.False(t, ok)
}

func TestGroupItemsByPhysicalRange_SingleRange(t *testing.T) {
	pkDef := PartitionKeyDefinition{
		Paths:   []string{"/pk"},
		Kind:    PartitionKeyKindHash,
		Version: 2,
	}
	// Single range covers everything
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}

	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString("pkA")},
		{ID: "2", PartitionKey: NewPartitionKeyString("pkB")},
		{ID: "3", PartitionKey: NewPartitionKeyString("pkA")},
	}

	orderedIDs, groups, err := groupItemsByPhysicalRange(items, pkDef, ranges)
	require.NoError(t, err)
	// All items should land in range "0"
	require.Equal(t, []string{"0"}, orderedIDs)
	require.Len(t, groups["0"], 3)
	require.Equal(t, "1", groups["0"][0].ID)
	require.Equal(t, "2", groups["0"][1].ID)
	require.Equal(t, "3", groups["0"][2].ID)
}

func TestGroupItemsByPhysicalRange_MultipleRanges(t *testing.T) {
	pkDef := PartitionKeyDefinition{
		Paths:   []string{"/pk"},
		Kind:    PartitionKeyKindHash,
		Version: 2,
	}

	// Compute EPKs to set up ranges that actually split the items
	pkAVal := NewPartitionKeyString("pkA")
	pkBVal := NewPartitionKeyString("pkB")
	epkA := pkAVal.computeEffectivePartitionKey(PartitionKeyKindHash, 2)
	epkB := pkBVal.computeEffectivePartitionKey(PartitionKeyKindHash, 2)

	// Determine which EPK is smaller
	var lowEPK, highEPK string
	var lowPK, highPK string
	if epkA.EPK < epkB.EPK {
		lowEPK = epkA.EPK
		highEPK = epkB.EPK
		lowPK = "pkA"
		highPK = "pkB"
	} else {
		lowEPK = epkB.EPK
		highEPK = epkA.EPK
		lowPK = "pkB"
		highPK = "pkA"
	}

	// Create ranges that separate the two EPKs.
	// midpoint: use highEPK as the boundary (it goes to range 1).
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: highEPK},
		{ID: "1", MinInclusive: highEPK, MaxExclusive: "FF"},
	}
	_ = lowEPK // used implicitly

	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString(lowPK)},
		{ID: "2", PartitionKey: NewPartitionKeyString(highPK)},
		{ID: "3", PartitionKey: NewPartitionKeyString(lowPK)},
	}

	orderedIDs, groups, err := groupItemsByPhysicalRange(items, pkDef, ranges)
	require.NoError(t, err)
	require.Len(t, orderedIDs, 2)

	// lowPK items should be in range "0", highPK in range "1"
	require.Len(t, groups["0"], 2)
	require.Len(t, groups["1"], 1)
	require.Equal(t, "1", groups["0"][0].ID)
	require.Equal(t, "3", groups["0"][1].ID)
	require.Equal(t, "2", groups["1"][0].ID)
}

func TestGroupItemsByPhysicalRange_DefaultVersion(t *testing.T) {
	// Version 0 should default to 1
	pkDef := PartitionKeyDefinition{
		Paths:   []string{"/pk"},
		Kind:    PartitionKeyKindHash,
		Version: 0,
	}
	ranges := []partitionKeyRange{
		{ID: "0", MinInclusive: "", MaxExclusive: "FF"},
	}

	items := []ItemIdentity{
		{ID: "1", PartitionKey: NewPartitionKeyString("test")},
	}

	orderedIDs, groups, err := groupItemsByPhysicalRange(items, pkDef, ranges)
	require.NoError(t, err)
	require.Equal(t, []string{"0"}, orderedIDs)
	require.Len(t, groups["0"], 1)
}

func TestBuildQueryChunksForRanges_SingleRange(t *testing.T) {
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	orderedIDs := []string{"0"}
	groups := map[string][]ItemIdentity{
		"0": {
			{ID: "1", PartitionKey: NewPartitionKeyString("pkA")},
			{ID: "2", PartitionKey: NewPartitionKeyString("pkB")},
		},
	}

	chunks, err := buildQueryChunksForRanges(orderedIDs, groups, pkDef)
	require.NoError(t, err)
	require.Len(t, chunks, 1)
	require.Equal(t, "0", chunks[0].rangeID)
	// Items have different PKs → OR-of-conjunctions
	require.Equal(t, "SELECT * FROM c WHERE (c.id = @param_id0 AND c.pk = @param_pk00) OR (c.id = @param_id1 AND c.pk = @param_pk10)", chunks[0].query)
}

func TestBuildQueryChunksForRanges_Chunking(t *testing.T) {
	pkDef := PartitionKeyDefinition{Paths: []string{"/pk"}}
	items := make([]ItemIdentity, maxItemsPerQuery+1)
	for i := range items {
		items[i] = ItemIdentity{ID: "id", PartitionKey: NewPartitionKeyString("pk")}
	}
	orderedIDs := []string{"0"}
	groups := map[string][]ItemIdentity{"0": items}

	chunks, err := buildQueryChunksForRanges(orderedIDs, groups, pkDef)
	require.NoError(t, err)
	require.Len(t, chunks, 2)
	require.Equal(t, "0", chunks[0].rangeID)
	require.Equal(t, "0", chunks[1].rangeID)
}

// cancelOnNthQueryPolicy is a pipeline policy that cancels a context after the
// Nth query request completes, providing deterministic mid-execution cancellation.
type cancelOnNthQueryPolicy struct {
	cancelAfterN int32         // cancel context after this many query requests complete
	count        atomic.Int32  // number of query requests seen so far
	firstDone    chan struct{} // closed after the Nth query completes
	gate         chan struct{} // blocks subsequent queries until signalled
}

func (p *cancelOnNthQueryPolicy) Do(req *policy.Request) (*http.Response, error) {
	isQuery := req.Raw().Header.Get(cosmosHeaderQuery) == "True"
	if !isQuery {
		return req.Next()
	}

	n := p.count.Add(1)
	if n <= p.cancelAfterN {
		// Let the first N queries proceed normally
		resp, err := req.Next()
		if n == p.cancelAfterN {
			close(p.firstDone)
		}
		return resp, err
	}

	// For queries beyond N, block until gate is opened (context cancellation)
	<-p.gate
	return req.Next()
}

func TestExecuteQueryChunks_CancelledContext(t *testing.T) {
	queryResp := []byte(`{"Documents":[{"id":"1","pk":"pkA"}]}`)

	srv, closeSrv := mock.NewTLSServer()
	defer closeSrv()
	srv.SetResponse(
		mock.WithBody(queryResp),
		mock.WithHeader(cosmosHeaderRequestCharge, "1.0"),
		mock.WithStatusCode(200))

	cancelPolicy := &cancelOnNthQueryPolicy{
		cancelAfterN: 1,
		firstDone:    make(chan struct{}),
		gate:         make(chan struct{}),
	}

	defaultEndpoint, _ := url.Parse(srv.URL())
	internalClient, _ := azcore.NewClient("azcosmostest", "v1.0.0",
		azruntime.PipelineOptions{PerCall: []policy.Policy{cancelPolicy}},
		&policy.ClientOptions{Transport: srv})

	gem := &globalEndpointManager{preferredLocations: []string{}}
	client := &Client{endpoint: srv.URL(), endpointUrl: defaultEndpoint, internal: internalClient, gem: gem}
	database, _ := newDatabase("databaseId", client)
	container, _ := newContainer("containerId", database)

	chunks := []queryChunk{
		{query: "SELECT * FROM c WHERE c.id = @id0", params: []QueryParameter{{Name: "@id0", Value: "1"}}, rangeID: "0"},
		{query: "SELECT * FROM c WHERE c.id = @id1", params: []QueryParameter{{Name: "@id1", Value: "2"}}, rangeID: "1"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel context after the first chunk completes, then unblock the gate
	go func() {
		<-cancelPolicy.firstDone
		cancel()
		close(cancelPolicy.gate)
	}()

	opCtx := pipelineRequestOptions{
		resourceType:    resourceTypeDocument,
		resourceAddress: "dbs/databaseId/colls/containerId",
	}

	results, err := container.executeQueryChunks(ctx, chunks, &QueryOptions{}, opCtx, 2)
	require.NoError(t, err)

	// One chunk should succeed, the other should have a context cancellation error
	var successCount, errorCount int
	for _, r := range results {
		if r.err != nil {
			errorCount++
			require.ErrorIs(t, r.err, context.Canceled)
		} else {
			successCount++
			require.Len(t, r.items, 1)
		}
	}
	require.Equal(t, 1, successCount, "expected exactly 1 successful chunk")
	require.Equal(t, 1, errorCount, "expected exactly 1 cancelled chunk")

	// collectChunkResults should surface the cancellation as an operation-level error
	resp, err := collectChunkResults(results)
	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, resp.Items)
}

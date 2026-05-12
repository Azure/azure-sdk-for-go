// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/internal/epk"
)

const maxItemsPerQuery = 1000
const maxPKRangeGoneRetries = 3

// queryChunk is a single parameterized query targeting one physical partition key range.
type queryChunk struct {
	query   string
	params  []QueryParameter
	rangeID string // physical partition key range ID for x-ms-documentdb-partitionkeyrangeid header
}

// chunkResult holds the outcome of executing a single queryChunk.
type chunkResult struct {
	items         [][]byte
	requestCharge float32
	err           error
}

// findPhysicalRangeForEPK finds the partition key range that contains the given
// EPK value. The ranges must be sorted by MinInclusive ascending. Returns the
// range ID and true if found, or ("", false) otherwise.
func findPhysicalRangeForEPK(epkValue string, ranges []partitionKeyRange) (string, bool) {
	// Binary search: find the last range whose MinInclusive <= epkValue.
	// Uses length-aware comparison for HPK containers with mixed-length EPK boundaries.
	idx := sort.Search(len(ranges), func(i int) bool {
		return epk.CompareEPK(ranges[i].MinInclusive, epkValue) > 0
	}) - 1
	if idx < 0 {
		return "", false
	}
	// Verify epkValue < MaxExclusive.
	// Empty MaxExclusive or "FF" means unbounded (the last partition).
	r := ranges[idx]
	if r.MaxExclusive != "" && r.MaxExclusive != "FF" && epk.CompareEPK(epkValue, r.MaxExclusive) >= 0 {
		return "", false
	}
	return r.ID, true
}

// groupItemsByPhysicalRange computes the EPK for each item and groups them by
// physical partition range. It returns the range IDs in first-seen order and
// the groups keyed by range ID.
//
// For MultiHash containers, items with partial partition keys (fewer components
// than paths) are fanned out to all overlapping physical partitions via EPK
// range computation.
func groupItemsByPhysicalRange(items []ItemIdentity, pkDef PartitionKeyDefinition, ranges []partitionKeyRange) ([]string, map[string][]ItemIdentity, error) {
	// Build a routing map for efficient lookups.
	routingMap := newCollectionRoutingMap(ranges, "")

	order := make([]string, 0)
	seen := make(map[string]bool)
	groups := make(map[string][]ItemIdentity)

	for _, item := range items {
		epkR, err := computeEPKRange(&item.PartitionKey, pkDef)
		if err != nil {
			return nil, nil, err
		}

		if epkR.isRange() {
			// Prefix key: fan out to all overlapping ranges
			overlapping := routingMap.getOverlappingRanges(epkR.Min, epkR.Max)
			if len(overlapping) == 0 {
				return nil, nil, errors.New("could not find physical partition range for item EPK range")
			}
			for _, r := range overlapping {
				if !seen[r.ID] {
					order = append(order, r.ID)
					seen[r.ID] = true
				}
				groups[r.ID] = append(groups[r.ID], item)
			}
		} else {
			// Point key: direct lookup
			rangeID, ok := findPhysicalRangeForEPK(epkR.Min, routingMap.orderedRanges)
			if !ok {
				return nil, nil, errors.New("could not find physical partition range for item EPK")
			}
			if !seen[rangeID] {
				order = append(order, rangeID)
				seen[rangeID] = true
			}
			groups[rangeID] = append(groups[rangeID], item)
		}
	}

	return order, groups, nil
}

// buildQueryChunksForRanges creates queryChunk values by splitting each physical
// range's item list into slices of at most maxItemsPerQuery and building
// parameterized SQL for each.
func buildQueryChunksForRanges(orderedRangeIDs []string, groups map[string][]ItemIdentity, pkDef PartitionKeyDefinition) ([]queryChunk, error) {
	qb := queryBuilder{}
	var chunks []queryChunk

	for _, rangeID := range orderedRangeIDs {
		items := groups[rangeID]
		for start := 0; start < len(items); start += maxItemsPerQuery {
			end := start + maxItemsPerQuery
			if end > len(items) {
				end = len(items)
			}
			q, params := qb.buildParameterizedQueryForItems(items[start:end], pkDef)
			chunks = append(chunks, queryChunk{
				query:   q,
				params:  params,
				rangeID: rangeID,
			})
		}
	}
	return chunks, nil
}

// executeQueryChunks runs the provided query chunks concurrently using a
// goroutine pool and returns the per-chunk results. It creates a child context
// so that when any chunk fails, all in-flight sibling requests are cancelled.
func (c *ContainerClient) executeQueryChunks(
	ctx context.Context,
	chunks []queryChunk,
	queryOpts *QueryOptions,
	operationContext pipelineRequestOptions,
	concurrency int,
) ([]chunkResult, error) {
	results := make([]chunkResult, len(chunks))
	jobs := make(chan int, len(chunks))
	var wg sync.WaitGroup

	// Create a child context that can be cancelled when any chunk encounters an error, to stop in-flight sibling requests.
	chunksCtx, cancelChunks := context.WithCancel(ctx)
	defer cancelChunks()

	workerCount := concurrency
	if workerCount > len(chunks) {
		workerCount = len(chunks)
	}
	if workerCount < 1 {
		workerCount = 1
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	if err != nil {
		// Unlikely since we are specifying the resource type, but handle just in case.
		return nil, errors.New("invalid resource address in operation context: " + operationContext.resourceAddress)
	}

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				if chunksCtx.Err() != nil {
					return
				}

				// Cancellation errors are stored in results and flow to collectChunkResults,
				// which reports the entire operation as failed.
				// Cancellation can come from EITHER the parent context, or the `cancelChunks` function cancelling our child context when any chunk encounters an error.
				items, charge, err := c.executeOneChunk(chunksCtx, chunks[idx], queryOpts, operationContext, path)
				results[idx] = chunkResult{items: items, requestCharge: charge, err: err}
				if err != nil {
					cancelChunks()
					return
				}
			}
		}()
	}

	// Feed jobs
	go func() {
		for i := range chunks {
			if chunksCtx.Err() != nil {
				break
			}
			jobs <- i
		}
		close(jobs)
	}()

	wg.Wait()
	return results, nil
}

// executeOneChunk sends a single parameterized query for a chunk, paging
// through continuation tokens until all results are collected.
func (c *ContainerClient) executeOneChunk(
	ctx context.Context,
	chunk queryChunk,
	queryOpts *QueryOptions,
	operationContext pipelineRequestOptions,
	path string,
) ([][]byte, float32, error) {
	var allItems [][]byte
	var totalCharge float32
	continuation := ""

	for {
		localOpts := *queryOpts
		if continuation != "" {
			localOpts.ContinuationToken = &continuation
		}

		rangeID := chunk.rangeID
		azResponse, err := c.database.client.sendQueryRequest(
			path,
			ctx,
			chunk.query,
			chunk.params,
			operationContext,
			&localOpts,
			func(req *policy.Request) {
				req.Raw().Header.Set(cosmosHeaderPartitionKeyRangeId, rangeID)
			},
		)
		if err != nil {
			return nil, totalCharge, err
		}

		qResp, err := newQueryResponse(azResponse)
		if err != nil {
			return nil, totalCharge, err
		}

		totalCharge += qResp.RequestCharge
		allItems = append(allItems, qResp.Items...)

		ct := azResponse.Header.Get(cosmosHeaderContinuationToken)
		if ct == "" {
			break
		}
		continuation = ct
	}

	return allItems, totalCharge, nil
}

// collectChunkResults merges per-chunk results into a single ReadManyItemsResponse.
// Returns the first error encountered, if any.
func collectChunkResults(results []chunkResult) (ReadManyItemsResponse, error) {
	var totalCharge float32
	var allItems [][]byte
	for _, res := range results {
		if res.err != nil {
			return ReadManyItemsResponse{}, res.err
		}
		totalCharge += res.requestCharge
		allItems = append(allItems, res.items...)
	}
	return ReadManyItemsResponse{RequestCharge: totalCharge, Items: allItems}, nil
}

// hasAnyPKRangeGoneError scans all chunk results for a partition key range gone error.
// This is needed because concurrent chunk cancellation can cause a context.Canceled
// error to appear at a lower index than the actual 410/Gone error, masking it from
// collectChunkResults which returns the first error it encounters.
func hasAnyPKRangeGoneError(results []chunkResult) bool {
	for _, res := range results {
		if res.err != nil && isPKRangeGoneResponseError(res.err) {
			return true
		}
	}
	return false
}

// executeReadManyWithQueries groups items by physical partition range using EPK
// hashing, builds parameterized SQL queries (one per physical range, chunked at
// maxItemsPerQuery), and executes them concurrently. This replaces the previous
// per-logical-PK strategy with fewer HTTP round-trips when multiple logical PKs
// map to the same physical partition range.
func (c *ContainerClient) executeReadManyWithQueries(
	ctx context.Context,
	items []ItemIdentity,
	readManyOptions *ReadManyOptions,
	operationContext pipelineRequestOptions,
) (ReadManyItemsResponse, error) {
	var pkDef PartitionKeyDefinition
	if c.database.client.containerCache != nil {
		containerProps, err := c.database.client.containerCache.getProperties(ctx, c)
		if err != nil {
			return ReadManyItemsResponse{}, err
		}
		pkDef = containerProps.PartitionKeyDefinition
	} else {
		// Fallback: direct fetch without caching
		containerResp, err := c.Read(ctx, nil)
		if err != nil {
			return ReadManyItemsResponse{}, err
		}
		pkDef = containerResp.ContainerProperties.PartitionKeyDefinition
	}

	concurrency := determineConcurrency(nil)
	if readManyOptions != nil {
		concurrency = determineConcurrency(readManyOptions.MaxConcurrency)
	}

	queryOpts := &QueryOptions{}
	if readManyOptions != nil {
		queryOpts.ConsistencyLevel = readManyOptions.ConsistencyLevel
		queryOpts.SessionToken = readManyOptions.SessionToken
		queryOpts.DedicatedGatewayRequestOptions = readManyOptions.DedicatedGatewayRequestOptions
	}

	// Retry loop for partition key range gone (splits/merges)
	for attempt := 0; attempt <= maxPKRangeGoneRetries; attempt++ {
		pkRangeResp, err := c.getPartitionKeyRanges(ctx, nil)
		if err != nil {
			return ReadManyItemsResponse{}, err
		}

		orderedRangeIDs, groups, err := groupItemsByPhysicalRange(items, pkDef, pkRangeResp.PartitionKeyRanges)
		if err != nil {
			return ReadManyItemsResponse{}, err
		}

		chunks, err := buildQueryChunksForRanges(orderedRangeIDs, groups, pkDef)
		if err != nil {
			return ReadManyItemsResponse{}, err
		}

		results, err := c.executeQueryChunks(ctx, chunks, queryOpts, operationContext, concurrency)
		if err != nil {
			if attempt < maxPKRangeGoneRetries && isPKRangeGoneResponseError(err) {
				if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
					return ReadManyItemsResponse{}, refreshErr
				}
				continue
			}
			return ReadManyItemsResponse{}, err
		}

		resp, err := collectChunkResults(results)
		if err != nil {
			// Check all results for 410/Gone, not just the first error returned by
			// collectChunkResults. Concurrent chunk cancellation can cause a
			// context.Canceled error at a lower index to mask the actual 410 error.
			if attempt < maxPKRangeGoneRetries && (isPKRangeGoneResponseError(err) || hasAnyPKRangeGoneError(results)) {
				if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
					return ReadManyItemsResponse{}, refreshErr
				}
				continue
			}
			return ReadManyItemsResponse{}, err
		}
		return resp, nil
	}

	return ReadManyItemsResponse{}, errors.New("exhausted retries for partition key range gone")
}

// refreshPKRangeCache forces a refresh of the partition key range cache for this container.
// Returns an error if the refresh fails, allowing the caller to fail fast.
func (c *ContainerClient) refreshPKRangeCache(ctx context.Context) error {
	if c.database.client.pkRangeCache != nil {
		containerRID, err := c.getContainerRID(ctx)
		if err != nil {
			return err
		}
		_, err = c.database.client.pkRangeCache.forceRefresh(ctx, containerRID, c.link, c.database.client)
		if err != nil {
			return err
		}
	}
	return nil
}

// isPKRangeGoneResponseError checks if an error is an azcore.ResponseError
// indicating a partition key range gone condition (HTTP 410 with split-related substatus).
func isPKRangeGoneResponseError(err error) bool {
	var respErr *azcore.ResponseError
	if !errors.As(err, &respErr) {
		return false
	}
	if respErr.StatusCode != http.StatusGone {
		return false
	}
	// Extract substatus from the raw response if available
	if respErr.RawResponse != nil {
		subStatus := respErr.RawResponse.Header.Get(cosmosHeaderSubstatus)
		return isPartitionKeyRangeGoneError(respErr.StatusCode, subStatus)
	}
	// If no raw response, any 410 Gone could be a PKRange gone
	return true
}

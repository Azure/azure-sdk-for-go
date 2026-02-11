// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// executeReadManyWithEngine executes a query using the provided query engine.
func (c *ContainerClient) executeReadManyWithEngine(queryEngine queryengine.QueryEngine, items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	// get the partition key ranges for the container
	rawPartitionKeyRanges, err := c.getPartitionKeyRangesRaw(ctx, operationContext)
	if err != nil {
		// if we can't get the partition key ranges, return empty response
		return ReadManyItemsResponse{}, err
	}

	// get the container properties
	containerRsp, err := c.Read(ctx, nil)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}

	// create the item identities for the query engine with json string
	newItemIdentities := make([]queryengine.ItemIdentity, len(items))
	for i := range items {
		pkStr, err := items[i].PartitionKey.toJsonString()
		if err != nil {
			return ReadManyItemsResponse{}, err
		}
		newItemIdentities[i] = queryengine.ItemIdentity{
			PartitionKeyValue: pkStr,
			ID:                items[i].ID,
		}
	}
	var pkVersion uint8
	pkDefinition := containerRsp.ContainerProperties.PartitionKeyDefinition
	if pkDefinition.Version == 0 {
		pkVersion = uint8(1)
	} else {
		pkVersion = uint8(pkDefinition.Version)
	}

	readManyPipeline, err := queryEngine.CreateReadManyPipeline(newItemIdentities, string(rawPartitionKeyRanges), string(pkDefinition.Kind), pkVersion, pkDefinition.Paths)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}
	log.Writef(EventQueryEngine, "Created readMany pipeline")
	// Initial run to get any requests.
	log.Writef(EventQueryEngine, "Fetching more data from readMany pipeline")
	result, err := readManyPipeline.Run()
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	concurrency := determineConcurrency(nil)
	if readManyOptions != nil {
		concurrency = determineConcurrency(readManyOptions.MaxConcurrency)
	}
	totalRequestCharge, err := runEngineRequests(ctx, c, path, readManyPipeline, operationContext, result.Requests, concurrency, func(req queryengine.QueryRequest) (string, []QueryParameter, bool) {
		// ReadMany pipeline requests carry a Query (optional override). No parameters and we always page until continuation exhausted.
		return req.Query, nil, true /* treat like drain for full pagination */
	})
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	// Final run to gather merged items.
	result, err = readManyPipeline.Run()
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	if readManyPipeline.IsComplete() {
		log.Writef(EventQueryEngine, "ReadMany pipeline is complete")
		readManyPipeline.Close()
		return ReadManyItemsResponse{
			Items:         result.Items,
			RequestCharge: totalRequestCharge,
		}, nil
	} else {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, errors.New("illegal state readMany pipeline did not complete")
	}
}

const maxItemsPerQuery = 1000

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
	idx := sort.Search(len(ranges), func(i int) bool {
		return ranges[i].MinInclusive > epkValue
	}) - 1
	if idx < 0 {
		return "", false
	}
	// Verify epkValue < MaxExclusive (empty MaxExclusive means unbounded).
	r := ranges[idx]
	if r.MaxExclusive != "" && epkValue >= r.MaxExclusive {
		return "", false
	}
	return r.ID, true
}

// groupItemsByPhysicalRange computes the EPK for each item and groups them by
// physical partition range. It returns the range IDs in first-seen order and
// the groups keyed by range ID.
func groupItemsByPhysicalRange(items []ItemIdentity, pkDef PartitionKeyDefinition, ranges []partitionKeyRange) ([]string, map[string][]ItemIdentity, error) {
	// Sort ranges by MinInclusive for binary search.
	sorted := make([]partitionKeyRange, len(ranges))
	copy(sorted, ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MinInclusive < sorted[j].MinInclusive
	})
	
	order := make([]string, 0)
	seen := make(map[string]bool)
	groups := make(map[string][]ItemIdentity)

	pkVersion := pkDef.Version
	if pkVersion == 0 {
		pkVersion = 1
	}

	for _, item := range items {
		epkVal := item.PartitionKey.computeEffectivePartitionKey(pkDef.Kind, pkVersion)
		rangeID, ok := findPhysicalRangeForEPK(epkVal.EPK, sorted)
		if !ok {
			return nil, nil, errors.New("could not find physical partition range for item EPK")
		}

		if !seen[rangeID] {
			order = append(order, rangeID)
			seen[rangeID] = true
		}
		groups[rangeID] = append(groups[rangeID], item)
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
// goroutine pool and returns the per-chunk results.
func (c *ContainerClient) executeQueryChunks(
	ctx context.Context,
	chunks []queryChunk,
	queryOpts *QueryOptions,
	operationContext pipelineRequestOptions,
	concurrency int,
) []chunkResult {
	results := make([]chunkResult, len(chunks))
	jobs := make(chan int, len(chunks))
	done := make(chan struct{})
	var wg sync.WaitGroup

	workerCount := concurrency
	if workerCount > len(chunks) {
		workerCount = len(chunks)
	}
	if workerCount < 1 {
		workerCount = 1
	}

	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				default:
				}

				items, charge, err := c.executeOneChunk(ctx, chunks[idx], queryOpts, operationContext, path, done)
				results[idx] = chunkResult{items: items, requestCharge: charge, err: err}
				if err != nil {
					select {
					case <-done:
					default:
						close(done)
					}
					return
				}
			}
		}()
	}

	// Feed jobs
	go func() {
		for i := range chunks {
			select {
			case <-done:
				close(jobs)
				return
			default:
			}
			jobs <- i
		}
		close(jobs)
	}()

	wg.Wait()
	return results
}

// executeOneChunk sends a single parameterized query for a chunk, paging
// through continuation tokens until all results are collected.
func (c *ContainerClient) executeOneChunk(
	ctx context.Context,
	chunk queryChunk,
	queryOpts *QueryOptions,
	operationContext pipelineRequestOptions,
	path string,
	done <-chan struct{},
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
				req.Raw().Header.Set(cosmosHeaderEnableCrossPartitionQuery, "True")
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
	containerResp, err := c.Read(ctx, nil)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}
	pkDef := containerResp.ContainerProperties.PartitionKeyDefinition

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

	results := c.executeQueryChunks(ctx, chunks, queryOpts, operationContext, concurrency)
	return collectChunkResults(results)
}

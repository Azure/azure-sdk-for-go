// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
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

// queryChunk is a single parameterized query targeting one logical partition key group.
type queryChunk struct {
	query  string
	params []QueryParameter
	pk     PartitionKey // used for x-ms-documentdb-partitionkey header routing
}

// chunkResult holds the outcome of executing a single queryChunk.
type chunkResult struct {
	items         [][]byte
	requestCharge float32
	err           error
}

// groupItemsByLogicalPK groups ItemIdentity values by their serialised partition
// key. It returns the groups in first-seen order.
func groupItemsByLogicalPK(items []ItemIdentity) ([]PartitionKey, map[string][]indexedItem, error) {
	order := make([]string, 0)
	pkForJSON := make(map[string]PartitionKey)
	groups := make(map[string][]indexedItem)

	for _, item := range items {
		pkJSON, err := item.PartitionKey.toJsonString()
		if err != nil {
			return nil, nil, err
		}
		if _, exists := groups[pkJSON]; !exists {
			order = append(order, pkJSON)
			pkForJSON[pkJSON] = item.PartitionKey
		}
		groups[pkJSON] = append(groups[pkJSON], indexedItem{
			id: item.ID,
			pk: item.PartitionKey,
		})
	}

	// Build ordered PK list
	pks := make([]PartitionKey, len(order))
	for i, j := range order {
		pks[i] = pkForJSON[j]
	}
	return pks, groups, nil
}

// buildQueryChunks creates queryChunk values by splitting each logical PK group
// into slices of at most maxItemsPerQuery and building parameterized SQL for each.
func buildQueryChunks(orderedPKs []PartitionKey, groups map[string][]indexedItem, pkDef PartitionKeyDefinition) ([]queryChunk, error) {
	qb := queryBuilder{}
	var chunks []queryChunk

	for _, pk := range orderedPKs {
		pkJSON, err := pk.toJsonString()
		if err != nil {
			return nil, err
		}
		items := groups[pkJSON]
		for start := 0; start < len(items); start += maxItemsPerQuery {
			end := start + maxItemsPerQuery
			if end > len(items) {
				end = len(items)
			}
			q, params := qb.buildParameterizedQueryForItems(items[start:end], pkDef)
			chunks = append(chunks, queryChunk{
				query:  q,
				params: params,
				pk:     pk,
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
	pkHeader, err := chunk.pk.toJsonString()
	if err != nil {
		return nil, 0, err
	}

	var allItems [][]byte
	var totalCharge float32
	continuation := ""

	for {
		localOpts := *queryOpts
		if continuation != "" {
			localOpts.ContinuationToken = &continuation
		}

		pkHeaderCopy := pkHeader
		azResponse, err := c.database.client.sendQueryRequest(
			path,
			ctx,
			chunk.query,
			chunk.params,
			operationContext,
			&localOpts,
			func(req *policy.Request) {
				req.Raw().Header.Set(cosmosHeaderPartitionKey, pkHeaderCopy)
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

// executeReadManyWithQueries groups items by logical partition key, builds
// parameterized SQL queries (one per PK group, chunked at maxItemsPerQuery),
// and executes them concurrently. This replaces the previous per-item point-read
// strategy with far fewer HTTP round-trips.
//
// V1 routes each query using the x-ms-documentdb-partitionkey header (one query
// per logical PK group). A future V2 can add EPK hashing to coalesce groups that
// map to the same physical partition range into a single OR-of-conjunctions query.
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

	orderedPKs, groups, err := groupItemsByLogicalPK(items)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}

	chunks, err := buildQueryChunks(orderedPKs, groups, pkDef)
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

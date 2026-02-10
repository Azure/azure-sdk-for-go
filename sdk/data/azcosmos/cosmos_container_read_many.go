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
	// 1. Fetch container properties → PartitionKeyDefinition
	containerResp, err := c.Read(ctx, nil)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}
	pkDef := containerResp.ContainerProperties.PartitionKeyDefinition

	// 2. Group items by logical partition key value.
	type pkGroup struct {
		pk    PartitionKey
		pkJSON string
		items []indexedItem
	}
	groupOrder := make([]string, 0)            // preserves first-seen order
	groupMap := make(map[string]*pkGroup)       // pkJSON → group

	for _, item := range items {
		pkJSON, jsonErr := item.PartitionKey.toJsonString()
		if jsonErr != nil {
			return ReadManyItemsResponse{}, jsonErr
		}
		g, exists := groupMap[pkJSON]
		if !exists {
			g = &pkGroup{pk: item.PartitionKey, pkJSON: pkJSON}
			groupMap[pkJSON] = g
			groupOrder = append(groupOrder, pkJSON)
		}
		g.items = append(g.items, indexedItem{
			id: item.ID,
			pk: item.PartitionKey,
		})
	}

	// 3. Build chunks (≤ maxItemsPerQuery) and corresponding queries.
	type queryChunk struct {
		query  string
		params []QueryParameter
		pk     PartitionKey // used for x-ms-documentdb-partitionkey header
	}

	qb := queryBuilder{}
	var chunks []queryChunk

	for _, pkJSON := range groupOrder {
		g := groupMap[pkJSON]
		for start := 0; start < len(g.items); start += maxItemsPerQuery {
			end := start + maxItemsPerQuery
			if end > len(g.items) {
				end = len(g.items)
			}
			slice := g.items[start:end]

			q, params := qb.buildParameterizedQueryForItems(slice, pkDef)
			chunks = append(chunks, queryChunk{
				query:  q,
				params: params,
				pk:     g.pk,
			})
		}
	}

	// 4. Execute chunks concurrently.
	concurrency := determineConcurrency(nil)
	if readManyOptions != nil {
		concurrency = determineConcurrency(readManyOptions.MaxConcurrency)
	}

	type chunkResult struct {
		items         [][]byte
		requestCharge float32
		err           error
	}

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

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}

	// Build query options from ReadManyOptions
	queryOpts := &QueryOptions{}
	if readManyOptions != nil {
		queryOpts.ConsistencyLevel = readManyOptions.ConsistencyLevel
		queryOpts.SessionToken = readManyOptions.SessionToken
		queryOpts.DedicatedGatewayRequestOptions = readManyOptions.DedicatedGatewayRequestOptions
	}

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

				chunk := chunks[idx]

				// Use the partition key header for routing
				pkHeader, jsonErr := chunk.pk.toJsonString()
				if jsonErr != nil {
					results[idx].err = jsonErr
					select {
					case <-done:
					default:
						close(done)
					}
					return
				}

				// Paginate through continuation tokens
				var allItems [][]byte
				var totalCharge float32
				continuation := ""

				for {
					localOpts := *queryOpts
					if continuation != "" {
						localOpts.ContinuationToken = &continuation
					}

					pkHeaderCopy := pkHeader
					azResponse, qErr := c.database.client.sendQueryRequest(
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
					if qErr != nil {
						results[idx].err = qErr
						select {
						case <-done:
						default:
							close(done)
						}
						return
					}

					qResp, qErr := newQueryResponse(azResponse)
					if qErr != nil {
						results[idx].err = qErr
						select {
						case <-done:
						default:
							close(done)
						}
						return
					}

					totalCharge += qResp.RequestCharge
					allItems = append(allItems, qResp.Items...)

					ct := azResponse.Header.Get(cosmosHeaderContinuationToken)
					if ct == "" {
						break
					}
					continuation = ct
				}

				results[idx].items = allItems
				results[idx].requestCharge = totalCharge
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

	// 5. Collect results — no ordering guarantee.
	var totalRequestCharge float32
	var allItems [][]byte
	for _, res := range results {
		if res.err != nil {
			return ReadManyItemsResponse{}, res.err
		}
		totalRequestCharge += res.requestCharge
		allItems = append(allItems, res.items...)
	}

	return ReadManyItemsResponse{
		RequestCharge: totalRequestCharge,
		Items:         allItems,
	}, nil
}

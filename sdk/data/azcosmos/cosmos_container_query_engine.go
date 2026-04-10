// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore Writef

package azcosmos

import (
	"bytes"
	"context"
	"runtime"
	"sync"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// EventQueryEngine contains logs related to the query engine.
const EventQueryEngine log.Event = "QueryEngine"

func (c *ContainerClient) getQueryPlanFromGateway(ctx context.Context, query string, supportedFeatures string, queryOptions *QueryOptions, operationContext pipelineRequestOptions) ([]byte, error) {
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	azResponse, err := c.database.client.sendQueryRequest(
		path,
		ctx,
		query,
		queryOptions.QueryParameters,
		operationContext,
		queryOptions,
		func(req *policy.Request) {
			req.Raw().Header.Set(cosmosHeaderIsQueryPlanRequest, "True")
			req.Raw().Header.Set(cosmosHeaderSupportedQueryFeatures, supportedFeatures)
		})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(azResponse.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *ContainerClient) getPartitionKeyRangesRaw(ctx context.Context, operationContext pipelineRequestOptions) ([]byte, error) {
	path, _ := generatePathForNameBased(resourceTypePartitionKeyRange, operationContext.resourceAddress, true)
	azResponse, err := c.database.client.sendGetRequest(
		path,
		ctx,
		pipelineRequestOptions{
			resourceType:          resourceTypePartitionKeyRange,
			resourceAddress:       operationContext.resourceAddress,
			headerOptionsOverride: operationContext.headerOptionsOverride,
		},
		nil,
		nil)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(azResponse.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Executes a query using the provided query engine.
func (c *ContainerClient) executeQueryWithEngine(queryEngine queryengine.QueryEngine, query string, queryOptions *QueryOptions, operationContext pipelineRequestOptions) *azruntime.Pager[QueryItemsResponse] {
	// NOTE: The current interface for runtime.Pager means we're probably going to risk leaking the pipeline, if it's provided by a native query engine.
	// There's no "Close" method, which means we can't call `queryengine.QueryPipeline.Close()` when we're done.
	// We _do_ close the pipeline if the user iterates the entire pager, but if they don't we don't have a way to clean up.
	// To mitigate that, we expect the queryengine.QueryPipeline to handle setting up a Go finalizer to clean up any native resources it holds.
	// Finalizers aren't deterministic though, so we should consider making the pager "closable" in the future, so we have a clear signal to free the native resources.

	var queryPipeline queryengine.QueryPipeline
	var lastResponse Response
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	log.Writef(EventQueryEngine, "Executing query using query engine")
	return azruntime.NewPager(azruntime.PagingHandler[QueryItemsResponse]{
		More: func(page QueryItemsResponse) bool {
			if queryPipeline == nil {
				// We haven't started yet, so there's certainly more to do.
				return true
			}

			if queryPipeline.IsComplete() {
				// If it's not already closed, close the pipeline.
				// Close is expected to be idempotent, so we can call it multiple times.
				queryPipeline.Close()
				return false
			}

			// The pipeline isn't complete, so we can keep going.
			return true
		},
		Fetcher: func(ctx context.Context, page *QueryItemsResponse) (QueryItemsResponse, error) {
			if queryPipeline == nil {
				// First page, we need to fetch the query plan and PK ranges
				// TODO: We could proactively try to run this query against the gateway and then fall back to the engine. That's what Python does.
				plan, err := c.getQueryPlanFromGateway(ctx, query, queryEngine.SupportedFeatures(), queryOptions, operationContext)
				if err != nil {
					return QueryItemsResponse{}, err
				}
				pkranges, err := c.getPartitionKeyRangesRaw(ctx, operationContext)
				if err != nil {
					return QueryItemsResponse{}, err
				}

				// Create a query pipeline
				queryPipeline, err = queryEngine.CreateQueryPipeline(query, string(plan), string(pkranges))
				if err != nil {
					return QueryItemsResponse{}, err
				}
				log.Writef(EventQueryEngine, "Created query pipeline")

				// The gateway may have rewritten the query, which would be encoded in the query plan.
				// The pipeline parsed the query plan, so we can ask it for the rewritten query.
				query = queryPipeline.Query()
			}

			for {
				if queryPipeline.IsComplete() {
					log.Writef(EventQueryEngine, "Query pipeline is complete")
					queryPipeline.Close()
					return QueryItemsResponse{
						Response: lastResponse,
						Items:    nil,
					}, nil
				}
				// Fetch more data from the pipeline
				log.Writef(EventQueryEngine, "Fetching more data from query pipeline")
				result, err := queryPipeline.Run()
				if err != nil {
					queryPipeline.Close()
					return QueryItemsResponse{}, err
				}

				// If we got items, we can return them, and we should do so now, to avoid making unnecessary requests.
				// Even if there are requests in the queue, the pipeline should return the same requests again on the next call to NextBatch.
				if len(result.Items) > 0 {
					log.Writef(EventQueryEngine, "Query pipeline returned %d items", len(result.Items))
					return QueryItemsResponse{
						Response: lastResponse,
						Items:    result.Items,
					}, nil
				}

				// If we didn't have any items to return, we need to make requests for the items in the queue.
				// If there are no requests, the pipeline should return true for IsComplete, so we'll stop on the next iteration.
				// Parallelize request execution using shared driver.
				concurrency := determineConcurrency(nil)
				charge, err := runEngineRequests(ctx, c, path, queryPipeline, operationContext, result.Requests, concurrency, func(qryRequest queryengine.QueryRequest) (string, []QueryParameter, bool) {
					// Override query if present;
					localQuery := query
					if qryRequest.Query != "" {
						localQuery = qryRequest.Query
					}
					var queryParameters []QueryParameter
					if qryRequest.IncludeParameters || qryRequest.Query == "" {
						// use query options parameters only if IncludeParameters is true or no override query is specified
						queryParameters = queryOptions.QueryParameters
					}
					// Drain if request.Drain is true.
					return localQuery, queryParameters, qryRequest.Drain
				})
				_ = charge // totalRequestCharge currently unused for query path;
				if err != nil {
					queryPipeline.Close()
					return QueryItemsResponse{}, err
				}
				// Loop again to attempt to produce items.
			}
		},
	})
}

// runEngineRequests concurrently executes per-partition QueryRequests for either query or readMany pipelines.
// prepareFn returns the query text, parameters, and a drain flag for each request.
// Collects all results and calls ProvideData once with a single batch to reduce CGo overhead.
func runEngineRequests(
	ctx context.Context,
	c *ContainerClient,
	path string,
	pipeline queryengine.QueryPipeline,
	operationContext pipelineRequestOptions,
	requests []queryengine.QueryRequest,
	concurrency int,
	prepareFn func(req queryengine.QueryRequest) (query string, params []QueryParameter, drain bool),
) (float32, error) {
	if len(requests) == 0 {
		return 0, nil
	}

	jobs := make(chan queryengine.QueryRequest, len(requests))
	errCh := make(chan error, 1)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Adjust concurrency.
	workerCount := concurrency
	if workerCount > len(requests) {
		workerCount = len(requests)
	}
	if workerCount < 1 {
		workerCount = 1
	}

	// Per-worker request charge slots and result slices (lock-free updates)
	charges := make([]float32, workerCount)
	resultsSlices := make([][]queryengine.QueryResult, workerCount)

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			localResults := make([]queryengine.QueryResult, 0, 8)
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case req, ok := <-jobs:
					if !ok {
						// jobs exhausted
						resultsSlices[workerIndex] = localResults
						return
					}
					log.Writef(azlog.EventRequest, "Engine pipeline requested data for PKRange: %s", req.PartitionKeyRangeID)
					queryText, params, drain := prepareFn(req)
					fetchMorePages := true
					for fetchMorePages {
						qr := queryRequest(req)
						azResponse, err := c.database.client.sendQueryRequest(
							path,
							ctx,
							queryText,
							params,
							operationContext,
							&qr,
							nil,
						)
						if err != nil {
							select {
							case errCh <- err:
							default:
							}
							return
						}
						qResp, err := newQueryResponse(azResponse)
						if err != nil {
							select {
							case errCh <- err:
							default:
							}
							return
						}
						charges[workerIndex] += qResp.RequestCharge
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(azResponse.Body); err != nil {
							select {
							case errCh <- err:
							default:
							}
							return
						}
						continuation := azResponse.Header.Get(cosmosHeaderContinuationToken)
						data := buf.Bytes()
						fetchMorePages = continuation != "" && drain
						localResults = append(localResults, queryengine.QueryResult{
							PartitionKeyRangeID: req.PartitionKeyRangeID,
							NextContinuation:    continuation,
							Data:                data,
							RequestId:           req.Id,
						})
						log.Writef(EventQueryEngine, "Received response for PKRange: %s. Continuation present: %v", req.PartitionKeyRangeID, continuation != "")
					}
				}
			}
		}(w)
	}

	// Feed jobs
	go func() {
		for _, r := range requests {
			select {
			case <-done:
				return
			default:
			}
			jobs <- r
		}
		close(jobs)
	}()

	// Wait for workers to finish (or error/cancel)
	workersDone := make(chan struct{})
	go func() { wg.Wait(); close(workersDone) }()

	// Helper to sum charges
	sumCharges := func() float32 {
		var total float32
		for _, cval := range charges {
			total += cval
		}
		return total
	}

	// Wait for completion / error / cancellation
	select {
	case e := <-errCh:
		select {
		case <-done:
		default:
			close(done)
		}
		return sumCharges(), e
	case <-ctx.Done():
		select {
		case <-done:
		default:
			close(done)
		}
		return sumCharges(), ctx.Err()
	case <-workersDone:
	}

	totalCharge := sumCharges()

	// Merge per-worker result slices deterministically
	// Pre-size combined slice for efficiency
	var combinedCount int
	for _, rs := range resultsSlices {
		combinedCount += len(rs)
	}
	if combinedCount > 0 {
		all := make([]queryengine.QueryResult, 0, combinedCount)
		for _, rs := range resultsSlices {
			all = append(all, rs...)
		}
		if err := pipeline.ProvideData(all); err != nil {
			return totalCharge, err
		}
	}

	return totalCharge, nil
}

// determineConcurrency returns either the provided positive max or NumCPU (>=1).
func determineConcurrency(max *int32) int {
	if max != nil && *max > 0 {
		return int(*max)
	}
	c := runtime.NumCPU()
	if c <= 0 {
		c = 1
	}
	return c
}

// Wrapper type because we can't define 'toHeaders' on DataRequest directly, nor can we define it in the queryengine package (because it's not a public method).
type queryRequest queryengine.QueryRequest

func (r *queryRequest) toHeaders() *map[string]string {
	headers := make(map[string]string)

	if r.Continuation != "" {
		headers[cosmosHeaderContinuationToken] = r.Continuation
	}
	headers[cosmosHeaderPartitionKeyRangeId] = r.PartitionKeyRangeID
	return &headers
}

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

func (c *ContainerClient) executeQueryWithEngine(queryEngine queryengine.QueryEngine, query string, queryOptions *QueryOptions, operationContext pipelineRequestOptions) *azruntime.Pager[QueryItemsResponse] {
	// NOTE: The current interface for runtime.Pager means we're probably going to risk leaking the pipeline, if it's provided by a native query engine.
	// There's no "Close" method, which means we can't call `queryengine.QueryPipeline.Close()` when we're done.
	// We _do_ close the pipeline if the user iterates the entire pager, but if they don't we don't have a way to clean up.
	// To mitigate that, we expect the queryengine.QueryPipeline to handle setting up a Go finalizer to clean up any native resources it holds.
	// Finalizers aren't deterministic though, so we should consider making the pager "closable" in the future, so we have a clear signal to free the native resources.

	var queryPipeline queryengine.QueryPipeline
	var lastResponse Response
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
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
					return QueryItemsResponse{Response: lastResponse, Items: nil}, nil
				}

				log.Writef(EventQueryEngine, "Fetching more data from query pipeline")
				result, err := queryPipeline.Run()
				if err != nil {
					queryPipeline.Close()
					return QueryItemsResponse{}, err
				}

				if len(result.Items) > 0 {
					log.Writef(EventQueryEngine, "Query pipeline returned %d items", len(result.Items))
					return QueryItemsResponse{Response: lastResponse, Items: result.Items}, nil
				}

				// Parallelize request execution using shared driver.
				concurrency := determineConcurrency(nil)
				charge, err := runEngineRequests(ctx, c, path, queryPipeline, operationContext, result.Requests, concurrency, func(req queryengine.QueryRequest) (string, []QueryParameter, bool) {
					// Override query if present; decide parameters usage same as previous logic.
					localQuery := query
					if req.Query != "" {
						localQuery = req.Query
					}
					var params []QueryParameter
					if req.IncludeParameters || req.Query == "" {
						params = queryOptions.QueryParameters
					}
					// Drain if request.Drain is true.
					return localQuery, params, req.Drain
				})
				_ = charge // totalRequestCharge currently unused for query path; could accumulate into lastResponse later.
				if err != nil {
					queryPipeline.Close()
					return QueryItemsResponse{}, err
				}
				// Loop again to attempt to produce items.
			}
		},
	})
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

// runEngineRequests concurrently executes per-partition QueryRequests for either query or readMany pipelines.
// prepareFn returns the query text, parameters, and a drain flag for each request.
// It serializes ProvideData calls through a single goroutine to preserve ordering guarantees required by the pipeline.
func runEngineRequests(
	ctx context.Context,
	c *ContainerClient,
	path string,
	pipeline queryengine.QueryPipeline,
	operationContext pipelineRequestOptions,
	requests []queryengine.QueryRequest,
	concurrency int,
	prepareFn func(req queryengine.QueryRequest) (query string, params []QueryParameter, drain bool),
) (totalCharge float32, err error) {
	if len(requests) == 0 {
		return 0, nil
	}

	jobs := make(chan queryengine.QueryRequest, len(requests))
	provideCh := make(chan []queryengine.QueryResult)
	errCh := make(chan error, 1)
	done := make(chan struct{})
	providerDone := make(chan struct{})
	var wg sync.WaitGroup
	var chargeMu sync.Mutex

	// Provider goroutine ensures only one ProvideData executes at a time.
	go func() {
		defer close(providerDone)
		for batch := range provideCh {
			if perr := pipeline.ProvideData(batch); perr != nil {
				select {
				case errCh <- perr:
				default:
				}
				return
			}
		}
	}()

	// Adjust concurrency.
	workerCount := concurrency
	if workerCount > len(requests) {
		workerCount = len(requests)
	}
	if workerCount < 1 {
		workerCount = 1
	}

	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case req, ok := <-jobs:
					if !ok {
						return
					}

					log.Writef(azlog.EventRequest, "Engine pipeline requested data for PKRange: %s", req.PartitionKeyRangeID)
					queryText, params, drain := prepareFn(req)
					// Pagination loop
					fetchMore := true
					for fetchMore {
						qr := queryRequest(req)
						azResponse, rerr := c.database.client.sendQueryRequest(
							path,
							ctx,
							queryText,
							params,
							operationContext,
							&qr,
							nil,
						)
						if rerr != nil {
							select {
							case errCh <- rerr:
							default:
							}
							return
						}

						qResp, qErr := newQueryResponse(azResponse)
						if qErr != nil {
							select {
							case errCh <- qErr:
							default:
							}
							return
						}
						chargeMu.Lock()
						totalCharge += qResp.RequestCharge
						chargeMu.Unlock()

						buf := new(bytes.Buffer)
						if _, rdErr := buf.ReadFrom(azResponse.Body); rdErr != nil {
							select {
							case errCh <- rdErr:
							default:
							}
							return
						}
						continuation := azResponse.Header.Get(cosmosHeaderContinuationToken)
						fetchMore = continuation != "" && drain

						qres := queryengine.QueryResult{
							PartitionKeyRangeID: req.PartitionKeyRangeID,
							NextContinuation:    continuation,
							RequestId:           req.Id,
							Data:                buf.Bytes(),
						}
						select {
						case <-done:
							return
						case provideCh <- []queryengine.QueryResult{qres}:
						}
					}
				}
			}
		}()
	}

	// Feed jobs
	go func() {
		for _, r := range requests {
			select {
			case <-done:
				break
			default:
			}
			jobs <- r
		}
		close(jobs)
	}()

	// Close provider after workers finish
	go func() { wg.Wait(); close(provideCh) }()

	// Wait for completion / error / cancellation
	select {
	case e := <-errCh:
		select {
		case <-done:
		default:
			close(done)
		}
		return totalCharge, e
	case <-ctx.Done():
		select {
		case <-done:
		default:
			close(done)
		}
		return totalCharge, ctx.Err()
	case <-providerDone:
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

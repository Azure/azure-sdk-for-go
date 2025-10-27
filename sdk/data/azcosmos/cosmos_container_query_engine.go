// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore Writef

package azcosmos

import (
	"bytes"
	"context"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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
func (c *ContainerClient) executeQueryWithEngine(queryEngine queryengine.QueryEngine, query string, queryOptions *QueryOptions, operationContext pipelineRequestOptions) *runtime.Pager[QueryItemsResponse] {
	// NOTE: The current interface for runtime.Pager means we're probably going to risk leaking the pipeline, if it's provided by a native query engine.
	// There's no "Close" method, which means we can't call `queryengine.QueryPipeline.Close()` when we're done.
	// We _do_ close the pipeline if the user iterates the entire pager, but if they don't we don't have a way to clean up.
	// To mitigate that, we expect the queryengine.QueryPipeline to handle setting up a Go finalizer to clean up any native resources it holds.
	// Finalizers aren't deterministic though, so we should consider making the pager "closable" in the future, so we have a clear signal to free the native resources.

	var queryPipeline queryengine.QueryPipeline
	var lastResponse Response
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	log.Writef(EventQueryEngine, "Executing query using query engine")
	return runtime.NewPager(runtime.PagingHandler[QueryItemsResponse]{
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
				// TODO: We can absolutely parallelize these requests
				for _, request := range result.Requests {
					log.Writef(azlog.EventRequest, "Query pipeline requested data for PKRange: %s", request.PartitionKeyRangeID)
					// Make the single-partition query request
					qryRequest := queryRequest(request) // Cast to our type, which has toHeaders defined on it.
					// if the query request has an override query, use it
					if qryRequest.Query != "" {
						query = qryRequest.Query
					}

					fetchMorePages := true
					for fetchMorePages {

						azResponse, err := c.database.client.sendQueryRequest(
							path,
							ctx,
							query,
							queryOptions.QueryParameters,
							operationContext,
							&qryRequest,
							nil)
						if err != nil {
							queryPipeline.Close()
							return QueryItemsResponse{}, err
						}
						lastResponse = newResponse(azResponse)

						// Load the data into a buffer to send it to the pipeline
						buf := new(bytes.Buffer)
						_, err = buf.ReadFrom(azResponse.Body)
						if err != nil {
							queryPipeline.Close()
							return QueryItemsResponse{}, err
						}
						data := buf.Bytes()
						continuation := azResponse.Header.Get(cosmosHeaderContinuationToken)

						fetchMorePages = continuation != "" && request.Drain

						// Provide the data to the pipeline, make sure it's tagged with the partition key range ID so the pipeline can merge it into the correct partition.
						result := queryengine.QueryResult{
							PartitionKeyRangeID: request.PartitionKeyRangeID,
							RequestIndex:        request.Index,
							NextContinuation:    continuation,
							Data:                data,
						}
						log.Writef(EventQueryEngine, "Received response for PKRange: %s. Continuation present: %v", request.PartitionKeyRangeID, continuation != "")
						if err = queryPipeline.ProvideData(result); err != nil {
							queryPipeline.Close()
							return QueryItemsResponse{}, err
						}
					}
				}

				// No items, but we provided more data, so let's continue the loop.
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

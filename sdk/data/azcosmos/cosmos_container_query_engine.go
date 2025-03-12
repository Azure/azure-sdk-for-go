package azcosmos

import (
	"bytes"
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/unstable/queryengine"
)

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

func (c *ContainerClient) getPartitionKeyRanges(ctx context.Context, operationContext pipelineRequestOptions) ([]byte, error) {
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
	// We _do_ call that if the user iterates the entire pager, but if they don't we don't have a way to clean up.
	// To mitigate that, we expect the queryengine.QueryPipeline to handle setting up a finalizer to clean up any native resources it holds.

	var queryPipeline queryengine.QueryPipeline
	var lastResponse Response
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	return runtime.NewPager(runtime.PagingHandler[QueryItemsResponse]{
		More: func(page QueryItemsResponse) bool {
			return queryPipeline == nil || !queryPipeline.IsComplete()
		},
		Fetcher: func(ctx context.Context, page *QueryItemsResponse) (QueryItemsResponse, error) {
			if queryPipeline == nil {
				// First page, we need to fetch the query plan and PK ranges
				// TODO: We could proactively try to run this against the gateway and then fall back to the engine.
				plan, err := c.getQueryPlanFromGateway(ctx, query, queryEngine.SupportedFeatures(), queryOptions, operationContext)
				if err != nil {
					return QueryItemsResponse{}, err
				}
				pkranges, err := c.getPartitionKeyRanges(ctx, operationContext)
				if err != nil {
					return QueryItemsResponse{}, err
				}

				// Create a query pipeline
				queryPipeline, err = queryEngine.CreateQueryPipeline(query, string(plan), string(pkranges))
				if err != nil {
					return QueryItemsResponse{}, err
				}

				query = queryPipeline.Query()
			}

			for {
				if queryPipeline.IsComplete() {
					queryPipeline.Close()
					return QueryItemsResponse{
						Response: lastResponse,
						Items:    nil,
					}, nil
				}
				// Fetch more data from the pipeline
				items, requests, err := queryPipeline.NextBatch()
				if err != nil {
					queryPipeline.Close()
					return QueryItemsResponse{}, err
				}

				if len(items) > 0 {
					return QueryItemsResponse{
						Response: lastResponse,
						Items:    items,
					}, nil
				}

				// We didn't get any items, so we need to put more data in.
				for _, request := range requests {
					qryRequest := queryRequest(request) // Cast to our type, which has toHeaders defined on it.
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

					buf := new(bytes.Buffer)
					_, err = buf.ReadFrom(azResponse.Body)
					if err != nil {
						queryPipeline.Close()
						return QueryItemsResponse{}, err
					}
					data := buf.Bytes()
					continuation := azResponse.Header.Get(cosmosHeaderContinuationToken)

					result := queryengine.QueryResult{
						PartitionKeyRangeID: request.PartitionKeyRangeID,
						NextContinuation:    continuation,
						Data:                data,
					}
					if err = queryPipeline.ProvideData(result); err != nil {
						queryPipeline.Close()
						return QueryItemsResponse{}, err
					}
				}

				// We've provided more data to the pipeline, so let's try to run the pipeline forward again.
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

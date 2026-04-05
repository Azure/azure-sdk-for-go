// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"strings"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const EventQueryOptimization log.Event = "QueryOptimization"

type queryExecutionMode int

const (
	queryModeGateway queryExecutionMode = iota
	queryModeODE
	queryModeEngine
)

func (c *ContainerClient) selectQueryExecutionMode(partitionKey PartitionKey, queryOptions *QueryOptions) queryExecutionMode {
	if queryOptions.QueryEngine != nil {
		return queryModeEngine
	}

	if c.database.client.directTransport == nil {
		return queryModeGateway
	}

	odeEnabled := queryOptions.EnableOptimisticDirectExecution == nil || *queryOptions.EnableOptimisticDirectExecution

	if !odeEnabled {
		return queryModeGateway
	}

	if len(partitionKey.values) == 0 {
		return queryModeGateway
	}

	return queryModeODE
}

func (c *ContainerClient) executeQueryWithODE(query string, partitionKey PartitionKey, queryOptions *QueryOptions, operationContext pipelineRequestOptions) *azruntime.Pager[QueryItemsResponse] {
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	log.Writef(EventQueryOptimization, "Executing query with ODE (Optimistic Direct Execution)")

	return azruntime.NewPager(azruntime.PagingHandler[QueryItemsResponse]{
		More: func(page QueryItemsResponse) bool {
			return page.ContinuationToken != nil
		},
		Fetcher: func(ctx context.Context, page *QueryItemsResponse) (QueryItemsResponse, error) {
			var err error
			spanName, err := c.getSpanForItems(operationTypeQuery)
			if err != nil {
				return QueryItemsResponse{}, err
			}
			ctx, endSpan := azruntime.StartSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
			defer func() { endSpan(err) }()

			if page != nil {
				if page.ContinuationToken != nil {
					queryOptions.ContinuationToken = page.ContinuationToken
				}
			}

			h := operationContext.headerOptionsOverride
			if h == nil {
				h = &headerOptionsOverride{partitionKey: &partitionKey}
			}
			c.applyDirectModeRouting(ctx, &partitionKey, h)
			operationContext.headerOptionsOverride = h

			azResponse, err := c.database.client.sendQueryRequest(
				path,
				ctx,
				query,
				queryOptions.QueryParameters,
				operationContext,
				queryOptions,
				nil)

			if err != nil {
				return QueryItemsResponse{}, err
			}

			return newQueryResponse(azResponse)
		},
	})
}

func (c *ContainerClient) getQueryPlanWithCache(ctx context.Context, query string, supportedFeatures string, queryOptions *QueryOptions, operationContext pipelineRequestOptions) ([]byte, error) {
	cache := c.database.client.queryPlanCache
	if cache == nil || queryOptions.DisableQueryPlanCache {
		return c.getQueryPlanFromGateway(ctx, query, supportedFeatures, queryOptions, operationContext)
	}

	if cached, ok := cache.Get(c.link, query); ok {
		log.Writef(EventQueryOptimization, "Query plan cache hit")
		return cached, nil
	}

	log.Writef(EventQueryOptimization, "Query plan cache miss, fetching from gateway")
	plan, err := c.getQueryPlanFromGateway(ctx, query, supportedFeatures, queryOptions, operationContext)
	if err != nil {
		return nil, err
	}

	cache.Set(c.link, query, plan)
	return plan, nil
}

func isSimpleQuery(query string) bool {
	upper := strings.ToUpper(query)

	hasOrderBy := strings.Contains(upper, "ORDER BY")
	hasGroupBy := strings.Contains(upper, "GROUP BY")
	hasDistinct := strings.Contains(upper, "DISTINCT")
	hasTop := strings.Contains(upper, "TOP")
	hasOffset := strings.Contains(upper, "OFFSET")
	hasAggregate := strings.Contains(upper, "COUNT(") ||
		strings.Contains(upper, "SUM(") ||
		strings.Contains(upper, "AVG(") ||
		strings.Contains(upper, "MIN(") ||
		strings.Contains(upper, "MAX(")

	return !hasOrderBy && !hasGroupBy && !hasDistinct && !hasTop && !hasOffset && !hasAggregate
}

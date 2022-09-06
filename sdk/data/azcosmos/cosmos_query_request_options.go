// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "fmt"

// QueryOptions includes options for query operations on items.
type QueryOptions struct {
	// SessionToken to be used when using Session consistency on the account.
	// When working with Session consistency, each new write request to Azure Cosmos DB is assigned a new SessionToken.
	// The client instance will use this token internally with each read/query request to ensure that the set consistency level is maintained.
	// In some scenarios you need to manage this Session yourself: Consider a web application with multiple nodes, each node will have its own client instance.
	// If you wanted these nodes to participate in the same session (to be able read your own writes consistently across web tiers),
	// you would have to send the SessionToken from the response of the write action on one node to the client tier, using a cookie or some other mechanism, and have that token flow back to the web tier for subsequent reads.
	// If you are using a round-robin load balancer which does not maintain session affinity between requests, such as the Azure Load Balancer,the read could potentially land on a different node to the write request, where the session was created.
	SessionToken string
	// ConsistencyLevel overrides the account defined consistency level for this operation.
	// Consistency can only be relaxed.
	ConsistencyLevel *ConsistencyLevel
	// PopulateIndexMetrics is used to obtain the index metrics to understand how the query engine used existing indexes and how it could use potential new indexes.
	// Please note that this options will incur overhead, so it should be enabled only when debugging slow queries and not in production.
	PopulateIndexMetrics bool
	// ResponseContinuationTokenLimitInKB is used to limit the length of continuation token in the query response. Valid values are >= 0.
	ResponseContinuationTokenLimitInKB int32
	// PageSizeHint determines the maximum number of items to be retrieved in a query result page.
	// '-1' Used for dynamic page size. This is a maximum. Query can return 0 items in the page.
	PageSizeHint int32
	// EnableScanInQuery Allow scan on the queries which couldn't be served as indexing was opted out on the requested paths.
	EnableScanInQuery bool
	// ContinuationToken to be used to continue a previous query execution.
	// Obtained from QueryItemsResponse.ContinuationToken.
	ContinuationToken string
	// QueryParameters allows execution of parametrized queries.
	// See https://docs.microsoft.com/azure/cosmos-db/sql/sql-query-parameterized-queries
	QueryParameters []QueryParameter
}

func (options *QueryOptions) toHeaders() *map[string]string {
	headers := make(map[string]string)

	if options.ConsistencyLevel != nil {
		headers[cosmosHeaderConsistencyLevel] = string(*options.ConsistencyLevel)
	}

	if options.SessionToken != "" {
		headers[cosmosHeaderSessionToken] = options.SessionToken
	}

	if options.ResponseContinuationTokenLimitInKB > 0 {
		headers[cosmosHeaderResponseContinuationTokenLimitInKb] = fmt.Sprint(options.ResponseContinuationTokenLimitInKB)
	}

	if options.PageSizeHint != 0 {
		headers[cosmosHeaderMaxItemCount] = fmt.Sprint(options.PageSizeHint)
	}

	if options.EnableScanInQuery {
		headers[cosmosHeaderEnableScanInQuery] = "true"
	}

	if options.PopulateIndexMetrics {
		headers[cosmosHeaderPopulateIndexMetrics] = "true"
	}

	if options.ContinuationToken != "" {
		headers[cosmosHeaderContinuationToken] = options.ContinuationToken
	}

	headers[cosmosHeaderPopulateQueryMetrics] = "true"

	return &headers
}

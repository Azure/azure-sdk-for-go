// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// TransactionalBatchOptions includes options for transactional batch operations.
type TransactionalBatchOptions struct {
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
	// When EnableContentResponseOnWrite is false, the operations in the batch response will have no body, except when they are Read operations.
	// The default is false.
	EnableContentResponseOnWrite bool
}

// TransactionalBatchItemOptions includes options for the specific operation inside a TransactionalBatch
type TransactionalBatchItemOptions struct {
	// IfMatchETag is used to ensure optimistic concurrency control.
	// https://docs.microsoft.com/azure/cosmos-db/sql/database-transactions-optimistic-concurrency#optimistic-concurrency-control
	IfMatchETag *azcore.ETag
}

func (options *TransactionalBatchOptions) toHeaders() *map[string]string {
	headers := make(map[string]string, 2)

	if options.ConsistencyLevel != nil {
		headers[cosmosHeaderConsistencyLevel] = string(*options.ConsistencyLevel)
	}

	if options.SessionToken != "" {
		headers[cosmosHeaderSessionToken] = options.SessionToken
	}

	headers[cosmosHeaderIsBatchRequest] = "True"
	headers[cosmosHeaderIsBatchAtomic] = "True"
	headers[cosmosHeaderIsBatchOrdered] = "True"

	return &headers
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "time"

// DedicatedGatewayRequestOptions includes options for operations in the dedicated gateway.
type DedicatedGatewayRequestOptions struct {
	// Gets or sets the staleness value associated with the request in the Azure Cosmos DB service.
	// For requests where the ConsistencyLevel is ConsistencyLevel.Eventual or ConsistencyLevel.Session,
	// responses from the integrated cache are guaranteed to be no staler than value indicated by this MaxIntegratedCacheStaleness.
	// Cache Staleness is supported in milliseconds granularity. Anything smaller than milliseconds will be ignored.
	MaxIntegratedCacheStaleness *time.Duration

	// When set to true, the request will not be served from the integrated cache, and the response will not be cached either.
	BypassIntegratedCache bool
}

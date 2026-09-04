// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientOptions defines the options for the Cosmos client.
type ClientOptions struct {
	azcore.ClientOptions
	// When EnableContentResponseOnWrite is false will cause the response to have a null resource. This reduces networking and CPU load by not sending the resource back over the network and serializing it on the client.
	// The default is false.
	EnableContentResponseOnWrite bool
	// PreferredRegions is a list of regions to be used when initializing the client in case the default region fails.
	PreferredRegions []string
	// DisableEndpointDiscovery controls whether the SDK uses the account's advertised
	// writable/readable locations for request routing. When true, those advertised
	// locations are ignored for routing and every request is sent to the endpoint the
	// client was constructed with. The account metadata request itself is not disabled;
	// discovery may still run, but its result is not used to select the target endpoint.
	// This is required in environments where the account advertises an internal document
	// endpoint (e.g. *.docdb.azs) that is not reachable from the client's network and all
	// traffic must instead flow through the externally-routable endpoint (e.g. via a
	// reverse proxy). When true, cross-region failover is suppressed since all requests
	// resolve to the single client endpoint, but same-region (same-endpoint) transport
	// retries are preserved. PreferredRegions may still be set but has no effect on
	// routing while discovery is disabled. The default is false (routing follows
	// discovery), matching prior behavior.
	DisableEndpointDiscovery bool
	// PriorityLevel defines the default priority level for all requests made by this client.
	// This feature is currently in preview. For more information, see https://aka.ms/CosmosDB/PriorityBasedExecution
	// Valid values are PriorityLevelHigh and PriorityLevelLow.
	// Can be overridden per-request via the operation options.
	PriorityLevel *PriorityLevel
	// ThroughputBucket defines the default throughput bucket for all requests made by this client.
	// This feature is currently in preview. For more information, see https://aka.ms/CosmosDB/ThroughputBuckets
	// The valid range is 1 to 5 (inclusive).
	// Can be overridden per-request via the operation options.
	ThroughputBucket *int32
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"time"

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
	// ThrottlingRetryOptions configures how the client retries requests that fail with
	// HTTP 429 (Too Many Requests). When unset, defaults consistent with the other
	// Cosmos SDKs are used (9 attempts, 30s cumulative wait).
	ThrottlingRetryOptions ThrottlingRetryOptions
}

// ThrottlingRetryOptions configures the retry behavior for HTTP 429
// (Too Many Requests) responses. The Cosmos service indicates the recommended
// retry delay via the x-ms-retry-after-ms response header; the client respects
// that value subject to the limits in this struct.
type ThrottlingRetryOptions struct {
	// MaxRetryAttempts is the maximum number of times the client will retry a
	// throttled request. The default is 9. Set to a negative value to disable
	// throttling retries.
	MaxRetryAttempts int
	// MaxRetryWaitTime is the maximum cumulative time the client will spend
	// waiting between throttled retries for a single request. Once this budget
	// is exhausted, the most recent 429 response is returned to the caller.
	// The default is 30 seconds.
	MaxRetryWaitTime time.Duration
}

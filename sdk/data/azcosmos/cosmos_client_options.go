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
	// PriorityLevel defines the default priority level for all requests made by this client.
	// Valid values are PriorityLevelHigh and PriorityLevelLow.
	// Can be overridden per-request via the operation options.
	PriorityLevel *PriorityLevel
	// ThroughputBucket defines the default throughput bucket for all requests made by this client.
	// The valid range is 1 to 5 (inclusive).
	// Can be overridden per-request via the operation options.
	ThroughputBucket *int32
}

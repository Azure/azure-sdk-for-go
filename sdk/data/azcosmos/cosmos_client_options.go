// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// DirectModeOptions configures the Direct mode connection pool and TCP behavior.
// These options only apply when ConnectionMode is set to ConnectionModeDirect.
type DirectModeOptions struct {
	// MaxRequestsPerConnection is the maximum number of concurrent requests allowed per TCP connection.
	// Higher values improve throughput but increase memory usage per connection.
	// Default: 30.
	MaxRequestsPerConnection int
	// IdleConnectionTimeout is the duration after which an idle TCP connection is closed.
	// Setting to 0 uses the server-provided value from context negotiation.
	// Default: 0 (use server value).
	IdleConnectionTimeout time.Duration
	// MaxConnectionsPerEndpoint is the maximum number of TCP connections per backend endpoint.
	// Higher values allow more parallel requests but consume more resources.
	// Default: 10.
	MaxConnectionsPerEndpoint int
	// ConnectTimeout is the timeout for establishing a new TCP connection.
	// Default: 5 seconds.
	ConnectTimeout time.Duration
}

// ClientOptions defines the options for the Cosmos client.
type ClientOptions struct {
	azcore.ClientOptions
	// ConnectionMode specifies how the client connects to Azure Cosmos DB.
	// Gateway mode (default) routes all requests through the Cosmos DB gateway.
	// Direct mode connects directly to backend nodes using RNTBD protocol for data operations.
	// See ConnectionModeGateway and ConnectionModeDirect for details.
	ConnectionMode ConnectionMode
	// DirectModeOptions configures Direct mode connection pool behavior.
	// Only used when ConnectionMode is ConnectionModeDirect.
	DirectModeOptions *DirectModeOptions
	// When EnableContentResponseOnWrite is false will cause the response to have a null resource. This reduces networking and CPU load by not sending the resource back over the network and serializing it on the client.
	// The default is false.
	EnableContentResponseOnWrite bool
	// PreferredRegions is a list of regions to be used when initializing the client in case the default region fails.
	PreferredRegions []string
}

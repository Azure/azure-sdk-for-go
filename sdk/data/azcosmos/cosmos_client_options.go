// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientOptions defines the options for the Cosmos client.
type ClientOptions struct {
	azcore.ClientOptions
	// ConnectionMode specifies how the client connects to Azure Cosmos DB.
	// Gateway mode (default) routes all requests through the Cosmos DB gateway.
	// Direct mode connects directly to backend nodes using RNTBD protocol for data operations.
	// See ConnectionModeGateway and ConnectionModeDirect for details.
	ConnectionMode ConnectionMode
	// When EnableContentResponseOnWrite is false will cause the response to have a null resource. This reduces networking and CPU load by not sending the resource back over the network and serializing it on the client.
	// The default is false.
	EnableContentResponseOnWrite bool
	// PreferredRegions is a list of regions to be used when initializing the client in case the default region fails.
	PreferredRegions []string
}

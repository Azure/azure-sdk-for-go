// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// A CosmosClient is used to interact with the Azure Cosmos DB database service.
type CosmosClient struct {
	// Endpoint used to create the client.
	Endpoint   string
	connection *clientConnection
	cred       CosmosAccountCredential
}

// NewCosmosClient creates a new instance of CosmosClient with the specified values. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional CosmosClient options.  Pass nil to accept default values.
func NewCosmosClient(endpoint string, cred azcore.Credential, options *CosmosClientOptions) (*CosmosClient, error) {
	connection := newConnection(endpoint, cred, options)

	c, _ := cred.(*SharedKeyCredential)

	return &CosmosClient{Endpoint: endpoint, connection: connection, cred: c}, nil
}

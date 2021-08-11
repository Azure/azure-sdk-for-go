// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// A CosmosClient is used to interact with the Azure Cosmos DB database service.
type CosmosClient struct {
	// Endpoint used to create the client.
	Endpoint   string
	connection *cosmosClientConnection
	cred       CosmosAccountCredential
}

// NewCosmosClient creates a new instance of CosmosClient with the specified values. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional CosmosClient options.  Pass nil to accept default values.
func NewCosmosClient(endpoint string, cred azcore.Credential, options *CosmosClientOptions) (*CosmosClient, error) {
	connection := newCosmosClientConnection(endpoint, cred, options)

	c, _ := cred.(*SharedKeyCredential)

	return &CosmosClient{Endpoint: endpoint, connection: connection, cred: c}, nil
}

// GetCosmosDatabase returns a CosmosDatabase object.
// id - The id of the database.
func (c *CosmosClient) GetCosmosDatabase(id string) (*CosmosDatabase, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return newCosmosDatabase(id, c), nil
}

// GetCosmosContainer returns a CosmosContainer object.
// databaseId - The id of the database.
// containerId - The id of the container.
func (c *CosmosClient) GetCosmosContainer(databaseId string, containerId string) (*CosmosContainer, error) {
	if databaseId == "" {
		return nil, errors.New("databaseId is required")
	}

	if containerId == "" {
		return nil, errors.New("containerId is required")
	}

	return newCosmosDatabase(databaseId, c).GetContainer(containerId)
}

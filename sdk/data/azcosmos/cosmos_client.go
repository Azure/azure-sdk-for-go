// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
)

// A CosmosClient is used to interact with the Azure Cosmos DB database service.
type CosmosClient struct {
	// Endpoint used to create the client.
	Endpoint   string
	connection *cosmosClientConnection
	cred       *SharedKeyCredential
	options    *CosmosClientOptions
}

// NewClientWithSharedKey creates a new instance of CosmosClient with the specified values. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional CosmosClient options.  Pass nil to accept default values.
func NewClientWithSharedKey(endpoint string, cred *SharedKeyCredential, options *CosmosClientOptions) (*CosmosClient, error) {
	if options == nil {
		options = &CosmosClientOptions{}
	}

	connection := newCosmosClientConnection(endpoint, cred, options)

	return &CosmosClient{Endpoint: endpoint, connection: connection, cred: cred, options: options}, nil
}

// GetDatabase returns a Database object.
// id - The id of the database.
func (c *CosmosClient) GetDatabase(id string) (*Database, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return newDatabase(id, c), nil
}

// GetContainer returns a Container object.
// databaseId - The id of the database.
// containerId - The id of the container.
func (c *CosmosClient) GetContainer(databaseId string, containerId string) (*Container, error) {
	if databaseId == "" {
		return nil, errors.New("databaseId is required")
	}

	if containerId == "" {
		return nil, errors.New("containerId is required")
	}

	return newDatabase(databaseId, c).GetContainer(containerId)
}

// CreateDatabase creates a new database.
// ctx - The context for the request.
// databaseProperties - The definition of the database
// throughputProperties - Optional throughput configuration of the database
// requestOptions - Optional parameters for the request.
func (c *CosmosClient) CreateDatabase(
	ctx context.Context,
	databaseProperties DatabaseProperties,
	throughputProperties *ThroughputProperties,
	requestOptions *DatabaseRequestOptions) (DatabaseResponse, error) {
	if requestOptions == nil {
		requestOptions = &DatabaseRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	path, err := generatePathForNameBased(resourceTypeDatabase, "", true)
	if err != nil {
		return DatabaseResponse{}, err
	}

	database, err := c.GetDatabase(databaseProperties.Id)
	if err != nil {
		return DatabaseResponse{}, err
	}

	azResponse, err := c.connection.sendPostRequest(
		path,
		ctx,
		databaseProperties,
		operationContext,
		requestOptions,
		throughputProperties.addHeadersToRequest)
	if err != nil {
		return DatabaseResponse{}, err
	}

	return newDatabaseResponse(azResponse, database)
}

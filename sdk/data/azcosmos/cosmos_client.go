// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
)

// Cosmos client is used to interact with the Azure Cosmos DB database service.
type Client struct {
	endpoint   string
	connection *cosmosClientConnection
	cred       *KeyCredential
	options    *CosmosClientOptions
}

// Endpoint used to create the client.
func (c *Client) Endpoint() string {
	return c.endpoint
}

// NewClientWithKey creates a new instance of Cosmos client with the specified values. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional Cosmos client options.  Pass nil to accept default values.
func NewClientWithKey(endpoint string, cred *KeyCredential, options *CosmosClientOptions) (*Client, error) {
	if options == nil {
		options = &CosmosClientOptions{}
	}

	connection := newCosmosClientConnection(endpoint, cred, options)

	return &Client{endpoint: endpoint, connection: connection, cred: cred, options: options}, nil
}

// GetDatabase returns a Database object.
// id - The id of the database.
func (c *Client) GetDatabase(id string) (*Database, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return newDatabase(id, c), nil
}

// GetContainer returns a Container object.
// databaseId - The id of the database.
// containerId - The id of the container.
func (c *Client) GetContainer(databaseId string, containerId string) (*Container, error) {
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
// o - Options for the create database operation.
func (c *Client) CreateDatabase(
	ctx context.Context,
	databaseProperties DatabaseProperties,
	o *CreateDatabaseOptions) (DatabaseResponse, error) {
	if o == nil {
		o = &CreateDatabaseOptions{}
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
		nil,
		o.ThroughputProperties.addHeadersToRequest)
	if err != nil {
		return DatabaseResponse{}, err
	}

	return newDatabaseResponse(azResponse, database)
}

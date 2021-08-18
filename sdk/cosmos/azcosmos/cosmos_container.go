// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "context"

// A CosmosContainer lets you perform read, update, change throughput, and delete container operations.
// It also lets you perform read, update, change throughput, and delete item operations.
type CosmosContainer struct {
	// The Id of the Cosmos container
	Id string
	// The database that contains the container
	Database *CosmosDatabase
	// The resource link
	link string
}

func newCosmosContainer(id string, database *CosmosDatabase) *CosmosContainer {
	return &CosmosContainer{
		Id:       id,
		Database: database,
		link:     createLink(database.link, pathSegmentCollection, id)}
}

// Get reads a Cosmos container.
// ctx - The context for the request.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) Get(
	ctx context.Context,
	requestOptions *CosmosContainerRequestOptions) (CosmosContainerResponse, error) {
	if requestOptions == nil {
		requestOptions = &CosmosContainerRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendGetRequest(path, ctx, operationContext, requestOptions)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	return newCosmosContainerResponse(azResponse, c)
}

// Update a Cosmos container.
// ctx - The context for the request.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) Update(
	ctx context.Context,
	containerProperties CosmosContainerProperties,
	requestOptions *CosmosContainerRequestOptions) (CosmosContainerResponse, error) {
	if requestOptions == nil {
		requestOptions = &CosmosContainerRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPutRequest(path, ctx, containerProperties, operationContext, requestOptions)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	return newCosmosContainerResponse(azResponse, c)
}

// Delete a Cosmos container.
// ctx - The context for the request.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) Delete(
	ctx context.Context,
	requestOptions *CosmosContainerRequestOptions) (CosmosContainerResponse, error) {
	if requestOptions == nil {
		requestOptions = &CosmosContainerRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendDeleteRequest(path, ctx, operationContext, requestOptions)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	return newCosmosContainerResponse(azResponse, c)
}

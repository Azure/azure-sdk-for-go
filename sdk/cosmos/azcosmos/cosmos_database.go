// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
)

// A CosmosDatabase lets you perform read, update, change throughput, and delete database operations.
type CosmosDatabase struct {
	// The Id of the Cosmos database
	Id string
	// The client associated with the Cosmos database
	client *CosmosClient
	// The resource link
	link string
}

func newCosmosDatabase(id string, client *CosmosClient) *CosmosDatabase {
	return &CosmosDatabase{
		Id:     id,
		client: client,
		link:   createLink("", pathSegmentDatabase, id)}
}

// GetContainer returns a CosmosContainer object for the container.
// id - The id of the container.
func (db *CosmosDatabase) GetContainer(id string) (*CosmosContainer, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return newCosmosContainer(id, db), nil
}

func (db *CosmosDatabase) AddContainer(ctx context.Context, containerProperties CosmosContainerProperties, requestOptions *CosmosContainerRequestOptions) (CosmosContainerResponse, error) {
	if requestOptions == nil {
		requestOptions = &CosmosContainerRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: db.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, db.link, true)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	container, err := db.GetContainer(containerProperties.Id)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	azResponse, err := db.client.connection.sendPostRequest(path, ctx, containerProperties, operationContext, requestOptions)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	return newCosmosContainerResponse(azResponse, container)
}

func (db *CosmosDatabase) Get(ctx context.Context, requestOptions *CosmosDatabaseRequestOptions) (CosmosDatabaseResponse, error) {
	if requestOptions == nil {
		requestOptions = &CosmosDatabaseRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: db.link,
	}

	path, err := generatePathForNameBased(resourceTypeDatabase, db.link, false)
	if err != nil {
		return CosmosDatabaseResponse{}, err
	}

	azResponse, err := db.client.connection.sendGetRequest(path, ctx, operationContext, requestOptions)
	if err != nil {
		return CosmosDatabaseResponse{}, err
	}

	return newCosmosDatabaseResponse(azResponse, db)
}

func (db *CosmosDatabase) Delete(ctx context.Context, requestOptions *CosmosDatabaseRequestOptions) (CosmosDatabaseResponse, error) {
	if requestOptions == nil {
		requestOptions = &CosmosDatabaseRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: db.link,
	}

	path, err := generatePathForNameBased(resourceTypeDatabase, db.link, false)
	if err != nil {
		return CosmosDatabaseResponse{}, err
	}

	azResponse, err := db.client.connection.sendDeleteRequest(path, ctx, operationContext, requestOptions)
	if err != nil {
		return CosmosDatabaseResponse{}, err
	}

	return newCosmosDatabaseResponse(azResponse, db)
}

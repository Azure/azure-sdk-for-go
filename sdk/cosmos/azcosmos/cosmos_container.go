// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

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

// Read obtains the information for a Cosmos container.
// ctx - The context for the request.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) Read(
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

	azResponse, err := c.Database.client.connection.sendGetRequest(
		path,
		ctx,
		operationContext,
		requestOptions,
		nil)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	return newCosmosContainerResponse(azResponse, c)
}

// Replace a Cosmos container.
// ctx - The context for the request.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) Replace(
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

	azResponse, err := c.Database.client.connection.sendPutRequest(
		path,
		ctx,
		containerProperties,
		operationContext,
		requestOptions,
		nil)
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

	azResponse, err := c.Database.client.connection.sendDeleteRequest(
		path,
		ctx,
		operationContext,
		requestOptions,
		nil)
	if err != nil {
		return CosmosContainerResponse{}, err
	}

	return newCosmosContainerResponse(azResponse, c)
}

// Creates an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// item - The item to create.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) CreateItem(
	ctx context.Context,
	partitionKey PartitionKey,
	item interface{},
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeDocument, c.link, true)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPostRequest(
		path,
		ctx,
		item,
		operationContext,
		requestOptions,
		addHeader)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	return newCosmosItemResponse(azResponse)
}

// Replaces an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key of the item to replace.
// itemId - The id of the item to replace.
// item - The content to be used to replace.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) ReplaceItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	item interface{},
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: itemId,
	}

	path, err := generatePathForNameBased(resourceTypeDocument, createLink(c.link, pathSegmentDocument, itemId), false)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPutRequest(
		path,
		ctx,
		item,
		operationContext,
		requestOptions,
		addHeader)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	return newCosmosItemResponse(azResponse)
}

// Reads an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// itemId - The id of the item to read.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) ReadItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: itemId,
	}

	path, err := generatePathForNameBased(resourceTypeDocument, createLink(c.link, pathSegmentDocument, itemId), false)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendGetRequest(
		path,
		ctx,
		operationContext,
		requestOptions,
		addHeader)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	return newCosmosItemResponse(azResponse)
}

// Deletes an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// itemId - The id of the item to delete.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) DeleteItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: itemId,
	}

	path, err := generatePathForNameBased(resourceTypeDocument, createLink(c.link, pathSegmentDocument, itemId), false)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendDeleteRequest(
		path,
		ctx,
		operationContext,
		requestOptions,
		addHeader)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	return newCosmosItemResponse(azResponse)
}

func (c *CosmosContainer) buildRequestEnricher(
	partitionKey PartitionKey,
	requestOptions *CosmosItemRequestOptions) (func(r *azcore.Request), error) {
	pk, err := partitionKey.toJsonString()
	if err != nil {
		return nil, err
	}

	var enableContentResponseOnWrite bool
	if requestOptions == nil {
		enableContentResponseOnWrite = c.Database.client.options.EnableContentResponseOnWrite
	} else {
		enableContentResponseOnWrite = requestOptions.EnableContentResponseOnWrite
	}

	return func(r *azcore.Request) {
		r.Header.Add(cosmosHeaderPartitionKey, pk)
		if !enableContentResponseOnWrite {
			r.Header.Add(cosmosHeaderPrefer, cosmosHeaderValuesPreferMinimal)
		}
	}, nil
}

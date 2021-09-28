// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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

// ReadThroughput obtains the provisioned throughput information for the container.
// ctx - The context for the request.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) ReadThroughput(
	ctx context.Context,
	requestOptions *ThroughputRequestOptions) (ThroughputResponse, error) {
	if requestOptions == nil {
		requestOptions = &ThroughputRequestOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{connection: c.Database.client.connection}
	return offers.ReadThroughputIfExists(ctx, rid, requestOptions)
}

// ReplaceThroughput updates the provisioned throughput for the container.
// ctx - The context for the request.
// throughputProperties - The throughput configuration of the container.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) ReplaceThroughput(
	ctx context.Context,
	throughputProperties ThroughputProperties,
	requestOptions *ThroughputRequestOptions) (ThroughputResponse, error) {
	if requestOptions == nil {
		requestOptions = &ThroughputRequestOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{connection: c.Database.client.connection}
	return offers.ReadThroughputIfExists(ctx, rid, requestOptions)
}

// Creates an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// item - The item to create.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) CreateItem(
	ctx context.Context,
	partitionKey *PartitionKey,
	item interface{},
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions, true)
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

// Upserts (create or replace) an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// item - The item to upsert.
// requestOptions - Optional parameters for the request.
func (c *CosmosContainer) UpsertItem(
	ctx context.Context,
	partitionKey *PartitionKey,
	item interface{},
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeaderInternal, err := c.buildRequestEnricher(partitionKey, requestOptions, true)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	addHeader := func(r *policy.Request) {
		addHeaderInternal(r)
		r.Raw().Header.Add(cosmosHeaderIsUpsert, "true")
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
	partitionKey *PartitionKey,
	itemId string,
	item interface{},
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions, true)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: createLink(c.link, pathSegmentDocument, itemId),
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
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
	partitionKey *PartitionKey,
	itemId string,
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions, false)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: createLink(c.link, pathSegmentDocument, itemId),
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
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
	partitionKey *PartitionKey,
	itemId string,
	requestOptions *CosmosItemRequestOptions) (CosmosItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, requestOptions, true)
	if err != nil {
		return CosmosItemResponse{}, err
	}

	if requestOptions == nil {
		requestOptions = &CosmosItemRequestOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: createLink(c.link, pathSegmentDocument, itemId),
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
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
	partitionKey *PartitionKey,
	requestOptions *CosmosItemRequestOptions,
	writeOperation bool) (func(r *policy.Request), error) {
	if partitionKey == nil {
		return nil, errors.New("partitionKey is required")
	}

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

	return func(r *policy.Request) {
		r.Raw().Header.Add(cosmosHeaderPartitionKey, pk)
		if writeOperation && !enableContentResponseOnWrite {
			r.Raw().Header.Add(cosmosHeaderPrefer, cosmosHeaderValuesPreferMinimal)
		}
	}, nil
}

func (c *CosmosContainer) getRID(ctx context.Context) (string, error) {
	containerResponse, err := c.Read(ctx, nil)
	if err != nil {
		return "", err
	}

	return containerResponse.ContainerProperties.ResourceId, nil
}

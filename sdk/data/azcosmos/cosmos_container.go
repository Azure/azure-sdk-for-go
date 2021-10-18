// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// A Container lets you perform read, update, change throughput, and delete container operations.
// It also lets you perform read, update, change throughput, and delete item operations.
type Container struct {
	// The Id of the Cosmos container
	Id string
	// The database that contains the container
	Database *Database
	// The resource link
	link string
}

func newContainer(id string, database *Database) *Container {
	return &Container{
		Id:       id,
		Database: database,
		link:     createLink(database.link, pathSegmentCollection, id)}
}

// Read obtains the information for a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *Container) Read(
	ctx context.Context,
	o *ReadContainerOptions) (ContainerResponse, error) {
	if o == nil {
		o = &ReadContainerOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return ContainerResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
	if err != nil {
		return ContainerResponse{}, err
	}

	return newContainerResponse(azResponse, c)
}

// Replace a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *Container) Replace(
	ctx context.Context,
	containerProperties ContainerProperties,
	o *ReplaceContainerOptions) (ContainerResponse, error) {
	if o == nil {
		o = &ReplaceContainerOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return ContainerResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPutRequest(
		path,
		ctx,
		containerProperties,
		operationContext,
		o,
		nil)
	if err != nil {
		return ContainerResponse{}, err
	}

	return newContainerResponse(azResponse, c)
}

// Delete a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *Container) Delete(
	ctx context.Context,
	o *DeleteContainerOptions) (ContainerResponse, error) {
	if o == nil {
		o = &DeleteContainerOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return ContainerResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendDeleteRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
	if err != nil {
		return ContainerResponse{}, err
	}

	return newContainerResponse(azResponse, c)
}

// ReadThroughput obtains the provisioned throughput information for the container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *Container) ReadThroughput(
	ctx context.Context,
	o *ThroughputOptions) (ThroughputResponse, error) {
	if o == nil {
		o = &ThroughputOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{connection: c.Database.client.connection}
	return offers.ReadThroughputIfExists(ctx, rid, o)
}

// ReplaceThroughput updates the provisioned throughput for the container.
// ctx - The context for the request.
// throughputProperties - The throughput configuration of the container.
// o - Options for the operation.
func (c *Container) ReplaceThroughput(
	ctx context.Context,
	throughputProperties ThroughputProperties,
	o *ThroughputOptions) (ThroughputResponse, error) {
	if o == nil {
		o = &ThroughputOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{connection: c.Database.client.connection}
	return offers.ReadThroughputIfExists(ctx, rid, o)
}

// Creates an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// item - The item to create.
// o - Options for the operation.
func (c *Container) CreateItem(
	ctx context.Context,
	partitionKey PartitionKey,
	item interface{},
	o *ItemOptions) (ItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, o, true)
	if err != nil {
		return ItemResponse{}, err
	}

	if o == nil {
		o = &ItemOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeDocument, c.link, true)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPostRequest(
		path,
		ctx,
		item,
		operationContext,
		o,
		addHeader)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

// Upserts (create or replace) an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// item - The item to upsert.
// o - Options for the operation.
func (c *Container) UpsertItem(
	ctx context.Context,
	partitionKey PartitionKey,
	item interface{},
	o *ItemOptions) (ItemResponse, error) {

	addHeaderInternal, err := c.buildRequestEnricher(partitionKey, o, true)
	if err != nil {
		return ItemResponse{}, err
	}

	addHeader := func(r *policy.Request) {
		addHeaderInternal(r)
		r.Raw().Header.Add(cosmosHeaderIsUpsert, "true")
	}

	if o == nil {
		o = &ItemOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeDocument, c.link, true)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPostRequest(
		path,
		ctx,
		item,
		operationContext,
		o,
		addHeader)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

// Replaces an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key of the item to replace.
// itemId - The id of the item to replace.
// item - The content to be used to replace.
// o - Options for the operation.
func (c *Container) ReplaceItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	item interface{},
	o *ItemOptions) (ItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, o, true)
	if err != nil {
		return ItemResponse{}, err
	}

	if o == nil {
		o = &ItemOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: createLink(c.link, pathSegmentDocument, itemId),
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendPutRequest(
		path,
		ctx,
		item,
		operationContext,
		o,
		addHeader)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

// Reads an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// itemId - The id of the item to read.
// o - Options for the operation.
func (c *Container) ReadItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	o *ItemOptions) (ItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, o, false)
	if err != nil {
		return ItemResponse{}, err
	}

	if o == nil {
		o = &ItemOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: createLink(c.link, pathSegmentDocument, itemId),
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		addHeader)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

// Deletes an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// itemId - The id of the item to delete.
// o - Options for the operation.
func (c *Container) DeleteItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	o *ItemOptions) (ItemResponse, error) {

	addHeader, err := c.buildRequestEnricher(partitionKey, o, true)
	if err != nil {
		return ItemResponse{}, err
	}

	if o == nil {
		o = &ItemOptions{}
	}

	operationContext := cosmosOperationContext{
		resourceType:    resourceTypeDocument,
		resourceAddress: createLink(c.link, pathSegmentDocument, itemId),
	}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.Database.client.connection.sendDeleteRequest(
		path,
		ctx,
		operationContext,
		o,
		addHeader)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

func (c *Container) buildRequestEnricher(
	partitionKey PartitionKey,
	requestOptions *ItemOptions,
	writeOperation bool) (func(r *policy.Request), error) {
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

func (c *Container) getRID(ctx context.Context) (string, error) {
	containerResponse, err := c.Read(ctx, nil)
	if err != nil {
		return "", err
	}

	return containerResponse.ContainerProperties.ResourceId, nil
}

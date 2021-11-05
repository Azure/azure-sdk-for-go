// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// A ContainerClient lets you perform read, update, change throughput, and delete container operations.
// It also lets you perform read, update, change throughput, and delete item operations.
type ContainerClient struct {
	// The Id of the Cosmos container
	id string
	// The database that contains the container
	database *DatabaseClient
	// The resource link
	link string
}

func newContainer(id string, database *DatabaseClient) (*ContainerClient, error) {
	return &ContainerClient{
		id:       id,
		database: database,
		link:     createLink(database.link, pathSegmentCollection, id)}, nil
}

// ID returns the identifier of the Cosmos container.
func (c *ContainerClient) ID() string {
	return c.id
}

// Read obtains the information for a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *ContainerClient) Read(
	ctx context.Context,
	o *ReadContainerOptions) (ContainerResponse, error) {
	if o == nil {
		o = &ReadContainerOptions{}
	}

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return ContainerResponse{}, err
	}

	azResponse, err := c.database.client.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
	if err != nil {
		return ContainerResponse{}, err
	}

	return newContainerResponse(azResponse)
}

// Replace a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *ContainerClient) Replace(
	ctx context.Context,
	containerProperties ContainerProperties,
	o *ReplaceContainerOptions) (ContainerResponse, error) {
	if o == nil {
		o = &ReplaceContainerOptions{}
	}

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return ContainerResponse{}, err
	}

	azResponse, err := c.database.client.sendPutRequest(
		path,
		ctx,
		containerProperties,
		operationContext,
		o,
		nil)
	if err != nil {
		return ContainerResponse{}, err
	}

	return newContainerResponse(azResponse)
}

// Delete a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *ContainerClient) Delete(
	ctx context.Context,
	o *DeleteContainerOptions) (ContainerResponse, error) {
	if o == nil {
		o = &DeleteContainerOptions{}
	}

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeCollection,
		resourceAddress: c.link,
	}

	path, err := generatePathForNameBased(resourceTypeCollection, c.link, false)
	if err != nil {
		return ContainerResponse{}, err
	}

	azResponse, err := c.database.client.sendDeleteRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
	if err != nil {
		return ContainerResponse{}, err
	}

	return newContainerResponse(azResponse)
}

// ReadThroughput obtains the provisioned throughput information for the container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *ContainerClient) ReadThroughput(
	ctx context.Context,
	o *ThroughputOptions) (ThroughputResponse, error) {
	if o == nil {
		o = &ThroughputOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{client: c.database.client}
	return offers.ReadThroughputIfExists(ctx, rid, o)
}

// ReplaceThroughput updates the provisioned throughput for the container.
// ctx - The context for the request.
// throughputProperties - The throughput configuration of the container.
// o - Options for the operation.
func (c *ContainerClient) ReplaceThroughput(
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

	offers := &cosmosOffers{client: c.database.client}
	return offers.ReadThroughputIfExists(ctx, rid, o)
}

// Creates an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// item - The item to create.
// o - Options for the operation.
func (c *ContainerClient) CreateItem(
	ctx context.Context,
	partitionKey PartitionKey,
	item []byte,
	o *ItemOptions) (ItemResponse, error) {
	h := headerOptionsOverride{
		partitionKey: &partitionKey,
	}

	if o == nil {
		o = &ItemOptions{}
	} else {
		h.enableContentResponseOnWrite = &o.EnableContentResponseOnWrite
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       c.link,
		isWriteOperation:      true,
		headerOptionsOverride: &h}

	path, err := generatePathForNameBased(resourceTypeDocument, c.link, true)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.database.client.sendPostRequest(
		path,
		ctx,
		item,
		operationContext,
		o,
		nil)
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
func (c *ContainerClient) UpsertItem(
	ctx context.Context,
	partitionKey PartitionKey,
	item []byte,
	o *ItemOptions) (ItemResponse, error) {
	h := headerOptionsOverride{
		partitionKey: &partitionKey,
	}

	addHeader := func(r *policy.Request) {
		r.Raw().Header.Add(cosmosHeaderIsUpsert, "true")
	}

	if o == nil {
		o = &ItemOptions{}
	} else {
		h.enableContentResponseOnWrite = &o.EnableContentResponseOnWrite
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       c.link,
		isWriteOperation:      true,
		headerOptionsOverride: &h}

	path, err := generatePathForNameBased(resourceTypeDocument, c.link, true)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.database.client.sendPostRequest(
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
func (c *ContainerClient) ReplaceItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	item []byte,
	o *ItemOptions) (ItemResponse, error) {
	h := headerOptionsOverride{
		partitionKey: &partitionKey,
	}

	if o == nil {
		o = &ItemOptions{}
	} else {
		h.enableContentResponseOnWrite = &o.EnableContentResponseOnWrite
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       createLink(c.link, pathSegmentDocument, itemId),
		isWriteOperation:      true,
		headerOptionsOverride: &h}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.database.client.sendPutRequest(
		path,
		ctx,
		item,
		operationContext,
		o,
		nil)
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
func (c *ContainerClient) ReadItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	o *ItemOptions) (ItemResponse, error) {
	h := headerOptionsOverride{
		partitionKey: &partitionKey,
	}

	if o == nil {
		o = &ItemOptions{}
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       createLink(c.link, pathSegmentDocument, itemId),
		headerOptionsOverride: &h}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.database.client.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
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
func (c *ContainerClient) DeleteItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	o *ItemOptions) (ItemResponse, error) {
	h := headerOptionsOverride{
		partitionKey: &partitionKey,
	}

	if o == nil {
		o = &ItemOptions{}
	} else {
		h.enableContentResponseOnWrite = &o.EnableContentResponseOnWrite
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       createLink(c.link, pathSegmentDocument, itemId),
		isWriteOperation:      true,
		headerOptionsOverride: &h}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, false)
	if err != nil {
		return ItemResponse{}, err
	}

	azResponse, err := c.database.client.sendDeleteRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

func (c *ContainerClient) getRID(ctx context.Context) (string, error) {
	containerResponse, err := c.Read(ctx, nil)
	if err != nil {
		return "", err
	}

	return containerResponse.ContainerProperties.ResourceID, nil
}

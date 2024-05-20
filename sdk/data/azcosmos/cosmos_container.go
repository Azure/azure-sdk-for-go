// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

// ContainerClient lets you perform read, update, change throughput, and delete container operations.
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
		resourceType:     resourceTypeCollection,
		resourceAddress:  c.link,
		isWriteOperation: true,
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
		resourceType:     resourceTypeCollection,
		resourceAddress:  c.link,
		isWriteOperation: true,
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
	return offers.ReplaceThroughputIfExists(ctx, throughputProperties, rid, o)
}

// CreateItem creates an item in a Cosmos container.
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

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
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

// UpsertItem creates or replaces an item in a Cosmos container.
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

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
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

// ReplaceItem replaces an item in a Cosmos container.
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

// ReadItem reads an item in a Cosmos container.
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

// DeleteItem deletes an item in a Cosmos container.
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

// NewQueryItemsPager executes a single partition query in a Cosmos container.
// query - The SQL query to execute.
// partitionKey - The partition key to scope the query on.
// o - Options for the operation.
func (c *ContainerClient) NewQueryItemsPager(query string, partitionKey PartitionKey, o *QueryOptions) *runtime.Pager[QueryItemsResponse] {
	correlatedActivityId, _ := uuid.New()
	h := headerOptionsOverride{
		partitionKey:         &partitionKey,
		correlatedActivityId: &correlatedActivityId,
	}

	queryOptions := &QueryOptions{}
	if o != nil {
		originalOptions := *o
		queryOptions = &originalOptions
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       c.link,
		headerOptionsOverride: &h,
	}

	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	return runtime.NewPager(runtime.PagingHandler[QueryItemsResponse]{
		More: func(page QueryItemsResponse) bool {
			return page.ContinuationToken != nil
		},
		Fetcher: func(ctx context.Context, page *QueryItemsResponse) (QueryItemsResponse, error) {
			if page != nil {
				if page.ContinuationToken != nil {
					// Use the previous page continuation if available
					queryOptions.ContinuationToken = page.ContinuationToken
				}
			}

			azResponse, err := c.database.client.sendQueryRequest(
				path,
				ctx,
				query,
				queryOptions.QueryParameters,
				operationContext,
				queryOptions,
				nil)

			if err != nil {
				return QueryItemsResponse{}, err
			}

			return newQueryResponse(azResponse)
		},
	})
}

// PatchItem patches an item in a Cosmos container.
// ctx - The context for the request.
// partitionKey - The partition key for the item.
// itemId - The id of the item to patch.
// ops - Operations to perform on the patch
// o - Options for the operation.
func (c *ContainerClient) PatchItem(
	ctx context.Context,
	partitionKey PartitionKey,
	itemId string,
	ops PatchOperations,
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

	azResponse, err := c.database.client.sendPatchRequest(
		path,
		ctx,
		ops,
		operationContext,
		o,
		nil)
	if err != nil {
		return ItemResponse{}, err
	}

	return newItemResponse(azResponse)
}

// NewTransactionalBatch creates a batch of operations to be committed as a single unit.
// See https://docs.microsoft.com/azure/cosmos-db/sql/transactional-batch
func (c *ContainerClient) NewTransactionalBatch(partitionKey PartitionKey) TransactionalBatch {
	return TransactionalBatch{partitionKey: partitionKey}
}

// ExecuteTransactionalBatch executes a transactional batch.
// Once executed, verify the Success property of the response to determine if the batch was committed
func (c *ContainerClient) ExecuteTransactionalBatch(ctx context.Context, b TransactionalBatch, o *TransactionalBatchOptions) (TransactionalBatchResponse, error) {
	if len(b.operations) == 0 {
		return TransactionalBatchResponse{}, errors.New("no operations in batch")
	}

	h := headerOptionsOverride{
		partitionKey: &b.partitionKey,
	}

	if o == nil {
		o = &TransactionalBatchOptions{}
	} else {
		h.enableContentResponseOnWrite = &o.EnableContentResponseOnWrite
	}

	// If contentResponseOnWrite is not enabled at the client level the
	// service will not even send a batch response payload
	// Instead we should automatically enforce contentResponseOnWrite for all
	// batch requests whenever at least one of the item operations requires a content response (read operation)
	enableContentResponseOnWriteForReadOperations := true
	for _, op := range b.operations {
		if op.getOperationType() == operationTypeRead {
			h.enableContentResponseOnWrite = &enableContentResponseOnWriteForReadOperations
			break
		}
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDocument,
		resourceAddress:       c.link,
		isWriteOperation:      true,
		headerOptionsOverride: &h}

	path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}

	azResponse, err := c.database.client.sendBatchRequest(
		ctx,
		path,
		b.operations,
		operationContext,
		o,
		nil)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}

	return newTransactionalBatchResponse(azResponse)
}

func (c *ContainerClient) getRID(ctx context.Context) (string, error) {
	containerResponse, err := c.Read(ctx, nil)
	if err != nil {
		return "", err
	}

	return containerResponse.ContainerProperties.ResourceID, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
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

// ItemIdentity represents the identity of an item (its id plus its partition key value).
// This is useful for bulk/read-many style operations that need to address multiple
// items under (potentially) different partition key values.
//
// ID must match the 'id' property of the stored item. PartitionKey is the value (or
// composite/hierarchical set of values) the item was written with. For hierarchical
// partition keys create the PartitionKey with NewPartitionKey* helpers (e.g.
// NewPartitionKeyString, NewPartitionKeyInt, or NewPartitionKeyArray) following the
// order defined in the container. For hierarchical partition keys, all of the
// levels must be provided.
type ItemIdentity struct {
	ID           string       // Item id
	PartitionKey PartitionKey // Partition key value for the item
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
	var err error
	spanName, err := c.getSpanForContainer(operationTypeRead, resourceTypeCollection, c.id)
	if err != nil {
		return ContainerResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()

	response, err := c.readContainerRaw(ctx, o)
	if err == nil && c.database.client.getContainerCache() != nil && response.ContainerProperties != nil {
		// Populate the container properties cache on successful Read
		c.database.client.getContainerCache().set(c.link, response.ContainerProperties)
	}
	return response, err
}

// readContainerRaw performs the HTTP call to read container properties.
// It is the shared implementation used by both Read() and the container
// properties cache refresh, ensuring consistent request construction.
func (c *ContainerClient) readContainerRaw(
	ctx context.Context,
	o *ReadContainerOptions,
) (ContainerResponse, error) {
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
	var err error
	spanName, err := c.getSpanForContainer(operationTypeReplace, resourceTypeCollection, c.id)
	if err != nil {
		return ContainerResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newContainerResponse(azResponse)
	if err == nil && c.database.client.getContainerCache() != nil && response.ContainerProperties != nil {
		c.database.client.getContainerCache().set(c.link, response.ContainerProperties)
	}
	return response, err
}

// Delete a Cosmos container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *ContainerClient) Delete(
	ctx context.Context,
	o *DeleteContainerOptions) (ContainerResponse, error) {
	var err error
	spanName, err := c.getSpanForContainer(operationTypeDelete, resourceTypeCollection, c.id)
	if err != nil {
		return ContainerResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newContainerResponse(azResponse)
	return response, err
}

// ReadThroughput obtains the provisioned throughput information for the container.
// ctx - The context for the request.
// o - Options for the operation.
func (c *ContainerClient) ReadThroughput(
	ctx context.Context,
	o *ThroughputOptions) (ThroughputResponse, error) {
	var err error
	spanName, err := c.getSpanForContainer(operationTypeRead, resourceTypeOffer, c.id)
	if err != nil {
		return ThroughputResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
	if o == nil {
		o = &ThroughputOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{client: c.database.client}
	response, err := offers.ReadThroughputIfExists(ctx, rid, o)
	return response, err
}

// ReplaceThroughput updates the provisioned throughput for the container.
// ctx - The context for the request.
// throughputProperties - The throughput configuration of the container.
// o - Options for the operation.
func (c *ContainerClient) ReplaceThroughput(
	ctx context.Context,
	throughputProperties ThroughputProperties,
	o *ThroughputOptions) (ThroughputResponse, error) {
	var err error
	spanName, err := c.getSpanForContainer(operationTypeReplace, resourceTypeOffer, c.id)
	if err != nil {
		return ThroughputResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
	if o == nil {
		o = &ThroughputOptions{}
	}

	rid, err := c.getRID(ctx)
	if err != nil {
		return ThroughputResponse{}, err
	}

	offers := &cosmosOffers{client: c.database.client}
	response, err := offers.ReplaceThroughputIfExists(ctx, throughputProperties, rid, o)
	return response, err
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
	var err error
	spanName, err := c.getSpanForItems(operationTypeCreate)
	if err != nil {
		return ItemResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newItemResponse(azResponse)
	return response, err
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
	var err error
	spanName, err := c.getSpanForItems(operationTypeUpsert)
	if err != nil {
		return ItemResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newItemResponse(azResponse)
	return response, err
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
	var err error
	spanName, err := c.getSpanForItems(operationTypeReplace)
	if err != nil {
		return ItemResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newItemResponse(azResponse)
	return response, err
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
	var err error
	spanName, err := c.getSpanForItems(operationTypeRead)
	if err != nil {
		return ItemResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newItemResponse(azResponse)
	return response, err
}

// ReadManyItems reads multiple items in a Cosmos container. Note that the items returned in the response are unordered.
// ctx - The context for the request.
// itemIdentities - The identities of the items to read.
// o - Options for the operation.
func (c *ContainerClient) ReadManyItems(
	ctx context.Context,
	itemIdentities []ItemIdentity,
	o *ReadManyOptions) (ReadManyItemsResponse, error) {
	// if empty list of items, return empty list
	if len(itemIdentities) == 0 {
		return ReadManyItemsResponse{}, nil
	}

	// Validate all item IDs are non-empty
	for i := range itemIdentities {
		if itemIdentities[i].ID == "" {
			return ReadManyItemsResponse{}, errors.New("item identity at index " + fmt.Sprint(i) + " has an empty ID")
		}
	}

	readManyOptions := &ReadManyOptions{}
	if o != nil {
		originalOptions := *o
		readManyOptions = &originalOptions
	}

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDocument,
		resourceAddress: c.link,
	}

	ctx, endTrace := ensureOperationTrace(ctx, fmt.Sprintf("read_many_items %s", c.id))
	response, err := c.executeReadManyWithQueries(ctx, itemIdentities, readManyOptions, operationContext)
	endTrace()
	if err != nil {
		return ReadManyItemsResponse{}, err
	}

	response.Diagnostics = diagnosticsFromContext(ctx)
	return response, nil
}

// GetFeedRanges retrieves all the feed ranges for which changefeed could be fetched.
// ctx - The context for the request.
func (c *ContainerClient) GetFeedRanges(ctx context.Context) ([]FeedRange, error) {
	// Get the partition key ranges from the container
	response, err := c.getPartitionKeyRanges(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Convert partition key ranges to feed ranges
	feedRanges := make([]FeedRange, 0, len(response.PartitionKeyRanges))
	for _, pkr := range response.PartitionKeyRanges {
		feedRange := FeedRange{
			MinInclusive: pkr.MinInclusive,
			MaxExclusive: pkr.MaxExclusive,
		}
		feedRanges = append(feedRanges, feedRange)
	}

	return feedRanges, nil
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
	var err error
	spanName, err := c.getSpanForItems(operationTypeDelete)
	if err != nil {
		return ItemResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newItemResponse(azResponse)
	return response, err
}

// NewQueryItemsPager executes a single partition query in a Cosmos container.
// query - The SQL query to execute.
// partitionKey - The partition key to scope the query on. See below for more information on cross partition queries.
// o - Options for the operation.
//
// You can specify an empty list of partition keys by passing `NewPartitionKey()` to the `partitionKey` parameter, to indicate that the query WHERE clauses will specify which partitions to query.
//
// Limited cross partition queries ARE possible with the Go SDK.
// If you specify partition keys in the `partitionKey` parameter, you must specify ALL partition keys that the container has (in the case of hierarchical partitioning).
//
// If the query itself contains WHERE clauses that filter down to a single partition, the query will be executed on that partition.
// If the query does not filter down to a single partition (i.e. it does not filter on partition key at all, or filters on only some of the partition keys a container defines), the query will be executed as a cross partition query.
// The Azure Cosmos DB Gateway API, used by the Go SDK, can only perform a LIMITED set of cross-partition queries.
// Specifically, the gateway can only perform simple projections and filtering on cross partition queries.
// See https://learn.microsoft.com/rest/api/cosmos-db/querying-cosmosdb-resources-using-the-rest-api#queries-that-cannot-be-served-by-gateway for more details.
//
// When performing a cross-partition query, the Gateway may return pages of inconsistent size, or even empty pages (while still having a non-nil continuation token).
// Ensure you fully iterate the pager, even if you receive empty pages, to ensure you get all results.
//
// If you provide a query that the gateway cannot execute, it will return a BadRequest error.
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

	// For now, we short-cut straight to the preview query engine if provided.
	// In the future, we could consider running the normal pipeline until the Gateway fails due to an unsupported query and then switch over.
	// However, this logic could also just be handled in the query engine itself.
	if queryOptions.QueryEngine != nil {
		return c.executeQueryWithEngine(queryOptions.QueryEngine, query, queryOptions, operationContext)
	}

	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	return runtime.NewPager(runtime.PagingHandler[QueryItemsResponse]{
		More: func(page QueryItemsResponse) bool {
			return page.ContinuationToken != nil
		},
		Fetcher: func(ctx context.Context, page *QueryItemsResponse) (QueryItemsResponse, error) {
			var err error
			spanName, err := c.getSpanForItems(operationTypeQuery)
			if err != nil {
				return QueryItemsResponse{}, err
			}
			ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
			defer func() { endSpan(err) }()
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
	var err error
	spanName, err := c.getSpanForItems(operationTypePatch)
	if err != nil {
		return ItemResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newItemResponse(azResponse)
	return response, err
}

// NewTransactionalBatch creates a batch of operations to be committed as a single unit.
// See https://docs.microsoft.com/azure/cosmos-db/sql/transactional-batch
func (c *ContainerClient) NewTransactionalBatch(partitionKey PartitionKey) TransactionalBatch {
	return TransactionalBatch{partitionKey: partitionKey}
}

// ExecuteTransactionalBatch executes a transactional batch.
// Once executed, verify the Success property of the response to determine if the batch was committed
func (c *ContainerClient) ExecuteTransactionalBatch(ctx context.Context, b TransactionalBatch, o *TransactionalBatchOptions) (TransactionalBatchResponse, error) {
	var err error
	spanName, err := c.getSpanForContainer(operationTypeBatch, resourceTypeCollection, c.id)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()
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

	response, err := newTransactionalBatchResponse(azResponse)
	return response, err
}

// GetChangeFeed retrieves a single page of the change feed using the provided options.
// ctx - The context for the request.
// options - Options for the operation
//
// Resolution & retry behavior (overlap-match + 410 retry + multi-range queue):
//
//  1. If options.Continuation contains a multi-range composite continuation token, the
//     queued sub-ranges drive routing — options.FeedRange is ignored in favor of the
//     queue. The token's ResourceID is validated against the container's current
//     ResourceID; a mismatch surfaces a loud error rather than misrouting against the
//     wrong routing map.
//
//  2. If options.FeedRange is set (no token, or fresh start), the range is overlap-
//     matched against the current PK-range cache. A multi-overlap result (i.e., the
//     customer's range straddles a split) is expanded into one queue entry per child,
//     each inheriting the parent's ETag so no events are skipped at the boundary.
//
//  3. On 410/Gone with a PK-range substatus, the PK-range cache is refreshed and the
//     current queue head is re-resolved against the new routing map. Bounded at
//     maxPKRangeGoneRetries attempts; the final 410 is surfaced to the caller.
//
//  4. On 200 OK with documents, the head rotates to the tail with its new ETag and
//     the page is returned immediately so the caller can make progress.
//
//  5. On 304 Not Modified, the head rotates to the tail with any newly-issued ETag
//     and the next queue entry is tried. Drain bookkeeping ensures that if a queued
//     sub-range splits mid-drain (i.e., one head replacement adds N children), the
//     rotation budget grows so newly-inserted children get queried in this call.
//     If the whole queue drains with no documents, an empty page is returned with
//     the rotated token so the caller can poll again later.
//
// Returns ErrFeedRangeUnresolved (wrapped) when the customer's FeedRange/token
// doesn't overlap any current physical range even after a forced refresh — a
// signal to re-derive FeedRanges from GetFeedRanges.
//
// Returns an error wrapping *azcore.ResponseError on persistent 410/Gone or any
// non-retryable HTTP error.
func (c *ContainerClient) GetChangeFeed(
	ctx context.Context,
	options *ChangeFeedOptions,
) (ChangeFeedResponse, error) {
	if options == nil {
		options = &ChangeFeedOptions{}
	}

	var err error
	spanName, err := c.getSpanForItems(operationTypeRead)
	if err != nil {
		return ChangeFeedResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()

	// Cross-container token guard. We only need the container's current
	// ResourceID when a continuation token carries one to validate. Building
	// the initial queue takes care of that lazily so the no-token path stays
	// at one extra request (pk-ranges fetch).
	token, partitionKeyRanges, err := c.buildChangeFeedInitialQueue(ctx, options)
	if err != nil {
		return ChangeFeedResponse{}, err
	}

	// Capture into err so the deferred endSpan closure observes drain
	// failures (410 budget exhaustion, send errors, refresh failures,
	// serialization errors, ErrFeedRangeUnresolved) instead of recording
	// success on every drain that actually failed.
	var resp ChangeFeedResponse
	resp, err = c.getChangeFeedForQueue(ctx, options, token, partitionKeyRanges)
	return resp, err
}

// buildChangeFeedInitialQueue assembles the queue this call's drain loop will
// operate on, validating any provided continuation token against the current
// container and resolving any provided FeedRange against the current PK
// range cache. Returns the fetched PK-range snapshot alongside the token so
// the drain loop can reuse it without re-fetching on every iteration; the
// 410-retry path re-fetches on its own.
//
// Returns (token, snapshot, nil) on success. Returns ErrFeedRangeUnresolved
// (wrapped) when no overlap exists even after the cache is fresh.
func (c *ContainerClient) buildChangeFeedInitialQueue(
	ctx context.Context,
	options *ChangeFeedOptions,
) (*compositeContinuationToken, []partitionKeyRange, error) {
	// Path A: continuation token drives the queue.
	if options.Continuation != nil && *options.Continuation != "" {
		var compositeToken compositeContinuationToken
		if err := json.Unmarshal([]byte(*options.Continuation), &compositeToken); err == nil && len(compositeToken.Continuation) > 0 {
			// Reject cross-container token reuse loudly. Customers who hit this
			// have either pasted the wrong token, dropped a container and
			// recreated it under the same name, or fanned out a token to a
			// different client. Continuing would route against the wrong map.
			if compositeToken.ResourceID != "" {
				currentRID, ridErr := c.getContainerRID(ctx)
				if ridErr != nil {
					return nil, nil, ridErr
				}
				if currentRID != "" && compositeToken.ResourceID != currentRID {
					return nil, nil, fmt.Errorf(
						"continuation token ResourceID %q does not match the current container's ResourceID %q; the token was issued for a different container",
						compositeToken.ResourceID, currentRID,
					)
				}
			}
			queue := append([]changeFeedRange(nil), compositeToken.Continuation...)
			token := compositeToken
			token.Continuation = queue

			// Fetch a PK-range snapshot for the drain loop. Reused so the
			// loop doesn't issue an extra request per iteration.
			pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
			if err != nil {
				return nil, nil, err
			}
			return &token, pkrResp.PartitionKeyRanges, nil
		}
		// Not a composite token — fall through to FeedRange path; the legacy
		// ETag-only continuation is handled by buildRequestHeaders via the
		// queue head's ContinuationToken field.
	}

	// Path B: FeedRange drives the queue.
	if options.FeedRange == nil {
		return nil, nil, fmt.Errorf("GetChangeFeed requires a FeedRange to be set in the options, or a continuation token that contains a composite continuation token")
	}

	children, pkrs, err := c.resolveFeedRangeToChildren(ctx, *options.FeedRange)
	if err != nil {
		return nil, nil, err
	}
	entries := buildChildQueueEntries(children, nil)
	token := compositeContinuationToken{
		Version:      cosmosCompositeContinuationTokenVersion,
		Continuation: entries,
	}
	return &token, pkrs, nil
}

// resolveFeedRangeToChildren returns the routing-map ranges that overlap the
// given customer-supplied FeedRange. On no-overlap, performs a single forced
// refresh and retries; on still no overlap, returns ErrFeedRangeUnresolved.
//
// Also returns the PK-range snapshot it fetched, so the drain loop can reuse
// it for the rest of this GetChangeFeed call.
//
// Degraded fallback: when the PK-range cache is unavailable AND the direct
// fetch returns no ranges (e.g., test scaffolding), routing information is
// effectively missing. Rather than failing loudly, we drop any continuation
// context the caller may have provided, log a warning so the misroute is
// observable, and return a passthrough entry representing the customer's
// FeedRange (placed both in the children list and in the snapshot). The
// drain loop then issues the request fresh — without a PK-range-id header
// and without an If-None-Match ETag — and the server resolves routing.
func (c *ContainerClient) resolveFeedRangeToChildren(
	ctx context.Context,
	feedRange FeedRange,
) ([]partitionKeyRange, []partitionKeyRange, error) {
	pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	overlaps := overlappingPartitionKeyRanges(feedRange, pkrResp.PartitionKeyRanges)
	if len(overlaps) > 0 {
		return overlaps, pkrResp.PartitionKeyRanges, nil
	}

	// No overlap on the cached map. Try a forced refresh once if a cache
	// exists; if still no overlap, the customer's FeedRange genuinely doesn't
	// apply to this container.
	if c.database.client.getPKRangeCache() != nil {
		if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
			return nil, nil, refreshErr
		}
		pkrResp, err = c.getPartitionKeyRanges(ctx, nil)
		if err != nil {
			return nil, nil, err
		}
		overlaps = overlappingPartitionKeyRanges(feedRange, pkrResp.PartitionKeyRanges)
		if len(overlaps) > 0 {
			return overlaps, pkrResp.PartitionKeyRanges, nil
		}
		// Cache is fresh and still no overlap → unresolvable.
		return nil, nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}

	// No cache wired up. If we got any PK ranges in the direct fetch but none
	// overlap, that's a real customer mistake — same loud-fail semantics as
	// the cache path.
	if len(pkrResp.PartitionKeyRanges) > 0 {
		return nil, nil, &feedRangeUnresolvedError{feedRange: feedRange}
	}

	// Degraded fallback (cache absent AND direct fetch returned no ranges).
	// Log a warning, then issue the read fresh: passthrough range with no
	// continuation token. The passthrough is also placed in the snapshot so
	// the drain loop's overlap check is satisfied without re-entering this
	// branch.
	log.Writef(azlog.EventResponse,
		"azcosmos: GetChangeFeed: routing information unavailable for FeedRange [%s, %s); reading fresh without continuation token or PK-range header",
		feedRange.MinInclusive, feedRange.MaxExclusive,
	)
	passthrough := partitionKeyRange{
		MinInclusive: feedRange.MinInclusive,
		MaxExclusive: feedRange.MaxExclusive,
	}
	return []partitionKeyRange{passthrough}, []partitionKeyRange{passthrough}, nil
}

// overlappingPartitionKeyRanges returns the subset of partitionKeyRanges whose
// boundaries overlap the given feedRange. Order in the returned slice mirrors
// the input slice (the routing map is already sorted by MinInclusive).
func overlappingPartitionKeyRanges(feedRange FeedRange, partitionKeyRanges []partitionKeyRange) []partitionKeyRange {
	if len(partitionKeyRanges) == 0 {
		return nil
	}
	ids, err := findOverlappingPartitionKeyRangeIDs(feedRange, partitionKeyRanges)
	if err != nil || len(ids) == 0 {
		return nil
	}
	byID := make(map[string]partitionKeyRange, len(partitionKeyRanges))
	for _, r := range partitionKeyRanges {
		byID[r.ID] = r
	}
	out := make([]partitionKeyRange, 0, len(ids))
	for _, id := range ids {
		if r, ok := byID[id]; ok {
			out = append(out, r)
		}
	}
	return out
}

// buildChildQueueEntries materializes [].changeFeedRange entries for each
// child range, copying the inheritETag pointer onto every child so no events
// are skipped at the split boundary. inheritETag may be nil for fresh ranges
// that have never been read.
func buildChildQueueEntries(children []partitionKeyRange, inheritETag *azcore.ETag) []changeFeedRange {
	out := make([]changeFeedRange, 0, len(children))
	for _, ch := range children {
		entry := changeFeedRange{
			MinInclusive: ch.MinInclusive,
			MaxExclusive: ch.MaxExclusive,
		}
		if inheritETag != nil {
			etagCopy := *inheritETag
			entry.ContinuationToken = &etagCopy
		}
		out = append(out, entry)
	}
	return out
}

// getChangeFeedForQueue drains the queue, advancing on every response (200 or
// 304). On 200 with documents, returns immediately so the caller can process
// the page; on 304, rotates and tries the next entry until the original queue
// length is fully consumed (with budget bumps on splits). On 410, refreshes
// the cache, re-resolves the head, and retries — capped at maxPKRangeGoneRetries.
//
// partitionKeyRanges is the snapshot fetched once at the start of the call;
// the loop reuses it instead of re-fetching per iteration. The 410-retry path
// re-fetches and replaces the snapshot.
func (c *ContainerClient) getChangeFeedForQueue(
	ctx context.Context,
	options *ChangeFeedOptions,
	token *compositeContinuationToken,
	partitionKeyRanges []partitionKeyRange,
) (ChangeFeedResponse, error) {
	if token == nil || len(token.Continuation) == 0 {
		return ChangeFeedResponse{}, fmt.Errorf("GetChangeFeed has nothing to drain: no FeedRange and no continuation token entries")
	}

	// Drain budget: how many rotations we'll perform before we give up and
	// return an empty page so the caller can poll again. Starts at the queue
	// length and grows whenever a split-expansion inserts children.
	originalQueueLen := len(token.Continuation)
	rotations := 0
	pkRangeGoneAttempts := 0

	var lastResp ChangeFeedResponse

	for rotations < originalQueueLen {
		head := token.head()
		if head == nil {
			break
		}

		// Resolve the head's EPK range to a single PK-range ID against the
		// current routing-map snapshot.
		headFeedRange := FeedRange{MinInclusive: head.MinInclusive, MaxExclusive: head.MaxExclusive}
		overlaps := overlappingPartitionKeyRanges(headFeedRange, partitionKeyRanges)
		if len(overlaps) == 0 {
			// No overlap on the cached map. Force a refresh and re-fetch; if
			// still no overlap, the head is unresolvable. The cache-nil branch
			// is reachable only from test scaffolding (production always wires
			// up caches via acquireCaches) but we still attempt a direct fetch
			// fallback so handcrafted tests behave the same.
			if c.database.client.getPKRangeCache() != nil {
				if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
					return ChangeFeedResponse{}, refreshErr
				}
			}
			pkrResp, err := c.getPartitionKeyRanges(ctx, nil)
			if err != nil {
				return ChangeFeedResponse{}, err
			}
			partitionKeyRanges = pkrResp.PartitionKeyRanges
			overlaps = overlappingPartitionKeyRanges(headFeedRange, partitionKeyRanges)
			if len(overlaps) == 0 {
				// Degraded fallback: cache absent AND the direct fetch
				// returned nothing — routing information is unavailable for
				// this head. Drop the head's continuation token (an ETag
				// from a prior call could correspond to a now-defunct
				// physical range), log a warning so the misroute is
				// observable, and synthesize a passthrough so the request
				// is issued fresh — without a PK-range-id header and
				// without an If-None-Match ETag.
				if c.database.client.getPKRangeCache() == nil && len(partitionKeyRanges) == 0 {
					log.Writef(azlog.EventResponse,
						"azcosmos: GetChangeFeed: routing information unavailable for head [%s, %s); dropping continuation token and reading fresh",
						head.MinInclusive, head.MaxExclusive,
					)
					token.dropHeadContinuation()
					// head is a pointer into token.Continuation[0]; the
					// in-place mutation by dropHeadContinuation is already
					// visible through it. No reload needed.
					headFeedRange = FeedRange{MinInclusive: head.MinInclusive, MaxExclusive: head.MaxExclusive}
					passthrough := partitionKeyRange{
						MinInclusive: head.MinInclusive,
						MaxExclusive: head.MaxExclusive,
					}
					partitionKeyRanges = []partitionKeyRange{passthrough}
					overlaps = []partitionKeyRange{passthrough}
				} else {
					return ChangeFeedResponse{}, &feedRangeUnresolvedError{feedRange: headFeedRange}
				}
			}
		}

		var resolvedPKRangeID string
		if len(overlaps) > 1 {
			// Split-expansion. Replace the head with N children inheriting
			// the head's ETag, and bump the rotation budget so newly-inserted
			// children get visited in this call. Reset the 410 budget too —
			// each newly-inserted child is a fresh physical head and deserves
			// its own retry allowance.
			children := buildChildQueueEntries(overlaps, head.ContinuationToken)
			token.replaceHeadWithChildren(children)
			originalQueueLen += len(children) - 1
			pkRangeGoneAttempts = 0
			continue
		}
		resolvedPKRangeID = overlaps[0].ID

		headers, headerErr := options.buildRequestHeaders(*head, resolvedPKRangeID)
		if headerErr != nil {
			return ChangeFeedResponse{}, headerErr
		}

		addHeaders := func(r *policy.Request) {
			for k, v := range headers {
				r.Raw().Header.Set(k, v)
			}
		}

		operationContext := pipelineRequestOptions{
			resourceType:    resourceTypeDocument,
			resourceAddress: c.link,
		}
		path, err := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)
		if err != nil {
			return ChangeFeedResponse{}, err
		}

		azResponse, sendErr := c.database.client.sendGetRequest(
			path, ctx, operationContext, nil, addHeaders,
		)
		if sendErr != nil {
			// 410/Gone with a PK-range substatus → refresh + retry.
			if isPKRangeGoneResponseError(sendErr) {
				if pkRangeGoneAttempts >= maxPKRangeGoneRetries {
					return ChangeFeedResponse{}, sendErr
				}
				pkRangeGoneAttempts++
				if refreshErr := c.refreshPKRangeCache(ctx); refreshErr != nil {
					return ChangeFeedResponse{}, refreshErr
				}
				// Re-fetch the routing map after the cache was invalidated.
				pkrResp, fetchErr := c.getPartitionKeyRanges(ctx, nil)
				if fetchErr != nil {
					return ChangeFeedResponse{}, fetchErr
				}
				partitionKeyRanges = pkrResp.PartitionKeyRanges
				// Retry the same head against the refreshed snapshot.
				continue
			}
			return ChangeFeedResponse{}, sendErr
		}

		response, err := newChangeFeedResponse(azResponse)
		if err != nil {
			return response, err
		}

		// Capture the response body's _rid into the token's ResourceID on first
		// successful response. This keeps the cross-container guard meaningful
		// across resume — token-issued-by-this-container always carries the
		// container's RID — and matches pre-F1 PopulateCompositeContinuationToken
		// semantics that downstream tests rely on.
		if token.ResourceID == "" && response.ResourceID != "" {
			token.ResourceID = response.ResourceID
		}

		// Always rotate the head with the freshly-issued ETag, regardless of
		// status. This preserves drain progress even across 304s.
		newETag := response.ETag
		feedRangeForResp := &FeedRange{MinInclusive: head.MinInclusive, MaxExclusive: head.MaxExclusive}
		token.advance(newETag)
		rotations++
		// Head advanced to a new physical range; reset the 410 budget so the
		// next head gets its own allowance instead of inheriting prior 410s.
		pkRangeGoneAttempts = 0

		response.FeedRange = feedRangeForResp

		serialized, serErr := serializeCompositeContinuationToken(token)
		if serErr != nil {
			return response, serErr
		}
		response.ContinuationToken = serialized
		lastResp = response

		// 200 with documents → return immediately so the caller can process.
		if response.RawResponse != nil && response.RawResponse.StatusCode == http.StatusOK && response.Count > 0 {
			return response, nil
		}

		// 304 (or 200 with zero documents) → keep draining the rest of the queue.
	}

	// Whole queue drained without finding documents. Return the last (empty)
	// response with the rotated continuation token so the caller knows the
	// drain progressed and can poll again later.
	if lastResp.RawResponse == nil {
		// Nothing was issued (queue was empty). Synthesize an empty response.
		return ChangeFeedResponse{}, nil
	}
	return lastResp, nil
}

// serializeCompositeContinuationToken marshals the token as JSON for emission
// to the customer. Returns "" if the token is nil or has an empty queue.
func serializeCompositeContinuationToken(token *compositeContinuationToken) (string, error) {
	if token == nil || len(token.Continuation) == 0 {
		return "", nil
	}
	if token.Version == 0 {
		token.Version = cosmosCompositeContinuationTokenVersion
	}
	b, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *ContainerClient) getRID(ctx context.Context) (string, error) {
	containerResponse, err := c.Read(ctx, nil)
	if err != nil {
		return "", err
	}

	return containerResponse.ContainerProperties.ResourceID, nil
}

// getContainerRID resolves the container's ResourceID, using the container
// properties cache if available, otherwise falling back to a direct Read.
func (c *ContainerClient) getContainerRID(ctx context.Context) (string, error) {
	if c.database.client.getContainerCache() != nil {
		props, err := c.database.client.getContainerCache().getProperties(ctx, c)
		if err != nil {
			return "", err
		}
		return props.ResourceID, nil
	}
	return c.getRID(ctx)
}

func (c *ContainerClient) getSpanForContainer(operationType operationType, resourceType resourceType, id string) (span, error) {
	return getSpanNameForContainers(c.database.client.accountEndpointUrl(), operationType, resourceType, c.database.id, id)
}

func (c *ContainerClient) getSpanForItems(operationType operationType) (span, error) {
	return getSpanNameForItems(c.database.client.accountEndpointUrl(), operationType, c.database.id, c.id)
}

func (c *ContainerClient) getPartitionKeyRanges(ctx context.Context, o *partitionKeyRangeOptions) (partitionKeyRangeResponse, error) {
	var err error
	spanName, err := c.getSpanForContainer(operationTypeRead, resourceTypePartitionKeyRange, c.id)
	if err != nil {
		return partitionKeyRangeResponse{}, err
	}
	ctx, endSpan := startSpan(ctx, spanName.name, c.database.client.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()

	// Use the cache if available, otherwise fall back to direct fetch
	if c.database.client.getPKRangeCache() != nil {
		var containerRID string
		containerRID, err = c.getContainerRID(ctx)
		if err != nil {
			return partitionKeyRangeResponse{}, err
		}

		var routingMap *collectionRoutingMap
		routingMap, err = c.database.client.getPKRangeCache().getRoutingMap(ctx, containerRID, c.link, c.database.client)
		if err != nil {
			return partitionKeyRangeResponse{}, err
		}

		return partitionKeyRangeResponse{
			PartitionKeyRanges: routingMap.orderedRanges,
			Count:              len(routingMap.orderedRanges),
		}, nil
	}

	// Fallback: direct fetch without caching
	return c.fetchPartitionKeyRangesDirect(ctx, o)
}

// fetchPartitionKeyRangesDirect fetches partition key ranges directly from the service
// without using the cache.
func (c *ContainerClient) fetchPartitionKeyRangesDirect(ctx context.Context, o *partitionKeyRangeOptions) (partitionKeyRangeResponse, error) {
	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypePartitionKeyRange,
		resourceAddress: c.link,
	}

	if o == nil {
		o = &partitionKeyRangeOptions{}
	}

	path, err := generatePathForNameBased(resourceTypePartitionKeyRange, operationContext.resourceAddress, true)
	if err != nil {
		return partitionKeyRangeResponse{}, err
	}

	azResponse, err := c.database.client.sendGetRequest(
		path,
		ctx,
		operationContext,
		o,
		nil)
	if err != nil {
		return partitionKeyRangeResponse{}, err
	}

	response, err := newPartitionKeyRangeResponse(azResponse)
	if err != nil {
		return partitionKeyRangeResponse{}, err
	}
	return response, nil
}

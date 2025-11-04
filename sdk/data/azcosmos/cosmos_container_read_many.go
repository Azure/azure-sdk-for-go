// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// Executes a query using the provided query engine.
func (c *ContainerClient) executeReadManyWithEngine(queryEngine queryengine.QueryEngine, items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {
	path, _ := generatePathForNameBased(resourceTypeDocument, operationContext.resourceAddress, true)

	// if empty list of items, return empty list
	if len(items) == 0 {
		return ReadManyItemsResponse{}, nil
	}

	// get the partition key ranges for the container
	rawPartitionKeyRanges, err := c.getPartitionKeyRangesRaw(context.Background(), operationContext)
	if err != nil {
		// if we can't get the partition key ranges, return empty map
		return ReadManyItemsResponse{}, err
	}

	// get the container properties
	containerRsp, err := c.Read(nil, nil)
	if err != nil {
		return ReadManyItemsResponse{}, err
	}

	// create the item identities for the query engine
	newItemIdentities := make([]queryengine.ItemIdentity, len(items))
	for i := range items {
		pkStr, err := items[i].PartitionKey.toJsonString()
		if err != nil {
			return ReadManyItemsResponse{}, err
		}
		newItemIdentities[i] = queryengine.ItemIdentity{
			PartitionKeyValue: pkStr,
			ID:                items[i].ID,
		}
	}
	readManyPipeline, err := queryEngine.CreateReadManyPipeline(rawPartitionKeyRanges, newItemIdentities, string(containerRsp.ContainerProperties.PartitionKeyDefinition.Kind), int32(containerRsp.ContainerProperties.PartitionKeyDefinition.Version))
	if err != nil {
		return ReadManyItemsResponse{}, err
	}
	log.Writef(EventQueryEngine, "Created readMany pipeline")
	// Fetch more data from the pipeline
	log.Writef(EventQueryEngine, "Fetching more data from readMany pipeline")
	result, err := readManyPipeline.Run()
	if err != nil {
		readManyPipeline.Close()
		return ReadManyItemsResponse{}, err
	}

	// If we got items, we can return them, and we should do so now, to avoid making unnecessary requests.
	if len(result.Items) > 0 {
		log.Writef(EventQueryEngine, "ReadMany pipeline did not process any items", len(result.Items))
		return ReadManyItemsResponse{}, nil
	}

	// If we didn't have any items to return, we need to make requests for the items in the queue.
	// If there are no requests, the pipeline should return true for IsComplete, so we'll stop on the next iteration.
	for _, request := range result.Requests {
		log.Writef(azlog.EventRequest, "ReadMany pipeline requested data for PKRange: %s", request.PartitionKeyRangeID)
		// Make the single-partition query request
		qryRequest := queryRequest(request) // Cast to our type, which has toHeaders defined on it.
		azResponse, err := c.database.client.sendQueryRequest(
			path,
			ctx,
			query,
			nil,
			operationContext,
			&qryRequest,
			nil)
		if err != nil {
			readManyPipeline.Close()
			return ReadManyItemsResponse{}, err
		}
		lastResponse = newResponse(azResponse)

		// Load the data into a buffer to send it to the pipeline
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(azResponse.Body)
		if err != nil {
			readManyPipeline.Close()
			return ReadManyItemsResponse{}, err
		}
		data := buf.Bytes()
		continuation := azResponse.Header.Get(cosmosHeaderContinuationToken)

		// Provide the data to the pipeline, make sure it's tagged with the partition key range ID so the pipeline can merge it into the correct partition.
		result := queryengine.QueryResult{
			PartitionKeyRangeID: request.PartitionKeyRangeID,
			NextContinuation:    continuation,
			Data:                data,
		}
		log.Writef(EventQueryEngine, "Received response for PKRange: %s.", request.PartitionKeyRangeID)
		if err = readManyPipeline.ProvideData(result); err != nil {
			readManyPipeline.Close()
			return ReadManyItemsResponse{}, err
		}
	}

	if readManyPipeline.IsComplete() {
		log.Writef(EventQueryEngine, "ReadMany pipeline is complete")
		readManyPipeline.Close()
		return ReadManyItemsResponse{
			Response: lastResponse,
			Items:    nil,
		}, nil
	}
}

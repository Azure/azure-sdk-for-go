// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
)

// Executes a query using the provided query engine.
func (c *ContainerClient) executeReadManyWithEngine(queryEngine queryengine.QueryEngine, items []ItemIdentity, queryOptions *QueryOptions, operationContext pipelineRequestOptions) []ReadManyItemsResponse {
	// NOTE: The current interface for runtime.Pager means we're probably going to risk leaking the pipeline, if it's provided by a native query engine.
	// There's no "Close" method, which means we can't call `queryengine.QueryPipeline.Close()` when we're done.
	// We _do_ close the pipeline if the user iterates the entire pager, but if they don't we don't have a way to clean up.
	// To mitigate that, we expect the queryengine.QueryPipeline to handle setting up a Go finalizer to clean up any native resources it holds.
	// Finalizers aren't deterministic though, so we should consider making the pager "closable" in the future, so we have a clear signal to free the native resources.

	// if empty list of items, return empty list
	if len(items) == 0 {
		return []ReadManyItemsResponse{}
	}

	// get the partition key ranges for the container
	rawPartitionKeyRanges, err := c.getPartitionKeyRangesRaw(context.Background(), nil)
	if err != nil {
		// if we can't get the partition key ranges, return empty map
		return nil, errors.New("failed to get partition key ranges: " + err.Error())
	}

	// get the container properties
	containeRsp, err := c.Read(nil, nil)
	if err != nil {
		return nil, errors.New("failed to get container properties: " + err.Error())
	}

	// call client engine here to group the partition key ranges
	// create query chunks for each physical partition
	queries := queryEngine.createQueryChunks(rawPartitionKeyRanges, items, containeRsp.ContainerProperties.PartitionKeyDefinition.Kind, containeRsp.ContainerProperties.PartitionKeyDefinition.Version)

	// execute queries concurrently

}

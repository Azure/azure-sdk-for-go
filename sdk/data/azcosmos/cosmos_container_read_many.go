// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
)

// executeReadManyWithEngine executes a query using the provided query engine.
func (c *ContainerClient) executeReadManyWithEngine(queryEngine queryengine.QueryEngine, items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {
	// throw error that this is unsupported
	return ReadManyItemsResponse{}, errors.New("ReadMany with query engine is not supported yet.")
}

func (c *ContainerClient) executeReadManyWithPointReads(items []ItemIdentity, readManyOptions *ReadManyOptions, operationContext pipelineRequestOptions, ctx context.Context) (ReadManyItemsResponse, error) {

	// if empty list of items, return empty list
	if len(items) == 0 {
		return ReadManyItemsResponse{}, nil
	}
	var readManyResponse ReadManyItemsResponse
	for _, item := range items {
		itemOptions := ItemOptions{ConsistencyLevel: readManyOptions.ConsistencyLevel, SessionToken: readManyOptions.SessionToken}
		ItemResponse, err := c.ReadItem(ctx, item.PartitionKey, item.ID, &itemOptions)
		if err != nil {
			// on an error, bail out and return the error
			return ReadManyItemsResponse{}, err
		}
		// Append the item response to the list of items
		readManyResponse.Items = append(readManyResponse.Items, ItemResponse.Value)
		readManyResponse.RequestCharge += ItemResponse.RequestCharge
	}

	return readManyResponse, nil
}

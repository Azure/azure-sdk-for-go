// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ItemIdentity represents the identity (ID and partition key) of an item.
// This is useful for bulk/read-many style operations that need to address multiple
// items under (potentially) different partition key values.
//
// ID must match the 'id' property of the stored item. PartitionKey is the value (or
// composite/hierarchical set of values) the item was written with. For hierarchical
// partition keys create the PartitionKey with NewPartitionKey* helpers (e.g.
// NewPartitionKeyString, NewPartitionKeyInt, or NewPartitionKeyArray) following the
// order defined in the container.
type ItemIdentity struct {
	ID           string
	PartitionKey PartitionKey
}

// ReadManyOptions contains options for ReadMany operations.
type ReadManyOptions struct {
	// SessionToken to be used when using Session consistency on the account.
	// When working with Session consistency, each new write request to Azure Cosmos DB is assigned a new SessionToken.
	// The client instance will use this token internally with each read/query request to ensure that the set consistency level is maintained.
	SessionToken *string
	// ConsistencyLevel overrides the account defined consistency level for this operation.
	// Consistency can only be relaxed.
	ConsistencyLevel *ConsistencyLevel
	// Options for operations in the dedicated gateway.
	DedicatedGatewayRequestOptions *DedicatedGatewayRequestOptions
}

// ReadManyResponse represents the response from a ReadMany operation.
type ReadManyResponse struct {
	Items              [][]byte
	TotalRequestCharge float32
}

// ReadMany reads multiple items from a Cosmos container.
// This operation may be equivalent to performing a sequence of Read operations for each item,
// but the SDK will apply optimizations when possible. In the worst case, this is no worse
// than sequential point reads, but may benefit from current and future optimizations.
// Items that do not exist in the container will be silently skipped and not included in the response.
// ctx - The context for the request.
// itemIdentities - The identities (ID and partition key) of the items to read.
// options - Options for the operation.
func (c *ContainerClient) ReadMany(
	ctx context.Context,
	itemIdentities []ItemIdentity,
	options *ReadManyOptions) (ReadManyResponse, error) {
	if options == nil {
		options = &ReadManyOptions{}
	}

	items := make([][]byte, 0, len(itemIdentities))
	var totalCharge float32

	for _, identity := range itemIdentities {
		itemOptions := &ItemOptions{
			SessionToken:                   options.SessionToken,
			ConsistencyLevel:               options.ConsistencyLevel,
			DedicatedGatewayRequestOptions: options.DedicatedGatewayRequestOptions,
		}
		response, err := c.ReadItem(ctx, identity.PartitionKey, identity.ID, itemOptions)
		if err != nil {
			var responseErr *azcore.ResponseError
			if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusNotFound {
				continue
			}
			return ReadManyResponse{}, err
		}
		items = append(items, response.Value)
		totalCharge += response.RequestCharge
	}

	return ReadManyResponse{
		Items:              items,
		TotalRequestCharge: totalCharge,
	}, nil
}

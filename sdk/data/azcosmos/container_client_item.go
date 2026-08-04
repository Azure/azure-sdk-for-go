// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ReadItemOptions configures [ContainerClient.ReadItem]. A nil *ReadItemOptions selects the
// defaults for every field.
type ReadItemOptions struct {
	// ConsistencyStrategy relaxes how fresh the read must be. The zero value reads with whatever
	// the account's consistency level implies.
	ConsistencyStrategy ReadConsistencyStrategy

	// IfNoneMatchETag skips returning the item when its ETag matches, so an unchanged item costs
	// no payload. Pass the ETag from a previous response.
	IfNoneMatchETag *azcore.ETag

	// SessionToken is the session token to read under, for observing writes made by another
	// process. Empty uses the token the client captured itself.
	SessionToken string

	// ExcludedRegions removes regions from consideration for this operation, in addition to any
	// the client is already avoiding.
	ExcludedRegions []string

	// ThroughputControlGroup assigns the operation to a throughput control group registered on
	// the client.
	ThroughputControlGroup string

	// EndToEndTimeout bounds the whole operation, including the retries the driver performs on
	// the caller's behalf. Zero means no bound beyond the context's deadline.
	EndToEndTimeout time.Duration
}

// CreateItemOptions configures [ContainerClient.CreateItem]. A nil *CreateItemOptions selects the
// defaults for every field.
type CreateItemOptions struct {
	// EnableContentResponseOnWrite requests that the created item be returned in the response.
	// Nil uses the client-level setting, and a non-nil value overrides it in either direction.
	// Leaving it off reduces network and CPU cost.
	EnableContentResponseOnWrite *bool

	// SessionToken is the session token to write under. Empty uses the token the client captured
	// itself.
	SessionToken string

	// ExcludedRegions removes regions from consideration for this operation, in addition to any
	// the client is already avoiding.
	ExcludedRegions []string

	// ThroughputControlGroup assigns the operation to a throughput control group registered on
	// the client.
	ThroughputControlGroup string

	// EndToEndTimeout bounds the whole operation, including the retries the driver performs on
	// the caller's behalf. Zero means no bound beyond the context's deadline.
	EndToEndTimeout time.Duration
}

// ReadItem reads a single item.
//
// partitionKey is the item's partition key value, and must have one component per path in the
// container's partition key definition. itemID is the item's id property. options may be nil.
//
// When the item does not exist the returned error is an [Error] with [Error.Code] set to
// [CodeNotFound]:
//
//	response, err := container.ReadItem(ctx, pk, "item-id", nil)
//	var cosmosErr *azcosmos.Error
//	if errors.As(err, &cosmosErr) && cosmosErr.Code == azcosmos.CodeNotFound {
//		// no such item
//	}
func (c *ContainerClient) ReadItem(ctx context.Context, partitionKey PartitionKey, itemID string, options *ReadItemOptions) (ItemResponse, error) {
	if err := validateItemArguments(partitionKey); err != nil {
		return ItemResponse{}, err
	}
	if itemID == "" {
		return ItemResponse{}, errors.New("azcosmos: item id must not be empty")
	}
	if err := ctx.Err(); err != nil {
		return ItemResponse{}, err
	}

	release, err := c.database.client.acquire()
	if err != nil {
		return ItemResponse{}, err
	}
	defer release()

	_ = options
	return ItemResponse{}, errNotImplemented
}

// CreateItem creates a new item, failing if one with the same id already exists in the partition.
//
// partitionKey is the item's partition key value, and must have one component per path in the
// container's partition key definition and match the values in item. item is the JSON encoding of
// the item, which must include an id property. options may be nil.
//
// When an item with the same id already exists the returned error is an [Error] with [Error.Code]
// set to [CodeConflict]. The response carries the created item only when content responses are
// enabled, on the client or through [CreateItemOptions.EnableContentResponseOnWrite].
func (c *ContainerClient) CreateItem(ctx context.Context, partitionKey PartitionKey, item []byte, options *CreateItemOptions) (ItemResponse, error) {
	if err := validateItemArguments(partitionKey); err != nil {
		return ItemResponse{}, err
	}
	if len(item) == 0 {
		return ItemResponse{}, errors.New("azcosmos: item must not be empty")
	}
	if err := ctx.Err(); err != nil {
		return ItemResponse{}, err
	}

	release, err := c.database.client.acquire()
	if err != nil {
		return ItemResponse{}, err
	}
	defer release()

	_ = options
	return ItemResponse{}, errNotImplemented
}

// validateItemArguments rejects a partition key with no components. No container has a partition
// key definition with zero paths, so an empty value is always a caller mistake — most often a
// PartitionKey that was declared but never appended to.
func validateItemArguments(partitionKey PartitionKey) error {
	if partitionKey.Len() == 0 {
		return errors.New("azcosmos: partition key must have at least one component")
	}
	return nil
}

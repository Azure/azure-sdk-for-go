// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// cSpell:ignore azsdk upserted

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ReadItemOptions configures [ContainerClient.ReadItem]. A nil *ReadItemOptions selects the
// defaults for every field.
type ReadItemOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions

	// SessionToken is the session token to read under, for observing writes made by another
	// process. Empty uses the token the client captured itself.
	SessionToken SessionToken

	// IfNoneMatchETag skips returning the item when its ETag matches, so an unchanged item costs
	// no payload. Pass the ETag from a previous response.
	IfNoneMatchETag *azcore.ETag
}

// CreateItemOptions configures [ContainerClient.CreateItem]. A nil *CreateItemOptions selects the
// defaults for every field.
type CreateItemOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions

	// SessionToken is the session token to write under. Empty uses the token the client captured
	// itself.
	SessionToken SessionToken
}

// ReplaceItemOptions configures [ContainerClient.ReplaceItem]. A nil *ReplaceItemOptions selects
// the defaults for every field.
type ReplaceItemOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions

	// SessionToken is the session token to write under. Empty uses the token the client captured
	// itself.
	SessionToken SessionToken

	// IfMatchETag makes the replacement conditional on the item still having this ETag.
	IfMatchETag *azcore.ETag
}

// UpsertItemOptions configures [ContainerClient.UpsertItem]. A nil *UpsertItemOptions selects the
// defaults for every field.
type UpsertItemOptions struct {
	// Operation holds the settings every operation accepts.
	Operation OperationOptions

	// SessionToken is the session token to write under. Empty uses the token the client captured
	// itself.
	SessionToken SessionToken

	// IfMatchETag applies an If-Match precondition. A mismatching ETag fails when the item exists;
	// when the item does not exist the upsert creates it.
	IfMatchETag *azcore.ETag
}

// DeleteItemOptions configures [ContainerClient.DeleteItem]. A nil *DeleteItemOptions selects the
// defaults for every field.
type DeleteItemOptions struct {
	// Operation holds the settings every operation accepts. Content-response settings have no
	// effect because DeleteItem never returns an item body.
	Operation OperationOptions

	// SessionToken is the session token to write under. Empty uses the token the client captured
	// itself.
	SessionToken SessionToken

	// IfMatchETag makes the deletion conditional on the item still having this ETag.
	IfMatchETag *azcore.ETag
}

// PatchItemOptions configures [ContainerClient.PatchItem]. A nil *PatchItemOptions selects the
// defaults for every field.
//
// PatchItemOptions is provisional. It may change or be removed before azcosmos/v2 reaches a stable
// release.
type PatchItemOptions struct {
	// Operation holds the settings every operation accepts. PatchItem returns the updated item by
	// default unless content responses are explicitly disabled at the client or operation level.
	Operation OperationOptions

	// SessionToken is the session token to write under. Empty uses the token the client captured
	// itself.
	SessionToken SessionToken

	// IfMatchETag makes the patch conditional on the item still having this ETag.
	IfMatchETag *azcore.ETag

	// Strategy selects how the driver executes the patch. The zero value inherits the strategy
	// configured for the driver, runtime, or environment.
	Strategy PatchStrategy
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
	if options != nil {
		if err := options.Operation.ConsistencyStrategy.validate(); err != nil {
			return ItemResponse{}, err
		}
		if err := options.SessionToken.validate(); err != nil {
			return ItemResponse{}, err
		}
		if options.IfNoneMatchETag != nil {
			etag := string(*options.IfNoneMatchETag)
			if etag == "" {
				return ItemResponse{}, errors.New("azcosmos: IfNoneMatchETag must not be empty")
			}
			if strings.IndexByte(etag, 0) >= 0 {
				return ItemResponse{}, errors.New(
					"azcosmos: IfNoneMatchETag must not contain a NUL byte",
				)
			}
		}
	}
	req := itemRequest{
		kind:         operationKindReadItem,
		databaseID:   c.database.id,
		containerID:  c.id,
		itemID:       itemID,
		partitionKey: partitionKey,
	}
	if options != nil {
		req.options = options.Operation
		req.sessionToken = options.SessionToken
		if options.IfNoneMatchETag != nil {
			req.preconditionKind = preconditionKindIfNoneMatch
			req.preconditionETag = string(*options.IfNoneMatchETag)
		}
	}

	return c.executeItem(ctx, req, true)
}

// CreateItem creates a new item, failing if one with the same id already exists in the partition.
//
// partitionKey is the item's partition key value, and must have one component per path in the
// container's partition key definition and match the values in item. id addresses the item being
// created and must match the item's own id property. item is the JSON encoding of the item.
// options may be nil.
//
// The id is taken separately because the driver ABI requires an item reference. The SDK passes the
// id and item through unchanged and does not parse the payload to compare them; callers must provide
// matching values, and service behavior is undefined when they differ.
//
// When an item with the same id already exists the returned error is an [Error] with [Error.Code]
// set to [CodeConflict]. The response carries the created item only when content responses are
// enabled, on the client through [ClientOptions.EnableContentResponseOnWrite] or per operation
// through [OperationOptions.EnableContentResponseOnWrite].
func (c *ContainerClient) CreateItem(ctx context.Context, partitionKey PartitionKey, id string, item []byte, options *CreateItemOptions) (ItemResponse, error) {
	if err := validateItemArguments(partitionKey); err != nil {
		return ItemResponse{}, err
	}
	if id == "" {
		return ItemResponse{}, errors.New("azcosmos: item id must not be empty")
	}
	if len(item) == 0 {
		return ItemResponse{}, errors.New("azcosmos: item must not be empty")
	}
	if options != nil {
		if err := options.SessionToken.validate(); err != nil {
			return ItemResponse{}, err
		}
	}
	req := itemRequest{
		kind:         operationKindCreateItem,
		databaseID:   c.database.id,
		containerID:  c.id,
		itemID:       id,
		partitionKey: partitionKey,
		body:         item,
	}
	if options != nil {
		req.options = options.Operation
		req.sessionToken = options.SessionToken
	}

	return c.executeItem(ctx, req, true)
}

// ReplaceItem replaces an existing item.
//
// partitionKey is the item's partition key value. id is the item's id property and must match the
// id in item. item is the JSON encoding of the replacement. options may be nil. The SDK passes id
// and item through unchanged and does not parse the payload to compare them.
//
// The response carries the replaced item only when content responses are enabled. If IfMatchETag
// does not match, the returned error has [Error.Code] set to [CodePreconditionFailed].
func (c *ContainerClient) ReplaceItem(ctx context.Context, partitionKey PartitionKey, id string, item []byte, options *ReplaceItemOptions) (ItemResponse, error) {
	if err := validateItemWriteArguments(partitionKey, id, item, true); err != nil {
		return ItemResponse{}, err
	}
	if options != nil {
		if err := validateItemWriteOptions(options.SessionToken, options.IfMatchETag); err != nil {
			return ItemResponse{}, err
		}
	}

	req := itemRequest{
		kind:         operationKindReplaceItem,
		databaseID:   c.database.id,
		containerID:  c.id,
		itemID:       id,
		partitionKey: partitionKey,
		body:         item,
	}
	if options != nil {
		req.options = options.Operation
		req.sessionToken = options.SessionToken
		setIfMatchPrecondition(&req, options.IfMatchETag)
	}
	return c.executeItem(ctx, req, true)
}

// UpsertItem creates an item when it does not exist or replaces it when it does.
//
// partitionKey is the item's partition key value. id is the item's id property and must match the
// id in item. item is the JSON encoding of the item. options may be nil. The SDK passes id and item
// through unchanged and does not parse the payload to compare them.
//
// The response carries the upserted item only when content responses are enabled. IfMatchETag is
// evaluated when the item exists; a missing item is created even when IfMatchETag is set.
func (c *ContainerClient) UpsertItem(ctx context.Context, partitionKey PartitionKey, id string, item []byte, options *UpsertItemOptions) (ItemResponse, error) {
	if err := validateItemWriteArguments(partitionKey, id, item, true); err != nil {
		return ItemResponse{}, err
	}
	if options != nil {
		if err := validateItemWriteOptions(options.SessionToken, options.IfMatchETag); err != nil {
			return ItemResponse{}, err
		}
	}

	req := itemRequest{
		kind:         operationKindUpsertItem,
		databaseID:   c.database.id,
		containerID:  c.id,
		itemID:       id,
		partitionKey: partitionKey,
		body:         item,
	}
	if options != nil {
		req.options = options.Operation
		req.sessionToken = options.SessionToken
		setIfMatchPrecondition(&req, options.IfMatchETag)
	}
	return c.executeItem(ctx, req, true)
}

// DeleteItem deletes an item.
//
// partitionKey is the item's partition key value and id is its id property. options may be nil.
// The returned ItemResponse never carries a Value. If IfMatchETag does not match, the returned
// error has [Error.Code] set to [CodePreconditionFailed].
func (c *ContainerClient) DeleteItem(ctx context.Context, partitionKey PartitionKey, id string, options *DeleteItemOptions) (ItemResponse, error) {
	if err := validateItemWriteArguments(partitionKey, id, nil, false); err != nil {
		return ItemResponse{}, err
	}
	if options != nil {
		if err := validateItemWriteOptions(options.SessionToken, options.IfMatchETag); err != nil {
			return ItemResponse{}, err
		}
	}

	req := itemRequest{
		kind:         operationKindDeleteItem,
		databaseID:   c.database.id,
		containerID:  c.id,
		itemID:       id,
		partitionKey: partitionKey,
	}
	if options != nil {
		req.options = options.Operation
		req.sessionToken = options.SessionToken
		setIfMatchPrecondition(&req, options.IfMatchETag)
	}
	return c.executeItem(ctx, req, false)
}

// PatchItem applies an ordered set of partial updates to an existing item.
//
// PatchItem is provisional. It may change or be removed before azcosmos/v2 reaches a stable
// release.
//
// partitionKey is the item's partition key value and id is its id property. operations must
// contain at least one operation. options may be nil. PatchItem returns the updated item by default
// unless content responses are explicitly disabled.
//
// The driver chooses server-side PATCH or client-side read-modify-write execution by default.
// options can select an explicit strategy. Client-side execution can permanently add the internal
// _azsdkPatchTracking property described by [PatchOperations].
func (c *ContainerClient) PatchItem(ctx context.Context, partitionKey PartitionKey, id string, operations PatchOperations, options *PatchItemOptions) (ItemResponse, error) {
	if err := validateItemWriteArguments(partitionKey, id, nil, false); err != nil {
		return ItemResponse{}, err
	}
	body, err := operations.marshal()
	if err != nil {
		return ItemResponse{}, err
	}
	if options != nil {
		if err := validateItemWriteOptions(options.SessionToken, options.IfMatchETag); err != nil {
			return ItemResponse{}, err
		}
		if err := options.Strategy.validate(); err != nil {
			return ItemResponse{}, err
		}
	}

	req := itemRequest{
		kind:         operationKindPatchItem,
		databaseID:   c.database.id,
		containerID:  c.id,
		itemID:       id,
		partitionKey: partitionKey,
		body:         body,
	}
	if options != nil {
		req.options = options.Operation
		req.sessionToken = options.SessionToken
		req.patchStrategy = options.Strategy
		setIfMatchPrecondition(&req, options.IfMatchETag)
	}
	return c.executeItem(ctx, req, true)
}

func (c *ContainerClient) executeItem(ctx context.Context, req itemRequest, includeValue bool) (ItemResponse, error) {
	if c.database.client.beforeItemAcquire != nil {
		c.database.client.beforeItemAcquire()
	}
	// The client is consulted before the context so that a call made after Close reports
	// CodeClientClosed, which Close guarantees, rather than whatever the caller's context happens
	// to say. Shutdown is exactly when both are likely to be true at once.
	release, err := c.database.client.acquire()
	if err != nil {
		return ItemResponse{}, err
	}
	defer release()

	if err := ctx.Err(); err != nil {
		return ItemResponse{}, err
	}

	response, body, err := c.database.client.execute(ctx, req)
	if err != nil {
		return ItemResponse{}, normalizeItemOperationError(req.kind, err)
	}
	if includeValue {
		response.Value = body
	}
	return response, nil
}

func normalizeItemOperationError(kind operationKind, err error) error {
	// Automatic PATCH can enforce If-Match during its client-side read-modify-write path. The
	// published driver reports that as a local HTTP 412 without a synthetic sub-status, so retain
	// FromWire=false while giving callers the same CodePreconditionFailed classification as a
	// server-side PATCH.
	cosmosErr, ok := err.(*Error)
	if kind != operationKindPatchItem || !ok || cosmosErr.FromWire ||
		cosmosErr.StatusCode != 412 || cosmosErr.Code != CodeClientError {
		return err
	}
	normalized := cloneError(cosmosErr).(*Error)
	normalized.Code = CodePreconditionFailed
	return normalized
}

func validateItemWriteArguments(partitionKey PartitionKey, id string, item []byte, requireItem bool) error {
	if err := validateItemArguments(partitionKey); err != nil {
		return err
	}
	if id == "" {
		return errors.New("azcosmos: item id must not be empty")
	}
	if requireItem && len(item) == 0 {
		return errors.New("azcosmos: item must not be empty")
	}
	return nil
}

func validateItemWriteOptions(sessionToken SessionToken, ifMatchETag *azcore.ETag) error {
	if err := sessionToken.validate(); err != nil {
		return err
	}
	return validateETag("IfMatchETag", ifMatchETag)
}

func validateETag(name string, etag *azcore.ETag) error {
	if etag == nil {
		return nil
	}
	value := string(*etag)
	if value == "" {
		return fmt.Errorf("azcosmos: %s must not be empty", name)
	}
	if strings.IndexByte(value, 0) >= 0 {
		return fmt.Errorf("azcosmos: %s must not contain a NUL byte", name)
	}
	return nil
}

func setIfMatchPrecondition(req *itemRequest, etag *azcore.ETag) {
	if etag == nil {
		return
	}
	req.preconditionKind = preconditionKindIfMatch
	req.preconditionETag = string(*etag)
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

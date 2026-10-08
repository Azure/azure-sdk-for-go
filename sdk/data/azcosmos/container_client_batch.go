// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"slices"
)

// TransactionalBatchOptions configures [ContainerClient.ExecuteTransactionalBatch].
// A nil options pointer selects defaults.
type TransactionalBatchOptions struct {
	// Operation contains request-wide native settings. Content-response settings apply to
	// write operations, not reads. Nil inherits the client/native default; false suppresses
	// write bodies, and true requests them.
	Operation OperationOptions

	// SessionToken supplies the session for this request. Empty uses the client's captured token.
	SessionToken SessionToken
}

// ExecuteTransactionalBatch executes batch as one native transaction in a complete logical partition.
//
// The batch must contain 1-100 operations and its encoded operation envelope must not exceed 2 MiB.
// A complete hierarchical partition key is supported; a prefix is rejected before batch submission.
// Validation may read container metadata, whose request charge is included in the response/error.
// Errors can contain setup metadata even when the batch itself was never submitted.
// Item bodies must match the batch's partition key; the driver/service validates their fields.
//
// A received rollback response is not an execution error: inspect [TransactionalBatchResponse.Success]
// and [TransactionalBatchResponse.FailedOperationIndex]. Native admission, execution, and decoding
// failures return the zero response and an error; a decoding or transport error is not proof that
// writes did not commit. Error metadata includes the native body and status when available.
//
// After admission, this method waits for the authoritative native write outcome, even if ctx ends.
// The released native ABI cannot cancel an admitted operation. [Client.Close] waits for the operation
// before releasing its resources. options may be nil. Do not mutate batch while it is executing.
func (c *ContainerClient) ExecuteTransactionalBatch(ctx context.Context, batch TransactionalBatch, options *TransactionalBatchOptions) (TransactionalBatchResponse, error) {
	req, err := newTransactionalBatchRequest(batch, options)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	req.databaseID, req.containerID = c.database.id, c.id
	release, err := c.database.client.acquire()
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return TransactionalBatchResponse{}, err
	}
	return c.database.client.executeBatch(ctx, req, len(batch.operations))
}

func newTransactionalBatchRequest(batch TransactionalBatch, options *TransactionalBatchOptions) (itemRequest, error) {
	body, err := batch.marshal()
	if err != nil {
		return itemRequest{}, err
	}
	req := itemRequest{
		kind:         operationKindBatch,
		partitionKey: batch.partitionKey,
		body:         body,
	}
	if options != nil {
		if err := options.Operation.ConsistencyStrategy.validate(); err != nil {
			return itemRequest{}, err
		}
		if err := options.SessionToken.validate(); err != nil {
			return itemRequest{}, err
		}
		req.options = options.Operation
		req.options.ExcludedRegions = slices.Clone(options.Operation.ExcludedRegions)
		if enabled := options.Operation.EnableContentResponseOnWrite; enabled != nil {
			copied := *enabled
			req.options.EnableContentResponseOnWrite = &copied
		}
		req.sessionToken = options.SessionToken
	}
	return req, nil
}

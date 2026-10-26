// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
)

// TransactionalBatchOptions configures [ContainerClient.ExecuteTransactionalBatch].
// A nil options pointer selects defaults.
type TransactionalBatchOptions struct {
	// Operation contains request-wide native settings. EnableContentResponseOnWrite inherits
	// runtime/client/native defaults when nil. With native driver v0.2.0, explicit or inherited
	// false can also suppress HTTP 200 read bodies. Set it explicitly to true when read bodies
	// are required; this also requests write bodies.
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
// A received response can contain non-2xx operation statuses without an execution error. Inspect
// [TransactionalBatchResponse.OperationResults]; HTTP 207 or a conditional read's HTTP 304 is not
// interpreted as a commitment or rollback verdict. Native admission, execution, and decoding failures
// return the zero response and an error; a decoding or transport error is not proof that writes did
// not commit. Error metadata includes the native body and status when available.
//
// Context cancellation stops waiting, not native execution; an admitted batch may still commit.
// Allow a context and native end-to-end timeout budget appropriate for the complete batch,
// including setup and retries. A longer budget does not guarantee a known commit outcome.
// A native completion already available when cancellation is observed is returned unchanged.
// Otherwise this method returns the zero response and a cancellation error, and native work drains
// in the background. [Client.Close] waits for that work before releasing its resources.
// ctx must not be nil; options may be nil.
// Do not mutate batch while it is executing.
func (c *ContainerClient) ExecuteTransactionalBatch(ctx context.Context, batch TransactionalBatch, options *TransactionalBatchOptions) (TransactionalBatchResponse, error) {
	req, err := newTransactionalBatchRequest(batch, options)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	if ctx == nil {
		return TransactionalBatchResponse{}, errors.New("azcosmos: context must not be nil")
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
		if err := options.Operation.validate(); err != nil {
			return itemRequest{}, err
		}
		if err := options.SessionToken.validate(); err != nil {
			return itemRequest{}, err
		}
		req.options = options.Operation.clone()
		req.sessionToken = options.SessionToken
	}
	return req, nil
}

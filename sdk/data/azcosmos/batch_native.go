// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore gocritic

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import "context"

func (c *Client) executeBatch(ctx context.Context, req itemRequest, operationCount int) (TransactionalBatchResponse, error) {
	d := c.driver
	ctx, snapshot, releaseSnapshot, err := d.snapshot(ctx, req.options)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	defer releaseSnapshot()

	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	container, err := d.submitResolveContainer(ctx, driver, req.databaseID, req.containerID)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	defer C.cosmos_container_ref_free(container)

	setup, err := d.readContainerMetadata(ctx, driver, container, req.options, snapshot)
	if err != nil {
		return TransactionalBatchResponse{}, err
	}
	defer setup.release()
	if setup.err != nil {
		return TransactionalBatchResponse{}, setup.err
	}
	if err := validateCompleteBatchPartitionKey(setup.body, req.partitionKey); err != nil {
		return TransactionalBatchResponse{}, batchPartitionKeyError(err, setup)
	}

	result, err := d.awaitCompletion(ctx, "submitting transactional batch",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			request, release := buildNativeItemRequest(req, container)
			defer release()
			request.options_snapshot = snapshot
			return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr targets cgo-generated code.
		})
	if err != nil {
		return TransactionalBatchResponse{}, addSetupCharge(err, setup.response.Response, "executing transactional batch")
	}
	defer result.release()
	if result.err != nil {
		return TransactionalBatchResponse{}, addSetupCharge(result.err, setup.response.Response, "executing transactional batch")
	}
	return result.batchResponse(setup.response.Response, operationCount)
}

func (d *nativeDriver) readContainerMetadata(ctx context.Context, driver *C.cosmos_driver_t, container *C.cosmos_container_ref_t, options OperationOptions, snapshot *C.cosmos_operation_options_snapshot_t) (completionResult, error) {
	return d.awaitCompletion(ctx, "reading container partition key metadata",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			request := newOperationRequest(operationKind(C.COSMOS_OPERATION_KIND_READ_CONTAINER), container)
			nativeOptions, freeOptions := options.toNative()
			defer freeOptions()
			request.options = nativeOptions
			request.options_snapshot = snapshot
			return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr targets cgo-generated code.
		})
}

func (result completionResult) batchResponse(setup Response, operationCount int) (TransactionalBatchResponse, error) {
	response := result.response
	response.RequestCharge += setup.RequestCharge
	if response.ActivityID == "" {
		response.ActivityID = setup.ActivityID
	}
	return decodeTransactionalBatchResponse(response, result.body, result.retryAfter, result.fromWire, operationCount)
}

func batchPartitionKeyError(err error, setup completionResult) error {
	enriched := setupChargeError(err, setup.response.Response, "validating transactional batch partition key")
	enriched.SessionToken = setup.response.SessionToken
	enriched.ETag = setup.response.ETag
	enriched.RetryAfter = setup.retryAfter
	enriched.Body = append([]byte(nil), setup.body...)
	return enriched
}

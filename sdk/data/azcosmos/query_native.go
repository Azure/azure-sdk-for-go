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

import (
	"context"
	"unsafe"
)

func (c *Client) executeQuery(ctx context.Context, req *queryRequest) (QueryItemsResponse, error) {
	ctx, cancel := contextWithEndToEndTimeout(ctx, req.options.Operation.EndToEndTimeout)
	defer cancel()
	d := c.driver
	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return QueryItemsResponse{}, err
	}
	// Resolve through the native cache each page: Go's lifetime-cached handle can retain a
	// deleted container's RID and invalidate tokens issued for its replacement.
	container, err := d.submitResolveContainer(ctx, driver, req.databaseID, req.containerID)
	if err != nil {
		return QueryItemsResponse{}, err
	}
	defer C.cosmos_container_ref_free(container)
	options := req.options.Operation
	var setup Response
	if !req.scopeValidated {
		options.EndToEndTimeout = endToEndTimeout(ctx, 0)
		metadata, err := d.awaitCompletion(ctx, "reading query scope metadata",
			func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
				request := newOperationRequest(operationKind(C.COSMOS_OPERATION_KIND_READ_CONTAINER), container)
				nativeOptions, freeOptions := options.toNative()
				defer freeOptions()
				request.options = nativeOptions
				return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr targets cgo-generated code.
			})
		if err != nil {
			return QueryItemsResponse{}, err
		}
		if metadata.err != nil {
			return QueryItemsResponse{}, metadata.err
		}
		setup = metadata.response.Response
		if err := validateQueryPartitionKey(metadata.body, req.partitionKey); err != nil {
			return QueryItemsResponse{}, addQuerySetupCharge(err, setup)
		}
		req.scopeValidated = true
	}
	options.EndToEndTimeout = endToEndTimeout(ctx, 0)
	result, err := d.awaitCompletion(ctx, "submitting query",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			request, release, status := buildNativeQueryRequest(req, options, container)
			defer release()
			if status != 0 {
				*preError = status
				return nil
			}
			return C.cosmos_submit_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr targets cgo-generated code.
		})
	if err != nil {
		return QueryItemsResponse{}, addQuerySetupCharge(err, setup)
	}
	if result.err != nil {
		return QueryItemsResponse{}, addQuerySetupCharge(result.err, setup)
	}
	response := result.response.Response
	response.RequestCharge += setup.RequestCharge
	return decodeQueryPage(result.body, response, result.response.SessionToken,
		result.nextContinuation, result.httpStatus == 0)
}

func buildNativeQueryRequest(req *queryRequest, options OperationOptions, container *C.cosmos_container_ref_t) (C.cosmos_operation_request_t, func(), C.cosmos_status_code_t) {
	request := newOperationRequest(operationKind(C.COSMOS_OPERATION_KIND_QUERY_ITEMS), container)
	var releases []func()
	release := func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
	components, freeComponents := req.partitionKey.toNative()
	releases = append(releases, freeComponents)
	var pk *C.cosmos_partition_key_t
	if status := C.cosmos_partition_key_create(components, req.partitionKey.partitionKeyLen(), &pk); status != 0 { //nolint:gocritic // dupSubExpr targets cgo-generated code.
		return request, release, status
	}
	releases = append(releases, func() { C.cosmos_partition_key_free(pk) })
	var feedRange *C.cosmos_feed_range_t
	if status := C.cosmos_feed_range_for_partition_key(container, pk, &feedRange); status != 0 { //nolint:gocritic // dupSubExpr targets cgo-generated code.
		return request, release, status
	}
	releases = append(releases, func() { C.cosmos_feed_range_free(feedRange) })
	request.feed_range = feedRange

	body := C.CBytes(req.body)
	releases = append(releases, func() { C.free(body) })
	request.body = (*C.uint8_t)(body)
	request.body_len = C.uintptr_t(len(req.body))
	if token := req.options.Feed.ContinuationToken; token != "" {
		value, allocation := toNativeString(token)
		releases = append(releases, func() { C.free(allocation) })
		request.continuation_token = value
	}
	if token := req.options.SessionToken; token != "" {
		value, allocation := toNativeString(string(token))
		releases = append(releases, func() { C.free(allocation) })
		request.session_token = value
	}
	if req.options.Feed.PageSizeHint > 0 {
		request.max_item_count = C.int32_t(req.options.Feed.PageSizeHint)
	}
	nativeOptions, freeOptions := options.toNative()
	releases = append(releases, freeOptions)
	request.options = nativeOptions
	return request, release, 0
}

// syntheticQueryCompletion frees the native buffers before returning, testing copied ownership.
func syntheticQueryCompletion(body []byte, continuation string, status int) (QueryItemsResponse, error) {
	nativeBody := C.CBytes(body)
	defer C.free(nativeBody)
	nativeToken := C.CString(continuation)
	defer C.free(unsafe.Pointer(nativeToken))
	completion := C.cosmos_completion_t{
		outcome: C.COSMOS_COMPLETION_OUTCOME_OK, http_status_code: C.uint16_t(status),
		body: (*C.uint8_t)(nativeBody), body_len: C.uintptr_t(len(body)),
		next_continuation: nativeToken,
	}
	result := translateCompletionOutcome(&completion)
	return decodeQueryPage(result.body, result.response.Response, result.response.SessionToken,
		result.nextContinuation, result.httpStatus == 0)
}

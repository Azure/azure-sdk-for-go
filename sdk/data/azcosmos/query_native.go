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

func (c *Client) openQuery(ctx context.Context, req *queryRequest) (queryCursor, Response, error) {
	driver, container, setup, err := c.resolveFeedScope(ctx, req.databaseID, req.containerID,
		req.partitionKey, req.fullContainer, req.options.Operation)
	if err != nil {
		return nil, setup, err
	}
	defer C.cosmos_container_ref_free(container)
	request, release, status := buildNativeQueryRequest(req, req.options.Operation, container)
	defer release()
	if status != 0 {
		return nil, setup, statusError(status, nil, "building query")
	}
	cursor, err := c.driver.openCursor(ctx, driver, request)
	if err != nil {
		return nil, setup, err
	}
	return cursor, setup, nil
}

func (c *Client) resolveFeedScope(ctx context.Context, databaseID, containerID string, partitionKey PartitionKey,
	fullContainer bool, options OperationOptions) (*C.cosmos_driver_t, *C.cosmos_container_ref_t, Response, error) {
	d := c.driver
	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return nil, nil, Response{}, err
	}
	// Resolve through the native cache each page: Go's lifetime-cached handle can retain a
	// deleted container's RID and invalidate tokens issued for its replacement.
	container, err := d.submitResolveContainer(ctx, driver, databaseID, containerID)
	if err != nil {
		return nil, nil, Response{}, err
	}
	keepContainer := false
	defer func() {
		if !keepContainer {
			C.cosmos_container_ref_free(container)
		}
	}()
	var setup Response
	if !fullContainer {
		options.EndToEndTimeout = endToEndTimeout(ctx, 0)
		metadata, err := d.awaitCompletion(ctx, "reading feed scope metadata",
			func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
				request := newOperationRequest(operationKind(C.COSMOS_OPERATION_KIND_READ_CONTAINER), container)
				nativeOptions, freeOptions := options.toNative()
				defer freeOptions()
				request.options = nativeOptions
				return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr targets cgo-generated code.
			})
		if err != nil {
			return nil, nil, Response{}, err
		}
		if metadata.err != nil {
			return nil, nil, Response{}, metadata.err
		}
		setup = metadata.response.Response
		if err := validateFeedPartitionKey(metadata.body, partitionKey); err != nil {
			return nil, nil, setup, err
		}
	}
	keepContainer = true
	return driver, container, setup, nil
}

func (result completionResult) queryPage(setup Response) (QueryItemsResponse, error) {
	response := result.response.Response
	response.RequestCharge += setup.RequestCharge
	if response.ActivityID == "" && setup.RequestCharge != 0 {
		response.ActivityID = setup.ActivityID
	}
	return decodeQueryPage(result.body, response, result.response.SessionToken,
		result.nextContinuation, result.httpStatus == 0)
}

func buildNativeQueryRequest(req *queryRequest, options OperationOptions, container *C.cosmos_container_ref_t) (C.cosmos_operation_request_t, func(), C.cosmos_status_code_t) {
	request, freeFeed, status := buildNativeFeedRequest(req.partitionKey, req.fullContainer,
		req.options.Feed, req.options.SessionToken, options, container)
	if status != 0 {
		return request, freeFeed, status
	}
	request.kind = C.COSMOS_OPERATION_KIND_QUERY_ITEMS
	body := C.CBytes(req.body)
	request.body = (*C.uint8_t)(body)
	request.body_len = C.uintptr_t(len(req.body))
	release := func() {
		C.free(body)
		freeFeed()
	}
	switch req.options.QueryPlanMode {
	case QueryPlanModeGatewayOnly:
		request.options.query_plan_mode = C.COSMOS_QUERY_PLAN_MODE_GATEWAY_ONLY
	case QueryPlanModeLocalPreferred:
		request.options.query_plan_mode = C.COSMOS_QUERY_PLAN_MODE_LOCAL_PREFERRED
	}
	request.populate_index_metrics = nativeOptionalBool(req.options.PopulateIndexMetrics)
	request.populate_query_metrics = nativeOptionalBool(req.options.PopulateQueryMetrics)
	return request, release, 0
}

func buildNativeFeedRequest(partitionKey PartitionKey, fullContainer bool, feed FeedOptions, sessionToken SessionToken,
	options OperationOptions, container *C.cosmos_container_ref_t) (C.cosmos_operation_request_t, func(), C.cosmos_status_code_t) {
	request := newOperationRequest(0, container)
	var releases []func()
	release := func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
	if !fullContainer {
		components, freeComponents := partitionKey.toNative()
		releases = append(releases, freeComponents)
		var pk *C.cosmos_partition_key_t
		if status := C.cosmos_partition_key_create(components, partitionKey.partitionKeyLen(), &pk); status != 0 { //nolint:gocritic // dupSubExpr targets cgo-generated code.
			return request, release, status
		}
		releases = append(releases, func() { C.cosmos_partition_key_free(pk) })
		var feedRange *C.cosmos_feed_range_t
		if status := C.cosmos_feed_range_for_partition_key(container, pk, &feedRange); status != 0 { //nolint:gocritic // dupSubExpr targets cgo-generated code.
			return request, release, status
		}
		releases = append(releases, func() { C.cosmos_feed_range_free(feedRange) })
		request.feed_range = feedRange
	}

	if token := feed.ContinuationToken; token != "" {
		value, allocation := toNativeString(token)
		releases = append(releases, func() { C.free(allocation) })
		request.continuation_token = value
	}
	if token := sessionToken; token != "" {
		value, allocation := toNativeString(string(token))
		releases = append(releases, func() { C.free(allocation) })
		request.session_token = value
	}
	if feed.PageSizeHint > 0 {
		request.max_item_count = C.int32_t(feed.PageSizeHint)
	}
	nativeOptions, freeOptions := options.toNative()
	releases = append(releases, freeOptions)
	request.max_fan_out = C.uint32_t(feed.MaxFanOut)
	request.options = nativeOptions
	return request, release, 0
}

func nativeOptionalBool(value *bool) C.int8_t {
	if value == nil {
		return 0
	}
	if *value {
		return 2
	}
	return 1
}

type nativeQueryOptions struct {
	full          bool
	fanOut        uint32
	pageSize      int32
	mode          int32
	indexMetrics  int8
	queryMetrics  int8
	timeoutMillis int64
}

func inspectNativeFullQuery(req *queryRequest) (nativeQueryOptions, error) {
	request, release, status := buildNativeQueryRequest(req, req.options.Operation, nil)
	defer release()
	if status != 0 {
		return nativeQueryOptions{}, statusError(status, nil, "inspecting query")
	}
	return nativeQueryOptions{
		full:   request.feed_range == nil,
		fanOut: uint32(request.max_fan_out), pageSize: int32(request.max_item_count),
		mode:         int32(request.options.query_plan_mode),
		indexMetrics: int8(request.populate_index_metrics), queryMetrics: int8(request.populate_query_metrics),
		timeoutMillis: int64(request.options.end_to_end_timeout_ms),
	}, nil
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
	result := translateCompletionOutcome(&completion, DiagnosticsVerbosityDefault)
	return result.queryPage(Response{})
}

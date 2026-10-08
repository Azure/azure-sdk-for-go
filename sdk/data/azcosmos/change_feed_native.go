// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"context"
	"errors"
	"time"
)

func (c *Client) openChangeFeed(ctx context.Context, req *changeFeedRequest) (changeFeedCursor, Response, error) {
	driver, container, setup, err := c.resolveFeedScope(ctx, req.databaseID, req.containerID,
		req.partitionKey, req.fullContainer, req.options.Operation)
	if err != nil {
		return nil, setup, err
	}
	defer C.cosmos_container_ref_free(container)
	request, release, status := buildNativeChangeFeedRequest(req, container)
	defer release()
	if status != 0 {
		return nil, setup, statusError(status, nil, "building change feed")
	}
	cursor, err := c.driver.openCursorRequest(ctx, driver, request)
	if err != nil {
		return nil, setup, err
	}
	return &nativeChangeFeedCursor{cursor}, setup, nil
}

func buildNativeChangeFeedRequest(req *changeFeedRequest, container *C.cosmos_container_ref_t) (C.cosmos_cursor_request_t, func(), C.cosmos_status_code_t) {
	var request C.cosmos_cursor_request_t
	C.cosmos_cursor_request_init(&request) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	common, freeCommon, status := buildNativeFeedRequest(req.partitionKey, req.fullContainer,
		req.options.Feed, req.options.SessionToken, req.options.Operation, container)
	request.operation = common
	request.change_feed_mode = 1
	if req.options.Mode == ChangeFeedModeAllVersionsAndDeletes {
		request.change_feed_mode = 2
	}
	request.start_from = C.uint32_t(req.startFrom.kind)
	if status != 0 || req.startFrom.kind != 3 {
		return request, freeCommon, status
	}
	timestamp, allocation := toNativeString(req.startFrom.time.Format(time.RFC3339Nano))
	request.start_time = timestamp
	return request, func() {
		C.free(allocation)
		freeCommon()
	}, 0
}

type nativeChangeFeedCursor struct {
	*nativeQueryCursor
}

func (q *nativeChangeFeedCursor) next(ctx context.Context) (ChangeFeedResponse, error) {
	var page ChangeFeedResponse
	err := q.advance(ctx, func(completion *C.cosmos_cursor_completion_t, verbosity DiagnosticsVerbosity) error {
		var err error
		page, err = decodeChangeFeedCursorPage(completion, verbosity)
		return err
	})
	return page, err
}

func decodeChangeFeedCursorPage(completion *C.cosmos_cursor_completion_t, verbosity DiagnosticsVerbosity) (ChangeFeedResponse, error) {
	common := translateCompletionOutcome(&completion.common, verbosity)
	page := ChangeFeedResponse{
		Response: common.response.Response, ETag: common.response.ETag, SessionToken: common.response.SessionToken,
	}
	if completion.result_kind != 2 {
		return ChangeFeedResponse{}, &Error{
			Code: CodeClientError, Message: "azcosmos: expected a change feed page, not native end-of-feed",
			Diagnostics: page.Diagnostics, StatusCode: page.StatusCode, SubStatus: page.SubStatus,
			AttemptCount: page.AttemptCount, RequestCharge: page.RequestCharge, ActivityID: page.ActivityID,
			ETag: page.ETag, SessionToken: page.SessionToken,
		}
	}
	switch completion.body_kind {
	case 0:
		if len(common.body) != 0 || completion.items_len != 0 {
			return ChangeFeedResponse{}, changeFeedResponseError(page, errors.New("no-payload page carried items"))
		}
		return page, nil
	case 1:
		return decodeChangeFeedPage(common.body, page)
	case 2:
		items, err := copyCursorItems(completion)
		if err != nil {
			return ChangeFeedResponse{}, changeFeedResponseError(page, err)
		}
		page.Items = items
		return page, nil
	default:
		return ChangeFeedResponse{}, changeFeedResponseError(page, errors.New("unexpected change feed body kind"))
	}
}

func changeFeedCancellationError(completion *C.cosmos_completion_t, verbosity DiagnosticsVerbosity, cause error) *Error {
	result := translateCompletionOutcome(completion, verbosity)
	cancelled := newOperationCancelledError(cause, result.response.RequestCharge, result.response.ActivityID)
	var nativeError *Error
	if errors.As(result.err, &nativeError) {
		*cancelled = *nativeError
		cancelled.Code, cancelled.Message, cancelled.cause = CodeOperationCancelled, "azcosmos: the Go wait was cancelled", cause
		return cancelled
	}
	cancelled.Diagnostics = result.response.Diagnostics
	cancelled.StatusCode, cancelled.SubStatus = result.response.StatusCode, result.response.SubStatus
	cancelled.AttemptCount = result.response.AttemptCount
	cancelled.SessionToken, cancelled.ETag = result.response.SessionToken, result.response.ETag
	return cancelled
}

type nativeChangeFeedOptions struct {
	full                       bool
	kind, mode, start, fanOut  uint32
	pageSize                   int32
	timeoutMillis              int64
	timestamp, token, session  string
	indexMetrics, queryMetrics int8
}

func inspectNativeFullChangeFeed(req *changeFeedRequest) (nativeChangeFeedOptions, error) {
	request, release, status := buildNativeChangeFeedRequest(req, nil)
	defer release()
	if status != 0 {
		return nativeChangeFeedOptions{}, statusError(status, nil, "inspecting change feed")
	}
	return nativeChangeFeedOptions{
		full: request.operation.feed_range == nil,
		kind: uint32(request.operation.kind), mode: uint32(request.change_feed_mode),
		start: uint32(request.start_from), fanOut: uint32(request.operation.max_fan_out),
		pageSize:      int32(request.operation.max_item_count),
		timeoutMillis: int64(request.operation.options.end_to_end_timeout_ms),
		timestamp:     fromNativeString(request.start_time),
		token:         fromNativeString(request.operation.continuation_token),
		session:       fromNativeString(request.operation.session_token),
		indexMetrics:  int8(request.operation.populate_index_metrics),
		queryMetrics:  int8(request.operation.populate_query_metrics),
	}, nil
}

func syntheticChangeFeedCursorPage(body []byte, items [][]byte, kind, result uint32, status int) (ChangeFeedResponse, error) {
	var page ChangeFeedResponse
	var err error
	withSyntheticCursorPage(body, items, kind, result, status, func(completion *C.cosmos_cursor_completion_t) {
		page, err = decodeChangeFeedCursorPage(completion, DiagnosticsVerbosityDefault)
	})
	return page, err
}

func syntheticChangeFeedCancellation(cause error) *Error {
	var err *Error
	withSyntheticCursorPage(nil, nil, 0, 2, 304, func(completion *C.cosmos_cursor_completion_t) {
		err = changeFeedCancellationError(&completion.common, DiagnosticsVerbosityDefault, cause)
	})
	return err
}

func syntheticInvalidChangeFeedItems(kind uint32) (ChangeFeedResponse, error) {
	completion := C.cosmos_cursor_completion_t{result_kind: 2, body_kind: 2, items_len: 1}
	item := C.cosmos_cursor_bytes_t{len: 1}
	switch kind {
	case 0:
	case 1:
		completion.items_len = ^C.uintptr_t(0)
	case 2:
		completion.items = &item
	case 3:
		item.len = ^C.uintptr_t(0)
		completion.items = &item
	}
	return decodeChangeFeedCursorPage(&completion, DiagnosticsVerbosityDefault)
}

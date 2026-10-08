// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

// cSpell:ignore gocritic

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"

static cosmos_value_t cosmos_cursor_test_string(const char *value) {
	cosmos_value_t out = {0};
	out.kind = COSMOS_VALUE_KIND_STRING;
	out.payload.string_value = value;
	return out;
}

static cosmos_value_t cosmos_cursor_test_charge(void) {
	cosmos_value_t out = {0};
	out.kind = COSMOS_VALUE_KIND_F64;
	out.payload.f64_value = 3.5;
	return out;
}
*/
import "C"

import (
	"context"
	"sync"
	"unsafe"
)

// Cursors share a client-owned, format-isolated completion reactor.
// Client shutdown owns idle cursors; canceled waits retain resources until delivery.
type nativeQueryCursor struct {
	mu         sync.Mutex
	owner      *nativeDriver
	queue      *C.cosmos_completion_queue_t
	reactor    *cursorReactor
	handle     *C.cosmos_cursor_t
	changeFeed bool

	// beforeCompletionWait controls delivery timing in boundary tests.
	beforeCompletionWait func(chan cursorDelivery)
}

func (d *nativeDriver) openCursor(ctx context.Context, driver *C.cosmos_driver_t, operation C.cosmos_operation_request_t) (*nativeQueryCursor, error) {
	var request C.cosmos_cursor_request_t
	C.cosmos_cursor_request_init(&request) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	request.operation = operation
	return d.openCursorRequest(ctx, driver, request)
}

func (d *nativeDriver) openCursorRequest(ctx context.Context, driver *C.cosmos_driver_t, request C.cosmos_cursor_request_t) (*nativeQueryCursor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	reactor, err := d.queryReactor()
	if err != nil {
		return nil, err
	}
	q := &nativeQueryCursor{owner: d, queue: reactor.queue, reactor: reactor, changeFeed: request.change_feed_mode != 0}
	d.mu.Lock()
	if d.cursors == nil {
		d.cursors = make(map[*nativeQueryCursor]struct{})
	}
	d.cursors[q] = struct{}{}
	d.mu.Unlock()
	q.mu.Lock()
	completion, err := q.submit(ctx, func(cookie C.intptr_t, status *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		return C.cosmos_cursor_open_submit(driver, &request, q.queue, cookie, status) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	})
	if err == nil {
		defer C.cosmos_cursor_completion_free(completion)
		err = cursorCompletionError(completion, d.reactor.verbosity)
		if err == nil {
			if completion.result_kind != 1 {
				err = invalidCursorCompletion("expected opened cursor")
			} else {
				q.handle = C.cosmos_cursor_completion_take_cursor(completion)
				if q.handle == nil {
					err = invalidCursorCompletion("opened completion has no cursor")
				}
			}
		}
	}
	q.mu.Unlock()
	if err != nil {
		q.close()
		return nil, err
	}
	return q, nil
}

func (q *nativeQueryCursor) next(ctx context.Context) (QueryItemsResponse, bool, error) {
	var page QueryItemsResponse
	var end bool
	err := q.advance(ctx, func(completion *C.cosmos_cursor_completion_t, verbosity DiagnosticsVerbosity) error {
		var err error
		page, end, err = decodeCursorPage(completion, verbosity)
		return err
	})
	return page, end, err
}

func (q *nativeQueryCursor) advance(ctx context.Context, decode func(*C.cosmos_cursor_completion_t, DiagnosticsVerbosity) error) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if q.handle == nil {
		return invalidCursorCompletion("cursor is closed")
	}
	completion, err := q.submit(ctx, func(cookie C.intptr_t, status *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		return C.cosmos_cursor_next_submit(q.handle, cookie, status)
	})
	if err != nil {
		return err
	}
	defer C.cosmos_cursor_completion_free(completion)
	if err := cursorCompletionError(completion, q.owner.reactor.verbosity); err != nil {
		return err
	}
	return decode(completion, q.owner.reactor.verbosity)
}

func (q *nativeQueryCursor) checkpoint(ctx context.Context) (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if q.handle == nil {
		return "", invalidCursorCompletion("cursor is closed")
	}
	completion, err := q.submit(ctx, func(cookie C.intptr_t, status *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		return C.cosmos_cursor_checkpoint_submit(q.handle, cookie, status)
	})
	if err != nil {
		return "", err
	}
	defer C.cosmos_cursor_completion_free(completion)
	if err := cursorCompletionError(completion, q.owner.reactor.verbosity); err != nil {
		return "", err
	}
	if completion.result_kind != 3 {
		return "", invalidCursorCompletion("expected checkpoint")
	}
	return fromNativeString(completion.checkpoint), nil
}

func (q *nativeQueryCursor) submit(ctx context.Context, submit func(C.intptr_t, *C.cosmos_status_code_t) *C.cosmos_operation_handle_t) (*C.cosmos_cursor_completion_t, error) {
	cookie, waiter, err := q.reactor.register()
	if err != nil {
		return nil, err
	}
	var status C.cosmos_status_code_t
	op := submit(C.intptr_t(cookie), &status)
	if op == nil {
		q.reactor.forget(cookie)
		if status == 0 {
			return nil, invalidCursorCompletion("submission returned no operation")
		}
		return nil, statusError(status, nil, "submitting feed cursor operation")
	}
	if q.beforeCompletionWait != nil {
		q.beforeCompletionWait(waiter)
	}
	finish := func(result cursorDelivery) (*C.cosmos_cursor_completion_t, error) {
		C.cosmos_operation_handle_free(op)
		if result.err != nil {
			return nil, result.err
		}
		if err := ctx.Err(); err != nil {
			headers := readCompletionHeaders(&result.completion.common)
			cancelled := newOperationCancelledError(err, headers.requestCharge, headers.activityID)
			if q.changeFeed {
				cancelled = changeFeedCancellationError(&result.completion.common, q.owner.reactor.verbosity, err)
			}
			C.cosmos_cursor_completion_free(result.completion)
			return nil, cancelled
		}
		return result.completion, nil
	}
	select {
	case result := <-waiter:
		return finish(result)
	case <-ctx.Done():
		// Preserve metadata if delivery and cancellation became ready together.
		select {
		case result := <-waiter:
			return finish(result)
		default:
		}
		q.owner.pending.Add(1)
		go func() {
			defer q.owner.pending.Done()
			defer C.cosmos_operation_handle_free(op)
			result := <-waiter
			if result.completion != nil {
				C.cosmos_cursor_completion_free(result.completion)
			}
		}()
		return nil, newOperationCancelledError(ctx.Err(), 0, "")
	}
}

func (q *nativeQueryCursor) close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	C.cosmos_cursor_free(q.handle)
	q.handle = nil
	q.queue = nil
	q.owner.mu.Lock()
	delete(q.owner.cursors, q)
	q.owner.mu.Unlock()
}

func invalidCursorCompletion(message string) error {
	return &Error{Code: CodeClientError, Message: "azcosmos: " + message}
}

func cursorCompletionError(completion *C.cosmos_cursor_completion_t, verbosity DiagnosticsVerbosity) error {
	if completion.abi_version != 1 || completion.struct_size_bytes < C.uint32_t(unsafe.Sizeof(*completion)) {
		return invalidCursorCompletion("unsupported cursor completion layout")
	}
	if completion.common.outcome != C.COSMOS_COMPLETION_OUTCOME_OK {
		return translateCompletionOutcome(&completion.common, verbosity).err
	}
	return nil
}

func decodeCursorPage(completion *C.cosmos_cursor_completion_t, verbosity DiagnosticsVerbosity) (QueryItemsResponse, bool, error) {
	common := translateCompletionOutcome(&completion.common, verbosity)
	headers := readCompletionHeaders(&completion.common)
	page := QueryItemsResponse{
		Response: common.response.Response, SessionToken: common.response.SessionToken,
		IndexMetrics: headers.indexMetrics, QueryMetrics: headers.queryMetrics,
	}

	if completion.result_kind == 4 {
		if completion.body_kind != 0 || completion.items_len != 0 || len(common.body) != 0 {
			return QueryItemsResponse{}, false, queryResponseError(page.Response, page.SessionToken, invalidCursorCompletion("end carried items"))
		}
		return page, true, nil
	}
	if completion.result_kind != 2 {
		return QueryItemsResponse{}, false, invalidCursorCompletion("expected query page")
	}
	switch completion.body_kind {
	case 1:
		decoded, err := decodeQueryPage(common.body, page.Response, page.SessionToken, "", false)
		if err != nil {
			return QueryItemsResponse{}, false, err
		}
		page.Items = decoded.Items
	case 2:
		var err error
		page.Items, err = copyCursorItems(completion)
		if err != nil {
			return QueryItemsResponse{}, false, queryResponseError(page.Response, page.SessionToken, err)
		}
	default:
		return QueryItemsResponse{}, false, queryResponseError(page.Response, page.SessionToken, invalidCursorCompletion("unexpected query body kind"))
	}
	return page, false, nil
}

func copyCursorItems(completion *C.cosmos_cursor_completion_t) ([][]byte, error) {
	if completion.items_len > C.uintptr_t(int(^uint(0)>>1)) || (completion.items_len != 0 && completion.items == nil) {
		return nil, invalidCursorCompletion("invalid item array")
	}
	items := unsafe.Slice(completion.items, int(completion.items_len))
	copied := make([][]byte, len(items))
	for i, item := range items {
		if item.len > C.uintptr_t(int(^uint(0)>>1)) || (item.len != 0 && item.data == nil) {
			return nil, invalidCursorCompletion("invalid item buffer")
		}
		copied[i] = append([]byte{}, unsafe.Slice((*byte)(unsafe.Pointer(item.data)), int(item.len))...)
	}
	return copied, nil
}

// syntheticCursorPage releases all C-owned input before returning the translated page.
func syntheticCursorPage(body []byte, items [][]byte, kind, result uint32) (QueryItemsResponse, bool, error) {
	var page QueryItemsResponse
	var end bool
	var err error
	withSyntheticCursorPage(body, items, kind, result, 0, func(completion *C.cosmos_cursor_completion_t) {
		page, end, err = decodeCursorPage(completion, DiagnosticsVerbosityDefault)
	})
	return page, end, err
}

func withSyntheticCursorPage(body []byte, items [][]byte, kind, result uint32, status int, decode func(*C.cosmos_cursor_completion_t)) {
	completion := C.cosmos_cursor_completion_t{
		struct_size_bytes: C.uint32_t(unsafe.Sizeof(C.cosmos_cursor_completion_t{})),
		abi_version:       1, result_kind: C.uint32_t(result), body_kind: C.uint32_t(kind),
	}
	completion.common.outcome = C.COSMOS_COMPLETION_OUTCOME_OK
	completion.common.http_status_code = C.uint16_t(status)
	completion.common.body = (*C.uint8_t)(C.CBytes(body))
	completion.common.body_len = C.uintptr_t(len(body))
	defer C.free(unsafe.Pointer(completion.common.body))
	if len(items) != 0 {
		memory := C.malloc(C.size_t(len(items)) * C.size_t(unsafe.Sizeof(C.cosmos_cursor_bytes_t{})))
		defer C.free(memory)
		nativeItems := unsafe.Slice((*C.cosmos_cursor_bytes_t)(memory), len(items))
		for i, item := range items {
			data := C.CBytes(item)
			defer C.free(data)
			nativeItems[i] = C.cosmos_cursor_bytes_t{data: (*C.uint8_t)(data), len: C.uintptr_t(len(item))}
		}
		completion.items = (*C.cosmos_cursor_bytes_t)(memory)
		completion.items_len = C.uintptr_t(len(items))
	}
	index := C.CString(`{"UtilizedIndexes":[]}`)
	query := C.CString("retrievedDocumentCount=2")
	defer C.free(unsafe.Pointer(index))
	defer C.free(unsafe.Pointer(query))
	headers := []C.cosmos_response_header_t{
		{id: C.COSMOS_HEADER_ID_INDEX_METRICS, value: C.cosmos_cursor_test_string(index)},
		{id: C.COSMOS_HEADER_ID_QUERY_METRICS, value: C.cosmos_cursor_test_string(query)},
	}
	if status != 0 {
		etag, session := C.CString("native-position"), C.CString("0:2")
		activity := C.CString("native-change-activity")
		defer C.free(unsafe.Pointer(etag))
		defer C.free(unsafe.Pointer(session))
		defer C.free(unsafe.Pointer(activity))
		headers = append(headers,
			C.cosmos_response_header_t{id: C.COSMOS_HEADER_ID_ETAG, value: C.cosmos_cursor_test_string(etag)},
			C.cosmos_response_header_t{id: C.COSMOS_HEADER_ID_SESSION_TOKEN, value: C.cosmos_cursor_test_string(session)},
			C.cosmos_response_header_t{id: C.COSMOS_HEADER_ID_ACTIVITY_ID, value: C.cosmos_cursor_test_string(activity)},
			C.cosmos_response_header_t{id: C.COSMOS_HEADER_ID_REQUEST_CHARGE, value: C.cosmos_cursor_test_charge()},
		)
	}
	completion.common.headers = &headers[0]
	completion.common.headers_len = C.uintptr_t(len(headers))
	decode(&completion)
}

func (d *nativeDriver) testIdleCursor() (*nativeQueryCursor, error) {
	reactor, err := d.queryReactor()
	if err != nil {
		return nil, err
	}
	q := &nativeQueryCursor{owner: d, queue: reactor.queue, reactor: reactor}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.cursors == nil {
		d.cursors = make(map[*nativeQueryCursor]struct{})
	}
	d.cursors[q] = struct{}{}
	return q, nil
}

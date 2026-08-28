// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_driver

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"context"
	"runtime/cgo"
	"unsafe"
)

// unsetMaxItemCount is what the driver reads as "no page-size hint".
//
// The request struct's numeric fields mostly treat zero as unset, so building it as a Go composite
// literal leaves them correct by default. This one is the exception, and getting it wrong is not a
// silent degradation: the driver rejects a zero hint with an invalid-option-value status before any
// network I/O, so every operation would fail. See newOperationRequest and the test that pins it.
const unsetMaxItemCount = -1

// execute runs one item operation against the driver.
//
// The caller already holds the client's read lock through acquire, which is what keeps the driver
// handles alive for the operation's duration: Close cannot take the write lock until every
// operation has released it.
func (c *Client) execute(ctx context.Context, req itemRequest) (ItemResponse, []byte, error) {
	return c.driver.execute(ctx, req)
}

// execute runs one operation to completion and returns its result.
//
// It submits the operation, then waits for either the driver's completion or the caller's context.
// Cancelling the context cancels the operation at the driver rather than abandoning it, so the
// driver stops work that no longer has a caller.
func (d *nativeDriver) execute(ctx context.Context, req itemRequest) (ItemResponse, []byte, error) {
	driver, err := d.ensureDriver()
	if err != nil {
		return ItemResponse{}, nil, err
	}

	container, err := d.resolveContainer(req.databaseID, req.containerID)
	if err != nil {
		return ItemResponse{}, nil, err
	}

	req.options.EndToEndTimeout = endToEndTimeout(ctx, req.options.EndToEndTimeout)

	// Buffered so the reactor can always deliver without blocking, even after this goroutine has
	// stopped waiting because the context was cancelled.
	pending := &pendingOperation{result: make(chan completionResult, 1)}
	handle := cgo.NewHandle(pending)
	// Freed only after the operation is known to be finished, because the driver round-trips the
	// cookie onto the completion and the reactor dereferences it.
	defer handle.Delete()

	op, err := d.submit(driver, container, req, handle)
	if err != nil {
		return ItemResponse{}, nil, err
	}
	defer C.cosmos_operation_handle_free(op)

	select {
	case result := <-pending.result:
		return result.response, result.body, result.err

	case <-ctx.Done():
		// Ask the driver to stop, then still wait for the completion. The operation is what owns
		// the cookie, so returning before it has been posted would delete a handle the reactor is
		// about to dereference.
		C.cosmos_operation_handle_cancel(op)
		<-pending.result
		return ItemResponse{}, nil, ctx.Err()
	}
}

// newOperationRequest builds a request carrying only its identity and the driver's unset
// sentinels, ready for the caller to fill in.
//
// Most of the struct's numeric fields treat zero as unset, so Go's zero value already means "leave
// it alone" for them. The ones that do not have to be written here.
func newOperationRequest(kind operationKind, container *C.cosmos_container_ref_t) C.cosmos_operation_request_t {
	return C.cosmos_operation_request_t{
		kind:           C.int32_t(kind),
		container:      container,
		max_item_count: unsetMaxItemCount,
	}
}

// inspectRequestSentinels reports the request fields whose unset value is not zero, in Go types,
// so a test can pin them without cgo.
func inspectRequestSentinels(kind operationKind) (maxItemCount int32) {
	request := newOperationRequest(kind, nil)
	return int32(request.max_item_count)
}

// submit builds the request struct and hands it to the driver.
//
// Every pointer in the struct is borrowed for the duration of the call: the driver copies what it
// needs before returning, which is what makes the defers here safe.
func (d *nativeDriver) submit(
	driver *C.cosmos_driver_t,
	container *C.cosmos_container_ref_t,
	req itemRequest,
	handle cgo.Handle,
) (*C.cosmos_operation_handle_t, error) {
	request := newOperationRequest(req.kind, container)

	if req.itemID != "" {
		itemID := C.CString(req.itemID)
		defer C.free(unsafe.Pointer(itemID))
		request.item_id = itemID
	}
	if req.sessionToken != "" {
		sessionToken := C.CString(string(req.sessionToken))
		defer C.free(unsafe.Pointer(sessionToken))
		request.session_token = sessionToken
	}
	if req.ifNoneMatchETag != "" {
		etag := C.CString(req.ifNoneMatchETag)
		defer C.free(unsafe.Pointer(etag))
		request.precondition_kind = C.int32_t(C.COSMOS_PRECONDITION_KIND_IF_NONE_MATCH)
		request.precondition_etag = etag
	}
	if len(req.body) > 0 {
		// Go memory may not be passed to C when it can hold a Go pointer, and the driver copies
		// the bytes before returning, so a C buffer is both required and cheap here.
		body := C.CBytes(req.body)
		defer C.free(body)
		request.body = (*C.uint8_t)(body)
		request.body_len = C.uintptr_t(len(req.body))
	}

	pk, freePartitionKey := req.partitionKey.toNative()
	defer freePartitionKey()
	// The inline component array takes precedence over the handle field, which is what lets the
	// binding avoid constructing a partition key handle whose lifetime it would have to track.
	request.partition_key_components = pk
	request.partition_key_len = req.partitionKey.partitionKeyLen()

	options, freeOptions := req.options.toNative()
	defer freeOptions()
	request.options = options

	var preError C.cosmos_status_code_t
	op := C.cosmos_submit_singleton_operation(driver, &request, d.reactor.queue, C.intptr_t(handle), &preError) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	if op == nil {
		// A pre-flight rejection posts no completion, so it is reported here rather than through
		// the queue.
		httpStatus, subStatus := unpackStatus(preError)
		return nil, &Error{
			Code:       codeForRichError(false, httpStatus, subStatus),
			StatusCode: httpStatus,
			SubStatus:  subStatus,
			Message:    "azcosmos: submitting the operation",
		}
	}
	return op, nil
}

// resolveContainer returns the driver's handle for a container, resolving it on first use.
//
// Resolution reads the container's metadata from the gateway on a cache miss, so the handle is
// cached per client: an item operation would otherwise pay for that lookup every call.
func (d *nativeDriver) resolveContainer(databaseID, containerID string) (*C.cosmos_container_ref_t, error) {
	key := databaseID + "/" + containerID

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return nil, &Error{Code: CodeClientClosed, Message: "the client has been closed"}
	}
	if container, ok := d.containers[key]; ok {
		return container, nil
	}

	cDatabaseID := C.CString(databaseID)
	defer C.free(unsafe.Pointer(cDatabaseID))
	cContainerID := C.CString(containerID)
	defer C.free(unsafe.Pointer(cContainerID))

	var container *C.cosmos_container_ref_t
	var richErr *C.cosmos_error_t
	status := C.cosmos_driver_resolve_container_blocking(d.runtime, d.driver, cDatabaseID, cContainerID, &container, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	if err := statusError(status, richErr, "resolving the container"); err != nil {
		return nil, err
	}

	if d.containers == nil {
		d.containers = make(map[string]*C.cosmos_container_ref_t)
	}
	d.containers[key] = container
	return container, nil
}

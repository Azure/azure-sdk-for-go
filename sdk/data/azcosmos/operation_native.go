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
func (d *nativeDriver) execute(ctx context.Context, req itemRequest) (ItemResponse, []byte, error) {
	ctx, cancel := contextWithEndToEndTimeout(ctx, req.options.EndToEndTimeout)
	defer cancel()

	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return ItemResponse{}, nil, err
	}

	container, err := d.resolveContainer(ctx, driver, req.databaseID, req.containerID)
	if err != nil {
		return ItemResponse{}, nil, err
	}

	// Initialization and resolution have already spent part of the budget, so the operation gets
	// only what remains rather than restarting the configured duration.
	req.options.EndToEndTimeout = endToEndTimeout(ctx, 0)

	result, err := d.awaitCompletion(ctx, "submitting the operation",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			return d.submit(driver, container, req, queue, cookie, preError)
		})
	if err != nil {
		return ItemResponse{}, nil, err
	}
	return result.response, result.body, result.err
}

// awaitCompletion submits one operation and waits for its completion or the caller's context,
// whichever comes first. Driver creation, container resolution and item operations all go through
// it, so all three honor a context the same way.
//
// The submit closure receives what the driver needs to answer: the queue to post the completion
// to, the cookie to round-trip onto it, and somewhere to report a pre-flight rejection. It returns
// NULL when the operation was rejected before it started, which posts no completion.
//
// Cancelling the context cancels the operation at the driver rather than abandoning it, and then
// still waits for the completion. The operation owns the cookie, so returning before it has been
// posted would delete a handle the reactor is about to dereference.
func (d *nativeDriver) awaitCompletion(
	ctx context.Context,
	doing string,
	submit func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t,
) (completionResult, error) {
	if err := ctx.Err(); err != nil {
		return completionResult{}, err
	}

	// Buffered so the reactor can always deliver without blocking, even after this goroutine has
	// stopped waiting because the context was cancelled.
	pending := &pendingOperation{result: make(chan completionResult, 1)}
	handle := cgo.NewHandle(pending)
	// Deleted only once the operation is known to be finished, because the driver round-trips the
	// cookie onto the completion and the reactor dereferences it.
	defer handle.Delete()

	var preError C.cosmos_status_code_t
	op := submit(d.reactor.queue, C.intptr_t(handle), &preError)
	if op == nil {
		// A pre-flight rejection posts no completion, so it is reported here rather than through
		// the queue.
		httpStatus, subStatus := unpackStatus(preError)
		return completionResult{}, &Error{
			Code:       codeForRichError(false, httpStatus, subStatus),
			StatusCode: httpStatus,
			SubStatus:  subStatus,
			Message:    "azcosmos: " + doing,
		}
	}
	defer C.cosmos_operation_handle_free(op)

	select {
	case result := <-pending.result:
		if cause := ctx.Err(); cause != nil {
			terminal, err := resultAfterCancellation(cause, result)
			if err != nil {
				result.release()
			}
			return terminal, err
		}
		return result, nil

	case <-ctx.Done():
		C.cosmos_operation_handle_cancel(op)
		// The terminal result is authoritative when completion and cancellation race. In
		// particular, a successful write must not be reported as cancelled after it committed.
		result := <-pending.result
		terminal, err := resultAfterCancellation(ctx.Err(), result)
		if err != nil {
			result.release()
		}
		return terminal, err
	}
}

func resultAfterCancellation(cause error, result completionResult) (completionResult, error) {
	if !result.cancelled {
		return result, nil
	}
	requestCharge := result.response.RequestCharge
	activityID := result.response.ActivityID
	var completionErr *Error
	if errors.As(result.err, &completionErr) {
		requestCharge = completionErr.RequestCharge
		activityID = completionErr.ActivityID
	}
	return completionResult{}, newOperationCancelledError(cause, requestCharge, activityID)
}

// inspectAwaitCompletionSubmission reports whether awaitCompletion invoked its submit closure.
// Tests use it because cgo types cannot appear directly in _test.go files.
func (d *nativeDriver) inspectAwaitCompletionSubmission(ctx context.Context) (bool, error) {
	submitted := false
	_, err := d.awaitCompletion(ctx, "testing submission",
		func(*C.cosmos_completion_queue_t, C.intptr_t, *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			submitted = true
			return nil
		})
	return submitted, err
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

// submit builds the request struct and hands it to the driver. It reports a pre-flight rejection
// the way the C ABI does, by returning NULL and writing preError, which awaitCompletion turns into
// an error.
//
// Every pointer in the struct is borrowed for the duration of the call: the driver copies what it
// needs before returning, which is what makes the defers here safe.
func (d *nativeDriver) submit(
	driver *C.cosmos_driver_t,
	container *C.cosmos_container_ref_t,
	req itemRequest,
	queue *C.cosmos_completion_queue_t,
	cookie C.intptr_t,
	preError *C.cosmos_status_code_t,
) *C.cosmos_operation_handle_t {
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

	return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
}

// resolveContainer returns the driver's handle for a container, resolving it on first use.
//
// Resolution reads the container's metadata from the gateway on a cache miss, so the handle is
// cached per client: an item operation would otherwise pay for that lookup every call. It is
// awaited rather than blocked on, so a miss does not make an operation ignore its context.
//
// The lock is not held across the resolution, so two callers can miss on the same container at
// once. That is deliberate: holding it would serialize every miss behind one gateway round trip,
// including misses for unrelated containers. The loser of the race frees its own handle and takes
// the winner's, so the cache keeps exactly one handle per key.
func (d *nativeDriver) resolveContainer(
	ctx context.Context,
	driver *C.cosmos_driver_t,
	databaseID, containerID string,
) (*C.cosmos_container_ref_t, error) {
	key := databaseID + "/" + containerID

	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return nil, errClientClosed()
	}
	if container, ok := d.containers[key]; ok {
		d.mu.Unlock()
		return container, nil
	}
	d.mu.Unlock()

	container, err := d.submitResolveContainer(ctx, driver, databaseID, containerID)
	if err != nil {
		return nil, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		// Closed while the resolution was in flight, so close has already emptied the cache and
		// this handle has no other owner.
		C.cosmos_container_ref_free(container)
		return nil, errClientClosed()
	}
	if existing, ok := d.containers[key]; ok {
		C.cosmos_container_ref_free(container)
		return existing, nil
	}
	if d.containers == nil {
		d.containers = make(map[string]*C.cosmos_container_ref_t)
	}
	d.containers[key] = container
	return container, nil
}

// submitResolveContainer runs one container resolution against the driver.
func (d *nativeDriver) submitResolveContainer(
	ctx context.Context,
	driver *C.cosmos_driver_t,
	databaseID, containerID string,
) (*C.cosmos_container_ref_t, error) {
	cDatabaseID := C.CString(databaseID)
	defer C.free(unsafe.Pointer(cDatabaseID))
	cContainerID := C.CString(containerID)
	defer C.free(unsafe.Pointer(cContainerID))

	result, err := d.awaitCompletion(ctx, "resolving the container",
		func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
			return C.cosmos_driver_resolve_container_submit(driver, cDatabaseID, cContainerID, queue, cookie, preError) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
		})
	if err != nil {
		return nil, err
	}
	if result.err != nil {
		result.release()
		return nil, result.err
	}
	if result.container == nil {
		result.release()
		return nil, &Error{
			Code:    CodeClientError,
			Message: "azcosmos: the container was resolved but the completion carried no handle",
		}
	}
	return result.container, nil
}

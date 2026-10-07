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
	"errors"
	"math"
	"runtime/cgo"
	"time"
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
	switch req.kind {
	case operationKindCreateItem, operationKindReadItem, operationKindReplaceItem, operationKindUpsertItem:
		if req.options.BinaryEncoding == nil {
			req.options.BinaryEncoding = c.binaryEncoding.clone()
		}
	}
	return c.driver.execute(ctx, req)
}

// execute runs one operation to completion and returns its result.
func (d *nativeDriver) execute(ctx context.Context, req itemRequest) (ItemResponse, []byte, error) {
	if req.patchStrategy != PatchStrategyUnset {
		req.options.PatchStrategy = req.patchStrategy
	}
	ctx, snapshot, releaseSnapshot, err := d.snapshot(ctx, req.options)
	if err != nil {
		return ItemResponse{}, nil, err
	}
	defer releaseSnapshot()

	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return ItemResponse{}, nil, err
	}

	container, err := d.resolveContainer(ctx, driver, req.databaseID, req.containerID)
	if err != nil {
		return ItemResponse{}, nil, err
	}

	submit := func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		request, release := buildNativeItemRequest(req, container)
		defer release()
		request.options_snapshot = snapshot
		return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr is reported against cgo-generated code.
	}

	// Only a write can leave a side effect that must not be misreported as cancelled after it
	// already committed, so only writes pay for an authoritative wait. A read has nothing to
	// protect: abandoning it the moment ctx ends, rather than waiting out the native driver's own
	// retry budget, is what keeps cancellation prompt.
	var result completionResult
	if req.kind == operationKindReadItem {
		result, err = d.awaitCompletion(ctx, "submitting the operation", submit)
	} else {
		result, err = d.awaitAuthoritativeCompletion(ctx, "submitting the operation", submit)
	}
	if err != nil {
		return ItemResponse{}, nil, patchCancellationError(err, req.patchTrackingID)
	}
	return result.response, result.body, result.err
}

func patchCancellationError(err error, requestedID PatchTrackingID) error {
	var cancelled *Error
	if errors.As(err, &cancelled) && cancelled.Code == CodeOperationCancelled && cancelled.PatchTrackingID == "" {
		cancelled.PatchTrackingID = requestedID
	}
	return err
}

// awaitCompletion submits one operation and waits for its completion or the caller's context,
// whichever comes first. Driver creation, container resolution and reads all go through it, so
// all three honor a context the same way.
//
// The submit closure receives what the driver needs to answer: the queue to post the completion
// to, the cookie to round-trip onto it, and somewhere to report a pre-flight rejection. It returns
// NULL when the operation was rejected before it started, which posts no completion.
//
// Non-authoritative waits may return on cancellation; late delivery releases their result.
// Close drains those deliveries before freeing resources. Writes await authoritative results.
func (d *nativeDriver) awaitCompletion(
	ctx context.Context,
	doing string,
	submit func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t,
) (completionResult, error) {
	return d.awaitOperation(ctx, doing, submit, false)
}

func (d *nativeDriver) awaitAuthoritativeCompletion(
	ctx context.Context,
	doing string,
	submit func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t,
) (completionResult, error) {
	return d.awaitOperation(ctx, doing, submit, true)
}

func (d *nativeDriver) awaitOperation(
	ctx context.Context,
	doing string,
	submit func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t,
	authoritative bool,
) (completionResult, error) {
	if err := ctx.Err(); err != nil {
		return completionResult{}, err
	}

	// Buffered so the reactor can deliver without blocking after cancellation.
	pending := &pendingOperation{result: make(chan completionResult, 1), delivered: make(chan struct{})}
	handle := cgo.NewHandle(pending)

	var preError C.cosmos_status_code_t
	op := submit(d.reactor.queue, C.intptr_t(handle), &preError)
	if op == nil {
		handle.Delete()
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
	defer func() {
		release := func() {
			C.cosmos_operation_handle_free(op)
			handle.Delete()
		}
		select {
		case <-pending.delivered:
			release()
		default:
			d.pending.Add(1)
			go func() {
				defer d.pending.Done()
				<-pending.delivered
				release()
			}()
		}
	}()

	return awaitOperationResult(ctx, pending, authoritative)
}

func awaitOperationResult(ctx context.Context, pending *pendingOperation, authoritative bool) (completionResult, error) {
	finish := func(result completionResult) (completionResult, error) {
		if cause := ctx.Err(); cause != nil {
			if !authoritative {
				defer result.release()
				return completionResult{}, completionCancellationError(cause, result)
			}
			terminal, err := resultAfterCancellation(cause, result)
			if err != nil {
				result.release()
			}
			return terminal, err
		}
		return result, nil
	}
	select {
	case result := <-pending.result:
		return finish(result)
	case <-ctx.Done():
		if !authoritative {
			// abandon() marks the pending operation closed and claims whatever result is buffered
			// atomically, under its own lock, so this cannot race with the reactor's concurrent
			// deliver(): either the result was already delivered and abandon() hands it back here,
			// or deliver() observes closed and releases it itself. A result that lands in that
			// exact instant is never silently discarded in favor of a cancellation error.
			if result, ok := pending.abandon(); ok {
				return finish(result)
			}
			return completionResult{}, newOperationCancelledError(ctx.Err(), 0, "")
		}
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
	if !result.cancelled && !isClientOperationTimeout(result.err) {
		return result, nil
	}
	return completionResult{}, completionCancellationError(cause, result)
}

// isClientOperationTimeout reports whether err is the native driver's own end-to-end budget
// expiring (CodeClientOperationTimeout), as opposed to the COSMOS_COMPLETION_OUTCOME_CANCELLED
// outcome that result.cancelled already covers.
//
// The native driver's end-to-end timeout and the caller's context deadline are set to the same
// nominal duration but run on independent clocks, so either can fire first. When the native clock
// wins that race, the completion arrives classified as ClientOperationTimeout rather than
// cancelled, and an authoritative wait (which always awaits the real completion, since a write may
// have already committed) would otherwise return that raw error instead of one that satisfies
// errors.Is(err, context.DeadlineExceeded) as callers expect once the context has ended.
func isClientOperationTimeout(err error) bool {
	var completionErr *Error
	return errors.As(err, &completionErr) && completionErr.Code == CodeClientOperationTimeout
}

func completionCancellationError(cause error, result completionResult) error {
	requestCharge := result.response.RequestCharge
	activityID := result.response.ActivityID
	diagnostics := result.response.Diagnostics
	attemptCount := result.response.AttemptCount
	statusCode := result.response.StatusCode
	subStatus := result.response.SubStatus
	var completionErr *Error
	if errors.As(result.err, &completionErr) {
		requestCharge = completionErr.RequestCharge
		activityID = completionErr.ActivityID
		diagnostics = completionErr.Diagnostics
		attemptCount = completionErr.AttemptCount
		statusCode = completionErr.StatusCode
		subStatus = completionErr.SubStatus
	}
	err := newOperationCancelledError(cause, requestCharge, activityID)
	err.PatchTrackingID = result.response.PatchTrackingID
	if completionErr != nil {
		err.PatchTrackingID = completionErr.PatchTrackingID
	}
	err.Diagnostics = diagnostics
	err.AttemptCount = attemptCount
	err.StatusCode = statusCode
	err.SubStatus = subStatus
	return err
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
	request, release := buildNativeItemRequest(req, container)
	defer release()

	return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
}

// buildNativeItemRequest converts the Go request into the borrowed C request passed to submit. The
// returned function releases every allocation after the driver has copied the request.
func buildNativeItemRequest(req itemRequest, container *C.cosmos_container_ref_t) (C.cosmos_operation_request_t, func()) {
	request := newOperationRequest(req.kind, container)
	var releases []func()
	if req.patchMaxAttempts != nil {
		request.patch_max_attempts = C.uint8_t(*req.patchMaxAttempts)
	}
	if req.patchTrackingCapacity != nil {
		request.patch_tracking_capacity = C.uint16_t(*req.patchTrackingCapacity)
	}
	if req.patchTrackingRetention != nil {
		request.patch_tracking_retention_seconds = C.uint32_t(min(math.MaxUint32, max(1, *req.patchTrackingRetention/time.Second)))
	}
	if req.patchTrackingID != "" {
		id, allocation := toNativeString(string(req.patchTrackingID))
		releases = append(releases, func() { C.free(allocation) })
		request.patch_tracking_id = id
	}

	if req.itemID != "" {
		itemID, allocation := toNativeString(req.itemID)
		releases = append(releases, func() { C.free(allocation) })
		request.item_id = itemID
	}
	if req.sessionToken != "" {
		sessionToken, allocation := toNativeString(string(req.sessionToken))
		releases = append(releases, func() { C.free(allocation) })
		request.session_token = sessionToken
	}
	if req.preconditionKind != preconditionKindNone {
		etag, allocation := toNativeString(req.preconditionETag)
		releases = append(releases, func() { C.free(allocation) })
		request.precondition_kind = C.int32_t(req.preconditionKind)
		request.precondition_etag = etag
	}
	if len(req.body) > 0 {
		// Go memory may not be passed to C when it can hold a Go pointer, and the driver copies
		// the bytes before returning, so a C buffer is both required and cheap here.
		body := C.CBytes(req.body)
		releases = append(releases, func() { C.free(body) })
		request.body = (*C.uint8_t)(body)
		request.body_len = C.uintptr_t(len(req.body))
	}

	pk, freePartitionKey := req.partitionKey.toNative()
	releases = append(releases, freePartitionKey)
	// The inline component array takes precedence over the handle field, which is what lets the
	// binding avoid constructing a partition key handle whose lifetime it would have to track.
	request.partition_key_components = pk
	request.partition_key_len = req.partitionKey.partitionKeyLen()

	options, freeOptions := req.options.toNative()
	releases = append(releases, freeOptions)
	if strategy, ok := req.patchStrategy.toNative(); ok {
		options.patch_strategy = strategy
	}
	request.options = options

	return request, func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
}

// nativeItemRequest is a converted request in Go types, for tests that cannot import C.
type nativeItemRequest struct {
	kind                   int32
	itemID                 string
	partitionKeyLen        int
	body                   []byte
	sessionToken           string
	maxItemCount           int32
	preconditionKind       int32
	preconditionETag       string
	contentResponseWrite   int32
	patchStrategy          int32
	patchMaxAttempts       uint8
	patchTrackingID        string
	patchTrackingCapacity  uint16
	patchTrackingRetention uint32
}

// inspectNativeItemRequest converts a request and reads it back before releasing its C memory.
func inspectNativeItemRequest(req itemRequest) (nativeItemRequest, func()) {
	request, release := buildNativeItemRequest(req, nil)
	converted := nativeItemRequest{
		kind:                   int32(request.kind),
		itemID:                 fromNativeString(request.item_id),
		partitionKeyLen:        int(request.partition_key_len),
		sessionToken:           fromNativeString(request.session_token),
		maxItemCount:           int32(request.max_item_count),
		preconditionKind:       int32(request.precondition_kind),
		preconditionETag:       fromNativeString(request.precondition_etag),
		patchMaxAttempts:       uint8(request.patch_max_attempts),
		patchTrackingID:        fromNativeString(request.patch_tracking_id),
		patchTrackingCapacity:  uint16(request.patch_tracking_capacity),
		patchTrackingRetention: uint32(request.patch_tracking_retention_seconds),
	}
	if request.body != nil {
		converted.body = C.GoBytes(unsafe.Pointer(request.body), C.int(request.body_len))
	}
	if request.options != nil {
		converted.contentResponseWrite = int32(request.options.content_response_on_write)
		converted.patchStrategy = int32(request.options.patch_strategy)
	}
	return converted, release
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
	cDatabaseID, databaseAllocation := toNativeString(databaseID)
	defer C.free(databaseAllocation)
	cContainerID, containerAllocation := toNativeString(containerID)
	defer C.free(containerAllocation)

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

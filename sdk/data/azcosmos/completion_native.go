// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore gocritic

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"

static cosmos_value_t cosmos_test_string_value(const char *value) {
	cosmos_value_t out = {0};
	out.kind = COSMOS_VALUE_KIND_STRING;
	out.payload.string_value = value;
	return out;
}

static cosmos_value_t cosmos_test_i64_value(int64_t value) {
	cosmos_value_t out = {0};
	out.kind = COSMOS_VALUE_KIND_I64;
	out.payload.i64_value = value;
	return out;
}

static cosmos_value_t cosmos_test_f64_value(double value) {
	cosmos_value_t out = {0};
	out.kind = COSMOS_VALUE_KIND_F64;
	out.payload.f64_value = value;
	return out;
}
*/
import "C"

import (
	"time"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// completionResult is everything an operation needs from a completion, copied out of driver-owned
// memory so it stays valid after the completion is freed.
type completionResult struct {
	response         ItemResponse
	body             []byte
	err              error
	cancelled        bool
	httpStatus       int
	nextContinuation string

	// driver and container are detached from the completion rather than copied, because they are
	// handles and not data. A get_or_create completion carries the first and a resolve_container
	// completion the second; every other completion carries neither.
	//
	// Detaching transfers ownership, so whoever ends up with this result owns them, and release
	// is what closes that loop when nobody does.
	driver    *C.cosmos_driver_t
	container *C.cosmos_container_ref_t
}

// release frees the handles a result owns. It is for the paths that drop a result rather than
// returning it, which would otherwise leak whatever the completion carried.
func (r completionResult) release() {
	// Freeing is a no-op on NULL, so the common case of a result carrying neither is free.
	C.cosmos_driver_free(r.driver)
	C.cosmos_container_ref_free(r.container)
}

// translateCompletion copies a completion into Go memory.
//
// Everything it reads is borrowed from the driver and reclaimed when the completion is freed at the
// end of the drain, so every string and byte slice is copied rather than referenced.
func translateCompletion(completion *C.cosmos_completion_t) completionResult {
	result := translateCompletionOutcome(completion)
	diagnostics := copyNativeDiagnostics(completion.diagnostics)
	result.response.Diagnostics = diagnostics
	if err, ok := result.err.(*Error); ok {
		err.Diagnostics = diagnostics
	}

	// Taken here rather than by the waiter, because the completion is freed at the end of this
	// drain and a handle left on it would be reclaimed with it. Both return NULL when the
	// completion carries nothing, so this is unconditional.
	result.driver = C.cosmos_completion_take_driver(completion)
	result.container = C.cosmos_completion_take_container(completion)
	return result
}

func copyNativeDiagnostics(native *C.cosmos_diagnostics_t) Diagnostics {
	if native == nil {
		return Diagnostics{}
	}
	const maxMicros = uint64((1<<63 - 1) / time.Microsecond)
	snapshot := Diagnostics{
		elapsed:       time.Duration(min(uint64(C.cosmos_diagnostics_total_elapsed_micros(native)), maxMicros)) * time.Microsecond,
		requestCharge: float64(C.cosmos_diagnostics_total_request_charge(native)),
		requestCount:  uint32(C.cosmos_diagnostics_request_count(native)),
		compacted:     bool(C.cosmos_diagnostics_is_compacted(native)),
		failure:       bool(C.cosmos_diagnostics_is_failure(native)),
	}
	var data *C.uint8_t
	var length C.uintptr_t
	status := C.cosmos_diagnostics_to_json(native, C.cosmos_diagnostics_verbosity_t_DEFAULT, &data, &length) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	snapshot.err = statusError(status, nil, "rendering native diagnostics")
	if snapshot.err == nil && data != nil {
		snapshot.json = string(unsafe.Slice((*byte)(unsafe.Pointer(data)), int(length)))
	}
	return snapshot
}

// translateCompletionOutcome copies the data half of a completion, leaving the handles to
// translateCompletion.
func translateCompletionOutcome(completion *C.cosmos_completion_t) completionResult {
	headers := readCompletionHeaders(completion)
	var trackingID PatchTrackingID
	if id := C.cosmos_completion_patch_tracking_id(completion); id != nil {
		trackingID = PatchTrackingID(C.GoString(id))
	}

	response := ItemResponse{
		Response: Response{
			RequestCharge: headers.requestCharge,
			ActivityID:    headers.activityID,
		},
		ETag:            headers.etag,
		SessionToken:    headers.sessionToken,
		PatchTrackingID: trackingID,
	}

	switch completion.outcome {
	case C.COSMOS_COMPLETION_OUTCOME_OK:
		return completionResult{
			response:         response,
			body:             copyCompletionBody(completion),
			httpStatus:       int(completion.http_status_code),
			nextContinuation: C.GoString(completion.next_continuation),
		}

	case C.COSMOS_COMPLETION_OUTCOME_CANCELLED:
		return completionResult{
			cancelled: true,
			err: &Error{
				Code:            CodeOperationCancelled,
				Message:         "azcosmos: the operation was cancelled",
				RequestCharge:   headers.requestCharge,
				ActivityID:      headers.activityID,
				PatchTrackingID: trackingID,
			},
		}

	default:
		// ERROR, and UNKNOWN, which the driver documents as a state the host should treat as a
		// failure rather than assume anything about.
		err := completionError(completion, headers)
		err.PatchTrackingID = trackingID
		return completionResult{err: err}
	}
}

// completionError builds the [Error] for a failed completion.
func completionError(completion *C.cosmos_completion_t, headers completionHeaders) *Error {
	httpStatus := int(completion.http_status_code)
	fromWire := completion.is_from_wire == 1

	// The completion's status is packed, unlike the sync out_error paths. Its sub-status is the
	// authoritative one; the header carries the same value but only when the service sent it.
	_, packedSubStatus := unpackStatus(completion.status)
	subStatus := packedSubStatus
	if subStatus == 0 {
		subStatus = headers.subStatus
	}

	err := &Error{
		Code:          codeForRichError(fromWire, httpStatus, subStatus),
		StatusCode:    httpStatus,
		SubStatus:     subStatus,
		RequestCharge: headers.requestCharge,
		ActivityID:    headers.activityID,
		SessionToken:  headers.sessionToken,
		ETag:          headers.etag,
		RetryAfter:    headers.retryAfter,
		FromWire:      fromWire,
		Body:          copyCompletionBody(completion),
	}
	if completion.message != nil {
		err.Message = C.GoString(completion.message)
	}
	return err
}

// copyCompletionBody copies the completion's body into Go memory, or returns nil when there is
// none. The driver reclaims the buffer when the completion is freed.
func copyCompletionBody(completion *C.cosmos_completion_t) []byte {
	if completion.body == nil || completion.body_len == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(completion.body), C.int(completion.body_len))
}

// completionHeaders is the subset of a completion's response headers this package surfaces.
type completionHeaders struct {
	requestCharge float64
	activityID    string
	sessionToken  SessionToken
	etag          azcore.ETag
	subStatus     int
	retryAfter    time.Duration
	indexMetrics  string
	queryMetrics  string
}

// readCompletionHeaders pulls the headers this package surfaces out of a completion.
//
// The driver reports headers as a typed array rather than as struct fields, so each is identified
// by its id and carries a tagged value. Reading them by id rather than by position is what keeps
// this correct as the driver adds headers.
func readCompletionHeaders(completion *C.cosmos_completion_t) completionHeaders {
	var headers completionHeaders
	if completion.headers == nil || completion.headers_len == 0 {
		return headers
	}

	all := unsafe.Slice(completion.headers, int(completion.headers_len))
	for i := range all {
		header := &all[i]
		switch header.id {
		case C.COSMOS_HEADER_ID_INDEX_METRICS:
			headers.indexMetrics = headerString(&header.value)
		case C.COSMOS_HEADER_ID_QUERY_METRICS:
			headers.queryMetrics = headerString(&header.value)
		case C.COSMOS_HEADER_ID_REQUEST_CHARGE:
			headers.requestCharge = headerFloat(&header.value)
		case C.COSMOS_HEADER_ID_ACTIVITY_ID:
			headers.activityID = headerString(&header.value)
		case C.COSMOS_HEADER_ID_SESSION_TOKEN:
			headers.sessionToken = SessionToken(headerString(&header.value))
		case C.COSMOS_HEADER_ID_ETAG:
			headers.etag = azcore.ETag(headerString(&header.value))
		case C.COSMOS_HEADER_ID_SUB_STATUS:
			headers.subStatus = normalizeSubStatus(int32(headerInt(&header.value)))
		case C.COSMOS_HEADER_ID_RETRY_AFTER_MS:
			if ms := headerInt(&header.value); ms > 0 {
				headers.retryAfter = time.Duration(ms) * time.Millisecond
			}
		}
	}
	return headers
}

// The three header accessors below each read the payload for one Go type. A header's value is a
// tagged union, and the driver is free to report a numeric header as any of its numeric kinds, so
// each accessor accepts every kind it can convert without loss rather than assuming one.

// headerString reads a string-valued header, or "" when the value is not a string.
func headerString(value *C.cosmos_value_t) string {
	if value.kind != C.COSMOS_VALUE_KIND_STRING {
		return ""
	}
	// The union's first member is the string pointer.
	ptr := *(**C.char)(unsafe.Pointer(&value.payload))
	if ptr == nil {
		return ""
	}
	return C.GoString(ptr)
}

// headerFloat reads a numeric header as a float64.
func headerFloat(value *C.cosmos_value_t) float64 {
	switch value.kind {
	case C.COSMOS_VALUE_KIND_F64:
		return float64(*(*C.double)(unsafe.Pointer(&value.payload)))
	case C.COSMOS_VALUE_KIND_I64:
		return float64(*(*C.int64_t)(unsafe.Pointer(&value.payload)))
	case C.COSMOS_VALUE_KIND_U64:
		return float64(*(*C.uint64_t)(unsafe.Pointer(&value.payload)))
	default:
		return 0
	}
}

// headerInt reads a numeric header as an int64.
func headerInt(value *C.cosmos_value_t) int64 {
	switch value.kind {
	case C.COSMOS_VALUE_KIND_I64:
		return int64(*(*C.int64_t)(unsafe.Pointer(&value.payload)))
	case C.COSMOS_VALUE_KIND_U64:
		return int64(*(*C.uint64_t)(unsafe.Pointer(&value.payload)))
	case C.COSMOS_VALUE_KIND_F64:
		return int64(*(*C.double)(unsafe.Pointer(&value.payload)))
	default:
		return 0
	}
}

// syntheticThrottledCompletion builds and translates a C-owned completion for native tests. Every
// borrowed buffer is freed before it returns, so assertions also prove the translator copied out.
func syntheticThrottledCompletion(packedSubStatus int) *Error {
	var allocations []unsafe.Pointer
	cString := func(value string) *C.char {
		ptr := C.CString(value)
		allocations = append(allocations, unsafe.Pointer(ptr))
		return ptr
	}
	defer func() {
		for _, allocation := range allocations {
			C.free(allocation)
		}
	}()

	const headerCount = 6
	headerMemory := C.malloc(C.size_t(headerCount) * C.size_t(unsafe.Sizeof(C.cosmos_response_header_t{})))
	defer C.free(headerMemory)
	headers := unsafe.Slice((*C.cosmos_response_header_t)(headerMemory), headerCount)
	headers[0] = C.cosmos_response_header_t{
		id:    C.COSMOS_HEADER_ID_REQUEST_CHARGE,
		value: C.cosmos_test_f64_value(4.5),
	}
	headers[1] = C.cosmos_response_header_t{
		id:    C.COSMOS_HEADER_ID_ACTIVITY_ID,
		value: C.cosmos_test_string_value(cString("activity-123")),
	}
	headers[2] = C.cosmos_response_header_t{
		id:    C.COSMOS_HEADER_ID_SESSION_TOKEN,
		value: C.cosmos_test_string_value(cString("1:2")),
	}
	headers[3] = C.cosmos_response_header_t{
		id:    C.COSMOS_HEADER_ID_ETAG,
		value: C.cosmos_test_string_value(cString("\"etag\"")),
	}
	headers[4] = C.cosmos_response_header_t{
		id:    C.COSMOS_HEADER_ID_SUB_STATUS,
		value: C.cosmos_test_i64_value(3200),
	}
	headers[5] = C.cosmos_response_header_t{
		id:    C.COSMOS_HEADER_ID_RETRY_AFTER_MS,
		value: C.cosmos_test_i64_value(125),
	}

	body := []byte(`{"code":"TooManyRequests","message":"slow down"}`)
	bodyMemory := C.CBytes(body)
	defer C.free(bodyMemory)

	completion := C.cosmos_completion_t{
		outcome:          C.COSMOS_COMPLETION_OUTCOME_ERROR,
		status:           C.cosmos_status_code_t((429 << 16) | packedSubStatus),
		http_status_code: 429,
		is_from_wire:     1,
		message:          cString("Request rate is large."),
		headers:          (*C.cosmos_response_header_t)(headerMemory),
		headers_len:      headerCount,
		body:             (*C.uint8_t)(bodyMemory),
		body_len:         C.uintptr_t(len(body)),
	}
	result := translateCompletion(&completion)
	return result.err.(*Error)
}

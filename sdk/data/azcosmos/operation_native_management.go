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

// execute runs one management operation (database or container lifecycle) against the driver.
//
// The caller already holds the client's read lock through acquire, which is what keeps the driver
// handles alive for the operation's duration: Close cannot take the write lock until every
// operation has released it.
func (c *Client) executeManagement(
	ctx context.Context,
	kind operationKind,
	databaseID, containerID string,
	body []byte,
	operation OperationOptions,
) (ItemResponse, []byte, error) {
	return c.driver.executeManagement(ctx, kind, databaseID, containerID, body, operation)
}

// executeManagement resolves the account/database/container scope the operation needs, submits
// it, and converts the result.
//
// Database references are pure value-type constructions. Container references are resolved through
// the native driver for each operation so Rust owns metadata caching and stale-resource refresh.
func (d *nativeDriver) executeManagement(
	ctx context.Context,
	kind operationKind,
	databaseID, containerID string,
	body []byte,
	operation OperationOptions,
) (ItemResponse, []byte, error) {
	ctx, cancel := contextWithEndToEndTimeout(ctx, operation.EndToEndTimeout)
	defer cancel()

	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return ItemResponse{}, nil, err
	}

	var database *C.cosmos_database_ref_t
	var container *C.cosmos_container_ref_t

	switch kind {
	case operationKindReadDatabase, operationKindDeleteDatabase, operationKindCreateContainer:
		database, err = d.createDatabaseRef(databaseID)
		if err != nil {
			return ItemResponse{}, nil, err
		}
		defer C.cosmos_database_ref_free(database)
	case operationKindReadContainer, operationKindDeleteContainer:
		container, err = d.submitResolveContainer(ctx, driver, databaseID, containerID)
		if err != nil {
			return ItemResponse{}, nil, err
		}
		defer C.cosmos_container_ref_free(container)
	}

	// Resolution has already spent part of the budget, so the operation gets only what remains
	// rather than restarting the configured duration.
	operation.EndToEndTimeout = endToEndTimeout(ctx, 0)

	submit := func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		request, release := buildNativeManagementRequest(kind, d.account, database, container, body, operation)
		defer release()
		return C.cosmos_submit_singleton_operation(driver, &request, queue, cookie, preError) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	}

	// Only a write can leave a side effect that must not be misreported as cancelled after it
	// already committed, so only a read pays for the cheaper, non-authoritative wait.
	var result completionResult
	if kind == operationKindReadDatabase || kind == operationKindReadContainer {
		result, err = d.awaitCompletion(ctx, "submitting the operation", submit)
	} else {
		result, err = d.awaitAuthoritativeCompletion(ctx, "submitting the operation", submit)
	}
	if err != nil {
		return ItemResponse{}, nil, err
	}
	return result.response, result.body, result.err
}

// nativeManagementRequest is a converted management request in Go types, for tests that cannot
// import C.
type nativeManagementRequest struct {
	kind                 int32
	hasAccount           bool
	hasDatabase          bool
	hasContainer         bool
	body                 []byte
	maxItemCount         int32
	contentResponseWrite int32
}

// inspectNativeManagementRequest converts a request and reads it back before releasing its C
// memory.
func inspectNativeManagementRequest(
	kind operationKind,
	account *C.cosmos_account_ref_t,
	database *C.cosmos_database_ref_t,
	container *C.cosmos_container_ref_t,
	body []byte,
) (nativeManagementRequest, func()) {
	return inspectNativeManagementRequestWithOptions(kind, account, database, container, body, OperationOptions{})
}

func inspectNativeManagementRequestWithOptions(
	kind operationKind,
	account *C.cosmos_account_ref_t,
	database *C.cosmos_database_ref_t,
	container *C.cosmos_container_ref_t,
	body []byte,
	operation OperationOptions,
) (nativeManagementRequest, func()) {
	request, release := buildNativeManagementRequest(kind, account, database, container, body, operation)
	converted := nativeManagementRequest{
		kind:         int32(request.kind),
		hasAccount:   request.account != nil,
		hasDatabase:  request.database != nil,
		hasContainer: request.container != nil,
		maxItemCount: int32(request.max_item_count),
	}
	if request.body != nil {
		converted.body = C.GoBytes(unsafe.Pointer(request.body), C.int(request.body_len))
	}
	if request.options != nil {
		converted.contentResponseWrite = int32(request.options.content_response_on_write)
	}
	return converted, release
}

// createDatabaseRef builds a pure value-type database reference. It performs no network I/O.
func (d *nativeDriver) createDatabaseRef(databaseID string) (*C.cosmos_database_ref_t, error) {
	cDatabaseID, allocation := toNativeString(databaseID)
	defer C.free(allocation)

	var out *C.cosmos_database_ref_t
	status := C.cosmos_database_ref_create(d.account, cDatabaseID, &out) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	if status != 0 {
		httpStatus, subStatus := unpackStatus(status)
		return nil, &Error{
			Code:       codeForRichError(false, httpStatus, subStatus),
			StatusCode: httpStatus,
			SubStatus:  subStatus,
			Message:    "azcosmos: creating the database reference",
		}
	}
	return out, nil
}

// buildNativeManagementRequest converts a management operation into the borrowed C request passed
// to submit. The returned function releases every allocation after the driver has copied the
// request.
func buildNativeManagementRequest(
	kind operationKind,
	account *C.cosmos_account_ref_t,
	database *C.cosmos_database_ref_t,
	container *C.cosmos_container_ref_t,
	body []byte,
	operation OperationOptions,
) (C.cosmos_operation_request_t, func()) {
	request := C.cosmos_operation_request_t{
		kind:           C.int32_t(kind),
		account:        account,
		database:       database,
		container:      container,
		max_item_count: unsetMaxItemCount,
	}
	var releases []func()

	if len(body) > 0 {
		// Go memory may not be passed to C when it can hold a Go pointer, and the driver copies
		// the bytes before returning, so a C buffer is both required and cheap here.
		cBody := C.CBytes(body)
		releases = append(releases, func() { C.free(cBody) })
		request.body = (*C.uint8_t)(cBody)
		request.body_len = C.uintptr_t(len(body))
	}

	nativeOptions, freeOptions := operation.toNative()
	releases = append(releases, freeOptions)
	request.options = nativeOptions

	return request, func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
}

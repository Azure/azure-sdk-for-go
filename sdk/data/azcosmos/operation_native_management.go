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
	"encoding/json"
	"errors"
	"strings"
	"unsafe"
)

type managementOperationScope uint8

const (
	managementOperationScopeInvalid managementOperationScope = iota
	managementOperationScopeAccount
	managementOperationScopeDatabase
	managementOperationScopeContainer
)

func scopeForManagementOperation(kind operationKind) managementOperationScope {
	switch kind {
	case operationKindCreateDatabase, operationKindQueryDatabases, operationKindQueryOffers, operationKindReadOffer, operationKindReplaceOffer:
		return managementOperationScopeAccount
	case operationKindReadDatabase, operationKindDeleteDatabase, operationKindCreateContainer, operationKindQueryContainers:
		return managementOperationScopeDatabase
	case operationKindReadContainer, operationKindDeleteContainer, operationKindReplaceContainer:
		return managementOperationScopeContainer
	default:
		return managementOperationScopeInvalid
	}
}

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
	result, err := c.executeManagementRequest(ctx, managementRequest{kind: kind, databaseID: databaseID, containerID: containerID, body: body, operation: operation})
	return result.response, result.body, err
}

func (c *Client) executeManagementRequest(ctx context.Context, req managementRequest) (managementResult, error) {
	var err error
	req, err = c.prepareManagementThroughput(req)
	if err != nil {
		return managementResult{}, err
	}
	d := c.driver
	ctx, snapshot, releaseSnapshot, err := d.snapshot(ctx, req.operation)
	if err != nil {
		return managementResult{}, err
	}
	defer releaseSnapshot()

	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return managementResult{}, err
	}

	var account *C.cosmos_account_ref_t
	var database *C.cosmos_database_ref_t
	var container *C.cosmos_container_ref_t

	switch scopeForManagementOperation(req.kind) {
	case managementOperationScopeAccount:
		account = d.account
	case managementOperationScopeDatabase:
		database, err = d.createDatabaseRef(req.databaseID)
		if err != nil {
			return managementResult{}, err
		}
		defer C.cosmos_database_ref_free(database)
	case managementOperationScopeContainer:
		container, err = d.submitResolveContainer(ctx, driver, req.databaseID, req.containerID)
		if err != nil {
			return managementResult{}, err
		}
		defer C.cosmos_container_ref_free(container)
	}

	submit := func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		request, release := buildNativeManagementRequest(req.kind, account, database, container, req.body, req.operation)
		defer release()
		if req.resourceLink != "" {
			value, allocation := toNativeString(req.resourceLink)
			defer C.free(allocation)
			request.resource_link = value
		}
		if req.precondition.kind != preconditionKindNone {
			value, allocation := toNativeString(string(req.precondition.etag))
			defer C.free(allocation)
			request.precondition_kind = C.int32_t(req.precondition.kind)
			request.precondition_etag = value
		}
		request.options_snapshot = snapshot
		return submitNativeManagement(driver, &request, req.throughput, queue, cookie, preError)
	}

	result, err := d.awaitCompletion(ctx, "submitting the operation", submit)
	if err != nil {
		return managementResult{}, err
	}
	if result.err != nil {
		return managementResult{}, result.err
	}
	return managementResult{response: result.response, body: result.body, pending: result.httpStatus == 202 || result.offerReplacePending}, nil
}

func (c *Client) openManagementQuery(ctx context.Context, req managementRequest, options ManagementQueryOptions) (queryCursor, error) {
	d := c.driver
	ctx, _, releaseSnapshot, err := d.snapshot(ctx, options.Operation)
	if err != nil {
		return nil, err
	}
	defer releaseSnapshot()
	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return nil, err
	}
	var database *C.cosmos_database_ref_t
	if scopeForManagementOperation(req.kind) == managementOperationScopeDatabase {
		database, err = d.createDatabaseRef(req.databaseID)
		if err != nil {
			return nil, err
		}
		defer C.cosmos_database_ref_free(database)
	}
	request, release := buildNativeManagementRequest(req.kind, d.account, database, nil, req.body, options.Operation)
	defer release()
	if options.PageSizeHint > 0 {
		request.max_item_count = C.int32_t(options.PageSizeHint)
	}
	if options.ContinuationToken != "" {
		value, allocation := toNativeString(options.ContinuationToken)
		defer C.free(allocation)
		request.continuation_token = value
	}
	cursor, err := d.openCursor(ctx, driver, request)
	if err != nil {
		return nil, err
	}
	switch req.kind {
	case operationKindQueryDatabases:
		cursor.managementEnvelope = "Databases"
	case operationKindQueryContainers:
		cursor.managementEnvelope = "DocumentCollections"
	case operationKindQueryOffers:
		cursor.managementEnvelope = "Offers"
	}
	return cursor, nil
}

func decodeManagementCursorPage(completion *C.cosmos_cursor_completion_t, verbosity DiagnosticsVerbosity, field string) (queryPage, bool, error) {
	if completion.result_kind != 2 || completion.body_kind != 1 {
		return decodeCursorPage(completion, verbosity)
	}
	common := translateCompletionOutcome(&completion.common, verbosity)
	var envelope map[string]json.RawMessage
	err := json.Unmarshal(common.body, &envelope)
	var items []json.RawMessage
	if err == nil {
		value, exists := envelope[field]
		if !exists {
			err = errors.New("management query response is missing its resource array")
		} else {
			err = json.Unmarshal(value, &items)
			if err == nil && items == nil {
				err = errors.New("management query resource array must not be null")
			}
		}
	}
	if err != nil {
		return queryPage{}, false, newManagementSerializationError(common.response, common.body, "management query", err)
	}
	page := queryPage{Response: common.response.Response, SessionToken: common.response.SessionToken, Items: make([][]byte, len(items))}
	for i := range items {
		page.Items[i] = items[i]
	}
	return page, false, nil
}

func managementOfferReplacePending(completion *C.cosmos_completion_t) bool {
	if completion.headers == nil {
		return false
	}
	for _, header := range unsafe.Slice(completion.headers, int(completion.headers_len)) {
		if header.id == C.COSMOS_HEADER_ID_OFFER_REPLACE_PENDING {
			if header.value.kind == C.COSMOS_VALUE_KIND_BOOL {
				return *(*uint8)(unsafe.Pointer(&header.value.payload)) != 0
			}
			return strings.EqualFold(headerString(&header.value), "true")
		}
	}
	return false
}

func managementContinuationToken(completion *C.cosmos_completion_t) string {
	if completion.headers != nil {
		for _, header := range unsafe.Slice(completion.headers, int(completion.headers_len)) {
			if header.id == C.COSMOS_HEADER_ID_CONTINUATION {
				return headerString(&header.value)
			}
		}
	}
	return ""
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

// inspectNativeManagementRequestWithAllScopes supplies a non-nil candidate for every scope so
// tests can verify that request construction keeps only the handle required by the operation.
func inspectNativeManagementRequestWithAllScopes(kind operationKind) (nativeManagementRequest, func()) {
	handle := C.malloc(1)
	request, release := inspectNativeManagementRequest(
		kind,
		(*C.cosmos_account_ref_t)(handle),
		(*C.cosmos_database_ref_t)(handle),
		(*C.cosmos_container_ref_t)(handle),
		nil,
	)
	return request, func() {
		release()
		C.free(handle)
	}
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
		max_item_count: unsetMaxItemCount,
	}
	switch scopeForManagementOperation(kind) {
	case managementOperationScopeAccount:
		request.account = account
	case managementOperationScopeDatabase:
		request.database = database
	case managementOperationScopeContainer:
		request.container = container
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

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
	"unsafe"
)

func (c *Client) openReadMany(ctx context.Context, req *readManyRequest) (queryCursor, error) {
	d := c.driver
	ctx, _, releaseSnapshot, err := d.snapshot(ctx, req.options.Operation)
	if err != nil {
		return nil, err
	}
	defer releaseSnapshot()
	driver, err := d.ensureDriver(ctx)
	if err != nil {
		return nil, err
	}
	container, err := d.submitResolveContainer(ctx, driver, req.databaseID, req.containerID)
	if err != nil {
		return nil, err
	}
	defer C.cosmos_container_ref_free(container)
	request, release := buildNativeReadManyRequest(req, container)
	defer release()
	cursor, err := d.openCursorWith(ctx, func(queue *C.cosmos_completion_queue_t, cookie C.intptr_t, status *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
		return C.cosmos_read_many_open_submit(driver, request, queue, cookie, status)
	})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// buildNativeReadManyRequest allocates the request and its nested arrays in C memory, because the
// identity and parameter records hold pointers that cgo forbids in Go memory.
func buildNativeReadManyRequest(req *readManyRequest, container *C.cosmos_container_ref_t) (*C.cosmos_read_many_request_t, func()) {
	var releases []func()
	release := func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
	request := (*C.cosmos_read_many_request_t)(C.calloc(1, C.size_t(unsafe.Sizeof(C.cosmos_read_many_request_t{}))))
	releases = append(releases, func() { C.free(unsafe.Pointer(request)) })
	C.cosmos_read_many_request_init(request)

	// The read-many open request carries no operation kind of its own.
	request.operation = newOperationRequest(0, container)
	if token := req.options.SessionToken; token != "" {
		value, allocation := toNativeString(string(token))
		releases = append(releases, func() { C.free(allocation) })
		request.operation.session_token = value
	}
	if req.options.PageSizeHint > 0 {
		request.operation.max_item_count = C.int32_t(req.options.PageSizeHint)
	}
	request.operation.max_fan_out = C.uint32_t(req.options.MaxFanOut)
	nativeOptions, freeOptions := req.options.Operation.toNative()
	releases = append(releases, freeOptions)
	request.operation.options = nativeOptions

	selection := req.selection
	request.selection_kind = C.uint32_t(selection.kind)
	if n := len(selection.identities); n > 0 {
		array := (*C.cosmos_read_many_identity_t)(C.calloc(C.size_t(n), C.size_t(unsafe.Sizeof(C.cosmos_read_many_identity_t{}))))
		releases = append(releases, func() { C.free(unsafe.Pointer(array)) })
		identities := unsafe.Slice(array, n)
		for i, identity := range selection.identities {
			components, freeComponents := identity.partitionKey.toNative()
			releases = append(releases, freeComponents)
			identities[i].partition_key = components
			identities[i].partition_key_len = identity.partitionKey.partitionKeyLen()
			if selection.kind == readManySelectionItems {
				value, allocation := toNativeString(identity.id)
				releases = append(releases, func() { C.free(allocation) })
				identities[i].item_id = value
			}
		}
		request.identities = array
		request.identities_len = C.uintptr_t(n)
	}
	if selection.filter != "" {
		value, allocation := toNativeString(selection.filter)
		releases = append(releases, func() { C.free(allocation) })
		request.filter = value
	}
	if n := len(selection.parameters); n > 0 {
		array := (*C.cosmos_read_many_parameter_t)(C.calloc(C.size_t(n), C.size_t(unsafe.Sizeof(C.cosmos_read_many_parameter_t{}))))
		releases = append(releases, func() { C.free(unsafe.Pointer(array)) })
		parameters := unsafe.Slice(array, n)
		for i, parameter := range selection.parameters {
			name, nameAllocation := toNativeString(parameter.Name)
			releases = append(releases, func() { C.free(nameAllocation) })
			value, valueAllocation := toNativeString(string(parameter.Value))
			releases = append(releases, func() { C.free(valueAllocation) })
			parameters[i] = C.cosmos_read_many_parameter_t{name: name, json_value: value}
		}
		request.parameters = array
		request.parameters_len = C.uintptr_t(n)
	}
	return request, release
}

type nativeReadManyRequest struct {
	selectionKind uint32
	identities    int
	itemIDs       []string
	pkLens        []int
	filter        string
	parameters    map[string]string
	pageSize      int32
	fanOut        uint32
	sessionToken  string
	kind          int32
	abiVersion    uint32
	structSize    uint32
}

// inspectNativeReadMany exposes the built native request in Go types for tests.
func inspectNativeReadMany(req *readManyRequest) nativeReadManyRequest {
	request, release := buildNativeReadManyRequest(req, nil)
	defer release()
	result := nativeReadManyRequest{
		selectionKind: uint32(request.selection_kind),
		identities:    int(request.identities_len),
		filter:        fromNativeString(request.filter),
		pageSize:      int32(request.operation.max_item_count),
		fanOut:        uint32(request.operation.max_fan_out),
		sessionToken:  fromNativeString(request.operation.session_token),
		kind:          int32(request.operation.kind),
		abiVersion:    uint32(request.abi_version),
		structSize:    uint32(request.struct_size_bytes),
	}
	if request.identities_len > 0 {
		for _, identity := range unsafe.Slice(request.identities, int(request.identities_len)) {
			result.itemIDs = append(result.itemIDs, fromNativeString(identity.item_id))
			result.pkLens = append(result.pkLens, int(identity.partition_key_len))
		}
	}
	if request.parameters_len > 0 {
		result.parameters = map[string]string{}
		for _, parameter := range unsafe.Slice(request.parameters, int(request.parameters_len)) {
			result.parameters[fromNativeString(parameter.name)] = fromNativeString(parameter.json_value)
		}
	}
	return result
}

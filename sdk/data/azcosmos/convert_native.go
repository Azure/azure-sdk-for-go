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
	"unsafe"
)

// The converters here build the C representations of the option and partition-key types. Each
// returns a release function rather than relying on a finalizer, because the driver borrows the
// memory only for the duration of the submit call and the caller knows exactly when that ends.
//
// C memory is used rather than Go memory throughout: cgo forbids passing Go memory that itself
// holds Go pointers, which a slice of structs containing strings does.

// toNative builds the driver's partition key value. The returned function releases it.
//
// The components are passed inline rather than through cosmos_partition_key_create, because the
// request struct accepts an array directly and that avoids a handle whose lifetime would have to
// be tracked separately.
func (pk PartitionKey) toNative() (*C.cosmos_partition_key_component_t, func()) {
	if len(pk.components) == 0 {
		return nil, func() {}
	}

	size := C.size_t(len(pk.components)) * C.size_t(unsafe.Sizeof(C.cosmos_partition_key_component_t{}))
	array := (*C.cosmos_partition_key_component_t)(C.malloc(size))
	components := unsafe.Slice(array, len(pk.components))

	// Strings are allocated separately and freed by the returned function, since the component
	// only borrows the pointer.
	var strings []unsafe.Pointer

	for i, component := range pk.components {
		components[i] = C.cosmos_partition_key_component_t{kind: C.uint8_t(component.kind)}
		switch component.kind {
		case partitionKeyKindString:
			value := unsafe.Pointer(C.CString(component.stringValue))
			strings = append(strings, value)
			*(**C.char)(unsafe.Pointer(&components[i].value)) = (*C.char)(value)
		case partitionKeyKindNumber:
			*(*C.double)(unsafe.Pointer(&components[i].value)) = C.double(component.numberValue)
		case partitionKeyKindBool:
			*(*C.bool)(unsafe.Pointer(&components[i].value)) = C.bool(component.boolValue)
		case partitionKeyKindNull, partitionKeyKindUndefined:
			// Carry no value; the kind alone says what they are.
		}
	}

	return array, func() {
		for _, value := range strings {
			C.free(value)
		}
		C.free(unsafe.Pointer(array))
	}
}

// partitionKeyLen reports how many components the driver should read from the array toNative built.
func (pk PartitionKey) partitionKeyLen() C.uintptr_t {
	return C.uintptr_t(len(pk.components))
}

// toNative builds the driver's per-operation options. The returned function releases them.
func (o OperationOptions) toNative() (*C.cosmos_operation_options_t, func()) {
	// Starting from the driver's defaults rather than a zero struct, so that every field this
	// package does not set keeps its documented default rather than becoming zero.
	options := (*C.cosmos_operation_options_t)(C.malloc(C.size_t(unsafe.Sizeof(C.cosmos_operation_options_t{}))))
	*options = C.cosmos_operation_options_default()

	var allocations []unsafe.Pointer
	release := func() {
		for _, value := range allocations {
			C.free(value)
		}
		C.free(unsafe.Pointer(options))
	}

	if strategy, ok := o.ConsistencyStrategy.toNative(); ok {
		options.read_consistency_strategy = strategy
	}
	if o.EnableContentResponseOnWrite != nil {
		// Tri-state: 0 unset, 1 false, 2 true.
		if *o.EnableContentResponseOnWrite {
			options.content_response_on_write = 2
		} else {
			options.content_response_on_write = 1
		}
	}
	if o.EndToEndTimeout > 0 {
		options.end_to_end_timeout_ms = C.int64_t(o.EndToEndTimeout.Milliseconds())
	}
	if len(o.ExcludedRegions) > 0 {
		size := C.size_t(len(o.ExcludedRegions)) * C.size_t(unsafe.Sizeof(uintptr(0)))
		array := (**C.char)(C.malloc(size))
		allocations = append(allocations, unsafe.Pointer(array))

		regions := unsafe.Slice(array, len(o.ExcludedRegions))
		for i, region := range o.ExcludedRegions {
			value := C.CString(string(region))
			allocations = append(allocations, unsafe.Pointer(value))
			regions[i] = value
		}
		options.excluded_regions = (**C.char)(array)
		options.excluded_regions_len = C.uintptr_t(len(o.ExcludedRegions))
	}

	return options, release
}

// toNative maps a read consistency strategy onto the driver's discriminant. The second result is
// false for the unset strategy, which leaves the driver's default in place.
func (s ReadConsistencyStrategy) toNative() (C.int32_t, bool) {
	switch s {
	case ReadConsistencyStrategyDefault:
		return C.COSMOS_READ_CONSISTENCY_STRATEGY_DEFAULT, true
	case ReadConsistencyStrategyEventual:
		return C.COSMOS_READ_CONSISTENCY_STRATEGY_EVENTUAL, true
	case ReadConsistencyStrategySession:
		return C.COSMOS_READ_CONSISTENCY_STRATEGY_SESSION, true
	case ReadConsistencyStrategyGlobalStrong:
		return C.COSMOS_READ_CONSISTENCY_STRATEGY_GLOBAL_STRONG, true
	default:
		return 0, false
	}
}

// The inspectors below read back what the converters wrote, in Go types. They exist because cgo is
// not permitted in _test.go files, so a test cannot dereference these structs itself — without
// them the converters could only be tested by observing their effect on a live service.

// nativePartitionKeyComponent is one converted partition key component, in Go types.
type nativePartitionKeyComponent struct {
	kind        uint8
	stringValue string
	numberValue float64
	boolValue   bool
}

// inspectNativePartitionKey converts a partition key and reads the result back.
func inspectNativePartitionKey(pk PartitionKey) ([]nativePartitionKeyComponent, func()) {
	array, release := pk.toNative()
	if array == nil {
		return nil, release
	}

	components := unsafe.Slice(array, len(pk.components))
	out := make([]nativePartitionKeyComponent, len(components))
	for i := range components {
		out[i].kind = uint8(components[i].kind)
		switch components[i].kind {
		case C.COSMOS_PARTITION_KEY_COMPONENT_KIND_STRING:
			if ptr := *(**C.char)(unsafe.Pointer(&components[i].value)); ptr != nil {
				out[i].stringValue = C.GoString(ptr)
			}
		case C.COSMOS_PARTITION_KEY_COMPONENT_KIND_NUMBER:
			out[i].numberValue = float64(*(*C.double)(unsafe.Pointer(&components[i].value)))
		case C.COSMOS_PARTITION_KEY_COMPONENT_KIND_BOOL:
			out[i].boolValue = bool(*(*C.bool)(unsafe.Pointer(&components[i].value)))
		}
	}
	return out, release
}

// The partition key component kinds, in Go types, so tests can assert against the ABI's values
// rather than against literals that could drift from them.
var (
	nativeKindString    = uint8(C.COSMOS_PARTITION_KEY_COMPONENT_KIND_STRING)
	nativeKindNumber    = uint8(C.COSMOS_PARTITION_KEY_COMPONENT_KIND_NUMBER)
	nativeKindBool      = uint8(C.COSMOS_PARTITION_KEY_COMPONENT_KIND_BOOL)
	nativeKindNull      = uint8(C.COSMOS_PARTITION_KEY_COMPONENT_KIND_NULL)
	nativeKindUndefined = uint8(C.COSMOS_PARTITION_KEY_COMPONENT_KIND_UNDEFINED)
)

// nativeOperationOptions is a converted option set, in Go types.
type nativeOperationOptions struct {
	readConsistencyStrategy int32
	contentResponseOnWrite  int32
	endToEndTimeoutMillis   int64
	excludedRegions         []string
}

// inspectNativeOperationOptions converts an option set and reads the result back.
func inspectNativeOperationOptions(o OperationOptions) (nativeOperationOptions, func()) {
	options, release := o.toNative()

	out := nativeOperationOptions{
		readConsistencyStrategy: int32(options.read_consistency_strategy),
		contentResponseOnWrite:  int32(options.content_response_on_write),
		endToEndTimeoutMillis:   int64(options.end_to_end_timeout_ms),
	}
	if options.excluded_regions != nil && options.excluded_regions_len > 0 {
		regions := unsafe.Slice(options.excluded_regions, int(options.excluded_regions_len))
		out.excludedRegions = make([]string, len(regions))
		for i, region := range regions {
			out.excludedRegions[i] = C.GoString(region)
		}
	}
	return out, release
}

// defaultNativeOperationOptions reports the driver's own defaults, so a test can assert that an
// unset field was left alone rather than hardcoding what the default happens to be.
func defaultNativeOperationOptions() nativeOperationOptions {
	defaults := C.cosmos_operation_options_default()
	return nativeOperationOptions{
		readConsistencyStrategy: int32(defaults.read_consistency_strategy),
		contentResponseOnWrite:  int32(defaults.content_response_on_write),
		endToEndTimeoutMillis:   int64(defaults.end_to_end_timeout_ms),
	}
}

// nativeReadConsistencyStrategy reports the discriminant a strategy maps to, in Go types.
func nativeReadConsistencyStrategy(s ReadConsistencyStrategy) (int32, bool) {
	value, ok := s.toNative()
	return int32(value), ok
}

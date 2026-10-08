// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore finalizer

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

// The converters here build the C representations of the option and partition-key types. Each
// returns a release function rather than relying on a finalizer, because the driver borrows the
// memory only for the duration of the submit call and the caller knows exactly when that ends.
//
// C memory is used rather than Go memory throughout: cgo forbids passing Go memory that itself
// holds Go pointers, which a slice of structs containing strings does.

func toNativeString(value string) (C.cosmos_string_view_t, unsafe.Pointer) {
	allocation := unsafe.Pointer(C.CString(value))
	return C.cosmos_string_view_t{
		data: (*C.uint8_t)(allocation),
		len:  C.uintptr_t(len(value)),
	}, allocation
}

func fromNativeString(value C.cosmos_string_view_t) string {
	if value.data == nil {
		return ""
	}
	return string(C.GoBytes(unsafe.Pointer(value.data), C.int(value.len)))
}

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
			value, allocation := toNativeString(component.stringValue)
			strings = append(strings, allocation)
			*(*C.cosmos_string_view_t)(unsafe.Pointer(&components[i].value)) = value
		case partitionKeyKindNumber:
			*(*C.double)(unsafe.Pointer(&components[i].value)) = C.double(component.numberValue)
		case partitionKeyKindBool:
			if component.boolValue {
				*(*C.uint8_t)(unsafe.Pointer(&components[i].value)) = 1
			}
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
	if strategy, ok := o.PatchStrategy.toNative(); ok {
		options.patch_strategy = strategy
	}
	if o.SessionCapturingDisabled != nil {
		options.session_capturing_disabled = nativeBool(*o.SessionCapturingDisabled)
	}
	if o.MaxFailoverRetryCount != nil {
		options.max_failover_retry_count = C.int64_t(*o.MaxFailoverRetryCount)
	}
	if o.MaxSessionRetryCount != nil {
		options.max_session_retry_count = C.int64_t(*o.MaxSessionRetryCount)
	}
	if o.EndpointUnavailabilityTTL != nil {
		options.endpoint_unavailability_ttl_ms = C.int64_t(o.EndpointUnavailabilityTTL.Milliseconds())
	}
	if o.BinaryEncoding != nil {
		options.binary_encoding_enabled = nativeBool(o.BinaryEncoding.enabled())
		options.binary_encoding_request_text_response = nativeBool(o.BinaryEncoding.enabled())
	}
	if o.ThroughputControl.ThroughputBucket != nil {
		options.throughput_bucket = C.int64_t(*o.ThroughputControl.ThroughputBucket)
	}
	switch o.ThroughputControl.PriorityLevel {
	case PriorityLevelHigh:
		options.priority_level = 1
	case PriorityLevelLow:
		options.priority_level = 2
	}
	if o.ThrottlingRetry.MaxRetryCount != nil {
		options.max_throttle_retry_count = C.int64_t(*o.ThrottlingRetry.MaxRetryCount)
	}
	if o.ThrottlingRetry.MaxRetryWaitTime != nil {
		options.max_throttle_retry_wait_time_ms = C.int64_t(o.ThrottlingRetry.MaxRetryWaitTime.Milliseconds())
	}
	if o.HedgingEnabled != nil {
		options.hedging_enabled = nativeBool(*o.HedgingEnabled)
	}
	options.availability_strategy = C.int32_t(o.AvailabilityStrategy.kind)
	if o.AvailabilityStrategy.kind == 2 {
		options.hedge_threshold_ms = C.int64_t(max(1, o.AvailabilityStrategy.threshold.Milliseconds()))
	}
	if o.CustomHeaders != nil {
		array := (*C.cosmos_header_kv_t)(C.malloc(C.size_t(max(1, len(o.CustomHeaders))) * C.size_t(unsafe.Sizeof(C.cosmos_header_kv_t{}))))
		allocations = append(allocations, unsafe.Pointer(array))
		headers := unsafe.Slice(array, len(o.CustomHeaders))
		i := 0
		for name, value := range o.CustomHeaders {
			n, na := toNativeString(strings.ToLower(name))
			v, va := toNativeString(value)
			allocations = append(allocations, na, va)
			headers[i] = C.cosmos_header_kv_t{name: n, value: v}
			i++
		}
		options.custom_headers = array
		options.custom_headers_len = C.uintptr_t(len(headers))
	}
	if o.EnableContentResponseOnWrite != nil {
		// Tri-state: 0 unset, 1 false, 2 true.
		if *o.EnableContentResponseOnWrite {
			options.content_response_on_write = 2
		} else {
			options.content_response_on_write = 1
		}

	}
	if o.EndToEndTimeout != nil {
		milliseconds := max(o.EndToEndTimeout.Milliseconds(), 1000)
		options.end_to_end_timeout_ms = C.int64_t(milliseconds)
	}
	if o.ExcludedRegions != nil {
		size := C.size_t(max(1, len(o.ExcludedRegions))) * C.size_t(unsafe.Sizeof(C.cosmos_string_view_t{}))
		array := (*C.cosmos_string_view_t)(C.malloc(size))
		allocations = append(allocations, unsafe.Pointer(array))

		regions := unsafe.Slice(array, len(o.ExcludedRegions))
		for i, region := range o.ExcludedRegions {
			value, allocation := toNativeString(string(region))
			allocations = append(allocations, allocation)
			regions[i] = value
		}
		options.excluded_regions = array
		options.excluded_regions_len = C.uintptr_t(len(o.ExcludedRegions))
	}

	return options, release
}

func nativeBool(value bool) C.int8_t {
	if value {
		return 2
	}
	return 1
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
	case ReadConsistencyStrategyLatestCommitted:
		return C.COSMOS_READ_CONSISTENCY_STRATEGY_LATEST_COMMITTED, true
	case ReadConsistencyStrategyGlobalStrong:
		return C.COSMOS_READ_CONSISTENCY_STRATEGY_GLOBAL_STRONG, true
	default:
		return 0, false
	}
}

// toNative maps a PATCH strategy onto the driver's discriminant. The second result is false for
// the unset strategy, which leaves the driver's default in place.
func (s PatchStrategy) toNative() (C.int32_t, bool) {
	switch s {
	case PatchStrategyAuto:
		return C.COSMOS_PATCH_STRATEGY_AUTO, true
	case PatchStrategyClientSide:
		return C.COSMOS_PATCH_STRATEGY_CLIENT_SIDE, true
	case PatchStrategyServerSide:
		return C.COSMOS_PATCH_STRATEGY_SERVER_SIDE, true
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
			value := *(*C.cosmos_string_view_t)(unsafe.Pointer(&components[i].value))
			out[i].stringValue = fromNativeString(value)
		case C.COSMOS_PARTITION_KEY_COMPONENT_KIND_NUMBER:
			out[i].numberValue = float64(*(*C.double)(unsafe.Pointer(&components[i].value)))
		case C.COSMOS_PARTITION_KEY_COMPONENT_KIND_BOOL:
			out[i].boolValue = *(*C.uint8_t)(unsafe.Pointer(&components[i].value)) != 0
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
	binaryEncodingEnabled      int8
	binaryTextResponse         int8
	readConsistencyStrategy    int32
	contentResponseOnWrite     int32
	endToEndTimeoutMillis      int64
	excludedRegions            []string
	patchStrategy              int32
	sessionCapturingDisabled   int8
	maxFailoverRetryCount      int64
	maxSessionRetryCount       int64
	endpointTTLMillis          int64
	customHeaders              map[string]string
	throughputBucket           int64
	priorityLevel              int32
	maxThrottleRetryCount      int64
	maxThrottleRetryWaitMillis int64
	hedgingEnabled             int8
	availabilityStrategy       int32
	hedgeThresholdMillis       int64
}

// inspectNativeOperationOptions converts an option set and reads the result back.
func inspectNativeOperationOptions(o OperationOptions) (nativeOperationOptions, func()) {
	options, release := o.toNative()
	return readNativeOperationOptions(options), release
}

func readNativeOperationOptions(options *C.cosmos_operation_options_t) nativeOperationOptions {
	out := nativeOperationOptions{
		binaryEncodingEnabled:      int8(options.binary_encoding_enabled),
		binaryTextResponse:         int8(options.binary_encoding_request_text_response),
		readConsistencyStrategy:    int32(options.read_consistency_strategy),
		contentResponseOnWrite:     int32(options.content_response_on_write),
		endToEndTimeoutMillis:      int64(options.end_to_end_timeout_ms),
		patchStrategy:              int32(options.patch_strategy),
		sessionCapturingDisabled:   int8(options.session_capturing_disabled),
		maxFailoverRetryCount:      int64(options.max_failover_retry_count),
		maxSessionRetryCount:       int64(options.max_session_retry_count),
		endpointTTLMillis:          int64(options.endpoint_unavailability_ttl_ms),
		throughputBucket:           int64(options.throughput_bucket),
		priorityLevel:              int32(options.priority_level),
		maxThrottleRetryCount:      int64(options.max_throttle_retry_count),
		maxThrottleRetryWaitMillis: int64(options.max_throttle_retry_wait_time_ms),
		hedgingEnabled:             int8(options.hedging_enabled),
		availabilityStrategy:       int32(options.availability_strategy),
		hedgeThresholdMillis:       int64(options.hedge_threshold_ms),
	}
	if options.excluded_regions != nil {
		regions := unsafe.Slice(options.excluded_regions, int(options.excluded_regions_len))
		out.excludedRegions = make([]string, len(regions))
		for i, region := range regions {
			out.excludedRegions[i] = fromNativeString(region)
		}
	}
	if options.custom_headers != nil {
		out.customHeaders = make(map[string]string)
		for _, h := range unsafe.Slice(options.custom_headers, int(options.custom_headers_len)) {
			out.customHeaders[fromNativeString(h.name)] = fromNativeString(h.value)
		}
	}
	return out
}

// defaultNativeOperationOptions reports the driver's own defaults, so a test can assert that an
// unset field was left alone rather than hardcoding what the default happens to be.
func defaultNativeOperationOptions() nativeOperationOptions {
	defaults := C.cosmos_operation_options_default()
	return readNativeOperationOptions(&defaults)
}

// nativeReadConsistencyStrategy reports the discriminant a strategy maps to, in Go types.
func nativeReadConsistencyStrategy(s ReadConsistencyStrategy) (int32, bool) {
	value, ok := s.toNative()
	return int32(value), ok
}

// nativePatchStrategy reports the discriminant a strategy maps to, in Go types.
func nativePatchStrategy(s PatchStrategy) (int32, bool) {
	value, ok := s.toNative()
	return int32(value), ok
}

// toNative builds the driver's per-client options config. The returned function releases it.
//
// The config is flat: preferred regions plus the operation options every operation starts from.
// Application identity is not here because the ABI carries the user agent on the runtime;
// see nativeDriver.buildRuntime.
func (o ClientOptions) toNative() (*C.cosmos_driver_options_config_t, func(), error) {
	config := (*C.cosmos_driver_options_config_t)(C.malloc(C.size_t(unsafe.Sizeof(C.cosmos_driver_options_config_t{}))))
	*config = C.cosmos_driver_options_config_default()

	operationOptions, releaseOperationOptions := o.Operation.toNative()
	config.operation_options = operationOptions

	var allocations []unsafe.Pointer
	release := func() {
		for _, value := range allocations {
			C.free(value)
		}
		releaseOperationOptions()
		C.free(unsafe.Pointer(config))
	}

	regions, err := o.Routing.preferredRegionOrder()
	if err != nil {
		release()
		return nil, func() {}, err
	}
	if len(regions) > 0 {
		size := C.size_t(len(regions)) * C.size_t(unsafe.Sizeof(C.cosmos_string_view_t{}))
		array := (*C.cosmos_string_view_t)(C.malloc(size))
		allocations = append(allocations, unsafe.Pointer(array))

		values := unsafe.Slice(array, len(regions))
		for i, region := range regions {
			value, allocation := toNativeString(string(region))
			allocations = append(allocations, allocation)
			values[i] = value
		}
		config.preferred_regions = array
		config.preferred_regions_len = C.uintptr_t(len(regions))
	}

	return config, release, nil
}

// nativeClientOptions is a converted client option set, in Go types.
type nativeClientOptions struct {
	preferredRegions []string
	operationOptions nativeOperationOptions
}

// inspectNativeClientOptions converts a client option set and reads the result back.
func inspectNativeClientOptions(o ClientOptions) (nativeClientOptions, func(), error) {
	config, release, err := o.toNative()
	if err != nil {
		return nativeClientOptions{}, func() {}, err
	}

	out := nativeClientOptions{}
	if config.preferred_regions != nil && config.preferred_regions_len > 0 {
		regions := unsafe.Slice(config.preferred_regions, int(config.preferred_regions_len))
		out.preferredRegions = make([]string, len(regions))
		for i, region := range regions {
			out.preferredRegions[i] = fromNativeString(region)
		}
	}
	if config.operation_options != nil {
		out.operationOptions = readNativeOperationOptions(config.operation_options)
	}
	return out, release, nil
}

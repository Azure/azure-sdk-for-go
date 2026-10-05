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
	"time"
	"unsafe"
)

// inspectFaultInjectionOptionsBuild exercises native validation without network I/O.
// Tests cannot import C, so this helper also releases the built options.
func (d *nativeDriver) inspectFaultInjectionOptionsBuild() error {
	options, err := d.buildDriverOptions()
	if err != nil {
		return err
	}
	C.cosmos_driver_options_free(options)
	return nil
}

func nativeOptionalMilliseconds(duration *time.Duration) C.int64_t {
	if duration == nil {
		return -1
	}
	return C.int64_t(duration.Milliseconds())
}

func nativeFaultInjectionRules(rules []FaultInjectionRule) (*C.cosmos_fault_injection_rule_t, func()) {
	if len(rules) == 0 {
		return nil, func() {}
	}
	array := (*C.cosmos_fault_injection_rule_t)(C.malloc(C.size_t(len(rules)) * C.size_t(unsafe.Sizeof(C.cosmos_fault_injection_rule_t{}))))
	releases := []func(){func() { C.free(unsafe.Pointer(array)) }}
	stringView := func(value string) C.cosmos_string_view_t {
		if value == "" {
			return C.cosmos_string_view_t{}
		}
		view, allocation := toNativeString(value)
		releases = append(releases, func() { C.free(allocation) })
		return view
	}
	for i, r := range rules {
		native := &unsafe.Slice(array, len(rules))[i]
		*native = C.cosmos_fault_injection_rule_default()
		native.id = stringView(r.ID)
		native.start_delay_ms = nativeOptionalMilliseconds(r.StartDelay)
		native.expire_after_ms = nativeOptionalMilliseconds(r.ExpireAfter)
		if r.HitLimit != nil {
			native.hit_limit = C.int64_t(*r.HitLimit)
		}
		condition := (*C.cosmos_fault_injection_condition_t)(C.malloc(C.size_t(unsafe.Sizeof(C.cosmos_fault_injection_condition_t{}))))
		*condition = C.cosmos_fault_injection_condition_default()
		releases = append(releases, func() { C.free(unsafe.Pointer(condition)) })
		condition.operation_type = C.int32_t(r.Condition.Operation)
		condition.transport_kind = C.int32_t(r.Condition.Transport)
		condition.region = stringView(string(r.Condition.Region))
		condition.container_id = stringView(r.Condition.ContainerID)
		native.condition = condition
		result := (*C.cosmos_fault_injection_result_t)(C.malloc(C.size_t(unsafe.Sizeof(C.cosmos_fault_injection_result_t{}))))
		*result = C.cosmos_fault_injection_result_default()
		releases = append(releases, func() { C.free(unsafe.Pointer(result)) })
		result.error_type = C.int32_t(r.Result.Error)
		result.delay_ms = nativeOptionalMilliseconds(r.Result.Delay)
		result.retry_after_ms = nativeOptionalMilliseconds(r.Result.RetryAfter)
		result.custom_status_code = C.int32_t(r.Result.CustomStatusCode)
		if r.Result.CustomSubStatus != nil {
			result.custom_sub_status = C.int32_t(*r.Result.CustomSubStatus)
		}
		if r.Result.Probability != nil {
			result.probability = C.float(*r.Result.Probability)
		}
		// A non-null pointer selects a custom response, even when its header count is zero.
		if len(r.Result.Headers) != 0 {
			headers, releaseHeaders := (OperationOptions{CustomHeaders: r.Result.Headers}).toNative()
			releases = append(releases, releaseHeaders)
			result.custom_headers = headers.custom_headers
			result.custom_headers_len = headers.custom_headers_len
		}
		if len(r.Result.Body) != 0 {
			body := C.CBytes(r.Result.Body)
			releases = append(releases, func() { C.free(body) })
			result.body = (*C.uint8_t)(body)
			result.body_len = C.uintptr_t(len(r.Result.Body))
		}
		native.result = result
	}
	return array, func() {
		for i := len(releases) - 1; i >= 0; i-- {
			releases[i]()
		}
	}
}

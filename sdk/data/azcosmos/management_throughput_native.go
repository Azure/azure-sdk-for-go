// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && azcosmos_native_preview && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include "azurecosmosdriver.h"
*/
import "C"

const nativeManagementThroughputAvailable = true

func (c *Client) prepareManagementThroughput(req managementRequest) (managementRequest, error) {
	return req, nil
}

func submitNativeManagement(driver *C.cosmos_driver_t, request *C.cosmos_operation_request_t, throughput *ThroughputProperties, queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
	if throughput == nil {
		return C.cosmos_submit_singleton_operation(driver, request, queue, cookie, preError)
	}
	var options C.cosmos_management_options_t
	C.cosmos_management_options_init(&options)
	if throughput.ManualThroughput != nil {
		options.throughput_kind = 1
		options.manual_throughput = C.uint64_t(*throughput.ManualThroughput)
	} else {
		options.throughput_kind = 2
		options.autoscale_max_throughput = C.uint64_t(throughput.Autoscale.MaxThroughput)
		if upgrade := throughput.Autoscale.AutoUpgradePolicy; upgrade != nil {
			options.has_auto_upgrade_policy = 1
			if policy := upgrade.ThroughputPolicy; policy != nil {
				options.has_throughput_policy = 1
				options.increment_percent = C.uint32_t(policy.IncrementPercent)
			}
		}
	}
	return C.cosmos_submit_singleton_operation_v2(driver, request, &options, queue, cookie, preError)
}

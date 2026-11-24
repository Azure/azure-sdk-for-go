// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && !azcosmos_native_preview && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"encoding/json"
	"strconv"
)

const nativeManagementThroughputAvailable = false

func (c *Client) prepareManagementThroughput(req managementRequest) (managementRequest, error) {
	if req.throughput == nil {
		return req, nil
	}
	if req.operation.CustomHeaders == nil {
		headers := c.options.Operation.CustomHeaders
		inherited := len(headers) > 0
		if headers == nil && c.options.Runtime != nil {
			inherited = c.options.Runtime.hasCustomHeaders
		}
		if inherited {
			return managementRequest{}, &Error{Code: CodeClientError, Message: "azcosmos: typed creation throughput with inherited custom headers requires an unpublished native extension and matching native preview artifacts"}
		}
	}
	req.operation = req.operation.clone()
	if req.operation.CustomHeaders == nil {
		req.operation.CustomHeaders = make(map[string]string)
	}
	if req.throughput.ManualThroughput != nil {
		req.operation.CustomHeaders["x-ms-offer-throughput"] = strconv.FormatUint(*req.throughput.ManualThroughput, 10)
	} else {
		body, err := json.Marshal(req.throughput.Autoscale)
		if err != nil {
			return managementRequest{}, err
		}
		req.operation.CustomHeaders["x-ms-cosmos-offer-autopilot-settings"] = string(body)
	}
	return req, nil
}

func submitNativeManagement(driver *C.cosmos_driver_t, request *C.cosmos_operation_request_t, _ *ThroughputProperties, queue *C.cosmos_completion_queue_t, cookie C.intptr_t, preError *C.cosmos_status_code_t) *C.cosmos_operation_handle_t {
	return C.cosmos_submit_singleton_operation(driver, request, queue, cookie, preError)
}

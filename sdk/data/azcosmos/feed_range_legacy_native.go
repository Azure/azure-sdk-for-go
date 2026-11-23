// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !azcosmos_native_preview && cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include "azurecosmosdriver.h"
*/
import "C"

func translateFeedRanges(*C.cosmos_completion_t, *completionResult) {}

func bindNativeFeedRange(_ FeedRange, request C.cosmos_operation_request_t, release func()) (C.cosmos_operation_request_t, func(), C.cosmos_status_code_t) {
	// Published builds reject range scopes during argument validation.
	return request, release, C.cosmos_status_code_t((400 << 16) | C.COSMOS_SUB_STATUS_CLIENT_FFI_INVALID_OPTION_VALUE)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !azcosmos_native_preview && cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include "azurecosmosdriver.h"
*/
import "C"

func translateFeedRanges(*C.cosmos_completion_t, *completionResult) {}

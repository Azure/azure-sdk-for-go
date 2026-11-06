//go:build linux && cgo

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

/*
#include <stdint.h>
#include <stdlib.h>
uint64_t cosmos_perf_flow_phase(uint32_t phase);
uint64_t cosmos_perf_flow_epoch_seconds(void);
int32_t cosmos_perf_flow_write(const char *path);
*/
import "C"

import "unsafe"

func flowNativeEpoch() uint64 {
	return uint64(C.cosmos_perf_flow_epoch_seconds())
}

func flowNativePhase(phase uint32) uint64 {
	return uint64(C.cosmos_perf_flow_phase(C.uint32_t(phase)))
}

func flowNativeWrite(path string) int {
	value := C.CString(path)
	defer C.free(unsafe.Pointer(value))
	return int(C.cosmos_perf_flow_write(value))
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include <stdlib.h>
#include <string.h>
#include "azurecosmosdriver.h"

typedef struct {
	char *endpoint;
	char *region;
	uint16_t status;
	int32_t sub_status;
	uint64_t latency_ms;
	double request_charge;
	double server_duration_ms;
} cosmos_go_attempt_t;

typedef struct {
	cosmos_go_attempt_t *entries;
	uintptr_t count;
	uintptr_t capacity;
	int failed;
} cosmos_go_attempts_t;

typedef struct {
	char **entries;
	uintptr_t count;
	uintptr_t capacity;
	int failed;
} cosmos_go_regions_t;

static void cosmos_go_count_region(void *user_data, const char *region) {
	(void)region;
	(*(uintptr_t *)user_data)++;
}

static void cosmos_go_collect_region(void *user_data, const char *region) {
	cosmos_go_regions_t *out = user_data;
	if (out->count < out->capacity) {
		char *copy = region == NULL ? NULL : strdup(region);
		if (region != NULL && copy == NULL) {
			out->failed = 1;
		}
		out->entries[out->count++] = copy;
	}
}

static void cosmos_go_collect_attempt(void *user_data, const char *endpoint,
		const char *region, uint16_t status, int32_t sub_status,
		uint64_t latency_ms, double request_charge, double server_duration_ms) {
	cosmos_go_attempts_t *out = user_data;
	if (out->count < out->capacity) {
		char *endpoint_copy = endpoint == NULL ? NULL : strdup(endpoint);
		char *region_copy = region == NULL ? NULL : strdup(region);
		if ((endpoint != NULL && endpoint_copy == NULL) ||
			(region != NULL && region_copy == NULL)) {
			out->failed = 1;
		}
		out->entries[out->count++] = (cosmos_go_attempt_t){
			endpoint_copy, region_copy, status, sub_status, latency_ms,
			request_charge, server_duration_ms
		};
	}
}

static uintptr_t cosmos_go_region_count(const cosmos_diagnostics_t *d) {
	uintptr_t count = 0;
	cosmos_diagnostics_iter_regions_contacted(d, cosmos_go_count_region, &count);
	return count;
}

static void cosmos_go_read_regions(const cosmos_diagnostics_t *d, cosmos_go_regions_t *out) {
	cosmos_diagnostics_iter_regions_contacted(d, cosmos_go_collect_region, out);
}

static void cosmos_go_read_attempts(const cosmos_diagnostics_t *d, cosmos_go_attempts_t *out) {
	cosmos_diagnostics_iter_attempts(d, cosmos_go_collect_attempt, out);
}

static void cosmos_go_free_regions(cosmos_go_regions_t *out) {
	for (uintptr_t i = 0; i < out->count; i++) {
		free(out->entries[i]);
	}
	free(out->entries);
}

static void cosmos_go_free_attempts(cosmos_go_attempts_t *out) {
	for (uintptr_t i = 0; i < out->count; i++) {
		free(out->entries[i].endpoint);
		free(out->entries[i].region);
	}
	free(out->entries);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

// copyDiagnostics copies the native diagnostics for one completion into Go memory.
//
// It is best effort: diagnostics are optional and must never block an otherwise successful or
// cancelled operation, so a failure copying one section (regions, attempts, or the JSON
// rendering) does not discard the scalar fields or any other section that copied cleanly. The
// returned Diagnostics is non-nil whenever d is non-nil; the returned error, if any, joins every
// section that could not be copied.
func copyDiagnostics(d *C.cosmos_diagnostics_t) (*Diagnostics, error) {
	if d == nil {
		return nil, nil
	}
	out := &Diagnostics{
		AttemptCount:       uint32(C.cosmos_diagnostics_request_count(d)),
		TotalRequestCharge: float64(C.cosmos_diagnostics_total_request_charge(d)),
		Elapsed:            durationMicros(uint64(C.cosmos_diagnostics_total_elapsed_micros(d))),
		Completed:          bool(C.cosmos_diagnostics_is_completed(d)),
		Failed:             bool(C.cosmos_diagnostics_is_failure(d)),
		Compacted:          bool(C.cosmos_diagnostics_is_compacted(d)),
	}

	return out, errors.Join(
		copyDiagnosticRegions(d, out),
		copyDiagnosticAttempts(d, out),
		copyDiagnosticJSON(d, out),
	)
}

// copyDiagnosticRegions copies the regions-contacted list into out.RegionsContacted. It leaves
// that field nil on failure rather than touching fields the other sections own.
func copyDiagnosticRegions(d *C.cosmos_diagnostics_t, out *Diagnostics) error {
	count := int(C.cosmos_go_region_count(d))
	if count == 0 {
		return nil
	}
	ptr := C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof((*C.char)(nil))))
	if ptr == nil {
		return fmt.Errorf("azcosmos: allocating diagnostic regions")
	}
	regions := C.cosmos_go_regions_t{entries: (**C.char)(ptr), capacity: C.uintptr_t(count)}
	defer C.cosmos_go_free_regions(&regions) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	C.cosmos_go_read_regions(d, &regions)    //nolint:gocritic // dupSubExpr targets cgo-generated code.
	if regions.failed != 0 {
		return fmt.Errorf("azcosmos: copying diagnostic regions")
	}
	out.RegionsContacted = make([]string, int(regions.count))
	for i, region := range unsafe.Slice(regions.entries, int(regions.count)) {
		out.RegionsContacted[i] = C.GoString(region)
	}
	return nil
}

// copyDiagnosticAttempts copies the retained per-attempt timeline into out.Attempts. It leaves
// that field nil on failure rather than touching fields the other sections own.
func copyDiagnosticAttempts(d *C.cosmos_diagnostics_t, out *Diagnostics) error {
	count := int(C.cosmos_diagnostics_retained_request_count(d))
	if count == 0 {
		return nil
	}
	ptr := C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.cosmos_go_attempt_t{})))
	if ptr == nil {
		return fmt.Errorf("azcosmos: allocating diagnostic attempts")
	}
	attempts := C.cosmos_go_attempts_t{entries: (*C.cosmos_go_attempt_t)(ptr), capacity: C.uintptr_t(count)}
	defer C.cosmos_go_free_attempts(&attempts) //nolint:gocritic // dupSubExpr targets cgo-generated code.
	C.cosmos_go_read_attempts(d, &attempts)    //nolint:gocritic // dupSubExpr targets cgo-generated code.
	if attempts.failed != 0 {
		return fmt.Errorf("azcosmos: copying diagnostic attempts")
	}
	out.Attempts = make([]DiagnosticAttempt, int(attempts.count))
	for i, attempt := range unsafe.Slice(attempts.entries, int(attempts.count)) {
		out.Attempts[i] = DiagnosticAttempt{
			Endpoint:         C.GoString(attempt.endpoint),
			Region:           C.GoString(attempt.region),
			StatusCode:       int(attempt.status),
			SubStatus:        normalizeSubStatus(int32(attempt.sub_status)),
			Latency:          durationMillis(uint64(attempt.latency_ms)),
			RequestCharge:    float64(attempt.request_charge),
			ServerDurationMS: float64(attempt.server_duration_ms),
		}
	}
	return nil
}

// copyDiagnosticJSON copies the driver's detailed diagnostics rendering into out.JSON. It leaves
// that field empty on failure rather than touching fields the other sections own.
func copyDiagnosticJSON(d *C.cosmos_diagnostics_t, out *Diagnostics) error {
	var data *C.uint8_t
	var length C.uintptr_t
	if status := C.cosmos_diagnostics_to_json(d, C.cosmos_diagnostics_verbosity_t_DETAILED, &data, &length); status != 0 { //nolint:gocritic // dupSubExpr targets cgo-generated code.
		return fmt.Errorf("azcosmos: rendering native diagnostics (status %d)", int32(status))
	}
	if length == 0 {
		return nil
	}
	if data == nil {
		return fmt.Errorf("azcosmos: native diagnostics returned a null JSON buffer")
	}
	out.JSON = string(unsafe.Slice((*byte)(unsafe.Pointer(data)), int(length)))
	return nil
}

func durationMicros(value uint64) time.Duration {
	const maxDuration = uint64(1<<63 - 1)
	if value > maxDuration/uint64(time.Microsecond) {
		return time.Duration(maxDuration)
	}
	return time.Duration(value) * time.Microsecond
}

func durationMillis(value uint64) time.Duration {
	const maxDuration = uint64(1<<63 - 1)
	if value > maxDuration/uint64(time.Millisecond) {
		return time.Duration(maxDuration)
	}
	return time.Duration(value) * time.Millisecond
}

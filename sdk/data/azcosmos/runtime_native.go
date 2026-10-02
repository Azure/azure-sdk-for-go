// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

// cSpell:ignore gocritic

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"
*/
import "C"

import (
	"context"
	"time"
)

type nativeRuntime struct {
	handle *C.cosmos_runtime_t
}

func openRuntime(config RuntimeOptions) (*nativeRuntime, error) {
	if err := verifyDriverVersion(); err != nil {
		return nil, err
	}
	options := C.cosmos_runtime_options_default()
	identifier, allocation := toNativeString(wrappingSDKIdentifier())
	defer C.free(allocation)
	options.wrapping_sdk_identifier = identifier
	common, release := config.Operation.toNative()
	defer release()
	options.operation_options = common
	if config.CPURefreshInterval != nil {
		options.cpu_refresh_interval_ms = C.uint64_t(config.CPURefreshInterval.Milliseconds())
	}
	if config.ApplicationID != "" {
		value, allocation := toNativeString(config.ApplicationID)
		defer C.free(allocation)
		options.user_agent_suffix = value
	}
	native := &nativeRuntime{}
	var richErr *C.cosmos_error_t
	status := C.cosmos_runtime_build(&options, &native.handle, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code.
	if err := statusError(status, richErr, "building runtime (check SDK identity, ApplicationID and Operation options)"); err != nil {
		return nil, err
	}
	return native, nil
}

func (r *nativeRuntime) close() {
	if r != nil {
		C.cosmos_runtime_free(r.handle)
		r.handle = nil
	}
}

func (d *nativeDriver) snapshot(ctx context.Context, request OperationOptions) (context.Context, *C.cosmos_operation_options_snapshot_t, func(), error) {
	client, releaseClient := d.cfg.options.Operation.toNative()
	defer releaseClient()
	options, releaseOptions := request.toNative()
	defer releaseOptions()
	var snapshot *C.cosmos_operation_options_snapshot_t
	var timeout C.int64_t
	status := C.cosmos_operation_options_snapshot_create(d.runtime, client, options, &snapshot, &timeout) //nolint:gocritic // dupSubExpr is reported against cgo-generated code.
	if err := statusError(status, nil, "capturing operation options"); err != nil {
		return ctx, nil, func() {}, err
	}
	var budget time.Duration
	if timeout > 0 {
		// Native milliseconds can exceed Go's duration range through environment configuration.
		const maxMillis = int64((1<<63 - 1) / time.Millisecond)
		budget = time.Duration(min(int64(timeout), maxMillis)) * time.Millisecond
	}
	ctx, cancelSnapshot := contextWithEndToEndTimeout(ctx, budget)
	return ctx, snapshot, func() {
		cancelSnapshot()
		C.cosmos_operation_options_snapshot_free(snapshot)
	}, nil
}

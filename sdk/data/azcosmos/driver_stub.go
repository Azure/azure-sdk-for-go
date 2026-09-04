// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import "context"

// This is the build of the package that is not bound to the Cosmos driver. It is a diagnostic
// path, not an alternative implementation: v2 executes every operation through the driver, so
// there is no pure-Go mode, no degraded mode and no fallback.
//
// It exists so that a build without cgo fails with something that names the cause. Without it a
// caller would get `build constraints exclude all Go files`, which says nothing about what is
// wrong or how to fix it. Construction still succeeds, so the API can be compiled and explored
// against, and operations report why they cannot run.

// driverAvailable reports whether this build can reach the Cosmos driver.
const driverAvailable = false

// nativeDriver holds no resources in this build.
type nativeDriver struct{}

// cancel is a no-op because the diagnostic build owns no credential callback.
func (d *nativeDriver) cancel() {}

// initialize reports that the diagnostic build has no driver to initialize.
func (d *nativeDriver) initialize(context.Context) error {
	return newDriverUnavailableError()
}

// openDriver reports success without acquiring anything. Operations fail when they are called,
// with a message that says what the build is missing, rather than at construction.
func openDriver(driverConfig) (*nativeDriver, error) {
	return nil, nil
}

// close is a no-op. It tolerates a nil receiver because openDriver returns one.
func (d *nativeDriver) close() error {
	return nil
}

// execute reports that this build cannot reach the driver. The driver-backed build in
// driver_native.go runs the operation instead.
func (c *Client) execute(context.Context, itemRequest) (ItemResponse, []byte, error) {
	return ItemResponse{}, nil, newDriverUnavailableError()
}

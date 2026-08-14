// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !cgo || !azcosmos_driver

package azcosmos

// This is the build of the package that is not bound to the Cosmos driver. It is what `go build`
// selects by default, so the package stays usable for compiling against the API and for tooling
// that cannot link a native library: no C toolchain is needed and CGO_ENABLED=0 works.
//
// Operations report [CodeClientError] with a not-implemented message rather than reaching a
// driver. Build with `-tags azcosmos_driver` and CGO_ENABLED=1 to select the binding in
// driver_native.go instead.

// driverAvailable reports whether this build can reach the Cosmos driver.
const driverAvailable = false

// nativeDriver holds no resources in this build.
type nativeDriver struct{}

// openDriver reports success without acquiring anything, so that a client can still be constructed
// and its surface explored. Operations fail when they are called, not here.
func openDriver(driverConfig) (*nativeDriver, error) {
	return nil, nil
}

// close is a no-op. It tolerates a nil receiver because openDriver returns one.
func (d *nativeDriver) close() error {
	return nil
}

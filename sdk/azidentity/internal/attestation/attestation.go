// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package attestation

import "sync/atomic"

var registered atomic.Bool

// Register enables attestation integration.
// It must only be called during package initialization.
func Register() {
	registered.Store(true)
}

// Registered reports whether the optional integration has been imported.
func Registered() bool {
	return registered.Load()
}

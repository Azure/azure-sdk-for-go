// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import "runtime/cgo"

func cgoHandleValue(handle cgo.Handle) (value any, valid bool) {
	defer func() {
		if recover() != nil {
			value = nil
			valid = false
		}
	}()
	return handle.Value(), true
}

// cgo.Handle panics on duplicate deletion; native cleanup callbacks must remain idempotent.
func cgoHandleDelete(handle cgo.Handle) {
	defer func() {
		_ = recover()
	}()
	handle.Delete()
}

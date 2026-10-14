// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !cgo || !((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

type nativeRuntime struct{}

func openRuntime(RuntimeOptions) (*nativeRuntime, error) {
	return nil, newDriverUnavailableError()
}

func (r *nativeRuntime) close() {}

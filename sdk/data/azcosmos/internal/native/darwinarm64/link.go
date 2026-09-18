// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && darwin && !ios && arm64

// Package darwinarm64 provides the Go token callback shim for macOS on Apple silicon.
package darwinarm64

import "C"

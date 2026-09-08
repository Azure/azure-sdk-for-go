// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && darwin && !ios && arm64

// Package darwinarm64 links the native Cosmos driver for macOS on Apple silicon.
package darwinarm64

/*
#cgo LDFLAGS: -framework Security -framework CoreFoundation -liconv -lc -lm
*/
import "C"

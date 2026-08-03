// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build darwin

package main

import "golang.org/x/sys/unix"

// availableSystemMemoryBytes returns the total physical memory on macOS, in
// bytes. Darwin does not expose a simple "available" counter comparable to
// Linux's MemAvailable, so total physical RAM is used as a conservative budget
// base; the caller's safety fraction provides headroom. Returns 0 on error.
func availableSystemMemoryBytes() uint64 {
	memSize, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return 0
	}
	return memSize
}

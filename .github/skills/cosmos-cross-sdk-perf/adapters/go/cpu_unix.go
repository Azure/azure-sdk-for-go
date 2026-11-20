// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build unix

package main

// cSpell:words Stime Usec Utime

import "syscall"

func processCPUSeconds() float64 {
	var usage syscall.Rusage
	if syscall.Getrusage(syscall.RUSAGE_SELF, &usage) != nil {
		return 0
	}
	return float64(usage.Utime.Sec+usage.Stime.Sec) +
		float64(usage.Utime.Usec+usage.Stime.Usec)/1_000_000
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build !darwin && !windows

package main

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// availableSystemMemoryBytes returns the physical memory currently available
// to user processes, in bytes, for every platform except macOS and Windows
// (which have their own build-tagged implementations because they require
// platform-only syscalls that don't compile elsewhere).
//
// The OS is selected with a switch: Linux is read from /proc/meminfo's
// MemAvailable line; other Unix variants without a uniform query report 0
// ("unknown"), in which case callers skip the memory budget check rather than
// guess.
func availableSystemMemoryBytes() uint64 {
	switch runtime.GOOS {
	case "linux":
		if v, ok := memAvailableFromProc(); ok {
			return v
		}
	}
	return 0
}

// memAvailableFromProc parses the MemAvailable line from /proc/meminfo and
// returns its value in bytes. The bool is false when the value can't be read.
func memAvailableFromProc() (uint64, bool) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "MemAvailable:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return 0, false
		}
		kb, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0, false
		}
		return kb * 1024, true
	}
	return 0, false
}

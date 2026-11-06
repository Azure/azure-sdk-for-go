//go:build linux

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

func monotonicNS() (uint64, error) {
	var value syscall.Timespec
	_, _, errno := syscall.Syscall(syscall.SYS_CLOCK_GETTIME, 1, uintptr(unsafe.Pointer(&value)), 0)
	if errno != 0 {
		return 0, fmt.Errorf("clock_gettime(CLOCK_MONOTONIC): %w", errno)
	}
	return uint64(value.Sec)*1_000_000_000 + uint64(value.Nsec), nil
}

func measurementBarrier(name string) error {
	directory := os.Getenv("COSMOS_PERF_BARRIER_DIR")
	if directory == "" {
		return nil
	}
	file, err := os.OpenFile(filepath.Join(directory, "request-"+name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create measurement barrier: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close measurement barrier: %w", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(directory, "ready-"+name)); err == nil {
			return nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read measurement barrier: %w", err)
		}
		if time.Now().After(deadline) {
			return errors.New("measurement barrier timed out after 30 seconds")
		}
		time.Sleep(time.Millisecond)
	}
}

func measurementBounds(start, end uint64) error {
	directory := os.Getenv("COSMOS_PERF_BARRIER_DIR")
	if directory == "" {
		return nil
	}
	value := map[string]uint64{"monotonicStartNs": start, "monotonicEndNs": end}
	file, err := os.OpenFile(filepath.Join(directory, "measurement-bounds.json"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open measurement bounds: %w", err)
	}
	writeErr := json.NewEncoder(file).Encode(value)
	return errors.Join(writeErr, file.Close())
}

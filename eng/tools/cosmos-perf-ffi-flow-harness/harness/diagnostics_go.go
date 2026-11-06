// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"errors"
	"os"
	"runtime"
	"runtime/pprof"
	runtimeTrace "runtime/trace"
	"sync"
	"time"
)

var diagnosticState struct {
	mu         sync.Mutex
	cpu        *os.File
	trace      *os.File
	timer      *time.Timer
	traceError error
	stopped    bool
}

func startDiagnostics() error {
	if os.Getenv("COSMOS_PERF_BLOCK_PROFILE") != "" {
		runtime.SetBlockProfileRate(1)
	}
	if os.Getenv("COSMOS_PERF_MUTEX_PROFILE") != "" {
		runtime.SetMutexProfileFraction(1)
	}
	if path := os.Getenv("COSMOS_PERF_CPU_PROFILE"); path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		if err := pprof.StartCPUProfile(file); err != nil {
			return errors.Join(err, file.Close())
		}
		diagnosticState.cpu = file
	}
	if path := os.Getenv("COSMOS_PERF_RUNTIME_TRACE"); path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		if err := runtimeTrace.Start(file); err != nil {
			return errors.Join(err, file.Close())
		}
		diagnosticState.mu.Lock()
		diagnosticState.trace = file
		diagnosticState.mu.Unlock()
		diagnosticState.timer = time.AfterFunc(5*time.Second, stopRuntimeTrace)
	}
	return nil
}

func stopRuntimeTrace() {
	diagnosticState.mu.Lock()
	defer diagnosticState.mu.Unlock()
	if diagnosticState.trace != nil {
		runtimeTrace.Stop()
		diagnosticState.traceError = diagnosticState.trace.Close()
		diagnosticState.trace = nil
	}
}

func stopDiagnostics() error {
	if diagnosticState.stopped {
		return nil
	}
	diagnosticState.stopped = true
	if diagnosticState.timer != nil {
		diagnosticState.timer.Stop()
	}
	stopRuntimeTrace()
	var result error
	if diagnosticState.cpu != nil {
		pprof.StopCPUProfile()
		result = diagnosticState.cpu.Close()
		diagnosticState.cpu = nil
	}
	diagnosticState.mu.Lock()
	result = errors.Join(result, diagnosticState.traceError)
	diagnosticState.mu.Unlock()
	for _, profile := range []struct{ name, variable string }{
		{"heap", "COSMOS_PERF_HEAP_PROFILE"},
		{"block", "COSMOS_PERF_BLOCK_PROFILE"},
		{"mutex", "COSMOS_PERF_MUTEX_PROFILE"},
	} {
		if path := os.Getenv(profile.variable); path != "" {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
			if err != nil {
				result = errors.Join(result, err)
				continue
			}
			result = errors.Join(result, pprof.Lookup(profile.name).WriteTo(file, 0), file.Close())
		}
	}
	return result
}

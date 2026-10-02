// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// RuntimeOptions configures a runtime. Defaults and environment settings are captured at construction.
type RuntimeOptions struct {
	// Operation supplies defaults below client and request options.
	Operation OperationOptions
	// ApplicationID is the application identifier shared by attached clients.
	ApplicationID string
	// CPURefreshInterval controls CPU/memory sampling. Nil inherits the native default;
	// explicit intervals must be between one and sixty seconds.
	CPURefreshInterval *time.Duration
}

// Runtime owns shared execution resources with construction-time defaults.
// Close releases the owner's reference; attached clients keep the resources alive.
// Same-account clients share native caches, including cached credentials. Use separate
// runtimes when accounts must be accessed under isolated identities.
type Runtime struct {
	mu      sync.Mutex
	native  *nativeRuntime
	clients int
	closed  bool
}

var globalRuntime struct {
	sync.Mutex
	runtime *Runtime
}

func defaultRuntime() (*Runtime, error) {
	globalRuntime.Lock()
	defer globalRuntime.Unlock()
	if globalRuntime.runtime == nil {
		runtime, err := NewRuntime(nil)
		if err != nil {
			return nil, err
		}
		globalRuntime.runtime = runtime
	}
	return globalRuntime.runtime, nil
}

// NewRuntime creates an isolated runtime without network I/O. Nil selects driver defaults.
// Builds without native support return a driver-unavailable error.
func NewRuntime(options *RuntimeOptions) (*Runtime, error) {
	var config RuntimeOptions
	if options != nil {
		config = *options
		config.Operation = options.Operation.clone()
		config.CPURefreshInterval = clonePointer(options.CPURefreshInterval)
	}
	if err := config.Operation.validate(); err != nil {
		return nil, err
	}
	if strings.ContainsRune(config.ApplicationID, 0) {
		return nil, errors.New("azcosmos: RuntimeOptions.ApplicationID must not contain a NUL byte")
	}
	if value := config.CPURefreshInterval; value != nil && (*value < time.Second || *value > time.Minute) {
		return nil, errors.New("azcosmos: CPURefreshInterval must be between one and sixty seconds")
	}
	native, err := openRuntime(config)
	if err != nil {
		return nil, err
	}
	return &Runtime{native: native}, nil
}

// Close releases the owner's runtime reference and prevents new client attachments.
// Existing clients remain usable until their own Close calls drain and release them.
// Close is idempotent and concurrency-safe.
func (r *Runtime) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	r.releaseIfUnused()
	return nil
}

func (r *Runtime) detach() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients--
	r.releaseIfUnused()
}

func (r *Runtime) releaseIfUnused() {
	if r.closed && r.clients == 0 {
		r.native.close()
		r.native = nil
	}
}

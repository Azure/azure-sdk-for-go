// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"strings"
	"sync"
)

// RuntimeOptions configures a shared runtime. Environment settings are captured at construction.
type RuntimeOptions struct {
	// Operation supplies defaults below client and request options.
	Operation OperationOptions
	// ApplicationID is the immutable application identifier shared by attached clients.
	ApplicationID string
}

// Runtime owns shared native execution resources and operation defaults.
// Construct it with NewRuntime; the zero value is not usable. Runtime is concurrency-safe.
// Native v0.2.0 retains credential-bearing container references, so each account hostname
// may be attached only once for the lifetime of a Runtime, even after its client closes.
type Runtime struct {
	mu            sync.Mutex
	native        *nativeRuntime
	applicationID string
	clients       map[*Client]struct{}
	// Reservations survive client closure because native account caches survive it too.
	accountHosts map[string]struct{}
	closing      bool
	inflight     sync.WaitGroup
	closeOnce    sync.Once
	closeErr     error
}

// NewRuntime creates a runtime without network I/O. Nil options selects driver defaults.
// Builds without native support return a driver-unavailable error.
func NewRuntime(options *RuntimeOptions) (*Runtime, error) {
	var config RuntimeOptions
	if options != nil {
		config = *options
		config.Operation = options.Operation.clone()
	}
	if err := config.Operation.validate(); err != nil {
		return nil, err
	}
	if strings.ContainsRune(config.ApplicationID, 0) {
		return nil, errors.New("azcosmos: RuntimeOptions.ApplicationID must not contain a NUL byte")
	}
	native, err := openRuntime(config)
	if err != nil {
		return nil, err
	}
	return &Runtime{
		native: native, applicationID: config.ApplicationID,
		clients: make(map[*Client]struct{}), accountHosts: make(map[string]struct{}),
	}, nil
}

// SetOperationOptions atomically replaces all runtime defaults; it does not patch individual fields.
// Admitted operations retain their snapshot. Caller-owned pointers, maps, and slices are copied.
func (r *Runtime) SetOperationOptions(options OperationOptions) error {
	if err := options.validate(); err != nil {
		return err
	}
	options = options.clone()
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closing || r.native == nil {
		return &Error{Code: CodeClientClosed, Message: "the runtime is closed or uninitialized"}
	}
	return r.native.setOperationOptions(options)
}

// Close rejects new work across all attached clients, drains admitted operations, and releases
// the clients and runtime. Concurrent callers wait for the same result. Close is idempotent.
func (r *Runtime) Close() error {
	r.closeOnce.Do(func() {
		r.mu.Lock()
		r.closing = true
		clients := make([]*Client, 0, len(r.clients))
		for client := range r.clients {
			clients = append(clients, client)
		}
		r.mu.Unlock()
		// Closing clients cancels token acquisition before waiting for their admitted operations.
		for _, client := range clients {
			r.closeErr = errors.Join(r.closeErr, client.close())
		}
		r.inflight.Wait()
		r.mu.Lock()
		defer r.mu.Unlock()
		r.native.close()
		r.native = nil
	})
	return r.closeErr
}

func (r *Runtime) acquire() (func(), error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closing || r.native == nil {
		return nil, &Error{Code: CodeClientClosed, Message: "the runtime is closed or uninitialized"}
	}
	r.inflight.Add(1)
	return r.inflight.Done, nil
}

func (r *Runtime) detach(client *Client) {
	r.mu.Lock()
	delete(r.clients, client)
	r.mu.Unlock()
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "time"

// Diagnostics is an immutable native operation snapshot copied before its completion is freed.
// The zero value has no native snapshot, for example when the caller stops waiting before completion.
type Diagnostics struct {
	json          string
	elapsed       time.Duration
	requestCharge float64
	requestCount  uint32
	compacted     bool
	failure       bool
	err           error
}

// JSON returns the native default-verbosity JSON snapshot. Its schema is driver-defined.
func (d Diagnostics) JSON() string { return d.json }

// Elapsed returns the native operation duration, including retries.
func (d Diagnostics) Elapsed() time.Duration { return d.elapsed }

// RequestCharge returns the aggregate native request charge.
func (d Diagnostics) RequestCharge() float64 { return d.requestCharge }

// RequestCount returns the total number of native attempts, including retries.
func (d Diagnostics) RequestCount() uint32 { return d.requestCount }

// IsCompacted reports whether the native driver discarded older attempt records.
func (d Diagnostics) IsCompacted() bool { return d.compacted }

// IsFailure reports whether the native snapshot records a failure.
func (d Diagnostics) IsFailure() bool { return d.failure }

// Err returns a diagnostic rendering error without hiding the operation's own result.
func (d Diagnostics) Err() error { return d.err }

// Available reports whether this value contains a native diagnostic snapshot or rendering error.
func (d Diagnostics) Available() bool { return d.json != "" || d.err != nil }

// OperationDiagnostic describes a completed Go item call. Native fields can be absent on validation
// failure or when context cancellation ends the Go wait before native completion.
type OperationDiagnostic struct {
	// Operation is the item operation name.
	Operation string
	// DatabaseID is the database addressed by the call.
	DatabaseID string
	// ContainerID is the container addressed by the call.
	ContainerID string
	// Diagnostics contains the native snapshot when available.
	Diagnostics Diagnostics
	// Error is the error returned by the Go call.
	Error error
}

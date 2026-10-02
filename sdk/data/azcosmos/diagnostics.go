// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "time"

// Diagnostics is a Go-owned snapshot of the native driver's operation diagnostics.
// It is available only when the driver attaches diagnostics to a completion. In
// particular, cancellation can return before a native completion is available.
type Diagnostics struct {
	// StatusCode and SubStatus describe the final completion. They are zero
	// when no HTTP status or sub-status is available.
	StatusCode int
	SubStatus  int

	// AttemptCount includes the initial request and every native retry. It can
	// exceed len(Attempts) if the driver compacted the timeline.
	AttemptCount uint32

	// TotalRequestCharge is the charge across all attempts, including retries.
	TotalRequestCharge float64

	// Elapsed is the native driver's total wall-clock time for the operation.
	Elapsed time.Duration

	// Completed reports whether the driver recorded a terminal outcome.
	Completed bool

	// Failed reports whether that terminal outcome was a failure.
	Failed bool

	// Compacted reports whether the driver discarded older attempt records.
	Compacted bool

	// RegionsContacted lists distinct regions in first-contact order.
	RegionsContacted []string

	// Attempts contains the retained per-attempt records in execution order.
	Attempts []DiagnosticAttempt

	// JSON is the driver's detailed diagnostics rendering, copied into Go memory.
	// It is diagnostic data, not a stable schema; do not parse it for SDK behavior.
	JSON string
}

// DiagnosticAttempt describes one native request retained in the timeline.
type DiagnosticAttempt struct {
	Endpoint string
	Region   string

	// StatusCode and SubStatus are the status of this attempt. SubStatus is zero
	// when the driver did not record one.
	StatusCode int
	SubStatus  int

	Latency       time.Duration
	RequestCharge float64

	// ServerDurationMS is negative when the service did not report a duration.
	ServerDurationMS float64
}

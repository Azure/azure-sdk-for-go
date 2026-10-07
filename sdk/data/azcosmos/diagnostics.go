// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// EventDiagnostics entries report when an operation's native diagnostics snapshot could not be
// copied. Diagnostics are informational only: the operation's own outcome never depends on this
// copy succeeding, so a failure here is logged rather than reported as an operation failure.
//
// Declared here rather than in completion_native.go, which is platform-tagged, so the constant
// stays part of the package's API on every platform: a caller using azcosmos.EventDiagnostics
// otherwise fails to compile on Windows, unsupported architectures, or CGO_ENABLED=0, where that
// file is excluded from the build. [EventRouting] in routing_strategy.go follows the same pattern.
const EventDiagnostics log.Event = "CosmosDiagnostics"

// DiagnosticsVerbosity controls how much detail the native driver renders into
// [Diagnostics.JSON]. Set it on [ClientOptions.DiagnosticsVerbosity].
//
// It does not affect [Diagnostics.Attempts], [Diagnostics.RegionsContacted] or any scalar field:
// those are always copied losslessly regardless of verbosity. It only governs the native driver's
// own JSON rendering, which is diagnostic text rather than a stable schema.
//
// The zero value, DiagnosticsVerbosityDefault, resolves to DiagnosticsVerbositySummary.
type DiagnosticsVerbosity int32

const (
	// DiagnosticsVerbosityDefault leaves the rendering choice to the SDK. It currently resolves
	// to DiagnosticsVerbositySummary; that resolution may change in a future version without
	// being considered a breaking change.
	DiagnosticsVerbosityDefault DiagnosticsVerbosity = iota

	// DiagnosticsVerbositySummary renders the native driver's compact, deduplicated diagnostics:
	// requests are grouped by region, with the first and last request per region kept in full
	// detail and the rest summarized by count and duration statistics. This is cheap enough that
	// every completion can afford to pay for it.
	DiagnosticsVerbositySummary

	// DiagnosticsVerbosityDetailed renders every individual request the native driver recorded,
	// with no deduplication or truncation. This duplicates data already available losslessly and
	// more cheaply through [Diagnostics.Attempts]: prefer that field unless something specifically
	// needs the driver's own JSON rendering of the same data. Every completion pays the
	// serialization and allocation cost of this rendering whether or not it is ever read, since
	// the native diagnostics handle cannot be retained past the completion to defer that work.
	DiagnosticsVerbosityDetailed
)

// String returns the name of the verbosity level, for logging and error messages.
func (v DiagnosticsVerbosity) String() string {
	switch v {
	case DiagnosticsVerbosityDefault:
		return "Default"
	case DiagnosticsVerbositySummary:
		return "Summary"
	case DiagnosticsVerbosityDetailed:
		return "Detailed"
	default:
		return fmt.Sprintf("DiagnosticsVerbosity(%d)", int32(v))
	}
}

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

	// JSON is the native driver's diagnostics rendering, copied into Go memory at the verbosity
	// configured by [ClientOptions.DiagnosticsVerbosity] (DiagnosticsVerbositySummary by
	// default). Use [Diagnostics.Attempts] for the full per-attempt detail instead of parsing
	// JSON, which is diagnostic data, not a stable schema.
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

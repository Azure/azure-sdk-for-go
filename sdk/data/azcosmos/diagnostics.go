// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"fmt"
	"sync"
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

// DiagnosticsVerbosity controls how much detail [Diagnostics.JSON] renders.
// Set it on [ClientOptions.DiagnosticsVerbosity].
//
// It does not affect [Diagnostics.Attempts], [Diagnostics.RegionsContacted] or any scalar field:
// those are always copied from the native completion regardless of verbosity. It only governs the
// lazy Go JSON rendering, which is diagnostic text rather than a stable schema.
//
// The zero value, DiagnosticsVerbosityDefault, resolves to DiagnosticsVerbositySummary.
type DiagnosticsVerbosity int32

const (
	// DiagnosticsVerbosityDefault leaves the rendering choice to the SDK. It currently resolves
	// to DiagnosticsVerbositySummary; that resolution may change in a future version without
	// being considered a breaking change.
	DiagnosticsVerbosityDefault DiagnosticsVerbosity = iota

	// DiagnosticsVerbositySummary renders totals, final status, completion flags, and contacted
	// regions. It omits retained per-attempt records.
	DiagnosticsVerbositySummary

	// DiagnosticsVerbosityDetailed renders the summary fields plus every retained record in
	// [Diagnostics.Attempts]. When [Diagnostics.Compacted] is true, these attempts are only the
	// native driver's retained subset, not the full timeline.
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

	json *diagnosticsJSON
}

// JSON renders and caches this Go-owned diagnostics snapshot. The rendering uses the
// [DiagnosticsVerbosity] configured when the client was created and is safe for concurrent calls.
//
// The current native ABI exposes totals, elapsed time, request counts, completion flags, contacted
// regions, and retained attempts. JSON contains only those copied fields; it is not the native
// driver's full diagnostics schema and callers should not treat its shape as a stable contract.
// Use the typed fields on Diagnostics when programmatic access is required.
func (d *Diagnostics) JSON() (string, error) {
	if d == nil {
		return "", nil
	}
	cache := d.json
	if cache == nil {
		cache = newDiagnosticsJSON(d, DiagnosticsVerbosityDefault)
	}
	cache.once.Do(func() {
		data, err := marshalDiagnosticsJSON(cache.renderValue())
		if err != nil {
			cache.err = fmt.Errorf("azcosmos: rendering diagnostics JSON: %w", err)
			return
		}
		cache.value = string(data)
	})
	return cache.value, cache.err
}

type diagnosticsJSON struct {
	once      sync.Once
	verbosity DiagnosticsVerbosity
	snapshot  diagnosticsJSONSnapshot
	value     string
	err       error
}

type diagnosticsJSONSnapshot struct {
	StatusCode         int
	SubStatus          int
	AttemptCount       uint32
	TotalRequestCharge float64
	Elapsed            time.Duration
	Completed          bool
	Failed             bool
	Compacted          bool
	RegionsContacted   []string
	Attempts           []DiagnosticAttempt
}

type diagnosticsJSONSummary struct {
	StatusCode           int      `json:"statusCode"`
	SubStatus            int      `json:"subStatus"`
	AttemptCount         uint32   `json:"attemptCount"`
	RetainedAttemptCount int      `json:"retainedAttemptCount"`
	TotalRequestCharge   float64  `json:"totalRequestCharge"`
	ElapsedTimeMicros    int64    `json:"elapsedTimeMicros"`
	Completed            bool     `json:"completed"`
	Failed               bool     `json:"failed"`
	Compacted            bool     `json:"compacted"`
	RegionsContacted     []string `json:"regionsContacted"`
}

type diagnosticsJSONDetailed struct {
	diagnosticsJSONSummary
	Attempts []diagnosticAttemptJSON `json:"attempts"`
}

type diagnosticAttemptJSON struct {
	Endpoint         string  `json:"endpoint"`
	Region           string  `json:"region"`
	StatusCode       int     `json:"statusCode"`
	SubStatus        int     `json:"subStatus"`
	LatencyMS        int64   `json:"latencyMS"`
	RequestCharge    float64 `json:"requestCharge"`
	ServerDurationMS float64 `json:"serverDurationMS"`
}

var marshalDiagnosticsJSON = json.Marshal

func (d *Diagnostics) initializeJSON(verbosity DiagnosticsVerbosity) {
	d.json = newDiagnosticsJSON(d, verbosity)
}

func (d *Diagnostics) jsonVerbosity() DiagnosticsVerbosity {
	if d.json == nil {
		return DiagnosticsVerbosityDefault
	}
	return d.json.verbosity
}

func newDiagnosticsJSON(d *Diagnostics, verbosity DiagnosticsVerbosity) *diagnosticsJSON {
	return &diagnosticsJSON{
		verbosity: verbosity,
		snapshot: diagnosticsJSONSnapshot{
			StatusCode:         d.StatusCode,
			SubStatus:          d.SubStatus,
			AttemptCount:       d.AttemptCount,
			TotalRequestCharge: d.TotalRequestCharge,
			Elapsed:            d.Elapsed,
			Completed:          d.Completed,
			Failed:             d.Failed,
			Compacted:          d.Compacted,
			RegionsContacted:   append([]string(nil), d.RegionsContacted...),
			Attempts:           append([]DiagnosticAttempt(nil), d.Attempts...),
		},
	}
}

func (d *diagnosticsJSON) renderValue() any {
	summary := diagnosticsJSONSummary{
		StatusCode:           d.snapshot.StatusCode,
		SubStatus:            d.snapshot.SubStatus,
		AttemptCount:         d.snapshot.AttemptCount,
		RetainedAttemptCount: len(d.snapshot.Attempts),
		TotalRequestCharge:   d.snapshot.TotalRequestCharge,
		ElapsedTimeMicros:    d.snapshot.Elapsed.Microseconds(),
		Completed:            d.snapshot.Completed,
		Failed:               d.snapshot.Failed,
		Compacted:            d.snapshot.Compacted,
		RegionsContacted:     d.snapshot.RegionsContacted,
	}
	if d.verbosity != DiagnosticsVerbosityDetailed {
		return summary
	}
	attempts := make([]diagnosticAttemptJSON, len(d.snapshot.Attempts))
	for i, attempt := range d.snapshot.Attempts {
		attempts[i] = diagnosticAttemptJSON{
			Endpoint:         attempt.Endpoint,
			Region:           attempt.Region,
			StatusCode:       attempt.StatusCode,
			SubStatus:        attempt.SubStatus,
			LatencyMS:        attempt.Latency.Milliseconds(),
			RequestCharge:    attempt.RequestCharge,
			ServerDurationMS: attempt.ServerDurationMS,
		}
	}
	return diagnosticsJSONDetailed{diagnosticsJSONSummary: summary, Attempts: attempts}
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

// OperationDiagnostic describes an item call that passed validation and acquired its client.
// Diagnostics can be nil when context cancellation returns before native completion.
type OperationDiagnostic struct {
	// Operation is the item operation name.
	Operation string
	// DatabaseID is the database addressed by the call.
	DatabaseID string
	// ContainerID is the container addressed by the call.
	ContainerID string
	// Diagnostics is a Go-owned native snapshot, when available.
	Diagnostics *Diagnostics
	// Error is the error returned by the Go call.
	Error error
}

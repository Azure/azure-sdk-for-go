// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/stretchr/testify/require"
)

// TestEventDiagnosticsIsBuildNeutral guards EventDiagnostics staying declared in this
// build-neutral file rather than completion_native.go, which is platform-tagged: a caller using
// azcosmos.EventDiagnostics must compile and behave identically on every platform/cgo
// configuration, and this test file itself carries no build tag.
func TestEventDiagnosticsIsBuildNeutral(t *testing.T) {
	require.Equal(t, azlog.Event("CosmosDiagnostics"), EventDiagnostics)

	var event azlog.Event
	var message string
	azlog.SetEvents(EventDiagnostics)
	azlog.SetListener(func(receivedEvent azlog.Event, receivedMessage string) {
		event = receivedEvent
		message = receivedMessage
	})
	t.Cleanup(func() {
		azlog.SetListener(nil)
		azlog.SetEvents()
	})

	log.Write(EventDiagnostics, "copying operation diagnostics failed")
	require.Equal(t, EventDiagnostics, event)
	require.Equal(t, "copying operation diagnostics failed", message)
}

func TestDiagnosticsVerbosityString(t *testing.T) {
	for _, tt := range []struct {
		verbosity DiagnosticsVerbosity
		want      string
	}{
		{DiagnosticsVerbosityDefault, "Default"},
		{DiagnosticsVerbositySummary, "Summary"},
		{DiagnosticsVerbosityDetailed, "Detailed"},
		{DiagnosticsVerbosity(99), "DiagnosticsVerbosity(99)"},
		{DiagnosticsVerbosity(-1), "DiagnosticsVerbosity(-1)"},
	} {
		t.Run(tt.want, func(t *testing.T) {
			require.Equal(t, tt.want, tt.verbosity.String())
		})
	}
}

func TestDiagnosticsJSONIsLazyAndCached(t *testing.T) {
	original := marshalDiagnosticsJSON
	var calls atomic.Int32
	marshalDiagnosticsJSON = func(value any) ([]byte, error) {
		calls.Add(1)
		return original(value)
	}
	t.Cleanup(func() { marshalDiagnosticsJSON = original })

	diagnostics := sampleDiagnostics()
	diagnostics.initializeJSON(DiagnosticsVerbosityDetailed)
	require.Zero(t, calls.Load(), "capturing a completion must not render diagnostics JSON")

	first, err := diagnostics.JSON()
	require.NoError(t, err)
	second, err := diagnostics.JSON()
	require.NoError(t, err)
	require.Equal(t, first, second)
	require.Equal(t, int32(1), calls.Load())
	require.JSONEq(t, `{
		"statusCode":429,
		"subStatus":3200,
		"attemptCount":2,
		"retainedAttemptCount":1,
		"totalRequestCharge":3.5,
		"elapsedTimeMicros":1250,
		"completed":true,
		"failed":true,
		"compacted":true,
		"regionsContacted":["East US"],
		"attempts":[{
			"endpoint":"https://example.com/",
			"region":"East US",
			"statusCode":429,
			"subStatus":3200,
			"latencyMS":7,
			"requestCharge":1.25,
			"serverDurationMS":2.5
		}]
	}`, first)
}

func TestDiagnosticsJSONSummaryOmitsAttempts(t *testing.T) {
	diagnostics := sampleDiagnostics()
	diagnostics.initializeJSON(DiagnosticsVerbositySummary)

	rendered, err := diagnostics.JSON()
	require.NoError(t, err)
	var value map[string]any
	require.NoError(t, json.Unmarshal([]byte(rendered), &value))
	require.NotContains(t, value, "attempts")
	require.Equal(t, float64(1), value["retainedAttemptCount"])
}

func TestDiagnosticsJSONConcurrentCallsShareOneRendering(t *testing.T) {
	original := marshalDiagnosticsJSON
	var calls atomic.Int32
	marshalDiagnosticsJSON = func(value any) ([]byte, error) {
		calls.Add(1)
		return original(value)
	}
	t.Cleanup(func() { marshalDiagnosticsJSON = original })

	diagnostics := sampleDiagnostics()
	diagnostics.initializeJSON(DiagnosticsVerbosityDetailed)
	const callers = 32
	values := make([]string, callers)
	errs := make([]error, callers)
	var wg sync.WaitGroup
	for i := range callers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			values[i], errs[i] = diagnostics.JSON()
		}(i)
	}
	wg.Wait()
	for i := range callers {
		require.NoError(t, errs[i])
		require.Equal(t, values[0], values[i])
	}
	require.Equal(t, int32(1), calls.Load())
}

func TestDiagnosticsJSONCachesRenderingError(t *testing.T) {
	original := marshalDiagnosticsJSON
	var calls atomic.Int32
	marshalDiagnosticsJSON = func(any) ([]byte, error) {
		calls.Add(1)
		return nil, errors.New("boom")
	}
	t.Cleanup(func() { marshalDiagnosticsJSON = original })

	diagnostics := sampleDiagnostics()
	diagnostics.initializeJSON(DiagnosticsVerbositySummary)
	for range 2 {
		rendered, err := diagnostics.JSON()
		require.Empty(t, rendered)
		require.EqualError(t, err, "azcosmos: rendering diagnostics JSON: boom")
	}
	require.Equal(t, int32(1), calls.Load())
}

func sampleDiagnostics() *Diagnostics {
	return &Diagnostics{
		StatusCode:         429,
		SubStatus:          3200,
		AttemptCount:       2,
		TotalRequestCharge: 3.5,
		Elapsed:            1250 * time.Microsecond,
		Completed:          true,
		Failed:             true,
		Compacted:          true,
		RegionsContacted:   []string{"East US"},
		Attempts: []DiagnosticAttempt{{
			Endpoint:         "https://example.com/",
			Region:           "East US",
			StatusCode:       429,
			SubStatus:        3200,
			Latency:          7 * time.Millisecond,
			RequestCharge:    1.25,
			ServerDurationMS: 2.5,
		}},
	}
}

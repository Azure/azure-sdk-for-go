// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNativeDiagnosticsVerbosityMapsEveryEnumValue walks every DiagnosticsVerbosity constant plus
// an unrecognized value, verifying each materializes into the native driver's
// cosmos_diagnostics_verbosity_t the package's doc comments promise: DiagnosticsVerbosityDetailed
// is the only value that renders at a different native verbosity than
// DiagnosticsVerbosityDefault, which resolves to the same native SUMMARY code as
// DiagnosticsVerbositySummary and as any unrecognized value.
func TestNativeDiagnosticsVerbosityMapsEveryEnumValue(t *testing.T) {
	summaryCode := nativeDiagnosticsVerbosityCode(DiagnosticsVerbositySummary)
	detailedCode := nativeDiagnosticsVerbosityCode(DiagnosticsVerbosityDetailed)
	require.NotEqual(t, summaryCode, detailedCode, "SUMMARY and DETAILED must be distinct native codes")

	for _, tt := range []struct {
		name      string
		verbosity DiagnosticsVerbosity
		want      uint32
	}{
		{"Default resolves to the same native code as Summary", DiagnosticsVerbosityDefault, summaryCode},
		{"Summary", DiagnosticsVerbositySummary, summaryCode},
		{"Detailed", DiagnosticsVerbosityDetailed, detailedCode},
		{"unrecognized value falls back to Summary's native code", DiagnosticsVerbosity(99), summaryCode},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, nativeDiagnosticsVerbosityCode(tt.verbosity))
		})
	}
}

// TestDiagnosticsVerbosityString covers every named constant plus an out-of-range value, since the
// error message ClientOptions.validate produces for an invalid value embeds this string.
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

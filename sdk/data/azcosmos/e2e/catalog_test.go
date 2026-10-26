// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"os/exec"
	"strings"
	"testing"

	internale2e "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2/internal/e2e"
	"github.com/stretchr/testify/require"
)

// TestImplementationMapMatchesCompiledTests is the validation evidence internal/e2e.Validate
// needs: the set of Go test function names actually discovered by `go test -list` for this
// build, not merely text in this package that resembles a function declaration. It runs
// unconditionally (no EMULATOR requirement) so catalog/implementation-map drift fails CI even
// when the emulator stage is unavailable.
func TestImplementationMapMatchesCompiledTests(t *testing.T) {
	compiled := compiledTestNames(t)

	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)

	err = internale2e.Validate(catalog, internale2e.Implementations, compiled)
	require.NoError(t, err)
}

// TestCoverageReportListsEveryScenario documents the exact covered/blocked/not-applicable matrix
// in test output, satisfying "the task documents covered, blocked, and intentionally
// non-applicable scenario IDs" without requiring a human to read implementations.go directly.
func TestCoverageReportListsEveryScenario(t *testing.T) {
	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)
	require.NoError(t, internale2e.Validate(catalog, internale2e.Implementations, compiledTestNames(t)))

	report := internale2e.Coverage(catalog, internale2e.Implementations)
	require.Equal(t, len(catalog.Scenarios), len(report.Active)+len(report.Blocked)+len(report.NotApplicable))
	t.Log("\n" + report.String())
}

// compiledTestNames runs `go test -list` against this package so the implementation map is
// checked against tests the current build actually compiled, matching the current package's
// build constraints (platform, cgo, build tags) rather than a static source scan.
func compiledTestNames(t *testing.T) map[string]bool {
	t.Helper()

	output, err := exec.Command("go", "test", "-list", "^Test", ".").CombinedOutput() //nolint:gosec // fixed args, no injection
	require.NoErrorf(t, err, "go test -list failed: %s", output)

	names := make(map[string]bool)
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || line == "ok" || strings.HasPrefix(line, "ok  \t") || !strings.HasPrefix(line, "Test") {
			continue
		}
		names[line] = true
	}
	require.NotEmpty(t, names, "go test -list reported no compiled tests; output was: %s", output)
	return names
}

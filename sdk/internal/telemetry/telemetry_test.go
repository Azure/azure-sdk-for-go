// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package telemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	userAgent := Format("azservicebus", "v1.0.0")

	// Examples:
	// * azsdk-go-azservicebus/v1.0.0 (go1.19.3; linux)
	// * azsdk-go-azservicebus/v1.0.0 (go1.19; Windows_NT)
	// * azsdk-go-azservicebus/v1.0.0 (go1.21rc3; linux)
	// * azsdk-go-azservicebus/v1.0.0 (go1.22.0 X:nocoverageredesign; linux)
	//
	// The OS varies based on the actual platform but it's a small set.
	re := `^azsdk-go-azservicebus/v1.0.0` +
		` ` +
		`\(` +
		`go\d+\.\d+(?:|\.\d+|rc\d+)(?:\s[a-zA-Z0-9|:]+)?; (?:Windows_NT|linux|freebsd)` +
		`\)$`

	require.Regexp(t, re, userAgent)
	require.Regexp(t, re, "azsdk-go-azservicebus/v1.0.0 (go1.21rc3; linux)")
	require.Regexp(t, re, "azsdk-go-azservicebus/v1.0.0 (go1.22.0 X:nocoverageredesign; linux)")
}

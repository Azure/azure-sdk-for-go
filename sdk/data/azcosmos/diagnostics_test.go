// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

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

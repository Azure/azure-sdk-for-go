// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// telemetryClient embeds an appinsights TelemtryClient but adds some convenience
// functions to make sending events with baggage/properties one-liners.
type telemetryClient struct {
	appinsights.TelemetryClient
}

func (tc *telemetryClient) TrackMetric(name string, value float64, properties map[string]string) {
	mt := appinsights.NewMetricTelemetry(name, value)

	for k, v := range properties {
		mt.Properties[k] = v
	}

	tc.TelemetryClient.Track(mt)
}

func (tc *telemetryClient) TrackEvent(name string, properties map[string]string) {
	et := appinsights.NewEventTelemetry(name)

	for k, v := range properties {
		et.Properties[k] = v
	}

	tc.TelemetryClient.Track(et)
}

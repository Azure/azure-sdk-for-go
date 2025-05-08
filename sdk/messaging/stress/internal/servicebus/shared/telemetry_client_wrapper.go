// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

func NewTelemetryClientWrapper[MetricT ~string, EventT ~string]() *TelemetryClientWrapper[MetricT, EventT] {
	return &TelemetryClientWrapper[MetricT, EventT]{}
}

// TelemetryClientWrapper is a wrapper for telemetry client, once we get that phased back in.
type TelemetryClientWrapper[MetricT ~string, EventT ~string] struct {
	context TelemetryClientWrapperContext
}

type TelemetryClientWrapperContext struct {
	CommonProperties map[string]string
}

func (tc *TelemetryClientWrapper[MetricT, EventT]) TrackException(err error) {
	// will be replaced with a real telemetry client once we enable OTEL reporting.
}

func (tc *TelemetryClientWrapper[MetricT, EventT]) TrackEvent(name EventT) {
	// will be replaced with a real telemetry client once we enable OTEL reporting.
}

func (tc *TelemetryClientWrapper[MetricT, EventT]) TrackMetricWithProps(name MetricT, value float64, properties map[string]string) {
	// will be replaced with a real telemetry client once we enable OTEL reporting.
}

func (tc *TelemetryClientWrapper[MetricT, EventT]) TrackEventWithProps(name EventT, properties map[string]string) {
	// will be replaced with a real telemetry client once we enable OTEL reporting.
}

func (tc *TelemetryClientWrapper[MetricT, EventT]) TrackExceptionWithProps(err error, properties map[string]string) {
	// will be replaced with a real telemetry client once we enable OTEL reporting.
}

func (tc *TelemetryClientWrapper[MetricT, EventT]) Flush() {
	// tc.TC.Channel().Flush()
	// <-tc.TC.Channel().Close()
}

// Context returns the context that is included for each reported event or metric.
func (tc *TelemetryClientWrapper[MetricT, EventT]) Context() *TelemetryClientWrapperContext {
	return &tc.context
}

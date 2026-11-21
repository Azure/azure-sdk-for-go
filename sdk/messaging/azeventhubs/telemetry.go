// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
)

// TelemetryOptions configures optional telemetry for a client and its partition clients.
// Each signal is disabled until its provider is supplied. The application owns provider shutdown.
type TelemetryOptions struct {
	// TracingProvider creates operation and transport spans. Use azotel.NewTracingProvider to provide an OpenTelemetry tracer.
	TracingProvider tracing.Provider
	// MeterProvider records operation counts, event counts, durations, and AMQP lifecycle transitions.
	MeterProvider metric.MeterProvider
	// LoggerProvider receives structured AMQP lifecycle logs.
	LoggerProvider otellog.LoggerProvider
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal"
	"go.opentelemetry.io/otel"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
)

// TelemetryOptions configures optional telemetry for a client and its partition clients.
// The zero value disables all signals. The application owns provider configuration and shutdown.
type TelemetryOptions struct {
	// TracingProvider creates operation and transport spans. Use azotel.NewTracingProvider to provide an OpenTelemetry tracer.
	TracingProvider tracing.Provider
	// EnableMetrics records counts and durations with the global OpenTelemetry meter provider.
	EnableMetrics bool
	// EnableLogs emits AMQP lifecycle logs through the global OpenTelemetry logger provider.
	EnableLogs bool
}

func (o TelemetryOptions) diagnostics(host, hub string) *internal.Diagnostics {
	var meter metric.MeterProvider
	if o.EnableMetrics {
		meter = otel.GetMeterProvider()
	}
	var logger otellog.LoggerProvider
	if o.EnableLogs {
		logger = otel.GetLoggerProvider()
	}
	return internal.NewDiagnostics(host, hub, o.TracingProvider, meter, logger)
}

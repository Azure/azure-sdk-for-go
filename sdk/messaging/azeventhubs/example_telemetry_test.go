// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
)

func ExampleTelemetryOptions() {
	// Use providers that the application has configured with its exporters.
	telemetry := azeventhubs.TelemetryOptions{
		TracingProvider: azotel.NewTracingProvider(otel.GetTracerProvider(), nil),
		LoggerProvider:  global.GetLoggerProvider(),
		MeterProvider:   otel.GetMeterProvider(),
	}
	producer, err := azeventhubs.NewProducerClientFromConnectionString(
		os.Getenv("EVENTHUB_CONNECTION_STRING"), "my-hub", &azeventhubs.ProducerClientOptions{Telemetry: telemetry})
	if err != nil {
		panic(err)
	}
	defer producer.Close(context.Background())
}

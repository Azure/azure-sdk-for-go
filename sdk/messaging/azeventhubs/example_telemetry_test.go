// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel"
	"go.opentelemetry.io/otel"
)

func ExampleTelemetryOptions() {
	// Use providers that the application has configured with its exporters.
	telemetry := azeventhubs.TelemetryOptions{
		TracingProvider: azotel.NewTracingProvider(otel.GetTracerProvider(), nil),
		EnableLogs:      true,
		EnableMetrics:   true,
	}
	producer, err := azeventhubs.NewProducerClientFromConnectionString(
		os.Getenv("EVENTHUB_CONNECTION_STRING"), "my-hub", &azeventhubs.ProducerClientOptions{Telemetry: telemetry})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	if err := producer.Close(context.TODO()); err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
}

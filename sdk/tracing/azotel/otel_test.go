//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azotel_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel/internal"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func ExampleNewTracingProvider() {
	// end-to-end example creating a simple console-based trace exporter
	// then using the azotel adapter to conenct it to a generic azcore.Client

	// create a basic otel TracerProvider
	otelTP := tracesdk.NewTracerProvider(tracesdk.WithBatcher(ConsoleExporter{}))

	// connect the otel TracerProvider to the Azure SDK client
	options := azcore.ClientOptions{}
	options.TracingProvider = azotel.NewTracingProvider(otelTP, nil)

	// create client with the above options
	client, err := azcore.NewClient("azotel.SampleClient", internal.Version, azruntime.PipelineOptions{}, &options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start a new span
	ctx, span := client.Tracer().Start(context.TODO(), "sample_span", nil)

	// perform a simple HTTP GET
	req, err := azruntime.NewRequest(ctx, http.MethodGet, "https://www.microsoft.com/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// send the request with our client
	_, err = client.Pipeline().Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// span is over
	span.End()

	// shut down the tracing provider to flush all spans
	otelTP.Shutdown(context.TODO())

	// Output:
	// sample_span
}

// ConsoleExporter dumps span names to stdout.
type ConsoleExporter struct{}

// ExportSpans implements the tracesdk.SpanExporter interface for the ConsoleExporter type.
func (c ConsoleExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	for _, span := range spans {
		fmt.Println(span.Name())
	}
	return nil
}

// Shutdown implements the tracesdk.SpanExporter interface for the ConsoleExporter type.
func (c ConsoleExporter) Shutdown(ctx context.Context) error {
	return nil
}

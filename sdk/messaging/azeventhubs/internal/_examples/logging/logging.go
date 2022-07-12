// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

// This sample shows how to use OpenCensus to log internal logging from the Event Hubs
// package.
// NOTE: the OpenCensus PrintExporter is used for illustrative purposes but should be replaced
// by a proper exporter when used in production.
//  More information about OpenCensus and exporters can be found here: https://opencensus.io/quickstart/go/tracing/#configure-exporter

import (
	"context"
	"log"
	"os"

	eventhubs "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	_ "github.com/devigned/tab/opencensus"
	"github.com/joho/godotenv"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/trace"
)

func main() {
	godotenv.Load(".env")

	// NOTE: connection strings for Event Hubs must point to an entity path
	// Example: 'Endpoint=<sb://<host>;SharedAccessKeyName=<key name>;SharedAccessKey=<key>;EntityPath=<event hub name>'
	cs := os.Getenv("EVENTHUB_CONNECTION_STRING")

	// The OpenCensus PrintExporter is just an example. It's output is very verbose - you'll
	// want to choose a more suitable exporter for production use.
	// More information about OpenCensus exporters can be found here:
	// https://opencensus.io/quickstart/go/tracing/#configure-exporter
	trace.RegisterExporter(&exporter.PrintExporter{})

	// For debugging, it can be useful to disable sampling.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	hub, err := eventhubs.NewHubFromConnectionString(cs)

	if err != nil {
		log.Fatalf("Failed to create Event Hub client using a connection string: %s", err.Error())
	}

	// An example of trace output you'll see (among others):
	//
	// TraceID:      4b6c2c3cdbf8e8096e262d2770ff9b15
	// SpanID:       6682e012dbfdca72
	//
	// Span:    eh.Hub.Send
	// Status:   [0]
	// Elapsed: 485ms
	//
	// Attributes:
	//   - component=github.com/Azure/azure-event-hubs-go
	//   - version=3.3.9
	//   - peer.hostname=<your hostname>
	err = hub.Send(context.Background(), &eventhubs.Event{
		Data: []byte("hello world"),
	})

	if err != nil {
		log.Fatalf("Failed to send message: %s", err.Error())
	}

	err = hub.Close(context.Background())

	if err != nil {
		log.Fatalf("Failed to close Event Hub client: %s", err.Error())
	}
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azotel_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/tracing/azotel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Example_otlpExporter() {
	// end-to-end example creating an OTel TracerProvider that exports to Jaeger
	// via OTLP then uses the azotel adapter to connect it to an Azure SDK client.

	ctx := context.Background()

	// create the OTLP exporter, using the local Jaeger OTLP endpoint
	// NOTE: this assumes you have a local instance of Jaeger running on the same machine
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://localhost:4318/v1/traces"))
	if err != nil {
		log.Fatal(err)
	}

	// create an OTel TracerProvider that uses the OTLP exporter
	otelTP := otelsdk.NewTracerProvider(
		otelsdk.WithBatcher(exp),
		otelsdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Example_otlpExporter"),
		)),
	)

	// create a credential for the Azure SDK client
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	// create an Azure SDK client, connecting the OTel TracerProvider to it
	client, err := armresources.NewClientFactory("<subscription ID>", credential, &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			TracingProvider: azotel.NewTracingProvider(otelTP, nil),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// make various API calls with the client.  each one will create its own span
	_, err = client.NewClient().CheckExistenceByID(context.TODO(), "<resource ID>", "<api-version>", nil)
	if err != nil {
		log.Fatal(err)
	}

	// shut down the tracing provider to flush all spans to OTLP
	if err = otelTP.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

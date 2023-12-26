//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package publisher_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher"
)

// PublishEvents publishes events using the EventGrid schema to a topic. The
// topic must be configured to use the EventGrid schema or this will fail.
func ExampleClient_PublishEvents() {
	// ex: https://<topic-name>.<region>.eventgrid.azure.net/api/events
	endpoint := os.Getenv("EVENTGRID_TOPIC_ENDPOINT")
	key := os.Getenv("EVENTGRID_TOPIC_KEY")

	if endpoint == "" || key == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	// Other authentication methods:
	// - publisher.NewClient(): authenticate using a TokenCredential from azidentity.
	// - publisher.NewClientWithSAS(): authenticate using a SAS token.
	client, err := publisher.NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential(key), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	events := []publisher.Event{
		{
			Data:        "data for this event",
			DataVersion: to.Ptr("1.0"),
			EventType:   to.Ptr("event-type"),
			EventTime:   to.Ptr(time.Now()),
			ID:          to.Ptr("unique-id"),
			Subject:     to.Ptr("subject"),
		},
	}

	_, err = client.PublishEvents(context.TODO(), events, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Output:
}

// PublishCloudEvents publishes events using the CloudEvent schema to a topic. The
// topic must be configured to use the CloudEvent schema or this will fail.
func ExampleClient_PublishCloudEvents() {
	// ex: https://<topic-name>.<region>.eventgrid.azure.net/api/events
	endpoint := os.Getenv("EVENTGRID_CE_TOPIC_ENDPOINT")
	key := os.Getenv("EVENTGRID_CE_TOPIC_KEY")

	if endpoint == "" || key == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	// Other authentication methods:
	// - publisher.NewClient(): authenticate using a TokenCredential from azidentity.
	// - publisher.NewClientWithSAS(): authenticate using a SAS token.
	client, err := publisher.NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential(key), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	cloudEvent, err := messaging.NewCloudEvent("source", "eventtype", "data", nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	events := []messaging.CloudEvent{
		cloudEvent,
	}

	_, err = client.PublishCloudEvents(context.TODO(), events, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Output:
}

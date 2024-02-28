// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces"
)

func ExampleNewClientWithSharedKeyCredential() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	sharedKey := os.Getenv("EVENTGRID_KEY")

	if endpoint == "" || sharedKey == "" {
		return
	}

	client, err := aznamespaces.NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential(sharedKey), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleClient_PublishCloudEvents() {
	client := getEventGridClient()

	if client == nil {
		return
	}

	topic := os.Getenv("EVENTGRID_TOPIC")

	// CloudEvent is in github.com/Azure/azure-sdk-for-go/azcore/messaging and can be
	// used to transport

	// you can send a variety of different payloads, all of which can be encoded by messaging.CloudEvent
	var payloads = []any{
		[]byte{1, 2, 3},
		"hello world",
		struct{ Value string }{Value: "hello world"},
	}

	var eventsToSend []messaging.CloudEvent

	for _, payload := range payloads {
		event, err := messaging.NewCloudEvent("source", "eventType", payload, &messaging.CloudEventOptions{
			DataContentType: to.Ptr("application/octet-stream"),
		})

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		eventsToSend = append(eventsToSend, event)
	}

	_, err := client.PublishCloudEvents(context.TODO(), topic, eventsToSend, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Output:
}

func ExampleClient_ReceiveCloudEvents() {
	client := getEventGridClient()

	if client == nil {
		return
	}

	topic := os.Getenv("EVENTGRID_TOPIC")
	subscription := os.Getenv("EVENTGRID_SUBSCRIPTION")

	resp, err := client.ReceiveCloudEvents(context.TODO(), topic, subscription, &aznamespaces.ReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10), // in seconds
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, rd := range resp.Value {
		lockToken := rd.BrokerProperties.LockToken

		// NOTE: See the documentation for CloudEvent.Data on how your data
		// is deserialized.
		data := rd.Event.Data

		fmt.Fprintf(os.Stderr, "Event ID:%s, data: %#v, lockToken: %s\n", rd.Event.ID, data, *lockToken)

		// This will complete the message, deleting it from the subscription.
		resp, err := client.AcknowledgeCloudEvents(context.TODO(), topic, subscription, []string{*lockToken}, nil)

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		if len(resp.FailedLockTokens) > 0 {
			log.Fatalf("ERROR: %d events were not acknowledged", len(resp.FailedLockTokens))
		}
	}

	// Output:
}

func getEventGridClient() *aznamespaces.Client {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	sharedKey := os.Getenv("EVENTGRID_KEY")

	if endpoint == "" || sharedKey == "" {
		return nil
	}

	client, err := aznamespaces.NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential(sharedKey), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	return client
}

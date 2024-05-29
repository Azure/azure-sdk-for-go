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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces"
)

func ExampleNewReceiverClientWithSharedKeyCredential() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	sharedKey := os.Getenv("EVENTGRID_KEY")
	topic := os.Getenv("EVENTGRID_TOPIC")
	subscription := os.Getenv("EVENTGRID_SUBSCRIPTION")

	if endpoint == "" || sharedKey == "" || topic == "" || subscription == "" {
		return
	}

	client, err := aznamespaces.NewReceiverClientWithSharedKeyCredential(endpoint, topic, subscription, azcore.NewKeyCredential(sharedKey), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleNewReceiverClient() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	topic := os.Getenv("EVENTGRID_TOPIC")
	subscription := os.Getenv("EVENTGRID_SUBSCRIPTION")

	if endpoint == "" || topic == "" || subscription == "" {
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := aznamespaces.NewReceiverClient(endpoint, topic, subscription, tokenCredential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleNewSenderClientWithSharedKeyCredential() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	sharedKey := os.Getenv("EVENTGRID_KEY")
	topic := os.Getenv("EVENTGRID_TOPIC")

	if endpoint == "" || sharedKey == "" || topic == "" {
		return
	}

	client, err := aznamespaces.NewSenderClientWithSharedKeyCredential(endpoint, topic, azcore.NewKeyCredential(sharedKey), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleNewSenderClient() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	topic := os.Getenv("EVENTGRID_TOPIC")

	if endpoint == "" || topic == "" {
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := aznamespaces.NewSenderClient(endpoint, topic, tokenCredential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // ignore

	// Output:
}

func ExampleSenderClient_SendEvents() {
	sender, receiver := getEventGridClients()

	if sender == nil || receiver == nil {
		return
	}

	// CloudEvent is in github.com/Azure/azure-sdk-for-go/azcore/messaging and can be
	// used to transport

	// you can send a variety of different payloads, all of which can be encoded by messaging.CloudEvent
	var payloads = []any{
		[]byte{1, 2, 3},
		"hello world",
		struct{ Value string }{Value: "hello world"},
	}

	var eventsToSend []*messaging.CloudEvent

	for _, payload := range payloads {
		event, err := messaging.NewCloudEvent("source", "eventType", payload, &messaging.CloudEventOptions{
			DataContentType: to.Ptr("application/octet-stream"),
		})

		if err != nil {
			//  TODO: Update the following line with your application specific error handling logic
			log.Fatalf("ERROR: %s", err)
		}

		eventsToSend = append(eventsToSend, &event)
	}

	_, err := sender.SendEvents(context.TODO(), eventsToSend, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Output:
}

func ExampleReceiverClient_Receive() {
	sender, receiver := getEventGridClients()

	if sender == nil || receiver == nil {
		return
	}

	resp, err := receiver.Receive(context.TODO(), &aznamespaces.ReceiveOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10), // in seconds
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, rd := range resp.Details {
		lockToken := rd.BrokerProperties.LockToken

		// NOTE: See the documentation for CloudEvent.Data on how your data
		// is deserialized.
		data := rd.Event.Data

		fmt.Fprintf(os.Stderr, "Event ID:%s, data: %#v, lockToken: %s\n", rd.Event.ID, data, *lockToken)

		// This will complete the message, deleting it from the subscription.
		resp, err := receiver.Acknowledge(context.TODO(), []string{*lockToken}, nil)

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

func ExampleSenderClient_Send() {
	sender, receiver := getEventGridClients()

	if sender == nil || receiver == nil {
		return
	}

	// CloudEvent is in github.com/Azure/azure-sdk-for-go/azcore/messaging and can be
	// used to transport

	// you can send a variety of different payloads, all of which can be encoded by messaging.CloudEvent
	payload := []byte{1, 2, 3}
	eventToSend, err := messaging.NewCloudEvent("source", "eventType", payload, &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_, err = sender.Send(context.TODO(), &eventToSend, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Output:
}

// BinaryMode sends a CloudEvent more efficiently by avoiding unnecessary encoding of the Body.
// There are some caveats to be aware of:
//
//   - CloudEvent.Data must be of type []byte.
//   - CloudEvent.DataContentType will be used as the Content-Type for the HTTP request.
//   - CloudEvent.Extensions fields are converted to strings.
func ExampleSenderClient_Send_binaryMode() {
	sender, _ := getEventGridClients()

	if sender == nil {
		return
	}

	event, err := messaging.NewCloudEvent("source", "eventType", []byte{1, 2, 3}, &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// BinaryMode sends a CloudEvent more efficiently by avoiding unnecessary encoding of the Body.
	// There are some caveats to be aware of:
	// - [CloudEvent.Data] must be of type []byte.
	// - [CloudEvent.DataContentType] will be used as the Content-Type for the HTTP request.
	// - [CloudEvent.Extensions] fields are converted to strings.
	_, err = sender.Send(context.TODO(), &event, &aznamespaces.SendOptions{
		BinaryMode: true,
	})

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	// Output:
}

func getEventGridClients() (*aznamespaces.SenderClient, *aznamespaces.ReceiverClient) {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	sharedKey := os.Getenv("EVENTGRID_KEY")
	topic := os.Getenv("EVENTGRID_TOPIC")
	subscription := os.Getenv("EVENTGRID_SUBSCRIPTION")

	if endpoint == "" || sharedKey == "" || topic == "" || subscription == "" {
		return nil, nil
	}

	sender, err := aznamespaces.NewSenderClientWithSharedKeyCredential(endpoint, topic, azcore.NewKeyCredential(sharedKey), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	receiver, err := aznamespaces.NewReceiverClientWithSharedKeyCredential(endpoint, topic, subscription, azcore.NewKeyCredential(sharedKey), nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	return sender, receiver
}

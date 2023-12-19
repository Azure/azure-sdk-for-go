// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
)

func Example_publishAndReceiveCloudEvents() {
	endpoint := os.Getenv("EVENTGRID_ENDPOINT")
	key := os.Getenv("EVENTGRID_KEY")
	topicName := os.Getenv("EVENTGRID_TOPIC")
	subscriptionName := os.Getenv("EVENTGRID_SUBSCRIPTION")

	if endpoint == "" || key == "" || topicName == "" || subscriptionName == "" {
		return
	}

	client, err := azeventgrid.NewClientWithSharedKeyCredential(endpoint, azcore.NewKeyCredential(key), nil)

	if err != nil {
		panic(err)
	}

	//
	// Publish an event with a string payload
	//
	fmt.Fprintf(os.Stderr, "Published event with a string payload 'hello world'\n")
	eventWithString, err := publishAndReceiveEvent(client, topicName, subscriptionName, "application/json", "hello world")

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "Received an event with a string payload\n")
	fmt.Fprintf(os.Stderr, "ID: %s\n", eventWithString.Event.ID)

	var str *string

	if err := json.Unmarshal(eventWithString.Event.Data.([]byte), &str); err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "  Body: %s\n", *str) // prints 'Body: hello world'
	fmt.Fprintf(os.Stderr, "  Delivery count: %d\n", eventWithString.BrokerProperties.DeliveryCount)

	//
	// Publish an event with a []byte payload
	//
	eventWithBytes, err := publishAndReceiveEvent(client, topicName, subscriptionName, "application/octet-stream", []byte{0, 1, 2})

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "ID: %s\n", eventWithBytes.Event.ID)
	fmt.Fprintf(os.Stderr, "  Body: %#v\n", eventWithBytes.Event.Data.([]byte)) // prints 'Body: []byte{0x0, 0x1, 0x2}'
	fmt.Fprintf(os.Stderr, "  Delivery count: %d\n", eventWithBytes.BrokerProperties.DeliveryCount)

	//
	// Publish an event with a struct as the payload
	//
	type SampleData struct {
		Name string `json:"name"`
	}

	eventWithStruct, err := publishAndReceiveEvent(client, topicName, subscriptionName, "application/json", SampleData{Name: "hello"})

	if err != nil {
		panic(err)
	}

	var sampleData *SampleData
	if err := json.Unmarshal(eventWithStruct.Event.Data.([]byte), &sampleData); err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "ID: %s\n", eventWithStruct.Event.ID)
	fmt.Fprintf(os.Stderr, "  Body: %#v\n", sampleData) // prints 'Body: &azeventgrid_test.SampleData{Name:"hello"}'
	fmt.Fprintf(os.Stderr, "  Delivery count: %d\n", eventWithStruct.BrokerProperties.DeliveryCount)

	// Output:
}

func publishAndReceiveEvent(client *azeventgrid.Client, topicName string, subscriptionName string, dataContentType string, payload any) (azeventgrid.ReceiveDetails, error) {
	event, err := messaging.NewCloudEvent("source", "eventType", payload, &messaging.CloudEventOptions{
		DataContentType: &dataContentType,
	})

	if err != nil {
		return azeventgrid.ReceiveDetails{}, err
	}

	eventsToSend := []messaging.CloudEvent{
		event,
	}

	// NOTE: we're sending a single event as an example. For better efficiency it's best if you send
	// multiple events at a time.
	_, err = client.PublishCloudEvents(context.TODO(), topicName, eventsToSend, nil)

	if err != nil {
		return azeventgrid.ReceiveDetails{}, err
	}

	events, err := client.ReceiveCloudEvents(context.TODO(), topicName, subscriptionName, &azeventgrid.ReceiveCloudEventsOptions{
		MaxEvents: to.Ptr(int32(1)),

		// Wait for 60 seconds for events.
		MaxWaitTime: to.Ptr[int32](60),
	})

	if err != nil {
		return azeventgrid.ReceiveDetails{}, err
	}

	if len(events.Value) == 0 {
		return azeventgrid.ReceiveDetails{}, errors.New("no events received")
	}

	// We can (optionally) renew the lock (multiple times) if we want to continue to
	// extend the lock time on the event.
	_, err = client.RenewCloudEventLocks(context.TODO(), topicName, subscriptionName, []string{
		*events.Value[0].BrokerProperties.LockToken,
	}, nil)

	if err != nil {
		return azeventgrid.ReceiveDetails{}, err
	}

	// This acknowledges the event and causes it to be deleted from the subscription.
	// Other options are:
	// - client.ReleaseCloudEvents, which invalidates our event lock and allows another subscriber to receive the event.
	// - client.RejectCloudEvents, which rejects the event.
	//     If dead-lettering is configured, the event will be moved into the dead letter queue.
	//     Otherwise the event is deleted.
	ackResp, err := client.AcknowledgeCloudEvents(context.TODO(), topicName, subscriptionName, []string{
		*events.Value[0].BrokerProperties.LockToken,
	}, nil)

	if err != nil {
		return azeventgrid.ReceiveDetails{}, err
	}

	if len(ackResp.FailedLockTokens) > 0 {
		// some events failed when we tried to acknowledge them.
		for _, failed := range ackResp.FailedLockTokens {
			fmt.Printf("Failed to acknowledge event with lock token %s: %s\n", *failed.LockToken, failed.Error)
		}

		return azeventgrid.ReceiveDetails{}, errors.New("failed to acknowledge event")
	}

	return events.Value[0], nil
}

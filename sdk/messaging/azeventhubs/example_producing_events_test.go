// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

// Shows how to send events to an Event Hub partition using the [ProducerClient]
// and [EventDataBatch].
func Example_producingEventsUsingProducerClient() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	// Can also use a connection string:
	//
	// producerClient, err := azeventhubs.NewProducerClientFromConnectionString(connectionString, eventHubName, nil)
	//
	producerClient, err := azeventhubs.NewProducerClient(eventHubNamespace, eventHubName, defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}

	defer func() { _ = producerClient.Close(context.TODO()) }()

	events := createEventsForSample()

	newBatchOptions := &azeventhubs.EventDataBatchOptions{
		// The options allow you to control the size of the batch, as well as the partition it will get sent to.

		// PartitionID can be used to target a specific partition ID.
		// specific partition ID.
		//
		// PartitionID: partitionID,

		// PartitionKey can be used to ensure that messages that have the same key
		// will go to the same partition without requiring your application to specify
		// that partition ID.
		//
		// PartitionKey: partitionKey,

		//
		// Or, if you leave both PartitionID and PartitionKey nil, the service will choose a partition.
	}

	// Creates an EventDataBatch, which you can use to pack multiple events together, allowing for efficient transfer.
	batch, err := producerClient.NewEventDataBatch(context.TODO(), newBatchOptions)

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(events); i++ {
		err = batch.AddEventData(events[i], nil)

		if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
			if batch.NumEvents() == 0 {
				// This one event is too large for this batch, even on its own. No matter what we do it
				// will not be sendable at its current size.
				panic(err)
			}

			// This batch is full - we can send it and create a new one and continue
			// packaging and sending events.
			if err := producerClient.SendEventDataBatch(context.TODO(), batch, nil); err != nil {
				panic(err)
			}

			// create the next batch we'll use for events, ensuring that we use the same options
			// each time so all the messages go the same target.
			tmpBatch, err := producerClient.NewEventDataBatch(context.TODO(), newBatchOptions)

			if err != nil {
				panic(err)
			}

			batch = tmpBatch

			// rewind so we can retry adding this event to a batch
			i--
		} else if err != nil {
			panic(err)
		}
	}

	// if we have any events in the last batch, send it
	if batch.NumEvents() > 0 {
		if err := producerClient.SendEventDataBatch(context.TODO(), batch, nil); err != nil {
			panic(err)
		}
	}
}

func createEventsForSample() []*azeventhubs.EventData {
	return []*azeventhubs.EventData{
		{
			Body: []byte("hello"),
		},
		{
			Body: []byte("world"),
		},
	}
}

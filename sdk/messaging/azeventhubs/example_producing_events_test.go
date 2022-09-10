// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

func Example_producing() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	// Can also use a connection string:
	//
	// producerClient, err = azeventhubs.NewProducerClientFromConnectionString(connectionString, eventHubName, nil)
	//
	producerClient, err = azeventhubs.NewProducerClient(eventHubNamespace, eventHubName, defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}

	defer producerClient.Close(context.TODO())

	events := createEventsForSample()

	// Other examples:
	//
	// sending a batch to a specific partition:
	// batch, err := producerClient.NewEventDataBatch(context.TODO(), &azeventhubs.NewEventDataBatchOptions{
	// 	PartitionID: to.Ptr("0"),
	// })
	//
	// targeting a batch using a partition key
	// batch, err := producerClient.NewEventDataBatch(context.TODO(), &azeventhubs.NewEventDataBatchOptions{
	// 	PartitionKey: to.Ptr("partition key"),
	// })
	newBatchOptions := &azeventhubs.NewEventDataBatchOptions{}

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
			if err := producerClient.SendEventBatch(context.TODO(), batch, nil); err != nil {
				panic(err)
			}

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
		if err := producerClient.SendEventBatch(context.TODO(), batch, nil); err != nil {
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

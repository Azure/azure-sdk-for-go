// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

// Example_consumingEventsUsingConsumerClient shows how to start consuming events in partitions
// in an Event Hub.
//
// If you have an Azure Storage account you can use the [Processor] type instead, which will handle
// distributing partitions between multiple consumers. See example_processor_test.go for usage of
// that type.
//
// For an example of that see [example_processor_test.go].
//
// [example_processor_test.go]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go
func Example_consumingEventsUsingConsumerClient() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")

	fmt.Printf("Event Hub Namespace: %s, hubname: %s\n", eventHubNamespace, eventHubName)

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	// Can also use a connection string:
	//
	// consumerClient, err = azeventhubs.NewConsumerClientFromConnectionString(connectionString, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
	//
	consumerClient, err := azeventhubs.NewConsumerClient(eventHubNamespace, eventHubName, azeventhubs.DefaultConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}

	defer consumerClient.Close(context.TODO())

	// Consuming events from Azure Event Hubs requires an Event Hub AND a partition to receive from.
	// This differs from sending events where a partition ID is not required, although events are
	// assigned a partition when stored on the server side.
	//
	// For this example, we'll assume we don't know which partitions will have events, so we'll just
	// receive from all of them.
	eventHubProperties, err := consumerClient.GetEventHubProperties(context.Background(), nil)

	if err != nil {
		panic(err)
	}

	overallReceiveTime := 5 * time.Minute
	appCtx, cancelApp := context.WithTimeout(context.Background(), overallReceiveTime)
	defer cancelApp()

	fmt.Printf("Starting PartitionClients, will receive for %s...\n", overallReceiveTime)

	wg := sync.WaitGroup{}

	for _, tmpPartitionID := range eventHubProperties.PartitionIDs {
		wg.Add(1)

		go func(partitionID string) {
			defer wg.Done()

			// NOTE: If you're not changing any options for this function you can pass `nil` for the options
			// parameter.
			partitionClient, err := consumerClient.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
				// The default start position is Latest, which means this partition client will only start receiving
				// events stored AFTER the receiver has initialized.
				//
				// If you need to read events earlier than that or require a more deterministic start, you can use
				// the other fields in StartPosition, which allow you to choose a starting sequence number, offset
				// or even a timestamp.
				StartPosition: azeventhubs.StartPosition{
					Latest: to.Ptr(true),
				},
			})

			if err != nil {
				panic(err)
			}

			defer partitionClient.Close(context.TODO())

			fmt.Printf("[Partition: %s] Starting receive loop for partition\n", partitionID)

			for {
				// Using a context with a timeout will allow ReceiveEvents() to return with events it
				// collected in a minute, or earlier if it actually gets all 100 events we requested.
				receiveCtx, cancel := context.WithTimeout(appCtx, time.Minute)
				events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
				cancel()

				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					if appCtx.Err() != nil {
						// Gracefully exit. the 'defer partitionClient.Close(context.Background())' above
						// will take care of closing the PartitionClient.
						fmt.Printf("[Partition: %s] Application is stopping, stopping receive for partition\n", partitionID)
						break
					}

					// We didn't get any events before our context cancelled. Let's loop again.
					fmt.Printf("[Partition: %s] No events arrived in 1m, trying to receive again\n", partitionID)
				} else if err != nil {
					panic(err)
				}

				processConsumedEvents(partitionID, events)
			}
		}(tmpPartitionID)
	}

	fmt.Printf("Waiting for %s, for events to arrive on any partition", overallReceiveTime)
	wg.Wait()
	fmt.Printf("Done receiving events\n")
}

func processConsumedEvents(partitionID string, events []*azeventhubs.ReceivedEventData) {
	for _, event := range events {
		// We're assuming the Body is a byte-encoded string. EventData.Body supports any payload
		// that can be encoded to []byte.
		fmt.Printf("  [Partition: %s] Event received with body '%s'\n", partitionID, string(event.Body))
	}
}

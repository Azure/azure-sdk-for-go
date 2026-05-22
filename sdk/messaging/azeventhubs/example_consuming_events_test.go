// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2"
)

// Shows how to start consuming events in partitions in an Event Hub using the [ConsumerClient].
//
// If you have an Azure Storage account you can use the [Processor] type instead, which will handle
// distributing partitions between multiple consumers and storing progress using checkpoints.
// See [example_consuming_with_checkpoints_test.go] for an example.
//
// [example_consuming_with_checkpoints_test.go]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_consuming_with_checkpoints_test.go
func Example_consumingEventsUsingConsumerClient() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")
	partitionID := os.Getenv("EVENTHUB_PARTITION_ID")

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

	defer func() { _ = consumerClient.Close(context.TODO()) }()

	partitionClient, err := consumerClient.NewPartitionClient(partitionID, &azeventhubs.PartitionClientOptions{
		StartPosition: azeventhubs.StartPosition{
			Earliest: to.Ptr(true),
		},
	})

	if err != nil {
		panic(err)
	}

	defer func() { _ = partitionClient.Close(context.TODO()) }()

	// Will wait up to 1 minute for 100 events. If the context is cancelled (or expires)
	// you'll get any events that have been collected up to that point.
	receiveCtx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
	cancel()

	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		panic(err)
	}

	for _, event := range events {
		// We're assuming the Body is a byte-encoded string. EventData.Body supports any payload
		// that can be encoded to []byte.
		fmt.Printf("Event received with body '%s'\n", string(event.Body))
	}

	fmt.Printf("Done receiving events\n")
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

func Example_consuming() {
	eventHubNamespace := os.Getenv("EVENTHUB_NAMESPACE") // <ex: myeventhubnamespace.servicebus.windows.net>
	eventHubName := os.Getenv("EVENTHUB_NAME")
	eventHubPartitionID := os.Getenv("EVENTHUB_PARTITION")

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		panic(err)
	}

	// Can also use a connection string:
	//
	// consumerClient, err = azeventhubs.NewConsumerClientFromConnectionString(connectionString, "", "partition id", consumerGroup, nil)
	//
	consumerClient, err = azeventhubs.NewConsumerClient(eventHubNamespace, eventHubName, eventHubPartitionID, azeventhubs.DefaultConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		panic(err)
	}

	defer consumerClient.Close(context.TODO())

	for {
		// ReceiveEvents will wait until it either receives the # of events requested (100, in this call)
		// or if the context is cancelled, in which case it'll return any messages it has received.
		ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)

		events, err := consumerClient.ReceiveEvents(ctx, 100, nil)
		cancel()

		if err != nil {
			panic(err)
		}

		for _, event := range events {
			// process the event in some way
			fmt.Printf("Event received with body %v\n", event.Body)
		}
	}
}

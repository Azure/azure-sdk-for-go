// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
)

func Example_processor() {
	ehCS := os.Getenv("EVENTHUB_CONNECTION_STRING")
	eventHubName := os.Getenv("EVENTHUB_NAME")

	storageCS := os.Getenv("CHECKPOINTSTORE_STORAGE_CONNECTION_STRING")
	containerName := os.Getenv("CHECKPOINTSTORE_STORAGE_CONTAINER_NAME")

	// Create the checkpoint store
	// NOTE: the container must exist before the checkpoint store can be used.
	checkpointStore, err := checkpoints.NewBlobStoreFromConnectionString(storageCS, containerName, nil)

	if err != nil {
		panic(err)
	}

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(ehCS, eventHubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		panic(err)
	}

	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)

	if err != nil {
		panic(err)
	}

	go func() {
		// Loop continually - each time we acquire a new partition NextPartitionClient() will
		// return it.
		for {
			partitionClient := processor.NextPartitionClient(context.TODO())

			if partitionClient == nil {
				break
			}

			go func() {
				if err := processEvents(partitionClient); err != nil {
					panic(err)
				}
			}()
		}
	}()

	processorCtx, processorCancel := context.WithCancel(context.TODO())
	defer processorCancel()

	// Launch the load balancer - this will create new ProcessorPartitionClient's, which you can
	// retrieve by calls to NextPartitionClient, in a loop. This is demonstrated below.
	//
	// To stop the processor cancel the context that you passed in to Run().
	if err := processor.Run(processorCtx); err != nil {
		panic(err)
	}
}

func processEvents(partitionClient *azeventhubs.ProcessorPartitionClient) error {
	// initialize any resources needed to process the partition
	// This is the equivalent to PartitionOpen

	defer func() {
		// Do cleanup here, like shutting down database connections
		// or other resources used for processing this partition.
		partitionClient.Close(context.TODO())
	}()

	for {
		// wait for a minute for up to 100 events to arrive.
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), time.Minute)
		events, err := partitionClient.ReceiveEvents(receiveCtx, 100, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			if eventHubError := (*azeventhubs.Error)(nil); errors.As(err, &eventHubError) && eventHubError.Code == exported.CodeOwnershipLost {
				// This means that the partition was "stolen" - this can happen as partitions are balanced between
				// consumers.

				// the 'defer'd function we did above will take of any resource cleanup so we'll just exit the
				// function at this point.
				return nil
			}

			return err
		}

		fmt.Printf("Processing %d event(s)\n", len(events))

		if len(events) != 0 {
			// Update the checkpoint with the last event received. If the processor is restarted
			// it will resume from this point in the partition.
			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1]); err != nil {
				return err
			}
		}
	}
}

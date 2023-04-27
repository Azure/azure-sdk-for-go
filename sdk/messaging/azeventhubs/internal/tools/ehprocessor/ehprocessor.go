// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// This tool lets you consume events in two ways using the Processor. The Processor
// tracks progress and can balance load between itself and other Processors,
// storing checkpoint information to Azure Storage Blobs.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
)

func printProcessorExamples() {
	fmt.Fprintf(os.Stderr, "\n"+
		"Examples for processor:\n"+
		"  # Consume from multiple partitions, using the Processor and checkpointing\n"+
		"  ehprocessor -namespace <your event hub namespace>.servicebus.windows.net -eventhub tests -storageaccount https://<your storage account>.blob.core.windows.net -container <your storage container>\n"+
		"\n")
}

func processorCmd() error {
	eventHubNamespace := flag.String("namespace", "", "The fully qualified hostname of your Event Hub namespace (ex: <your event hub>.servicebus.windows.net)")
	eventHubName := flag.String("eventhub", "", "The name of your Event Hub")
	eventHubConsumerGroup := flag.String("consumergroup", azeventhubs.DefaultConsumerGroup, "The Event Hub consumer group used by your application")

	maxBatchWaitTime := flag.Duration("wait", 30*time.Second, "Max wait time for events, per batch")
	maxBatchSize := flag.Int("count", 1, "Maximum number of events to receive, per batch")

	storageAccountURL := flag.String("storageaccount", "", "The storage account URL used by your blob store (ex: https://<storage account name>.blob.core.windows.net)")
	storageContainerName := flag.String("container", "", "The storage container used by your checkpoints")

	verbose := flag.Bool("v", false, "Enable Azure SDK verbose logging")

	flag.Parse()

	if *eventHubName == "" || *eventHubNamespace == "" || *eventHubConsumerGroup == "" || *storageAccountURL == "" || *storageContainerName == "" {
		flag.PrintDefaults()
		printProcessorExamples()

		return errors.New("missing command line arguments")
	}

	if *verbose {
		azlog.SetEvents(azeventhubs.EventConsumer, azeventhubs.EventConn, azeventhubs.EventAuth, azeventhubs.EventProducer)
		azlog.SetListener(func(e azlog.Event, s string) {
			log.Printf("[%s] %s", e, s)
		})
	}

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return err
	}

	checkpointBlobStore, err := createCheckpointStore(storageAccountURL, defaultAzureCred, storageContainerName)

	if err != nil {
		return err
	}

	consumerClient, err := azeventhubs.NewConsumerClient(*eventHubNamespace, *eventHubName, *eventHubConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		return err
	}

	defer consumerClient.Close(context.Background())

	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointBlobStore, &azeventhubs.ProcessorOptions{
		LoadBalancingStrategy: azeventhubs.ProcessorStrategyGreedy,
	})

	if err != nil {
		return err
	}

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	dispatchPartitionClients := func() {
		log.Printf("PartitionClient dispatcher has started...")
		defer log.Printf("PartitionClient dispatcher has stopped.")

		for {
			pc := processor.NextPartitionClient(appCtx)

			if pc == nil {
				log.Println("Processor has stopped, stopping partition client dispatch loop")
				break
			}

			log.Printf("Acquired partition %s, receiving", pc.PartitionID())

			go processPartition(appCtx, pc, *maxBatchWaitTime, *maxBatchSize)
		}
	}

	go dispatchPartitionClients()

	log.Printf("Starting processor.")
	if err := processor.Run(appCtx); err != nil {
		return err
	}

	return nil
}

func createCheckpointStore(storageAccountURL *string, defaultAzureCred *azidentity.DefaultAzureCredential, storageContainerName *string) (azeventhubs.CheckpointStore, error) {
	blobClient, err := azblob.NewClient(*storageAccountURL, defaultAzureCred, nil)

	if err != nil {
		return nil, err
	}

	containerClient := blobClient.ServiceClient().NewContainerClient(*storageContainerName)

	log.Printf("Creating storage container %q, if it doesn't already exist", *storageContainerName)

	if _, err := containerClient.Create(context.Background(), nil); err != nil {
		if !bloberror.HasCode(err, bloberror.ContainerAlreadyExists) {
			return nil, err
		}
	}

	return checkpoints.NewBlobStore(containerClient, nil)
}

func processPartition(ctx context.Context, pc *azeventhubs.ProcessorPartitionClient, eventHubMaxTime time.Duration, eventHubMaxSize int) {
	defer pc.Close(ctx)

	for {
		receiveCtx, cancelReceive := context.WithTimeout(ctx, eventHubMaxTime)
		events, err := pc.ReceiveEvents(receiveCtx, eventHubMaxSize, nil)
		cancelReceive()

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			if ctx.Err() != nil { // parent cancelled
				break
			}

			// timing out without any events is fine. Continue receiving...
			continue
		} else if err != nil {
			log.Printf("ERROR while processing partition %q: %s", pc.PartitionID(), err)
			break
		}

		if len(events) > 0 {
			if err := printEventsAsJSON(pc.PartitionID(), events); err != nil {
				log.Printf("ERROR: failed when printing events: %s", err)
				break
			}

			latestEvent := events[len(events)-1]

			log.Printf("[%s] Updating checkpoint with offset: %d, sequenceNumber: %d", pc.PartitionID(), latestEvent.SequenceNumber, latestEvent.Offset)

			if err := pc.UpdateCheckpoint(ctx, latestEvent, nil); err != nil {
				log.Printf("ERROR: failed when updating checkpoint: %s", err)
			}
		}
	}
}

func printEventsAsJSON(partitionID string, events []*azeventhubs.ReceivedEventData) error {
	for _, evt := range events {
		var bodyBytes []int

		for _, b := range evt.Body {
			bodyBytes = append(bodyBytes, int(b))
		}

		// pick out some of the common fields
		jsonBytes, err := json.Marshal(struct {
			PartitionID    string
			MessageID      any
			BodyAsString   string
			Body           []int
			SequenceNumber int64
			Offset         int64
		}{partitionID, evt.MessageID, string(evt.Body), bodyBytes, evt.SequenceNumber, evt.Offset})

		if err != nil {
			return fmt.Errorf("Failed to marshal received event with message ID %v: %s", evt.MessageID, err.Error())
		}

		fmt.Printf("%s\n", string(jsonBytes))
	}

	return nil
}

func main() {
	if err := processorCmd(); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

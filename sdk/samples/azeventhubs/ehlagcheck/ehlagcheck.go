// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// This tool queries metadata from Event Hubs and checks it against information stored in the checkpoint
// store to calculate the "lag" between our Processors and the service. It's best used as a rough approximation
// of state as the data sources are not necessarily in-sync when updates occur frequently.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	if err := checkpointLagTool(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func checkpointLagTool(commandLineArgs []string) error {
	clients, err := createClients()

	if err != nil {
		return err
	}

	defer clients.ConsumerClient.Close(context.Background())

	eventHubProps, err := clients.ConsumerClient.GetEventHubProperties(context.Background(), nil)

	if err != nil {
		return err
	}

	checkpoints, err := clients.CheckpointStore.ListCheckpoints(context.Background(), clients.EventHubNamespace, clients.EventHubName, clients.EventHubConsumerGroup, nil)

	if err != nil {
		return err
	}

	checkpointsMap := map[string]*azeventhubs.Checkpoint{}

	for _, cp := range checkpoints {
		cp := cp
		checkpointsMap[cp.PartitionID] = &cp
	}
	ownerships, err := clients.CheckpointStore.ListOwnership(context.Background(), clients.EventHubNamespace, clients.EventHubName, clients.EventHubConsumerGroup, nil)

	if err != nil {
		return err
	}

	ownersMap := map[string]*azeventhubs.Ownership{}

	for _, o := range ownerships {
		o := o
		ownersMap[o.PartitionID] = &o
	}

	sort.Strings(eventHubProps.PartitionIDs)

	fmt.Fprintf(os.Stderr, "WARNING: Excessive querying of the checkpoint store/Event Hubs can impact application performance.\n")

	for _, partID := range eventHubProps.PartitionIDs {
		partID := partID

		cp, o := checkpointsMap[partID], ownersMap[partID]

		partProps, err := clients.ConsumerClient.GetPartitionProperties(context.Background(), partID, nil)

		if err != nil {
			return err
		}

		fmt.Printf("Partition ID %q\n", partID)

		if o != nil {
			fmt.Printf("  Owner ID: %q, last updated: %s\n", o.OwnerID, o.LastModifiedTime.Format(time.RFC3339))
		} else {
			fmt.Printf("  Owner ID: <no owner>\n")
		}

		fmt.Printf("  Last enqueued sequence number is %d\n", partProps.LastEnqueuedSequenceNumber)

		if cp != nil && cp.SequenceNumber != nil {
			fmt.Printf("  Delta (between service and checkpoint): %d\n", partProps.LastEnqueuedSequenceNumber-*cp.SequenceNumber)
		}
	}

	return nil
}

type clients struct {
	ConsumerClient  *azeventhubs.ConsumerClient
	CheckpointStore *checkpoints.BlobStore

	EventHubNamespace     string
	EventHubName          string
	EventHubConsumerGroup string
}

func createClients() (*clients, error) {
	eventHubNamespace := flag.String("namespace", "", "The fully qualified hostname of your Event Hub namespace (ex: <your event hub>.servicebus.windows.net)")
	eventHubName := flag.String("eventhub", "", "The name of your Event Hub")
	eventHubConsumerGroup := flag.String("consumergroup", azeventhubs.DefaultConsumerGroup, "The Event Hub consumer group used by your application")

	storageAccountURL := flag.String("storageaccount", "", "The storage account URL used by your blob store (ex: https://<storage account name>.blob.core.windows.net/)")
	storageContainerName := flag.String("container", "", "The storage container used by your checkpoints")

	flag.Parse()

	if *eventHubNamespace == "" || *eventHubName == "" || *eventHubConsumerGroup == "" || *storageAccountURL == "" || *storageContainerName == "" {
		flag.PrintDefaults()
		return nil, errors.New("Missing command line arguments")
	}

	defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	blobClient, err := azblob.NewClient(*storageAccountURL, defaultAzureCred, nil)

	if err != nil {
		return nil, err
	}

	blobStore, err := checkpoints.NewBlobStore(blobClient.ServiceClient().NewContainerClient(*storageContainerName), nil)

	if err != nil {
		return nil, err
	}

	// Both ProducerClient and ConsumerClient can query event hub partition properties.
	cc, err := azeventhubs.NewConsumerClient(*eventHubNamespace, *eventHubName, *eventHubConsumerGroup, defaultAzureCred, nil)

	if err != nil {
		return nil, err
	}

	return &clients{
		ConsumerClient:        cc,
		CheckpointStore:       blobStore,
		EventHubNamespace:     *eventHubNamespace,
		EventHubName:          *eventHubName,
		EventHubConsumerGroup: *eventHubConsumerGroup,
	}, nil
}

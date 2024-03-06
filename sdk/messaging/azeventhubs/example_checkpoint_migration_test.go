// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type LegacyCheckpoint struct {
	PartitionID string `json:"partitionID"`
	Epoch       int    `json:"epoch"`
	Owner       string `json:"owner"`
	Checkpoint  struct {
		Offset         string `json:"offset"`
		SequenceNumber int64  `json:"sequenceNumber"`
		EnqueueTime    string `json:"enqueueTime"` // ": "0001-01-01T00:00:00Z"
	} `json:"checkpoint"`
}

// Shows how to migrate from the older `github.com/Azure/azure-event-hubs-go` checkpointer to to
// the format used by this package, `github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints/BlobStore`
//
// NOTE: This example is not safe to run while either the old or new checkpoint store is in-use as it doesn't
// respect locking or ownership.
func Example_migrateCheckpoints() {
	// Azure Event Hubs connection string. You can get this from the Azure Portal.
	// For example: youreventhub.servicebus.windows.net
	var EventHubNamespace = os.Getenv("EVENTHUB_NAMESPACE")

	// Name of your Event Hub that these checkpoints reference.
	var EventHubName = os.Getenv("EVENTHUB_NAME")

	// Name of your Event Hub consumer group
	// Example: $Default
	var EventHubConsumerGroup = os.Getenv("EVENTHUB_CONSUMER_GROUP")

	// Azure Storage account connection string. You can get this from the Azure Portal.
	// For example: DefaultEndpointsProtocol=https;AccountName=accountname;AccountKey=account-key;EndpointSuffix=core.windows.net
	var StorageConnectionString = os.Getenv("STORAGE_CONNECTION_STRING")

	// Optional: If you used `eventhub.WithPrefixInBlobPath()` configuration option for your Event Processor Host
	// then you'll need to set this value.
	//
	// NOTE: This is no longer needed with the new checkpoint store as it automatically makes the path unique
	// for each combination of eventhub + hubname + consumergroup + partition.
	var BlobPrefix = os.Getenv("OLD_STORAGE_BLOB_PREFIX")

	// Name of the checkpoint store's Azure Storage container.
	var OldStorageContainerName = os.Getenv("OLD_STORAGE_CONTAINER_NAME")

	// Name of the Azure Storage container to place new checkpoints in.
	var NewStorageContainerName = os.Getenv("NEW_STORAGE_CONTAINER_NAME")

	if EventHubNamespace == "" || EventHubName == "" || EventHubConsumerGroup == "" ||
		StorageConnectionString == "" || OldStorageContainerName == "" || NewStorageContainerName == "" {
		fmt.Printf("Skipping migration, missing parameters\n")
		return
	}

	blobClient, err := azblob.NewClientFromConnectionString(StorageConnectionString, nil)

	if err != nil {
		panic(err)
	}

	oldCheckpoints, err := loadOldCheckpoints(blobClient, OldStorageContainerName, BlobPrefix)

	if err != nil {
		panic(err)
	}

	newCheckpointStore, err := checkpoints.NewBlobStore(blobClient.ServiceClient().NewContainerClient(NewStorageContainerName), nil)

	if err != nil {
		panic(err)
	}

	for _, oldCheckpoint := range oldCheckpoints {
		newCheckpoint := azeventhubs.Checkpoint{
			ConsumerGroup:           EventHubConsumerGroup,
			EventHubName:            EventHubName,
			FullyQualifiedNamespace: EventHubNamespace,
			PartitionID:             oldCheckpoint.PartitionID,
		}

		offset, err := strconv.ParseInt(oldCheckpoint.Checkpoint.Offset, 10, 64)

		if err != nil {
			panic(err)
		}

		newCheckpoint.Offset = &offset
		newCheckpoint.SequenceNumber = &oldCheckpoint.Checkpoint.SequenceNumber

		if err := newCheckpointStore.SetCheckpoint(context.Background(), newCheckpoint, nil); err != nil {
			panic(err)
		}
	}
}

func loadOldCheckpoints(blobClient *azblob.Client, containerName string, customBlobPrefix string) ([]*LegacyCheckpoint, error) {
	blobPrefix := &customBlobPrefix

	if customBlobPrefix == "" {
		blobPrefix = nil
	}

	pager := blobClient.NewListBlobsFlatPager(containerName, &container.ListBlobsFlatOptions{
		Prefix: blobPrefix,
	})

	var checkpoints []*LegacyCheckpoint

	for pager.More() {
		page, err := pager.NextPage(context.Background())

		if err != nil {
			return nil, err
		}

		for _, item := range page.Segment.BlobItems {
			buff := [4000]byte{}

			len, err := blobClient.DownloadBuffer(context.Background(), containerName, *item.Name, buff[:], nil)

			if err != nil {
				return nil, err
			}

			var legacyCheckpoint *LegacyCheckpoint

			if err := json.Unmarshal(buff[0:len], &legacyCheckpoint); err != nil {
				return nil, err
			}

			checkpoints = append(checkpoints, legacyCheckpoint)
		}
	}

	return checkpoints, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package checkpoints_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestBlobStore_Checkpoints(t *testing.T) {
	testData := getContainerClient(t)
	defer testData.Cleanup()

	cc, err := container.NewClientFromConnectionString(testData.ConnectionString, testData.ContainerName, nil)
	require.NoError(t, err)

	store, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	checkpoints, err := store.ListCheckpoints(context.Background(), "fully-qualified-namespace", "event-hub-name", "consumer-group", nil)
	require.NoError(t, err)
	require.Empty(t, checkpoints)

	err = store.UpdateCheckpoint(context.Background(), azeventhubs.Checkpoint{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		Offset:                  to.Ptr[int64](101),
		SequenceNumber:          to.Ptr[int64](202),
	}, nil)
	require.NoError(t, err)

	checkpoints, err = store.ListCheckpoints(context.Background(), "ns.servicebus.windows.net", "event-hub-name", "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, azeventhubs.Checkpoint{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		Offset:                  to.Ptr[int64](101),
		SequenceNumber:          to.Ptr[int64](202),
	}, checkpoints[0])

	// There's a code path to allow updating the blob after it's been created but without an etag
	// in which case it just updates it.
	err = store.UpdateCheckpoint(context.Background(), azeventhubs.Checkpoint{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		Offset:                  to.Ptr[int64](102),
		SequenceNumber:          to.Ptr[int64](203),
	}, nil)
	require.NoError(t, err)

	checkpoints, err = store.ListCheckpoints(context.Background(), "ns.servicebus.windows.net", "event-hub-name", "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, azeventhubs.Checkpoint{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		Offset:                  to.Ptr[int64](102),
		SequenceNumber:          to.Ptr[int64](203),
	}, checkpoints[0])
}

func TestBlobStore_Ownership(t *testing.T) {
	testData := getContainerClient(t)
	defer testData.Cleanup()

	cc, err := container.NewClientFromConnectionString(testData.ConnectionString, testData.ContainerName, nil)
	require.NoError(t, err)

	store, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	ownerships, err := store.ListOwnership(context.Background(), "fully-qualified-namespace", "event-hub-name", "consumer-group", nil)
	require.NoError(t, err)
	require.Empty(t, ownerships, "no ownerships yet")

	ownerships, err = store.ClaimOwnership(context.Background(), nil, nil)
	require.NoError(t, err)
	require.Empty(t, ownerships)

	ownerships, err = store.ClaimOwnership(context.Background(), []azeventhubs.Ownership{}, nil)
	require.NoError(t, err)
	require.Empty(t, ownerships)

	ownerships, err = store.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
		{
			ConsumerGroup:           "$Default",
			EventHubName:            "event-hub-name",
			FullyQualifiedNamespace: "ns.servicebus.windows.net",
			PartitionID:             "partition-id",
			OwnerID:                 "owner-id",
		},
	}, nil)
	require.NoError(t, err)

	etagAfterFirstClaim := ownerships[0].ETag
	require.NotEmpty(t, ownerships[0].ETag)
	require.NotZero(t, ownerships[0].LastModifiedTime)

	require.Equal(t, azeventhubs.Ownership{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		OwnerID:                 "owner-id",
		ETag:                    ownerships[0].ETag,
		LastModifiedTime:        ownerships[0].LastModifiedTime,
	}, ownerships[0])

	// if we attempt to claim it with a non-matching etag it will fail to claim
	// but not fail the call.
	ownerships, err = store.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
		{
			ConsumerGroup:           "$Default",
			EventHubName:            "event-hub-name",
			FullyQualifiedNamespace: "ns.servicebus.windows.net",
			PartitionID:             "partition-id",
			OwnerID:                 "owner-id",
			ETag:                    to.Ptr(azcore.ETag("non-matching-etag")),
		},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ownerships, "we're out of date (based on the non-matching etag), so no ownerships were claimed")

	// now we'll use the actual etag
	ownerships, err = store.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
		{
			ConsumerGroup:           "$Default",
			EventHubName:            "event-hub-name",
			FullyQualifiedNamespace: "ns.servicebus.windows.net",
			PartitionID:             "partition-id",
			OwnerID:                 "owner-id",
			ETag:                    etagAfterFirstClaim,
		},
	}, nil)
	require.NoError(t, err)

	require.Equal(t, azeventhubs.Ownership{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		OwnerID:                 "owner-id",
		ETag:                    ownerships[0].ETag,
		LastModifiedTime:        ownerships[0].LastModifiedTime,
	}, ownerships[0])

	// etag definitely got updated.
	require.NotEqual(t, etagAfterFirstClaim, ownerships[0].ETag)
	require.NotZero(t, ownerships[0].LastModifiedTime)
}

func TestBlobStore_ListAndClaim(t *testing.T) {
	// listing ownerships is a slightly different code path
	testData := getContainerClient(t)
	defer testData.Cleanup()

	cc, err := container.NewClientFromConnectionString(testData.ConnectionString, testData.ContainerName, nil)
	require.NoError(t, err)

	store, err := checkpoints.NewBlobStore(cc, nil)
	require.NoError(t, err)

	claimedOwnerships, err := store.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
		{
			ConsumerGroup:           "$Default",
			EventHubName:            "event-hub-name",
			FullyQualifiedNamespace: "ns.servicebus.windows.net",
			PartitionID:             "partition-id",
			OwnerID:                 "first-client",
		},
	}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, claimedOwnerships)

	listedOwnerships, err := store.ListOwnership(context.Background(), "ns.servicebus.windows.net", "event-hub-name", "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, "first-client", listedOwnerships[0].OwnerID)
	require.NotEmpty(t, listedOwnerships[0].ETag)
	require.NotZero(t, listedOwnerships[0].LastModifiedTime)

	require.Equal(t, "$Default", listedOwnerships[0].ConsumerGroup)
	require.Equal(t, "event-hub-name", listedOwnerships[0].EventHubName)
	require.Equal(t, "ns.servicebus.windows.net", listedOwnerships[0].FullyQualifiedNamespace)
	require.Equal(t, "partition-id", listedOwnerships[0].PartitionID)

	// update using the etag
	claimedOwnerships, err = store.ClaimOwnership(context.Background(), listedOwnerships, nil)
	require.NoError(t, err)

	require.Equal(t, "partition-id", claimedOwnerships[0].PartitionID)

	// try to do it again and it'll fail since we don't have an updated etag
	claimedOwnerships, err = store.ClaimOwnership(context.Background(), listedOwnerships, nil)
	require.NoError(t, err)
	require.Empty(t, claimedOwnerships)
}

func getContainerClient(t *testing.T) struct {
	ConnectionString string
	ContainerName    string
	Cleanup          func()
} {
	_ = godotenv.Load("../.env")

	storageCS := os.Getenv("CHECKPOINTSTORE_STORAGE_CONNECTION_STRING")

	if storageCS == "" {
		t.Skipf("CHECKPOINTSTORE_STORAGE_CONNECTION_STRING is not defined in the environment. Skipping blob checkpoint store live tests")
		return struct {
			ConnectionString string
			ContainerName    string
			Cleanup          func()
		}{}
	}

	nano := time.Now().UTC().UnixNano()

	containerName := strconv.FormatInt(nano, 10)
	client, err := container.NewClientFromConnectionString(storageCS, containerName, nil)
	require.NoError(t, err)

	_, err = client.Create(context.Background(), nil)
	require.NoError(t, err)

	return struct {
		ConnectionString string
		ContainerName    string
		Cleanup          func()
	}{
		ConnectionString: storageCS,
		ContainerName:    containerName,
		Cleanup: func() {
			_, err := client.Delete(context.Background(), nil)
			require.NoError(t, err)
		},
	}
}

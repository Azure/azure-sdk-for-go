// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package checkpoints_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestBlobStore_Checkpoints(t *testing.T) {
	testData := newBlobStoreTestData(t)

	checkpoints, err := testData.BlobStore.ListCheckpoints(context.Background(), "fully-qualified-namespace", "event-hub-name", "consumer-group", nil)
	require.NoError(t, err)
	require.Empty(t, checkpoints)

	err = testData.BlobStore.SetCheckpoint(context.Background(), azeventhubs.Checkpoint{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		Offset:                  to.Ptr[int64](101),
		SequenceNumber:          to.Ptr[int64](202),
	}, nil)
	require.NoError(t, err)

	checkpoints, err = testData.BlobStore.ListCheckpoints(context.Background(), "ns.servicebus.windows.net", "event-hub-name", "$Default", nil)
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
	err = testData.BlobStore.SetCheckpoint(context.Background(), azeventhubs.Checkpoint{
		ConsumerGroup:           "$Default",
		EventHubName:            "event-hub-name",
		FullyQualifiedNamespace: "ns.servicebus.windows.net",
		PartitionID:             "partition-id",
		Offset:                  to.Ptr[int64](102),
		SequenceNumber:          to.Ptr[int64](203),
	}, nil)
	require.NoError(t, err)

	checkpoints, err = testData.BlobStore.ListCheckpoints(context.Background(), "ns.servicebus.windows.net", "event-hub-name", "$Default", nil)
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
	testData := newBlobStoreTestData(t)

	ownerships, err := testData.BlobStore.ListOwnership(context.Background(), "fully-qualified-namespace", "event-hub-name", "consumer-group", nil)
	require.NoError(t, err)
	require.Empty(t, ownerships, "no ownerships yet")

	ownerships, err = testData.BlobStore.ClaimOwnership(context.Background(), nil, nil)
	require.NoError(t, err)
	require.Empty(t, ownerships)

	ownerships, err = testData.BlobStore.ClaimOwnership(context.Background(), []azeventhubs.Ownership{}, nil)
	require.NoError(t, err)
	require.Empty(t, ownerships)

	ownerships, err = testData.BlobStore.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
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
	ownerships, err = testData.BlobStore.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
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
	ownerships, err = testData.BlobStore.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
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
	testData := newBlobStoreTestData(t)

	claimedOwnerships, err := testData.BlobStore.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
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

	listedOwnerships, err := testData.BlobStore.ListOwnership(context.Background(), "ns.servicebus.windows.net", "event-hub-name", "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, "first-client", listedOwnerships[0].OwnerID)
	require.NotEmpty(t, listedOwnerships[0].ETag)
	require.NotZero(t, listedOwnerships[0].LastModifiedTime)

	require.Equal(t, "$Default", listedOwnerships[0].ConsumerGroup)
	require.Equal(t, "event-hub-name", listedOwnerships[0].EventHubName)
	require.Equal(t, "ns.servicebus.windows.net", listedOwnerships[0].FullyQualifiedNamespace)
	require.Equal(t, "partition-id", listedOwnerships[0].PartitionID)

	// update using the etag
	claimedOwnerships, err = testData.BlobStore.ClaimOwnership(context.Background(), listedOwnerships, nil)
	require.NoError(t, err)

	require.Equal(t, "partition-id", claimedOwnerships[0].PartitionID)

	// try to do it again and it'll fail since we don't have an updated etag
	claimedOwnerships, err = testData.BlobStore.ClaimOwnership(context.Background(), listedOwnerships, nil)
	require.NoError(t, err)
	require.Empty(t, claimedOwnerships)
}

func TestBlobStore_OnlyOneOwnershipClaimSucceeds(t *testing.T) {
	testData := newBlobStoreTestData(t)

	// we're going to make multiple calls to the blob store but only _one_ should succeed
	// since it's "first one in wins"
	claimsCh := make(chan []azeventhubs.Ownership, 20)

	t.Logf("Starting %d goroutines to claim ownership without an etag", cap(claimsCh))

	// attempt to claim the same partition from multiple goroutines. Only _one_ of the
	// goroutines should walk away thinking it claimed the partition.
	for i := 0; i < cap(claimsCh); i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			ownerships, err := testData.BlobStore.ClaimOwnership(ctx, []azeventhubs.Ownership{
				{ConsumerGroup: azeventhubs.DefaultConsumerGroup, EventHubName: "name", FullyQualifiedNamespace: "ns", PartitionID: "0", OwnerID: "ownerID"},
			}, nil)

			if err != nil {
				claimsCh <- nil
				require.NoError(t, err)
			} else {
				claimsCh <- ownerships
			}
		}()
	}

	claimed := map[string]bool{}
	numFailedClaims := 0

	for i := 0; i < cap(claimsCh); i++ {
		claims := <-claimsCh

		if claims == nil {
			numFailedClaims++
			continue
		}

		for _, claim := range claims {
			require.False(t, claimed[claim.PartitionID], fmt.Sprintf("Partition ID %s was claimed more than once", claim.PartitionID))
			require.NotNil(t, claim.ETag)
			claimed[claim.PartitionID] = true
		}
	}

	require.Equal(t, cap(claimsCh)-1, numFailedClaims, fmt.Sprintf("One of the 1/%d wins and the rest all fail to claim", cap(claimsCh)))
}

func TestBlobStore_OnlyOneOwnershipUpdateSucceeds(t *testing.T) {
	testData := newBlobStoreTestData(t)

	// we're going to make multiple calls to the blob store but only _one_ should succeed
	// since it's "first one in wins"
	claimsCh := make(chan []azeventhubs.Ownership, 20)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ownerships, err := testData.BlobStore.ClaimOwnership(ctx, []azeventhubs.Ownership{
		{ConsumerGroup: azeventhubs.DefaultConsumerGroup, EventHubName: "name", FullyQualifiedNamespace: "ns", PartitionID: "0", OwnerID: "ownerID"},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, "0", ownerships[0].PartitionID)
	require.NotNil(t, ownerships[0].ETag)

	t.Logf("Starting %d goroutines to claim ownership without an etag", cap(claimsCh))

	// attempt to claim the same partition from multiple goroutines. Only _one_ of the
	// goroutines should walk away thinking it claimed the partition.
	for i := 0; i < cap(claimsCh); i++ {
		go func() {

			ownerships, err := testData.BlobStore.ClaimOwnership(ctx, ownerships, nil)

			if err != nil {
				claimsCh <- nil
				require.NoError(t, err)
			} else {
				claimsCh <- ownerships
			}
		}()
	}

	claimed := map[string]bool{}
	numFailedClaims := 0

	for i := 0; i < cap(claimsCh); i++ {
		claims := <-claimsCh

		if claims == nil {
			numFailedClaims++
			continue
		}

		for _, claim := range claims {
			require.False(t, claimed[claim.PartitionID], fmt.Sprintf("Partition ID %s was claimed more than once", claim.PartitionID))
			require.NotNil(t, claim.ETag)
			claimed[claim.PartitionID] = true
		}
	}

	require.Equal(t, cap(claimsCh)-1, numFailedClaims, fmt.Sprintf("One of the 1/%d wins and the rest all fail to claim", cap(claimsCh)))
}

func TestBlobStore_RelinquishClaim(t *testing.T) {
	testData := newBlobStoreTestData(t)

	initialClaims, err := testData.BlobStore.ClaimOwnership(context.Background(), []azeventhubs.Ownership{
		{
			ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
			EventHubName:            "eventhubname",
			FullyQualifiedNamespace: "fullyQualifiedNamespace",
			PartitionID:             "partitionID",
			OwnerID:                 "ownerID",
			LastModifiedTime:        time.Now().UTC(),
		},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, "ownerID", initialClaims[0].OwnerID)

	// relinquish our ownership claim
	initialClaims[0].OwnerID = ""
	relinquishedClaims, err := testData.BlobStore.ClaimOwnership(context.Background(), initialClaims, nil)
	require.NoError(t, err)
	require.Empty(t, relinquishedClaims[0].OwnerID)

	// now be some other person and claim it.
	relinquishedClaims[0].OwnerID = "new owner!"
	lastClaimed, err := testData.BlobStore.ClaimOwnership(context.Background(), relinquishedClaims, nil)
	require.NoError(t, err)
	require.Equal(t, "new owner!", lastClaimed[0].OwnerID)
}

type blobStoreTestData struct {
	CC        *container.Client
	BlobStore *checkpoints.BlobStore
}

// newBlobStoreTestData creates an Azure Blob storage container
// and returns the associated ContainerClient and BlobStore instance.
func newBlobStoreTestData(t *testing.T) blobStoreTestData {
	_ = godotenv.Load("../.env")

	storageEndpoint := os.Getenv("CHECKPOINTSTORE_STORAGE_ENDPOINT")

	if storageEndpoint == "" {
		t.Skipf("CHECKPOINTSTORE_STORAGE_ENDPOINT is not defined in the environment. Skipping blob checkpoint store live tests")
		return blobStoreTestData{}
	}

	nano := time.Now().UTC().UnixNano()

	containerName := strconv.FormatInt(nano, 10)

	cred, err := credential.New(nil)
	require.NoError(t, err)

	containerURL := test.URLJoinPaths(storageEndpoint, containerName)
	require.NoError(t, err)

	client, err := container.NewClient(containerURL, cred, nil)
	require.NoError(t, err)

	_, err = client.Create(context.Background(), nil)
	require.NoError(t, err)

	blobStore, err := checkpoints.NewBlobStore(client, nil)
	require.NoError(t, err)

	t.Cleanup(func() {

	})

	return blobStoreTestData{
		CC:        client,
		BlobStore: blobStore,
	}
}

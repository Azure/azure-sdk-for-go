// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package checkpoints_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
)

// The blob checkpoint store SHOULD have been storing blobs with lowercased names (this is the de-facto standard amongst
// all the checkpoint stores in other languages), but we were allowing potential mixed-case names.
//
// To fix this we added an in-place migration that occurs as the checkpoint store is used. These tests cover the code paths
// for migration. See `checkpoint_stress_tester.go` for a stress test that runs with larger concurrency numbers and larger
// numbers of checkpoints.

func TestOwnershipBlob_IncorrectCaseExists(t *testing.T) {
	testData := newBlobStoreTestData(t)

	// Note that all of these have some form of mixed case
	orig := azeventhubs.Ownership{
		EventHubName:            "eventHub",
		ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
		FullyQualifiedNamespace: "eventhHubsTest.servicebus.windows.net",
		PartitionID:             "0",
		OwnerID:                 "owner id",
	}

	legacyName := uploadLegacyOwnershipBlob(t, testData.CC, orig)

	ownerships, err := testData.BlobStore.ListOwnership(context.Background(), orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, nil)
	require.NoError(t, err)

	// the ownership that comes back here will NOT have an etag (despite the legacy blob existing), which forces anyone that wants
	// to claim this ownership to _create_ a new blob.
	assertOwnerships(t, []azeventhubs.Ownership{orig}, ownerships, false)

	//
	// Now we'll simulate a load balancing cycle of "list ownerships, attempt to claim ownerships"
	//

	claimed, err := testData.BlobStore.ClaimOwnership(context.Background(), ownerships, nil)
	require.NoError(t, err)
	require.NotEmpty(t, claimed)
	require.NotNil(t, claimed[0].ETag)

	// at this point the old blob with the mixed-case name should be deleted, and a new blob, with a completely lowercase
	// name will be there.
	allBlobs := listAllBlobs(t, testData.CC)
	require.Equal(t, 1, len(allBlobs), "Old blob is deleted, new blob is the only one that exists")
	require.Equal(t, strings.ToLower(legacyName), *allBlobs[0].Name)

	// and now things are "normal" - we should get an etag, based on the new lowercased blob.
	ownerships, err = testData.BlobStore.ListOwnership(context.Background(), orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, nil)
	require.NoError(t, err)

	assertOwnerships(t, []azeventhubs.Ownership{
		{
			ConsumerGroup:           strings.ToLower(orig.ConsumerGroup),
			EventHubName:            strings.ToLower(orig.EventHubName),
			FullyQualifiedNamespace: strings.ToLower(orig.FullyQualifiedNamespace),
			PartitionID:             strings.ToLower(orig.PartitionID),
			OwnerID:                 orig.OwnerID, // this casing doesn't affect the blob, it's specified by the user.
		},
	}, ownerships, true)
	require.Equal(t, *allBlobs[0].Properties.ETag, *ownerships[0].ETag)

	// simple check - try to claim ownership again.
	newOwnerships, err := testData.BlobStore.ClaimOwnership(context.Background(), ownerships, nil)
	require.NoError(t, err)

	assertOwnerships(t, ownerships, newOwnerships, true)
}

func TestOwnershipBlob_IncorrectCaseExists_Parallel(t *testing.T) {
	testData := newBlobStoreTestData(t)

	// Note that all of these have some form of mixed case
	orig := azeventhubs.Ownership{
		EventHubName:            "eventHub",
		ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
		FullyQualifiedNamespace: "eventhHubsTest.servicebus.windows.net",
		PartitionID:             "0",
		OwnerID:                 "owner id",
	}

	// note that we don't normalize the case in our "legacy" blob names
	legacyName := fmt.Sprintf("%s/%s/%s/ownership/%s", orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, orig.PartitionID)

	blobClient := testData.CC.NewBlockBlobClient(legacyName)

	_, err := blobClient.UploadBuffer(context.Background(), nil, &blockblob.UploadBufferOptions{
		Metadata: map[string]*string{
			"ownerid": &orig.OwnerID,
		},
	})
	require.NoError(t, err)

	allReady := &sync.WaitGroup{}
	allDone := &sync.WaitGroup{}
	startCh := make(chan struct{})
	claimantCh := make(chan int, 1)

	ownerships, err := testData.BlobStore.ListOwnership(context.Background(), orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, nil)
	require.NoError(t, err)
	require.NotEmpty(t, ownerships)

	// now attempt to claim in multiple goroutines
	for i := 0; i < 100; i++ {
		allReady.Add(1)
		allDone.Add(1)

		go func(i int) {
			allReady.Done()
			defer allDone.Done()

			<-startCh
			claimed, err := testData.BlobStore.ClaimOwnership(context.Background(), ownerships, nil)
			require.NoError(t, err)

			if len(claimed) == 1 {
				select {
				case claimantCh <- i:
				default:
					require.FailNow(t, "Another goroutine has claimed ownership")
				}
			}
		}(i)
	}

	allReady.Wait()
	close(startCh)

	allDone.Wait()

	select {
	case i := <-claimantCh:
		t.Logf("%d got the ownership!", i)
	default:
		require.FailNow(t, "No claimers were able to claim ownership")
	}
}

func TestCheckpointBlob_IncorrectCaseExists(t *testing.T) {
	testData := newBlobStoreTestData(t)

	orig := azeventhubs.Checkpoint{
		EventHubName:            "eventHub",
		ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
		FullyQualifiedNamespace: "eventhHubsTest.servicebus.windows.net",
		PartitionID:             "0",
		SequenceNumber:          to.Ptr[int64](99),
		Offset:                  to.Ptr[int64](101),
	}

	legacyName := uploadLegacyCheckpointBlob(t, testData.CC, orig)

	checkpoints, err := testData.BlobStore.ListCheckpoints(context.Background(), orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, nil)
	require.NoError(t, err)

	require.Equal(t, []azeventhubs.Checkpoint{
		orig,
	}, checkpoints)

	// let's update the checkpoint - it'll delete the old one and write out a new checkpoint to the correct
	// name.
	err = testData.BlobStore.SetCheckpoint(context.Background(), orig, nil)
	require.NoError(t, err)

	blobs := listAllBlobs(t, testData.CC)
	require.Equal(t, 1, len(blobs), "Old blob deleted, new blob written")
	require.Equal(t, strings.ToLower(legacyName), *blobs[0].Name)

	// if for some reason we attempt to do this again it's fine - nothing breaks since
	// deleting the blob is "best effort".
	err = testData.BlobStore.SetCheckpoint(context.Background(), orig, nil)
	require.NoError(t, err)

	checkpoints, err = testData.BlobStore.ListCheckpoints(context.Background(), orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, nil)
	require.NoError(t, err)

	require.Equal(t, []azeventhubs.Checkpoint{
		{
			EventHubName:            strings.ToLower(orig.EventHubName),
			ConsumerGroup:           strings.ToLower(orig.ConsumerGroup),
			FullyQualifiedNamespace: strings.ToLower(orig.FullyQualifiedNamespace),
			PartitionID:             strings.ToLower(orig.PartitionID),
			// it's critical that when the checkpoint is migrated that we do not lose our
			// current progress.
			SequenceNumber: orig.SequenceNumber,
			Offset:         orig.Offset,
		},
	}, checkpoints)

	// advance one of our checkpoints
	*checkpoints[0].SequenceNumber += 1
	*checkpoints[0].Offset += 1

	err = testData.BlobStore.SetCheckpoint(context.Background(), checkpoints[0], nil)
	require.NoError(t, err)

	checkpointsAfterUpdate, err := testData.BlobStore.ListCheckpoints(context.Background(), orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, nil)
	require.NoError(t, err)

	require.Equal(t, checkpoints, checkpointsAfterUpdate)
}

func listAllBlobs(t *testing.T, cc *container.Client) []*container.BlobItem {
	var blobs []*container.BlobItem

	pager := cc.NewListBlobsFlatPager(nil)

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		blobs = append(blobs, page.Segment.BlobItems...)
	}

	return blobs
}

func assertOwnerships(t *testing.T, expected []azeventhubs.Ownership, actual []azeventhubs.Ownership, checkETag bool) {
	t.Helper()

	expected = sanitizeOwnerships(t, expected, false, false)
	actual = sanitizeOwnerships(t, actual, true, checkETag)

	require.Equal(t, expected, actual)
}

func sanitizeOwnerships(t *testing.T, source []azeventhubs.Ownership, checkTimestamp bool, checkETag bool) []azeventhubs.Ownership {
	t.Helper()

	var dest []azeventhubs.Ownership

	for i, o := range source {
		dest = append(dest, o)

		if checkTimestamp {
			require.NotZero(t, dest[i].LastModifiedTime)
		}

		// clear it out so our asserts work.
		dest[i].LastModifiedTime = time.Time{}

		if checkETag {
			require.NotEmpty(t, dest[i].ETag)
		}

		dest[i].ETag = nil // make our simple .Equal work below.
	}

	return dest
}

// uploadLegacyOwnershipBlob creates a blob without lowercasing the name - this is "incorrect" but what we did for the first
// release.
func uploadLegacyOwnershipBlob(t *testing.T, cc *container.Client, orig azeventhubs.Ownership) string {
	// note that we don't normalize the case in our "legacy" blob names
	legacyName := fmt.Sprintf("%s/%s/%s/ownership/%s", orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, orig.PartitionID)

	blobClient := cc.NewBlockBlobClient(legacyName)

	_, err := blobClient.UploadBuffer(context.Background(), nil, &blockblob.UploadBufferOptions{
		Metadata: map[string]*string{
			"ownerid": &orig.OwnerID,
		},
	})
	require.NoError(t, err)
	return legacyName
}

// uploadLegacyCheckpointBlob creates a blob without lowercasing the name - this is "incorrect" but what we did for the first
// release.
func uploadLegacyCheckpointBlob(t *testing.T, cc *container.Client, orig azeventhubs.Checkpoint) string {
	// note that we don't normalize the case in our "legacy" blob names
	legacyName := fmt.Sprintf("%s/%s/%s/checkpoint/%s", orig.FullyQualifiedNamespace, orig.EventHubName, orig.ConsumerGroup, orig.PartitionID)

	blobClient := cc.NewBlockBlobClient(legacyName)

	_, err := blobClient.UploadBuffer(context.Background(), nil, &blockblob.UploadBufferOptions{
		Metadata: map[string]*string{
			"sequencenumber": to.Ptr(fmt.Sprintf("%d", *orig.SequenceNumber)),
			"offset":         to.Ptr(fmt.Sprintf("%d", *orig.Offset)),
		},
	})
	require.NoError(t, err)
	return legacyName
}

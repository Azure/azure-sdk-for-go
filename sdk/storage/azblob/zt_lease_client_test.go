// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//var headersToIgnoreForLease = []string {"X-Ms-Proposed-Lease-Id", "X-Ms-Lease-Id"}
var proposedLeaseIDs = []*string{to.StringPtr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.StringPtr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func TestContainerAcquireLease(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)
}

func TestContainerDeleteContainerWithoutLeaseId(t *testing.T) {
	_assert := assert.New(t)
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerClient.Delete(ctx, nil)
	_assert.Error(err)

	leaseID := containerLeaseClient.leaseID
	_, err = containerClient.Delete(ctx, &DeleteContainerOptions{
		LeaseAccessConditions: &LeaseAccessConditions{
			LeaseID: leaseID,
		},
	})
	_assert.NoError(err)
}

func TestContainerReleaseLease(t *testing.T) {
	_assert := assert.New(t)
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerClient.Delete(ctx, nil)
	_assert.Error(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)

	_, err = containerClient.Delete(ctx, nil)
	_assert.NoError(err)
}

func TestContainerRenewLease(t *testing.T) {
	_assert := assert.New(t)
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	_assert.NoError(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)
}

func TestContainerChangeLease(t *testing.T) {
	_assert := assert.New(t)
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	changeLeaseResp, err := containerLeaseClient.ChangeLease(ctx, &ChangeLeaseContainerOptions{
		ProposedLeaseID: proposedLeaseIDs[1],
	})
	_assert.NoError(err)
	_assert.EqualValues(changeLeaseResp.LeaseID, proposedLeaseIDs[1])
	_assert.EqualValues(containerLeaseClient.leaseID, proposedLeaseIDs[1])

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	_assert.NoError(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)
}

func TestBlobAcquireLease(t *testing.T) {
	_assert := assert.New(t)

	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(t.Name())
	bbClient := createNewBlockBlob(assert.New(t), blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)
}

func TestDeleteBlobWithoutLeaseId(t *testing.T) {
	_assert := assert.New(t)

	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(t.Name())
	bbClient := createNewBlockBlob(assert.New(t), blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	_assert.Error(err)

	leaseID := blobLeaseClient.leaseID
	_, err = blobLeaseClient.Delete(ctx, &DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			LeaseAccessConditions: &LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_assert.NoError(err)
}

func TestBlobReleaseLease(t *testing.T) {
	_assert := assert.New(t)

	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(t.Name())
	bbClient := createNewBlockBlob(assert.New(t), blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	_assert.Error(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)

	_, err = blobLeaseClient.Delete(ctx, nil)
	_assert.NoError(err)
}

func TestBlobRenewLease(t *testing.T) {
	_assert := assert.New(t)

	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(t.Name())
	bbClient := createNewBlockBlob(assert.New(t), blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	_assert.NoError(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)
}

func TestBlobChangeLease(t *testing.T) {
	_assert := assert.New(t)

	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(t.Name())
	bbClient := createNewBlockBlob(assert.New(t), blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	_assert.NoError(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.Equal(*acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	changeLeaseResp, err := blobLeaseClient.ChangeLease(ctx, &ChangeLeaseBlobOptions{
		ProposedLeaseID: proposedLeaseIDs[1],
	})
	_assert.NoError(err)
	_assert.Equal(*changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	_assert.NoError(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.NoError(err)
}

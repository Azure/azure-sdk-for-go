// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

var proposedLeaseIDs = []*string{to.StringPtr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.StringPtr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func TestContainerAcquireLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
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
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerClient.Delete(ctx, nil)
	require.NotNil(t, err)

	leaseID := containerLeaseClient.leaseID
	_, err = containerClient.Delete(ctx, &DeleteContainerOptions{
		LeaseAccessConditions: &LeaseAccessConditions{
			LeaseID: leaseID,
		},
	})
	require.NoError(t, err)
}

func TestContainerReleaseLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerClient.Delete(ctx, nil)
	require.NotNil(t, err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)

	_, err = containerClient.Delete(ctx, nil)
	require.NoError(t, err)
}

func TestContainerRenewLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	require.NoError(t, err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)
}

func TestContainerChangeLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	changeLeaseResp, err := containerLeaseClient.ChangeLease(ctx, &ChangeLeaseContainerOptions{
		ProposedLeaseID: proposedLeaseIDs[1],
	})
	require.NoError(t, err)
	require.EqualValues(t, changeLeaseResp.LeaseID, proposedLeaseIDs[1])
	require.EqualValues(t, containerLeaseClient.leaseID, proposedLeaseIDs[1])

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	require.NoError(t, err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)
}

func TestBlobAcquireLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)
}

func TestDeleteBlobWithoutLeaseId(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	require.Error(t, err)

	leaseID := blobLeaseClient.leaseID
	_, err = blobLeaseClient.Delete(ctx, &DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			LeaseAccessConditions: &LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	require.NoError(t, err)
}

func TestBlobReleaseLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	require.NotNil(t, err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)

	_, err = blobLeaseClient.Delete(ctx, nil)
	require.NoError(t, err)
}

func TestBlobRenewLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.EqualValues(t, acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	require.NoError(t, err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)
}

func TestBlobChangeLease(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(t, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	require.NoError(t, err)
	require.NotNil(t, acquireLeaseResponse.LeaseID)
	require.Equal(t, *acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	changeLeaseResp, err := blobLeaseClient.ChangeLease(ctx, &ChangeLeaseBlobOptions{
		ProposedLeaseID: proposedLeaseIDs[1],
	})
	require.NoError(t, err)
	require.Equal(t, *changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	require.NoError(t, err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	require.NoError(t, err)
}

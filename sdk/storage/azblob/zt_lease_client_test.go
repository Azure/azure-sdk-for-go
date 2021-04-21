// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	chk "gopkg.in/check.v1"
)

func (s *aztestsSuite) TestContainerAcquireLease(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestContainerDeleteContainerWithoutLeaseId(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)

	_, err = containerClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)

	leaseID := containerLeaseClient.LeaseID
	_, err = containerClient.Delete(ctx, &DeleteContainerOptions{
		LeaseAccessConditions: &LeaseAccessConditions{
			LeaseID: leaseID,
		},
	})
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestContainerReleaseLease(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)

	_, err = containerClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = containerClient.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestContainerRenewLease(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestContainerChangeLease(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)

	proposedLeaseID := newUUID().String()
	changeLeaseResp, err := containerLeaseClient.ChangeLease(ctx, &ChangeLeaseContainerOptions{
		ProposedLeaseId: proposedLeaseID,
	})
	c.Assert(err, chk.IsNil)
	c.Assert(changeLeaseResp.LeaseID, chk.DeepEquals, proposedLeaseID)
	c.Assert(containerLeaseClient.LeaseID, chk.DeepEquals, proposedLeaseID)

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobAcquireLease(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient, _ := createNewBlockBlob(c, containerClient)
	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestDeleteBlobWithoutLeaseId(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient, _ := createNewBlockBlob(c, containerClient)
	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)

	leaseID := blobLeaseClient.LeaseID
	_, err = blobLeaseClient.Delete(ctx, &DeleteBlobOptions{
		LeaseAccessConditions: &LeaseAccessConditions{
			LeaseID: leaseID,
		},
	})
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobReleaseLease(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient, _ := createNewBlockBlob(c, containerClient)
	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = blobLeaseClient.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobRenewLease(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient, _ := createNewBlockBlob(c, containerClient)
	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobChangeLease(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient, _ := createNewBlockBlob(c, containerClient)
	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	c.Assert(err, chk.IsNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)

	proposedLeaseID := to.StringPtr(newUUID().String())
	changeLeaseResp, err := blobLeaseClient.ChangeLease(ctx, &ChangeLeaseBlobOptions{
		ProposedLeaseID: proposedLeaseID,
	})
	c.Assert(err, chk.IsNil)
	c.Assert(changeLeaseResp.LeaseID, chk.DeepEquals, proposedLeaseID)
	c.Assert(blobLeaseClient.LeaseID, chk.DeepEquals, proposedLeaseID)

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	c.Assert(err, chk.IsNil)
}

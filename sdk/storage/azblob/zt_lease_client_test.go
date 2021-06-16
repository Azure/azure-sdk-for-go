// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

//
//import (
//	"context"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/uuid"
//	"github.com/Azure/azure-sdk-for-go/sdk/to"
//	chk "gopkg.in/check.v1"
//)
//
//func (s *azblobTestSuite) TestContainerAcquireLease() {
//	bsu := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)
//
//	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestContainerDeleteContainerWithoutLeaseId() {
//	bsu := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)
//
//	_, err = containerClient.Delete(ctx, nil)
//	_assert.NotNil(err)
//
//	leaseID := containerLeaseClient.LeaseID
//	_, err = containerClient.Delete(ctx, &DeleteContainerOptions{
//		LeaseAccessConditions: &LeaseAccessConditions{
//			LeaseID: leaseID,
//		},
//	})
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestContainerReleaseLease() {
//	bsu := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)
//
//	_, err = containerClient.Delete(ctx, nil)
//	_assert.NotNil(err)
//
//	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//
//	_, err = containerClient.Delete(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestContainerRenewLease() {
//	bsu := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)
//
//	_, err = containerLeaseClient.RenewLease(ctx, nil)
//	_assert.Nil(err)
//
//	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestContainerChangeLease() {
//	bsu := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	containerLeaseClient := bsu.NewContainerLeaseClient(containerName, nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, containerLeaseClient.LeaseID)
//
//	proposedLeaseID := to.StringPtr(uuid.New().String())
//	changeLeaseResp, err := containerLeaseClient.ChangeLease(ctx, &ChangeLeaseContainerOptions{
//		ProposedLeaseID: proposedLeaseID,
//	})
//	_assert.Nil(err)
//	c.Assert(changeLeaseResp.LeaseID, chk.DeepEquals, proposedLeaseID)
//	c.Assert(containerLeaseClient.LeaseID, chk.DeepEquals, proposedLeaseID)
//
//	_, err = containerLeaseClient.RenewLease(ctx, nil)
//	_assert.Nil(err)
//
//	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestBlobAcquireLease() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	bbClient, _ := createNewBlockBlob(c, containerClient)
//	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)
//
//	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestDeleteBlobWithoutLeaseId() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	bbClient, _ := createNewBlockBlob(c, containerClient)
//	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)
//
//	_, err = blobLeaseClient.Delete(ctx, nil)
//	_assert.NotNil(err)
//
//	leaseID := blobLeaseClient.LeaseID
//	_, err = blobLeaseClient.Delete(ctx, &DeleteBlobOptions{
//		LeaseAccessConditions: &LeaseAccessConditions{
//			LeaseID: leaseID,
//		},
//	})
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestBlobReleaseLease() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	bbClient, _ := createNewBlockBlob(c, containerClient)
//	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)
//
//	_, err = blobLeaseClient.Delete(ctx, nil)
//	_assert.NotNil(err)
//
//	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//
//	_, err = blobLeaseClient.Delete(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestBlobRenewLease() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	bbClient, _ := createNewBlockBlob(c, containerClient)
//	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)
//
//	_, err = blobLeaseClient.RenewLease(ctx, nil)
//	_assert.Nil(err)
//
//	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//}
//
//func (s *azblobTestSuite) TestBlobChangeLease() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//
//	bbClient, _ := createNewBlockBlob(c, containerClient)
//	blobLeaseClient := bbClient.NewBlobLeaseClient(nil)
//
//	ctx := context.Background()
//	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
//	_assert.Nil(err)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.NotNil)
//	c.Assert(acquireLeaseResponse.LeaseID, chk.DeepEquals, blobLeaseClient.LeaseID)
//
//	proposedLeaseID := to.StringPtr(uuid.New().String())
//	changeLeaseResp, err := blobLeaseClient.ChangeLease(ctx, &ChangeLeaseBlobOptions{
//		ProposedLeaseID: proposedLeaseID,
//	})
//	_assert.Nil(err)
//	c.Assert(changeLeaseResp.LeaseID, chk.DeepEquals, proposedLeaseID)
//	c.Assert(blobLeaseClient.LeaseID, chk.DeepEquals, proposedLeaseID)
//
//	_, err = blobLeaseClient.RenewLease(ctx, nil)
//	_assert.Nil(err)
//
//	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
//	_assert.Nil(err)
//}

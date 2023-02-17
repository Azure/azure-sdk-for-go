//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running lease Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &LeaseRecordedTestsSuite{})
		suite.Run(t, &LeaseUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &LeaseRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &LeaseRecordedTestsSuite{})
	}
}

func (s *LeaseRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *LeaseRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *LeaseUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *LeaseUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type LeaseRecordedTestsSuite struct {
	suite.Suite
}

type LeaseUnrecordedTestsSuite struct {
	suite.Suite
}

// var headersToIgnoreForLease = []string {"X-Ms-Proposed-Lease-Id", "X-Ms-Lease-Id"}
var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func (s *LeaseRecordedTestsSuite) TestContainerAcquireLease() {
	_require := require.New(s.T())
	testName := s.T().Name()
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	containerLeaseClient, _ := lease.NewContainerClient(containerClient, &lease.ContainerClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *containerLeaseClient.LeaseID())

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestContainerDeleteContainerWithoutLeaseId() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	containerLeaseClient, _ := lease.NewContainerClient(containerClient, &lease.ContainerClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *containerLeaseClient.LeaseID())

	_, err = containerClient.Delete(ctx, nil)
	_require.NotNil(err)

	leaseID := containerLeaseClient.LeaseID()
	_, err = containerClient.Delete(ctx, &container.DeleteOptions{
		AccessConditions: &container.AccessConditions{
			LeaseAccessConditions: &container.LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestContainerReleaseLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	containerLeaseClient, _ := lease.NewContainerClient(containerClient, &lease.ContainerClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *containerLeaseClient.LeaseID())

	_, err = containerClient.Delete(ctx, nil)
	_require.NotNil(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)

	_, err = containerClient.Delete(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestContainerRenewLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	containerLeaseClient, _ := lease.NewContainerClient(containerClient, &lease.ContainerClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *containerLeaseClient.LeaseID())

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	_require.Nil(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestContainerChangeLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	containerLeaseClient, _ := lease.NewContainerClient(containerClient, &lease.ContainerClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(*acquireLeaseResponse.LeaseID, *containerLeaseClient.LeaseID())

	changeLeaseResp, err := containerLeaseClient.ChangeLease(ctx, *proposedLeaseIDs[1], nil)
	_require.Nil(err)
	_require.EqualValues(changeLeaseResp.LeaseID, proposedLeaseIDs[1])
	_require.EqualValues(containerLeaseClient.LeaseID(), proposedLeaseIDs[1])

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	_require.Nil(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestBlobAcquireLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobLeaseClient, err := lease.NewBlobClient(bbClient, &lease.BlobClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.LeaseID())

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestDeleteBlobWithoutLeaseId() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobLeaseClient, err := lease.NewBlobClient(bbClient, &lease.BlobClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.LeaseID())

	_, err = bbClient.Delete(ctx, nil)
	_require.NotNil(err)

	leaseID := blobLeaseClient.LeaseID()
	_, err = bbClient.Delete(ctx, &blob.DeleteOptions{
		AccessConditions: &blob.AccessConditions{
			LeaseAccessConditions: &blob.LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestBlobReleaseLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobLeaseClient, _ := lease.NewBlobClient(bbClient, &lease.BlobClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.LeaseID())

	_, err = bbClient.Delete(ctx, nil)
	_require.NotNil(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)

	_, err = bbClient.Delete(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestBlobRenewLease() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobLeaseClient, _ := lease.NewBlobClient(bbClient, &lease.BlobClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.LeaseID())

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	_require.Nil(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)
}

func (s *LeaseRecordedTestsSuite) TestBlobChangeLease() {
	_require := require.New(s.T())
	testName := s.T().Name()

	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobLeaseClient, _ := lease.NewBlobClient(bbClient, &lease.BlobClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(15), nil)
	_require.Nil(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.Equal(*acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	changeLeaseResp, err := blobLeaseClient.ChangeLease(ctx, *proposedLeaseIDs[1], nil)
	_require.Nil(err)
	_require.Equal(*changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	_require.Nil(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_require.Nil(err)
}

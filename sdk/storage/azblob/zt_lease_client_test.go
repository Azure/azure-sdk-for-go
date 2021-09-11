// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
)

//var headersToIgnoreForLease = []string {"X-Ms-Proposed-Lease-Id", "X-Ms-Lease-Id"}
var proposedLeaseIDs = []*string{to.StringPtr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.StringPtr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func (s *azblobTestSuite) TestContainerAcquireLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestContainerDeleteContainerWithoutLeaseId() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerClient.Delete(ctx, nil)
	_assert.NotNil(err)

	leaseID := containerLeaseClient.leaseID
	_, err = containerClient.Delete(ctx, &DeleteContainerOptions{
		LeaseAccessConditions: &LeaseAccessConditions{
			LeaseID: leaseID,
		},
	})
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestContainerReleaseLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(60)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerClient.Delete(ctx, nil)
	_assert.NotNil(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)

	_, err = containerClient.Delete(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestContainerRenewLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	_assert.Nil(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestContainerChangeLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	containerLeaseClient, _ := containerClient.NewContainerLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := containerLeaseClient.AcquireLease(ctx, &AcquireLeaseContainerOptions{Duration: to.Int32Ptr(15)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, containerLeaseClient.leaseID)

	changeLeaseResp, err := containerLeaseClient.ChangeLease(ctx, &ChangeLeaseContainerOptions{
		ProposedLeaseID: proposedLeaseIDs[1],
	})
	_assert.Nil(err)
	_assert.EqualValues(changeLeaseResp.LeaseID, proposedLeaseIDs[1])
	_assert.EqualValues(containerLeaseClient.leaseID, proposedLeaseIDs[1])

	_, err = containerLeaseClient.RenewLease(ctx, nil)
	_assert.Nil(err)

	_, err = containerLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobAcquireLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestDeleteBlobWithoutLeaseId() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	_assert.NotNil(err)

	leaseID := blobLeaseClient.leaseID
	_, err = blobLeaseClient.Delete(ctx, &DeleteBlobOptions{
		BlobAccessConditions: &BlobAccessConditions{
			LeaseAccessConditions: &LeaseAccessConditions{
				LeaseID: leaseID,
			},
		},
	})
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobReleaseLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(60)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.Delete(ctx, nil)
	_assert.NotNil(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)

	_, err = blobLeaseClient.Delete(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobRenewLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.leaseID)

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	_assert.Nil(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestBlobChangeLease() {
	_assert := assert.New(s.T())
	testName := s.T().Name()

	_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, headersToIgnoreForLease)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	bbClient := createNewBlockBlob(_assert, blobName, containerClient)
	blobLeaseClient, _ := bbClient.NewBlobLeaseClient(proposedLeaseIDs[0])

	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, &AcquireLeaseBlobOptions{Duration: to.Int32Ptr(15)})
	_assert.Nil(err)
	_assert.NotNil(acquireLeaseResponse.LeaseID)
	_assert.Equal(*acquireLeaseResponse.LeaseID, *proposedLeaseIDs[0])

	changeLeaseResp, err := blobLeaseClient.ChangeLease(ctx, &ChangeLeaseBlobOptions{
		ProposedLeaseID: proposedLeaseIDs[1],
	})
	_assert.Nil(err)
	_assert.Equal(*changeLeaseResp.LeaseID, *proposedLeaseIDs[1])

	_, err = blobLeaseClient.RenewLease(ctx, nil)
	_assert.Nil(err)

	_, err = blobLeaseClient.ReleaseLease(ctx, nil)
	_assert.Nil(err)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/stretchr/testify/require"
	"strconv"
	"time"
)

//nolint
func (s *azblobUnrecordedTestSuite) TestSetEmptyAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, &container.SetAccessPolicyOptions{})
	_require.Nil(err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*container.SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		AccessPolicy: &container.AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})

	param := container.SetAccessPolicyOptions{
		ContainerACL: signedIdentifiers,
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	_require.Nil(err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetMultipleAccessPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	id := "empty"

	signedIdentifiers := make([]*container.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id2,
		AccessPolicy: &container.AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id3,
		AccessPolicy: &container.AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})

	param := container.SetAccessPolicyOptions{
		ContainerACL: signedIdentifiers,
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	_require.Nil(err)

	// Make a Get to assert two access policies
	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 3)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetNullAccessPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	id := "null"

	signedIdentifiers := make([]*container.SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &container.SignedIdentifier{
		ID: &id,
	})

	param := container.SetAccessPolicyOptions{
		ContainerACL: signedIdentifiers,
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Equal(len(resp.SignedIdentifiers), 1)
}

func (s *azblobTestSuite) TestContainerGetSetPermissionsMultiplePolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := to.Ptr(container.AccessPolicyPermission{Read: true, Write: true}).String()
	readOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	id1, id2 := "0000", "0001"
	permissions := []*container.SignedIdentifier{
		{ID: &id1,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &readWrite,
			},
		},
		{ID: &id2,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry2,
				Permission: &readOnly,
			},
		},
	}

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)

	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *azblobTestSuite) TestContainerGetPermissionsPublicAccessNotNone() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	access := container.PublicAccessTypeBlob
	createContainerOptions := container.CreateOptions{
		Access: &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	_require.Nil(err)
	defer deleteContainer(_require, containerClient)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)

	_require.Nil(err)
	_require.Equal(*resp.BlobPublicAccess, container.PublicAccessTypeBlob)
}

//func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessNone() {
//	// Test the basic one by making an anonymous request to ensure it's actually doing it and also with GetPermissions
//	// For all the others, can just use GetPermissions since we've validated that it at least registers on the server correctly
//	svcClient := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, svcClient)
//	defer deleteContainer(_require, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	// Container is created with PublicAccessTypeBlob, so setting it to None will actually test that it is changed through this method
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_require.Nil(err)
//
//	_require.Nil(err)
//	bsu2, err := NewServiceClient(svcClient.URL(), azcore.AnonymousCredential(), nil)
//	_require.Nil(err)
//
//	containerClient2 := bsu2.NewContainerClient(containerName)
//	blobURL2 := containerClient2.NewBlockBlobClient(blobName)
//
//	// Get permissions via the original container URL so the request succeeds
//	resp, err := containerClient.GetAccessPolicy(ctx, nil)
//	_assert(resp.BlobPublicAccess, chk.IsNil)
//	_require.Nil(err)
//
//	// If we cannot access a blob's data, we will also not be able to enumerate blobs
//	p := containerClient2.NewListBlobsFlatPager(nil)
//	p.NextPage(ctx)
//	err = p.Err() // grab the next page
//	validateBlobErrorCode(c, err, StorageErrorCodeNoAuthenticationInformation)
//
//	_, err = blobURL2.Download(ctx, nil)
//	validateBlobErrorCode(c, err, StorageErrorCodeNoAuthenticationInformation)
//}

func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessTypeBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: &access,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobPublicAccess, container.PublicAccessTypeBlob)
}

func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessContainer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	access := container.PublicAccessTypeContainer
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access: &access,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobPublicAccess, container.PublicAccessTypeContainer)
}

////// TODO: After Pacer is ready
////func (s *azblobTestSuite) TestContainerSetPermissionsACLSinglePolicy() {
////	svcClient := getServiceClient()
////	credential, err := getGenericCredential("")
////	if err != nil {
////		c.Fatal("Invalid credential")
////	}
////	containerClient, containerName := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	_, blobName := createNewBlockBlob(c, containerClient)
////
////	start := time.Now().UTC().Add(-15 * time.Second)
////	expiry := start.Add(5 * time.Minute).UTC()
////	listOnly := AccessPolicyPermission{List: true}.String()
////	id := "0000"
////	permissions := []container.SignedIdentifier{{
////		ID: &id,
////		AccessPolicy: &container.AccessPolicy{
////			Start:      &start,
////			Expiry:     &expiry,
////			Permission: &listOnly,
////		},
////	}}
////
////	setAccessPolicyOptions := container.SetAccessPolicyOptions{
////		ContainerAcquireLeaseOptions: ContainerAcquireLeaseOptions{
////			ContainerACL: permissions,
////		},
////	}
////	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
////	_require.Nil(err)
////
////	serviceSASValues := BlobSASSignatureValues{Identifier: "0000", ContainerName: containerName}
////	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////
////	sasURL := svcClient.URL()
////	sasURL.RawQuery = queryParams.Encode()
////	sasPipeline := (NewAnonymousCredential(), PipelineOptions{})
////	sasBlobServiceURL := NewServiceURL(sasURL, sasPipeline)
////
////	// Verifies that the SAS can access the resource
////	sasContainer := sasBlobServiceURL.NewContainerClient(containerName)
////	resp, err := sasContainer.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////
////	// Verifies that successful sas access is not just because it's public
////	anonymousBlobService := NewServiceURL(svcClient.URL(), sasPipeline)
////	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
////	_, err = anonymousContainer.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{})
////	validateBlobErrorCode(c, err, StorageErrorCodeNoAuthenticationInformation)
////}

func (s *azblobTestSuite) TestContainerSetPermissionsACLMoreThanFive() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*container.SignedIdentifier, 6)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
}

func (s *azblobTestSuite) TestContainerSetPermissionsDeleteAndModifyACL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	permissions := make([]*container.SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions1)
	_require.Nil(err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, 1)
	_require.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *azblobTestSuite) TestContainerSetPermissionsDeleteAllPolicies() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*container.SignedIdentifier, 2)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Len(resp.SignedIdentifiers, len(permissions))
	_require.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: []*container.SignedIdentifier{},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.SignedIdentifiers)
}

func (s *azblobTestSuite) TestContainerSetPermissionsInvalidPolicyTimes() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.Nil(err)
	permissions := make([]*container.SignedIdentifier, 2)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestContainerSetPermissionsNilPolicySlice() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	_require.Nil(err)
}

func (s *azblobTestSuite) TestContainerSetPermissionsSignedIdentifierTooLong() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	defer deleteContainer(_require, containerClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_require.Nil(err)
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]*container.SignedIdentifier, 2)
	listOnly := to.Ptr(container.AccessPolicyPermission{Read: true}).String()
	for i := 0; i < 2; i++ {
		permissions[i] = &container.SignedIdentifier{
			ID: &id,
			AccessPolicy: &container.AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := container.PublicAccessTypeBlob
	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		Access:       &access,
		ContainerACL: permissions,
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.InvalidXMLDocument)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.BlobPublicAccess)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.BlobPublicAccess)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	//_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := container.SetAccessPolicyOptions{
		AccessConditions: &container.AccessConditions{
			ModifiedAccessConditions: &container.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_require.NotNil(err)

	validateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"time"
)

//nolint
func (s *azblobUnrecordedTestSuite) TestSetEmptyAccessPolicy() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetAccessPolicy(ctx, &SetAccessPolicyOptions{})
	_assert.Nil(err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetAccessPolicy() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expiration := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	permission := "r"
	id := "1"

	signedIdentifiers := make([]*SignedIdentifier, 0)

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		AccessPolicy: &AccessPolicy{
			Expiry:     &expiration,
			Start:      &start,
			Permission: &permission,
		},
		ID: &id,
	})

	param := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: signedIdentifiers,
		},
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	_assert.Nil(err)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetMultipleAccessPolicies() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	id := "empty"

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id,
	})

	permission2 := "r"
	id2 := "partial"

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id2,
		AccessPolicy: &AccessPolicy{
			Permission: &permission2,
		},
	})

	id3 := "full"
	permission3 := "r"
	start := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)
	expiry := time.Date(2021, 6, 8, 2, 10, 9, 0, time.UTC)

	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id3,
		AccessPolicy: &AccessPolicy{
			Start:      &start,
			Expiry:     &expiry,
			Permission: &permission3,
		},
	})

	param := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: signedIdentifiers,
		},
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	_assert.Nil(err)

	// Make a Get to assert two access policies
	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.SignedIdentifiers, 3)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSetNullAccessPolicy() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	id := "null"

	signedIdentifiers := make([]*SignedIdentifier, 0)
	signedIdentifiers = append(signedIdentifiers, &SignedIdentifier{
		ID: &id,
	})

	param := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: signedIdentifiers,
		},
	}

	_, err = containerClient.SetAccessPolicy(ctx, &param)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(len(resp.SignedIdentifiers), 1)
}

func (s *azblobTestSuite) TestContainerGetSetPermissionsMultiplePolicies() {
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

	// Define the policies
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := AccessPolicyPermission{Read: true, Write: true}.String()
	readOnly := AccessPolicyPermission{Read: true}.String()
	id1, id2 := "0000", "0001"
	permissions := []*SignedIdentifier{
		{ID: &id1,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &readWrite,
			},
		},
		{ID: &id2,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry2,
				Permission: &readOnly,
			},
		},
	}

	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)

	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *azblobTestSuite) TestContainerGetPermissionsPublicAccessNotNone() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access: &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	_assert.Nil(err)
	defer deleteContainer(_assert, containerClient)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(*resp.BlobPublicAccess, PublicAccessTypeBlob)
}

//func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessNone() {
//	// Test the basic one by making an anonymous request to ensure it's actually doing it and also with GetPermissions
//	// For all the others, can just use GetPermissions since we've validated that it at least registers on the server correctly
//	svcClient := getServiceClient(nil)
//	containerClient, containerName := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	// Container is created with PublicAccessTypeBlob, so setting it to None will actually test that it is changed through this method
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_assert.Nil(err)
//
//	_assert.Nil(err)
//	bsu2, err := NewServiceClient(svcClient.URL(), azcore.AnonymousCredential(), nil)
//	_assert.Nil(err)
//
//	containerClient2 := bsu2.NewContainerClient(containerName)
//	blobURL2 := containerClient2.NewBlockBlobClient(blobName)
//
//	// Get permissions via the original container URL so the request succeeds
//	resp, err := containerClient.GetAccessPolicy(ctx, nil)
//	_assert(resp.BlobPublicAccess, chk.IsNil)
//	_assert.Nil(err)
//
//	// If we cannot access a blob's data, we will also not be able to enumerate blobs
//	p := containerClient2.ListBlobsFlat(nil)
//	p.NextPage(ctx)
//	err = p.Err() // grab the next page
//	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
//
//	_, err = blobURL2.Download(ctx, nil)
//	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
//}

func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessTypeBlob() {
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

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.BlobPublicAccess, PublicAccessTypeBlob)
}

func (s *azblobTestSuite) TestContainerSetPermissionsPublicAccessContainer() {
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

	access := PublicAccessTypeContainer
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.BlobPublicAccess, PublicAccessTypeContainer)
}

////// TODO: After Pacer is ready
////func (s *azblobTestSuite) TestContainerSetPermissionsACLSinglePolicy() {
////	svcClient := getServiceClient()
////	credential, err := getGenericCredential("")
////	if err != nil {
////		c.Fatal("Invalid credential")
////	}
////	containerClient, containerName := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	_, blobName := createNewBlockBlob(c, containerClient)
////
////	start := time.Now().UTC().Add(-15 * time.Second)
////	expiry := start.Add(5 * time.Minute).UTC()
////	listOnly := AccessPolicyPermission{List: true}.String()
////	id := "0000"
////	permissions := []SignedIdentifier{{
////		ID: &id,
////		AccessPolicy: &AccessPolicy{
////			Start:      &start,
////			Expiry:     &expiry,
////			Permission: &listOnly,
////		},
////	}}
////
////	setAccessPolicyOptions := SetAccessPolicyOptions{
////		ContainerAcquireLeaseOptions: ContainerAcquireLeaseOptions{
////			ContainerACL: permissions,
////		},
////	}
////	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
////	_assert.Nil(err)
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
////	resp, err := sasContainer.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.Nil(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////
////	// Verifies that successful sas access is not just because it's public
////	anonymousBlobService := NewServiceURL(svcClient.URL(), sasPipeline)
////	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
////	_, err = anonymousContainer.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{})
////	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
////}

func (s *azblobTestSuite) TestContainerSetPermissionsACLMoreThanFive() {
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

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)
	permissions := make([]*SignedIdentifier, 6)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

func (s *azblobTestSuite) TestContainerSetPermissionsDeleteAndModifyACL() {
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

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)
	listOnly := AccessPolicyPermission{Read: true}.String()
	permissions := make([]*SignedIdentifier, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp.SignedIdentifiers, permissions)

	permissions = resp.SignedIdentifiers[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions1)
	_assert.Nil(err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.SignedIdentifiers, 1)
	_assert.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *azblobTestSuite) TestContainerSetPermissionsDeleteAllPolicies() {
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

	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)
	permissions := make([]*SignedIdentifier, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.SignedIdentifiers, len(permissions))
	_assert.EqualValues(resp.SignedIdentifiers, permissions)

	setAccessPolicyOptions = SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: []*SignedIdentifier{},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.SignedIdentifiers)
}

func (s *azblobTestSuite) TestContainerSetPermissionsInvalidPolicyTimes() {
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

	// Swap start and expiry
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	start, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)
	permissions := make([]*SignedIdentifier, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestContainerSetPermissionsNilPolicySlice() {
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

	_, err = containerClient.SetAccessPolicy(ctx, nil)
	_assert.Nil(err)
}

func (s *azblobTestSuite) TestContainerSetPermissionsSignedIdentifierTooLong() {
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

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]*SignedIdentifier, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		permissions[i] = &SignedIdentifier{
			ID: &id,
			AccessPolicy: &AccessPolicy{
				Start:      &start,
				Expiry:     &expiry,
				Permission: &listOnly,
			},
		}
	}

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.BlobPublicAccess)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfUnModifiedSinceTrue() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.BlobPublicAccess)
}

func (s *azblobTestSuite) TestContainerSetPermissionsIfUnModifiedSinceFalse() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_assert, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		AccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

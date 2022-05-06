//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
)

//nolint
//func (s *azblobUnrecordedTestSuite) TestNewContainerClientValidName() {
//	_require := require.New(s.T())
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	testURL := svcClient.NewContainerClient(containerPrefix)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_require.Nil(err)
//	correctURL := "https://" + accountName + "." + DefaultBlobEndpointSuffix + containerPrefix
//	_require.Equal(testURL.URL(), correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateRootContainerURL() {
//	_require := require.New(s.T())
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	testURL := svcClient.NewContainerClient(ContainerNameRoot)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	_require.Nil(err)
//	correctURL := "https://" + accountName + ".blob.core.windows.net/$root"
//	_require.Equal(testURL.URL(), correctURL)
//}

func (s *azblobTestSuite) TestContainerCreateInvalidName() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient, _ := svcClient.NewContainerClient("foo bar")

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_require.NotNil(err)
	validateStorageError(_require, err, StorageErrorCodeInvalidResourceName)
}

func (s *azblobTestSuite) TestContainerCreateEmptyName() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient, _ := svcClient.NewContainerClient("")

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *azblobTestSuite) TestContainerCreateNameCollision() {
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

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}

	containerClient, _ = svcClient.NewContainerClient(containerName)
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeContainerAlreadyExists)
}

func (s *azblobTestSuite) TestContainerCreateInvalidMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: map[string]string{"1 foo": "bar"},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)

	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *azblobTestSuite) TestContainerCreateNilMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(_require, containerClient)
	_require.Nil(err)

	response, err := containerClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(response.Metadata)
}

func (s *azblobTestSuite) TestContainerCreateEmptyMetadata() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)

	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(_require, containerClient)
	_require.Nil(err)

	response, err := containerClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(response.Metadata)
}

//func (s *azblobTestSuite) TestContainerCreateAccessContainer() {
//	// TOD0: NotWorking
//	_require := require.New(s.T())
//testName := s.T().Name()
//	_context := getTestContext(testName)
//
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	credential, err := getGenericCredential("")
//	_require.Nil(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := ContainerCreateOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(_require, containerClient)
//	_require.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := BlockBlobUploadOptions{
//		Metadata: basicMetadata,
//	}
//	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	_require.Nil(err)
//
//	// Anonymous enumeration should be valid with container access
//	containerClient2, _ := NewContainerClient(containerClient.URL(), credential, nil)
//	pager := containerClient2.ListBlobsFlat(nil)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			_require.Equal(*blob.Name, blobPrefix)
//		}
//	}
//
//	_require.Nil(pager.Err())
//
//	// Getting blob data anonymously should still be valid with container access
//	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.EqualValues(resp.Metadata, basicMetadata)
//}

//func (s *azblobTestSuite) TestContainerCreateAccessBlob() {
//	// TODO: Not Working
//	_require := require.New(s.T())
// testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := ContainerCreateOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(_require, containerClient)
//	_require.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := BlockBlobUploadOptions{
//		Metadata: basicMetadata,
//	}
//	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	_require.Nil(err)
//
//	// Reference the same container URL but with anonymous credentials
//	containerClient2, err := NewContainerClient(containerClient.URL(), azcore.AnonymousCredential(), nil)
//	_require.Nil(err)
//
//	pager := containerClient2.ListBlobsFlat(nil)
//
//	_require.Equal(pager.NextPage(ctx), false)
//	_require.NotNil(pager.Err())
//
//	// Accessing blob specific data should be public
//	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.EqualValues(resp.Metadata, basicMetadata)
//}

func (s *azblobTestSuite) TestContainerCreateAccessNone() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	// Public Access Type None
	_, err = containerClient.Create(ctx, nil)
	defer deleteContainer(_require, containerClient)
	_require.Nil(err)

	bbClient, _ := containerClient.NewBlockBlobClient(blobPrefix)
	uploadBlockBlobOptions := BlockBlobUploadOptions{
		Metadata: basicMetadata,
	}
	_, err = bbClient.Upload(ctx, internal.NopCloser(strings.NewReader("Content")), &uploadBlockBlobOptions)
	_require.Nil(err)

	// Reference the same container URL but with anonymous credentials
	containerClient2, err := NewContainerClientWithNoCredential(containerClient.URL(), nil)
	_require.Nil(err)

	pager := containerClient2.ListBlobsFlat(nil)

	_require.Equal(pager.NextPage(ctx), false)
	_require.NotNil(pager.Err())

	// Blob data is not public
	blobURL2, _ := containerClient2.NewBlockBlobClient(blobPrefix)
	_, err = blobURL2.GetProperties(ctx, nil)
	_require.NotNil(err)

	//serr := err.(StorageError)
	//_assert(serr.Response().StatusCode, chk.Equals, 401) // HEAD request does not return a status code
}

//func (s *azblobTestSuite) TestContainerCreateIfExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	// Public Access Type None
//	_, err = containerClient.Create(ctx, nil)
//	defer deleteContainer(_require, containerClient)
//	_require.Nil(err)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := ContainerCreateOptions{
//		Access:   &access,
//		Metadata: nil,
//	}
//	_, err = containerClient.CreateIfNotExists(ctx, &createContainerOptions)
//	_require.Nil(err)
//
//	// Ensure that next create call doesn't update the properties of already created container
//	getResp, err := containerClient.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.Nil(getResp.BlobPublicAccess)
//	_require.Nil(getResp.Metadata)
//}
//
//func (s *azblobTestSuite) TestContainerCreateIfNotExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := ContainerCreateOptions{
//		Access:   &access,
//		Metadata: basicMetadata,
//	}
//	_, err = containerClient.CreateIfNotExists(ctx, &createContainerOptions)
//	_require.Nil(err)
//	defer deleteContainer(_require, containerClient)
//
//	// Ensure that next create call doesn't update the properties of already created container
//	getResp, err := containerClient.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_require.EqualValues(*getResp.BlobPublicAccess, PublicAccessTypeBlob)
//	_require.EqualValues(getResp.Metadata, basicMetadata)
//}

func validateContainerDeleted(_require *require.Assertions, containerClient *ContainerClient) {
	_, err := containerClient.GetAccessPolicy(ctx, nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeContainerNotFound)
}

func (s *azblobTestSuite) TestContainerDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)

	_, err = containerClient.Delete(ctx, nil)
	_require.Nil(err)

	validateContainerDeleted(_require, containerClient)
}

//func (s *azblobTestSuite) TestContainerDeleteIfExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	// Public Access Type None
//	_, err = containerClient.Create(ctx, nil)
//	defer deleteContainer(_require, containerClient)
//	_require.Nil(err)
//
//	_, err = containerClient.DeleteIfExists(ctx, nil)
//	_require.Nil(err)
//
//	validateContainerDeleted(_require, containerClient)
//}
//
//func (s *azblobTestSuite) TestContainerDeleteIfNotExists() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	_, err = containerClient.DeleteIfExists(ctx, nil)
//	_require.Nil(err)
//}

func (s *azblobTestSuite) TestContainerDeleteNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	_, err = containerClient.Delete(ctx, nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeContainerNotFound)
}

func (s *azblobTestSuite) TestContainerDeleteIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := ContainerDeleteOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_require.Nil(err)
	validateContainerDeleted(_require, containerClient)
}

func (s *azblobTestSuite) TestContainerDeleteIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := ContainerDeleteOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestContainerDeleteIfUnModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := ContainerDeleteOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_require.Nil(err)

	validateContainerDeleted(_require, containerClient)
}

func (s *azblobTestSuite) TestContainerDeleteIfUnModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	_require.Nil(err)
	_require.Equal(cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(_require, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := ContainerDeleteOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
}

////func (s *azblobTestSuite) TestContainerAccessConditionsUnsupportedConditions() {
////	// This test defines that the library will panic if the user specifies conditional headers
////	// that will be ignored by the service
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////
////	invalidEtag := "invalid"
////	deleteContainerOptions := ContainerSetMetadataOptions{
////		Metadata: basicMetadata,
////		ModifiedAccessConditions: &ModifiedAccessConditions{
////			IfMatch: &invalidEtag,
////		},
////	}
////	_, err := containerClient.SetMetadata(ctx, &deleteContainerOptions)
////	_require.NotNil(err)
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsNonexistentPrefix() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	prefix := blobPrefix + blobPrefix
////	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
////		Prefix: &prefix,
////	}
////	listResponse, errChan := containerClient.ListBlobsFlat(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(listResponse, chk.IsNil)
////}
//
//func (s *azblobTestSuite) TestContainerListBlobsSpecificValidPrefix() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_require, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	prefix := blobPrefix
//	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
//		Prefix: &prefix,
//	}
//	pager := containerClient.ListBlobsFlat(&containerListBlobFlatSegmentOptions)
//
//	count := 0
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			count++
//			_assert(*blob.Name, chk.Equals, blobName)
//		}
//	}
//
//	_assert(pager.Err(), chk.IsNil)
//
//	_assert(count, chk.Equals, 1)
//}
//
//func (s *azblobTestSuite) TestContainerListBlobsValidDelimiter() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_require, containerClient)
//	prefixes := []string{"a/1", "a/2", "b/2", "blob"}
//	blobNames := make([]string, 4)
//	for idx, prefix := range prefixes {
//		_, blobNames[idx] = createNewBlockBlobWithPrefix(c, containerClient, prefix)
//	}
//
//	pager := containerClient.ListBlobsHierarchy("/", nil)
//
//	count := 0
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			count++
//			_assert(*blob.Name, chk.Equals, blobNames[3])
//		}
//	}
//
//	_assert(pager.Err(), chk.IsNil)
//	_assert(count, chk.Equals, 1)
//
//	// TODO: Ask why the output is BlobItemInternal and why other fields are not there for ex: prefix array
//	//_require.Nil(err)
//	//_assert(len(resp.Segment.BlobItems), chk.Equals, 1)
//	//_assert(len(resp.Segment.BlobPrefixes), chk.Equals, 2)
//	//_assert(resp.Segment.BlobPrefixes[0].Name, chk.Equals, "a/")
//	//_assert(resp.Segment.BlobPrefixes[1].Name, chk.Equals, "b/")
//	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

func (s *azblobTestSuite) TestContainerListBlobsWithSnapshots() {
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

	// initialize a blob and create a snapshot of it
	snapBlobName := generateBlobName(testName)
	snapBlob := createNewBlockBlob(_require, snapBlobName, containerClient)
	snap, err := snapBlob.CreateSnapshot(ctx, nil)
	// snap.
	_require.Nil(err)

	listBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots},
	}
	pager := containerClient.ListBlobsFlat(&listBlobFlatSegmentOptions)

	wasFound := false // hold the for loop accountable for finding the blob and it's snapshot
	for pager.NextPage(ctx) {
		_require.Nil(pager.Err())

		resp := pager.PageResponse()

		for _, blob := range resp.Segment.BlobItems {
			if *blob.Name == snapBlobName && blob.Snapshot != nil {
				wasFound = true
				_require.Equal(*blob.Snapshot, *snap.Snapshot)
			}
		}
	}
	_require.Equal(wasFound, true)
}

func (s *azblobTestSuite) TestContainerListBlobsInvalidDelimiter() {
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
	prefixes := []string{"a/1", "a/2", "b/1", "blob"}
	for _, prefix := range prefixes {
		blobName := prefix + generateBlobName(testName)
		createNewBlockBlob(_require, blobName, containerClient)
	}

	pager := containerClient.ListBlobsHierarchy("^", nil)

	pager.NextPage(ctx)
	_require.Nil(pager.Err())
	_require.Nil(pager.PageResponse().Segment.BlobPrefixes)
}

////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeMetadata() {
////	svcClient := getServiceClient()
////	container, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(container)
////	_, blobNameNoMetadata := createNewBlockBlobWithPrefix(c, container, "a")
////	blobMetadata, blobNameMetadata := createNewBlockBlobWithPrefix(c, container, "b")
////	_, err := blobMetadata.SetMetadata(ctx, Metadata{"field": "value"}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////
////	resp, err := container.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Metadata: true}})
////
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobNameNoMetadata)
////	_assert(resp.Segment.BlobItems[0].Metadata, chk.HasLen, 0)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobNameMetadata)
////	_assert(resp.Segment.BlobItems[1].Metadata["field"], chk.Equals, "value")
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeSnapshots() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	blob, blobName := createNewBlockBlob(c, containerClient)
////	_, err := blob.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
////
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 2)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[0].Snapshot, chk.NotNil)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[1].Snapshot, chk.Equals, "")
////}
////
////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeCopy() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	bbClient, blobName := createNewBlockBlob(c, containerClient)
////	blobCopyURL, blobCopyName := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	_, err := blobCopyURL.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	_require.Nil(err)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Copy: true}})
////
////	// These are sufficient to show that the blob copy was in fact included
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 2)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobCopyName)
////	_assert(*resp.Segment.BlobItems[0].Properties.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
////	temp := bbClient.URL()
////	_assert(*resp.Segment.BlobItems[0].Properties.CopySource, chk.Equals, temp.String())
////	_assert(resp.Segment.BlobItems[0].Properties.CopyStatus, chk.Equals, CopyStatusTypeSuccess)
////}
////
////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeUncommitted() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	bbClient, blobName := getBlockBlobURL(c, containerClient)
////	_, err := bbClient.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), LeaseAccessConditions{}, nil, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{UncommittedBlobs: true}})
////
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////}
//
////func testContainerListBlobsIncludeTypeDeletedImpl(, svcClient ServiceURL) error {
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_require.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
////	_require.Nil(err)
////
////	resp, err = containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_require.Nil(err)
////	if len(resp.Segment.BlobItems) != 1 {
////		return errors.New("DeletedBlobNotFound")
////	}
////
////	// resp.Segment.BlobItems[0].Deleted == true/false if versioning is disabled/enabled.
////	_assert(resp.Segment.BlobItems[0].Deleted, chk.Equals, false)
////	return nil
////}
////
////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeDeleted() {
////	svcClient := getServiceClient()
////
////	runTestRequiringServiceProperties(c, svcClient, "DeletedBlobNotFound", enableSoftDelete,
////		testContainerListBlobsIncludeTypeDeletedImpl, disableSoftDelete)
////}
////
////func testContainerListBlobsIncludeMultipleImpl(, svcClient ServiceURL) error {
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////
////	bbClient, _ := createNewBlockBlobWithPrefix(c, containerClient, "z")
////	_, err := bbClient.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_require.Nil(err)
////	blobURL2, _ := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	resp2, err := blobURL2.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	_require.Nil(err)
////	waitForCopy(c, blobURL2, resp2)
////	blobURL3, _ := createNewBlockBlobWithPrefix(c, containerClient, "deleted")
////
////	_, err = blobURL3.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Copy: true, Deleted: true, Versions: true}})
////
////	_require.Nil(err)
////	if len(resp.Segment.BlobItems) != 6 {
////		// If there are fewer blobs in the container than there should be, it will be because one was permanently deleted.
////		return errors.New("DeletedBlobNotFound")
////	}
////
////	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName2)
////	//_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName) // With soft delete, the overwritten blob will have a backup snapshot
////	//_assert(resp.Segment.BlobItems[2].Name, chk.Equals, blobName)
////	return nil
////}
////
////func (s *azblobTestSuite) TestContainerListBlobsIncludeMultiple() {
////	svcClient := getServiceClient()
////
////	runTestRequiringServiceProperties(c, svcClient, "DeletedBlobNotFound", enableSoftDelete,
////		testContainerListBlobsIncludeMultipleImpl, disableSoftDelete)
////}
////
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsNegative() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////
////	defer deleteContainer(_require, containerClient)
////	_, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: -2})
////	_assert(err, chk.Not(chk.IsNil))
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsZero() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	maxResults := int32(0)
////	resp, errChan := containerClient.ListBlobsFlat(ctx, 1, 0, &ContainerListBlobsFlatOptions{MaxResults: &maxResults})
////
////	_assert(<-errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////}
//
//// TODO: Adele: Case failing
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsInsufficient() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_require, containerClient)
////	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
////	createNewBlockBlobWithPrefix(c, containerClient, "b")
////
////	maxResults := int32(1)
////	resp, errChan := containerClient.ListBlobsFlat(ctx, 3, 0, &ContainerListBlobsFlatOptions{MaxResults: &maxResults})
////	_assert(<- errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////	_assert((<- resp).Name, chk.Equals, blobName)
////}

func (s *azblobTestSuite) TestContainerListBlobsMaxResultsExact() {
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
	blobNames := make([]string, 2)
	blobName := generateBlobName(testName)
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	createNewBlockBlob(_require, blobNames[0], containerClient)
	createNewBlockBlob(_require, blobNames[1], containerClient)

	maxResult := int32(2)
	pager := containerClient.ListBlobsFlat(&ContainerListBlobsFlatOptions{
		MaxResults: &maxResult,
	})

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.Segment.BlobItems {
			_require.Equal(nameMap[*blob.Name], true)
		}
	}

	_require.Nil(pager.Err())
}

func (s *azblobTestSuite) TestContainerListBlobsMaxResultsSufficient() {
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

	blobNames := make([]string, 2)
	blobName := generateBlobName(testName)
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	createNewBlockBlob(_require, blobNames[0], containerClient)
	createNewBlockBlob(_require, blobNames[1], containerClient)

	maxResult := int32(3)
	containerListBlobFlatSegmentOptions := ContainerListBlobsFlatOptions{
		MaxResults: &maxResult,
	}
	pager := containerClient.ListBlobsFlat(&containerListBlobFlatSegmentOptions)

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.Segment.BlobItems {
			_require.Equal(nameMap[*blob.Name], true)
		}
	}

	_require.Nil(pager.Err())
}

func (s *azblobTestSuite) TestContainerListBlobsNonExistentContainer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	pager := containerClient.ListBlobsFlat(nil)

	pager.NextPage(ctx)
	_require.NotNil(pager.Err())
}

func (s *azblobTestSuite) TestContainerGetPropertiesAndMetadataNoMetadata() {
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

	resp, err := containerClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestContainerGetPropsAndMetaNonExistentContainer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	_, err = containerClient.GetProperties(ctx, nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeContainerNotFound)
}

func (s *azblobTestSuite) TestContainerSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Metadata: basicMetadata,
		Access:   &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(_require, containerClient)
	_require.Nil(err)

	setMetadataContainerOptions := ContainerSetMetadataOptions{
		Metadata: map[string]string{},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	_require.Nil(err)

	resp, err := containerClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestContainerSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)
	access := PublicAccessTypeBlob
	createContainerOptions := ContainerCreateOptions{
		Access:   &access,
		Metadata: basicMetadata,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_require.Nil(err)
	defer deleteContainer(_require, containerClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	_require.Nil(err)

	resp, err := containerClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestContainerSetMetadataInvalidField() {
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

	setMetadataContainerOptions := ContainerSetMetadataOptions{
		Metadata: map[string]string{"!nval!d Field!@#%": "value"},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	_require.NotNil(err)
	_require.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *azblobTestSuite) TestContainerSetMetadataNonExistent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	_require.NotNil(err)

	validateStorageError(_require, err, StorageErrorCodeContainerNotFound)
}

//
//func (s *azblobTestSuite) TestContainerSetMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//
//	defer deleteContainer(_require, containerClient)
//
//	setMetadataContainerOptions := ContainerSetMetadataOptions{
//		Metadata: basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	_require.Nil(err)
//
//	resp, err := containerClient.GetProperties(ctx, nil)
//	_require.Nil(err)
//	_assert(resp.Metadata, chk.NotNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//
//}

//func (s *azblobTestSuite) TestContainerSetMetadataIfModifiedSinceFalse() {
//	// TODO: NotWorking
//	_require := require.New(s.T())
// testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	containerClient, _ := createNewContainer(_require, testName, svcClient)
//
//	defer deleteContainer(_require, containerClient)
//
//	//currentTime := getRelativeTimeGMT(10)
//	//currentTime, err := time.Parse(time.UnixDate, "Wed Jan 07 11:11:11 PST 2099")
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	_require.Nil(err)
//	setMetadataContainerOptions := ContainerSetMetadataOptions{
//		Metadata: basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	_require.NotNil(err)
//
//	validateStorageError(_require, err, StorageErrorCodeConditionNotMet)
//}

func (s *azblobTestSuite) TestContainerNewBlobURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	bbClient, _ := containerClient.NewBlobClient(blobPrefix)

	_require.Equal(bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	_require.IsTypef(bbClient, &BlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlobClient{}))
}

func (s *azblobTestSuite) TestContainerNewBlockBlobClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient, _ := getContainerClient(containerName, svcClient)

	bbClient, _ := containerClient.NewBlockBlobClient(blobPrefix)

	_require.Equal(bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	_require.IsTypef(bbClient, &BlockBlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlockBlobClient{}))
}

func (s *azblobTestSuite) TestListBlobIncludeMetadata() {
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

	blobName := generateBlobName(testName)
	for i := 0; i < 6; i++ {
		bbClient, _ := getBlockBlobClient(blobName+strconv.Itoa(i), containerClient)
		cResp, err := bbClient.Upload(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &BlockBlobUploadOptions{Metadata: basicMetadata})
		_require.Nil(err)
		_require.Equal(cResp.RawResponse.StatusCode, 201)
	}

	pager := containerClient.ListBlobsFlat(&ContainerListBlobsFlatOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemMetadata},
	})

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		_require.Len(resp.ListBlobsFlatSegmentResponse.Segment.BlobItems, 6)
		for _, blob := range resp.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			_require.NotNil(blob.Metadata)
			_require.Len(blob.Metadata, len(basicMetadata))
		}
	}
	_require.Nil(pager.Err())

	//----------------------------------------------------------

	pager1 := containerClient.ListBlobsHierarchy("/", &ContainerListBlobsHierarchyOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemMetadata, ListBlobsIncludeItemTags},
	})

	for pager1.NextPage(ctx) {
		resp := pager1.PageResponse()
		_require.Len(resp.ListBlobsHierarchySegmentResponse.Segment.BlobItems, 6)
		for _, blob := range resp.ListBlobsHierarchySegmentResponse.Segment.BlobItems {
			_require.NotNil(blob.Metadata)
			_require.Len(blob.Metadata, len(basicMetadata))
		}
	}
	_require.Nil(pager1.Err())
}

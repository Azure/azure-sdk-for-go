// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

//nolint
//func (s *azblobUnrecordedTestSuite) TestNewContainerClientValidName() {
//	_assert := assert.New(s.T())
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	testURL := svcClient.NewContainerClient(containerPrefix)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	require.NoError(t, err)
//	correctURL := "https://" + accountName + "." + DefaultBlobEndpointSuffix + containerPrefix
//	require.Equal(t, testURL.URL(), correctURL)
//}

//nolint
//func (s *azblobUnrecordedTestSuite) TestCreateRootContainerURL() {
//	_assert := assert.New(s.T())
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	testURL := svcClient.NewContainerClient(ContainerNameRoot)
//
//	accountName, err := getRequiredEnv(AccountNameEnvVar)
//	require.NoError(t, err)
//	correctURL := "https://" + accountName + ".blob.core.windows.net/$root"
//	require.Equal(t, testURL.URL(), correctURL)
//}

func TestContainerCreateInvalidName(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := svcClient.NewContainerClient("foo bar")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	require.Error(t, err)
	validateStorageError(t, err, StorageErrorCodeInvalidResourceName)
}

func TestContainerCreateEmptyName(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := svcClient.NewContainerClient("")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidQueryParameterValue)
}

func TestContainerCreateNameCollision(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}

	containerClient = svcClient.NewContainerClient(containerName)
	_, err = containerClient.Create(ctx, &createContainerOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeContainerAlreadyExists)
}

func TestContainerCreateInvalidMetadata(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{"1 foo": "bar"},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)

	require.Error(t, err)
	require.Equal(t, strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func TestContainerCreateNilMetadata(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(t, containerClient)
	require.NoError(t, err)

	response, err := containerClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, response.Metadata)
}

func TestContainerCreateEmptyMetadata(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(t, containerClient)
	require.NoError(t, err)

	response, err := containerClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, response.Metadata)
}

//func (s *azblobTestSuite) TestContainerCreateAccessContainer() {
//	// TOD0: NotWorking
//	_assert := assert.New(s.T())
//// testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	credential, err := getGenericCredential(t)
//	require.NoError(t, err)
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, svcClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := CreateContainerOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(t, containerClient)
//	require.NoError(t, err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	}
//	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	require.NoError(t, err)
//
//	// Anonymous enumeration should be valid with container access
//	containerClient2, _ := NewContainerClient(containerClient.URL(), credential, nil)
//	pager := containerClient2.ListBlobsFlat(nil)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			require.Equal(t, *blob.Name, blobPrefix)
//		}
//	}
//
//	require.Nil(t, pager.Err())
//
//	// Getting blob data anonymously should still be valid with container access
//	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, nil)
//	require.NoError(t, err)
//	_assert.EqualValues(resp.Metadata, basicMetadata)
//}

//func (s *azblobTestSuite) TestContainerCreateAccessBlob() {
//	// TODO: Not Working
//	_assert := assert.New(s.T())
// // testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, svcClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := CreateContainerOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(t, containerClient)
//	require.NoError(t, err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	}
//	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	require.NoError(t, err)
//
//	// Reference the same container URL but with anonymous credentials
//	containerClient2, err := NewContainerClient(containerClient.URL(), azcore.AnonymousCredential(), nil)
//	require.NoError(t, err)
//
//	pager := containerClient2.ListBlobsFlat(nil)
//
//	require.Equal(t, pager.NextPage(ctx), false)
//	require.NotNil(t, pager.Err())
//
//	// Accessing blob specific data should be public
//	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, nil)
//	require.NoError(t, err)
//	_assert.EqualValues(resp.Metadata, basicMetadata)
//}

func TestContainerCreateAccessNone(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	// Public Access Type None
	_, err = containerClient.Create(ctx, nil)
	defer deleteContainer(t, containerClient)
	require.NoError(t, err)

	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: basicMetadata,
	}
	_, err = bbClient.Upload(ctx, internal.NopCloser(strings.NewReader("Content")), &uploadBlockBlobOptions)
	require.NoError(t, err)

	// Reference the same container URL but with anonymous credentials
	containerClient2, err := NewContainerClientWithNoCredential(containerClient.URL(), nil)
	require.NoError(t, err)

	pager := containerClient2.ListBlobsFlat(nil)

	require.Equal(t, pager.NextPage(ctx), false)
	require.NotNil(t, pager.Err())

	// Blob data is not public
	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
	_, err = blobURL2.GetProperties(ctx, nil)
	require.Error(t, err)

	//serr := err.(StorageError)
	//_assert(serr.Response().StatusCode, chk.Equals, 401) // HEAD request does not return a status code
}

//func (s *azblobTestSuite) TestContainerCreateIfExists() {
//	_assert := assert.New(s.T())
//	// testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	// Public Access Type None
//	_, err = containerClient.Create(ctx, nil)
//	defer deleteContainer(t, containerClient)
//	require.NoError(t, err)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := CreateContainerOptions{
//		Access:   &access,
//		Metadata: nil,
//	}
//	_, err = containerClient.CreateIfNotExists(ctx, &createContainerOptions)
//	require.NoError(t, err)
//
//	// Ensure that next create call doesn't update the properties of already created container
//	getResp, err := containerClient.GetProperties(ctx, nil)
//	require.NoError(t, err)
//	require.Nil(t, getResp.BlobPublicAccess)
//	require.Nil(t, getResp.Metadata)
//}
//
//func (s *azblobTestSuite) TestContainerCreateIfNotExists() {
//	_assert := assert.New(s.T())
//	// testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := CreateContainerOptions{
//		Access:   &access,
//		Metadata: basicMetadata,
//	}
//	_, err = containerClient.CreateIfNotExists(ctx, &createContainerOptions)
//	require.NoError(t, err)
//	defer deleteContainer(t, containerClient)
//
//	// Ensure that next create call doesn't update the properties of already created container
//	getResp, err := containerClient.GetProperties(ctx, nil)
//	require.NoError(t, err)
//	_assert.EqualValues(*getResp.BlobPublicAccess, PublicAccessTypeBlob)
//	_assert.EqualValues(getResp.Metadata, basicMetadata)
//}

func validateContainerDeleted(t *testing.T, containerClient ContainerClient) {
	_, err := containerClient.GetAccessPolicy(ctx, nil)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeContainerNotFound)
}

func TestContainerDelete(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	_, err = containerClient.Delete(ctx, nil)
	require.NoError(t, err)

	validateContainerDeleted(t, containerClient)
}

//func (s *azblobTestSuite) TestContainerDeleteIfExists() {
//	_assert := assert.New(s.T())
//	// testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	// Public Access Type None
//	_, err = containerClient.Create(ctx, nil)
//	defer deleteContainer(t, containerClient)
//	require.NoError(t, err)
//
//	_, err = containerClient.DeleteIfExists(ctx, nil)
//	require.NoError(t, err)
//
//	validateContainerDeleted(_assert, containerClient)
//}
//
//func (s *azblobTestSuite) TestContainerDeleteIfNotExists() {
//	_assert := assert.New(s.T())
//	// testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, serviceClient)
//
//	_, err = containerClient.DeleteIfExists(ctx, nil)
//	require.NoError(t, err)
//}

func TestContainerDeleteNonExistent(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Delete(ctx, nil)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeContainerNotFound)
}

func TestContainerDeleteIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	require.NoError(t, err)
	validateContainerDeleted(t, containerClient)
}

func TestContainerDeleteIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

func TestContainerDeleteIfUnModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	require.NoError(t, err)

	validateContainerDeleted(t, containerClient)
}

func TestContainerDeleteIfUnModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	cResp, err := containerClient.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp.RawResponse.StatusCode, 201)
	defer deleteContainer(t, containerClient)

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
}

////func (s *azblobTestSuite) TestContainerAccessConditionsUnsupportedConditions() {
////	// This test defines that the library will panic if the user specifies conditional headers
////	// that will be ignored by the service
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(t, containerClient)
////
////	invalidEtag := "invalid"
////	deleteContainerOptions := SetMetadataContainerOptions{
////		Metadata: basicMetadata,
////		ModifiedAccessConditions: &ModifiedAccessConditions{
////			IfMatch: &invalidEtag,
////		},
////	}
////	_, err := containerClient.SetMetadata(ctx, &deleteContainerOptions)
////	require.Error(t, err)
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsNonexistentPrefix() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(t, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	prefix := blobPrefix + blobPrefix
////	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
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
//	defer deleteContainer(t, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	prefix := blobPrefix
//	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
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
//	defer deleteContainer(t, containerClient)
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
//	//require.NoError(t, err)
//	//_assert(len(resp.Segment.BlobItems), chk.Equals, 1)
//	//_assert(len(resp.Segment.BlobPrefixes), chk.Equals, 2)
//	//_assert(resp.Segment.BlobPrefixes[0].Name, chk.Equals, "a/")
//	//_assert(resp.Segment.BlobPrefixes[1].Name, chk.Equals, "b/")
//	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

func TestContainerListBlobsWithSnapshots(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	// initialize a blob and create a snapshot of it
	snapBlobName := generateBlobName(t.Name())
	snapBlob := createNewBlockBlob(t, snapBlobName, containerClient)
	snap, err := snapBlob.CreateSnapshot(ctx, nil)
	// snap.
	require.NoError(t, err)

	listBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots},
	}
	pager := containerClient.ListBlobsFlat(&listBlobFlatSegmentOptions)

	wasFound := false // hold the for loop accountable for finding the blob and it's snapshot
	for pager.NextPage(ctx) {
		require.Nil(t, pager.Err())

		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			if *blob.Name == snapBlobName && blob.Snapshot != nil {
				wasFound = true
				require.Equal(t, *blob.Snapshot, *snap.Snapshot)
			}
		}
	}
	require.Equal(t, wasFound, true)
}

func TestContainerListBlobsInvalidDelimiter(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)
	prefixes := []string{"a/1", "a/2", "b/1", "blob"}
	for _, prefix := range prefixes {
		blobName := prefix + generateBlobName(t.Name())
		createNewBlockBlob(t, blobName, containerClient)
	}

	pager := containerClient.ListBlobsHierarchy("^", nil)

	pager.NextPage(ctx)
	require.Nil(t, pager.Err())
	require.Nil(t, pager.PageResponse().ContainerListBlobHierarchySegmentResult.Segment.BlobPrefixes)
}

////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeMetadata() {
////	svcClient := getServiceClient()
////	container, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(container)
////	_, blobNameNoMetadata := createNewBlockBlobWithPrefix(c, container, "a")
////	blobMetadata, blobNameMetadata := createNewBlockBlobWithPrefix(c, container, "b")
////	_, err := blobMetadata.SetMetadata(ctx, Metadata{"field": "value"}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	require.NoError(t, err)
////
////	resp, err := container.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Metadata: true}})
////
////	require.NoError(t, err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobNameNoMetadata)
////	_assert(resp.Segment.BlobItems[0].Metadata, chk.HasLen, 0)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobNameMetadata)
////	_assert(resp.Segment.BlobItems[1].Metadata["field"], chk.Equals, "value")
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeSnapshots() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(t, containerClient)
////	blob, blobName := createNewBlockBlob(c, containerClient)
////	_, err := blob.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	require.NoError(t, err)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
////
////	require.NoError(t, err)
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
////	defer deleteContainer(t, containerClient)
////	bbClient, blobName := createNewBlockBlob(c, containerClient)
////	blobCopyURL, blobCopyName := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	_, err := blobCopyURL.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	require.NoError(t, err)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Copy: true}})
////
////	// These are sufficient to show that the blob copy was in fact included
////	require.NoError(t, err)
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
////	defer deleteContainer(t, containerClient)
////	bbClient, blobName := getBlockBlobURL(c, containerClient)
////	_, err := bbClient.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), LeaseAccessConditions{}, nil, ClientProvidedKeyOptions{})
////	require.NoError(t, err)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{UncommittedBlobs: true}})
////
////	require.NoError(t, err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////}
//
////func testContainerListBlobsIncludeTypeDeletedImpl(, svcClient ServiceURL) error {
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(t, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	require.NoError(t, err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
////	require.NoError(t, err)
////
////	resp, err = containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	require.NoError(t, err)
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
////	defer deleteContainer(t, containerClient)
////
////	bbClient, _ := createNewBlockBlobWithPrefix(c, containerClient, "z")
////	_, err := bbClient.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	require.NoError(t, err)
////	blobURL2, _ := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	resp2, err := blobURL2.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	require.NoError(t, err)
////	waitForCopy(c, blobURL2, resp2)
////	blobURL3, _ := createNewBlockBlobWithPrefix(c, containerClient, "deleted")
////
////	_, err = blobURL3.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
////
////	resp, err := containerClient.ListBlobsFlat(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Copy: true, Deleted: true, Versions: true}})
////
////	require.NoError(t, err)
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
////	defer deleteContainer(t, containerClient)
////	_, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: -2})
////	_assert(err, chk.Not(chk.IsNil))
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsZero() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(t, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	maxResults := int32(0)
////	resp, errChan := containerClient.ListBlobsFlat(ctx, 1, 0, &ContainerListBlobFlatSegmentOptions{Maxresults: &maxResults})
////
////	_assert(<-errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////}
//
//// TODO: Adele: Case failing
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsInsufficient() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(t, containerClient)
////	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
////	createNewBlockBlobWithPrefix(c, containerClient, "b")
////
////	maxResults := int32(1)
////	resp, errChan := containerClient.ListBlobsFlat(ctx, 3, 0, &ContainerListBlobFlatSegmentOptions{Maxresults: &maxResults})
////	_assert(<- errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////	_assert((<- resp).Name, chk.Equals, blobName)
////}

func TestContainerListBlobsMaxResultsExact(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)
	blobNames := make([]string, 2)
	blobName := generateBlobName(t.Name())
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	createNewBlockBlob(t, blobNames[0], containerClient)
	createNewBlockBlob(t, blobNames[1], containerClient)

	maxResult := int32(2)
	pager := containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
		Maxresults: &maxResult,
	})

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			require.Equal(t, nameMap[*blob.Name], true)
		}
	}

	require.Nil(t, pager.Err())
}

func TestContainerListBlobsMaxResultsSufficient(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	blobNames := make([]string, 2)
	blobName := generateBlobName(t.Name())
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	createNewBlockBlob(t, blobNames[0], containerClient)
	createNewBlockBlob(t, blobNames[1], containerClient)

	maxResult := int32(3)
	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
		Maxresults: &maxResult,
	}
	pager := containerClient.ListBlobsFlat(&containerListBlobFlatSegmentOptions)

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			require.Equal(t, nameMap[*blob.Name], true)
		}
	}

	require.Nil(t, pager.Err())
}

func TestContainerListBlobsNonExistentContainer(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	pager := containerClient.ListBlobsFlat(nil)

	pager.NextPage(ctx)
	require.NotNil(t, pager.Err())
}

func TestContainerGetPropertiesAndMetadataNoMetadata(t *testing.T) {
	t.Skip("Error: 'System.InvalidCastException: Unable to cast object of type 'System.Net.Http.EmptyReadStream' to type 'System.IO.MemoryStream'.'")
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	resp, err := containerClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

func TestContainerGetPropsAndMetaNonExistentContainer(t *testing.T) {
	t.Skip("Error: 'System.InvalidCastException: Unable to cast object of type 'System.Net.Http.EmptyReadStream' to type 'System.IO.MemoryStream'.'")
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.GetProperties(ctx, nil)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeContainerNotFound)
}

func TestContainerSetMetadataEmpty(t *testing.T) {
	t.Skip("Error: 'System.InvalidCastException: Unable to cast object of type 'System.Net.Http.EmptyReadStream' to type 'System.IO.MemoryStream'.'")
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Metadata: basicMetadata,
		Access:   &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(t, containerClient)
	require.NoError(t, err)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: map[string]string{},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	require.NoError(t, err)

	resp, err := containerClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

func TestContainerSetMetadataNil(t *testing.T) {
	t.Skip("Error: 'System.InvalidCastException: Unable to cast object of type 'System.Net.Http.EmptyReadStream' to type 'System.IO.MemoryStream'.'")
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)
	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: basicMetadata,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	require.NoError(t, err)
	defer deleteContainer(t, containerClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	require.NoError(t, err)

	resp, err := containerClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

func TestContainerSetMetadataInvalidField(t *testing.T) {
	t.Skip("Error: 'System.InvalidCastException: Unable to cast object of type 'System.Net.Http.EmptyReadStream' to type 'System.IO.MemoryStream'.'")
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := createNewContainer(t, containerName, svcClient)

	defer deleteContainer(t, containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: map[string]string{"!nval!d Field!@#%": "value"},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	require.Error(t, err)
	require.Equal(t, strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func TestContainerSetMetadataNonExistent(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeContainerNotFound)
}

//
//func (s *azblobTestSuite) TestContainerSetMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//
//	defer deleteContainer(t, containerClient)
//
//	setMetadataContainerOptions := SetMetadataContainerOptions{
//		Metadata: basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	require.NoError(t, err)
//
//	resp, err := containerClient.GetProperties(ctx, nil)
//	require.NoError(t, err)
//	_assert(resp.Metadata, chk.NotNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//
//}

//func (s *azblobTestSuite) TestContainerSetMetadataIfModifiedSinceFalse() {
//	// TODO: NotWorking
//	_assert := assert.New(s.T())
// // testName := s.T().Name()
//	_context := getTestContext(s.T().Name())
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	containerClient, _ := createNewContainer(t, testName, svcClient)
//
//	defer deleteContainer(t, containerClient)
//
//	//currentTime := getRelativeTimeGMT(10)
//	//currentTime, err := time.Parse(time.UnixDate, "Wed Jan 07 11:11:11 PST 2099")
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	require.NoError(t, err)
//	setMetadataContainerOptions := SetMetadataContainerOptions{
//		Metadata: basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	require.Error(t, err)
//
//	validateStorageError(t, err, StorageErrorCodeConditionNotMet)
//}

func TestContainerNewBlobURL(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	bbClient := containerClient.NewBlobClient(blobPrefix)

	require.Equal(t, bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	require.IsTypef(t, bbClient, BlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlobClient{}))
}

func TestContainerNewBlockBlobClient(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	bbClient := containerClient.NewBlockBlobClient(blobPrefix)

	require.Equal(t, bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	require.IsTypef(t, bbClient, BlockBlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlockBlobClient{}))
}

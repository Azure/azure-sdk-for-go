// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *azblobTestSuite) TestNewContainerClientValidName() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	testURL := svcClient.NewContainerClient(containerPrefix)

	correctURL := "https://" + os.Getenv(AccountNameEnvVar) + "." + DefaultBlobEndpointSuffix + containerPrefix
	_assert.Equal(testURL.URL(), correctURL)
}

func (s *azblobTestSuite) TestCreateRootContainerURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	testURL := svcClient.NewContainerClient(ContainerNameRoot)

	correctURL := "https://" + os.Getenv(AccountNameEnvVar) + ".blob.core.windows.net/$root"
	_assert.Equal(testURL.URL(), correctURL)
}

func (s *azblobTestSuite) TestContainerCreateInvalidName() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := svcClient.NewContainerClient("foo bar")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidResourceName)
}

func (s *azblobTestSuite) TestContainerCreateEmptyName() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := svcClient.NewContainerClient("")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *azblobTestSuite) TestContainerCreateNameCollision() {
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
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: map[string]string{},
	}

	containerClient = svcClient.NewContainerClient(containerName)
	_, err = containerClient.Create(ctx, &createContainerOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerAlreadyExists)
}

func (s *azblobTestSuite) TestContainerCreateInvalidMetadata() {
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
		Access:   &access,
		Metadata: map[string]string{"1 foo": "bar"},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)

	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *azblobTestSuite) TestContainerCreateNilMetadata() {
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
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(_assert, containerClient)
	_assert.Nil(err)

	response, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(response.Metadata)
}

func (s *azblobTestSuite) TestContainerCreateEmptyMetadata() {
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
		Access:   &access,
		Metadata: map[string]string{},
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(_assert, containerClient)
	_assert.Nil(err)

	response, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(response.Metadata)
}

//func (s *azblobTestSuite) TestContainerCreateAccessContainer() {
//	// TOD0: NotWorking
//	_assert := assert.New(s.T())
//testName := s.T().Name()
//	_context := getTestContext(testName)
//
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	credential, err := getGenericCredential("")
//	_assert.Nil(err)
//
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, svcClient)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := CreateContainerOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(_assert, containerClient)
//	_assert.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	}
//	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	_assert.Nil(err)
//
//	// Anonymous enumeration should be valid with container access
//	containerClient2, _ := NewContainerClient(containerClient.URL(), credential, nil)
//	pager := containerClient2.ListBlobsFlatSegment(nil)
//
//	for pager.NextPage(ctx) {
//		resp := pager.PageResponse()
//
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			_assert.Equal(*blob.Name, blobPrefix)
//		}
//	}
//
//	_assert.Nil(pager.Err())
//
//	// Getting blob data anonymously should still be valid with container access
//	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.EqualValues(resp.Metadata, basicMetadata)
//}

//func (s *azblobTestSuite) TestContainerCreateAccessBlob() {
//	// TODO: Not Working
//	_assert := assert.New(s.T())
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
//	createContainerOptions := CreateContainerOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(_assert, containerClient)
//	_assert.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	}
//	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
//	_assert.Nil(err)
//
//	// Reference the same container URL but with anonymous credentials
//	containerClient2, err := NewContainerClient(containerClient.URL(), azcore.AnonymousCredential(), nil)
//	_assert.Nil(err)
//
//	pager := containerClient2.ListBlobsFlatSegment(nil)
//
//	_assert.Equal(pager.NextPage(ctx), false)
//	_assert.NotNil(pager.Err())
//
//	// Accessing blob specific data should be public
//	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert.EqualValues(resp.Metadata, basicMetadata)
//}

func (s *azblobTestSuite) TestContainerCreateAccessNone() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	// Public Access Type None
	_, err = containerClient.Create(ctx, nil)
	defer deleteContainer(_assert, containerClient)

	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: basicMetadata,
	}
	_, err = bbClient.Upload(ctx, bytes.NewReader([]byte("Content")), &uploadBlockBlobOptions)
	_assert.Nil(err)

	// Reference the same container URL but with anonymous credentials
	containerClient2, err := NewContainerClient(containerClient.URL(), azcore.AnonymousCredential(), nil)
	_assert.Nil(err)

	pager := containerClient2.ListBlobsFlatSegment(nil)

	_assert.Equal(pager.NextPage(ctx), false)
	_assert.NotNil(pager.Err())

	// Blob data is not public
	blobURL2 := containerClient2.NewBlockBlobClient(blobPrefix)
	_, err = blobURL2.GetProperties(ctx, nil)
	_assert.NotNil(err)

	//serr := err.(StorageError)
	//_assert(serr.Response().StatusCode, chk.Equals, 401) // HEAD request does not return a status code
}

func validateContainerDeleted(_assert *assert.Assertions, containerClient ContainerClient) {
	_, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

func (s *azblobTestSuite) TestContainerDelete() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)

	_, err = containerClient.Delete(ctx, nil)
	_assert.Nil(err)

	validateContainerDeleted(_assert, containerClient)
}

func (s *azblobTestSuite) TestContainerDeleteNonExistent() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Delete(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

func (s *azblobTestSuite) TestContainerDeleteIfModifiedSinceTrue() {
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

	currentTime := getRelativeTimeFromAnchor(cResp.Date, -10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_assert.Nil(err)
	validateContainerDeleted(_assert, containerClient)
}

func (s *azblobTestSuite) TestContainerDeleteIfModifiedSinceFalse() {
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

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestContainerDeleteIfUnModifiedSinceTrue() {
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

	currentTime := getRelativeTimeFromAnchor(cResp.Date, 10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_assert.Nil(err)

	validateContainerDeleted(_assert, containerClient)
}

func (s *azblobTestSuite) TestContainerDeleteIfUnModifiedSinceFalse() {
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

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

////func (s *azblobTestSuite) TestContainerAccessConditionsUnsupportedConditions() {
////	// This test defines that the library will panic if the user specifies conditional headers
////	// that will be ignored by the service
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////
////	invalidEtag := "invalid"
////	deleteContainerOptions := SetMetadataContainerOptions{
////		Metadata: basicMetadata,
////		ModifiedAccessConditions: &ModifiedAccessConditions{
////			IfMatch: &invalidEtag,
////		},
////	}
////	_, err := containerClient.SetMetadata(ctx, &deleteContainerOptions)
////	_assert.NotNil(err)
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsNonexistentPrefix() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	prefix := blobPrefix + blobPrefix
////	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
////		Prefix: &prefix,
////	}
////	listResponse, errChan := containerClient.ListBlobsFlatSegment(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(listResponse, chk.IsNil)
////}
//
//func (s *azblobTestSuite) TestContainerListBlobsSpecificValidPrefix() {
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	prefix := blobPrefix
//	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
//		Prefix: &prefix,
//	}
//	pager := containerClient.ListBlobsFlatSegment(&containerListBlobFlatSegmentOptions)
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
//	defer deleteContainer(_assert, containerClient)
//	prefixes := []string{"a/1", "a/2", "b/2", "blob"}
//	blobNames := make([]string, 4)
//	for idx, prefix := range prefixes {
//		_, blobNames[idx] = createNewBlockBlobWithPrefix(c, containerClient, prefix)
//	}
//
//	pager := containerClient.ListBlobsHierarchySegment("/", nil)
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
//	//_assert.Nil(err)
//	//_assert(len(resp.Segment.BlobItems), chk.Equals, 1)
//	//_assert(len(resp.Segment.BlobPrefixes), chk.Equals, 2)
//	//_assert(resp.Segment.BlobPrefixes[0].Name, chk.Equals, "a/")
//	//_assert(resp.Segment.BlobPrefixes[1].Name, chk.Equals, "b/")
//	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

func (s *azblobTestSuite) TestContainerListBlobsWithSnapshots() {
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

	// initialize a blob and create a snapshot of it
	snapBlobName := generateBlobName(testName)
	snapBlob := createNewBlockBlob(_assert, snapBlobName, containerClient)
	snap, err := snapBlob.CreateSnapshot(ctx, nil)
	// snap.
	_assert.Nil(err)

	listBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots},
	}
	pager := containerClient.ListBlobsFlatSegment(&listBlobFlatSegmentOptions)

	wasFound := false // hold the for loop accountable for finding the blob and it's snapshot
	for pager.NextPage(ctx) {
		_assert.Nil(pager.Err())

		resp := pager.PageResponse()

		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
			if *blob.Name == snapBlobName && blob.Snapshot != nil {
				wasFound = true
				_assert.Equal(*blob.Snapshot, *snap.Snapshot)
			}
		}
	}
	_assert.Equal(wasFound, true)
}

func (s *azblobTestSuite) TestContainerListBlobsInvalidDelimiter() {
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
	prefixes := []string{"a/1", "a/2", "b/1", "blob"}
	for _, prefix := range prefixes {
		blobName := prefix + generateBlobName(testName)
		createNewBlockBlob(_assert, blobName, containerClient)
	}

	pager := containerClient.ListBlobsHierarchySegment("^", nil)

	pager.NextPage(ctx)
	_assert.Nil(pager.Err())
	_assert.Nil(pager.PageResponse().EnumerationResults.Segment.BlobPrefixes)
}

////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeMetadata() {
////	svcClient := getServiceClient()
////	container, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(container)
////	_, blobNameNoMetadata := createNewBlockBlobWithPrefix(c, container, "a")
////	blobMetadata, blobNameMetadata := createNewBlockBlobWithPrefix(c, container, "b")
////	_, err := blobMetadata.SetMetadata(ctx, Metadata{"field": "value"}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////
////	resp, err := container.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Metadata: true}})
////
////	_assert.Nil(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobNameNoMetadata)
////	_assert(resp.Segment.BlobItems[0].Metadata, chk.HasLen, 0)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobNameMetadata)
////	_assert(resp.Segment.BlobItems[1].Metadata["field"], chk.Equals, "value")
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsIncludeTypeSnapshots() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	blob, blobName := createNewBlockBlob(c, containerClient)
////	_, err := blob.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
////
////	_assert.Nil(err)
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
////	defer deleteContainer(_assert, containerClient)
////	bbClient, blobName := createNewBlockBlob(c, containerClient)
////	blobCopyURL, blobCopyName := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	_, err := blobCopyURL.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	_assert.Nil(err)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Copy: true}})
////
////	// These are sufficient to show that the blob copy was in fact included
////	_assert.Nil(err)
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
////	defer deleteContainer(_assert, containerClient)
////	bbClient, blobName := getBlockBlobURL(c, containerClient)
////	_, err := bbClient.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), LeaseAccessConditions{}, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{UncommittedBlobs: true}})
////
////	_assert.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////}
//
////func testContainerListBlobsIncludeTypeDeletedImpl(, svcClient ServiceURL) error {
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_assert.Nil(err)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
////	_assert.Nil(err)
////
////	resp, err = containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_assert.Nil(err)
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
////	defer deleteContainer(_assert, containerClient)
////
////	bbClient, _ := createNewBlockBlobWithPrefix(c, containerClient, "z")
////	_, err := bbClient.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	blobURL2, _ := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	resp2, err := blobURL2.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	_assert.Nil(err)
////	waitForCopy(c, blobURL2, resp2)
////	blobURL3, _ := createNewBlockBlobWithPrefix(c, containerClient, "deleted")
////
////	_, err = blobURL3.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Copy: true, Deleted: true, Versions: true}})
////
////	_assert.Nil(err)
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
////	defer deleteContainer(_assert, containerClient)
////	_, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: -2})
////	_assert(err, chk.Not(chk.IsNil))
////}
//
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsZero() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	createNewBlockBlob(c, containerClient)
////
////	maxResults := int32(0)
////	resp, errChan := containerClient.ListBlobsFlatSegment(ctx, 1, 0, &ContainerListBlobFlatSegmentOptions{Maxresults: &maxResults})
////
////	_assert(<-errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////}
//
//// TODO: Adele: Case failing
////func (s *azblobTestSuite) TestContainerListBlobsMaxResultsInsufficient() {
////	svcClient := getServiceClient()
////	containerClient, _ := createNewContainer(c, svcClient)
////	defer deleteContainer(_assert, containerClient)
////	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
////	createNewBlockBlobWithPrefix(c, containerClient, "b")
////
////	maxResults := int32(1)
////	resp, errChan := containerClient.ListBlobsFlatSegment(ctx, 3, 0, &ContainerListBlobFlatSegmentOptions{Maxresults: &maxResults})
////	_assert(<- errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////	_assert((<- resp).Name, chk.Equals, blobName)
////}

func (s *azblobTestSuite) TestContainerListBlobsMaxResultsExact() {
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
	blobNames := make([]string, 2)
	blobName := generateBlobName(testName)
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	createNewBlockBlob(_assert, blobNames[0], containerClient)
	createNewBlockBlob(_assert, blobNames[1], containerClient)

	maxResult := int32(2)
	pager := containerClient.ListBlobsFlatSegment(&ContainerListBlobFlatSegmentOptions{
		Maxresults: &maxResult,
	})

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
			_assert.Equal(nameMap[*blob.Name], true)
		}
	}

	_assert.Nil(pager.Err())
}

func (s *azblobTestSuite) TestContainerListBlobsMaxResultsSufficient() {
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

	blobNames := make([]string, 2)
	blobName := generateBlobName(testName)
	blobNames[0], blobNames[1] = "a"+blobName, "b"+blobName
	createNewBlockBlob(_assert, blobNames[0], containerClient)
	createNewBlockBlob(_assert, blobNames[1], containerClient)

	maxResult := int32(3)
	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
		Maxresults: &maxResult,
	}
	pager := containerClient.ListBlobsFlatSegment(&containerListBlobFlatSegmentOptions)

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
			_assert.Equal(nameMap[*blob.Name], true)
		}
	}

	_assert.Nil(pager.Err())
}

func (s *azblobTestSuite) TestContainerListBlobsNonExistentContainer() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	pager := containerClient.ListBlobsFlatSegment(nil)

	pager.NextPage(ctx)
	_assert.NotNil(pager.Err())
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
//	p := containerClient2.ListBlobsFlatSegment(nil)
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
////	resp, err := sasContainer.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert.Nil(err)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////
////	// Verifies that successful sas access is not just because it's public
////	anonymousBlobService := NewServiceURL(svcClient.URL(), sasPipeline)
////	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
////	_, err = anonymousContainer.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
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
	permissions := make([]*SignedIdentifier, 6, 6)
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
	permissions := make([]*SignedIdentifier, 2, 2)
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
	permissions := make([]*SignedIdentifier, 2, 2)
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
	permissions := make([]*SignedIdentifier, 2, 2)
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
	permissions := make([]*SignedIdentifier, 2, 2)
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
		ContainerAccessConditions: &ContainerAccessConditions{
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
		ContainerAccessConditions: &ContainerAccessConditions{
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
		ContainerAccessConditions: &ContainerAccessConditions{
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
		ContainerAccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func (s *azblobTestSuite) TestContainerGetPropertiesAndMetadataNoMetadata() {
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

	resp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestContainerGetPropsAndMetaNonExistentContainer() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.GetProperties(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

func (s *azblobTestSuite) TestContainerSetMetadataEmpty() {
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
		Metadata: basicMetadata,
		Access:   &access,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)

	defer deleteContainer(_assert, containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: map[string]string{},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestContainerSetMetadataNil() {
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
		Access:   &access,
		Metadata: basicMetadata,
	}
	_, err = containerClient.Create(ctx, &createContainerOptions)

	defer deleteContainer(_assert, containerClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	_assert.Nil(err)

	resp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *azblobTestSuite) TestContainerSetMetadataInvalidField() {
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

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: map[string]string{"!nval!d Field!@#%": "value"},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *azblobTestSuite) TestContainerSetMetadataNonExistent() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

//
//func (s *azblobTestSuite) TestContainerSetMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	svcClient := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, svcClient)
//
//	defer deleteContainer(_assert, containerClient)
//
//	setMetadataContainerOptions := SetMetadataContainerOptions{
//		Metadata: basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	_assert.Nil(err)
//
//	resp, err := containerClient.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert(resp.Metadata, chk.NotNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//
//}

//func (s *azblobTestSuite) TestContainerSetMetadataIfModifiedSinceFalse() {
//	// TODO: NotWorking
//	_assert := assert.New(s.T())
// testName := s.T().Name()
//	_context := getTestContext(testName)
//	svcClient := getServiceClient(&ClientOptions{
//		HTTPClient: _context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	containerClient, _ := createNewContainer(_assert, testName, svcClient)
//
//	defer deleteContainer(_assert, containerClient)
//
//	//currentTime := getRelativeTimeGMT(10)
//	//currentTime, err := time.Parse(time.UnixDate, "Wed Jan 07 11:11:11 PST 2099")
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	_assert.Nil(err)
//	setMetadataContainerOptions := SetMetadataContainerOptions{
//		Metadata: basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

func (s *azblobTestSuite) TestContainerNewBlobURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	bbClient := containerClient.NewBlobClient(blobPrefix)

	_assert.Equal(bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	_assert.IsTypef(bbClient, BlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlobClient{}))
}

func (s *azblobTestSuite) TestContainerNewBlockBlobClient() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	bbClient := containerClient.NewBlockBlobClient(blobPrefix)

	_assert.Equal(bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	_assert.IsTypef(bbClient, BlockBlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlockBlobClient{}))
}

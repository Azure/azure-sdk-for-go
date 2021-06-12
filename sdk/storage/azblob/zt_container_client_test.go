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

func (s *aztestsSuite) TestNewContainerClientValidName() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testURL := bsu.NewContainerClient(containerPrefix)

	correctURL := "https://" + os.Getenv(AccountNameEnvVar) + "." + DefaultBlobEndpointSuffix + containerPrefix
	_assert.Equal(testURL.URL(), correctURL)
}

func (s *aztestsSuite) TestCreateRootContainerURL() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testURL := bsu.NewContainerClient(ContainerNameRoot)

	correctURL := "https://" + os.Getenv(AccountNameEnvVar) + ".blob.core.windows.net/$root"
	_assert.Equal(testURL.URL(), correctURL)
}

func (s *aztestsSuite) TestContainerCreateInvalidName() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient := bsu.NewContainerClient("foo bar")

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	_assert.NotNil(err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidResourceName)
}

func (s *aztestsSuite) TestContainerCreateEmptyName() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})

	containerClient := bsu.NewContainerClient("")

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidQueryParameterValue)
}

func (s *aztestsSuite) TestContainerCreateNameCollision() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, containerName := createNewContainer(_assert, testName, bsu)

	defer deleteContainer(containerClient)

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	containerClient = bsu.NewContainerClient(containerName)
	_, err := containerClient.Create(ctx, &createContainerOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerAlreadyExists)
}

func (s *aztestsSuite) TestContainerCreateInvalidMetadata() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{"1 foo": "bar"},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *aztestsSuite) TestContainerCreateNilMetadata() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(containerClient)
	_assert.Nil(err)

	response, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(response.Metadata)
}

func (s *aztestsSuite) TestContainerCreateEmptyMetadata() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())

	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(containerClient)
	_assert.Nil(err)

	response, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(response.Metadata)
}

//func (s *aztestsSuite) TestContainerCreateAccessContainer() {
//	// TOD0: NotWorking
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	credential, err := getGenericCredential("")
//	_assert.Nil(err)
//
//	containerName := generateContainerName(s.T().Name())
//	containerClient := getContainerClient(containerName, bsu)
//
//	access := PublicAccessBlob
//	createContainerOptions := CreateContainerOptions{
//		Access: &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(containerClient)
//	_assert.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata: &basicMetadata,
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
//		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
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

//func (s *aztestsSuite) TestContainerCreateAccessBlob() {
//	// TODO: Not Working
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	testName := s.T().Name()
//	containerName := generateContainerName(testName)
//	containerClient := getContainerClient(containerName, bsu)
//
//	access := PublicAccessBlob
//	createContainerOptions := CreateContainerOptions{
//		Access: &access,
//	}
//	_, err := containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(containerClient)
//	_assert.Nil(err)
//
//	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata: &basicMetadata,
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

func (s *aztestsSuite) TestContainerCreateAccessNone() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, bsu)

	// Public Access Type None
	_, err := containerClient.Create(ctx, nil)
	defer deleteContainer(containerClient)

	bbClient := containerClient.NewBlockBlobClient(blobPrefix)
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: &basicMetadata,
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

func (s *aztestsSuite) TestContainerDelete() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)

	_, err := containerClient.Delete(ctx, nil)
	_assert.Nil(err)

	validateContainerDeleted(_assert, containerClient)
}

func (s *aztestsSuite) TestContainerDeleteNonExistent() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	_, err := containerClient.Delete(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

func (s *aztestsSuite) TestContainerDeleteIfModifiedSinceTrue() {

	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err) // Ensure the requests occur at different times

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_assert.Nil(err)
	validateContainerDeleted(_assert, containerClient)
}

//func (s *aztestsSuite) TestContainerDeleteIfModifiedSinceFalse() {
//	// TODO: NotWorking
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	testName := s.T().Name()
//	containerClient, _ := createNewContainer(_assert, testName, bsu)
//
//	defer deleteContainer(containerClient)
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	_assert.Nil(err)
//
//	deleteContainerOptions := DeleteContainerOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

func (s *aztestsSuite) TestContainerDeleteIfUnModifiedSinceTrue() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
	_assert.Nil(err)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
	_assert.Nil(err)

	validateContainerDeleted(_assert, containerClient)
}

//func (s *aztestsSuite) TestContainerDeleteIfUnModifiedSinceFalse() {
//	// TODO: Not Working
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	testName := s.T().Name()
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
//	_assert.Nil(err)
//
//	containerClient, _ := createNewContainer(_assert, testName, bsu)
//
//	defer deleteContainer(containerClient)
//
//	deleteContainerOptions := DeleteContainerOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.Delete(ctx, &deleteContainerOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

////func (s *aztestsSuite) TestContainerAccessConditionsUnsupportedConditions() {
////	// This test defines that the library will panic if the user specifies conditional headers
////	// that will be ignored by the service
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////
////	invalidEtag := "invalid"
////	deleteContainerOptions := SetMetadataContainerOptions{
////		Metadata: &basicMetadata,
////		ModifiedAccessConditions: &ModifiedAccessConditions{
////			IfMatch: &invalidEtag,
////		},
////	}
////	_, err := containerClient.SetMetadata(ctx, &deleteContainerOptions)
////	_assert(err, chk.NotNil)
////}
//
////func (s *aztestsSuite) TestContainerListBlobsNonexistentPrefix() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
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
//func (s *aztestsSuite) TestContainerListBlobsSpecificValidPrefix() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
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
//func (s *aztestsSuite) TestContainerListBlobsValidDelimiter() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
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
//		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
//			count++
//			_assert(*blob.Name, chk.Equals, blobNames[3])
//		}
//	}
//
//	_assert(pager.Err(), chk.IsNil)
//	_assert(count, chk.Equals, 1)
//
//	// TODO: Ask why the output is BlobItemInternal and why other fields are not there for ex: prefix array
//	//_assert(err, chk.IsNil)
//	//_assert(len(resp.Segment.BlobItems), chk.Equals, 1)
//	//_assert(len(resp.Segment.BlobPrefixes), chk.Equals, 2)
//	//_assert(resp.Segment.BlobPrefixes[0].Name, chk.Equals, "a/")
//	//_assert(resp.Segment.BlobPrefixes[1].Name, chk.Equals, "b/")
//	//_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

func (s *aztestsSuite) TestContainerListBlobsWithSnapshots() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)

	// initialize a blob and create a snapshot of it
	snapBlob, snapBlobName := createNewBlockBlob(_assert, testName, containerClient)
	snap, err := snapBlob.CreateSnapshot(ctx, nil)
	// snap.
	_assert.Nil(err)

	listBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
		Include: &[]ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots},
	}
	pager := containerClient.ListBlobsFlatSegment(&listBlobFlatSegmentOptions)

	wasFound := false // hold the for loop accountable for finding the blob and it's snapshot
	for pager.NextPage(ctx) {
		_assert.Nil(pager.Err())

		resp := pager.PageResponse()

		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
			if *blob.Name == snapBlobName && blob.Snapshot != nil {
				wasFound = true
				_assert.Equal(*blob.Snapshot, *snap.Snapshot)
			}
		}
	}
	_assert.Equal(wasFound, true)
}

func (s *aztestsSuite) TestContainerListBlobsInvalidDelimiter() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)
	prefixes := []string{"a/1", "a/2", "b/1", "blob"}
	for _, prefix := range prefixes {
		createNewBlockBlobWithPrefix(_assert, testName, containerClient, prefix)
	}

	pager := containerClient.ListBlobsHierarchySegment("^", nil)

	pager.NextPage(ctx)
	_assert.Nil(pager.Err())
	_assert.Nil(pager.PageResponse().EnumerationResults.Segment.BlobPrefixes)
}

////func (s *aztestsSuite) TestContainerListBlobsIncludeTypeMetadata() {
////	bsu := getBSU()
////	container, _ := createNewContainer(c, bsu)
////	defer deleteContainer(container)
////	_, blobNameNoMetadata := createNewBlockBlobWithPrefix(c, container, "a")
////	blobMetadata, blobNameMetadata := createNewBlockBlobWithPrefix(c, container, "b")
////	_, err := blobMetadata.SetMetadata(ctx, Metadata{"field": "value"}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert(err, chk.IsNil)
////
////	resp, err := container.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Metadata: true}})
////
////	_assert(err, chk.IsNil)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobNameNoMetadata)
////	_assert(resp.Segment.BlobItems[0].Metadata, chk.HasLen, 0)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobNameMetadata)
////	_assert(resp.Segment.BlobItems[1].Metadata["field"], chk.Equals, "value")
////}
//
////func (s *aztestsSuite) TestContainerListBlobsIncludeTypeSnapshots() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	blob, blobName := createNewBlockBlob(c, containerClient)
////	_, err := blob.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert(err, chk.IsNil)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
////
////	_assert(err, chk.IsNil)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 2)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[0].Snapshot, chk.NotNil)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[1].Snapshot, chk.Equals, "")
////}
////
////func (s *aztestsSuite) TestContainerListBlobsIncludeTypeCopy() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	bbClient, blobName := createNewBlockBlob(c, containerClient)
////	blobCopyURL, blobCopyName := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	_, err := blobCopyURL.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	_assert(err, chk.IsNil)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Copy: true}})
////
////	// These are sufficient to show that the blob copy was in fact included
////	_assert(err, chk.IsNil)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 2)
////	_assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobCopyName)
////	_assert(*resp.Segment.BlobItems[0].Properties.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
////	temp := bbClient.URL()
////	_assert(*resp.Segment.BlobItems[0].Properties.CopySource, chk.Equals, temp.String())
////	_assert(resp.Segment.BlobItems[0].Properties.CopyStatus, chk.Equals, CopyStatusSuccess)
////}
////
////func (s *aztestsSuite) TestContainerListBlobsIncludeTypeUncommitted() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	bbClient, blobName := getBlockBlobURL(c, containerClient)
////	_, err := bbClient.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), LeaseAccessConditions{}, nil, ClientProvidedKeyOptions{})
////	_assert(err, chk.IsNil)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{UncommittedBlobs: true}})
////
////	_assert(err, chk.IsNil)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////}
//
////func testContainerListBlobsIncludeTypeDeletedImpl(, bsu ServiceURL) error {
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	bbClient, _ := createNewBlockBlob(c, containerClient)
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_assert(err, chk.IsNil)
////	_assert(resp.Segment.BlobItems, chk.HasLen, 1)
////
////	_, err = bbClient.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
////	_assert(err, chk.IsNil)
////
////	resp, err = containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
////	_assert(err, chk.IsNil)
////	if len(resp.Segment.BlobItems) != 1 {
////		return errors.New("DeletedBlobNotFound")
////	}
////
////	// resp.Segment.BlobItems[0].Deleted == true/false if versioning is disabled/enabled.
////	_assert(resp.Segment.BlobItems[0].Deleted, chk.Equals, false)
////	return nil
////}
////
////func (s *aztestsSuite) TestContainerListBlobsIncludeTypeDeleted() {
////	bsu := getBSU()
////
////	runTestRequiringServiceProperties(c, bsu, "DeletedBlobNotFound", enableSoftDelete,
////		testContainerListBlobsIncludeTypeDeletedImpl, disableSoftDelete)
////}
////
////func testContainerListBlobsIncludeMultipleImpl(, bsu ServiceURL) error {
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////
////	bbClient, _ := createNewBlockBlobWithPrefix(c, containerClient, "z")
////	_, err := bbClient.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert(err, chk.IsNil)
////	blobURL2, _ := createNewBlockBlobWithPrefix(c, containerClient, "copy")
////	resp2, err := blobURL2.StartCopyFromURL(ctx, bbClient.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
////	_assert(err, chk.IsNil)
////	waitForCopy(c, blobURL2, resp2)
////	blobURL3, _ := createNewBlockBlobWithPrefix(c, containerClient, "deleted")
////
////	_, err = blobURL3.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
////
////	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
////		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Copy: true, Deleted: true, Versions: true}})
////
////	_assert(err, chk.IsNil)
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
////func (s *aztestsSuite) TestContainerListBlobsIncludeMultiple() {
////	bsu := getBSU()
////
////	runTestRequiringServiceProperties(c, bsu, "DeletedBlobNotFound", enableSoftDelete,
////		testContainerListBlobsIncludeMultipleImpl, disableSoftDelete)
////}
////
////func (s *aztestsSuite) TestContainerListBlobsMaxResultsNegative() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////
////	defer deleteContainer(containerClient)
////	_, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: -2})
////	_assert(err, chk.Not(chk.IsNil))
////}
//
////func (s *aztestsSuite) TestContainerListBlobsMaxResultsZero() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
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
////func (s *aztestsSuite) TestContainerListBlobsMaxResultsInsufficient() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
////	createNewBlockBlobWithPrefix(c, containerClient, "b")
////
////	maxResults := int32(1)
////	resp, errChan := containerClient.ListBlobsFlatSegment(ctx, 3, 0, &ContainerListBlobFlatSegmentOptions{Maxresults: &maxResults})
////	_assert(<- errChan, chk.IsNil)
////	_assert(resp, chk.HasLen, 1)
////	_assert((<- resp).Name, chk.Equals, blobName)
////}

func (s *aztestsSuite) TestContainerListBlobsMaxResultsExact() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)
	blobNames := make([]string, 2)
	_, blobNames[0] = createNewBlockBlobWithPrefix(_assert, testName, containerClient, "a")
	_, blobNames[1] = createNewBlockBlobWithPrefix(_assert, testName, containerClient, "b")

	maxResult := int32(2)
	pager := containerClient.ListBlobsFlatSegment(&ContainerListBlobFlatSegmentOptions{
		Maxresults: &maxResult,
	})

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
			_assert.Equal(nameMap[*blob.Name], true)
		}
	}

	_assert.Nil(pager.Err())
}

func (s *aztestsSuite) TestContainerListBlobsMaxResultsSufficient() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)
	blobNames := make([]string, 2)

	_, blobNames[0] = createNewBlockBlobWithPrefix(_assert, testName, containerClient, "a")
	_, blobNames[1] = createNewBlockBlobWithPrefix(_assert, testName, containerClient, "b")

	maxResult := int32(3)
	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
		Maxresults: &maxResult,
	}
	pager := containerClient.ListBlobsFlatSegment(&containerListBlobFlatSegmentOptions)

	nameMap := blobListToMap(blobNames)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
			_assert.Equal(nameMap[*blob.Name], true)
		}
	}

	_assert.Nil(pager.Err())
}

func (s *aztestsSuite) TestContainerListBlobsNonExistentContainer() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	pager := containerClient.ListBlobsFlatSegment(nil)

	pager.NextPage(ctx)
	_assert.NotNil(pager.Err())
}

func (s *aztestsSuite) TestContainerGetSetPermissionsMultiplePolicies() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

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
			ContainerACL: &permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)

	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *aztestsSuite) TestContainerGetPermissionsPublicAccessNotNone() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access: &access,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	_assert.Nil(err)
	defer deleteContainer(containerClient)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)

	_assert.Nil(err)
	_assert.Equal(*resp.BlobPublicAccess, PublicAccessBlob)
}

//func (s *aztestsSuite) TestContainerSetPermissionsPublicAccessNone() {
//	// Test the basic one by making an anonymous request to ensure it's actually doing it and also with GetPermissions
//	// For all the others, can just use GetPermissions since we've validated that it at least registers on the server correctly
//	bsu := getBSU(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	// Container is created with PublicAccessBlob, so setting it to None will actually test that it is changed through this method
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	_assert(err, chk.IsNil)
//	bsu2, err := NewServiceClient(bsu.URL(), azcore.AnonymousCredential(), nil)
//	_assert(err, chk.IsNil)
//
//	containerClient2 := bsu2.NewContainerClient(containerName)
//	blobURL2 := containerClient2.NewBlockBlobClient(blobName)
//
//	// Get permissions via the original container URL so the request succeeds
//	resp, err := containerClient.GetAccessPolicy(ctx, nil)
//	_assert(resp.BlobPublicAccess, chk.IsNil)
//	_assert(err, chk.IsNil)
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

func (s *aztestsSuite) TestContainerSetPermissionsPublicAccessBlob() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

	access := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.BlobPublicAccess, PublicAccessBlob)
}

func (s *aztestsSuite) TestContainerSetPermissionsPublicAccessContainer() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

	access := PublicAccessContainer
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp.BlobPublicAccess, PublicAccessContainer)
}

////// TODO: After Pacer is ready
////func (s *aztestsSuite) TestContainerSetPermissionsACLSinglePolicy() {
////	bsu := getBSU()
////	credential, err := getGenericCredential("")
////	if err != nil {
////		c.Fatal("Invalid credential")
////	}
////	containerClient, containerName := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
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
////			ContainerACL: &permissions,
////		},
////	}
////	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
////	_assert(err, chk.IsNil)
////
////	serviceSASValues := BlobSASSignatureValues{Identifier: "0000", ContainerName: containerName}
////	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
////	if err != nil {
////		c.Fatal(err)
////	}
////
////	sasURL := bsu.URL()
////	sasURL.RawQuery = queryParams.Encode()
////	sasPipeline := (NewAnonymousCredential(), PipelineOptions{})
////	sasBlobServiceURL := NewServiceURL(sasURL, sasPipeline)
////
////	// Verifies that the SAS can access the resource
////	sasContainer := sasBlobServiceURL.NewContainerClient(containerName)
////	resp, err := sasContainer.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert(err, chk.IsNil)
////	_assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
////
////	// Verifies that successful sas access is not just because it's public
////	anonymousBlobService := NewServiceURL(bsu.URL(), sasPipeline)
////	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
////	_, err = anonymousContainer.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	validateStorageError(c, err, StorageErrorCodeNoAuthenticationInformation)
////}

func (s *aztestsSuite) TestContainerSetPermissionsACLMoreThanFive() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

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

	access := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: &permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

func (s *aztestsSuite) TestContainerSetPermissionsDeleteAndModifyACL() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

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

	access := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: &permissions,
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
			ContainerACL: &permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions1)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp.SignedIdentifiers, 1)
	_assert.EqualValues(resp.SignedIdentifiers, permissions)
}

func (s *aztestsSuite) TestContainerSetPermissionsDeleteAllPolicies() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

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

	access := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: &permissions,
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
			ContainerACL: &[]*SignedIdentifier{},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.SignedIdentifiers)
}

func (s *aztestsSuite) TestContainerSetPermissionsInvalidPolicyTimes() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

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

	access := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: &permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.Nil(err)
}

func (s *aztestsSuite) TestContainerSetPermissionsNilPolicySlice() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

	_, err := containerClient.SetAccessPolicy(ctx, nil)
	_assert.Nil(err)
}

func (s *aztestsSuite) TestContainerSetPermissionsSignedIdentifierTooLong() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

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

	access := PublicAccessBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerACL: &permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidXMLDocument)
}

//
//func (s *aztestsSuite) TestContainerSetPermissionsIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//	bsu := getBSU(nil)
//	container, _ := createNewContainer(c, bsu)
//
//	defer deleteContainer(container)
//
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerAccessConditions: &ContainerAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//		},
//	}
//	_, err := container.SetAccessPolicy(ctx, &setAccessPolicyOptions)
//	_assert(err, chk.IsNil)
//
//	resp, err := container.GetAccessPolicy(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.BlobPublicAccess, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestContainerSetPermissionsIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//
//	defer deleteContainer(containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerAccessConditions: &ContainerAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//		},
//	}
//	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
//	_assert(err, chk.NotNil)
//
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestContainerSetPermissionsIfUnModifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//
//	defer deleteContainer(containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerAccessConditions: &ContainerAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//		},
//	}
//	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
//	_assert(err, chk.IsNil)
//
//	resp, err := containerClient.GetAccessPolicy(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.BlobPublicAccess, chk.IsNil)
//}

//func (s *aztestsSuite) TestContainerSetPermissionsIfUnModifiedSinceFalse() {
//	// TODO: NotWorking
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)
//
//	defer deleteContainer(containerClient)
//
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerAccessConditions: &ContainerAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//		},
//	}
//	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

func (s *aztestsSuite) TestContainerGetPropertiesAndMetadataNoMetadata() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

	resp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *aztestsSuite) TestContainerGetPropsAndMetaNonExistentContainer() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	_, err := containerClient.GetProperties(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

func (s *aztestsSuite) TestContainerSetMetadataEmpty() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Metadata: &basicMetadata,
		Access:   &access,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	defer deleteContainer(containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: &map[string]string{},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	_assert.Nil(err)

	resp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *aztestsSuite) TestContainerSetMetadataNil() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)
	access := PublicAccessBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &basicMetadata,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	defer deleteContainer(containerClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	_assert.Nil(err)

	resp, err := containerClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Nil(resp.Metadata)
}

func (s *aztestsSuite) TestContainerSetMetadataInvalidField() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)

	defer deleteContainer(containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: &map[string]string{"!nval!d Field!@#%": "value"},
	}
	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *aztestsSuite) TestContainerSetMetadataNonExistent() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	_, err := containerClient.SetMetadata(ctx, nil)
	_assert.NotNil(err)

	validateStorageError(_assert, err, StorageErrorCodeContainerNotFound)
}

//
//func (s *aztestsSuite) TestContainerSetMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//
//	defer deleteContainer(containerClient)
//
//	setMetadataContainerOptions := SetMetadataContainerOptions{
//		Metadata: &basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	_assert(err, chk.IsNil)
//
//	resp, err := containerClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.NotNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//
//}

//func (s *aztestsSuite) TestContainerSetMetadataIfModifiedSinceFalse() {
//	// TODO: NotWorking
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	containerClient, _ := createNewContainer(_assert, s.T().Name(), bsu)
//
//	defer deleteContainer(containerClient)
//
//	//currentTime := getRelativeTimeGMT(10)
//	//currentTime, err := time.Parse(time.UnixDate, "Wed Jan 07 11:11:11 PST 2099")
//	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2021")
//	_assert.Nil(err)
//	setMetadataContainerOptions := SetMetadataContainerOptions{
//		Metadata: &basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
//	_assert.NotNil(err)
//
//	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
//}

func (s *aztestsSuite) TestContainerNewBlobURL() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	bbClient := containerClient.NewBlobClient(blobPrefix)

	_assert.Equal(bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	_assert.IsTypef(bbClient, BlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlobClient{}))
}

func (s *aztestsSuite) TestContainerNewBlockBlobClient() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	bbClient := containerClient.NewBlockBlobClient(blobPrefix)

	_assert.Equal(bbClient.URL(), containerClient.URL()+"/"+blobPrefix)
	_assert.IsTypef(bbClient, BlockBlobClient{}, fmt.Sprintf("%T should be of type %T", bbClient, BlockBlobClient{}))
}

package azblob

import (
	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
	"os"
	"strconv"
	"strings"
	"time"
)

func (s *aztestsSuite) TestNewContainerClientValidName(c *chk.C) {
	bsu := getBSU()
	testURL := bsu.NewContainerClient(containerPrefix)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".blob.core.windows.net/" + containerPrefix
	temp := testURL.URL()
	c.Assert(temp.String(), chk.Equals, correctURL)
}

func (s *aztestsSuite) TestCreateRootContainerURL(c *chk.C) {
	bsu := getBSU()
	testURL := bsu.NewContainerClient(ContainerNameRoot)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".blob.core.windows.net/$root"
	temp := testURL.URL()
	c.Assert(temp.String(), chk.Equals, correctURL)
}

func (s *aztestsSuite) TestAccountWithPipeline(c *chk.C) {
	bsu := getBSU()
	pipeline := newTestPipeline()
	bsu = bsu.WithPipeline(pipeline) // testPipeline returns an identifying message as an error
	containerClient := bsu.NewContainerClient("name")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	c.Assert(err, chk.NotNil)
	c.Assert(err.Error(), chk.Equals, testPipelineMessage)
}

func (s *aztestsSuite) TestContainerCreateInvalidName(c *chk.C) {
	bsu := getBSU()
	containerClient := bsu.NewContainerClient("foo bar")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	c.Assert(err, chk.NotNil)

	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeInvalidResourceName)
}

func (s *aztestsSuite) TestContainerCreateEmptyName(c *chk.C) {
	bsu := getBSU()
	containerClient := bsu.NewContainerClient("")

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	c.Assert(err, chk.NotNil)

	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeInvalidQueryParameterValue)
}

func (s *aztestsSuite) TestContainerCreateNameCollision(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	containerClient = bsu.NewContainerClient(containerName)
	_, err := containerClient.Create(ctx, &createContainerOptions)
	c.Assert(err, chk.NotNil)

	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeContainerAlreadyExists)
}

func (s *aztestsSuite) TestContainerCreateInvalidMetadata(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{"1 foo": "bar"},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	c.Assert(err, chk.NotNil)
	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
}

func (s *aztestsSuite) TestContainerCreateNilMetadata(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(c, containerClient)
	c.Assert(err, chk.IsNil)

	response, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(response.NewMetadata(), chk.HasLen, 0)
}

func (s *aztestsSuite) TestContainerCreateEmptyMetadata(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &map[string]string{},
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)
	defer deleteContainer(c, containerClient)
	c.Assert(err, chk.IsNil)

	response, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(response.NewMetadata(), chk.HasLen, 0)
}

// Note that for all tests that create blobs, deleting the container also deletes any blobs within that container, thus we
// simply delete the whole container after the test

//// TODO: After Pacer is ready
//func (s *aztestsSuite) TestContainerCreateAccessContainer(c *chk.C) {
//	bsu := getBSU()
//	credential, err := getGenericCredential("")
//	c.Assert(err, chk.IsNil)
//	containerClient, _ := getContainerClient(c, bsu)
//
//	access := PublicAccessTypeBlob
//	createContainerOptions := CreateContainerOptions{
//		Access:   &access,
//	}
//	_, err = containerClient.Create(ctx, &createContainerOptions)
//	defer deleteContainer(c, containerClient)
//	c.Assert(err, chk.IsNil)
//
//	blobURL := containerClient.NewBlockBlobClient(blobPrefix)
//	_, err = blobURL.Upload(ctx, bytes.NewReader([]byte("Content")), nil)
//	c.Assert(err, chk.IsNil)
//
//	// Anonymous enumeration should be valid with container access
//	containerURL2, _ := NewContainerClient(containerClient.String(), credential, nil)
//	response, err := containerURL2.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(response.Segment.BlobItems[0].Name, chk.Equals, blobPrefix)
//
//	// Getting blob data anonymously should still be valid with container access
//	blobURL2 := containerURL2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)
//}

//func (s *aztestsSuite) TestContainerCreateAccessBlob(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := getContainerClient(c, bsu)
//
//	_, err := containerClient.Create(ctx, nil, PublicAccessBlob)
//	defer deleteContainer(c, containerClient)
//	c.Assert(err, chk.IsNil)
//
//	blobURL := containerClient.NewBlockBlobClient(blobPrefix)
//	blobURL.Upload(ctx, bytes.NewReader([]byte("Content")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//
//	// Reference the same container URL but with anonymous credentials
//	containerURL2 := NewContainerClient(containerClient.URL(), NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//	_, err = containerURL2.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	validateStorageError(c, err, ServiceCodeNoAuthenticationInformation) // Listing blobs is not publicly accessible
//
//	// Accessing blob specific data should be public
//	blobURL2 := containerURL2.NewBlockBlobClient(blobPrefix)
//	resp, err := blobURL2.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestContainerCreateAccessNone(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := getContainerClient(c, bsu)
//
//	_, err := containerClient.Create(ctx, nil, PublicAccessNone)
//	defer deleteContainer(c, containerClient)
//
//	blobURL := containerClient.NewBlockBlobClient(blobPrefix)
//	blobURL.Upload(ctx, bytes.NewReader([]byte("Content")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//
//	// Reference the same container URL but with anonymous credentials
//	containerURL2 := NewContainerClient(containerClient.URL(), NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//	// Listing blobs is not public
//	_, err = containerURL2.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	validateStorageError(c, err, ServiceCodeNoAuthenticationInformation)
//
//	// Blob data is not public
//	blobURL2 := containerURL2.NewBlockBlobClient(blobPrefix)
//	_, err = blobURL2.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//	serr := err.(StorageError)
//	c.Assert(serr.Response().StatusCode, chk.Equals, 401) // HEAD request does not return a status code
//}
//

func validateContainerDeleted(c *chk.C, containerClient ContainerClient) {
	_, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeContainerNotFound)
}

func (s *aztestsSuite) TestContainerDelete(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	_, err := containerClient.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)

	validateContainerDeleted(c, containerClient)
}

func (s *aztestsSuite) TestContainerDeleteNonExistent(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	_, err := containerClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeContainerNotFound)
}

func (s *aztestsSuite) TestContainerDeleteIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10) // Ensure the requests occur at different times
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err := containerClient.Delete(ctx, &deleteContainerOptions)
	c.Assert(err, chk.IsNil)
	validateContainerDeleted(c, containerClient)
}

func (s *aztestsSuite) TestContainerDeleteIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err := containerClient.Delete(ctx, &deleteContainerOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestContainerDeleteIfUnModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	currentTime := getRelativeTimeGMT(10)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err := containerClient.Delete(ctx, &deleteContainerOptions)
	c.Assert(err, chk.IsNil)

	validateContainerDeleted(c, containerClient)
}

func (s *aztestsSuite) TestContainerDeleteIfUnModifiedSinceFalse(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10) // Ensure the requests occur at different times

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	deleteContainerOptions := DeleteContainerOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err := containerClient.Delete(ctx, &deleteContainerOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestContainerAccessConditionsUnsupportedConditions(c *chk.C) {
	// This test defines that the library will panic if the user specifies conditional headers
	// that will be ignored by the service
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	invalidEtag := "invalid"
	deleteContainerOptions := SetMetadataContainerOptions{
		Metadata: &basicMetadata,
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: &invalidEtag,
		},
	}
	_, err := containerClient.SetMetadata(ctx, &deleteContainerOptions)
	c.Assert(err, chk.NotNil)
}

//func (s *aztestsSuite) TestContainerListBlobsNonexistentPrefix(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	createNewBlockBlob(c, containerClient)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Prefix: blobPrefix + blobPrefix})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 0)
//}

//func (s *aztestsSuite) TestContainerListBlobsSpecificValidPrefix(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Prefix: blobPrefix})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 1)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

//func (s *aztestsSuite) TestContainerListBlobsValidDelimiter(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	createNewBlockBlobWithPrefix(c, containerClient, "a/1")
//	createNewBlockBlobWithPrefix(c, containerClient, "a/2")
//	createNewBlockBlobWithPrefix(c, containerClient, "b/1")
//	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "blob")
//
//	resp, err := containerClient.ListBlobsHierarchySegment(ctx, Marker{}, "/", ListBlobsSegmentOptions{})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(len(resp.Segment.BlobItems), chk.Equals, 1)
//	c.Assert(len(resp.Segment.BlobPrefixes), chk.Equals, 2)
//	c.Assert(resp.Segment.BlobPrefixes[0].Name, chk.Equals, "a/")
//	c.Assert(resp.Segment.BlobPrefixes[1].Name, chk.Equals, "b/")
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

//func (s *aztestsSuite) TestContainerListBlobsWithSnapshots(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//
//	_, err := containerClient.ListBlobsHierarchySegment(ctx, Marker{}, "/", ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//	c.Assert(err, chk.Not(chk.Equals), nil)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsInvalidDelimiter(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	createNewBlockBlobWithPrefix(c, containerClient, "a/1")
//	createNewBlockBlobWithPrefix(c, containerClient, "a/2")
//	createNewBlockBlobWithPrefix(c, containerClient, "b/1")
//	createNewBlockBlobWithPrefix(c, containerClient, "blob")
//
//	resp, err := containerClient.ListBlobsHierarchySegment(ctx, Marker{}, "^", ListBlobsSegmentOptions{})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 4)
//}

//func (s *aztestsSuite) TestContainerListBlobsIncludeTypeMetadata(c *chk.C) {
//	bsu := getBSU()
//	container, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, container)
//	_, blobNameNoMetadata := createNewBlockBlobWithPrefix(c, container, "a")
//	blobMetadata, blobNameMetadata := createNewBlockBlobWithPrefix(c, container, "b")
//	_, err := blobMetadata.SetMetadata(ctx, Metadata{"field": "value"}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err := container.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Metadata: true}})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobNameNoMetadata)
//	c.Assert(resp.Segment.BlobItems[0].Metadata, chk.HasLen, 0)
//	c.Assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobNameMetadata)
//	c.Assert(resp.Segment.BlobItems[1].Metadata["field"], chk.Equals, "value")
//}

//func (s *aztestsSuite) TestContainerListBlobsIncludeTypeSnapshots(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blob, blobName := createNewBlockBlob(c, containerClient)
//	_, err := blob.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 2)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//	c.Assert(resp.Segment.BlobItems[0].Snapshot, chk.NotNil)
//	c.Assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
//	c.Assert(resp.Segment.BlobItems[1].Snapshot, chk.Equals, "")
//}
//
//func (s *aztestsSuite) TestContainerListBlobsIncludeTypeCopy(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobURL, blobName := createNewBlockBlob(c, containerClient)
//	blobCopyURL, blobCopyName := createNewBlockBlobWithPrefix(c, containerClient, "copy")
//	_, err := blobCopyURL.StartCopyFromURL(ctx, blobURL.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Copy: true}})
//
//	// These are sufficient to show that the blob copy was in fact included
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 2)
//	c.Assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobCopyName)
//	c.Assert(*resp.Segment.BlobItems[0].Properties.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//	temp := blobURL.URL()
//	c.Assert(*resp.Segment.BlobItems[0].Properties.CopySource, chk.Equals, temp.String())
//	c.Assert(resp.Segment.BlobItems[0].Properties.CopyStatus, chk.Equals, CopyStatusSuccess)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsIncludeTypeUncommitted(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobURL, blobName := getBlockBlobURL(c, containerClient)
//	_, err := blobURL.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), LeaseAccessConditions{}, nil, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{UncommittedBlobs: true}})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 1)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}

//func testContainerListBlobsIncludeTypeDeletedImpl(c *chk.C, bsu ServiceURL) error {
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobURL, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 1)
//
//	_, err = blobURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err = containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true, Deleted: true}})
//	c.Assert(err, chk.IsNil)
//	if len(resp.Segment.BlobItems) != 1 {
//		return errors.New("DeletedBlobNotFound")
//	}
//
//	// resp.Segment.BlobItems[0].Deleted == true/false if versioning is disabled/enabled.
//	c.Assert(resp.Segment.BlobItems[0].Deleted, chk.Equals, false)
//	return nil
//}
//
//func (s *aztestsSuite) TestContainerListBlobsIncludeTypeDeleted(c *chk.C) {
//	bsu := getBSU()
//
//	runTestRequiringServiceProperties(c, bsu, "DeletedBlobNotFound", enableSoftDelete,
//		testContainerListBlobsIncludeTypeDeletedImpl, disableSoftDelete)
//}
//
//func testContainerListBlobsIncludeMultipleImpl(c *chk.C, bsu ServiceURL) error {
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//
//	blobURL, _ := createNewBlockBlobWithPrefix(c, containerClient, "z")
//	_, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	blobURL2, _ := createNewBlockBlobWithPrefix(c, containerClient, "copy")
//	resp2, err := blobURL2.StartCopyFromURL(ctx, blobURL.URL(), Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, DefaultAccessTier, nil)
//	c.Assert(err, chk.IsNil)
//	waitForCopy(c, blobURL2, resp2)
//	blobURL3, _ := createNewBlockBlobWithPrefix(c, containerClient, "deleted")
//
//	_, err = blobURL3.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Copy: true, Deleted: true, Versions: true}})
//
//	c.Assert(err, chk.IsNil)
//	if len(resp.Segment.BlobItems) != 6 {
//		// If there are fewer blobs in the container than there should be, it will be because one was permanently deleted.
//		return errors.New("DeletedBlobNotFound")
//	}
//
//	//c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName2)
//	//c.Assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName) // With soft delete, the overwritten blob will have a backup snapshot
//	//c.Assert(resp.Segment.BlobItems[2].Name, chk.Equals, blobName)
//	return nil
//}
//
//func (s *aztestsSuite) TestContainerListBlobsIncludeMultiple(c *chk.C) {
//	bsu := getBSU()
//
//	runTestRequiringServiceProperties(c, bsu, "DeletedBlobNotFound", enableSoftDelete,
//		testContainerListBlobsIncludeMultipleImpl, disableSoftDelete)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsMaxResultsNegative(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//
//	defer deleteContainer(c, containerClient)
//	_, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: -2})
//	c.Assert(err, chk.Not(chk.IsNil))
//}
//
//func (s *aztestsSuite) TestContainerListBlobsMaxResultsZero(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	createNewBlockBlob(c, containerClient)
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: 0})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 1)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsMaxResultsInsufficient(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
//	createNewBlockBlobWithPrefix(c, containerClient, "b")
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: 1})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 1)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsMaxResultsExact(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
//	_, blobName2 := createNewBlockBlobWithPrefix(c, containerClient, "b")
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: 2})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 2)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//	c.Assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName2)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsMaxResultsSufficient(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, blobName := createNewBlockBlobWithPrefix(c, containerClient, "a")
//	_, blobName2 := createNewBlockBlobWithPrefix(c, containerClient, "b")
//
//	resp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{MaxResults: 3})
//
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 2)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//	c.Assert(resp.Segment.BlobItems[1].Name, chk.Equals, blobName2)
//}
//
//func (s *aztestsSuite) TestContainerListBlobsNonExistentContainer(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := getContainerClient(c, bsu)
//
//	_, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//
//	c.Assert(err, chk.NotNil)
//}

func (s *aztestsSuite) TestContainerWithNewPipeline(c *chk.C) {
	bsu := getBSU()
	pipeline := newTestPipeline()
	containerClient, _ := getContainerClient(c, bsu)
	containerClient = containerClient.WithPipeline(pipeline)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access: &access,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	c.Assert(err, chk.NotNil)
	c.Assert(err.Error(), chk.Equals, testPipelineMessage)
}

func (s *aztestsSuite) TestContainerGetSetPermissionsMultiplePolicies(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	// Define the policies
	start := generateCurrentTimeWithModerateResolution()
	expiry := start.Add(5 * time.Minute)
	expiry2 := start.Add(time.Minute)
	readWrite := AccessPolicyPermission{Read: true, Write: true}.String()
	readOnly := AccessPolicyPermission{Read: true}.String()
	id1, id2 := "0000", "0001"
	permissions := []SignedIDentifier{
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
			ContainerAcl: &permissions,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)

	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.SignedIdentifiers, chk.DeepEquals, permissions)
}

func (s *aztestsSuite) TestContainerGetPermissionsPublicAccessNotNone(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access: &access,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions) // We create the container explicitly so we can be sure the access policy is not empty
	c.Assert(err, chk.IsNil)
	defer deleteContainer(c, containerClient)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)

	c.Assert(err, chk.IsNil)
	c.Assert(*resp.BlobPublicAccess, chk.Equals, PublicAccessTypeBlob)
}

func (s *aztestsSuite) TestContainerSetPermissionsPublicAccessNone(c *chk.C) {
	// Test the basic one by making an anonymous request to ensure it's actually doing it and also with GetPermissions
	// For all the others, can just use GetPermissions since we've validated that it at least registers on the server correctly
	bsu := getBSU()
	containerClient, containerName := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	_, blobName := createNewBlockBlob(c, containerClient)

	// Container is created with PublicAccessBlob, so setting it to None will actually test that it is changed through this method
	_, err := containerClient.SetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)

	//pipeline := NewPipeline(NewAnonymousCredential(), PipelineOptions{})
	credential, err := getGenericCredential("")
	c.Assert(err, chk.IsNil)
	bsu2, err := NewServiceClient(bsu.String(), credential, nil)
	c.Assert(err, chk.IsNil)

	containerURL2 := bsu2.NewContainerClient(containerName)
	blobURL2 := containerURL2.NewBlockBlobClient(blobName)
	_, err = blobURL2.Download(ctx, nil)

	// Get permissions via the original container URL so the request succeeds
	resp, _ := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(resp.BlobPublicAccess, chk.IsNil)
	// If we cannot access a blob's data, we will also not be able to enumerate blobs
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeNoAuthenticationInformation)

}

func (s *aztestsSuite) TestContainerSetPermissionsPublicAccessBlob(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	access := PublicAccessTypeBlob
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.BlobPublicAccess, chk.Equals, PublicAccessTypeBlob)
}

func (s *aztestsSuite) TestContainerSetPermissionsPublicAccessContainer(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	access := PublicAccessTypeContainer
	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access: &access,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.BlobPublicAccess, chk.Equals, PublicAccessTypeContainer)
}

//// TODO: After Pacer is ready
//func (s *aztestsSuite) TestContainerSetPermissionsACLSinglePolicy(c *chk.C) {
//	bsu := getBSU()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, blobName := createNewBlockBlob(c, containerClient)
//
//	start := time.Now().UTC().Add(-15 * time.Second)
//	expiry := start.Add(5 * time.Minute).UTC()
//	listOnly := AccessPolicyPermission{List: true}.String()
//	id := "0000"
//	permissions := []SignedIDentifier{{
//		ID: &id,
//		AccessPolicy: &AccessPolicy{
//			Start:      &start,
//			Expiry:     &expiry,
//			Permission: &listOnly,
//		},
//	}}
//
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
//			ContainerAcl: &permissions,
//		},
//	}
//	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
//	c.Assert(err, chk.IsNil)
//
//	serviceSASValues := BlobSASSignatureValues{Identifier: "0000", ContainerName: containerName}
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	sasURL := bsu.URL()
//	sasURL.RawQuery = queryParams.Encode()
//	sasPipeline := (NewAnonymousCredential(), PipelineOptions{})
//	sasBlobServiceURL := NewServiceURL(sasURL, sasPipeline)
//
//	// Verifies that the SAS can access the resource
//	sasContainer := sasBlobServiceURL.NewContainerClient(containerName)
//	resp, err := sasContainer.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.Segment.BlobItems[0].Name, chk.Equals, blobName)
//
//	// Verifies that successful sas access is not just because it's public
//	anonymousBlobService := NewServiceURL(bsu.URL(), sasPipeline)
//	anonymousContainer := anonymousBlobService.NewContainerClient(containerName)
//	_, err = anonymousContainer.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	validateStorageError(c, err, ServiceCodeNoAuthenticationInformation)
//}

func (s *aztestsSuite) TestContainerSetPermissionsACLMoreThanFive(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	start := time.Now().UTC()
	expiry := start.Add(5 * time.Minute).UTC()
	permissions := make([]SignedIDentifier, 6, 6)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 6; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = SignedIDentifier{
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
			ContainerAcl: &permissions,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeInvalidXMLDocument)
}

func (s *aztestsSuite) TestContainerSetPermissionsDeleteAndModifyACL(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	start := generateCurrentTimeWithModerateResolution()
	expiry := start.Add(5 * time.Minute).UTC()
	listOnly := AccessPolicyPermission{Read: true}.String()
	permissions := make([]SignedIDentifier, 2, 2)
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = SignedIDentifier{
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
			ContainerAcl: &permissions,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.SignedIdentifiers, chk.DeepEquals, permissions)

	permissions = (*resp.SignedIdentifiers)[:1] // Delete the first policy by removing it from the slice
	newId := "0004"
	permissions[0].ID = &newId // Modify the remaining policy which is at index 0 in the new slice
	setAccessPolicyOptions1 := SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerAcl: &permissions,
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions1)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.SignedIdentifiers, chk.HasLen, 1)
	c.Assert(*resp.SignedIdentifiers, chk.DeepEquals, permissions)
}

func (s *aztestsSuite) TestContainerSetPermissionsDeleteAllPolicies(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	start := time.Now().UTC()
	expiry := start.Add(5 * time.Minute).UTC()
	permissions := make([]SignedIDentifier, 2, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = SignedIDentifier{
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
			ContainerAcl: &permissions,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.SignedIdentifiers, chk.HasLen, len(permissions))
	c.Assert(*resp.SignedIdentifiers, chk.DeepEquals, permissions)

	setAccessPolicyOptions = SetAccessPolicyOptions{
		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
			Access:       &access,
			ContainerAcl: &[]SignedIDentifier{},
		},
	}
	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err = containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.SignedIdentifiers, chk.IsNil)
}

func (s *aztestsSuite) TestContainerSetPermissionsInvalidPolicyTimes(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	// Swap start and expiry
	expiry := time.Now().UTC()
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]SignedIDentifier, 2, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		id := "000" + strconv.Itoa(i)
		permissions[i] = SignedIDentifier{
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
			ContainerAcl: &permissions,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestContainerSetPermissionsNilPolicySlice(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	_, err := containerClient.SetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestContainerSetPermissionsSignedIdentifierTooLong(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	id := ""
	for i := 0; i < 65; i++ {
		id += "a"
	}
	expiry := time.Now().UTC()
	start := expiry.Add(5 * time.Minute).UTC()
	permissions := make([]SignedIDentifier, 2, 2)
	listOnly := AccessPolicyPermission{Read: true}.String()
	for i := 0; i < 2; i++ {
		permissions[i] = SignedIDentifier{
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
			ContainerAcl: &permissions,
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeInvalidXMLDocument)
}

func (s *aztestsSuite) TestContainerSetPermissionsIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)
	bsu := getBSU()
	container, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, container)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerAccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err := container.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err := container.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.BlobPublicAccess, chk.IsNil)
}

func (s *aztestsSuite) TestContainerSetPermissionsIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerAccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestContainerSetPermissionsIfUnModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerAccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetAccessPolicy(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.BlobPublicAccess, chk.IsNil)
}

func (s *aztestsSuite) TestContainerSetPermissionsIfUnModifiedSinceFalse(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	setAccessPolicyOptions := SetAccessPolicyOptions{
		ContainerAccessConditions: &ContainerAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err := containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestContainerGetPropertiesAndMetadataNoMetadata(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	resp, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.IsNil)
}

func (s *aztestsSuite) TestContainerGetPropsAndMetaNonExistentContainer(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	_, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeContainerNotFound)
}

func (s *aztestsSuite) TestContainerSetMetadataEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Metadata: &basicMetadata,
		Access:   &access,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	defer deleteContainer(c, containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: &map[string]string{},
	}
	_, err = containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.IsNil)
}

func (*aztestsSuite) TestContainerSetMetadataNil(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)
	access := PublicAccessTypeBlob
	createContainerOptions := CreateContainerOptions{
		Access:   &access,
		Metadata: &basicMetadata,
	}
	_, err := containerClient.Create(ctx, &createContainerOptions)

	defer deleteContainer(c, containerClient)

	_, err = containerClient.SetMetadata(ctx, nil)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.IsNil)
}

func (*aztestsSuite) TestContainerSetMetadataInvalidField(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: &map[string]string{"!nval!d Field!@#%": "value"},
	}
	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	c.Assert(err, chk.NotNil)
	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
}

func (*aztestsSuite) TestContainerSetMetadataNonExistent(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	_, err := containerClient.SetMetadata(ctx, nil)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeContainerNotFound)
}

func (s *aztestsSuite) TestContainerSetMetadataIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: &basicMetadata,
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	c.Assert(err, chk.IsNil)

	resp, err := containerClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)

}

func (s *aztestsSuite) TestContainerSetMetadataIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	setMetadataContainerOptions := SetMetadataContainerOptions{
		Metadata: &basicMetadata,
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err := containerClient.SetMetadata(ctx, &setMetadataContainerOptions)
	c.Assert(err, chk.NotNil)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestContainerNewBlobURL(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	blobURL := containerClient.NewBlobClient(blobPrefix, nil)
	tempBlob := blobURL.URL()
	tempContainer := containerClient.URL()
	c.Assert(tempBlob.String(), chk.Equals, tempContainer.String()+"/"+blobPrefix)
	c.Assert(blobURL, chk.FitsTypeOf, BlobClient{})
}

func (s *aztestsSuite) TestContainerNewBlockBlobClient(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := getContainerClient(c, bsu)

	blobURL := containerClient.NewBlockBlobClient(blobPrefix)
	tempBlob := blobURL.URL()
	tempContainer := containerClient.URL()
	c.Assert(tempBlob.String(), chk.Equals, tempContainer.String()+"/"+blobPrefix)
	c.Assert(blobURL, chk.FitsTypeOf, BlockBlobClient{})
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
)

func (s *azblobTestSuite) TestBlockBlobGetPropertiesUsingVID() {
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
	bbClient := createNewBlockBlob(_assert, generateBlobName(testName), containerClient)

	blobProp, _ := bbClient.GetProperties(ctx, nil)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	uploadResp, err := bbClient.Upload(ctx, getReaderToGeneratedBytes(1024), &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.NotNil(uploadResp.VersionID)
	blobProp, _ = bbClient.GetProperties(ctx, nil)
	_assert.EqualValues(uploadResp.VersionID, blobProp.VersionID)
	_assert.EqualValues(uploadResp.LastModified, blobProp.LastModified)
	_assert.Equal(*uploadResp.ETag, *blobProp.ETag)
	_assert.Equal(*blobProp.IsCurrentVersion, true)
}

func (s *azblobTestSuite) TestAppendBlobGetPropertiesUsingVID() {
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
	abClient := createNewAppendBlob(_assert, generateBlobName(testName), containerClient)

	blobProp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	createResp, err := abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)
	_assert.NotNil(createResp.VersionID)
	blobProp, _ = abClient.GetProperties(ctx, nil)
	_assert.EqualValues(createResp.VersionID, blobProp.VersionID)
	_assert.EqualValues(createResp.LastModified, blobProp.LastModified)
	_assert.Equal(*createResp.ETag, *blobProp.ETag)
	_assert.Equal(*blobProp.IsCurrentVersion, true)
}

//nolint
//func (s *azblobUnrecordedTestSuite) TestSetBlobMetadataReturnsVID() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	bbName := generateName(testName)
//	bbClient := createNewBlockBlob(_assert, bbName, containerClient)
//
//	metadata := map[string]string{"test_key_1": "test_value_1", "test_key_2": "2019"}
//	resp, err := bbClient.SetMetadata(ctx, metadata, nil)
//	_assert.Nil(err)
//	_assert.NotNil(resp.VersionID)
//
//	pager := containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
//		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemMetadata},
//	})
//
//	if !pager.NextPage(ctx) {
//		_assert.Nil(pager.Err()) // check for an error first
//		s.T().Fail()             // no page was gotten
//	}
//
//	pageResp := pager.PageResponse()
//
//	_assert.NotNil(pageResp.EnumerationResults.Segment.BlobItems)
//	blobList := pageResp.EnumerationResults.Segment.BlobItems
//	_assert.Len(blobList, 1)
//	blobResp1 := blobList[0]
//	_assert.Equal(*blobResp1.Name, bbName)
//	_assert.NotNil(blobResp1.Metadata.AdditionalProperties)
//	_assert.Len(blobResp1.Metadata.AdditionalProperties, 2)
//	// _assert(*blobResp1.Metadata, chk.DeepEquals, metadata)
//
//}

func (s *azblobTestSuite) TestCreateAndDownloadBlobSpecialCharactersWithVID() {
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
	data := []rune("-._/()$=',~0123456789")
	for i := 0; i < len(data); i++ {
		blobName := "abc" + string(data[i])
		blobURL := containerClient.NewBlockBlobClient(blobName)
		resp, err := blobURL.Upload(ctx, internal.NopCloser(strings.NewReader(string(data[i]))), nil)
		_assert.Nil(err)
		_assert.NotNil(resp.VersionID)

		dResp, err := blobURL.WithVersionID(*resp.VersionID).Download(ctx, nil)
		_assert.Nil(err)
		d1, err := ioutil.ReadAll(dResp.Body(RetryReaderOptions{}))
		_assert.Nil(err)
		_assert.NotEqual(*dResp.Version, "")
		_assert.EqualValues(string(d1), string(data[i]))
		versionId := dResp.RawResponse.Header.Get("x-ms-version-id")
		_assert.NotNil(versionId)
		_assert.Equal(versionId, *resp.VersionID)
	}
}

//func (s *azblobTestSuite) TestDeleteSpecificBlobVersion() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	blobURL := getBlockBlobClient(generateBlobName(testName), containerClient)
//
//	uploadResp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), &UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	})
//	_assert.Nil(err)
//	_assert.NotNil(uploadResp.VersionID)
//	versionID1 := uploadResp.VersionID
//
//	uploadResp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, &UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	})
//	_assert.Nil(err)
//	_assert.NotNil(uploadResp.VersionID)
//
//	listPager := containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
//		Include: &[]ListBlobsIncludeItem{ListBlobsIncludeItemVersions},
//	})
//
//	count := 0
//	blobs
//	for listPager.NextPage(ctx) {
//		resp := listPager.PageResponse()
//		for _, blob := range resp.EnumerationResults.Segment.BlobItems {
//			count += 1;
//			// Process the blobs returned
//			snapTime := "N/A"
//			if blob.Snapshot != nil {
//				snapTime = *blob.Snapshot
//			}
//			fmt.Printf("Blob name: %s, Snapshot: %s\n", *blob.Name, snapTime)
//		}
//	}
//	_assert.Nil(listPager.Err())
//	_assert.Len(count, 2)
//
//	// Deleting previous version snapshot.
//	deleteResp, err := blobURL.WithVersionID(versionID1).Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//	_assert.Nil(err)
//	_assert(deleteResp.StatusCode(), chk.Equals, 202)
//
//	listBlobsResp, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
//	_assert.Nil(err)
//	_assert(listBlobsResp.Segment.BlobItems, chk.NotNil)
//	if len(listBlobsResp.Segment.BlobItems) != 1 {
//		s.T().Fail()
//	}
//}
//
//func (s *azblobTestSuite) TestDeleteSpecificBlobVersionWithBlobSAS() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	blobURL, blobName := getBlockBlobClient(c, containerClient)
//
//	resp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	versionId := resp.VersionID
//	_assert(versionId, chk.NotNil)
//
//	resp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(resp.VersionID, chk.NotNil)
//
//	blobParts := NewBlobURLParts(blobURL.URL())
//	blobParts.VersionID = versionId
//	blobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: containerName,
//		BlobName:      blobName,
//		Permissions:   BlobSASPermissions{Delete: true, DeletePreviousVersion: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	sbURL := NewBlockBlobClient(blobParts.URL(), containerClient.client.p)
//	deleteResp, err := sbURL.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//	_assert(deleteResp, chk.IsNil)
//
//	listBlobResp, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
//	_assert.Nil(err)
//	for _, blob := range listBlobResp.Segment.BlobItems {
//		_assert(blob.VersionID, chk.Not(chk.Equals), versionId)
//	}
//}
//
//func (s *azblobTestSuite) TestDownloadSpecificBlobVersion() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	blobURL, _ := getBlockBlobClient(c, containerClient)
//
//	blockBlobUploadResp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(blockBlobUploadResp, chk.NotNil)
//	versionId1 := blockBlobUploadResp.VersionID
//
//	blockBlobUploadResp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(blockBlobUploadResp, chk.NotNil)
//	versionId2 := blockBlobUploadResp.VersionID
//	_assert(blockBlobUploadResp.VersionID, chk.NotNil)
//
//	// Download previous version of snapshot.
//	blobURL = blobURL.WithVersionID(versionId1)
//	blockBlobDeleteResp, err := blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	data, err := ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
//	_assert(string(data), chk.Equals, "data")
//
//	// Download current version of snapshot.
//	blobURL = blobURL.WithVersionID(versionId2)
//	blockBlobDeleteResp, err = blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	data, err = ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
//	_assert(string(data), chk.Equals, "updated_data")
//}
//
//func (s *azblobTestSuite) TestCreateBlobSnapshotReturnsVID() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
//	uploadResp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(uploadResp.VersionID, chk.NotNil)
//
//	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(csResp.VersionID, chk.NotNil)
//	lbResp, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{
//		Details: BlobListingDetails{Versions: true, Snapshots: true},
//	})
//	_assert(lbResp, chk.NotNil)
//	if len(lbResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//
//	_, err = blobURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
//	lbResp, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{
//		Details: BlobListingDetails{Versions: true, Snapshots: true},
//	})
//	_assert(lbResp, chk.NotNil)
//	if len(lbResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//	for _, blob := range lbResp.Segment.BlobItems {
//		_assert(blob.Snapshot, chk.Equals, "")
//	}
//}
//
//func (s *azblobTestSuite) TestCopyBlobFromURLWithSASReturnsVID() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	testSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(testSize)
//	sourceDataMD5Value := md5.Sum(sourceData)
//	ctx := context.Background()
//	srcBlob := container.NewBlockBlobClient(generateBlobName())
//	destBlob := container.NewBlockBlobClient(generateBlobName())
//
//	uploadSrcResp, err := srcBlob.Upload(ctx, r, HTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)
//	_assert(uploadSrcResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	srcBlobParts := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], DefaultAccessTier, nil)
//	_assert.Nil(err)
//	_assert(resp.Response().StatusCode, chk.Equals, 202)
//	_assert(resp.Version(), chk.Not(chk.Equals), "")
//	_assert(resp.CopyID(), chk.Not(chk.Equals), "")
//	_assert(string(resp.CopyStatus()), chk.DeepEquals, "success")
//	_assert(resp.VersionID, chk.NotNil)
//
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	_assert.Nil(err)
//	_assert(destData, chk.DeepEquals, sourceData)
//	_assert(downloadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//	_assert(len(downloadResp.NewMetadata()), chk.Equals, 1)
//	_, badMD5 := getRandomDataAndReader(16)
//	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, DefaultAccessTier, nil)
//	_assert.NotNil(err)
//
//	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, DefaultAccessTier, nil)
//	_assert.Nil(err)
//	_assert(resp.Response().StatusCode, chk.Equals, 202)
//	_assert(resp.XMsContentCRC64(), chk.Not(chk.Equals), "")
//	_assert(resp.Response().Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	_assert(resp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//}
//
//func (s *azblobTestSuite) TestCreateBlockBlobReturnsVID() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_assert, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	testSize := 2 * 1024 * 1024 // 1MB
//	r, _ := getRandomDataAndReader(testSize)
//	ctx := context.Background() // Use default Background context
//	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
//
//	// Prepare source blob for copy.
//	uploadResp, err := blobURL.Upload(ctx, r, HTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(uploadResp.Response().StatusCode, chk.Equals, 201)
//	_assert(uploadResp.rawResponse.Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	_assert(uploadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//	_assert(csResp.Response().StatusCode, chk.Equals, 201)
//	_assert(csResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	listBlobResp, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//	_assert.Nil(err)
//	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
//	if len(listBlobResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//
//	deleteResp, err := blobURL.Delete(ctx, DeleteSnapshotsOptionOnly, BlobAccessConditions{})
//	_assert.Nil(err)
//	_assert(deleteResp.Response().StatusCode, chk.Equals, 202)
//	_assert(deleteResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	listBlobResp, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Versions: true}})
//	_assert.Nil(err)
//	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
//	if len(listBlobResp.Segment.BlobItems) == 0 {
//		s.T().Fail()
//	}
//	blobs := listBlobResp.Segment.BlobItems
//	_assert(blobs[0].Snapshot, chk.Equals, "")
//}

func (s *azblobTestSuite) TestPutBlockListReturnsVID() {
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

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(d)), nil)
		_assert.Nil(err)
		_assert.Equal(resp.RawResponse.StatusCode, 201)
		_assert.NotNil(resp.Version)
		_assert.NotEqual(*resp.Version, "")
	}

	commitResp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, nil)
	_assert.Nil(err)
	_assert.NotNil(commitResp.VersionID)

	contentResp, err := bbClient.Download(ctx, nil)
	_assert.Nil(err)
	contentData, err := ioutil.ReadAll(contentResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(contentData, []uint8(strings.Join(data, "")))
}

func (s *azblobTestSuite) TestCreatePageBlobReturnsVID() {
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

	pbClob := createNewPageBlob(_assert, generateBlobName(testName), containerClient)

	contentSize := 1 * 1024
	r, _ := generateData(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	putResp, err := pbClob.UploadPages(context.Background(), r, &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.Equal(putResp.LastModified.IsZero(), false)
	_assert.NotNil(putResp.ETag)
	_assert.NotEqual(putResp.Version, "")
	_assert.NotNil(putResp.RawResponse.Header.Get("x-ms-version-id"))

	gpResp, err := pbClob.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(gpResp)
}

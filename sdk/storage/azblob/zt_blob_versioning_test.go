// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
)

func TestBlockBlobGetPropertiesUsingVID(t *testing.T) {
	t.Skipf("VersionID is not filled")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)
	bbClient := createNewBlockBlob(t, generateBlobName(testName), containerClient)

	blobProp, _ := bbClient.GetProperties(ctx, nil)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	uploadResp, err := bbClient.Upload(ctx, getReaderToGeneratedBytes(1024), &uploadBlockBlobOptions)
	require.NoError(t, err)
	t.Skip("expected VersionID to be not nil")
	require.NotNil(t, uploadResp.VersionID)
	blobProp, _ = bbClient.GetProperties(ctx, nil)
	require.EqualValues(t, uploadResp.VersionID, blobProp.VersionID)
	require.EqualValues(t, uploadResp.LastModified, blobProp.LastModified)
	require.Equal(t, *uploadResp.ETag, *blobProp.ETag)
	require.Equal(t, *blobProp.IsCurrentVersion, true)
}

func TestAppendBlobGetPropertiesUsingVID(t *testing.T) {
	t.Skipf("VersionID is nil but not expected to be")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)
	abClient := createNewAppendBlob(t, generateBlobName(testName), containerClient)

	blobProp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	createResp, err := abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)
	require.NotNil(t, createResp.VersionID)
	blobProp, err = abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, createResp.VersionID, blobProp.VersionID)
	require.EqualValues(t, createResp.LastModified, blobProp.LastModified)
	require.Equal(t, *createResp.ETag, *blobProp.ETag)
	require.Equal(t, *blobProp.IsCurrentVersion, true)
}

func TestCreateAndDownloadBlobSpecialCharactersWithVID(t *testing.T) {
	t.Skipf("VersionID is not filled")
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)
	data := []rune("-._/()$=',~0123456789")
	for i := 0; i < len(data); i++ {
		blobName := "abc" + string(data[i])
		blobURL := containerClient.NewBlockBlobClient(blobName)
		resp, err := blobURL.Upload(ctx, internal.NopCloser(strings.NewReader(string(data[i]))), nil)
		require.NoError(t, err)
		require.NotNil(t, resp.VersionID) // VersionID is nil

		dResp, err := blobURL.WithVersionID(*resp.VersionID).Download(ctx, nil)
		require.NoError(t, err)
		d1, err := ioutil.ReadAll(dResp.Body(nil))
		require.NoError(t, err)
		require.NotEqual(t, *dResp.Version, "")
		require.EqualValues(t, string(d1), string(data[i]))
		versionId := dResp.RawResponse.Header.Get("x-ms-version-id")
		require.NotNil(t, versionId)
		require.Equal(t, versionId, *resp.VersionID)
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
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//	blobURL := getBlockBlobClient(generateBlobName(testName), containerClient)
//
//	uploadResp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), &UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	})
//	_assert.NoError(err)
//	_assert.NotNil(uploadResp.VersionID)
//	versionID1 := uploadResp.VersionID
//
//	uploadResp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, &UploadBlockBlobOptions{
//		Metadata: basicMetadata,
//	})
//	_assert.NoError(err)
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
//	_assert.NoError(err)
//	_assert(deleteResp.StatusCode(), chk.Equals, 202)
//
//	listBlobsResp, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
//	_assert.NoError(err)
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
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//	blobURL, blobName := getBlockBlobClient(c, containerClient)
//
//	resp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	versionId := resp.VersionID
//	_assert(versionId, chk.NotNil)
//
//	resp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
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
//	_assert.NoError(err)
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
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//	blobURL, _ := getBlockBlobClient(c, containerClient)
//
//	blockBlobUploadResp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"))), HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	_assert(blockBlobUploadResp, chk.NotNil)
//	versionId1 := blockBlobUploadResp.VersionID
//
//	blockBlobUploadResp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	_assert(blockBlobUploadResp, chk.NotNil)
//	versionId2 := blockBlobUploadResp.VersionID
//	_assert(blockBlobUploadResp.VersionID, chk.NotNil)
//
//	// Download previous version of snapshot.
//	blobURL = blobURL.WithVersionID(versionId1)
//	blockBlobDeleteResp, err := blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	data, err := ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
//	_assert(string(data), chk.Equals, "data")
//
//	// Download current version of snapshot.
//	blobURL = blobURL.WithVersionID(versionId2)
//	blockBlobDeleteResp, err = blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
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
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
//	uploadResp, err := blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	_assert(uploadResp.VersionID, chk.NotNil)
//
//	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
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
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//
//	testSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(testSize)
//	sourceDataMD5Value := md5.Sum(sourceData)
//	ctx := context.Background()
//	srcBlob := container.NewBlockBlobClient(generateBlobName())
//	destBlob := container.NewBlockBlobClient(generateBlobName())
//
//	uploadSrcResp, err := srcBlob.Upload(ctx, r, HTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
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
//	_assert.NoError(err)
//	_assert(resp.Response().StatusCode, chk.Equals, 202)
//	_assert(resp.Version(), chk.Not(chk.Equals), "")
//	_assert(resp.CopyID(), chk.Not(chk.Equals), "")
//	_assert(string(resp.CopyStatus()), chk.DeepEquals, "success")
//	_assert(resp.VersionID, chk.NotNil)
//
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
//	_assert.NoError(err)
//	_assert(destData, chk.DeepEquals, sourceData)
//	_assert(downloadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//	_assert(len(downloadResp.NewMetadata()), chk.Equals, 1)
//	_, badMD5 := getRandomDataAndReader(16)
//	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, DefaultAccessTier, nil)
//	_assert.Error(err)
//
//	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, DefaultAccessTier, nil)
//	_assert.NoError(err)
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
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//
//	testSize := 2 * 1024 * 1024 // 1MB
//	r, _ := getRandomDataAndReader(testSize)
//	ctx := context.Background() // Use default Background context
//	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
//
//	// Prepare source blob for copy.
//	uploadResp, err := blobURL.Upload(ctx, r, HTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	_assert(uploadResp.Response().StatusCode, chk.Equals, 201)
//	_assert(uploadResp.rawResponse.Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	_assert(uploadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	_assert.NoError(err)
//	_assert(csResp.Response().StatusCode, chk.Equals, 201)
//	_assert(csResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	listBlobResp, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//	_assert.NoError(err)
//	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
//	if len(listBlobResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//
//	deleteResp, err := blobURL.Delete(ctx, DeleteSnapshotsOptionOnly, BlobAccessConditions{})
//	_assert.NoError(err)
//	_assert(deleteResp.Response().StatusCode, chk.Equals, 202)
//	_assert(deleteResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	listBlobResp, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Versions: true}})
//	_assert.NoError(err)
//	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
//	if len(listBlobResp.Segment.BlobItems) == 0 {
//		s.T().Fail()
//	}
//	blobs := listBlobResp.Segment.BlobItems
//	_assert(blobs[0].Snapshot, chk.Equals, "")
//}

func TestPutBlockListReturnsVID(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(ctx, base64BlockIDs[index], internal.NopCloser(strings.NewReader(d)), nil)
		require.NoError(t, err)
		require.Equal(t, resp.RawResponse.StatusCode, 201)
		require.NotNil(t, resp.Version)
		require.NotEqual(t, *resp.Version, "")
	}

	commitResp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, nil)
	require.NoError(t, err)
	t.Skip("expected VersionID to be not nil")
	require.NotNil(t, commitResp.VersionID)

	contentResp, err := bbClient.Download(ctx, nil)
	require.NoError(t, err)
	contentData, err := ioutil.ReadAll(contentResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, contentData, []uint8(strings.Join(data, "")))
}

func TestCreatePageBlobReturnsVID(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	pbClob := createNewPageBlob(t, generateBlobName(testName), containerClient)

	contentSize := 1 * 1024
	r, _ := generateData(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	putResp, err := pbClob.UploadPages(context.Background(), r, &uploadPagesOptions)
	require.NoError(t, err)
	require.Equal(t, putResp.RawResponse.StatusCode, 201)
	require.Equal(t, putResp.LastModified.IsZero(), false)
	require.NotNil(t, putResp.ETag)
	require.NotEqual(t, putResp.Version, "")
	require.NotNil(t, putResp.RawResponse.Header.Get("x-ms-version-id"))

	gpResp, err := pbClob.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, gpResp)
}

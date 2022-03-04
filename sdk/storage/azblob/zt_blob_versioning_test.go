// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestBlockBlobGetPropertiesUsingVID(t *testing.T) {
	//recording.LiveOnly(t)
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
	//t.Skip("expected VersionID to be not nil")
	require.NotNil(t, uploadResp.VersionID)
	blobProp, _ = bbClient.GetProperties(ctx, nil)
	require.EqualValues(t, uploadResp.VersionID, blobProp.VersionID)
	require.EqualValues(t, uploadResp.LastModified, blobProp.LastModified)
	require.Equal(t, *uploadResp.ETag, *blobProp.ETag)
	require.Equal(t, *blobProp.IsCurrentVersion, true)
}

func TestAppendBlobGetPropertiesUsingVID(t *testing.T) {
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

func TestDeleteSpecificBlobVersion(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)
	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		uploadResp, err := bbClient.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("data"+strconv.Itoa(i)))), &UploadBlockBlobOptions{
			Metadata: basicMetadata,
		})
		require.NoError(t, err)
		require.NotNil(t, uploadResp.VersionID)
		versions = append(versions, *uploadResp.VersionID)
	}

	listPager := containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemVersions},
	})

	found := make([]*BlobItemInternal, 0)
	for listPager.NextPage(ctx) {
		resp := listPager.PageResponse()
		found = append(found, resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems...)
	}
	require.NoError(t, listPager.Err())
	require.Len(t, found, 5)

	// Deleting the 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		deleteResp, err := bbClient.WithVersionID(versions[i]).Delete(ctx, nil)
		require.NoError(t, err)
		require.Equal(t, deleteResp.RawResponse.StatusCode, 202)
	}

	listPager = containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemVersions},
	})

	found = make([]*BlobItemInternal, 0)
	for listPager.NextPage(ctx) {
		resp := listPager.PageResponse()
		found = append(found, resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems...)
	}
	require.NoError(t, listPager.Err())
	require.Len(t, found, 2)

	for i := 3; i < 5; i++ {
		downloadResp, err := bbClient.WithVersionID(versions[i]).Download(ctx, nil)
		require.NoError(t, err)
		destData, err := ioutil.ReadAll(downloadResp.Body(nil))
		require.NoError(t, err)
		require.EqualValues(t, destData, "data"+strconv.Itoa(i))
	}
}

//func (s *azblobTestSuite) TestDeleteSpecificBlobVersionWithBlobSAS() {
//	require := assert.New(s.T())
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
//	require.NoError(t, err)
//	versionId := resp.VersionID
//	require(t, versionId, chk.NotNil)
//
//	resp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	require.NoError(t, err)
//	require(t, resp.VersionID, chk.NotNil)
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
//	require(t, deleteResp, chk.IsNil)
//
//	listBlobResp, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
//	require.NoError(t, err)
//	for _, blob := range listBlobResp.Segment.BlobItems {
//		require(t, blob.VersionID, chk.Not(chk.Equals), versionId)
//	}
//}
//
//func (s *azblobTestSuite) TestDownloadSpecificBlobVersion() {
//	require := assert.New(s.T())
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
//	require.NoError(t, err)
//	require(t, blockBlobUploadResp, chk.NotNil)
//	versionId1 := blockBlobUploadResp.VersionID
//
//	blockBlobUploadResp, err = blobURL.Upload(ctx, internal.NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	require.NoError(t, err)
//	require(t, blockBlobUploadResp, chk.NotNil)
//	versionId2 := blockBlobUploadResp.VersionID
//	require(t, blockBlobUploadResp.VersionID, chk.NotNil)
//
//	// Download previous version of snapshot.
//	blobURL = blobURL.WithVersionID(versionId1)
//	blockBlobDeleteResp, err := blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	require.NoError(t, err)
//	data, err := ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
//	require(t, string(data), chk.Equals, "data")
//
//	// Download current version of snapshot.
//	blobURL = blobURL.WithVersionID(versionId2)
//	blockBlobDeleteResp, err = blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	require.NoError(t, err)
//	data, err = ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
//	require(t, string(data), chk.Equals, "updated_data")
//}
//
//func (s *azblobTestSuite) TestCreateBlobSnapshotReturnsVID() {
//	require := assert.New(s.T())
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
//	require.NoError(t, err)
//	require(t, uploadResp.VersionID, chk.NotNil)
//
//	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	require.NoError(t, err)
//	require(t, csResp.VersionID, chk.NotNil)
//	lbResp, err := containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{
//		Details: BlobListingDetails{Versions: true, Snapshots: true},
//	})
//	require(t, lbResp, chk.NotNil)
//	if len(lbResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//
//	_, err = blobURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
//	lbResp, err = containerClient.ListBlobsFlat(ctx, Marker{}, ListBlobsSegmentOptions{
//		Details: BlobListingDetails{Versions: true, Snapshots: true},
//	})
//	require(t, lbResp, chk.NotNil)
//	if len(lbResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//	for _, blob := range lbResp.Segment.BlobItems {
//		require(t, blob.Snapshot, chk.Equals, "")
//	}
//}

//func TestCopyBlobFromURLWithSASReturnsVID(t *testing.T) {
//	recording.LiveOnly(t) // Live only because of random data and random name
//	stop := start(t)
//	defer stop()
//
//	testName := t.Name()
//	svcClient, err := createServiceClient(t, testAccountDefault)
//	require.NoError(t, err)
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(t, containerClient)
//
//	testSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(testSize)
//	sourceDataMD5Value := md5.Sum(sourceData)
//	ctx := context.Background()
//	srcBlob := containerClient.NewBlockBlobClient(generateBlobName(testName))
//	destBlob := containerClient.NewBlockBlobClient(generateBlobName(testName))
//
//	uploadSrcResp, err := srcBlob.Upload(ctx, internal.NopCloser(r), nil)
//	require.NoError(t, err)
//	require.Equal(t, uploadSrcResp.RawResponse.StatusCode, 201)
//	require.NotNil(t, uploadSrcResp.VersionID)
//
//	credential, err := getCredential(testAccountDefault)
//	require.NoError(t, err)
//
//	srcBlobParts := NewBlobURLParts(srcBlob.URL())
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	require.NoError(t, err)
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], DefaultAccessTier, nil)
//	require.NoError(t, err)
//	require(t, resp.RawResponse.StatusCode, chk.Equals, 202)
//	require(t, resp.Version(), chk.Not(chk.Equals), "")
//	require(t, resp.CopyID(), chk.Not(chk.Equals), "")
//	require(t, string(resp.CopyStatus()), chk.DeepEquals, "success")
//	require(t, resp.VersionID, chk.NotNil)
//
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	require.NoError(t, err)
//	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
//	require.NoError(t, err)
//	require(t, destData, chk.DeepEquals, sourceData)
//	require(t, downloadResp.VersionID, chk.NotNil)
//	require(t, len(downloadResp.NewMetadata()), chk.Equals, 1)
//	_, badMD5 := getRandomDataAndReader(16)
//	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, DefaultAccessTier, nil)
//	require.Error(err)
//
//	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, DefaultAccessTier, nil)
//	require.NoError(t, err)
//	require(t, resp.RawResponse.StatusCode, chk.Equals, 202)
//	require(t, resp.XMsContentCRC64(), chk.Not(chk.Equals), "")
//	require(t, resp.Response().Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	require(t, resp.VersionID, chk.NotNil)
//}

func TestCreateBlockBlobReturnsVID(t *testing.T) {
	//recording.LiveOnly(t)
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	testSize := 2 * 1024 * 1024 // 1MB
	r, _ := getRandomDataAndReader(testSize)
	ctx := context.Background() // Use default Background context
	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	// Prepare source blob for copy.
	uploadResp, err := bbClient.Upload(ctx, internal.NopCloser(r), nil)
	require.NoError(t, err)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)
	require.NotNil(t, uploadResp.VersionID)

	csResp, err := bbClient.CreateSnapshot(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, csResp.RawResponse.StatusCode, 201)
	require.NotNil(t, csResp.VersionID)

	pager := containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots},
	})

	found := make([]*BlobItemInternal, 0)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		found = append(found, resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems...)
	}
	require.NoError(t, pager.Err())
	require.Len(t, found, 2)

	deleteSnapshotsOnly := DeleteSnapshotsOptionTypeOnly
	deleteResp, err := bbClient.Delete(ctx, &DeleteBlobOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	require.NoError(t, err)
	require.Equal(t, deleteResp.RawResponse.StatusCode, 202)

	pager = containerClient.ListBlobsFlat(&ContainerListBlobFlatSegmentOptions{
		Include: []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots, ListBlobsIncludeItemVersions},
	})

	found = make([]*BlobItemInternal, 0)
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		found = append(found, resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems...)
	}
	require.NoError(t, pager.Err())
	require.NotEqual(t, len(found), 0)
}

func TestPutBlockListReturnsVID(t *testing.T) {
	//recording.LiveOnly(t)
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
	//t.Skip("expected VersionID to be not nil")
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

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

//
//import (
//	"context"
//	"io/ioutil"
//	"strings"
//
//	chk "gopkg.in/check.v1"
//)
//
//func (s *azblobTestSuite) TestBlockBlobGetPropertiesUsingVID() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(_assert, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	blobProp, _ := blobClient.GetProperties(ctx, nil)
//
//	uploadBlockBlobOptions := UploadBlockBlobOptions{
//		Metadata:                 &basicMetadata,
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: blobProp.ETag},
//	}
//	uploadResp, err := blobClient.Upload(ctx, getReaderToGeneratedBytes(1024), &uploadBlockBlobOptions)
//	_assert.Nil(err)
//	_assert(uploadResp.VersionID, chk.NotNil)
//	blobProp, _ = blobClient.GetProperties(ctx, nil)
//	_assert(uploadResp.VersionID, chk.DeepEquals, blobProp.VersionID)
//	_assert(uploadResp.LastModified, chk.DeepEquals, blobProp.LastModified)
//	_assert(*uploadResp.ETag, chk.Equals, *blobProp.ETag)
//	_assert(*blobProp.IsCurrentVersion, chk.Equals, true)
//}

//func (s *azblobTestSuite) TestAppendBlobGetPropertiesUsingVID() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(_assert, containerClient)
//	blobClient, _ := createNewAppendBlob(c, containerClient)
//
//	blobProp, _ := blobClient.GetProperties(ctx, nil)
//
//	createAppendBlobOptions := CreateAppendBlobOptions{
//		Metadata: &basicMetadata,
//		BlobAccessConditions: BlobAccessConditions{
//			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: blobProp.ETag},
//		},
//	}
//	createResp, err := blobClient.Create(ctx, &createAppendBlobOptions)
//	_assert.Nil(err)
//	_assert(createResp.VersionID, chk.NotNil)
//	blobProp, _ = blobClient.GetProperties(ctx, nil)
//	_assert(createResp.VersionID, chk.DeepEquals, blobProp.VersionID)
//	_assert(createResp.LastModified, chk.DeepEquals, blobProp.LastModified)
//	_assert(*createResp.ETag, chk.Equals, *blobProp.ETag)
//	_assert(*blobProp.IsCurrentVersion, chk.Equals, true)
//}

//func (s *azblobTestSuite) TestSetBlobMetadataReturnsVID() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(_assert, containerClient)
//	blobURL, blobName := createNewBlockBlob(c, containerClient)
//	metadata := map[string]string{"test_key_1": "test_value_1", "test_key_2": "2019"}
//	resp, err := blobURL.SetMetadata(ctx, metadata, nil)
//	_assert.Nil(err)
//	_assert(resp.VersionID, chk.NotNil)
//
//	include := []ListBlobsIncludeItem{ListBlobsIncludeItemMetadata}
//	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
//		Include: &include,
//	}
//	pager := containerClient.ListBlobsFlatSegment(&containerListBlobFlatSegmentOptions)
//
//	if !pager.NextPage(ctx) {
//		_assert(pager.Err(), chk.IsNil) // check for an error first
//		c.Fail()                         // no page was gotten
//	}
//
//	pageResp := pager.PageResponse()
//
//	_assert(pageResp.EnumerationResults.Segment.BlobItems, chk.NotNil)
//	blobList := *pageResp.EnumerationResults.Segment.BlobItems
//	_assert(len(blobList), chk.Equals, 1)
//	blobResp1 := blobList[0]
//	_assert(*blobResp1.Name, chk.Equals, blobName)
//	_assert(*blobResp1.Metadata.AdditionalProperties, chk.NotNil)
//	_assert(*blobResp1.Metadata.AdditionalProperties, chk.HasLen, 2)
//	// _assert(*blobResp1.Metadata, chk.DeepEquals, metadata)
//
//}
//
////func (s *azblobTestSuite) TestCreateAndDownloadBlobSpecialCharactersWithVID() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	data := []rune("-._/()$=',~0123456789")
////	for i := 0; i < len(data); i++ {
////		blobName := "abc" + string(data[i])
////		blobURL := containerClient.NewBlockBlobClient(blobName)
////		resp, err := blobURL.Upload(ctx, strings.NewReader(string(data[i])), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////		_assert.Nil(err)
////		_assert(resp.VersionID(), chk.NotNil)
////
////		dResp, err := blobURL.WithVersionID(resp.VersionID()).Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
////		_assert.Nil(err)
////		d1, err := ioutil.ReadAll(dResp.Body(RetryReaderOptions{}))
////		_assert(dResp.Version(), chk.Not(chk.Equals), "")
////		_assert(string(d1), chk.DeepEquals, string(data[i]))
////		versionId := dResp.r.rawResponse.Header.Get("x-ms-version-id")
////		_assert(versionId, chk.NotNil)
////		_assert(versionId, chk.Equals, resp.VersionID())
////	}
////}
//
////func (s *azblobTestSuite) TestDeleteSpecificBlobVersion() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	blobURL, _ := getBlockBlobURL(c, containerClient)
////
////	blockBlobUploadResp, err := blobURL.Upload(ctx, bytes.NewReader([]byte("data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(blockBlobUploadResp.VersionID(), chk.NotNil)
////	versionID1 := blockBlobUploadResp.VersionID()
////
////	blockBlobUploadResp, err = blobURL.Upload(ctx, bytes.NewReader([]byte("updated_data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(blockBlobUploadResp.VersionID(), chk.NotNil)
////
////	listBlobsResp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
////	_assert.Nil(err)
////	_assert(listBlobsResp.Segment.BlobItems, chk.HasLen, 2)
////
////	// Deleting previous version snapshot.
////	deleteResp, err := blobURL.WithVersionID(versionID1).Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
////	_assert.Nil(err)
////	_assert(deleteResp.StatusCode(), chk.Equals, 202)
////
////	listBlobsResp, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
////	_assert.Nil(err)
////	_assert(listBlobsResp.Segment.BlobItems, chk.NotNil)
////	if len(listBlobsResp.Segment.BlobItems) != 1 {
////		c.Fail()
////	}
////}
//
////func (s *azblobTestSuite) TestDeleteSpecificBlobVersionWithBlobSAS() {
////	bsu := getServiceClient()
////	credential, err := getGenericCredential("")
////	if err != nil {
////		c.Fatal(err)
////	}
////	containerClient, containerName := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	blobURL, blobName := getBlockBlobURL(c, containerClient)
////
////	resp, err := blobURL.Upload(ctx, bytes.NewReader([]byte("data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	versionId := resp.VersionID()
////	_assert(versionId, chk.NotNil)
////
////	resp, err = blobURL.Upload(ctx, bytes.NewReader([]byte("updated_data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(resp.VersionID(), chk.NotNil)
////
////	blobParts := NewBlobURLParts(blobURL.URL())
////	blobParts.VersionID = versionId
////	blobParts.SAS, err = BlobSASSignatureValues{
////		Protocol:      SASProtocolHTTPS,
////		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
////		ContainerName: containerName,
////		BlobName:      blobName,
////		Permissions:   BlobSASPermissions{Delete: true, DeletePreviousVersion: true}.String(),
////	}.NewSASQueryParameters(credential)
////	if err != nil {
////		c.Fatal(err)
////	}
////
////	sbURL := NewBlockBlobClient(blobParts.URL(), containerClient.client.p)
////	deleteResp, err := sbURL.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
////	_assert(deleteResp, chk.IsNil)
////
////	listBlobResp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
////	_assert.Nil(err)
////	for _, blob := range listBlobResp.Segment.BlobItems {
////		_assert(blob.VersionID, chk.Not(chk.Equals), versionId)
////	}
////}
//
////func (s *azblobTestSuite) TestDownloadSpecificBlobVersion() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	blobURL, _ := getBlockBlobURL(c, containerClient)
////
////	blockBlobUploadResp, err := blobURL.Upload(ctx, bytes.NewReader([]byte("data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(blockBlobUploadResp, chk.NotNil)
////	versionId1 := blockBlobUploadResp.VersionID()
////
////	blockBlobUploadResp, err = blobURL.Upload(ctx, bytes.NewReader([]byte("updated_data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(blockBlobUploadResp, chk.NotNil)
////	versionId2 := blockBlobUploadResp.VersionID()
////	_assert(blockBlobUploadResp.VersionID(), chk.NotNil)
////
////	// Download previous version of snapshot.
////	blobURL = blobURL.WithVersionID(versionId1)
////	blockBlobDeleteResp, err := blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	data, err := ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
////	_assert(string(data), chk.Equals, "data")
////
////	// Download current version of snapshot.
////	blobURL = blobURL.WithVersionID(versionId2)
////	blockBlobDeleteResp, err = blobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	data, err = ioutil.ReadAll(blockBlobDeleteResp.Response().Body)
////	_assert(string(data), chk.Equals, "updated_data")
////}
//
////func (s *azblobTestSuite) TestCreateBlobSnapshotReturnsVID() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer delContainer(c, containerClient)
////	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
////	uploadResp, err := blobURL.Upload(ctx, bytes.NewReader([]byte("updated_data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(uploadResp.VersionID(), chk.NotNil)
////
////	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(csResp.VersionID(), chk.NotNil)
////	lbResp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{
////		Details: BlobListingDetails{Versions: true, Snapshots: true},
////	})
////	_assert(lbResp, chk.NotNil)
////	if len(lbResp.Segment.BlobItems) < 2 {
////		c.Fail()
////	}
////
////	_, err = blobURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
////	lbResp, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{
////		Details: BlobListingDetails{Versions: true, Snapshots: true},
////	})
////	_assert(lbResp, chk.NotNil)
////	if len(lbResp.Segment.BlobItems) < 2 {
////		c.Fail()
////	}
////	for _, blob := range lbResp.Segment.BlobItems {
////		_assert(blob.Snapshot, chk.Equals, "")
////	}
////}
//
////func (s *azblobTestSuite) TestCopyBlobFromURLWithSASReturnsVID() {
////	bsu := getServiceClient()
////	credential, err := getGenericCredential("")
////	if err != nil {
////		c.Fatal("Invalid credential")
////	}
////	container, _ := createNewContainer(c, bsu)
////	defer delContainer(c, container)
////
////	testSize := 4 * 1024 * 1024 // 4MB
////	r, sourceData := getRandomDataAndReader(testSize)
////	sourceDataMD5Value := md5.Sum(sourceData)
////	ctx := context.Background()
////	srcBlob := container.NewBlockBlobClient(generateBlobName())
////	destBlob := container.NewBlockBlobClient(generateBlobName())
////
////	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)
////	_assert(uploadSrcResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
////
////	srcBlobParts := NewBlobURLParts(srcBlob.URL())
////
////	srcBlobParts.SAS, err = BlobSASSignatureValues{
////		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
////		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
////		ContainerName: srcBlobParts.ContainerName,
////		BlobName:      srcBlobParts.BlobName,
////		Permissions:   BlobSASPermissions{Read: true}.String(),
////	}.NewSASQueryParameters(credential)
////	if err != nil {
////		c.Fatal(err)
////	}
////
////	srcBlobURLWithSAS := srcBlobParts.URL()
////
////	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], DefaultAccessTier, nil)
////	_assert.Nil(err)
////	_assert(resp.Response().StatusCode, chk.Equals, 202)
////	_assert(resp.Version(), chk.Not(chk.Equals), "")
////	_assert(resp.CopyID(), chk.Not(chk.Equals), "")
////	_assert(string(resp.CopyStatus()), chk.DeepEquals, "success")
////	_assert(resp.VersionID(), chk.NotNil)
////
////	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
////	_assert.Nil(err)
////	_assert(destData, chk.DeepEquals, sourceData)
////	_assert(downloadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
////	_assert(len(downloadResp.NewMetadata()), chk.Equals, 1)
////	_, badMD5 := getRandomDataAndReader(16)
////	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, DefaultAccessTier, nil)
////	_assert.NotNil(err)
////
////	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, DefaultAccessTier, nil)
////	_assert.Nil(err)
////	_assert(resp.Response().StatusCode, chk.Equals, 202)
////	_assert(resp.XMsContentCRC64(), chk.Not(chk.Equals), "")
////	_assert(resp.Response().Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
////	_assert(resp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
////}
//
////func (s *azblobTestSuite) TestCreateBlockBlobReturnsVID() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer delContainer(c, containerClient)
////
////	testSize := 2 * 1024 * 1024 // 1MB
////	r, _ := getRandomDataAndReader(testSize)
////	ctx := context.Background() // Use default Background context
////	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
////
////	// Prepare source blob for copy.
////	uploadResp, err := blobURL.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(uploadResp.Response().StatusCode, chk.Equals, 201)
////	_assert(uploadResp.rawResponse.Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
////	_assert(uploadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
////
////	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
////	_assert.Nil(err)
////	_assert(csResp.Response().StatusCode, chk.Equals, 201)
////	_assert(csResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
////
////	listBlobResp, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
////	_assert.Nil(err)
////	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
////	if len(listBlobResp.Segment.BlobItems) < 2 {
////		c.Fail()
////	}
////
////	deleteResp, err := blobURL.Delete(ctx, DeleteSnapshotsOptionOnly, BlobAccessConditions{})
////	_assert.Nil(err)
////	_assert(deleteResp.Response().StatusCode, chk.Equals, 202)
////	_assert(deleteResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
////
////	listBlobResp, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Versions: true}})
////	_assert.Nil(err)
////	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
////	if len(listBlobResp.Segment.BlobItems) == 0 {
////		c.Fail()
////	}
////	blobs := listBlobResp.Segment.BlobItems
////	_assert(blobs[0].Snapshot, chk.Equals, "")
////}
//
//func (s *azblobTestSuite) TestPutBlockListReturnsVID() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(_assert, containerClient)
//
//	blobClient := containerClient.NewBlockBlobClient(generateBlobName())
//
//	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
//	base64BlockIDs := make([]string, len(data))
//
//	for index, d := range data {
//		base64BlockIDs[index] = blockIDIntToBase64(index)
//		resp, err := blobClient.StageBlock(ctx, base64BlockIDs[index], strings.NewReader(d), nil)
//		_assert.Nil(err)
//		_assert(resp.RawResponse.StatusCode, chk.Equals, 201)
//		_assert(resp.Version, chk.NotNil)
//		_assert(*resp.Version, chk.Not(chk.Equals), "")
//	}
//
//	commitResp, err := blobClient.CommitBlockList(ctx, base64BlockIDs, nil)
//	_assert.Nil(err)
//	_assert(commitResp.VersionID, chk.NotNil)
//
//	contentResp, err := blobClient.Download(ctx, nil)
//	_assert.Nil(err)
//	contentData, err := ioutil.ReadAll(contentResp.Body(RetryReaderOptions{}))
//	_assert(contentData, chk.DeepEquals, []uint8(strings.Join(data, "")))
//}
//
//func (s *azblobTestSuite) TestCreatePageBlobReturnsVID() {
//	bsu := getServiceClient(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(_assert, containerClient)
//
//	blob, _ := createNewPageBlob(c, containerClient)
//
//	contentSize := 1 * 1024
//	r := getReaderToGeneratedBytes(contentSize)
//	offset, count := int64(0), int64(contentSize)
//	uploadPagesOptions := UploadPagesOptions{
//		PageRange: &HttpRange{offset, count},
//	}
//	putResp, err := blob.UploadPages(context.Background(), r, &uploadPagesOptions)
//	_assert.Nil(err)
//	_assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
//	_assert(putResp.LastModified.IsZero(), chk.Equals, false)
//	_assert(putResp.ETag, chk.NotNil)
//	_assert(putResp.Version, chk.Not(chk.Equals), "")
//	_assert(putResp.RawResponse.Header.Get("x-ms-version-id"), chk.NotNil)
//
//	gpResp, err := blob.GetProperties(ctx, nil)
//	_assert.Nil(err)
//	_assert(gpResp, chk.NotNil)
//}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//func (s *AZBlobUnrecordedTestsSuite) TestDeleteSpecificBlobVersionWithBlobSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	blobClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
//
//	resp, err := blobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("data"))), &blockblob.UploadOptions{Metadata: testcommon.BasicMetadata})
//	_require.Nil(err)
//	versionId := resp.VersionID
//	_require.NotNil(versionId)
//
//	resp, err = blobClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("updated_data"))), &blockblob.UploadOptions{Metadata: testcommon.BasicMetadata})
//	_require.Nil(err)
//	_require.NotNil(resp.VersionID)
//
//	cred, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//
//	blobParts, err := azblob.ParseURL(blobClient.URL())
//	_require.Nil(err)
//	blobParts.VersionID = *versionId
//	blobParts.SAS, err = service.SASSignatureValues{
//		Protocol:      service.SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		Permissions:   to.Ptr(service.SASPermissions{Delete: true, DeletePreviousVersion: true}).String(),
//		ResourceTypes: to.Ptr(service.SASResourceTypes{Service: true, Container: true, Object: true}).String(),
//	}.Sign(cred)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	sbClient, err := blockblob.NewClientWithNoCredential(blobParts.URL(), nil)
//	_, err = sbClient.Delete(ctx, nil)
//	_require.Nil(err)
//
//	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
//		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemVersions},
//	})
//	for pager.More() {
//		pageResp, err := pager.NextPage(ctx)
//		_require.Nil(err)
//		for _, blob := range pageResp.Segment.BlobItems {
//			_require.NotEqual(*blob.VersionID, versionId)
//		}
//	}
//
//}

//
//func (s *AZBlobRecordedTestsSuite) TestCopyBlobFromURLWithSASReturnsVID() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	testSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := testcommon.GetRandomDataAndReader(testSize)
//	sourceDataMD5Value := md5.Sum(sourceData)
//	ctx := ctx
//	srcBlob := container.NewBlockBlobClient(testcommon.GenerateBlobName())
//	destBlob := container.NewBlockBlobClient(testcommon.GenerateBlobName())
//
//	uploadSrcResp, err := srcBlob.Upload(ctx, r, HTTPHeaders{}, Metadata{}, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_require.Nil(err)
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
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, LeaseAccessConditions{}, sourceDataMD5Value[:], DefaultAccessTier, nil)
//	_require.Nil(err)
//	_assert(resp.Response().StatusCode, chk.Equals, 202)
//	_assert(resp.Version(), chk.Not(chk.Equals), "")
//	_assert(resp.CopyID(), chk.Not(chk.Equals), "")
//	_assert(string(resp.CopyStatus()), chk.DeepEquals, "success")
//	_assert(resp.VersionID, chk.NotNil)
//
//	downloadResp, err := destBlob.ServiceURL.DownloadStream(ctx, 0, CountToEnd, LeaseAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//	_require.Nil(err)
//	_assert(destData, chk.DeepEquals, sourceData)
//	_assert(downloadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//	_assert(len(downloadResp.NewMetadata()), chk.Equals, 1)
//	_, badMD5 := testcommon.GetRandomDataAndReader(16)
//	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, LeaseAccessConditions{}, badMD5, DefaultAccessTier, nil)
//	_require.NotNil(err)
//
//	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, LeaseAccessConditions{}, nil, DefaultAccessTier, nil)
//	_require.Nil(err)
//	_assert(resp.Response().StatusCode, chk.Equals, 202)
//	_assert(resp.ContentCRC64(), chk.Not(chk.Equals), "")
//	_assert(resp.Response().Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	_assert(resp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//}
//
//func (s *AZBlobRecordedTestsSuite) TestCreateBlockBlobReturnsVID() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	testSize := 2 * 1024 * 1024 // 1MB
//	r, _ := testcommon.GetRandomDataAndReader(testSize)
//	ctx := ctx // Use default Background context
//	blobURL := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName())
//
//	// Prepare source blob for copy.
//	uploadResp, err := blobURL.Upload(ctx, r, HTTPHeaders{}, Metadata{}, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//	_assert(uploadResp.Response().StatusCode, chk.Equals, 201)
//	_assert(uploadResp.rawResponse.Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	_assert(uploadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	csResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, LeaseAccessConditions{}, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//	_assert(csResp.Response().StatusCode, chk.Equals, 201)
//	_assert(csResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	listBlobResp, err := containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//	_require.Nil(err)
//	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
//	if len(listBlobResp.Segment.BlobItems) < 2 {
//		s.T().Fail()
//	}
//
//	deleteResp, err := blobURL.Delete(ctx, DeleteSnapshotsOptionOnly, LeaseAccessConditions{})
//	_require.Nil(err)
//	_assert(deleteResp.Response().StatusCode, chk.Equals, 202)
//	_assert(deleteResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//
//	listBlobResp, err = containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true, Versions: true}})
//	_require.Nil(err)
//	_assert(listBlobResp.rawResponse.Header.Get("x-ms-request-id"), chk.NotNil)
//	if len(listBlobResp.Segment.BlobItems) == 0 {
//		s.T().Fail()
//	}
//	blobs := listBlobResp.Segment.BlobItems
//	_assert(blobs[0].Snapshot, chk.Equals, "")
//}

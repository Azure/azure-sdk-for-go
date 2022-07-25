//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strconv"
	"strings"
)

//
//import (
//	"bytes"
//	"context"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
//	"github.com/stretchr/testify/require"
//	"io/ioutil"
//	"strconv"
//	"strings"
//)

func (s *azblobTestSuite) TestBlockBlobGetPropertiesUsingVID() {
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
	bbClient := createNewBlockBlob(_require, generateBlobName(testName), containerClient)

	blobProp, _ := bbClient.GetProperties(ctx, nil)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		Metadata: basicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	uploadResp, err := bbClient.Upload(ctx, getReaderToGeneratedBytes(1024), &uploadBlockBlobOptions)
	_require.Nil(err)
	_require.NotNil(uploadResp.VersionID)
	blobProp, _ = bbClient.GetProperties(ctx, nil)
	_require.EqualValues(uploadResp.VersionID, blobProp.VersionID)
	_require.EqualValues(uploadResp.LastModified, blobProp.LastModified)
	_require.Equal(*uploadResp.ETag, *blobProp.ETag)
	_require.Equal(*blobProp.IsCurrentVersion, true)
}

func (s *azblobTestSuite) TestAppendBlobGetPropertiesUsingVID() {
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
	abClient := createNewAppendBlob(_require, generateBlobName(testName), containerClient)

	blobProp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: basicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	createResp, err := abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)
	_require.NotNil(createResp.VersionID)
	blobProp, _ = abClient.GetProperties(ctx, nil)
	_require.EqualValues(createResp.VersionID, blobProp.VersionID)
	_require.EqualValues(createResp.LastModified, blobProp.LastModified)
	_require.Equal(*createResp.ETag, *blobProp.ETag)
	_require.Equal(*blobProp.IsCurrentVersion, true)
}

//nolint
//func (s *azblobUnrecordedTestSuite) TestSetBlobMetadataReturnsVID() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	bbName := generateName(testName)
//	bbClient := createNewBlockBlob(_require, bbName, containerClient)
//
//	metadata := map[string]string{"test_key_1": "test_value_1", "test_key_2": "2019"}
//	resp, err := bbClient.SetMetadata(ctx, metadata, nil)
//	_require.Nil(err)
//	_require.NotNil(resp.VersionID)
//
//	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
//		Include: []container.ListBlobsIncludeItem{ListBlobsIncludeItemMetadata},
//	})
//
//	if !pager.NextPage(ctx) {
//		_require.Nil(pager.Err()) // check for an error first
//		s.T().Fail()             // no page was gotten
//	}
//
//	pageResp := pager.PageResponse()
//
//	_require.NotNil(pageResp.EnumerationResults.Segment.BlobItems)
//	blobList := pageResp.EnumerationResults.Segment.BlobItems
//	_require.Len(blobList, 1)
//	blobResp1 := blobList[0]
//	_require.Equal(*blobResp1.Name, bbName)
//	_require.NotNil(blobResp1.Metadata.AdditionalProperties)
//	_require.Len(blobResp1.Metadata.AdditionalProperties, 2)
//	// _assert(*blobResp1.Metadata, chk.DeepEquals, metadata)
//
//}

func (s *azblobTestSuite) TestCreateAndDownloadBlobSpecialCharactersWithVID() {
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
	data := []rune("-._/()$=',~0123456789")
	for i := 0; i < len(data); i++ {
		blobName := "abc" + string(data[i])
		blobClient := containerClient.NewBlockBlobClient(blobName)
		resp, err := blobClient.Upload(ctx, NopCloser(strings.NewReader(string(data[i]))), nil)
		_require.Nil(err)
		_require.NotNil(resp.VersionID)

		blobClientWithVersionID, err := blobClient.WithVersionID(*resp.VersionID)
		_require.Nil(err)
		dResp, err := blobClientWithVersionID.Download(ctx, nil)
		_require.Nil(err)
		d1, err := ioutil.ReadAll(dResp.BodyReader(nil))
		_require.Nil(err)
		_require.NotEqual(*dResp.Version, "")
		_require.EqualValues(string(d1), string(data[i]))
		_require.NotNil(dResp.VersionID)
		_require.Equal(*dResp.VersionID, *resp.VersionID)
	}
}

//func (s *azblobTestSuite) TestDeleteSpecificBlobVersion() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//	blobURL := getBlockBlobClient(generateBlobName(testName), containerClient)
//
//	uploadResp, err := blobURL.Upload(ctx, NopCloser(bytes.NewReader([]byte("data"))), &blockblob.UploadOptions{
//		Metadata: basicMetadata,
//	})
//	_require.Nil(err)
//	_require.NotNil(uploadResp.VersionID)
//	versionID1 := uploadResp.VersionID
//
//	uploadResp, err = blobURL.Upload(ctx, NopCloser(bytes.NewReader([]byte("updated_data"))),, &blockblob.UploadOptions{
//		Metadata: basicMetadata,
//	})
//	_require.Nil(err)
//	_require.NotNil(uploadResp.VersionID)
//
//	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
//		Include: &[]container.ListBlobsIncludeItem{container.ListBlobsIncludeItemVersions},
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
//	_require.Nil(listPager.Err())
//	_require.Len(count, 2)
//
//	// Deleting previous version snapshot.
//	deleteResp, err := blobURL.WithVersionID(versionID1).Delete(ctx, DeleteSnapshotsOptionNone, LeaseAccessConditions{})
//	_require.Nil(err)
//	_assert(deleteResp.StatusCode(), chk.Equals, 202)
//
//	listBlobsResp, err = containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
//	_require.Nil(err)
//	_assert(listBlobsResp.Segment.BlobItems, chk.NotNil)
//	if len(listBlobsResp.Segment.BlobItems) != 1 {
//		s.T().Fail()
//	}
//}
//
//func (s *azblobTestSuite) TestDeleteSpecificBlobVersionWithBlobSAS() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//	blobURL, blobName := getBlockBlobClient(c, containerClient)
//
//	resp, err := blobURL.Upload(ctx, NopCloser(bytes.NewReader([]byte("data"))), HTTPHeaders{}, basicMetadata, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//	versionId := resp.VersionID
//	_assert(versionId, chk.NotNil)
//
//	resp, err = blobURL.Upload(ctx, NopCloser(bytes.NewReader([]byte("updated_data"))),, HTTPHeaders{}, basicMetadata, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_require.Nil(err)
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
//	deleteResp, err := sbURL.Delete(ctx, DeleteSnapshotsOptionNone, LeaseAccessConditions{})
//	_assert(deleteResp, chk.IsNil)
//
//	listBlobResp, err := containerClient.NewListBlobsFlatPager(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Versions: true}})
//	_require.Nil(err)
//	for _, blob := range listBlobResp.Segment.BlobItems {
//		_assert(blob.VersionID, chk.Not(chk.Equals), versionId)
//	}
//}
func (s *azblobTestSuite) TestDeleteSpecificBlobVersion() {
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
	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		uploadResp, err := bbClient.Upload(ctx, NopCloser(bytes.NewReader([]byte("data"+strconv.Itoa(i)))), &blockblob.UploadOptions{
			Metadata: basicMetadata,
		})
		_require.Nil(err)
		_require.NotNil(uploadResp.VersionID)
		versions = append(versions, *uploadResp.VersionID)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemVersions},
	})

	found := make([]*container.BlobItem, 0)
	for listPager.More() {
		resp, err := listPager.NextPage(ctx)
		_require.Nil(err)
		if err != nil {
			break
		}
		found = append(found, resp.Segment.BlobItems...)
	}
	_require.Len(found, 5)

	// Deleting the 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.Nil(err)
		_, err = bbClientWithVersionID.Delete(ctx, nil)
		_require.Nil(err)
		//_require.Equal(deleteResp.RawResponse.StatusCode, 202)
	}

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemVersions},
	})

	found = make([]*container.BlobItem, 0)
	for listPager.More() {
		resp, err := listPager.NextPage(ctx)
		_require.Nil(err)
		if err != nil {
			break
		}
		found = append(found, resp.Segment.BlobItems...)
	}
	_require.Len(found, 2)

	for i := 3; i < 5; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.Nil(err)
		downloadResp, err := bbClientWithVersionID.Download(ctx, nil)
		_require.Nil(err)
		destData, err := ioutil.ReadAll(downloadResp.BodyReader(nil))
		_require.Nil(err)
		_require.EqualValues(destData, "data"+strconv.Itoa(i))
	}
}

//
//func (s *azblobTestSuite) TestCopyBlobFromURLWithSASReturnsVID() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	testSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(testSize)
//	sourceDataMD5Value := md5.Sum(sourceData)
//	ctx := ctx
//	srcBlob := container.NewBlockBlobClient(generateBlobName())
//	destBlob := container.NewBlockBlobClient(generateBlobName())
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
//	}.NewSASQueryParameters(credential)
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
//	downloadResp, err := destBlob.ServiceURL.Download(ctx, 0, CountToEnd, LeaseAccessConditions{}, false, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//	destData, err := ioutil.ReadAll(downloadresp.BodyReader(nil))
//	_require.Nil(err)
//	_assert(destData, chk.DeepEquals, sourceData)
//	_assert(downloadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//	_assert(len(downloadResp.NewMetadata()), chk.Equals, 1)
//	_, badMD5 := getRandomDataAndReader(16)
//	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, LeaseAccessConditions{}, badMD5, DefaultAccessTier, nil)
//	_require.NotNil(err)
//
//	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, LeaseAccessConditions{}, nil, DefaultAccessTier, nil)
//	_require.Nil(err)
//	_assert(resp.Response().StatusCode, chk.Equals, 202)
//	_assert(resp.XMsContentCRC64(), chk.Not(chk.Equals), "")
//	_assert(resp.Response().Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
//	_assert(resp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
//}
//
//func (s *azblobTestSuite) TestCreateBlockBlobReturnsVID() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//
//	_context := getTestContext(testName)
//	svcClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(_require, containerName, svcClient)
//	defer deleteContainer(_require, containerClient)
//
//	testSize := 2 * 1024 * 1024 // 1MB
//	r, _ := getRandomDataAndReader(testSize)
//	ctx := ctx // Use default Background context
//	blobURL := containerClient.NewBlockBlobClient(generateBlobName())
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

func (s *azblobTestSuite) TestPutBlockListReturnsVID() {
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

	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(ctx, base64BlockIDs[index], NopCloser(strings.NewReader(d)), nil)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.NotNil(resp.Version)
		_require.NotEqual(*resp.Version, "")
	}

	commitResp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, nil)
	_require.Nil(err)
	_require.NotNil(commitResp.VersionID)

	contentResp, err := bbClient.Download(ctx, nil)
	_require.Nil(err)
	contentData, err := ioutil.ReadAll(contentResp.BodyReader(nil))
	_require.Nil(err)
	_require.EqualValues(contentData, []uint8(strings.Join(data, "")))
}

//nolint
func (s *azblobUnrecordedTestSuite) TestCreateBlockBlobReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_require, containerName, svcClient)
	defer deleteContainer(_require, containerClient)

	testSize := 2 * 1024 * 1024 // 1MB
	r, _ := getRandomDataAndReader(testSize)
	ctx := ctx // Use default Background context
	bbClient := containerClient.NewBlockBlobClient(generateBlobName(testName))

	// Prepare source blob for copy.
	uploadResp, err := bbClient.Upload(ctx, NopCloser(r), nil)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.NotNil(uploadResp.VersionID)

	csResp, err := bbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)
	//_require.Equal(csResp.RawResponse.StatusCode, 201)
	_require.NotNil(csResp.VersionID)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemSnapshots},
	})

	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = bbClient.Delete(ctx, &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.Nil(err)
	//_require.Equal(deleteResp.RawResponse.StatusCode, 202)

	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemSnapshots, container.ListBlobsIncludeItemVersions},
	})

	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.NotEqual(len(found), 0)
}

func (s *azblobTestSuite) TestCreatePageBlobReturnsVID() {
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

	pbClob := createNewPageBlob(_require, generateBlobName(testName), containerClient)

	contentSize := 1 * 1024
	r, _ := generateData(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(offset),
		Count:  to.Ptr(count),
	}
	putResp, err := pbClob.UploadPages(ctx, r, &uploadPagesOptions)
	_require.Nil(err)
	//_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotEqual(*putResp.Version, "")

	gpResp, err := pbClob.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/require"
)

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestSetBlobTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"azure":    "bbClient",
		"bbClient": "sdk",
		"sdk":      "go",
	}

	contentSize := 4 * 1024 * 1024 // 4MB
	r, _ := testcommon.GenerateData(contentSize)

	_, err = bbClient.Upload(ctx, r, nil)
	_require.Nil(err)
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	_, err = bbClient.SetTags(ctx, blobTagsMap, nil)
	_require.Nil(err)
	// _require.Equal(blobSetTagsResponse.RawResponse.StatusCode, 204)

	blobGetTagsResponse, err := bbClient.GetTags(ctx, nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, 3)
	for _, blobTag := range blobTagsSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestSetBlobTagsWithVID() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"Go":         "CPlusPlus",
		"Python":     "CSharp",
		"Javascript": "Android",
	}

	blockBlobUploadResp, err := bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.Nil(err)
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	versionId1 := blockBlobUploadResp.VersionID

	blockBlobUploadResp, err = bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("updated_data"))), nil)
	_require.Nil(err)
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	versionId2 := blockBlobUploadResp.VersionID

	setTagsBlobOptions := blob.SetTagsOptions{
		VersionID: versionId1,
	}
	_, err = bbClient.SetTags(ctx, blobTagsMap, &setTagsBlobOptions)
	_require.Nil(err)
	// _require.Equal(blobSetTagsResponse.RawResponse.StatusCode, 204)

	getTagsBlobOptions1 := blob.GetTagsOptions{
		VersionID: versionId1,
	}
	blobGetTagsResponse, err := bbClient.GetTags(ctx, &getTagsBlobOptions1)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_require.NotNil(blobGetTagsResponse.BlobTagSet)
	_require.Len(blobGetTagsResponse.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResponse.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	getTagsBlobOptions2 := blob.GetTagsOptions{
		VersionID: versionId2,
	}
	blobGetTagsResponse, err = bbClient.GetTags(ctx, &getTagsBlobOptions2)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_require.Nil(blobGetTagsResponse.BlobTagSet)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadBlockBlobWithSpecialCharactersInTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	//_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"+-./:=_ ": "firsttag",
		"tag2":     "+-./:=_",
		"+-./:=_1": "+-./:=_",
	}

	uploadBlockBlobOptions := blockblob.UploadOptions{
		Metadata:    testcommon.BasicMetadata,
		HTTPHeaders: &testcommon.BasicHeaders,
		Tags:        blobTagsMap,
	}
	_, err = bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("data"))), &uploadBlockBlobOptions)
	_require.Nil(err)
	// TODO: Check for metadata and header
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	blobGetTagsResponse, err := bbClient.GetTags(ctx, nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_require.Len(blobGetTagsResponse.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResponse.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestStageBlockWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	//_context := getTestContext(testName)
	//ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = testcommon.BlockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(ctx, base64BlockIDs[index], streaming.NopCloser(strings.NewReader(d)), nil)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.NotEqual(*resp.Version, "")
	}

	blobTagsMap := map[string]string{
		"azure":    "bbClient",
		"bbClient": "sdk",
		"sdk":      "go",
	}

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		Tags: blobTagsMap,
	}
	commitResp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_require.Nil(err)
	_require.NotNil(commitResp.VersionID)
	versionId := commitResp.VersionID

	contentResp, err := bbClient.DownloadStream(ctx, nil)
	_require.Nil(err)
	contentData, err := io.ReadAll(contentResp.Body)
	_require.Nil(err)
	_require.EqualValues(contentData, []uint8(strings.Join(data, "")))

	getTagsBlobOptions := blob.GetTagsOptions{
		VersionID: versionId,
	}
	blobGetTagsResp, err := bbClient.GetTags(ctx, &getTagsBlobOptions)
	_require.Nil(err)
	_require.NotNil(blobGetTagsResp)
	_require.Len(blobGetTagsResp.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResp.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	blobGetTagsResp, err = bbClient.GetTags(ctx, nil)
	_require.Nil(err)
	_require.NotNil(blobGetTagsResp)
	_require.Len(blobGetTagsResp.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResp.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

//
////nolint
//func (s *AZBlobUnrecordedTestsSuite) TestStageBlockFromURLWithTags() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	if err != nil {
//		s.T().Fatal("Invalid credential")
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := testcommon.GenerateData(contentSize)
//	ctx := ctx // Use default Background context
//	srcBlob := containerClient.NewBlockBlobClient("sourceBlob")
//	destBlob := containerClient.NewBlockBlobClient("destBlob")
//
//	blobTagsMap := map[string]string{
//		"Go":         "CPlusPlus",
//		"Python":     "CSharp",
//		"Javascript": "Android",
//	}
//
//	uploadBlockBlobOptions := blockblob.UploadOptions{
//		Tags: blobTagsMap,
//	}
//	uploadSrcResp, err := srcBlob.Upload(ctx, r, &uploadBlockBlobOptions)
//	_require.Nil(err)
//	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//	uploadDate := uploadSrcResp.Date
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := azblob.ParseURL(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    uploadDate.UTC().Add(1 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fail()
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
//	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
//
//	offset1, count1 := int64(0), int64(contentSize/2)
//	options1 := BlockBlobStageBlockFromURLOptions{
//		Offset: &offset1,
//		Count:  &count1,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.Nil(err)
//	// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//	_require.NotEqual(*stageResp1.RequestID, "")
//	_require.NotEqual(*stageResp1.Version, "")
//	_require.NotNil(stageResp1.Date)
//	_require.Equal((*stageResp1.Date).IsZero(), false)
//
//	offset2, count2 := int64(contentSize/2), int64(CountToEnd)
//	options2 := BlockBlobStageBlockFromURLOptions{
//		Offset: &offset2,
//		Count:  &count2,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.Nil(err)
//	// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(*stageResp2.RequestID, "")
//	_require.NotEqual(*stageResp2.Version, "")
//	_require.NotNil(stageResp2.Date)
//	_require.Equal((*stageResp2.Date).IsZero(), false)
//
//	blockList, err := destBlob.GetBlockList(ctx, BlockListTypeAll, nil)
//	_require.Nil(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	commitBlockListOptions := blockblob.CommitBlockListOptions{
//		Tags: blobTagsMap,
//	}
//	_, err = destBlob.CommitBlockList(ctx, []string{blockID1, blockID2}, &commitBlockListOptions)
//	_require.Nil(err)
//	// _require.Equal(listResp.RawResponse.StatusCode,  201)
//	//versionId := listResp.VersionID()
//
//	blobGetTagsResp, err := destBlob.GetTags(ctx, nil)
//	_require.Nil(err)
//	_require.Len(blobGetTagsResp.BlobTagSet, 3)
//	for _, blobTag := range blobGetTagsResp.BlobTagSet {
//		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
//	}
//
//	downloadResp, err := destBlob.DownloadStream(ctx, nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(nil))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestCopyBlockBlobFromURLWithTags() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	if err != nil {
//		s.T().Fatal("Invalid credential")
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 1 * 1024 * 1024 // 1MB
//	r, sourceData := testcommon.GenerateData(contentSize)
//	sourceDataMD5Value := md5.Sum(sourceData)
//	srcBlob := containerClient.NewBlockBlobClient("srcBlob")
//	destBlob := containerClient.NewBlockBlobClient("destBlob")
//
//	blobTagsMap := map[string]string{
//		"Go":         "CPlusPlus",
//		"Python":     "CSharp",
//		"Javascript": "Android",
//	}
//
//	uploadBlockBlobOptions := blockblob.UploadOptions{
//		Tags: blobTagsMap,
//	}
//	_, err = srcBlob.Upload(ctx, r, &uploadBlockBlobOptions)
//	_require.Nil(err)
//	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
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
//	sourceContentMD5 := sourceDataMD5Value[:]
//	copyBlockBlobFromURLOptions1 := BlockBlobCopyFromURLOptions{
//		Tags:         map[string]string{"foo": "bar"},
//		SourceContentMD5: sourceContentMD5,
//	}
//	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 202)
//	_require.NotEqual(*resp.ETag, "")
//	_require.NotEqual(*resp.RequestID, "")
//	_require.NotEqual(*resp.Version, "")
//	_require.Equal((*resp.Date).IsZero(), false)
//	_require.NotEqual(*resp.CopyID, "")
//	_require.EqualValues(resp.ContentMD5, sourceDataMD5Value[:])
//	_require.EqualValues(*resp.CopyStatus, "success")
//
//	downloadResp, err := destBlob.DownloadStream(ctx, nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//	_require.Equal(*downloadResp.TagCount, int64(1))
//
//	_, badMD5 := getRandomDataAndReader(16)
//	copyBlockBlobFromURLOptions2 := BlockBlobCopyFromURLOptions{
//		Tags:         blobTagsMap,
//		SourceContentMD5: badMD5,
//	}
//	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
//	_require.NotNil(err)
//
//	copyBlockBlobFromURLOptions3 := BlockBlobCopyFromURLOptions{
//		Tags: blobTagsMap,
//	}
//	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions3)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 202)
//}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestGetPropertiesReturnsTagsCount() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		Tags:        testcommon.BasicBlobTagsMap,
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
	}
	_, err = bbClient.Upload(ctx, streaming.NopCloser(bytes.NewReader([]byte("data"))), &uploadBlockBlobOptions)
	_require.Nil(err)

	getPropertiesResponse, err := bbClient.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*getPropertiesResponse.TagCount, int64(3))

	downloadResp, err := bbClient.DownloadStream(ctx, nil)
	_require.Nil(err)
	_require.NotNil(downloadResp)
	_require.Equal(*downloadResp.TagCount, int64(3))
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestSetBlobTagForSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	// _context := getTestContext(testName)
	// ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"Microsoft Azure": "Azure Storage",
		"Storage+SDK":     "SDK/GO",
		"GO ":             ".Net",
	}
	_, err = bbClient.SetTags(ctx, blobTagsMap, nil)
	_require.Nil(err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_require.Nil(err)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_require.Nil(err)
	_require.Equal(*resp2.TagCount, int64(3))
}

// TODO: Once new pacer is done.
// nolint
func (s *AZBlobUnrecordedTestsSuite) TestListBlobReturnsTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	blobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobTagsMap := map[string]string{
		"+-./:=_ ": "firsttag",
		"tag2":     "+-./:=_",
		"+-./:=_1": "+-./:=_",
	}

	_, err = blobClient.SetTags(ctx, blobTagsMap, nil)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode,204)

	include := []container.ListBlobsIncludeItem{container.ListBlobsIncludeItemTags}
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: include,
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

	_require.Equal(*(found[0].Name), blobName)
	_require.Len(found[0].BlobTags.BlobTagSet, 3)
	for _, blobTag := range found[0].BlobTags.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

//
////nolint
//func (s *AZBlobUnrecordedTestsSuite) TestFindBlobsByTags() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerClient1 := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName) + "1", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient1)
//
//	containerClient2 := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName) + "2", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient2)
//
//	containerClient3 := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName) + "3", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient3)
//
//	blobTagsMap1 := map[string]string{
//		"tag2": "tagsecond",
//		"tag3": "tagthird",
//	}
//	blobTagsMap2 := map[string]string{
//		"tag1": "firsttag",
//		"tag2": "secondtag",
//		"tag3": "thirdtag",
//	}
//
//	blobURL11 := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName) + "11", containerClient1)
//	_, err = blobURL11.Upload(ctx, bytes.NewReader([]byte("random data")), &blockblob.UploadOptions{
//		Metadata: testcommon.BasicMetadata,
//		Tags: blobTagsMap1,
//	})
//	_require.Nil(err)
//
//	blobURL12 := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName) + "12", containerClient1)
//	_, err = blobURL12.Upload(ctx, bytes.NewReader([]byte("another random data")), &blockblob.UploadOptions{
//		Metadata: testcommon.BasicMetadata,
//		Tags: blobTagsMap2,
//	})
//	_require.Nil(err)
//
//	blobURL21 := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName) + "21", containerClient2)
//	_, err = blobURL21.Upload(ctx, bytes.NewReader([]byte("random data")), HTTPHeaders{}, testcommon.BasicMetadata, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//
//	blobURL22 := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName) + "22", containerClient2)
//	_, err = blobURL22.Upload(ctx, bytes.NewReader([]byte("another random data")), HTTPHeaders{}, testcommon.BasicMetadata, LeaseAccessConditions{}, DefaultAccessTier, blobTagsMap2, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//
//	blobURL31 := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName) + "31", containerClient3)
//	_, err = blobURL31.Upload(ctx, bytes.NewReader([]byte("random data")), HTTPHeaders{}, testcommon.BasicMetadata, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_require.Nil(err)
//
//	where := "\"tag4\"='fourthtag'"
//	lResp, err := svcClient.FindBlobByTags(ctx, nil, nil, &where, Marker{}, nil)
//	_require.Nil(err)
//	_assert(lResp.Blobs, chk.HasLen, 0)
//
//	//where = "\"tag1\"='firsttag'AND\"tag2\"='secondtag'AND\"@container\"='"+ containerName1 + "'"
//	//TODO: Figure out how to do a composite query based on container.
//	where = "\"tag1\"='firsttag'AND\"tag2\"='secondtag'"
//
//	lResp, err = svcClient.FindBlobsByTags(ctx, nil, nil, &where, Marker{}, nil)
//	_require.Nil(err)
//
//	for _, blob := range lResp.Blobs {
//		_assert(blob.TagValue, chk.Equals, "firsttag")
//	}
//}
//
//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestFilterBlobsUsingAccountSAS() {
//	accountName, accountKey := accountInfo()
//	credential, err := NewSharedKeyCredential(accountName, accountKey)
//	if err != nil {
//		s.T().Fail()
//	}
//
//	sasQueryParams, err := AccountSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
//		Permissions:   AccountSASPermissions{Read: true, List: true, Write: true, DeletePreviousVersion: true, Tag: true, FilterByTags: true, Create: true}.String(),
//		Services:      AccountSASServices{Blob: true}.String(),
//		ResourceTypes: AccountSASResourceTypes{Service: true, Container: true, Object: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	qp := sasQueryParams.Encode()
//	urlToSendToSomeone := fmt.Sprintf("https://%s.blob.core.windows.net?%s", accountName, qp)
//	u, _ := url.Parse(urlToSendToSomeone)
//	serviceURL := NewServiceURL(*u, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//
//	containerName := testcommon.GenerateContainerName()
//	containerClient := serviceURL.NewcontainerClient(containerName)
//	_, err = containerClient.Create(ctx, Metadata{}, PublicAccessNone)
//	defer containerClient.Delete(ctx, LeaseAccessConditions{})
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	blobClient := containerClient.NewBlockBlobURL("temp")
//	_, err = blobClient.Upload(ctx, bytes.NewReader([]byte("random data")), HTTPHeaders{}, testcommon.BasicMetadata, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	if err != nil {
//		s.T().Fail()
//	}
//
//	blobTagsMap := BlobTags{"tag1": "firsttag", "tag2": "secondtag", "tag3": "thirdtag"}
//	setBlobTagsResp, err := blobClient.SetTags(ctx, nil, nil, nil, nil, nil, nil, blobTagsMap)
//	_require.Nil(err)
//	_assert(setBlobTagsResp.StatusCode(), chk.Equals, 204)
//
//	blobGetTagsResp, err := blobClient.GetTags(ctx, nil, nil, nil, nil, nil)
//	_require.Nil(err)
//	_assert(blobGetTagsResp.StatusCode(), chk.Equals, 200)
//	_assert(blobGetTagsResp.BlobTagSet, chk.HasLen, 3)
//	for _, blobTag := range blobGetTagsResp.BlobTagSet {
//		_assert(blobTagsMap[blobTag.Key], chk.Equals, blobTag.Value)
//	}
//
//	time.Sleep(30 * time.Second)
//	where := "\"tag1\"='firsttag'AND\"tag2\"='secondtag'AND@container='" + containerName + "'"
//	_, err = serviceURL.FindBlobsByTags(ctx, nil, nil, &where, Marker{}, nil)
//	_require.Nil(err)
//}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"time"
)

func (s *azblobTestSuite) TestSetBlobTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"azure":    "bbClient",
		"bbClient": "sdk",
		"sdk":      "go",
	}

	contentSize := 4 * 1024 * 1024 // 4MB
	r, _ := generateData(contentSize)

	blockBlobUploadResp, err := bbClient.Upload(ctx, r, nil)
	_assert.Nil(err)
	_assert.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	setTagsBlobOptions := SetTagsBlobOptions{
		BlobTagsMap: &blobTagsMap,
	}
	blobSetTagsResponse, err := bbClient.SetTags(ctx, &setTagsBlobOptions)
	_assert.Nil(err)
	_assert.Equal(blobSetTagsResponse.RawResponse.StatusCode, 204)

	blobGetTagsResponse, err := bbClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.Tags.BlobTagSet
	_assert.NotNil(blobTagsSet)
	_assert.Len(*blobTagsSet, 3)
	for _, blobTag := range *blobTagsSet {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *azblobTestSuite) TestSetBlobTagsWithVID() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"Go":         "CPlusPlus",
		"Python":     "CSharp",
		"Javascript": "Android",
	}

	blockBlobUploadResp, err := bbClient.Upload(ctx, bytes.NewReader([]byte("data")), nil)
	_assert.Nil(err)
	_assert.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	versionId1 := blockBlobUploadResp.VersionID

	blockBlobUploadResp, err = bbClient.Upload(ctx, bytes.NewReader([]byte("updated_data")), nil)
	_assert.Nil(err)
	_assert.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	versionId2 := blockBlobUploadResp.VersionID

	setTagsBlobOptions := SetTagsBlobOptions{
		BlobTagsMap: &blobTagsMap,
		VersionID:   versionId1,
	}
	blobSetTagsResponse, err := bbClient.SetTags(ctx, &setTagsBlobOptions)
	_assert.Nil(err)
	_assert.Equal(blobSetTagsResponse.RawResponse.StatusCode, 204)

	getTagsBlobOptions1 := GetTagsBlobOptions{
		VersionID: versionId1,
	}
	blobGetTagsResponse, err := bbClient.GetTags(ctx, &getTagsBlobOptions1)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_assert.NotNil(blobGetTagsResponse.Tags.BlobTagSet)
	_assert.Len(*blobGetTagsResponse.Tags.BlobTagSet, 3)
	for _, blobTag := range *blobGetTagsResponse.Tags.BlobTagSet {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	getTagsBlobOptions2 := GetTagsBlobOptions{
		VersionID: versionId2,
	}
	blobGetTagsResponse, err = bbClient.GetTags(ctx, &getTagsBlobOptions2)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_assert.Nil(blobGetTagsResponse.Tags.BlobTagSet)
}

func (s *azblobTestSuite) TestUploadBlockBlobWithSpecialCharactersInTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"+-./:=_ ": "firsttag",
		"tag2":     "+-./:=_",
		"+-./:=_1": "+-./:=_",
	}

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata:        &basicMetadata,
		BlobHTTPHeaders: &basicHeaders,
		BlobTagsMap:     &blobTagsMap,
	}
	blockBlobUploadResp, err := bbClient.Upload(ctx, bytes.NewReader([]byte("data")), &uploadBlockBlobOptions)
	_assert.Nil(err)
	// TODO: Check for metadata and header
	_assert.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	blobGetTagsResponse, err := bbClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_assert.Len(*blobGetTagsResponse.Tags.BlobTagSet, 3)
	for _, blobTag := range *blobGetTagsResponse.Tags.BlobTagSet {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *azblobTestSuite) TestStageBlockWithTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(ctx, base64BlockIDs[index], strings.NewReader(d), nil)
		_assert.Nil(err)
		_assert.Equal(resp.RawResponse.StatusCode, 201)
		_assert.NotEqual(*resp.Version, "")
	}

	blobTagsMap := map[string]string{
		"azure":    "bbClient",
		"bbClient": "sdk",
		"sdk":      "go",
	}

	commitBlockListOptions := CommitBlockListOptions{
		BlobTagsMap: &blobTagsMap,
	}
	commitResp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_assert.Nil(err)
	_assert.NotNil(commitResp.VersionID)
	versionId := commitResp.VersionID

	contentResp, err := bbClient.Download(ctx, nil)
	_assert.Nil(err)
	contentData, err := ioutil.ReadAll(contentResp.Body(RetryReaderOptions{}))
	_assert.EqualValues(contentData, []uint8(strings.Join(data, "")))

	getTagsBlobOptions := GetTagsBlobOptions{
		VersionID: versionId,
	}
	blobGetTagsResp, err := bbClient.GetTags(ctx, &getTagsBlobOptions)
	_assert.Nil(err)
	_assert.NotNil(blobGetTagsResp)
	_assert.Len(*blobGetTagsResp.Tags.BlobTagSet, 3)
	for _, blobTag := range *blobGetTagsResp.Tags.BlobTagSet {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	blobGetTagsResp, err = bbClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(blobGetTagsResp)
	_assert.Len(*blobGetTagsResp.Tags.BlobTagSet, 3)
	for _, blobTag := range *blobGetTagsResp.Tags.BlobTagSet {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *azblobUnrecordedTestSuite) TestStageBlockFromURLWithTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	serviceClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	credential, err := getGenericCredential(nil, testAccountDefault)
	if err != nil {
		s.T().Fatal("Invalid credential")
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := generateData(contentSize)
	ctx := ctx // Use default Background context
	srcBlob := containerClient.NewBlockBlobClient("sourceBlob")
	destBlob := containerClient.NewBlockBlobClient("destBlob")

	blobTagsMap := map[string]string{
		"Go":         "CPlusPlus",
		"Python":     "CSharp",
		"Javascript": "Android",
	}

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobTagsMap: &blobTagsMap,
	}
	uploadSrcResp, err := srcBlob.Upload(ctx, r, &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
	uploadDate := uploadSrcResp.Date

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    uploadDate.UTC().Add(1 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fail()
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))

	offset1, count1 := int64(0), int64(contentSize/2)
	options1 := StageBlockFromURLOptions{
		Offset: &offset1,
		Count:  &count1,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	_assert.Nil(err)
	_assert.Equal(stageResp1.RawResponse.StatusCode, 201)
	_assert.NotEqual(*stageResp1.RequestID, "")
	_assert.NotEqual(*stageResp1.Version, "")
	_assert.NotNil(stageResp1.Date)
	_assert.Equal((*stageResp1.Date).IsZero(), false)

	offset2, count2 := int64(contentSize/2), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset: &offset2,
		Count:  &count2,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	_assert.Nil(err)
	_assert.Equal(stageResp2.RawResponse.StatusCode, 201)
	_assert.NotEqual(*stageResp2.RequestID, "")
	_assert.NotEqual(*stageResp2.Version, "")
	_assert.NotNil(stageResp2.Date)
	_assert.Equal((*stageResp2.Date).IsZero(), false)

	blockList, err := destBlob.GetBlockList(ctx, BlockListAll, nil)
	_assert.Nil(err)
	_assert.Equal(blockList.RawResponse.StatusCode, 200)
	_assert.Nil(blockList.BlockList.CommittedBlocks)
	_assert.Len(*blockList.BlockList.UncommittedBlocks, 2)

	commitBlockListOptions := CommitBlockListOptions{
		BlobTagsMap: &blobTagsMap,
	}
	listResp, err := destBlob.CommitBlockList(ctx, []string{blockID1, blockID2}, &commitBlockListOptions)
	_assert.Nil(err)
	_assert.Equal(listResp.RawResponse.StatusCode, 201)
	//versionId := listResp.VersionID()

	blobGetTagsResp, err := destBlob.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Len(*blobGetTagsResp.Tags.BlobTagSet, 3)
	for _, blobTag := range *blobGetTagsResp.Tags.BlobTagSet {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	downloadResp, err := destBlob.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, sourceData)
}

func (s *azblobUnrecordedTestSuite) TestCopyBlockBlobFromURLWithTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	serviceClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	credential, err := getGenericCredential(nil, testAccountDefault)
	if err != nil {
		s.T().Fatal("Invalid credential")
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 1 * 1024 * 1024 // 1MB
	r, sourceData := generateData(contentSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	srcBlob := containerClient.NewBlockBlobClient("srcBlob")
	destBlob := containerClient.NewBlockBlobClient("destBlob")

	blobTagsMap := map[string]string{
		"Go":         "CPlusPlus",
		"Python":     "CSharp",
		"Javascript": "Android",
	}

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobTagsMap: &blobTagsMap,
	}
	uploadSrcResp, err := srcBlob.Upload(ctx, r, &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.T().Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	sourceContentMD5 := sourceDataMD5Value[:]
	copyBlockBlobFromURLOptions1 := CopyBlockBlobFromURLOptions{
		BlobTagsMap:      &map[string]string{"foo": "bar"},
		SourceContentMD5: &sourceContentMD5,
	}
	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
	_assert.NotEqual(*resp.ETag, "")
	_assert.NotEqual(*resp.RequestID, "")
	_assert.NotEqual(*resp.Version, "")
	_assert.Equal((*resp.Date).IsZero(), false)
	_assert.NotEqual(*resp.CopyID, "")
	_assert.EqualValues(*resp.ContentMD5, sourceDataMD5Value[:])
	_assert.EqualValues(*resp.CopyStatus, "success")

	downloadResp, err := destBlob.Download(ctx, nil)
	_assert.Nil(err)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	_assert.Nil(err)
	_assert.EqualValues(destData, sourceData)
	_assert.Equal(*downloadResp.TagCount, int64(1))

	_, badMD5 := getRandomDataAndReader(16)
	copyBlockBlobFromURLOptions2 := CopyBlockBlobFromURLOptions{
		BlobTagsMap:      &blobTagsMap,
		SourceContentMD5: &badMD5,
	}
	_, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
	_assert.NotNil(err)

	copyBlockBlobFromURLOptions3 := CopyBlockBlobFromURLOptions{
		BlobTagsMap: &blobTagsMap,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions3)
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 202)
}

func (s *azblobTestSuite) TestGetPropertiesReturnsTagsCount() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobTagsMap:     &basicBlobTagsMap,
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
	}
	blockBlobUploadResp, err := bbClient.Upload(ctx, bytes.NewReader([]byte("data")), &uploadBlockBlobOptions)
	_assert.Nil(err)
	_assert.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	getPropertiesResponse, err := bbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*getPropertiesResponse.TagCount, int64(3))

	downloadResp, err := bbClient.Download(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(downloadResp)
	_assert.Equal(*downloadResp.TagCount, int64(3))
}

func (s *azblobTestSuite) TestSetBlobTagForSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	bbClient := createNewBlockBlob(_assert, generateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"Microsoft Azure": "Azure Storage",
		"Storage+SDK":     "SDK/GO",
		"GO ":             ".Net",
	}
	setTagsBlobOptions := SetTagsBlobOptions{
		BlobTagsMap: &blobTagsMap,
	}
	_, err = bbClient.SetTags(ctx, &setTagsBlobOptions)
	_assert.Nil(err)

	resp, err := bbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp2.TagCount, int64(3))
}

// TODO: Once new pacer is done.
func (s *azblobTestSuite) TestListBlobReturnsTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	blobName := generateBlobName(testName)
	blobClient := createNewBlockBlob(_assert, blobName, containerClient)
	blobTagsMap := map[string]string{
		"+-./:=_ ": "firsttag",
		"tag2":     "+-./:=_",
		"+-./:=_1": "+-./:=_",
	}

	resp, err := blobClient.SetTags(ctx, &SetTagsBlobOptions{
		BlobTagsMap: &blobTagsMap,
	})
	_assert.Nil(err)
	_assert.Equal(resp.RawResponse.StatusCode, 204)

	include := []ListBlobsIncludeItem{ListBlobsIncludeItemTags}
	pager := containerClient.ListBlobsFlatSegment(&ContainerListBlobFlatSegmentOptions{
		Include: &include,
	})

	found := make([]*BlobItemInternal, 0)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range *resp.EnumerationResults.Segment.BlobItems {
			found = append(found, blob)
		}
	}
	_assert.Nil(pager.Err())

	_assert.Equal(*(found[0].Name), blobName)
	_assert.Len(*(found[0].BlobTags.BlobTagSet), 3)
	for _, blobTag := range *(found[0].BlobTags.BlobTagSet) {
		_assert.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

//func (s *azblobTestSuite) TestFindBlobsByTags() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	_context := getTestContext(testName)
//	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
//	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerClient1 := createNewContainer(_assert, generateContainerName(testName) + "1", serviceClient)
//	defer deleteContainer(_assert, containerClient1)
//
//	containerClient2 := createNewContainer(_assert, generateContainerName(testName) + "2", serviceClient)
//	defer deleteContainer(_assert, containerClient2)
//
//	containerClient3 := createNewContainer(_assert, generateContainerName(testName) + "3", serviceClient)
//	defer deleteContainer(_assert, containerClient3)
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
//	blobURL11 := getBlockBlobClient(generateBlobName(testName) + "11", containerClient1)
//	_, err = blobURL11.Upload(ctx, bytes.NewReader([]byte("random data")), &UploadBlockBlobOptions{
//		Metadata: &basicMetadata,
//		BlobTagsMap: &blobTagsMap1,
//	})
//	_assert.Nil(err)
//
//	blobURL12 := getBlockBlobClient(generateBlobName(testName) + "12", containerClient1)
//	_, err = blobURL12.Upload(ctx, bytes.NewReader([]byte("another random data")), &UploadBlockBlobOptions{
//		Metadata: &basicMetadata,
//		BlobTagsMap: &blobTagsMap2,
//	})
//	_assert.Nil(err)
//
//	blobURL21 := getBlockBlobClient(generateBlobName(testName) + "21", containerClient2)
//	_, err = blobURL21.Upload(ctx, bytes.NewReader([]byte("random data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//
//	blobURL22 := getBlockBlobClient(generateBlobName(testName) + "22", containerClient2)
//	_, err = blobURL22.Upload(ctx, bytes.NewReader([]byte("another random data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, blobTagsMap2, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//
//	blobURL31 := getBlockBlobClient(generateBlobName(testName) + "31", containerClient3)
//	_, err = blobURL31.Upload(ctx, bytes.NewReader([]byte("random data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	_assert.Nil(err)
//
//	where := "\"tag4\"='fourthtag'"
//	lResp, err := serviceClient.FindBlobByTags(ctx, nil, nil, &where, Marker{}, nil)
//	_assert.Nil(err)
//	_assert(lResp.Blobs, chk.HasLen, 0)
//
//	//where = "\"tag1\"='firsttag'AND\"tag2\"='secondtag'AND\"@container\"='"+ containerName1 + "'"
//	//TODO: Figure out how to do a composite query based on container.
//	where = "\"tag1\"='firsttag'AND\"tag2\"='secondtag'"
//
//	lResp, err = serviceClient.FindBlobsByTags(ctx, nil, nil, &where, Marker{}, nil)
//	_assert.Nil(err)
//
//	for _, blob := range lResp.Blobs {
//		_assert(blob.TagValue, chk.Equals, "firsttag")
//	}
//}

//func (s *azblobTestSuite) TestFilterBlobsUsingAccountSAS() {
//	accountName, accountKey := accountInfo()
//	credential, err := NewSharedKeyCredential(accountName, accountKey)
//	if err != nil {
//		c.Fail()
//	}
//
//	sasQueryParams, err := AccountSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
//		Permissions:   AccountSASPermissions{Read: true, List: true, Write: true, DeletePreviousVersion: true, Tag: true, FilterByTags: true, Create: true}.String(),
//		Services:      AccountSASServices{Blob: true}.String(),
//		ResourceTypes: AccountSASResourceTypes{Service: true, Container: true, Object: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	qp := sasQueryParams.Encode()
//	urlToSendToSomeone := fmt.Sprintf("https://%s.blob.core.windows.net?%s", accountName, qp)
//	u, _ := url.Parse(urlToSendToSomeone)
//	serviceURL := NewServiceURL(*u, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
//
//	containerName := generateContainerName()
//	containerClient := serviceURL.NewcontainerClient(containerName)
//	_, err = containerClient.Create(ctx, Metadata{}, PublicAccessNone)
//	defer containerClient.Delete(ctx, ContainerAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	blobClient := containerClient.NewBlockBlobURL("temp")
//	_, err = blobClient.Upload(ctx, bytes.NewReader([]byte("random data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	if err != nil {
//		c.Fail()
//	}
//
//	blobTagsMap := BlobTagsMap{"tag1": "firsttag", "tag2": "secondtag", "tag3": "thirdtag"}
//	setBlobTagsResp, err := blobClient.SetTags(ctx, nil, nil, nil, nil, nil, nil, blobTagsMap)
//	_assert.Nil(err)
//	_assert(setBlobTagsResp.StatusCode(), chk.Equals, 204)
//
//	blobGetTagsResp, err := blobClient.GetTags(ctx, nil, nil, nil, nil, nil)
//	_assert.Nil(err)
//	_assert(blobGetTagsResp.StatusCode(), chk.Equals, 200)
//	_assert(blobGetTagsResp.BlobTagSet, chk.HasLen, 3)
//	for _, blobTag := range blobGetTagsResp.BlobTagSet {
//		_assert(blobTagsMap[blobTag.Key], chk.Equals, blobTag.Value)
//	}
//
//	time.Sleep(30 * time.Second)
//	where := "\"tag1\"='firsttag'AND\"tag2\"='secondtag'AND@container='" + containerName + "'"
//	_, err = serviceURL.FindBlobsByTags(ctx, nil, nil, &where, Marker{}, nil)
//	_assert.Nil(err)
//}

func (s *azblobTestSuite) TestCreatePageBlobWithTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	pbClient := createNewPageBlob(_assert, generateBlobName(testName), containerClient)

	contentSize := 1 * 1024
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	putResp, err := pbClient.UploadPages(ctx, getReaderToGeneratedBytes(1024), &uploadPagesOptions)
	_assert.Nil(err)
	_assert.Equal(putResp.RawResponse.StatusCode, 201)
	_assert.Equal(putResp.LastModified.IsZero(), false)
	_assert.NotEqual(putResp.ETag, ETagNone)
	_assert.NotEqual(putResp.Version, "")

	setTagsBlobOptions := SetTagsBlobOptions{
		BlobTagsMap: &basicBlobTagsMap,
	}
	setTagResp, err := pbClient.SetTags(ctx, &setTagsBlobOptions)
	_assert.Nil(err)
	_assert.Equal(setTagResp.RawResponse.StatusCode, 204)

	gpResp, err := pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(gpResp)
	_assert.Equal(*gpResp.TagCount, int64(len(basicBlobTagsMap)))

	blobGetTagsResponse, err := pbClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.Tags.BlobTagSet
	_assert.NotNil(blobTagsSet)
	_assert.Len(*blobTagsSet, len(basicBlobTagsMap))
	for _, blobTag := range *blobTagsSet {
		_assert.Equal(basicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	modifiedBlobTags := map[string]string{
		"a0z1u2r3e4": "b0l1o2b3",
		"b0l1o2b3":   "s0d1k2",
	}

	setTagsBlobOptions2 := SetTagsBlobOptions{
		BlobTagsMap: &modifiedBlobTags,
	}
	setTagResp, err = pbClient.SetTags(ctx, &setTagsBlobOptions2)
	_assert.Nil(err)
	_assert.Equal(setTagResp.RawResponse.StatusCode, 204)

	gpResp, err = pbClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.NotNil(gpResp)
	_assert.Equal(*gpResp.TagCount, int64(len(modifiedBlobTags)))

	blobGetTagsResponse, err = pbClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet = blobGetTagsResponse.Tags.BlobTagSet
	_assert.NotNil(blobTagsSet)
	_assert.Len(*blobTagsSet, len(modifiedBlobTags))
	for _, blobTag := range *blobTagsSet {
		_assert.Equal(modifiedBlobTags[*blobTag.Key], *blobTag.Value)
	}
}

func (s *azblobTestSuite) TestPageBlobSetBlobTagForSnapshot() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	pbClient := createNewPageBlob(_assert, generateBlobName(testName), containerClient)

	setTagsBlobOptions := SetTagsBlobOptions{
		BlobTagsMap: &specialCharBlobTagsMap,
	}
	_, err = pbClient.SetTags(ctx, &setTagsBlobOptions)
	_assert.Nil(err)

	resp, err := pbClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)

	snapshotURL := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp2.TagCount, int64(len(specialCharBlobTagsMap)))

	blobGetTagsResponse, err := pbClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.Tags.BlobTagSet
	_assert.NotNil(blobTagsSet)
	_assert.Len(*blobTagsSet, len(specialCharBlobTagsMap))
	for _, blobTag := range *blobTagsSet {
		_assert.Equal(specialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *azblobTestSuite) TestCreateAppendBlobWithTags() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	_context := getTestContext(testName)
	ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	serviceClient, err := getServiceClient(_context.recording, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), serviceClient)
	defer deleteContainer(_assert, containerClient)

	abClient := getAppendBlobClient(generateBlobName(testName), containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobTagsMap: &specialCharBlobTagsMap,
	}
	createResp, err := abClient.Create(ctx, &createAppendBlobOptions)
	_assert.Nil(err)
	_assert.NotNil(createResp.VersionID)

	_, err = abClient.GetProperties(ctx, nil)
	_assert.Nil(err)

	blobGetTagsResponse, err := abClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.Tags.BlobTagSet
	_assert.NotNil(blobTagsSet)
	_assert.Len(*blobTagsSet, len(specialCharBlobTagsMap))
	for _, blobTag := range *blobTagsSet {
		_assert.Equal(specialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	resp, err := abClient.CreateSnapshot(ctx, nil)
	_assert.Nil(err)

	snapshotURL := abClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(*resp2.TagCount, int64(len(specialCharBlobTagsMap)))

	blobGetTagsResponse, err = abClient.GetTags(ctx, nil)
	_assert.Nil(err)
	_assert.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet = blobGetTagsResponse.Tags.BlobTagSet
	_assert.NotNil(blobTagsSet)
	_assert.Len(*blobTagsSet, len(specialCharBlobTagsMap))
	for _, blobTag := range *blobTagsSet {
		_assert.Equal(specialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

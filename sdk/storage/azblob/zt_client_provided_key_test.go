//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"io"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/require"
)

/*
Azure Storage supports following operations support of sending customer-provided encryption keys on a request:
Put Blob, Put Block List, Put Block, Put Block from URL, Put Page, Put Page from URL, Append Block,
Set Blob Properties, Set Blob Metadata, Get Blob, Get Blob Properties, Get Blob Metadata, Snapshot Blob.
*/

func (s *AZBlobRecordedTestsSuite) TestPutBlockAndPutBlockListWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = testcommon.BlockIDIntToBase64(index)
		stageBlockOptions := blockblob.StageBlockOptions{
			CpkInfo: &testcommon.TestCPKByValue,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_require.Nil(err)
	}

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_require.Nil(err)

	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(*resp.EncryptionKeySHA256, *(testcommon.TestCPKByValue.EncryptionKeySHA256))

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.DownloadStream(ctx, nil)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	getResp, err := bbClient.DownloadStream(ctx, &downloadBlobOptions)
	_require.Nil(err)
	b := bytes.Buffer{}
	_, err = b.ReadFrom(getResp.Body)
	_require.NoError(err)
	err = getResp.Body.Close()
	_require.NoError(err)
	_require.Equal(b.String(), "AAA BBB CCC ")
	_require.EqualValues(*getResp.ETag, *resp.ETag)
	_require.EqualValues(*getResp.LastModified, *resp.LastModified)
}

func (s *AZBlobRecordedTestsSuite) TestPutBlockAndPutBlockListWithCPKByScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = testcommon.BlockIDIntToBase64(index)
		stageBlockOptions := blockblob.StageBlockOptions{
			CpkScopeInfo: &testcommon.TestCPKByScope,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_require.Nil(err)
	}

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	_require.Nil(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(resp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)

	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	_, err = bbClient.DownloadStream(ctx, &downloadBlobOptions)
	_require.NotNil(err)

	downloadBlobOptions = blob.DownloadStreamOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	getResp, err := bbClient.DownloadStream(ctx, &downloadBlobOptions)
	_require.Nil(err)
	b := bytes.Buffer{}
	reader := getResp.Body
	_, err = b.ReadFrom(reader)
	_require.Nil(err)
	_ = reader.Close() // The client must close the response body when finished with it
	_require.Equal(b.String(), "AAA BBB CCC ")
	_require.EqualValues(*getResp.ETag, *resp.ETag)
	_require.EqualValues(*getResp.LastModified, *resp.LastModified)
	_require.Equal(*getResp.IsServerEncrypted, true)
	_require.EqualValues(*getResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
}

//
////nolint
//func (s *AZBlobUnrecordedTestsSuite) TestPutBlockFromURLAndCommitWithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024 // 8 KB
//	content := make([]byte, contentSize)
//	body := bytes.NewReader(content)
//	rsc := NopCloser(body)
//	ctx := ctx
//	srcBlob := containerClient.NewBlockBlobClient("srcblob")
//	destBlob := containerClient.NewBlockBlobClient("destblob")
//	_, err = srcBlob.Upload(ctx, rsc, nil)
//	_require.Nil(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
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
//	// Stage blocks from URL.
//	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
//	offset1, count1 := int64(0), int64(4*1024)
//	options1 := BlockBlobStageBlockFromURLOptions{
//		Offset:  &offset1,
//		Count:   &count1,
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.Nil(err)
//	// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp1.ContentMD5, "")
//	_require.NotEqual(stageResp1.RequestID, "")
//	_require.NotEqual(stageResp1.Version, "")
//	_require.Equal(stageResp1.Date.IsZero(), false)
//
//	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
//	offset2, count2 := int64(4*1024), int64(CountToEnd)
//	options2 := BlockBlobStageBlockFromURLOptions{
//		Offset:  &offset2,
//		Count:   &count2,
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.Nil(err)
//	// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp2.ContentMD5, "")
//	_require.NotEqual(stageResp2.RequestID, "")
//	_require.NotEqual(stageResp2.Version, "")
//	_require.Equal(stageResp2.Date.IsZero(), false)
//
//	// Check block list.
//	blockList, err := destBlob.GetBlockList(ctx, BlockListTypeAll, nil)
//	_require.Nil(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.NotNil(blockList.BlockList.UncommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	// Commit block list.
//	commitBlockListOptions := blockblob.CommitBlockListOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	listResp, err := destBlob.CommitBlockList(ctx, []string{blockID1, blockID2}, &commitBlockListOptions)
//	_require.Nil(err)
//	// _require.Equal(listResp.RawResponse.StatusCode,  201)
//	_require.NotNil(listResp.LastModified)
//	_require.Equal((*listResp.LastModified).IsZero(), false)
//	_require.NotNil(listResp.ETag)
//	_require.NotNil(listResp.RequestID)
//	_require.NotNil(listResp.Version)
//	_require.NotNil(listResp.Date)
//	_require.Equal((*listResp.Date).IsZero(), false)
//
//	// Check block list.
//	blockList, err = destBlob.GetBlockList(ctx, BlockListTypeAll, nil)
//	_require.Nil(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.UncommittedBlocks)
//	_require.NotNil(blockList.BlockList.CommittedBlocks)
//	_require.Len(blockList.BlockList.CommittedBlocks, 2)
//
//	// Check data integrity through downloading.
//	_, err = destBlob.BlobClient.DownloadStream(ctx, nil)
//	_require.NotNil(err)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.BlobClient.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.Body)
//	_require.Nil(err)
//	_require.EqualValues(destData, content)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestPutBlockFromURLAndCommitWithCPKWithScope() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024 // 8 KB
//	content := make([]byte, contentSize)
//	body := bytes.NewReader(content)
//	rsc := NopCloser(body)
//	srcBlob := containerClient.NewBlockBlobClient("srcblob")
//	destBlob := containerClient.NewBlockBlobClient("destblob")
//	_, err = srcBlob.Upload(ctx, rsc, nil)
//	_require.Nil(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
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
//	// Stage blocks from URL.
//	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
//	offset1, count1 := int64(0), int64(4*1024)
//	options1 := BlockBlobStageBlockFromURLOptions{
//		Offset:       &offset1,
//		Count:        &count1,
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.Nil(err)
//	// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp1.ContentMD5, "")
//	_require.NotEqual(stageResp1.RequestID, "")
//	_require.NotEqual(stageResp1.Version, "")
//	_require.Equal(stageResp1.Date.IsZero(), false)
//
//	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
//	offset2, count2 := int64(4*1024), int64(CountToEnd)
//	options2 := BlockBlobStageBlockFromURLOptions{
//		Offset:       &offset2,
//		Count:        &count2,
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.Nil(err)
//	//_require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp2.ContentMD5, "")
//	_require.NotEqual(stageResp2.RequestID, "")
//	_require.NotEqual(stageResp2.Version, "")
//	_require.Equal(stageResp2.Date.IsZero(), false)
//
//	// Check block list.
//	blockList, err := destBlob.GetBlockList(ctx, BlockListTypeAll, nil)
//	_require.Nil(err)
//	//_require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.NotNil(blockList.BlockList.UncommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	// Commit block list.
//	commitBlockListOptions := blockblob.CommitBlockListOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	listResp, err := destBlob.CommitBlockList(ctx, []string{blockID1, blockID2}, &commitBlockListOptions)
//	_require.Nil(err)
//	//_require.Equal(listResp.RawResponse.StatusCode, 201)
//	_require.NotNil(listResp.LastModified)
//	_require.Equal((*listResp.LastModified).IsZero(), false)
//	_require.NotNil(listResp.ETag)
//	_require.NotNil(listResp.RequestID)
//	_require.NotNil(listResp.Version)
//	_require.NotNil(listResp.Date)
//	_require.Equal((*listResp.Date).IsZero(), false)
//
//	// Check block list.
//	blockList, err = destBlob.GetBlockList(ctx, BlockListTypeAll, nil)
//	_require.Nil(err)
//	//_require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.UncommittedBlocks)
//	_require.NotNil(blockList.BlockList.CommittedBlocks)
//	_require.Len(blockList.BlockList.CommittedBlocks, 2)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := destBlob.BlobClient.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.Body)
//	_require.Nil(err)
//	_require.EqualValues(destData, content)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
//}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadBlobWithMD5WithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 8 * 1024
	r, srcData := testcommon.GenerateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	uploadBlockBlobOptions := blockblob.UploadOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.Equal(*uploadResp.IsServerEncrypted, true)
	_require.EqualValues(uploadResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.DownloadStream(ctx, nil)
	_require.NotNil(err)

	_, err = bbClient.DownloadStream(ctx, &blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestInvalidCPKByValue,
	})
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadResp, err := bbClient.DownloadStream(ctx, &blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	})
	_require.Nil(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(downloadResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *AZBlobRecordedTestsSuite) TestUploadBlobWithMD5WithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 8 * 1024
	r, srcData := testcommon.GenerateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	uploadBlockBlobOptions := blockblob.UploadOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	uploadResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.Equal(*uploadResp.IsServerEncrypted, true)
	_require.EqualValues(uploadResp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	downloadResp, err := bbClient.DownloadStream(ctx, &downloadBlobOptions)
	_require.Nil(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
}

func (s *AZBlobRecordedTestsSuite) TestAppendBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := appendblob.AppendBlockOptions{
			CpkInfo: &testcommon.TestCPKByValue,
		}
		resp, err := abClient.AppendBlock(ctx, streaming.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal((*resp.Date).IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = abClient.DownloadStream(ctx, nil)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	downloadResp, err := abClient.DownloadStream(ctx, &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *AZBlobRecordedTestsSuite) TestAppendBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := appendblob.AppendBlockOptions{
			CpkScopeInfo: &testcommon.TestCPKByScope,
		}
		resp, err := abClient.AppendBlock(ctx, streaming.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.Nil(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal((*resp.Date).IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(resp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	downloadResp, err := abClient.DownloadStream(ctx, &downloadBlobOptions)
	_require.Nil(err)

	data, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	ctx := ctx
//	srcABClient := containerClient.NewAppendBlobClient(generateName("src"))
//	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))
//
//	_, err = srcABClient.Create(ctx, nil)
//	_require.Nil(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	resp, err := srcABClient.AppendBlock(ctx, streaming.NopCloser(r), nil)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.Equal(*resp.BlobAppendOffset, "0")
//	_require.Equal(*resp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.Equal((*resp.LastModified).IsZero(), false)
//	_require.Nil(resp.ContentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//
//	srcBlobParts, _ := NewBlobURLParts(srcABClient.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
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
//	createAppendBlobOptions := appendblob.CreateOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	_, err = destBlob.Create(ctx, &createAppendBlobOptions)
//	_require.Nil(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
//		Offset:  &offset,
//		Count:   &count,
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.Nil(err)
//	//_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
//	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendFromURLResp.ETag)
//	_require.NotNil(appendFromURLResp.LastModified)
//	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
//	_require.NotNil(appendFromURLResp.ContentMD5)
//	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
//	_require.NotNil(appendFromURLResp.RequestID)
//	_require.NotNil(appendFromURLResp.Version)
//	_require.NotNil(appendFromURLResp.Date)
//	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
//	_require.Equal(*appendFromURLResp.IsServerEncrypted, true)
//
//	// Get blob content without encryption key should fail the request.
//	_, err = destBlob.DownloadStream(ctx, nil)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//
//	_require.Equal(*downloadResp.IsServerEncrypted, true)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCPKScope() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	ctx := ctx
//	srcClient := containerClient.NewAppendBlobClient(generateName("src"))
//	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))
//
//	_, err = srcClient.Create(ctx, nil)
//	_require.Nil(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	resp, err := srcClient.AppendBlock(ctx, streaming.NopCloser(r), nil)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.Equal(*resp.BlobAppendOffset, "0")
//	_require.Equal(*resp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.Equal((*resp.LastModified).IsZero(), false)
//	_require.Nil(resp.ContentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//
//	srcBlobParts, _ := NewBlobURLParts(srcClient.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
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
//	createAppendBlobOptions := appendblob.CreateOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	_, err = destBlob.Create(ctx, &createAppendBlobOptions)
//	_require.Nil(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
//		Offset:       &offset,
//		Count:        &count,
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.Nil(err)
//	//_require.Equal(appendFromURLResp.RawResponse.StatusCode, 201)
//	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
//	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
//	_require.NotNil(appendFromURLResp.ETag)
//	_require.NotNil(appendFromURLResp.LastModified)
//	_require.Equal((*appendFromURLResp.LastModified).IsZero(), false)
//	_require.NotNil(appendFromURLResp.ContentMD5)
//	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5)
//	_require.NotNil(appendFromURLResp.RequestID)
//	_require.NotNil(appendFromURLResp.Version)
//	_require.NotNil(appendFromURLResp.Date)
//	_require.Equal((*appendFromURLResp.Date).IsZero(), false)
//	_require.Equal(*appendFromURLResp.IsServerEncrypted, true)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.Equal(*downloadResp.IsServerEncrypted, true)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestPageBlockFromURLWithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024 // 1MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	ctx := ctx // Use default Background context
//	srcPBName := "src" + testcommon.GenerateBlobName(testName)
//	bbClient := createNewPageBlobWithSize(_require, srcPBName, containerClient, int64(contentSize))
//	dstPBName := "dst" + testcommon.GenerateBlobName(testName)
//	destBlob := createNewPageBlobWithCPK(_require, dstPBName, containerClient, int64(contentSize), &testcommon.TestCPKByValue, nil)
//
//	offset, count := int64(0), int64(contentSize)
//	uploadPagesOptions := pageblob.UploadPagesOptions{
//		Offset: to.Ptr(offset), Count: to.Ptr(count),
//	}
//	_, err = bbClient.UploadPages(ctx, streaming.NopCloser(r), &uploadPagesOptions)
//	_require.Nil(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//	srcBlobParts, _ := NewBlobURLParts(bbClient.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: contentMD5,
//		CpkInfo:          &testcommon.TestCPKByValue,
//	}
//	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.NotNil(resp.ContentMD5)
//	_require.EqualValues(resp.ContentMD5, contentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//	_require.Equal(*resp.BlobSequenceNumber, int64(0))
//	_require.Equal(*resp.IsServerEncrypted, true)
//	_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	_, err = destBlob.DownloadStream(ctx, nil)
//	_require.NotNil(err)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestPageBlockFromURLWithCPKScope() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024 // 1MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	ctx := ctx // Use default Background context
//	srcPBName := "src" + testcommon.GenerateBlobName(testName)
//	srcPBClient := createNewPageBlobWithSize(_require, srcPBName, containerClient, int64(contentSize))
//	dstPBName := "dst" + testcommon.GenerateBlobName(testName)
//	dstPBBlob := createNewPageBlobWithCPK(_require, dstPBName, containerClient, int64(contentSize), nil, &testcommon.TestCPKByScope)
//
//	offset, count := int64(0), int64(contentSize)
//	uploadPagesOptions := pageblob.UploadPagesOptions{
//		Offset: to.Ptr(offset), Count: to.Ptr(count),
//	}
//	_, err = srcPBClient.UploadPages(ctx, streaming.NopCloser(r), &uploadPagesOptions)
//	_require.Nil(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//	srcBlobParts, _ := NewBlobURLParts(srcPBClient.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: contentMD5,
//		CpkScopeInfo:     &testcommon.TestCPKByScope,
//	}
//	resp, err := dstPBBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.NotNil(resp.ContentMD5)
//	_require.EqualValues(resp.ContentMD5, contentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//	_require.Equal(*resp.BlobSequenceNumber, int64(0))
//	_require.Equal(*resp.IsServerEncrypted, true)
//	_require.EqualValues(resp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := dstPBBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestUploadPagesFromURLWithMD5WithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Sum := md5.Sum(srcData)
//	contentMD5 := md5Sum[:]
//	srcPBName := "src" + testcommon.GenerateBlobName(testName)
//	srcBlob := createNewPageBlobWithSize(_require, srcPBName, containerClient, int64(contentSize))
//
//	offset, count := int64(0), int64(contentSize)
//	uploadPagesOptions := pageblob.UploadPagesOptions{
//		Offset: to.Ptr(offset), Count: to.Ptr(count),
//	}
//	_, err = srcBlob.UploadPages(ctx, streaming.NopCloser(r), &uploadPagesOptions)
//	_require.Nil(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//	dstPBName := "dst" + testcommon.GenerateBlobName(testName)
//	destPBClient := createNewPageBlobWithCPK(_require, dstPBName, containerClient, int64(contentSize), &testcommon.TestCPKByValue, nil)
//	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: contentMD5,
//		CpkInfo:          &testcommon.TestCPKByValue,
//	}
//	resp, err := destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.Nil(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 201)
//	_require.NotNil(resp.ETag)
//	_require.NotNil(resp.LastModified)
//	_require.NotNil(resp.ContentMD5)
//	_require.EqualValues(resp.ContentMD5, contentMD5)
//	_require.NotNil(resp.RequestID)
//	_require.NotNil(resp.Version)
//	_require.NotNil(resp.Date)
//	_require.Equal((*resp.Date).IsZero(), false)
//	_require.Equal(*resp.BlobSequenceNumber, int64(0))
//	_require.Equal(*resp.IsServerEncrypted, true)
//	_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	_, err = destPBClient.DownloadStream(ctx, nil)
//	_require.NotNil(err)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destPBClient.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destPBClient.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CpkInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
//
//	_, badMD5 := getRandomDataAndReader(16)
//	badContentMD5 := badMD5[:]
//	uploadPagesFromURLOptions1 := pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: badContentMD5,
//	}
//	_, err = destPBClient.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions1)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, StorageErrorCodeMD5Mismatch)
//}

//func (s *AZBlobRecordedTestsSuite) TestClearDiffPagesWithCPK() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	pbName := testcommon.GenerateBlobName(testName)
//	pbClient := createNewPageBlobWithCPK(_require, pbName, containerClient, pageblob.PageBytes*10, &testcommon.TestCPKByValue, nil)
//
//	contentSize := 2 * 1024
//	r := getReaderToGeneratedBytes(contentSize)
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := pageblob.UploadPagesOptions{Range: &HttpRange{offset, count}, CpkInfo: &testcommon.TestCPKByValue}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_require.Nil(err)
//
//	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	snapshotResp, err := pbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
//	_require.Nil(err)
//
//	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
//	uploadPagesOptions1 := pageblob.UploadPagesOptions{Range: &HttpRange{offset1, count1}, CpkInfo: &testcommon.TestCPKByValue}
//	_, err = pbClient.UploadPages(ctx, getReaderToGeneratedBytes(2048), &uploadPagesOptions1)
//	_require.Nil(err)
//
//	pageListResp, err := pbClient.NewGetPageRangesDiffPager(ctx, HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
//	_require.Nil(err)
//	pageRangeResp := pageListResp.PageList.Range
//	_require.NotNil(pageRangeResp)
//	_require.Len(pageRangeResp, 1)
//	rawStart, rawEnd := pageRangeResp[0].Raw()
//	_require.Equal(rawStart, offset1)
//	_require.Equal(rawEnd, end1)
//
//	clearPagesOptions := PageBlobClearPagesOptions{
//		CpkInfo: &testcommon.TestCPKByValue,
//	}
//	clearResp, err := pbClient.ClearPages(ctx, HttpRange{2048, 2048}, &clearPagesOptions)
//	_require.Nil(err)
//	_require.Equal(clearResp.RawResponse.StatusCode, 201)
//
//	pageListResp, err = pbClient.NewGetPageRangesDiffPager(ctx, HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
//	_require.Nil(err)
//	_require.Nil(pageListResp.PageList.Range)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlobBlobPropertiesWithCPKKey() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	_require.NoError(err)
//
//	blobSize := 1024
//	bufferSize := 8 * 1024
//	maxBuffers := 3
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	// Set up test blob
//	blobName := testcommon.GenerateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//	_require.Nil(err)
//
//	// Create some data to test the upload stream
//	blobContentReader, blobData := testcommon.GenerateData(blobSize)
//
//	// Perform UploadStream
//	_, err = bbClient.UploadStream(ctx, blobContentReader,
//		&UploadStreamOptions{
//			BufferSize:  bufferSize,
//			MaxBuffers:  maxBuffers,
//			Metadata:    testcommon.BasicMetadata,
//			BlobTags:    basicBlobTagsMap,
//			HTTPHeaders: &basicHeaders,
//			CpkInfo:     &testcommon.TestCPKByValue,
//		})
//
//	// Assert that upload was successful
//	_require.Equal(err, nil)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	getPropertiesResp, err := bbClient.GetProperties(ctx, &blob.GetPropertiesOptions{CpkInfo: &testcommon.TestCPKByValue})
//	_require.NoError(err)
//	_require.EqualValues(getPropertiesResp.Metadata, testcommon.BasicMetadata)
//	_require.Equal(*getPropertiesResp.TagCount, int64(len(basicBlobTagsMap)))
//	_require.Equal(blob.ParseHTTPHeaders(getPropertiesResp), basicHeaders)
//
//	getTagsResp, err := bbClient.GetTags(ctx, nil)
//	_require.NoError(err)
//	_require.Len(getTagsResp.BlobTagSet, 3)
//	for _, blobTag := range getTagsResp.BlobTagSet {
//		_require.Equal(basicBlobTagsMap[*blobTag.Key], *blobTag.Value)
//	}
//
//	// Download the blob to verify
//	downloadResponse, err := bbClient.DownloadStream(ctx, &blob.downloadWriterAtOptions{CpkInfo: &testcommon.TestCPKByValue})
//	_require.NoError(err)
//
//	// Assert that the content is correct
//	actualBlobData, err := io.ReadAll(downloadResponse.Body(nil))
//	_require.NoError(err)
//	_require.Equal(len(actualBlobData), blobSize)
//	_require.EqualValues(actualBlobData, blobData)
//}

//nolint
//func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlobBlobPropertiesWithCPKScope() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//	_require.NoError(err)
//
//	blobSize := 1024
//	bufferSize := 8 * 1024
//	maxBuffers := 3
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	// Set up test blob
//	blobName := testcommon.GenerateBlobName(testName)
//	bbClient := getBlockBlobClient(blobName, containerClient)
//	_require.NoError(err)
//
//	// Create some data to test the upload stream
//	blobContentReader, blobData := testcommon.GenerateData(blobSize)
//
//	// Perform UploadStream
//	_, err = bbClient.UploadStream(ctx, blobContentReader,
//		&UploadStreamOptions{
//			BufferSize:   bufferSize,
//			MaxBuffers:   maxBuffers,
//			Metadata:     testcommon.BasicMetadata,
//			BlobTags:     basicBlobTagsMap,
//			HTTPHeaders:  &basicHeaders,
//			CpkScopeInfo: &testcommon.TestCPKByScope,
//		})
//
//	// Assert that upload was successful
//	_require.Equal(err, nil)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	getPropertiesResp, err := bbClient.GetProperties(ctx, nil)
//	_require.NoError(err)
//	_require.EqualValues(getPropertiesResp.Metadata, testcommon.BasicMetadata)
//	_require.Equal(*getPropertiesResp.TagCount, int64(len(basicBlobTagsMap)))
//	_require.Equal(blob.ParseHTTPHeaders(getPropertiesResp), basicHeaders)
//
//	getTagsResp, err := bbClient.GetTags(ctx, nil)
//	_require.NoError(err)
//	_require.Len(getTagsResp.BlobTagSet, 3)
//	for _, blobTag := range getTagsResp.BlobTagSet {
//		_require.Equal(basicBlobTagsMap[*blobTag.Key], *blobTag.Value)
//	}
//
//	// Download the blob to verify
//	downloadResponse, err := bbClient.DownloadStream(ctx, &blob.downloadWriterAtOptions{CpkScopeInfo: &testcommon.TestCPKByScope})
//	_require.NoError(err)
//
//	// Assert that the content is correct
//	actualBlobData, err := io.ReadAll(downloadResponse.Body(nil))
//	_require.NoError(err)
//	_require.Equal(len(actualBlobData), blobSize)
//	_require.EqualValues(actualBlobData, blobData)
//}

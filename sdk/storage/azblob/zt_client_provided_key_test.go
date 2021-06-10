// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/ioutil"
	"strconv"
	"time"

	//"encoding/binary"
	chk "gopkg.in/check.v1"
	"strings"
)

/*
Azure Storage supports following operations support of sending customer-provided encryption keys on a request:
Put Blob, Put Block List, Put Block, Put Block from URL, Put Page, Put Page from URL, Append Block,
Set Blob Properties, Set Blob Metadata, Get Blob, Get Blob Properties, Get Blob Metadata, Snapshot Blob.
*/
var testEncryptedKey = "MDEyMzQ1NjcwMTIzNDU2NzAxMjM0NTY3MDEyMzQ1Njc="
var testEncryptedHash = "3QFFFpRA5+XANHqwwbT4yXDmrT/2JaLt/FKHjzhOdoE="
var testEncryptionAlgorithm = "AES256"
var testCPKByValue = CpkInfo{
	EncryptionKey:       &testEncryptedKey,
	EncryptionKeySHA256: &testEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

var testInvalidEncryptedKey = "mumbojumbo"
var testInvalidEncryptedHash = "mumbojumbohash"
var testInvalidCPKByValue = CpkInfo{
	EncryptionKey:       &testInvalidEncryptedKey,
	EncryptionKeySHA256: &testInvalidEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

var testEncryptedScope = "blobgokeytestscope"
var testCPKByScope = CpkScopeInfo{
	EncryptionScope: &testEncryptedScope,
}

var testInvalidEncryptedScope = "mumbojumbo"
var testInvalidCPKByScope = CpkScopeInfo{
	EncryptionScope: &testInvalidEncryptedScope,
}

func (s *aztestsSuite) TestPutBlockAndPutBlockListWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := StageBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], strings.NewReader(word), &stageBlockOptions)
		c.Assert(err, chk.IsNil)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)

	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
	c.Assert(*resp.EncryptionKeySHA256, chk.DeepEquals, *(testCPKByValue.EncryptionKeySHA256))

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	b := bytes.Buffer{}
	reader := getResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue})
	_, _ = b.ReadFrom(reader)
	_ = reader.Close()
	c.Assert(b.String(), chk.Equals, "AAA BBB CCC ")
	c.Assert(*getResp.ETag, chk.DeepEquals, *resp.ETag)
	c.Assert(*getResp.LastModified, chk.DeepEquals, *resp.LastModified)
}

func (s *aztestsSuite) TestPutBlockAndPutBlockListWithCPKByScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		stageBlockOptions := StageBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], strings.NewReader(word), &stageBlockOptions)
		c.Assert(err, chk.IsNil)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
	c.Assert(resp.EncryptionScope, chk.DeepEquals, testCPKByScope.EncryptionScope)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.NotNil)

	//storageErr := err.(*StorageError)
	//c.Assert(storageErr.RawResponse.StatusCode, chk.Equals, 409)
	//c.Assert(storageErr.ErrorCode, chk.Equals, StorageErrorCodeFeatureEncryptionMismatch)

	downloadBlobOptions = DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	getResp, err = bbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	b := bytes.Buffer{}
	reader := getResp.Body(RetryReaderOptions{})
	_, err = b.ReadFrom(reader)
	c.Assert(err, chk.IsNil)
	_ = reader.Close() // The client must close the response body when finished with it
	c.Assert(b.String(), chk.Equals, "AAA BBB CCC ")
	c.Assert(*getResp.ETag, chk.DeepEquals, *resp.ETag)
	c.Assert(*getResp.LastModified, chk.DeepEquals, *resp.LastModified)
	c.Assert(*getResp.IsServerEncrypted, chk.Equals, true)
	c.Assert(*getResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)
}

func (s *aztestsSuite) TestPutBlockFromURLAndCommitWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential("")
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	offset1, count1 := int64(0), int64(4*1024)
	options1 := StageBlockFromURLOptions{
		Offset:  &offset1,
		Count:   &count1,
		CpkInfo: &testCPKByValue,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp1.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Date.IsZero(), chk.Equals, false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset:  &offset2,
		Count:   &count2,
		CpkInfo: &testCPKByValue,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp2.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp2.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Date.IsZero(), chk.Equals, false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.UncommittedBlocks, chk.HasLen, 2)

	// Commit block list.
	commitBlockListOptions := CommitBlockListOptions{
		CpkInfo: &testCPKByValue,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(listResp.LastModified, chk.NotNil)
	c.Assert((*listResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(listResp.ETag, chk.NotNil)
	c.Assert(listResp.RequestID, chk.NotNil)
	c.Assert(listResp.Version, chk.NotNil)
	c.Assert(listResp.Date, chk.NotNil)
	c.Assert((*listResp.Date).IsZero(), chk.Equals, false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.CommittedBlocks, chk.HasLen, 2)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err = destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, content)
	c.Assert(*downloadResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)
}

func (s *aztestsSuite) TestPutBlockFromURLAndCommitWithCPKWithScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential("")
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	offset1, count1 := int64(0), int64(4*1024)
	options1 := StageBlockFromURLOptions{
		Offset:       &offset1,
		Count:        &count1,
		CpkScopeInfo: &testCPKByScope,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp1.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Date.IsZero(), chk.Equals, false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset:       &offset2,
		Count:        &count2,
		CpkScopeInfo: &testCPKByScope,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp2.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp2.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Date.IsZero(), chk.Equals, false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.UncommittedBlocks, chk.HasLen, 2)

	// Commit block list.
	commitBlockListOptions := CommitBlockListOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(listResp.LastModified, chk.NotNil)
	c.Assert((*listResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(listResp.ETag, chk.NotNil)
	c.Assert(listResp.RequestID, chk.NotNil)
	c.Assert(listResp.Version, chk.NotNil)
	c.Assert(listResp.Date, chk.NotNil)
	c.Assert((*listResp.Date).IsZero(), chk.Equals, false)

	// Check block list.
	blockList, err = destBlob.GetBlockList(context.Background(), BlockListAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.CommittedBlocks, chk.HasLen, 2)

	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.BlobClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, content)
	c.Assert(*downloadResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)
}

func (s *aztestsSuite) TestUploadBlobWithMD5WithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	uploadSrcResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*uploadSrcResp.IsServerEncrypted, chk.Equals, true)
	c.Assert(uploadSrcResp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)

	// Get blob content without encryption key should fail the request.
	downloadResp, err := bbClient.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	downloadResp, err = bbClient.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err = bbClient.BlobClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*downloadResp.ContentMD5, chk.DeepEquals, md5Val[:])
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
	c.Assert(downloadResp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)
}

func (s *aztestsSuite) TestUploadBlobWithMD5WithCPKScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	uploadSrcResp, err := bbClient.Upload(ctx, r, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*uploadSrcResp.IsServerEncrypted, chk.Equals, true)
	c.Assert(uploadSrcResp.EncryptionScope, chk.DeepEquals, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := bbClient.BlobClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*downloadResp.ContentMD5, chk.DeepEquals, md5Val[:])
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
	c.Assert(*downloadResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)
}

func (s *aztestsSuite) TestAppendBlockWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	appendBlobURL := containerClient.NewAppendBlobClient(generateBlobName())

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := appendBlobURL.Create(context.Background(), &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlockOptions{
			CpkInfo: &testCPKByValue,
		}
		resp, err := appendBlobURL.AppendBlock(context.Background(), strings.NewReader(word), &appendBlockOptions)
		c.Assert(err, chk.IsNil)
		c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
		c.Assert(*resp.BlobAppendOffset, chk.Equals, strconv.Itoa(index*4))
		c.Assert(*resp.BlobCommittedBlockCount, chk.Equals, int32(index+1))
		c.Assert(resp.ETag, chk.NotNil)
		c.Assert(resp.LastModified, chk.NotNil)
		c.Assert(resp.LastModified.IsZero(), chk.Equals, false)
		c.Assert(resp.ContentMD5, chk.Not(chk.Equals), "")
		c.Assert(resp.RequestID, chk.Not(chk.Equals), "")
		c.Assert(resp.Version, chk.Not(chk.Equals), "")
		c.Assert(resp.Date, chk.NotNil)
		c.Assert((*resp.Date).IsZero(), chk.Equals, false)
		c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
		c.Assert(resp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = appendBlobURL.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err := appendBlobURL.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)

	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.DeepEquals, "AAA BBB CCC ")
	c.Assert(*downloadResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)
}

func (s *aztestsSuite) TestAppendBlockWithCPKScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	appendBlobURL := containerClient.NewAppendBlobClient(generateBlobName())

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := appendBlobURL.Create(context.Background(), &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := AppendBlockOptions{
			CpkScopeInfo: &testCPKByScope,
		}
		resp, err := appendBlobURL.AppendBlock(context.Background(), strings.NewReader(word), &appendBlockOptions)
		c.Assert(err, chk.IsNil)
		c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
		c.Assert(*resp.BlobAppendOffset, chk.Equals, strconv.Itoa(index*4))
		c.Assert(*resp.BlobCommittedBlockCount, chk.Equals, int32(index+1))
		c.Assert(resp.ETag, chk.NotNil)
		c.Assert(resp.LastModified, chk.NotNil)
		c.Assert(resp.LastModified.IsZero(), chk.Equals, false)
		c.Assert(resp.ContentMD5, chk.Not(chk.Equals), "")
		c.Assert(resp.RequestID, chk.Not(chk.Equals), "")
		c.Assert(resp.Version, chk.Not(chk.Equals), "")
		c.Assert(resp.Date, chk.NotNil)
		c.Assert((*resp.Date).IsZero(), chk.Equals, false)
		c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
		c.Assert(resp.EncryptionScope, chk.DeepEquals, testCPKByScope.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := appendBlobURL.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)

	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.DeepEquals, "AAA BBB CCC ")
	c.Assert(*downloadResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)
}

func (s *aztestsSuite) TestAppendBlockFromURLWithCPK(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(4 * 1024 * 1024) // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcClient := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcClient.Create(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp1.RawResponse.StatusCode, chk.Equals, 201)

	resp, err := srcClient.AppendBlock(context.Background(), r, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*resp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*resp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert((*resp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(resp.ContentMD5, chk.IsNil)
	c.Assert(resp.RequestID, chk.NotNil)
	c.Assert(resp.Version, chk.NotNil)
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)

	srcBlobParts := NewBlobURLParts(srcClient.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp2.RawResponse.StatusCode, chk.Equals, 201)

	offset := int64(0)
	count := contentSize
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:  &offset,
		Count:   &count,
		CpkInfo: &testCPKByValue,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(appendFromURLResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendFromURLResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendFromURLResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendFromURLResp.ETag, chk.NotNil)
	c.Assert(appendFromURLResp.LastModified, chk.NotNil)
	c.Assert((*appendFromURLResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendFromURLResp.ContentMD5, chk.NotNil)
	c.Assert(*appendFromURLResp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(appendFromURLResp.RequestID, chk.NotNil)
	c.Assert(appendFromURLResp.Version, chk.NotNil)
	c.Assert(appendFromURLResp.Date, chk.NotNil)
	c.Assert((*appendFromURLResp.Date).IsZero(), chk.Equals, false)
	c.Assert(*appendFromURLResp.IsServerEncrypted, chk.Equals, true)

	// Get blob content without encryption key should fail the request.
	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	downloadResp, err = destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err = destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)

	c.Assert(*downloadResp.IsServerEncrypted, chk.Equals, true)
	c.Assert(*downloadResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
}

func (s *aztestsSuite) TestAppendBlockFromURLWithCPKScope(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(4 * 1024 * 1024) // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background()
	srcClient := containerClient.NewAppendBlobClient(generateName("src"))
	destBlob := containerClient.NewAppendBlobClient(generateName("dest"))

	cResp1, err := srcClient.Create(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp1.RawResponse.StatusCode, chk.Equals, 201)

	resp, err := srcClient.AppendBlock(context.Background(), r, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*resp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*resp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert((*resp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(resp.ContentMD5, chk.IsNil)
	c.Assert(resp.RequestID, chk.NotNil)
	c.Assert(resp.Version, chk.NotNil)
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)

	srcBlobParts := NewBlobURLParts(srcClient.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	createAppendBlobOptions := CreateAppendBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	cResp2, err := destBlob.Create(context.Background(), &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp2.RawResponse.StatusCode, chk.Equals, 201)

	offset := int64(0)
	count := contentSize
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:       &offset,
		Count:        &count,
		CpkScopeInfo: &testCPKByScope,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(appendFromURLResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendFromURLResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendFromURLResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendFromURLResp.ETag, chk.NotNil)
	c.Assert(appendFromURLResp.LastModified, chk.NotNil)
	c.Assert((*appendFromURLResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendFromURLResp.ContentMD5, chk.NotNil)
	c.Assert(*appendFromURLResp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(appendFromURLResp.RequestID, chk.NotNil)
	c.Assert(appendFromURLResp.Version, chk.NotNil)
	c.Assert(appendFromURLResp.Date, chk.NotNil)
	c.Assert((*appendFromURLResp.Date).IsZero(), chk.Equals, false)
	c.Assert(*appendFromURLResp.IsServerEncrypted, chk.Equals, true)

	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*downloadResp.IsServerEncrypted, chk.Equals, true)
	c.Assert(*downloadResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
}

func createNewPageBlobWithCPK(c *chk.C, container ContainerClient, sizeInBytes int64, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (pbClient PageBlobClient, name string) {
	pbClient, name = getPageBlobClient(c, container)

	createPageBlobOptions := CreatePageBlobOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	}
	resp, err := pbClient.Create(ctx, sizeInBytes, &createPageBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	return
}

func (s *aztestsSuite) TestPageBlockWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(4 * 1024 * 1024) // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	pbClient, _ := createNewPageBlobWithCPK(c, containerClient, contentSize, &testCPKByValue, nil)

	offset, count := int64(0), contentSize
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
		CpkInfo:   &testCPKByValue,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(uploadResp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	c.Assert(err, chk.IsNil)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), contentSize-1
	rawStart, rawEnd := (*pageListResp)[0].Raw()
	c.Assert(rawStart, chk.Equals, start)
	c.Assert(rawEnd, chk.Equals, end)

	// Get blob content without encryption key should fail the request.
	downloadResp, err := pbClient.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	downloadResp, err = pbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err = pbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
	c.Assert(*downloadResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)
}

func (s *aztestsSuite) TestPageBlockWithCPKScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(4 * 1024 * 1024) // 4MB
	r, srcData := getRandomDataAndReader(contentSize)
	pbClient, _ := createNewPageBlobWithCPK(c, containerClient, contentSize, nil, &testCPKByScope)

	offset, count := int64(0), contentSize
	uploadPagesOptions := UploadPagesOptions{
		PageRange:    &HttpRange{offset, count},
		CpkScopeInfo: &testCPKByScope,
	}
	uploadResp, err := pbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(uploadResp.EncryptionScope, chk.DeepEquals, testCPKByScope.EncryptionScope)

	resp, err := pbClient.GetPageRanges(ctx, HttpRange{0, CountToEnd}, nil)
	c.Assert(err, chk.IsNil)
	pageListResp := resp.PageList.PageRange
	start, end := int64(0), contentSize-1
	// c.Assert((*pageListResp)[0], chk.DeepEquals, PageRange{Start: &start, End: &end})
	rawStart, rawEnd := (*pageListResp)[0].Raw()
	c.Assert(rawStart, chk.Equals, start)
	c.Assert(rawEnd, chk.Equals, end)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := pbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
	c.Assert(*downloadResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)
}

func (s *aztestsSuite) TestPageBlockFromURLWithCPK(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	bbClient, _ := createNewPageBlobWithSize(c, containerClient, contentSize)
	destBlob, _ := createNewPageBlobWithCPK(c, containerClient, contentSize, &testCPKByValue, nil)

	offset, count := int64(0), contentSize
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := bbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadResp.RawResponse.StatusCode, chk.Equals, 201)
	srcBlobParts := NewBlobURLParts(bbClient.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: &contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, contentSize, &uploadPagesFromURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert(resp.ContentMD5, chk.NotNil)
	c.Assert(*resp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(resp.RequestID, chk.NotNil)
	c.Assert(resp.Version, chk.NotNil)
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)
	c.Assert(*resp.BlobSequenceNumber, chk.Equals, int64(0))
	c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
	c.Assert(resp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)

	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	downloadResp, err = destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err = destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*downloadResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
}

func (s *aztestsSuite) TestPageBlockFromURLWithCPKScope(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	ctx := context.Background() // Use default Background context
	bbClient, _ := createNewPageBlobWithSize(c, containerClient, contentSize)
	destBlob, _ := createNewPageBlobWithCPK(c, containerClient, contentSize, nil, &testCPKByScope)

	offset, count := int64(0), contentSize
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadResp, err := bbClient.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadResp.RawResponse.StatusCode, chk.Equals, 201)
	srcBlobParts := NewBlobURLParts(bbClient.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: &contentMD5,
		CpkScopeInfo:     &testCPKByScope,
	}
	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, contentSize, &uploadPagesFromURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert(resp.ContentMD5, chk.NotNil)
	c.Assert(*resp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(resp.RequestID, chk.NotNil)
	c.Assert(resp.Version, chk.NotNil)
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)
	c.Assert(*resp.BlobSequenceNumber, chk.Equals, int64(0))
	c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
	c.Assert(resp.EncryptionScope, chk.DeepEquals, testCPKByScope.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	downloadResp, err := destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*downloadResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)
}

func (s *aztestsSuite) TestUploadPagesFromURLWithMD5WithCPK(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := int64(8 * 1024) // 8 KB
	r, srcData := getRandomDataAndReader(contentSize)
	md5Sum := md5.Sum(srcData)
	contentMD5 := md5Sum[:]
	srcBlob, _ := createNewPageBlobWithSize(c, containerClient, contentSize)

	offset, count := int64(0), contentSize
	uploadPagesOptions := UploadPagesOptions{
		PageRange: &HttpRange{offset, count},
	}
	uploadSrcResp1, err := srcBlob.UploadPages(ctx, r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp1.RawResponse.StatusCode, chk.Equals, 201)

	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	destBlob, _ := createNewPageBlobWithCPK(c, containerClient, contentSize, &testCPKByValue, nil)
	uploadPagesFromURLOptions := UploadPagesFromURLOptions{
		SourceContentMD5: &contentMD5,
		CpkInfo:          &testCPKByValue,
	}
	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, contentSize, &uploadPagesFromURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert(resp.ContentMD5, chk.NotNil)
	c.Assert(*resp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(resp.RequestID, chk.NotNil)
	c.Assert(resp.Version, chk.NotNil)
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)
	c.Assert(*resp.BlobSequenceNumber, chk.Equals, int64(0))
	c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
	c.Assert(resp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)

	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.NotNil)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	downloadResp, err = destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.NotNil)

	// Download blob to do data integrity check.
	downloadBlobOptions = DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	downloadResp, err = destBlob.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*downloadResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)

	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{CpkInfo: &testCPKByValue}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, srcData)

	_, badMD5 := getRandomDataAndReader(16)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions1 := UploadPagesFromURLOptions{
		SourceContentMD5: &badContentMD5,
	}
	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, contentSize, &uploadPagesFromURLOptions1)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeMD5Mismatch)
}

func (s *aztestsSuite) TestClearDiffPagesWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	pbClient, _ := createNewPageBlobWithCPK(c, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	contentSize := int64(2 * 1024)
	r := getReaderToRandomBytes(contentSize)
	offset, _, count := int64(0), contentSize-1, contentSize
	uploadPagesOptions := UploadPagesOptions{PageRange: &HttpRange{offset, count}, CpkInfo: &testCPKByValue}
	_, err := pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	c.Assert(err, chk.IsNil)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	c.Assert(err, chk.IsNil)

	offset1, end1, count1 := contentSize, 2*contentSize-1, contentSize
	uploadPagesOptions1 := UploadPagesOptions{PageRange: &HttpRange{offset1, count1}, CpkInfo: &testCPKByValue}
	_, err = pbClient.UploadPages(context.Background(), getReaderToRandomBytes(2048), &uploadPagesOptions1)
	c.Assert(err, chk.IsNil)

	pageListResp, err := pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4096}, *snapshotResp.Snapshot, nil)
	c.Assert(err, chk.IsNil)
	pageRangeResp := pageListResp.PageList.PageRange
	c.Assert(pageRangeResp, chk.NotNil)
	c.Assert(*pageRangeResp, chk.HasLen, 1)
	rawStart, rawEnd := (*pageRangeResp)[0].Raw()
	c.Assert(rawStart, chk.Equals, offset1)
	c.Assert(rawEnd, chk.Equals, end1)

	clearPagesOptions := ClearPagesOptions{
		CpkInfo: &testCPKByValue,
	}
	clearResp, err := pbClient.ClearPages(context.Background(), HttpRange{2048, 2048}, &clearPagesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(clearResp.RawResponse.StatusCode, chk.Equals, 201)

	pageListResp, err = pbClient.GetPageRangesDiff(context.Background(), HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(pageListResp.PageList.PageRange, chk.IsNil)
}

func (s *aztestsSuite) TestBlobResizeWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	pbClient, _ := createNewPageBlobWithCPK(c, containerClient, PageBlobPageBytes*10, &testCPKByValue, nil)

	resizePageBlobOptions := ResizePageBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	_, err := pbClient.Resize(ctx, PageBlobPageBytes, &resizePageBlobOptions)
	c.Assert(err, chk.IsNil)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, _ := pbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	c.Assert(*resp.ContentLength, chk.Equals, int64(PageBlobPageBytes))
}

func createNewBlockBlobWithCPK(c *chk.C, container ContainerClient, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (bbClient BlockBlobClient, name string) {
	bbClient, name = getBlockBlobClient(c, container)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	}
	cResp, err := bbClient.Upload(ctx, strings.NewReader(blockBlobDefaultData), &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*cResp.IsServerEncrypted, chk.Equals, true)
	if cpkInfo != nil {
		c.Assert(cResp.EncryptionKeySHA256, chk.DeepEquals, cpkInfo.EncryptionKeySHA256)
	}
	if cpkScopeInfo != nil {
		c.Assert(cResp.EncryptionScope, chk.DeepEquals, cpkScopeInfo.EncryptionScope)
	}
	return
}

func (s *aztestsSuite) TestGetSetBlobMetadataWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, &testCPKByValue, nil)

	// Set blob metadata without encryption key should fail the request.
	_, err := bbClient.SetMetadata(ctx, basicMetadata, nil)
	c.Assert(err, chk.NotNil)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.EncryptionKeySHA256, chk.DeepEquals, testCPKByValue.EncryptionKeySHA256)

	// Get blob properties without encryption key should fail the request.
	getResp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.NotNil)

	getBlobPropertiesOptions := GetBlobPropertiesOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(getResp.Metadata, chk.NotNil)
	c.Assert(getResp.Metadata, chk.HasLen, len(basicMetadata))
	c.Assert(getResp.Metadata, chk.DeepEquals, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	c.Assert(err, chk.IsNil)

	getResp, err = bbClient.GetProperties(ctx, &getBlobPropertiesOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(getResp.Metadata, chk.IsNil)
}

func (s *aztestsSuite) TestGetSetBlobMetadataWithCPKScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, nil, &testCPKByScope)

	// Set blob metadata without encryption key should fail the request.
	_, err := bbClient.SetMetadata(ctx, basicMetadata, nil)
	c.Assert(err, chk.NotNil)

	setBlobMetadataOptions := SetBlobMetadataOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err := bbClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.EncryptionScope, chk.DeepEquals, testCPKByScope.EncryptionScope)

	getResp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(getResp.Metadata, chk.NotNil)
	c.Assert(getResp.Metadata, chk.HasLen, len(basicMetadata))
	c.Assert(getResp.Metadata, chk.DeepEquals, basicMetadata)

	_, err = bbClient.SetMetadata(ctx, map[string]string{}, &setBlobMetadataOptions)
	c.Assert(err, chk.IsNil)

	getResp, err = bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(getResp.Metadata, chk.IsNil)
}

func (s *aztestsSuite) TestBlobSnapshotWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, &testCPKByValue, nil)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	resp, err := bbClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.NotNil)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkInfo: &testInvalidCPKByValue,
	}
	resp, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	c.Assert(err, chk.NotNil)

	createBlobSnapshotOptions1 := CreateBlobSnapshotOptions{
		CpkInfo: &testCPKByValue,
	}
	resp, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.IsServerEncrypted, chk.Equals, false)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*dResp.EncryptionKeySHA256, chk.DeepEquals, *testCPKByValue.EncryptionKeySHA256)

	_, err = snapshotURL.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)

	// Get blob properties of snapshot without encryption key should fail the request.
	_, err = snapshotURL.GetProperties(ctx, nil)
	c.Assert(err, chk.NotNil)

	//c.Assert(err.(*StorageError).Response().StatusCode, chk.Equals, 404)
}

func (s *aztestsSuite) TestBlobSnapshotWithCPKScope(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, nil, &testCPKByScope)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	resp, err := bbClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.NotNil)

	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
		CpkScopeInfo: &testInvalidCPKByScope,
	}
	resp, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
	c.Assert(err, chk.NotNil)

	createBlobSnapshotOptions1 := CreateBlobSnapshotOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	resp, err = bbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions1)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.IsServerEncrypted, chk.Equals, false)

	snapshotURL := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := DownloadBlobOptions{
		CpkScopeInfo: &testCPKByScope,
	}
	dResp, err := snapshotURL.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(*dResp.EncryptionScope, chk.DeepEquals, *testCPKByScope.EncryptionScope)

	_, err = snapshotURL.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)
}

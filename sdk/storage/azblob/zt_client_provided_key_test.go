package azblob

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/ioutil"
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
	EncryptionKeySha256: &testEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

var testEncryptedScope = "blobgokeytestscope"
var testCPKByName = CpkScopeInfo{
	EncryptionScope: &testEncryptedScope,
}

var testInvalidEncryptedScope = "mumbojumbo"
var testInvalidCPKByName = CpkScopeInfo{
	EncryptionScope: &testInvalidEncryptedScope,
}

func blockIDBinaryToBase64(blockID []byte) string {
	return base64.StdEncoding.EncodeToString(blockID)
}

func blockIDBase64ToBinary(blockID string) []byte {
	binaryStr, _ := base64.StdEncoding.DecodeString(blockID)
	return binaryStr
}

func blockIDBase64ToInt(blockID string) int {
	blockIDBase64ToBinary(blockID)
	return int(binary.LittleEndian.Uint32(blockIDBase64ToBinary(blockID)))
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
	c.Assert(*resp.EncryptionKeySHA256, chk.DeepEquals, *(testCPKByValue.EncryptionKeySha256))

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
	c.Assert(getResp.ETag(), chk.DeepEquals, resp.ETag)
	c.Assert(getResp.LastModified(), chk.DeepEquals, resp.LastModified)
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
			CpkScopeInfo: &testCPKByName,
		}
		_, err := bbClient.StageBlock(ctx, base64BlockIDs[index], strings.NewReader(word), &stageBlockOptions)
		c.Assert(err, chk.IsNil)
	}

	commitBlockListOptions := CommitBlockListOptions{
		CpkScope: &testCPKByName,
	}
	resp, err := bbClient.CommitBlockList(ctx, base64BlockIDs, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.LastModified, chk.NotNil)
	c.Assert(*resp.IsServerEncrypted, chk.Equals, true)
	c.Assert(resp.EncryptionScope, chk.DeepEquals, testCPKByName.EncryptionScope)

	downloadBlobOptions := DownloadBlobOptions{
		CpkInfo: &testCPKByValue,
	}
	getResp, err := bbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.NotNil)

	// TODO: Fix issue with storage error interface
	//storageErr := err.(StorageError)
	//c.Assert(storageErr.RawResponse.StatusCode, chk.Equals, 409)
	//c.Assert(storageErr.ServiceCode(), chk.Equals, ServiceCodeFeatureEncryptionMismatch)

	downloadBlobOptions = DownloadBlobOptions{
		CpkScopeInfo: &testCPKByName,
	}
	getResp, err = bbClient.Download(ctx, &downloadBlobOptions)
	c.Assert(err, chk.IsNil)
	b := bytes.Buffer{}
	reader := getResp.Body(RetryReaderOptions{})
	_, err = b.ReadFrom(reader)
	c.Assert(err, chk.IsNil)
	_ = reader.Close() // The client must close the response body when finished with it
	c.Assert(b.String(), chk.Equals, "AAA BBB CCC ")
	c.Assert(getResp.ETag(), chk.DeepEquals, resp.ETag)
	c.Assert(getResp.LastModified(), chk.DeepEquals, resp.LastModified)
	c.Assert(*getResp.IsServerEncrypted(), chk.Equals, true)
	c.Assert(getResp.r.EncryptionScope, chk.DeepEquals, testCPKByName.EncryptionScope)
}

func (s *aztestsSuite) TestPutBlockFromURLAndCommitWithCPK(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := 8 * 1024 // 8 KB
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
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
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

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, content)

	// Get properties to validate the tier
	//destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	//c.Assert(*destBlobPropResp.AccessTier, chk.Equals, string(tier))
}

//func (s *aztestsSuite) TestPutBlockFromURLAndCommitWithCPKWithScope(c *chk.C) {
//	bsu := getBSU()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	contentSize := 2 * 1024 // 2KB
//	r, srcData := getRandomDataAndReader(contentSize)
//	ctx := context.Background()
//	bbClient := containerClient.NewBlockBlobURL(generateBlobName())
//
//	uploadSrcResp, err := bbClient.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)
//
//	srcBlobParts := NewBlobURLParts(bbClient.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//	destBlob := containerClient.NewBlockBlobURL(generateBlobName())
//	blockID1, blockID2 := blockIDIntToBase64(0), blockIDIntToBase64(1)
//	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, 1*1024, LeaseAccessConditions{}, ModifiedAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(stageResp1.Response().StatusCode, chk.Equals, 201)
//	c.Assert(stageResp1.ContentMD5(), chk.Not(chk.Equals), "")
//	c.Assert(stageResp1.RequestID(), chk.Not(chk.Equals), "")
//	c.Assert(stageResp1.Version(), chk.Not(chk.Equals), "")
//	c.Assert(stageResp1.Date().IsZero(), chk.Equals, false)
//	c.Assert(stageResp1.IsServerEncrypted(), chk.Equals, "true")
//
//	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 1*1024, CountToEnd, LeaseAccessConditions{}, ModifiedAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(stageResp2.Response().StatusCode, chk.Equals, 201)
//	c.Assert(stageResp2.ContentMD5(), chk.Not(chk.Equals), "")
//	c.Assert(stageResp2.RequestID(), chk.Not(chk.Equals), "")
//	c.Assert(stageResp2.Version(), chk.Not(chk.Equals), "")
//	c.Assert(stageResp2.Date().IsZero(), chk.Equals, false)
//	c.Assert(stageResp2.IsServerEncrypted(), chk.Equals, "true")
//
//	blockList, err := destBlob.GetBlockList(ctx, BlockListAll, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(blockList.Response().StatusCode, chk.Equals, 200)
//	c.Assert(blockList.UncommittedBlocks, chk.HasLen, 2)
//	c.Assert(blockList.CommittedBlocks, chk.HasLen, 0)
//
//	listResp, err := destBlob.CommitBlockList(ctx, []string{blockID1, blockID2}, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, DefaultAccessTier, nil, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(listResp.Response().StatusCode, chk.Equals, 201)
//	c.Assert(listResp.IsServerEncrypted(), chk.Equals, "true")
//	c.Assert(listResp.EncryptionScope(), chk.Equals, *(testCPK1.EncryptionScope))
//
//	blockList, err = destBlob.GetBlockList(ctx, BlockListAll, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(blockList.Response().StatusCode, chk.Equals, 200)
//	c.Assert(blockList.UncommittedBlocks, chk.HasLen, 0)
//	c.Assert(blockList.CommittedBlocks, chk.HasLen, 2)
//
//	// Download blob to do data integrity check.
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPK1)
//	c.Assert(err, chk.IsNil)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, srcData)
//	c.Assert(downloadResp.r.rawResponse.Header.Get("x-ms-encryption-scope"), chk.Equals, *(testCPK1.EncryptionScope))
//}
//
//func (s *aztestsSuite) TestUploadBlobWithMD5WithCPK(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	contentSize := 1 * 1024 * 1024
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Val := md5.Sum(srcData)
//	bbClient := containerClient.NewBlockBlobURL(generateBlobName())
//
//	uploadSrcResp, err := bbClient.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)
//
//	// Get blob content without encryption key should fail the request.
//	downloadResp, err := bbClient.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	// Download blob to do data integrity check.
//	downloadResp, err = bbClient.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(downloadResp.ContentMD5(), chk.DeepEquals, md5Val[:])
//	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(data, chk.DeepEquals, srcData)
//}
//
//func (s *aztestsSuite) TestAppendBlockWithCPK(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	appendBlobURL := containerClient.NewAppendBlobURL(generateBlobName())
//
//	resp, err := appendBlobURL.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, nil, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.StatusCode(), chk.Equals, 201)
//
//	words := []string{"AAA ", "BBB ", "CCC "}
//	for index, word := range words {
//		resp, err := appendBlobURL.AppendBlock(context.Background(), strings.NewReader(word), AppendBlobAccessConditions{}, nil, testCPKByValue)
//		c.Assert(err, chk.IsNil)
//		c.Assert(err, chk.IsNil)
//		c.Assert(resp.Response().StatusCode, chk.Equals, 201)
//		c.Assert(resp.BlobAppendOffset(), chk.Equals, strconv.Itoa(index*4))
//		c.Assert(resp.BlobCommittedBlockCount(), chk.Equals, int32(index+1))
//		c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
//		c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
//		c.Assert(resp.ContentMD5(), chk.Not(chk.Equals), "")
//		c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
//		c.Assert(resp.Version(), chk.Not(chk.Equals), "")
//		c.Assert(resp.Date().IsZero(), chk.Equals, false)
//		c.Assert(resp.IsServerEncrypted(), chk.Equals, "true")
//		c.Assert(resp.EncryptionKeySha256(), chk.Equals, *(testCPKByValue.EncryptionKeySha256))
//	}
//
//	// Get blob content without encryption key should fail the request.
//	_, err = appendBlobURL.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	// Download blob to do data integrity check.
//	downloadResp, err := appendBlobURL.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(string(data), chk.DeepEquals, "AAA BBB CCC ")
//}
//
//func (s *aztestsSuite) TestAppendBlockWithCPKByScope(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	appendBlobURL := containerClient.NewAppendBlobURL(generateBlobName())
//
//	resp, err := appendBlobURL.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, nil, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.StatusCode(), chk.Equals, 201)
//
//	words := []string{"AAA ", "BBB ", "CCC "}
//	for index, word := range words {
//		resp, err := appendBlobURL.AppendBlock(context.Background(), strings.NewReader(word), AppendBlobAccessConditions{}, nil, testCPK1)
//		c.Assert(err, chk.IsNil)
//		c.Assert(err, chk.IsNil)
//		c.Assert(resp.Response().StatusCode, chk.Equals, 201)
//		c.Assert(resp.BlobAppendOffset(), chk.Equals, strconv.Itoa(index*4))
//		c.Assert(resp.BlobCommittedBlockCount(), chk.Equals, int32(index+1))
//		c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
//		c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
//		c.Assert(resp.ContentMD5(), chk.Not(chk.Equals), "")
//		c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
//		c.Assert(resp.Version(), chk.Not(chk.Equals), "")
//		c.Assert(resp.Date().IsZero(), chk.Equals, false)
//		c.Assert(resp.IsServerEncrypted(), chk.Equals, "true")
//		c.Assert(resp.EncryptionScope(), chk.Equals, *(testCPK1.EncryptionScope))
//	}
//
//	// Download blob to do data integrity check.
//	downloadResp, err := appendBlobURL.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(downloadResp.IsServerEncrypted(), chk.Equals, "true")
//
//	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{ClientProvidedKeyOptions: testCPK1}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(string(data), chk.DeepEquals, "AAA BBB CCC ")
//	c.Assert(downloadResp.r.rawResponse.Header.Get("x-ms-encryption-scope"), chk.Equals, *(testCPK1.EncryptionScope))
//}
//
//func (s *aztestsSuite) TestAppendBlockFromURLWithCPK(c *chk.C) {
//	bsu := getBSU()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	contentSize := 2 * 1024 * 1024 // 2MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	ctx := context.Background() // Use default Background context
//	bbClient := containerClient.NewAppendBlobURL(generateName("src"))
//	destBlob := containerClient.NewAppendBlobURL(generateName("dest"))
//
//	cResp1, err := bbClient.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, nil, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(cResp1.StatusCode(), chk.Equals, 201)
//
//	resp, err := bbClient.AppendBlock(context.Background(), r, AppendBlobAccessConditions{}, nil, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
//	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
//	c.Assert(resp.ContentMD5(), chk.Not(chk.Equals), "")
//
//	srcBlobParts := NewBlobURLParts(bbClient.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	cResp2, err := destBlob.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, nil, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(cResp2.StatusCode(), chk.Equals, 201)
//
//	appendResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, 0, int64(contentSize), AppendBlobAccessConditions{}, ModifiedAccessConditions{}, nil, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(appendResp.ETag(), chk.Not(chk.Equals), ETagNone)
//	c.Assert(appendResp.LastModified().IsZero(), chk.Equals, false)
//	c.Assert(appendResp.IsServerEncrypted(), chk.Equals, "true")
//
//	// Get blob content without encryption key should fail the request.
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	// Download blob to do data integrity check.
//	downloadResp, err = destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{ClientProvidedKeyOptions: testCPKByValue}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, srcData)
//}
//
//func (s *aztestsSuite) TestPageBlockWithCPK(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	contentSize := 1 * 1024 * 1024
//	r, srcData := getRandomDataAndReader(contentSize)
//	bbClient, _ := createNewPageBlobWithCPK(c, containerClient, int64(contentSize), testCPKByValue)
//
//	uploadResp, err := bbClient.UploadPages(ctx, 0, r, PageBlobAccessConditions{}, nil, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadResp.Response().StatusCode, chk.Equals, 201)
//
//	// Get blob content without encryption key should fail the request.
//	downloadResp, err := bbClient.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	// Download blob to do data integrity check.
//	downloadResp, err = bbClient.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{ClientProvidedKeyOptions: testCPKByValue}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, srcData)
//}
//
//func (s *aztestsSuite) TestPageBlockWithCPKByScope(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	// defer delContainer(c, containerClient)
//
//	contentSize := 1 * 1024 * 1024
//	r, srcData := getRandomDataAndReader(contentSize)
//	bbClient, _ := createNewPageBlobWithCPK(c, containerClient, int64(contentSize), testCPK1)
//
//	uploadResp, err := bbClient.UploadPages(ctx, 0, r, PageBlobAccessConditions{}, nil, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadResp.Response().StatusCode, chk.Equals, 201)
//	c.Assert(uploadResp.EncryptionScope(), chk.Equals, *(testCPK1.EncryptionScope))
//
//	// Download blob to do data integrity check.
//	downloadResp, err := bbClient.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPK1)
//	c.Assert(err, chk.IsNil)
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{ClientProvidedKeyOptions: testCPK1}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, srcData)
//	c.Assert(downloadResp.r.rawResponse.Header.Get("x-ms-encryption-scope"), chk.Equals, *(testCPK1.EncryptionScope))
//}
//
//func (s *aztestsSuite) TestPageBlockFromURLWithCPK(c *chk.C) {
//	bsu := getBSU()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	contentSize := 1 * 1024 * 1024 // 1MB
//	r, srcData := getRandomDataAndReader(contentSize)
//	ctx := context.Background() // Use default Background context
//	bbClient, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
//	destBlob, _ := createNewPageBlobWithCPK(c, containerClient, int64(contentSize), testCPKByValue)
//
//	uploadResp, err := bbClient.UploadPages(ctx, 0, r, PageBlobAccessConditions{}, nil, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadResp.Response().StatusCode, chk.Equals, 201)
//	srcBlobParts := NewBlobURLParts(bbClient.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	resp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil, PageBlobAccessConditions{}, ModifiedAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.ETag(), chk.NotNil)
//	c.Assert(resp.LastModified(), chk.NotNil)
//	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
//	c.Assert(resp.IsServerEncrypted(), chk.Equals, "true")
//
//	// Download blob to do data integrity check.
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(downloadResp.r.EncryptionKeySha256(), chk.Equals, *(testCPKByValue.EncryptionKeySha256))
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{ClientProvidedKeyOptions: testCPKByValue}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, srcData)
//}
//
//func (s *aztestsSuite) TestUploadPagesFromURLWithMD5WithCPK(c *chk.C) {
//	bsu := getBSU()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	containerClient, _ := createNewContainer(c, bsu)
//	defer delContainer(c, containerClient)
//
//	contentSize := 1 * 1024 * 1024
//	r, srcData := getRandomDataAndReader(contentSize)
//	md5Value := md5.Sum(srcData)
//	srcBlob, _ := createNewPageBlobWithSize(c, containerClient, int64(contentSize))
//
//	uploadSrcResp1, err := srcBlob.UploadPages(ctx, 0, r, PageBlobAccessConditions{}, nil, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadSrcResp1.Response().StatusCode, chk.Equals, 201)
//
//	srcBlobParts := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//	destBlob, _ := createNewPageBlobWithCPK(c, containerClient, int64(contentSize), testCPKByValue)
//	uploadResp, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), md5Value[:], PageBlobAccessConditions{}, ModifiedAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(uploadResp.ETag(), chk.NotNil)
//	c.Assert(uploadResp.LastModified(), chk.NotNil)
//	c.Assert(uploadResp.EncryptionKeySha256(), chk.Equals, *(testCPKByValue.EncryptionKeySha256))
//	c.Assert(uploadResp.ContentMD5(), chk.DeepEquals, md5Value[:])
//	c.Assert(uploadResp.BlobSequenceNumber(), chk.Equals, int64(0))
//
//	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(downloadResp.r.EncryptionKeySha256(), chk.Equals, *(testCPKByValue.EncryptionKeySha256))
//	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{ClientProvidedKeyOptions: testCPKByValue}))
//	c.Assert(err, chk.IsNil)
//	c.Assert(destData, chk.DeepEquals, srcData)
//
//	_, badMD5 := getRandomDataAndReader(16)
//	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), badMD5[:], PageBlobAccessConditions{}, ModifiedAccessConditions{}, ClientProvidedKeyOptions{})
//	validateStorageError(c, err, ServiceCodeMd5Mismatch)
//}
//
//func (s *aztestsSuite) TestGetSetBlobMetadataWithCPK(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, testCPKByValue)
//
//	metadata := Metadata{"key": "value", "another_key": "1234"}
//
//	// Set blob metadata without encryption key should fail the request.
//	_, err := bbClient.SetMetadata(ctx, metadata, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	resp, err := bbClient.SetMetadata(ctx, metadata, BlobAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.EncryptionKeySha256(), chk.Equals, *(testCPKByValue.EncryptionKeySha256))
//
//	// Get blob properties without encryption key should fail the request.
//	getResp, err := bbClient.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	getResp, err = bbClient.GetProperties(ctx, BlobAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(getResp.NewMetadata(), chk.HasLen, 2)
//	c.Assert(getResp.NewMetadata(), chk.DeepEquals, metadata)
//
//	_, err = bbClient.SetMetadata(ctx, Metadata{}, BlobAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//
//	getResp, err = bbClient.GetProperties(ctx, BlobAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(getResp.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *aztestsSuite) TestGetSetBlobMetadataWithCPKByScope(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, testCPK1)
//
//	metadata := Metadata{"key": "value", "another_key": "1234"}
//
//	// Set blob metadata without encryption key should fail the request.
//	_, err := bbClient.SetMetadata(ctx, metadata, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	_, err = bbClient.SetMetadata(ctx, metadata, BlobAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//
//	getResp, err := bbClient.GetProperties(ctx, BlobAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(getResp.NewMetadata(), chk.HasLen, 2)
//	c.Assert(getResp.NewMetadata(), chk.DeepEquals, metadata)
//
//	_, err = bbClient.SetMetadata(ctx, Metadata{}, BlobAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//
//	getResp, err = bbClient.GetProperties(ctx, BlobAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(getResp.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotWithCPK(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, testCPKByValue)
//	_, err := bbClient.Upload(ctx, strings.NewReader("113333555555"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, testCPKByValue)
//
//	// Create Snapshot of an encrypted blob without encryption key should fail the request.
//	resp, err := bbClient.CreateSnapshot(ctx, nil, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	resp, err = bbClient.CreateSnapshot(ctx, nil, BlobAccessConditions{}, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.IsServerEncrypted(), chk.Equals, "false")
//	snapshotURL := bbClient.WithSnapshot(resp.Snapshot())
//
//	dResp, err := snapshotURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPKByValue)
//	c.Assert(err, chk.IsNil)
//	c.Assert(dResp.r.EncryptionKeySha256(), chk.Equals, *(testCPKByValue.EncryptionKeySha256))
//	_, err = snapshotURL.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//	c.Assert(err, chk.IsNil)
//
//	// Get blob properties of snapshot without encryption key should fail the request.
//	_, err = snapshotURL.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//	c.Assert(err.(StorageError).Response().StatusCode, chk.Equals, 404)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotWithCPKByScope(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	bbClient, _ := createNewBlockBlobWithCPK(c, containerClient, testCPKByValue)
//	_, err := bbClient.Upload(ctx, strings.NewReader("113333555555"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, testCPK1)
//
//	// Create Snapshot of an encrypted blob without encryption key should fail the request.
//	resp, err := bbClient.CreateSnapshot(ctx, nil, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//
//	resp, err = bbClient.CreateSnapshot(ctx, nil, BlobAccessConditions{}, testCPK1)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.IsServerEncrypted(), chk.Equals, "false")
//	snapshotURL := bbClient.WithSnapshot(resp.Snapshot())
//
//	_, err = snapshotURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, testCPK1)
//	c.Assert(err, chk.IsNil)
//	_, err = snapshotURL.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//	c.Assert(err, chk.IsNil)
//
//	// Get blob properties of snapshot without encryption key should fail the request.
//	_, err = snapshotURL.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
//	c.Assert(err, chk.NotNil)
//	c.Assert(err.(StorageError).Response().StatusCode, chk.Equals, 404)
//}

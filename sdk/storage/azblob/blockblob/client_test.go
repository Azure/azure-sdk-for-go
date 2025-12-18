//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var proposedLeaseIDs = []*string{to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"), to.Ptr("326cc5e1-746e-4af8-4811-a50e6629a8ca")}

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running blockblob Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &BlockBlobRecordedTestsSuite{})
		suite.Run(t, &BlockBlobUnrecordedTestsSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &BlockBlobRecordedTestsSuite{})
	case recording.RecordingMode:
		suite.Run(t, &BlockBlobRecordedTestsSuite{})
	}
}

func (s *BlockBlobRecordedTestsSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *BlockBlobRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *BlockBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *BlockBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *BlockBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *BlockBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type BlockBlobRecordedTestsSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type BlockBlobUnrecordedTestsSuite struct {
	suite.Suite
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlockBlobClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testName
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	cred, err := credential.New(nil)
	_require.NoError(err)

	bbClient, err := blockblob.NewClient(blobURL, cred, nil)
	_require.NoError(err)

	contentSize := 4 * 1024 // 4 KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)

	// Prepare bbClient for copy.
	resp, err := bbClient.Upload(context.Background(), rsc, nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlockBlobClientSharedKey() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	blobName := testName
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	cred, err := blob.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	bbClient, err := blockblob.NewClientWithSharedKeyCredential(blobURL, cred, nil)
	_require.NoError(err)

	contentSize := 4 * 1024 // 4 KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)

	// Prepare bbClient for copy.
	resp, err := bbClient.Upload(context.Background(), rsc, nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlockBlobClientConnectionString() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testName
	connectionString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)

	bbClient, err := blockblob.NewClientFromConnectionString(*connectionString, containerName, blobName, nil)
	_require.NoError(err)

	contentSize := 4 * 1024 // 4 KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)

	// Prepare bbClient for copy.
	resp, err := bbClient.Upload(context.Background(), rsc, nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

//	func (s *BlockBlobRecordedTestsSuite) TestStageGetBlocks() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//
// //		svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		blobName := testcommon.GenerateBlobName(testName)
//		bbClient := containerClient.NewBlockBlobClient(blobName)
//
//		data := []string{"Azure ", "Storage ", "Block ", "Blob."}
//		base64BlockIDs := make([]string, len(data))
//
//		for index, d := range data {
//			base64BlockIDs[index] = blockIDIntToBase64(index)
//			io.NopCloser(strings.NewReader("hello world"))
//			putResp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(d)), nil)
//			_require.NoError(err)
//			//_require.Equal(putResp.RawResponse.StatusCode, 201)
//			_require.Nil(putResp.ContentMD5)
//			_require.NotNil(putResp.RequestID)
//			_require.NotNil(putResp.Version)
//			_require.NotNil(putResp.Date)
//			_require.Equal((*putResp.Date).IsZero(), false)
//		}
//
//		blockList, err := bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
//		_require.NoError(err)
//		// _require.Equal(blockList.RawResponse.StatusCode, 200)
//		_require.Nil(blockList.LastModified)
//		_require.Nil(blockList.ETag)
//		_require.NotNil(blockList.ContentType)
//		_require.Nil(blockList.BlobContentLength)
//		_require.NotNil(blockList.RequestID)
//		_require.NotNil(blockList.Version)
//		_require.NotNil(blockList.Date)
//		_require.Equal((*blockList.Date).IsZero(), false)
//		_require.NotNil(blockList.BlockList)
//		_require.Nil(blockList.BlockList.CommittedBlocks)
//		_require.NotNil(blockList.BlockList.UncommittedBlocks)
//		_require.Len(blockList.BlockList.UncommittedBlocks, len(data))
//
//		listResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
//		_require.NoError(err)
//		// _require.Equal(listResp.RawResponse.StatusCode,  201)
//		_require.NotNil(listResp.LastModified)
//		_require.Equal((*listResp.LastModified).IsZero(), false)
//		_require.NotNil(listResp.ETag)
//		_require.NotNil(listResp.RequestID)
//		_require.NotNil(listResp.Version)
//		_require.NotNil(listResp.Date)
//		_require.Equal((*listResp.Date).IsZero(), false)
//
//		blockList, err = bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
//		_require.NoError(err)
//		// _require.Equal(blockList.RawResponse.StatusCode, 200)
//		_require.NotNil(blockList.LastModified)
//		_require.Equal((*blockList.LastModified).IsZero(), false)
//		_require.NotNil(blockList.ETag)
//		_require.NotNil(blockList.ContentType)
//		_require.Equal(*blockList.BlobContentLength, int64(25))
//		_require.NotNil(blockList.RequestID)
//		_require.NotNil(blockList.Version)
//		_require.NotNil(blockList.Date)
//		_require.Equal((*blockList.Date).IsZero(), false)
//		_require.NotNil(blockList.BlockList)
//		_require.NotNil(blockList.BlockList.CommittedBlocks)
//		_require.Nil(blockList.BlockList.UncommittedBlocks)
//		_require.Len(blockList.BlockList.CommittedBlocks, len(data))
//	}
//
//	func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockFromURL() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		contentSize := 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		body := bytes.NewReader(content)
//		rsc := streaming.NopCloser(body)
//
//		ctx := context.Background() // Use default Background context
//		srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
//
//		destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))
//
//		// Prepare source bbClient for copy.
//		_, err = srcBlob.Upload(context.Background(), rsc, nil)
//		_require.NoError(err)
//		//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get source blob url with SAS for StageFromURL.
//		srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//		credential, err := testcommon.GetGenericSharedKeyCredential(nil, testcommon.TestAccountDefault)
//		_require.NoError(err)
//
//		srcBlobParts.SAS, err = BlobSASSignatureValues{
//			Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//			ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//			ContainerName: srcBlobParts.ContainerName,
//			BlobName:      srcBlobParts.BlobName,
//			Permissions:   BlobSASPermissions{Read: true}.String(),
//		}.Sign(credential)
//		_require.NoError(err)
//
//		srcBlobURLWithSAS := srcBlobParts.URL()
//
//		// Stage blocks from URL.
//		blockIDs := testcommon.GenerateBlockIDsList(2)
//
//		stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[0], srcBlobURLWithSAS, 0, &BlockBlobStageBlockFromURLOptions{
//			Offset: to.Ptr[int64](0),
//			Count:  to.Ptr(int64(contentSize / 2)),
//		})
//		_require.NoError(err)
//		// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//		_require.NotEqual(stageResp1.ContentMD5, "")
//		_require.NotEqual(stageResp1.RequestID, "")
//		_require.NotEqual(stageResp1.Version, "")
//		_require.Equal(stageResp1.Date.IsZero(), false)
//
//		stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[1], srcBlobURLWithSAS, 0, &BlockBlobStageBlockFromURLOptions{
//			Offset: to.Ptr(int64(contentSize / 2)),
//			Count:  to.Ptr(int64(CountToEnd)),
//		})
//		_require.NoError(err)
//		// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//		_require.NotEqual(stageResp2.ContentMD5, "")
//		_require.NotEqual(stageResp2.RequestID, "")
//		_require.NotEqual(stageResp2.Version, "")
//		_require.Equal(stageResp2.Date.IsZero(), false)
//
//		// Check block list.
//		blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//		_require.NoError(err)
//		// _require.Equal(blockList.RawResponse.StatusCode, 200)
//		_require.NotNil(blockList.BlockList)
//		_require.Nil(blockList.BlockList.CommittedBlocks)
//		_require.NotNil(blockList.BlockList.UncommittedBlocks)
//		_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//		// Commit block list.
//		listResp, err := destBlob.CommitBlockList(context.Background(), blockIDs, nil)
//		_require.NoError(err)
//		// _require.Equal(listResp.RawResponse.StatusCode,  201)
//		_require.NotNil(listResp.LastModified)
//		_require.Equal((*listResp.LastModified).IsZero(), false)
//		_require.NotNil(listResp.ETag)
//		_require.NotNil(listResp.RequestID)
//		_require.NotNil(listResp.Version)
//		_require.NotNil(listResp.Date)
//		_require.Equal((*listResp.Date).IsZero(), false)
//
//		// Check data integrity through downloading.
//		downloadResp, err := destBlob.BlobClient.DownloadStream(context.Background(), nil)
//		_require.NoError(err)
//		destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//		_require.NoError(err)
//		_require.EqualValues(destData, content)
//	}

func (s *BlockBlobRecordedTestsSuite) TestPutBlobCrcResponseHeader() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testName

	bbClient := containerClient.NewBlockBlobClient(blobName)

	contentSize := 4 * 1024 // 4 KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	resp, err := bbClient.Upload(context.Background(), rsc, nil)
	_require.NoError(err)
	_require.NotNil(resp)
	_require.NotNil(resp.ContentCRC64)
	_require.Equal(resp.ContentCRC64, crc)
}

func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockFromURLWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4 KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	contentMD5 := md5.Sum(sourceData)
	rsc := streaming.NopCloser(r)

	srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
	destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))

	// Prepare source bbClient for copy.
	_, err = srcBlob.Upload(context.Background(), rsc, nil)
	_require.NoError(err)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)
	perms := sas.BlobPermissions{Read: true}

	srcBlobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   perms.String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobURLWithSAS := srcBlobParts.String()

	// Stage blocks from URL.
	blockIDs := testcommon.GenerateBlockIDsList(2)

	opts := blockblob.StageBlockFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeMD5(contentMD5[:]),
	}

	stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[0], srcBlobURLWithSAS, &opts)
	_require.NoError(err)
	_require.Equal(stageResp1.ContentMD5, contentMD5[:])

	stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[1], srcBlobURLWithSAS, &opts)
	_require.NoError(err)
	_require.Equal(stageResp2.ContentMD5, contentMD5[:])

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.CommittedBlocks)
	_require.NotNil(blockList.UncommittedBlocks)
	_require.Len(blockList.UncommittedBlocks, 2)

	// Commit block list.
	_, err = destBlob.CommitBlockList(context.Background(), blockIDs, nil)
	_require.NoError(err)

	// Check data integrity through downloading.
	destBuffer := make([]byte, 4*1024)
	downloadBufferOptions := blob.DownloadBufferOptions{Range: blob.HTTPRange{Offset: 0, Count: 4096}}
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, &downloadBufferOptions)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)

	// Test stage block from URL with bad MD5 value
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	var badMD5Validator blob.SourceContentValidationTypeMD5 = badMD5
	opts = blockblob.StageBlockFromURLOptions{
		SourceContentValidation: badMD5Validator,
	}
	_, err = destBlob.StageBlockFromURL(context.Background(), blockIDs[0], srcBlobURLWithSAS, &opts)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)

	_, err = destBlob.StageBlockFromURL(context.Background(), blockIDs[1], srcBlobURLWithSAS, &opts)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
}

func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockFromURLWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4 KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
	destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))

	// Prepare source bbClient for copy.
	_, err = srcBlob.Upload(context.Background(), rsc, nil)
	_require.NoError(err)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)
	perms := sas.BlobPermissions{Read: true}

	srcBlobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   perms.String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobURLWithSAS := srcBlobParts.String()

	// Stage blocks from URL.
	blockIDs := testcommon.GenerateBlockIDsList(2)

	opts := blockblob.StageBlockFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeCRC64(crc),
	}

	stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[0], srcBlobURLWithSAS, &opts)
	_require.NoError(err)
	_require.Equal(stageResp1.ContentCRC64, crc)

	stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[1], srcBlobURLWithSAS, &opts)
	_require.NoError(err)
	_require.Equal(stageResp2.ContentCRC64, crc)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.CommittedBlocks)
	_require.NotNil(blockList.UncommittedBlocks)
	_require.Len(blockList.UncommittedBlocks, 2)

	// Commit block list.
	_, err = destBlob.CommitBlockList(context.Background(), blockIDs, nil)
	_require.NoError(err)

	// Check data integrity through downloading.
	destBuffer := make([]byte, 4*1024)
	downloadBufferOptions := blob.DownloadBufferOptions{Range: blob.HTTPRange{Offset: 0, Count: 4096}}
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, &downloadBufferOptions)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)

	// Test stage block from URL with bad CRC64 value
	_, sourceData = testcommon.GetDataAndReader(testName, 16)
	crc64Value = crc64.Checksum(sourceData, shared.CRC64Table)
	badCRC := make([]byte, 8)
	binary.LittleEndian.PutUint64(badCRC, crc64Value)
	opts = blockblob.StageBlockFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeCRC64(badCRC),
	}
	_, err = destBlob.StageBlockFromURL(context.Background(), blockIDs[0], srcBlobURLWithSAS, &opts)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CRC64Mismatch)

	_, err = destBlob.StageBlockFromURL(context.Background(), blockIDs[1], srcBlobURLWithSAS, &opts)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CRC64Mismatch)
}

func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockFromURLWithRequestIntent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4 KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)

	srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
	destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))

	// Prepare source bbClient for copy.
	_, err = srcBlob.Upload(context.Background(), rsc, nil)
	_require.NoError(err)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	sharedKeyCredential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)
	perms := sas.BlobPermissions{Read: true}

	srcBlobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   perms.String(),
	}.SignWithSharedKey(sharedKeyCredential)
	_require.NoError(err)

	srcBlobURLWithSAS := srcBlobParts.String()

	// Stage blocks from URL.
	blockIDs := testcommon.GenerateBlockIDsList(2)
	requestIntent := blob.FileRequestIntentTypeBackup

	opts := blockblob.StageBlockFromURLOptions{
		FileRequestIntent: &requestIntent,
	}

	stageResponse, err := destBlob.StageBlockFromURL(context.Background(), blockIDs[0], srcBlobURLWithSAS, &opts)
	_require.NoError(err)
	_require.NotNil(stageResponse)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.NotNil(blockList.BlockList)
	_require.Nil(blockList.CommittedBlocks)
}

//
//	func (s *BlockBlobUnrecordedTestsSuite) TestCopyBlockBlobFromURL() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		const contentSize = 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		contentMD5 := md5.Sum(content)
//		body := bytes.NewReader(content)
//		ctx := context.Background()
//
//		srcBlob := containerClient.NewBlockBlobClient("srcblob")
//		destBlob := containerClient.NewBlockBlobClient("destblob")
//
//		// Prepare source bbClient for copy.
//		_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(body), nil)
//		_require.NoError(err)
//		// _require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get source blob url with SAS for StageFromURL.
//		srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//		credential, err := testcommon.GetGenericSharedKeyCredential(nil, testcommon.TestAccountDefault)
//		_require.NoError(err)
//
//		srcBlobParts.SAS, err = BlobSASSignatureValues{
//			Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//			ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//			ContainerName: srcBlobParts.ContainerName,
//			BlobName:      srcBlobParts.BlobName,
//			Permissions:   BlobSASPermissions{Read: true}.String(),
//		}.Sign(credential)
//		if err != nil {
//			s.T().Fatal(err)
//		}
//
//		srcBlobURLWithSAS := srcBlobParts.URL()
//
//		// Invoke copy bbClient from URL.
//		sourceContentMD5 := contentMD5[:]
//		resp, err := destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &BlockBlobCopyFromURLOptions{
//			Metadata:         map[string]string{"foo": "bar"},
//			SourceContentMD5: sourceContentMD5,
//		})
//		_require.NoError(err)
//		// _require.Equal(resp.RawResponse.StatusCode, 202)
//		_require.NotNil(resp.ETag)
//		_require.NotNil(resp.RequestID)
//		_require.NotNil(resp.Version)
//		_require.NotNil(resp.Date)
//		_require.Equal((*resp.Date).IsZero(), false)
//		_require.NotNil(resp.CopyID)
//		_require.EqualValues(resp.ContentMD5, sourceContentMD5)
//		_require.Equal(*resp.CopyStatus, "success")
//
//		// Make sure the metadata got copied over
//		getPropResp, err := destBlob.GetProperties(context.Background(), nil)
//		_require.NoError(err)
//		metadata := getPropResp.Metadata
//		_require.NotNil(metadata)
//		_require.Len(metadata, 1)
//		_require.EqualValues(metadata, map[string]string{"Foo": "bar"})
//
//		// Check data integrity through downloading.
//		downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
//		_require.NoError(err)
//		destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//		_require.NoError(err)
//		_require.EqualValues(destData, content)
//
//		// Edge case 1: Provide bad MD5 and make sure the copy fails
//		_, badMD5 := testcommon.GetRandomDataAndReader(16)
//		copyBlockBlobFromURLOptions1 := BlockBlobCopyFromURLOptions{
//			SourceContentMD5: badMD5,
//		}
//		resp, err = destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
//		_require.Error(err)
//
//		// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
//		copyBlockBlobFromURLOptions2 := BlockBlobCopyFromURLOptions{
//			SourceContentMD5: sourceContentMD5,
//		}
//		resp, err = destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
//		_require.NoError(err)
//		// _require.Equal(resp.RawResponse.StatusCode, 202)
//		_require.EqualValues(*resp.CopyStatus, "success")
//	}
//
//	func (s *BlockBlobUnrecordedTestsSuite) TestBlobSASQueryParamOverrideResponseHeaders() {
//		_require := require.New(s.T())
//		testName := s.T().Name()
//		svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//		if err != nil {
//			s.Fail("Unable to fetch service client because " + err.Error())
//		}
//
//		containerName := testcommon.GenerateContainerName(testName)
//		containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//		defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//		const contentSize = 8 * 1024 // 8 KB
//		content := make([]byte, contentSize)
//		body := bytes.NewReader(content)
//		// contentMD5 := md5.Sum(content)
//
//		ctx := context.Background()
//
//		bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))
//
//		_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), nil)
//		_require.NoError(err)
//		// _require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//		// Get blob url with SAS.
//		blobParts, _ := NewBlobURLParts(bbClient.URL())
//
//		cacheControlVal := "cache-control-override"
//		contentDispositionVal := "content-disposition-override"
//		contentEncodingVal := "content-encoding-override"
//		contentLanguageVal := "content-language-override"
//		contentTypeVal := "content-type-override"
//
//		credential, err := testcommon.GetGenericSharedKeyCredential(nil, testcommon.TestAccountDefault)
//		_require.NoError(err)
//		// Append User Delegation SAS token to URL
//		blobParts.SAS, err = BlobSASSignatureValues{
//			Protocol:           SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//			ExpiryTime:         time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//			ContainerName:      blobParts.ContainerName,
//			BlobName:           blobParts.BlobName,
//			Permissions:        BlobSASPermissions{Read: true}.String(),
//			CacheControl:       cacheControlVal,
//			ContentDisposition: contentDispositionVal,
//			ContentEncoding:    contentEncodingVal,
//			ContentLanguage:    contentLanguageVal,
//			ContentType:        contentTypeVal,
//		}.Sign(credential)
//		_require.NoError(err)
//
//		// Generate new bbClient client
//		blobURLWithSAS := blobParts.URL()
//		_require.NotNil(blobURLWithSAS)
//
//		blobClientWithSAS, err := NewBlockBlobClientWithNoCredential(blobURLWithSAS, nil)
//		_require.NoError(err)
//
//		gResp, err := blobClientWithSAS.GetProperties(context.Background(), nil)
//		_require.NoError(err)
//		_require.Equal(*gResp.CacheControl, cacheControlVal)
//		_require.Equal(*gResp.ContentDisposition, contentDispositionVal)
//		_require.Equal(*gResp.ContentEncoding, contentEncodingVal)
//		_require.Equal(*gResp.ContentLanguage, contentLanguageVal)
//		_require.Equal(*gResp.ContentType, contentTypeVal)
//	}

// nolint
func (s *BlockBlobRecordedTestsSuite) TestStageBlockWithGeneratedCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	// test stage block with valid CRC64 value
	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	contentCrc64 := crc64.Checksum(content, shared.CRC64Table)
	rsc := streaming.NopCloser(body)

	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	putResp, err := bbClient.StageBlock(context.Background(), blockID1, rsc, &blockblob.StageBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
	})
	_require.NoError(err)
	_require.NotNil(putResp.ContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(putResp.ContentCRC64), contentCrc64)
	_require.NotNil(putResp.RequestID)
	_require.NotNil(putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)
}

func (s *BlockBlobRecordedTestsSuite) TestStageBlockWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return
	}
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	// test stage block with valid CRC64 value
	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	contentCrc64 := crc64.Checksum(content, shared.CRC64Table)
	rsc := streaming.NopCloser(body)

	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	putResp, err := bbClient.StageBlock(context.Background(), blockID1, rsc, &blockblob.StageBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(contentCrc64),
	})
	_require.NoError(err)
	_require.EqualValues(binary.LittleEndian.Uint64(putResp.ContentCRC64), contentCrc64)

	// test put block with bad CRC64 value
	badContentCrc64 := binary.LittleEndian.Uint64(b[:])

	_, _ = rsc.Seek(0, io.SeekStart)
	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	_, err = bbClient.StageBlock(context.Background(), blockID2, rsc, &blockblob.StageBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(badContentCrc64),
	})
	_require.Error(err)
}

// nolint
func (s *BlockBlobRecordedTestsSuite) TestStageBlockWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := containerClient.NewBlockBlobClient(blobName)

	// test stage block with valid MD5 value
	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	putResp, err := bbClient.StageBlock(context.Background(), blockID1, rsc, &blockblob.StageBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(contentMD5),
	})
	_require.NoError(err)
	_require.EqualValues(putResp.ContentMD5, contentMD5)
	_require.NotNil(putResp.RequestID)
	_require.NotNil(putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	// test stage block with bad MD5 value
	_, badContent := testcommon.GetDataAndReader(testName, contentSize)
	badMD5Value := md5.Sum(badContent)
	badContentMD5 := badMD5Value[:]

	_, _ = rsc.Seek(0, io.SeekStart)
	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	_, err = bbClient.StageBlock(context.Background(), blockID2, rsc, &blockblob.StageBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(badContentMD5),
	})
	_require.Error(err)
	_require.Contains(err.Error(), bloberror.MD5Mismatch)
}

func (s *BlockBlobRecordedTestsSuite) TestPutBlobWithCRC64() {
	s.T().Skip("Content CRC64 cannot be validated in Upload()")
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := "test" + testcommon.GenerateContainerName(testName) + "1"
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4 KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	_, err = bbClient.StageBlock(context.Background(), blockID, streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(crc64Value),
	})
	_require.Error(err, bloberror.UnsupportedChecksum)
	// TODO: UploadResponse does not return ContentCRC64
	//	_require.Equal(resp.ContentCRC64, crc)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	})
	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	h := blob.ParseHTTPHeaders(resp)
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func (s *BlockBlobRecordedTestsSuite) TestUploadBlockWithImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.NoError(err)
	policy := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.NoError(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	legalHold := true
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		ImmutabilityPolicyExpiryTime: &currentTime,
		ImmutabilityPolicyMode:       &policy,
		LegalHold:                    &legalHold,
		HTTPHeaders:                  &testcommon.BasicHeaders,
	})

	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	policy1 := blob.ImmutabilityPolicyMode("unlocked")
	_require.Equal(resp.ImmutabilityPolicyMode, &policy1)

	_, err = bbClient.SetLegalHold(context.Background(), false, nil)
	_require.NoError(err)

	_, err = bbClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.NoError(err)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func setUpPutBlobFromURLTest(testName string, _require *require.Assertions, svcClient *service.Client) (*container.Client, *blockblob.Client, *blockblob.Client, string, []byte) {
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	srcBlob := testcommon.GenerateBlobName("src" + testName)
	srcBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, srcBlob, containerClient)

	dest := testcommon.GenerateBlobName("dest" + testName)
	destBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, dest, containerClient)

	// Upload some data to source
	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	_, err := srcBBClient.Upload(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)

	// Create SAS for source and get SAS URL
	expiryTime := time.Now().UTC().Add(15 * time.Minute)
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true, Tag: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobParts, _ := blob.ParseURL(srcBBClient.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	return containerClient, srcBBClient, destBBClient, srcBlobURLWithSAS, sourceData
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, sourceData := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Invoke UploadBlobFromURL
	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, nil)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Download data from destination
	destBuffer := make([]byte, 4*1024)
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, nil)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLWithCopySourceTagsDefault() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Set tags to source
	srcBlobTagsMap := map[string]string{
		"source": "tags",
	}
	_, err = srcBlob.SetTags(context.Background(), srcBlobTagsMap, nil)
	_require.NoError(err)

	// Dest tags
	destBlobTagsMap := map[string]string{
		"dest": "tags",
	}

	// By default, the CopySourceTag header is Replace
	options := blockblob.UploadBlobFromURLOptions{
		Tags: destBlobTagsMap,
	}

	// Invoke UploadBlobFromURL
	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Get tags from dest and check if tags got replaced with dest tags
	resp, err := destBlob.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.BlobTagSet[0].Key, "dest")
	_require.Equal(*resp.BlobTagSet[0].Value, "tags")
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLWithCopySourceTagsReplace() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Set tags to source
	srcBlobTagsMap := map[string]string{
		"source": "tags",
	}
	_, err = srcBlob.SetTags(context.Background(), srcBlobTagsMap, nil)
	_require.NoError(err)

	// Dest tags
	destBlobTagsMap := map[string]string{
		"dest": "tags",
	}

	options := blockblob.UploadBlobFromURLOptions{
		Tags:           destBlobTagsMap,
		CopySourceTags: to.Ptr(blockblob.BlobCopySourceTagsReplace),
	}

	// Invoke UploadBlobFromURL
	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Get tags from dest and check if tags got replaced with dest tags
	resp, err := destBlob.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.BlobTagSet[0].Key, "dest")
	_require.Equal(*resp.BlobTagSet[0].Value, "tags")
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLWithCopySourceTagsCopy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Set tags to source
	srcBlobTagsMap := map[string]string{
		"source": "tags",
	}
	_, err = srcBlob.SetTags(context.Background(), srcBlobTagsMap, nil)
	_require.NoError(err)

	// Set tags to dest to ensure that COPY works
	destBlobTagsMap := map[string]string{
		"dest": "tags",
	}
	_, err = destBlob.SetTags(context.Background(), destBlobTagsMap, nil)
	_require.NoError(err)

	options := blockblob.UploadBlobFromURLOptions{
		CopySourceTags: to.Ptr(blockblob.BlobCopySourceTagsCopy),
	}

	// Invoke UploadBlobFromURL
	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Get tags from dest and check if it matches source tags
	resp, err := destBlob.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.BlobTagSet[0].Key, "source")
	_require.Equal(*resp.BlobTagSet[0].Value, "tags")
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, _, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Invoke UploadBlobFromURL without SAS
	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlob.URL(), nil)
	_require.Error(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLWithHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Invoke UploadBlobFromURL
	tier := blob.AccessTierCool
	options := blockblob.UploadBlobFromURLOptions{
		Tags:        testcommon.BasicBlobTagsMap,
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		Tier:        &tier,
	}

	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Check dest and source properties
	resp, err := destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
	h := blob.ParseHTTPHeaders(resp)
	h.BlobContentMD5 = nil // the service generates a MD5 value, omit before comparing
	_require.EqualValues(h, testcommon.BasicHeaders)
	_require.EqualValues(resp.AccessTier, &tier)
	tagcount := int64(len(testcommon.BasicBlobTagsMap))
	_require.EqualValues(resp.TagCount, &tagcount)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLWithIntent() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	requestIntent := blob.FileRequestIntentTypeBackup

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	options := blockblob.UploadBlobFromURLOptions{
		Tags:              testcommon.BasicBlobTagsMap,
		HTTPHeaders:       &testcommon.BasicHeaders,
		Metadata:          testcommon.BasicMetadata,
		FileRequestIntent: &requestIntent,
	}

	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Check dest and source properties
	_, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Invoke UploadBlobFromURL
	options := blockblob.UploadBlobFromURLOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Get CPKInfo and compare
	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}

	getResp, err := destBlob.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.NoError(err)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(getResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, _, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create Blob with CPK
	bbName := testcommon.GenerateBlobName(testName)
	srcBlob := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, nil, &encryptionScope)
	expiryTime := time.Now().UTC().Add(15 * time.Minute)
	_require.NoError(err)

	// Create SAS credentials to get SAS URL for source
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}

	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	// Invoke UploadBlobFromURL
	options := blockblob.UploadBlobFromURLOptions{
		CPKScopeInfo: &encryptionScope,
	}

	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NotNil(pbResp)
	_require.NoError(err)

	// Compare EncryptionScope info
	getResp, err := destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(*getResp.EncryptionScope, *encryptionScope.EncryptionScope)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlSourceContentMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, sourceData := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Invoke UploadBlobFromURL
	sourceDataMD5Value := md5.Sum(sourceData)
	sourceContentMD5 := sourceDataMD5Value[:]
	options := blockblob.UploadBlobFromURLOptions{
		SourceContentMD5: sourceContentMD5,
	}

	resp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)
	_require.NotEqual(*resp.ETag, "")
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.Equal(resp.Date.IsZero(), false)
	_require.EqualValues(resp.ContentMD5, sourceDataMD5Value[:])

	// Try UploadBlobFromURL with bad MD5
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	options2 := blockblob.UploadBlobFromURLOptions{
		SourceContentMD5: badMD5,
	}
	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options2)
	_require.Error(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlSourceIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, sourceData := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Get source properties
	resp, err := srcBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL
	options := blockblob.UploadBlobFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfMatch: resp.ETag,
		},
	}

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)

	// Get dest properties
	_, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Download data from destination
	destBuffer := make([]byte, 4*1024)
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, nil)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlSourceIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Getting random etag
	randomEtag := azcore.ETag("a")
	accessConditions := blob.SourceModifiedAccessConditions{
		SourceIfMatch: &randomEtag,
	}

	// Invoke UploadBlobFromURL, should fail so validate error
	options := blockblob.UploadBlobFromURLOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SourceConditionNotMet)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlSourceIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, sourceData := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Invoke UploadBlobFromURL
	options := blockblob.UploadBlobFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfNoneMatch: to.Ptr(azcore.ETag("a")),
		},
	}
	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)

	// Download data from destination
	destBuffer := make([]byte, 4*1024)
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, nil)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlSourceIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Get source properties
	resp, err := srcBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL, should fail and validate error
	options := blockblob.UploadBlobFromURLOptions{
		SourceModifiedAccessConditions: &blob.SourceModifiedAccessConditions{
			SourceIfNoneMatch: resp.ETag,
		},
	}

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CannotVerifyCopySource)
	_require.ErrorContains(err, "304")
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	cResp, err := srcBlob.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL
	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)

	_, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	cResp, err := srcBlob.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL, should fail and validate error
	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	cResp, err := srcBlob.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL
	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)

	_, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	cResp, err := srcBlob.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL
	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.Error(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlobPutBlobFromUrlDestIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Get ETag from dest blob
	resp, err := destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)

	resp, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Get ETag from dest blob
	resp, err := destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Invoke UploadBlobFromURL
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}

	// Set metadata on dest blob
	metadata := make(map[string]*string)
	metadata["bla"] = to.Ptr("bla")
	_, err = destBlob.SetMetadata(context.Background(), metadata, nil)
	_require.NoError(err)

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Get Etag from dest blob
	resp, err := destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}

	_, err = destBlob.SetMetadata(context.Background(), nil, nil) // SetMetadata changes the blob's etag
	_require.NoError(err)

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)

	resp, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromUrlDestIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, _, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Get ETag from dest blob
	resp, err := destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)

	// Invoke PutBlobFromURL, should fail and validate error
	options := blockblob.UploadBlobFromURLOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}

	_, err = destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLCopySourceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient, srcBlob, destBlob, srcBlobURLWithSAS, _ := setUpPutBlobFromURLTest(testName, _require, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Set tier to Cool and check tier has been set
	_, err = srcBlob.SetTier(context.Background(), blob.AccessTierCool, nil)
	_require.NoError(err)

	resp, err := srcBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(resp.AccessTier, to.Ptr("Cool"))

	// Invoke UploadBlobForURL
	// CopySourceBlobProperties is true by default, trying false here
	options := blockblob.UploadBlobFromURLOptions{
		CopySourceBlobProperties: to.Ptr(false),
	}

	pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &options)
	_require.NoError(err)
	_require.NotNil(pbResp)

	// Access Tier for dest blob will not be Cool
	resp, err = destBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotEqual(resp.AccessTier, to.Ptr("Cool"))
}

func (s *BlockBlobRecordedTestsSuite) TestPutBlobFromURLCopySourceAuth() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Getting AAD Authentication
	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create source and destination blobs
	srcBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "src "+testName, containerClient)
	destBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "dest"+testName, containerClient)

	// Upload some data to source
	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	_, err = srcBBClient.Upload(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)

	// Getting token
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	_require.NoError(err)

	options := blockblob.UploadBlobFromURLOptions{
		CopySourceAuthorization: to.Ptr("Bearer " + token.Token),
	}

	pbResp, err := destBBClient.UploadBlobFromURL(context.Background(), srcBBClient.URL(), &options)
	_require.NoError(err)
	_require.NotNil(pbResp)

	// Download data from destination
	destBuffer := make([]byte, 4*1024)
	_, err = srcBBClient.DownloadBuffer(context.Background(), destBuffer, nil)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)

}

func (s *BlockBlobRecordedTestsSuite) TestPutBlobFromURLCopySourceAuthNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create source and destination blobs
	srcBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "src "+testName, containerClient)
	destBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, "dest"+testName, containerClient)

	// Upload some data to source
	contentSize := 4 * 1024 // 4KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	_, err = srcBBClient.Upload(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)

	options := blockblob.UploadBlobFromURLOptions{
		CopySourceAuthorization: to.Ptr("Bearer XXXXXXXXXXXXXX"),
	}

	_, err = destBBClient.UploadBlobFromURL(context.Background(), srcBBClient.URL(), &options)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CannotVerifyCopySource)
}

func (s *BlockBlobUnrecordedTestsSuite) TestPutBlobFromURLWithTier() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	src := testcommon.GenerateBlobName("src" + testName)
	srcBlob := testcommon.CreateNewBlockBlob(context.Background(), _require, src, containerClient)

	// Create SAS for source and get SAS URL
	expiryTime := time.Now().UTC().Add(15 * time.Minute)
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot, blob.AccessTierCold} {
		dest := testcommon.GenerateBlobName("dest" + string(tier) + testName)
		destBlob := testcommon.CreateNewBlockBlob(context.Background(), _require, dest, containerClient)

		opts := blockblob.UploadBlobFromURLOptions{
			Tier: &tier,
		}
		// Invoke UploadBlobFromURL
		pbResp, err := destBlob.UploadBlobFromURL(context.Background(), srcBlobURLWithSAS, &opts)
		_require.NotNil(pbResp)
		_require.NoError(err)

		getResp, err := destBlob.GetProperties(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(*getResp.AccessTier, string(tier))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestPutBlockListWithImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.NoError(err)
	policy := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.NoError(err)

	blockIDs := testcommon.GenerateBlockIDsList(1)
	_, err = bbClient.StageBlock(context.Background(), blockIDs[0], streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	legalHold := true
	options := blockblob.CommitBlockListOptions{
		ImmutabilityPolicyExpiryTime: &currentTime,
		ImmutabilityPolicyMode:       &policy,
		LegalHold:                    &legalHold,
	}

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &options)
	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	policy1 := blob.ImmutabilityPolicyMode("unlocked")
	_require.Equal(resp.ImmutabilityPolicyMode, &policy1)

	time.Sleep(time.Second * 7)

	_, err = bbClient.SetLegalHold(context.Background(), false, nil)
	_require.NoError(err)

	_, err = bbClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.NoError(err)

	_, err = bbClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobMetadataNotEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		Metadata: testcommon.BasicMetadata,
	})
	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	actualMetadata := resp.Metadata
	_require.NotNil(actualMetadata)
	_require.EqualValues(actualMetadata, testcommon.BasicMetadata)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, nil)
	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.Metadata)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		Metadata: map[string]*string{"In valid!": to.Ptr("bar")},
	})
	_require.Error(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)
	// _require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	loc, err := time.LoadLocation("Asia/Kolkata")
	_require.NoError(err)
	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, -10)
	currentTime = currentTime.In(loc) // converting to IST

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NoError(err)
	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)
	// _require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	loc, err := time.LoadLocation("EST")
	_require.NoError(err)
	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, 10)
	currentTime = currentTime.In(loc) // converting to EST

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}

	_, err = bbClient.Upload(context.Background(), rsc, &uploadBlockBlobOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)
	// _require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	loc, err := time.LoadLocation("Asia/Kolkata")
	_require.NoError(err)
	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, 10)
	currentTime = currentTime.In(loc) // converting to IST

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(context.Background(), rsc, &uploadBlockBlobOptions)
	_require.NoError(err)

	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blockBlobName, containerClient)

	createResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)
	// _require.Equal(createResp.RawResponse.StatusCode, 201)
	_require.NotNil(createResp.Date)

	loc, err := time.LoadLocation("EST")
	_require.NoError(err)
	currentTime := testcommon.GetRelativeTimeFromAnchor(createResp.Date, -10)
	currentTime = currentTime.In(loc) // converting to EST

	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader(nil)), &uploadBlockBlobOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.NoError(err)

	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)

	ifMatch := azcore.ETag("garbage")
	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &ifMatch,
			},
		},
	}
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(body), &uploadBlockBlobOptions)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	_, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	ifNoneMatch := azcore.ETag("garbage")
	uploadBlockBlobOptions := blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &ifNoneMatch,
			},
		},
	}

	_, err = bbClient.Upload(context.Background(), rsc, &uploadBlockBlobOptions)
	_require.NoError(err)

	testcommon.ValidateUpload(context.Background(), _require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlobIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.Upload(context.Background(), rsc, &blockblob.UploadOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func validateBlobCommitted(_require *require.Assertions, bbClient *blockblob.Client) {
	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.Len(resp.CommittedBlocks, 1)
}

func setupPutBlockListTest(t *testing.T, _require *require.Assertions, testName string) (*container.Client, *blockblob.Client, []string) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	blockIDs := testcommon.GenerateBlockIDsList(1)
	_, err = bbClient.StageBlock(context.Background(), blockIDs[0], streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)
	return containerClient, bbClient, blockIDs
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListHTTPHeadersEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		HTTPHeaders: &blob.HTTPHeaders{BlobContentDisposition: &testcommon.BlobContentDisposition},
	})
	_require.NoError(err)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, nil)
	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ContentDisposition)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_require.NoError(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertyResp, err := containerClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getPropertyResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertyResp.Date, 10)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(commitBlockListResp.Date, 10)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)
	_require.NoError(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	commitBlockListResp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)
	_require.NotNil(commitBlockListResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(commitBlockListResp.Date, -10)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)

	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag}},
	})
	_require.NoError(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)

	eTag := azcore.ETag("garbage")
	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)

	eTag := azcore.ETag("garbage")
	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)
	_require.NoError(err)

	validateBlobCommitted(_require, bbClient)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil) // The bbClient must actually exist to have a modifed time
	_require.NoError(err)

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		AccessConditions: &blob.AccessConditions{ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(context.Background(), blockIDs, &commitBlockListOptions)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListValidateData() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil)
	_require.NoError(err)

	resp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	data, err := io.ReadAll(resp.Body)
	_require.NoError(err)
	_require.Equal(string(data), testcommon.BlockBlobDefaultData)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobPutBlockListModifyBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, bbClient, blockIDs := setupPutBlockListTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := bbClient.CommitBlockList(context.Background(), blockIDs, nil)
	_require.NoError(err)

	_, err = bbClient.StageBlock(context.Background(), "0001", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.NoError(err)
	_, err = bbClient.StageBlock(context.Background(), "0010", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.NoError(err)
	_, err = bbClient.StageBlock(context.Background(), "0011", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.NoError(err)
	_, err = bbClient.StageBlock(context.Background(), "0100", streaming.NopCloser(bytes.NewReader([]byte("new data"))), nil)
	_require.NoError(err)

	_, err = bbClient.CommitBlockList(context.Background(), []string{"0001", "0011"}, nil)
	_require.NoError(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.Len(resp.CommittedBlocks, 2)
	committed := resp.CommittedBlocks
	_require.Equal(*(committed[0].Name), "0001")
	_require.Equal(*(committed[1].Name), "0011")
	_require.Nil(resp.UncommittedBlocks)
}

func (s *BlockBlobRecordedTestsSuite) TestSetTierOnBlobUpload() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot, blob.AccessTierCold} {
		blobName := strings.ToLower(string(tier)) + testcommon.GenerateBlobName(testName)
		bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

		uploadBlockBlobOptions := blockblob.UploadOptions{
			HTTPHeaders: &testcommon.BasicHeaders,
			Tier:        &tier,
		}
		_, err := bbClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &uploadBlockBlobOptions)
		_require.NoError(err)

		resp, err := bbClient.GetProperties(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(*resp.AccessTier, string(tier))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSetTierOnCommit() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := "test" + testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierCool, blob.AccessTierHot, blob.AccessTierCold} {
		blobName := strings.ToLower(string(tier)) + testcommon.GenerateBlobName(testName)
		bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
		_, err := bbClient.StageBlock(context.Background(), blockID, streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
		_require.NoError(err)

		_, err = bbClient.CommitBlockList(context.Background(), []string{blockID}, &blockblob.CommitBlockListOptions{
			Tier: &tier,
		})
		_require.NoError(err)

		resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeCommitted, nil)
		_require.NoError(err)
		_require.NotNil(resp.BlockList)
		_require.NotNil(resp.CommittedBlocks)
		_require.Nil(resp.UncommittedBlocks)
		_require.Len(resp.CommittedBlocks, 1)

		getResp, err := bbClient.GetProperties(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(*getResp.AccessTier, string(tier))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestCommitBlockListWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := "test" + testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4 KB
	_, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	contentMD5 := md5.Sum(sourceData)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	_, err = bbClient.StageBlock(context.Background(), blockID, streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	// CommitBlockList is a multipart upload, user generated checksum cannot be passed
	_, err = bbClient.CommitBlockList(context.Background(), []string{blockID}, &blockblob.CommitBlockListOptions{
		TransactionalContentMD5: contentMD5[:],
	})
	_require.Error(err, bloberror.UnsupportedChecksum)
}

func (s *BlockBlobRecordedTestsSuite) TestCommitBlockListWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := "test" + testcommon.GenerateContainerName(testName) + "1"
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4 KB
	_, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)

	blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	_, err = bbClient.StageBlock(context.Background(), blockID, streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	// CommitBlockList is a multipart upload, user generated checksum cannot be passed
	_, err = bbClient.CommitBlockList(context.Background(), []string{blockID}, &blockblob.CommitBlockListOptions{
		TransactionalContentCRC64: crc,
	})
	_require.Error(err, bloberror.UnsupportedChecksum)
}

func (s *BlockBlobUnrecordedTestsSuite) TestCopyBlockBlobFromURLWithEncryptionScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// set up source blob
	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := testcommon.GetDataAndReader(testName, contentSize)

	srcBlob := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))
	_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(contentReader), nil)
	_require.NoError(err)

	// Get source blob url with SAS for StageFromURL.
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	destBlobName := testcommon.GenerateBlobName(testName)
	destBlob := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(destBlobName))

	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.EncryptionScopeEnvVar)
	_require.NoError(err)
	cpk := blob.CPKScopeInfo{
		EncryptionScope: to.Ptr(encryptionScope),
	}
	copyBlockBlobFromURLOptions := blob.CopyFromURLOptions{
		CPKScopeInfo: &cpk,
	}
	resp, err := destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
	_require.NoError(err)
	_require.Equal(*resp.CopyStatus, "success")
	_require.Equal(*resp.EncryptionScope, encryptionScope)
}

func (s *BlockBlobUnrecordedTestsSuite) TestGetSASURLBlockBlobClient() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)

	// Creating service client with credentials
	serviceClient, err := service.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	_require.NoError(err)

	// Creating container client
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, serviceClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Creating block blob client with credentials
	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	// Adding SAS and options
	permissions := sas.BlobPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(5 * time.Minute)

	sasUrl, err := bbClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	// Get new blob client with sasUrl and attempt GetProperties
	newClient, err := blob.NewClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	_, err = newClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestSetTierOnCopyBlockBlobFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	const contentSize = 4 * 1024 * 1024 // 4 MB
	contentReader, _ := testcommon.GetDataAndReader(testName, contentSize)

	srcBlob := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	tier := blob.AccessTierCool
	_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(contentReader), &blockblob.UploadOptions{Tier: &tier})
	_require.NoError(err)
	// _require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)

	// Get source blob url with SAS for StageFromURL.
	expiryTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	if err != nil {
		s.T().Fatal("Couldn't fetch credential because " + err.Error())
	}
	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    expiryTime,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot, blob.AccessTierCold} {
		destBlobName := strings.ToLower(string(tier)) + testcommon.GenerateBlobName(testName)
		destBlob := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(destBlobName))

		copyBlockBlobFromURLOptions := blob.CopyFromURLOptions{
			Tier:     &tier,
			Metadata: testcommon.BasicMetadata,
		}
		resp, err := destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
		_require.NoError(err)
		// _require.Equal(resp.RawResponse.StatusCode, 202)
		_require.Equal(*resp.CopyStatus, "success")

		destBlobPropResp, err := destBlob.GetProperties(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(*destBlobPropResp.AccessTier, string(tier))
	}
}

func (s *BlockBlobUnrecordedTestsSuite) TestCopyBlockBlobFromUrlSourceContentMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 8 * 1024 // 8 KB
	body, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	contentMD5 := md5.Sum(sourceData)

	srcBlob := containerClient.NewBlockBlobClient("src" + testName)
	destBlob := containerClient.NewBlockBlobClient("dest" + testName)

	// Prepare source bbClient for copy.
	_, err = srcBlob.Upload(context.Background(), streaming.NopCloser(body), nil)
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Get source blob url with SAS for CopyFromURL.
	sasQueryParams, err := sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration,
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.String()

	// Invoke CopyFromURL.
	sourceContentMD5 := contentMD5[:]
	resp, err := destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &blob.CopyFromURLOptions{
		SourceContentMD5: sourceContentMD5,
	})
	_require.NoError(err)
	_require.EqualValues(resp.ContentMD5, sourceContentMD5)

	// Provide bad MD5 and make sure the copy fails
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	resp, err = destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &blob.CopyFromURLOptions{
		SourceContentMD5: badMD5,
	})
	_require.Error(err)
}

// func (s *BlockBlobUnrecordedTestsSuite) TestSetTierOnStageBlockFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 8 * 1024 // 8 KB
//	content := make([]byte, contentSize)
//	body := bytes.NewReader(content)
//	rsc := streaming.NopCloser(body)
//	ctx := context.Background()
//	srcBlob := containerClient.NewBlockBlobClient("src" + testcommon.GenerateBlobName(testName))
//	destBlob := containerClient.NewBlockBlobClient("dst" + testcommon.GenerateBlobName(testName))
//	tier := AccessTierCool
//	_, err = srcBlob.Upload(context.Background(), rsc, &blockblob.UploadOptions{Tier: &tier})
//	_require.NoError(err)
//	//_require.Equal(uploadSrcResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := testcommon.GetGenericSharedKeyCredential(nil, testcommon.TestAccountDefault)
//	_require.NoError(err)
//	srcBlobParts.SAS, err = blob.SASSignatureValues{
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
//	options1 := blockblob.StageBlockFromURLOptions{
//		Offset: &offset1,
//		Count:  &count1,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.NoError(err)
//	// _require.Equal(stageResp1.RawResponse.StatusCode, 201)
//	_require.Nil(stageResp1.ContentMD5)
//	_require.NotEqual(*stageResp1.RequestID, "")
//	_require.NotEqual(*stageResp1.Version, "")
//	_require.Equal(stageResp1.Date.IsZero(), false)
//
//	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
//	offset2, count2 := int64(4*1024), int64(blob.CountToEnd)
//	options2 := blockblob.StageBlockFromURLOptions{
//		Offset: &offset2,
//		Count:  &count2,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.NoError(err)
//	// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp2.ContentMD5, "")
//	_require.NotEqual(stageResp2.RequestID, "")
//	_require.NotEqual(stageResp2.Version, "")
//	_require.Equal(stageResp2.Date.IsZero(), false)
//
//	// Check block list.
//	blockList, err := destBlob.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
//	_require.NoError(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.NotNil(blockList.BlockList.UncommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	// Commit block list.
//	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &blockblob.CommitBlockListOptions{
//		Tier: &tier,
//	})
//	_require.NoError(err)
//	// _require.Equal(listResp.RawResponse.StatusCode,  201)
//	_require.NotNil(listResp.LastModified)
//	_require.Equal((*listResp.LastModified).IsZero(), false)
//	_require.NotNil(listResp.ETag)
//	_require.NotNil(listResp.RequestID)
//	_require.NotNil(listResp.Version)
//	_require.NotNil(listResp.Date)
//	_require.Equal((*listResp.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
//	_require.NoError(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(nil))
//	_require.NoError(err)
//	_require.EqualValues(destData, content)
//
//	// Get properties to validate the tier
//	destBlobPropResp, err := destBlob.GetProperties(context.Background(), nil)
//	_require.NoError(err)
//	_require.Equal(*destBlobPropResp.AccessTier, string(tier))
// }

func (s *BlockBlobRecordedTestsSuite) TestSetStandardBlobTierWithRehydratePriority() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	standardTier, rehydrateTier, rehydratePriority := blob.AccessTierArchive, blob.AccessTierCool, blob.RehydratePriorityStandard
	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)

	_, err = bbClient.SetTier(context.Background(), standardTier, &blob.SetTierOptions{
		RehydratePriority: &rehydratePriority,
	})
	_require.NoError(err)

	getResp1, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp1.AccessTier, string(standardTier))

	_, err = bbClient.SetTier(context.Background(), rehydrateTier, nil)
	_require.NoError(err)

	getResp2, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCool))
}

func (s *BlockBlobRecordedTestsSuite) TestRehydrateStatus() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName1 := "rehydration_test_blob_1"
	blobName2 := "rehydration_test_blob_2"
	blobName3 := "rehydration_test_blob_3"

	bbClient1 := testcommon.GetBlockBlobClient(blobName1, containerClient)
	reader1, _ := testcommon.GenerateData(1024)
	_, err = bbClient1.Upload(context.Background(), reader1, nil)
	_require.NoError(err)
	_, err = bbClient1.SetTier(context.Background(), blob.AccessTierArchive, nil)
	_require.NoError(err)
	_, err = bbClient1.SetTier(context.Background(), blob.AccessTierCool, nil)
	_require.NoError(err)

	getResp1, err := bbClient1.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp1.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp1.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCool))

	pager := containerClient.NewListBlobsFlatPager(nil)
	var blobs []*container.BlobItem
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		blobs = append(blobs, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.GreaterOrEqual(len(blobs), 1)
	_require.Equal(*blobs[0].Properties.AccessTier, blob.AccessTierArchive)
	_require.Equal(*blobs[0].Properties.ArchiveStatus, blob.ArchiveStatusRehydratePendingToCool)

	// ------------------------------------------

	bbClient2 := testcommon.GetBlockBlobClient(blobName2, containerClient)
	reader2, _ := testcommon.GenerateData(1024)
	_, err = bbClient2.Upload(context.Background(), reader2, nil)
	_require.NoError(err)
	_, err = bbClient2.SetTier(context.Background(), blob.AccessTierArchive, nil)
	_require.NoError(err)
	_, err = bbClient2.SetTier(context.Background(), blob.AccessTierHot, nil)
	_require.NoError(err)

	getResp2, err := bbClient2.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp2.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToHot))

	// ------------------------------------------

	bbClient3 := testcommon.GetBlockBlobClient(blobName3, containerClient)
	reader3, _ := testcommon.GenerateData(1024)
	_, err = bbClient3.Upload(context.Background(), reader3, nil)
	_require.NoError(err)
	_, err = bbClient3.SetTier(context.Background(), blob.AccessTierArchive, nil)
	_require.NoError(err)
	_, err = bbClient3.SetTier(context.Background(), blob.AccessTierCold, nil)
	_require.NoError(err)

	getResp3, err := bbClient3.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp3.AccessTier, string(blob.AccessTierArchive))
	_require.Equal(*getResp3.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToCold))
}

func (s *BlockBlobRecordedTestsSuite) TestCopyBlobWithRehydratePriority() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sourceBlobName := testcommon.GenerateBlobName(testName)
	sourceBBClient := testcommon.CreateNewBlockBlob(context.Background(), _require, sourceBlobName, containerClient)

	blobTier, rehydratePriority := blob.AccessTierArchive, blob.RehydratePriorityHigh

	copyBlobName := "copy" + sourceBlobName
	destBBClient := testcommon.GetBlockBlobClient(copyBlobName, containerClient)
	_, err = destBBClient.StartCopyFromURL(context.Background(), sourceBBClient.URL(), &blob.StartCopyFromURLOptions{
		RehydratePriority: &rehydratePriority,
		Tier:              &blobTier,
	})
	_require.NoError(err)

	getResp1, err := destBBClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp1.AccessTier, string(blobTier))

	_, err = destBBClient.SetTier(context.Background(), blob.AccessTierHot, nil)
	_require.NoError(err)

	getResp2, err := destBBClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp2.ArchiveStatus, string(blob.ArchiveStatusRehydratePendingToHot))
}

func (s *BlockBlobRecordedTestsSuite) TestCopyWithTier() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	for _, tier := range []blob.AccessTier{blob.AccessTierArchive, blob.AccessTierCool, blob.AccessTierHot, blob.AccessTierCold} {
		src := testcommon.GenerateBlobName("src" + string(tier) + testName)
		srcBlob := testcommon.CreateNewBlockBlob(context.Background(), _require, src, containerClient)

		dest := testcommon.GenerateBlobName("dest" + string(tier) + testName)
		destBlob := testcommon.CreateNewBlockBlob(context.Background(), _require, dest, containerClient)

		_, err = destBlob.StartCopyFromURL(context.Background(), srcBlob.URL(), &blob.StartCopyFromURLOptions{
			Tier: &tier,
		})
		_require.NoError(err)

		getResp, err := destBlob.GetProperties(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(*getResp.AccessTier, string(tier))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestBlobServiceClientDelete() {
	_require := require.New(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	code := 404
	testcommon.RunTestRequiringServiceProperties(context.Background(), _require, svcClient, string(rune(code)),
		testcommon.EnableSoftDelete, func(context.Context, *require.Assertions, *service.Client) error { return nil }, testcommon.DisableSoftDelete)
}

func setAndCheckBlockBlobTier(_require *require.Assertions, bbClient *blockblob.Client, tier blob.AccessTier) {
	_, err := bbClient.SetTier(context.Background(), tier, nil)
	_require.NoError(err)

	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSetTierAllTiersOnBlockBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierHot)
	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierCool)
	setAndCheckBlockBlobTier(_require, bbClient, blob.AccessTierArchive)

}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobGetPropertiesUsingVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	blobProp, _ := bbClient.GetProperties(context.Background(), nil)

	uploadBlockBlobOptions := blockblob.UploadOptions{
		Metadata: testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	uploadResp, err := bbClient.Upload(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), &uploadBlockBlobOptions)
	_require.NoError(err)
	_require.NotNil(uploadResp.VersionID)
	blobProp, _ = bbClient.GetProperties(context.Background(), nil)
	_require.EqualValues(uploadResp.VersionID, blobProp.VersionID)
	_require.EqualValues(uploadResp.LastModified, blobProp.LastModified)
	_require.Equal(*uploadResp.ETag, *blobProp.ETag)
	_require.Equal(*blobProp.IsCurrentVersion, true)
}

func (s *BlockBlobRecordedTestsSuite) TestGetSetBlobMetadataWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, &testcommon.TestCPKByValue, nil)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Error(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	resp, err := bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.NoError(err)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}

	// Get blob properties without encryption key should fail the request.
	_, err = bbClient.GetProperties(context.Background(), nil)
	_require.Error(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	getResp, err := bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.NoError(err)
	_require.NotNil(getResp.Metadata)
	_require.Len(getResp.Metadata, len(testcommon.BasicMetadata))
	_require.EqualValues(getResp.Metadata, testcommon.BasicMetadata)

	_, err = bbClient.SetMetadata(context.Background(), map[string]*string{}, &setBlobMetadataOptions)
	_require.NoError(err)

	getResp, err = bbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.NoError(err)
	_require.Nil(getResp.Metadata)
}

func (s *BlockBlobRecordedTestsSuite) TestGetSetBlobMetadataWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, nil, &encryptionScope)

	// Set blob metadata without encryption key should fail the request.
	_, err = bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, nil)
	_require.Error(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		CPKScopeInfo: &encryptionScope,
	}
	resp, err := bbClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.NoError(err)
	_require.EqualValues(*resp.EncryptionScope, *encryptionScope.EncryptionScope)

	getResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(getResp.Metadata)
	_require.Len(getResp.Metadata, len(testcommon.BasicMetadata))
	_require.EqualValues(getResp.Metadata, testcommon.BasicMetadata)

	_, err = bbClient.SetMetadata(context.Background(), map[string]*string{}, &setBlobMetadataOptions)
	_require.NoError(err)

	getResp, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(getResp.Metadata)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSnapshotWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, &testcommon.TestCPKByValue, nil)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.Error(err)

	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
		CPKInfo: &testcommon.TestInvalidCPKByValue,
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_require.Error(err)

	createBlobSnapshotOptions1 := blob.CreateSnapshotOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions1)
	_require.NoError(err)
	_require.Equal(*resp.IsServerEncrypted, false)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	dResp, err := snapshotURL.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(*dResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
	}
	_, err = snapshotURL.Delete(context.Background(), nil)
	_require.NoError(err)

	// Get blob properties of snapshot without encryption key should fail the request.
	_, err = snapshotURL.GetProperties(context.Background(), nil)
	_require.Error(err)

	// _assert(err.(StorageError).Response().StatusCode, chk.Equals, 404)
}

func (s *BlockBlobRecordedTestsSuite) TestBlobSnapshotWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlobWithCPK(context.Background(), _require, bbName, containerClient, nil, &encryptionScope)

	// Create Snapshot of an encrypted blob without encryption key should fail the request.
	_, err = bbClient.CreateSnapshot(context.Background(), nil)
	_require.Error(err)

	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
		CPKScopeInfo: &testcommon.TestInvalidCPKByScope,
	}
	_, err = bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions)
	_require.Error(err)

	createBlobSnapshotOptions1 := blob.CreateSnapshotOptions{
		CPKScopeInfo: &encryptionScope,
	}
	resp, err := bbClient.CreateSnapshot(context.Background(), &createBlobSnapshotOptions1)
	_require.NoError(err)
	_require.Equal(*resp.IsServerEncrypted, false)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKScopeInfo: &encryptionScope,
	}
	dResp, err := snapshotURL.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)
	_require.EqualValues(*dResp.EncryptionScope, *encryptionScope.EncryptionScope)

	_, err = snapshotURL.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobRecordedTestsSuite) TestCreateAndDownloadBlobSpecialCharactersWithVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	data := []rune("-._/()$=',~0123456789")
	for i := 0; i < len(data); i++ {
		blobName := "abc" + string(data[i])
		blobClient := containerClient.NewBlockBlobClient(blobName)
		resp, err := blobClient.Upload(context.Background(), streaming.NopCloser(strings.NewReader(string(data[i]))), nil)
		_require.NoError(err)
		_require.NotNil(resp.VersionID)

		blobClientWithVersionID, err := blobClient.WithVersionID(*resp.VersionID)
		_require.NoError(err)
		dResp, err := blobClientWithVersionID.DownloadStream(context.Background(), nil)
		_require.NoError(err)
		d1, err := io.ReadAll(dResp.Body)
		_require.NoError(err)
		_require.NotEqual(*dResp.Version, "")
		_require.EqualValues(string(d1), string(data[i]))
		_require.NotNil(dResp.VersionID)
		_require.Equal(*dResp.VersionID, *resp.VersionID)
	}
}

func (s *BlockBlobRecordedTestsSuite) TestDeleteSpecificBlobVersion() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		uploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"+strconv.Itoa(i)))), &blockblob.UploadOptions{
			Metadata: testcommon.BasicMetadata,
		})
		_require.NoError(err)
		_require.NotNil(uploadResp.VersionID)
		versions = append(versions, *uploadResp.VersionID)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})

	found := make([]*container.BlobItem, 0)
	for listPager.More() {
		resp, err := listPager.NextPage(context.Background())
		_require.NoError(err)
		if err != nil {
			break
		}
		found = append(found, resp.Segment.BlobItems...)
	}
	_require.Len(found, 5)

	// Deleting the 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.NoError(err)
		_, err = bbClientWithVersionID.Delete(context.Background(), nil)
		_require.NoError(err)
		// _require.Equal(deleteResp.RawResponse.StatusCode, 202)
	}

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})

	found = make([]*container.BlobItem, 0)
	for listPager.More() {
		resp, err := listPager.NextPage(context.Background())
		_require.NoError(err)
		if err != nil {
			break
		}
		found = append(found, resp.Segment.BlobItems...)
	}
	_require.Len(found, 2)

	for i := 3; i < 5; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.NoError(err)
		downloadResp, err := bbClientWithVersionID.DownloadStream(context.Background(), nil)
		_require.NoError(err)
		destData, err := io.ReadAll(downloadResp.Body)
		_require.NoError(err)
		_require.EqualValues(destData, "data"+strconv.Itoa(i))
	}
}

func (s *BlockBlobRecordedTestsSuite) TestUndeleteBlockBlobVersion() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		uploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"+strconv.Itoa(i)))), &blockblob.UploadOptions{
			Metadata: testcommon.BasicMetadata,
		})
		_require.NoError(err)
		_require.NotNil(uploadResp.VersionID)
		versions = append(versions, *uploadResp.VersionID)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 5)

	// Deleting the 1st, 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.NoError(err)
		_, err = bbClientWithVersionID.Delete(context.Background(), nil)
		_require.NoError(err)
	}

	// adding wait after delete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 2)

	_, err = bbClient.Undelete(context.Background(), nil)
	_require.NoError(err)

	// adding wait after undelete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 5)
}

func (s *BlockBlobRecordedTestsSuite) TestUndeleteBlockBlobSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	snapshots := make([]string, 0)
	for i := 0; i < 5; i++ {
		resp, err := bbClient.CreateSnapshot(context.Background(), nil)
		_require.NoError(err)
		_require.NotNil(resp.Snapshot)
		snapshots = append(snapshots, *resp.Snapshot)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6) // 5 snapshots and 1 current version

	// Deleting the 1st, 2nd and 3rd snapshots
	for i := 0; i < 3; i++ {
		bbClientWithSnapshot, err := bbClient.WithSnapshot(snapshots[i])
		_require.NoError(err)
		_, err = bbClientWithSnapshot.Delete(context.Background(), nil)
		_require.NoError(err)
	}

	// adding wait after delete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 3) // 2 snapshots and 1 current version

	_, err = bbClient.Undelete(context.Background(), nil)
	_require.NoError(err)

	// adding wait after undelete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6) // 5 snapshots and 1 current version
}

func (s *BlockBlobRecordedTestsSuite) TestPutBlockListReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = testcommon.BlockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(d)), nil)
		_require.NoError(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.NotNil(resp.Version)
		_require.NotEqual(*resp.Version, "")
	}

	commitResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
	_require.NoError(err)
	_require.NotNil(commitResp.VersionID)

	contentResp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	contentData, err := io.ReadAll(contentResp.Body)
	_require.NoError(err)
	_require.EqualValues(contentData, []uint8(strings.Join(data, "")))
}

func (s *BlockBlobUnrecordedTestsSuite) TestCreateBlockBlobReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	testSize := 2 * 1024 * 1024 // 1MB
	r, _ := testcommon.GetDataAndReader(testName, testSize)
	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	// Prepare source blob for copy.
	uploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.NotNil(uploadResp.VersionID)

	csResp, err := bbClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(csResp.RawResponse.StatusCode, 201)
	_require.NotNil(csResp.VersionID)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})

	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = bbClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.NoError(err)
	// _require.Equal(deleteResp.RawResponse.StatusCode, 202)

	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.NotEqual(len(found), 0)
}

func (s *BlockBlobRecordedTestsSuite) TestORSSource() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)

	getResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(getResp.ObjectReplicationRules)
}

// func (s *azblobTestSuite) TestSnapshotSAS() {
//	//Generate URLs ----------------------------------------------------------------------------------------------------
//	bsu := getServiceClient(nil)
//	containerClient, containerName := getContainerClient(bsu)
//	blobURL, blobName := getBlockBlobClient(c, containerClient)
//
//	_, err := containerClient.Create(context.Background(), nil)
//	defer containerClient.Delete(context.Background(), nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	//Create file in container, download from snapshot to test. --------------------------------------------------------
//	blobClient := containerClient.NewBlockBlobClient(blobName)
//	data := "Hello world!"
//
//	contentType := "text/plain"
//	uploadBlockBlobOptions := BlockBlobUploadOptions{
//		HTTPHeaders: &HTTPHeaders{
//			BlobContentType: &contentType,
//		},
//	}
//	_, err = blobClient.Upload(context.Background(), strings.NewReader(data), &uploadBlockBlobOptions)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	//Create a snapshot & URL
//	createSnapshot, err := blobClient.CreateSnapshot(context.Background(), nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	_assert(createSnapshot.Snapshot, chk.NotNil)
//
//	//Format snapshot time
//	snapTime, err := time.Parse(SnapshotTimeFormat, *createSnapshot.Snapshot)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	//Get credentials & current time
//	currentTime := time.Now().UTC()
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//
//	//Create SAS query
//	snapSASQueryParams, err := BlobSASSignatureValues{
//		StartTime:     currentTime,
//		ExpiryTime:    currentTime.Add(48 * time.Hour),
//		SnapshotTime:  snapTime,
//		Permissions:   "racwd",
//		ContainerName: containerName,
//		BlobName:      blobName,
//		Protocol:      SASProtocolHTTPS,
//	}.Sign(credential)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	time.Sleep(time.Second * 2)
//
//	//Attach SAS query to block blob URL
//	snapParts := NewBlobURLParts(blobURL.URL())
//	snapParts.SAS = snapSASQueryParams
//	sbUrl, err := NewBlockBlobClient(snapParts.URL(), azcore.AnonymousCredential(), nil)
//
//	//Test the snapshot
//	downloadResponse, err := sbUrl.Download(context.Background(), nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	downloadedData := &bytes.Buffer{}
//	reader := downloadResponse.Body(RetryReaderOptions{})
//	downloadedData.ReadFrom(reader)
//	reader.Close()
//
//	_assert(data, chk.Equals, downloadedData.String())
//
//	//Try to delete snapshot -------------------------------------------------------------------------------------------
//	_, err = sbUrl.Delete(context.Background(), nil)
//	if err != nil { //This shouldn't fail.
//		s.T().Fatal(err)
//	}
//
//	//Create a normal blob and attempt to use the snapshot SAS against it (assuming failure) ---------------------------
//	//If this succeeds, it means a normal SAS token was created.
//
//	uploadBlockBlobOptions1 := BlockBlobUploadOptions{
//		HTTPHeaders: &HTTPHeaders{
//			BlobContentType: &contentType,
//		},
//	}
//	fsbUrl := containerClient.NewBlockBlobClient("failsnap")
//	_, err = fsbUrl.Upload(context.Background(), strings.NewReader(data), &uploadBlockBlobOptions1)
//	if err != nil {
//		s.T().Fatal(err) //should succeed to create the blob via normal auth means
//	}
//
//	fsbUrlParts := NewBlobURLParts(fsbUrl.URL())
//	fsbUrlParts.SAS = snapSASQueryParams
//	fsbUrl, err = NewBlockBlobClient(fsbUrlParts.URL(), azcore.AnonymousCredential(), nil) //re-use fsbUrl as we don't need the sharedkey version anymore
//
//	resp, err := fsbUrl.Delete(context.Background(), nil)
//	if err == nil {
//		c.Fatal(resp) //This SHOULD fail. Otherwise we have a normal SAS token...
//	}
// }

func (s *BlockBlobUnrecordedTestsSuite) TestSetBlobTags() {
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

	_, err = bbClient.Upload(context.Background(), r, nil)
	_require.NoError(err)
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	_, err = bbClient.SetTags(context.Background(), blobTagsMap, nil)
	_require.NoError(err)
	// _require.Equal(blobSetTagsResponse.RawResponse.StatusCode, 204)

	blobGetTagsResponse, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, 3)
	for _, blobTag := range blobTagsSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *BlockBlobUnrecordedTestsSuite) TestSetBlobTagsWithLeaseId() {
	_require := require.New(s.T())
	testName := "bb" + s.T().Name()
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

	_, err = bbClient.Upload(context.Background(), r, nil)
	_require.NoError(err)
	blobLeaseClient, err := lease.NewBlobClient(bbClient, &lease.BlobClientOptions{
		LeaseID: proposedLeaseIDs[0],
	})
	_require.NoError(err)
	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.LeaseID())

	_, err = bbClient.SetTags(ctx, blobTagsMap, nil)
	_require.Error(err)

	_, err = bbClient.SetTags(ctx, blobTagsMap, &blob.SetTagsOptions{AccessConditions: &blob.AccessConditions{
		LeaseAccessConditions: &blob.LeaseAccessConditions{LeaseID: blobLeaseClient.LeaseID()}}})
	_require.NoError(err)

	_, err = bbClient.GetTags(ctx, nil)
	_require.NoError(err)

	blobGetTagsResponse, err := bbClient.GetTags(ctx, &blob.GetTagsOptions{BlobAccessConditions: &blob.AccessConditions{
		LeaseAccessConditions: &blob.LeaseAccessConditions{LeaseID: blobLeaseClient.LeaseID()}}})
	_require.NoError(err)

	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, 3)
	for _, blobTag := range blobTagsSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *BlockBlobUnrecordedTestsSuite) TestSetBlobTagsWithVID() {
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

	blockBlobUploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.NoError(err)
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	versionId1 := blockBlobUploadResp.VersionID

	blockBlobUploadResp, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("updated_data"))), nil)
	_require.NoError(err)
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)
	versionId2 := blockBlobUploadResp.VersionID

	setTagsBlobOptions := blob.SetTagsOptions{
		VersionID: versionId1,
	}
	_, err = bbClient.SetTags(context.Background(), blobTagsMap, &setTagsBlobOptions)
	_require.NoError(err)
	// _require.Equal(blobSetTagsResponse.RawResponse.StatusCode, 204)

	getTagsBlobOptions1 := blob.GetTagsOptions{
		VersionID: versionId1,
	}
	blobGetTagsResponse, err := bbClient.GetTags(context.Background(), &getTagsBlobOptions1)
	_require.NoError(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_require.NotNil(blobGetTagsResponse.BlobTagSet)
	_require.Len(blobGetTagsResponse.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResponse.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	getTagsBlobOptions2 := blob.GetTagsOptions{
		VersionID: versionId2,
	}
	blobGetTagsResponse, err = bbClient.GetTags(context.Background(), &getTagsBlobOptions2)
	_require.NoError(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_require.Nil(blobGetTagsResponse.BlobTagSet)
}

func (s *BlockBlobUnrecordedTestsSuite) TestUploadBlockBlobWithSpecialCharactersInTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	// _context := getTestContext(testName)
	// ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"+-./:=_ ": "firsttag", //nolint
		"tag2":     "+-./:=_",
		"+-./:=_1": "+-./:=_",
	}

	uploadBlockBlobOptions := blockblob.UploadOptions{
		Metadata:    testcommon.BasicMetadata,
		HTTPHeaders: &testcommon.BasicHeaders,
		Tags:        blobTagsMap,
	}
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"))), &uploadBlockBlobOptions)
	_require.NoError(err)
	// TODO: Check for metadata and header
	// _require.Equal(blockBlobUploadResp.RawResponse.StatusCode, 201)

	blobGetTagsResponse, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	_require.Len(blobGetTagsResponse.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResponse.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *BlockBlobRecordedTestsSuite) TestRemoveBlobTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	blobTagsMap := map[string]string{
		"sdk": "go",
	}

	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_require.NoError(err)

	// set blob tags
	_, err = bbClient.SetTags(context.Background(), blobTagsMap, nil)
	_require.NoError(err)

	blobGetTagsResponse1, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(blobGetTagsResponse1.BlobTagSet)
	_require.Len(blobGetTagsResponse1.BlobTagSet, 1)
	_require.Equal(blobTagsMap[*blobGetTagsResponse1.BlobTagSet[0].Key], *blobGetTagsResponse1.BlobTagSet[0].Value)

	// remove tags by passing nil blob tags map
	_, err = bbClient.SetTags(context.Background(), nil, nil)
	_require.NoError(err)

	blobGetTagsResponse2, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.Len(blobGetTagsResponse2.BlobTagSet, 0)

	// set blob tags
	_, err = bbClient.SetTags(context.Background(), blobTagsMap, nil)
	_require.NoError(err)

	blobGetTagsResponse3, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(blobGetTagsResponse3.BlobTagSet)
	_require.Len(blobGetTagsResponse3.BlobTagSet, 1)
	_require.Equal(blobTagsMap[*blobGetTagsResponse3.BlobTagSet[0].Key], *blobGetTagsResponse3.BlobTagSet[0].Value)

	// remove tags by passing empty blob tags map
	_, err = bbClient.SetTags(context.Background(), make(map[string]string), nil)
	_require.NoError(err)

	blobGetTagsResponse4, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.Len(blobGetTagsResponse4.BlobTagSet, 0)
}

func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	// _context := getTestContext(testName)
	// ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = testcommon.BlockIDIntToBase64(index)
		resp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(d)), nil)
		_require.NoError(err)
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
	commitResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, &commitBlockListOptions)
	_require.NoError(err)
	_require.NotNil(commitResp.VersionID)
	versionId := commitResp.VersionID

	contentResp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	contentData, err := io.ReadAll(contentResp.Body)
	_require.NoError(err)
	_require.EqualValues(contentData, []uint8(strings.Join(data, "")))

	getTagsBlobOptions := blob.GetTagsOptions{
		VersionID: versionId,
	}
	blobGetTagsResp, err := bbClient.GetTags(context.Background(), &getTagsBlobOptions)
	_require.NoError(err)
	_require.NotNil(blobGetTagsResp)
	_require.Len(blobGetTagsResp.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResp.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	blobGetTagsResp, err = bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(blobGetTagsResp)
	_require.Len(blobGetTagsResp.BlobTagSet, 3)
	for _, blobTag := range blobGetTagsResp.BlobTagSet {
		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

// func (s *BlockBlobUnrecordedTestsSuite) TestStageBlockFromURLWithTags() {
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
//	uploadSrcResp, err := srcBlob.Upload(context.Background(), r, &uploadBlockBlobOptions)
//	_require.NoError(err)
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
//	stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.NoError(err)
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
//	stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.NoError(err)
//	// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(*stageResp2.RequestID, "")
//	_require.NotEqual(*stageResp2.Version, "")
//	_require.NotNil(stageResp2.Date)
//	_require.Equal((*stageResp2.Date).IsZero(), false)
//
//	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//	_require.NoError(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	commitBlockListOptions := blockblob.CommitBlockListOptions{
//		Tags: blobTagsMap,
//	}
//	_, err = destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
//	_require.NoError(err)
//	// _require.Equal(listResp.RawResponse.StatusCode,  201)
//	//versionId := listResp.VersionID()
//
//	blobGetTagsResp, err := destBlob.GetTags(context.Background(), nil)
//	_require.NoError(err)
//	_require.Len(blobGetTagsResp.BlobTagSet, 3)
//	for _, blobTag := range blobGetTagsResp.BlobTagSet {
//		_require.Equal(blobTagsMap[*blobTag.Key], *blobTag.Value)
//	}
//
//	downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
//	_require.NoError(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(nil))
//	_require.NoError(err)
//	_require.EqualValues(destData, sourceData)
// }

// func (s *BlockBlobUnrecordedTestsSuite) TestCopyBlockBlobFromURLWithTags() {
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
//	_, err = srcBlob.Upload(context.Background(), r, &uploadBlockBlobOptions)
//	_require.NoError(err)
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
//	resp, err := destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
//	_require.NoError(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 202)
//	_require.NotEqual(*resp.ETag, "")
//	_require.NotEqual(*resp.RequestID, "")
//	_require.NotEqual(*resp.Version, "")
//	_require.Equal((*resp.Date).IsZero(), false)
//	_require.NotEqual(*resp.CopyID, "")
//	_require.EqualValues(resp.ContentMD5, sourceDataMD5Value[:])
//	_require.EqualValues(*resp.CopyStatus, "success")
//
//	downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
//	_require.NoError(err)
//	destData, err := io.ReadAll(downloadresp.BodyReader(nil))
//	_require.NoError(err)
//	_require.EqualValues(destData, sourceData)
//	_require.Equal(*downloadResp.TagCount, int64(1))
//
//	_, badMD5 := getRandomDataAndReader(16)
//	copyBlockBlobFromURLOptions2 := BlockBlobCopyFromURLOptions{
//		Tags:         blobTagsMap,
//		SourceContentMD5: badMD5,
//	}
//	_, err = destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
//	_require.Error(err)
//
//	copyBlockBlobFromURLOptions3 := BlockBlobCopyFromURLOptions{
//		Tags: blobTagsMap,
//	}
//	resp, err = destBlob.CopyFromURL(context.Background(), srcBlobURLWithSAS, &copyBlockBlobFromURLOptions3)
//	_require.NoError(err)
//	// _require.Equal(resp.RawResponse.StatusCode, 202)
// }

func (s *BlockBlobUnrecordedTestsSuite) TestGetPropertiesReturnsTagsCount() {
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
	_, err = bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"))), &uploadBlockBlobOptions)
	_require.NoError(err)

	getPropertiesResponse, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getPropertiesResponse.TagCount, int64(3))

	downloadResp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(downloadResp)
	_require.Equal(*downloadResp.TagCount, int64(3))
}

func (s *BlockBlobUnrecordedTestsSuite) TestSetBlobTagForSnapshot() {
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
	_, err = bbClient.SetTags(context.Background(), blobTagsMap, nil)
	_require.NoError(err)

	resp, err := bbClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)

	snapshotURL, _ := bbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp2.TagCount, int64(3))
}

// TODO: Once new pacer is done.
func (s *BlockBlobUnrecordedTestsSuite) TestListBlobReturnsTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	blobClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName, containerClient)
	blobTagsMap := map[string]string{
		"+-./:=_ ": "firsttag", //nolint
		"tag2":     "+-./:=_",
		"+-./:=_1": "+-./:=_",
	}

	_, err = blobClient.SetTags(context.Background(), blobTagsMap, nil)
	_require.NoError(err)
	// _require.Equal(resp.RawResponse.StatusCode,204)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Tags: true},
	})

	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
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

func (s *BlockBlobUnrecordedTestsSuite) TestFilterBlobsWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"1", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobTagsMap1 := map[string]string{
		"tag2": "tagsecond",
		"tag3": "tagthird",
	}
	blobTagsMap2 := map[string]string{
		"tag1":    "firsttag",
		"tag2":    "secondtag",
		"tag3":    "thirdtag",
		"tag key": "tag value", // tags with spaces
	}

	blobName1 := testcommon.GenerateBlobName(testName) + "1"
	blobClient1 := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName1, containerClient)
	_, err = blobClient1.SetTags(context.Background(), blobTagsMap1, nil)
	_require.NoError(err)

	blobName2 := testcommon.GenerateBlobName(testName) + "2"
	blobClient2 := testcommon.CreateNewBlockBlob(context.Background(), _require, blobName2, containerClient)
	_, err = blobClient2.SetTags(context.Background(), blobTagsMap2, nil)
	_require.NoError(err)
	time.Sleep(10 * time.Second)

	blobTagsResp, err := blobClient2.GetTags(context.Background(), nil)
	_require.NoError(err)
	blobTagsSet := blobTagsResp.BlobTagSet
	_require.NotNil(blobTagsSet)

	// Test invalid tag
	where := "\"tag4\"='fourthtag'"
	lResp, err := svcClient.FilterBlobs(context.Background(), where, nil)
	_require.NoError(err)
	_require.Equal(len(lResp.Blobs), 0)

	// Test multiple valid tags
	where = "\"tag1\"='firsttag'AND\"tag2\"='secondtag'"
	// where := "foo=\"value 1\""
	lResp, err = svcClient.FilterBlobs(context.Background(), where, nil)
	_require.NoError(err)
	_require.Len(lResp.Blobs[0].Tags.BlobTagSet, 2)
	_require.Equal(lResp.Blobs[0].Tags.BlobTagSet[0], blobTagsSet[1])
	_require.Equal(lResp.Blobs[0].Tags.BlobTagSet[1], blobTagsSet[2])

	// Test tags with spaces
	where = "\"tag key\"='tag value'"
	lResp, err = svcClient.FilterBlobs(context.Background(), where, nil)
	_require.NoError(err)
	_require.Len(lResp.Blobs[0].Tags.BlobTagSet, 1)
	_require.Equal(lResp.Blobs[0].Tags.BlobTagSet[0], blobTagsSet[0])

}

// func (s *BlockBlobUnrecordedTestsSuite) TestFilterBlobsUsingAccountSAS() {
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
//	_, err = containerClient.Create(context.Background(), Metadata{}, PublicAccessNone)
//	defer containerClient.Delete(context.Background(), LeaseAccessConditions{})
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	blobClient := containerClient.NewBlockBlobURL("temp")
//	_, err = blobClient.Upload(context.Background(), bytes.NewReader([]byte("random data")), HTTPHeaders{}, testcommon.BasicMetadata, LeaseAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
//	if err != nil {
//		s.T().Fail()
//	}
//
//	blobTagsMap := BlobTags{"tag1": "firsttag", "tag2": "secondtag", "tag3": "thirdtag"}
//	setBlobTagsResp, err := blobClient.SetTags(context.Background(), nil, nil, nil, nil, nil, nil, blobTagsMap)
//	_require.NoError(err)
//	_assert(setBlobTagsResp.StatusCode(), chk.Equals, 204)
//
//	blobGetTagsResp, err := blobClient.GetTags(context.Background(), nil, nil, nil, nil, nil)
//	_require.NoError(err)
//	_assert(blobGetTagsResp.StatusCode(), chk.Equals, 200)
//	_assert(blobGetTagsResp.BlobTagSet, chk.HasLen, 3)
//	for _, blobTag := range blobGetTagsResp.BlobTagSet {
//		_assert(blobTagsMap[blobTag.Key], chk.Equals, blobTag.Value)
//	}
//
//	time.Sleep(30 * time.Second)
//	where := "\"tag1\"='firsttag'AND\"tag2\"='secondtag'AND@container='" + containerName + "'"
//	_, err = serviceURL.FindBlobsByTags(context.Background(), nil, nil, &where, Marker{}, nil)
//	_require.NoError(err)
// }

func (s *BlockBlobRecordedTestsSuite) TestPutBlockAndPutBlockListWithCPK() {
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
			CPKInfo: &testcommon.TestCPKByValue,
		}
		_, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_require.NoError(err)
	}

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	resp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, &commitBlockListOptions)
	_require.NoError(err)

	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal(*resp.IsServerEncrypted, true)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(*resp.EncryptionKeySHA256, *(testcommon.TestCPKByValue.EncryptionKeySHA256))
	}

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.DownloadStream(context.Background(), nil)
	_require.Error(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	getResp, err := bbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)
	b := bytes.Buffer{}
	_, err = b.ReadFrom(getResp.Body)
	_require.NoError(err)
	err = getResp.Body.Close()
	_require.NoError(err)
	_require.Equal(b.String(), "AAA BBB CCC ")
	_require.EqualValues(*getResp.ETag, *resp.ETag)
	_require.EqualValues(*getResp.LastModified, *resp.LastModified)
}

func (s *BlockBlobRecordedTestsSuite) TestPutBlockAndPutBlockListWithCPKByScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
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
			CPKScopeInfo: &encryptionScope,
		}
		_, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], streaming.NopCloser(strings.NewReader(word)), &stageBlockOptions)
		_require.NoError(err)
	}

	commitBlockListOptions := blockblob.CommitBlockListOptions{
		CPKScopeInfo: &encryptionScope,
	}
	resp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, &commitBlockListOptions)
	_require.NoError(err)
	_require.NotNil(resp.ETag)
	_require.NotNil(resp.LastModified)
	_require.Equal(*resp.IsServerEncrypted, true)
	_require.EqualValues(*encryptionScope.EncryptionScope, *resp.EncryptionScope)

	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	_, err = bbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Error(err)

	downloadBlobOptions = blob.DownloadStreamOptions{
		CPKScopeInfo: &encryptionScope,
	}
	getResp, err := bbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)
	b := bytes.Buffer{}
	reader := getResp.Body
	_, err = b.ReadFrom(reader)
	_require.NoError(err)
	_ = reader.Close() // The client must close the response body when finished with it
	_require.Equal(b.String(), "AAA BBB CCC ")
	_require.EqualValues(*getResp.ETag, *resp.ETag)
	_require.EqualValues(*getResp.LastModified, *resp.LastModified)
	_require.Equal(*getResp.IsServerEncrypted, true)
	_require.EqualValues(*getResp.EncryptionScope, *encryptionScope.EncryptionScope)
}

// func (s *BlockBlobUnrecordedTestsSuite) TestPutBlockFromURLAndCommitWithCPK() {
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
//	_, err = srcBlob.Upload(context.Background(), rsc, nil)
//	_require.NoError(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.NoError(err)
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
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.NoError(err)
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
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.NoError(err)
//	// _require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp2.ContentMD5, "")
//	_require.NotEqual(stageResp2.RequestID, "")
//	_require.NotEqual(stageResp2.Version, "")
//	_require.Equal(stageResp2.Date.IsZero(), false)
//
//	// Check block list.
//	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//	_require.NoError(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.NotNil(blockList.BlockList.UncommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	// Commit block list.
//	commitBlockListOptions := blockblob.CommitBlockListOptions{
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
//	_require.NoError(err)
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
//	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//	_require.NoError(err)
//	// _require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.UncommittedBlocks)
//	_require.NotNil(blockList.BlockList.CommittedBlocks)
//	_require.Len(blockList.BlockList.CommittedBlocks, 2)
//
//	// Check data integrity through downloading.
//	_, err = destBlob.BlobClient.DownloadStream(context.Background(), nil)
//	_require.Error(err)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.BlobClient.DownloadStream(context.Background(), &downloadBlobOptions)
//	_require.NoError(err)
//	destData, err := io.ReadAll(downloadResp.Body)
//	_require.NoError(err)
//	_require.EqualValues(destData, content)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
// }

// func (s *BlockBlobUnrecordedTestsSuite) TestPutBlockFromURLAndCommitWithCPKWithScope() {
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
//	_, err = srcBlob.Upload(context.Background(), rsc, nil)
//	_require.NoError(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	// Get source blob url with SAS for StageFromURL.
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.NoError(err)
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
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	stageResp1, err := destBlob.StageBlockFromURL(context.Background(), blockID1, srcBlobURLWithSAS, 0, &options1)
//	_require.NoError(err)
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
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	stageResp2, err := destBlob.StageBlockFromURL(context.Background(), blockID2, srcBlobURLWithSAS, 0, &options2)
//	_require.NoError(err)
//	//_require.Equal(stageResp2.RawResponse.StatusCode, 201)
//	_require.NotEqual(stageResp2.ContentMD5, "")
//	_require.NotEqual(stageResp2.RequestID, "")
//	_require.NotEqual(stageResp2.Version, "")
//	_require.Equal(stageResp2.Date.IsZero(), false)
//
//	// Check block list.
//	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//	_require.NoError(err)
//	//_require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.CommittedBlocks)
//	_require.NotNil(blockList.BlockList.UncommittedBlocks)
//	_require.Len(blockList.BlockList.UncommittedBlocks, 2)
//
//	// Commit block list.
//	commitBlockListOptions := blockblob.CommitBlockListOptions{
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
//	_require.NoError(err)
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
//	blockList, err = destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
//	_require.NoError(err)
//	//_require.Equal(blockList.RawResponse.StatusCode, 200)
//	_require.NotNil(blockList.BlockList)
//	_require.Nil(blockList.BlockList.UncommittedBlocks)
//	_require.NotNil(blockList.BlockList.CommittedBlocks)
//	_require.Len(blockList.BlockList.CommittedBlocks, 2)
//
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := destBlob.BlobClient.DownloadStream(context.Background(), &downloadBlobOptions)
//	_require.NoError(err)
//	destData, err := io.ReadAll(downloadResp.Body)
//	_require.NoError(err)
//	_require.EqualValues(destData, content)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
// }

func (s *BlockBlobRecordedTestsSuite) TestUploadBlobWithMD5() {
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
		TransactionalValidation: blob.TransferValidationTypeMD5(md5Val[:]),
	}
	uploadResp, err := bbClient.Upload(context.Background(), r, &uploadBlockBlobOptions)
	_require.NoError(err)
	_require.Equal(uploadResp.ContentMD5, md5Val[:])

	// Download blob to do data integrity check.
	downloadResp, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.EqualValues(destData, srcData)

	// Test Upload with bad MD5
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	var badMD5Validator blob.TransferValidationTypeMD5 = badMD5

	uploadBlockBlobOptions = blockblob.UploadOptions{
		TransactionalValidation: badMD5Validator,
	}
	uploadResp, err = bbClient.Upload(context.Background(), r, &uploadBlockBlobOptions)
	_require.Error(err)
}

func (s *BlockBlobRecordedTestsSuite) TestUploadBlobWithMD5WithCPK() {
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
		CPKInfo: &testcommon.TestCPKByValue,
	}
	uploadResp, err := bbClient.Upload(context.Background(), r, &uploadBlockBlobOptions)
	_require.NoError(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.Equal(*uploadResp.IsServerEncrypted, true)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(uploadResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}

	// Get blob content without encryption key should fail the request.
	_, err = bbClient.DownloadStream(context.Background(), nil)
	_require.Error(err)

	_, err = bbClient.DownloadStream(context.Background(), &blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestInvalidCPKByValue,
	})
	_require.Error(err)

	// Download blob to do data integrity check.
	downloadResp, err := bbClient.DownloadStream(context.Background(), &blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	})
	_require.NoError(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.EqualValues(destData, srcData)
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(downloadResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
	}
}

func (s *BlockBlobRecordedTestsSuite) TestUploadBlobWithMD5WithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 8 * 1024
	r, srcData := testcommon.GenerateData(contentSize)
	md5Val := md5.Sum(srcData)
	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	uploadBlockBlobOptions := blockblob.UploadOptions{
		CPKScopeInfo: &encryptionScope,
	}
	uploadResp, err := bbClient.Upload(context.Background(), r, &uploadBlockBlobOptions)
	_require.NoError(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.Equal(*uploadResp.IsServerEncrypted, true)
	_require.EqualValues(*encryptionScope.EncryptionScope, *uploadResp.EncryptionScope)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKScopeInfo: &encryptionScope,
	}
	downloadResp, err := bbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)
	_require.EqualValues(downloadResp.ContentMD5, md5Val[:])
	destData, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionScope, *encryptionScope.EncryptionScope)
}

// func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlobBlobPropertiesWithCPKKey() {
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
//			BufferSize:  bufferSize,
//			MaxBuffers:  maxBuffers,
//			Metadata:    testcommon.BasicMetadata,
//			BlobTags:    basicBlobTagsMap,
//			HTTPHeaders: &basicHeaders,
//			CPKInfo:     &testcommon.TestCPKByValue,
//		})
//
//	// Assert that upload was successful
//	_require.NoError(err)
//	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
//
//	getPropertiesResp, err := bbClient.GetProperties(ctx, &blob.GetPropertiesOptions{CPKInfo: &testcommon.TestCPKByValue})
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
//	downloadResponse, err := bbClient.DownloadStream(ctx, &blob.downloadWriterAtOptions{CPKInfo: &testcommon.TestCPKByValue})
//	_require.NoError(err)
//
//	// Assert that the content is correct
//	actualBlobData, err := io.ReadAll(downloadResponse.Body(nil))
//	_require.NoError(err)
//	_require.Equal(len(actualBlobData), blobSize)
//	_require.EqualValues(actualBlobData, blobData)
// }

// func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlobBlobPropertiesWithCPKScope() {
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
//			CPKScopeInfo: &testcommon.TestCPKByScope,
//		})
//
//	// Assert that upload was successful
//	_require.NoError(err)
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
//	downloadResponse, err := bbClient.DownloadStream(ctx, &blob.downloadWriterAtOptions{CPKScopeInfo: &testcommon.TestCPKByScope})
//	_require.NoError(err)
//
//	// Assert that the content is correct
//	actualBlobData, err := io.ReadAll(downloadResponse.Body(nil))
//	_require.NoError(err)
//	_require.Equal(len(actualBlobData), blobSize)
//	_require.EqualValues(actualBlobData, blobData)
// }

func (s *BlockBlobUnrecordedTestsSuite) TestUploadStreamToBlobProperties() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)
	_require.NoError(err)
	// Create some data to test the upload stream
	blobContentReader, blobData := testcommon.GenerateData(blobSize)

	// Perform UploadStream
	_, err = bbClient.UploadStream(context.Background(), blobContentReader,
		&blockblob.UploadStreamOptions{
			BlockSize:   int64(bufferSize),
			Concurrency: maxBuffers,
			Metadata:    testcommon.BasicMetadata,
			Tags:        testcommon.BasicBlobTagsMap,
			HTTPHeaders: &testcommon.BasicHeaders,
		})

	// Assert that upload was successful
	_require.NoError(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)

	getPropertiesResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(getPropertiesResp.Metadata, testcommon.BasicMetadata)
	_require.Equal(*getPropertiesResp.TagCount, int64(len(testcommon.BasicBlobTagsMap)))
	respHeaders := testcommon.BasicHeaders
	calcMD5 := md5.Sum(blobData)
	respHeaders.BlobContentMD5 = calcMD5[:]
	_require.Equal(respHeaders, blob.ParseHTTPHeaders(getPropertiesResp))

	getTagsResp, err := bbClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	_require.Len(getTagsResp.BlobTagSet, 3)
	for _, blobTag := range getTagsResp.BlobTagSet {
		_require.Equal(testcommon.BasicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	// Download the blob to verify
	downloadResponse, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.Body)
	_require.NoError(err)
	_require.Equal(len(actualBlobData), blobSize)
	_require.EqualValues(actualBlobData, blobData)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlobUploadDownloadStream() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	blobSize := 11 * 1024 * 1024
	bufferSize := 2 * 1024 * 1024
	maxBuffers := 2

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.GetBlockBlobClient(blobName, containerClient)
	blobContentReader, blobData := testcommon.GenerateData(blobSize)

	_, err = bbClient.UploadStream(context.Background(), blobContentReader, &blockblob.UploadStreamOptions{
		BlockSize:   int64(bufferSize),
		Concurrency: maxBuffers,
		Metadata:    testcommon.BasicMetadata,
		Tags:        testcommon.BasicBlobTagsMap,
		HTTPHeaders: &testcommon.BasicHeaders,
	})
	_require.NoError(err)

	downloadResponse, err := bbClient.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	bbClient2 := testcommon.GetBlockBlobClient("blobName2", containerClient)

	// UploadStream using http.Response.Body as the reader
	_, err = bbClient2.UploadStream(context.Background(), downloadResponse.Body, &blockblob.UploadStreamOptions{
		BlockSize:   int64(bufferSize),
		Concurrency: maxBuffers,
	})
	_require.NoError(err)

	downloadResp2, err := bbClient2.DownloadStream(context.Background(), nil)
	_require.NoError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResp2.Body)
	_require.NoError(err)
	_require.Equal(len(actualBlobData), len(blobData))
	_require.EqualValues(actualBlobData, blobData)
}

// This test simulates UploadStream and DownloadBuffer methods,
// and verifies length and content of file
func (s *BlockBlobUnrecordedTestsSuite) TestBlobUploadStreamDownloadBuffer() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	const MiB = 1024 * 1024
	testUploadDownload := func(contentSize int) {
		content := make([]byte, contentSize)
		_, err := rand.Read(content)
		if err != nil {
			return
		}
		contentMD5 := md5.Sum(content)
		body := streaming.NopCloser(bytes.NewReader(content))

		srcBlob := containerClient.NewBlockBlobClient("srcblob")

		// Prepare source bbClient for copy.
		_, err = srcBlob.UploadStream(context.Background(), body, &blockblob.UploadStreamOptions{
			BlockSize:   4 * MiB,
			Concurrency: 5,
		})
		_require.NoError(err)

		// Download to a buffer and verify contents
		buff := make([]byte, contentSize)
		b := blob.DownloadBufferOptions{
			BlockSize:   5 * MiB,
			Concurrency: 4,
		}
		n, err := srcBlob.DownloadBuffer(context.Background(), buff, &b)
		_require.NoError(err)
		_require.Equal(int64(contentSize), n)
		_require.Equal(contentMD5, md5.Sum(buff))
	}

	testUploadDownload(0) // zero byte blob
	testUploadDownload(5 * MiB)
	testUploadDownload(20 * MiB)
	testUploadDownload(199 * MiB)
}

type fakeReader struct {
	cnt int
}

func (a *fakeReader) Read(bytes []byte) (count int, err error) {
	if a.cnt < 5 {
		_, buf := testcommon.GenerateData(1024)
		n := copy(bytes, buf)
		a.cnt++
		return n, nil
	}
	return 0, io.ErrUnexpectedEOF
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlobUploadStreamUsingCustomReader() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	r := &fakeReader{}
	_, err = bbClient.UploadStream(context.Background(), r, nil)
	_require.Error(err)
	_require.Equal(err, io.ErrUnexpectedEOF)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobSetTierOnVersions() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	archiveTier, rehydrateTier := blob.AccessTierArchive, blob.AccessTierCool
	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		uploadResp, err := bbClient.Upload(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("data"+strconv.Itoa(i)))), nil)
		_require.NoError(err)
		_require.NotNil(uploadResp.VersionID)
		versions = append(versions, *uploadResp.VersionID)
	}

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, b := range resp.Segment.BlobItems {
			_require.Equal(*b.Properties.AccessTier, blob.AccessTierHot)
		}
		if err != nil {
			break
		}
	}

	// set tier to archive for first three versions
	for i := 0; i < 3; i++ {
		bbClientWithVersionID, err := bbClient.WithVersionID(versions[i])
		_require.NoError(err)
		_, err = bbClientWithVersionID.SetTier(context.Background(), archiveTier, nil)
		_require.NoError(err)
	}

	// check access tier of versions
	for i, v := range versions {
		bbClientWithVersionID, err := bbClient.WithVersionID(v)
		_require.NoError(err)
		resp, err := bbClientWithVersionID.GetProperties(context.Background(), nil)
		_require.NoError(err)
		if i < 3 {
			_require.Equal(*resp.AccessTier, string(archiveTier))
		} else {
			_require.Equal(*resp.AccessTier, string(blob.AccessTierHot))
		}
	}

	// Versions tiered to archive cannot be rehydrated back to hot/cool tier
	// For detailed information refer this, https://learn.microsoft.com/en-us/rest/api/storageservices/set-blob-tier
	bbClientWithVersionID, err := bbClient.WithVersionID(versions[0])
	_require.NoError(err)
	_, err = bbClientWithVersionID.SetTier(context.Background(), rehydrateTier, nil)
	_require.Error(err)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobSetTierOnSnapshots() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	archiveTier, rehydrateTier := blob.AccessTierArchive, blob.AccessTierCool
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	snapshots := make([]string, 0)
	for i := 0; i < 5; i++ {
		resp, err := bbClient.CreateSnapshot(context.Background(), nil)
		_require.NoError(err)
		_require.NotNil(resp.Snapshot)
		snapshots = append(snapshots, *resp.Snapshot)
	}

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		for _, b := range resp.Segment.BlobItems {
			_require.Equal(*b.Properties.AccessTier, blob.AccessTierHot)
		}
		if err != nil {
			break
		}
	}

	// set tier to archive for first three snapshots
	for i := 0; i < 3; i++ {
		bbClientWithSnapshot, err := bbClient.WithSnapshot(snapshots[i])
		_require.NoError(err)
		_, err = bbClientWithSnapshot.SetTier(context.Background(), archiveTier, nil)
		_require.NoError(err)
	}

	// check access tier of snapshots
	for i, snap := range snapshots {
		bbClientWithSnapshot, err := bbClient.WithSnapshot(snap)
		_require.NoError(err)
		resp, err := bbClientWithSnapshot.GetProperties(context.Background(), nil)
		_require.NoError(err)
		if i < 3 {
			_require.Equal(*resp.AccessTier, string(archiveTier))
		} else {
			_require.Equal(*resp.AccessTier, string(blob.AccessTierHot))
		}
	}

	// Snapshots tiered to archive cannot be rehydrated back to hot/cool tier
	// For detailed information refer this, https://learn.microsoft.com/en-us/rest/api/storageservices/set-blob-tier
	bbClientWithSnapshot, err := bbClient.WithSnapshot(snapshots[0])
	_require.NoError(err)
	_, err = bbClientWithSnapshot.SetTier(context.Background(), rehydrateTier, nil)
	_require.Error(err)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobSetExpiryOnHnsDisabledAccount() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = bbClient.SetExpiry(context.Background(), nil, nil)
	testcommon.ValidateBlobErrorCode(_require, err, "HierarchicalNamespaceNotEnabled")
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobSetExpiryToNeverExpire() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = bbClient.SetExpiry(context.Background(), blockblob.ExpiryTypeNever{}, nil)
	_require.NoError(err)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobSetExpiryRelativeToNow() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = bbClient.SetExpiry(context.Background(), blockblob.ExpiryTypeRelativeToNow(8*time.Second), nil)
	_require.NoError(err)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ExpiresOn)

	time.Sleep(time.Second * 10)

	_, err = bbClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobSetExpiryRelativeToCreation() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = bbClient.SetExpiry(context.Background(), blockblob.ExpiryTypeRelativeToCreation(8*time.Second), nil)
	_require.NoError(err)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ExpiresOn)

	time.Sleep(time.Second * 10)

	_, err = bbClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlockBlobSetExpiryToAbsolute() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	expiryTimeAbsolute := time.Now().Add(8 * time.Second)
	_, err = bbClient.SetExpiry(context.Background(), blockblob.ExpiryTypeAbsolute(expiryTimeAbsolute), nil)
	_require.NoError(err)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ExpiresOn)
	_require.Equal(expiryTimeAbsolute.UTC().Format(http.TimeFormat), resp.ExpiresOn.UTC().Format(http.TimeFormat))

	time.Sleep(time.Second * 10)

	_, err = bbClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlockBlobSetExpiryToPast() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
	resp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	expiryTimeAbsolute := time.Now().Add(8 * time.Second)
	time.Sleep(time.Second * 10)
	_, err = bbClient.SetExpiry(context.Background(), blockblob.ExpiryTypeAbsolute(expiryTimeAbsolute), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)

	resp, err = bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)
}

func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockBlobStage() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	var largeBlockSize int64 = blockblob.MaxStageBlockBytes
	content := make([]byte, largeBlockSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	_, err = bbClient.StageBlock(context.Background(), blockID, rsc, nil)
	_require.NoError(err)

	_, err = bbClient.CommitBlockList(context.Background(), []string{blockID}, nil)
	_require.NoError(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.Len(resp.CommittedBlocks, 1)
	committed := resp.CommittedBlocks
	_require.Equal(*(committed[0].Name), blockID)
	_require.Equal(*(committed[0].Size), largeBlockSize)
	_require.Nil(resp.UncommittedBlocks)
}

func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockStreamUploadWithDifferentBlockSize() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	var firstBlockSize, secondBlockSize int64 = 2500 * 1024 * 1024, 10 * 1024 * 1024
	content := make([]byte, firstBlockSize+secondBlockSize)
	body := bytes.NewReader(content)
	rsc := streaming.NopCloser(body)

	_, err = bbClient.UploadStream(context.Background(), rsc, &blockblob.UploadStreamOptions{
		BlockSize:   firstBlockSize,
		Concurrency: 2,
	})
	_require.NoError(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.Len(resp.CommittedBlocks, 2)
	_require.Equal(*resp.BlobContentLength, firstBlockSize+secondBlockSize)
	committed := resp.CommittedBlocks
	_require.Equal(*(committed[0].Size), firstBlockSize)
	_require.Equal(*(committed[1].Size), secondBlockSize)
}

func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockBufferedUploadInParallel() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	var largeBlockSize, numberOfBlocks int64 = 2500 * 1024 * 1024, 2
	content := make([]byte, numberOfBlocks*largeBlockSize)

	_, err = bbClient.UploadBuffer(context.Background(), content, &blockblob.UploadBufferOptions{
		BlockSize:   largeBlockSize,
		Concurrency: 2,
	})
	_require.NoError(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.Len(resp.CommittedBlocks, 2)
	_require.Equal(*resp.BlobContentLength, numberOfBlocks*largeBlockSize)
	committed := resp.CommittedBlocks
	_require.Equal(*(committed[0].Size), largeBlockSize)
	_require.Equal(*(committed[1].Size), largeBlockSize)
}

func (s *BlockBlobUnrecordedTestsSuite) TestUploadFromReader() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := []int{398457897, 332398592, 19922944, 314572800, 269484032} // 380 MB, 317 MB, 19 MB, 300 MB, 257 MB

	for i, cs := range contentSize {
		bbClient := testcommon.GetBlockBlobClient(fmt.Sprintf("%v%v", testcommon.GenerateBlobName(testName), i), containerClient)

		_, content := testcommon.GetDataAndReader(testName, cs)
		md5Value := md5.Sum(content)
		contentMD5 := md5Value[:]

		_, err = bbClient.UploadBuffer(context.Background(), content, nil)
		_require.NoError(err)

		destBuffer := make([]byte, cs)
		cnt, err := bbClient.DownloadBuffer(context.Background(), destBuffer, nil)
		_require.NoError(err)
		_require.Equal(cnt, int64(cs))

		downloadedMD5Value := md5.Sum(destBuffer)
		downloadedContentMD5 := downloadedMD5Value[:]

		_require.EqualValues(downloadedContentMD5, contentMD5)

		gResp2, err := bbClient.GetProperties(context.Background(), nil)
		_require.NoError(err)
		_require.Equal(*gResp2.ContentLength, int64(cs))
	}
}

/* siminsavani: This test has a large allocation and blocks other tests from running that's why this test is commented out
func (s *BlockBlobUnrecordedTestsSuite) TestLargeBlockBufferedUploadInParallelWithGeneratedCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	var largeBlockSize, numberOfBlocks int64 = 2500 * 1024 * 1024, 2
	_, sourceData := testcommon.GetDataAndReader(testName, int(numberOfBlocks*largeBlockSize))
	// rsc := streaming.NopCloser(r)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	_, err = bbClient.UploadBuffer(context.Background(), sourceData, &blockblob.UploadBufferOptions{
		TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
		BlockSize:               largeBlockSize,
		Concurrency:             2,
	})
	_require.NoError(err)

	resp, err := bbClient.GetBlockList(context.Background(), blockblob.BlockListTypeAll, nil)
	_require.NoError(err)
	_require.Len(resp.CommittedBlocks, 2)
	_require.Equal(*resp.BlobContentLength, numberOfBlocks*largeBlockSize)
	committed := resp.CommittedBlocks
	_require.Equal(*(committed[0].Size), largeBlockSize)
	_require.Equal(*(committed[1].Size), largeBlockSize)
}*/

func (s *BlockBlobRecordedTestsSuite) TestUploadBufferWithCRC64OrMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := testcommon.GetBlockBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	_, content := testcommon.GetDataAndReader(testName, 1024)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	crc64Value := crc64.Checksum(content, shared.CRC64Table)

	_, err = bbClient.UploadBuffer(context.Background(), content, &blockblob.UploadBufferOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(crc64Value),
	})
	_require.Error(err)
	_require.Error(err, bloberror.UnsupportedChecksum)

	_, err = bbClient.UploadBuffer(context.Background(), content, &blockblob.UploadBufferOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(contentMD5),
	})
	_require.Error(err)
	_require.Error(err, bloberror.UnsupportedChecksum)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockGetAccountInfo() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blockBlobName := testcommon.GenerateBlobName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, blockBlobName, containerClient)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	bAccInfo, err := bbClient.GetAccountInfo(context.Background(), nil)
	_require.NoError(err)
	_require.NotZero(bAccInfo)
}

type fakeBlockBlob struct {
	totalStaged int64
}

func (f *fakeBlockBlob) Do(req *http.Request) (*http.Response, error) {
	// verify that the number of bytes read matches what was specified
	data := make([]byte, req.ContentLength)
	read, err := req.Body.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	} else if int64(read) < req.ContentLength {
		return nil, fmt.Errorf("expected %d bytes, read %d", req.ContentLength, read)
	}
	qp := req.URL.Query()
	if comp := qp.Get("comp"); comp == "block" {
		// staging a block, record its size
		f.totalStaged += int64(read)
	}
	return &http.Response{
		Request:    req,
		Status:     "Created",
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       http.NoBody,
	}, nil
}

func TestUploadBufferUnevenBlockSize(t *testing.T) {
	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/path", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's not evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 263*1024*1024)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(buffer)), fbb.totalStaged)
}

func TestUploadBufferEvenBlockSize(t *testing.T) {
	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/path", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 270000000)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(buffer)), fbb.totalStaged)
}

func TestUploadLogEvent(t *testing.T) {
	listnercalled := false
	log.SetEvents(azblob.EventUpload, log.EventRequest, log.EventResponse)
	log.SetListener(func(cls log.Event, msg string) {
		if cls == azblob.EventUpload {
			listnercalled = true
			require.Equal(t, msg, "blob name path1/path2 actual size 270000000 block-size 4194304 block-count 65")
		}
	})

	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/path1/path2", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: fbb,
			Telemetry: policy.TelemetryOptions{ApplicationID: "testApp/1.0.0-preview.2"},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 270000000)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
		Progress: func(bytesTransferred int64) {
			t.Logf("%v percent job done", (bytesTransferred*100)/270000000)
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(len(buffer)), fbb.totalStaged)
	require.Equal(t, listnercalled, true)
}

type trequestIDPolicy struct{}

// NewRequestIDPolicy returns a policy that add the x-ms-client-request-id header
func newTestRequestIDPolicy() policy.Policy {
	return &trequestIDPolicy{}
}

func (r *trequestIDPolicy) Do(req *policy.Request) (*http.Response, error) {
	const requestIdHeader = "x-ms-client-request-id"
	req.Raw().Header.Set(requestIdHeader, "azblob-test-request-id")
	return req.Next()
}

func TestRequestIDGeneration(t *testing.T) {
	requestIdMatch := false
	log.SetEvents(log.EventRequest)
	log.SetListener(func(cls log.Event, msg string) {
		require.Contains(t, msg, "X-Ms-Client-Request-Id: azblob-test-request-id")
		require.Contains(t, msg, "User-Agent: testApp/1.0.0-preview.2")
		requestIdMatch = true
	})

	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/testpath", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Telemetry:       policy.TelemetryOptions{ApplicationID: "testApp/1.0.0-preview.2"},
			Transport:       fbb,
			PerCallPolicies: []policy.Policy{newTestRequestIDPolicy()},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// create fake source that's evenly divisible by 50000 (max number of blocks)
	// and greater than MaxUploadBlobBytes (256MB) so that it doesn't fit into a single upload.

	buffer := make([]byte, 10)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 1
	}

	_, err = client.UploadBuffer(context.Background(), buffer, &blockblob.UploadBufferOptions{
		Concurrency: 1,
	})
	require.NoError(t, err)
	require.Equal(t, requestIdMatch, true)
}

type serviceVersionTest struct{}

// newServiceVersionTestPolicy returns a policy that checks the x-ms-version header
func newServiceVersionTestPolicy() policy.Policy {
	return &serviceVersionTest{}
}

func (m serviceVersionTest) Do(req *policy.Request) (*http.Response, error) {
	const versionHeader = "x-ms-version"

	currentVersion := map[string][]string(req.Raw().Header)[versionHeader]
	if currentVersion[0] != generated.ServiceVersion {
		return nil, fmt.Errorf("%s service version doesn't match expected version: %s", currentVersion[0], generated.ServiceVersion)

	}

	return &http.Response{
		Request:    req.Raw(),
		Status:     "Created",
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       http.NoBody,
	}, nil
}

func TestServiceVersion(t *testing.T) {
	fbb := &fakeBlockBlob{}
	client, err := blockblob.NewClientWithNoCredential("https://fake/blob/testpath", &blockblob.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport:       fbb,
			PerCallPolicies: []policy.Policy{newServiceVersionTestPolicy()},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)

	// Upload some data to source
	contentSize := 4 * 1024 // 4KB
	r, _ := testcommon.GetDataAndReader(t.Name(), contentSize)

	_, err = client.Upload(context.Background(), streaming.NopCloser(r), nil)
	require.NoError(t, err)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobClientDefaultAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	options := &blockblob.ClientOptions{
		Audience: "https://storage.azure.com/",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	bbClientAudience, err := blockblob.NewClient(blobURL, cred, options)
	_require.NoError(err)

	contentSize := 4 * 1024 // 4 KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)

	_, err = bbClientAudience.Upload(context.Background(), rsc, nil)
	_require.NoError(err)

	_, err = bbClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobRecordedTestsSuite) TestBlockBlobClientCustomAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)

	options := &blockblob.ClientOptions{
		Audience: "https://" + accountName + ".blob.core.windows.net",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	bbClientAudience, err := blockblob.NewClient(blobURL, cred, options)
	_require.NoError(err)

	contentSize := 4 * 1024 // 4 KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	rsc := streaming.NopCloser(r)

	_, err = bbClientAudience.Upload(context.Background(), rsc, nil)
	_require.NoError(err)

	_, err = bbClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *BlockBlobUnrecordedTestsSuite) TestBlockBlobClientUploadDownloadFile() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbClient := containerClient.NewBlockBlobClient(testcommon.GenerateBlobName(testName))

	// create local file
	var fileSize int64 = 401 * 1024 * 1024
	content := make([]byte, fileSize)
	_, err = rand.Read(content)
	_require.NoError(err)
	err = os.WriteFile("testFile", content, 0644)
	_require.NoError(err)

	defer func() {
		err = os.Remove("testFile")
		_require.NoError(err)
	}()

	fh, err := os.Open("testFile")
	_require.NoError(err)

	defer func(fh *os.File) {
		err := fh.Close()
		_require.NoError(err)
	}(fh)

	srcHash := md5.New()
	_, err = io.Copy(srcHash, fh)
	_require.NoError(err)
	contentMD5 := srcHash.Sum(nil)

	_, err = bbClient.UploadFile(context.Background(), fh, &blockblob.UploadFileOptions{
		Concurrency: 5,
		BlockSize:   4 * 1024 * 1024,
	})
	_require.NoError(err)

	// download to a temp file and verify contents
	tmp, err := os.CreateTemp("", "")
	_require.NoError(err)
	defer func() { _ = tmp.Close() }()

	n, err := bbClient.DownloadFile(context.Background(), tmp, &blob.DownloadFileOptions{BlockSize: 4 * 1024 * 1024})
	_require.NoError(err)
	_require.Equal(fileSize, n)

	stat, err := tmp.Stat()
	_require.NoError(err)
	_require.Equal(fileSize, stat.Size())

	destHash := md5.New()
	_, err = io.Copy(destHash, tmp)
	_require.NoError(err)
	downloadedContentMD5 := destHash.Sum(nil)

	_require.EqualValues(contentMD5, downloadedContentMD5)

	gResp, err := bbClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(gResp.ContentLength)
	_require.Equal(fileSize, *gResp.ContentLength)
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"hash/crc64"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running pageblob Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &PageBlobRecordedTestsSuite{})
		suite.Run(t, &PageBlobUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &PageBlobRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &PageBlobRecordedTestsSuite{})
	}
}

func (s *PageBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *PageBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *PageBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *PageBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type PageBlobRecordedTestsSuite struct {
	suite.Suite
}

type PageBlobUnrecordedTestsSuite struct {
	suite.Suite
}

func getPageBlobClient(pageBlobName string, containerClient *container.Client) *pageblob.Client {
	return containerClient.NewPageBlobClient(pageBlobName)
}

func createNewPageBlob(ctx context.Context, _require *require.Assertions, pageBlobName string, containerClient *container.Client) *pageblob.Client {
	return createNewPageBlobWithSize(ctx, _require, pageBlobName, containerClient, pageblob.PageBytes*10)
}

func createNewPageBlobWithSize(ctx context.Context, _require *require.Assertions, pageBlobName string, containerClient *container.Client, sizeInBytes int64) *pageblob.Client {
	pbClient := getPageBlobClient(pageBlobName, containerClient)

	_, err := pbClient.Create(ctx, sizeInBytes, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	return pbClient
}

func createNewPageBlobWithCPK(ctx context.Context, _require *require.Assertions, pageBlobName string, container *container.Client, sizeInBytes int64, cpkInfo *blob.CPKInfo, cpkScopeInfo *blob.CPKScopeInfo) (pbClient *pageblob.Client) {
	pbClient = getPageBlobClient(pageBlobName, container)

	_, err := pbClient.Create(ctx, sizeInBytes, &pageblob.CreateOptions{
		CPKInfo:      cpkInfo,
		CPKScopeInfo: cpkScopeInfo,
	})
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)
	return
}

func (s *PageBlobRecordedTestsSuite) TestPutGetPages() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	offset, count := int64(0), int64(1024)
	reader, _ := testcommon.GenerateData(1024)
	putResp, err := pbClient.UploadPages(context.Background(), reader, blob.HTTPRange{Count: count}, nil)
	_require.Nil(err)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.Nil(putResp.ContentMD5)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Range: blob.HTTPRange{
			Count: 1023,
		},
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.NotNil(pageListResp.LastModified)
		_require.Equal((*pageListResp.LastModified).IsZero(), false)
		_require.NotNil(pageListResp.ETag)
		_require.Equal(*pageListResp.BlobContentLength, int64(512*10))
		_require.NotNil(*pageListResp.RequestID)
		_require.NotNil(*pageListResp.Version)
		_require.NotNil(pageListResp.Date)
		_require.Equal((*pageListResp.Date).IsZero(), false)
		_require.NotNil(pageListResp.PageList)
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := rawPageRange((pageRangeResp)[0])
		_require.Equal(rawStart, offset)
		_require.Equal(rawEnd, count-1)
		if err != nil {
			break
		}
	}
}

// func (s *PageBlobUnrecordedTestsSuite) TestUploadPagesFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
//	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))
//
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadSrcResp1, err := srcBlob.UploadPages(context.Background(), streaming.NopCloser(r), &pageblob.UploadPagesOptions{
//		Offset: to.Ptr(offset),
//		Count: to.Ptr(count),
//	}
//	_require.Nil(err)
//	_require.NotNil(uploadSrcResp1.LastModified)
//	_require.Equal((*uploadSrcResp1.LastModified).IsZero(), false)
//	_require.NotNil(uploadSrcResp1.ETag)
//	_require.Nil(uploadSrcResp1.ContentMD5)
//	_require.Equal(*uploadSrcResp1.BlobSequenceNumber, int64(0))
//	_require.NotNil(*uploadSrcResp1.RequestID)
//	_require.NotNil(*uploadSrcResp1.Version)
//	_require.NotNil(uploadSrcResp1.Date)
//	_require.Equal((*uploadSrcResp1.Date).IsZero(), false)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
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
//		_require.Error(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL.
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
//	_require.Nil(err)
//	// _require.Equal(pResp1.RawResponse.StatusCode, 201)
//	_require.NotNil(pResp1.ETag)
//	_require.NotNil(pResp1.LastModified)
//	_require.NotNil(pResp1.ContentMD5)
//	_require.NotNil(pResp1.RequestID)
//	_require.NotNil(pResp1.Version)
//	_require.NotNil(pResp1.Date)
//	_require.Equal((*pResp1.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{}))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
// }
//

func (s *PageBlobUnrecordedTestsSuite) TestUploadPagesFromURLWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	md5Value := md5.Sum(sourceData)
	contentMD5 := md5Value[:]
	srcBlob := createNewPageBlobWithSize(context.Background(), _require, "srcblob"+testName, containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(context.Background(), _require, "dstblob"+testName, containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	_, err = srcBlob.UploadPages(context.Background(), streaming.NopCloser(r), blob.HTTPRange{Offset: offset, Count: count}, nil)
	_require.Nil(err)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())

	srcBlobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                      // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(15 * time.Minute), // 15 minutes before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   to.Ptr(sas.BlobPermissions{Read: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	srcBlobURLWithSAS := srcBlobParts.String()

	// Upload page from URL with MD5.
	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeMD5(contentMD5),
	}
	pResp1, err := destBlob.UploadPagesFromURL(context.Background(), srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.Nil(err)
	_require.EqualValues(pResp1.ContentMD5, contentMD5)

	// Download blob to do data integrity check.
	downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, sourceData)

	// Upload page from URL with bad MD5
	_, badMD5 := testcommon.GetDataAndReader(testName+"bad-md5", contentSize)
	badContentMD5 := badMD5[:]
	uploadPagesFromURLOptions = pageblob.UploadPagesFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeMD5(badContentMD5),
	}
	_, err = destBlob.UploadPagesFromURL(context.Background(), srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.NotNil(err)
	testcommon.ValidateHTTPErrorCode(_require, err, 400) // Fails with 400 (Bad Request)
}

func (s *PageBlobUnrecordedTestsSuite) TestUploadPagesFromURLWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)
	srcBlob := createNewPageBlobWithSize(context.Background(), _require, "srcblob"+testName, containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(context.Background(), _require, "dstblob"+testName, containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	_, err = srcBlob.UploadPages(context.Background(), streaming.NopCloser(r), blob.HTTPRange{Offset: offset, Count: count}, nil)
	_require.Nil(err)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())

	srcBlobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                      // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(15 * time.Minute), // 15 minutes before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   to.Ptr(sas.BlobPermissions{Read: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	srcBlobURLWithSAS := srcBlobParts.String()

	// Upload page from URL with CRC64.
	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeCRC64(crc),
	}
	_, err = destBlob.UploadPagesFromURL(context.Background(), srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.Nil(err)
	// TODO: This does not work... ContentCRC64 is not returned. Fix this later.
	// _require.EqualValues(pResp1.ContentCRC64, crc)

	// Download blob to do data integrity check.
	downloadResp, err := destBlob.DownloadStream(context.Background(), nil)
	_require.Nil(err)
	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, sourceData)
}

func (s *PageBlobUnrecordedTestsSuite) TestUploadPagesFromURLWithCRC64Negative() {
	s.T().Skip("This test is skipped because of issues in the service.")

	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)
	srcBlob := createNewPageBlobWithSize(context.Background(), _require, "srcblob"+testName, containerClient, int64(contentSize))
	destBlob := createNewPageBlobWithSize(context.Background(), _require, "dstblob"+testName, containerClient, int64(contentSize))

	// Prepare source pbClient for copy.
	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
	_, err = srcBlob.UploadPages(context.Background(), streaming.NopCloser(r), blob.HTTPRange{Offset: offset, Count: count}, nil)
	_require.Nil(err)

	// Get source pbClient URL with SAS for UploadPagesFromURL.
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	srcBlobParts, _ := blob.ParseURL(srcBlob.URL())

	srcBlobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                      // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(15 * time.Minute), // 15 minutes before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   to.Ptr(sas.BlobPermissions{Read: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	srcBlobURLWithSAS := srcBlobParts.String()

	// Upload page from URL with bad CRC64
	badCRC64 := rand.Uint64()
	badcrc := make([]byte, 8)
	binary.LittleEndian.PutUint64(badcrc, badCRC64)
	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeCRC64(badcrc),
	}
	_, err = destBlob.UploadPagesFromURL(context.Background(), srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
	_require.NotNil(err) // TODO: UploadPagesFromURL should fail, but is currently not working due to service issue.
}

func (s *PageBlobUnrecordedTestsSuite) TestClearDiffPages() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	contentSize := 2 * 1024
	r := testcommon.GetReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{Count: int64(contentSize)}, nil)
	_require.Nil(err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	r1 := testcommon.GetReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r1, blob.HTTPRange{Offset: int64(contentSize), Count: int64(contentSize)}, nil)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Range: blob.HTTPRange{
			Count: 4096,
		},
		PrevSnapshot: snapshotResp.Snapshot,
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		pageRangeResp := pageListResp.PageList.PageRange
		_require.NotNil(pageRangeResp)
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := rawPageRange((pageRangeResp)[0])
		_require.Equal(rawStart, int64(2048))
		_require.Equal(rawEnd, int64(4095))
		if err != nil {
			break
		}
	}

	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Offset: 2048, Count: 2048}, nil)
	_require.Nil(err)

	pager = pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Range: blob.HTTPRange{
			Count: 4096,
		},
		PrevSnapshot: snapshotResp.Snapshot,
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 0)
		if err != nil {
			break
		}
	}
}

func waitForIncrementalCopy(_require *require.Assertions, copyBlobClient *pageblob.Client, blobCopyResponse *pageblob.CopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult blob.GetPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != blob.CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(context.Background(), nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_require.Fail("")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

func (s *PageBlobUnrecordedTestsSuite) TestIncrementalCopy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err = containerClient.SetAccessPolicy(context.Background(), &container.SetAccessPolicyOptions{Access: to.Ptr(container.PublicAccessTypeBlob)})
	_require.Nil(err)

	srcBlob := createNewPageBlob(context.Background(), _require, "src"+testcommon.GenerateBlobName(testName), containerClient)

	contentSize := 1024
	r := testcommon.GetReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	_, err = srcBlob.UploadPages(context.Background(), r, blob.HTTPRange{
		Offset: offset,
		Count:  count,
	}, nil)
	_require.Nil(err)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	dstBlob := containerClient.NewPageBlobClient("dst" + testcommon.GenerateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	_require.Nil(err)
	_require.NotNil(resp.LastModified)
	_require.Equal((*resp.LastModified).IsZero(), false)
	_require.NotNil(resp.ETag)
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.NotEqual(*resp.CopyID, "")
	_require.Equal(*resp.CopyStatus, blob.CopyStatusTypePending)

	waitForIncrementalCopy(_require, dstBlob, &resp)
}

func (s *PageBlobRecordedTestsSuite) TestResizePageBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.Resize(context.Background(), 2048, nil)
	_require.Nil(err)

	_, err = pbClient.Resize(context.Background(), 8192, nil)
	_require.Nil(err)

	resp2, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(8192))
}

func (s *PageBlobRecordedTestsSuite) TestPageSequenceNumbers() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(0)
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	sequenceNumber = int64(7)
	actionType = pageblob.SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	updateSequenceNumberPageBlob = pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: to.Ptr(int64(11)),
		ActionType:     to.Ptr(pageblob.SequenceNumberActionTypeUpdate),
	}

	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)
}

func (s *PageBlobRecordedTestsSuite) TestPutPagesWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	// put page with valid auto-generated CRC64
	contentSize := 1024
	readerToBody, body := testcommon.GetDataAndReader(testName, contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	crc64Value := crc64.Checksum(body, shared.CRC64Table)
	_ = body

	putResp, err := pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), blob.HTTPRange{Offset: offset, Count: count}, &pageblob.UploadPagesOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(crc64Value),
	})
	_require.Nil(err)
	_require.NotNil(putResp.ContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(putResp.ContentCRC64), crc64Value)

	// put page with bad CRC64
	readerToBody, _ = testcommon.GetDataAndReader(testName, 1024)
	badCRC64 := rand.Uint64()
	putResp, err = pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), blob.HTTPRange{Offset: offset, Count: count}, &pageblob.UploadPagesOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(badCRC64),
	})
	_require.NotNil(err)

	// testcommon.ValidateBlobErrorCode(_require, err, bloberror.CRC64Mismatch)
}

// nolint
func (s *PageBlobRecordedTestsSuite) TestPutPagesWithAutoGeneratedCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	// put page with valid auto-generated CRC64
	contentSize := 1024
	readerToBody, body := testcommon.GetDataAndReader(testName, contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	crc64Value := crc64.Checksum(body, shared.CRC64Table)
	_ = body

	putResp, err := pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), blob.HTTPRange{Offset: offset, Count: count}, &pageblob.UploadPagesOptions{
		TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
	})
	_require.Nil(err)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotNil(putResp.ContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(putResp.ContentCRC64), crc64Value)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)
}

// nolint
func (s *PageBlobRecordedTestsSuite) TestPutPagesWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	// put page with valid MD5
	contentSize := 1024
	readerToBody, body := testcommon.GetDataAndReader(testName, contentSize)
	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]

	putResp, err := pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), blob.HTTPRange{Offset: offset, Count: count}, &pageblob.UploadPagesOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(contentMD5),
	})
	_require.Nil(err)
	// _require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotNil(putResp.ContentMD5)
	_require.EqualValues(putResp.ContentMD5, contentMD5)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	// put page with bad MD5
	readerToBody, _ = testcommon.GetDataAndReader(testName, 1024)
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	badContentMD5 := badMD5[:]
	putResp, err = pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), blob.HTTPRange{Offset: offset, Count: count}, &pageblob.UploadPagesOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(badContentMD5),
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageSizeInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(context.Background(), 1, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageSequenceInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       map[string]*string{},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       map[string]*string{"In valid1": to.Ptr("bar")},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)

}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		HTTPHeaders:    &testcommon.BasicHeaders,
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	h := blob.ParseHTTPHeaders(resp)
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func validatePageBlobPut(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
	_require.EqualValues(blob.ParseHTTPHeaders(resp), testcommon.BasicHeaders)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := azcore.ETag("garbage")
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := azcore.ETag("garbage")
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobPutPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	contentSize := 1024
	r := testcommon.GetReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{Count: int64(contentSize / 2)}, nil)
	_require.NotNil(err)
}

//// Body cannot be nil check already added in the request preparer
////func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesNilBody() {
////  svcClient := testcommon.GetServiceClient()
////  containerClient, _ := createNewContainer(c, svcClient)
////  defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////  pbClient, _ := createNewPageBlob(c, containerClient)
////
////  _, err := pbClient.UploadPages(context.Background(), nil, nil)
////  _require.NotNil(err)
////}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r := bytes.NewReader([]byte{})
	_, err = pbClient.UploadPages(context.Background(), streaming.NopCloser(r), blob.HTTPRange{Offset: 0, Count: 0}, nil)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{Count: pageblob.PageBytes}, nil)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func validateUploadPages(_require *require.Assertions, pbClient *pageblob.Client) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		start, end := int64(0), int64(pageblob.PageBytes-1)
		rawStart, rawEnd := *(pageListResp.PageList.PageRange[0].Start), *(pageListResp.PageList.PageRange[0].End)
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	eTag := azcore.ETag("garbage")
	uploadPagesOptions := pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	eTag := azcore.ETag("garbage")
	uploadPagesOptions := pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberLessThan := int64(10)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberLessThan := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}

	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

// func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blobName := testcommon.GenerateBlobName(testName)
//	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)
//
//	r, _ := testcommon.GenerateData(pageblob.PageBytes)
//	offset, count := int64(0), int64(pageblob.PageBytes)
//	ifSequenceNumberEqualTo := int64(-1)
//	uploadPagesOptions := pageblob.UploadPagesOptions{
//		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
//		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
//	_require.Nil(err)
// }

func setupClearPagesTest(t *testing.T, _require *require.Assertions, testName string) (*container.Client, *pageblob.Client) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, nil)
	_require.Nil(err)

	return containerClient, pbClient
}

func validateClearPagesTest(_require *require.Assertions, pbClient *pageblob.Client) {
	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{})
	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(pageListResp.PageRange)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes + 1}, nil)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		}})
	_require.Nil(err)
	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	eTag := azcore.ETag("garbage")
	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	eTag := azcore.ETag("garbage")
	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), blob.HTTPRange{Count: pageblob.PageBytes}, &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func setupGetPageRangesTest(t *testing.T, _require *require.Assertions, testName string) (containerClient *container.Client, pbClient *pageblob.Client) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient = testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient = createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, nil)
	_require.Nil(err)
	return
}

func validateBasicGetPageRanges(_require *require.Assertions, resp pageblob.PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	start, end := int64(0), int64(pageblob.PageBytes-1)
	rawStart, rawEnd := rawPageRange((resp.PageRange)[0])
	_require.Equal(rawStart, start)
	_require.Equal(rawEnd, end)
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesEmptyBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(resp.PageRange)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesEmptyRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Range: blob.HTTPRange{
			Offset: -2,
			Count:  500,
		},
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesNonContiguousRanges() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(2*pageblob.PageBytes), int64(pageblob.PageBytes)
	_, err := pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Offset: offset,
		Count:  count,
	}, nil)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		_require.NotNil(pageListResp)
		_require.Len(pageListResp, 2)

		start, end := int64(0), int64(pageblob.PageBytes-1)
		rawStart, rawEnd := rawPageRange(pageListResp[0])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)

		start, end = int64(pageblob.PageBytes*2), int64((pageblob.PageBytes*3)-1)
		rawStart, rawEnd = rawPageRange(pageListResp[1])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesNotPageAligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Range: blob.HTTPRange{
			Count: 2000,
		},
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)
	pager := snapshotURL.NewGetPageRangesPager(nil)
	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)

		validateBasicGetPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}})
	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfMatch: to.Ptr(azcore.ETag("garbage")),
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: to.Ptr(azcore.ETag("garbage")),
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}

	// serr := err.(StorageError)
	// _require.(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

func setupDiffPageRangesTest(t *testing.T, _require *require.Assertions, testName string) (containerClient *container.Client, pbClient *pageblob.Client, snapshot string) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient = testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient = createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r := testcommon.GetReaderToGeneratedBytes(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, nil)
	_require.Nil(err)

	resp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	snapshot = *resp.Snapshot

	r = testcommon.GetReaderToGeneratedBytes(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, nil)
	_require.Nil(err)
	return
}

func rawPageRange(pr *pageblob.PageRange) (start, end int64) {
	if pr.Start != nil {
		start = *pr.Start
	}
	if pr.End != nil {
		end = *pr.End
	}
	return
}

func validateDiffPageRanges(_require *require.Assertions, resp pageblob.PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	rawStart, rawEnd := rawPageRange(resp.PageRange[0])
	_require.EqualValues(rawStart, int64(0))
	_require.EqualValues(rawEnd, int64(pageblob.PageBytes-1))
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangesNonExistentSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	snapshotTime, _ := time.Parse(blob.SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		PrevSnapshot: to.Ptr(snapshotTime.Format(blob.SnapshotTimeFormat))})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.PreviousSnapshotNotFound)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Range: blob.HTTPRange{
			Count:  14,
			Offset: -22,
		},
		Snapshot: &snapshot,
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(-10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateDiffPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateDiffPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(-10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

// func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	resp, err := pbClient.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//
//	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
//		Snapshot: to.Ptr(snapshot),
//		LeaseAccessConditions: &blob.LeaseAccessConditions{
//			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	})
//	for pager.More() {
//		resp2, err := pager.NextPage(context.Background())
//		_require.Nil(err)
//		validateDiffPageRanges(_require, resp2.PageList, err)
//		if err != nil {
//			break
//		}
//	}
// }

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Snapshot: to.Ptr(snapshotStr),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr(azcore.ETag("garbage")),
			},
		}})

	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}

	}
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		PrevSnapshot: to.Ptr(snapshotStr),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: to.Ptr(azcore.ETag("garbage")),
			},
		}})

	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateDiffPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		PrevSnapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	})

	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err = pbClient.Resize(context.Background(), 0, nil)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(0))
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeInvalidSizeNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.Resize(context.Background(), -4, nil)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeInvalidSizeMisaligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.Resize(context.Background(), 12, nil)
	_require.NotNil(err)
}

func validateResize(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, _ := pbClient.GetProperties(context.Background(), nil)
	_require.Equal(*resp.ContentLength, int64(pageblob.PageBytes))
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	eTag := azcore.ETag("garbage")
	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	eTag := azcore.ETag("garbage")
	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberActionTypeInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := pageblob.SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberSequenceNumberInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	defer func() { // Invalid sequence number should panic
		_ = recover()
	}()

	sequenceNumber := int64(-1)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func validateSequenceNumberSet(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobSequenceNumber, int64(1))
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestPageSetImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.Nil(err)
	policy := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.Nil(err)

	setImmutabilityPolicyOptions := &blob.SetImmutabilityPolicyOptions{
		Mode:                     &policy,
		ModifiedAccessConditions: nil,
	}
	_, err = pbClient.SetImmutabilityPolicy(context.Background(), currentTime, setImmutabilityPolicyOptions)
	_require.Nil(err)

	_, err = pbClient.SetLegalHold(context.Background(), false, nil)
	_require.Nil(err)

	_, err = pbClient.Delete(context.Background(), nil)
	_require.NotNil(err)

	_, err = pbClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.Nil(err)

	_, err = pbClient.Delete(context.Background(), nil)
	_require.Nil(err)
}

func (s *PageBlobRecordedTestsSuite) TestPageDeleteImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.Nil(err)

	policy := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.Nil(err)

	setImmutabilityPolicyOptions := &blob.SetImmutabilityPolicyOptions{
		Mode:                     &policy,
		ModifiedAccessConditions: nil,
	}
	_, err = pbClient.SetImmutabilityPolicy(context.Background(), currentTime, setImmutabilityPolicyOptions)
	_require.Nil(err)

	_, err = pbClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.Nil(err)

	_, err = pbClient.Delete(context.Background(), nil)
	_require.Nil(err)
}

func (s *PageBlobRecordedTestsSuite) TestPageSetLegalHold() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	_, err = pbClient.SetLegalHold(context.Background(), true, nil)
	_require.Nil(err)

	// should fail since time has not passed yet
	_, err = pbClient.Delete(context.Background(), nil)
	_require.NotNil(err)

	_, err = pbClient.SetLegalHold(context.Background(), false, nil)
	_require.Nil(err)

	_, err = pbClient.Delete(context.Background(), nil)
	_require.Nil(err)

}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	eTag := azcore.ETag("garbage")
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, "src"+blobName, containerClient)

	eTag := azcore.ETag("garbage")
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, "src"+blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

// func setupStartIncrementalCopyTest(_require *require.Assertions, testName string) (containerClient *container.Client,
//	pbClient *pageblob.Client, copyPBClient *pageblob.Client, snapshot string) {
////	var recording *testframework.Recording
//	if _context != nil {
//		recording = _context.recording
//	}
//	svcClient, err := testcommon.GetServiceClient(recording, testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient = testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	accessType := container.PublicAccessTypeBlob
//	setAccessPolicyOptions := container.SetAccessPolicyOptions{
//		Access: &accessType,
//	}
//	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	_require.Nil(err)
//
//	pbClient = createNewPageBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
//	resp, _ := pbClient.CreateSnapshot(context.Background(), nil)
//
//	copyPBClient = getPageBlobClient("copy"+testcommon.GenerateBlobName(testName), containerClient)
//
//	// Must create the incremental copy pbClient so that the access conditions work on it
//	resp2, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), *resp.Snapshot, nil)
//	_require.Nil(err)
//	waitForIncrementalCopy(_require, copyPBClient, &resp2)
//
//	resp, _ = pbClient.CreateSnapshot(context.Background(), nil) // Take a new snapshot so the next copy will succeed
//	snapshot = *resp.Snapshot
//	return
// }

// func validateIncrementalCopy(_require *require.Assertions, copyPBClient *pageblob.Client, resp *pageblob.CopyIncrementalResponse) {
//	t := waitForIncrementalCopy(_require, copyPBClient, resp)
//
//	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
//	copySnapshotURL, err := copyPBClient.WithSnapshot(*t)
//	_require.Nil(err)
//	_, err = copySnapshotURL.GetProperties(context.Background(), nil)
//	_require.Nil(err)
// }

// func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopySnapshotNotExist() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blobName := testcommon.GenerateBlobName(testName)
//	pbClient := createNewPageBlob(context.Background(), _require, "src"+blobName, containerClient)
//	copyPBClient := getPageBlobClient("dst"+blobName, containerClient)
//
//	snapshot := time.Now().UTC().Format(blob.SnapshotTimeFormat)
//	_, err = copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, nil)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CannotVerifyCopySource)
// }

// func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
// }
//
// func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
// }
//
// func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
// }
//
// func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
// }
//

// func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	resp, _ := copyPBClient.GetProperties(context.Background(), nil)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfMatch: resp.ETag,
//		},
//	}
//	resp2, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp2)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
// }
//

// func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfMatchFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfMatch: &eTag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.TargetConditionNotMet)
// }

// func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfNoneMatch: &eTag,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
// }

// func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	resp, _ := copyPBClient.GetProperties(context.Background(), nil)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfNoneMatch: resp.ETag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
// }

func setAndCheckPageBlobTier(_require *require.Assertions, pbClient *pageblob.Client, tier blob.AccessTier) {
	_, err := pbClient.SetTier(context.Background(), tier, nil)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetTierAllTiersOnPageBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	premiumServiceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	premContainerName := "prem" + testcommon.GenerateContainerName(testName)
	premContainerClient := testcommon.CreateNewContainer(context.Background(), _require, premContainerName, premiumServiceClient)
	defer testcommon.DeleteContainer(context.Background(), _require, premContainerClient)

	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, pbName, premContainerClient)

	possibleTiers := []blob.AccessTier{
		blob.AccessTierP4,
		blob.AccessTierP6,
		blob.AccessTierP10,
		blob.AccessTierP20,
		blob.AccessTierP30,
		blob.AccessTierP40,
		blob.AccessTierP50,
	}
	for _, possibleTier := range possibleTiers {
		setAndCheckPageBlobTier(_require, pbClient, possibleTier)
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestPageBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := testcommon.GenerateData(contentSize)
	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(context.Background(), _require, pbName, containerClient, int64(contentSize), &testcommon.TestCPKByValue, nil)

	uploadPagesOptions := pageblob.UploadPagesOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	uploadResp, err := pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: int64(contentSize),
	}, &uploadPagesOptions)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.EqualValues(uploadResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)

	pager := pbClient.NewGetPageRangesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(contentSize-1)
		rawStart, rawEnd := rawPageRange(pageListResp[0])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

	// Get blob content without encryption key should fail the request.
	_, err = pbClient.DownloadStream(context.Background(), nil)
	_require.NotNil(err)

	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestInvalidCPKByValue,
	}
	_, err = pbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	downloadResp, err := pbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *PageBlobUnrecordedTestsSuite) TestPageBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := testcommon.GenerateData(contentSize)
	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(context.Background(), _require, pbName, containerClient, int64(contentSize), nil, &encryptionScope)

	uploadPagesOptions := pageblob.UploadPagesOptions{
		CPKScopeInfo: &encryptionScope,
	}
	uploadResp, err := pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: int64(contentSize),
	}, &uploadPagesOptions)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.EqualValues(*encryptionScope.EncryptionScope, *uploadResp.EncryptionScope)

	pager := pbClient.NewGetPageRangesPager(nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(contentSize-1)
		rawStart, rawEnd := rawPageRange(pageListResp[0])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKScopeInfo: &encryptionScope,
	}
	downloadResp, err := pbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionScope, *encryptionScope.EncryptionScope)
}

func (s *PageBlobUnrecordedTestsSuite) TestCreatePageBlobWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := createNewPageBlob(context.Background(), _require, "src"+testcommon.GenerateBlobName(testName), containerClient)

	putResp, err := pbClient.UploadPages(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), blob.HTTPRange{
		Count: 1024,
	}, nil)
	_require.Nil(err)
	// _require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotEqual(putResp.ETag, "")
	_require.NotEqual(putResp.Version, "")

	_, err = pbClient.SetTags(context.Background(), testcommon.BasicBlobTagsMap, nil)
	_require.Nil(err)
	time.Sleep(10 * time.Second)
	// _require.Equal(setTagResp.RawResponse.StatusCode, 204)

	gpResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
	_require.Equal(*gpResp.TagCount, int64(len(testcommon.BasicBlobTagsMap)))

	blobGetTagsResponse, err := pbClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.BasicBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.BasicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	modifiedBlobTags := map[string]string{
		"a0z1u2r3e4": "b0l1o2b3",
		"b0l1o2b3":   "s0d1k2",
	}

	_, err = pbClient.SetTags(context.Background(), modifiedBlobTags, nil)
	_require.Nil(err)
	// _require.Equal(setTagResp.RawResponse.StatusCode, 204)

	gpResp, err = pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
	_require.Equal(*gpResp.TagCount, int64(len(modifiedBlobTags)))

	blobGetTagsResponse, err = pbClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet = blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(modifiedBlobTags))
	for _, blobTag := range blobTagsSet {
		_require.Equal(modifiedBlobTags[*blobTag.Key], *blobTag.Value)
	}

	// Test FilterBlobs API
	where := "\"azure\"='blob'"
	lResp, err := svcClient.FilterBlobs(context.Background(), where, nil)
	_require.Nil(err)
	_require.Equal(*lResp.FilterBlobSegment.Blobs[0].Tags.BlobTagSet[0].Key, "azure")
	_require.Equal(*lResp.FilterBlobSegment.Blobs[0].Tags.BlobTagSet[0].Value, "blob")
}

func (s *PageBlobUnrecordedTestsSuite) TestPageBlobSetBlobTagForSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := createNewPageBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	_, err = pbClient.SetTags(context.Background(), testcommon.SpecialCharBlobTagsMap, nil)
	_require.Nil(err)
	time.Sleep(10 * time.Second)

	resp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp2.TagCount, int64(len(testcommon.SpecialCharBlobTagsMap)))

	blobGetTagsResponse, err := pbClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.SpecialCharBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.SpecialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	// Tags with spaces
	where := "\"GO \"='.Net'"
	lResp, err := svcClient.FilterBlobs(context.Background(), where, nil)
	_require.Nil(err)
	_require.Equal(*lResp.FilterBlobSegment.Blobs[0].Tags.BlobTagSet[0].Key, "GO ")
	_require.Equal(*lResp.FilterBlobSegment.Blobs[0].Tags.BlobTagSet[0].Value, ".Net")
}

func (s *PageBlobRecordedTestsSuite) TestCreatePageBlobReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClob := createNewPageBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	const contentSize = 1 * 1024
	r, _ := testcommon.GenerateData(contentSize)
	putResp, err := pbClob.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: contentSize,
	}, nil)
	_require.Nil(err)
	// _require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotEqual(*putResp.Version, "")

	gpResp, err := pbClob.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(context.Background(), _require, pbName, containerClient, pageblob.PageBytes*10, &testcommon.TestCPKByValue, nil)

	resizePageBlobOptions := pageblob.ResizeOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	resp, _ := pbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Equal(*resp.ContentLength, int64(pageblob.PageBytes))
}

func (s *PageBlobRecordedTestsSuite) TestPageBlockPermanentDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.Nil(err)

	// Create container and blob, upload blob to container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	count := int64(1024)
	reader, _ := testcommon.GenerateData(1024)
	_, err = pbClient.UploadPages(context.Background(), reader, blob.HTTPRange{
		Count: count,
	}, nil)
	_require.Nil(err)

	parts, err := sas.ParseURL(pbClient.URL()) // Get parts for BlobURL
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	// Set Account SAS and set Permanent Delete to true
	parts.SAS, err = sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true, PermanentDelete: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	// Create snapshot of Blob and get snapshot URL
	resp, err := pbClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})
	_require.Nil(err)
	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)

	// Check that there are two items in the container: one snapshot, one blob
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Include: container.ListBlobsInclude{Snapshots: true}})
	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	// Delete snapshot (snapshot will be soft deleted)
	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = pbClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.Nil(err)

	// Check that only blob exists (snapshot is soft-deleted)
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)

	// Check that soft-deleted snapshot exists by including deleted items
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Deleted: true},
	})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	// Options for PermanentDeleteOptions
	perm := blob.DeleteTypePermanent
	deleteBlobOptions := blob.DeleteOptions{
		BlobDeleteType: &perm,
	}
	time.Sleep(time.Second * 30)

	// Execute Delete with DeleteTypePermanent
	pdResp, err := snapshotURL.Delete(context.Background(), &deleteBlobOptions)
	_require.Nil(err)
	_require.NotNil(pdResp)
	_require.NotNil(pdResp)

	// Check that only blob exists even after including snapshots and deleted items
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Deleted: true}})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)
}

func (s *PageBlobRecordedTestsSuite) TestPageBlockPermanentDeleteWithoutPermission() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.Nil(err)

	// Create container and blob, upload blob to container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	count := int64(1024)
	reader, _ := testcommon.GenerateData(1024)
	_, err = pbClient.UploadPages(context.Background(), reader, blob.HTTPRange{
		Count: count,
	}, nil)
	_require.Nil(err)

	parts, err := sas.ParseURL(pbClient.URL()) // Get parts for BlobURL
	_require.Nil(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.Nil(err)

	// Set Account SAS
	parts.SAS, err = sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.Nil(err)

	// Create snapshot of Blob and get snapshot URL
	resp, err := pbClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})
	_require.Nil(err)
	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)

	// Check that there are two items in the container: one snapshot, one blob
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Include: container.ListBlobsInclude{Snapshots: true}})
	found := make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 2)

	// Delete snapshot
	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = pbClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.Nil(err)

	// Check that only blob exists
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)

	// Options for PermanentDeleteOptions
	perm := blob.DeleteTypePermanent
	deleteBlobOptions := blob.DeleteOptions{
		BlobDeleteType: &perm,
	}
	// Execute Delete with DeleteTypePermanent,should fail because permissions are not set and snapshot is not soft-deleted
	_, err = snapshotURL.Delete(context.Background(), &deleteBlobOptions)
	_require.NotNil(err)
}

// func (s *AZBlobUnrecordedTestsSuite) TestPageBlockFromURLWithCPK() {
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
//		CPKInfo:          &testcommon.TestCPKByValue,
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
//		CPKInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CPKInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
// }

// func (s *AZBlobUnrecordedTestsSuite) TestPageBlockFromURLWithCPKScope() {
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
//		CPKScopeInfo:     &testcommon.TestCPKByScope,
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
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := dstPBBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CPKInfo: &testcommon.TestCPKByValue}))
//	_require.Nil(err)
//	_require.EqualValues(destData, srcData)
// }

// func (s *AZBlobUnrecordedTestsSuite) TestUploadPagesFromURLWithMD5WithCPK() {
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
//		CPKInfo:          &testcommon.TestCPKByValue,
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
//		CPKInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destPBClient.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NotNil(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destPBClient.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Nil(err)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CPKInfo: &testcommon.TestCPKByValue}))
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
// }

// func (s *AZBlobRecordedTestsSuite) TestClearDiffPagesWithCPK() {
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
//	uploadPagesOptions := pageblob.UploadPagesOptions{Range: &HttpRange{offset, count}, CPKInfo: &testcommon.TestCPKByValue}
//	_, err = pbClient.UploadPages(ctx, r, &uploadPagesOptions)
//	_require.Nil(err)
//
//	createBlobSnapshotOptions := blob.CreateSnapshotOptions{
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	snapshotResp, err := pbClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
//	_require.Nil(err)
//
//	offset1, end1, count1 := int64(contentSize), int64(2*contentSize-1), int64(contentSize)
//	uploadPagesOptions1 := pageblob.UploadPagesOptions{Range: &HttpRange{offset1, count1}, CPKInfo: &testcommon.TestCPKByValue}
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
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	clearResp, err := pbClient.ClearPages(ctx, HttpRange{2048, 2048}, &clearPagesOptions)
//	_require.Nil(err)
//	_require.Equal(clearResp.RawResponse.StatusCode, 201)
//
//	pageListResp, err = pbClient.NewGetPageRangesDiffPager(ctx, HttpRange{0, 4095}, *snapshotResp.Snapshot, nil)
//	_require.Nil(err)
//	_require.Nil(pageListResp.PageList.Range)
// }

func (s *PageBlobRecordedTestsSuite) TestUndeletePageBlobVersion() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := getPageBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, nil)
	_require.Nil(err)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		resp, err := pbClient.CreateSnapshot(context.Background(), nil)
		_require.Nil(err)
		_require.NotNil(resp.VersionID)
		versions = append(versions, *resp.VersionID)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6)

	// Deleting the 1st, 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		pbClientWithVersionID, err := pbClient.WithVersionID(versions[i])
		_require.Nil(err)
		_, err = pbClientWithVersionID.Delete(context.Background(), nil)
		_require.Nil(err)
	}

	// adding wait after delete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 3)

	_, err = pbClient.Undelete(context.Background(), nil)
	_require.Nil(err)

	// adding wait after undelete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6)
}

func (s *PageBlobRecordedTestsSuite) TestUndeletePageBlobSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := getPageBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, blob.HTTPRange{
		Count: pageblob.PageBytes,
	}, nil)
	_require.Nil(err)

	snapshots := make([]string, 0)
	for i := 0; i < 5; i++ {
		resp, err := pbClient.CreateSnapshot(context.Background(), nil)
		_require.Nil(err)
		_require.NotNil(resp.Snapshot)
		snapshots = append(snapshots, *resp.Snapshot)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6) // 5 snapshots and 1 current version

	// Deleting the 1st, 2nd and 3rd snapshots
	for i := 0; i < 3; i++ {
		pbClientWithSnapshot, err := pbClient.WithSnapshot(snapshots[i])
		_require.Nil(err)
		_, err = pbClientWithSnapshot.Delete(context.Background(), nil)
		_require.Nil(err)
	}

	// adding wait after delete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 3) // 2 snapshots and 1 current version

	_, err = pbClient.Undelete(context.Background(), nil)
	_require.Nil(err)

	// adding wait after undelete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6) // 5 snapshots and 1 current version
}

func (s *PageBlobRecordedTestsSuite) TestPageGetAccountInfo() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := getPageBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	bAccInfo, err := pbClient.GetAccountInfo(context.Background(), nil)
	_require.Nil(err)
	_require.NotZero(bAccInfo)
}

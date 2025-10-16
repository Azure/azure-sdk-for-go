//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package appendblob_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	afservice "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/service"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/appendblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running appendblob Tests in %s mode\n", recordMode)
	if recordMode == recording.LiveMode {
		suite.Run(t, &AppendBlobRecordedTestsSuite{})
		suite.Run(t, &AppendBlobUnrecordedTestsSuite{})
	} else if recordMode == recording.PlaybackMode {
		suite.Run(t, &AppendBlobRecordedTestsSuite{})
	} else if recordMode == recording.RecordingMode {
		suite.Run(t, &AppendBlobRecordedTestsSuite{})
	}
}

func (s *AppendBlobRecordedTestsSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *AppendBlobRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *AppendBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *AppendBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *AppendBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *AppendBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type AppendBlobRecordedTestsSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type AppendBlobUnrecordedTestsSuite struct {
	suite.Suite
}

func getAppendBlobClient(appendBlobName string, containerClient *container.Client) *appendblob.Client {
	return containerClient.NewAppendBlobClient(appendBlobName)
}

func createNewAppendBlob(ctx context.Context, _require *require.Assertions, appendBlobName string, containerClient *container.Client) *appendblob.Client {
	abClient := getAppendBlobClient(appendBlobName, containerClient)

	_, err := abClient.Create(ctx, nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	return abClient
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlock() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	appendResp, err := abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NoError(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal(appendResp.LastModified.IsZero(), false)
	_require.Nil(appendResp.ContentMD5)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal(appendResp.Date.IsZero(), false)

	appendResp, err = abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NoError(err)
	_require.Equal(*appendResp.BlobAppendOffset, "1024")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(2))
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlobClient() {
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

	abClient, err := appendblob.NewClient(blobURL, cred, nil)
	_require.NoError(err)

	resp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlobClientSharedKey() {
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

	abClient, err := appendblob.NewClientWithSharedKeyCredential(blobURL, cred, nil)
	_require.NoError(err)

	resp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlobClientConnectionString() {
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

	abClient, err := appendblob.NewClientFromConnectionString(*connectionString, containerName, blobName, nil)
	_require.NoError(err)

	resp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockHighThroughput() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(abName))

	// Create AppendBlob with 5MB data
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	contentSize := 5 * 1024 * 1024 // 5MB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	appendResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)
	_require.NotNil(appendResp.ETag)

	// Check data integrity through downloading.
	destBuffer := make([]byte, contentSize)
	downloadBufferOptions := blob.DownloadBufferOptions{Range: blob.HTTPRange{Offset: 0, Count: int64(contentSize)}}
	_, err = abClient.DownloadBuffer(context.Background(), destBuffer, &downloadBufferOptions)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockWithAutoGeneratedCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	// test append block with valid CRC64 value
	readerToBody, body := testcommon.GetDataAndReader(testName, 1024)
	contentCRC64 := crc64.Checksum(body, shared.CRC64Table)
	appendBlockOptions := appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
	}

	appendResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.NoError(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal(appendResp.LastModified.IsZero(), false)
	_require.NotNil(appendResp.ContentCRC64)
	_require.EqualValues(binary.LittleEndian.Uint64(appendResp.ContentCRC64), contentCRC64)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal(appendResp.Date.IsZero(), false)

	// test bad CRC64
	readerToBody, body = testcommon.GetDataAndReader(testName, 1024)
	badCRC64 := rand.Uint64()
	_ = body
	appendBlockOptions = appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(badCRC64),
	}

	appendResp, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CRC64Mismatch)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	// test append block with valid MD5 value
	readerToBody, body := testcommon.GetDataAndReader(testName, 1024)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]
	appendBlockOptions := appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(contentMD5),
	}
	appendResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.NoError(err)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal(appendResp.LastModified.IsZero(), false)
	_require.EqualValues(appendResp.ContentMD5, contentMD5)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal(appendResp.Date.IsZero(), false)

	// test append block with bad MD5 value
	readerToBody, body = testcommon.GetDataAndReader(testName, 1024)
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	_ = body
	appendBlockOptions = appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeMD5(badMD5),
	}
	appendResp, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	// test append block with valid CRC64 value
	readerToBody, body := testcommon.GetDataAndReader(testName, 1024)
	crc64Value := crc64.Checksum(body, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)
	appendBlockOptions := appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(crc64Value),
	}
	appendResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.NoError(err)
	_require.EqualValues(appendResp.ContentCRC64, crc)

	// test append block with bad CRC64 value
	readerToBody, body = testcommon.GetDataAndReader(testName, 1024)
	badCRC64 := rand.Uint64()
	_ = body
	appendBlockOptions = appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeCRC64(badCRC64),
	}
	appendResp, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.Error(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithSDKGeneratedCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	// test append block with SDK generated CRC64 value
	readerToBody, body := testcommon.GetDataAndReader(testName, 1024)
	crc64Value := crc64.Checksum(body, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)
	appendBlockOptions := appendblob.AppendBlockOptions{
		TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
	}
	appendResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(readerToBody), &appendBlockOptions)
	_require.NoError(err)
	_require.EqualValues(appendResp.ContentCRC64, crc)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURL() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	contentMD5 := md5.Sum(sourceData)
	srcBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))

	// Prepare source abClient for copy.
	_, err = srcBlob.Create(context.Background(), nil)
	_require.NoError(err)

	appendResp, err := srcBlob.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)

	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal(appendResp.LastModified.IsZero(), false)
	_require.Nil(appendResp.ContentMD5)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal(appendResp.Date.IsZero(), false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
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

	// Append block from URL.
	_, err = destBlob.Create(context.Background(), nil)
	_require.NoError(err)

	appendFromURLResp, err := destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendblob.AppendBlockFromURLOptions{
		SourceContentValidation: blob.SourceContentValidationTypeMD5(contentMD5[:]),
	})

	_require.NoError(err)
	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendFromURLResp.ETag)
	_require.NotNil(appendFromURLResp.LastModified)
	_require.Equal(appendFromURLResp.LastModified.IsZero(), false)
	_require.NotNil(appendFromURLResp.ContentMD5)
	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5[:])
	_require.NotNil(appendFromURLResp.RequestID)
	_require.NotNil(appendFromURLResp.Version)
	_require.NotNil(appendFromURLResp.Date)
	_require.Equal(appendFromURLResp.Date.IsZero(), false)

	// Check data integrity through downloading.
	destBuffer := make([]byte, 4*1024)
	downloadBufferOptions := blob.DownloadBufferOptions{Range: blob.HTTPRange{Offset: 0, Count: 4096}}
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, &downloadBufferOptions)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithRequestIntentHeader() {
	_require := require.New(s.T())
	ctx := context.Background()

	blobSvcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSecondary, nil)
	_require.NoError(err)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)
	accessToken, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	_require.NoError(err)

	containerName := "testcontainer"
	containerClient := testcommon.CreateNewContainer(ctx, _require, containerName, blobSvcClient)
	defer testcommon.DeleteContainer(ctx, _require, containerClient)

	destBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))
	_, err = destBlob.Create(ctx, nil)
	_require.NoError(err)

	contentSize := 1024
	_, sourceData := testcommon.GenerateData(contentSize)
	accountName, accountKey := testcommon.GetGenericAccountInfo(testcommon.TestAccountSecondary)

	fileSvcURL := "https://" + accountName + ".file.core.windows.net/"
	sharedKeyCred, err := afservice.NewSharedKeyCredential(accountName, accountKey)
	_require.NoError(err)
	fileSvcClient, err := afservice.NewClientWithSharedKeyCredential(fileSvcURL, sharedKeyCred, nil)
	_require.NoError(err)

	shareName := "testshare"
	shareClient := fileSvcClient.NewShareClient(shareName)
	_, err = shareClient.Create(ctx, nil)
	_require.NoError(err)
	defer func() { _, _ = shareClient.Delete(ctx, nil) }()

	dirClient := shareClient.NewRootDirectoryClient()
	fileName := "testfile"
	srcFile := dirClient.NewFileClient(fileName)

	fileSize := int64(contentSize)
	_, err = srcFile.Create(ctx, fileSize, nil)
	_require.NoError(err)
	_, err = srcFile.UploadRange(ctx, 0, streaming.NopCloser(bytes.NewReader(sourceData)), nil)
	_require.NoError(err)

	// Append block from file source via OAuth with file request intent header
	requestIntent := blob.FileRequestIntentTypeBackup
	_, err = destBlob.AppendBlockFromURL(ctx, srcFile.URL(), &appendblob.AppendBlockFromURLOptions{
		CopySourceAuthorization: to.Ptr("Bearer " + accessToken.Token),
		FileRequestIntent:       &requestIntent,
	})
	_require.NoError(err)
}

func (s *AppendBlobUnrecordedTestsSuite) TestBlobEncryptionScopeSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))

	// Get source abClient URL with SAS for AppendBlockFromURL.
	blobParts, _ := blob.ParseURL(blobClient.URL())

	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.EncryptionScopeEnvVar)
	_require.NoError(err)
	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)
	perms := sas.BlobPermissions{Read: true, Create: true, Write: true, Delete: true}

	blobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:        sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:      time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName:   blobParts.ContainerName,
		BlobName:        blobParts.BlobName,
		Permissions:     perms.String(),
		EncryptionScope: encryptionScope,
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	blobURLWithSAS := blobParts.String()

	// create new client with sas url
	blobClient, err = appendblob.NewClientWithNoCredential(blobURLWithSAS, nil)
	_require.NoError(err)

	createResponse, err := blobClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*createResponse.EncryptionScope, encryptionScope)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAccountEncryptionScopeSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName("appendsrc")
	blobClient := containerClient.NewAppendBlobClient(blobName)

	// Get blob URL with SAS for AppendBlockFromURL.
	blobParts, _ := blob.ParseURL(blobClient.URL())

	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.EncryptionScopeEnvVar)
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	blobParts.SAS, err = sas.AccountSignatureValues{
		Protocol:        sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:      time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:     to.Ptr(sas.AccountPermissions{Read: true, Create: true, Write: true, Delete: true}).String(),
		ResourceTypes:   to.Ptr(sas.AccountResourceTypes{Service: true, Container: true, Object: true}).String(),
		EncryptionScope: encryptionScope,
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	blobURLWithSAS := blobParts.String()
	blobClient, err = appendblob.NewClientWithNoCredential(blobURLWithSAS, nil)
	_require.NoError(err)

	createResp, err := blobClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(createResp)
	_require.Equal(*createResp.EncryptionScope, encryptionScope)
}

func (s *AppendBlobUnrecordedTestsSuite) TestGetUserDelegationEncryptionScopeSAS() {
	_require := require.New(s.T())
	testName := s.T().Name()
	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	svcClient, err := service.NewClient("https://"+accountName+".blob.core.windows.net/", cred, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	cntClientTokenCred := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, cntClientTokenCred)

	blobName := testcommon.GenerateBlobName("appendsrc")
	blobClient := cntClientTokenCred.NewAppendBlobClient(blobName)

	// Set current and past time and create key
	now := time.Now().UTC().Add(-10 * time.Second)
	expiry := now.Add(2 * time.Hour)
	info := service.KeyInfo{
		Start:  to.Ptr(now.UTC().Format(sas.TimeFormat)),
		Expiry: to.Ptr(expiry.UTC().Format(sas.TimeFormat)),
	}

	udc, err := svcClient.GetUserDelegationCredential(context.Background(), info, nil)
	_require.NoError(err)

	// get permissions and details for sas
	encryptionScope, err := testcommon.GetRequiredEnv(testcommon.EncryptionScopeEnvVar)
	_require.NoError(err)

	permissions := sas.BlobPermissions{Read: true, Create: true, Write: true, List: true, Add: true, Delete: true}

	blobParts, _ := blob.ParseURL(blobClient.URL())

	// Create Blob Signature Values with desired permissions and sign with user delegation credential
	blobParts.SAS, err = sas.BlobSignatureValues{
		Protocol:        sas.ProtocolHTTPS,
		StartTime:       time.Now().UTC().Add(time.Second * -10),
		ExpiryTime:      time.Now().UTC().Add(15 * time.Minute),
		Permissions:     permissions.String(),
		ContainerName:   containerName,
		EncryptionScope: encryptionScope,
	}.SignWithUserDelegation(udc)
	_require.NoError(err)

	blobURLWithSAS := blobParts.String()
	blobClient, err = appendblob.NewClientWithNoCredential(blobURLWithSAS, nil)
	_require.NoError(err)

	createResp, err := blobClient.Create(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(createResp)
	_require.Equal(*createResp.EncryptionScope, encryptionScope)

}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithMD5() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	contentMD5 := md5.Sum(sourceData)
	srcBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))

	// Prepare source abClient for copy.
	_, err = srcBlob.Create(context.Background(), nil)
	_require.NoError(err)

	appendResp, err := srcBlob.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)

	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendResp.ETag)
	_require.NotNil(appendResp.LastModified)
	_require.Equal(appendResp.LastModified.IsZero(), false)
	_require.Nil(appendResp.ContentMD5)
	_require.NotNil(appendResp.RequestID)
	_require.NotNil(appendResp.Version)
	_require.NotNil(appendResp.Date)
	_require.Equal(appendResp.Date.IsZero(), false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
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

	// Append block from URL.
	_, err = destBlob.Create(context.Background(), nil)
	_require.NoError(err)

	count := int64(contentSize)
	var md5Validator blob.SourceContentValidationTypeMD5 = contentMD5[:]
	appendBlockURLOptions := appendblob.AppendBlockFromURLOptions{
		Range:                   blob.HTTPRange{Offset: 0, Count: count},
		SourceContentValidation: blob.SourceContentValidationType(md5Validator),
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)

	_require.NoError(err)
	_require.Equal(*appendFromURLResp.BlobAppendOffset, "0")
	_require.Equal(*appendFromURLResp.BlobCommittedBlockCount, int32(1))
	_require.NotNil(appendFromURLResp.ETag)
	_require.NotNil(appendFromURLResp.LastModified)
	_require.Equal(appendFromURLResp.LastModified.IsZero(), false)
	_require.NotNil(appendFromURLResp.ContentMD5)
	_require.EqualValues(appendFromURLResp.ContentMD5, contentMD5[:])
	_require.NotNil(appendFromURLResp.RequestID)
	_require.NotNil(appendFromURLResp.Version)
	_require.NotNil(appendFromURLResp.Date)
	_require.Equal(appendFromURLResp.Date.IsZero(), false)

	// Check data integrity through downloading.
	destBuffer := make([]byte, 4*1024)
	downloadBufferOptions := blob.DownloadBufferOptions{Range: blob.HTTPRange{Offset: 0, Count: 4096}}
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, &downloadBufferOptions)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)

	// Test append block from URL with bad MD5 value
	_, badMD5 := testcommon.GetDataAndReader(testName, 16)
	var badMD5Validator blob.SourceContentValidationTypeMD5 = badMD5
	appendBlockURLOptions = appendblob.AppendBlockFromURLOptions{
		Range:                   blob.HTTPRange{Offset: 0, Count: count},
		SourceContentValidation: badMD5Validator,
	}
	_, err = destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCRC64() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)
	srcBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))

	// Prepare source abClient for copy.
	_, err = srcBlob.Create(context.Background(), nil)
	_require.NoError(err)

	appendResp, err := srcBlob.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)
	_require.EqualValues(appendResp.ContentCRC64, crc)

	// Get source abClient URL with SAS for AppendBlockFromURL.
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

	// Append block from URL.
	_, err = destBlob.Create(context.Background(), nil)
	_require.NoError(err)

	count := int64(contentSize)
	appendBlockURLOptions := appendblob.AppendBlockFromURLOptions{
		Range:                   blob.HTTPRange{Offset: 0, Count: count},
		SourceContentValidation: blob.SourceContentValidationTypeCRC64(crc),
	}
	_, err = destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)
	_require.NoError(err)

	// TODO: This does not work... ContentCRC64 is not returned. Fix this later.
	// _require.EqualValues(appendFromURLResp.ContentCRC64, crc)

	// Check data integrity through downloading.
	destBuffer := make([]byte, 4*1024)
	downloadBufferOptions := blob.DownloadBufferOptions{Range: blob.HTTPRange{Offset: 0, Count: 4096}}
	_, err = destBlob.DownloadBuffer(context.Background(), destBuffer, &downloadBufferOptions)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCRC64Negative() {
	s.T().Skip("This test is skipped because of issues in the service.")

	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	crc64Value := crc64.Checksum(sourceData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)
	srcBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))

	// Prepare source abClient for copy.
	_, err = srcBlob.Create(context.Background(), nil)
	_require.NoError(err)

	appendResp, err := srcBlob.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)
	_require.EqualValues(appendResp.ContentCRC64, crc)

	// Get source abClient URL with SAS for AppendBlockFromURL.
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

	// Append block from URL.
	_, err = destBlob.Create(context.Background(), nil)
	_require.NoError(err)
	count := int64(contentSize)

	// Test append block from URL with bad CRC64 value
	_, sourceData = testcommon.GetDataAndReader(testName, 16)
	crc64Value = crc64.Checksum(sourceData, shared.CRC64Table)
	badCRC := make([]byte, 8)
	binary.LittleEndian.PutUint64(badCRC, crc64Value)
	appendBlockURLOptions := appendblob.AppendBlockFromURLOptions{
		Range:                   blob.HTTPRange{Offset: 0, Count: count},
		SourceContentValidation: blob.SourceContentValidationTypeCRC64(crc),
	}
	_, err = destBlob.AppendBlockFromURL(context.Background(), srcBlobURLWithSAS, &appendBlockURLOptions)

	// TODO: AppendBlockFromURL should fail, but is currently not working due to service issue.
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CRC64Mismatch)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockFromURLCopySourceAuth() {
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
	srcABClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))
	destABClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))

	// Upload some data to source
	_, err = srcABClient.Create(context.Background(), nil)
	_require.NoError(err)
	contentSize := 4 * 1024 // 4KB
	r, sourceData := testcommon.GetDataAndReader(testName, contentSize)
	_, err = srcABClient.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)
	_, err = destABClient.Create(context.Background(), nil)
	_require.NoError(err)

	// Getting token
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	_require.NoError(err)

	options := appendblob.AppendBlockFromURLOptions{
		CopySourceAuthorization: to.Ptr("Bearer " + token.Token),
	}

	pbResp, err := destABClient.AppendBlockFromURL(context.Background(), srcABClient.URL(), &options)
	_require.NoError(err)
	_require.NotNil(pbResp)

	// Download data from destination
	destBuffer := make([]byte, 4*1024)
	_, err = destABClient.DownloadBuffer(context.Background(), destBuffer, nil)
	_require.NoError(err)
	_require.Equal(destBuffer, sourceData)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockFromURLCopySourceAuthNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	// Create source and destination blobs
	srcABClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appendsrc"))
	destABClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName("appenddest"))

	// Upload some data to source
	_, err = srcABClient.Create(context.Background(), nil)
	_require.NoError(err)
	contentSize := 4 * 1024 // 4KB
	r, _ := testcommon.GetDataAndReader(testName, contentSize)
	_, err = srcABClient.AppendBlock(context.Background(), streaming.NopCloser(r), nil)
	_require.NoError(err)
	_, err = destABClient.Create(context.Background(), nil)
	_require.NoError(err)

	options := appendblob.AppendBlockFromURLOptions{
		CopySourceAuthorization: to.Ptr("Bearer faketoken"),
	}

	_, err = destABClient.AppendBlockFromURL(context.Background(), srcABClient.URL(), &options)
	_require.Error(err)
	_require.True(bloberror.HasCode(err, bloberror.CannotVerifyCopySource))
}

func (s *AppendBlobUnrecordedTestsSuite) TestGetSASURLAppendBlobClient() {
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

	// Creating append blob client with credentials
	blockBlobName := testcommon.GenerateBlobName(testName)
	apClient := testcommon.CreateNewAppendBlob(context.Background(), _require, blockBlobName, containerClient)

	// Adding SAS and options
	permissions := sas.BlobPermissions{
		Read:   true,
		Add:    true,
		Write:  true,
		Create: true,
		Delete: true,
	}
	expiry := time.Now().Add(5 * time.Minute)

	sasUrl, err := apClient.GetSASURL(permissions, expiry, nil)
	_require.NoError(err)

	// Get new blob client with sasUrl and attempt GetProperties
	newClient, err := blob.NewClientWithNoCredential(sasUrl, nil)
	_require.NoError(err)

	_, err = newClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(context.Background(), &appendblob.CreateOptions{
		Metadata: testcommon.BasicMetadata,
	})
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: map[string]*string{},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.Metadata)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: map[string]*string{"In valid!": to.Ptr("bar")},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Error(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	h := blob.ParseHTTPHeaders(resp)
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func validateAppendBlobPut(_require *require.Assertions, abClient *appendblob.Client) {
	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
	_require.EqualValues(blob.ParseHTTPHeaders(resp), testcommon.BasicHeaders)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)

	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendSetImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.NoError(err)
	immutabilityPolicySetting := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.NoError(err)

	setImmutabilityPolicyOptions := &blob.SetImmutabilityPolicyOptions{
		Mode:                     &immutabilityPolicySetting,
		ModifiedAccessConditions: nil,
	}
	_, err = abClient.SetImmutabilityPolicy(context.Background(), currentTime, setImmutabilityPolicyOptions)
	_require.NoError(err)

	_, err = abClient.SetLegalHold(context.Background(), false, nil)
	_require.NoError(err)

	_, err = abClient.Delete(context.Background(), nil)
	_require.Error(err)

	_, err = abClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendDeleteImmutabilityPolicy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 GMT 2049")
	_require.NoError(err)

	immutabilityPolicySetting := blob.ImmutabilityPolicySetting(blob.ImmutabilityPolicySettingUnlocked)
	_require.NoError(err)

	setImmutabilityPolicyOptions := &blob.SetImmutabilityPolicyOptions{
		Mode:                     &immutabilityPolicySetting,
		ModifiedAccessConditions: nil,
	}
	_, err = abClient.SetImmutabilityPolicy(context.Background(), currentTime, setImmutabilityPolicyOptions)
	_require.NoError(err)

	_, err = abClient.DeleteImmutabilityPolicy(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.Delete(context.Background(), nil)
	_require.NoError(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendSetLegalHold() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountImmutable, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainerUsingManagementClient(_require, testcommon.TestAccountImmutable, containerName)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetLegalHold(context.Background(), true, nil)
	_require.NoError(err)

	// should fail since time has not passed yet
	_, err = abClient.Delete(context.Background(), nil)
	_require.Error(err)

	_, err = abClient.SetLegalHold(context.Background(), false, nil)
	_require.NoError(err)

	_, err = abClient.Delete(context.Background(), nil)
	_require.NoError(err)

}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	eTag := azcore.ETag("garbage")
	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	validateAppendBlobPut(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobCreateAppendIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		HTTPHeaders: &testcommon.BasicHeaders,
		Metadata:    testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockNilBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(bytes.NewReader(nil)), nil)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("")), nil)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func validateBlockAppended(_require *require.Assertions, abClient *appendblob.Client, expectedSize int) {
	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentLength, int64(expectedSize))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.NoError(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.NoError(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(appendBlobCreateResp.RawResponse.StatusCode, 201)
	_require.NotNil(appendBlobCreateResp.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
	_require.Error(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.NoError(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	})
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: to.Ptr(azcore.ETag("garbage")),
			},
		},
	})
	_require.NoError(err)
	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	resp, _ := abClient.GetProperties(context.Background(), nil)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

// TODO: Fix this
// func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNegOne() {
//	bsu := testcommon.GetServiceClient()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	abClient, _ := createNewAppendBlob(c, containerClient)
//
//	appendPosition := int64(-1)
//	appendBlockOptions := appendblob.AppendBlockOptions{
//		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
//			AppendPosition: &appendPosition,
//		},
//	}
//	_, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions) // This will cause the library to set the value of the header to 0
//	_require.Error(err)
//
//	validateBlockAppended(c, abClient, len(testcommon.BlockBlobDefaultData))
// }

// func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchZero() {
//	bsu := testcommon.GetServiceClient()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//	abClient, _ := createNewAppendBlob(c, containerClient)
//
//	_, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil) // The position will not match, but the condition should be ignored
//	_require.NoError(err)
//
//	appendPosition := int64(0)
//	appendBlockOptions := appendblob.AppendBlockOptions{
//		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
//			AppendPosition: &appendPosition,
//		},
//	}
//	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendBlockOptions)
//	_require.NoError(err)
//
//	validateBlockAppended(c, abClient, 2*len(testcommon.BlockBlobDefaultData))
// }

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNonZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			AppendPosition: to.Ptr(int64(len(testcommon.BlockBlobDefaultData))),
		},
	})
	_require.NoError(err)

	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData)*2)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), nil)
	_require.NoError(err)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			AppendPosition: to.Ptr[int64](-1),
		},
	})
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNonZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			AppendPosition: to.Ptr[int64](12),
		},
	})
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.AppendPositionConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMaxSizeTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			MaxSize: to.Ptr(int64(len(testcommon.BlockBlobDefaultData) + 1)),
		},
	})
	_require.NoError(err)
	validateBlockAppended(_require, abClient, len(testcommon.BlockBlobDefaultData))
}

func (s *AppendBlobRecordedTestsSuite) TestBlobAppendBlockIfMaxSizeFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(testcommon.BlockBlobDefaultData)), &appendblob.AppendBlockOptions{
		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
			MaxSize: to.Ptr(int64(len(testcommon.BlockBlobDefaultData) - 1)),
		},
	})
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MaxBlobSizeConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestSeal() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	appendResp, err := abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NoError(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))

	resp, err := abClient.Seal(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.IsSealed, true)

	appendResp, err = abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, "BlobIsSealed")

	getPropResp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getPropResp.IsSealed, true)
}

func (s *AppendBlobRecordedTestsSuite) TestSealWithAccessConditions() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	appendResp, err := abClient.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NoError(err)
	// _require.Equal(appendResp.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp.BlobAppendOffset, "0")
	_require.Equal(*appendResp.BlobCommittedBlockCount, int32(1))

	futureTime := testcommon.GetRelativeTimeFromAnchor(appendResp.Date, 10).AddDate(1, 1, 1)
	pastTime := futureTime.AddDate(-3, -3, -3)
	sealOpts := &appendblob.SealOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &futureTime},
	}, AppendPositionAccessConditions: nil}

	_, err = abClient.Seal(context.Background(), sealOpts)
	// seal should fail on the condition
	_require.Error(err)

	sealOpts = &appendblob.SealOptions{AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &pastTime},
	}, AppendPositionAccessConditions: nil}

	resp, err := abClient.Seal(context.Background(), sealOpts)
	_require.NoError(err)
	_require.Equal(*resp.IsSealed, true)

	_, err = abClient.Delete(context.Background(), nil)
	_require.NoError(err)

}

// TODO: Learn about the behaviour of AppendPosition
// func (s *AppendBlobUnrecordedTestsSuite) TestSealWithAppendConditions() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	abName := testcommon.GenerateBlobName(testName)
//	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)
//
//	sealResp, err := abClient.Seal(context.Background(), &AppendBlobSealOptions{
//		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
//			AppendPosition: to.Ptr(1),
//		},
//	})
//	_require.Error(err)
//	_ = sealResp
//
//	sealResp, err = abClient.Seal(context.Background(), &AppendBlobSealOptions{
//		AppendPositionAccessConditions: &appendblob.AppendPositionAccessConditions{
//			AppendPosition: to.Ptr(0),
//		},
//	})
// }

func (s *AppendBlobRecordedTestsSuite) TestCopySealedBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	_, err = abClient.Seal(context.Background(), nil)
	_require.NoError(err)

	copiedBlob1 := getAppendBlobClient("copy1"+abName, containerClient)
	// copy sealed blob will get a sealed blob
	_, err = copiedBlob1.StartCopyFromURL(context.Background(), abClient.URL(), nil)
	_require.NoError(err)

	getResp1, err := copiedBlob1.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp1.IsSealed, true)

	_, err = copiedBlob1.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, "BlobIsSealed")

	copiedBlob2 := getAppendBlobClient("copy2"+abName, containerClient)
	_, err = copiedBlob2.StartCopyFromURL(context.Background(), abClient.URL(), &blob.StartCopyFromURLOptions{
		SealBlob: to.Ptr(true),
	})
	_require.NoError(err)

	getResp2, err := copiedBlob2.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp2.IsSealed, true)

	_, err = copiedBlob2.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.Error(err)
	testcommon.ValidateBlobErrorCode(_require, err, "BlobIsSealed")

	copiedBlob3 := getAppendBlobClient("copy3"+abName, containerClient)
	_, err = copiedBlob3.StartCopyFromURL(context.Background(), abClient.URL(), &blob.StartCopyFromURLOptions{
		SealBlob: to.Ptr(false),
	})
	_require.NoError(err)

	getResp3, err := copiedBlob3.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(getResp3.IsSealed)

	appendResp3, err := copiedBlob3.AppendBlock(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), nil)
	_require.NoError(err)
	// _require.Equal(appendResp3.RawResponse.StatusCode, 201)
	_require.Equal(*appendResp3.BlobAppendOffset, "0")
	_require.Equal(*appendResp3.BlobCommittedBlockCount, int32(1))
}

func (s *AppendBlobRecordedTestsSuite) TestCopyUnsealedBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := createNewAppendBlob(context.Background(), _require, abName, containerClient)

	copiedBlob := getAppendBlobClient("copy"+abName, containerClient)
	_, err = copiedBlob.StartCopyFromURL(context.Background(), abClient.URL(), &blob.StartCopyFromURLOptions{
		SealBlob: to.Ptr(true),
	})
	_require.NoError(err)

	getResp, err := copiedBlob.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*getResp.IsSealed, true)
}

func (s *AppendBlobUnrecordedTestsSuite) TestCreateAppendBlobWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	// _context := getTestContext(testName)
	// ignoreHeaders(_context.recording, []string{"x-ms-tags", "X-Ms-Tags"})
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)

	createAppendBlobOptions := appendblob.CreateOptions{
		Tags: testcommon.SpecialCharBlobTagsMap,
	}
	createResp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)
	_require.NotNil(createResp.VersionID)
	time.Sleep(10 * time.Second)

	_, err = abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	blobGetTagsResponse, err := abClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.SpecialCharBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.SpecialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	resp, err := abClient.CreateSnapshot(context.Background(), nil)
	_require.NoError(err)

	snapshotURL, _ := abClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp2.TagCount, int64(len(testcommon.SpecialCharBlobTagsMap)))

	blobGetTagsResponse, err = abClient.GetTags(context.Background(), nil)
	_require.NoError(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet = blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.SpecialCharBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.SpecialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	// Tags with spaces
	where := "\"GO \"='.Net'"
	lResp, err := svcClient.FilterBlobs(context.Background(), where, nil)
	_require.NoError(err)
	_require.Len(lResp.FilterBlobSegment.Blobs[0].Tags.BlobTagSet, 1)
	_require.Equal(lResp.FilterBlobSegment.Blobs[0].Tags.BlobTagSet[0], blobTagsSet[2])
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobGetPropertiesUsingVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	abClient := createNewAppendBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	blobProp, _ := abClient.GetProperties(context.Background(), nil)

	createAppendBlobOptions := appendblob.CreateOptions{
		Metadata: testcommon.BasicMetadata,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: blobProp.ETag},
		},
	}
	createResp, err := abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)
	_require.NotNil(createResp.VersionID)
	blobProp, _ = abClient.GetProperties(context.Background(), nil)
	_require.EqualValues(createResp.VersionID, blobProp.VersionID)
	_require.EqualValues(createResp.LastModified, blobProp.LastModified)
	_require.Equal(*createResp.ETag, *blobProp.ETag)
	_require.Equal(*blobProp.IsCurrentVersion, true)
}

func (s *AppendBlobUnrecordedTestsSuite) TestSetBlobMetadataReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	bbName := testcommon.GenerateName(testName)
	bbClient := testcommon.CreateNewBlockBlob(context.Background(), _require, bbName, containerClient)

	metadata := map[string]*string{"test_key_1": to.Ptr("test_value_1"), "test_key_2": to.Ptr("2019")}
	resp, err := bbClient.SetMetadata(context.Background(), metadata, nil)
	_require.NoError(err)
	_require.NotNil(resp.VersionID)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Metadata: true},
	})

	if pager.More() {
		pageResp, err := pager.NextPage(context.Background())
		_require.NoError(err) // check for an error first
		// s.T().Fail()      // no page was gotten

		_require.NotNil(pageResp.Segment.BlobItems)
		blobList := pageResp.Segment.BlobItems
		_require.Len(blobList, 1)
		blobResp1 := blobList[0]
		_require.Equal(*blobResp1.Name, bbName)
		_require.NotNil(blobResp1.Metadata)
		_require.Len(blobResp1.Metadata, 2)
	}
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := appendblob.AppendBlockOptions{
			CPKInfo: &testcommon.TestCPKByValue,
		}
		resp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.NoError(err)
		// _require.Equal(resp.RawResponse.StatusCode, 201)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal(resp.Date.IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		if recording.GetRecordMode() != recording.PlaybackMode {
			_require.EqualValues(resp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)
		}
	}

	// Get blob content without encryption key should fail the request.
	_, err = abClient.DownloadStream(context.Background(), nil)
	_require.Error(err)

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKInfo: &testcommon.TestCPKByValue,
	}
	downloadResp, err := abClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)

	data, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	if recording.GetRecordMode() != recording.PlaybackMode {
		_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
	}
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	encryptionScope := testcommon.GetCPKScopeInfo(s.T())
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{
		CPKScopeInfo: &encryptionScope,
	}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		appendBlockOptions := appendblob.AppendBlockOptions{
			CPKScopeInfo: &encryptionScope,
		}
		resp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader(word)), &appendBlockOptions)
		_require.NoError(err)
		_require.Equal(*resp.BlobAppendOffset, strconv.Itoa(index*4))
		_require.Equal(*resp.BlobCommittedBlockCount, int32(index+1))
		_require.NotNil(resp.ETag)
		_require.NotNil(resp.LastModified)
		_require.Equal(resp.LastModified.IsZero(), false)
		_require.NotEqual(resp.ContentMD5, "")

		_require.NotEqual(resp.Version, "")
		_require.NotNil(resp.Date)
		_require.Equal(resp.Date.IsZero(), false)
		_require.Equal(*resp.IsServerEncrypted, true)
		_require.EqualValues(*encryptionScope.EncryptionScope, *resp.EncryptionScope)
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CPKScopeInfo: &encryptionScope,
	}
	downloadResp, err := abClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NoError(err)

	data, err := io.ReadAll(downloadResp.Body)
	_require.NoError(err)
	_require.EqualValues(string(data), "AAA BBB CCC ")
	_require.EqualValues(*downloadResp.EncryptionScope, *encryptionScope.EncryptionScope)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockPermanentDelete() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	// Create container and blob, upload blob to container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	createAppendBlobOptions := appendblob.CreateOptions{}
	_, err = abClient.Create(context.Background(), &createAppendBlobOptions)
	_require.NoError(err)

	parts, err := sas.ParseURL(abClient.URL()) // Get parts for BlobURL
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Set Account SAS and set Permanent Delete to true
	parts.SAS, err = sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true, PermanentDelete: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	// Create snapshot of Blob and get snapshot URL
	resp, err := abClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})
	_require.NoError(err)
	snapshotURL, _ := abClient.WithSnapshot(*resp.Snapshot)

	// Check that there are two items in the container: one snapshot, one blob
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Include: container.ListBlobsInclude{Snapshots: true}})
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

	// Delete snapshot (snapshot will be soft deleted)
	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = abClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.NoError(err)

	// Check that only blob exists (snapshot is soft-deleted)
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
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
	_require.Len(found, 1)

	// Check that soft-deleted snapshot exists by including deleted items
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Deleted: true},
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
	_require.Len(found, 2)

	// Options for PermanentDeleteOptions
	perm := blob.DeleteTypePermanent
	deleteBlobOptions := blob.DeleteOptions{
		BlobDeleteType: &perm,
	}
	// Execute Delete with DeleteTypePermanent
	pdResp, err := snapshotURL.Delete(context.Background(), &deleteBlobOptions)
	_require.NoError(err)
	_require.NotNil(pdResp)

	// Check that only blob exists even after including snapshots and deleted items
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true, Deleted: true}})
	found = make([]*container.BlobItem, 0)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		found = append(found, resp.Segment.BlobItems...)
		if err != nil {
			break
		}
	}
	_require.Len(found, 1)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockPermanentDeleteWithoutPermission() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	// Create container and blob, upload blob to container
	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := containerClient.NewAppendBlobClient(testcommon.GenerateBlobName(testName))

	_, err = abClient.Create(context.Background(), &appendblob.CreateOptions{})
	_require.NoError(err)

	parts, err := sas.ParseURL(abClient.URL()) // Get parts for BlobURL
	_require.NoError(err)

	credential, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	// Set Account SAS
	parts.SAS, err = sas.AccountSignatureValues{
		Protocol:      sas.ProtocolHTTPS,                    // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		Permissions:   to.Ptr(sas.AccountPermissions{Read: true, List: true}).String(),
		ResourceTypes: to.Ptr(sas.AccountResourceTypes{Container: true, Object: true}).String(),
	}.SignWithSharedKey(credential)
	_require.NoError(err)

	// Create snapshot of Blob and get snapshot URL
	resp, err := abClient.CreateSnapshot(context.Background(), &blob.CreateSnapshotOptions{})
	_require.NoError(err)
	snapshotURL, _ := abClient.WithSnapshot(*resp.Snapshot)

	// Check that there are two items in the container: one snapshot, one blob
	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{Include: container.ListBlobsInclude{Snapshots: true}})
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

	// Delete snapshot
	deleteSnapshotsOnly := blob.DeleteSnapshotsOptionTypeOnly
	_, err = abClient.Delete(context.Background(), &blob.DeleteOptions{DeleteSnapshots: &deleteSnapshotsOnly})
	_require.NoError(err)

	// Check that only blob exists
	pager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
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
	_require.Len(found, 1)

	// Options for PermanentDeleteOptions
	perm := blob.DeleteTypePermanent
	deleteBlobOptions := blob.DeleteOptions{
		BlobDeleteType: &perm,
	}
	// Execute Delete with DeleteTypePermanent,should fail because permissions are not set and snapshot is not soft-deleted
	_, err = snapshotURL.Delete(context.Background(), &deleteBlobOptions)
	_require.Error(err)
}

// nolint

// func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCPK() {
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
//	_require.NoError(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	resp, err := srcABClient.AppendBlock(ctx, streaming.NopCloser(r), nil)
//	_require.NoError(err)
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
//	_require.NoError(err)
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
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	_, err = destBlob.Create(ctx, &createAppendBlobOptions)
//	_require.NoError(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
//		Offset:  &offset,
//		Count:   &count,
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.NoError(err)
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
//	_require.Error(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions := blob.downloadWriterAtOptions{
//		CPKInfo: &testcommon.TestInvalidCPKByValue,
//	}
//	_, err = destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.Error(err)
//
//	// Download blob to do data integrity check.
//	downloadBlobOptions = blob.downloadWriterAtOptions{
//		CPKInfo: &testcommon.TestCPKByValue,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NoError(err)
//
//	_require.Equal(*downloadResp.IsServerEncrypted, true)
//	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CPKInfo: &testcommon.TestCPKByValue}))
//	_require.NoError(err)
//	_require.EqualValues(destData, srcData)
// }

// func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlockFromURLWithCPKScope() {
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
//	_require.NoError(err)
//	//_require.Equal(cResp1.RawResponse.StatusCode, 201)
//
//	resp, err := srcClient.AppendBlock(ctx, streaming.NopCloser(r), nil)
//	_require.NoError(err)
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
//	_require.NoError(err)
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
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	_, err = destBlob.Create(ctx, &createAppendBlobOptions)
//	_require.NoError(err)
//	//_require.Equal(cResp2.RawResponse.StatusCode, 201)
//
//	offset := int64(0)
//	count := int64(contentSize)
//	appendBlockURLOptions := AppendBlobAppendBlockFromURLOptions{
//		Offset:       &offset,
//		Count:        &count,
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
//	_require.NoError(err)
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
//		CPKScopeInfo: &testcommon.TestCPKByScope,
//	}
//	downloadResp, err := destBlob.DownloadStream(ctx, &downloadBlobOptions)
//	_require.NoError(err)
//	_require.Equal(*downloadResp.IsServerEncrypted, true)
//	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
//
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{CPKInfo: &testcommon.TestCPKByValue}))
//	_require.NoError(err)
//	_require.EqualValues(destData, srcData)
// }

func (s *AppendBlobRecordedTestsSuite) TestUndeleteAppendBlobVersion() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	versions := make([]string, 0)
	for i := 0; i < 5; i++ {
		resp, err := abClient.CreateSnapshot(context.Background(), nil)
		_require.NoError(err)
		_require.NotNil(resp.VersionID)
		versions = append(versions, *resp.VersionID)
	}

	listPager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6)

	// Deleting the 1st, 2nd and 3rd versions
	for i := 0; i < 3; i++ {
		abClientWithVersionID, err := abClient.WithVersionID(versions[i])
		_require.NoError(err)
		_, err = abClientWithVersionID.Delete(context.Background(), nil)
		_require.NoError(err)
	}

	// adding wait after delete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 3)

	_, err = abClient.Undelete(context.Background(), nil)
	_require.NoError(err)

	// adding wait after undelete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Versions: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6)
}

func (s *AppendBlobRecordedTestsSuite) TestUndeleteAppendBlobSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountSoftDelete, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	snapshots := make([]string, 0)
	for i := 0; i < 5; i++ {
		resp, err := abClient.CreateSnapshot(context.Background(), nil)
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
		abClientWithSnapshot, err := abClient.WithSnapshot(snapshots[i])
		_require.NoError(err)
		_, err = abClientWithSnapshot.Delete(context.Background(), nil)
		_require.NoError(err)
	}

	// adding wait after delete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 3) // 2 snapshots and 1 current version

	_, err = abClient.Undelete(context.Background(), nil)
	_require.NoError(err)

	// adding wait after undelete
	time.Sleep(time.Second * 10)

	listPager = containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		Include: container.ListBlobsInclude{Snapshots: true},
	})
	testcommon.ListBlobsCount(context.Background(), _require, listPager, 6) // 5 snapshots and 1 current version
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetExpiryToNeverExpire() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = abClient.SetExpiry(context.Background(), appendblob.ExpiryTypeNever{}, nil)
	_require.NoError(err)

	resp, err = abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetExpiryRelativeToNow() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = abClient.SetExpiry(context.Background(), appendblob.ExpiryTypeRelativeToNow(8*time.Second), nil)
	_require.NoError(err)

	resp, err = abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ExpiresOn)

	time.Sleep(time.Second * 10)

	_, err = abClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetExpiryRelativeToCreation() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	_, err = abClient.SetExpiry(context.Background(), appendblob.ExpiryTypeRelativeToCreation(8*time.Second), nil)
	_require.NoError(err)

	resp, err = abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ExpiresOn)

	time.Sleep(time.Second * 10)

	_, err = abClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *AppendBlobUnrecordedTestsSuite) TestAppendBlobSetExpiryToAbsolute() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDatalake, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Nil(resp.ExpiresOn)

	expiryTimeAbsolute := time.Now().Add(8 * time.Second)
	_, err = abClient.SetExpiry(context.Background(), appendblob.ExpiryTypeAbsolute(expiryTimeAbsolute), nil)
	_require.NoError(err)

	resp, err = abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.NotNil(resp.ExpiresOn)
	_require.Equal(expiryTimeAbsolute.UTC().Format(http.TimeFormat), resp.ExpiresOn.UTC().Format(http.TimeFormat))

	time.Sleep(time.Second * 10)

	_, err = abClient.GetProperties(context.Background(), nil)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataNil() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetMetadata(context.Background(), map[string]*string{"not": to.Ptr("nil")}, nil)
	_require.NoError(err)

	_, err = abClient.SetMetadata(context.Background(), nil, nil)
	_require.NoError(err)

	blobGetResp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(blobGetResp.Metadata, 0)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetMetadata(context.Background(), map[string]*string{"not": to.Ptr("nil")}, nil)
	_require.NoError(err)

	_, err = abClient.SetMetadata(context.Background(), map[string]*string{}, nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(resp.Metadata, 0)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataInvalidField() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetMetadata(context.Background(), map[string]*string{"Invalid field!": to.Ptr("value")}, nil)
	_require.Error(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)
}

func validateMetadataSet(_require *require.Assertions, abClient *appendblob.Client) {
	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.NoError(err)

	validateMetadataSet(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.NoError(err)

	validateMetadataSet(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)
	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: to.Ptr(currentTime)},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.NoError(err)

	validateMetadataSet(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: to.Ptr(azcore.ETag("garbage"))},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(azcore.ETag("garbage"))},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	_require.NoError(err)

	validateMetadataSet(_require, abClient)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetMetadataIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	setBlobMetadataOptions := blob.SetMetadataOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	}
	_, err = abClient.SetMetadata(context.Background(), testcommon.BasicMetadata, &setBlobMetadataOptions)
	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func validatePropertiesSet(_require *require.Assertions, abClient *appendblob.Client, disposition string) {
	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(*resp.ContentDisposition, disposition)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
			},
		})
	_require.NoError(err)

	validatePropertiesSet(_require, abClient, "my_disposition")
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{
			AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime},
			}})
	_require.Error(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, 10)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		}})
	_require.NoError(err)

	validatePropertiesSet(_require, abClient, "my_disposition")
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	cResp, err := abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(cResp.Date, -10)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		}})
	_require.Error(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: resp.ETag},
		}})
	_require.NoError(err)

	validatePropertiesSet(_require, abClient, "my_disposition")
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: to.Ptr(azcore.ETag("garbage"))},
		}})
	_require.Error(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: to.Ptr(azcore.ETag("garbage"))},
		}})
	_require.NoError(err)

	validatePropertiesSet(_require, abClient, "my_disposition")
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetHTTPHeaderIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	resp, err := abClient.GetProperties(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetHTTPHeaders(context.Background(), blob.HTTPHeaders{BlobContentDisposition: to.Ptr("my_disposition")},
		&blob.SetHTTPHeadersOptions{AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		}})
	_require.Error(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobSetBlobTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)
	_, err = abClient.AppendBlock(context.Background(), streaming.NopCloser(strings.NewReader("Appending block\n")), nil)
	_require.NoError(err)

	var tagsMap = map[string]string{
		"azure": "blob",
	}

	_, err = abClient.SetTags(context.Background(), tagsMap, nil)
	_require.NoError(err)
	time.Sleep(10 * time.Second)

	blobGetTagsResponse, err := abClient.GetTags(context.Background(), nil)
	_require.NoError(err)

	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, 1)
	for _, blobTag := range blobTagsSet {
		_require.Equal(tagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *AppendBlobUnrecordedTestsSuite) TestSetBlobTagsWithLeaseId() {
	_require := require.New(s.T())
	testName := "ab" + s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	blobLeaseClient, err := lease.NewBlobClient(abClient, &lease.BlobClientOptions{
		LeaseID: to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c"),
	})
	_require.NoError(err)
	ctx := context.Background()
	acquireLeaseResponse, err := blobLeaseClient.AcquireLease(ctx, int32(60), nil)
	_require.NoError(err)
	_require.NotNil(acquireLeaseResponse.LeaseID)
	_require.EqualValues(acquireLeaseResponse.LeaseID, blobLeaseClient.LeaseID())

	_, err = abClient.SetTags(ctx, testcommon.BasicBlobTagsMap, nil)
	_require.Error(err)
	time.Sleep(10 * time.Second)

	// add lease conditions
	_, err = abClient.SetTags(ctx, testcommon.BasicBlobTagsMap, &blob.SetTagsOptions{AccessConditions: &blob.AccessConditions{
		LeaseAccessConditions: &blob.LeaseAccessConditions{LeaseID: blobLeaseClient.LeaseID()}}})
	_require.NoError(err)

	_, err = abClient.GetTags(ctx, nil)
	_require.NoError(err)

	blobGetTagsResponse, err := abClient.GetTags(ctx, &blob.GetTagsOptions{BlobAccessConditions: &blob.AccessConditions{
		LeaseAccessConditions: &blob.LeaseAccessConditions{LeaseID: blobLeaseClient.LeaseID()}}})
	_require.NoError(err)

	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, 3)
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.BasicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *AppendBlobRecordedTestsSuite) TestAppendGetAccountInfo() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abClient := getAppendBlobClient(testcommon.GenerateBlobName(testName), containerClient)
	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	bAccInfo, err := abClient.GetAccountInfo(context.Background(), nil)
	_require.NoError(err)
	_require.NotZero(bAccInfo)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlockSetTier() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	abName := testcommon.GenerateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClient.SetTier(context.Background(), blob.AccessTierHot, nil)
	_require.ErrorContains(err, "operation will not work on this blob type. SetTier only works for page blob in premium storage account and block blob in blob storage account")
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobClientDefaultAudience() {
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

	options := &appendblob.ClientOptions{
		Audience: "https://storage.azure.com/",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	abClientAudience, err := appendblob.NewClient(blobURL, cred, options)
	_require.NoError(err)

	_, err = abClientAudience.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

func (s *AppendBlobRecordedTestsSuite) TestAppendBlobClientCustomAudience() {
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

	options := &appendblob.ClientOptions{
		Audience: "https://" + accountName + ".blob.core.windows.net",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	abClientAudience, err := appendblob.NewClient(blobURL, cred, options)
	_require.NoError(err)

	_, err = abClientAudience.Create(context.Background(), nil)
	_require.NoError(err)

	_, err = abClientAudience.GetProperties(context.Background(), nil)
	_require.NoError(err)
}

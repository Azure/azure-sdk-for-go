//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

type AZBlobRecordedTestsSuite struct {
	suite.Suite
	proxy *recording.TestProxyInstance
}

type AZBlobUnrecordedTestsSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	recordMode := recording.GetRecordMode()
	t.Logf("Running azblob Tests in %s mode\n", recordMode)
	switch recordMode {
	case recording.LiveMode:
		suite.Run(t, &AZBlobRecordedTestsSuite{})
		suite.Run(t, &AZBlobUnrecordedTestsSuite{})
	case recording.PlaybackMode:
		suite.Run(t, &AZBlobRecordedTestsSuite{})
	case recording.RecordingMode:
		suite.Run(t, &AZBlobRecordedTestsSuite{})
	}
}

func (s *AZBlobRecordedTestsSuite) SetupSuite() {
	s.proxy = testcommon.SetupSuite(&s.Suite)
}

func (s *AZBlobRecordedTestsSuite) TearDownSuite() {
	testcommon.TearDownSuite(&s.Suite, s.proxy)
}

func (s *AZBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

func (s *AZBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

func (s *AZBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

func (s *AZBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

func (s *AZBlobRecordedTestsSuite) TestAzBlobClientSharedKey() {
	_require := require.New(s.T())

	cred, err := testcommon.GetGenericSharedKeyCredential(testcommon.TestAccountDefault)
	_require.NoError(err)

	svcURL := "https://" + cred.AccountName() + ".blob.core.windows.net/"
	azClient, err := azblob.NewClientWithSharedKeyCredential(svcURL, cred, nil)
	_require.NoError(err)
	_require.NotNil(azClient)
	_require.Equal(azClient.URL(), svcURL)
}

func (s *AZBlobRecordedTestsSuite) TestAzBlobClientConnectionString() {
	_require := require.New(s.T())

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	connString, err := testcommon.GetGenericConnectionString(testcommon.TestAccountDefault)
	_require.NoError(err)
	_require.NotNil(connString)

	svcURL := "https://" + accountName + ".blob.core.windows.net/"
	azClient, err := azblob.NewClientFromConnectionString(*connString, nil)
	_require.NoError(err)
	_require.NotNil(azClient)
	_require.Equal(azClient.URL(), svcURL)
}

// create a test file
func generateFile(fileName string, fileSize int) []byte {
	// generate random data
	_, bigBuff := testcommon.GenerateData(fileSize)

	// write to file and return the data
	_ = os.WriteFile(fileName, bigBuff, 0666)
	return bigBuff
}

func performUploadStreamToBlockBlobTestWithChecksums(t *testing.T, _require *require.Assertions, testName string, blobSize, bufferSize, maxBuffers int) {
	client, err := testcommon.GetClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)
	defer func() {
		_, err := client.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)

	// Create some data to test the upload stream
	blobContentReader, blobData := testcommon.GenerateData(blobSize)
	crc64Value := crc64.Checksum(blobData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	// Perform UploadStream
	_, err = client.UploadStream(ctx, containerName, blobName, blobContentReader,
		&blockblob.UploadStreamOptions{BlockSize: int64(bufferSize), Concurrency: maxBuffers, TransactionalValidation: blob.TransferValidationTypeComputeCRC64()})
	_require.NoError(err)

	// TODO: UploadResponse does not return ContentCRC64
	// // Assert that upload was successful
	// _require.Equal(uploadResp.ContentCRC64, crc)

	// UploadStream does not accept user generated checksum, should return UnsupportedChecksum
	_, err = client.UploadStream(ctx, containerName, blobName, blobContentReader,
		&blockblob.UploadStreamOptions{BlockSize: int64(bufferSize), Concurrency: maxBuffers, TransactionalValidation: blob.TransferValidationTypeCRC64(crc64Value)})
	_require.Error(err)
	_require.Error(err, bloberror.UnsupportedChecksum)

	md5Value := md5.Sum(blobData)
	contentMD5 := md5Value[:]

	_, err = client.UploadStream(ctx, containerName, blobName, blobContentReader,
		&blockblob.UploadStreamOptions{BlockSize: int64(bufferSize), Concurrency: maxBuffers, TransactionalValidation: blob.TransferValidationTypeMD5(contentMD5)})
	_require.Error(err)
	_require.Error(err, bloberror.UnsupportedChecksum)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobInChunksChecksums() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTestWithChecksums(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

func performUploadStreamToBlockBlobTest(t *testing.T, _require *require.Assertions, testName string, blobSize, bufferSize, maxBuffers int) {
	client, err := testcommon.GetClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)
	defer func() {
		_, err := client.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)

	// Create some data to test the upload stream
	blobContentReader, blobData := testcommon.GenerateData(blobSize)

	// Perform UploadStream
	_, err = client.UploadStream(ctx, containerName, blobName, blobContentReader,
		&blockblob.UploadStreamOptions{BlockSize: int64(bufferSize), Concurrency: maxBuffers})

	// Assert that upload was successful
	_require.NoError(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)

	// Download the blob to verify
	downloadResponse, err := client.DownloadStream(ctx, containerName, blobName, nil)
	_require.NoError(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.Body)
	_require.NoError(err)
	_require.Equal(len(actualBlobData), blobSize)
	_require.EqualValues(actualBlobData, blobData)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobInChunks() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobSingleBuffer() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 1
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobSingleIO() {
	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobSingleIOEdgeCase() {
	blobSize := 8 * 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobEmpty() {
	blobSize := 0
	bufferSize := 8 * 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

func performUploadAndDownloadFileTestWithChecksums(t *testing.T, _require *require.Assertions, testName string, fileSize, blockSize, concurrency, downloadOffset, downloadCount int) {
	// Set up file to upload
	fileName := "BigFile.bin"
	fileData := generateFile(fileName, fileSize)

	// Open the file to upload
	file, err := os.Open(fileName)
	_require.NoError(err)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	defer func(name string) {
		_ = os.Remove(name)
	}(fileName)

	// body := bytes.NewReader(fileData)
	crc64Value := crc64.Checksum(fileData, shared.CRC64Table)
	crc := make([]byte, 8)
	binary.LittleEndian.PutUint64(crc, crc64Value)

	client, err := testcommon.GetClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)
	defer func() {
		_, err := client.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)

	// Upload the file to a block blob
	var errTransferred error
	_, err = client.UploadFile(context.Background(), containerName, blobName, file,
		&blockblob.UploadFileOptions{
			BlockSize:               int64(blockSize),
			Concurrency:             uint16(concurrency),
			TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(fileSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})
	assert.NoError(t, errTransferred)
	_require.NoError(err)

	// UploadFile does not accept user generated checksum, should return UnsupportedChecksum
	_, err = client.UploadFile(context.Background(), containerName, blobName, file,
		&blockblob.UploadFileOptions{
			BlockSize:               int64(blockSize),
			Concurrency:             uint16(concurrency),
			TransactionalValidation: blob.TransferValidationTypeCRC64(crc64Value),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(fileSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})
	_require.Error(err)
	_require.Error(err, bloberror.UnsupportedChecksum)

	md5Value := md5.Sum(fileData)
	contentMD5 := md5Value[:]

	// UploadFile does not accept user generated checksum, should return UnsupportedChecksum
	_, err = client.UploadFile(context.Background(), containerName, blobName, file,
		&blockblob.UploadFileOptions{
			BlockSize:               int64(blockSize),
			Concurrency:             uint16(concurrency),
			TransactionalValidation: blob.TransferValidationTypeMD5(contentMD5),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(fileSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})
	_require.Error(err)
	_require.Error(err, bloberror.UnsupportedChecksum)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadFileInChunksChecksum() {
	fileSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTestWithChecksums(s.T(), _require, testName, fileSize, blockSize, concurrency, 0, 0)
}

func performUploadAndDownloadFileTest(t *testing.T, _require *require.Assertions, testName string, fileSize, blockSize, concurrency, downloadOffset, downloadCount int) {
	// Set up file to upload
	fileName := "BigFile.bin"
	fileData := generateFile(fileName, fileSize)

	// Open the file to upload
	file, err := os.Open(fileName)
	_require.NoError(err)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	defer func(name string) {
		_ = os.Remove(name)
	}(fileName)

	client, err := testcommon.GetClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)
	defer func() {
		_, err := client.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)

	// Upload the file to a block blob
	var errTransferred error
	_, err = client.UploadFile(context.Background(), containerName, blobName, file,
		&blockblob.UploadFileOptions{
			BlockSize:   int64(blockSize),
			Concurrency: uint16(concurrency),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(fileSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})
	assert.NoError(t, errTransferred)
	_require.NoError(err)

	// Set up file to download the blob to
	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_require.NoError(err)
	defer func(destFile *os.File) {
		_ = destFile.Close()

	}(destFile)
	defer func(name string) {
		_ = os.Remove(name)

	}(destFileName)

	filePointer, err := file.Seek(0, io.SeekCurrent)
	_require.NoError(err)

	// Perform download
	_, err = client.DownloadFile(context.Background(),
		containerName,
		blobName,
		destFile,
		&blob.DownloadFileOptions{
			Range: azblob.HTTPRange{
				Count:  int64(downloadCount),
				Offset: int64(downloadOffset),
			},
			BlockSize:   int64(blockSize),
			Concurrency: uint16(concurrency),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(fileSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})

	// Assert download was successful
	assert.NoError(t, errTransferred)
	_require.NoError(err)

	currPointer, err := destFile.Seek(0, io.SeekCurrent)
	_require.Equal(filePointer, currPointer)
	_require.NoError(err)

	// Assert downloaded data is consistent
	var destBuffer []byte
	if downloadCount == blob.CountToEnd {
		destBuffer = make([]byte, fileSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	n, err := destFile.Read(destBuffer)
	_require.NoError(err)

	if downloadOffset == 0 && downloadCount == 0 {
		_require.EqualValues(destBuffer, fileData)
	} else {
		if downloadCount == 0 {
			_require.Equal(n, fileSize-downloadOffset)
			_require.EqualValues(destBuffer, fileData[downloadOffset:])
		} else {
			_require.Equal(n, downloadCount)
			_require.EqualValues(destBuffer, fileData[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileInChunks() {
	fileSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileSingleIO() {
	fileSize := 1024
	blockSize := 2048
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileSingleRoutine() {
	fileSize := 8 * 1024
	blockSize := 1024
	concurrency := 1
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileEmpty() {
	fileSize := 0
	blockSize := 1024
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileNonZeroOffset() {
	fileSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	downloadOffset := 1000
	downloadCount := 0
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileNonZeroCount() {
	fileSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	downloadOffset := 0
	downloadCount := 6000
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileNonZeroOffsetAndCount() {
	fileSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	downloadOffset := 1000
	downloadCount := 6000
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func performUploadAndDownloadBufferTest(t *testing.T, _require *require.Assertions, testName string, blobSize, blockSize, concurrency, downloadOffset, downloadCount int) {
	// Set up buffer to upload
	_, bytesToUpload := testcommon.GenerateData(blobSize)

	// Set up test container
	client, err := testcommon.GetClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerName := testcommon.GenerateContainerName(testName)
	_, err = client.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)
	defer func() {
		_, err := client.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	// Set up test blob
	blobName := testcommon.GenerateBlobName(testName)

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	var errTransferred error
	_, err = client.UploadBuffer(context.Background(), containerName, blobName, bytesToUpload,
		&blockblob.UploadBufferOptions{
			BlockSize:   int64(blockSize),
			Concurrency: uint16(concurrency),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(blobSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})
	assert.NoError(t, errTransferred)
	_require.NoError(err)

	// Set up buffer to download the blob to
	var destBuffer []byte
	if downloadCount == blob.CountToEnd {
		destBuffer = make([]byte, blobSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	// Download the blob to a buffer
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	n, err := client.DownloadBuffer(ctx,
		containerName,
		blobName,
		destBuffer, &blob.DownloadBufferOptions{
			Range: azblob.HTTPRange{
				Count:  int64(downloadCount),
				Offset: int64(downloadOffset),
			},
			BlockSize:   int64(blockSize),
			Concurrency: uint16(concurrency),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				if bytesTransferred <= 0 || bytesTransferred > int64(blobSize) {
					errTransferred = fmt.Errorf("invalid bytes transferred %d", bytesTransferred)
				}
			},
		})
	cancel()

	assert.NoError(t, errTransferred)
	_require.NoError(err)

	if downloadOffset == 0 && downloadCount == 0 {
		_require.EqualValues(destBuffer, bytesToUpload)
	} else {
		if downloadCount == 0 {
			_require.EqualValues(destBuffer, bytesToUpload[downloadOffset:])
		} else {
			_require.EqualValues(destBuffer[:int(n)], bytesToUpload[downloadOffset:downloadOffset+int(n)])
			if downloadCount > blobSize {
				_require.EqualValues(blobSize-downloadOffset, n)
			} else {
				_require.EqualValues(downloadCount, n)
			}

		}
	}
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferInChunks() {
	blobSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferSingleIO() {
	blobSize := 1024
	blockSize := 8 * 1024
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferSingleRoutine() {
	blobSize := 8 * 1024
	blockSize := 1024
	concurrency := 1
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferEmpty() {
	blobSize := 0
	blockSize := 1024
	concurrency := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, 0, 0)
}

func (s *AZBlobUnrecordedTestsSuite) TestDownloadBufferWithNonZeroOffset() {
	blobSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	downloadOffset := 1000
	downloadCount := 0
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func (s *AZBlobUnrecordedTestsSuite) TestDownloadBufferWithNonZeroCount() {
	blobSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	downloadOffset := 0
	downloadCount := 6000
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func (s *AZBlobUnrecordedTestsSuite) TestDownloadBufferWithNonZeroOffsetAndCount() {
	blobSize := 8 * 1024
	blockSize := 1024
	concurrency := 3
	downloadOffset := 2000
	downloadCount := 6 * 1024
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func (s *AZBlobUnrecordedTestsSuite) TestBasicDoBatchTransfer() {
	_require := require.New(s.T())
	// test the basic multi-routine processing
	type testInstance struct {
		transferSize int64
		chunkSize    int64
		concurrency  uint16
		expectError  bool
	}

	testMatrix := []testInstance{
		{transferSize: 100, chunkSize: 10, concurrency: 5, expectError: false},
		{transferSize: 100, chunkSize: 9, concurrency: 4, expectError: false},
		{transferSize: 100, chunkSize: 8, concurrency: 15, expectError: false},
		{transferSize: 100, chunkSize: 1, concurrency: 3, expectError: false},
		{transferSize: 0, chunkSize: 100, concurrency: 5, expectError: false}, // empty file works
		{transferSize: 100, chunkSize: 0, concurrency: 5, expectError: true},  // 0 chunk size on the other hand must fail
		{transferSize: 0, chunkSize: 0, concurrency: 5, expectError: true},
	}

	for _, test := range testMatrix {
		ctx := context.Background()
		// maintain some counts to make sure the right number of chunks were queued, and the total size is correct
		totalSizeCount := int64(0)
		runCount := int64(0)

		numChunks := uint64(0)
		if test.chunkSize != 0 {
			numChunks = uint64(((test.transferSize - 1) / test.chunkSize) + 1)
		}

		err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
			TransferSize: test.transferSize,
			ChunkSize:    test.chunkSize,
			NumChunks:    numChunks,
			Concurrency:  test.concurrency,
			Operation: func(ctx context.Context, offset int64, chunkSize int64) error {
				atomic.AddInt64(&totalSizeCount, chunkSize)
				atomic.AddInt64(&runCount, 1)
				return nil
			},
			OperationName: "TestHappyPath",
		})

		if test.expectError {
			_require.Error(err)
		} else {
			_require.NoError(err)
			_require.Equal(totalSizeCount, test.transferSize)
			_require.Equal(runCount, ((test.transferSize-1)/test.chunkSize)+1)
		}
	}
}

// mock a memory mapped file (low-quality mock, meant to simulate the scenario only)
type mockMMF struct {
	isClosed   bool
	failHandle *require.Assertions
}

// accept input
func (m *mockMMF) write(_ string) {
	if m.isClosed {
		// simulate panic
		m.failHandle.Fail("")
	}
}

func (s *AZBlobUnrecordedTestsSuite) TestDoBatchTransferWithError() {
	_require := require.New(s.T())
	ctx := context.Background()
	mmf := mockMMF{failHandle: _require}
	expectedFirstError := errors.New("#3 means trouble")

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		TransferSize: 5,
		ChunkSize:    1,
		NumChunks:    5,
		Concurrency:  5,
		Operation: func(ctx context.Context, offset int64, chunkSize int64) error {
			// simulate doing some work (HTTP call in real scenarios)
			// later chunks later longer to finish
			time.Sleep(time.Second * time.Duration(offset))
			// simulate having gotten data and write it to the memory mapped file
			mmf.write("input")

			// with one of the chunks, pretend like an error occurred (like the network connection breaks)
			if offset == 3 {
				return expectedFirstError
			} else if offset > 3 {
				// anything after offset=3 are canceled
				// so verify that the context indeed got canceled
				ctxErr := ctx.Err()
				_require.Equal(ctxErr, context.Canceled)
				return ctxErr
			}

			// anything before offset=3 should be done without problem
			return nil
		},
		OperationName: "TestErrorPath",
	})

	_require.Equal(err, expectedFirstError)

	// simulate closing the mmf and make sure no panic occurs (as reported in #139)
	mmf.isClosed = true
	time.Sleep(time.Second * 5)
}

func (s *AZBlobRecordedTestsSuite) TestAzBlobClientDefaultAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &azblob.ClientOptions{
		Audience: "https://storage.azure.com/",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	azClientAudience, err := azblob.NewClient("https://"+accountName+".blob.core.windows.net/", cred, options)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	_, err = azClientAudience.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)

	defer func() {
		_, err = azClientAudience.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	pager := azClientAudience.NewListContainersPager(&azblob.ListContainersOptions{
		Prefix: &containerName,
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(len(resp.ContainerItems), 1)
		_require.NotNil(resp.ContainerItems[0].Name)
		_require.Equal(*resp.ContainerItems[0].Name, containerName)
	}
}

func (s *AZBlobRecordedTestsSuite) TestAzBlobClientCustomAudience() {
	_require := require.New(s.T())
	testName := s.T().Name()

	accountName, _ := testcommon.GetGenericAccountInfo(testcommon.TestAccountDefault)
	_require.Greater(len(accountName), 0)

	cred, err := testcommon.GetGenericTokenCredential()
	_require.NoError(err)

	options := &azblob.ClientOptions{
		Audience: "https://" + accountName + ".blob.core.windows.net",
	}
	testcommon.SetClientOptions(s.T(), &options.ClientOptions)
	azClientAudience, err := azblob.NewClient("https://"+accountName+".blob.core.windows.net/", cred, options)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	_, err = azClientAudience.CreateContainer(context.Background(), containerName, nil)
	_require.NoError(err)

	defer func() {
		_, err = azClientAudience.DeleteContainer(context.Background(), containerName, nil)
		_require.NoError(err)
	}()

	pager := azClientAudience.NewListContainersPager(&azblob.ListContainersOptions{
		Prefix: &containerName,
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.NoError(err)
		_require.Equal(len(resp.ContainerItems), 1)
		_require.NotNil(resp.ContainerItems[0].Name)
		_require.Equal(*resp.ContainerItems[0].Name, containerName)
	}
}

func (s *AZBlobRecordedTestsSuite) TestDownloadBufferWithLessData() {
	blobSize := 15
	blockSize := 1024
	concurrency := 3
	downloadOffset := 0
	downloadCount := 100
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, downloadOffset, downloadCount)
}

func (s *AZBlobRecordedTestsSuite) TestDownloadBufferWithLargeData() {
	blobSize := 10 * 1024 * 1024
	blockSize := 1024 * 1024
	concurrency := 3
	downloadOffset := 0
	downloadCount := 7 * 1024 * 1024
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, concurrency, downloadOffset, downloadCount)
}

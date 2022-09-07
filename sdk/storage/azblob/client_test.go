//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

import (
	"context"
	"errors"
	"io"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ctx = context.Background()

type AZBlobRecordedTestsSuite struct {
	suite.Suite
}

type AZBlobUnrecordedTestsSuite struct {
	suite.Suite
}

// Hookup to the testing framework
func Test(t *testing.T) {
	suite.Run(t, &AZBlobRecordedTestsSuite{})
	//suite.Run(t, &AZBlobUnrecordedTestsSuite{})
}

// nolint
func (s *AZBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

// nolint
func (s *AZBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

// nolint
func (s *AZBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

// create a test file
// nolint
func generateFile(fileName string, fileSize int) []byte {
	// generate random data
	_, bigBuff := testcommon.GenerateData(fileSize)

	// write to file and return the data
	_ = os.WriteFile(fileName, bigBuff, 0666)
	return bigBuff
}

// nolint
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
		&blockblob.UploadStreamOptions{BufferSize: bufferSize, MaxBuffers: maxBuffers})

	// Assert that upload was successful
	_require.Equal(err, nil)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)

	// Download the blob to verify
	downloadResponse, err := client.DownloadStream(ctx, containerName, blobName, nil)
	_require.Nil(err)

	// Assert that the content is correct
	actualBlobData, err := io.ReadAll(downloadResponse.Body)
	_require.Nil(err)
	_require.Equal(len(actualBlobData), blobSize)
	_require.EqualValues(actualBlobData, blobData)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobInChunks() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobSingleBuffer() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 1
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobSingleIO() {
	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobSingleIOEdgeCase() {
	blobSize := 8 * 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadStreamToBlockBlobEmpty() {
	blobSize := 0
	bufferSize := 8 * 1024
	maxBuffers := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(s.T(), _require, testName, blobSize, bufferSize, maxBuffers)
}

// nolint
func performUploadAndDownloadFileTest(t *testing.T, _require *require.Assertions, testName string, fileSize, blockSize, parallelism, downloadOffset, downloadCount int) {
	// Set up file to upload
	fileName := "BigFile.bin"
	fileData := generateFile(fileName, fileSize)

	// Open the file to upload
	file, err := os.Open(fileName)
	_require.Equal(err, nil)
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
	_, err = client.UploadFile(context.Background(), containerName, blobName, file,
		&blockblob.UploadFileOptions{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_require.Equal(bytesTransferred > 0 && bytesTransferred <= int64(fileSize), true)
			},
		})
	_require.NoError(err)
	//_require.Equal(response.StatusCode, 201)

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
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_require.Equal(bytesTransferred > 0 && bytesTransferred <= int64(fileSize), true)
			},
		})

	// Assert download was successful
	_require.NoError(err)

	// Assert downloaded data is consistent
	var destBuffer []byte
	if downloadCount == blob.CountToEnd {
		destBuffer = make([]byte, fileSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	n, err := destFile.Read(destBuffer)
	_require.Equal(err, nil)

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

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileInChunks() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileSingleIO() {
	fileSize := 1024
	blockSize := 2048
	parallelism := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileSingleRoutine() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 1
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileEmpty() {
	fileSize := 0
	blockSize := 1024
	parallelism := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileNonZeroOffset() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 0
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileNonZeroCount() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 0
	downloadCount := 6000
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadFileNonZeroOffsetAndCount() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 6000
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(s.T(), _require, testName, fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

// nolint
func performUploadAndDownloadBufferTest(t *testing.T, _require *require.Assertions, testName string, blobSize, blockSize, parallelism, downloadOffset, downloadCount int) {
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
	_, err = client.UploadBuffer(context.Background(), containerName, blobName, bytesToUpload,
		&blockblob.UploadBufferOptions{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_require.Equal(bytesTransferred > 0 && bytesTransferred <= int64(blobSize), true)
			},
		})
	_require.Equal(err, nil)
	//_require.Equal(response.StatusCode, 201)

	// Set up buffer to download the blob to
	var destBuffer []byte
	if downloadCount == blob.CountToEnd {
		destBuffer = make([]byte, blobSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	// Download the blob to a buffer
	_, err = client.DownloadBuffer(context.Background(),
		containerName,
		blobName,
		destBuffer, &blob.DownloadBufferOptions{
			Range: azblob.HTTPRange{
				Count:  int64(downloadCount),
				Offset: int64(downloadOffset),
			},
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_require.Equal(bytesTransferred > 0 && bytesTransferred <= int64(blobSize), true)
			},
		})

	_require.Equal(err, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		_require.EqualValues(destBuffer, bytesToUpload)
	} else {
		if downloadCount == 0 {
			_require.EqualValues(destBuffer, bytesToUpload[downloadOffset:])
		} else {
			_require.EqualValues(destBuffer, bytesToUpload[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferInChunks() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferSingleIO() {
	blobSize := 1024
	blockSize := 8 * 1024
	parallelism := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferSingleRoutine() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 1
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestUploadAndDownloadBufferEmpty() {
	blobSize := 0
	blockSize := 1024
	parallelism := 3
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, 0, 0)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestDownloadBufferWithNonZeroOffset() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 0
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestDownloadBufferWithNonZeroCount() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 0
	downloadCount := 6000
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestDownloadBufferWithNonZeroOffsetAndCount() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 2000
	downloadCount := 6 * 1024
	_require := require.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(s.T(), _require, testName, blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestBasicDoBatchTransfer() {
	_require := require.New(s.T())
	// test the basic multi-routine processing
	type testInstance struct {
		transferSize int64
		chunkSize    int64
		parallelism  uint16
		expectError  bool
	}

	testMatrix := []testInstance{
		{transferSize: 100, chunkSize: 10, parallelism: 5, expectError: false},
		{transferSize: 100, chunkSize: 9, parallelism: 4, expectError: false},
		{transferSize: 100, chunkSize: 8, parallelism: 15, expectError: false},
		{transferSize: 100, chunkSize: 1, parallelism: 3, expectError: false},
		{transferSize: 0, chunkSize: 100, parallelism: 5, expectError: false}, // empty file works
		{transferSize: 100, chunkSize: 0, parallelism: 5, expectError: true},  // 0 chunk size on the other hand must fail
		{transferSize: 0, chunkSize: 0, parallelism: 5, expectError: true},
	}

	for _, test := range testMatrix {
		ctx := context.Background()
		// maintain some counts to make sure the right number of chunks were queued, and the total size is correct
		totalSizeCount := int64(0)
		runCount := int64(0)

		err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
			TransferSize: test.transferSize,
			ChunkSize:    test.chunkSize,
			Parallelism:  test.parallelism,
			Operation: func(offset int64, chunkSize int64, ctx context.Context) error {
				atomic.AddInt64(&totalSizeCount, chunkSize)
				atomic.AddInt64(&runCount, 1)
				return nil
			},
			OperationName: "TestHappyPath",
		})

		if test.expectError {
			_require.NotNil(err)
		} else {
			_require.Nil(err)
			_require.Equal(totalSizeCount, test.transferSize)
			_require.Equal(runCount, ((test.transferSize-1)/test.chunkSize)+1)
		}
	}
}

// mock a memory mapped file (low-quality mock, meant to simulate the scenario only)
// nolint
type mockMMF struct {
	isClosed   bool
	failHandle *require.Assertions
}

// accept input
// nolint
func (m *mockMMF) write(_ string) {
	if m.isClosed {
		// simulate panic
		m.failHandle.Fail("")
	}
}

// nolint
func (s *AZBlobUnrecordedTestsSuite) TestDoBatchTransferWithError() {
	_require := require.New(s.T())
	ctx := context.Background()
	mmf := mockMMF{failHandle: _require}
	expectedFirstError := errors.New("#3 means trouble")

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		TransferSize: 5,
		ChunkSize:    1,
		Parallelism:  5,
		Operation: func(offset int64, chunkSize int64, ctx context.Context) error {
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

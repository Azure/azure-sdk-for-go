// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	testframework "github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"
)

// create a test file
//nolint
func generateFile(fileName string, fileSize int) []byte {
	// generate random data
	_, bigBuff := generateData(fileSize)

	// write to file and return the data
	_ = ioutil.WriteFile(fileName, bigBuff, 0666)
	return bigBuff
}

//nolint
func performUploadStreamToBlockBlobTest(_assert *assert.Assertions, testName string, blobSize, bufferSize, maxBuffers int) {
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(_assert, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	// Set up test blob
	blobName := generateBlobName(testName)
	blobClient := getBlockBlobClient(blobName, containerClient)

	// Create some data to test the upload stream
	blobContentReader, blobData := generateData(blobSize)

	// Perform UploadStreamToBlockBlob
	uploadResp, err := blobClient.UploadStreamToBlockBlob(ctx, blobContentReader,
		UploadStreamToBlockBlobOptions{BufferSize: bufferSize, MaxBuffers: maxBuffers})

	// Assert that upload was successful
	_assert.Equal(err, nil)
	_assert.Equal(uploadResp.RawResponse.StatusCode, 201)

	// Download the blob to verify
	downloadResponse, err := blobClient.Download(ctx, nil)
	_assert.Nil(err)

	// Assert that the content is correct
	actualBlobData, err := ioutil.ReadAll(downloadResponse.RawResponse.Body)
	_assert.Nil(err)
	_assert.Equal(len(actualBlobData), blobSize)
	_assert.EqualValues(actualBlobData, blobData)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlockBlobInChunks() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(_assert, testName, blobSize, bufferSize, maxBuffers)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlockBlobSingleBuffer() {
	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 1
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(_assert, testName, blobSize, bufferSize, maxBuffers)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlockBlobSingleIO() {
	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(_assert, testName, blobSize, bufferSize, maxBuffers)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlockBlobSingleIOEdgeCase() {
	blobSize := 8 * 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(_assert, testName, blobSize, bufferSize, maxBuffers)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadStreamToBlockBlobEmpty() {
	blobSize := 0
	bufferSize := 8 * 1024
	maxBuffers := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadStreamToBlockBlobTest(_assert, testName, blobSize, bufferSize, maxBuffers)
}

//nolint
func performUploadAndDownloadFileTest(_assert *assert.Assertions, testName string, fileSize, blockSize, parallelism, downloadOffset, downloadCount int) {
	// Set up file to upload
	fileName := "BigFile.bin"
	fileData := generateFile(fileName, fileSize)

	// Open the file to upload
	file, err := os.Open(fileName)
	_assert.Equal(err, nil)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	defer func(name string) {
		_ = os.Remove(name)
	}(fileName)

	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}

	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	// Set up test blob
	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	// Upload the file to a block blob
	response, err := bbClient.UploadFileToBlockBlob(context.Background(), file,
		HighLevelUploadToBlockBlobOption{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_assert.Equal(bytesTransferred > 0 && bytesTransferred <= int64(fileSize), true)
			},
		})
	_assert.Equal(err, nil)
	_assert.Equal(response.StatusCode, 201)

	// Set up file to download the blob to
	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	_assert.Equal(err, nil)
	defer func(destFile *os.File) {
		_ = destFile.Close()

	}(destFile)
	defer func(name string) {
		_ = os.Remove(name)

	}(destFileName)

	// Perform download
	err = bbClient.DownloadBlobToFile(context.Background(), int64(downloadOffset), int64(downloadCount),
		destFile,
		HighLevelDownloadFromBlobOptions{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_assert.Equal(bytesTransferred > 0 && bytesTransferred <= int64(fileSize), true)
			},
		})

	// Assert download was successful
	_assert.Equal(err, nil)

	// Assert downloaded data is consistent
	var destBuffer []byte
	if downloadCount == CountToEnd {
		destBuffer = make([]byte, fileSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	n, err := destFile.Read(destBuffer)
	_assert.Equal(err, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		_assert.EqualValues(destBuffer, fileData)
	} else {
		if downloadCount == 0 {
			_assert.Equal(n, fileSize-downloadOffset)
			_assert.EqualValues(destBuffer, fileData[downloadOffset:])
		} else {
			_assert.Equal(n, downloadCount)
			_assert.EqualValues(destBuffer, fileData[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileInChunks() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileSingleIO() {
	fileSize := 1024
	blockSize := 2048
	parallelism := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileSingleRoutine() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 1
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileEmpty() {
	fileSize := 0
	blockSize := 1024
	parallelism := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileNonZeroOffset() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 0
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileNonZeroCount() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 0
	downloadCount := 6000
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadFileNonZeroOffsetAndCount() {
	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 6000
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadFileTest(_assert, testName, fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

//nolint
func performUploadAndDownloadBufferTest(_assert *assert.Assertions, testName string, blobSize, blockSize, parallelism, downloadOffset, downloadCount int) {
	// Set up buffer to upload
	_, bytesToUpload := generateData(blobSize)

	// Set up test container
	_context := getTestContext(testName)
	var recording *testframework.Recording
	if _context != nil {
		recording = _context.recording
	}
	svcClient, err := getServiceClient(recording, testAccountDefault, nil)
	if err != nil {
		_assert.Fail("Unable to fetch service client because " + err.Error())
	}
	containerClient := createNewContainer(_assert, generateContainerName(testName), svcClient)
	defer deleteContainer(_assert, containerClient)

	// Set up test blob
	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	response, err := bbClient.UploadBufferToBlockBlob(context.Background(), bytesToUpload,
		HighLevelUploadToBlockBlobOption{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_assert.Equal(bytesTransferred > 0 && bytesTransferred <= int64(blobSize), true)
			},
		})
	_assert.Equal(err, nil)
	_assert.Equal(response.StatusCode, 201)

	// Set up buffer to download the blob to
	var destBuffer []byte
	if downloadCount == CountToEnd {
		destBuffer = make([]byte, blobSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	// Download the blob to a buffer
	err = bbClient.DownloadBlobToBuffer(context.Background(), int64(downloadOffset), int64(downloadCount),
		destBuffer, HighLevelDownloadFromBlobOptions{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				_assert.Equal(bytesTransferred > 0 && bytesTransferred <= int64(blobSize), true)
			},
		})

	_assert.Equal(err, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		_assert.EqualValues(destBuffer, bytesToUpload)
	} else {
		if downloadCount == 0 {
			_assert.EqualValues(destBuffer, bytesToUpload[downloadOffset:])
		} else {
			_assert.EqualValues(destBuffer, bytesToUpload[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadBufferInChunks() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadBufferSingleIO() {
	blobSize := 1024
	blockSize := 8 * 1024
	parallelism := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadBufferSingleRoutine() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 1
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestUploadAndDownloadBufferEmpty() {
	blobSize := 0
	blockSize := 1024
	parallelism := 3
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, 0, 0)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestDownloadBufferWithNonZeroOffset() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 0
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestDownloadBufferWithNonZeroCount() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 0
	downloadCount := 6000
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestDownloadBufferWithNonZeroOffsetAndCount() {
	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 2000
	downloadCount := 6 * 1024
	_assert := assert.New(s.T())
	testName := s.T().Name()
	performUploadAndDownloadBufferTest(_assert, testName, blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

//nolint
func (s *azblobUnrecordedTestSuite) TestBasicDoBatchTransfer() {
	_assert := assert.New(s.T())
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

		err := DoBatchTransfer(ctx, BatchTransferOptions{
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
			_assert.NotNil(err)
		} else {
			_assert.Nil(err)
			_assert.Equal(totalSizeCount, test.transferSize)
			_assert.Equal(runCount, ((test.transferSize-1)/test.chunkSize)+1)
		}
	}
}

// mock a memory mapped file (low-quality mock, meant to simulate the scenario only)
//nolint
type mockMMF struct {
	isClosed   bool
	failHandle *assert.Assertions
}

// accept input
//nolint
func (m *mockMMF) write(_ string) {
	if m.isClosed {
		// simulate panic
		m.failHandle.Fail("")
	}
}

//nolint
func (s *azblobUnrecordedTestSuite) TestDoBatchTransferWithError() {
	_assert := assert.New(s.T())
	ctx := context.Background()
	mmf := mockMMF{failHandle: _assert}
	expectedFirstError := errors.New("#3 means trouble")

	err := DoBatchTransfer(ctx, BatchTransferOptions{
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
				_assert.Equal(ctxErr, context.Canceled)
				return ctxErr
			}

			// anything before offset=3 should be done without problem
			return nil
		},
		OperationName: "TestErrorPath",
	})

	_assert.Equal(err, expectedFirstError)

	// simulate closing the mmf and make sure no panic occurs (as reported in #139)
	mmf.isClosed = true
	time.Sleep(time.Second * 5)
}

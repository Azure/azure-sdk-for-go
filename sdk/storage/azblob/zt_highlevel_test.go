// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// create a test file
func generateFile(fileName string, fileSize int) []byte {
	// generate random data
	_, bigBuff := generateData(fileSize)

	// write to file and return the data
	_ = ioutil.WriteFile(fileName, bigBuff, 0666)
	return bigBuff
}

func performUploadStreamToBlockBlobTest(t *testing.T, testName string, blobSize, bufferSize, maxBuffers int) {
	recording.LiveOnly(t)
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(t, containerClient)

	// Set up test blob
	blobName := generateBlobName(testName)
	blobClient := getBlockBlobClient(blobName, containerClient)

	// Create some data to test the upload stream
	blobContentReader, blobData := generateData(blobSize)

	// Perform UploadStreamToBlockBlob
	uploadResp, err := blobClient.UploadStreamToBlockBlob(ctx, blobContentReader,
		UploadStreamToBlockBlobOptions{BufferSize: bufferSize, MaxBuffers: maxBuffers})

	// Assert that upload was successful
	require.Equal(t, err, nil)
	require.Equal(t, uploadResp.RawResponse.StatusCode, 201)

	// Download the blob to verify
	downloadResponse, err := blobClient.Download(ctx, nil)
	require.NoError(t, err)

	// Assert that the content is correct
	actualBlobData, err := ioutil.ReadAll(downloadResponse.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, len(actualBlobData), blobSize)
	require.EqualValues(t, actualBlobData, blobData)
}

func TestUploadStreamToBlockBlobInChunks(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(t, t.Name(), blobSize, bufferSize, maxBuffers)
}

func TestUploadStreamToBlockBlobSingleBuffer(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	bufferSize := 1024
	maxBuffers := 1
	performUploadStreamToBlockBlobTest(t, t.Name(), blobSize, bufferSize, maxBuffers)
}

func TestUploadStreamToBlockBlobSingleIO(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(t, t.Name(), blobSize, bufferSize, maxBuffers)
}

func TestUploadStreamToBlockBlobSingleIOEdgeCase(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	bufferSize := 8 * 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(t, t.Name(), blobSize, bufferSize, maxBuffers)
}

func TestUploadStreamToBlockBlobEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 0
	bufferSize := 8 * 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(t, t.Name(), blobSize, bufferSize, maxBuffers)
}

func performUploadAndDownloadFileTest(t *testing.T, testName string, fileSize, blockSize, parallelism, downloadOffset, downloadCount int) {
	recording.LiveOnly(t)
	// Set up file to upload
	fileName := "BigFile.bin"
	fileData := generateFile(fileName, fileSize)

	// Open the file to upload
	file, err := os.Open(fileName)
	require.NoError(t, err)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	defer func(name string) {
		_ = os.Remove(name)
	}(fileName)

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	// Set up test blob
	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	// Upload the file to a block blob
	response, err := bbClient.UploadFileToBlockBlob(context.Background(), file,
		HighLevelUploadToBlockBlobOption{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				require.Equal(t, bytesTransferred > 0 && bytesTransferred <= int64(fileSize), true)
			},
		})
	require.Equal(t, err, nil)
	require.Equal(t, response.StatusCode, 201)

	// Set up file to download the blob to
	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	require.Equal(t, err, nil)
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
				require.Equal(t, bytesTransferred > 0 && bytesTransferred <= int64(fileSize), true)
			},
		})

	// Assert download was successful
	require.Equal(t, err, nil)

	// Assert downloaded data is consistent
	var destBuffer []byte
	if downloadCount == CountToEnd {
		destBuffer = make([]byte, fileSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	n, err := destFile.Read(destBuffer)
	require.Equal(t, err, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		require.EqualValues(t, destBuffer, fileData)
	} else {
		if downloadCount == 0 {
			require.Equal(t, n, fileSize-downloadOffset)
			require.EqualValues(t, destBuffer, fileData[downloadOffset:])
		} else {
			require.Equal(t, n, downloadCount)
			require.EqualValues(t, destBuffer, fileData[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

func TestUploadAndDownloadFileInChunks(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadFileSingleIO(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 1024
	blockSize := 2048
	parallelism := 3
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadFileSingleRoutine(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 1
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadFileEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 0
	blockSize := 1024
	parallelism := 3
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadFileNonZeroOffset(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 0
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

func TestUploadAndDownloadFileNonZeroCount(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 0
	downloadCount := 6000
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

func TestUploadAndDownloadFileNonZeroOffsetAndCount(t *testing.T) {
	stop := start(t)
	defer stop()

	fileSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 6000
	performUploadAndDownloadFileTest(t, t.Name(), fileSize, blockSize, parallelism, downloadOffset, downloadCount)
}

func performUploadAndDownloadBufferTest(t *testing.T, testName string, blobSize, blockSize, parallelism, downloadOffset, downloadCount int) {
	recording.LiveOnly(t)
	// Set up buffer to upload
	_, bytesToUpload := generateData(blobSize)

	// Set up test container
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)
	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)

	// Set up test blob
	bbClient := getBlockBlobClient(generateBlobName(testName), containerClient)

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	response, err := bbClient.UploadBufferToBlockBlob(context.Background(), bytesToUpload,
		HighLevelUploadToBlockBlobOption{
			BlockSize:   int64(blockSize),
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				require.Equal(t, bytesTransferred > 0 && bytesTransferred <= int64(blobSize), true)
			},
		})
	require.Equal(t, err, nil)
	require.Equal(t, response.StatusCode, 201)

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
				require.Equal(t, bytesTransferred > 0 && bytesTransferred <= int64(blobSize), true)
			},
		})

	require.Equal(t, err, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		require.EqualValues(t, destBuffer, bytesToUpload)
	} else {
		if downloadCount == 0 {
			require.EqualValues(t, destBuffer, bytesToUpload[downloadOffset:])
		} else {
			require.EqualValues(t, destBuffer, bytesToUpload[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

func TestUploadAndDownloadBufferInChunks(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadBufferSingleIO(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 1024
	blockSize := 8 * 1024
	parallelism := 3
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadBufferSingleRoutine(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 1
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, 0, 0)
}

func TestUploadAndDownloadBufferEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 0
	blockSize := 1024
	parallelism := 3
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, 0, 0)
}

func TestDownloadBufferWithNonZeroOffset(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 1000
	downloadCount := 0
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

func TestDownloadBufferWithNonZeroCount(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 0
	downloadCount := 6000
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

func TestDownloadBufferWithNonZeroOffsetAndCount(t *testing.T) {
	stop := start(t)
	defer stop()

	blobSize := 8 * 1024
	blockSize := 1024
	parallelism := 3
	downloadOffset := 2000
	downloadCount := 6 * 1024
	performUploadAndDownloadBufferTest(t, t.Name(), blobSize, blockSize, parallelism, downloadOffset, downloadCount)
}

func TestBasicDoBatchTransfer(t *testing.T) {
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
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, totalSizeCount, test.transferSize)
			require.Equal(t, runCount, ((test.transferSize-1)/test.chunkSize)+1)
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

func TestDoBatchTransferWithError(t *testing.T) {
	recording.LiveOnly(t)
	ctx := context.Background()
	mmf := mockMMF{failHandle: require.New(t)}
	expectedFirstError := errors.New("#3 means trouble")

	err := DoBatchTransfer(ctx, BatchTransferOptions{
		TransferSize: 5,
		ChunkSize:    1,
		Parallelism:  5,
		Operation: func(offset int64, chunkSize int64, ctx context.Context) error {
			// simulate doing some work (HTTP call in real scenarios)
			// later chunks later longer to finish
			recording.Sleep(time.Second * time.Duration(offset))
			// simulate having gotten data and write it to the memory mapped file
			mmf.write("input")

			// with one of the chunks, pretend like an error occurred (like the network connection breaks)
			if offset == 3 {
				return expectedFirstError
			} else if offset > 3 {
				// anything after offset=3 are canceled
				// so verify that the context indeed got canceled
				ctxErr := ctx.Err()
				require.Equal(t, ctxErr, context.Canceled)
				return ctxErr
			}

			// anything before offset=3 should be done without problem
			return nil
		},
		OperationName: "TestErrorPath",
	})

	require.Equal(t, err, expectedFirstError)

	// simulate closing the mmf and make sure no panic occurs (as reported in #139)
	mmf.isClosed = true
	recording.Sleep(time.Second * 5)
}

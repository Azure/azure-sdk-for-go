// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"

	chk "gopkg.in/check.v1"
)

// create a test file
func generateFile(fileName string, fileSize int64) []byte {
	// generate random data
	_, bigBuff := getRandomDataAndReader(fileSize)

	// write to file and return the data
	_ = ioutil.WriteFile(fileName, bigBuff, 0666)
	return bigBuff
}

func performUploadStreamToBlockBlobTest(c *chk.C, blobSize int64, bufferSize, maxBuffers int) {
	// Set up test container
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	// Set up test blob
	blobClient, _ := getBlockBlobClient(c, containerClient)

	// Create some data to test the upload stream
	blobContentReader, blobData := getRandomDataAndReader(blobSize)

	// Perform UploadStreamToBlockBlob
	uploadResp, err := UploadStreamToBlockBlob(ctx, blobContentReader, blobClient,
		UploadStreamToBlockBlobOptions{BufferSize: bufferSize, MaxBuffers: maxBuffers})

	// Assert that upload was successful
	c.Assert(err, chk.Equals, nil)
	c.Assert(uploadResp.RawResponse.StatusCode, chk.Equals, 201)

	// Download the blob to verify
	downloadResponse, err := blobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)

	// Assert that the content is correct
	actualBlobData, err := ioutil.ReadAll(downloadResponse.RawResponse.Body)
	c.Assert(err, chk.IsNil)
	c.Assert(int64(len(actualBlobData)), chk.Equals, blobSize)
	c.Assert(actualBlobData, chk.DeepEquals, blobData)
}

func (s *aztestsSuite) TestUploadStreamToBlockBlobInChunks(c *chk.C) {
	blobSize := int64(8 * 1024)
	bufferSize := 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(c, blobSize, bufferSize, maxBuffers)
}

func (s *aztestsSuite) TestUploadStreamToBlockBlobSingleBuffer(c *chk.C) {
	blobSize := int64(8 * 1024)
	bufferSize := 1024
	maxBuffers := 1
	performUploadStreamToBlockBlobTest(c, blobSize, bufferSize, maxBuffers)
}

func (s *aztestsSuite) TestUploadStreamToBlockBlobSingleIO(c *chk.C) {
	blobSize := int64(1024)
	bufferSize := 8 * 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(c, blobSize, bufferSize, maxBuffers)
}

func (s *aztestsSuite) TestUploadStreamToBlockBlobSingleIOEdgeCase(c *chk.C) {
	blobSize := int64(8 * 1024)
	bufferSize := 8 * 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(c, blobSize, bufferSize, maxBuffers)
}

func (s *aztestsSuite) TestUploadStreamToBlockBlobEmpty(c *chk.C) {
	blobSize := int64(0)
	bufferSize := 8 * 1024
	maxBuffers := 3
	performUploadStreamToBlockBlobTest(c, blobSize, bufferSize, maxBuffers)
}

func performUploadAndDownloadFileTest(c *chk.C, fileSize, blockSize, downloadOffset, downloadCount int64, parallelism int) {
	// Set up file to upload
	fileName := "BigFile.bin"
	fileData := generateFile(fileName, fileSize)

	// Open the file to upload
	file, err := os.Open(fileName)
	c.Assert(err, chk.Equals, nil)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	defer func(name string) {
		_ = os.Remove(name)
	}(fileName)

	// Set up test container
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	// Set up test blob
	blockBlobClient, _ := getBlockBlobClient(c, containerClient)

	// Upload the file to a block blob
	response, err := UploadFileToBlockBlob(context.Background(), file, blockBlobClient,
		HighLevelUploadToBlockBlobOption{
			BlockSize:   blockSize,
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				c.Assert(bytesTransferred > 0 && bytesTransferred <= fileSize, chk.Equals, true)
			},
		})
	c.Assert(err, chk.Equals, nil)
	c.Assert(response.StatusCode, chk.Equals, 201)

	// Set up file to download the blob to
	destFileName := "BigFile-downloaded.bin"
	destFile, err := os.Create(destFileName)
	c.Assert(err, chk.Equals, nil)
	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)
	defer func(name string) {
		_ = os.Remove(name)
	}(destFileName)

	// Perform download
	err = DownloadBlobToFile(context.Background(), blockBlobClient.BlobClient, downloadOffset, downloadCount,
		destFile,
		HighLevelDownloadFromBlobOptions{
			BlockSize:   blockSize,
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				c.Assert(bytesTransferred > 0 && bytesTransferred <= fileSize, chk.Equals, true)
			},
		})

	// Assert download was successful
	c.Assert(err, chk.Equals, nil)

	// Assert downloaded data is consistent
	var destBuffer []byte
	if downloadCount == CountToEnd {
		destBuffer = make([]byte, fileSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	n, err := destFile.Read(destBuffer)
	c.Assert(err, chk.Equals, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		c.Assert(destBuffer, chk.DeepEquals, fileData)
	} else {
		if downloadCount == 0 {
			c.Assert(int64(n), chk.Equals, fileSize-downloadOffset)
			c.Assert(destBuffer, chk.DeepEquals, fileData[downloadOffset:])
		} else {
			c.Assert(int64(n), chk.Equals, downloadCount)
			c.Assert(destBuffer, chk.DeepEquals, fileData[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

func (s *aztestsSuite) TestUploadAndDownloadFileInChunks(c *chk.C) {
	fileSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	performUploadAndDownloadFileTest(c, fileSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadFileSingleIO(c *chk.C) {
	fileSize := int64(1024)
	blockSize := int64(2048)
	parallelism := 3
	performUploadAndDownloadFileTest(c, fileSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadFileSingleRoutine(c *chk.C) {
	fileSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 1
	performUploadAndDownloadFileTest(c, fileSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadFileEmpty(c *chk.C) {
	fileSize := int64(0)
	blockSize := int64(1024)
	parallelism := 3
	performUploadAndDownloadFileTest(c, fileSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadFileNonZeroOffset(c *chk.C) {
	fileSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	downloadOffset := int64(1000)
	downloadCount := int64(0)
	performUploadAndDownloadFileTest(c, fileSize, blockSize, downloadOffset, downloadCount, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadFileNonZeroCount(c *chk.C) {
	fileSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	downloadOffset := int64(0)
	downloadCount := int64(6000)
	performUploadAndDownloadFileTest(c, fileSize, blockSize, downloadOffset, downloadCount, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadFileNonZeroOffsetAndCount(c *chk.C) {
	fileSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	downloadOffset := int64(1000)
	downloadCount := int64(6000)
	performUploadAndDownloadFileTest(c, fileSize, blockSize, downloadOffset, downloadCount, parallelism)
}

func performUploadAndDownloadBufferTest(c *chk.C, blobSize, blockSize, downloadOffset, downloadCount int64, parallelism int) {
	// Set up buffer to upload
	_, bytesToUpload := getRandomDataAndReader(blobSize)

	// Set up test container
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	// Set up test blob
	blockBlobClient, _ := getBlockBlobClient(c, containerClient)

	// Pass the Context, stream, stream size, block blob URL, and options to StreamToBlockBlob
	response, err := UploadBufferToBlockBlob(context.Background(), bytesToUpload, blockBlobClient,
		HighLevelUploadToBlockBlobOption{
			BlockSize:   blockSize,
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			Progress: func(bytesTransferred int64) {
				c.Assert(bytesTransferred > 0 && bytesTransferred <= blobSize, chk.Equals, true)
			},
		})
	c.Assert(err, chk.Equals, nil)
	c.Assert(response.StatusCode, chk.Equals, 201)

	// Set up buffer to download the blob to
	var destBuffer []byte
	if downloadCount == CountToEnd {
		destBuffer = make([]byte, blobSize-downloadOffset)
	} else {
		destBuffer = make([]byte, downloadCount)
	}

	// Download the blob to a buffer
	err = DownloadBlobToBuffer(context.Background(), blockBlobClient.BlobClient, downloadOffset, downloadCount,
		destBuffer, HighLevelDownloadFromBlobOptions{
			BlockSize:   blockSize,
			Parallelism: uint16(parallelism),
			// If Progress is non-nil, this function is called periodically as bytes are uploaded.
			//Progress: func(bytesTransferred int64) {
			//	c.Assert(bytesTransferred > 0 && bytesTransferred <= int64(blobSize), chk.Equals, true)
			//},
		})

	c.Assert(err, chk.Equals, nil)

	if downloadOffset == 0 && downloadCount == 0 {
		c.Assert(destBuffer, chk.DeepEquals, bytesToUpload)
	} else {
		if downloadCount == 0 {
			c.Assert(destBuffer, chk.DeepEquals, bytesToUpload[downloadOffset:])
		} else {
			c.Assert(destBuffer, chk.DeepEquals, bytesToUpload[downloadOffset:downloadOffset+downloadCount])
		}
	}
}

func (s *aztestsSuite) TestUploadAndDownloadBufferInChunks(c *chk.C) {
	blobSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadBufferSingleIO(c *chk.C) {
	blobSize := int64(1024)
	blockSize := int64(8 * 1024)
	parallelism := 3
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadBufferSingleRoutine(c *chk.C) {
	blobSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 1
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestUploadAndDownloadBufferEmpty(c *chk.C) {
	blobSize := int64(0)
	blockSize := int64(1024)
	parallelism := 3
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, 0, 0, parallelism)
}

func (s *aztestsSuite) TestDownloadBufferWithNonZeroOffset(c *chk.C) {
	blobSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	downloadOffset := int64(1000)
	downloadCount := int64(0)
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, downloadOffset, downloadCount, parallelism)
}

func (s *aztestsSuite) TestDownloadBufferWithNonZeroCount(c *chk.C) {
	blobSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	downloadOffset := int64(0)
	downloadCount := int64(6000)
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, downloadOffset, downloadCount, parallelism)
}

func (s *aztestsSuite) TestDownloadBufferWithNonZeroOffsetAndCount(c *chk.C) {
	blobSize := int64(8 * 1024)
	blockSize := int64(1024)
	parallelism := 3
	downloadOffset := int64(2000)
	downloadCount := int64(6 * 1024)
	performUploadAndDownloadBufferTest(c, blobSize, blockSize, downloadOffset, downloadCount, parallelism)
}

func (s *aztestsSuite) TestBasicDoBatchTransfer(c *chk.C) {
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
			c.Assert(err, chk.NotNil)
		} else {
			c.Assert(err, chk.IsNil)
			c.Assert(totalSizeCount, chk.Equals, test.transferSize)
			c.Assert(runCount, chk.Equals, ((test.transferSize-1)/test.chunkSize)+1)
		}
	}
}

// mock a memory mapped file (low-quality mock, meant to simulate the scenario only)
type mockMMF struct {
	isClosed   bool
	failHandle *chk.C
}

// accept input
func (m *mockMMF) write(_ string) {
	if m.isClosed {
		// simulate panic
		m.failHandle.Fail()
	}
}

func (s *aztestsSuite) TestDoBatchTransferWithError(c *chk.C) {
	ctx := context.Background()
	mmf := mockMMF{failHandle: c}
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
				c.Assert(ctxErr, chk.Equals, context.Canceled)
				return ctxErr
			}

			// anything before offset=3 should be done without problem
			return nil
		},
		OperationName: "TestErrorPath",
	})

	c.Assert(err, chk.Equals, expectedFirstError)

	// simulate closing the mmf and make sure no panic occurs (as reported in #139)
	mmf.isClosed = true
	time.Sleep(time.Second * 5)
}

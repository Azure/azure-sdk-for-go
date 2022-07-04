//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

const defaultDownloadBlockSize = int64(4 * 1024 * 1024) // 4MB

// NewWriteAtBuffer creates a new WriteAtBuffer with the specified buffer.
func NewWriteAtBuffer(b []byte) *WriteAtBuffer {
	return &WriteAtBuffer{
		buffer: b,
		mu:     &sync.Mutex{},
	}
}

// WriteAtBuffer is an in-memory buffer that supports the io.WriteAt interface.
// Safe for concurrent use.
type WriteAtBuffer struct {
	buffer []byte
	mu     *sync.Mutex
}

// Bytes returns the slice of bytes written to the buffer.
func (w *WriteAtBuffer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer
}

// WriteAt implements the io.WriteAt interface for the WriteAtBuffer type.
// If the underlying buffer is too small, it will grow the buffer by the required amount.
func (w *WriteAtBuffer) WriteAt(data []byte, off int64) (int, error) {
	dataLen := len(data)
	neededLen := off + int64(dataLen)

	w.mu.Lock()
	defer w.mu.Unlock()

	// check if we need to grow the buffer
	if int64(len(w.buffer)) < neededLen {
		// length is sufficient, check capacity
		if int64(cap(w.buffer)) < neededLen {
			// not enough capacity, expand
			newBuf := make([]byte, neededLen)
			copy(newBuf, w.buffer)
			w.buffer = newBuf
		}
		w.buffer = w.buffer[:neededLen]
	}

	copy(w.buffer[off:], data)
	return dataLen, nil
}

type DownloadToWriterOptions struct {
	// BlockSize specifies the block size to use for each parallel download; the default size is BlobDefaultDownloadBlockSize.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are received.
	Progress func(bytesTransferred int64)

	// BlobAccessConditions indicates the access conditions used when making HTTP GET requests against the blob.
	BlobAccessConditions *blob.AccessConditions

	// ClientProvidedKeyOptions indicates the client provided key by name and/or by value to encrypt/decrypt data.
	CpkInfo      *blob.CpkInfo
	CpkScopeInfo *blob.CpkScopeInfo

	// Parallelism indicates the maximum number of blocks to download in parallel (0=default)
	Parallelism uint16

	// RetryReaderOptionsPerBlock is used when downloading each block.
	RetryReaderOptionsPerBlock blob.RetryReaderOptions
}

func DownloadToWriterAt(ctx context.Context, blob blob.DownloadResponse, writer io.WriterAt, o *DownloadToWriterOptions) error {
	o = shared.CopyOptions(o)
	if o.BlockSize == 0 {
		o.BlockSize = defaultDownloadBlockSize
	}

	if *blob.ContentLength <= 0 {
		// The file is empty, there is nothing to download.
		return nil
	}

	// Prepare and do parallel download.
	progress := int64(0)
	progressLock := &sync.Mutex{}

	err := DoBatchTransfer(ctx, o.BlockSize, *blob.ContentLength, func(ctx context.Context, chunkStart int64, count int64) error {
		body := blob.NewRetryReader(ctx, &o.RetryReaderOptionsPerBlock)
		if o.Progress != nil {
			rangeProgress := int64(0)
			body = streaming.NewResponseProgress(
				body,
				func(bytesTransferred int64) {
					diff := bytesTransferred - rangeProgress
					rangeProgress = bytesTransferred
					progressLock.Lock()
					progress += diff
					o.Progress(progress)
					progressLock.Unlock()
				})
		}
		_, err := io.Copy(newSectionWriter(writer, chunkStart, count), body)
		if err != nil {
			return err
		}
		err = body.Close()
		return err
	}, &BatchTransferOptions{
		Parallelism: o.Parallelism,
	})
	if err != nil {
		return err
	}
	return nil
}

// BatchTransferOptions contains the optional values for DoBatchTransfer.
type BatchTransferOptions struct {
	Parallelism uint16
}

// DoBatchTransfer helps to execute operations in a batch manner.
// Can be used by users to customize batch works (for other scenarios that the SDK does not provide)
func DoBatchTransfer(ctx context.Context, chunkSize, transferSize int64, operation func(ctx context.Context, offset int64, chunkSize int64) error, o *BatchTransferOptions) error {
	if chunkSize == 0 {
		return errors.New("chunkSize cannot be 0")
	} else if transferSize == 0 {
		return errors.New("transferSize cannot be 0")
	}

	o = shared.CopyOptions(o)
	if o.Parallelism == 0 {
		o.Parallelism = 5 // default Parallelism
	}

	// Prepare and do parallel operations.
	numChunks := uint16(((transferSize - 1) / chunkSize) + 1)
	operationChannel := make(chan func() error, o.Parallelism) // Create the channel that release 'Parallelism' goroutines concurrently
	operationResponseChannel := make(chan error, numChunks)    // Holds each response
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create the goroutines that process each operation (in parallel).
	for g := uint16(0); g < o.Parallelism; g++ {
		//grIndex := g
		go func() {
			for f := range operationChannel {
				err := f()
				operationResponseChannel <- err
			}
		}()
	}

	// Add each chunk's operation to the channel.
	for chunkNum := uint16(0); chunkNum < numChunks; chunkNum++ {
		curChunkSize := chunkSize

		if chunkNum == numChunks-1 { // Last chunk
			curChunkSize = transferSize - (int64(chunkNum) * chunkSize) // Remove size of all transferred chunks from total
		}
		offset := int64(chunkNum) * chunkSize

		operationChannel <- func() error {
			return operation(ctx, offset, curChunkSize)
		}
	}
	close(operationChannel)

	// Wait for the operations to complete.
	var firstErr error = nil
	for chunkNum := uint16(0); chunkNum < numChunks; chunkNum++ {
		responseError := <-operationResponseChannel
		// record the first error (the original error which should cause the other chunks to fail with canceled context)
		if responseError != nil && firstErr == nil {
			cancel() // As soon as any operation fails, cancel all remaining operation calls
			firstErr = responseError
		}
	}
	return firstErr
}

type sectionWriter struct {
	count    int64
	offset   int64
	position int64
	writerAt io.WriterAt
}

func newSectionWriter(c io.WriterAt, off int64, count int64) *sectionWriter {
	return &sectionWriter{
		count:    count,
		offset:   off,
		writerAt: c,
	}
}

func (c *sectionWriter) Write(p []byte) (int, error) {
	remaining := c.count - c.position

	if remaining <= 0 {
		return 0, errors.New("end of section reached")
	}

	slice := p

	if int64(len(slice)) > remaining {
		slice = slice[:remaining]
	}

	n, err := c.writerAt.WriteAt(slice, c.offset+c.position)
	c.position += int64(n)
	if err != nil {
		return n, err
	}

	if len(p) > n {
		return n, errors.New("not enough space for all bytes")
	}

	return n, nil
}

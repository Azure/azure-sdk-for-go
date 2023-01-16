//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeBlockBlob struct {
	uncommittedBlocksMu sync.Mutex
	UncommittedBlocks   map[string][]byte
	CommittedBlocks     []byte

	ErrOnBlockID int           // when > -1, StageBlock returns an error on a blockID >= ErrOnBlockID
	StageDelay   time.Duration // when > 0, StageBlock sleeps for the specified duration before returning
}

func newFakeBlockBlob() *fakeBlockBlob {
	return &fakeBlockBlob{
		UncommittedBlocks: map[string][]byte{},
		ErrOnBlockID:      -1,
	}
}

func (f *fakeBlockBlob) StageBlock(ctx context.Context, blockID string, body io.ReadSeekCloser, _ *StageBlockOptions) (StageBlockResponse, error) {
	if err := ctx.Err(); err != nil {
		return StageBlockResponse{}, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return StageBlockResponse{}, err
	}

	var blockCount int
	f.uncommittedBlocksMu.Lock()
	if _, ok := f.UncommittedBlocks[blockID]; ok {
		f.uncommittedBlocksMu.Unlock()
		return StageBlockResponse{}, fmt.Errorf("duplicate block ID %s", blockID)
	}
	f.UncommittedBlocks[blockID] = data
	blockCount = len(f.UncommittedBlocks)
	f.uncommittedBlocksMu.Unlock()

	if f.StageDelay > 0 {
		select {
		case <-ctx.Done():
			return StageBlockResponse{}, ctx.Err()
		case <-time.After(f.StageDelay):
			// delay elapsed
		}
	}

	if f.ErrOnBlockID > -1 && blockCount >= f.ErrOnBlockID {
		return StageBlockResponse{}, io.ErrNoProgress
	}

	return StageBlockResponse{}, nil
}

func (f *fakeBlockBlob) CommitBlockList(ctx context.Context, base64BlockIDs []string, _ *CommitBlockListOptions) (CommitBlockListResponse, error) {
	if err := ctx.Err(); err != nil {
		return CommitBlockListResponse{}, err
	}

	for _, blockID := range base64BlockIDs {
		toCommit, ok := f.UncommittedBlocks[blockID]
		if !ok {
			return CommitBlockListResponse{}, fmt.Errorf("didn't find block with ID %s", blockID)
		}

		f.CommittedBlocks = append(f.CommittedBlocks, toCommit...)
		delete(f.UncommittedBlocks, blockID)
	}

	return CommitBlockListResponse{}, nil
}

func (f *fakeBlockBlob) Upload(ctx context.Context, body io.ReadSeekCloser, _ *UploadOptions) (UploadResponse, error) {
	if err := ctx.Err(); err != nil {
		return UploadResponse{}, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return UploadResponse{}, err
	}

	f.CommittedBlocks = data
	return UploadResponse{}, nil
}

func calcMD5(data []byte) string {
	h := md5.New()
	if _, err := io.Copy(h, bytes.NewReader(data)); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// used to track proper acquisition and closing of buffers
type bufMgrTracker struct {
	inner bufferManager[mmb]

	Count int  // total count of allocated buffers
	Freed bool // buffers were freed
}

func newBufMgrTracker(maxBuffers int, bufferSize int64) *bufMgrTracker {
	return &bufMgrTracker{
		inner: newMMBPool(maxBuffers, bufferSize),
	}
}

func (pool *bufMgrTracker) Acquire() <-chan mmb {
	return pool.inner.Acquire()
}

func (pool *bufMgrTracker) Grow() (int, error) {
	n, err := pool.inner.Grow()
	if err != nil {
		return 0, err
	}
	pool.Count = n
	return n, nil
}

func (pool *bufMgrTracker) Release(buffer mmb) {
	pool.inner.Release(buffer)
}

func (pool *bufMgrTracker) Free() {
	pool.inner.Free()
	pool.Freed = true
}

func TestSlowDestCopyFrom(t *testing.T) {
	bigSrc := make([]byte, _1MiB+500*1024) //This should cause 2 reads

	fakeBB := newFakeBlockBlob()

	fakeBB.StageDelay = 200 * time.Millisecond
	fakeBB.ErrOnBlockID = 0

	var tracker *bufMgrTracker

	errs := make(chan error, 1)
	go func() {
		_, err := copyFromReader(context.Background(), bytes.NewReader(bigSrc), fakeBB, UploadStreamOptions{}, func(maxBuffers int, bufferSize int64) bufferManager[mmb] {
			tracker = newBufMgrTracker(maxBuffers, bufferSize)
			return tracker
		})
		errs <- err
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		failMsg := "TestSlowDestCopyFrom(slow writes shouldn't cause deadlock) failed: Context expired, copy deadlocked"
		t.Fatal(failMsg)
	case err := <-errs:
		require.ErrorIs(t, err, io.ErrNoProgress)
	}

	require.Equal(t, 1, tracker.Count)
	require.True(t, tracker.Freed)
}

func TestCopyFromReader(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		desc      string
		ctx       context.Context
		o         UploadStreamOptions
		fileSize  int
		uploadErr bool
		err       bool
	}{
		{
			desc: "context was cancelled",
			ctx:  canceled,
			err:  true,
		},
		{
			desc:     "Send file(0 KiB) with default UploadStreamOptions",
			ctx:      context.Background(),
			fileSize: 0,
		},
		{
			desc:     "Send file(10 KiB) with default UploadStreamOptions",
			ctx:      context.Background(),
			fileSize: 10 * 1024,
		},
		{
			desc:     "Send file(10 KiB) with default UploadStreamOptions set to azcopy settings",
			ctx:      context.Background(),
			fileSize: 10 * 1024,
			o:        UploadStreamOptions{Concurrency: 5, BlockSize: 8 * 1024 * 1024},
		},
		{
			desc:     "Send file(1 MiB) with default UploadStreamOptions",
			ctx:      context.Background(),
			fileSize: _1MiB,
		},
		{
			desc:     "Send file(1 MiB) with default UploadStreamOptions set to azcopy settings",
			ctx:      context.Background(),
			fileSize: _1MiB,
			o:        UploadStreamOptions{Concurrency: 5, BlockSize: 8 * 1024 * 1024},
		},
		{
			desc:     "Send file(1.5 MiB) with default UploadStreamOptions",
			ctx:      context.Background(),
			fileSize: _1MiB + 500*1024,
		},
		{
			desc:     "Send file(1.5 MiB) with 2 writers",
			ctx:      context.Background(),
			fileSize: _1MiB + 500*1024 + 1,
			o:        UploadStreamOptions{Concurrency: 2},
		},
		{
			desc:      "Send file(12 MiB) with 3 writers and 1 MiB buffer and a write error",
			ctx:       context.Background(),
			fileSize:  12 * _1MiB,
			o:         UploadStreamOptions{Concurrency: 2, BlockSize: _1MiB},
			uploadErr: true,
			err:       true,
		},
		{
			desc:     "Send file(12 MiB) with 3 writers and 1.5 MiB buffer",
			ctx:      context.Background(),
			fileSize: 12 * _1MiB,
			o:        UploadStreamOptions{Concurrency: 2, BlockSize: _1MiB + .5*_1MiB},
		},
		{
			desc:     "Send file(12 MiB) with default UploadStreamOptions set to azcopy settings",
			ctx:      context.Background(),
			fileSize: 12 * _1MiB,
			o:        UploadStreamOptions{Concurrency: 5, BlockSize: 8 * 1024 * 1024},
		},
	}

	for _, test := range tests {
		from := make([]byte, test.fileSize)
		fakeBB := newFakeBlockBlob()

		if test.uploadErr {
			fakeBB.ErrOnBlockID = 1
		}

		var tracker *bufMgrTracker

		_, err := copyFromReader(test.ctx, bytes.NewReader(from), fakeBB, test.o, func(maxBuffers int, bufferSize int64) bufferManager[mmb] {
			tracker = newBufMgrTracker(maxBuffers, bufferSize)
			return tracker
		})

		// assert that at least one buffer was allocated
		assert.Greater(t, tracker.Count, 0, test.desc)
		require.True(t, tracker.Freed)

		switch {
		case err == nil && test.err:
			t.Errorf("TestCopyFromReader(%s): got err == nil, want err != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestCopyFromReader(%s): got err == %s, want err == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		want := calcMD5(from)
		got := calcMD5(fakeBB.CommittedBlocks)

		if got != want {
			t.Errorf("TestCopyFromReader(%s): MD5 not the same: got %s, want %s", test.desc, got, want)
		}
	}
}

type readFailer struct {
	reader io.Reader

	reads  int
	failOn int
}

func (r *readFailer) Read(b []byte) (int, error) {
	r.reads++
	if r.reads == r.failOn {
		return 0, io.ErrNoProgress
	}
	return r.reader.Read(b)
}

func TestCopyFromReaderReadError(t *testing.T) {
	fakeBB := newFakeBlockBlob()
	var tracker *bufMgrTracker

	rf := readFailer{
		reader: bytes.NewReader(make([]byte, 5*_1MiB)),
		failOn: 2,
	}
	_, err := copyFromReader(context.Background(), &rf, fakeBB, UploadStreamOptions{}, func(maxBuffers int, bufferSize int64) bufferManager[mmb] {
		tracker = newBufMgrTracker(maxBuffers, bufferSize)
		return tracker
	})

	require.ErrorIs(t, err, io.ErrNoProgress)
	assert.Greater(t, tracker.Count, 0)
	require.True(t, tracker.Freed)
}

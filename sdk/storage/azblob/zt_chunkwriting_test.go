//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
//import (
//	"context"
//	"crypto/md5"
//	"errors"
//	"fmt"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
//	"io"
//	"math/rand"
//	"os"
//	"path/filepath"
//	"strings"
//	"sync/atomic"
//	"time"
//
//	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
//)
//
//const finalFileName = "final"
//
////nolint
//type fakeBlockWriter struct {
//	path       string
//	block      int32
//	errOnBlock int32
//}
//
////nolint
//func newFakeBlockWriter() *fakeBlockWriter {
//	generatedUuid, _ := uuid.New()
//	f := &fakeBlockWriter{
//		path:       filepath.Join(os.TempDir(), generatedUuid.String()),
//		block:      -1,
//		errOnBlock: -1,
//	}
//
//	if err := os.MkdirAll(f.path, 0700); err != nil {
//		panic(err)
//	}
//
//	return f
//}
//
////nolint
//func (f *fakeBlockWriter) StageBlock(_ context.Context, blockID string, body io.ReadSeekCloser, _ *BlockBlobStageBlockOptions) (BlockBlobStageBlockResponse, error) {
//	n := atomic.AddInt32(&f.block, 1)
//	if n == f.errOnBlock {
//		return BlockBlobStageBlockResponse{}, io.ErrNoProgress
//	}
//
//	blockID = strings.Replace(blockID, "/", "slash", -1)
//
//	fp, err := os.OpenFile(filepath.Join(f.path, blockID), os.O_CREATE+os.O_WRONLY, 0600)
//	if err != nil {
//		return BlockBlobStageBlockResponse{}, fmt.Errorf("could not create a stage block file: %s", err)
//	}
//	defer func(fp *os.File) {
//		_ = fp.Close()
//	}(fp)
//
//	if _, err := io.Copy(fp, body); err != nil {
//		return BlockBlobStageBlockResponse{}, err
//	}
//
//	return BlockBlobStageBlockResponse{}, nil
//}
//
////nolint
//func (f *fakeBlockWriter) CommitBlockList(_ context.Context, base64BlockIDs []string, _ *BlockBlobCommitBlockListOptions) (BlockBlobCommitBlockListResponse, error) {
//	dst, err := os.OpenFile(filepath.Join(f.path, finalFileName), os.O_CREATE+os.O_WRONLY, 0600)
//	if err != nil {
//		return BlockBlobCommitBlockListResponse{}, err
//	}
//	defer func(dst *os.File) {
//		_ = dst.Close()
//	}(dst)
//
//	for _, id := range base64BlockIDs {
//		id = strings.Replace(id, "/", "slash", -1)
//		src, err := os.Open(filepath.Join(f.path, id))
//		if err != nil {
//			return BlockBlobCommitBlockListResponse{}, fmt.Errorf("could not combine chunk %s: %s", id, err)
//		}
//		_, err = io.Copy(dst, src)
//		if err != nil {
//			return BlockBlobCommitBlockListResponse{}, fmt.Errorf("problem writing final file from chunks: %s", err)
//		}
//		err = src.Close()
//		if err != nil {
//			return BlockBlobCommitBlockListResponse{}, fmt.Errorf("problem closing the source : %s", err)
//		}
//	}
//	return BlockBlobCommitBlockListResponse{}, nil
//}
//
////nolint
//func (f *fakeBlockWriter) cleanup() {
//	err := os.RemoveAll(f.path)
//	if err != nil {
//		return
//	}
//}
//
////nolint
//func (f *fakeBlockWriter) final() string {
//	return filepath.Join(f.path, finalFileName)
//}
//
////nolint
//func createSrcFile(size int) (string, error) {
//	generatedUuid, err := uuid.New()
//	if err != nil {
//		return "", err
//	}
//	p := filepath.Join(os.TempDir(), generatedUuid.String())
//	fp, err := os.OpenFile(p, os.O_CREATE+os.O_WRONLY, 0600)
//	if err != nil {
//		return "", fmt.Errorf("could not create source file: %s", err)
//	}
//	defer func(fp *os.File) {
//		_ = fp.Close()
//	}(fp)
//
//	lr := &io.LimitedReader{R: rand.New(rand.NewSource(time.Now().UnixNano())), N: int64(size)}
//	copied, err := io.Copy(fp, lr)
//	switch {
//	case err != nil && err != io.EOF:
//		return "", fmt.Errorf("copying %v: %s", size, err)
//	case copied != int64(size):
//		return "", fmt.Errorf("copying %v: copied %d bytes, expected %d", size, copied, size)
//	}
//	return p, nil
//}
//
////nolint
//func fileMD5(p string) string {
//	f, err := os.Open(p)
//	if err != nil {
//		panic(err)
//	}
//	defer func(f *os.File) {
//		_ = f.Close()
//
//	}(f)
//
//	h := md5.New()
//	if _, err := io.Copy(h, f); err != nil {
//		panic(err)
//	}
//
//	return fmt.Sprintf("%x", h.Sum(nil))
//}
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestGetErr() {
//	s.T().Parallel()
//
//	canceled, cancel := context.WithCancel(context.Background())
//	cancel()
//	err := errors.New("error")
//
//	tests := []struct {
//		desc string
//		ctx  context.Context
//		err  error
//		want error
//	}{
//		{"No errors", context.Background(), nil, nil},
//		{"Context was cancelled", canceled, nil, context.Canceled},
//		{"Context was cancelled but had error", canceled, err, err},
//		{"Err returned", context.Background(), err, err},
//	}
//
//	tm, err := internal.NewStaticBuffer(_1MiB, 1)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, test := range tests {
//		c := copier{
//			errCh: make(chan error, 1),
//			ctx:   test.ctx,
//			o:     UploadStreamOptions{TransferManager: tm},
//		}
//		if test.err != nil {
//			c.errCh <- test.err
//		}
//
//		got := c.getErr()
//		if test.want != got {
//			s.T().Errorf("TestGetErr(%s): got %v, want %v", test.desc, got, test.want)
//		}
//	}
//}
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestCopyFromReader() {
//	s.T().Parallel()
//
//	canceled, cancel := context.WithCancel(context.Background())
//	cancel()
//
//	spm, err := internal.NewSyncPool(_1MiB, 2)
//	if err != nil {
//		panic(err)
//	}
//	defer spm.Close()
//
//	tests := []struct {
//		desc      string
//		ctx       context.Context
//		o         UploadStreamOptions
//		fileSize  int
//		uploadErr bool
//		err       bool
//	}{
//		{
//			desc: "context was cancelled",
//			ctx:  canceled,
//			err:  true,
//		},
//		{
//			desc:     "Send file(0 KiB) with default UploadStreamOptions",
//			ctx:      context.Background(),
//			fileSize: 0,
//		},
//		{
//			desc:     "Send file(10 KiB) with default UploadStreamOptions",
//			ctx:      context.Background(),
//			fileSize: 10 * 1024,
//		},
//		{
//			desc:     "Send file(10 KiB) with default UploadStreamOptions set to azcopy settings",
//			ctx:      context.Background(),
//			fileSize: 10 * 1024,
//			o:        UploadStreamOptions{MaxBuffers: 5, BufferSize: 8 * 1024 * 1024},
//		},
//		{
//			desc:     "Send file(1 MiB) with default UploadStreamOptions",
//			ctx:      context.Background(),
//			fileSize: _1MiB,
//		},
//		{
//			desc:     "Send file(1 MiB) with default UploadStreamOptions set to azcopy settings",
//			ctx:      context.Background(),
//			fileSize: _1MiB,
//			o:        UploadStreamOptions{MaxBuffers: 5, BufferSize: 8 * 1024 * 1024},
//		},
//		{
//			desc:     "Send file(1.5 MiB) with default UploadStreamOptions",
//			ctx:      context.Background(),
//			fileSize: _1MiB + 500*1024,
//		},
//		{
//			desc:     "Send file(1.5 MiB) with 2 writers",
//			ctx:      context.Background(),
//			fileSize: _1MiB + 500*1024 + 1,
//			o:        UploadStreamOptions{MaxBuffers: 2},
//		},
//		{
//			desc:      "Send file(12 MiB) with 3 writers and 1 MiB buffer and a write error",
//			ctx:       context.Background(),
//			fileSize:  12 * _1MiB,
//			o:         UploadStreamOptions{MaxBuffers: 2, BufferSize: _1MiB},
//			uploadErr: true,
//			err:       true,
//		},
//		{
//			desc:     "Send file(12 MiB) with 3 writers and 1.5 MiB buffer",
//			ctx:      context.Background(),
//			fileSize: 12 * _1MiB,
//			o:        UploadStreamOptions{MaxBuffers: 2, BufferSize: _1MiB + .5*_1MiB},
//		},
//		{
//			desc:     "Send file(12 MiB) with default UploadStreamOptions set to azcopy settings",
//			ctx:      context.Background(),
//			fileSize: 12 * _1MiB,
//			o:        UploadStreamOptions{MaxBuffers: 5, BufferSize: 8 * 1024 * 1024},
//		},
//		{
//			desc:     "Send file(12 MiB) with default UploadStreamOptions using SyncPool manager",
//			ctx:      context.Background(),
//			fileSize: 12 * _1MiB,
//			o: UploadStreamOptions{
//				TransferManager: spm,
//			},
//		},
//	}
//
//	for _, test := range tests {
//		p, err := createSrcFile(test.fileSize)
//		if err != nil {
//			panic(err)
//		}
//		defer func(name string) {
//			_ = os.Remove(name)
//		}(p)
//
//		from, err := os.Open(p)
//		if err != nil {
//			panic(err)
//		}
//
//		br := newFakeBlockWriter()
//		defer br.cleanup()
//		if test.uploadErr {
//			br.errOnBlock = 1
//		}
//
//		_, err = copyFromReader(test.ctx, from, br, test.o)
//		switch {
//		case err == nil && test.err:
//			s.T().Errorf("TestCopyFromReader(%s): got err == nil, want err != nil", test.desc)
//			continue
//		case err != nil && !test.err:
//			s.T().Errorf("TestCopyFromReader(%s): got err == %s, want err == nil", test.desc, err)
//			continue
//		case err != nil:
//			continue
//		}
//
//		want := fileMD5(p)
//		got := fileMD5(br.final())
//
//		if got != want {
//			s.T().Errorf("TestCopyFromReader(%s): MD5 not the same: got %s, want %s", test.desc, got, want)
//		}
//	}
//}

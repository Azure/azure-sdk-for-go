// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"crypto/rand"
	"errors"
	"io"
)

const (
	defaultLength = 1024 * 1024
)

type randomStream struct {
	offset   int
	baseData []byte
}

func (r *randomStream) Read(p []byte) (n int, err error) {
	if len(p) < len(r.baseData)-r.offset {
		n = copy(p, r.baseData[r.offset:r.offset+len(p)])
		r.offset += n
		return n, nil
	}

	n = copy(p, r.baseData[r.offset:])
	r.offset += n
	return n, io.EOF
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	var i int64
	switch whence {
	case io.SeekStart:
		i = offset
	case io.SeekCurrent:
		i = int64(r.offset) + offset
	case io.SeekEnd:
		i = int64(len(r.baseData)) + offset
	default:
		return 0, errors.New("randomStreamSeek: invalid whence")
	}
	if i < 0 {
		return 0, errors.New("randomStreamSeek: negative position")
	}
	r.offset = int(i)
	return i, nil
}

func (r randomStream) Close() error {
	return nil
}

func getRandomBytes(i int) ([]byte, error) {
	ret := make([]byte, i)
	_, err := rand.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func NewRandomStream(size int) (io.ReadSeekCloser, error) {
	baseData, err := getRandomBytes(defaultLength)
	if err != nil {
		return nil, err
	}
	return &randomStream{
		offset:   0,
		baseData: baseData,
	}, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

const (
	defaultLength = 1024 * 1024
)

type randomStream struct {
	baseData         []byte
	dataLength       int64
	baseBufferLength int
	position         int64
	remaining        int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *randomStream) Read(p []byte) (int, error) {
	if r.remaining == 0 {
		return 0, io.EOF
	}

	size := len(p)
	e := min(size, r.remaining)
	if e > r.baseBufferLength {
		// Need to create a larger buffer
		b, err := getRandomBytes(e)
		if err != nil {
			return 0, err
		}
		r.baseData = b
		r.baseBufferLength = e
	}

	r.remaining -= e
	r.position += int64(e)

	n := copy(p, r.baseData[:e])
	return n, nil
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.position = offset
	case io.SeekCurrent:
		r.position += offset
	case io.SeekEnd:
		r.position = r.dataLength + offset
	default:
		return 0, fmt.Errorf("randomStream: did not understand whence: %d", whence)
	}
	r.remaining = int(r.dataLength - r.position)
	if r.position < 0 {
		return 0, errors.New("randomStream: negative position")
	}
	return r.position, nil
}

func (r randomStream) Close() error {
	return nil
}

func getRandomBytes(i int) ([]byte, error) {
	ret := make([]byte, i)
	n, err := rand.Read(ret)
	if err != nil {
		return nil, err
	}
	if n != i {
		return nil, fmt.Errorf("did not create a byte slice of size %d, got %d", i, n)
	}
	return ret, nil
}

func NewRandomStream(length int) (io.ReadSeekCloser, error) {
	base, err := getRandomBytes(min(length, defaultLength))
	if err != nil {
		return nil, err
	}
	return &randomStream{
		baseData:         base,
		dataLength:       int64(length),
		baseBufferLength: defaultLength,
		position:         0,
		remaining:        length,
	}, nil
}

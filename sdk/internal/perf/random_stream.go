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
	dataLength       int
	baseBufferLength int
	position         int
	remaining        int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *randomStream) Read(p []byte) (int, error) {
	fmt.Println("read")

	// At the end of the buffer
	if r.remaining == 0 {
		return 0, io.EOF
	}

	var e int
	if len(p) == 0 {
		// p has no more area to fill
		return 0, io.EOF
	} else {
		// bytes to copy into p
		e = len(p)
	}

	e = min(e, r.remaining)
	if e > r.baseBufferLength {
		newBase, err := getRandomBytes(e)
		if err != nil {
			return 0, err
		}
		r.baseData = newBase
		r.baseBufferLength = 0
	}
	n := copy(p, r.baseData[r.position:r.position + e])
	r.remaining -= n
	r.position += n
	fmt.Println("e: ", e)
	fmt.Println("position: ", r.position)
	fmt.Println("len(baseData): ", len(r.baseData))

	return n, nil
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	fmt.Println("seek")
	fmt.Println(offset, whence)
	switch whence {
	case io.SeekStart:
		r.position = int(offset)
	case io.SeekCurrent:
		r.position += int(offset)
	case io.SeekEnd:
		r.position = r.dataLength - 1 + int(offset)
	default:
		return 0, errors.New("randomStream: invalid whence")
	}
	if r.position < 0 {
		return 0, errors.New("randomStream: negative position")
	}
	r.remaining = r.dataLength - r.position
	return int64(r.position), nil
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

func NewRandomStream(size int) (io.ReadSeekCloser, error) {
	baseData, err := getRandomBytes(defaultLength)
	if err != nil {
		return nil, err
	}
	return &randomStream{
		baseData:         baseData,
		dataLength:       size,
		baseBufferLength: len(baseData),
		position:         0,
		remaining:        size,
	}, nil
}

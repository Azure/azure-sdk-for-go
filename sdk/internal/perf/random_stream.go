// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	defaultLength = 1024 * 1024
)

type randomStream struct {
	readSeeker io.ReadSeeker
}

func (r *randomStream) Read(p []byte) (int, error) {
	return r.readSeeker.Read(p)
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	return r.readSeeker.Seek(offset, whence)
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
	baseData, err := getRandomBytes(size)
	if err != nil {
		return nil, err
	}
	return &randomStream{
		readSeeker: bytes.NewReader(baseData),
	}, nil
}

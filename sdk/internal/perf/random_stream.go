// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package perf

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

const (
	defaultLength = 1024 * 1024
)

type randomStream struct {
	reader    *bytes.Reader
	remaining int
}

func (r *randomStream) Read(p []byte) (n int, err error) {
	fmt.Println("read")
	return r.reader.Read(p)
}

func (r *randomStream) Seek(offset int64, whence int) (int64, error) {
	fmt.Println("seek")
	return r.reader.Seek(offset, whence)
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
	reader := bytes.NewReader(baseData)
	return streaming.NopCloser(reader), nil
}

//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package hashing

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"hash/crc64"
	"io"
)

// ReadWrapper is an io.ReadSeekCloser wrapper to prevent reading twice for SDK-generated hashes.
type ReadWrapper struct {
	internal *bytes.Reader
	crc64    uint64
}

func (h *ReadWrapper) CRC64Hash() uint64 {
	return h.crc64
}

func (h *ReadWrapper) Read(p []byte) (n int, err error) {
	return h.internal.Read(p)
}

func (h *ReadWrapper) Seek(offset int64, whence int) (int64, error) {
	return h.internal.Seek(offset, whence)
}

func (h *ReadWrapper) Close() error {
	return nil
}

func NewReadWrapper(r io.ReadSeekCloser, validationOption blob.TransferValidationType) (*ReadWrapper, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var crc uint64
	if validationOption&blob.TransferValidationTypeCRC64 == blob.TransferValidationTypeCRC64 {
		crc = crc64.Checksum(buf, CRC64Table)
	}

	return &ReadWrapper{
		internal: bytes.NewReader(buf),
		crc64:    crc,
	}, nil
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"bytes"
	"encoding/binary"
	"hash/crc64"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/structuredmessage"
)

// TransferValidationType abstracts the various mechanisms used to verify a transfer.
type TransferValidationType interface {
	Apply(io.ReadSeekCloser, generated.TransactionalContentSetter) (io.ReadSeekCloser, error)
	notPubliclyImplementable()
}

// StructuredMessageSetter defines the interface for setting structured message headers
type StructuredMessageSetter interface {
	SetStructuredBodyType(bodyType string)
	SetStructuredContentLength(length int64)
}

// StructuredTransferValidationType extends TransferValidationType for structured message support
type StructuredTransferValidationType interface {
	ApplyStructured(io.ReadSeekCloser, StructuredMessageSetter) (io.ReadSeekCloser, error)
	notPubliclyImplementable()
}

// TransferValidationTypeCRC64 is a TransferValidationType used to provide a precomputed CRC64.
type TransferValidationTypeCRC64 uint64

func (c TransferValidationTypeCRC64) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(c))
	cfg.SetCRC64(buf)
	return rsc, nil
}

func (TransferValidationTypeCRC64) notPubliclyImplementable() {}

// TransferValidationTypeComputeCRC64 is a TransferValidationType that indicates a CRC64 should be computed during transfer.
func TransferValidationTypeComputeCRC64() TransferValidationType {
	return transferValidationTypeFn(func(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
		buf, err := io.ReadAll(rsc)
		if err != nil {
			return nil, err
		}

		crc := crc64.Checksum(buf, shared.CRC64Table)
		return TransferValidationTypeCRC64(crc).Apply(streaming.NopCloser(bytes.NewReader(buf)), cfg)
	})
}

type transferValidationTypeFn func(io.ReadSeekCloser, generated.TransactionalContentSetter) (io.ReadSeekCloser, error)

func (t transferValidationTypeFn) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	return t(rsc, cfg)
}

func (transferValidationTypeFn) notPubliclyImplementable() {}

// TransferValidationTypeStructuredMessage is a transfer validation type that uses structured message format with CRC64
type TransferValidationTypeStructuredMessage struct{}

func (t TransferValidationTypeStructuredMessage) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	// Structured messages don't use transactional content headers, so this should not be called directly
	// Use ApplyStructured instead
	return rsc, nil
}

func (t TransferValidationTypeStructuredMessage) ApplyStructured(rsc io.ReadSeekCloser, cfg StructuredMessageSetter) (io.ReadSeekCloser, error) {
	// Create encoder reader to wrap the data in structured message format
	encoder, err := structuredmessage.NewEncoderReader(rsc)
	if err != nil {
		return nil, err
	}

	// Set structured message headers
	bodyType, contentLength := encoder.GetStructuredHeaders()
	cfg.SetStructuredBodyType(bodyType)
	cfg.SetStructuredContentLength(contentLength)

	// Return a ReadSeekCloser that provides the encoded data
	return &structuredMessageReadSeekCloser{
		encoder: encoder,
	}, nil
}

func (TransferValidationTypeStructuredMessage) notPubliclyImplementable() {}

// structuredMessageReadSeekCloser wraps EncoderReader to provide ReadSeekCloser interface
type structuredMessageReadSeekCloser struct {
	encoder *structuredmessage.EncoderReader
	reader  io.ReadSeeker
	closed  bool
}

func (sm *structuredMessageReadSeekCloser) Read(p []byte) (n int, err error) {
	if sm.closed {
		return 0, io.EOF
	}
	if sm.reader == nil {
		// Create reader from encoded data on first read
		encodedData, err := io.ReadAll(sm.encoder)
		if err != nil {
			return 0, err
		}
		sm.reader = bytes.NewReader(encodedData)
	}
	return sm.reader.Read(p)
}

func (sm *structuredMessageReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	if sm.closed {
		return 0, io.EOF
	}
	if sm.reader == nil {
		// Create reader from encoded data on first seek
		encodedData, err := io.ReadAll(sm.encoder)
		if err != nil {
			return 0, err
		}
		sm.reader = bytes.NewReader(encodedData)
	}
	return sm.reader.Seek(offset, whence)
}

func (sm *structuredMessageReadSeekCloser) Close() error {
	sm.closed = true
	return nil
}

// TransferValidationTypeStructuredMessageDownload handles structured message decoding for downloads
type TransferValidationTypeStructuredMessageDownload struct {
	StructuredContentLength int64
}

func (t TransferValidationTypeStructuredMessageDownload) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	// For downloads, we need to decode the structured message and validate CRC64
	decoder, err := structuredmessage.NewDecoderReader(rsc, t.StructuredContentLength)
	if err != nil {
		return nil, err
	}

	// Return a ReadSeekCloser that provides the decoded data
	return &structuredMessageDecodeReadSeekCloser{
		decoder: decoder,
	}, nil
}

func (TransferValidationTypeStructuredMessageDownload) notPubliclyImplementable() {}

// structuredMessageDecodeReadSeekCloser wraps DecoderReader to provide ReadSeekCloser interface
type structuredMessageDecodeReadSeekCloser struct {
	decoder *structuredmessage.DecoderReader
	reader  io.ReadSeeker
	closed  bool
}

func (sm *structuredMessageDecodeReadSeekCloser) Read(p []byte) (n int, err error) {
	if sm.closed {
		return 0, io.EOF
	}
	if sm.reader == nil {
		// Create reader from decoded data on first read
		decodedData, err := io.ReadAll(sm.decoder)
		if err != nil {
			return 0, err
		}
		sm.reader = bytes.NewReader(decodedData)
	}
	return sm.reader.Read(p)
}

func (sm *structuredMessageDecodeReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	if sm.closed {
		return 0, io.EOF
	}
	if sm.reader == nil {
		// Create reader from decoded data on first seek
		decodedData, err := io.ReadAll(sm.decoder)
		if err != nil {
			return 0, err
		}
		sm.reader = bytes.NewReader(decodedData)
	}
	return sm.reader.Seek(offset, whence)
}

func (sm *structuredMessageDecodeReadSeekCloser) Close() error {
	sm.closed = true
	return nil
}

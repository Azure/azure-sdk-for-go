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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal/structuredmsg"
)

// TransferValidationType abstracts the various mechanisms used to verify a transfer.
type TransferValidationType interface {
	Apply(io.ReadSeekCloser, generated.TransactionalContentSetter) (io.ReadSeekCloser, error)
	notPubliclyImplementable()
	supportsMultiBlock() bool
}

// SupportsMultiBlock reports whether the validation type can be used with multi-block uploads.
func SupportsMultiBlock(tv TransferValidationType) bool {
	return tv.supportsMultiBlock()
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
func (TransferValidationTypeCRC64) supportsMultiBlock() bool  { return false }

// TransferValidationTypeComputeCRC64 is a TransferValidationType that indicates a CRC64 should be computed during transfer.
func TransferValidationTypeComputeCRC64() TransferValidationType {
	return transferValidationTypeFn(func(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
		buf, err := io.ReadAll(rsc)
		if err != nil {
			return nil, err
		}

		crc := crc64.Checksum(buf, structuredmsg.CRC64Table)
		return TransferValidationTypeCRC64(crc).Apply(streaming.NopCloser(bytes.NewReader(buf)), cfg)
	})
}

// TransferValidationTypeComputeStructuredMessageCRC64 is a TransferValidationType that computes
// per-segment CRC64 checksums using the structured message binary format.
// The body is wrapped in a streaming SMEncoder that produces SM-encoded output on Read().
// segmentSize specifies the maximum segment size in bytes. Values <= 0 use the default (4 MB).
func TransferValidationTypeComputeStructuredMessageCRC64(segmentSize int) TransferValidationType {
	return &transferValidationTypeSMCRC64{segmentSize: segmentSize}
}

type transferValidationTypeSMCRC64 struct {
	segmentSize int
}

func (t *transferValidationTypeSMCRC64) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	contentLen, err := internal.ValidateSeekableStreamAt0AndGetCount(rsc)
	if err != nil {
		return nil, err
	}

	encoder := structuredmsg.NewSMEncoder(rsc, contentLen, t.segmentSize)
	cfg.SetStructuredBody(structuredmsg.SMHeaderValue, encoder.OriginalContentLength())
	return encoder, nil
}

func (*transferValidationTypeSMCRC64) notPubliclyImplementable() {}
func (*transferValidationTypeSMCRC64) supportsMultiBlock() bool  { return true }

func (t *transferValidationTypeSMCRC64) StructuredBodyHeaderValue() string {
	return structuredmsg.SMHeaderValue
}

// GetStructuredBodyType returns the structured body header value if the given TransferValidationType
// is a structured message type, or empty string otherwise.
func GetStructuredBodyType(tv TransferValidationType) string {
	if sm, ok := tv.(*transferValidationTypeSMCRC64); ok {
		return sm.StructuredBodyHeaderValue()
	}
	return ""
}

type transferValidationTypeFn func(io.ReadSeekCloser, generated.TransactionalContentSetter) (io.ReadSeekCloser, error)

func (t transferValidationTypeFn) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	return t(rsc, cfg)
}

func (transferValidationTypeFn) notPubliclyImplementable() {}
func (transferValidationTypeFn) supportsMultiBlock() bool  { return true }

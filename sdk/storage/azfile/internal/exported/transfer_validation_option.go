// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
)

// TransferValidationType abstracts the various mechanisms used to verify a transfer.
type TransferValidationType interface {
	Apply(io.ReadSeekCloser, generated.TransactionalContentSetter) (io.ReadSeekCloser, error)
	notPubliclyImplementable()
}

// TransferValidationTypeMD5 is a TransferValidationType used to provide a precomputed MD5.
type TransferValidationTypeMD5 []byte

func (c TransferValidationTypeMD5) Apply(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	cfg.SetMD5(c)
	return rsc, nil
}

func (TransferValidationTypeMD5) notPubliclyImplementable() {}

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
	contentLen, err := shared.ValidateSeekableStreamAt0AndGetCount(rsc)
	if err != nil {
		return nil, err
	}

	encoder := shared.NewSMEncoder(rsc, contentLen, t.segmentSize)
	cfg.SetStructuredBody(shared.SMHeaderValue, encoder.OriginalContentLength())
	return encoder, nil
}

func (*transferValidationTypeSMCRC64) notPubliclyImplementable() {}

func (t *transferValidationTypeSMCRC64) StructuredBodyHeaderValue() string {
	return shared.SMHeaderValue
}

// GetStructuredBodyType returns the structured body header value if the given TransferValidationType
// is a structured message type, or empty string otherwise.
func GetStructuredBodyType(tv TransferValidationType) string {
	if sm, ok := tv.(*transferValidationTypeSMCRC64); ok {
		return sm.StructuredBodyHeaderValue()
	}
	return ""
}

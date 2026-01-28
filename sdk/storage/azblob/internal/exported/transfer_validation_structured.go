//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"bytes"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmsg"
)

const (
	// StructuredMessageThreshold is the size threshold above which structured message format is used
	StructuredMessageThreshold = 4 * 1024 * 1024
)

// TransferValidationTypeStructuredMessage is a TransferValidationType that wraps data in structured message format
func TransferValidationTypeStructuredMessage() TransferValidationType {
	return transferValidationTypeFn(func(rsc io.ReadSeekCloser, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
		// Read all data to determine size
		buf, err := io.ReadAll(rsc)
		if err != nil {
			return nil, err
		}

		size := int64(len(buf))

		// For small uploads (<4MB), use simple CRC64 header approach
		if size < StructuredMessageThreshold {
			return TransferValidationTypeComputeCRC64().Apply(streaming.NopCloser(bytes.NewReader(buf)), cfg)
		}

		// For large uploads, use structured message format
		return wrapInStructuredMessage(buf, cfg)
	})
}

// wrapInStructuredMessage wraps the data in structured message format and sets appropriate headers
func wrapInStructuredMessage(data []byte, cfg generated.TransactionalContentSetter) (io.ReadSeekCloser, error) {
	dataSize := int64(len(data))

	// Calculate segment count (max 4MB per segment)
	segmentSize := int64(structuredmsg.MaxSegmentSize)
	numSegments := uint16((dataSize + segmentSize - 1) / segmentSize)

	// Calculate total message length with CRC64 (required for transfer validation):
	// HeaderSize + sum of (SegmentHeaderSize + dataLength + SegmentCRC64Size for each segment) + TrailerSize
	messageLength := uint64(structuredmsg.HeaderSize) // Header
	offset := int64(0)
	for i := uint16(0); i < numSegments; i++ {
		remaining := dataSize - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}
		messageLength += uint64(structuredmsg.SegmentHeaderSize) + uint64(chunkSize) + uint64(structuredmsg.SegmentCRC64Size)
		offset += chunkSize
	}
	messageLength += uint64(structuredmsg.TrailerSize) // Trailer

	// Buffer to hold the structured message
	var structuredBuffer bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&structuredBuffer)

	// Write header with CRC64 enabled (required for transfer validation)
	if err := writer.WriteHeader(numSegments, messageLength, true); err != nil {
		return nil, err
	}

	// Write segments
	offset = 0
	for i := uint16(0); i < numSegments; i++ {
		remaining := dataSize - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}

		segmentData := data[offset : offset+chunkSize]
		if err := writer.WriteSegment(segmentData); err != nil {
			return nil, err
		}

		offset += chunkSize
	}

	// Write trailer
	if err := writer.WriteTrailer(); err != nil {
		return nil, err
	}

	// Set headers
	cfg.SetStructuredBody("XSM/1.0;CRC64", dataSize)

	return streaming.NopCloser(bytes.NewReader(structuredBuffer.Bytes())), nil
}

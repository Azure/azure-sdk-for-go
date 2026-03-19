//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"io"
)

// Structured Message (XSM/1.0) constants
const (
	SMVersion   uint8  = 1
	SMFlagCRC64 uint16 = 0x0001

	SMDefaultSegmentSize = 4 * 1024 * 1024 // 4MB

	SMHeaderSize         = 13 // version(1) + msgLen(8) + flags(2) + numSegments(2)
	SMSegmentHeaderSize  = 10 // segNum(2) + dataLen(8)
	SMSegmentFooterSize  = 8  // CRC64
	SMMessageTrailerSize = 8  // CRC64

	SMHeaderValue = "XSM/1.0; properties=crc64"
)

// SMEncodeResult holds the result of encoding data into structured message format.
type SMEncodeResult struct {
	// EncodedData is the complete structured message binary payload.
	EncodedData []byte
	// OriginalContentLength is the length of the original unframed data.
	OriginalContentLength int64
}

// SMEncode encodes raw data into structured message format.
// segmentSize specifies the maximum segment size (use 0 for default 4MB).
// Returns the full SM binary payload and the original content length.
func SMEncode(data []byte, segmentSize int) SMEncodeResult {
	if segmentSize <= 0 {
		segmentSize = SMDefaultSegmentSize
	}

	totalDataLen := len(data)
	numSegments := totalDataLen / segmentSize
	if totalDataLen%segmentSize != 0 {
		numSegments++
	}
	if numSegments == 0 {
		numSegments = 1
	}

	// Calculate total message length
	msgLen := int64(SMHeaderSize)
	for i := 0; i < numSegments; i++ {
		segStart := i * segmentSize
		segEnd := segStart + segmentSize
		if segEnd > totalDataLen {
			segEnd = totalDataLen
		}
		segDataLen := segEnd - segStart
		msgLen += int64(SMSegmentHeaderSize) + int64(segDataLen) + int64(SMSegmentFooterSize)
	}
	msgLen += int64(SMMessageTrailerSize)

	var buf bytes.Buffer
	buf.Grow(int(msgLen))

	// Message Header (13 bytes)
	buf.WriteByte(SMVersion)
	_ = binary.Write(&buf, binary.LittleEndian, msgLen)
	_ = binary.Write(&buf, binary.LittleEndian, SMFlagCRC64)
	_ = binary.Write(&buf, binary.LittleEndian, uint16(numSegments))

	// Compute message-level CRC64 over all raw data
	messageCRC := crc64.Checksum(data, CRC64Table)

	for i := 0; i < numSegments; i++ {
		segStart := i * segmentSize
		segEnd := segStart + segmentSize
		if segEnd > totalDataLen {
			segEnd = totalDataLen
		}
		segData := data[segStart:segEnd]

		// Segment Header (10 bytes)
		_ = binary.Write(&buf, binary.LittleEndian, uint16(i+1))
		_ = binary.Write(&buf, binary.LittleEndian, int64(len(segData)))

		// Segment Data
		buf.Write(segData)

		// Segment Footer - CRC64 of segment data (8 bytes)
		segCRC := crc64.Checksum(segData, CRC64Table)
		_ = binary.Write(&buf, binary.LittleEndian, segCRC)
	}

	// Message Trailer - CRC64 of all raw content (8 bytes)
	_ = binary.Write(&buf, binary.LittleEndian, messageCRC)

	return SMEncodeResult{
		EncodedData:           buf.Bytes(),
		OriginalContentLength: int64(totalDataLen),
	}
}

// SMDecodeResult holds the result of decoding a structured message.
type SMDecodeResult struct {
	// Data is the extracted raw content.
	Data []byte
	// Version is the SM protocol version.
	Version uint8
	// Flags are the SM flags (e.g., SMFlagCRC64).
	Flags uint16
	// NumSegments is the number of segments in the message.
	NumSegments uint16
}

// SMDecode decodes a structured message binary payload and validates CRC64 checksums.
// Returns the extracted raw data or an error if the SM is malformed or CRC validation fails.
func SMDecode(smData []byte) (SMDecodeResult, error) {
	if len(smData) < SMHeaderSize {
		return SMDecodeResult{}, fmt.Errorf("structured message too short for header: %d bytes (minimum %d)", len(smData), SMHeaderSize)
	}

	// Parse Message Header
	version := smData[0]
	if version != SMVersion {
		return SMDecodeResult{}, fmt.Errorf("unsupported structured message version: %d (expected %d)", version, SMVersion)
	}

	msgLen := binary.LittleEndian.Uint64(smData[1:9])
	if uint64(len(smData)) != msgLen {
		return SMDecodeResult{}, fmt.Errorf("structured message length mismatch: header says %d, got %d bytes", msgLen, len(smData))
	}

	flags := binary.LittleEndian.Uint16(smData[9:11])
	numSegments := binary.LittleEndian.Uint16(smData[11:13])

	hasCRC := flags&SMFlagCRC64 != 0

	offset := SMHeaderSize
	var result []byte

	if hasCRC {
		msgHasher := crc64.New(CRC64Table)

		for i := 0; i < int(numSegments); i++ {
			if offset+SMSegmentHeaderSize > len(smData) {
				return SMDecodeResult{}, fmt.Errorf("segment %d: insufficient data for segment header at offset %d", i+1, offset)
			}

			segNum := binary.LittleEndian.Uint16(smData[offset : offset+2])
			segDataLen := int64(binary.LittleEndian.Uint64(smData[offset+2 : offset+10]))
			offset += SMSegmentHeaderSize

			if segNum != uint16(i+1) {
				return SMDecodeResult{}, fmt.Errorf("segment number mismatch: expected %d, got %d", i+1, segNum)
			}

			if segDataLen < 0 || offset+int(segDataLen) > len(smData) {
				return SMDecodeResult{}, fmt.Errorf("segment %d: insufficient data for segment content (need %d bytes at offset %d, have %d)", segNum, segDataLen, offset, len(smData))
			}

			segData := smData[offset : offset+int(segDataLen)]
			offset += int(segDataLen)
			result = append(result, segData...)

			if offset+SMSegmentFooterSize > len(smData) {
				return SMDecodeResult{}, fmt.Errorf("segment %d: insufficient data for CRC64 footer at offset %d", segNum, offset)
			}

			expectedCRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
			actualCRC := crc64.Checksum(segData, CRC64Table)
			offset += SMSegmentFooterSize

			if expectedCRC != actualCRC {
				return SMDecodeResult{}, fmt.Errorf("segment %d: CRC64 mismatch (expected 0x%016x, got 0x%016x)", segNum, expectedCRC, actualCRC)
			}

			msgHasher.Write(segData)
		}

		// Validate message trailer CRC64
		if offset+SMMessageTrailerSize > len(smData) {
			return SMDecodeResult{}, fmt.Errorf("insufficient data for message trailer CRC64 at offset %d", offset)
		}

		expectedMsgCRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
		actualMsgCRC := msgHasher.Sum64()

		if expectedMsgCRC != actualMsgCRC {
			return SMDecodeResult{}, fmt.Errorf("message trailer CRC64 mismatch (expected 0x%016x, got 0x%016x)", expectedMsgCRC, actualMsgCRC)
		}
	} else {
		// No CRC — just extract segment data
		for i := 0; i < int(numSegments); i++ {
			if offset+SMSegmentHeaderSize > len(smData) {
				return SMDecodeResult{}, fmt.Errorf("segment %d: insufficient data for segment header at offset %d", i+1, offset)
			}

			segNum := binary.LittleEndian.Uint16(smData[offset : offset+2])
			segDataLen := int64(binary.LittleEndian.Uint64(smData[offset+2 : offset+10]))
			offset += SMSegmentHeaderSize

			if segNum != uint16(i+1) {
				return SMDecodeResult{}, fmt.Errorf("segment number mismatch: expected %d, got %d", i+1, segNum)
			}

			if segDataLen < 0 || offset+int(segDataLen) > len(smData) {
				return SMDecodeResult{}, fmt.Errorf("segment %d: insufficient data for segment content", segNum)
			}

			segData := smData[offset : offset+int(segDataLen)]
			offset += int(segDataLen)
			result = append(result, segData...)
		}
	}

	return SMDecodeResult{
		Data:        result,
		Version:     version,
		Flags:       flags,
		NumSegments: numSegments,
	}, nil
}

// SMEncoder wraps raw data as an io.ReadSeekCloser that produces SM-encoded output.
// This is used for uploads — the body content is replaced with the SM-encoded form.
type SMEncoder struct {
	reader *bytes.Reader
	result SMEncodeResult
}

// NewSMEncoder creates a new encoder that wraps the given raw data.
// segmentSize specifies the max segment size; use 0 for the default (4MB).
func NewSMEncoder(data []byte, segmentSize int) *SMEncoder {
	result := SMEncode(data, segmentSize)
	return &SMEncoder{
		reader: bytes.NewReader(result.EncodedData),
		result: result,
	}
}

func (e *SMEncoder) Read(p []byte) (int, error) {
	return e.reader.Read(p)
}

func (e *SMEncoder) Seek(offset int64, whence int) (int64, error) {
	return e.reader.Seek(offset, whence)
}

func (e *SMEncoder) Close() error {
	return nil
}

// OriginalContentLength returns the length of the original unframed data.
func (e *SMEncoder) OriginalContentLength() int64 {
	return e.result.OriginalContentLength
}

// EncodedLength returns the total length of the SM-encoded data.
func (e *SMEncoder) EncodedLength() int64 {
	return int64(len(e.result.EncodedData))
}

// SMDecoder wraps an SM-encoded io.ReadCloser and produces raw data on Read().
// It reads the full SM payload, validates CRC64 checksums, and serves the extracted raw data.
type SMDecoder struct {
	source    io.ReadCloser
	reader    *bytes.Reader
	decoded   bool
	rawData   []byte
	decResult SMDecodeResult
	err       error
}

// NewSMDecoder creates a decoder that wraps an SM-encoded body.
// The full body is read and validated on the first call to Read().
func NewSMDecoder(body io.ReadCloser) *SMDecoder {
	return &SMDecoder{
		source: body,
	}
}

func (d *SMDecoder) decode() error {
	if d.decoded {
		return d.err
	}
	d.decoded = true

	smData, err := io.ReadAll(d.source)
	if err != nil {
		d.err = fmt.Errorf("failed to read structured message body: %w", err)
		return d.err
	}

	result, err := SMDecode(smData)
	if err != nil {
		d.err = err
		return d.err
	}

	d.decResult = result
	d.rawData = result.Data
	d.reader = bytes.NewReader(d.rawData)
	return nil
}

func (d *SMDecoder) Read(p []byte) (int, error) {
	if err := d.decode(); err != nil {
		return 0, err
	}
	return d.reader.Read(p)
}

func (d *SMDecoder) Close() error {
	if d.source != nil {
		return d.source.Close()
	}
	return nil
}

// DecodeResult returns the decode result after the first Read().
// Returns nil if not yet decoded.
func (d *SMDecoder) DecodeResult() *SMDecodeResult {
	if !d.decoded {
		return nil
	}
	return &d.decResult
}

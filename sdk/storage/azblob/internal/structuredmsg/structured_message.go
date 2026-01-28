//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Package structuredmsg implements the XSM/1.0 structured message format for Azure Storage.
// This format provides CRC64 checksums for data integrity validation during uploads and downloads.
//
// All numeric fields are encoded in LittleEndian byte order, consistent with Azure Storage conventions.
//
// The structured message format consists of:
// - Header: version, flags, properties, and segment count
// - Segments: each with segment number, data length, data, and CRC64
// - Trailer: version, flags, properties, and cumulative CRC64
package structuredmsg

import (
	"encoding/binary"
	"errors"
	"hash/crc64"
	"io"
	"math"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// Constants for the structured message format
const (
	// MessageVersion is the Message format version
	MessageVersion = 1

	// HeaderSize is the Header size in bytes: version(1) + message-length(8) + message-flags(2) + num-segments(2) = 13 bytes
	HeaderSize = 13

	// TrailerSize is the Total trailer size in bytes: message-crc64(8) = 8 bytes
	TrailerSize = 8

	// SegmentHeaderSize is segment-num(2) + data-length(8) = 10 bytes
	// Note: segment-crc64(8) comes AFTER the data, not in the header
	SegmentHeaderSize = 10

	// SegmentCRC64Size is the size of segment CRC64 in bytes
	SegmentCRC64Size = 8

	// MaxSegments is Maximum number of segments allowed (uint16 max)
	MaxSegments = math.MaxUint16

	// MaxSegmentSize Maximum segment size (4MB)
	MaxSegmentSize = 4 * 1024 * 1024
)

const (
	FlagNone = 0x00
)

// Properties bitmask
const (
	PropCRC64 = 0x0001 // CRC64 checksum is present
)

// Errors
var (
	ErrInvalidVersion       = errors.New("invalid message version")
	ErrInvalidSegmentCount  = errors.New("invalid segment count")
	ErrSegmentTooLarge      = errors.New("segment too large")
	ErrCRC64Mismatch        = errors.New("CRC64 checksum mismatch")
	ErrUnexpectedEOF        = errors.New("unexpected end of file")
	ErrInvalidSegmentNumber = errors.New("invalid segment number")
	ErrInvalidDataLength    = errors.New("invalid data length")
)

// StructuredMessageWriter encodes data into the binary message format
type StructuredMessageWriter struct {
	writer      io.Writer
	crc64Table  *crc64.Table
	segmentNum  uint16
	totalCRC64  uint64
	enableCRC64 bool // Whether CRC64 is enabled (flags = 1) or disabled (flags = 0)
}

// NewStructuredMessageWriter creates a new StructuredMessageWriter
func NewStructuredMessageWriter(writer io.Writer) *StructuredMessageWriter {
	return &StructuredMessageWriter{
		writer:     writer,
		crc64Table: shared.CRC64Table,
		segmentNum: 1, // Segments are 1-indexed
		totalCRC64: 0,
	}
}

// WriteHeader writes the message header
// messageLength is the total size of the structured message in bytes:
//   - With CRC64: HeaderSize + sum of (SegmentHeaderSize + dataLength + SegmentCRC64Size for each segment) + TrailerSize
//   - Without CRC64: HeaderSize + sum of (SegmentHeaderSize + dataLength for each segment)
//
// enableCRC64: if true, flags = 1 (PropCRC64), segments and trailer include CRC64; if false, flags = 0, no CRC64
func (w *StructuredMessageWriter) WriteHeader(numSegments uint16, messageLength uint64, enableCRC64 bool) error {
	if numSegments == 0 || numSegments > MaxSegments {
		return ErrInvalidSegmentCount
	}

	w.enableCRC64 = enableCRC64

	header := make([]byte, HeaderSize)
	header[0] = MessageVersion                                // byte 0: version (1 byte)
	binary.LittleEndian.PutUint64(header[1:9], messageLength) // bytes 1-8: message-length (8 bytes)

	var flags uint16
	if enableCRC64 {
		flags = PropCRC64 // flags = 1 if CRC64 is included
	} else {
		flags = FlagNone // flags = 0 if CRC64 is excluded
	}
	binary.LittleEndian.PutUint16(header[9:11], flags)        // bytes 9-10: message-flags (2 bytes)
	binary.LittleEndian.PutUint16(header[11:13], numSegments) // bytes 11-12: num-segments (2 bytes)

	_, err := w.writer.Write(header)
	return err
}

// WriteSegment writes a segment with its data and optionally CRC64
func (w *StructuredMessageWriter) WriteSegment(data []byte) error {
	if len(data) > MaxSegmentSize {
		return ErrSegmentTooLarge
	}

	// Calculate segment CRC64 if enabled
	var segmentCRC64 uint64
	if w.enableCRC64 {
		segmentCRC64 = crc64.Checksum(data, w.crc64Table)
	}

	// Write segment header (segment-num + data-length)
	header := make([]byte, SegmentHeaderSize)
	binary.LittleEndian.PutUint16(header[0:2], w.segmentNum)       // bytes 0-1: segment-num (2 bytes)
	binary.LittleEndian.PutUint64(header[2:10], uint64(len(data))) // bytes 2-9: data-length (8 bytes)

	if _, err := w.writer.Write(header); err != nil {
		return err
	}

	// Write segment data
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	// Write segment CRC64 AFTER the data (8 bytes) only if CRC64 is enabled
	if w.enableCRC64 {
		crcBytes := make([]byte, SegmentCRC64Size)
		binary.LittleEndian.PutUint64(crcBytes, segmentCRC64)
		if _, err := w.writer.Write(crcBytes); err != nil {
			return err
		}
	}

	// Update cumulative CRC64 if enabled
	if w.enableCRC64 {
		w.totalCRC64 = crc64.Update(w.totalCRC64, w.crc64Table, data)
	}
	w.segmentNum++

	return nil
}

// WriteTrailer writes the message trailer with cumulative CRC64 (if CRC64 is enabled)
// If CRC64 is enabled: Trailer is 8 bytes (message-crc64)
// If CRC64 is disabled: No trailer is written
func (w *StructuredMessageWriter) WriteTrailer() error {
	if !w.enableCRC64 {
		// No trailer when CRC64 is disabled
		return nil
	}

	trailer := make([]byte, TrailerSize)
	binary.LittleEndian.PutUint64(trailer[0:8], w.totalCRC64)
	_, err := w.writer.Write(trailer)
	return err
}

// GetTotalCRC64 returns the cumulative CRC64 of all segments written
func (w *StructuredMessageWriter) GetTotalCRC64() uint64 {
	return w.totalCRC64
}

// StructuredMessageReader decodes the binary message format
type StructuredMessageReader struct {
	reader     io.Reader
	crc64Table *crc64.Table
	header     *MessageHeader
	segmentNum uint16
	totalCRC64 uint64
	closed     bool
}

// MessageHeader represents the message header
type MessageHeader struct {
	Version       uint8
	MessageLength uint64 // 8 bytes - total SM length
	MessageFlags  uint16 // 2 bytes - bitmask (CRC64 = 0x0001)
	NumSegments   uint16 // 2 bytes - number of segments
}

// NewStructuredMessageReader creates a new StructuredMessageReader
func NewStructuredMessageReader(reader io.Reader) *StructuredMessageReader {
	return &StructuredMessageReader{
		reader:     reader,
		crc64Table: shared.CRC64Table,
		segmentNum: 1, // Segments are 1-indexed
		totalCRC64: 0,
		closed:     false,
	}
}

// ReadHeader reads and validates the message header
func (r *StructuredMessageReader) ReadHeader() (*MessageHeader, error) {
	if r.header != nil {
		return r.header, nil
	}

	headerBytes := make([]byte, HeaderSize)
	n, err := io.ReadFull(r.reader, headerBytes)
	if err != nil {
		if err == io.EOF {
			return nil, ErrUnexpectedEOF
		}
		return nil, err
	}
	if n != HeaderSize {
		return nil, ErrUnexpectedEOF
	}

	header := &MessageHeader{
		Version:       headerBytes[0],
		MessageLength: binary.LittleEndian.Uint64(headerBytes[1:9]),
		MessageFlags:  binary.LittleEndian.Uint16(headerBytes[9:11]),
		NumSegments:   binary.LittleEndian.Uint16(headerBytes[11:13]),
	}

	if header.Version != MessageVersion {
		return nil, ErrInvalidVersion
	}

	// Flags can be 0 (CRC64 disabled) or 1 (CRC64 enabled)
	// Both are valid, so we don't require PropCRC64 to be set

	if header.NumSegments == 0 || header.NumSegments > MaxSegments {
		return nil, ErrInvalidSegmentCount
	}

	r.header = header
	return header, nil
}

// ReadSegment reads the next segment and validates its CRC64
func (r *StructuredMessageReader) ReadSegment() ([]byte, error) {
	if r.closed {
		return nil, io.EOF
	}

	if r.header == nil {
		if _, err := r.ReadHeader(); err != nil {
			return nil, err
		}
	}

	if r.segmentNum > r.header.NumSegments {
		return nil, io.EOF
	}

	// Read segment header (segment-num + data-length)
	headerBytes := make([]byte, SegmentHeaderSize)
	n, err := io.ReadFull(r.reader, headerBytes)
	if err != nil {
		if err == io.EOF {
			return nil, ErrUnexpectedEOF
		}
		return nil, err
	}
	if n != SegmentHeaderSize {
		return nil, ErrUnexpectedEOF
	}

	segmentNum := binary.LittleEndian.Uint16(headerBytes[0:2])
	dataLength := binary.LittleEndian.Uint64(headerBytes[2:10])

	if segmentNum != r.segmentNum {
		return nil, ErrInvalidSegmentNumber
	}

	if dataLength > MaxSegmentSize {
		return nil, ErrSegmentTooLarge
	}

	// Read segment data
	data := make([]byte, dataLength)
	n, err = io.ReadFull(r.reader, data)
	if err != nil {
		if err == io.EOF {
			return nil, ErrUnexpectedEOF
		}
		return nil, err
	}
	if n != int(dataLength) {
		return nil, ErrUnexpectedEOF
	}

	// Read and validate segment CRC64 (8 bytes AFTER the data) only if CRC64 is enabled
	enableCRC64 := (r.header.MessageFlags & PropCRC64) != 0
	if enableCRC64 {
		crcBytes := make([]byte, SegmentCRC64Size)
		n, err = io.ReadFull(r.reader, crcBytes)
		if err != nil {
			if err == io.EOF {
				return nil, ErrUnexpectedEOF
			}
			return nil, err
		}
		if n != SegmentCRC64Size {
			return nil, ErrUnexpectedEOF
		}
		expectedCRC64 := binary.LittleEndian.Uint64(crcBytes)

		// Validate segment CRC64
		actualCRC64 := crc64.Checksum(data, r.crc64Table)
		if actualCRC64 != expectedCRC64 {
			return nil, ErrCRC64Mismatch
		}

		// Update cumulative CRC64
		r.totalCRC64 = crc64.Update(r.totalCRC64, r.crc64Table, data)
	}

	r.segmentNum++

	return data, nil
}

// ReadAllSegments reads all remaining segments and returns them as a single byte slice
func (r *StructuredMessageReader) ReadAllSegments() ([]byte, error) {
	var allData []byte

	for {
		data, err := r.ReadSegment()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		allData = append(allData, data...)
	}

	return allData, nil
}

// ReadTrailer reads and validates the message trailer
func (r *StructuredMessageReader) ReadTrailer() (uint64, error) {
	if r.closed {
		return 0, io.EOF
	}

	// Validate segment count before reading trailer
	// With 1-indexed segments, after reading all segments, segmentNum = NumSegments + 1
	if r.segmentNum != r.header.NumSegments+1 {
		return 0, ErrInvalidSegmentCount
	}

	// Check if CRC64 is enabled
	enableCRC64 := (r.header.MessageFlags & PropCRC64) != 0
	if !enableCRC64 {
		// No trailer when CRC64 is disabled
		r.closed = true
		return 0, nil
	}

	// Trailer is ONLY 8 bytes: message-crc64
	trailerBytes := make([]byte, TrailerSize)
	n, err := io.ReadFull(r.reader, trailerBytes)
	if err != nil {
		if err == io.EOF {
			return 0, ErrUnexpectedEOF
		}
		return 0, err
	}
	if n != TrailerSize {
		return 0, ErrUnexpectedEOF
	}

	expectedCRC64 := binary.LittleEndian.Uint64(trailerBytes[0:8])
	if expectedCRC64 != r.totalCRC64 {
		return 0, ErrCRC64Mismatch
	}

	r.closed = true
	return expectedCRC64, nil
}

// GetTotalCRC64 returns the cumulative CRC64 of all segments read
func (r *StructuredMessageReader) GetTotalCRC64() uint64 {
	return r.totalCRC64
}

// Close marks the reader as closed
func (r *StructuredMessageReader) Close() error {
	r.closed = true
	return nil
}

// IsClosed returns true if the reader is closed
func (r *StructuredMessageReader) IsClosed() bool {
	return r.closed
}

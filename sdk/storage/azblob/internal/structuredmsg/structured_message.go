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

	// HeaderSize is the Header size in bytes: version(1) + flags(1) + properties(2) + numSegments(4) = 8 bytes
	HeaderSize = 8

	// TrailerHeaderSize is the Trailer header size in bytes: version(1) + flags(1) + properties(2) = 4 bytes
	TrailerHeaderSize = 4

	// TrailerCRC64Size is the Trailer CRC64 size in bytes: cumulativeCRC64(8) = 8 bytes
	TrailerCRC64Size = 8

	// TrailerSize is the Total trailer size in bytes: header(4) + crc64(8) = 12 bytes
	TrailerSize = TrailerHeaderSize + TrailerCRC64Size

	// SegmentHeaderSize is segmentNum(4) + dataLength(4) + segmentCRC64(8) = 16 bytes
	SegmentHeaderSize = 16

	// MaxSegments is Maximum number of segments allowed
	MaxSegments = math.MaxUint32

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
	writer     io.Writer
	crc64Table *crc64.Table
	segmentNum uint32
	totalCRC64 uint64
}

// NewStructuredMessageWriter creates a new StructuredMessageWriter
func NewStructuredMessageWriter(writer io.Writer) *StructuredMessageWriter {
	return &StructuredMessageWriter{
		writer:     writer,
		crc64Table: shared.CRC64Table,
		segmentNum: 0,
		totalCRC64: 0,
	}
}

// WriteHeader writes the message header
func (w *StructuredMessageWriter) WriteHeader(numSegments uint32) error {
	if numSegments == 0 || numSegments > MaxSegments {
		return ErrInvalidSegmentCount
	}

	header := make([]byte, HeaderSize)
	header[0] = MessageVersion
	header[1] = FlagNone
	binary.LittleEndian.PutUint16(header[2:], PropCRC64)
	binary.LittleEndian.PutUint32(header[4:], numSegments)

	_, err := w.writer.Write(header)
	return err
}

// WriteSegment writes a segment with its data and CRC64
func (w *StructuredMessageWriter) WriteSegment(data []byte) error {
	if len(data) > MaxSegmentSize {
		return ErrSegmentTooLarge
	}

	// Calculate segment CRC64
	segmentCRC64 := crc64.Checksum(data, w.crc64Table)

	// Write segment header
	header := make([]byte, SegmentHeaderSize)
	binary.LittleEndian.PutUint32(header[0:4], w.segmentNum)
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(data)))
	binary.LittleEndian.PutUint64(header[8:16], segmentCRC64)

	if _, err := w.writer.Write(header); err != nil {
		return err
	}

	// Write segment data
	if _, err := w.writer.Write(data); err != nil {
		return err
	}

	// Update cumulative CRC64
	w.totalCRC64 = crc64.Update(w.totalCRC64, w.crc64Table, data)
	w.segmentNum++

	return nil
}

// WriteTrailer writes the message trailer with cumulative CRC64
func (w *StructuredMessageWriter) WriteTrailer() error {
	trailer := make([]byte, TrailerSize)
	trailer[0] = MessageVersion
	trailer[1] = FlagNone
	binary.LittleEndian.PutUint16(trailer[2:], PropCRC64)
	binary.LittleEndian.PutUint64(trailer[4:], w.totalCRC64)
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
	segmentNum uint32
	totalCRC64 uint64
	closed     bool
}

// MessageHeader represents the message header
type MessageHeader struct {
	Version     uint8
	Flags       uint8
	Properties  uint16
	NumSegments uint32
}

// NewStructuredMessageReader creates a new StructuredMessageReader
func NewStructuredMessageReader(reader io.Reader) *StructuredMessageReader {
	return &StructuredMessageReader{
		reader:     reader,
		crc64Table: shared.CRC64Table,
		segmentNum: 0,
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
		Version:     headerBytes[0],
		Flags:       headerBytes[1],
		Properties:  binary.LittleEndian.Uint16(headerBytes[2:]),
		NumSegments: binary.LittleEndian.Uint32(headerBytes[4:]),
	}

	if header.Version != MessageVersion {
		return nil, ErrInvalidVersion
	}

	if header.Properties&PropCRC64 == 0 {
		return nil, errors.New("CRC64 property not set")
	}

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

	if r.segmentNum >= r.header.NumSegments {
		return nil, io.EOF
	}

	// Read segment header
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

	segmentNum := binary.LittleEndian.Uint32(headerBytes[0:])
	dataLength := binary.LittleEndian.Uint32(headerBytes[4:])
	expectedCRC64 := binary.LittleEndian.Uint64(headerBytes[8:])

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

	// Validate segment CRC64
	actualCRC64 := crc64.Checksum(data, r.crc64Table)
	if actualCRC64 != expectedCRC64 {
		return nil, ErrCRC64Mismatch
	}

	// Update cumulative CRC64
	r.totalCRC64 = crc64.Update(r.totalCRC64, r.crc64Table, data)
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
	if r.segmentNum != r.header.NumSegments {
		return 0, ErrInvalidSegmentCount
	}

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

	// Validate trailer header
	if trailerBytes[0] != MessageVersion {
		return 0, ErrInvalidVersion
	}
	if trailerBytes[1] != FlagNone {
		return 0, errors.New("invalid trailer flags")
	}
	properties := binary.LittleEndian.Uint16(trailerBytes[2:])
	if properties&PropCRC64 == 0 {
		return 0, errors.New("CRC64 property not set in trailer")
	}

	expectedCRC64 := binary.LittleEndian.Uint64(trailerBytes[4:])
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

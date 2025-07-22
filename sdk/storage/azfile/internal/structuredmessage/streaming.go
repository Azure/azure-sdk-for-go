//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package structuredmessage

import (
	"bytes"
	"encoding/binary"
	"hash"
	"hash/crc64"
	"io"
)

// EncoderReader wraps a reader and provides structured message encoding with CRC64 validation
type EncoderReader struct {
	originalData   []byte
	encodedData    []byte
	reader         io.Reader
	contentLength  int64
	encodedLength  int64
}

// NewEncoderReader creates a new EncoderReader that encodes the provided data into structured message format
func NewEncoderReader(reader io.ReadSeeker) (*EncoderReader, error) {
	// Ensure we start from the beginning
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// Read all data first (this is required for CRC64 computation)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Reset reader position again (not needed since we read all, but good practice)
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// Encode the data
	encoded, err := EncodeMessage(data)
	if err != nil {
		return nil, err
	}

	return &EncoderReader{
		originalData:  data,
		encodedData:   encoded,
		reader:        bytes.NewReader(encoded),
		contentLength: int64(len(data)),
		encodedLength: int64(len(encoded)),
	}, nil
}

// Read implements io.Reader for the encoded data
func (er *EncoderReader) Read(p []byte) (n int, err error) {
	return er.reader.Read(p)
}

// ContentLength returns the original content length (before encoding)
func (er *EncoderReader) ContentLength() int64 {
	return er.contentLength
}

// EncodedLength returns the encoded message length
func (er *EncoderReader) EncodedLength() int64 {
	return er.encodedLength
}

// GetStructuredHeaders returns the headers needed for structured message requests
func (er *EncoderReader) GetStructuredHeaders() (structuredBodyType string, structuredContentLength int64) {
	return "XSM/1.0; properties=crc64", er.contentLength
}

// DecoderReader wraps a reader and provides structured message decoding with CRC64 validation
type DecoderReader struct {
	originalReader  io.Reader
	decodedData     []byte
	reader          io.Reader
	contentLength   int64
	originalLength  int64
	validated       bool
}

// NewDecoderReader creates a new DecoderReader that decodes structured message format
func NewDecoderReader(reader io.Reader, structuredContentLength int64) (*DecoderReader, error) {
	// Read all encoded data
	encodedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Decode and validate the structured message
	decodedData, err := DecodeMessage(encodedData)
	if err != nil {
		return nil, err
	}

	// Validate that the decoded length matches the expected content length
	if int64(len(decodedData)) != structuredContentLength {
		return nil, io.ErrUnexpectedEOF
	}

	return &DecoderReader{
		originalReader: reader,
		decodedData:    decodedData,
		reader:         bytes.NewReader(decodedData),
		contentLength:  int64(len(decodedData)),
		originalLength: int64(len(encodedData)),
		validated:      true,
	}, nil
}

// Read implements io.Reader for the decoded data
func (dr *DecoderReader) Read(p []byte) (n int, err error) {
	return dr.reader.Read(p)
}

// ContentLength returns the decoded content length
func (dr *DecoderReader) ContentLength() int64 {
	return dr.contentLength
}

// OriginalLength returns the original encoded message length
func (dr *DecoderReader) OriginalLength() int64 {
	return dr.originalLength
}

// IsValidated returns whether the content has been successfully validated
func (dr *DecoderReader) IsValidated() bool {
	return dr.validated
}

// StreamingEncoder provides streaming encoding with CRC64 for upload scenarios
type StreamingEncoder struct {
	data           []byte
	crc64Hasher    hash.Hash64
	buffer         *bytes.Buffer
	headerWritten  bool
	dataWritten    bool
	trailerWritten bool
	contentLength  int64
	segmentNum     uint16
}

// NewStreamingEncoder creates a new streaming encoder for structured message format
func NewStreamingEncoder(contentLength int64) *StreamingEncoder {
	hasher := crc64.New(CRC64Table)
	return &StreamingEncoder{
		crc64Hasher:   hasher,
		buffer:        &bytes.Buffer{},
		contentLength: contentLength,
		segmentNum:    1,
	}
}

// WriteData writes data to the encoder and updates CRC64
func (se *StreamingEncoder) WriteData(data []byte) error {
	se.data = append(se.data, data...)
	_, err := se.crc64Hasher.Write(data)
	return err
}

// GetEncodedData returns the complete encoded structured message
func (se *StreamingEncoder) GetEncodedData() ([]byte, error) {
	if !se.trailerWritten {
		if err := se.writeAll(); err != nil {
			return nil, err
		}
	}
	return se.buffer.Bytes(), nil
}

func (se *StreamingEncoder) writeAll() error {
	// Calculate final CRC64
	finalCRC64 := se.crc64Hasher.Sum64()

	// Calculate total message length
	headerSize := uint64(1 + 8 + 2 + 2) // version + length + flags + segments
	segmentHeaderSize := uint64(2 + 8)  // segment num + data length
	segmentDataSize := uint64(len(se.data))
	segmentCRC64Size := uint64(CRC64Size)
	messageCRC64Size := uint64(CRC64Size)
	totalLength := headerSize + segmentHeaderSize + segmentDataSize + segmentCRC64Size + messageCRC64Size

	// Write header
	header := Header{
		MessageVersion: MessageVersion,
		MessageLength:  totalLength,
		MessageFlags:   FlagIncludeCRC64,
		NumSegments:    1,
	}

	if err := binary.Write(se.buffer, binary.LittleEndian, header.MessageVersion); err != nil {
		return err
	}
	if err := binary.Write(se.buffer, binary.LittleEndian, header.MessageLength); err != nil {
		return err
	}
	if err := binary.Write(se.buffer, binary.LittleEndian, header.MessageFlags); err != nil {
		return err
	}
	if err := binary.Write(se.buffer, binary.LittleEndian, header.NumSegments); err != nil {
		return err
	}

	// Write segment
	if err := binary.Write(se.buffer, binary.LittleEndian, se.segmentNum); err != nil {
		return err
	}
	if err := binary.Write(se.buffer, binary.LittleEndian, uint64(len(se.data))); err != nil {
		return err
	}
	if _, err := se.buffer.Write(se.data); err != nil {
		return err
	}

	// Write segment CRC64
	segmentCRC64Bytes := make([]byte, CRC64Size)
	binary.LittleEndian.PutUint64(segmentCRC64Bytes, finalCRC64)
	if _, err := se.buffer.Write(segmentCRC64Bytes); err != nil {
		return err
	}

	// Write message CRC64 (trailer)
	messageCRC64Bytes := make([]byte, CRC64Size)
	binary.LittleEndian.PutUint64(messageCRC64Bytes, finalCRC64)
	if _, err := se.buffer.Write(messageCRC64Bytes); err != nil {
		return err
	}

	se.trailerWritten = true
	return nil
}

// GetContentLength returns the original content length
func (se *StreamingEncoder) GetContentLength() int64 {
	return se.contentLength
}

// GetStructuredHeaders returns the headers needed for structured message requests
func (se *StreamingEncoder) GetStructuredHeaders() (structuredBodyType string, structuredContentLength int64) {
	return "XSM/1.0; properties=crc64", se.contentLength
}
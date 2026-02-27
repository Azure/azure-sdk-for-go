//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Package structuredmessage provides encoding and decoding for Azure Storage Structured Message v1 format.
// This format is used for CRC64 content validation in upload and download operations.
package structuredmessage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc64"
	"io"
)

const (
	// MessageVersion is the structured message version supported
	MessageVersion uint8 = 1
	// FlagIncludeCRC64 indicates that CRC64 checksums should be included
	FlagIncludeCRC64 uint16 = 0x0001
	// CRC64Size is the size of a CRC64 checksum in bytes
	CRC64Size = 8
)

// CRC64Polynomial is the custom polynomial used by Azure Storage for CRC64 computation
const CRC64Polynomial uint64 = 0x9A6C9329AC4BC9B5

// CRC64Table is the lookup table for the Azure Storage CRC64 polynomial
var CRC64Table = crc64.MakeTable(CRC64Polynomial)

// Header represents the structured message header
type Header struct {
	MessageVersion     uint8
	MessageLength      uint64
	MessageFlags       uint16
	NumSegments        uint16
}

// Segment represents a structured message segment
type Segment struct {
	SegmentNum        uint16
	SegmentDataLength uint64
	SegmentData       []byte
	SegmentDataCRC64  []byte // Optional, 8 bytes when present
}

// Message represents a complete structured message
type Message struct {
	Header          Header
	Segments        []Segment
	MessageDataCRC64 []byte // Optional, 8 bytes when present
}

// EncodeMessage encodes data into a structured message v1 format with CRC64 validation
func EncodeMessage(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	// Create single segment for the data
	segmentCRC64 := crc64.Checksum(data, CRC64Table)
	segmentCRC64Bytes := make([]byte, CRC64Size)
	binary.LittleEndian.PutUint64(segmentCRC64Bytes, segmentCRC64)

	segment := Segment{
		SegmentNum:        1,
		SegmentDataLength: uint64(len(data)),
		SegmentData:       data,
		SegmentDataCRC64:  segmentCRC64Bytes,
	}

	// Create message-level CRC64 (same as segment CRC64 for single segment)
	messageCRC64Bytes := make([]byte, CRC64Size)
	binary.LittleEndian.PutUint64(messageCRC64Bytes, segmentCRC64)

	// Calculate total message length
	headerSize := uint64(1 + 8 + 2 + 2) // version + length + flags + segments
	segmentHeaderSize := uint64(2 + 8)  // segment num + data length
	segmentDataSize := uint64(len(data))
	segmentCRC64Size := uint64(CRC64Size)
	messageCRC64Size := uint64(CRC64Size)
	
	totalLength := headerSize + segmentHeaderSize + segmentDataSize + segmentCRC64Size + messageCRC64Size

	header := Header{
		MessageVersion: MessageVersion,
		MessageLength:  totalLength,
		MessageFlags:   FlagIncludeCRC64,
		NumSegments:    1,
	}

	message := Message{
		Header:           header,
		Segments:         []Segment{segment},
		MessageDataCRC64: messageCRC64Bytes,
	}

	return encodeMessageToBytes(message)
}

// DecodeMessage decodes a structured message v1 format and validates CRC64 checksums
func DecodeMessage(encodedData []byte) ([]byte, error) {
	if len(encodedData) < 13 { // Minimum header size
		return nil, errors.New("data too short to be a valid structured message")
	}

	reader := bytes.NewReader(encodedData)
	
	// Read header
	header, err := readHeader(reader)
	if err != nil {
		return nil, err
	}

	if header.MessageVersion != MessageVersion {
		return nil, errors.New("unsupported message version")
	}

	if header.MessageLength != uint64(len(encodedData)) {
		return nil, errors.New("message length mismatch")
	}

	// Read segments
	segments := make([]Segment, header.NumSegments)
	var allData []byte
	
	for i := uint16(0); i < header.NumSegments; i++ {
		segment, err := readSegment(reader, header.MessageFlags&FlagIncludeCRC64 != 0)
		if err != nil {
			return nil, err
		}
		
		if segment.SegmentNum != i+1 {
			return nil, errors.New("invalid segment number")
		}
		
		// Validate segment CRC64 if present
		if header.MessageFlags&FlagIncludeCRC64 != 0 && len(segment.SegmentDataCRC64) == CRC64Size {
			expectedCRC64 := binary.LittleEndian.Uint64(segment.SegmentDataCRC64)
			actualCRC64 := crc64.Checksum(segment.SegmentData, CRC64Table)
			if expectedCRC64 != actualCRC64 {
				return nil, errors.New("segment CRC64 validation failed")
			}
		}
		
		segments[i] = segment
		allData = append(allData, segment.SegmentData...)
	}

	// Read trailer (message CRC64) if present
	if header.MessageFlags&FlagIncludeCRC64 != 0 {
		messageCRC64 := make([]byte, CRC64Size)
		n, err := reader.Read(messageCRC64)
		if err != nil || n != CRC64Size {
			return nil, errors.New("failed to read message CRC64")
		}
		
		// Validate message CRC64
		expectedCRC64 := binary.LittleEndian.Uint64(messageCRC64)
		actualCRC64 := crc64.Checksum(allData, CRC64Table)
		if expectedCRC64 != actualCRC64 {
			return nil, errors.New("message CRC64 validation failed")
		}
	}

	return allData, nil
}

func encodeMessageToBytes(msg Message) ([]byte, error) {
	var buf bytes.Buffer

	// Write header
	if err := binary.Write(&buf, binary.LittleEndian, msg.Header.MessageVersion); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, msg.Header.MessageLength); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, msg.Header.MessageFlags); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, msg.Header.NumSegments); err != nil {
		return nil, err
	}

	// Write segments
	for _, segment := range msg.Segments {
		if err := binary.Write(&buf, binary.LittleEndian, segment.SegmentNum); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.LittleEndian, segment.SegmentDataLength); err != nil {
			return nil, err
		}
		if _, err := buf.Write(segment.SegmentData); err != nil {
			return nil, err
		}
		if len(segment.SegmentDataCRC64) > 0 {
			if _, err := buf.Write(segment.SegmentDataCRC64); err != nil {
				return nil, err
			}
		}
	}

	// Write trailer
	if len(msg.MessageDataCRC64) > 0 {
		if _, err := buf.Write(msg.MessageDataCRC64); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func readHeader(reader io.Reader) (Header, error) {
	var header Header
	
	if err := binary.Read(reader, binary.LittleEndian, &header.MessageVersion); err != nil {
		return header, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &header.MessageLength); err != nil {
		return header, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &header.MessageFlags); err != nil {
		return header, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &header.NumSegments); err != nil {
		return header, err
	}
	
	return header, nil
}

func readSegment(reader io.Reader, includeCRC64 bool) (Segment, error) {
	var segment Segment
	
	if err := binary.Read(reader, binary.LittleEndian, &segment.SegmentNum); err != nil {
		return segment, err
	}
	if err := binary.Read(reader, binary.LittleEndian, &segment.SegmentDataLength); err != nil {
		return segment, err
	}
	
	// Read segment data
	segment.SegmentData = make([]byte, segment.SegmentDataLength)
	n, err := reader.Read(segment.SegmentData)
	if err != nil || uint64(n) != segment.SegmentDataLength {
		return segment, errors.New("failed to read segment data")
	}
	
	// Read segment CRC64 if present
	if includeCRC64 {
		segment.SegmentDataCRC64 = make([]byte, CRC64Size)
		n, err := reader.Read(segment.SegmentDataCRC64)
		if err != nil || n != CRC64Size {
			return segment, errors.New("failed to read segment CRC64")
		}
	}
	
	return segment, nil
}
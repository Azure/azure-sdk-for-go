//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package structuredmsg

import (
	"bytes"
	"encoding/binary"
	"hash/crc64"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// calculateMessageLength calculates the total message length for given segments
// enableCRC64: if true, includes SegmentCRC64Size for each segment and TrailerSize
func calculateMessageLength(segmentSizes []int, enableCRC64 bool) uint64 {
	messageLength := uint64(HeaderSize)
	for _, size := range segmentSizes {
		messageLength += uint64(SegmentHeaderSize) + uint64(size)
		if enableCRC64 {
			messageLength += uint64(SegmentCRC64Size)
		}
	}
	if enableCRC64 {
		messageLength += uint64(TrailerSize)
	}
	return messageLength
}

func TestStructuredMessageWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name          string
		numSegments   uint16
		messageLength uint64
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "Valid single segment",
			numSegments:   1,
			messageLength: calculateMessageLength([]int{100}, true),
			expectError:   false,
		},
		{
			name:          "Valid multiple segments",
			numSegments:   5,
			messageLength: calculateMessageLength([]int{100, 200, 300, 400, 500}, true),
			expectError:   false,
		},
		{
			name:          "Valid maximum segments",
			numSegments:   MaxSegments,
			messageLength: calculateMessageLength([]int{100}, true), // Just a placeholder
			expectError:   false,
		},
		{
			name:          "Zero segments",
			numSegments:   0,
			messageLength: 0,
			expectError:   true,
			errorMessage:  "invalid segment count",
		},
		{
			name:          "Too many segments",
			numSegments:   MaxSegments,
			messageLength: calculateMessageLength([]int{100}, true),
			expectError:   false, // MaxSegments should be valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewStructuredMessageWriter(&buf)

			// Default to CRC64 enabled for tests
			err := writer.WriteHeader(tt.numSegments, tt.messageLength, true)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
			} else {
				require.NoError(t, err)
				require.Equal(t, HeaderSize, buf.Len())

				// Verify header content
				headerBytes := buf.Bytes()
				assert.Equal(t, uint8(MessageVersion), headerBytes[0])
				assert.Equal(t, tt.messageLength, binary.LittleEndian.Uint64(headerBytes[1:9]))
				assert.Equal(t, uint16(PropCRC64), binary.LittleEndian.Uint16(headerBytes[9:11]))
				assert.Equal(t, tt.numSegments, binary.LittleEndian.Uint16(headerBytes[11:13]))
			}
		})
	}
}

func TestStructuredMessageWriter_WriteSegment(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Empty segment",
			data:        []byte{},
			expectError: false,
		},
		{
			name:        "Small segment",
			data:        []byte("hello"),
			expectError: false,
		},
		{
			name:        "Large segment",
			data:        make([]byte, 1024*1024), // 1MB
			expectError: false,
		},
		{
			name:         "Segment too large",
			data:         make([]byte, MaxSegmentSize+1),
			expectError:  true,
			errorMessage: "segment too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewStructuredMessageWriter(&buf)

			// Write header first
			messageLength := calculateMessageLength([]int{len(tt.data)}, true)
			err := writer.WriteHeader(1, messageLength, true)
			require.NoError(t, err)

			// Write segment
			err = writer.WriteSegment(tt.data)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
			} else {
				require.NoError(t, err)

				// Verify segment was written correctly
				expectedSize := HeaderSize + SegmentHeaderSize + len(tt.data) + SegmentCRC64Size
				require.Equal(t, expectedSize, buf.Len())

				// Verify segment header
				segmentHeaderStart := HeaderSize
				segmentHeaderBytes := buf.Bytes()[segmentHeaderStart : segmentHeaderStart+SegmentHeaderSize]
				segmentNum := binary.LittleEndian.Uint16(segmentHeaderBytes[0:2])
				dataLength := binary.LittleEndian.Uint64(segmentHeaderBytes[2:10])

				assert.Equal(t, uint16(1), segmentNum) // Segments are 1-indexed
				assert.Equal(t, uint64(len(tt.data)), dataLength)

				// Verify data
				dataStart := segmentHeaderStart + SegmentHeaderSize
				actualData := buf.Bytes()[dataStart : dataStart+len(tt.data)]
				assert.Equal(t, tt.data, actualData)

				// Verify CRC64 (after data)
				crcStart := dataStart + len(tt.data)
				crcBytes := buf.Bytes()[crcStart : crcStart+SegmentCRC64Size]
				segmentCRC64 := binary.LittleEndian.Uint64(crcBytes)
				expectedCRC64 := crc64.Checksum(tt.data, shared.CRC64Table)
				assert.Equal(t, expectedCRC64, segmentCRC64)
			}
		})
	}
}

func TestStructuredMessageWriter_WriteTrailer(t *testing.T) {
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)

	data1 := []byte("hello")
	data2 := []byte("world")

	// Write header and segments
	messageLength := calculateMessageLength([]int{len(data1), len(data2)}, true)
	err := writer.WriteHeader(2, messageLength, true)
	require.NoError(t, err)

	err = writer.WriteSegment(data1)
	require.NoError(t, err)

	err = writer.WriteSegment(data2)
	require.NoError(t, err)

	// Write trailer
	err = writer.WriteTrailer()
	require.NoError(t, err)

	// Verify trailer
	expectedSize := HeaderSize + (SegmentHeaderSize + len(data1) + SegmentCRC64Size) + (SegmentHeaderSize + len(data2) + SegmentCRC64Size) + TrailerSize
	require.Equal(t, expectedSize, buf.Len())

	// Verify trailer content (trailer is ONLY 8 bytes: message-crc64)
	trailerStart := expectedSize - TrailerSize
	trailerBytes := buf.Bytes()[trailerStart:]

	expectedCRC64 := binary.LittleEndian.Uint64(trailerBytes[0:8])

	// Calculate expected cumulative CRC64
	allData := append(data1, data2...)
	expectedCumulativeCRC64 := crc64.Checksum(allData, shared.CRC64Table)
	assert.Equal(t, expectedCumulativeCRC64, expectedCRC64)
	assert.Equal(t, expectedCumulativeCRC64, writer.GetTotalCRC64())
}

func TestStructuredMessageReader_ReadHeader(t *testing.T) {
	// Helper to create valid 13-byte header
	createHeader := func(version uint8, messageLength uint64, flags uint16, numSegments uint16) []byte {
		header := make([]byte, HeaderSize)
		header[0] = version
		binary.LittleEndian.PutUint64(header[1:9], messageLength)
		binary.LittleEndian.PutUint16(header[9:11], flags)
		binary.LittleEndian.PutUint16(header[11:13], numSegments)
		return header
	}

	tests := []struct {
		name         string
		headerBytes  []byte
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid header",
			headerBytes: createHeader(MessageVersion, 100, PropCRC64, 1),
			expectError: false,
		},
		{
			name:         "Invalid version",
			headerBytes:  createHeader(2, 100, PropCRC64, 1), // version=2
			expectError:  true,
			errorMessage: "invalid message version",
		},
		{
			name:         "Zero segments",
			headerBytes:  createHeader(MessageVersion, 100, PropCRC64, 0), // numSegments=0
			expectError:  true,
			errorMessage: "invalid segment count",
		},
		{
			name:         "Missing CRC64 property",
			headerBytes:  createHeader(MessageVersion, 100, 0, 1), // flags=0 (no CRC64)
			expectError:  true,
			errorMessage: "CRC64 property not set",
		},
		{
			name:        "Maximum segments",
			headerBytes: createHeader(MessageVersion, 100, PropCRC64, MaxSegments),
			expectError: false,
		},
		{
			name:         "Incomplete header",
			headerBytes:  []byte{MessageVersion, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, 0x00, 1}, // missing 1 byte
			expectError:  true,
			errorMessage: "unexpected EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewStructuredMessageReader(bytes.NewReader(tt.headerBytes))

			header, err := reader.ReadHeader()

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
				require.Nil(t, header)
			} else {
				require.NoError(t, err)
				require.NotNil(t, header)
				assert.Equal(t, uint8(MessageVersion), header.Version)
				if tt.name == "Maximum segments" {
					assert.Equal(t, uint16(MaxSegments), header.NumSegments)
				} else {
					assert.Equal(t, uint16(1), header.NumSegments)
				}
			}
		})
	}
}

func TestStructuredMessageReader_ReadSegment(t *testing.T) {
	// Create test data
	data1 := []byte("hello")
	data2 := []byte("world")

	// Create structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)
	messageLength := calculateMessageLength([]int{len(data1), len(data2)}, true)
	writer.WriteHeader(2, messageLength, true)
	writer.WriteSegment(data1)
	writer.WriteSegment(data2)
	writer.WriteTrailer()

	// Test reading segments
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read first segment
	segment1, err := reader.ReadSegment()
	require.NoError(t, err)
	assert.Equal(t, data1, segment1)

	// Read second segment
	segment2, err := reader.ReadSegment()
	require.NoError(t, err)
	assert.Equal(t, data2, segment2)

	// Try to read beyond available segments
	_, err = reader.ReadSegment()
	assert.Equal(t, io.EOF, err)
}

func TestStructuredMessageReader_ReadSegment_CRC64Mismatch(t *testing.T) {
	// Create test data
	data := []byte("hello")

	// Create structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)
	messageLength := calculateMessageLength([]int{len(data)}, true)
	writer.WriteHeader(1, messageLength, true)
	writer.WriteSegment(data)
	writer.WriteTrailer()

	// Corrupt the CRC64 after the segment data
	bufBytes := buf.Bytes()
	segmentHeaderStart := HeaderSize
	dataStart := segmentHeaderStart + SegmentHeaderSize
	dataEnd := dataStart + len(data)
	crc64Offset := dataEnd        // CRC64 is after the data
	bufBytes[crc64Offset] ^= 0xFF // Flip some bits

	// Test reading with corrupted CRC64
	reader := NewStructuredMessageReader(bytes.NewReader(bufBytes))

	_, err := reader.ReadSegment()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CRC64 checksum mismatch")
}

func TestStructuredMessageReader_ReadAllSegments(t *testing.T) {
	// Create test data
	data1 := []byte("hello")
	data2 := []byte("world")
	data3 := []byte("test")

	// Create structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)
	messageLength := calculateMessageLength([]int{len(data1), len(data2), len(data3)}, true)
	writer.WriteHeader(3, messageLength, true)
	writer.WriteSegment(data1)
	writer.WriteSegment(data2)
	writer.WriteSegment(data3)
	writer.WriteTrailer()

	// Test reading all segments
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	allData, err := reader.ReadAllSegments()
	require.NoError(t, err)

	expectedData := append(append(data1, data2...), data3...)
	assert.Equal(t, expectedData, allData)
}

func TestStructuredMessageReader_ReadTrailer(t *testing.T) {
	// Create test data
	data1 := []byte("hello")
	data2 := []byte("world")

	// Create structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)
	messageLength := calculateMessageLength([]int{len(data1), len(data2)}, true)
	writer.WriteHeader(2, messageLength, true)
	writer.WriteSegment(data1)
	writer.WriteSegment(data2)
	writer.WriteTrailer()

	// Test reading trailer
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read all segments first
	_, err := reader.ReadAllSegments()
	require.NoError(t, err)

	// Read trailer
	trailerCRC64, err := reader.ReadTrailer()
	require.NoError(t, err)

	// Verify CRC64
	allData := append(data1, data2...)
	expectedCRC64 := crc64.Checksum(allData, shared.CRC64Table)
	assert.Equal(t, expectedCRC64, trailerCRC64)
	assert.Equal(t, expectedCRC64, reader.GetTotalCRC64())
}

func TestStructuredMessageReader_ReadTrailer_CRC64Mismatch(t *testing.T) {
	// Create test data
	data := []byte("hello")

	// Create structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)
	messageLength := calculateMessageLength([]int{len(data)}, true)
	writer.WriteHeader(1, messageLength, true)
	writer.WriteSegment(data)
	writer.WriteTrailer()

	// Corrupt the trailer CRC64
	bufBytes := buf.Bytes()
	trailerStart := len(bufBytes) - TrailerSize
	crc64Offset := trailerStart   // CRC64 is at offset 0 in trailer (trailer is only CRC64)
	bufBytes[crc64Offset] ^= 0xFF // Flip some bits in CRC64

	// Test reading with corrupted trailer
	reader := NewStructuredMessageReader(bytes.NewReader(bufBytes))

	// Read segment first
	_, err := reader.ReadSegment()
	require.NoError(t, err)

	// Read trailer
	_, err = reader.ReadTrailer()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CRC64 checksum mismatch")
}

func TestStructuredMessageReader_Close(t *testing.T) {
	reader := NewStructuredMessageReader(bytes.NewReader([]byte{}))

	assert.False(t, reader.IsClosed())

	err := reader.Close()
	require.NoError(t, err)

	assert.True(t, reader.IsClosed())
}

func TestStructuredMessageReader_ReadAfterClose(t *testing.T) {
	reader := NewStructuredMessageReader(bytes.NewReader([]byte{}))
	reader.Close()

	_, err := reader.ReadSegment()
	assert.Equal(t, io.EOF, err)

	_, err = reader.ReadTrailer()
	assert.Equal(t, io.EOF, err)
}

func TestStructuredMessageWriter_GetTotalCRC64(t *testing.T) {
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)

	// Initially should be 0
	assert.Equal(t, uint64(0), writer.GetTotalCRC64())

	// Write segments
	messageLength := calculateMessageLength([]int{5, 5}, true)
	writer.WriteHeader(2, messageLength, true)
	writer.WriteSegment([]byte("hello"))
	writer.WriteSegment([]byte("world"))

	// Calculate expected CRC64
	allData := []byte("helloworld")
	expectedCRC64 := crc64.Checksum(allData, shared.CRC64Table)
	assert.Equal(t, expectedCRC64, writer.GetTotalCRC64())
}

func TestStructuredMessageReader_GetTotalCRC64(t *testing.T) {
	// Create test data
	data1 := []byte("hello")
	data2 := []byte("world")

	// Create structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)
	messageLength := calculateMessageLength([]int{len(data1), len(data2)}, true)
	writer.WriteHeader(2, messageLength, true)
	writer.WriteSegment(data1)
	writer.WriteSegment(data2)
	writer.WriteTrailer()

	// Test reading and CRC64 tracking
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Initially should be 0
	assert.Equal(t, uint64(0), reader.GetTotalCRC64())

	// Read first segment
	_, err := reader.ReadSegment()
	require.NoError(t, err)

	expectedCRC64AfterFirst := crc64.Checksum(data1, shared.CRC64Table)
	assert.Equal(t, expectedCRC64AfterFirst, reader.GetTotalCRC64())

	// Read second segment
	_, err = reader.ReadSegment()
	require.NoError(t, err)

	allData := append(data1, data2...)
	expectedCRC64AfterSecond := crc64.Checksum(allData, shared.CRC64Table)
	assert.Equal(t, expectedCRC64AfterSecond, reader.GetTotalCRC64())
}

func TestStructuredMessage_EndToEnd(t *testing.T) {
	// Test data
	testData := [][]byte{
		[]byte("segment1"),
		[]byte("segment2"),
		[]byte("segment3"),
	}

	// Write structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)

	// Calculate message length
	segmentSizes := make([]int, len(testData))
	for i, data := range testData {
		segmentSizes[i] = len(data)
	}
	messageLength := calculateMessageLength(segmentSizes, true)
	err := writer.WriteHeader(uint16(len(testData)), messageLength, true)
	require.NoError(t, err)

	for _, data := range testData {
		err = writer.WriteSegment(data)
		require.NoError(t, err)
	}

	err = writer.WriteTrailer()
	require.NoError(t, err)

	// Read structured message
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read header
	header, err := reader.ReadHeader()
	require.NoError(t, err)
	assert.Equal(t, uint16(len(testData)), header.NumSegments)

	// Read all segments
	allData, err := reader.ReadAllSegments()
	require.NoError(t, err)

	// Verify data
	expectedData := make([]byte, 0)
	for _, data := range testData {
		expectedData = append(expectedData, data...)
	}
	assert.Equal(t, expectedData, allData)

	// Read trailer
	trailerCRC64, err := reader.ReadTrailer()
	require.NoError(t, err)

	// Verify CRC64
	expectedCRC64 := crc64.Checksum(expectedData, shared.CRC64Table)
	assert.Equal(t, expectedCRC64, trailerCRC64)
	assert.Equal(t, expectedCRC64, reader.GetTotalCRC64())
	assert.Equal(t, expectedCRC64, writer.GetTotalCRC64())
}

func TestStructuredMessage_EmptySegments(t *testing.T) {
	// Test with empty segments
	testData := [][]byte{
		[]byte(""),
		[]byte("non-empty"),
		[]byte(""),
	}

	// Write structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)

	// Calculate message length
	segmentSizes := make([]int, len(testData))
	for i, data := range testData {
		segmentSizes[i] = len(data)
	}
	messageLength := calculateMessageLength(segmentSizes, true)
	err := writer.WriteHeader(uint16(len(testData)), messageLength, true)
	require.NoError(t, err)

	for _, data := range testData {
		err = writer.WriteSegment(data)
		require.NoError(t, err)
	}

	err = writer.WriteTrailer()
	require.NoError(t, err)

	// Read structured message
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read all segments
	allData, err := reader.ReadAllSegments()
	require.NoError(t, err)

	// Verify data
	expectedData := make([]byte, 0)
	for _, data := range testData {
		expectedData = append(expectedData, data...)
	}
	assert.Equal(t, expectedData, allData)

	// Read trailer
	_, err = reader.ReadTrailer()
	require.NoError(t, err)
}

func TestStructuredMessage_LargeData(t *testing.T) {
	// Test with large data
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	// Write structured message
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)

	messageLength := calculateMessageLength([]int{len(largeData)}, true)
	err := writer.WriteHeader(1, messageLength, true)
	require.NoError(t, err)

	err = writer.WriteSegment(largeData)
	require.NoError(t, err)

	err = writer.WriteTrailer()
	require.NoError(t, err)

	// Read structured message
	reader := NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read segment
	readData, err := reader.ReadSegment()
	require.NoError(t, err)
	assert.Equal(t, largeData, readData)

	// Read trailer
	_, err = reader.ReadTrailer()
	require.NoError(t, err)
}

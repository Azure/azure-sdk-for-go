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

func TestStructuredMessageWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name         string
		numSegments  uint32
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid single segment",
			numSegments: 1,
			expectError: false,
		},
		{
			name:        "Valid multiple segments",
			numSegments: 5,
			expectError: false,
		},
		{
			name:        "Valid maximum segments",
			numSegments: MaxSegments,
			expectError: false,
		},
		{
			name:         "Zero segments",
			numSegments:  0,
			expectError:  true,
			errorMessage: "invalid segment count",
		},
		{
			name:        "Too many segments",
			numSegments: MaxSegments,
			expectError: false, // MaxSegments should be valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewStructuredMessageWriter(&buf)

			err := writer.WriteHeader(tt.numSegments)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
			} else {
				require.NoError(t, err)
				require.Equal(t, HeaderSize, buf.Len())

				// Verify header content
				headerBytes := buf.Bytes()
				assert.Equal(t, uint8(MessageVersion), headerBytes[0])
				assert.Equal(t, uint8(FlagNone), headerBytes[1])
				assert.Equal(t, uint16(PropCRC64), binary.LittleEndian.Uint16(headerBytes[2:]))
				assert.Equal(t, tt.numSegments, binary.LittleEndian.Uint32(headerBytes[4:]))
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
			err := writer.WriteHeader(1)
			require.NoError(t, err)

			// Write segment
			err = writer.WriteSegment(tt.data)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMessage)
			} else {
				require.NoError(t, err)

				// Verify segment was written correctly
				expectedSize := HeaderSize + SegmentHeaderSize + len(tt.data)
				require.Equal(t, expectedSize, buf.Len())

				// Verify segment header
				segmentHeaderStart := HeaderSize
				segmentHeaderBytes := buf.Bytes()[segmentHeaderStart : segmentHeaderStart+SegmentHeaderSize]
				segmentNum := binary.LittleEndian.Uint32(segmentHeaderBytes[0:])
				dataLength := binary.LittleEndian.Uint32(segmentHeaderBytes[4:])
				segmentCRC64 := binary.LittleEndian.Uint64(segmentHeaderBytes[8:])

				assert.Equal(t, uint32(0), segmentNum)
				assert.Equal(t, uint32(len(tt.data)), dataLength)

				// Verify CRC64
				expectedCRC64 := crc64.Checksum(tt.data, shared.CRC64Table)
				assert.Equal(t, expectedCRC64, segmentCRC64)

				// Verify data
				dataStart := segmentHeaderStart + SegmentHeaderSize
				actualData := buf.Bytes()[dataStart : dataStart+len(tt.data)]
				assert.Equal(t, tt.data, actualData)
			}
		})
	}
}

func TestStructuredMessageWriter_WriteTrailer(t *testing.T) {
	var buf bytes.Buffer
	writer := NewStructuredMessageWriter(&buf)

	// Write header and segments
	err := writer.WriteHeader(2)
	require.NoError(t, err)

	data1 := []byte("hello")
	data2 := []byte("world")

	err = writer.WriteSegment(data1)
	require.NoError(t, err)

	err = writer.WriteSegment(data2)
	require.NoError(t, err)

	// Write trailer
	err = writer.WriteTrailer()
	require.NoError(t, err)

	// Verify trailer
	expectedSize := HeaderSize + (SegmentHeaderSize + len(data1)) + (SegmentHeaderSize + len(data2)) + TrailerSize
	require.Equal(t, expectedSize, buf.Len())

	// Verify trailer content
	trailerStart := expectedSize - TrailerSize
	trailerBytes := buf.Bytes()[trailerStart:]

	// Verify trailer header
	assert.Equal(t, uint8(MessageVersion), trailerBytes[0])
	assert.Equal(t, uint8(FlagNone), trailerBytes[1])
	assert.Equal(t, uint16(PropCRC64), binary.LittleEndian.Uint16(trailerBytes[2:]))

	expectedCRC64 := binary.LittleEndian.Uint64(trailerBytes[4:])

	// Calculate expected cumulative CRC64
	allData := append(data1, data2...)
	expectedCumulativeCRC64 := crc64.Checksum(allData, shared.CRC64Table)
	assert.Equal(t, expectedCumulativeCRC64, expectedCRC64)
	assert.Equal(t, expectedCumulativeCRC64, writer.GetTotalCRC64())
}

func TestStructuredMessageReader_ReadHeader(t *testing.T) {
	tests := []struct {
		name         string
		headerBytes  []byte
		expectError  bool
		errorMessage string
	}{
		{
			name:        "Valid header",
			headerBytes: []byte{MessageVersion, FlagNone, 0x01, 0x00, 1, 0, 0, 0}, // version=1, flags=0, properties=PropCRC64, numSegments=1
			expectError: false,
		},
		{
			name:         "Invalid version",
			headerBytes:  []byte{2, FlagNone, 0x01, 0x00, 1, 0, 0, 0}, // version=2
			expectError:  true,
			errorMessage: "invalid message version",
		},
		{
			name:         "Zero segments",
			headerBytes:  []byte{MessageVersion, FlagNone, 0x01, 0x00, 0, 0, 0, 0}, // numSegments=0
			expectError:  true,
			errorMessage: "invalid segment count",
		},
		{
			name:         "Missing CRC64 property",
			headerBytes:  []byte{MessageVersion, FlagNone, 0x00, 0x00, 1, 0, 0, 0}, // properties=0
			expectError:  true,
			errorMessage: "CRC64 property not set",
		},
		{
			name:        "Maximum segments",
			headerBytes: []byte{MessageVersion, FlagNone, 0x01, 0x00, 0xFF, 0xFF, 0xFF, 0xFF}, // numSegments=MaxUint32
			expectError: false,                                                                // MaxUint32 should be valid
		},
		{
			name:         "Incomplete header",
			headerBytes:  []byte{MessageVersion, FlagNone, 0x01, 0x00, 1, 0, 0}, // missing 1 byte
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
				assert.Equal(t, uint8(FlagNone), header.Flags)
				assert.Equal(t, uint16(PropCRC64), header.Properties)
				if tt.name == "Maximum segments" {
					assert.Equal(t, uint32(0xFFFFFFFF), header.NumSegments)
				} else {
					assert.Equal(t, uint32(1), header.NumSegments)
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
	writer.WriteHeader(2)
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
	writer.WriteHeader(1)
	writer.WriteSegment(data)
	writer.WriteTrailer()

	// Corrupt the CRC64 in the segment header
	bufBytes := buf.Bytes()
	segmentHeaderStart := HeaderSize
	crc64Offset := segmentHeaderStart + 8 // CRC64 is at offset 8 in segment header
	bufBytes[crc64Offset] ^= 0xFF         // Flip some bits

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
	writer.WriteHeader(3)
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
	writer.WriteHeader(2)
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
	writer.WriteHeader(1)
	writer.WriteSegment(data)
	writer.WriteTrailer()

	// Corrupt the trailer CRC64
	bufBytes := buf.Bytes()
	trailerStart := len(bufBytes) - TrailerSize
	crc64Offset := trailerStart + 4 // CRC64 is at offset 4 in trailer
	bufBytes[crc64Offset] ^= 0xFF   // Flip some bits in CRC64

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
	writer.WriteHeader(2)
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
	writer.WriteHeader(2)
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

	err := writer.WriteHeader(uint32(len(testData)))
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
	assert.Equal(t, uint32(len(testData)), header.NumSegments)

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

	err := writer.WriteHeader(uint32(len(testData)))
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

	err := writer.WriteHeader(1)
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

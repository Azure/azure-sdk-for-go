//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package structuredmessage

import (
	"bytes"
	"encoding/binary"
	"hash/crc64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeMessage(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "small data",
			data: []byte("hello world"),
		},
		{
			name: "empty string",
			data: []byte(""),
		},
		{
			name: "binary data",
			data: []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
		},
		{
			name: "large data",
			data: bytes.Repeat([]byte("test data 1234567890"), 100),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.data) == 0 {
				// Empty data should return error
				_, err := EncodeMessage(tc.data)
				require.Error(t, err)
				return
			}

			// Encode the message
			encoded, err := EncodeMessage(tc.data)
			require.NoError(t, err)
			require.NotEmpty(t, encoded)

			// Decode the message
			decoded, err := DecodeMessage(encoded)
			require.NoError(t, err)
			require.Equal(t, tc.data, decoded)
		})
	}
}

func TestEncodeMessageStructure(t *testing.T) {
	data := []byte("test data")
	encoded, err := EncodeMessage(data)
	require.NoError(t, err)

	reader := bytes.NewReader(encoded)

	// Check header
	var messageVersion uint8
	var messageLength uint64
	var messageFlags uint16
	var numSegments uint16

	require.NoError(t, binary.Read(reader, binary.LittleEndian, &messageVersion))
	require.Equal(t, MessageVersion, messageVersion)

	require.NoError(t, binary.Read(reader, binary.LittleEndian, &messageLength))
	require.Equal(t, uint64(len(encoded)), messageLength)

	require.NoError(t, binary.Read(reader, binary.LittleEndian, &messageFlags))
	require.Equal(t, FlagIncludeCRC64, messageFlags)

	require.NoError(t, binary.Read(reader, binary.LittleEndian, &numSegments))
	require.Equal(t, uint16(1), numSegments)

	// Check segment
	var segmentNum uint16
	var segmentDataLength uint64

	require.NoError(t, binary.Read(reader, binary.LittleEndian, &segmentNum))
	require.Equal(t, uint16(1), segmentNum)

	require.NoError(t, binary.Read(reader, binary.LittleEndian, &segmentDataLength))
	require.Equal(t, uint64(len(data)), segmentDataLength)

	// Read segment data
	segmentData := make([]byte, segmentDataLength)
	n, err := reader.Read(segmentData)
	require.NoError(t, err)
	require.Equal(t, len(data), n)
	require.Equal(t, data, segmentData)

	// Read segment CRC64
	segmentCRC64 := make([]byte, CRC64Size)
	n, err = reader.Read(segmentCRC64)
	require.NoError(t, err)
	require.Equal(t, CRC64Size, n)

	expectedCRC64 := crc64.Checksum(data, CRC64Table)
	actualCRC64 := binary.LittleEndian.Uint64(segmentCRC64)
	require.Equal(t, expectedCRC64, actualCRC64)

	// Read message CRC64 (trailer)
	messageCRC64 := make([]byte, CRC64Size)
	n, err = reader.Read(messageCRC64)
	require.NoError(t, err)
	require.Equal(t, CRC64Size, n)

	expectedMessageCRC64 := crc64.Checksum(data, CRC64Table)
	actualMessageCRC64 := binary.LittleEndian.Uint64(messageCRC64)
	require.Equal(t, expectedMessageCRC64, actualMessageCRC64)

	// Should be at end of message
	remaining := make([]byte, 1)
	n, err = reader.Read(remaining)
	require.Equal(t, 0, n)
}

func TestDecodeMessageValidation(t *testing.T) {
	data := []byte("test data")
	encoded, err := EncodeMessage(data)
	require.NoError(t, err)

	t.Run("corrupt segment CRC64", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Corrupt the segment CRC64 (last 16 bytes are segment CRC64 + message CRC64)
		corruptOffset := len(corrupted) - 16
		corrupted[corruptOffset] ^= 0xFF // Flip bits

		_, err := DecodeMessage(corrupted)
		require.Error(t, err)
		require.Contains(t, err.Error(), "segment CRC64 validation failed")
	})

	t.Run("corrupt message CRC64", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Corrupt the message CRC64 (last 8 bytes)
		corruptOffset := len(corrupted) - 8
		corrupted[corruptOffset] ^= 0xFF // Flip bits

		_, err := DecodeMessage(corrupted)
		require.Error(t, err)
		require.Contains(t, err.Error(), "message CRC64 validation failed")
	})

	t.Run("corrupt data", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Corrupt the data (somewhere in the middle)
		dataOffset := 13 + 2 + 8 + 2 // header + segment header + 2 bytes into data
		corrupted[dataOffset] ^= 0xFF // Flip bits

		_, err := DecodeMessage(corrupted)
		require.Error(t, err)
		require.Contains(t, err.Error(), "CRC64 validation failed")
	})

	t.Run("wrong message version", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Change message version
		corrupted[0] = 2

		_, err := DecodeMessage(corrupted)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported message version")
	})

	t.Run("wrong message length", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Change message length (bytes 1-8)
		binary.LittleEndian.PutUint64(corrupted[1:9], uint64(len(corrupted)+1))

		_, err := DecodeMessage(corrupted)
		require.Error(t, err)
		require.Contains(t, err.Error(), "message length mismatch")
	})

	t.Run("truncated data", func(t *testing.T) {
		truncated := encoded[:len(encoded)-5] // Remove last 5 bytes

		_, err := DecodeMessage(truncated)
		require.Error(t, err)
	})

	t.Run("too short", func(t *testing.T) {
		_, err := DecodeMessage([]byte{1, 2, 3})
		require.Error(t, err)
		require.Contains(t, err.Error(), "data too short")
	})
}

func TestCRC64Polynomial(t *testing.T) {
	// Ensure our polynomial matches the expected Azure Storage polynomial
	require.Equal(t, uint64(0x9A6C9329AC4BC9B5), CRC64Polynomial)
	
	// Test CRC64 computation with known data
	data := []byte("hello world")
	crc := crc64.Checksum(data, CRC64Table)
	
	// The exact value depends on the polynomial, just ensure it's consistent
	crc2 := crc64.Checksum(data, CRC64Table)
	require.Equal(t, crc, crc2)
}

func TestMessageSizeCalculation(t *testing.T) {
	data := []byte("test")
	encoded, err := EncodeMessage(data)
	require.NoError(t, err)

	// Expected size calculation:
	// Header: 1 + 8 + 2 + 2 = 13 bytes
	// Segment: 2 + 8 + 4 + 8 = 22 bytes (num + length + data + crc64)
	// Trailer: 8 bytes (message crc64)
	// Total: 13 + 22 + 8 = 43 bytes
	expectedSize := 43
	require.Equal(t, expectedSize, len(encoded))
}

func BenchmarkEncodeMessage(b *testing.B) {
	data := bytes.Repeat([]byte("benchmark data test"), 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeMessage(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMessage(b *testing.B) {
	data := bytes.Repeat([]byte("benchmark data test"), 100)
	encoded, err := EncodeMessage(data)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecodeMessage(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
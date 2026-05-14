//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"bytes"
	"hash/crc64"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmessage"
)

// TestStructuredMessageCRC64Consistency ensures our CRC64 implementation matches the expected polynomial
func TestStructuredMessageCRC64Consistency(t *testing.T) {
	// Verify we're using the correct Azure Storage polynomial
	require.Equal(t, uint64(0x9A6C9329AC4BC9B5), structuredmessage.CRC64Polynomial)
	require.Equal(t, structuredmessage.CRC64Polynomial, uint64(0x9A6C9329AC4BC9B5))
	
	// Verify our table matches the shared table
	require.Equal(t, shared.CRC64Table, structuredmessage.CRC64Table)
	
	// Test consistent CRC64 computation
	testData := []byte("Azure Storage CRC64 test data")
	
	sharedCRC64 := crc64.Checksum(testData, shared.CRC64Table)
	structuredCRC64 := crc64.Checksum(testData, structuredmessage.CRC64Table)
	
	require.Equal(t, sharedCRC64, structuredCRC64)
}

// TestStructuredMessageBinaryFormat validates the structured message binary format
func TestStructuredMessageBinaryFormat(t *testing.T) {
	testData := []byte("Binary format validation test")
	
	// Encode the message
	encoded, err := structuredmessage.EncodeMessage(testData)
	require.NoError(t, err)
	require.Greater(t, len(encoded), len(testData))
	
	// The encoded message should have a specific structure:
	// Header (13 bytes): version(1) + length(8) + flags(2) + segments(2)
	// Segment: segnum(2) + datalen(8) + data + crc64(8) 
	// Trailer: crc64(8)
	expectedMinSize := 13 + 2 + 8 + len(testData) + 8 + 8
	require.Equal(t, expectedMinSize, len(encoded))
	
	// Validate header structure
	require.Equal(t, uint8(1), encoded[0]) // version = 1
	require.Equal(t, uint16(1), uint16(encoded[11])|(uint16(encoded[12])<<8)) // num segments = 1
	
	// Decode and validate
	decoded, err := structuredmessage.DecodeMessage(encoded)
	require.NoError(t, err)
	require.Equal(t, testData, decoded)
}

// TestStructuredMessageHeaders validates the headers generated for structured messages
func TestStructuredMessageHeaders(t *testing.T) {
	testData := []byte("Header validation test data")
	reader := shared.NopCloser(bytes.NewReader(testData))
	
	validation := TransferValidationTypeStructuredMessage{}
	setter := &mockStructuredMessageSetter{}
	
	result, err := validation.ApplyStructured(reader, setter)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Validate headers match specification
	require.Equal(t, "XSM/1.0; properties=crc64", setter.bodyType)
	require.Equal(t, int64(len(testData)), setter.contentLength)
	
	// Cleanup
	err = result.Close()
	require.NoError(t, err)
}

// TestStructuredMessageDataIntegrity validates end-to-end data integrity
func TestStructuredMessageDataIntegrity(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"small data", []byte("hello")},
		{"medium data", []byte("This is a medium sized test data for CRC64 validation")},
		{"large data", bytes.Repeat([]byte("Large test data block "), 100)},
		{"binary data", []byte{0x00, 0x01, 0xFF, 0xFE, 0x42, 0x24, 0x80, 0x7F}},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Upload simulation (encode)
			uploadReader := shared.NopCloser(bytes.NewReader(tc.data))
			uploadValidation := TransferValidationTypeStructuredMessage{}
			uploadSetter := &mockStructuredMessageSetter{}
			
			encodedReader, err := uploadValidation.ApplyStructured(uploadReader, uploadSetter)
			require.NoError(t, err)
			
			encodedData, err := io.ReadAll(encodedReader)
			require.NoError(t, err)
			require.Greater(t, len(encodedData), len(tc.data))
			
			// Download simulation (decode)
			downloadReader := shared.NopCloser(bytes.NewReader(encodedData))
			downloadValidation := TransferValidationTypeStructuredMessageDownload{
				StructuredContentLength: int64(len(tc.data)),
			}
			
			decodedReader, err := downloadValidation.Apply(downloadReader, nil)
			require.NoError(t, err)
			
			decodedData, err := io.ReadAll(decodedReader)
			require.NoError(t, err)
			require.Equal(t, tc.data, decodedData)
			
			// Cleanup
			err = encodedReader.Close()
			require.NoError(t, err)
			err = decodedReader.Close()
			require.NoError(t, err)
		})
	}
}

// TestStructuredMessageCorruptionDetection validates that CRC64 validation detects corruption
func TestStructuredMessageCorruptionDetection(t *testing.T) {
	testData := []byte("Corruption detection test data")
	
	// Encode the message
	encoded, err := structuredmessage.EncodeMessage(testData)
	require.NoError(t, err)
	
	// Test various corruption scenarios
	t.Run("corrupt data", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Corrupt a byte in the middle of the data section
		dataOffset := 13 + 2 + 8 + 5 // header + segment header + 5 bytes into data
		corrupted[dataOffset] ^= 0xFF
		
		downloadReader := shared.NopCloser(bytes.NewReader(corrupted))
		downloadValidation := TransferValidationTypeStructuredMessageDownload{
			StructuredContentLength: int64(len(testData)),
		}
		
		_, err := downloadValidation.Apply(downloadReader, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "CRC64 validation failed")
	})
	
	t.Run("corrupt checksum", func(t *testing.T) {
		corrupted := make([]byte, len(encoded))
		copy(corrupted, encoded)
		
		// Corrupt the last byte (in the message CRC64)
		corrupted[len(corrupted)-1] ^= 0xFF
		
		downloadReader := shared.NopCloser(bytes.NewReader(corrupted))
		downloadValidation := TransferValidationTypeStructuredMessageDownload{
			StructuredContentLength: int64(len(testData)),
		}
		
		_, err := downloadValidation.Apply(downloadReader, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "validation failed")
	})
}
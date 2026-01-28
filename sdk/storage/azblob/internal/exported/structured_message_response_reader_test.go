//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmsg"
)

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func TestStructuredMessageResponseReader_NoHeader(t *testing.T) {
	// Test that reader returns body as-is when no structured body header
	data := []byte("test data")
	body := streaming.NopCloser(bytes.NewReader(data))
	structuredBodyType := ""

	reader, err := NewStructuredMessageResponseReader(body, &structuredBodyType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	readData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error reading: %v", err)
	}

	if !bytes.Equal(readData, data) {
		t.Errorf("expected %v, got %v", data, readData)
	}
}

func TestStructuredMessageResponseReader_NilHeader(t *testing.T) {
	// Test that reader returns body as-is when structured body header is nil
	data := []byte("test data")
	body := streaming.NopCloser(bytes.NewReader(data))

	reader, err := NewStructuredMessageResponseReader(body, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	readData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error reading: %v", err)
	}

	if !bytes.Equal(readData, data) {
		t.Errorf("expected %v, got %v", data, readData)
	}
}

func TestStructuredMessageResponseReader_ValidStructuredMessage(t *testing.T) {
	// Create a valid structured message
	originalData := make([]byte, 5*1024*1024) // 5MB
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// Write structured message
	var structuredBuffer bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&structuredBuffer)

	// Calculate segments
	segmentSize := int64(structuredmsg.MaxSegmentSize)
	numSegments := uint16((int64(len(originalData)) + segmentSize - 1) / segmentSize)

	// Calculate message length
	messageLength := uint64(structuredmsg.HeaderSize)
	offset := int64(0)
	for i := uint16(0); i < numSegments; i++ {
		remaining := int64(len(originalData)) - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}
		messageLength += uint64(structuredmsg.SegmentHeaderSize) + uint64(chunkSize) + uint64(structuredmsg.SegmentCRC64Size)
		offset += chunkSize
	}
	messageLength += uint64(structuredmsg.TrailerSize)

	// Write header (CRC64 enabled for transfer validation)
	if err := writer.WriteHeader(numSegments, messageLength, true); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}

	// Write segments
	offset = 0
	for i := uint16(0); i < numSegments; i++ {
		remaining := int64(len(originalData)) - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}

		segmentData := originalData[offset : offset+chunkSize]
		if err := writer.WriteSegment(segmentData); err != nil {
			t.Fatalf("failed to write segment: %v", err)
		}

		offset += chunkSize
	}

	// Write trailer
	if err := writer.WriteTrailer(); err != nil {
		t.Fatalf("failed to write trailer: %v", err)
	}

	structuredBodyType := "XSM/1.0;CRC64"
	body := streaming.NopCloser(bytes.NewReader(structuredBuffer.Bytes()))

	// Create reader
	reader, err := NewStructuredMessageResponseReader(body, &structuredBodyType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read all data
	readData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error reading: %v", err)
	}

	// Verify data matches
	if len(readData) != len(originalData) {
		t.Errorf("expected length %d, got %d", len(originalData), len(readData))
	}

	if !bytes.Equal(readData, originalData) {
		t.Error("read data does not match original data")
	}
}

func TestStructuredMessageResponseReader_ValidStructuredMessage_PropertiesFormat(t *testing.T) {
	// Test that "XSM/1.0;properties=crc64" format is also accepted
	originalData := make([]byte, 5*1024*1024) // 5MB
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// Write structured message
	var structuredBuffer bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&structuredBuffer)

	// Calculate segments
	segmentSize := int64(structuredmsg.MaxSegmentSize)
	numSegments := uint16((int64(len(originalData)) + segmentSize - 1) / segmentSize)

	// Calculate message length
	messageLength := uint64(structuredmsg.HeaderSize)
	offset := int64(0)
	for i := uint16(0); i < numSegments; i++ {
		remaining := int64(len(originalData)) - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}
		messageLength += uint64(structuredmsg.SegmentHeaderSize) + uint64(chunkSize) + uint64(structuredmsg.SegmentCRC64Size)
		offset += chunkSize
	}
	messageLength += uint64(structuredmsg.TrailerSize)

	// Write header (CRC64 enabled for transfer validation)
	if err := writer.WriteHeader(numSegments, messageLength, true); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}

	// Write segments
	offset = 0
	for i := uint16(0); i < numSegments; i++ {
		remaining := int64(len(originalData)) - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}

		segmentData := originalData[offset : offset+chunkSize]
		if err := writer.WriteSegment(segmentData); err != nil {
			t.Fatalf("failed to write segment: %v", err)
		}

		offset += chunkSize
	}

	// Write trailer
	if err := writer.WriteTrailer(); err != nil {
		t.Fatalf("failed to write trailer: %v", err)
	}

	// Test with "properties=crc64" format (case-insensitive)
	structuredBodyType := "XSM/1.0;properties=crc64"
	body := streaming.NopCloser(bytes.NewReader(structuredBuffer.Bytes()))

	// Create reader
	reader, err := NewStructuredMessageResponseReader(body, &structuredBodyType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read all data
	readData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error reading: %v", err)
	}

	// Verify data matches
	if len(readData) != len(originalData) {
		t.Errorf("expected length %d, got %d", len(originalData), len(readData))
	}

	if !bytes.Equal(readData, originalData) {
		t.Error("read data does not match original data")
	}
}

func TestStructuredMessageResponseReader_CRC64Mismatch(t *testing.T) {
	// Create a structured message with corrupted data
	originalData := make([]byte, 1024*1024) // 1MB
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// Write structured message
	var structuredBuffer bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&structuredBuffer)

	// Calculate message length
	messageLength := uint64(structuredmsg.HeaderSize) +
		uint64(structuredmsg.SegmentHeaderSize) + uint64(len(originalData)) + uint64(structuredmsg.SegmentCRC64Size) +
		uint64(structuredmsg.TrailerSize)

	// Write header (CRC64 enabled for transfer validation)
	if err := writer.WriteHeader(1, messageLength, true); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}

	// Write segment
	if err := writer.WriteSegment(originalData); err != nil {
		t.Fatalf("failed to write segment: %v", err)
	}

	// Corrupt the trailer CRC64 by writing an invalid one
	// First, get the correct trailer from the writer
	correctCRC64 := writer.GetTotalCRC64()

	// Now write a corrupted trailer (trailer is only 8 bytes: CRC64)
	trailerBytes := make([]byte, structuredmsg.TrailerSize)
	// Set invalid CRC64 (different from correct one)
	invalidCRC64 := correctCRC64 + 1
	binary.LittleEndian.PutUint64(trailerBytes[0:8], invalidCRC64)
	structuredBuffer.Write(trailerBytes)

	structuredBodyType := "XSM/1.0;CRC64"
	body := streaming.NopCloser(bytes.NewReader(structuredBuffer.Bytes()))

	// Create reader
	reader, err := NewStructuredMessageResponseReader(body, &structuredBodyType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to read - should fail with CRC64 mismatch
	_, err = io.ReadAll(reader)
	if err == nil {
		t.Fatal("expected error for CRC64 mismatch, got nil")
	}

	// Verify error message contains CRC64 mismatch
	if err.Error() == "" || !contains(err.Error(), "CRC64") {
		t.Errorf("expected error about CRC64 mismatch, got: %v", err)
	}

	// Close should also return the error
	closeErr := reader.Close()
	if closeErr == nil {
		t.Error("expected error on close, got nil")
	}
}

func TestStructuredMessageResponseReader_InvalidFormat(t *testing.T) {
	// Test that reader returns body as-is for unrecognized format
	data := []byte("test data")

	invalidFormats := []string{
		"UNKNOWN/1.0",                // Different prefix
		"XSM/1.0",                    // Missing semicolon and property
		"XSM/1.0;",                   // Missing property
		"XSM/1.0;UNKNOWN",            // Unknown property
		"XSM/1.0;CRC64EXTRA",         // Malformed (no space/semicolon after CRC64)
		"XSM/1.0;CRC64;EXTRA",        // Extra content after CRC64
		"XSM/1.0;CRC65",              // Typo in property name
		"XSM/1.0;properties=unknown", // Unknown property value
		"XSM/1.0;properties=CRC65",   // Typo in property value
	}

	for _, invalidFormat := range invalidFormats {
		t.Run(invalidFormat, func(t *testing.T) {
			body := streaming.NopCloser(bytes.NewReader(data))
			structuredBodyType := invalidFormat

			reader, err := NewStructuredMessageResponseReader(body, &structuredBodyType)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			readData, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("unexpected error reading: %v", err)
			}

			if !bytes.Equal(readData, data) {
				t.Errorf("expected %v, got %v", data, readData)
			}
		})
	}
}

func TestStructuredMessageResponseReader_PartialRead(t *testing.T) {
	// Test reading partial data and then closing
	originalData := make([]byte, 2*1024*1024) // 2MB
	for i := range originalData {
		originalData[i] = byte(i % 256)
	}

	// Write structured message
	var structuredBuffer bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&structuredBuffer)

	segmentSize := int64(structuredmsg.MaxSegmentSize)
	numSegments := uint16((int64(len(originalData)) + segmentSize - 1) / segmentSize)

	// Calculate message length
	messageLength := uint64(structuredmsg.HeaderSize)
	offset := int64(0)
	for i := uint16(0); i < numSegments; i++ {
		remaining := int64(len(originalData)) - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}
		messageLength += uint64(structuredmsg.SegmentHeaderSize) + uint64(chunkSize) + uint64(structuredmsg.SegmentCRC64Size)
		offset += chunkSize
	}
	messageLength += uint64(structuredmsg.TrailerSize)

	if err := writer.WriteHeader(numSegments, messageLength); err != nil {
		t.Fatalf("failed to write header: %v", err)
	}

	offset = 0
	for i := uint16(0); i < numSegments; i++ {
		remaining := int64(len(originalData)) - offset
		chunkSize := segmentSize
		if remaining < chunkSize {
			chunkSize = remaining
		}

		segmentData := originalData[offset : offset+chunkSize]
		if err := writer.WriteSegment(segmentData); err != nil {
			t.Fatalf("failed to write segment: %v", err)
		}

		offset += chunkSize
	}

	if err := writer.WriteTrailer(); err != nil {
		t.Fatalf("failed to write trailer: %v", err)
	}

	structuredBodyType := "XSM/1.0;CRC64"
	body := streaming.NopCloser(bytes.NewReader(structuredBuffer.Bytes()))

	reader, err := NewStructuredMessageResponseReader(body, &structuredBodyType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read only part of the data
	partialData := make([]byte, 1024)
	n, err := reader.Read(partialData)
	if err != nil && err != io.EOF {
		t.Fatalf("unexpected error reading: %v", err)
	}
	if n == 0 {
		t.Error("expected to read some data")
	}

	// Close should validate remaining segments and trailer
	err = reader.Close()
	if err != nil {
		t.Errorf("unexpected error on close: %v", err)
	}
}

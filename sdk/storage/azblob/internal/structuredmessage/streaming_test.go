//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package structuredmessage

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoderReader(t *testing.T) {
	testData := []byte("hello world test data")
	reader := bytes.NewReader(testData)

	encoder, err := NewEncoderReader(reader)
	require.NoError(t, err)
	require.NotNil(t, encoder)

	// Check lengths
	require.Equal(t, int64(len(testData)), encoder.ContentLength())
	require.Greater(t, encoder.EncodedLength(), encoder.ContentLength())

	// Check headers
	bodyType, contentLength := encoder.GetStructuredHeaders()
	require.Equal(t, "XSM/1.0; properties=crc64", bodyType)
	require.Equal(t, int64(len(testData)), contentLength)

	// Read encoded data
	encodedData, err := io.ReadAll(encoder)
	require.NoError(t, err)
	require.Equal(t, int(encoder.EncodedLength()), len(encodedData))

	// Verify it can be decoded back to original
	decodedData, err := DecodeMessage(encodedData)
	require.NoError(t, err)
	require.Equal(t, testData, decodedData)
}

func TestDecoderReader(t *testing.T) {
	testData := []byte("hello world test data")

	// First encode the data
	encodedData, err := EncodeMessage(testData)
	require.NoError(t, err)

	// Create decoder reader
	reader := bytes.NewReader(encodedData)
	decoder, err := NewDecoderReader(reader, int64(len(testData)))
	require.NoError(t, err)
	require.NotNil(t, decoder)

	// Check properties
	require.Equal(t, int64(len(testData)), decoder.ContentLength())
	require.Equal(t, int64(len(encodedData)), decoder.OriginalLength())
	require.True(t, decoder.IsValidated())

	// Read decoded data
	decodedData, err := io.ReadAll(decoder)
	require.NoError(t, err)
	require.Equal(t, testData, decodedData)
}

func TestDecoderReaderValidation(t *testing.T) {
	testData := []byte("hello world test data")
	encodedData, err := EncodeMessage(testData)
	require.NoError(t, err)

	t.Run("wrong content length", func(t *testing.T) {
		reader := bytes.NewReader(encodedData)
		_, err := NewDecoderReader(reader, int64(len(testData)+1))
		require.Error(t, err)
	})

	t.Run("corrupted data", func(t *testing.T) {
		corrupted := make([]byte, len(encodedData))
		copy(corrupted, encodedData)
		corrupted[len(corrupted)-1] ^= 0xFF // Corrupt last byte

		reader := bytes.NewReader(corrupted)
		_, err := NewDecoderReader(reader, int64(len(testData)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "validation failed")
	})
}

func TestStreamingEncoder(t *testing.T) {
	testData := []byte("hello world streaming test")
	
	encoder := NewStreamingEncoder(int64(len(testData)))
	require.NotNil(t, encoder)

	// Write data
	err := encoder.WriteData(testData)
	require.NoError(t, err)

	// Check headers
	bodyType, contentLength := encoder.GetStructuredHeaders()
	require.Equal(t, "XSM/1.0; properties=crc64", bodyType)
	require.Equal(t, int64(len(testData)), contentLength)

	// Get encoded data
	encodedData, err := encoder.GetEncodedData()
	require.NoError(t, err)
	require.Greater(t, len(encodedData), len(testData))

	// Verify it can be decoded
	decodedData, err := DecodeMessage(encodedData)
	require.NoError(t, err)
	require.Equal(t, testData, decodedData)
}

func TestStreamingEncoderMultipleWrites(t *testing.T) {
	chunks := [][]byte{
		[]byte("hello "),
		[]byte("world "),
		[]byte("streaming "),
		[]byte("test"),
	}
	
	expectedData := []byte("hello world streaming test")
	totalLength := int64(len(expectedData))
	
	encoder := NewStreamingEncoder(totalLength)
	
	// Write chunks
	for _, chunk := range chunks {
		err := encoder.WriteData(chunk)
		require.NoError(t, err)
	}

	// Get encoded data
	encodedData, err := encoder.GetEncodedData()
	require.NoError(t, err)

	// Verify it can be decoded
	decodedData, err := DecodeMessage(encodedData)
	require.NoError(t, err)
	require.Equal(t, expectedData, decodedData)
}

func TestEncoderDecoderRoundTrip(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"small text", "hello"},
		{"medium text", "hello world test data with more content"},
		{"large text", strings.Repeat("test data 1234567890 ", 100)},
		{"binary data", string([]byte{0x00, 0x01, 0xFF, 0xFE, 0x42, 0x24})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testData := []byte(tc.data)

			// Encode using EncoderReader
			reader := bytes.NewReader(testData)
			encoder, err := NewEncoderReader(reader)
			require.NoError(t, err)

			encodedData, err := io.ReadAll(encoder)
			require.NoError(t, err)

			// Decode using DecoderReader
			encodedReader := bytes.NewReader(encodedData)
			decoder, err := NewDecoderReader(encodedReader, int64(len(testData)))
			require.NoError(t, err)

			decodedData, err := io.ReadAll(decoder)
			require.NoError(t, err)

			// Verify round trip
			require.Equal(t, testData, decodedData)
			require.True(t, decoder.IsValidated())
		})
	}
}

func TestEncoderReaderWithSeeker(t *testing.T) {
	testData := []byte("seekable test data")
	reader := bytes.NewReader(testData)
	
	// Seek to middle to test that encoder resets position
	_, err := reader.Seek(5, io.SeekStart)
	require.NoError(t, err)

	// Create encoder (should reset reader position)
	encoder, err := NewEncoderReader(reader)
	require.NoError(t, err)

	// Verify full data is encoded
	require.Equal(t, int64(len(testData)), encoder.ContentLength())

	// Read and decode to verify
	encodedData, err := io.ReadAll(encoder)
	require.NoError(t, err)

	decodedData, err := DecodeMessage(encodedData)
	require.NoError(t, err)
	require.Equal(t, testData, decodedData)
}

func BenchmarkEncoderReader(b *testing.B) {
	testData := bytes.Repeat([]byte("benchmark test data"), 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(testData)
		encoder, err := NewEncoderReader(reader)
		if err != nil {
			b.Fatal(err)
		}
		_, err = io.ReadAll(encoder)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecoderReader(b *testing.B) {
	testData := bytes.Repeat([]byte("benchmark test data"), 100)
	encodedData, err := EncodeMessage(testData)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(encodedData)
		decoder, err := NewDecoderReader(reader, int64(len(testData)))
		if err != nil {
			b.Fatal(err)
		}
		_, err = io.ReadAll(decoder)
		if err != nil {
			b.Fatal(err)
		}
	}
}
//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmessage"
)

// Mock structured message setter for testing
type mockStructuredMessageSetter struct {
	bodyType      string
	contentLength int64
}

func (m *mockStructuredMessageSetter) SetStructuredBodyType(bodyType string) {
	m.bodyType = bodyType
}

func (m *mockStructuredMessageSetter) SetStructuredContentLength(length int64) {
	m.contentLength = length
}

func TestTransferValidationTypeStructuredMessage(t *testing.T) {
	testData := []byte("hello world structured message test")
	reader := shared.NopCloser(bytes.NewReader(testData))
	
	// Create the transfer validation type
	validation := TransferValidationTypeStructuredMessage{}
	
	// Create mock setter
	mockSetter := &mockStructuredMessageSetter{}
	
	// Apply structured message validation
	result, err := validation.ApplyStructured(reader, mockSetter)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Check that headers were set correctly
	require.Equal(t, "XSM/1.0; properties=crc64", mockSetter.bodyType)
	require.Equal(t, int64(len(testData)), mockSetter.contentLength)
	
	// Read the encoded result
	encodedData, err := io.ReadAll(result)
	require.NoError(t, err)
	require.Greater(t, len(encodedData), len(testData)) // Should be larger due to structured message format
	
	// Close the result
	err = result.Close()
	require.NoError(t, err)
}

func TestTransferValidationTypeStructuredMessageDownload(t *testing.T) {
	testData := []byte("hello world download test")
	
	// First encode the data to simulate a structured message response
	encoded, err := structuredmessage.EncodeMessage(testData)
	require.NoError(t, err)
	
	reader := shared.NopCloser(bytes.NewReader(encoded))
	
	// Create the download validation type
	validation := TransferValidationTypeStructuredMessageDownload{
		StructuredContentLength: int64(len(testData)),
	}
	
	// Apply validation (this should decode the structured message)
	result, err := validation.Apply(reader, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Read the decoded result
	decodedData, err := io.ReadAll(result)
	require.NoError(t, err)
	require.Equal(t, testData, decodedData)
	
	// Close the result
	err = result.Close()
	require.NoError(t, err)
}

func TestStructuredMessageReadSeekCloser(t *testing.T) {
	testData := []byte("seekable structured message test")
	reader := shared.NopCloser(bytes.NewReader(testData))
	
	validation := TransferValidationTypeStructuredMessage{}
	mockSetter := &mockStructuredMessageSetter{}
	
	result, err := validation.ApplyStructured(reader, mockSetter)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Test seeking
	pos, err := result.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
	
	// Read some data
	buf := make([]byte, 10)
	n, err := result.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 10, n)
	
	// Seek back to start
	pos, err = result.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
	
	// Read all data
	allData, err := io.ReadAll(result)
	require.NoError(t, err)
	require.Greater(t, len(allData), len(testData))
	
	err = result.Close()
	require.NoError(t, err)
}

func TestStructuredMessageDecodeReadSeekCloser(t *testing.T) {
	testData := []byte("seekable decode test")
	encoded, err := structuredmessage.EncodeMessage(testData)
	require.NoError(t, err)
	
	reader := shared.NopCloser(bytes.NewReader(encoded))
	
	validation := TransferValidationTypeStructuredMessageDownload{
		StructuredContentLength: int64(len(testData)),
	}
	
	result, err := validation.Apply(reader, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Test seeking
	pos, err := result.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
	
	// Read some data
	buf := make([]byte, 5)
	n, err := result.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, testData[:5], buf)
	
	// Seek to end
	pos, err = result.Seek(0, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(len(testData)), pos)
	
	// Seek back to start
	pos, err = result.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
	
	// Read all data
	allData, err := io.ReadAll(result)
	require.NoError(t, err)
	require.Equal(t, testData, allData)
	
	err = result.Close()
	require.NoError(t, err)
}

// Helper function for tests - we need to expose EncodeMessage for testing
func EncodeMessage(data []byte) ([]byte, error) {
	return structuredmessage.EncodeMessage(data)
}
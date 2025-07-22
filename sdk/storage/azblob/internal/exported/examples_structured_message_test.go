//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmessage"
)

// ExampleStructuredMessageSetter shows how to implement the StructuredMessageSetter interface
type ExampleStructuredMessageSetter struct {
	BodyType      string
	ContentLength int64
}

func (e *ExampleStructuredMessageSetter) SetStructuredBodyType(bodyType string) {
	e.BodyType = bodyType
}

func (e *ExampleStructuredMessageSetter) SetStructuredContentLength(length int64) {
	e.ContentLength = length
}

// ExampleTransferValidationTypeStructuredMessage demonstrates upload with structured message CRC64 validation
func Example_transferValidationTypeStructuredMessage() {
	// Sample data to upload
	data := []byte("Hello, Azure Storage with CRC64 validation!")
	reader := shared.NopCloser(bytes.NewReader(data))

	// Create structured message transfer validation
	validation := exported.TransferValidationTypeStructuredMessage{}
	setter := &ExampleStructuredMessageSetter{}

	// Apply structured message encoding
	encodedReader, err := validation.ApplyStructured(reader, setter)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// The setter now contains the headers needed for the HTTP request
	fmt.Printf("Structured Body Type: %s\n", setter.BodyType)
	fmt.Printf("Structured Content Length: %d\n", setter.ContentLength)

	// Read the encoded data (this would be sent in the HTTP request body)
	encodedData, err := io.ReadAll(encodedReader)
	if err != nil {
		fmt.Printf("Error reading encoded data: %v\n", err)
		return
	}

	fmt.Printf("Original data size: %d bytes\n", len(data))
	fmt.Printf("Encoded data size: %d bytes\n", len(encodedData))
	fmt.Printf("Encoded data includes CRC64 validation and structured message format\n")

	// Verify that the encoded data can be decoded back to original
	decoded, err := structuredmessage.DecodeMessage(encodedData)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return
	}

	fmt.Printf("Successfully validated: %t\n", bytes.Equal(data, decoded))

	// Output:
	// Structured Body Type: XSM/1.0; properties=crc64
	// Structured Content Length: 43
	// Original data size: 43 bytes
	// Encoded data size: 82 bytes
	// Encoded data includes CRC64 validation and structured message format
	// Successfully validated: true
}

// ExampleTransferValidationTypeStructuredMessageDownload demonstrates download with structured message CRC64 validation
func Example_transferValidationTypeStructuredMessageDownload() {
	// Original data (what was stored)
	originalData := []byte("Hello, Azure Storage download with CRC64 validation!")

	// Simulate server response with structured message (this would come from the server)
	encodedResponse, err := structuredmessage.EncodeMessage(originalData)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}

	// Simulate reading the structured response
	responseReader := shared.NopCloser(bytes.NewReader(encodedResponse))

	// Create download validation with expected content length
	validation := exported.TransferValidationTypeStructuredMessageDownload{
		StructuredContentLength: int64(len(originalData)),
	}

	// Apply structured message decoding and validation
	decodedReader, err := validation.Apply(responseReader, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Read the decoded and validated data
	validatedData, err := io.ReadAll(decodedReader)
	if err != nil {
		fmt.Printf("Error reading validated data: %v\n", err)
		return
	}

	fmt.Printf("Original data size: %d bytes\n", len(originalData))
	fmt.Printf("Encoded response size: %d bytes\n", len(encodedResponse))
	fmt.Printf("Validated data size: %d bytes\n", len(validatedData))
	fmt.Printf("Data integrity verified: %t\n", bytes.Equal(originalData, validatedData))
	fmt.Printf("CRC64 validation passed automatically during decoding\n")

	// Output:
	// Original data size: 52 bytes
	// Encoded response size: 91 bytes
	// Validated data size: 52 bytes
	// Data integrity verified: true
	// CRC64 validation passed automatically during decoding
}

// ExampleStructuredMessageRoundTrip demonstrates a complete round trip with validation
func Example_structuredMessageRoundTrip() {
	originalData := []byte("Round trip test with structured message CRC64 validation!")

	fmt.Printf("=== Upload (Encode) ===\n")
	// Upload: encode data with structured message
	uploadReader := shared.NopCloser(bytes.NewReader(originalData))
	uploadValidation := exported.TransferValidationTypeStructuredMessage{}
	uploadSetter := &ExampleStructuredMessageSetter{}

	encodedReader, err := uploadValidation.ApplyStructured(uploadReader, uploadSetter)
	if err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return
	}

	encodedData, err := io.ReadAll(encodedReader)
	if err != nil {
		fmt.Printf("Error reading encoded data: %v\n", err)
		return
	}

	fmt.Printf("Upload headers: %s, Content-Length: %d\n", uploadSetter.BodyType, uploadSetter.ContentLength)
	fmt.Printf("Encoded size: %d bytes (includes CRC64 and metadata)\n", len(encodedData))

	fmt.Printf("\n=== Download (Decode) ===\n")
	// Download: decode data and validate structured message
	downloadReader := shared.NopCloser(bytes.NewReader(encodedData))
	downloadValidation := exported.TransferValidationTypeStructuredMessageDownload{
		StructuredContentLength: int64(len(originalData)),
	}

	decodedReader, err := downloadValidation.Apply(downloadReader, nil)
	if err != nil {
		fmt.Printf("Download error: %v\n", err)
		return
	}

	validatedData, err := io.ReadAll(decodedReader)
	if err != nil {
		fmt.Printf("Error reading validated data: %v\n", err)
		return
	}

	fmt.Printf("Decoded size: %d bytes\n", len(validatedData))
	fmt.Printf("Round trip successful: %t\n", bytes.Equal(originalData, validatedData))
	fmt.Printf("CRC64 validation ensures data integrity throughout the process\n")

	// Output:
	// === Upload (Encode) ===
	// Upload headers: XSM/1.0; properties=crc64, Content-Length: 57
	// Encoded size: 96 bytes (includes CRC64 and metadata)
	//
	// === Download (Decode) ===
	// Decoded size: 57 bytes
	// Round trip successful: true
	// CRC64 validation ensures data integrity throughout the process
}
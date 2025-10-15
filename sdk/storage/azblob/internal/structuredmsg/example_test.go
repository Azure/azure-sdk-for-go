//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package structuredmsg_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmsg"
)

// ExampleStructuredMessageWriter demonstrates how to use StructuredMessageWriter
func ExampleStructuredMessageWriter() {
	// Create a buffer to write the structured message to
	var buf bytes.Buffer

	// Create a new structured message writer
	writer := structuredmsg.NewStructuredMessageWriter(&buf)

	// Prepare some test data
	data1 := []byte("Hello, ")
	data2 := []byte("World!")
	data3 := []byte(" This is a test.")

	// Write the message header (specify number of segments)
	err := writer.WriteHeader(3)
	if err != nil {
		panic(err)
	}

	// Write each segment
	err = writer.WriteSegment(data1)
	if err != nil {
		panic(err)
	}

	err = writer.WriteSegment(data2)
	if err != nil {
		panic(err)
	}

	err = writer.WriteSegment(data3)
	if err != nil {
		panic(err)
	}

	// Write the trailer with cumulative CRC64
	err = writer.WriteTrailer()
	if err != nil {
		panic(err)
	}

	// Get the total CRC64 for verification
	totalCRC64 := writer.GetTotalCRC64()
	fmt.Printf("Total CRC64: %x\n", totalCRC64)
	fmt.Printf("Message size: %d bytes\n", buf.Len())
}

// ExampleStructuredMessageReader demonstrates how to use StructuredMessageReader
func ExampleStructuredMessageReader() {
	// Create some test data
	data1 := []byte("Hello, ")
	data2 := []byte("World!")
	data3 := []byte(" This is a test.")

	// First, create a structured message
	var buf bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&buf)
	writer.WriteHeader(3)
	writer.WriteSegment(data1)
	writer.WriteSegment(data2)
	writer.WriteSegment(data3)
	writer.WriteTrailer()

	// Now read it back
	reader := structuredmsg.NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read the header
	header, err := reader.ReadHeader()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version: %d, Segments: %d\n", header.Version, header.NumSegments)

	// Read all segments
	allData, err := reader.ReadAllSegments()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Combined data: %s\n", string(allData))

	// Read the trailer and verify CRC64
	trailerCRC64, err := reader.ReadTrailer()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Trailer CRC64: %x\n", trailerCRC64)
	fmt.Printf("Reader CRC64: %x\n", reader.GetTotalCRC64())
}

// ExampleStructuredMessageReader_ReadSegmentBySegment demonstrates reading segments one by one
func ExampleStructuredMessageReader_readSegmentBySegment() {
	// Create some test data
	data1 := []byte("Segment 1")
	data2 := []byte("Segment 2")
	data3 := []byte("Segment 3")

	// First, create a structured message
	var buf bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&buf)
	writer.WriteHeader(3)
	writer.WriteSegment(data1)
	writer.WriteSegment(data2)
	writer.WriteSegment(data3)
	writer.WriteTrailer()

	// Now read it back segment by segment
	reader := structuredmsg.NewStructuredMessageReader(bytes.NewReader(buf.Bytes()))

	// Read header
	header, err := reader.ReadHeader()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of segments: %d\n", header.NumSegments)

	// Read segments one by one
	for i := 0; i < int(header.NumSegments); i++ {
		segment, err := reader.ReadSegment()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		fmt.Printf("Segment %d: %s\n", i, string(segment))
	}

	// Read trailer
	trailerCRC64, err := reader.ReadTrailer()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Trailer CRC64: %x\n", trailerCRC64)
}

// ExampleStructuredMessageReader_errorHandling demonstrates error handling
func ExampleStructuredMessageReader_errorHandling() {
	// Create a reader with invalid data
	invalidData := []byte{0x02, 0x00, 0x01, 0x00, 0x00, 0x00} // Invalid version
	reader := structuredmsg.NewStructuredMessageReader(bytes.NewReader(invalidData))

	// Try to read header - this should fail
	_, err := reader.ReadHeader()
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Create a reader with corrupted CRC64
	data := []byte("test data")
	var buf bytes.Buffer
	writer := structuredmsg.NewStructuredMessageWriter(&buf)
	writer.WriteHeader(1)
	writer.WriteSegment(data)
	writer.WriteTrailer()

	// Corrupt the CRC64 in the segment header
	bufBytes := buf.Bytes()
	segmentHeaderStart := 6               // HeaderSize
	crc64Offset := segmentHeaderStart + 8 // CRC64 is at offset 8 in segment header
	bufBytes[crc64Offset] ^= 0xFF         // Flip some bits

	// Try to read the corrupted data
	reader2 := structuredmsg.NewStructuredMessageReader(bytes.NewReader(bufBytes))
	_, err = reader2.ReadSegment()
	if err != nil {
		fmt.Printf("Expected CRC64 error: %v\n", err)
	}
}

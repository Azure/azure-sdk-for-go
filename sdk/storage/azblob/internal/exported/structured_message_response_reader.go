//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/structuredmsg"
)

// StructuredMessageResponseReader wraps a response body and validates structured messages
// if the x-ms-structured-body header is present. It strips the structured message framing
// and returns only the validated data to the caller.
type StructuredMessageResponseReader struct {
	reader              *structuredmsg.StructuredMessageReader
	underlyingBody      io.ReadCloser
	currentSegmentData  []byte
	currentSegmentIndex int
	headerRead          bool
	closed              bool
}

// NewStructuredMessageResponseReader creates a new reader that validates structured messages
// if the structuredBodyType header indicates a structured message format.
// If structuredBodyType is nil or empty, it returns the original body unchanged.
func NewStructuredMessageResponseReader(body io.ReadCloser, structuredBodyType *string) (io.ReadCloser, error) {
	if structuredBodyType == nil || *structuredBodyType == "" {
		// No structured message, return body as-is
		return body, nil
	}

	// Check if it's a structured message format we support
	// Format: "XSM/1.0;CRC64" or "XSM/1.0;properties=crc64"
	bodyType := strings.ToUpper(*structuredBodyType)
	if !strings.HasPrefix(bodyType, "XSM/1.0") {
		return body, nil
	}

	// Create structured message reader
	reader := structuredmsg.NewStructuredMessageReader(body)

	return &StructuredMessageResponseReader{
		reader:             reader,
		underlyingBody:     body,
		currentSegmentData: nil,
		headerRead:         false,
		closed:             false,
	}, nil
}

// Read implements io.Reader. It reads data from structured message segments,
// validates CRC64, and returns only the data (stripping framing).
func (r *StructuredMessageResponseReader) Read(p []byte) (int, error) {
	if r.closed {
		return 0, io.EOF
	}

	// Read header on first read
	if !r.headerRead {
		_, err := r.reader.ReadHeader()
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			// Map structured message errors to blob errors
			if errors.Is(err, structuredmsg.ErrInvalidVersion) {
				return 0, fmt.Errorf("invalid structured message version: %w", err)
			}
			if errors.Is(err, structuredmsg.ErrUnexpectedEOF) {
				return 0, fmt.Errorf("unexpected end of structured message: %w", err)
			}
			return 0, err
		}
		r.headerRead = true
	}

	// Read data from current segment, or get next segment
	totalRead := 0
	for totalRead < len(p) {
		// If we've exhausted current segment, read next one
		if r.currentSegmentData == nil || r.currentSegmentIndex >= len(r.currentSegmentData) {
			segmentData, err := r.reader.ReadSegment()
			if err != nil {
				if err == io.EOF {
					// No more segments, validate trailer
					_, trailerErr := r.reader.ReadTrailer()
					if trailerErr != nil {
						if trailerErr == structuredmsg.ErrCRC64Mismatch {
							return 0, fmt.Errorf("CRC64 checksum mismatch in structured message trailer: %w", trailerErr)
						}
						if trailerErr == structuredmsg.ErrUnexpectedEOF {
							return 0, fmt.Errorf("unexpected end of structured message trailer: %w", trailerErr)
						}
						return totalRead, trailerErr
					}
					// Successfully validated all segments and trailer
					r.closed = true
					if totalRead == 0 {
						return 0, io.EOF
					}
					return totalRead, nil
				}
				if err == structuredmsg.ErrCRC64Mismatch {
					return 0, fmt.Errorf("CRC64 checksum mismatch in structured message segment: %w", err)
				}
				if err == structuredmsg.ErrInvalidSegmentNumber {
					return 0, fmt.Errorf("invalid segment number in structured message: %w", err)
				}
				if err == structuredmsg.ErrUnexpectedEOF {
					return 0, fmt.Errorf("unexpected end of structured message segment: %w", err)
				}
				return totalRead, err
			}
			r.currentSegmentData = segmentData
			r.currentSegmentIndex = 0
		}

		// Copy data from current segment to output buffer
		available := len(r.currentSegmentData) - r.currentSegmentIndex
		needed := len(p) - totalRead
		copySize := available
		if copySize > needed {
			copySize = needed
		}

		copy(p[totalRead:totalRead+copySize], r.currentSegmentData[r.currentSegmentIndex:r.currentSegmentIndex+copySize])
		totalRead += copySize
		r.currentSegmentIndex += copySize

		// If we've consumed the entire segment, mark it as nil for next iteration
		if r.currentSegmentIndex >= len(r.currentSegmentData) {
			r.currentSegmentData = nil
			r.currentSegmentIndex = 0
		}
	}

	return totalRead, nil
}

// Close implements io.Closer. It closes the underlying body and marks the reader as closed.
// If not all segments have been read, it attempts to read and validate the remaining data.
func (r *StructuredMessageResponseReader) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true

	// If we haven't finished reading, try to read remaining segments and validate trailer
	// This ensures we catch any CRC64 mismatches even if the user doesn't read all data
	if r.headerRead {
		// Read any remaining segments
		for {
			_, err := r.reader.ReadSegment()
			if err == io.EOF {
				// All segments read, validate trailer
				_, trailerErr := r.reader.ReadTrailer()
				if trailerErr != nil {
					// Close underlying body first, then return error
					closeErr := r.underlyingBody.Close()
					if closeErr != nil {
						return closeErr
					}
					if trailerErr == structuredmsg.ErrCRC64Mismatch {
						return fmt.Errorf("CRC64 checksum mismatch in structured message trailer: %w", trailerErr)
					}
					return trailerErr
				}
				break
			}
			if err != nil {
				// Error reading segment, close and return
				closeErr := r.underlyingBody.Close()
				if closeErr != nil {
					return closeErr
				}
				if err == structuredmsg.ErrCRC64Mismatch {
					return fmt.Errorf("CRC64 checksum mismatch in structured message segment: %w", err)
				}
				return err
			}
		}
	}

	return r.underlyingBody.Close()
}

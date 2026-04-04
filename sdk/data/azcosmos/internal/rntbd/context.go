// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// Client Version
// -----------------------------------------------------------------------------

// ClientVersion is the RNTBD client version sent in context requests.
// This should match the SDK version.
const ClientVersion = "2.0.0"

// -----------------------------------------------------------------------------
// Context Request
// -----------------------------------------------------------------------------

// ContextRequest represents an RNTBD context negotiation request.
// This is sent after TLS handshake to negotiate protocol parameters.
//
// Wire format:
//
//	┌─────────────────┬──────────────┬──────────────┬─────────────────────────┬───────────────┐
//	│ length (32-bit) │ resType (16) │ opType (16)  │ activityId (128-bit)    │ headers (var) │
//	└─────────────────┴──────────────┴──────────────┴─────────────────────────┴───────────────┘
//
// ResourceType and OperationType are both "Connection" (0x0000).
type ContextRequest struct {
	ActivityID uuid.UUID
	Headers    *contextRequestHeaderValues
}

// contextRequestHeaderValues holds the headers for a context request.
type contextRequestHeaderValues struct {
	ProtocolVersion uint32
	ClientVersion   string
	UserAgent       string
}

// NewContextRequest creates a new context request with the given user agent.
func NewContextRequest(activityID uuid.UUID, userAgent string) *ContextRequest {
	return &ContextRequest{
		ActivityID: activityID,
		Headers: &contextRequestHeaderValues{
			ProtocolVersion: CurrentProtocolVersion,
			ClientVersion:   ClientVersion,
			UserAgent:       userAgent,
		},
	}
}

// ComputeLength returns the total encoded length of the context request.
func (cr *ContextRequest) ComputeLength() int {
	// Length prefix (4) + resourceType (2) + operationType (2) + activityId (16) + headers
	return RequestFrameLength + cr.Headers.computeLength()
}

// Encode writes the context request to the writer.
func (cr *ContextRequest) Encode(w io.Writer) error {
	headerLen := cr.Headers.computeLength()
	totalLength := RequestFrameLength + headerLen

	// Write length prefix (4 bytes, little-endian)
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(totalLength))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}

	// Write resourceType = Connection (0x0000)
	binary.LittleEndian.PutUint16(buf[:2], uint16(ResourceConnection))
	if _, err := w.Write(buf[:2]); err != nil {
		return err
	}

	// Write operationType = Connection (0x0000)
	binary.LittleEndian.PutUint16(buf[:2], uint16(OperationConnection))
	if _, err := w.Write(buf[:2]); err != nil {
		return err
	}

	// Write activityId (16 bytes, MS-GUID format)
	if err := WriteUUID(cr.ActivityID, w); err != nil {
		return err
	}

	// Write headers
	return cr.Headers.encode(w)
}

// EncodeToBytes encodes the context request to a byte slice.
func (cr *ContextRequest) EncodeToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := cr.Encode(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeContextRequest reads a context request from the reader.
func DecodeContextRequest(r io.Reader) (*ContextRequest, error) {
	// Read length prefix (4 bytes, little-endian)
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	expectedLength := binary.LittleEndian.Uint32(buf[:])

	// Validate length
	if expectedLength < RequestFrameLength {
		return nil, newCorruptedFrameError("context request length too small: %d", expectedLength)
	}

	// Read resourceType and operationType (verify they're both Connection)
	if _, err := io.ReadFull(r, buf[:4]); err != nil {
		return nil, err
	}
	resourceOperationCode := binary.LittleEndian.Uint32(buf[:])
	if resourceOperationCode != 0 {
		return nil, newCorruptedFrameError("context request must have resourceOperationType=0, got 0x%08X", resourceOperationCode)
	}

	// Read activityId (16 bytes, MS-GUID format)
	activityID, err := ReadUUID(r)
	if err != nil {
		return nil, err
	}

	// Read headers
	headersLength := int(expectedLength) - RequestFrameLength
	headers, err := decodeContextRequestHeaders(r, headersLength)
	if err != nil {
		return nil, err
	}

	return &ContextRequest{
		ActivityID: activityID,
		Headers:    headers,
	}, nil
}

// computeLength returns the encoded length of the headers.
func (h *contextRequestHeaderValues) computeLength() int {
	ts := h.toTokenStream()
	return ts.ComputeLength()
}

// encode writes the headers to the writer.
func (h *contextRequestHeaderValues) encode(w io.Writer) error {
	ts := h.toTokenStream()
	return ts.Encode(w)
}

// toTokenStream converts headers to a token stream.
func (h *contextRequestHeaderValues) toTokenStream() *TokenStream {
	ts := NewTokenStream()

	// ProtocolVersion - ULong (required)
	ts.SetValue(uint16(ContextRequestHeaderProtocolVersion), TokenULong, h.ProtocolVersion)

	// ClientVersion - SmallString (required)
	ts.SetValue(uint16(ContextRequestHeaderClientVersion), TokenSmallString, h.ClientVersion)

	// UserAgent - SmallString (required)
	ts.SetValue(uint16(ContextRequestHeaderUserAgent), TokenSmallString, h.UserAgent)

	return ts
}

// decodeContextRequestHeaders reads headers from the reader.
func decodeContextRequestHeaders(r io.Reader, length int) (*contextRequestHeaderValues, error) {
	ts, err := DecodeTokenStream(r, length)
	if err != nil {
		return nil, err
	}

	headers := &contextRequestHeaderValues{}

	// ProtocolVersion - ULong codec returns int64
	if t := ts.Get(uint16(ContextRequestHeaderProtocolVersion)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if n, ok := v.(int64); ok {
				headers.ProtocolVersion = uint32(n)
			}
		}
	}

	// ClientVersion
	if t := ts.Get(uint16(ContextRequestHeaderClientVersion)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if s, ok := v.(string); ok {
				headers.ClientVersion = s
			}
		}
	}

	// UserAgent
	if t := ts.Get(uint16(ContextRequestHeaderUserAgent)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if s, ok := v.(string); ok {
				headers.UserAgent = s
			}
		}
	}

	return headers, nil
}

// -----------------------------------------------------------------------------
// Context (Response)
// -----------------------------------------------------------------------------

// Context represents the RNTBD context negotiation response.
// This is received after sending a ContextRequest and contains
// server capabilities and connection parameters.
type Context struct {
	ActivityID                      uuid.UUID
	Status                          int32
	ProtocolVersion                 uint32
	ClientVersion                   string
	ServerAgent                     string
	ServerVersion                   string
	IdleTimeoutInSeconds            uint32
	UnauthenticatedTimeoutInSeconds uint32
}

// ServerProperties returns the server agent and version as a formatted string.
func (c *Context) ServerProperties() string {
	return fmt.Sprintf("%s/%s", c.ServerAgent, c.ServerVersion)
}

// IsSuccess returns true if the context negotiation succeeded.
func (c *Context) IsSuccess() bool {
	return c.Status >= 200 && c.Status < 300
}

// DecodeContext reads a context response from the reader.
// If the status code indicates an error (>= 400), returns a ContextException.
func DecodeContext(r io.Reader) (*Context, error) {
	// Read response frame (length + status + activityId)
	frame, err := DecodeResponseFrame(r)
	if err != nil {
		return nil, err
	}

	// Read headers
	headersLength := frame.HeadersLength()
	headers, err := decodeContextHeaders(r, headersLength)
	if err != nil {
		return nil, err
	}

	// Build context from frame and headers
	ctx := &Context{
		ActivityID:                      frame.ActivityID,
		Status:                          frame.Status,
		ProtocolVersion:                 headers.ProtocolVersion,
		ClientVersion:                   headers.ClientVersion,
		ServerAgent:                     headers.ServerAgent,
		ServerVersion:                   headers.ServerVersion,
		IdleTimeoutInSeconds:            headers.IdleTimeoutInSeconds,
		UnauthenticatedTimeoutInSeconds: headers.UnauthenticatedTimeoutInSeconds,
	}

	// Check for error status
	if frame.Status < 200 || frame.Status >= 300 {
		return nil, &ContextException{
			Status:          frame.Status,
			ActivityID:      frame.ActivityID,
			ClientVersion:   headers.ClientVersion,
			ProtocolVersion: headers.ProtocolVersion,
			ServerAgent:     headers.ServerAgent,
			ServerVersion:   headers.ServerVersion,
		}
	}

	return ctx, nil
}

// Encode writes the context response to the writer.
// This is primarily used for testing (mock server responses).
func (c *Context) Encode(w io.Writer) error {
	headers := &contextHeaders{
		ProtocolVersion:                 c.ProtocolVersion,
		ClientVersion:                   c.ClientVersion,
		ServerAgent:                     c.ServerAgent,
		ServerVersion:                   c.ServerVersion,
		IdleTimeoutInSeconds:            c.IdleTimeoutInSeconds,
		UnauthenticatedTimeoutInSeconds: c.UnauthenticatedTimeoutInSeconds,
	}

	headersLength := headers.computeLength()
	totalLength := ResponseFrameLength + headersLength

	// Write frame
	frame := &ResponseFrame{
		Length:     uint32(totalLength),
		Status:     c.Status,
		ActivityID: c.ActivityID,
	}
	if err := frame.Encode(w); err != nil {
		return err
	}

	// Write headers
	return headers.encode(w)
}

// EncodeToBytes encodes the context to a byte slice.
func (c *Context) EncodeToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := c.Encode(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// contextHeaders holds the decoded context response headers.
type contextHeaders struct {
	ProtocolVersion                 uint32
	ClientVersion                   string
	ServerAgent                     string
	ServerVersion                   string
	IdleTimeoutInSeconds            uint32
	UnauthenticatedTimeoutInSeconds uint32
}

// computeLength returns the encoded length of the headers.
func (h *contextHeaders) computeLength() int {
	ts := h.toTokenStream()
	return ts.ComputeLength()
}

// encode writes the headers to the writer.
func (h *contextHeaders) encode(w io.Writer) error {
	ts := h.toTokenStream()
	return ts.Encode(w)
}

// toTokenStream converts headers to a token stream.
func (h *contextHeaders) toTokenStream() *TokenStream {
	ts := NewTokenStream()

	// ProtocolVersion - ULong
	ts.SetValue(uint16(ContextHeaderProtocolVersion), TokenULong, h.ProtocolVersion)

	// ClientVersion - SmallString
	if h.ClientVersion != "" {
		ts.SetValue(uint16(ContextHeaderClientVersion), TokenSmallString, h.ClientVersion)
	}

	// ServerAgent - SmallString (required)
	ts.SetValue(uint16(ContextHeaderServerAgent), TokenSmallString, h.ServerAgent)

	// ServerVersion - SmallString (required)
	ts.SetValue(uint16(ContextHeaderServerVersion), TokenSmallString, h.ServerVersion)

	// IdleTimeoutInSeconds - ULong
	ts.SetValue(uint16(ContextHeaderIdleTimeoutInSeconds), TokenULong, h.IdleTimeoutInSeconds)

	// UnauthenticatedTimeoutInSeconds - ULong
	ts.SetValue(uint16(ContextHeaderUnauthenticatedTimeoutInSeconds), TokenULong, h.UnauthenticatedTimeoutInSeconds)

	return ts
}

// decodeContextHeaders reads context response headers from the reader.
func decodeContextHeaders(r io.Reader, length int) (*contextHeaders, error) {
	ts, err := DecodeTokenStream(r, length)
	if err != nil {
		return nil, err
	}

	headers := &contextHeaders{}

	// ProtocolVersion - ULong codec returns int64
	if t := ts.Get(uint16(ContextHeaderProtocolVersion)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if n, ok := v.(int64); ok {
				headers.ProtocolVersion = uint32(n)
			}
		}
	}

	// ClientVersion
	if t := ts.Get(uint16(ContextHeaderClientVersion)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if s, ok := v.(string); ok {
				headers.ClientVersion = s
			}
		}
	}

	// ServerAgent
	if t := ts.Get(uint16(ContextHeaderServerAgent)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if s, ok := v.(string); ok {
				headers.ServerAgent = s
			}
		}
	}

	// ServerVersion
	if t := ts.Get(uint16(ContextHeaderServerVersion)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if s, ok := v.(string); ok {
				headers.ServerVersion = s
			}
		}
	}

	// IdleTimeoutInSeconds - ULong codec returns int64
	if t := ts.Get(uint16(ContextHeaderIdleTimeoutInSeconds)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if n, ok := v.(int64); ok {
				headers.IdleTimeoutInSeconds = uint32(n)
			}
		}
	}

	// UnauthenticatedTimeoutInSeconds - ULong codec returns int64
	if t := ts.Get(uint16(ContextHeaderUnauthenticatedTimeoutInSeconds)); t != nil && t.IsPresent() {
		if v, err := t.GetValue(); err == nil {
			if n, ok := v.(int64); ok {
				headers.UnauthenticatedTimeoutInSeconds = uint32(n)
			}
		}
	}

	return headers, nil
}

// -----------------------------------------------------------------------------
// Context Exception
// -----------------------------------------------------------------------------

// ContextException is returned when context negotiation fails.
type ContextException struct {
	Status          int32
	ActivityID      uuid.UUID
	Details         map[string]interface{}
	ClientVersion   string
	ProtocolVersion uint32
	ServerAgent     string
	ServerVersion   string
}

// Error implements the error interface.
func (e *ContextException) Error() string {
	if e.Details != nil {
		detailsJSON, _ := json.Marshal(e.Details)
		return fmt.Sprintf("context negotiation failed: status=%d, activityId=%s, details=%s",
			e.Status, e.ActivityID, string(detailsJSON))
	}
	return fmt.Sprintf("context negotiation failed: status=%d, activityId=%s",
		e.Status, e.ActivityID)
}

// ResponseHeaders returns the headers received in the error response.
func (e *ContextException) ResponseHeaders() map[string]interface{} {
	headers := make(map[string]interface{})
	if e.ClientVersion != "" {
		headers["requiredClientVersion"] = e.ClientVersion
	}
	if e.ProtocolVersion != 0 {
		headers["requiredProtocolVersion"] = e.ProtocolVersion
	}
	if e.ServerAgent != "" {
		headers["serverAgent"] = e.ServerAgent
	}
	if e.ServerVersion != "" {
		headers["serverVersion"] = e.ServerVersion
	}
	return headers
}

// IsContextException returns true if the error is a ContextException.
func IsContextException(err error) bool {
	_, ok := err.(*ContextException)
	return ok
}

// -----------------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------------

// ContextFrom creates a Context response from a ContextRequest.
// This is used for testing (mock server creating responses).
func ContextFrom(request *ContextRequest, serverAgent, serverVersion string, status int32) *Context {
	return &Context{
		ActivityID:                      request.ActivityID,
		Status:                          status,
		ProtocolVersion:                 CurrentProtocolVersion,
		ClientVersion:                   request.Headers.ClientVersion,
		ServerAgent:                     serverAgent,
		ServerVersion:                   serverVersion,
		IdleTimeoutInSeconds:            0,
		UnauthenticatedTimeoutInSeconds: 0,
	}
}

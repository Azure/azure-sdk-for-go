// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sort"

	"github.com/google/uuid"
)

// -----------------------------------------------------------------------------
// Frame Constants
// -----------------------------------------------------------------------------

// Frame size constants matching Java SDK
const (
	// RequestFrameLength is the fixed header size for request frames:
	// 4 bytes (length) + 2 bytes (resourceType) + 2 bytes (operationType) + 16 bytes (activityId)
	RequestFrameLength = 4 + 2 + 2 + 16 // 24 bytes

	// ResponseFrameLength is the fixed header size for response frames:
	// 4 bytes (length) + 4 bytes (status) + 16 bytes (activityId)
	ResponseFrameLength = 4 + 4 + 16 // 24 bytes
)

// -----------------------------------------------------------------------------
// Request Frame
// -----------------------------------------------------------------------------

// RequestFrame represents the fixed header portion of an RNTBD request.
// Wire format:
//
//	┌─────────────┬─────────────┬─────────────┬─────────────────────────────┐
//	│ length (4B) │ resType(2B) │ opType (2B) │ activityId (16B, MS-GUID)   │
//	└─────────────┴─────────────┴─────────────┴─────────────────────────────┘
type RequestFrame struct {
	ResourceType  ResourceType
	OperationType OperationType
	ActivityID    uuid.UUID
}

// Encode writes the request frame to the writer.
// Note: This does NOT write the length prefix - that's handled by RequestMessage.
func (f *RequestFrame) Encode(w io.Writer) error {
	// Write resourceType (2 bytes, little-endian)
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(f.ResourceType))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}

	// Write operationType (2 bytes, little-endian)
	binary.LittleEndian.PutUint16(buf[:], uint16(f.OperationType))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}

	// Write activityId (16 bytes, MS-GUID format)
	return WriteUUID(f.ActivityID, w)
}

// DecodeRequestFrame reads a request frame from the reader.
// The length prefix should already be consumed by the framer.
func DecodeRequestFrame(r io.Reader) (*RequestFrame, error) {
	// Read resourceType (2 bytes, little-endian)
	var buf [2]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	resourceType := ResourceType(binary.LittleEndian.Uint16(buf[:]))

	// Read operationType (2 bytes, little-endian)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	operationType := OperationType(binary.LittleEndian.Uint16(buf[:]))

	// Read activityId (16 bytes, MS-GUID format)
	activityID, err := ReadUUID(r)
	if err != nil {
		return nil, err
	}

	return &RequestFrame{
		ResourceType:  resourceType,
		OperationType: operationType,
		ActivityID:    activityID,
	}, nil
}

// -----------------------------------------------------------------------------
// Response Frame (Status)
// -----------------------------------------------------------------------------

// ResponseFrame represents the fixed header portion of an RNTBD response.
// Wire format:
//
//	┌─────────────┬─────────────┬─────────────────────────────────────────────┐
//	│ length (4B) │ status (4B) │ activityId (16B, MS-GUID)                   │
//	└─────────────┴─────────────┴─────────────────────────────────────────────┘
type ResponseFrame struct {
	Length     uint32    // Total length of headers (including this frame, excluding payload)
	Status     int32     // HTTP-like status code
	ActivityID uuid.UUID // Activity ID in MS-GUID format
}

// HeadersLength returns the length of the headers portion (excluding the fixed frame).
func (f *ResponseFrame) HeadersLength() int {
	return int(f.Length) - ResponseFrameLength
}

// Encode writes the response frame to the writer.
func (f *ResponseFrame) Encode(w io.Writer) error {
	// Write length (4 bytes, little-endian)
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], f.Length)
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}

	// Write status (4 bytes, little-endian, signed)
	binary.LittleEndian.PutUint32(buf[:], uint32(f.Status))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}

	// Write activityId (16 bytes, MS-GUID format)
	return WriteUUID(f.ActivityID, w)
}

// DecodeResponseFrame reads a response frame from the reader.
func DecodeResponseFrame(r io.Reader) (*ResponseFrame, error) {
	// Read length (4 bytes, little-endian, unsigned)
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(buf[:])

	// Validate length
	if length < ResponseFrameLength {
		return nil, newCorruptedFrameError("response frame length too small: %d", length)
	}

	// Read status (4 bytes, little-endian, signed)
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	status := int32(binary.LittleEndian.Uint32(buf[:]))

	// Read activityId (16 bytes, MS-GUID format)
	activityID, err := ReadUUID(r)
	if err != nil {
		return nil, err
	}

	return &ResponseFrame{
		Length:     length,
		Status:     status,
		ActivityID: activityID,
	}, nil
}

// -----------------------------------------------------------------------------
// Token Stream
// -----------------------------------------------------------------------------

// TokenStream represents a sequence of tokens in an RNTBD frame.
type TokenStream struct {
	tokens map[uint16]*Token // Map from token ID to token
}

// NewTokenStream creates a new empty token stream.
func NewTokenStream() *TokenStream {
	return &TokenStream{
		tokens: make(map[uint16]*Token),
	}
}

// Get returns the token with the given ID, or nil if not present.
func (ts *TokenStream) Get(id uint16) *Token {
	return ts.tokens[id]
}

// Set adds or updates a token in the stream.
func (ts *TokenStream) Set(token *Token) {
	ts.tokens[token.ID] = token
}

// SetValue creates and sets a token with the given ID, type, and value.
func (ts *TokenStream) SetValue(id uint16, tokenType TokenType, value interface{}) error {
	token, err := NewToken(id, tokenType, value)
	if err != nil {
		return err
	}
	ts.Set(token)
	return nil
}

// GetValue returns the value of the token with the given ID.
// Returns the default value for the token type if not present.
func (ts *TokenStream) GetValue(id uint16, tokenType TokenType) (interface{}, error) {
	token := ts.tokens[id]
	if token == nil || !token.IsPresent() {
		return GetCodec(tokenType).DefaultValue(), nil
	}
	return token.GetValue()
}

// -----------------------------------------------------------------------------
// TokenStream Helper Methods - Type-Safe Setters
// -----------------------------------------------------------------------------

// SetByte sets a byte token value.
func (ts *TokenStream) SetByte(id uint16, value byte) {
	ts.SetValue(id, TokenByte, value) //nolint:errcheck
}

// SetUShort sets an unsigned 16-bit token value.
func (ts *TokenStream) SetUShort(id uint16, value uint16) {
	ts.SetValue(id, TokenUShort, value) //nolint:errcheck
}

// SetULong sets an unsigned 32-bit token value.
func (ts *TokenStream) SetULong(id uint16, value uint32) {
	ts.SetValue(id, TokenULong, value) //nolint:errcheck
}

// SetLong sets a signed 32-bit token value.
func (ts *TokenStream) SetLong(id uint16, value int32) {
	ts.SetValue(id, TokenLong, value) //nolint:errcheck
}

// SetULongLong sets an unsigned 64-bit token value.
func (ts *TokenStream) SetULongLong(id uint16, value uint64) {
	ts.SetValue(id, TokenULongLong, value) //nolint:errcheck
}

// SetLongLong sets a signed 64-bit token value.
func (ts *TokenStream) SetLongLong(id uint16, value int64) {
	ts.SetValue(id, TokenLongLong, value) //nolint:errcheck
}

// SetFloat sets a 32-bit floating point token value.
func (ts *TokenStream) SetFloat(id uint16, value float32) {
	ts.SetValue(id, TokenFloat, value) //nolint:errcheck
}

// SetDouble sets a 64-bit floating point token value.
func (ts *TokenStream) SetDouble(id uint16, value float64) {
	ts.SetValue(id, TokenDouble, value) //nolint:errcheck
}

// SetString sets a string token value.
func (ts *TokenStream) SetString(id uint16, value string) {
	ts.SetValue(id, TokenString, value) //nolint:errcheck
}

// SetBytes sets a byte array token value.
func (ts *TokenStream) SetBytes(id uint16, value []byte) {
	ts.SetValue(id, TokenBytes, value) //nolint:errcheck
}

// SetGUID sets a GUID token value.
func (ts *TokenStream) SetGUID(id uint16, value uuid.UUID) {
	ts.SetValue(id, TokenGuid, value) //nolint:errcheck
}

// -----------------------------------------------------------------------------
// TokenStream Helper Methods - Type-Safe Getters
// -----------------------------------------------------------------------------

// GetByte returns the byte value of a token, or 0 if not present.
func (ts *TokenStream) GetByte(id uint16) byte {
	val, err := ts.GetValue(id, TokenByte)
	if err != nil || val == nil {
		return 0
	}
	if b, ok := val.(byte); ok {
		return b
	}
	return 0
}

// GetUShort returns the unsigned 16-bit value of a token, or 0 if not present.
func (ts *TokenStream) GetUShort(id uint16) uint16 {
	val, err := ts.GetValue(id, TokenUShort)
	if err != nil || val == nil {
		return 0
	}
	if v, ok := val.(uint16); ok {
		return v
	}
	return 0
}

// GetULong returns the unsigned 32-bit value of a token, or 0 if not present.
// Handles both uint32 (from SetULong) and int64 (from codec decode).
func (ts *TokenStream) GetULong(id uint16) uint32 {
	val, err := ts.GetValue(id, TokenULong)
	if err != nil || val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return uint32(v)
	case uint32:
		return v
	case int32:
		return uint32(v)
	case int:
		return uint32(v)
	}
	return 0
}

// GetLong returns the signed 32-bit value of a token, or 0 if not present.
func (ts *TokenStream) GetLong(id uint16) int32 {
	val, err := ts.GetValue(id, TokenLong)
	if err != nil || val == nil {
		return 0
	}
	if v, ok := val.(int32); ok {
		return v
	}
	return 0
}

// GetULongLong returns the unsigned 64-bit value of a token, or 0 if not present.
// Handles both uint64 (from SetULongLong) and int64 (from codec decode).
func (ts *TokenStream) GetULongLong(id uint16) int64 {
	val, err := ts.GetValue(id, TokenULongLong)
	if err != nil || val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case uint64:
		return int64(v)
	}
	return 0
}

// GetLongLong returns the signed 64-bit value of a token, or 0 if not present.
func (ts *TokenStream) GetLongLong(id uint16) int64 {
	val, err := ts.GetValue(id, TokenLongLong)
	if err != nil || val == nil {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	}
	return 0
}

// GetFloat returns the 32-bit floating point value of a token, or 0 if not present.
func (ts *TokenStream) GetFloat(id uint16) float32 {
	val, err := ts.GetValue(id, TokenFloat)
	if err != nil || val == nil {
		return 0
	}
	if f, ok := val.(float32); ok {
		return f
	}
	return 0
}

// GetDouble returns the 64-bit floating point value of a token, or 0 if not present.
func (ts *TokenStream) GetDouble(id uint16) float64 {
	val, err := ts.GetValue(id, TokenDouble)
	if err != nil || val == nil {
		return 0
	}
	if f, ok := val.(float64); ok {
		return f
	}
	return 0
}

// GetString returns the string value of a token, or empty string if not present.
func (ts *TokenStream) GetString(id uint16) string {
	val, err := ts.GetValue(id, TokenString)
	if err != nil || val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// GetBytes returns the byte array value of a token, or nil if not present.
func (ts *TokenStream) GetBytes(id uint16) []byte {
	val, err := ts.GetValue(id, TokenBytes)
	if err != nil || val == nil {
		return nil
	}
	if b, ok := val.([]byte); ok {
		return b
	}
	return nil
}

// GetGUID returns the GUID value of a token, or empty UUID if not present.
func (ts *TokenStream) GetGUID(id uint16) uuid.UUID {
	val, err := ts.GetValue(id, TokenGuid)
	if err != nil || val == nil {
		return EmptyUUID
	}
	if g, ok := val.(uuid.UUID); ok {
		return g
	}
	return EmptyUUID
}

// Count returns the number of present tokens in the stream.
func (ts *TokenStream) Count() int {
	count := 0
	for _, t := range ts.tokens {
		if t.IsPresent() {
			count++
		}
	}
	return count
}

// Tokens returns all present tokens in the stream.
func (ts *TokenStream) Tokens() []*Token {
	result := make([]*Token, 0, len(ts.tokens))
	for _, t := range ts.tokens {
		if t.IsPresent() {
			result = append(result, t)
		}
	}
	return result
}

// ComputeLength returns the total encoded length of all tokens.
func (ts *TokenStream) ComputeLength() int {
	length := 0
	for _, t := range ts.tokens {
		length += t.ComputeLength()
	}
	return length
}

// Encode writes all present tokens to the writer in sorted order by token ID.
// Sorting ensures deterministic wire format, which is required by the Cosmos DB server.
func (ts *TokenStream) Encode(w io.Writer) error {
	ids := make([]uint16, 0, len(ts.tokens))
	for id, t := range ts.tokens {
		if t.IsPresent() {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		if err := ts.tokens[id].Encode(w); err != nil {
			return err
		}
	}
	return nil
}

// DecodeTokenStream reads tokens from the reader until the specified length is consumed.
func DecodeTokenStream(r io.Reader, length int) (*TokenStream, error) {
	ts := NewTokenStream()

	if length <= 0 {
		return ts, nil
	}

	// Wrap in a limited reader to track consumption
	lr := io.LimitReader(r, int64(length))

	for {
		token, err := DecodeToken(lr)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		ts.Set(token)
	}

	return ts, nil
}

// -----------------------------------------------------------------------------
// Request Message
// -----------------------------------------------------------------------------

// RequestMessage represents a complete RNTBD request (frame + headers + optional payload).
type RequestMessage struct {
	Frame   *RequestFrame
	Headers *TokenStream
	Payload []byte
}

// NewRequestMessage creates a new request message with the given parameters.
func NewRequestMessage(resourceType ResourceType, operationType OperationType, activityID uuid.UUID) *RequestMessage {
	return &RequestMessage{
		Frame: &RequestFrame{
			ResourceType:  resourceType,
			OperationType: operationType,
			ActivityID:    activityID,
		},
		Headers: NewTokenStream(),
	}
}

// ComputeLength returns the total encoded length of the request message.
func (m *RequestMessage) ComputeLength() int {
	// Frame: 4 (length prefix) + 2 (resourceType) + 2 (operationType) + 16 (activityId) = 24
	// Plus headers
	headLength := RequestFrameLength + m.Headers.ComputeLength()

	// If there's a payload, add 4 bytes for payload length + payload bytes
	if len(m.Payload) > 0 {
		return headLength + 4 + len(m.Payload)
	}
	return headLength
}

// Encode writes the complete request message to the writer.
func (m *RequestMessage) Encode(w io.Writer) error {
	// Calculate head length (frame + headers, NOT including payload)
	headLength := RequestFrameLength + m.Headers.ComputeLength()

	// Write length prefix (4 bytes, little-endian)
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(headLength))
	if _, err := w.Write(buf[:]); err != nil {
		return err
	}

	// Write frame (resourceType, operationType, activityId)
	if err := m.Frame.Encode(w); err != nil {
		return err
	}

	// Write headers
	if err := m.Headers.Encode(w); err != nil {
		return err
	}

	// Write payload if present
	if len(m.Payload) > 0 {
		// Write payload length (4 bytes, little-endian)
		binary.LittleEndian.PutUint32(buf[:], uint32(len(m.Payload)))
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
		// Write payload
		if _, err := w.Write(m.Payload); err != nil {
			return err
		}
	}

	return nil
}

// DecodeRequestMessage reads a complete request message from the reader.
func DecodeRequestMessage(r io.Reader) (*RequestMessage, error) {
	// Read length prefix (4 bytes, little-endian)
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return nil, err
	}
	headLength := binary.LittleEndian.Uint32(buf[:])

	// Validate length
	if headLength < RequestFrameLength {
		return nil, newCorruptedFrameError("request head length too small: %d", headLength)
	}

	// Read frame
	frame, err := DecodeRequestFrame(r)
	if err != nil {
		return nil, err
	}

	// Read headers
	headersLength := int(headLength) - RequestFrameLength
	headers, err := DecodeTokenStream(r, headersLength)
	if err != nil {
		return nil, err
	}

	// Check for payload (look for PayloadPresent header)
	var payload []byte
	payloadPresent := headers.Get(uint16(RequestHeaderPayloadPresent))
	if payloadPresent != nil && payloadPresent.IsPresent() {
		val, err := payloadPresent.GetValue()
		if err != nil {
			return nil, err
		}
		if b, ok := val.(byte); ok && b != 0 {
			// Read payload length
			if _, err := io.ReadFull(r, buf[:]); err != nil {
				return nil, err
			}
			payloadLength := binary.LittleEndian.Uint32(buf[:])

			// Read payload
			payload = make([]byte, payloadLength)
			if _, err := io.ReadFull(r, payload); err != nil {
				return nil, err
			}
		}
	}

	return &RequestMessage{
		Frame:   frame,
		Headers: headers,
		Payload: payload,
	}, nil
}

// -----------------------------------------------------------------------------
// Response Message
// -----------------------------------------------------------------------------

// ResponseMessage represents a complete RNTBD response (frame + headers + optional payload).
type ResponseMessage struct {
	Frame   *ResponseFrame
	Headers *TokenStream
	Payload []byte
}

// NewResponseMessage creates a new response message with the given parameters.
func NewResponseMessage(status int32, activityID uuid.UUID) *ResponseMessage {
	return &ResponseMessage{
		Frame: &ResponseFrame{
			Status:     status,
			ActivityID: activityID,
		},
		Headers: NewTokenStream(),
	}
}

// ComputeLength returns the total encoded length of the response message.
func (m *ResponseMessage) ComputeLength() int {
	// Frame: 4 (length) + 4 (status) + 16 (activityId) = 24
	// Plus headers
	headLength := ResponseFrameLength + m.Headers.ComputeLength()

	// If there's a payload, add 4 bytes for payload length + payload bytes
	if len(m.Payload) > 0 {
		return headLength + 4 + len(m.Payload)
	}
	return headLength
}

// Encode writes the complete response message to the writer.
func (m *ResponseMessage) Encode(w io.Writer) error {
	// Calculate and set head length (frame + headers, NOT including payload)
	headLength := ResponseFrameLength + m.Headers.ComputeLength()
	m.Frame.Length = uint32(headLength)

	// Write frame (length, status, activityId)
	if err := m.Frame.Encode(w); err != nil {
		return err
	}

	// Write headers
	if err := m.Headers.Encode(w); err != nil {
		return err
	}

	// Write payload if present
	if len(m.Payload) > 0 {
		// Write payload length (4 bytes, little-endian)
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], uint32(len(m.Payload)))
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
		// Write payload
		if _, err := w.Write(m.Payload); err != nil {
			return err
		}
	}

	return nil
}

// DecodeResponseMessage reads a complete response message from the reader.
// Note: The caller should verify there's enough data available before calling.
func DecodeResponseMessage(r io.Reader) (*ResponseMessage, error) {
	// Read frame
	frame, err := DecodeResponseFrame(r)
	if err != nil {
		return nil, err
	}

	// Read headers
	headersLength := frame.HeadersLength()
	headers, err := DecodeTokenStream(r, headersLength)
	if err != nil {
		return nil, err
	}

	// Check for payload (look for PayloadPresent header)
	var payload []byte
	payloadPresent := headers.Get(uint16(ResponseHeaderPayloadPresent))
	if payloadPresent != nil && payloadPresent.IsPresent() {
		val, err := payloadPresent.GetValue()
		if err != nil {
			return nil, err
		}
		if b, ok := val.(byte); ok && b != 0 {
			// Read payload length
			var buf [4]byte
			if _, err := io.ReadFull(r, buf[:]); err != nil {
				return nil, err
			}
			payloadLength := binary.LittleEndian.Uint32(buf[:])

			// Read payload
			payload = make([]byte, payloadLength)
			if _, err := io.ReadFull(r, payload); err != nil {
				return nil, err
			}
		}
	}

	return &ResponseMessage{
		Frame:   frame,
		Headers: headers,
		Payload: payload,
	}, nil
}

// -----------------------------------------------------------------------------
// Utility Functions
// -----------------------------------------------------------------------------

// EncodeRequestToBytes encodes a request message to a byte slice.
func EncodeRequestToBytes(m *RequestMessage) ([]byte, error) {
	var buf bytes.Buffer
	if err := m.Encode(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeResponseToBytes encodes a response message to a byte slice.
func EncodeResponseToBytes(m *ResponseMessage) ([]byte, error) {
	var buf bytes.Buffer
	if err := m.Encode(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeRequestFromBytes decodes a request message from a byte slice.
func DecodeRequestFromBytes(data []byte) (*RequestMessage, error) {
	return DecodeRequestMessage(bytes.NewReader(data))
}

// DecodeResponseFromBytes decodes a response message from a byte slice.
func DecodeResponseFromBytes(data []byte) (*ResponseMessage, error) {
	return DecodeResponseMessage(bytes.NewReader(data))
}

// String returns a debug string for the request frame.
func (f *RequestFrame) String() string {
	return fmt.Sprintf("RequestFrame{ResourceType: %s, OperationType: %s, ActivityID: %s}",
		f.ResourceType, f.OperationType, f.ActivityID)
}

// String returns a debug string for the response frame.
func (f *ResponseFrame) String() string {
	return fmt.Sprintf("ResponseFrame{Length: %d, Status: %d, ActivityID: %s}",
		f.Length, f.Status, f.ActivityID)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// Request Frame Tests
// -----------------------------------------------------------------------------

func TestRequestFrame_Encode(t *testing.T) {
	// Test encoding a basic request frame
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")

	frame := &RequestFrame{
		ResourceType:  ResourceDocument,
		OperationType: OperationRead,
		ActivityID:    activityID,
	}

	var buf bytes.Buffer
	err := frame.Encode(&buf)
	require.NoError(t, err)

	// Frame should be 20 bytes (2 + 2 + 16) - length prefix NOT included by Frame.Encode
	require.Equal(t, 20, buf.Len())

	data := buf.Bytes()

	// Verify resourceType (little-endian)
	require.Equal(t, uint16(ResourceDocument), binary.LittleEndian.Uint16(data[0:2]))

	// Verify operationType (little-endian)
	require.Equal(t, uint16(OperationRead), binary.LittleEndian.Uint16(data[2:4]))

	// Verify activityID (16 bytes in MS-GUID format)
	// We'll decode it back and compare
	decodedID, err := DecodeUUID(data[4:20])
	require.NoError(t, err)
	require.Equal(t, activityID, decodedID)
}

func TestDecodeRequestFrame(t *testing.T) {
	activityID := uuid.MustParse("abcdef12-3456-7890-abcd-ef1234567890")

	// Manually construct the binary representation
	var buf bytes.Buffer

	// Write resourceType (little-endian)
	binary.Write(&buf, binary.LittleEndian, uint16(ResourceCollection))

	// Write operationType (little-endian)
	binary.Write(&buf, binary.LittleEndian, uint16(OperationCreate))

	// Write activityID (MS-GUID format)
	WriteUUID(activityID, &buf)

	// Decode the frame
	frame, err := DecodeRequestFrame(&buf)
	require.NoError(t, err)
	require.NotNil(t, frame)

	require.Equal(t, ResourceCollection, frame.ResourceType)
	require.Equal(t, OperationCreate, frame.OperationType)
	require.Equal(t, activityID, frame.ActivityID)
}

func TestRequestFrame_RoundTrip(t *testing.T) {
	testCases := []struct {
		name          string
		resourceType  ResourceType
		operationType OperationType
		activityID    uuid.UUID
	}{
		{
			name:          "document read",
			resourceType:  ResourceDocument,
			operationType: OperationRead,
			activityID:    uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:          "collection create",
			resourceType:  ResourceCollection,
			operationType: OperationCreate,
			activityID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		{
			name:          "database delete",
			resourceType:  ResourceDatabase,
			operationType: OperationDelete,
			activityID:    uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		},
		{
			name:          "sproc execute",
			resourceType:  ResourceStoredProcedure,
			operationType: OperationExecuteJavaScript,
			activityID:    uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		},
		{
			name:          "nil UUID",
			resourceType:  ResourceDocument,
			operationType: OperationRead,
			activityID:    uuid.Nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := &RequestFrame{
				ResourceType:  tc.resourceType,
				OperationType: tc.operationType,
				ActivityID:    tc.activityID,
			}

			// Encode
			var buf bytes.Buffer
			err := original.Encode(&buf)
			require.NoError(t, err)

			// Decode
			decoded, err := DecodeRequestFrame(&buf)
			require.NoError(t, err)

			// Compare
			require.Equal(t, original.ResourceType, decoded.ResourceType)
			require.Equal(t, original.OperationType, decoded.OperationType)
			require.Equal(t, original.ActivityID, decoded.ActivityID)
		})
	}
}

// -----------------------------------------------------------------------------
// Response Frame Tests
// -----------------------------------------------------------------------------

func TestResponseFrame_Encode(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")

	frame := &ResponseFrame{
		Length:     ResponseFrameLength + 100, // header + 100 bytes of tokens
		Status:     200,
		ActivityID: activityID,
	}

	var buf bytes.Buffer
	err := frame.Encode(&buf)
	require.NoError(t, err)

	// Full response frame is 24 bytes
	require.Equal(t, ResponseFrameLength, buf.Len())

	data := buf.Bytes()

	// Verify length (little-endian)
	require.Equal(t, uint32(ResponseFrameLength+100), binary.LittleEndian.Uint32(data[0:4]))

	// Verify status (little-endian, signed)
	require.Equal(t, int32(200), int32(binary.LittleEndian.Uint32(data[4:8])))

	// Verify activityID
	decodedID, err := DecodeUUID(data[8:24])
	require.NoError(t, err)
	require.Equal(t, activityID, decodedID)
}

func TestDecodeResponseFrame(t *testing.T) {
	activityID := uuid.MustParse("abcdef12-3456-7890-abcd-ef1234567890")

	var buf bytes.Buffer

	// Write length (little-endian)
	binary.Write(&buf, binary.LittleEndian, uint32(ResponseFrameLength+50))

	// Write status (little-endian, signed)
	binary.Write(&buf, binary.LittleEndian, int32(404))

	// Write activityID (MS-GUID format)
	WriteUUID(activityID, &buf)

	// Decode the frame
	frame, err := DecodeResponseFrame(&buf)
	require.NoError(t, err)
	require.NotNil(t, frame)

	require.Equal(t, uint32(ResponseFrameLength+50), frame.Length)
	require.Equal(t, int32(404), frame.Status)
	require.Equal(t, activityID, frame.ActivityID)
}

func TestDecodeResponseFrame_TooSmall(t *testing.T) {
	var buf bytes.Buffer

	// Write length smaller than minimum
	binary.Write(&buf, binary.LittleEndian, uint32(10)) // Less than ResponseFrameLength
	binary.Write(&buf, binary.LittleEndian, int32(200))
	WriteUUID(uuid.Nil, &buf)

	_, err := DecodeResponseFrame(&buf)
	require.Error(t, err)
	require.Contains(t, err.Error(), "too small")
}

func TestResponseFrame_RoundTrip(t *testing.T) {
	testCases := []struct {
		name       string
		length     uint32
		status     int32
		activityID uuid.UUID
	}{
		{
			name:       "success",
			length:     ResponseFrameLength,
			status:     200,
			activityID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		{
			name:       "created",
			length:     ResponseFrameLength + 100,
			status:     201,
			activityID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		{
			name:       "not found",
			length:     ResponseFrameLength + 50,
			status:     404,
			activityID: uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		},
		{
			name:       "server error",
			length:     ResponseFrameLength + 200,
			status:     500,
			activityID: uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := &ResponseFrame{
				Length:     tc.length,
				Status:     tc.status,
				ActivityID: tc.activityID,
			}

			// Encode
			var buf bytes.Buffer
			err := original.Encode(&buf)
			require.NoError(t, err)

			// Decode
			decoded, err := DecodeResponseFrame(&buf)
			require.NoError(t, err)

			// Compare
			require.Equal(t, original.Length, decoded.Length)
			require.Equal(t, original.Status, decoded.Status)
			require.Equal(t, original.ActivityID, decoded.ActivityID)
		})
	}
}

// -----------------------------------------------------------------------------
// Token Stream Tests
// -----------------------------------------------------------------------------

func TestTokenStream_Basic(t *testing.T) {
	ts := NewTokenStream()
	require.NotNil(t, ts)
	require.Equal(t, 0, ts.Count())

	// Set a token
	err := ts.SetValue(0x0001, TokenByte, byte(1))
	require.NoError(t, err)
	require.Equal(t, 1, ts.Count())

	// Get the token
	token := ts.Get(0x0001)
	require.NotNil(t, token)
	require.True(t, token.IsPresent())

	// Get the value
	val, err := ts.GetValue(0x0001, TokenByte)
	require.NoError(t, err)
	require.Equal(t, byte(1), val)

	// Get missing token - should return default
	val, err = ts.GetValue(0x9999, TokenByte)
	require.NoError(t, err)
	require.Equal(t, byte(0), val) // default
}

func TestTokenStream_MultipleTokens(t *testing.T) {
	ts := NewTokenStream()

	// Add various tokens
	require.NoError(t, ts.SetValue(0x0001, TokenByte, byte(42)))
	require.NoError(t, ts.SetValue(0x0002, TokenUShort, uint16(1234)))
	require.NoError(t, ts.SetValue(0x0003, TokenULong, uint32(56789)))
	require.NoError(t, ts.SetValue(0x0004, TokenString, "hello"))

	require.Equal(t, 4, ts.Count())

	// Verify values
	val, _ := ts.GetValue(0x0001, TokenByte)
	require.Equal(t, byte(42), val)

	val, _ = ts.GetValue(0x0002, TokenUShort)
	require.Equal(t, uint16(1234), val)

	val, _ = ts.GetValue(0x0003, TokenULong)
	require.Equal(t, uint32(56789), val)

	val, _ = ts.GetValue(0x0004, TokenString)
	require.Equal(t, "hello", val)
}

func TestTokenStream_RoundTrip(t *testing.T) {
	// Create and populate token stream
	ts := NewTokenStream()
	require.NoError(t, ts.SetValue(0x0001, TokenByte, byte(42)))
	require.NoError(t, ts.SetValue(0x0002, TokenUShort, uint16(1234)))
	require.NoError(t, ts.SetValue(0x0003, TokenString, "test"))

	// Encode
	var buf bytes.Buffer
	length := ts.ComputeLength()
	err := ts.Encode(&buf)
	require.NoError(t, err)
	require.Equal(t, length, buf.Len())

	// Decode
	decoded, err := DecodeTokenStream(&buf, length)
	require.NoError(t, err)
	require.Equal(t, ts.Count(), decoded.Count())

	// Verify values
	val, _ := decoded.GetValue(0x0001, TokenByte)
	require.Equal(t, byte(42), val)

	val, _ = decoded.GetValue(0x0002, TokenUShort)
	require.Equal(t, uint16(1234), val)

	val, _ = decoded.GetValue(0x0003, TokenString)
	require.Equal(t, "test", val)
}

// -----------------------------------------------------------------------------
// Request Message Tests
// -----------------------------------------------------------------------------

func TestNewRequestMessage(t *testing.T) {
	activityID := uuid.New()
	msg := NewRequestMessage(ResourceDocument, OperationRead, activityID)

	require.NotNil(t, msg)
	require.NotNil(t, msg.Frame)
	require.NotNil(t, msg.Headers)
	require.Nil(t, msg.Payload)

	require.Equal(t, ResourceDocument, msg.Frame.ResourceType)
	require.Equal(t, OperationRead, msg.Frame.OperationType)
	require.Equal(t, activityID, msg.Frame.ActivityID)
}

func TestRequestMessage_EncodeSimple(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	msg := NewRequestMessage(ResourceDocument, OperationRead, activityID)

	// No headers, no payload
	data, err := EncodeRequestToBytes(msg)
	require.NoError(t, err)

	// Should be exactly RequestFrameLength (24 bytes)
	require.Equal(t, RequestFrameLength, len(data))

	// Verify length prefix
	length := binary.LittleEndian.Uint32(data[0:4])
	require.Equal(t, uint32(RequestFrameLength), length)
}

func TestRequestMessage_EncodeWithHeaders(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	msg := NewRequestMessage(ResourceDocument, OperationRead, activityID)

	// Add some headers
	require.NoError(t, msg.Headers.SetValue(uint16(RequestHeaderResourceId), TokenString, "test-resource"))

	data, err := EncodeRequestToBytes(msg)
	require.NoError(t, err)

	// Should be more than base frame length
	require.Greater(t, len(data), RequestFrameLength)

	// Length prefix should match header portion
	length := binary.LittleEndian.Uint32(data[0:4])
	require.Equal(t, uint32(len(data)), length) // No payload
}

func TestRequestMessage_EncodeWithPayload(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	msg := NewRequestMessage(ResourceDocument, OperationCreate, activityID)

	// Add payload
	msg.Payload = []byte(`{"id": "doc1", "data": "test"}`)

	// Add PayloadPresent header (required for decoding)
	require.NoError(t, msg.Headers.SetValue(uint16(RequestHeaderPayloadPresent), TokenByte, byte(1)))

	data, err := EncodeRequestToBytes(msg)
	require.NoError(t, err)

	// Total = header + 4 (payload length) + payload
	expectedLength := msg.ComputeLength()
	require.Equal(t, expectedLength, len(data))
}

func TestRequestMessage_RoundTrip(t *testing.T) {
	activityID := uuid.MustParse("abcdef12-3456-7890-abcd-ef1234567890")
	msg := NewRequestMessage(ResourceCollection, OperationQuery, activityID)

	// Add headers
	require.NoError(t, msg.Headers.SetValue(uint16(RequestHeaderResourceId), TokenString, "my-collection"))
	require.NoError(t, msg.Headers.SetValue(uint16(RequestHeaderPayloadPresent), TokenByte, byte(1)))

	// Add payload
	msg.Payload = []byte(`SELECT * FROM c`)

	// Encode
	data, err := EncodeRequestToBytes(msg)
	require.NoError(t, err)

	// Decode
	decoded, err := DecodeRequestFromBytes(data)
	require.NoError(t, err)

	// Verify frame
	require.Equal(t, msg.Frame.ResourceType, decoded.Frame.ResourceType)
	require.Equal(t, msg.Frame.OperationType, decoded.Frame.OperationType)
	require.Equal(t, msg.Frame.ActivityID, decoded.Frame.ActivityID)

	// Verify headers
	val, _ := decoded.Headers.GetValue(uint16(RequestHeaderResourceId), TokenString)
	require.Equal(t, "my-collection", val)

	// Verify payload
	require.Equal(t, msg.Payload, decoded.Payload)
}

// -----------------------------------------------------------------------------
// Response Message Tests
// -----------------------------------------------------------------------------

func TestNewResponseMessage(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)

	require.NotNil(t, msg)
	require.NotNil(t, msg.Frame)
	require.NotNil(t, msg.Headers)
	require.Nil(t, msg.Payload)

	require.Equal(t, int32(200), msg.Frame.Status)
	require.Equal(t, activityID, msg.Frame.ActivityID)
}

func TestResponseMessage_EncodeSimple(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	msg := NewResponseMessage(200, activityID)

	// No headers, no payload
	data, err := EncodeResponseToBytes(msg)
	require.NoError(t, err)

	// Should be exactly ResponseFrameLength (24 bytes)
	require.Equal(t, ResponseFrameLength, len(data))

	// Verify length prefix (should equal full header length)
	length := binary.LittleEndian.Uint32(data[0:4])
	require.Equal(t, uint32(ResponseFrameLength), length)
}

func TestResponseMessage_EncodeWithHeaders(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	msg := NewResponseMessage(200, activityID)

	// Add some headers
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderLSN), TokenLongLong, int64(12345)))

	data, err := EncodeResponseToBytes(msg)
	require.NoError(t, err)

	// Should be more than base frame length
	require.Greater(t, len(data), ResponseFrameLength)
}

func TestResponseMessage_EncodeWithPayload(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	msg := NewResponseMessage(200, activityID)

	// Add payload
	msg.Payload = []byte(`{"id": "doc1", "data": "response data"}`)

	// Add PayloadPresent header
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderPayloadPresent), TokenByte, byte(1)))

	data, err := EncodeResponseToBytes(msg)
	require.NoError(t, err)

	// Total = header + 4 (payload length) + payload
	expectedLength := msg.ComputeLength()
	require.Equal(t, expectedLength, len(data))
}

func TestResponseMessage_RoundTrip(t *testing.T) {
	activityID := uuid.MustParse("fedcba98-7654-3210-fedc-ba9876543210")
	msg := NewResponseMessage(201, activityID)

	// Add headers
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderLSN), TokenLongLong, int64(99999)))
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderPayloadPresent), TokenByte, byte(1)))

	// Add payload
	msg.Payload = []byte(`{"_rid": "xyz", "_etag": "\"abc\""}`)

	// Encode
	data, err := EncodeResponseToBytes(msg)
	require.NoError(t, err)

	// Decode
	decoded, err := DecodeResponseFromBytes(data)
	require.NoError(t, err)

	// Verify frame
	require.Equal(t, msg.Frame.Status, decoded.Frame.Status)
	require.Equal(t, msg.Frame.ActivityID, decoded.Frame.ActivityID)

	// Verify headers
	val, _ := decoded.Headers.GetValue(uint16(ResponseHeaderLSN), TokenLongLong)
	require.Equal(t, int64(99999), val)

	// Verify payload
	require.Equal(t, msg.Payload, decoded.Payload)
}

func TestResponseMessage_RoundTripErrorStatus(t *testing.T) {
	testCases := []struct {
		name   string
		status int32
	}{
		{"BadRequest", 400},
		{"Unauthorized", 401},
		{"Forbidden", 403},
		{"NotFound", 404},
		{"Conflict", 409},
		{"Gone", 410},
		{"TooManyRequests", 429},
		{"InternalServerError", 500},
		{"ServiceUnavailable", 503},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			activityID := uuid.New()
			msg := NewResponseMessage(tc.status, activityID)

			// Add error payload
			msg.Payload = []byte(`{"code": "` + tc.name + `", "message": "Error occurred"}`)
			require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderPayloadPresent), TokenByte, byte(1)))

			// Round trip
			data, err := EncodeResponseToBytes(msg)
			require.NoError(t, err)

			decoded, err := DecodeResponseFromBytes(data)
			require.NoError(t, err)

			require.Equal(t, tc.status, decoded.Frame.Status)
			require.Equal(t, msg.Payload, decoded.Payload)
		})
	}
}

// -----------------------------------------------------------------------------
// Error Mapping Tests
// -----------------------------------------------------------------------------

func TestErrorFromResponse_Success(t *testing.T) {
	// 2xx status codes should return nil
	msg := NewResponseMessage(200, uuid.New())
	err := ErrorFromResponse(msg)
	require.Nil(t, err)

	msg = NewResponseMessage(201, uuid.New())
	err = ErrorFromResponse(msg)
	require.Nil(t, err)

	msg = NewResponseMessage(204, uuid.New())
	err = ErrorFromResponse(msg)
	require.Nil(t, err)
}

func TestErrorFromResponse_BadRequest(t *testing.T) {
	msg := NewResponseMessage(400, uuid.New())
	msg.Payload = []byte(`{"code": "BadRequest", "message": "Invalid request"}`)

	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsBadRequest(err))

	badReq, ok := err.(*BadRequestError)
	require.True(t, ok)
	require.Equal(t, int32(400), badReq.StatusCode)
}

func TestErrorFromResponse_NotFound(t *testing.T) {
	msg := NewResponseMessage(404, uuid.New())
	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsNotFound(err))
}

func TestErrorFromResponse_Conflict(t *testing.T) {
	msg := NewResponseMessage(409, uuid.New())
	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsConflict(err))
}

func TestErrorFromResponse_Gone_Basic(t *testing.T) {
	msg := NewResponseMessage(410, uuid.New())
	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsGone(err))

	_, ok := err.(*GoneError)
	require.True(t, ok)
}

func TestErrorFromResponse_Gone_InvalidPartition(t *testing.T) {
	msg := NewResponseMessage(410, uuid.New())
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderSubStatus), TokenULong, uint32(SubStatusNameCacheIsStale)))

	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsGone(err))

	_, ok := err.(*InvalidPartitionError)
	require.True(t, ok)
}

func TestErrorFromResponse_Gone_PartitionKeyRangeGone(t *testing.T) {
	msg := NewResponseMessage(410, uuid.New())
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderSubStatus), TokenULong, uint32(SubStatusPartitionKeyRangeGone)))

	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsGone(err))

	_, ok := err.(*PartitionKeyRangeGoneError)
	require.True(t, ok)
}

func TestErrorFromResponse_Gone_PartitionSplitting(t *testing.T) {
	msg := NewResponseMessage(410, uuid.New())
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderSubStatus), TokenULong, uint32(SubStatusCompletingSplit)))

	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsGone(err))

	_, ok := err.(*PartitionKeyRangeIsSplittingError)
	require.True(t, ok)
}

func TestErrorFromResponse_Gone_PartitionMigrating(t *testing.T) {
	msg := NewResponseMessage(410, uuid.New())
	require.NoError(t, msg.Headers.SetValue(uint16(ResponseHeaderSubStatus), TokenULong, uint32(SubStatusCompletingPartitionMigrate)))

	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsGone(err))

	_, ok := err.(*PartitionIsMigratingError)
	require.True(t, ok)
}

func TestErrorFromResponse_TooManyRequests(t *testing.T) {
	msg := NewResponseMessage(429, uuid.New())
	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsRequestRateTooLarge(err))

	rateLimitErr, ok := err.(*RequestRateTooLargeError)
	require.True(t, ok)
	require.Equal(t, int32(429), rateLimitErr.StatusCode)
}

func TestErrorFromResponse_ServiceUnavailable(t *testing.T) {
	msg := NewResponseMessage(503, uuid.New())
	err := ErrorFromResponse(msg)
	require.Error(t, err)
	require.True(t, IsServiceUnavailable(err))
}

func TestError_IsRetriable(t *testing.T) {
	testCases := []struct {
		name      string
		status    int32
		retriable bool
	}{
		{"BadRequest", 400, false},
		{"Unauthorized", 401, false},
		{"Forbidden", 403, false},
		{"NotFound", 404, false},
		{"Conflict", 409, false},
		{"Gone", 410, true},
		{"PreconditionFailed", 412, false},
		{"TooManyRequests", 429, true},
		{"InternalServerError", 500, false},
		{"ServiceUnavailable", 503, true},
		{"RequestTimeout", 408, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg := NewResponseMessage(tc.status, uuid.New())
			err := ErrorFromResponse(msg)
			require.Error(t, err)
			require.Equal(t, tc.retriable, IsRetriable(err))
		})
	}
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func TestFrameConstants(t *testing.T) {
	// Verify frame size constants match expected values
	require.Equal(t, 24, RequestFrameLength)  // 4 + 2 + 2 + 16
	require.Equal(t, 24, ResponseFrameLength) // 4 + 4 + 16
}

func TestResponseFrame_HeadersLength(t *testing.T) {
	frame := &ResponseFrame{
		Length: ResponseFrameLength + 100,
	}
	require.Equal(t, 100, frame.HeadersLength())

	frame = &ResponseFrame{
		Length: ResponseFrameLength,
	}
	require.Equal(t, 0, frame.HeadersLength())
}

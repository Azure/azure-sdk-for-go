// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParseResponseMessage_BasicResponse(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-5678-1234-567812345678")
	msg := NewResponseMessage(200, activityID)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, activityID, resp.ActivityID)
	require.Equal(t, "https://test.cosmos.azure.com", resp.Endpoint)
	require.True(t, resp.IsSuccessStatusCode())
}

func TestParseResponseMessage_ErrorResponse(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(404, activityID)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.Equal(t, 404, resp.StatusCode)
	require.False(t, resp.IsSuccessStatusCode())
}

func TestParseResponseMessage_WithPayload(t *testing.T) {
	activityID := uuid.New()
	payload := []byte(`{"id":"doc1","content":"test"}`)

	msg := NewResponseMessage(200, activityID)
	msg.Payload = payload
	msg.Headers.SetByte(uint16(ResponseHeaderPayloadPresent), 1)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.Equal(t, payload, resp.Content)
}

func TestParseResponseMessage_RequestCharge(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)
	msg.Headers.SetDouble(uint16(ResponseHeaderRequestCharge), 3.14)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	charge := resp.GetRequestCharge()
	require.InDelta(t, 3.14, charge, 0.001)
}

func TestParseResponseMessage_SessionToken(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)
	msg.Headers.SetString(uint16(ResponseHeaderSessionToken), "0:1#1234#56=789")

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	sessionToken := resp.GetSessionTokenString()
	require.Equal(t, "0:1#1234#56=789", sessionToken)
}

func TestParseResponseMessage_LSN(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)
	msg.Headers.SetLongLong(uint16(ResponseHeaderLSN), 12345)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	lsn := resp.GetLSN()
	require.Equal(t, int64(12345), lsn)
}

func TestParseResponseMessage_ItemCount(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)
	msg.Headers.SetULong(uint16(ResponseHeaderItemCount), 42)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	count := resp.GetItemCount()
	require.Equal(t, 42, count)
}

func TestParseResponseMessage_Continuation(t *testing.T) {
	activityID := uuid.New()
	continuationToken := `{"token":"abc123","range":{"min":"","max":"FF"}}`

	msg := NewResponseMessage(200, activityID)
	msg.Headers.SetString(uint16(ResponseHeaderContinuationToken), continuationToken)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	continuation := resp.GetContinuation()
	require.Equal(t, continuationToken, continuation)
}

func TestParseResponseMessage_ETag(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)
	msg.Headers.SetString(uint16(ResponseHeaderETag), "\"etag-value-12345\"")

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	etag := resp.GetETag()
	require.Equal(t, "\"etag-value-12345\"", etag)
}

func TestParseResponseMessage_SubStatus(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(429, activityID)
	msg.Headers.SetULong(uint16(ResponseHeaderSubStatus), 3200)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	subStatus := resp.GetSubStatusCode()
	require.Equal(t, 3200, subStatus)
}

func TestParseResponseMessage_RetryAfter(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(429, activityID)
	msg.Headers.SetULong(uint16(ResponseHeaderRetryAfterMilliseconds), 1000)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	retryAfter := resp.GetRetryAfterMs()
	require.Equal(t, int64(1000), retryAfter)
}

func TestParseResponseMessage_AllHeaders(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)

	msg.Headers.SetDouble(uint16(ResponseHeaderRequestCharge), 2.5)
	msg.Headers.SetString(uint16(ResponseHeaderSessionToken), "session-token")
	msg.Headers.SetString(uint16(ResponseHeaderContinuationToken), "continuation")
	msg.Headers.SetString(uint16(ResponseHeaderETag), "etag")
	msg.Headers.SetLongLong(uint16(ResponseHeaderLSN), 100)
	msg.Headers.SetLongLong(uint16(ResponseHeaderGlobalCommittedLSN), 99)
	msg.Headers.SetULong(uint16(ResponseHeaderItemCount), 10)
	msg.Headers.SetString(uint16(ResponseHeaderPartitionKeyRangeId), "0")
	msg.Headers.SetString(uint16(ResponseHeaderOwnerId), "owner123")
	msg.Headers.SetString(uint16(ResponseHeaderOwnerFullName), "dbs/db/colls/coll")
	msg.Headers.SetByte(uint16(ResponseHeaderIsRUPerMinuteUsed), 1)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.InDelta(t, 2.5, resp.GetRequestCharge(), 0.001)
	require.Equal(t, "session-token", resp.GetHeader(RespHeaderSessionToken))
	require.Equal(t, "continuation", resp.GetHeader(RespHeaderContinuation))
	require.Equal(t, "etag", resp.GetHeader(RespHeaderETag))
	require.Equal(t, "100", resp.GetHeader(RespHeaderLSN))
	require.Equal(t, "99", resp.GetHeader(RespHeaderGlobalCommittedLSN))
	require.Equal(t, "10", resp.GetHeader(RespHeaderItemCount))
	require.Equal(t, "0", resp.GetHeader(RespHeaderPartitionKeyRangeID))
	require.Equal(t, "owner123", resp.GetHeader(RespHeaderOwnerID))
	require.Equal(t, "dbs/db/colls/coll", resp.GetHeader(RespHeaderOwnerFullName))
	require.Equal(t, "true", resp.GetHeader(RespHeaderIsRUPerMinuteUsed))
}

func TestParseResponseMessage_RoundTrip(t *testing.T) {
	activityID := uuid.MustParse("abcd1234-abcd-1234-abcd-1234abcd5678")
	payload := []byte(`{"documents":[{"id":"1"}],"_count":1}`)

	original := NewResponseMessage(200, activityID)
	original.Payload = payload
	original.Headers.SetByte(uint16(ResponseHeaderPayloadPresent), 1)
	original.Headers.SetDouble(uint16(ResponseHeaderRequestCharge), 5.5)
	original.Headers.SetString(uint16(ResponseHeaderSessionToken), "0:1#100#5=200")
	original.Headers.SetLongLong(uint16(ResponseHeaderLSN), 500)
	original.Headers.SetULong(uint16(ResponseHeaderItemCount), 1)

	encoded, err := EncodeResponseToBytes(original)
	require.NoError(t, err)

	decoded, err := DecodeResponseFromBytes(encoded)
	require.NoError(t, err)

	require.Equal(t, original.Frame.Status, decoded.Frame.Status)
	require.Equal(t, original.Frame.ActivityID, decoded.Frame.ActivityID)
	require.Equal(t, original.Payload, decoded.Payload)

	resp, err := ParseResponseMessage(decoded, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, activityID, resp.ActivityID)
	require.InDelta(t, 5.5, resp.GetRequestCharge(), 0.001)
	require.Equal(t, "0:1#100#5=200", resp.GetSessionTokenString())
	require.Equal(t, int64(500), resp.GetLSN())
	require.Equal(t, 1, resp.GetItemCount())
	require.Equal(t, payload, resp.Content)
}

func TestIsSuccessStatus(t *testing.T) {
	require.True(t, IsSuccessStatus(200))
	require.True(t, IsSuccessStatus(201))
	require.True(t, IsSuccessStatus(204))
	require.True(t, IsSuccessStatus(299))
	require.False(t, IsSuccessStatus(199))
	require.False(t, IsSuccessStatus(300))
	require.False(t, IsSuccessStatus(400))
	require.False(t, IsSuccessStatus(404))
	require.False(t, IsSuccessStatus(500))
}

func TestIsRetriableStatus(t *testing.T) {
	require.True(t, IsRetriableStatus(408, 0))
	require.True(t, IsRetriableStatus(503, 0))
	require.True(t, IsRetriableStatus(410, 0))
	require.True(t, IsRetriableStatus(449, 0))
	require.True(t, IsRetriableStatus(429, 0))
	require.True(t, IsRetriableStatus(500, 0))

	require.False(t, IsRetriableStatus(500, 1002))
	require.False(t, IsRetriableStatus(404, 0))
	require.False(t, IsRetriableStatus(400, 0))
	require.False(t, IsRetriableStatus(200, 0))
}

func TestGetResponseHeaderHelpers(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)

	msg.Headers.SetByte(uint16(ResponseHeaderPayloadPresent), 1)
	msg.Headers.SetULong(uint16(ResponseHeaderItemCount), 42)
	msg.Headers.SetLongLong(uint16(ResponseHeaderLSN), 12345)
	msg.Headers.SetDouble(uint16(ResponseHeaderRequestCharge), 3.14)
	msg.Headers.SetString(uint16(ResponseHeaderSessionToken), "token")

	require.Equal(t, byte(1), GetResponseHeaderByte(msg, ResponseHeaderPayloadPresent))
	require.Equal(t, uint32(42), GetResponseHeaderULong(msg, ResponseHeaderItemCount))
	require.Equal(t, int64(12345), GetResponseHeaderLongLong(msg, ResponseHeaderLSN))
	require.InDelta(t, 3.14, GetResponseHeaderDouble(msg, ResponseHeaderRequestCharge), 0.001)
	require.Equal(t, "token", GetResponseHeaderString(msg, ResponseHeaderSessionToken))
}

func TestConvertTokenValueToString(t *testing.T) {
	tests := []struct {
		name     string
		headerID ResponseHeader
		value    interface{}
		expected string
	}{
		{"byte_0", ResponseHeaderPayloadPresent, byte(0), "0"},
		{"byte_1", ResponseHeaderPayloadPresent, byte(1), "1"},
		{"bool_true", ResponseHeaderIsRUPerMinuteUsed, byte(1), "true"},
		{"bool_false", ResponseHeaderIsRUPerMinuteUsed, byte(0), "false"},
		{"int64", ResponseHeaderLSN, int64(12345), "12345"},
		{"float64", ResponseHeaderRequestCharge, float64(3.14), "3.14"},
		{"string", ResponseHeaderSessionToken, "my-token", "my-token"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := convertTokenValueToString(tc.headerID, tc.value)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestParseResponseMessage_LongLSNValues(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)

	largeLSN := int64(9223372036854775807)
	largeGLSN := int64(8223372036854775807)

	msg.Headers.SetLongLong(uint16(ResponseHeaderLSN), largeLSN)
	msg.Headers.SetLongLong(uint16(ResponseHeaderGlobalCommittedLSN), largeGLSN)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.Equal(t, "9223372036854775807", resp.Headers[RespHeaderLSN])
	require.Equal(t, "8223372036854775807", resp.Headers[RespHeaderGlobalCommittedLSN])
}

func TestParseResponseMessage_SubStatusMapping(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int32
		subStatus     int
		expectRetry   bool
		expectSuccess bool
	}{
		{"410_PartitionKeyRangeGone", 410, 1002, true, false},
		{"410_CompletingSplit", 410, 1007, true, false},
		{"410_CompletingPartitionMigration", 410, 1008, true, false},
		{"410_GenericGone", 410, 0, true, false},
		{"429_TooManyRequests", 429, 0, true, false},
		{"449_RetryWith", 449, 0, true, false},
		{"503_ServiceUnavailable", 503, 0, true, false},
		{"200_Success", 200, 0, false, true},
		{"404_NotFound", 404, 0, false, false},
		{"400_BadRequest", 400, 0, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			activityID := uuid.New()
			msg := NewResponseMessage(tc.statusCode, activityID)
			if tc.subStatus > 0 {
				msg.Headers.SetULong(uint16(ResponseHeaderSubStatus), uint32(tc.subStatus))
			}

			resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
			require.NoError(t, err)

			require.Equal(t, tc.expectSuccess, resp.IsSuccessStatusCode())
			require.Equal(t, tc.expectRetry, IsRetriableStatus(int32(resp.StatusCode), tc.subStatus))
		})
	}
}

func TestParseResponseMessage_GlobalCommittedLSN(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)

	msg.Headers.SetLongLong(uint16(ResponseHeaderLSN), 100)
	msg.Headers.SetLongLong(uint16(ResponseHeaderGlobalCommittedLSN), 95)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.Equal(t, "100", resp.Headers[RespHeaderLSN])
	require.Equal(t, "95", resp.Headers[RespHeaderGlobalCommittedLSN])
}

func TestParseResponseMessage_TransportRequestIDCorrelation(t *testing.T) {
	activityID := uuid.New()
	msg := NewResponseMessage(200, activityID)

	transportID := uint32(987654321)
	msg.Headers.SetULong(uint16(ResponseHeaderTransportRequestID), transportID)

	resp, err := ParseResponseMessage(msg, "https://test.cosmos.azure.com")
	require.NoError(t, err)

	require.Equal(t, "987654321", resp.Headers[RespHeaderTransportRequestID])
}

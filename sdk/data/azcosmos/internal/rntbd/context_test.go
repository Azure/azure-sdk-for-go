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

func TestContextRequest_NewContextRequest(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	userAgent := "azcosmos-go/1.0.0"

	req := NewContextRequest(activityID, userAgent)

	require.NotNil(t, req)
	require.Equal(t, activityID, req.ActivityID)
	require.NotNil(t, req.Headers)
	require.Equal(t, CurrentProtocolVersion, req.Headers.ProtocolVersion)
	require.Equal(t, ClientVersion, req.Headers.ClientVersion)
	require.Equal(t, userAgent, req.Headers.UserAgent)
}

func TestContextRequest_Encode_WireFormat(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	req := NewContextRequest(activityID, "test-agent")

	data, err := req.EncodeToBytes()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	length := binary.LittleEndian.Uint32(data[0:4])
	require.Equal(t, uint32(len(data)), length)

	resourceType := binary.LittleEndian.Uint16(data[4:6])
	require.Equal(t, uint16(ResourceConnection), resourceType)
	require.Equal(t, uint16(0x0000), resourceType)

	operationType := binary.LittleEndian.Uint16(data[6:8])
	require.Equal(t, uint16(OperationConnection), operationType)
	require.Equal(t, uint16(0x0000), operationType)

	resourceOperationCode := binary.LittleEndian.Uint32(data[4:8])
	require.Equal(t, uint32(0), resourceOperationCode)
}

func TestContextRequest_RoundTrip(t *testing.T) {
	testCases := []struct {
		name       string
		activityID uuid.UUID
		userAgent  string
	}{
		{
			name:       "basic",
			activityID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			userAgent:  "azcosmos-go/1.0.0",
		},
		{
			name:       "long user agent",
			activityID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			userAgent:  "azcosmos-go/1.0.0 (Linux; x86_64) azure-cosmos-dotnet-sdk/3.0.0",
		},
		{
			name:       "nil UUID",
			activityID: uuid.Nil,
			userAgent:  "test",
		},
		{
			name:       "random UUID",
			activityID: uuid.New(),
			userAgent:  "agent",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			original := NewContextRequest(tc.activityID, tc.userAgent)

			data, err := original.EncodeToBytes()
			require.NoError(t, err)

			decoded, err := DecodeContextRequest(bytes.NewReader(data))
			require.NoError(t, err)

			require.Equal(t, original.ActivityID, decoded.ActivityID)
			require.Equal(t, original.Headers.ProtocolVersion, decoded.Headers.ProtocolVersion)
			require.Equal(t, original.Headers.ClientVersion, decoded.Headers.ClientVersion)
			require.Equal(t, original.Headers.UserAgent, decoded.Headers.UserAgent)
		})
	}
}

func TestContextRequest_ComputeLength(t *testing.T) {
	activityID := uuid.New()
	req := NewContextRequest(activityID, "test-agent")

	computedLength := req.ComputeLength()
	data, err := req.EncodeToBytes()
	require.NoError(t, err)
	require.Equal(t, computedLength, len(data))
}

func TestDecodeContextRequest_InvalidResourceOperationType(t *testing.T) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, uint32(RequestFrameLength))
	binary.Write(&buf, binary.LittleEndian, uint16(ResourceDocument))
	binary.Write(&buf, binary.LittleEndian, uint16(OperationRead))
	WriteUUID(uuid.New(), &buf)

	_, err := DecodeContextRequest(bytes.NewReader(buf.Bytes()))
	require.Error(t, err)
	require.Contains(t, err.Error(), "resourceOperationType=0")
}

func TestDecodeContextRequest_TooSmall(t *testing.T) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, uint32(10))
	binary.Write(&buf, binary.LittleEndian, uint32(0))
	WriteUUID(uuid.Nil, &buf)

	_, err := DecodeContextRequest(bytes.NewReader(buf.Bytes()))
	require.Error(t, err)
	require.Contains(t, err.Error(), "too small")
}

func TestContext_NewContext(t *testing.T) {
	req := NewContextRequest(uuid.New(), "test-agent")
	ctx := ContextFrom(req, "cosmos-server", "1.0.0", 200)

	require.NotNil(t, ctx)
	require.Equal(t, req.ActivityID, ctx.ActivityID)
	require.Equal(t, int32(200), ctx.Status)
	require.Equal(t, CurrentProtocolVersion, ctx.ProtocolVersion)
	require.Equal(t, req.Headers.ClientVersion, ctx.ClientVersion)
	require.Equal(t, "cosmos-server", ctx.ServerAgent)
	require.Equal(t, "1.0.0", ctx.ServerVersion)
}

func TestContext_IsSuccess(t *testing.T) {
	testCases := []struct {
		status    int32
		isSuccess bool
	}{
		{200, true},
		{201, true},
		{204, true},
		{299, true},
		{300, false},
		{399, false},
		{400, false},
		{401, false},
		{404, false},
		{500, false},
		{503, false},
	}

	for _, tc := range testCases {
		ctx := &Context{Status: tc.status}
		require.Equal(t, tc.isSuccess, ctx.IsSuccess())
	}
}

func TestContext_ServerProperties(t *testing.T) {
	ctx := &Context{
		ServerAgent:   "CosmosDB",
		ServerVersion: "2.0.0",
	}
	require.Equal(t, "CosmosDB/2.0.0", ctx.ServerProperties())
}

func TestContext_RoundTrip(t *testing.T) {
	testCases := []struct {
		name                   string
		status                 int32
		serverAgent            string
		serverVersion          string
		idleTimeout            uint32
		unauthenticatedTimeout uint32
	}{
		{
			name:                   "success",
			status:                 200,
			serverAgent:            "CosmosDB",
			serverVersion:          "1.0.0",
			idleTimeout:            300,
			unauthenticatedTimeout: 30,
		},
		{
			name:                   "success with zero timeouts",
			status:                 200,
			serverAgent:            "TestServer",
			serverVersion:          "2.0.0",
			idleTimeout:            0,
			unauthenticatedTimeout: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			activityID := uuid.New()
			original := &Context{
				ActivityID:                      activityID,
				Status:                          tc.status,
				ProtocolVersion:                 CurrentProtocolVersion,
				ClientVersion:                   ClientVersion,
				ServerAgent:                     tc.serverAgent,
				ServerVersion:                   tc.serverVersion,
				IdleTimeoutInSeconds:            tc.idleTimeout,
				UnauthenticatedTimeoutInSeconds: tc.unauthenticatedTimeout,
			}

			data, err := original.EncodeToBytes()
			require.NoError(t, err)

			decoded, err := DecodeContext(bytes.NewReader(data))
			require.NoError(t, err)

			require.Equal(t, original.ActivityID, decoded.ActivityID)
			require.Equal(t, original.Status, decoded.Status)
			require.Equal(t, original.ProtocolVersion, decoded.ProtocolVersion)
			require.Equal(t, original.ClientVersion, decoded.ClientVersion)
			require.Equal(t, original.ServerAgent, decoded.ServerAgent)
			require.Equal(t, original.ServerVersion, decoded.ServerVersion)
			require.Equal(t, original.IdleTimeoutInSeconds, decoded.IdleTimeoutInSeconds)
			require.Equal(t, original.UnauthenticatedTimeoutInSeconds, decoded.UnauthenticatedTimeoutInSeconds)
		})
	}
}

func TestDecodeContext_ErrorStatus_ReturnsContextException(t *testing.T) {
	testCases := []struct {
		name   string
		status int32
	}{
		{"bad request", 400},
		{"unauthorized", 401},
		{"forbidden", 403},
		{"not found", 404},
		{"internal error", 500},
		{"service unavailable", 503},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			activityID := uuid.New()
			ctx := &Context{
				ActivityID:    activityID,
				Status:        tc.status,
				ServerAgent:   "TestServer",
				ServerVersion: "1.0.0",
			}

			data, err := ctx.EncodeToBytes()
			require.NoError(t, err)

			_, err = DecodeContext(bytes.NewReader(data))
			require.Error(t, err)
			require.True(t, IsContextException(err))

			ctxErr, ok := err.(*ContextException)
			require.True(t, ok)
			require.Equal(t, tc.status, ctxErr.Status)
			require.Equal(t, activityID, ctxErr.ActivityID)
			require.Equal(t, "TestServer", ctxErr.ServerAgent)
			require.Equal(t, "1.0.0", ctxErr.ServerVersion)
		})
	}
}

func TestContextException_Error(t *testing.T) {
	activityID := uuid.MustParse("12345678-1234-1234-1234-123456789abc")
	err := &ContextException{
		Status:     401,
		ActivityID: activityID,
	}

	errMsg := err.Error()
	require.Contains(t, errMsg, "context negotiation failed")
	require.Contains(t, errMsg, "status=401")
	require.Contains(t, errMsg, activityID.String())
}

func TestContextException_ErrorWithDetails(t *testing.T) {
	activityID := uuid.New()
	err := &ContextException{
		Status:     400,
		ActivityID: activityID,
		Details: map[string]interface{}{
			"code":    "InvalidProtocol",
			"message": "Protocol version not supported",
		},
	}

	errMsg := err.Error()
	require.Contains(t, errMsg, "context negotiation failed")
	require.Contains(t, errMsg, "InvalidProtocol")
}

func TestContextException_ResponseHeaders(t *testing.T) {
	err := &ContextException{
		Status:          400,
		ClientVersion:   "1.0.0",
		ProtocolVersion: 1,
		ServerAgent:     "CosmosDB",
		ServerVersion:   "2.0.0",
	}

	headers := err.ResponseHeaders()
	require.Equal(t, "1.0.0", headers["requiredClientVersion"])
	require.Equal(t, uint32(1), headers["requiredProtocolVersion"])
	require.Equal(t, "CosmosDB", headers["serverAgent"])
	require.Equal(t, "2.0.0", headers["serverVersion"])
}

func TestIsContextException(t *testing.T) {
	ctxErr := &ContextException{Status: 400}
	require.True(t, IsContextException(ctxErr))

	otherErr := &BadRequestError{RntbdError: RntbdError{StatusCode: 400}}
	require.False(t, IsContextException(otherErr))

	require.False(t, IsContextException(nil))
}

func TestContextFrom(t *testing.T) {
	activityID := uuid.New()
	req := NewContextRequest(activityID, "test-agent")

	ctx := ContextFrom(req, "server", "1.0", 200)

	require.Equal(t, activityID, ctx.ActivityID)
	require.Equal(t, int32(200), ctx.Status)
	require.Equal(t, CurrentProtocolVersion, ctx.ProtocolVersion)
	require.Equal(t, ClientVersion, ctx.ClientVersion)
	require.Equal(t, "server", ctx.ServerAgent)
	require.Equal(t, "1.0", ctx.ServerVersion)
	require.Equal(t, uint32(0), ctx.IdleTimeoutInSeconds)
	require.Equal(t, uint32(0), ctx.UnauthenticatedTimeoutInSeconds)
}

func TestContextRequest_HeaderTokenTypes(t *testing.T) {
	req := NewContextRequest(uuid.New(), "test")
	ts := req.Headers.toTokenStream()

	protocolToken := ts.Get(uint16(ContextRequestHeaderProtocolVersion))
	require.NotNil(t, protocolToken)
	require.True(t, protocolToken.IsPresent())

	clientVersionToken := ts.Get(uint16(ContextRequestHeaderClientVersion))
	require.NotNil(t, clientVersionToken)
	require.True(t, clientVersionToken.IsPresent())

	userAgentToken := ts.Get(uint16(ContextRequestHeaderUserAgent))
	require.NotNil(t, userAgentToken)
	require.True(t, userAgentToken.IsPresent())
}

func TestContext_HeaderTokenTypes(t *testing.T) {
	ctx := &Context{
		ProtocolVersion:                 1,
		ClientVersion:                   "1.0",
		ServerAgent:                     "server",
		ServerVersion:                   "2.0",
		IdleTimeoutInSeconds:            300,
		UnauthenticatedTimeoutInSeconds: 30,
	}

	headers := &contextHeaders{
		ProtocolVersion:                 ctx.ProtocolVersion,
		ClientVersion:                   ctx.ClientVersion,
		ServerAgent:                     ctx.ServerAgent,
		ServerVersion:                   ctx.ServerVersion,
		IdleTimeoutInSeconds:            ctx.IdleTimeoutInSeconds,
		UnauthenticatedTimeoutInSeconds: ctx.UnauthenticatedTimeoutInSeconds,
	}

	ts := headers.toTokenStream()

	protocolToken := ts.Get(uint16(ContextHeaderProtocolVersion))
	require.NotNil(t, protocolToken)
	require.True(t, protocolToken.IsPresent())

	serverAgentToken := ts.Get(uint16(ContextHeaderServerAgent))
	require.NotNil(t, serverAgentToken)
	require.True(t, serverAgentToken.IsPresent())

	serverVersionToken := ts.Get(uint16(ContextHeaderServerVersion))
	require.NotNil(t, serverVersionToken)
	require.True(t, serverVersionToken.IsPresent())

	idleToken := ts.Get(uint16(ContextHeaderIdleTimeoutInSeconds))
	require.NotNil(t, idleToken)
	require.True(t, idleToken.IsPresent())
}

func TestContextRequest_ConnectionTypeValues(t *testing.T) {
	require.Equal(t, ResourceType(0x0000), ResourceConnection)
	require.Equal(t, OperationType(0x0000), OperationConnection)
}

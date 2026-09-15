// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestServerTimeoutMillis(t *testing.T) {
	tests := []struct {
		name      string
		ctx       func() (context.Context, context.CancelFunc)
		wantMin   uint
		wantMax   uint
		wantExact *uint // when non-nil, assert exact value instead of range
	}{
		{
			name: "no deadline returns default 60s",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			wantExact: ptrUint(60000),
		},
		{
			name: "deadline shorter than 60s drops the buffer second",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			// serverTimeoutBuffer comes off the remaining time so the broker's
			// timeout fires first. The upper bound is below the 10s deadline on
			// purpose: sending the full remaining time would land near 10000 and
			// fail here. The lower bound tolerates scheduler drift.
			wantMin: 8500,
			wantMax: 9000,
		},
		{
			name: "deadline equal to the buffer returns 0",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 1*time.Second)
			},
			// Remaining is already under the buffer, so the caller is about to give
			// up regardless and the value clamps rather than underflowing.
			wantExact: ptrUint(0),
		},
		{
			name: "deadline just above the buffer keeps what is left",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 1500*time.Millisecond)
			},
			wantMin: 300,
			wantMax: 500,
		},
		{
			name: "deadline longer than 60s respects user timeout",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 120*time.Second)
			},
			wantMin: 118500,
			wantMax: 119000,
		},
		{
			name: "deadline exactly 60s drops the buffer second",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 60*time.Second)
			},
			// Proves the deadline path is taken rather than defaultServerTimeout:
			// the default would give exactly 60000, which is outside this range.
			wantMin: 58500,
			wantMax: 59000,
		},
		{
			name: "expired context returns 0",
			ctx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
				return ctx, cancel
			},
			wantExact: ptrUint(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctx()
			defer cancel()

			got := serverTimeoutMillis(ctx)

			if tt.wantExact != nil {
				require.Equal(t, *tt.wantExact, got)
			} else {
				require.GreaterOrEqual(t, got, tt.wantMin, "timeout should be >= %d, got %d", tt.wantMin, got)
				require.LessOrEqual(t, got, tt.wantMax, "timeout should be <= %d, got %d", tt.wantMax, got)
			}
		})
	}
}

func ptrUint(v uint) *uint {
	return &v
}

// capturingRPCLink records the message handed to RPC so a test can assert on what
// the management operation built, and answers 204 (no messages available) so the
// caller returns without needing a payload.
type capturingRPCLink struct {
	Sent     *amqp.Message
	Response *amqpwrap.RPCResponse
	Err      error
}

func (l *capturingRPCLink) Close(ctx context.Context) error { return nil }

func (l *capturingRPCLink) RPC(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
	l.Sent = msg
	if l.Err != nil {
		return nil, l.Err
	}
	if l.Response != nil {
		return l.Response, nil
	}
	return &amqpwrap.RPCResponse{Code: 204}, nil
}

func TestBatchDeleteMessages(t *testing.T) {
	cutoff := time.Date(2026, time.August, 21, 12, 30, 0, 0, time.UTC)
	link := &capturingRPCLink{
		Response: &amqpwrap.RPCResponse{
			Code: 200,
			Message: &amqp.Message{Value: map[string]any{
				"message-count": int32(37),
			}},
		},
	}

	deletedCount, err := BatchDeleteMessages(context.Background(), link, "receiver-link", 50, cutoff, ptrString("session-1"))
	require.NoError(t, err)
	require.EqualValues(t, 37, deletedCount)
	require.Equal(t, "com.microsoft:batch-delete-messages", link.Sent.ApplicationProperties["operation"])
	require.Equal(t, "receiver-link", link.Sent.ApplicationProperties["associated-link-name"])

	body := link.Sent.Value.(map[string]any)
	require.Equal(t, int32(50), body["message-count"])
	require.Equal(t, cutoff, body["enqueued-time-utc"])
	require.Equal(t, "session-1", body["session-id"])
}

func TestBatchDeleteMessagesNoMessages(t *testing.T) {
	link := &capturingRPCLink{}

	deletedCount, err := BatchDeleteMessages(context.Background(), link, "receiver-link", 1, time.Now().UTC(), nil)
	require.NoError(t, err)
	require.Zero(t, deletedCount)

	body := link.Sent.Value.(map[string]any)
	require.NotContains(t, body, "session-id")
}

func TestBatchDeleteMessagesMessageNotFound(t *testing.T) {
	response := &amqpwrap.RPCResponse{
		Code: 404,
		Message: &amqp.Message{
			ApplicationProperties: map[string]any{"error-condition": "com.microsoft:message-not-found"},
			Value:                 map[string]any{"message-count": int32(2)},
		},
	}
	link := &capturingRPCLink{Err: RPCError{Resp: response, Message: "message not found"}}

	deletedCount, err := BatchDeleteMessages(context.Background(), link, "receiver-link", 10, time.Now().UTC(), nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, deletedCount)
}

func TestBatchDeleteMessagesRejectsMessageNotFoundWithoutValidCount(t *testing.T) {
	for _, value := range []any{nil, map[string]any{}, map[string]any{"message-count": "2"}} {
		response := &amqpwrap.RPCResponse{
			Code: 404,
			Message: &amqp.Message{
				ApplicationProperties: map[string]any{"error-condition": "com.microsoft:message-not-found"},
				Value:                 value,
			},
		}
		link := &capturingRPCLink{Err: RPCError{Resp: response, Message: "message not found"}}

		_, err := BatchDeleteMessages(context.Background(), link, "receiver-link", 10, time.Now().UTC(), nil)
		require.Error(t, err)
	}
}

func TestBatchDeleteMessagesRejectsInvalidResponseCount(t *testing.T) {
	for _, value := range []any{int32(-1), int32(11), float64(1)} {
		link := &capturingRPCLink{Response: &amqpwrap.RPCResponse{
			Code:    200,
			Message: &amqp.Message{Value: map[string]any{"message-count": value}},
		}}

		_, err := BatchDeleteMessages(context.Background(), link, "receiver-link", 10, time.Now().UTC(), nil)
		require.Error(t, err)
	}
}

func TestBatchDeleteMessagesRejectsUnexpectedStatus(t *testing.T) {
	link := &capturingRPCLink{Response: &amqpwrap.RPCResponse{
		Code:    201,
		Message: &amqp.Message{Value: map[string]any{"message-count": int32(1)}},
	}}

	_, err := BatchDeleteMessages(context.Background(), link, "receiver-link", 10, time.Now().UTC(), nil)
	require.Error(t, err)
}

func ptrString(value string) *string {
	return &value
}

// TestPeekMessagesDefersServerTimeoutToRPC pins the two conditions PeekMessages has
// to meet for rpcLink.RPC to set its server-timeout. PeekMessages no longer sets the
// value itself, so that RPC can compute it immediately before the send instead of
// while the message is still being built. That handoff is silent: if PeekMessages
// ever renamed its operation or preset a timeout key, the message would simply go
// out with no bound and every existing test would still pass.
//
// The other direction, that RPC does set the value under these conditions, is
// covered by TestRPCLinkServerTimeoutScoping.
func TestPeekMessagesDefersServerTimeoutToRPC(t *testing.T) {
	link := &capturingRPCLink{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := PeekMessages(ctx, link, "link-name", 0, 1)
	require.NoError(t, err)
	require.NotNil(t, link.Sent)

	op, ok := link.Sent.ApplicationProperties["operation"].(string)
	require.True(t, ok, "operation must be set")
	require.True(t, strings.HasPrefix(op, "com.microsoft:"),
		"RPC only sets server-timeout for com.microsoft: operations, got %q", op)

	require.NotContains(t, link.Sent.ApplicationProperties, "server-timeout",
		"PeekMessages must leave server-timeout to RPC")
	require.NotContains(t, link.Sent.ApplicationProperties, "com.microsoft:server-timeout",
		"a preset vendor-prefixed key would stop RPC from setting server-timeout")
}

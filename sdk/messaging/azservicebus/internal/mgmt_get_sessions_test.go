// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetMessageSessions_UpdatedAfterMode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
			// Verify request body keys
			body, ok := msg.Value.(map[string]any)
			require.True(t, ok)
			require.Contains(t, body, "last-updated-time")
			require.Contains(t, body, "skip")
			require.Contains(t, body, "top")

			// Verify a real timestamp is sent (not the far-future sentinel)
			ts, ok := body["last-updated-time"].(time.Time)
			require.True(t, ok)
			require.Equal(t, 2026, ts.Year())

			// Verify operation name
			require.Equal(t, "com.microsoft:get-message-sessions", msg.ApplicationProperties["operation"])

			// Verify no associated-link-name (entity-level operation)
			_, hasLinkName := msg.ApplicationProperties["associated-link-name"]
			require.False(t, hasLinkName)

			return &amqpwrap.RPCResponse{
				Code: 200,
				Message: &amqp.Message{
					Value: map[string]any{
						"sessions-ids": []any{"session-a", "session-b", "session-c"},
					},
				},
			}, nil
		})

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC), 0, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"session-a", "session-b", "session-c"}, result)
}

func TestGetMessageSessions_ActiveSessionsMode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// The sentinel for "active messages" mode: 9999-12-31T23:59:59.999999999.
	// On the AMQP wire this becomes 253402300799999 ms (DateTime.MaxValue at ms precision).
	sentinel := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
			body := msg.Value.(map[string]any)
			ts := body["last-updated-time"].(time.Time)
			require.Equal(t, 9999, ts.Year())
			require.Equal(t, time.December, ts.Month())
			require.Equal(t, 31, ts.Day())

			return &amqpwrap.RPCResponse{
				Code: 200,
				Message: &amqp.Message{
					Value: map[string]any{
						"sessions-ids": []any{"active-1", "active-2"},
					},
				},
			}, nil
		})

	result, err := GetMessageSessions(context.Background(), rpcLink, sentinel, 0, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"active-1", "active-2"}, result)
}

func TestGetMessageSessions_Returns204NoContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).Return(
		&amqpwrap.RPCResponse{
			Code:    204,
			Message: &amqp.Message{},
		}, nil)

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC), 0, 100)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestGetMessageSessions_ReturnsEmptySessionsIds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).Return(
		&amqpwrap.RPCResponse{
			Code: 200,
			Message: &amqp.Message{
				Value: map[string]any{
					"sessions-ids": nil,
				},
			},
		}, nil)

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC), 0, 100)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestGetMessageSessions_HandlesTypedStringSlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).Return(
		&amqpwrap.RPCResponse{
			Code: 200,
			Message: &amqp.Message{
				Value: map[string]any{
					"sessions-ids": []string{"typed-1", "typed-2"},
				},
			},
		}, nil)

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), 0, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"typed-1", "typed-2"}, result)
}

func TestGetMessageSessions_404SessionNotFoundReturnsEmpty(t *testing.T) {
	// 404 + a "session not found" description means the entity exists but has no
	// sessions. This is the only 404 case GetMessageSessions silently treats as
	// empty; all other 404s (covered by TestGetMessageSessions_404OtherSurfacesError)
	// must surface as an error so callers can distinguish "no sessions" from
	// "entity not found".
	descriptions := []string{
		"Session not found",
		"SessionNotFound",
		"session not found.",
	}

	for _, description := range descriptions {
		t.Run(description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rpcLink := mock.NewMockRPCLink(ctrl)
			rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).Return(
				&amqpwrap.RPCResponse{
					Code:        404,
					Description: description,
					Message:     &amqp.Message{},
				}, nil)

			result, err := GetMessageSessions(context.Background(), rpcLink,
				time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC), 0, 100)

			require.NoError(t, err)
			require.Nil(t, result)
		})
	}
}

func TestGetMessageSessions_404OtherSurfacesError(t *testing.T) {
	// 404s without a "session not found" description (e.g. entity not found) must
	// be surfaced as errors so callers can distinguish a missing entity from an
	// entity with no sessions.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).Return(
		&amqpwrap.RPCResponse{
			Code:        404,
			Description: "The messaging entity 'sb://ns/q' could not be found.",
			Message:     &amqp.Message{},
		}, nil)

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC), 0, 100)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "404")
}

func TestGetMessageSessions_NonStringSessionIDSurfacesError(t *testing.T) {
	// The wire format declares sessions-ids as a string array. If the service
	// returns a non-string element (e.g. nil or a number) we must fail fast
	// rather than coerce via fmt.Sprintf("%v") and produce surprising IDs like
	// "<nil>".
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).Return(
		&amqpwrap.RPCResponse{
			Code: 200,
			Message: &amqp.Message{
				Value: map[string]any{
					"sessions-ids": []any{"valid-session", nil},
				},
			},
		}, nil)

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), 0, 100)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetMessageSessions_SkipAndTop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rpcLink := mock.NewMockRPCLink(ctrl)
	rpcLink.EXPECT().RPC(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
			body := msg.Value.(map[string]any)
			require.Equal(t, int32(200), body["skip"])
			require.Equal(t, int32(50), body["top"])

			return &amqpwrap.RPCResponse{
				Code: 200,
				Message: &amqp.Message{
					Value: map[string]any{
						"sessions-ids": []any{"page3-1"},
					},
				},
			}, nil
		})

	result, err := GetMessageSessions(context.Background(), rpcLink,
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), 200, 50)

	require.NoError(t, err)
	require.Equal(t, []string{"page3-1"}, result)
}

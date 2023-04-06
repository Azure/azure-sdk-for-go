// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNegotiateClaimWithCloseTimeout(t *testing.T) {
	for _, errToReturn := range []error{context.Canceled, context.DeadlineExceeded} {
		t.Run(fmt.Sprintf("Close() cancels with error %v", errToReturn), func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tp := mock.NewMockTokenProvider(ctrl)
			receiver := mock.NewMockAMQPReceiverCloser(ctrl)
			sender := mock.NewMockAMQPSenderCloser(ctrl)
			session := mock.NewMockAMQPSession(ctrl)
			client := mock.NewMockAMQPClient(ctrl)

			client.EXPECT().NewSession(mock.NotCancelled, gomock.Any()).Return(session, nil)
			session.EXPECT().NewReceiver(mock.NotCancelled, gomock.Any(), gomock.Any()).Return(receiver, nil)
			session.EXPECT().NewSender(mock.NotCancelled, gomock.Any(), gomock.Any()).Return(sender, nil)
			tp.EXPECT().GetToken(gomock.Any()).Return(&auth.Token{}, nil)

			mock.SetupRPC(sender, receiver, 1, func(sent, response *amqp.Message) {
				response.ApplicationProperties = map[string]any{
					"status-code": int32(200),
				}
			})

			// the context passed to these calls are already cancelled since the parent
			// context was cancelled. This basically just falls through the error handling
			// but it's okay - each resource should close any local state they can before
			// returning and we're going to end up abandoning ship on the connection.
			session.EXPECT().Close(mock.CancelledAndHasTimeout)
			sender.EXPECT().Close(mock.CancelledAndHasTimeout)

			// When links fail to close in a timely manner it's either because the connection is (somehow)
			// no longer valid _or_ conditions are preventing us from closing the link. In either case we
			// have to be careful since it means that some resources (for instance, singleton links like $cbs)
			// might "leak" since they can't be closed.
			//
			// Rather than attempt to do some complicated piecemeal recovery, we instead invalidate the entire
			// connection, which is the only safe way to ensure the client and service agree on what is open and
			// active.
			receiver.EXPECT().Close(mock.NotCancelledAndHasTimeout).DoAndReturn(func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			})

			err := NegotiateClaim(context.Background(), "audience", client, tp, mock.NewContextWithTimeoutForTests)
			require.EqualError(t, err, "connection must be reset, link/connection state may be inconsistent")
			require.Equal(t, GetRecoveryKind(err), RecoveryKindConn)
		})
	}
}

func TestNegotiateClaimWithAuthFailure(t *testing.T) {
	ctrl := gomock.NewController(t)

	tp := mock.NewMockTokenProvider(ctrl)
	receiver := mock.NewMockAMQPReceiverCloser(ctrl)
	sender := mock.NewMockAMQPSenderCloser(ctrl)
	session := mock.NewMockAMQPSession(ctrl)
	client := mock.NewMockAMQPClient(ctrl)

	client.EXPECT().NewSession(mock.NotCancelled, gomock.Any()).Return(session, nil)
	session.EXPECT().NewReceiver(mock.NotCancelled, gomock.Any(), gomock.Any()).Return(receiver, nil)
	session.EXPECT().NewSender(mock.NotCancelled, gomock.Any(), gomock.Any()).Return(sender, nil)
	tp.EXPECT().GetToken(gomock.Any()).Return(&auth.Token{}, nil)

	session.EXPECT().Close(mock.NotCancelledAndHasTimeout)
	sender.EXPECT().Close(mock.NotCancelledAndHasTimeout)
	receiver.EXPECT().Close(mock.NotCancelledAndHasTimeout)

	mock.SetupRPC(sender, receiver, 1, func(sent, response *amqp.Message) {
		// this is the kind of error you get if your connection string is inconsistent
		// (ie, you tamper with the shared key, etc..)
		response.ApplicationProperties = map[string]any{
			"status-code":        int32(401),
			"status-description": "InvalidSignature: The token has an invalid signature.",
			"error-condition":    "com.microsoft:auth-failed",
		}
	})

	err := NegotiateClaim(context.Background(), "audience", client, tp, mock.NewContextWithTimeoutForTests)

	require.EqualError(t, err, "rpc: failed, status code 401 and description: InvalidSignature: The token has an invalid signature.")
	require.Equal(t, GetRecoveryKind(err), RecoveryKindFatal)
}

func TestNegotiateClaimSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	tp := mock.NewMockTokenProvider(ctrl)
	receiver := mock.NewMockAMQPReceiverCloser(ctrl)
	sender := mock.NewMockAMQPSenderCloser(ctrl)
	session := mock.NewMockAMQPSession(ctrl)
	client := mock.NewMockAMQPClient(ctrl)

	client.EXPECT().NewSession(mock.NotCancelled, gomock.Any()).Return(session, nil)
	session.EXPECT().NewReceiver(mock.NotCancelled, gomock.Any(), gomock.Any()).Return(receiver, nil)
	session.EXPECT().NewSender(mock.NotCancelled, gomock.Any(), gomock.Any()).Return(sender, nil)
	tp.EXPECT().GetToken(gomock.Any()).Return(&auth.Token{}, nil)

	session.EXPECT().Close(mock.NotCancelledAndHasTimeout)
	sender.EXPECT().Close(mock.NotCancelledAndHasTimeout)
	receiver.EXPECT().Close(mock.NotCancelledAndHasTimeout)

	mock.SetupRPC(sender, receiver, 1, func(sent, response *amqp.Message) {
		response.ApplicationProperties = map[string]any{
			"status-code": int32(200),
		}
	})

	err := NegotiateClaim(context.Background(), "audience", client, tp, mock.NewContextWithTimeoutForTests)
	require.NoError(t, err)
}

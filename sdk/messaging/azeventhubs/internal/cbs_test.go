// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/Azure/go-amqp"
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

			client.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(session, nil)
			session.EXPECT().NewReceiver(test.NotCancelled, gomock.Any(), gomock.Any(), gomock.Any()).Return(receiver, nil)
			session.EXPECT().NewSender(test.NotCancelled, gomock.Any(), gomock.Any(), gomock.Any()).Return(sender, nil)
			tp.EXPECT().GetToken(gomock.Any()).Return(&auth.Token{}, nil)

			mock.SetupRPC(sender, receiver, 1, func(sent, response *amqp.Message) {
				response.ApplicationProperties = map[string]any{
					"status-code": int32(200),
				}
			})

			callerCtx, cancelCallerCtx := context.WithCancel(context.Background())
			defer cancelCallerCtx()

			// the context passed to these calls are already cancelled since the parent
			// context was cancelled. This basically just falls through the error handling
			// but it's okay - each resource should close any local state they can before
			// returning and we're going to end up abandoning ship on the connection.
			session.EXPECT().Close(test.NotCancelled).DoAndReturn(func(ctx context.Context) error {
				cancelCallerCtx()
				<-ctx.Done()
				return errToReturn
			})

			err := NegotiateClaim(callerCtx, "audience", client, tp)
			require.ErrorIs(t, err, errToReturn)
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

	client.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(session, nil)

	sender.EXPECT().LinkName().Return("sender-link-name")

	session.EXPECT().NewReceiver(test.NotCancelled, gomock.Any(), gomock.Any(), gomock.Any()).Return(receiver, nil)
	session.EXPECT().NewSender(test.NotCancelled, gomock.Any(), gomock.Any(), gomock.Any()).Return(sender, nil)
	session.EXPECT().Close(test.NotCancelled)
	session.EXPECT().ConnID().Return(uint64(101))

	tp.EXPECT().GetToken(gomock.Any()).Return(&auth.Token{}, nil)

	mock.SetupRPC(sender, receiver, 1, func(sent, response *amqp.Message) {
		// this is the kind of error you get if your connection string is inconsistent
		// (ie, you tamper with the shared key, etc..)
		response.ApplicationProperties = map[string]any{
			"status-code":        int32(401),
			"status-description": "InvalidSignature: The token has an invalid signature.",
			"error-condition":    "com.microsoft:auth-failed",
		}
	})

	err := NegotiateClaim(context.Background(), "audience", client, tp)

	require.EqualError(t, err, "rpc: failed, status code 401 and description: InvalidSignature: The token has an invalid signature.")
	require.Equal(t, GetRecoveryKind(err), RecoveryKindFatal)

	var amqpwrapErr amqpwrap.Error
	require.ErrorAs(t, err, &amqpwrapErr)
	require.Equal(t, uint64(101), amqpwrapErr.ConnID)
	require.Equal(t, "sender-link-name", amqpwrapErr.LinkName)
	require.Empty(t, amqpwrapErr.PartitionID)
}

func TestNegotiateClaimSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	tp := mock.NewMockTokenProvider(ctrl)
	receiver := mock.NewMockAMQPReceiverCloser(ctrl)
	sender := mock.NewMockAMQPSenderCloser(ctrl)
	session := mock.NewMockAMQPSession(ctrl)
	client := mock.NewMockAMQPClient(ctrl)

	client.EXPECT().NewSession(test.NotCancelled, gomock.Any()).Return(session, nil)
	session.EXPECT().NewReceiver(test.NotCancelled, gomock.Any(), gomock.Any(), gomock.Any()).Return(receiver, nil)
	session.EXPECT().NewSender(test.NotCancelled, gomock.Any(), gomock.Any(), gomock.Any()).Return(sender, nil)
	tp.EXPECT().GetToken(gomock.Any()).Return(&auth.Token{}, nil)

	session.EXPECT().Close(test.NotCancelled)

	mock.SetupRPC(sender, receiver, 1, func(sent, response *amqp.Message) {
		response.ApplicationProperties = map[string]any{
			"status-code": int32(200),
		}
	})

	err := NegotiateClaim(context.Background(), "audience", client, tp)
	require.NoError(t, err)
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestRPCLink(t *testing.T) {
	initFn := func() *fakeAMQPClient {
		return &fakeAMQPClient{
			session: &FakeAMQPSession{
				NS: &FakeNSForPartClient{
					Receiver: &FakeAMQPReceiver{},
					Sender:   &FakeAMQPSender{},
				},
			},
		}
	}

	t.Run("everything works, RPCLink is created", func(t *testing.T) {
		fakeClient := initFn()

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})
		require.NoError(t, err)
		require.NotNil(t, rpcLink)

		require.Zero(t, fakeClient.session.CloseCalled)
		require.Zero(t, fakeClient.session.NS.Receiver.CloseCalled)
		require.Zero(t, fakeClient.session.NS.Sender.CloseCalled)
	})

	t.Run("session created, sender fails", func(t *testing.T) {
		fakeClient := initFn()

		fakeClient.session.NS.NewSenderErr = errors.New("test error")

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})
		require.EqualError(t, err, "test error")
		require.Nil(t, rpcLink)

		require.Equal(t, 1, fakeClient.session.CloseCalled, "session closed as part of cleanup")
		require.Equal(t, 1, fakeClient.session.NS.NewSenderCalled, "sender creation failed, but was called")
		require.Zero(t, fakeClient.session.NS.NewReceiverCalled, "receiver was never created")
	})

	t.Run("receiver fails to be created", func(t *testing.T) {
		// receiver is last in the list, so we'll have to close out sender and session.
		fakeClient := initFn()

		fakeClient.session.NS.NewReceiverErr = errors.New("test error")

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})
		require.EqualError(t, err, "test error")
		require.Nil(t, rpcLink)

		require.Equal(t, 1, fakeClient.session.NS.NewSenderCalled, "sender creation failed, but was called")
		require.Equal(t, 1, fakeClient.session.CloseCalled, "session closed as part of cleanup")
		require.Equal(t, 1, fakeClient.session.NS.Sender.CloseCalled, "sender was closed")
		require.Equal(t, 1, fakeClient.session.NS.NewReceiverCalled, "attempted to create receiver but will fail")
	})

	t.Run("close failures are logged", func(t *testing.T) {
		logsFn := test.CaptureLogsForTest()

		// receiver is last in the list, so we'll have to close out sender and session.
		fakeClient := initFn()

		fakeClient.session.NS.NewReceiverErr = errors.New("test error")
		// we want this failure to be logged, in case it's responsible for our "multiple $cbs links are open" error.
		fakeClient.session.NS.Sender.CloseError = errors.New("sender close failure")

		rpcLink, err := NewRPCLink(context.Background(), RPCLinkArgs{
			Client:   fakeClient,
			Address:  "fake-address",
			LogEvent: log.Event("testing"),
		})
		require.EqualError(t, err, "test error")
		require.Nil(t, rpcLink)

		require.Equal(t, 1, fakeClient.session.NS.NewSenderCalled, "sender creation failed, but was called")
		require.Equal(t, 1, fakeClient.session.CloseCalled, "session closed as part of cleanup")
		require.Equal(t, 1, fakeClient.session.NS.Sender.CloseCalled, "sender was closed")
		require.Equal(t, 1, fakeClient.session.NS.NewReceiverCalled, "attempted to create receiver but will fail")

		logMessages := logsFn()
		require.Equal(t, []string{
			"[azeh.Auth] Failed closing sender for RPC Link: sender close failure",
		}, logMessages)
	})
}

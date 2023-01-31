// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestReceiver_Simulated(t *testing.T) {
	md, client := newClientWithMockedConn(t, nil)
	defer test.RequireClose(t, client)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, sender)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	err = receiver.Close(context.Background())
	require.NoError(t, err)
	require.Equal(t, 3, len(md.Events.GetOpenLinks()), "Sender remains open")

	err = sender.Close(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(md.Events.GetOpenLinks()), "Sender remains open")

	require.Equal(t, 1, len(md.Events.GetOpenConns()), "Connection remains open")

	err = client.Close(context.Background())
	require.NoError(t, err)

	emulation.RequireNoLeaks(t, md.Events)
}

func TestReceiver_Simulated_CloseTopLevelClientClosesChildren(t *testing.T) {
	md, client := newClientWithMockedConn(t, nil)
	defer test.RequireClose(t, client)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, sender)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	err = client.Close(context.Background())
	require.NoError(t, err)

	emulation.RequireNoLeaks(t, md.Events)
}

func TestReceiver_Simulated_Recovery(t *testing.T) {
	md, client := newClientWithMockedConn(t, nil)
	defer test.RequireClose(t, client)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, sender)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	md.DetachSenders("queue")
	md.DetachReceivers("queue")

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello2"),
	}, nil)
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err, "We eat the error in this case since it's recoverable and we want them to try again")
	require.Empty(t, messages)

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3, len(md.Events.GetOpenLinks()), "Sender is still alive, but the Receiver is closed until we call it again...")

	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, 1, len(messages))

	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	err = client.Close(context.Background())
	require.NoError(t, err)

	emulation.RequireNoLeaks(t, md.Events)
}

func TestReceiver_ReceiveMessages_SomeMessagesAndCancelled(t *testing.T) {
	for _, mode := range receiveModesForTests {
		t.Run(mode.Name, func(t *testing.T) {
			md, client := newClientWithMockedConn(t, nil)
			defer test.RequireClose(t, client)

			sender, err := client.NewSender("queue", nil)
			require.NoError(t, err)

			err = sender.SendMessage(context.Background(), &Message{Body: []byte("hello")}, nil)
			require.NoError(t, err)

			test.RequireClose(t, sender)

			receiver, err := client.NewReceiverForQueue("queue", &ReceiverOptions{ReceiveMode: mode.Val})
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			messages, err := receiver.ReceiveMessages(ctx, 2, nil)
			require.NoError(t, err)
			require.Equal(t, []string{"hello"}, getSortedBodies(messages))

			sender.Close(context.Background())

			require.Equal(t, 3, len(md.Events.GetOpenLinks()))
			require.Equal(t, 1, len(md.Events.GetOpenConns()))
		})
	}
}

func TestReceiver_ReceiveMessages_NoMessagesReceivedAndError(t *testing.T) {
	type args struct {
		Name        string
		InternalErr error
		ReturnedErr error
		ConnAlive   bool
		LinkAlive   bool
	}

	fn := func(args args) {
		t.Run(args.Name, func(t *testing.T) {
			md := emulation.NewMockData(t, &emulation.MockDataOptions{
				PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
					if mr.Source == "queue" {
						mr.EXPECT().Receive(gomock.Any()).Return(nil, args.InternalErr)
					}

					return nil
				},
			})

			client, err := newClientImpl(clientCreds{
				connectionString: "Endpoint=sb://example.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=DEADBEEF",
			}, clientImplArgs{
				NSOptions: []internal.NamespaceOption{
					internal.NamespaceWithNewClientFn(md.NewConnection),
				},
			})

			defer test.RequireClose(t, client)

			require.NoError(t, err)
			require.NotNil(t, client)

			receiver, err := client.NewReceiverForQueue("queue", nil)
			require.NoError(t, err)

			messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
			require.EqualValues(t, args.ReturnedErr, err)
			require.Empty(t, messages)

			if args.ConnAlive {
				require.Equal(t, 1, len(md.Events.GetOpenConns()), "Connection is still alive")
			} else {
				require.Equal(t, 0, len(md.Events.GetOpenConns()), "Connection has been closed")
			}

			if args.LinkAlive {
				require.Equal(t, 3, len(md.Events.GetOpenLinks()), "Links are still alive")
			} else {
				require.Equal(t, 0, len(md.Events.GetOpenLinks()), "Links have been closed")
			}
		})
	}

	fn(args{
		Name:        "Non-retriable errors shut down the connection",
		InternalErr: internal.NewErrNonRetriable("non retriable error"),
		ReturnedErr: internal.NewErrNonRetriable("non retriable error"),
		ConnAlive:   false,
	})

	fn(args{
		Name:        "Cancel errors don't close the connection",
		InternalErr: context.Canceled,
		ReturnedErr: context.Canceled,
		ConnAlive:   true,
		LinkAlive:   true,
	})

	fn(args{
		Name:        "Connection level errors close link and connection",
		InternalErr: &amqp.ConnectionError{},
		ReturnedErr: nil,
		ConnAlive:   false,
		LinkAlive:   false,
	})

	fn(args{
		Name:        "Link level errors close the link",
		InternalErr: amqp.ErrLinkClosed,
		ReturnedErr: nil,
		ConnAlive:   true,
		LinkAlive:   false,
	})
}

func TestReceiver_ReceiveMessages_AllMessagesReceived(t *testing.T) {
	fn := func(receiveMode ReceiveMode) {
		t.Run(ReceiveModeString(receiveMode), func(t *testing.T) {
			md, client := newClientWithMockedConn(t, nil)
			defer test.RequireClose(t, client)

			sender, err := client.NewSender("queue", nil)
			require.NoError(t, err)

			err = sender.SendMessage(context.Background(), &Message{Body: []byte("hello")}, nil)
			require.NoError(t, err)

			err = sender.SendMessage(context.Background(), &Message{Body: []byte("world")}, nil)
			require.NoError(t, err)

			test.RequireClose(t, sender)

			receiver, err := client.NewReceiverForQueue("queue", &ReceiverOptions{
				ReceiveMode: receiveMode,
			})
			require.NoError(t, err)

			messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
			require.NoError(t, err)
			require.Equal(t, []string{"hello", "world"}, getSortedBodies(messages))

			require.Equal(t, 1, len(md.Events.GetOpenConns()))
			require.Equal(t, 3, len(md.Events.GetOpenLinks()), "Receive links are still open")
		})
	}

	fn(ReceiveModePeekLock)
	fn(ReceiveModeReceiveAndDelete)
}

func newClientWithMockedConn(t *testing.T, options *emulation.MockDataOptions) (*emulation.MockData, *Client) {
	md := emulation.NewMockData(t, nil)

	client, err := newClientImpl(clientCreds{
		connectionString: "Endpoint=sb://example.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=DEADBEEF",
	}, clientImplArgs{
		NSOptions: []internal.NamespaceOption{
			internal.NamespaceWithNewClientFn(md.NewConnection),
		},
	})
	require.NoError(t, err)

	return md, client
}

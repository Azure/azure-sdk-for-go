// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestReceiver_Simulated(t *testing.T) {
	md, client := newClientWithMockedConn(t, nil, nil)
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
	md, client := newClientWithMockedConn(t, nil, nil)
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
	md, client := newClientWithMockedConn(t, nil, nil)
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
			md, client := newClientWithMockedConn(t, nil, nil)
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
			md, client := newClientWithMockedConn(t, nil, nil)
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

func TestReceiver_ReceiveMessages_SomeMessagesAndError(t *testing.T) {
	md, client := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				mr.EXPECT().Receive(gomock.Any()).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					return mr.InternalReceive(ctx)
				})
				mr.EXPECT().Receive(gomock.Any()).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					require.NoError(t, ctx.Err())
					return nil, internal.NewErrNonRetriable("non-retriable error on second message")
				})
			}

			return nil
		},
	}, &ClientOptions{})

	defer test.RequireClose(t, client)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{Body: []byte("hello")}, nil)
	require.NoError(t, err)

	test.RequireClose(t, sender)

	messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
	require.Equal(t, []string{"hello"}, getSortedBodies(messages))
	require.NoError(t, err, "error is 'erased' when there are some messages to return")

	require.Equal(t, 0, len(md.Events.GetOpenConns()))
	require.Equal(t, 0, len(md.Events.GetOpenLinks()), "Receive links are still open")
}

func TestReceiver_UserFacingErrors(t *testing.T) {
	var receiveErr error

	_, client := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(mock.NotCancelled).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					return nil, receiveErr
				}).AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)

	var asSBError *Error

	receiveErr = amqp.ErrLinkClosed
	messages, err := receiver.PeekMessages(context.Background(), 1, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	receiveErr = &amqp.ConnectionError{}
	messages, err = receiver.ReceiveDeferredMessages(context.Background(), []int64{1}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	receiveErr = &amqp.ConnectionError{}
	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Empty(t, messages)
	// require.ErrorAs(t, err, &asSBError)
	// require.Equal(t, CodeConnectionLost, asSBError.Code)

	receiveErr = internal.RPCError{Resp: &amqpwrap.RPCResponse{Code: internal.RPCResponseCodeLockLost}}

	id, err := uuid.New()
	require.NoError(t, err)

	msg := &ReceivedMessage{
		LockToken: id,
		RawAMQPMessage: &AMQPAnnotatedMessage{
			linkName: "linkName",
		},
	}

	err = receiver.AbandonMessage(context.Background(), msg, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.CompleteMessage(context.Background(), msg, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.DeadLetterMessage(context.Background(), msg, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.DeferMessage(context.Background(), msg, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.RenewMessageLock(context.Background(), msg, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)
}

func TestReceiver_ReceiveMessages(t *testing.T) {
	_, client := newClientWithMockedConn(t, nil, nil)
	defer test.RequireClose(t, client)

	receiver, err := client.NewReceiverForQueue("queue", &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)
	defer test.RequireClose(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)
	defer test.RequireClose(t, sender)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("Received message 1"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 10, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"Received message 1"}, getSortedBodies(messages))

	links := receiver.amqpLinks.(*internal.AMQPLinksImpl)
	require.Equal(t, uint32(9), links.Receiver.Credits())
}

func TestReceive_ReuseExistingCredits(t *testing.T) {
	type contextKey string
	const key = contextKey("CalledFromReceiveMessages")

	_, client := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				mr.EXPECT().Receive(gomock.Any()).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					if ctx.Value(key) != nil {
						log.Printf("Doing receive, called from ReceiveMessages")
						return mr.InternalReceive(ctx)
					} else {
						log.Printf("Waiting for cancellation, called from Releaser loop")
						<-ctx.Done()
						log.Printf("Cancellation, we should exit from Releaser loop")
						return nil, ctx.Err()
					}
				}).AnyTimes()
			}

			return nil
		},
	}, nil)
	defer test.RequireClose(t, client)

	// we want to end up in a situation where we have excess credits.
	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("message 1"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.WithValue(context.Background(), key, ""), 5, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"message 1"}, getSortedBodies(messages))

	links := receiver.amqpLinks.(*internal.AMQPLinksImpl)
	require.Equal(t, uint32(4), links.Receiver.Credits())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("message 2"),
	}, nil)
	require.NoError(t, err)

	// now we've got credits left over so we won't have to issue _more_ credits
	messages, err = receiver.ReceiveMessages(context.WithValue(context.Background(), key, ""), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"message 2"}, getSortedBodies(messages))

	require.Equal(t, uint32(3), links.Receiver.Credits(), "We re-used our already issued credits")

	// now let's request _more_ than what we have. We'll issue enough credits, taking into account what
	// we already have.
	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("message 3"),
	}, nil)
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(context.WithValue(context.Background(), key, ""), 1001, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"message 3"}, getSortedBodies(messages))

	require.Equal(t, uint32(1001-1), links.Receiver.Credits(), "We re-used our already issued credits")
}

func TestReceiver_ReceiveMessages_MessageReleaser(t *testing.T) {
	md, client := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				mr.EXPECT().Receive(gomock.Any()).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					return mr.InternalReceive(ctx)
				}).AnyTimes()
			}

			return nil
		},
	}, nil)
	defer test.RequireClose(t, client)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("received message 1"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 3, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"received message 1"}, getSortedBodies(messages))

	// We were able to get one message during this ReceiveMessages() call
	// which means we still have 2 active credits. If messages arrive in
	// between they'll be consumed and released.
	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("message available again after being released by releaser"),
	}, nil)
	require.NoError(t, err)

	// keep running until the releaser receives and releases the message, since
	// we're in between ReceiveMessages() calls.
	for evt := range md.Events.Chan() {
		if evt.Type == emulation.EventTypeLinkDisposition {
			dispEvt := evt.Data.(emulation.DispositionEvent)

			if dispEvt.LinkEvent.Entity == "queue" && string(dispEvt.Data[0]) == "message available again after being released by releaser" {
				break
			}
		}
	}

	// we can receive now - the message will be consumed again (.Release() just lets the broker serve it up again)
	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"message available again after being released by releaser"}, getSortedBodies(messages))
}

func newClientWithMockedConn(t *testing.T, mockDataOptions *emulation.MockDataOptions, clientOptions *ClientOptions) (*emulation.MockData, *Client) {
	md := emulation.NewMockData(t, mockDataOptions)

	client, err := newClientImpl(clientCreds{
		connectionString: "Endpoint=sb://example.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=DEADBEEF",
	}, clientImplArgs{
		ClientOptions: clientOptions,
		NSOptions: []internal.NamespaceOption{
			internal.NamespaceWithNewClientFn(md.NewConnection),
		},
	})
	require.NoError(t, err)

	return md, client
}

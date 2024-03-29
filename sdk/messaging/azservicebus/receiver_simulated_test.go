// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"log"
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestReceiver_Simulated(t *testing.T) {
	md, client, cleanup := newClientWithMockedConn(t, nil, nil)
	defer cleanup()

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
	md, client, cleanup := newClientWithMockedConn(t, nil, nil)
	defer cleanup()

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
	md, client, cleanup := newClientWithMockedConn(t, nil, nil)
	defer cleanup()

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
	require.Equal(t, []string{"hello"}, getSortedBodies(messages))

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	md.DetachSenders("queue")
	md.DetachReceivers("queue")

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello2"),
	}, nil)
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err, "receiver was able to recover from error")
	require.Equal(t, []string{"hello2"}, getSortedBodies(messages))

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver both recover from their forced detach")

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello3"),
	}, nil)
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"hello3"}, getSortedBodies(messages))

	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")

	err = client.Close(context.Background())
	require.NoError(t, err)

	emulation.RequireNoLeaks(t, md.Events)
}

func TestReceiver_ReceiveMessages_SomeMessagesAndCancelled(t *testing.T) {
	for _, mode := range receiveModesForTests {
		t.Run(mode.Name, func(t *testing.T) {
			md, client, cleanup := newClientWithMockedConn(t, nil, nil)
			defer cleanup()

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
		ReceiveErr      error // error returned when AMQPReceiver.Receive() is called (for non-$cbs links)
		ExpectedErr     error // error that should be returned after all internal retries have occurred.
		ExpectConnAlive bool
		ExpectLinkAlive bool
	}

	testFn := func(t *testing.T, args args) {
		t.Parallel()

		// make it so any AMQPReceiver created here will return args.InternalErr when AMQPReceiver.Receive()
		// is called.
		md, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
			PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
				if mr.Source == "queue" {
					mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).Return(nil, args.ReceiveErr)
				}

				return nil
			},
		}, &ClientOptions{
			RetryOptions: exported.RetryOptions{
				RetryDelay:    0,
				MaxRetryDelay: 0,
				MaxRetries:    1,
			},
		})
		defer cleanup()

		receiver, err := client.NewReceiverForQueue("queue", nil)
		require.NoError(t, err)

		// internally we're handling args.InternalErr in our retry loop.
		messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
		require.EqualValues(t, args.ExpectedErr, err)
		require.Empty(t, messages)

		if args.ExpectConnAlive {
			require.Equal(t, 1, len(md.Events.GetOpenConns()), "Connection is still alive")
		} else {
			require.Equal(t, 0, len(md.Events.GetOpenConns()), "Connection has been closed")
		}

		if args.ExpectLinkAlive {
			require.Equal(t, 3, len(md.Events.GetOpenLinks()), "Links are still alive")
		} else {
			require.Equal(t, 0, len(md.Events.GetOpenLinks()), "Links have been closed")
		}
	}

	t.Run("Non-retriable errors shut down the connection", func(t *testing.T) {
		testFn(t, args{
			ReceiveErr:      internal.NewErrNonRetriable("non retriable error"),
			ExpectedErr:     internal.NewErrNonRetriable("non retriable error"),
			ExpectConnAlive: false,
		})
	})

	t.Run("Cancel errors don't close the connection", func(t *testing.T) {
		testFn(t, args{
			ReceiveErr:      context.Canceled,
			ExpectedErr:     context.Canceled,
			ExpectConnAlive: true,
			ExpectLinkAlive: true,
		})
	})

	t.Run("Connection level errors close link and connection", func(t *testing.T) {
		testFn(t, args{
			ReceiveErr: &amqp.ConnError{},
			ExpectedErr: exported.NewError(
				exported.CodeConnectionLost,
				&amqp.ConnError{},
			),
			ExpectConnAlive: false,
			ExpectLinkAlive: false,
		})
	})

	t.Run("Link level errors close the link", func(t *testing.T) {
		testFn(t, args{
			ReceiveErr: &amqp.LinkError{},
			ExpectedErr: exported.NewError(
				exported.CodeConnectionLost,
				&amqp.LinkError{},
			),
			ExpectConnAlive: true,
			ExpectLinkAlive: false,
		})
	})
}

func TestReceiver_ReceiveMessages_AllMessagesReceived(t *testing.T) {
	fn := func(receiveMode ReceiveMode) {
		t.Run(ReceiveModeString(receiveMode), func(t *testing.T) {
			md, client, cleanup := newClientWithMockedConn(t, nil, nil)
			defer cleanup()

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
	md, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return mr.InternalReceive(ctx, o)
				})
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					require.NoError(t, ctx.Err())
					return nil, internal.NewErrNonRetriable("non-retriable error on second message")
				})
			}

			return nil
		},
	}, &ClientOptions{})
	defer cleanup()

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

	// all AMQPReceivers created from this client will return receiveErr whenver AMQPReceiver.Receive() is called.
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(mock.NotCancelled, gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return nil, receiveErr
				}).AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})
	defer cleanup()

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)

	var asSBError *Error

	// forcing a link error to come back on first use of AMQPReceiver.Receive()
	receiveErr = &amqp.LinkError{}
	messages, err := receiver.PeekMessages(context.Background(), 1, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	receiveErr = &amqp.ConnError{}
	messages, err = receiver.ReceiveDeferredMessages(context.Background(), []int64{1}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
	require.Empty(t, messages)

	receiveErr = &amqp.ConnError{}
	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
	require.Empty(t, messages)

	receiveErr = internal.RPCError{Resp: &amqpwrap.RPCResponse{Code: internal.RPCResponseCodeLockLost}}

	id, err := uuid.New()
	require.NoError(t, err)

	msg := &ReceivedMessage{
		LockToken: id,
		RawAMQPMessage: &AMQPAnnotatedMessage{
			inner: &amqp.Message{},
		},
		linkName:         "link-name",
		settleOnMgmtLink: true,
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
	_, client, cleanup := newClientWithMockedConn(t, nil, nil)
	defer cleanup()

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

	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					if ctx.Value(key) != nil {
						log.Printf("Doing receive, called from ReceiveMessages")
						return mr.InternalReceive(ctx, o)
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
	defer cleanup()

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
	md, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return mr.InternalReceive(ctx, o)
				}).AnyTimes()
			}

			return nil
		},
	}, nil)
	defer cleanup()

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

func TestReceiver_ReceiveMessages_CreditValidation(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, nil, nil)
	defer cleanup()

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	messages, err := receiver.ReceiveMessages(context.Background(), 5001, nil)
	require.EqualError(t, err, "maxMessages cannot exceed 5000")
	require.Empty(t, messages)

	messages, err = receiver.ReceiveMessages(context.Background(), -1, nil)
	require.EqualError(t, err, "maxMessages should be greater than 0")
	require.Empty(t, messages)

	messages, err = receiver.ReceiveMessages(context.Background(), 0, nil)
	require.EqualError(t, err, "maxMessages should be greater than 0")
	require.Empty(t, messages)
}

func TestReceiver_CreditsDontExceedMax(t *testing.T) {
	type keyType string
	totalCreditIssued := uint32(0)

	md, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue" {
				fn := func(credit uint32) error {
					totalCreditIssued += credit
					return mr.InternalIssueCredit(credit)
				}

				// first actual request, 5000 fresh credits.
				mr.EXPECT().IssueCredit(uint32(5000)).DoAndReturn(fn)

				// we're going to eat up one credit with a Receive() call and then
				// issue 5000 again, and should only need to issue 1 new credit.
				mr.EXPECT().IssueCredit(uint32(1)).DoAndReturn(fn)

				mr.EXPECT().Receive(mock.NewContextWithValueMatcher(keyType("FromReceive"), true), gomock.Nil()).DoAndReturn(mr.InternalReceive).AnyTimes()

				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					// interaction with the releaser just makes this test harder to make predictable and doesn't
					// add or change anything.
					azlog.Writef(azlog.Event("testing"), "===> Releaser asking for message, blocking on cancel.")
					<-ctx.Done()
					return nil, ctx.Err()
				}).AnyTimes()

				require.EqualValues(t, -1, mr.Opts.Credit)
			}

			return nil
		},
	}, nil)
	defer cleanup()

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	// we've got a gomock matcher that looks at the context and activates our InternalReceive() function.
	baseReceiveCtx := context.WithValue(context.Background(), keyType("FromReceive"), true)

	ctx, cancel := context.WithTimeout(baseReceiveCtx, time.Second)
	defer cancel()

	messages, err := receiver.ReceiveMessages(ctx, 5000, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Empty(t, messages)

	err = sender.SendMessage(context.Background(), &Message{Body: []byte("hello world")}, nil)
	require.NoError(t, err)

	// no issue credit needed - we've still got the 5000 from last time since we didn't
	// receive any messages.
	messages, err = receiver.ReceiveMessages(baseReceiveCtx, 5000, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"hello world"}, getSortedBodies(messages))

	require.Equal(t, uint32(5000), totalCreditIssued)

	ctx, cancel = context.WithTimeout(baseReceiveCtx, time.Second)
	defer cancel()

	// we ate a credit last time since we received a single message, so this time we'll still
	// need to issue some more to backfill.
	messages, err = receiver.ReceiveMessages(ctx, 5000, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Empty(t, messages)

	require.Equal(t, uint32(5000+1), totalCreditIssued)

	require.Equal(t, 1, len(md.Events.GetOpenConns()))
	require.Equal(t, 3+3, len(md.Events.GetOpenLinks()), "Sender and Receiver each own 3 links apiece ($mgmt, actual link)")
}

func TestSessionReceiver_ConnectionDeadForAccept(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				return &amqp.ConnError{}
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})
	defer cleanup()

	receiver, err := client.AcceptSessionForQueue(context.Background(), "queue", "session ID", nil)
	var sbErr *Error
	require.ErrorAs(t, err, &sbErr)
	require.Nil(t, receiver)
}

func TestSessionReceiverUserFacingErrors_Methods(t *testing.T) {
	lockLost := false

	mgmtStub := func(ctx context.Context, o *amqp.ReceiveOptions, mr *emulation.MockReceiver) (*amqp.Message, error) {
		msg, _ := mr.InternalReceive(ctx, o)

		if lockLost {
			return nil, &amqp.Error{
				Condition: amqp.ErrCond("com.microsoft:message-lock-lost"),
			}
		}

		// TODO: this is hacky - we don't have a full mgmt link like we do with $cbs.
		return &amqp.Message{
			Properties: &amqp.MessageProperties{
				CorrelationID: msg.Properties.MessageID,
			},
			ApplicationProperties: map[string]any{
				"status-code": int32(200),
			},
			Value: map[string]any{
				"expiration": time.Now().Add(time.Hour),
			},
		}, nil
	}

	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source == "queue/$management" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return mgmtStub(ctx, o, mr)
				}).AnyTimes()
			} else if mr.Source != "$cbs" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return nil, &amqp.ConnError{}
				}).AnyTimes()

				mr.EXPECT().LinkSourceFilterValue("com.microsoft:session-filter").Return("session ID").AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})
	defer cleanup()

	// we'll return valid responses for the mgmt link since we need
	// that to get a session receiver.
	receiver, err := client.AcceptSessionForQueue(context.Background(), "queue", "session ID", nil)
	require.NoError(t, err)

	// now replace it so we get connection errors.
	var asSBError *Error

	lockLost = true

	state, err := receiver.GetSessionState(context.Background(), nil)
	require.Nil(t, state)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.SetSessionState(context.Background(), []byte{}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.RenewSessionLock(context.Background(), nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)
}

func newClientWithMockedConn(t *testing.T, mockDataOptions *emulation.MockDataOptions, clientOptions *ClientOptions) (*emulation.MockData, *Client, func()) {
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

	return md, client, func() {
		test.RequireClose(t, client)
		md.Close()
	}
}

var noRetriesNeeded = exported.RetryOptions{
	MaxRetries:    -1,
	RetryDelay:    0,
	MaxRetryDelay: 0,
}

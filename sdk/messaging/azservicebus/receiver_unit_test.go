// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func TestReceiver_ReceiveMessages_AMQPLinksFailure(t *testing.T) {
	fakeAMQPLinks := &internal.FakeAMQPLinks{
		Err: internal.NewErrNonRetriable("failed to create links"),
	}

	receiver := &Receiver{
		amqpLinks:   fakeAMQPLinks,
		receiveMode: ReceiveModePeekLock,
		// TODO: need to make this test rely less on stubbing.
		cancelReleaser: &atomic.Value{},
	}

	receiver.cancelReleaser.Store(emptyCancelFn)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.Equal(t, internal.RecoveryKindFatal, internal.GetRecoveryKind(err))
	require.Equal(t, "failed to create links", err.Error())
	require.Empty(t, messages)
}

var receiveModesForTests = []struct {
	Name string
	Val  ReceiveMode
}{
	{Name: "peekLock", Val: ReceiveModePeekLock},
	{Name: "receiveAndDelete", Val: ReceiveModeReceiveAndDelete},
}

func TestReceiver_ReceiveMessages_SomeMessagesAndCancelled(t *testing.T) {
	for _, mode := range receiveModesForTests {
		t.Run(mode.Name, func(t *testing.T) {
			fakeAMQPReceiver := &internal.FakeAMQPReceiver{
				ReceiveResults: []struct {
					M *amqp.Message
					E error
				}{
					{M: &amqp.Message{Data: [][]byte{[]byte("hello")}}},
					// after this the context will block until the cancellation context's deadline fires.
				},
			}

			fakeAMQPLinks := &internal.FakeAMQPLinks{
				Receiver: fakeAMQPReceiver,
			}

			receiver, err := newReceiver(newReceiverArgs{
				ns:     &internal.FakeNS{AMQPLinks: fakeAMQPLinks},
				entity: entity{Queue: "queue"},
			}, &ReceiverOptions{ReceiveMode: mode.Val})
			require.NoError(t, err)

			messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
			require.NoError(t, err)
			require.Equal(t, []string{"hello"}, getSortedBodies(messages))

			// and the links did not need to be closed for a cancellation
			require.Equal(t, 0, fakeAMQPLinks.Closed)
			require.Equal(t, 1, fakeAMQPLinks.CloseIfNeededCalled)
		})
	}
}

func TestReceiver_ReceiveMessages_NoMessagesReceivedAndError(t *testing.T) {
	var errors = []struct {
		Err      error
		Expected error
	}{
		// a fatal error is always returned in peekLock mode
		{Err: internal.NewErrNonRetriable("non retriable error"), Expected: internal.NewErrNonRetriable("non retriable error")},
		// non-fatal errors are "erased" and the error will be caught on the next iteration of the loop
		{Err: amqp.ErrLinkClosed, Expected: nil},
	}

	// all the receive modes work the same when there are no messages
	for _, mode := range receiveModesForTests {
		for i, data := range errors {
			t.Run(fmt.Sprintf("%s [%d] %s", mode.Name, i, data.Err), func(t *testing.T) {
				fakeAMQPReceiver := &internal.FakeAMQPReceiver{
					ReceiveResults: []struct {
						M *amqp.Message
						E error
					}{
						{E: data.Err},
					},
				}

				fakeAMQPLinks := &internal.FakeAMQPLinks{
					Receiver: fakeAMQPReceiver,
				}

				receiver, err := newReceiver(newReceiverArgs{
					ns:     &internal.FakeNS{AMQPLinks: fakeAMQPLinks},
					entity: entity{Queue: "queue"},
				}, &ReceiverOptions{ReceiveMode: mode.Val})
				require.NoError(t, err)

				messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
				require.EqualValues(t, data.Expected, err)
				require.Empty(t, messages)

				require.Equal(t, 1, fakeAMQPReceiver.PrefetchedCalled, "prefetched before throwing away the broken link")

				// a fatal error happened, links should be closed.
				require.Equal(t, 0, fakeAMQPLinks.Closed, "links are closed using CloseIfNeeded")
				require.Equal(t, 1, fakeAMQPLinks.CloseIfNeededCalled, "links are closed on receive errors")
			})
		}
	}
}

func TestReceiver_ReceiveMessages_AllMessagesReceived(t *testing.T) {
	for _, mode := range receiveModesForTests {
		t.Run(mode.Name, func(t *testing.T) {
			fakeAMQPReceiver := &internal.FakeAMQPReceiver{
				ReceiveResults: []struct {
					M *amqp.Message
					E error
				}{
					{M: &amqp.Message{Data: [][]byte{[]byte("hello")}}},
					{M: &amqp.Message{Data: [][]byte{[]byte("world")}}},
				},
			}

			fakeAMQPLinks := &internal.FakeAMQPLinks{
				Receiver: fakeAMQPReceiver,
			}

			receiver, err := newReceiver(newReceiverArgs{
				ns:     &internal.FakeNS{AMQPLinks: fakeAMQPLinks},
				entity: entity{Queue: "queue"},
			}, &ReceiverOptions{ReceiveMode: mode.Val})
			require.NoError(t, err)

			messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
			require.NoError(t, err)
			require.Equal(t, []string{"hello", "world"}, getSortedBodies(messages))

			// and the links did not need to be closed for a cancellation
			require.Equal(t, 0, fakeAMQPLinks.Closed)
			require.Equal(t, 1, fakeAMQPLinks.CloseIfNeededCalled, "called, but with a benign error")
		})
	}
}

func TestReceiver_ReceiveMessages_SomeMessagesAndError(t *testing.T) {
	fakeAMQPReceiver := &internal.FakeAMQPReceiver{
		ReceiveResults: []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte("hello")}}},
			{E: internal.NewErrNonRetriable("non-retriable error on second message")},
		},
	}

	fakeAMQPLinks := &internal.FakeAMQPLinks{
		Receiver: fakeAMQPReceiver,
	}

	receiver, err := newReceiver(newReceiverArgs{
		ns:     &internal.FakeNS{AMQPLinks: fakeAMQPLinks},
		entity: entity{Queue: "queue"},
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
	require.Equal(t, []string{"hello"}, getSortedBodies(messages))
	require.NoError(t, err, "error is 'erased' when there are some messages to return")

	// a fatal error happened, links should be closed.
	require.Equal(t, 0, fakeAMQPLinks.Closed, "links are closed using CloseIfNeeded")
	require.Equal(t, 1, fakeAMQPLinks.CloseIfNeededCalled, "prefetch is called")
}

func TestReceiverCancellationUnitTests(t *testing.T) {
	t.Run("ImmediatelyCancelled", func(t *testing.T) {
		r := &Receiver{
			amqpLinks: &internal.FakeAMQPLinks{
				Receiver: &internal.FakeAMQPReceiver{},
			},
			cancelReleaser: &atomic.Value{},
		}

		r.cancelReleaser.Store(emptyCancelFn)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		msgs, err := r.ReceiveMessages(ctx, 95, nil)
		require.Empty(t, msgs)
		require.True(t, internal.IsCancelError(err))
	})

	t.Run("CancelledWhileReceiving", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		r := &Receiver{
			defaultTimeAfterFirstMsg: time.Second,
			defaultDrainTimeout:      time.Second,
			amqpLinks: &internal.FakeAMQPLinks{
				Receiver: &internal.FakeAMQPReceiver{
					ReceiveFn: func(ctx context.Context) (*amqp.Message, error) {
						cancel()
						return nil, ctx.Err()
					},
				},
			},
			cancelReleaser: &atomic.Value{},
		}

		r.cancelReleaser.Store(emptyCancelFn)

		msgs, err := r.ReceiveMessages(ctx, 95, nil)
		require.Empty(t, msgs)
		require.ErrorIs(t, err, context.Canceled)
	})
}

func TestReceiverOptions(t *testing.T) {
	// defaults
	receiver := &Receiver{}
	e := &entity{Topic: "topic", Subscription: "subscription"}

	require.NoError(t, applyReceiverOptions(receiver, e, nil))

	require.EqualValues(t, ReceiveModePeekLock, receiver.receiveMode)
	path, err := e.String()
	require.NoError(t, err)
	require.EqualValues(t, "topic/Subscriptions/subscription", path)

	// using options
	receiver = &Receiver{}
	e = &entity{Topic: "topic", Subscription: "subscription"}

	require.NoError(t, applyReceiverOptions(receiver, e, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
		SubQueue:    SubQueueTransfer,
	}))

	require.EqualValues(t, ReceiveModeReceiveAndDelete, receiver.receiveMode)
	path, err = e.String()
	require.NoError(t, err)
	require.EqualValues(t, "topic/Subscriptions/subscription/$Transfer/$DeadLetterQueue", path)
}

func TestReceiver_UserFacingErrors(t *testing.T) {
	fakeAMQPLinks := &internal.FakeAMQPLinks{}

	receiver, err := newReceiver(newReceiverArgs{
		ns: &internal.FakeNS{
			AMQPLinks: fakeAMQPLinks,
		},
		entity: entity{
			Queue: "queue",
		},
		cleanupOnClose:      func() {},
		getRecoveryKindFunc: internal.GetRecoveryKind,
		newLinkFn: func(ctx context.Context, session amqpwrap.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
			return nil, nil, nil
		},
		retryOptions: RetryOptions{
			MaxRetries:    0,
			RetryDelay:    0,
			MaxRetryDelay: 0,
		},
	}, nil)
	require.NoError(t, err)

	var asSBError *Error

	fakeAMQPLinks.Err = amqp.ErrLinkClosed
	messages, err := receiver.PeekMessages(context.Background(), 1, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	fakeAMQPLinks.Err = &amqp.ConnectionError{}
	messages, err = receiver.ReceiveDeferredMessages(context.Background(), []int64{1}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	fakeAMQPLinks.Err = &amqp.ConnectionError{}
	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	fakeAMQPLinks.Err = internal.RPCError{Resp: &internal.RPCResponse{Code: internal.RPCResponseCodeLockLost}}

	err = receiver.AbandonMessage(context.Background(), &ReceivedMessage{}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.CompleteMessage(context.Background(), &ReceivedMessage{}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.DeadLetterMessage(context.Background(), &ReceivedMessage{}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.DeferMessage(context.Background(), &ReceivedMessage{}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)

	err = receiver.RenewMessageLock(context.Background(), &ReceivedMessage{}, nil)
	require.Empty(t, messages)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)
}

func TestReceiver_releaserFunc(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), nil)
	require.NoError(t, err)

	successfulReleases := 0

	amqpReceiver := internal.FakeAMQPReceiver{
		ReceiveFn: func(ctx context.Context) (*amqp.Message, error) {
			return &amqp.Message{
				Data: [][]byte{[]byte("hello")},
			}, nil
		},
		ReleaseMessageFn: func(ctx context.Context, msg *amqp.Message) error {
			require.Equal(t, "hello", string(msg.Data[0]))
			successfulReleases++
			go receiver.cancelReleaser.Load().(func() string)()
			return nil
		},
	}

	logsFn := test.CaptureLogsForTest()

	releaserFn := receiver.newReleaserFunc(&amqpReceiver)
	releaserFn()

	// some non-determinism since we launched the cancel asynchronously
	// but we'll get at least one.
	require.LessOrEqual(t, 1, amqpReceiver.ReceiveCalled)
	require.LessOrEqual(t, 1, amqpReceiver.ReleaseMessageCalled)

	require.Contains(t,
		logsFn(),
		fmt.Sprintf("[azsb.Receiver] [fakelink] Message releaser pausing. Released %d messages", successfulReleases),
	)
}

func TestReceiver_releaserFunc_errorOnFirstMessage(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), nil)
	require.NoError(t, err)

	amqpReceiver := internal.FakeAMQPReceiver{
		ReleaseMessageFn: func(ctx context.Context, msg *amqp.Message) error {
			panic("Not called for this test since Receive() is returning an error")
		},
	}

	amqpReceiver.ReceiveFn = func(ctx context.Context) (*amqp.Message, error) {
		if amqpReceiver.ReceiveCalled > 2 {
			return nil, amqp.ErrLinkClosed
		}

		// This is one of the few error types classified as RecoveryKindNone
		// in the releaser this means we'll just retry since the link is still
		// considered good at this point.
		return nil, &amqp.Error{
			Condition: amqp.ErrorCondition("com.microsoft:server-busy"),
		}
	}

	logsFn := test.CaptureLogsForTest()

	releaserFn := receiver.newReleaserFunc(&amqpReceiver)
	releaserFn()

	// we got called a few times, but none of them succeeded.
	require.Equal(t, 2+1, amqpReceiver.ReceiveCalled)

	require.Contains(t,
		logsFn(),
		fmt.Sprintf("[azsb.Receiver] [fakelink] Message releaser stopping because of link failure. Released 0 messages. Will start again after next receive: %s", amqp.ErrLinkClosed))
}

func TestReceiver_releaserFunc_receiveAndDeleteIsNoop(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	amqpReceiver := internal.FakeAMQPReceiver{
		ReceiveFn: func(ctx context.Context) (*amqp.Message, error) {
			panic("Should not be used in this test")
		},
		ReleaseMessageFn: func(ctx context.Context, msg *amqp.Message) error {
			panic("Should not be used in this test")
		},
	}

	releaserFn := receiver.newReleaserFunc(&amqpReceiver)

	// cancelling is still the empty function
	cancelFn := receiver.cancelReleaser.Load().(func() string)
	require.Equal(t, "empty", cancelFn())

	// in this case you don't need to cancel anything - it's just no-op
	// Note it'll just exit immediately, the "releaser" doesn't block here.
	releaserFn()

	require.LessOrEqual(t, 0, amqpReceiver.ReceiveCalled)
	require.LessOrEqual(t, 0, amqpReceiver.ReleaseMessageCalled)
}

func TestReceiver_fetchMessages_FirstMessageFailure(t *testing.T) {
	errors := []error{amqp.ErrLinkClosed, context.Canceled}

	for _, err := range errors {
		t.Run("error: "+err.Error(), func(t *testing.T) {
			receiver, err := newReceiver(defaultNewReceiverArgsForTest(), &ReceiverOptions{
				ReceiveMode: ReceiveModeReceiveAndDelete,
			})
			require.NoError(t, err)

			amqpReceiver := &internal.FakeAMQPReceiver{
				ReceiveResults: []struct {
					M *amqp.Message
					E error
				}{
					{
						M: nil,
						E: amqp.ErrLinkClosed,
					},
				},
				PrefetchedResults: []*amqp.Message{
					{Data: [][]byte{[]byte(("prefetched message 1"))}},
					{Data: [][]byte{[]byte(("prefetched message 2"))}},
					{Data: [][]byte{[]byte(("prefetched message 3"))}},
					// not used since we'd return too many results (they onlyu requested '3' below)
					{Data: [][]byte{[]byte(("prefetched message 4"))}},
				},
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			res := receiver.fetchMessages(ctx, amqpReceiver, 3, time.Hour)
			require.ErrorIs(t, res.Error, amqp.ErrLinkClosed)

			require.Equal(t, []*amqp.Message{
				{Data: [][]byte{[]byte(("prefetched message 1"))}},
				{Data: [][]byte{[]byte(("prefetched message 2"))}},
				{Data: [][]byte{[]byte(("prefetched message 3"))}},
			}, res.Messages)

			// and we should have messages remaining in our prefetch
			require.Equal(t, []*amqp.Message{
				{Data: [][]byte{[]byte(("prefetched message 4"))}},
			}, amqpReceiver.PrefetchedResults)
		})
	}
}

func TestReceiver_fetchMessages_DontOverflow(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	amqpReceiver := &internal.FakeAMQPReceiver{
		ReceiveResults: []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte(("received message 1"))}}},
			{M: &amqp.Message{Data: [][]byte{[]byte(("received message 2"))}}},
			{M: &amqp.Message{Data: [][]byte{[]byte(("received message 3"))}}},
			{M: &amqp.Message{Data: [][]byte{[]byte(("<receive: will not get received here>"))}}},
		},
		PrefetchedResults: []*amqp.Message{
			{Data: [][]byte{[]byte(("<prefetched: will not get used>"))}},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res := receiver.fetchMessages(ctx, amqpReceiver, 3, time.Hour)
	require.NoError(t, res.Error)

	require.Equal(t, []*amqp.Message{
		{Data: [][]byte{[]byte(("received message 1"))}},
		{Data: [][]byte{[]byte(("received message 2"))}},
		{Data: [][]byte{[]byte(("received message 3"))}},
	}, res.Messages)

	require.Equal(t, 1, len(amqpReceiver.ReceiveResults))
	require.Equal(t,
		&amqp.Message{Data: [][]byte{[]byte(("<receive: will not get received here>"))}},
		amqpReceiver.ReceiveResults[0].M)

	require.Equal(t,
		[]*amqp.Message{{Data: [][]byte{[]byte(("<prefetched: will not get used>"))}}},
		amqpReceiver.PrefetchedResults)
}

func TestReceiver_fetchMessages_TimeAfterFirstMessageCancels(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	amqpReceiver := &internal.FakeAMQPReceiver{
		ReceiveResults: []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte("Received message 1")}}},
			{M: &amqp.Message{Data: [][]byte{[]byte("Received message 2")}}},
		},
		PrefetchedResults: []*amqp.Message{
			{Data: [][]byte{[]byte("Prefetched message 1")}},
			{Data: [][]byte{[]byte("<will be ignored 1>")}},
			{Data: [][]byte{[]byte("<will be ignored 2>")}},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timeAfterFirstMessage := time.Second
	res := receiver.fetchMessages(ctx, amqpReceiver, 3, timeAfterFirstMessage)
	require.NoError(t, res.Error, "No error since it was just the timeAfterFirstMessage cancellation")

	require.Equal(t, []*amqp.Message{
		{Data: [][]byte{[]byte("Received message 1")}},
		{Data: [][]byte{[]byte("Received message 2")}},
		{Data: [][]byte{[]byte("Prefetched message 1")}},
	}, res.Messages)

	require.Empty(t, 0, len(amqpReceiver.ReceiveResults))
	require.Equal(t,
		[]*amqp.Message{
			{Data: [][]byte{[]byte("<will be ignored 1>")}},
			{Data: [][]byte{[]byte("<will be ignored 2>")}},
		},
		amqpReceiver.PrefetchedResults)
}

func TestReceiver_fetchMessages_UserCancelsAfterFirstMessage(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	testMessages := []*amqp.Message{
		{Data: [][]byte{[]byte("Received message 1")}},
		{Data: [][]byte{[]byte("Received message 2")}},
	}

	usersCtx, cancelUsersCtx := context.WithCancel(context.Background())
	defer cancelUsersCtx()

	amqpReceiver := &internal.FakeAMQPReceiver{
		ReceiveFn: func(ctx context.Context) (*amqp.Message, error) {
			msg := testMessages[0]
			testMessages = testMessages[1:]

			if len(testMessages) == 0 {
				cancelUsersCtx()
			}

			return msg, nil
		},
		PrefetchedResults: []*amqp.Message{
			{Data: [][]byte{[]byte("Prefetched message 1")}},
			{Data: [][]byte{[]byte("<will be ignored 1>")}},
			{Data: [][]byte{[]byte("<will be ignored 2>")}},
		},
	}

	timeAfterFirstMessage := time.Second
	res := receiver.fetchMessages(usersCtx, amqpReceiver, 3, timeAfterFirstMessage)
	require.ErrorIs(t, res.Error, context.Canceled, "Users cancellation error is propagated")

	require.Equal(t, []*amqp.Message{
		{Data: [][]byte{[]byte("Received message 1")}},
		{Data: [][]byte{[]byte("Received message 2")}},
		{Data: [][]byte{[]byte("Prefetched message 1")}},
	}, res.Messages)

	require.Empty(t, 0, len(amqpReceiver.ReceiveResults))
	require.Equal(t,
		[]*amqp.Message{
			{Data: [][]byte{[]byte("<will be ignored 1>")}},
			{Data: [][]byte{[]byte("<will be ignored 2>")}},
		},
		amqpReceiver.PrefetchedResults)
}

func TestReceiver_ReceiveMessages_RollingCredits_NoMessages(t *testing.T) {
	amqpReceiver := &internal.FakeAMQPReceiver{
		ReceiveResults: nil,
	}

	args := defaultNewReceiverArgsForTest()

	ns := args.ns.(*internal.FakeNS)
	ns.AMQPLinks = &internal.FakeAMQPLinks{
		Receiver: amqpReceiver,
	}

	receiver, err := newReceiver(args, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	t.Run("no initial credits, must issue new credit", func(t *testing.T) {
		amqpReceiver.ReceiveResults = []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte("Received message 1")}}},
			{E: context.Canceled},
		}

		messages, err := receiver.ReceiveMessages(context.Background(), 10, nil)
		require.NoError(t, err, "User didn't cancel so cancellation is erased")
		require.Equal(t, []string{"Received message 1"}, getSortedBodies(messages))

		// note, our AMQPReceiver is a stub so the credits here are whatever
		// RecieveMessages() actually just IssueCredit'd.
		require.Equal(t, uint32(10), amqpReceiver.Credits())
		require.Empty(t, amqpReceiver.ReceiveResults)
	})

	t.Run("existing credit, but not enough to cover all the requested credits", func(t *testing.T) {
		amqpReceiver.ReceiveResults = []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte("Received message 2")}}},
			{E: context.Canceled},
		}

		// now this time we have 10 credits, so our next receive won't need to add as many
		messages, err := receiver.ReceiveMessages(context.Background(), 101, nil)
		require.NoError(t, err, "User didn't cancel so cancellation is erased")
		require.Equal(t, []string{"Received message 2"}, getSortedBodies(messages))

		require.Equal(t, uint32(10)+uint32(101-10), amqpReceiver.Credits(), "Credits includes the existing credits on the line and adds more to make up the difference to get to 101")
		require.Empty(t, amqpReceiver.ReceiveResults)
	})

	t.Run("credits on the line are already enough to cover our request", func(t *testing.T) {
		amqpReceiver.ReceiveResults = []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte("Received message 3")}}},
			{E: context.Canceled},
		}

		// now we're going to request messages but this time our current credit should cover it
		messages, err := receiver.ReceiveMessages(context.Background(), 5, nil)
		require.NoError(t, err, "User didn't cancel so cancellation is erased")
		require.Equal(t, []string{"Received message 3"}, getSortedBodies(messages))

		require.Equal(t, uint32(101), amqpReceiver.Credits(), "No new credit needed to be issued - existing credits cover it all")

		require.NoError(t, receiver.Close(context.Background()))
	})
}

func TestReceiver_ReceiveMessages_MessagesArrivingBetweenReceiveMessagesCallsAreReleased(t *testing.T) {
	released := make(chan *amqp.Message, 2)

	amqpReceiver := &internal.FakeAMQPReceiver{
		ReceiveResults: []struct {
			M *amqp.Message
			E error
		}{
			{M: &amqp.Message{Data: [][]byte{[]byte(("received message 1"))}}},
			{M: &amqp.Message{Data: [][]byte{[]byte(("received message 2"))}}},

			// our expectation is that the messageReleaser will pick these up
			{M: &amqp.Message{Data: [][]byte{[]byte(("will be released 1"))}}},
			{M: &amqp.Message{Data: [][]byte{[]byte(("will be released 2"))}}},
		},
		ReleaseMessageFn: func(ctx context.Context, msg *amqp.Message) error {
			select {
			case released <- msg:
			default:
				require.Fail(t, "More messages were released than expected")
			}
			return nil
		},
	}

	args := defaultNewReceiverArgsForTest()

	ns := args.ns.(*internal.FakeNS)
	ns.AMQPLinks = &internal.FakeAMQPLinks{
		Receiver: amqpReceiver,
	}

	receiver, err := newReceiver(args, &ReceiverOptions{
		ReceiveMode: ReceiveModePeekLock,
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"received message 1", "received message 2"}, getSortedBodies(messages))

	// now, we can just wait - the messages will be drained by the background goroutine.
	for i := 0; i < 2; i++ {
		msg := <-released
		require.Equal(t, fmt.Sprintf("will be released %d", i+1), string(msg.Data[0]))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// there aren't any messages in there now since they were all released by
	// the messageReleaser.
	messages, err = receiver.ReceiveMessages(ctx, 1, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Empty(t, messages)

	// closing also shuts down the releaser
	err = receiver.Close(context.Background())
	require.NoError(t, err)
	require.Empty(t, amqpReceiver.ReceiveResults)
	require.Equal(t, 2, amqpReceiver.ReleaseMessageCalled)
}

func defaultNewReceiverArgsForTest() newReceiverArgs {
	return newReceiverArgs{
		entity: entity{
			Queue: "queue",
		},
		ns:                  &internal.FakeNS{},
		cleanupOnClose:      func() {},
		getRecoveryKindFunc: internal.GetRecoveryKind,
		newLinkFn: func(ctx context.Context, session amqpwrap.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
			return nil, nil, nil
		},
		retryOptions: exported.RetryOptions{},
	}
}

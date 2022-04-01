// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestReceiver_ReceiveMessages_AMQPLinksFailure(t *testing.T) {
	fakeAMQPLinks := &internal.FakeAMQPLinks{
		Err: internal.NewErrNonRetriable("failed to create links"),
	}

	receiver := &Receiver{
		amqpLinks:   fakeAMQPLinks,
		receiveMode: ReceiveModePeekLock,
	}

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

			// check that we did attempt to flush the prefetch cache
			require.Equal(t, 1, fakeAMQPReceiver.DrainCalled, "drain is called")
			require.Equal(t, 1, fakeAMQPReceiver.PrefetchedCalled, "prefetched is called")

			// and the links did not need to be closed for a cancellation
			require.Equal(t, 0, fakeAMQPLinks.Closed)
			require.Equal(t, 0, fakeAMQPLinks.CloseIfNeededCalled)
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

				require.Equal(t, 0, fakeAMQPReceiver.DrainCalled, "didn't drain on a broken link/connection")
				require.Equal(t, 1, fakeAMQPReceiver.PrefetchedCalled, "prefetched before throwing away the broken link")

				// a fatal error happened, links should be closed.
				require.Equal(t, 0, fakeAMQPLinks.Closed, "links are closed using CloseIfNeeded")
				require.Equal(t, 1, fakeAMQPLinks.CloseIfNeededCalled, "links are closed on receive errors")
			})
		}
	}
}

func TestReceiver_ReceiveMessages_DrainTimeout_SomeMessagesReceived(t *testing.T) {
	// all the receive modes work the same when there are no messages
	for _, mode := range receiveModesForTests {
		t.Run(mode.Name, func(t *testing.T) {
			fakeAMQPReceiver := &internal.FakeAMQPReceiver{
				ReceiveResults: []struct {
					M *amqp.Message
					E error
				}{
					// we're 1 message short, so we'll need to drain and we're going
					// to make that "hang" so it's cancelled.
					{M: &amqp.Message{Data: [][]byte{[]byte("hello")}}},
					{E: context.DeadlineExceeded},
				},
				DrainCreditImpl: func(ctx context.Context) error {
					// simulate the "drain never comes back" situation.
					// We have a timer now that stops it from hanging forever.
					<-ctx.Done()
					return ctx.Err()
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
			require.Equal(t, 1, fakeAMQPReceiver.DrainCalled, "drain was called")

			// in ReceiveAndDelete mode we make sure we return any messages we've retrieved,
			// even if that means "erasing" the error.
			require.NoError(t, err)
			require.Equal(t, []string{"hello"}, getSortedBodies(messages))

			// since ReceiveAndDelete messages would be lost we always flush any messages that might
			// be sitting in the link before we throw the instance away.
			require.Equal(t, 1, fakeAMQPReceiver.PrefetchedCalled, "prefetched before throwing away the broken link")

			// a drain timeout means we need to close our links as we no longer know the true state of it.
			require.Equal(t, 1, fakeAMQPLinks.Closed, "links are closed using Close(), not CloseIfNeeded()")
			require.Equal(t, 0, fakeAMQPLinks.CloseIfNeededCalled, "links are not closed using CloseIfNeeded()")
		})
	}
}

func TestReceiver_ReceiveMessages_DrainTimeout_NoMessagesReceived(t *testing.T) {
	// When no messages are received we enter into a separate flow where we are deciding
	// if we want to return the error or not.
	for _, mode := range receiveModesForTests {
		t.Run(mode.Name, func(t *testing.T) {
			fakeAMQPReceiver := &internal.FakeAMQPReceiver{
				ReceiveResults: []struct {
					M *amqp.Message
					E error
				}{
					{E: context.DeadlineExceeded},
				},
				DrainCreditImpl: func(ctx context.Context) error {
					// simulate the "drain never comes back" situation.
					// We have a timer now that stops it from hanging forever.
					<-ctx.Done()
					return ctx.Err()
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
			// the timeout occurred in drain. This isn't fatal to the user, even if we received
			// _no_ messages.
			require.NoError(t, err)
			require.Empty(t, messages)

			require.Equal(t, 1, fakeAMQPReceiver.DrainCalled, "drain was called")

			// since ReceiveAndDelete messages would be lost we always flush any messages that might
			// be sitting in the link before we throw the instance away.
			require.Equal(t, 1, fakeAMQPReceiver.PrefetchedCalled, "prefetched before throwing away the broken link")

			// a drain timeout means we need to close our links as we no longer know the true state of it.
			require.Equal(t, 1, fakeAMQPLinks.Closed, "links are closed using Close(), not CloseIfNeeded()")
			require.Equal(t, 0, fakeAMQPLinks.CloseIfNeededCalled, "links are not closed using CloseIfNeeded()")
		})
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

			// no flushing needed, all messages were received in the normal `Receive()` loop.
			require.Equal(t, 0, fakeAMQPReceiver.DrainCalled, "drain is called")
			require.Equal(t, 0, fakeAMQPReceiver.PrefetchedCalled, "prefetched is called")

			// and the links did not need to be closed for a cancellation
			require.Equal(t, 0, fakeAMQPLinks.Closed)
			require.Equal(t, 0, fakeAMQPLinks.CloseIfNeededCalled)
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

	// we should NOT bother flushing the prefetch cache here - we're going to invalidate
	// the link and we don't need to delay returning messages to the customer.
	require.Equal(t, 0, fakeAMQPReceiver.DrainCalled, "didn't drain on a broken link/connection")
	require.Equal(t, 1, fakeAMQPReceiver.PrefetchedCalled, "prefet")

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
		}

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
		}

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
		retryOptions: utils.RetryOptions{
			MaxRetries: 101,
		},
	}))

	require.EqualValues(t, ReceiveModeReceiveAndDelete, receiver.receiveMode)
	path, err = e.String()
	require.NoError(t, err)
	require.EqualValues(t, "topic/Subscriptions/subscription/$Transfer/$DeadLetterQueue", path)
	require.EqualValues(t, 101, receiver.retryOptions.MaxRetries)
}

func TestReceiverDeferUnitTests(t *testing.T) {
	r := &Receiver{
		amqpLinks: &internal.FakeAMQPLinks{
			Err: errors.New("links are dead"),
		},
	}

	messages, err := r.ReceiveDeferredMessages(context.Background(), []int64{1})
	require.EqualError(t, err, "links are dead")
	require.Nil(t, messages)

	r = &Receiver{
		amqpLinks: &internal.FakeAMQPLinks{
			RPC: &badRPCLink{},
		},
	}

	messages, err = r.ReceiveDeferredMessages(context.Background(), []int64{1})
	require.EqualError(t, err, "receive deferred messages failed")
	require.Nil(t, messages)
}

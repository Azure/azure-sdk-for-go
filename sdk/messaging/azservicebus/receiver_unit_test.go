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
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestReceiver_ReceiveMessages_AMQPLinksFailure(t *testing.T) {
	fakeAMQPLinks := &internal.FakeAMQPLinks{
		Err: internal.NewErrNonRetriable("failed to create links"),
	}

	receiver := &Receiver{
		tracer: tracing.NewSpanValidator(t, tracing.SpanMatcher{
			Name:   "Receiver.ReceiveMessages",
			Kind:   tracing.SpanKindConsumer,
			Status: tracing.SpanStatusError,
			Attributes: []tracing.Attribute{
				{Key: tracing.OperationName, Value: "receive"},
				{Key: tracing.OperationType, Value: "receive"},
				{Key: tracing.BatchMessageCount, Value: int64(1)},
			},
		}).NewTracer("module", "version"),
		amqpLinks:   fakeAMQPLinks,
		receiveMode: ReceiveModePeekLock,
		// TODO: need to make this test rely less on stubbing.
		cancelReleaser:    &atomic.Value{},
		maxAllowedCredits: defaultLinkRxBuffer,
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

func ReceiveModeString(mode ReceiveMode) string {
	switch mode {
	case ReceiveModePeekLock:
		return "peekLock"
	case ReceiveModeReceiveAndDelete:
		return "receiveAndDelete"
	default:
		panic(fmt.Sprintf("No string for receive mode %d", mode))
	}
}

func TestReceiverCancellationUnitTests(t *testing.T) {
	t.Run("ImmediatelyCancelled", func(t *testing.T) {
		r := &Receiver{
			tracer: tracing.NewSpanValidator(t, tracing.SpanMatcher{
				Name:   "Receiver.ReceiveMessages",
				Kind:   tracing.SpanKindConsumer,
				Status: tracing.SpanStatusError,
				Attributes: []tracing.Attribute{
					{Key: tracing.OperationName, Value: "receive"},
					{Key: tracing.OperationType, Value: "receive"},
					{Key: tracing.BatchMessageCount, Value: int64(95)},
				},
			}).NewTracer("module", "version"),
			amqpLinks: &internal.FakeAMQPLinks{
				Receiver: &internal.FakeAMQPReceiver{},
			},
			cancelReleaser:    &atomic.Value{},
			maxAllowedCredits: defaultLinkRxBuffer,
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
			tracer: tracing.NewSpanValidator(t, tracing.SpanMatcher{
				Name:   "Receiver.ReceiveMessages",
				Kind:   tracing.SpanKindConsumer,
				Status: tracing.SpanStatusError,
				Attributes: []tracing.Attribute{
					{Key: tracing.OperationName, Value: "receive"},
					{Key: tracing.OperationType, Value: "receive"},
					{Key: tracing.BatchMessageCount, Value: int64(95)},
				},
			}).NewTracer("module", "version"),
			amqpLinks: &internal.FakeAMQPLinks{
				Receiver: &internal.FakeAMQPReceiver{
					ReceiveFn: func(ctx context.Context) (*amqp.Message, error) {
						cancel()
						return nil, ctx.Err()
					},
				},
			},
			cancelReleaser:    &atomic.Value{},
			maxAllowedCredits: defaultLinkRxBuffer,
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
			return nil, &amqp.LinkError{}
		}

		// This is one of the few error types classified as RecoveryKindNone
		// in the releaser this means we'll just retry since the link is still
		// considered good at this point.
		return nil, &amqp.Error{
			Condition: amqp.ErrCond("com.microsoft:server-busy"),
		}
	}

	logsFn := test.CaptureLogsForTest(false)

	releaserFn := receiver.newReleaserFunc(&amqpReceiver)
	releaserFn()

	// we got called a few times, but none of them succeeded.
	require.Equal(t, 2+1, amqpReceiver.ReceiveCalled)

	_ = amqpReceiver.Close(context.Background())

	require.Contains(t,
		logsFn(),
		fmt.Sprintf("[azsb.Receiver] Message releaser stopping because of link failure. Released 0 messages. Will start again after next receive: %s", &amqp.LinkError{}))
}

func TestReceiver_releaserFunc_receiveAndDeleteIsNoop(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	amqpReceiver := mock.NewMockAMQPReceiver(ctrl)

	releaserFn := receiver.newReleaserFunc(amqpReceiver)

	// cancelling is still the empty function
	cancelFn := receiver.cancelReleaser.Load().(func() string)
	require.Equal(t, "empty", cancelFn())

	// in this case you don't need to cancel anything - it's just no-op
	// Note it'll just exit immediately, the "releaser" doesn't block here.
	releaserFn()
}

func TestReceiver_releaserFunc_cancelBetweenReceiveAndReleaseStillReleases(t *testing.T) {
	// Issue: https://github.com/Azure/azure-sdk-for-go/issues/23893
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), nil)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	amqpReceiver := mock.NewMockAMQPReceiver(ctrl)

	releaserFn := receiver.newReleaserFunc(amqpReceiver)
	cancelReleaser := receiver.cancelReleaser.Swap(emptyCancelFn).(func() string)

	amqpReceiver.EXPECT().LinkName().Return("link-name").AnyTimes()
	amqpReceiver.EXPECT().Receive(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// cancelReleaser blocks waiting for the releaser itself to exit. We can't block on waiting for it here or else we'll deadlock.
		go func() { cancelReleaser() }()

		for ctx.Err() == nil {
			time.Sleep(10 * time.Millisecond) // allow the context cancellation to happen
		}

		// simulate that our cancellation occurred _after_ the ReceiveMessages() call
		// NOTE: in the real world this can also happen if the amqp.Receiver is returning prefetched
		// messages since it ignores the context's cancellation state.
		return &amqp.Message{}, nil
	})

	amqpReceiver.EXPECT().ReleaseMessage(gomock.Any(), gomock.Any())

	releaserFn()
}

func TestReceiver_releaserFunc_cancelReceive(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), nil)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	amqpReceiver := mock.NewMockAMQPReceiver(ctrl)

	releaserFn := receiver.newReleaserFunc(amqpReceiver)
	cancelReleaser := receiver.cancelReleaser.Swap(emptyCancelFn).(func() string)

	amqpReceiver.EXPECT().LinkName().Return("link-name").AnyTimes()
	amqpReceiver.EXPECT().Receive(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// cancelReleaser blocks waiting for the releaser itself to exit. We can't block on waiting for it here or else we'll deadlock.
		go func() { cancelReleaser() }()

		for ctx.Err() == nil {
			time.Sleep(10 * time.Millisecond) // allow the context cancellation to happen
		}

		// simulates the cancellation occurring while we were receiving, thus cancelling the call.
		return nil, ctx.Err()
	})

	releaserFn()
}

func TestReceiver_releaserFunc_releaseTimesOut(t *testing.T) {
	receiver, err := newReceiver(defaultNewReceiverArgsForTest(), nil)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	amqpReceiver := mock.NewMockAMQPReceiver(ctrl)

	releaserFn := receiver.newReleaserFunc(amqpReceiver)

	amqpReceiver.EXPECT().LinkName().Return("link-name").AnyTimes()
	amqpReceiver.EXPECT().Receive(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		return &amqp.Message{}, nil
	})

	receiver.defaultReleaserTimeout = time.Millisecond // shorten this so the release call times out properly.

	amqpReceiver.EXPECT().ReleaseMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, msg *amqp.Message) error {
		for ctx.Err() == nil {
			time.Sleep(10 * time.Millisecond) // allow the context cancellation to happen
		}

		return ctx.Err()
	})

	releaserFn()
}

func TestReceiver_fetchMessages_FirstMessageFailure(t *testing.T) {
	errors := []error{&amqp.LinkError{}, context.Canceled}

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
						E: &amqp.LinkError{},
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
			var linkErr *amqp.LinkError
			require.ErrorAs(t, res.Error, &linkErr)

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

func defaultNewReceiverArgsForTest() newReceiverArgs {
	return newReceiverArgs{
		entity: entity{
			Queue: "queue",
		},
		ns:                  &internal.FakeNS{},
		cleanupOnClose:      func() {},
		getRecoveryKindFunc: internal.GetRecoveryKind,
		newLinkFn: func(ctx context.Context, session amqpwrap.AMQPSession) (amqpwrap.AMQPSenderCloser, amqpwrap.AMQPReceiverCloser, error) {
			return nil, nil, nil
		},
		retryOptions: exported.RetryOptions{},
	}
}

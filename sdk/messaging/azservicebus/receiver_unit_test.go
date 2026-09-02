// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
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

type batchDeleteRPCLink struct {
	responses []int32
	sent      []*amqp.Message
	err       error
}

func (l *batchDeleteRPCLink) Close(ctx context.Context) error { return nil }

func (l *batchDeleteRPCLink) RPC(ctx context.Context, msg *amqp.Message) (*amqpwrap.RPCResponse, error) {
	l.sent = append(l.sent, msg)
	if l.err != nil {
		return nil, l.err
	}
	responseIndex := len(l.sent) - 1
	if responseIndex >= len(l.responses) {
		return nil, fmt.Errorf("unexpected batch delete call %d", responseIndex+1)
	}

	return &amqpwrap.RPCResponse{
		Code: 200,
		Message: &amqp.Message{Value: map[string]any{
			"message-count": l.responses[responseIndex],
		}},
	}, nil
}

type countingBatchDeleteLinks struct {
	*internal.FakeAMQPLinks
	getCalls   int
	retryCalls int
	onGet      func()
}

func (l *countingBatchDeleteLinks) Get(ctx context.Context) (*internal.LinksWithID, error) {
	l.getCalls++
	if l.onGet != nil {
		l.onGet()
	}
	return l.FakeAMQPLinks.Get(ctx)
}

func (l *countingBatchDeleteLinks) Retry(ctx context.Context, eventName azlog.Event, operation string, fn internal.RetryWithLinksFn, options exported.RetryOptions) error {
	l.retryCalls++
	links, err := l.Get(ctx)
	if err != nil {
		return err
	}
	return fn(ctx, links, &utils.RetryFnArgs{})
}

func newBatchDeleteReceiver(responses ...int32) (*Receiver, *batchDeleteRPCLink, *countingBatchDeleteLinks) {
	rpcLink := &batchDeleteRPCLink{responses: responses}
	links := &countingBatchDeleteLinks{FakeAMQPLinks: &internal.FakeAMQPLinks{
		Receiver: &internal.FakeAMQPReceiver{},
		RPC:      rpcLink,
	}}
	receiver := &Receiver{
		amqpLinks: links,
	}
	return receiver, rpcLink, links
}

func TestReceiverDeleteMessages(t *testing.T) {
	for _, invalidCount := range []int{0, -1} {
		receiver, rpcLink, links := newBatchDeleteReceiver()

		result, err := receiver.DeleteMessages(context.Background(), invalidCount, nil)
		require.Error(t, err)
		require.Nil(t, result)
		require.Empty(t, rpcLink.sent)
		require.Zero(t, links.getCalls)
	}

	cutoff := time.Date(2026, time.August, 21, 12, 30, 0, 0, time.UTC)
	receiver, rpcLink, links := newBatchDeleteReceiver(37)

	result, err := receiver.DeleteMessages(context.Background(), 50, &DeleteMessagesOptions{
		BeforeEnqueueTime: &cutoff,
	})
	require.NoError(t, err)
	require.EqualValues(t, 37, result.DeletedCount)
	require.Len(t, rpcLink.sent, 1)
	require.Equal(t, 1, links.getCalls)
	require.Equal(t, 1, links.retryCalls)

	premiumReceiver, premiumRPC, _ := newBatchDeleteReceiver(4000)
	premiumResult, err := premiumReceiver.DeleteMessages(context.Background(), 4000, nil)
	require.NoError(t, err)
	require.EqualValues(t, 4000, premiumResult.DeletedCount)
	premiumBody := premiumRPC.sent[0].Value.(map[string]any)
	require.Equal(t, int32(4000), premiumBody["message-count"])
}

func TestReceiverDeleteMessagesDoesNotRedispatchAfterRPCError(t *testing.T) {
	receiver, rpcLink, links := newBatchDeleteReceiver()
	rpcLink.err = fmt.Errorf("response outcome is unknown")

	result, err := receiver.DeleteMessages(context.Background(), 50, nil)
	require.Error(t, err)
	require.Nil(t, result)
	require.Len(t, rpcLink.sent, 1)
	require.Equal(t, 1, links.getCalls)
	require.Equal(t, 1, links.retryCalls)
}

func TestReceiverPurgeMessagesUsesFixedCutoffAndStopsOnZero(t *testing.T) {
	cutoff := time.Date(2026, time.August, 21, 12, 30, 0, 0, time.UTC)
	receiver, rpcLink, links := newBatchDeleteReceiver(500, 2, 0)

	result, err := receiver.PurgeMessages(context.Background(), &PurgeMessagesOptions{
		BeforeEnqueueTime: &cutoff,
	})
	require.NoError(t, err)
	require.EqualValues(t, 502, result.DeletedCount)
	require.Len(t, rpcLink.sent, 3)
	require.Equal(t, 1, links.getCalls)
	require.Equal(t, 1, links.retryCalls)

	for _, sent := range rpcLink.sent {
		body := sent.Value.(map[string]any)
		require.Equal(t, int32(500), body["message-count"])
		require.Equal(t, cutoff, body["enqueued-time-utc"])
	}
}

func TestReceiverPurgeMessagesSupportsPremiumBatchSize(t *testing.T) {
	maxMessageCountPerBatch := 4000
	receiver, rpcLink, links := newBatchDeleteReceiver(4000, 2, 0)

	result, err := receiver.PurgeMessages(context.Background(), &PurgeMessagesOptions{
		MaxMessageCountPerBatch: &maxMessageCountPerBatch,
	})
	require.NoError(t, err)
	require.EqualValues(t, 4002, result.DeletedCount)
	require.Len(t, rpcLink.sent, 3)
	require.Equal(t, 1, links.getCalls)

	for _, sent := range rpcLink.sent {
		body := sent.Value.(map[string]any)
		require.Equal(t, int32(4000), body["message-count"])
	}
}

func TestReceiverPurgeMessagesAllowsServiceToEnforceBatchSizeLimit(t *testing.T) {
	maxMessageCountPerBatch := 4001
	receiver, rpcLink, links := newBatchDeleteReceiver(0)

	result, err := receiver.PurgeMessages(context.Background(), &PurgeMessagesOptions{
		MaxMessageCountPerBatch: &maxMessageCountPerBatch,
	})
	require.NoError(t, err)
	require.Zero(t, result.DeletedCount)
	require.Len(t, rpcLink.sent, 1)
	require.Equal(t, 1, links.getCalls)
	require.Equal(t, int32(4001), rpcLink.sent[0].Value.(map[string]any)["message-count"])
}

func TestReceiverPurgeMessagesRejectsInvalidBatchSizeBeforeLinkSetup(t *testing.T) {
	for _, invalidCount := range []int{0, -1, maxDirectDeleteMessageCount + 1} {
		receiver, rpcLink, links := newBatchDeleteReceiver()

		result, err := receiver.PurgeMessages(context.Background(), &PurgeMessagesOptions{
			MaxMessageCountPerBatch: &invalidCount,
		})
		require.Error(t, err)
		require.Nil(t, result)
		require.Empty(t, rpcLink.sent)
		require.Zero(t, links.getCalls)
	}
}

func TestSessionReceiverDeleteMessagesIncludesSessionID(t *testing.T) {
	receiver, rpcLink, _ := newBatchDeleteReceiver(1)
	sessionID := "session-1"
	sessionReceiver := &SessionReceiver{inner: receiver, sessionID: &sessionID}

	result, err := sessionReceiver.DeleteMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, result.DeletedCount)

	body := rpcLink.sent[0].Value.(map[string]any)
	require.Equal(t, sessionID, body["session-id"])
}

func TestSessionReceiverDeleteMessagesUsesSessionIDResolvedDuringSetup(t *testing.T) {
	receiver, rpcLink, links := newBatchDeleteReceiver(1)
	sessionReceiver := &SessionReceiver{inner: receiver}
	resolvedSessionID := "session-1"
	links.onGet = func() {
		sessionReceiver.sessionID = &resolvedSessionID
	}

	result, err := sessionReceiver.DeleteMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, result.DeletedCount)

	body := rpcLink.sent[0].Value.(map[string]any)
	require.Equal(t, resolvedSessionID, body["session-id"])
}

func TestSessionReceiverPurgeMessagesIncludesSessionID(t *testing.T) {
	receiver, rpcLink, _ := newBatchDeleteReceiver(2, 0)
	sessionID := "session-1"
	sessionReceiver := &SessionReceiver{inner: receiver, sessionID: &sessionID}

	result, err := sessionReceiver.PurgeMessages(context.Background(), nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, result.DeletedCount)
	require.Len(t, rpcLink.sent, 2)
	for _, request := range rpcLink.sent {
		body := request.Value.(map[string]any)
		require.Equal(t, sessionID, body["session-id"])
	}
}

func TestSessionReceiverPurgeMessagesUsesSessionIDResolvedDuringSetup(t *testing.T) {
	receiver, rpcLink, links := newBatchDeleteReceiver(2, 0)
	sessionReceiver := &SessionReceiver{inner: receiver}
	resolvedSessionID := "session-1"
	links.onGet = func() {
		sessionReceiver.sessionID = &resolvedSessionID
	}

	result, err := sessionReceiver.PurgeMessages(context.Background(), nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, result.DeletedCount)

	for _, request := range rpcLink.sent {
		body := request.Value.(map[string]any)
		require.Equal(t, resolvedSessionID, body["session-id"])
	}
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

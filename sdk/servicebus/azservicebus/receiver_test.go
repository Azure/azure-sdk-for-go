package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/require"
)

func TestReceiver(t *testing.T) {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceBusClient, err := NewClient(WithConnectionString(cs))
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)
	cleanupQueue := createQueue(t, cs, queueName)
	defer cleanupQueue()

	t.Run("SendFiveReceiveFive", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)
		defer sender.Close(ctx)

		for i := 0; i < 5; i++ {
			err = sender.SendMessage(ctx, &Message{
				Body: []byte(fmt.Sprintf("[%X,%d]: send five, receive five", nanoSeconds, i)),
			})
			require.NoError(t, err)
		}

		receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(ctx, 5)
		require.NoError(t, err)

		sort.Sort(receivedMessageSlice(messages))

		require.EqualValues(t, 5, len(messages))

		for i := 0; i < 5; i++ {
			require.EqualValues(t,
				fmt.Sprintf("[%X,%d]: send five, receive five", nanoSeconds, i),
				string(messages[i].Body))

			require.NoError(t, receiver.CompleteMessage(ctx, messages[i]))
		}
	})

	t.Run("ForceTimeoutWithTooFewMessages", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)
		defer sender.Close(ctx)

		for i := 0; i < 5; i++ {
			err = sender.SendMessage(ctx, &Message{
				Body: []byte(fmt.Sprintf("[%X,%d]: force timeout waiting for messages", nanoSeconds, i)),
			})
			require.NoError(t, err)
		}

		receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(ctx, 5+1, ReceiveWithMaxWaitTime(time.Second*10))
		require.NoError(t, err)

		sort.Sort(receivedMessageSlice(messages))

		require.EqualValues(t, 5, len(messages))

		for i := 0; i < 5; i++ {
			require.EqualValues(t,
				fmt.Sprintf("[%X,%d]: force timeout waiting for messages", nanoSeconds, i),
				string(messages[i].Body))

			require.NoError(t, receiver.CompleteMessage(ctx, messages[i]))
		}
	})

	t.Run("ReceiveAndAbandon", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)
		defer sender.Close(ctx)

		err = sender.SendMessage(ctx, &Message{
			Body: []byte(fmt.Sprintf("[%X]: send and abandon test", nanoSeconds)),
		})
		require.NoError(t, err)

		receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(ctx, 1)

		require.NoError(t, err)
		require.EqualValues(t, 1, len(messages))

		require.NoError(t, receiver.AbandonMessage(ctx, messages[0]))

		abandonedMessages, err := receiver.ReceiveMessages(ctx, 1)
		require.NoError(t, err)
		require.EqualValues(t, 1, len(abandonedMessages))

		require.NoError(t, receiver.CompleteMessage(ctx, abandonedMessages[0]))
	})

	// Receive has two timeouts - an explicit one (passed in via ReceiveWithMaxTimeout)
	// and an implicit one that kicks in as soon as we receive our first message.
	t.Run("ReceiveWithEarlyFirstMessageTimeout", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)
		defer sender.Close(ctx)

		err = sender.SendMessage(ctx, &Message{
			Body: []byte(fmt.Sprintf("[%X]: send and abandon test", nanoSeconds)),
		})
		require.NoError(t, err)

		receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
		require.NoError(t, err)

		startTime := time.Now()
		messages, err := receiver.ReceiveMessages(ctx, 1,
			ReceiveWithMaxWaitTime(time.Minute*10), // this is never meant to be hit since the first message time is so short.
			ReceiveWithMaxTimeAfterFirstMessage(time.Millisecond))

		require.NoError(t, err)
		require.EqualValues(t, 1, len(messages))

		// `time.Minute` to give some wiggle room for connection initialization
		require.WithinDuration(t, startTime, time.Now(), time.Minute)
	})

	t.Run("SendAndReceiveManyTimes", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)

		defer sender.Close(ctx)

		for i := 0; i < 100; i++ {
			err = sender.SendMessage(ctx, &Message{
				Body: []byte(fmt.Sprintf("[%X:%d]: many messages", nanoSeconds, i)),
			})
			require.NoError(t, err)
		}

		receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
		require.NoError(t, err)

		var allMessages []*ReceivedMessage

		for i := 0; i < 100; i++ {
			messages, err := receiver.ReceiveMessages(ctx, 1, ReceiveWithMaxWaitTime(time.Second*10))
			require.NoError(t, err)
			allMessages = append(allMessages, messages...)

			for _, message := range messages {
				require.NoError(t, receiver.CompleteMessage(ctx, message))
			}
		}

		sort.Sort(receivedMessageSlice(allMessages))

		require.EqualValues(t, len(allMessages), 100)
	})
}

func TestReceiverUnitTests(t *testing.T) {
	t.Run("ReceiverWillNotReopenAfterClose", func(t *testing.T) {
		receiver, err := newReceiver(internal.NewFakeNamespace(), ReceiverWithQueue("queue"))
		require.NoError(t, err)
		require.NoError(t, receiver.Close(context.Background()))

		messages, err := receiver.ReceiveMessages(context.Background(), 1)
		require.Nil(t, messages)

		_, ok := err.(errorinfo.NonRetriable)
		require.True(t, ok, "ErrClosed is a errorinfo.NonRetriable")
		require.ErrorIs(t, err, ErrClosed{"receiver"})
	})

	t.Run("CloseForwardsErrors", func(t *testing.T) {
		ns := internal.NewFakeNamespace()
		ns.NextReceiver = internal.NewFakeLegacyReceiver()

		ns.NextReceiver.CloseImpl = func(ctx context.Context) error {
			return errors.New("Close failed!")
		}

		receiver, err := newReceiver(ns, ReceiverWithSubscription("topic", "subscription"))
		require.NoError(t, err)

		go func() {
			// just cancel out of the listen operation entirely.
			registerEvent := <-ns.NextReceiver.ListenerRegisteredChan
			registerEvent.Cancel()
		}()

		// initializes the internal receiver (it's lazy)
		messages, err := receiver.ReceiveMessages(context.Background(), 1)
		require.EqualError(t, err, context.Canceled.Error())
		require.Empty(t, messages)

		err = receiver.Close(context.Background())
		require.EqualError(t, err, "Close failed!")
	})

	// If an error occurs and we have some messages accumulated in our internal
	// buffer we will still return them to the user.
	//
	// In ReceiveAndDelete _not_ returning these would mean they would be lost - our
	// receiver has the only copy of the message.
	// In PeekLock there is still a chance (if not using sessions, for instance) where
	// the user can still settle messages using the management link as a backup.
	//
	// NOTE: (this is a design item that needs discussion. Just documenting the current behavior)
	t.Run("MessagesAreStillReturnedOnErrors", func(t *testing.T) {
		ns := internal.NewFakeNamespace()
		ns.NextReceiver = internal.NewFakeLegacyReceiver()

		receiver, err := newReceiver(ns,
			ReceiverWithReceiveMode(ReceiveAndDelete),
			ReceiverWithSubscription("topic", "subscription"))
		require.NoError(t, err)

		go func() {
			// just cancel out of the listen operation entirely.
			registerEvent := <-ns.NextReceiver.ListenerRegisteredChan

			// funnel some messages in
			err := registerEvent.Handler.Handle(context.Background(), &internal.Message{
				ID:        "fakeID",
				LockToken: &uuid.UUID{},
				SystemProperties: &internal.SystemProperties{
					SequenceNumber: to.Int64Ptr(1),
				},
			})

			require.NoError(t, err)

			// now cancel the listening.
			registerEvent.Cancel()
		}()

		messages, err := receiver.ReceiveMessages(context.Background(), 2)
		require.EqualError(t, err, context.Canceled.Error())
		require.EqualValues(t, 1, len(messages), "Messages are still returned if we're in ReceiveAndDelete mode")
	})
}

type receivedMessageSlice []*ReceivedMessage

func (messages receivedMessageSlice) Len() int {
	return len(messages)
}

func (messages receivedMessageSlice) Less(i, j int) bool {
	return string(messages[i].Body) < string(messages[j].Body)
}

func (messages receivedMessageSlice) Swap(i, j int) {
	messages[i], messages[j] = messages[j], messages[i]
}

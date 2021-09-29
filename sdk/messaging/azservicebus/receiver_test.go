// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReceiverSendFiveReceiveFive(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("[%d]: send five, receive five", i)),
		})
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 5)
	require.NoError(t, err)

	sort.Sort(receivedMessageSlice(messages))

	require.EqualValues(t, 5, len(messages))

	for i := 0; i < 5; i++ {
		require.EqualValues(t,
			fmt.Sprintf("[%d]: send five, receive five", i),
			string(messages[i].Body))

		require.NoError(t, receiver.CompleteMessage(context.Background(), messages[i]))
	}
}

func TestReceiverForceTimeoutWithTooFewMessages(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	})
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	// there's only one message, requesting more messages will time out.
	messages, err := receiver.ReceiveMessages(context.Background(), 1+1, ReceiveWithMaxWaitTime(10*time.Second))
	require.NoError(t, err)

	require.EqualValues(t,
		[]string{"hello"},
		getSortedBodies(messages))

	require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0]))
}

func TestReceiverAbandon(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("send and abandon test"),
	})
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1)

	require.NoError(t, err)
	require.EqualValues(t, 1, len(messages))

	require.NoError(t, receiver.AbandonMessage(context.Background(), messages[0]))

	abandonedMessages, err := receiver.ReceiveMessages(context.Background(), 1)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(abandonedMessages))

	require.NoError(t, receiver.CompleteMessage(context.Background(), abandonedMessages[0]))
}

// Receive has two timeouts - an explicit one (passed in via ReceiveWithMaxTimeout)
// and an implicit one that kicks in as soon as we receive our first message.
func TestReceiveWithEarlyFirstMessageTimeout(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("send and abandon test"),
	})
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	startTime := time.Now()
	messages, err := receiver.ReceiveMessages(context.Background(), 1,
		ReceiveWithMaxWaitTime(10*time.Minute), // this is never meant to be hit since the first message time is so short.
		ReceiveWithMaxTimeAfterFirstMessage(time.Millisecond))

	require.NoError(t, err)
	require.EqualValues(t, 1, len(messages))

	// `time.Minute` to give some wiggle room for connection initialization
	require.WithinDuration(t, startTime, time.Now(), time.Minute)
}

func TestReceiverSendAndReceiveManyTimes(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)

	defer sender.Close(context.Background())

	for i := 0; i < 100; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("[%d]: many messages", i)),
		})
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	var allMessages []*ReceivedMessage

	for i := 0; i < 100; i++ {
		messages, err := receiver.ReceiveMessages(context.Background(), 1, ReceiveWithMaxWaitTime(10*time.Second))
		require.NoError(t, err)
		allMessages = append(allMessages, messages...)

		for _, message := range messages {
			require.NoError(t, receiver.CompleteMessage(context.Background(), message))
		}
	}

	sort.Sort(receivedMessageSlice(allMessages))

	require.EqualValues(t, len(allMessages), 100)
}

func TestReceiverDeferAndReceiveDeferredMessages(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	ctx := context.TODO()

	defer sender.Close(ctx)

	err = sender.SendMessage(ctx, &Message{
		Body: []byte("deferring a message"),
	})
	require.NoError(t, err)

	receiver, err := client.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 1)
	require.NoError(t, err)

	var sequenceNumbers []int64

	for _, m := range messages {
		err = receiver.DeferMessage(ctx, m)
		require.NoError(t, err)

		sequenceNumbers = append(sequenceNumbers, *m.SequenceNumber)
	}

	deferredMessages, err := receiver.ReceiveDeferredMessages(ctx, sequenceNumbers)
	require.NoError(t, err)

	require.EqualValues(t, []string{"deferring a message"}, getSortedBodies(deferredMessages))
	require.True(t, deferredMessages[0].deferred, "internal flag indicating it was from a deferred receiver method is set")

	for _, m := range deferredMessages {
		err = receiver.CompleteMessage(ctx, m)
		require.NoError(t, err)
	}
}

func TestReceiverPeek(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)

	ctx := context.TODO()

	defer sender.Close(ctx)

	batch, err := sender.NewMessageBatch(ctx)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = batch.Add(&Message{
			Body: []byte(fmt.Sprintf("Message %d", i)),
		})
		require.NoError(t, err)
	}

	err = sender.SendMessage(ctx, batch)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	// wait for a message to show up
	messages, err := receiver.ReceiveMessages(ctx, 1)
	require.NoError(t, err)

	// put them all back
	for _, m := range messages {
		require.NoError(t, receiver.AbandonMessage(ctx, m))
	}

	peekedMessages, err := receiver.PeekMessages(ctx, 2)
	require.NoError(t, err)
	require.EqualValues(t, 2, len(peekedMessages))

	peekedMessages2, err := receiver.PeekMessages(ctx, 2)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(peekedMessages2))

	require.EqualValues(t, []string{
		"Message 0", "Message 1", "Message 2",
	}, getSortedBodies(append(peekedMessages, peekedMessages2...)))
}

func TestReceiverUnitTests(t *testing.T) {
	receiver := &Receiver{}
	require.NoError(t, ReceiverWithSubQueue(SubQueueDeadLetter)(receiver))
	require.EqualValues(t, SubQueueDeadLetter, receiver.config.Entity.Subqueue)

	receiver = &Receiver{}
	require.NoError(t, ReceiverWithSubQueue(SubQueueTransfer)(receiver))
	require.EqualValues(t, SubQueueTransfer, receiver.config.Entity.Subqueue)

	receiver = &Receiver{}
	require.NoError(t, ReceiverWithQueue("queue1")(receiver))
	require.EqualValues(t, "queue1", receiver.config.Entity.Queue)

	receiver = &Receiver{}
	require.NoError(t, ReceiverWithSubscription("topic1", "subscription1")(receiver))
	require.EqualValues(t, "topic1", receiver.config.Entity.Topic)
	require.EqualValues(t, "subscription1", receiver.config.Entity.Subscription)

	receiver = &Receiver{}
	require.NoError(t, ReceiverWithReceiveMode(PeekLock)(receiver))
	require.EqualValues(t, PeekLock, receiver.config.ReceiveMode)

	receiver = &Receiver{}
	require.NoError(t, ReceiverWithReceiveMode(ReceiveAndDelete)(receiver))
	require.EqualValues(t, ReceiveAndDelete, receiver.config.ReceiveMode)

	// If an error occurs and we have some messages accumulated in our internal
	// buffer we will still return them to the user.
	//
	// In ReceiveAndDelete _not_ returning these would mean they would be lost - our
	// receiver has the only copy of the message.
	// In PeekLock there is still a chance (if not using sessions, for instance) where
	// the user can still settle messages using the management link as a backup.
	//
	// NOTE: (this is a design item that needs discussion. Just documenting the current behavior)
	// t.Run("MessagesAreStillReturnedOnErrors", func(t *testing.T) {
	// 	ns := newFakeNamespace()

	// 	ns.Links.Receiver.NextReceiveCalls <- receiveCallResponse{
	// 		message: (&ReceivedMessage{
	// 			Message: Message{
	// 				ID: "fakeID",
	// 			},
	// 			LockToken:      &amqp.UUID{},
	// 			SequenceNumber: to.Int64Ptr(1),
	// 		}).ToAMQPMessage(),
	// 		err: nil,
	// 	}

	// 	receiver, err := newReceiver(ns,
	// 		ReceiverWithReceiveMode(ReceiveAndDelete),
	// 		ReceiverWithSubscription("topic", "subscription"))
	// 	require.NoError(t, err)

	// 	messages, err := receiver.ReceiveMessages(context.Background(), 2)
	// 	require.EqualError(t, err, context.Canceled.Error())
	// 	require.EqualValues(t, 1, len(messages), "Messages are still returned if we're in ReceiveAndDelete mode")
	// })
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

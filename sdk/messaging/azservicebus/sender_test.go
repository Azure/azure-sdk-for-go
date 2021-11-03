// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func Test_Sender_SendBatchOfTwo(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	ctx := context.Background()

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	defer sender.Close(ctx)

	batch, err := sender.NewMessageBatch(ctx, nil)
	require.NoError(t, err)

	added, err := batch.Add(&Message{
		Body: []byte("[0] message in batch"),
	})
	require.NoError(t, err)
	require.True(t, added)

	added, err = batch.Add(&Message{
		Body: []byte("[1] message in batch"),
	})
	require.NoError(t, err)
	require.True(t, added)

	err = sender.SendMessageBatch(ctx, batch)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveModeReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(ctx)

	messages, err := receiver.ReceiveMessages(ctx, 2, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"[0] message in batch", "[1] message in batch"}, getSortedBodies(messages))
}

func Test_Sender_UsingPartitionedQueue(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		EnablePartitioning: to.BoolPtr(true),
	})
	defer cleanup()

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveModeReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		ID:           "message ID",
		Body:         []byte("1. single partitioned message"),
		PartitionKey: to.StringPtr("partitionKey1"),
	})
	require.NoError(t, err)

	batch, err := sender.NewMessageBatch(context.Background(), nil)
	require.NoError(t, err)

	added, err := batch.Add(&Message{
		Body:         []byte("2. Message in batch"),
		PartitionKey: to.StringPtr("partitionKey1"),
	})
	require.NoError(t, err)
	require.True(t, added)

	added, err = batch.Add(&Message{
		Body:         []byte("3. Message in batch"),
		PartitionKey: to.StringPtr("partitionKey1"),
	})
	require.NoError(t, err)
	require.True(t, added)

	err = sender.SendMessageBatch(context.Background(), batch)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1+2, nil)
	require.NoError(t, err)

	sort.Sort(receivedMessages(messages))

	require.EqualValues(t, 3, len(messages))

	require.EqualValues(t, "partitionKey1", *messages[0].PartitionKey)
	require.EqualValues(t, "partitionKey1", *messages[1].PartitionKey)
	require.EqualValues(t, "partitionKey1", *messages[2].PartitionKey)
}

func Test_Sender_SendMessages(t *testing.T) {
	ctx := context.Background()

	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		EnablePartitioning: to.BoolPtr(true),
	})
	defer cleanup()

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveModeReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(context.Background())

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessages(ctx, []*Message{
		{
			Body: []byte("hello"),
		},
		{
			Body: []byte("world"),
		},
	})

	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 2, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"hello", "world"}, getSortedBodies(messages))
}

func Test_Sender_SendMessages_resend(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	ctx := context.Background()

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	peekLockReceiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModePeekLock,
	})
	require.NoError(t, err)

	deletingReceiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	sendAndReceive := func(receiver *Receiver, complete bool) {
		msg := &Message{
			Body: []byte("ResendableMessage"),
			ApplicationProperties: map[string]interface{}{
				"Status": "first send",
			},
		}

		err = sender.SendMessage(ctx, msg)
		require.NoError(t, err)

		message, err := receiver.receiveMessage(ctx, nil)
		require.NoError(t, err)
		require.EqualValues(t, "first send", msg.ApplicationProperties["Status"])
		require.EqualValues(t, "ResendableMessage", string(msg.Body))

		if complete {
			require.NoError(t, receiver.CompleteMessage(ctx, message))
		}

		msg.ApplicationProperties["Status"] = "resend"
		err = sender.SendMessage(ctx, msg)
		require.NoError(t, err)

		message, err = receiver.receiveMessage(ctx, nil)
		require.NoError(t, err)
		require.EqualValues(t, "resend", msg.ApplicationProperties["Status"])
		require.EqualValues(t, "ResendableMessage", string(msg.Body))

		if complete {
			require.NoError(t, receiver.CompleteMessage(ctx, message))
		}
	}

	sendAndReceive(deletingReceiver, false)
	sendAndReceive(peekLockReceiver, true)
}

func Test_Sender_ScheduleMessages(t *testing.T) {
	ctx := context.Background()

	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveModeReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(context.Background())

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	now := time.Now()
	nearFuture := now.Add(20 * time.Second)

	// there are two ways to schedule a message - you can use the
	// `ScheduleMessages` API (in which case you get a sequence number that
	// you can use with CancelScheduledMessage(s)) or you can set the
	// `Scheduled`

	sequenceNumbers, err := sender.ScheduleMessages(ctx,
		[]SendableMessage{
			&Message{Body: []byte("To the future (that will be cancelled!)")},
			&Message{Body: []byte("To the future (not cancelled)")},
		},
		nearFuture)

	require.NoError(t, err)
	require.EqualValues(t, 2, len(sequenceNumbers))

	// cancel one of the ones scheduled using `ScheduleMessages`
	err = sender.CancelScheduledMessages(ctx, []int64{sequenceNumbers[0]})
	require.NoError(t, err)

	err = sender.SendMessage(ctx,
		&Message{
			Body:                 []byte("To the future (scheduled using the field)"),
			ScheduledEnqueueTime: &nearFuture,
		})

	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 2, nil)
	require.NoError(t, err)

	// we cancelled one of the messages so it won't get enqueued (this is the one that survived)
	require.EqualValues(t, []string{"To the future (not cancelled)", "To the future (scheduled using the field)"}, getSortedBodies(messages))

	for _, m := range messages {
		// and the scheduled enqueue time should match what we set pretty closely.
		diff := m.ScheduledEnqueueTime.Sub(nearFuture.UTC())

		// add a little wiggle room, but the scheduled time and the time we set when we scheduled it.
		require.LessOrEqual(t, diff, time.Second, "The requested scheduled time and the actual scheduled time should be close [%s]", m.ScheduledEnqueueTime)
	}
}

func getSortedBodies(messages []*ReceivedMessage) []string {
	sort.Sort(receivedMessages(messages))

	var bodies []string

	for _, msg := range messages {
		bodies = append(bodies, string(msg.Body))
	}

	return bodies
}

type receivedMessages []*ReceivedMessage

func (rm receivedMessages) Len() int {
	return len(rm)
}

// Less compares the messages assuming the .Body field is a valid string.
func (rm receivedMessages) Less(i, j int) bool {
	return string(rm[i].Body) < string(rm[j].Body)
}

func (rm receivedMessages) Swap(i, j int) {
	rm[i], rm[j] = rm[j], rm[i]
}

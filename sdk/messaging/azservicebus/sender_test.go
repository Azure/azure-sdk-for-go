// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func Test_Sender_MessageID(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		EnablePartitioning: to.BoolPtr(true),
	})
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		MessageID: to.StringPtr("message with a message ID"),
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, "message with a message ID", messages[0].MessageID)

	err = sender.SendMessage(context.Background(), &Message{
		// note if you don't explicitly send a message ID one will be auto-generated for you.
	})
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages[0].MessageID) // this is filled in by automatically.
}

func Test_Sender_SendBatchOfTwo(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	ctx := context.Background()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	defer sender.Close(ctx)

	batch, err := sender.NewMessageBatch(ctx, nil)
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		Body: []byte("[0] message in batch"),
	})
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		Body: []byte("[1] message in batch"),
	})
	require.NoError(t, err)

	err = sender.SendMessageBatch(ctx, batch)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveModeReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(ctx)

	messages := receiveAll(t, receiver, 2)
	require.NoError(t, err)

	require.EqualValues(t, []string{"[0] message in batch", "[1] message in batch"}, getSortedBodies(messages))
}

func receiveAll(t *testing.T, receiver *Receiver, expected int) []*ReceivedMessage {
	var all []*ReceivedMessage

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()

	for {
		messages, err := receiver.ReceiveMessages(ctx, 1+2, nil)
		require.NoError(t, err)

		if len(messages) == 0 {
			break
		}

		all = append(all, messages...)

		if len(all) == expected {
			break
		}
	}

	return all
}

func Test_Sender_UsingPartitionedQueue(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		EnablePartitioning: to.BoolPtr(true),
	})
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveModeReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(context.Background())

	batch, err := sender.NewMessageBatch(context.Background(), nil)
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		Body:         []byte("2. Message in batch"),
		PartitionKey: to.StringPtr("partitionKey1"),
	})
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		Body:         []byte("3. Message in batch"),
		PartitionKey: to.StringPtr("partitionKey1"),
	})
	require.NoError(t, err)

	err = sender.SendMessageBatch(context.Background(), batch)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		MessageID:    to.StringPtr("message ID"),
		Body:         []byte("1. single partitioned message"),
		PartitionKey: to.StringPtr("partitionKey1"),
	})
	require.NoError(t, err)

	messages := receiveAll(t, receiver, 3)
	sort.Sort(receivedMessages(messages))

	require.EqualValues(t, 3, len(messages))

	require.EqualValues(t, "partitionKey1", *messages[0].PartitionKey)
	require.EqualValues(t, "partitionKey1", *messages[1].PartitionKey)
	require.EqualValues(t, "partitionKey1", *messages[2].PartitionKey)
}

func Test_Sender_SendMessages_resend(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	ctx := context.Background()

	sender, err := client.NewSender(queueName, nil)
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

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	now := time.Now()
	nearFuture := now.Add(20 * time.Second)

	// there are two ways to schedule a message - you can use the
	// `ScheduleMessages` API (in which case you get a sequence number that
	// you can use with CancelScheduledMessage(s)) or you can set the
	// `Scheduled`

	sequenceNumbers, err := sender.ScheduleMessages(ctx,
		[]*Message{
			{Body: []byte("To the future (that will be cancelled!)")},
			{Body: []byte("To the future (not cancelled)")},
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

func TestSender_SendMessagesDetach(t *testing.T) {
	// NOTE: uncomment this to see some of the background reconnects
	// azlog.SetListener(func(e azlog.Event, s string) {
	// 	log.Printf("%s %s", e, s)
	// })

	sbc, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	sender, err := sbc.NewSender(queueName, nil)
	require.NoError(t, err)

	// make sure the sender link is open and active.
	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("0"),
	})
	require.NoError(t, err)

	// now force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	for i := 1; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("%d", i)),
		})
		require.NoError(t, err)
	}

	receiver, err := sbc.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	// get all the messages
	var all []*ReceivedMessage

	for {
		messages, err := receiver.ReceiveMessages(context.Background(), 5, nil)
		require.NoError(t, err)

		all = append(messages, all...)

		if len(all) == 5 {
			break
		}
	}

	require.EqualValues(t, []string{"0", "1", "2", "3", "4"}, getSortedBodies(all))
}

func TestSender_SendMessageBatchDetach(t *testing.T) {
	// NOTE: uncomment this to see some of the background reconnects
	// azlog.SetListener(func(e azlog.Event, s string) {
	// 	log.Printf("%s %s", e, s)
	// })

	sbc, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	sender, err := sbc.NewSender(queueName, nil)
	require.NoError(t, err)

	// make sure the sender link is open and active.
	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("0"),
	})
	require.NoError(t, err)

	// now force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	for i := 1; i < 5; i++ {
		batch, err := sender.NewMessageBatch(context.Background(), nil)
		require.NoError(t, err)
		require.NoError(t, batch.AddMessage(&Message{
			Body: []byte(fmt.Sprintf("%d", i)),
		}))

		err = sender.SendMessageBatch(context.Background(), batch)
		require.NoError(t, err)
	}

	receiver, err := sbc.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	// get all the messages
	var all []*ReceivedMessage

	for {
		messages, err := receiver.ReceiveMessages(context.Background(), 5, nil)
		require.NoError(t, err)

		all = append(messages, all...)

		if len(all) == 5 {
			break
		}
	}

	require.EqualValues(t, []string{"0", "1", "2", "3", "4"}, getSortedBodies(all))
}

func getSortedBodies(messages []*ReceivedMessage) []string {
	sort.Sort(receivedMessages(messages))

	var bodies []string

	for _, msg := range messages {
		body, err := msg.Body()

		if err != nil {
			panic(err)
		}

		bodies = append(bodies, string(body))
	}

	return bodies
}

type receivedMessages []*ReceivedMessage

func (rm receivedMessages) Len() int {
	return len(rm)
}

// Less compares the messages assuming the .Body field is a valid string.
func (rm receivedMessages) Less(i, j int) bool {
	bodyI, err := rm[i].Body()

	if err != nil {
		panic(err)
	}

	bodyJ, err := rm[j].Body()

	if err != nil {
		panic(err)
	}

	return string(bodyI) < string(bodyJ)
}

func (rm receivedMessages) Swap(i, j int) {
	rm[i], rm[j] = rm[j], rm[i]
}

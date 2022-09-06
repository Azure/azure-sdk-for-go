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
		EnablePartitioning: to.Ptr(true),
	})
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		MessageID: to.Ptr("message with a message ID"),
	}, nil)
	require.NoError(t, err)

	peekedMsg := peekSingleMessageForTest(t, receiver)
	require.EqualValues(t, MessageStateActive, peekedMsg.State)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, "message with a message ID", messages[0].MessageID)

	err = sender.SendMessage(context.Background(), &Message{
		// note if you don't explicitly send a message ID one will be auto-generated for you.
	}, nil)
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
	}, nil)
	require.NoError(t, err)

	err = batch.AddAMQPAnnotatedMessage(&AMQPAnnotatedMessage{
		Body: AMQPAnnotatedMessageBody{
			Data: [][]byte{
				[]byte("[1] message in batch"),
			},
		},
	}, nil)
	require.NoError(t, err)

	err = sender.SendMessageBatch(ctx, batch, nil)
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
		EnablePartitioning: to.Ptr(true),
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
		PartitionKey: to.Ptr("partitionKey1"),
	}, nil)
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		Body:         []byte("3. Message in batch"),
		PartitionKey: to.Ptr("partitionKey1"),
	}, nil)
	require.NoError(t, err)

	err = sender.SendMessageBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		MessageID:    to.Ptr("message ID"),
		Body:         []byte("1. single partitioned message"),
		PartitionKey: to.Ptr("partitionKey1"),
	}, nil)
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
		origSentMsg := &Message{
			Body: []byte("ResendableMessage"),
			ApplicationProperties: map[string]interface{}{
				"Status": "first send",
			},
		}

		err = sender.SendMessage(ctx, origSentMsg, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(ctx, 1, nil)
		require.NoError(t, err)

		require.EqualValues(t, "first send", messages[0].ApplicationProperties["Status"])
		require.EqualValues(t, "ResendableMessage", string(messages[0].Body))

		if complete {
			require.NoError(t, receiver.CompleteMessage(ctx, messages[0], nil))
		}

		messages[0].ApplicationProperties["Status"] = "resend"
		newMsg := messageFromReceivedMessage(t, messages[0])

		err = sender.SendMessage(ctx, newMsg, nil)
		require.NoError(t, err)

		messages, err = receiver.ReceiveMessages(ctx, 1, nil)
		require.NoError(t, err)
		require.EqualValues(t, "resend", messages[0].ApplicationProperties["Status"])
		require.EqualValues(t, "ResendableMessage", string(messages[0].Body))

		if complete {
			require.NoError(t, receiver.CompleteMessage(ctx, messages[0], nil))
		}
	}

	sendAndReceive(deletingReceiver, false)
	sendAndReceive(peekLockReceiver, true)
}

func messageFromReceivedMessage(t *testing.T, receivedMessage *ReceivedMessage) *Message {
	newMsg := &Message{
		MessageID:             &receivedMessage.MessageID,
		ContentType:           receivedMessage.ContentType,
		CorrelationID:         receivedMessage.CorrelationID,
		Body:                  receivedMessage.Body,
		SessionID:             receivedMessage.SessionID,
		Subject:               receivedMessage.Subject,
		ReplyTo:               receivedMessage.ReplyTo,
		ReplyToSessionID:      receivedMessage.ReplyToSessionID,
		To:                    receivedMessage.To,
		TimeToLive:            receivedMessage.TimeToLive,
		PartitionKey:          receivedMessage.PartitionKey,
		ScheduledEnqueueTime:  receivedMessage.ScheduledEnqueueTime,
		ApplicationProperties: receivedMessage.ApplicationProperties,
	}

	return newMsg
}

func Test_Sender_ScheduleAMQPMessages(t *testing.T) {
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
	sequenceNumbers, err := sender.ScheduleAMQPAnnotatedMessages(ctx,
		[]*AMQPAnnotatedMessage{
			{Body: AMQPAnnotatedMessageBody{Data: [][]byte{[]byte("To the future (that will be cancelled!)")}}},
			{Body: AMQPAnnotatedMessageBody{Data: [][]byte{[]byte("To the future (not cancelled)")}}},
		},
		nearFuture, nil)

	require.NoError(t, err)
	require.EqualValues(t, 2, len(sequenceNumbers))

	peekedMsg := peekSingleMessageForTest(t, receiver)
	require.EqualValues(t, MessageStateScheduled, peekedMsg.State)

	// cancel one of the ones scheduled using `ScheduleMessages`
	err = sender.CancelScheduledMessages(ctx, []int64{sequenceNumbers[0]}, nil)
	require.NoError(t, err)

	// this isn't a typical way of doing this, but it's possible to set the field directly
	// rather than using the simpler ScheduleAMQPMessages
	err = sender.SendAMQPAnnotatedMessage(ctx,
		&AMQPAnnotatedMessage{
			Body: AMQPAnnotatedMessageBody{
				Data: [][]byte{[]byte("To the future (scheduled using the field)")},
			},
			MessageAnnotations: map[any]any{
				"x-opt-scheduled-enqueue-time": &nearFuture,
			},
		}, nil)

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
		nearFuture, nil)

	require.NoError(t, err)
	require.EqualValues(t, 2, len(sequenceNumbers))

	peekedMsg := peekSingleMessageForTest(t, receiver)
	require.EqualValues(t, MessageStateScheduled, peekedMsg.State)

	// cancel one of the ones scheduled using `ScheduleMessages`
	err = sender.CancelScheduledMessages(ctx, []int64{sequenceNumbers[0]}, nil)
	require.NoError(t, err)

	err = sender.SendMessage(ctx,
		&Message{
			Body:                 []byte("To the future (scheduled using the field)"),
			ScheduledEnqueueTime: &nearFuture,
		}, nil)

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
	// test.EnableStdoutLogging

	sbc, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	sender, err := sbc.NewSender(queueName, nil)
	require.NoError(t, err)

	// make sure the sender link is open and active.
	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("0"),
	}, nil)
	require.NoError(t, err)

	// now force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	for i := 1; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("%d", i)),
		}, nil)
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
	}, nil)
	require.NoError(t, err)

	// now force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	for i := 1; i < 5; i++ {
		batch, err := sender.NewMessageBatch(context.Background(), nil)
		require.NoError(t, err)
		require.NoError(t, batch.AddMessage(&Message{
			Body: []byte(fmt.Sprintf("%d", i)),
		}, nil))

		err = sender.SendMessageBatch(context.Background(), batch, nil)
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

func TestSender_SendAMQPMessage(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendAMQPAnnotatedMessage(context.Background(), &AMQPAnnotatedMessage{
		Body: AMQPAnnotatedMessageBody{
			Data: [][]byte{
				[]byte("Hello World"),
			},
		},
		ApplicationProperties: map[string]interface{}{
			"hello": "world",
		},
	}, nil)

	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	require.Equal(t, "Hello World", string(messages[0].Body))
	require.Equal(t, "world", messages[0].ApplicationProperties["hello"])

	require.Equal(t, [][]byte{
		[]byte("Hello World"),
	}, messages[0].RawAMQPMessage.Body.Data)

	require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))
}
func TestSender_SendAMQPMessageWithMultipleByteSlicesInData(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendAMQPAnnotatedMessage(context.Background(), &AMQPAnnotatedMessage{
		DeliveryAnnotations: map[any]any{
			"x-opt-delivery-annotation-test": "test-value",
		},
		MessageAnnotations: map[any]any{
			"x-opt-message-annotation-test": "test-value",
		},
		Footer: map[any]any{
			"x-opt-footer-test": "footer-value",
		},
		Header: &AMQPAnnotatedMessageHeader{
			// These flags are just passed through - Service Bus doesn't
			// take any action.
			Priority: 100,
			Durable:  true,
		},
		Body: AMQPAnnotatedMessageBody{
			Data: [][]byte{
				// the defacto azservicebus.Message doesn't allow you to send multiple bytes in
				// the .Data section.
				[]byte("Hello World"),
				[]byte("And another message!"),
			}},
		ApplicationProperties: map[string]interface{}{
			"hello": "world",
		},
	}, nil)

	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	m := messages[0]

	// when the Body is incompatible (in this case because of multiple Data sections)
	// the Body is nil)
	require.Nil(t, m.Body)
	require.NotEmpty(t, m.RawAMQPMessage.DeliveryTag)

	// kill some fields that aren't important for this comparison
	m.RawAMQPMessage.linkName = ""
	m.RawAMQPMessage.inner = nil

	require.Equal(t, &AMQPAnnotatedMessage{
		Header: &AMQPAnnotatedMessageHeader{
			Priority: 100,
			Durable:  true,
		},
		DeliveryAnnotations: map[any]any{
			"x-opt-delivery-annotation-test": "test-value",
			"x-opt-lock-token":               m.RawAMQPMessage.DeliveryAnnotations["x-opt-lock-token"],
		},
		MessageAnnotations: map[any]any{
			"x-opt-message-annotation-test": "test-value",
			"x-opt-enqueued-time":           m.RawAMQPMessage.MessageAnnotations["x-opt-enqueued-time"],
			"x-opt-sequence-number":         m.RawAMQPMessage.MessageAnnotations["x-opt-sequence-number"],
			"x-opt-locked-until":            m.RawAMQPMessage.MessageAnnotations["x-opt-locked-until"],
		},
		Footer: map[any]any{
			"x-opt-footer-test": "footer-value",
		},
		Properties: &AMQPAnnotatedMessageProperties{
			MessageID: m.RawAMQPMessage.Properties.MessageID,
		},
		Body: AMQPAnnotatedMessageBody{Data: [][]byte{
			[]byte("Hello World"),
			[]byte("And another message!"),
		}},
		ApplicationProperties: map[string]any{
			"hello": "world",
		},
		// some items come from the service so we'll just copy them over
		// as we're not trying to test those necessarily
		DeliveryTag: m.RawAMQPMessage.DeliveryTag,
	}, messages[0].RawAMQPMessage)

	require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))
}

func TestSender_SendAMQPMessageWithValue(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	values := []any{
		100.0,
		int64(101),
		"hello world",
		[]byte("hello world as bytes"),
		[]string{"a", "b", "c"},
		nil,
	}

	for i, value := range values {
		err = sender.SendAMQPAnnotatedMessage(context.Background(), &AMQPAnnotatedMessage{
			Body: AMQPAnnotatedMessageBody{Value: value},
			ApplicationProperties: map[string]interface{}{
				"index": i,
			},
		}, nil)
		require.NoError(t, err)
	}

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages := receiveAll(t, receiver, len(values))

	sort.Slice(messages, func(i, j int) bool {
		left := messages[i].ApplicationProperties["index"].(int64)
		right := messages[j].ApplicationProperties["index"].(int64)
		return left < right
	})

	for i, msg := range messages {
		require.Equal(t, values[i], msg.RawAMQPMessage.Body.Value)
		require.NoError(t, receiver.CompleteMessage(context.Background(), msg, nil))
	}
}

func TestSender_SendAMQPMessageWithSequence(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sequences := [][][]any{
		{
			{int64(1), "hello"},
			{int64(2), "world"},
		},
		{},
	}

	for i, seq := range sequences {
		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		err = sender.SendAMQPAnnotatedMessage(context.Background(), &AMQPAnnotatedMessage{
			Body: AMQPAnnotatedMessageBody{Sequence: seq},
			ApplicationProperties: map[string]interface{}{
				"index": i,
			},
		}, nil)

		require.NoError(t, err)
	}

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages := receiveAll(t, receiver, len(sequences))

	sort.Slice(messages, func(i, j int) bool {
		left := messages[i].ApplicationProperties["index"].(int64)
		right := messages[j].ApplicationProperties["index"].(int64)
		return left < right
	})

	for i, msg := range messages {
		if len(sequences[i]) == 0 {
			// this is a bit special - when we send an empty sequence, as an optimization,
			// no value is actually sent for the section.
			require.Nil(t, msg.RawAMQPMessage.Body.Sequence)
			require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))
			continue
		}
		require.Equal(t, sequences[i], msg.RawAMQPMessage.Body.Sequence)
		require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))
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

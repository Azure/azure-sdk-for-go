// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestReceiverCancel(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	messages, err := receiver.ReceiveMessages(ctx, 5, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Empty(t, messages)
}

func TestReceiverSendFiveReceiveFive(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("[%d]: send five, receive five", i)),
		}, nil)
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 5, nil)
	require.NoError(t, err)

	sort.Sort(receivedMessageSlice(messages))

	require.EqualValues(t, 5, len(messages))

	for i := 0; i < 5; i++ {
		body, err := messages[i].Body()
		require.NoError(t, err)

		require.EqualValues(t,
			fmt.Sprintf("[%d]: send five, receive five", i),
			string(body))

		require.NoError(t, receiver.CompleteMessage(context.Background(), messages[i], nil))
	}
}

func TestReceiverForceTimeoutWithTooFewMessages(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// there's only one message, requesting more messages will time out.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	messages, err := receiver.ReceiveMessages(ctx, 1+1, nil)
	require.NoError(t, err)

	require.EqualValues(t,
		[]string{"hello"},
		getSortedBodies(messages))

	require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))
}

func TestReceiverDrainHasError(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// close the link _just_ before the drain credit to simulate connection loss.
	stubReceiver := &StubAMQPReceiver{
		stubDrainCredit: func(inner internal.AMQPReceiverCloser, ctx context.Context) error {
			require.NoError(t, inner.Close(context.Background()))
			return inner.DrainCredit(ctx)
		},
	}

	addStub(t, receiver, stubReceiver)

	// there's only one message, requesting more messages will time out.
	messages, err := receiver.ReceiveMessages(context.Background(), 1+1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"hello"}, getSortedBodies(messages))
}

func TestReceiverAbandon(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("send and abandon test"),
	}, nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)

	require.NoError(t, err)
	require.EqualValues(t, 1, len(messages))

	require.NoError(t, receiver.AbandonMessage(context.Background(), messages[0], nil))

	abandonedMessages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(abandonedMessages))

	require.NoError(t, receiver.CompleteMessage(context.Background(), abandonedMessages[0], nil))
}

// Receive has two timeouts - an explicit one (passed in via ReceiveOptions.MaxWaitTime)
// and an implicit one that kicks in as soon as we receive our first message.
func TestReceiveWithEarlyFirstMessageTimeout(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("send and abandon test"),
	}, nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// this is never meant to be hit since the first message time is so short.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	startTime := time.Now()
	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(messages))

	// `time.Minute` to give some wiggle room for connection initialization
	require.WithinDuration(t, startTime, time.Now(), time.Minute)
}

func TestReceiverSendAndReceiveManyTimes(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)

	defer sender.Close(context.Background())

	for i := 0; i < 100; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("[%d]: many messages", i)),
		}, nil)
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	var allMessages []*ReceivedMessage

	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		messages, err := receiver.ReceiveMessages(ctx, 1, nil)
		require.NoError(t, err)
		allMessages = append(allMessages, messages...)

		for _, message := range messages {
			require.NoError(t, receiver.CompleteMessage(context.Background(), message, nil))
		}
	}

	sort.Sort(receivedMessageSlice(allMessages))

	require.EqualValues(t, len(allMessages), 100)
}

func TestReceiverDeferAndReceiveDeferredMessages(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	ctx := context.TODO()

	defer sender.Close(ctx)

	err = sender.SendMessage(ctx, &Message{
		Body: []byte("deferring a message"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)

	var sequenceNumbers []int64

	for _, m := range messages {
		err = receiver.DeferMessage(ctx, m, nil)
		require.NoError(t, err)

		sequenceNumbers = append(sequenceNumbers, *m.SequenceNumber)
	}

	deferredMessages, err := receiver.ReceiveDeferredMessages(ctx, sequenceNumbers, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"deferring a message"}, getSortedBodies(deferredMessages))
	require.True(t, deferredMessages[0].deferred, "internal flag indicating it was from a deferred receiver method is set")

	for _, m := range deferredMessages {
		err = receiver.CompleteMessage(ctx, m, nil)
		require.NoError(t, err)
	}
}

func TestReceiverDeferWithReceiveAndDelete(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	ctx := context.TODO()

	defer sender.Close(ctx)

	err = sender.SendMessage(ctx, &Message{
		Body: []byte("deferring a message"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)

	var sequenceNumbers []int64

	for _, m := range messages {
		err = receiver.DeferMessage(ctx, m, nil)
		require.NoError(t, err)

		sequenceNumbers = append(sequenceNumbers, *m.SequenceNumber)
	}

	receiveAndDeleteReceiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	messages, err = receiveAndDeleteReceiver.ReceiveDeferredMessages(ctx, sequenceNumbers, nil)
	require.NoError(t, err)
	require.EqualValues(t, len(sequenceNumbers), len(messages))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	messages, err = receiveAndDeleteReceiver.ReceiveMessages(ctx, len(sequenceNumbers), nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Empty(t, messages)
}

func TestReceiverPeek(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)

	ctx := context.TODO()

	defer sender.Close(ctx)

	batch, err := sender.NewMessageBatch(ctx, nil)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err := batch.AddMessage(&Message{
			Body: []byte(fmt.Sprintf("Message %d", i)),
		}, nil)

		require.NoError(t, err)
	}

	err = sender.SendMessageBatch(ctx, batch, nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// wait for a message to show up
	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)

	// put them all back
	for _, m := range messages {
		require.NoError(t, receiver.AbandonMessage(ctx, m, nil))
	}

	peekedMessages, err := receiver.PeekMessages(ctx, 2, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, len(peekedMessages))

	peekedMessages2, err := receiver.PeekMessages(ctx, 2, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(peekedMessages2))

	// peek by seequence using one of our previous messages to prove
	// that we can peek at any arbitrary point in the messages
	require.EqualValues(t, []string{
		"Message 0", "Message 1", "Message 2",
	}, getSortedBodies(append(peekedMessages, peekedMessages2...)))

	repeekedMessages, err := receiver.PeekMessages(ctx, 1, &PeekMessagesOptions{
		FromSequenceNumber: peekedMessages2[0].SequenceNumber,
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, len(repeekedMessages))

	body, err := peekedMessages2[0].Body()
	require.NoError(t, err)

	require.EqualValues(t, []string{
		string(body),
	}, getSortedBodies(repeekedMessages))

	// and peek again (note it won't reset so there'll be "nothing")
	noMessagesExpected, err := receiver.PeekMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.Empty(t, noMessagesExpected)
}

func TestReceiverDetachWithPeekLock(t *testing.T) {
	// NOTE: uncomment this to see some of the background reconnects
	// azlog.SetListener(func(e azlog.Event, s string) {
	// 	log.Printf("%s %s", e, s)
	// })

	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// make sure the receiver link and connection are live.
	_, err = receiver.PeekMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)
	require.NoError(t, sender.Close(context.Background()))

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	var messages []*ReceivedMessage

	for i := 0; i < 5; i++ {
		// depending on how long it takes to rehydrate our links we might
		// have to call multiple times.
		tmpMessages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.NoError(t, err)

		if len(tmpMessages) == 1 {
			messages = tmpMessages
			break
		}
	}

	require.Equal(t, []string{"hello world"}, getSortedBodies(messages))

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	peekedMessages, err := receiver.PeekMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, []string{"hello world"}, getSortedBodies(peekedMessages))

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	// and last, check that the queue is properly empty
	peekedMessages, err = receiver.PeekMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Empty(t, peekedMessages)
}

func TestReceiverDetachWithReceiveAndDelete(t *testing.T) {
	// NOTE: uncomment this to see some of the background reconnects
	// test.EnableStdoutLogging

	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})

	require.NoError(t, err)

	// make sure the receiver link and connection are live.
	_, err = receiver.PeekMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)
	require.NoError(t, sender.Close(context.Background()))

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	var messages []*ReceivedMessage

	for i := 0; i < 5; i++ {
		// depending on how long it takes to rehydrate our links we might
		// have to call multiple times.
		tmpMessages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.NoError(t, err)

		if len(tmpMessages) == 1 {
			// NOTE: in ReceiveAndDelete mode we return any messages we've received, even in the face of connection
			// errors
			messages = tmpMessages
			break
		}
	}

	require.Equal(t, []string{"hello world"}, getSortedBodies(messages))

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{}, nil)
	require.NoError(t, err)

	peekedMessages, err := receiver.PeekMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Empty(t, peekedMessages)
}

func TestReceiver_RenewMessageLock(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)
	lockedUntilOld := messages[0].LockedUntil

	require.NoError(t, receiver.RenewMessageLock(context.Background(), messages[0], nil))

	// these should hopefully be unaffected by clock drift since both values come from
	// the service's times, not ours.
	require.Greater(t, messages[0].LockedUntil.UnixNano(), lockedUntilOld.UnixNano())

	// try renewing a bogus token
	for i := 0; i < len(messages[0].LockToken); i++ {
		messages[0].LockToken[i] = 0
	}

	endCaptureFn := test.CaptureLogsForTest()
	defer endCaptureFn()
	expectedLockBadError := receiver.RenewMessageLock(context.Background(), messages[0], nil)
	// String matching can go away once we fix #15644
	// For now it at least provides the user with good context that something is incorrect about their lock token.
	require.Contains(t, expectedLockBadError.Error(),
		"status code 410 and description: The lock supplied is invalid. Either the lock expired, or the message has already been removed from the queue",
		"error message from SB comes through")

	logMessages := endCaptureFn()

	failedOnFirstTry := false
	for _, msg := range logMessages {
		if strings.HasPrefix(msg, "[azsb.Receiver] (renewMessageLock) Retry attempt 0 returned non-retryable error") {
			failedOnFirstTry = true
		}
	}

	require.True(t, failedOnFirstTry, "No retries attempted for message locks being lost/invalid")
}

// TestReceiverAMQPDataTypes checks that we can send and receive all the AMQP primitive types that are supported
// in ApplicationProperties.
// http://docs.oasis-open.org/amqp/core/v1.0/os/amqp-core-messaging-v1.0-os.html#type-application-properties
//
// > The keys of this map are restricted to be of type string (which excludes the possibility of a null key) and the values
// > are restricted to be of simple types only, that is, excluding map, list, and array types.
func TestReceiverAMQPDataTypes(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	expectedTime, err := time.Parse(time.RFC3339, "2000-01-01T01:02:03Z")
	require.NoError(t, err)

	require.NoError(t, sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello, this is the body"),
		ApplicationProperties: map[string]interface{}{
			// Some primitive types are missing - it's a bit unclear what the right representation of this would be in Go:
			// - TypeCodeDecimal32
			// - TypeCodeDecimal64
			// - TypeCodeDecimal128
			// - TypeCodeChar  (although note below that a 'character' does work, although it's not a TypecodeChar value)
			// https://github.com/Azure/go-amqp/blob/e0c6c63fb01e6642686ee4f8e7412da042bf35dd/internal/encoding/decode.go#L568
			"timestamp": expectedTime,

			"byte":   byte(128),
			"uint8":  int8(101),
			"uint32": int32(400),
			"uint64": int64(400),

			"int":   400,
			"int8":  int8(-101),
			"int32": int32(-400),
			"int64": int64(-400),

			"float":   400.1,
			"float64": float64(400.1),

			"string": "hello world",
			// these aren't "true" chars in the amqp sense - they end up being int32's
			"char":  'g',
			"char2": '❤',

			"bool": true,
			"uuid": amqp.UUID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}),
		},
	}, nil))

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	actualProps := messages[0].ApplicationProperties

	require.Equal(t, map[string]interface{}{
		"timestamp": expectedTime,

		"byte":   byte(128),
		"uint8":  int8(101),
		"uint32": int32(400),
		"uint64": int64(400),

		"int":   int64(400),
		"int8":  int8(-101),
		"int32": int32(-400),
		"int64": int64(-400),

		"float":   float64(400.1),
		"float64": float64(400.1),

		"string": "hello world",
		"char":   'g',
		"char2":  '❤',

		"bool": true,
		"uuid": amqp.UUID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}),
	}, actualProps)
}

func TestReceiverMultiReceiver(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	sender2, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	receiver2, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"hello world"}, getSortedBodies(messages))
	require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))

	err = sender2.SendMessage(context.Background(), &Message{
		Body: []byte("hello world 2"),
	}, nil)
	require.NoError(t, err)

	messages, err = receiver2.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"hello world 2"}, getSortedBodies(messages))
	require.NoError(t, receiver2.CompleteMessage(context.Background(), messages[0], nil))
}

func TestReceiverMultiTopic(t *testing.T) {
	otherQueueName, cleanupOtherQueue := createQueue(t, test.GetConnectionString(t), nil)
	defer cleanupOtherQueue()

	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	queueSender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	queueReceiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	otherQueueSender, err := client.NewSender(otherQueueName, nil)
	require.NoError(t, err)

	otherQueueReceiver, err := client.NewReceiverForQueue(otherQueueName, nil)
	require.NoError(t, err)

	err = queueSender.SendMessage(context.Background(), &Message{
		Body: []byte("sent to queue"),
	}, nil)
	require.NoError(t, err)

	err = otherQueueSender.SendMessage(context.Background(), &Message{
		Body: []byte("sent to other queue"),
	}, nil)
	require.NoError(t, err)

	messages, err := queueReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"sent to queue"}, getSortedBodies(messages))
	require.NoError(t, queueReceiver.CompleteMessage(context.Background(), messages[0], nil))

	otherMessages, err := otherQueueReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"sent to other queue"}, getSortedBodies(otherMessages))
	require.NoError(t, otherQueueReceiver.CompleteMessage(context.Background(), otherMessages[0], nil))

	err = otherQueueSender.SendMessage(context.Background(), &Message{
		Body: []byte("sent to other queue2"),
	}, nil)
	require.NoError(t, err)

	err = queueSender.SendMessage(context.Background(), &Message{
		Body: []byte("sent to queue2"),
	}, nil)
	require.NoError(t, err)

	messages, err = queueReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"sent to queue2"}, getSortedBodies(messages))

	otherMessages, err = otherQueueReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, []string{"sent to other queue2"}, getSortedBodies(otherMessages))
}

type badRPCLink struct {
	internal.RPCLink
}

func (br *badRPCLink) RPC(ctx context.Context, msg *amqp.Message) (*internal.RPCResponse, error) {
	return nil, errors.New("receive deferred messages failed")
}

type receivedMessageSlice []*ReceivedMessage

func (messages receivedMessageSlice) Len() int {
	return len(messages)
}

func (messages receivedMessageSlice) Less(i, j int) bool {
	bodyI, err := messages[i].Body()

	if err != nil {
		panic(err)
	}

	bodyJ, err := messages[j].Body()

	if err != nil {
		panic(err)
	}

	return string(bodyI) < string(bodyJ)
}

func (messages receivedMessageSlice) Swap(i, j int) {
	messages[i], messages[j] = messages[j], messages[i]
}

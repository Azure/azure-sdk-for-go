// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestReceiverBackupSettlement(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			LockDuration: to.Ptr("PT5M"),
		},
	})
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	origReceiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)
	defer test.RequireClose(t, origReceiver)

	otherReceiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)
	defer test.RequireClose(t, otherReceiver)

	messages, err := origReceiver.ReceiveMessages(context.TODO(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	err = otherReceiver.CompleteMessage(context.Background(), messages[0], nil)
	require.NoError(t, err)
}

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
	getLogsFn := test.CaptureLogsForTest(false)

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

	messages := mustReceiveMessages(t, receiver, 5, time.Minute)

	for i := 0; i < 5; i++ {
		require.EqualValues(t,
			fmt.Sprintf("[%d]: send five, receive five", i),
			string(messages[i].Body))

		require.NoError(t, receiver.CompleteMessage(context.Background(), messages[i], nil))
	}

	logs := getLogsFn()
	checkForTokenRefresh(t, logs, queueName)
}

// checkForTokenRefresh just makes sure that background token refresh has been started
// and that we haven't somehow fallen into the trap of marking all tokens are expired.
func checkForTokenRefresh(t *testing.T, logs []string, queueName string) {
	require.NotContains(t, logs, backgroundRenewalDisabledMsg)
	for _, log := range logs {
		if strings.HasPrefix(log, fmt.Sprintf("[azsb.Auth] (%s) next refresh in ", queueName)) {
			return
		}
	}
	require.Fail(t, "No token negotiation log lines")
}

func TestReceiverSendFiveReceiveFive_Subscription(t *testing.T) {
	serviceBusClient, cleanup, topicName, subscriptionName := setupLiveTestWithSubscription(t, &liveTestOptionsWithSubscription{})
	defer cleanup()

	sender, err := serviceBusClient.NewSender(topicName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte(fmt.Sprintf("[%d]: send five, receive five", i)),
		}, nil)
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.NewReceiverForSubscription(topicName, subscriptionName, nil)
	require.NoError(t, err)

	messages := mustReceiveMessages(t, receiver, 5, time.Minute)

	for i := 0; i < 5; i++ {
		require.EqualValues(t,
			fmt.Sprintf("[%d]: send five, receive five", i),
			string(messages[i].Body))

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
	require.True(t, deferredMessages[0].settleOnMgmtLink, "deferred messages should always settle on the management link")

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

	require.EqualValues(t, []string{
		string(peekedMessages2[0].Body),
	}, getSortedBodies(repeekedMessages))

	// and peek again (note it won't reset so there'll be "nothing")
	noMessagesExpected, err := receiver.PeekMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.Empty(t, noMessagesExpected)
}

func TestReceiverDetachWithPeekLock(t *testing.T) {
	// NOTE: uncomment this to see some of the background reconnects
	// stopFn := test.EnableStdoutLogging()
	// defer stopFn()

	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	adminClient := newAdminClientForTest(t, nil)

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

	adminClient := newAdminClientForTest(t, nil)

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

	endCaptureFn := test.CaptureLogsForTest(false)
	defer endCaptureFn()
	expectedLockBadError := receiver.RenewMessageLock(context.Background(), messages[0], nil)

	var asSBError *Error
	require.ErrorAs(t, expectedLockBadError, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)
	require.Contains(t, expectedLockBadError.Error(),
		"status code 410 and description: The lock supplied is invalid. Either the lock expired, or the message has already been removed from the queue",
		"error message from SB comes through")

	logMessages := endCaptureFn()

	failedOnFirstTry := false

	re := regexp.MustCompile(`^\[azsb.Receiver\] \[c:1, l:1, r:name:[^\]]+\] \(renewMessageLock\) Retry attempt 0 returned non-retryable error`)

	for _, msg := range logMessages {
		if re.MatchString(msg) {
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
		ApplicationProperties: map[string]any{
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

	require.Equal(t, map[string]any{
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
	otherQueueName, cleanupOtherQueue := createQueue(t, nil, nil)
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

func TestReceiverMessageLockExpires(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			LockDuration: to.Ptr("PT5S"),
		}})
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{Body: []byte("hello")}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	// sleep so our message locks expire
	time.Sleep(6 * time.Second)

	err = receiver.CompleteMessage(context.Background(), messages[0], nil)

	var asSBError *Error
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeLockLost, asSBError.Code)
}

func TestReceiverUnauthorizedCreds(t *testing.T) {
	allPowerfulCS := test.GetConnectionString(t, test.EnvKeyConnectionString)
	queueName := "testqueue"

	t.Run("ListenOnly with Sender", func(t *testing.T) {
		cs := test.GetConnectionString(t, test.EnvKeyConnectionStringListenOnly)

		client, err := NewClientFromConnectionString(cs, nil) // allowed connection string
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello world"),
		}, nil)

		var sbErr *Error
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "Description: Unauthorized access. 'Send' claim(s) are required to perform this operation")
	})

	t.Run("SenderOnly with Receiver", func(t *testing.T) {
		cs := test.MustGetEnvVar(t, test.EnvKeyConnectionStringSendOnly)

		client, err := NewClientFromConnectionString(cs, nil) // allowed connection string
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.Empty(t, messages)

		var sbErr *Error
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "Description: Unauthorized access. 'Listen' claim(s) are required to perform this operation")
	})

	t.Run("Expired SAS", func(t *testing.T) {
		expiredCS, err := sas.CreateConnectionStringWithSASUsingExpiry(allPowerfulCS, time.Now().Add(-10*time.Minute))
		require.NoError(t, err)

		client, err := NewClientFromConnectionString(expiredCS, nil) // allowed connection string
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello world"),
		}, nil)

		var sbErr *Error
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "rpc: failed, status code 401 and description: ExpiredToken: The token is expired. Expiration time:")

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.Empty(t, messages)

		sbErr = nil
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "rpc: failed, status code 401 and description: ExpiredToken: The token is expired. Expiration time:")
	})

	t.Run("invalid identity creds", func(t *testing.T) {
		identityVars := test.GetIdentityVars(t)
		if identityVars == nil {
			return
		}

		cliCred, err := azidentity.NewClientSecretCredential("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", "bogus-client-secret", nil)
		require.NoError(t, err)

		client, err := NewClient(identityVars.Endpoint, cliCred, nil)
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		var authFailedErr *azidentity.AuthenticationFailedError
		require.ErrorAs(t, err, &authFailedErr)
		require.Empty(t, messages)
	})
}

func TestReceiveAndSendAndReceive(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	scheduledEnqueuedTime := time.Now()

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("body text"),
		ApplicationProperties: map[string]any{
			"hello": "world",
		},
		ContentType:          to.Ptr("application/text"),
		CorrelationID:        to.Ptr("correlation ID"),
		MessageID:            to.Ptr("message id"),
		PartitionKey:         to.Ptr("session id"),
		ReplyTo:              to.Ptr("reply to"),
		ReplyToSessionID:     to.Ptr("reply to session id"),
		ScheduledEnqueueTime: &scheduledEnqueuedTime,
		SessionID:            to.Ptr("session id"),
		Subject:              to.Ptr("subject"),
		TimeToLive:           to.Ptr(time.Minute),
		To:                   to.Ptr("to"),
	}, nil)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	msgs, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.Equal(t, "body text", string(msgs[0].Body))

	// re-send
	err = sender.SendMessage(context.Background(), msgs[0].Message(), nil)
	require.NoError(t, err)

	// re-receive
	rereceivedMsgs, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	require.Equal(t, msgs[0].Message(), rereceivedMsgs[0].Message(), "all sendable fields are preserved when resending")
}

func TestReceiveWithDifferentWaitTime(t *testing.T) {
	setup := func(t *testing.T, timeAfterFirstMessage *time.Duration) int {
		serviceBusClient, cleanup, queueName := setupLiveTest(t, nil)
		defer cleanup()

		sender, err := serviceBusClient.NewSender(queueName, nil)
		require.NoError(t, err)
		defer sender.Close(context.Background())

		batch, err := sender.NewMessageBatch(context.Background(), nil)
		require.NoError(t, err)

		bigBody := make([]byte, 1000)

		// send a bunch of messages
		for i := 0; i < 10000; i++ {
			err := batch.AddMessage(&Message{
				Body: bigBody,
			}, nil)

			if errors.Is(err, ErrMessageTooLarge) {
				err = sender.SendMessageBatch(context.Background(), batch, nil)
				require.NoError(t, err)

				batch, err = sender.NewMessageBatch(context.Background(), nil)
				require.NoError(t, err)

				i--
			}
		}

		if batch.NumMessages() > 0 {
			err = sender.SendMessageBatch(context.Background(), batch, nil)
			require.NoError(t, err)
		}

		receiver, err := serviceBusClient.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		var opts *ReceiveMessagesOptions

		if timeAfterFirstMessage != nil {
			opts = &ReceiveMessagesOptions{
				TimeAfterFirstMessage: *timeAfterFirstMessage,
			}

			t.Logf("Setting time after first message: %s", *timeAfterFirstMessage)
		} else {
			t.Log("Using default time after first message")
		}

		messages, err := receiver.ReceiveMessages(context.Background(), 1000, opts)
		require.NoError(t, err)

		return len(messages)
	}

	base := setup(t, nil)
	require.NotZero(t, base)
	t.Logf("Base case: %d messages", base)

	base2 := setup(t, to.Ptr[time.Duration](0))
	require.NotZero(t, base2)
	t.Logf("Base case2: %d messages", base2)

	bigger := setup(t, to.Ptr(20*time.Second))
	t.Logf("Bigger: %d messages", bigger)
	require.Greater(t, bigger, base)
	require.Greater(t, bigger, base2)
}

func TestReceiveCancelReleaser(t *testing.T) {
	getLogs := test.CaptureLogsForTest(false)

	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			// use a long lock time to really make it clear when we've accidentally
			// orphaned a message.
			LockDuration: to.Ptr("PT5M"),
		},
	})
	defer cleanup()

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	padding := make([]byte, 1)

	var batch *MessageBatch
	const numSent = 2000

	t.Logf("Sending messages")
SendLoop:
	for i := 0; i < numSent; i++ {
		if batch == nil {
			tmpBatch, err := sender.NewMessageBatch(context.Background(), nil)
			require.NoError(t, err)
			batch = tmpBatch
		}

		err := batch.AddMessage(&Message{
			MessageID: to.Ptr(fmt.Sprintf("%d", i)),
			Body:      padding}, nil)

		if errors.Is(err, ErrMessageTooLarge) {
			err := sender.SendMessageBatch(context.Background(), batch, nil)
			require.NoError(t, err)
			batch = nil
			i--
			continue SendLoop
		}

		if i == numSent-1 {
			err := sender.SendMessageBatch(context.Background(), batch, nil)
			require.NoError(t, err)
			break SendLoop
		}
	}

	t.Logf("Receiving small subset of messages")
	messages, err := receiver.ReceiveMessages(context.Background(), numSent, &ReceiveMessagesOptions{
		// Receive with a high credit count, but too little time to actually get them all
		// This will force us into a situation where the AMQP receiver will have a lot of messages
		// in its prefetch cache, giving us a high chance of triggering a previous bug where early
		// cancellation could result in us losing messages.
		TimeAfterFirstMessage: time.Nanosecond,
	})
	require.NoError(t, err)

	ids := &sync.Map{}

	for _, msg := range messages {
		require.NoError(t, receiver.CompleteMessage(context.Background(), msg, nil))
		_, exists := ids.LoadOrStore(msg.MessageID, true)
		require.False(t, exists)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	expected := numSent - len(messages) // remove any messages we've already received. Usually it's just one but it's not
	remaining := expected

	t.Logf("Receiving remaining messages (%d)", remaining)

	for remaining > 0 {
		messages, err := receiver.ReceiveMessages(ctx, remaining, nil)
		require.NoError(t, err)
		require.NotEmpty(t, messages)

		t.Logf("Received %d messages", len(messages))

		wg := sync.WaitGroup{}

		for _, msg := range messages {
			msg := msg
			wg.Add(1)

			go func() {
				defer wg.Done()
				require.NoError(t, receiver.CompleteMessage(context.Background(), msg, nil))
				_, exists := ids.LoadOrStore(msg.MessageID, true)
				require.False(t, exists)
			}()
		}

		wg.Wait()

		remaining -= len(messages)
	}

	count := 0
	ids.Range(func(_, _ any) bool {
		count++
		return true
	})

	require.Equal(t, numSent, count)

	logs := getLogs()

	found := 0

	for _, log := range logs {
		if strings.Contains(log, "Message releaser pausing. Released ") {
			found++
		}
	}

	// This is a bit of a non-deterministic bit so I'm not going to fail the overall test
	if found == 0 {
		t.Logf("Failed to find our 'messages released' log entry: %#v", logs)
	}
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

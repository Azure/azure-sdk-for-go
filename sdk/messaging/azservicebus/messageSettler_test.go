// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azservicebus

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/stretchr/testify/require"
)

func TestMessageSettlementUsingReceiver(t *testing.T) {
	testStuff := newTestStuff(t)
	defer testStuff.Close()

	receiver, deadLetterReceiver := testStuff.Receiver, testStuff.DeadLetterReceiver
	ctx := context.TODO()

	err := testStuff.Sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, messages[0].DeliveryCount)

	// message from queue -> Abandon -> back to the queue
	err = receiver.AbandonMessage(context.Background(), messages[0], &AbandonMessageOptions{
		PropertiesToModify: map[string]any{
			"hello": "world",
		},
	})
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, messages[0].DeliveryCount)
	require.EqualValues(t, "world", messages[0].ApplicationProperties["hello"].(string))

	// message from queue -> DeadLetter -> to the dead letter queue
	err = receiver.DeadLetterMessage(ctx, messages[0], &DeadLetterOptions{
		ErrorDescription: to.Ptr("the error description"),
		Reason:           to.Ptr("the error reason"),
	})
	require.NoError(t, err)

	messages, err = deadLetterReceiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, messages[0].DeliveryCount)

	require.EqualValues(t, "the error description", *messages[0].DeadLetterErrorDescription)
	require.EqualValues(t, "the error reason", *messages[0].DeadLetterReason)

	require.EqualValues(t, *messages[0].ExpiresAt, messages[0].EnqueuedTime.Add(*messages[0].TimeToLive))

	// TODO: introducing deferred messages into the chain seems to have broken something.
	// // message from dead letter queue -> Defer -> to the dead letter queue's deferred messages
	// err = deadLetterReceiver.DeferMessage(ctx, msg)
	// require.NoError(t, err)

	// msg, err = deadLetterReceiver.receiveDeferredMessage(ctx, *msg.SequenceNumber)
	// require.NoError(t, err)

	// deferred message from dead letter queue -> Abandon -> dead letter queue
	err = deadLetterReceiver.AbandonMessage(ctx, messages[0], &AbandonMessageOptions{
		PropertiesToModify: map[string]any{
			"hello": "world",
		},
	})
	require.NoError(t, err)

	messages, err = deadLetterReceiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, messages[0].DeliveryCount)
	require.EqualValues(t, "world", messages[0].ApplicationProperties["hello"].(string))

	// message from dead letter queue -> Complete -> (deleted from queue)
	err = deadLetterReceiver.CompleteMessage(ctx, messages[0], nil)
	require.NoError(t, err)
}

// TestMessageSettlementUsingReceiverWithReceiveAndDelete checks that we don't do anything
// bad if you attempt to settle a message received in ReceiveModeReceiveAndDelete. It should give
// back an error message, but otherwise cause no harm.
func TestMessageSettlementUsingReceiverWithReceiveAndDelete(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, nil)
	defer cleanup()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	ctx := context.Background()

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetRecoveryKind(receiver.AbandonMessage(ctx, messages[0], nil)))
	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetRecoveryKind(receiver.CompleteMessage(ctx, messages[0], nil)))
	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetRecoveryKind(receiver.DeferMessage(ctx, messages[0], nil)))
	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetRecoveryKind(receiver.DeadLetterMessage(ctx, messages[0], nil)))

	require.EqualError(t, receiver.DeadLetterMessage(ctx, messages[0], nil), "messages that are received in `ReceiveModeReceiveAndDelete` mode are not settleable")
}

func TestDeferredMessages(t *testing.T) {
	ctx := context.TODO()

	testStuff := newTestStuff(t)
	defer testStuff.Close()

	receiver := testStuff.Receiver

	t.Run("Abandon", func(t *testing.T) {
		originalDeferredMessage := testStuff.deferMessageForTest(t)

		// abandoning the deferred message will increment its delivery count
		err := receiver.AbandonMessage(ctx, originalDeferredMessage, &AbandonMessageOptions{
			PropertiesToModify: map[string]any{
				"hello": "world",
			},
		})
		require.NoError(t, err)

		// we can peek it without altering anything here.
		peekedMessage := peekSingleMessageForTest(t, receiver)
		require.Equal(t, originalDeferredMessage.DeliveryCount+1, peekedMessage.DeliveryCount, "Delivery count is incremented")
	})

	t.Run("Complete", func(t *testing.T) {
		msg := testStuff.deferMessageForTest(t)

		err := receiver.CompleteMessage(ctx, msg, nil)
		require.NoError(t, err)

		assertEntityEmpty(t, receiver)
	})

	t.Run("Defer", func(t *testing.T) {
		msg := testStuff.deferMessageForTest(t)
		require.EqualValues(t, MessageStateDeferred, msg.State)

		peekedMsg := peekSingleMessageForTest(t, receiver)
		require.EqualValues(t, MessageStateDeferred, peekedMsg.State)

		// double defer!
		err := receiver.DeferMessage(ctx, msg, &DeferMessageOptions{
			PropertiesToModify: map[string]any{
				"hello": "world",
			},
		})
		require.NoError(t, err)

		deferredMessages, err := receiver.ReceiveDeferredMessages(ctx, []int64{*msg.SequenceNumber}, nil)
		require.NoError(t, err)

		require.EqualValues(t, "world", deferredMessages[0].ApplicationProperties["hello"].(string))

		err = receiver.CompleteMessage(ctx, deferredMessages[0], nil)
		require.NoError(t, err)

		assertEntityEmpty(t, receiver)
	})
}

func TestDeferredMessage_DeadLettering(t *testing.T) {
	testStuff := newTestStuff(t)
	defer testStuff.Close()

	receiver, deadLetterReceiver := testStuff.Receiver, testStuff.DeadLetterReceiver

	msg := testStuff.deferMessageForTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := receiver.DeadLetterMessage(ctx, msg, nil)
	require.NoError(t, err)

	// check that the message made it to the dead letter queue
	messages, err := deadLetterReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	// remove it from the DLQ
	require.NoError(t, deadLetterReceiver.CompleteMessage(context.Background(), messages[0], nil))

	// and...everything should be clean
	assertEntityEmpty(t, deadLetterReceiver)
}

func TestMessageSettlementUsingOnlyBackupSettlement(t *testing.T) {
	testStuff := newTestStuff(t)
	defer testStuff.Close()

	actualSettler, _ := testStuff.Receiver.settler.(*messageSettler)
	actualSettler.onlyDoBackupSettlement = true

	actualSettler, _ = testStuff.DeadLetterReceiver.settler.(*messageSettler)
	actualSettler.onlyDoBackupSettlement = true

	receiver, deadLetterReceiver := testStuff.Receiver, testStuff.DeadLetterReceiver
	ctx := context.TODO()

	err := testStuff.Sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	// toggle the super secret switch
	actualSettler, _ = receiver.settler.(*messageSettler)
	actualSettler.onlyDoBackupSettlement = true

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, messages[0].DeliveryCount)

	err = receiver.AbandonMessage(context.Background(), messages[0], nil)
	require.NoError(t, err)

	messages, err = receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, messages[0].DeliveryCount)

	err = receiver.DeadLetterMessage(ctx, messages[0], nil)
	require.NoError(t, err)

	messages, err = deadLetterReceiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, messages[0].DeliveryCount)

	err = deadLetterReceiver.CompleteMessage(context.Background(), messages[0], nil)
	require.NoError(t, err)
}

func TestMessageSettlementWithDeferral(t *testing.T) {
	testStuff := newTestStuff(t)
	defer testStuff.Close()
}

type testStuff struct {
	DeadLetterReceiver *Receiver
	Receiver           *Receiver
	Sender             *Sender
	Require            *require.Assertions
	Client             *Client
	QueueName          string

	cleanup func()
}

func (t *testStuff) Close() {
	t.cleanup()
}

func (t *testStuff) First(messages []*ReceivedMessage, err error) *ReceivedMessage {
	t.Require.NoError(err)
	t.Require.EqualValues([]string{"hello"}, getSortedBodies(messages))
	return messages[0]
}

func newTestStuff(t *testing.T) *testStuff {
	client, cleanup, queueName := setupLiveTest(t, nil)

	testStuff := &testStuff{
		cleanup:   cleanup,
		Require:   require.New(t),
		Client:    client,
		QueueName: queueName,
	}

	var err error
	testStuff.Receiver, err = client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	testStuff.Sender, err = client.NewSender(queueName, nil)
	require.NoError(t, err)

	testStuff.DeadLetterReceiver, err = client.NewReceiverForQueue(
		queueName, &ReceiverOptions{SubQueue: SubQueueDeadLetter})
	require.NoError(t, err)

	return testStuff
}

func assertEntityEmpty(t *testing.T, receiver *Receiver) {
	messages, err := receiver.PeekMessages(context.TODO(), 1, nil)
	require.NoError(t, err)
	require.Empty(t, messages)
}

// deferMessageForTest defers a message with a message body of 'hello'.
func (testStuff *testStuff) deferMessageForTest(t *testing.T) *ReceivedMessage {
	err := testStuff.Sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	messages, err := testStuff.Receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.EqualValues(t, 1, messages[0].DeliveryCount)

	err = testStuff.Receiver.DeferMessage(context.Background(), messages[0], nil)
	require.NoError(t, err)

	messages, err = testStuff.Receiver.ReceiveDeferredMessages(context.Background(), []int64{*messages[0].SequenceNumber}, nil)
	require.NoError(t, err)

	return messages[0]
}

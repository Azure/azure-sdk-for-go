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
	})
	require.NoError(t, err)

	var msg *ReceivedMessage
	msg, err = receiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, msg.DeliveryCount)

	// message from queue -> Abandon -> back to the queue
	err = receiver.AbandonMessage(context.Background(), msg, &AbandonMessageOptions{
		PropertiesToModify: map[string]interface{}{
			"hello": "world",
		},
	})
	require.NoError(t, err)

	msg, err = receiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, msg.DeliveryCount)
	require.EqualValues(t, "world", msg.ApplicationProperties["hello"].(string))

	// message from queue -> DeadLetter -> to the dead letter queue
	err = receiver.DeadLetterMessage(ctx, msg, &DeadLetterOptions{
		ErrorDescription: to.StringPtr("the error description"),
		Reason:           to.StringPtr("the error reason"),
	})
	require.NoError(t, err)

	msg, err = deadLetterReceiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, msg.DeliveryCount)

	require.EqualValues(t, "the error description", *msg.DeadLetterErrorDescription)
	require.EqualValues(t, "the error reason", *msg.DeadLetterReason)

	require.EqualValues(t, *msg.ExpiresAt, msg.EnqueuedTime.Add(*msg.TimeToLive))

	// TODO: introducing deferred messages into the chain seems to have broken something.
	// // message from dead letter queue -> Defer -> to the dead letter queue's deferred messages
	// err = deadLetterReceiver.DeferMessage(ctx, msg)
	// require.NoError(t, err)

	// msg, err = deadLetterReceiver.receiveDeferredMessage(ctx, *msg.SequenceNumber)
	// require.NoError(t, err)

	// deferred message from dead letter queue -> Abandon -> dead letter queue
	err = deadLetterReceiver.AbandonMessage(ctx, msg, &AbandonMessageOptions{
		PropertiesToModify: map[string]interface{}{
			"hello": "world",
		},
	})
	require.NoError(t, err)

	msg, err = deadLetterReceiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, msg.DeliveryCount)
	require.EqualValues(t, "world", msg.ApplicationProperties["hello"].(string))

	// message from dead letter queue -> Complete -> (deleted from queue)
	err = deadLetterReceiver.CompleteMessage(ctx, msg)
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
	})
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queueName, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetSBErrInfo(receiver.AbandonMessage(ctx, messages[0], nil)).RecoveryKind)
	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetSBErrInfo(receiver.CompleteMessage(ctx, messages[0])).RecoveryKind)
	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetSBErrInfo(receiver.DeferMessage(ctx, messages[0], nil)).RecoveryKind)
	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetSBErrInfo(receiver.DeadLetterMessage(ctx, messages[0], nil)).RecoveryKind)

	require.EqualError(t, receiver.DeadLetterMessage(ctx, messages[0], nil), "messages that are received in `ReceiveModeReceiveAndDelete` mode are not settleable")
}

func TestDeferredMessages(t *testing.T) {
	ctx := context.TODO()

	testStuff := newTestStuff(t)
	defer testStuff.Close()

	receiver := testStuff.Receiver

	t.Run("Abandon", func(t *testing.T) {
		t.Skip("This test is currently broken, https://github.com/Azure/azure-sdk-for-go/issues/15626")

		msg := testStuff.deferMessageForTest(t)

		// abandon the deferred message, which should return
		// it to the queue.
		err := receiver.AbandonMessage(ctx, msg, &AbandonMessageOptions{
			PropertiesToModify: map[string]interface{}{
				"hello": "world",
			},
		})
		require.NoError(t, err)

		// BUG: we're timing out here, even though our abandon should have put the message
		// back into the queue. It appears that settlement methods don't work on messages
		// that have been received as deferred.
		msg, err = receiver.receiveMessage(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, msg)
		require.EqualValues(t, "world", msg.ApplicationProperties["hello"].(string))
	})

	t.Run("Complete", func(t *testing.T) {
		msg := testStuff.deferMessageForTest(t)

		err := receiver.CompleteMessage(ctx, msg)
		require.NoError(t, err)

		assertEntityEmpty(t, receiver)
	})

	t.Run("Defer", func(t *testing.T) {
		msg := testStuff.deferMessageForTest(t)

		// double defer!
		err := receiver.DeferMessage(ctx, msg, &DeferMessageOptions{
			PropertiesToModify: map[string]interface{}{
				"hello": "world",
			},
		})
		require.NoError(t, err)

		msg, err = receiver.receiveDeferredMessage(ctx, *msg.SequenceNumber)
		require.NoError(t, err)

		require.EqualValues(t, "world", msg.ApplicationProperties["hello"].(string))

		err = receiver.CompleteMessage(ctx, msg)
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
	msg, err = deadLetterReceiver.receiveMessage(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, msg)

	// remove it from the DLQ
	require.NoError(t, deadLetterReceiver.CompleteMessage(context.Background(), msg))

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
	})
	require.NoError(t, err)

	// toggle the super secret switch
	actualSettler, _ = receiver.settler.(*messageSettler)
	actualSettler.onlyDoBackupSettlement = true

	var msg *ReceivedMessage
	msg, err = receiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 1, msg.DeliveryCount)

	err = receiver.AbandonMessage(context.Background(), msg, nil)
	require.NoError(t, err)

	msg, err = receiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, msg.DeliveryCount)

	err = receiver.DeadLetterMessage(ctx, msg, nil)
	require.NoError(t, err)

	msg, err = deadLetterReceiver.receiveMessage(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, 2, msg.DeliveryCount)

	err = deadLetterReceiver.CompleteMessage(context.Background(), msg)
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

func (testStuff *testStuff) deferMessageForTest(t *testing.T) *ReceivedMessage {
	err := testStuff.Sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	})
	require.NoError(t, err)

	var msg *ReceivedMessage
	msg, err = testStuff.Receiver.receiveMessage(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, msg)

	require.EqualValues(t, 1, msg.DeliveryCount)

	err = testStuff.Receiver.DeferMessage(context.Background(), msg, nil)
	require.NoError(t, err)

	msg, err = testStuff.Receiver.receiveDeferredMessage(context.Background(), *msg.SequenceNumber)
	require.NoError(t, err)

	return msg
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azservicebus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageSettlementUsingReceiver(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	receiver, err := client.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	})
	require.NoError(t, err)

	getOneMessage := func() *ReceivedMessage {
		messages, err := receiver.ReceiveMessages(context.Background(), 1)
		require.NoError(t, err)

		bodies := getSortedBodies(messages)
		require.EqualValues(t, []string{"hello"}, bodies)

		return messages[0]
	}

	firstMessage := getOneMessage()
	require.EqualValues(t, 1, firstMessage.DeliveryCount)

	// bounce it back to the queue
	require.NoError(t, receiver.AbandonMessage(context.Background(), firstMessage))

	// get it again
	afterAbandon := getOneMessage()
	require.EqualValues(t, 2, afterAbandon.DeliveryCount)

	// now complete it
	require.NoError(t, receiver.CompleteMessage(context.Background(), afterAbandon))
}

func TestMessageSettlementUsingManagementLink(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t)
	defer cleanup()

	receiver, err := client.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	// toggle the super secret switch
	receiver.settler.onlyDoBackupSettlement = true

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello"),
	})
	require.NoError(t, err)

	getOneMessage := func() *ReceivedMessage {
		messages, err := receiver.ReceiveMessages(context.Background(), 1)
		require.NoError(t, err)

		bodies := getSortedBodies(messages)
		require.EqualValues(t, []string{"hello"}, bodies)

		return messages[0]
	}

	firstMessage := getOneMessage()
	require.EqualValues(t, 1, firstMessage.DeliveryCount)

	// bounce it back to the queue
	require.NoError(t, receiver.AbandonMessage(context.Background(), firstMessage))

	// get it again
	afterAbandon := getOneMessage()
	require.EqualValues(t, 2, afterAbandon.DeliveryCount)

	// now complete it
	require.NoError(t, receiver.CompleteMessage(context.Background(), afterAbandon))
}

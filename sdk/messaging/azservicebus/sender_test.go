// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"sort"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
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

	err = sender.SendMessage(ctx, batch)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveAndDelete})
	require.NoError(t, err)
	defer receiver.Close(ctx)

	messages, err := receiver.ReceiveMessages(ctx, 2, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"[0] message in batch", "[1] message in batch"}, getSortedBodies(messages))
}

func Test_Sender_UsingPartitionedQueue(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &internal.QueueDescription{
		EnablePartitioning: to.BoolPtr(true),
	})
	defer cleanup()

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	receiver, err := client.NewReceiverForQueue(
		queueName, &ReceiverOptions{ReceiveMode: ReceiveAndDelete})
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

	err = sender.SendMessage(context.Background(), batch)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1+2, nil)
	require.NoError(t, err)

	sort.Sort(receivedMessages(messages))

	require.EqualValues(t, 3, len(messages))

	require.EqualValues(t, "partitionKey1", *messages[0].PartitionKey)
	require.EqualValues(t, "partitionKey1", *messages[1].PartitionKey)
	require.EqualValues(t, "partitionKey1", *messages[2].PartitionKey)
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

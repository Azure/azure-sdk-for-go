package azservicebus

import (
	"context"
	"sort"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/require"
)

func TestSender(t *testing.T) {
	cs := getConnectionString(t)

	serviceBusClient, err := NewClient(WithConnectionString(cs))
	require.NoError(t, err)

	t.Run("testSendBatchOfTwo", func(t *testing.T) {
		queueName, cleanupQueue := createQueue(t, cs, nil)
		defer cleanupQueue()
		testSendBatchOfTwo(context.Background(), t, serviceBusClient, queueName)
	})

	t.Run("testUsingPartitionedQueue", func(t *testing.T) {
		queueName, cleanupQueue := createQueue(t, cs, &internal.QueueDescription{
			EnablePartitioning: to.BoolPtr(true),
		})
		defer cleanupQueue()

		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)
		defer sender.Close(context.Background())

		receiver, err := serviceBusClient.NewReceiver(
			ReceiverWithQueue(queueName),
			ReceiverWithReceiveMode(ReceiveAndDelete))
		require.NoError(t, err)
		defer receiver.Close(context.Background())

		err = sender.SendMessage(context.Background(), &Message{
			ID:           "message ID",
			Body:         []byte("1. single partitioned message"),
			PartitionKey: to.StringPtr("partitionKey1"),
		})
		require.NoError(t, err)

		batch, err := sender.NewMessageBatch(context.Background())
		require.NoError(t, err)

		err = batch.Add(&Message{
			Body:         []byte("2. Message in batch"),
			PartitionKey: to.StringPtr("partitionKey1"),
		})
		require.NoError(t, err)

		err = batch.Add(&Message{
			Body:         []byte("3. Message in batch"),
			PartitionKey: to.StringPtr("partitionKey1"),
		})
		require.NoError(t, err)

		err = sender.SendMessage(context.Background(), batch)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1+2)
		require.NoError(t, err)

		sort.Sort(receivedMessages(messages))

		require.EqualValues(t, 3, len(messages))

		require.EqualValues(t, "partitionKey1", *messages[0].PartitionKey)
		require.EqualValues(t, "partitionKey1", *messages[1].PartitionKey)
		require.EqualValues(t, "partitionKey1", *messages[2].PartitionKey)
	})
}

func TestSenderUnitTests(t *testing.T) {

}

func testSendBatchOfTwo(ctx context.Context, t *testing.T, serviceBusClient *Client, queueName string) {
	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)

	defer sender.Close(ctx)

	batch, err := sender.NewMessageBatch(ctx)
	require.NoError(t, err)

	err = batch.Add(&Message{
		Body: []byte("[0] message in batch"),
	})
	require.NoError(t, err)

	err = batch.Add(&Message{
		Body: []byte("[1] message in batch"),
	})
	require.NoError(t, err)

	err = sender.SendMessage(ctx, batch)
	require.NoError(t, err)

	receiver, err := serviceBusClient.NewReceiver(
		ReceiverWithQueue(queueName),
		ReceiverWithReceiveMode(ReceiveAndDelete))
	require.NoError(t, err)
	defer receiver.Close(ctx)

	messages, err := receiver.ReceiveMessages(ctx, 2)
	require.NoError(t, err)

	require.EqualValues(t, []string{"[0] message in batch", "[1] message in batch"}, getSortedBodies(messages))
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

package azservicebus

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/require"
)

func TestSender(t *testing.T) {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceBusClient, err := NewClient(WithConnectionString(cs))
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

	cleanupQueue := createQueue(t, cs, queueName)
	defer cleanupQueue()

	t.Run("testSendBatchOfTwo", func(t *testing.T) {
		testSendBatchOfTwo(ctx, t, serviceBusClient, queueName)
	})

	t.Run("testUsingPartitionedQueue", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender("partitionedQueue")
		require.NoError(t, err)
		defer sender.Close(ctx)

		receiver, err := serviceBusClient.NewReceiver(
			ReceiverWithQueue("partitionedQueue"),
			ReceiverWithReceiveMode(ReceiveAndDelete))
		require.NoError(t, err)
		defer receiver.Close(ctx)

		err = sender.SendMessage(ctx, &Message{
			Body:         []byte("1. single partitioned message"),
			PartitionKey: to.StringPtr("partitionKey1"),
		})
		require.NoError(t, err)

		batch, err := sender.NewMessageBatch(ctx)
		require.NoError(t, err)

		err = batch.Add(&Message{
			Body:         []byte("2. Message in batch"),
			PartitionKey: to.StringPtr("partitionKey1"),
		})
		require.NoError(t, err)

		err = batch.Add(&Message{
			Body:         []byte("3. Message in batch"),
			PartitionKey: to.StringPtr("partitionKey that gets ignored because first message in the batch wins"),
		})
		require.NoError(t, err)

		err = sender.SendMessage(ctx, batch)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(ctx, 1+2)
		require.NoError(t, err)

		sort.Sort(receivedMessages(messages))

		require.EqualValues(t, 3, len(messages))

		require.EqualValues(t, "partitionKey1", messages[0].PartitionKey)
		// (used the same partition key for the _first_ message in the batch, which was then also applied
		// to every message in the batch (ie, it override the second message's partition key)
		require.EqualValues(t, "partitionKey1", messages[1].PartitionKey)
		require.EqualValues(t, "partitionKey1", messages[2].PartitionKey)
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

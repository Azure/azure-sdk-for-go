package azservicebus

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

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

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)

	defer sender.Close(ctx)

	batch, err := sender.CreateMessageBatch(ctx)
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

	// we sent a single batch with two messages in it
	messages, err := receiver.ReceiveMessages(ctx, 2)
	require.NoError(t, err)

	require.EqualValues(t, []string{"[0] message in batch", "[1] message in batch"}, getSortedBodies(messages))
}

func getSortedBodies(messages []*ReceivedMessage) []string {
	var bodies []string

	for _, msg := range messages {
		bodies = append(bodies, string(msg.Body))
	}

	sort.Strings(bodies)
	return bodies
}

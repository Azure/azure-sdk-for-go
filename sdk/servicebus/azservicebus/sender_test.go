package azservicebus

import (
	"context"
	"fmt"
	"os"
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

	batch, err := sender.CreateMessageBatch(ctx)
	require.NoError(t, err)

	err = batch.Add(&Message{
		Body: []byte("hello world"),
	})
	require.NoError(t, err)

	err = sender.SendMessage(ctx, batch)
	require.NoError(t, err)
}

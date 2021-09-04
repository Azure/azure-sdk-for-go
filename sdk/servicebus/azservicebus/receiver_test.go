package azservicebus

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestReceiverReceiveMessagesExactly(t *testing.T) {
	godotenv.Load()
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceBusClient, err := NewServiceBusClient(ServiceBusWithConnectionString(cs))
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)
	cleanupQueue := createQueue(t, cs, queueName)
	defer cleanupQueue()

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(ctx, &Message{
			Body: []byte(fmt.Sprintf("hello %d", i)),
		})
	}

	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 5)
	require.NoError(t, err)

	require.EqualValues(t, 5, len(messages))
}

func TestReceiverReceiveFewerMessages(t *testing.T) {
	godotenv.Load()
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceBusClient, err := NewServiceBusClient(ServiceBusWithConnectionString(cs))
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)
	cleanupQueue := createQueue(t, cs, queueName)
	defer cleanupQueue()

	receiver, err := serviceBusClient.NewReceiver(ReceiverWithQueue(queueName))
	require.NoError(t, err)

	sender, err := serviceBusClient.NewSender(queueName)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(ctx, &Message{
			Body: []byte(fmt.Sprintf("hello %d", i)),
		})
	}

	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(ctx, 5+1, ReceiveWithMaxWaitTime(time.Second*10))
	require.NoError(t, err)

	require.EqualValues(t, 5, len(messages))
}

// TODO: there are some finer points of how receive messages works:

func TestReceiverReceiveMessagesWithFirstMessageTimer(t *testing.T) {
	// 1. We don't actually wait for the entirety of the maxWaitTime. There is a second "faster" timer
	//    that kicks in when the first message arrives. This favors immediacy.
	t.Fail()
}

func TestReceiverReceiveMessagesWithFailure(t *testing.T) {
	// 2. There are some cases that we could consider 'partial success'. There are some idiomacy issues
	//    to consider (ie, should we return a value _and_ an error) but in some cases, like when we
	//    receive messages in ReceiveAndDelete mode
	t.Fail()
}

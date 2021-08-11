package azservicebus

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func createQueue(ctx context.Context, t *testing.T, connectionString string) (string, func()) {
	ns, err := internal.NewNamespace(internal.NamespaceWithConnectionString(connectionString))
	require.NoError(t, err)

	qm := ns.NewQueueManager()

	// generate random queue
	queueName := fmt.Sprintf("test-%X", rand.Int63())

	_, err = qm.Put(ctx, queueName)
	require.NoError(t, err)

	return queueName, func() {
		require.NoError(t, qm.Delete(ctx, queueName))
	}
}

func TestProcessor(t *testing.T) {
	godotenv.Load()
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceBusClient, err := NewServiceBusClient(ServiceBusWithConnectionString(cs))
	require.NoError(t, err)

	// queueName, cleanup := createQueue(ctx, t, cs)
	// defer cleanup()
	queueName := "test1"

	t.Run("Receive messages using processor", func(t *testing.T) {
		sender, err := serviceBusClient.NewSender(queueName)
		require.NoError(t, err)

		err = sender.SendMessage(ctx, &Message{
			Body: []byte("hello world"),
		})

		require.NoError(t, err)

		processor, err := serviceBusClient.NewProcessor(ProcessorWithQueue(queueName))
		require.NoError(t, err)

		defer processor.Close(ctx)

		messagesCh := make(chan *ReceivedMessage, 1)

		err = processor.Start(func(message *ReceivedMessage) error {
			select {
			case messagesCh <- message:
				break
			default:
				return fmt.Errorf("More messages than expected")
			}
			return nil
		}, func(err error) {
			if err == context.Canceled {
				return
			}

			require.NoError(t, err)
		})

		require.NoError(t, err)

		// wait for a period of time, but let's be reasonable
		select {
		case message := <-messagesCh:
			require.EqualValues(t, "hello world", string(message.Body))
		case <-processor.Done():
			t.Fatal("Processor was closed before messages arrived")
			break
		case <-ctx.Done():
			t.Fatal("Test finished before any messages arrived")
			break
		}
	})

	t.Run("Stop waits for pending callbacks", func(t *testing.T) {
	})

	t.Run("Close waits for pending callbacks", func(t *testing.T) {
	})

	t.Run("Stop prevents new credits from being added to the link", func(t *testing.T) {
		// TODO: without this a user can call Stop() on the processor but, internally, we keep receiving messages
		// and keep forwarding them.
	})

	t.Run("Stop pauses the links but not close them", func(t *testing.T) {
		// Stop() should just 'drain' the links, and not outright close them. There are some operations
		// (with sessions, for instance) that will require the original link for settlement.
		//
		// For non-sessions we can fall back to just using the management link to settle.
	})

	t.Run("autoComplete (basic - complete on success)", func(t *testing.T) {
	})

	t.Run("autoComplete (basic - abandon on error)", func(t *testing.T) {
	})

	t.Run("autoComplete is a no-op if you manually settled", func(t *testing.T) {
	})
}

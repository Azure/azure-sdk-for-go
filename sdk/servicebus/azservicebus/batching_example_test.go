package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func Example_batchingMessages() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	qm := ns.NewQueueManager()
	qe, err := ensureQueue(ctx, qm, "MessageBatchingExample")
	if err != nil {
		fmt.Println(err)
		return
	}

	q, err := ns.NewQueue(qe.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = q.Close(ctx)
	}()

	msgs := make([]*servicebus.Message, 10)
	for i := 0; i < 10; i++ {
		msgs[i] = servicebus.NewMessageFromString(fmt.Sprintf("foo %d", i))
	}

	batcher := servicebus.NewMessageBatchIterator(servicebus.StandardMaxMessageSizeInBytes, msgs...)
	if err := q.SendBatch(ctx, batcher); err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 10; i++ {
		err := q.ReceiveOne(ctx, MessagePrinter{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Output:
	// foo 0
	// foo 1
	// foo 2
	// foo 3
	// foo 4
	// foo 5
	// foo 6
	// foo 7
	// foo 8
	// foo 9
}

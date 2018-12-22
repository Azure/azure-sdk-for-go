package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-amqp-common-go/uuid"

	"github.com/Azure/azure-service-bus-go"
)

func Example_duplicateMessageDetection() {
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

	window := 30 * time.Second
	qm := ns.NewQueueManager()
	qe, err := ensureQueue(ctx, qm, "DuplicateDetectionExample", servicebus.QueueEntityWithDuplicateDetection(&window))
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

	guid, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := servicebus.NewMessageFromString("foo")
	msg.ID = guid.String()

	// send the message twice with the same ID
	for i := 0; i < 2; i++ {
		if err := q.Send(ctx, msg); err != nil {
			fmt.Println(err)
			return
		}
	}

	// there should be only 1 message received from the queue
	go func() {
		if err := q.Receive(ctx, MessagePrinter{}); err != nil {
			if err.Error() != "context canceled" {
				fmt.Println(err)
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)

	// Output:
	// foo
}

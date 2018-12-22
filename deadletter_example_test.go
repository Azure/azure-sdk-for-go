package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func Example_deadletterQueues() {
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
	qe, err := ensureQueue(ctx, qm, "DeadletterExample")
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

	if err := q.Send(ctx, servicebus.NewMessageFromString("foo")); err != nil {
		fmt.Println(err)
		return
	}

	// Abandon the message 10 times simulating attempting to process the message 10 times. After the 10th time, the
	// message will be placed in the Deadletter Queue.
	for count := 0; count < 10; count++ {
		err = q.ReceiveOne(ctx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			fmt.Printf("count: %d\n", count+1)
			return msg.Abandon(ctx)
		}))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// receive one from the queue's deadletter queue. It should be the foo message.
	qdl := q.NewDeadLetter()
	if err := qdl.ReceiveOne(ctx, MessagePrinter{}); err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = qdl.Close(ctx)
	}()

	// Output:
	// count: 1
	// count: 2
	// count: 3
	// count: 4
	// count: 5
	// count: 6
	// count: 7
	// count: 8
	// count: 9
	// count: 10
	// foo
}

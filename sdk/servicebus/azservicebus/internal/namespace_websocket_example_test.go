package internal_test

import (
	"context"
	"fmt"
	"os"
	"time"

	servicebus "github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus"
)

func ExampleNamespaceWithWebSocket() {
	const queueName = "wssQueue"

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	// Create a Service Bus Namespace using a connection string over wss:// on port 443
	ns, err := servicebus.NewNamespace(
		servicebus.NamespaceWithConnectionString(connStr),
		servicebus.NamespaceWithWebSocket(),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qm := ns.NewQueueManager()
	if _, err := ensureQueue(ctx, qm, queueName); err != nil {
		fmt.Println(err)
		return
	}

	client, err := ns.NewQueue(queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send a message to the queue
	if err := client.Send(ctx, servicebus.NewMessageFromString("Hello World!!!")); err != nil {
		fmt.Println(err)
	}

	// Receive the message from the queue
	if err := client.ReceiveOne(ctx, MessagePrinter{}); err != nil {
		fmt.Println(err)
	}

	// Output: Hello World!!!
}

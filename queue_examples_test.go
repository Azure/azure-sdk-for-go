package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func ExampleQueue_getOrBuildQueue() {
	const queueName = "myqueue"

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("Fatal: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qm := ns.NewQueueManager()
	qe, err := qm.Get(ctx, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	if qe == nil {
		_, err := qm.Put(ctx, queueName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	q, err := ns.NewQueue(queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(q.Name)
	// Output: myqueue
}

func ExampleQueue_Send() {
	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString("<your connection string here>"))
	if err != nil {
		return
	}

	client, err := ns.NewQueue("myqueue")
	if err != nil {
		return
	}

	// Create a context to limit how long we will try to send, then push the message over the wire.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client.Send(ctx, servicebus.NewMessageFromString("Hello World!!!"))
}

func ExampleQueue_Receive() {
	// Define a function that should be executed when a message is received.
	var printMessage servicebus.HandlerFunc = func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
		fmt.Println(string(msg.Data))
		return msg.Complete()
	}

	// Instantiate the clients needed to communicate with a Service Bus Queue.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString("<your connection string here>"))
	if err != nil {
		return
	}

	client, err := ns.NewQueue("myqueue")
	if err != nil {
		return
	}

	// Define a context to limit how long we will block to receive messages, then start serving our function.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	client.Receive(ctx, printMessage)
}

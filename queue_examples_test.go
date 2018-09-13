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

	fmt.Println(q.Name)
	// Output: myqueue
}

func ExampleQueue_scheduledMessage() {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("Fatal: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Initialize a client to communicate with a Service Bus Queue named scheduledmessages
	q, err := ns.NewQueue("scheduledmessages", servicebus.QueueWithReceiveAndDelete())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// Send the message "Hello World!" to the Queue named helloworld to be delivered in 20 seconds
	future := time.Now().UTC().Add(1 * time.Minute)
	msg := servicebus.NewMessageFromString("Hello World!")
	msg.SystemProperties = &servicebus.SystemProperties{
		ScheduledEnqueueTime: &future,
	}
	err = q.Send(ctx, msg)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("sent message %q. It should arrive in about a minute\n", msg.Data)

	received, err := q.ReceiveOne(ctx)
	q.Close(ctx)
	fmt.Printf("received message: %q\n", string(received.Data))

	// Output:
	// sent message "Hello World!". It should arrive in about a minute
	// received message: "Hello World!"
}

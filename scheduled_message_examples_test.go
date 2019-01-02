package servicebus_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-service-bus-go"
	"os"
	"time"
)

func Example_scheduledMessage() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		fmt.Println("FATAL: expected environment variable SERVICEBUS_CONNECTION_STRING not set")
		return
	}

	// Create a client to communicate with a Service Bus Namespace.
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// Create a client to communicate with the queue. (The queue must have already been created, see `QueueManager`)
	client, err := ns.NewQueue("scheduledmessages")
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// purge all of the existing messages in the queue
	purgeMessages(ns)

	// The delay that we should schedule a message for.
	const waitTime = 1 * time.Minute
	// Service Bus guarantees roughly a one minute window. So that our tests aren't flaky, we'll buffer our expectations
	// on either side.
	const buffer = 20 * time.Second

	expectedTime := time.Now().Add(waitTime)
	msg := servicebus.NewMessageFromString("to the future!!")
	msg.ScheduleAt(expectedTime)

	err = client.Send(ctx, msg)
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	err = client.ReceiveOne(
		ctx,
		servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			received := time.Now()
			if received.Before(expectedTime.Add(buffer)) && received.After(expectedTime.Add(-buffer)) {
				fmt.Println("Received when expected!")
			} else {
				fmt.Println("Received outside the expected window.")
			}
			return msg.Complete(ctx)
		}))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// Output: Received when expected!
}

func purgeMessages(ns *servicebus.Namespace) {
	purgeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := ns.NewQueue("scheduledmessages")
	defer func() {
		_ = client.Close(purgeCtx)
	}()
	defer cancel()
	_ = client.Receive(purgeCtx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
		return msg.Complete(ctx)
	}))
}

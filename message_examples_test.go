package servicebus_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-service-bus-go"
	"os"
	"time"
)

func ExampleMessage_ScheduleAt() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute+40*time.Second)
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

	// Create a client to communicate with the queue. (The queue must have already been created, see `QueueManager`)
	client, err := ns.NewQueue("scheduledmessages")
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// The delay that we should schedule a message for.
	const waitTime = time.Duration(1 * time.Minute)
	// Service Bus guarantees roughly a one minute window. So that our tests aren't flaky, we'll buffer our expectations
	// on either side.
	const buffer = time.Duration(20 * time.Second)

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
		servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) servicebus.DispositionAction {
			received := time.Now()
			if received.Before(expectedTime.Add(buffer)) && received.After(expectedTime.Add(-buffer)) {
				fmt.Println("Received when expected!")
			} else {
				fmt.Println("Received outside the expected window.")
			}
			return msg.Complete()
		}))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// Output: Received when expected!
}

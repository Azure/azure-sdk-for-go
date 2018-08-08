package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
			)

func Example_helloWorld() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStr := mustGetenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	const queueName = "helloworld"
	q, err := getQueue(ctx, ns, queueName)
	if err != nil {
		fmt.Printf("failed to build a new queue named %q\n", queueName)
		return
	}

	errs := make(chan error, 2)

	messages := []string{"hello", "world"}

	// Start a receiver that will print messages that it is informed of by Service Bus.
	go func(ctx context.Context, client *servicebus.Queue, quitAfter int) {
		received := make(chan struct{})

		listenHandle, err := client.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
			fmt.Println(string(message.Data))
			received <- struct{}{}
			return message.Complete()
		})
		if err != nil {
			errs <- err
			return
		}
		defer listenHandle.Close(context.Background())
		defer fmt.Println("...no longer listening")
		fmt.Println("listening...")

		for i := 0; i < quitAfter; i++ {
			select {
			case <-received:
				// Intentionally Left Blank
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case <-listenHandle.Done():
				errs <- listenHandle.Err()
				return
			}
		}
		errs <- nil
	}(ctx, q, len(messages))

	// Publish messages to Service Bus so that the receiver defined above will print them
	go func(ctx context.Context, client *servicebus.Queue, messages ...string) {
		for i := range messages {
			messageSent := make(chan error, 1)

			go func() {
				messageSent <- client.Send(ctx, servicebus.NewMessageFromString(messages[i]))
			}()
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			case err := <-messageSent:
				if err != nil {
					errs <- err
					return
				}
			}
		}
		errs <- nil
	}(ctx, q, messages...)

	for i := 0; i < 2; i++ {
		select {
		case err := <-errs:
			if err != nil {
				fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}

	// Output:
	// listening...
	// hello
	// world
	// ...no longer listening
}

func getQueue(ctx context.Context, ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	qm := ns.NewQueueManager()
	qe, err := qm.Get(ctx, queueName)
	if err != nil {
		return nil, err
	}

	if qe == nil {
		_, err := qm.Put(ctx, queueName)
		if err != nil {
			return nil, err
		}
	}

	q, err := ns.NewQueue(queueName)
	return q, err
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}

package main

import (
	"context"
	"fmt"
		"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
			)

func main() {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Initialize and create a Service Bus Queue named helloworld if it doesn't exist
	q, err := getQueue(ns, "helloworld")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Send the message "Hello World!" to the Queue named helloworld to be delivered in 20 seconds
	future := time.Now().UTC().Add(1 * time.Minute)
	msg := servicebus.NewMessageFromString("Hello World!")
	msg.SystemProperties = &servicebus.SystemProperties{
		ScheduledEnqueueTime: &future,
	}
	err = q.Send(context.Background(), msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("sent message at %v and should arrive in about a minute\n", time.Now().UTC())

	received, err := q.ReceiveOne(context.Background())
	received.Complete()
	q.Close(context.Background())
	fmt.Printf("received message: %q at %v\n", string(received.Data), time.Now().UTC())
}

func getQueue(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q, err := ns.NewQueue(ctx, queueName, servicebus.QueueWithReceiveAndDelete())
	return q, err
}

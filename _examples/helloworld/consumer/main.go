package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func main() {
	// Connect
	connStr := mustGetenv("SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	queueName := "helloworld"
	// Create the queue if it doesn't exist
	err = ensureQueue(ns, queueName)
	q := ns.NewQueue(queueName)

	// Start listening to events on the queue
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exit := make(chan struct{})
	listenHandle, err := q.Receive(ctx, func(ctx context.Context, event *servicebus.Event) error {
		text := string(event.Data)
		if text == "exit\n" {
			fmt.Println("Oh snap!! Someone told me to exit!")
			exit <- *new(struct{})
		} else {
			fmt.Println(string(event.Data))
		}
		return nil
	})
	defer listenHandle.Close(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("I am listening...")

	select {
	case <-exit:
		fmt.Println("closing after 2 seconds")
		select {
		case <-time.After(2 * time.Second):
			listenHandle.Close(context.Background())
			return
		}
	}
}

func ensureQueue(ns *servicebus.Namespace, queueName string) error {
	qm := ns.NewQueueManager()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := qm.Put(ctx, queueName)
	return err
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}
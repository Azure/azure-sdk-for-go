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
	q, err := getQueue(ns, queueName)
	if err != nil {
		fmt.Printf("failed to build a new queue named %q\n", queueName)
		fmt.Printf("error %s\n", err)
		os.Exit(1)
	}

	exit := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Println("Setting up listener...")
	listenHandle, err := q.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
		text := string(message.Data)
		if text == "exit\n" {
			fmt.Println("Oh snap!! Someone told me to exit!")
			exit <- *new(struct{})
		} else {
			fmt.Println(string(message.Data))
		}
		return message.Complete()
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

func getQueue(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

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

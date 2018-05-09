package main

import (
	"bufio"
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		q.Send(ctx, servicebus.NewEventFromString(text))
		if text == "exit\n" {
			break
		}
		cancel()
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

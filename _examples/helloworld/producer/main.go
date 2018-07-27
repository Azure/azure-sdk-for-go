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
	q, err := getQueue(ns, queueName)
	if err != nil {
		fmt.Printf("failed to build a new queue named %q\n", queueName)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		q.Send(ctx, servicebus.NewMessageFromString(text))
		if text == "exit\n" {
			break
		}
		cancel()
	}
}

func getQueue(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
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

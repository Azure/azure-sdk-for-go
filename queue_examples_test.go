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

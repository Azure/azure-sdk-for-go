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
	q, err := ns.NewQueue("helloworld")
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	err = q.Send(ctx, servicebus.NewMessageFromString("Hello, World!!!"))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	err = q.ReceiveOne(
		ctx,
		servicebus.HandlerFunc(func(ctx context.Context, message *servicebus.Message) error {
			fmt.Println(string(message.Data))
			return message.Complete(ctx)
		}))
	if err != nil {
		fmt.Println("FATAL: ", err)
		return
	}

	// Output: Hello, World!!!
}

func ensureTopic(ctx context.Context, tm *servicebus.TopicManager, name string, opts ...servicebus.TopicManagementOption) (*servicebus.TopicEntity, error) {
	te, err := tm.Get(ctx, name)
	if err == nil {
		_ = tm.Delete(ctx, name)
	}

	te, err = tm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return te, nil
}

func ensureQueue(ctx context.Context, qm *servicebus.QueueManager, name string, opts ...servicebus.QueueManagementOption) (*servicebus.QueueEntity, error) {
	qe, err := qm.Get(ctx, name)
	if err == nil {
		_ = qm.Delete(ctx, name)
	}

	qe, err = qm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return qe, nil
}

func ensureSubscription(ctx context.Context, sm *servicebus.SubscriptionManager, name string, opts ...servicebus.SubscriptionManagementOption) (*servicebus.SubscriptionEntity, error) {
	subEntity, err := sm.Get(ctx, name)
	if err == nil {
		_ = sm.Delete(ctx, name)
	}

	subEntity, err = sm.Put(ctx, name, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return subEntity, nil
}

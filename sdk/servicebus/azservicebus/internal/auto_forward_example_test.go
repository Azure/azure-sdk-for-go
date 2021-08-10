package servicebus_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

type MessagePrinter struct{}

func (mp MessagePrinter) Handle(ctx context.Context, msg *servicebus.Message) error {
	fmt.Println(string(msg.Data))
	return msg.Complete(ctx)
}

func Example_autoForward() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
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

	qm := ns.NewQueueManager()
	target, err := ensureQueue(ctx, qm, "AutoForwardTargetQueue")
	if err != nil {
		fmt.Println(err)
		return
	}

	source, err := ensureQueue(ctx, qm, "AutoForwardSourceQueue", servicebus.QueueEntityWithAutoForward(target))
	if err != nil {
		fmt.Println(err)
		return
	}

	sourceQueue, err := ns.NewQueue(source.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = sourceQueue.Close(ctx)
	}()

	if err := sourceQueue.Send(ctx, servicebus.NewMessageFromString("forward me to target!")); err != nil {
		fmt.Println(err)
		return
	}

	targetQueue, err := ns.NewQueue(target.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = targetQueue.Close(ctx)
	}()

	if err := targetQueue.ReceiveOne(ctx, MessagePrinter{}); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// forward me to target!
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

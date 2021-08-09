package servicebus_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func Example_prefetch() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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
	prefetch1, err := ensureQueue(ctx, qm, "Prefetch1Example")
	if err != nil {
		fmt.Println(err)
		return
	}

	prefetch1000, err := ensureQueue(ctx, qm, "Prefetch1000Example")
	if err != nil {
		fmt.Println(err)
		return
	}

	// sendAndReceive will send to the queue and read from the queue
	sendAndReceive := func(ctx context.Context, name string, opt servicebus.QueueOption) error {
		messageCount := 200
		q, err := ns.NewQueue(name, opt, servicebus.QueueWithReceiveAndDelete())
		if err != nil {
			return err
		}

		buffer := make([]byte, 1000)
		if _, err := rand.Read(buffer); err != nil {
			return err
		}

		for i := 0; i < messageCount; i++ {
			if err := q.Send(ctx, servicebus.NewMessage(buffer)); err != nil {
				return err
			}
		}

		innerCtx, cancel := context.WithCancel(ctx)
		count := 0
		err = q.Receive(innerCtx, servicebus.HandlerFunc(func(ctx context.Context, msg *servicebus.Message) error {
			count++
			if count == messageCount-1 {
				defer cancel()
			}
			return msg.Complete(ctx)
		}))
		if err != nil {
			if err.Error() != "context canceled" {
				return err
			}
		}
		return nil
	}

	// run send and receive concurrently and compare the times
	totalPrefetch1 := make(chan time.Duration)
	go func() {
		start := time.Now()
		if err := sendAndReceive(ctx, prefetch1.Name, servicebus.QueueWithPrefetchCount(1)); err != nil {
			fmt.Println(err)
			return
		}
		totalPrefetch1 <- time.Now().Sub(start)
	}()

	totalPrefetch1000 := make(chan time.Duration)
	go func() {
		start := time.Now()
		if err := sendAndReceive(ctx, prefetch1000.Name, servicebus.QueueWithPrefetchCount(1000)); err != nil {
			fmt.Println(err)
			return
		}
		totalPrefetch1000 <- time.Now().Sub(start)
	}()

	tp1 := <-totalPrefetch1
	tp2 := <-totalPrefetch1000

	if tp1 > tp2 {
		fmt.Println("prefetch of 1000 took less time!")
	}

	// Output:
	// prefetch of 1000 took less time!

}

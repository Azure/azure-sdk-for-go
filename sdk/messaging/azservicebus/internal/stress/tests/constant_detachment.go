// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ConstantDetachment(remainingArgs []string) {
	sc := shared.MustCreateStressContext("ConstantDetachment")
	defer sc.End()

	queueName := fmt.Sprintf("detach-tester-%s", sc.Nano)

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	stats := sc.NewStat("stats")

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("create a sender", err)

	const maxMessages = 20000

	shared.MustGenerateMessages(sc, sender, maxMessages, 1024, stats)

	time.Sleep(10 * time.Second)

	adminClient, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("creating admin client to check properties", err)

	queueProps, err := adminClient.GetQueueRuntimeProperties(sc.Context, queueName, nil)
	sc.PanicOnError("failed to get queue runtime properties", err)

	if queueProps.TotalMessageCount != maxMessages {
		sc.PanicOnError("", fmt.Errorf("incorrect number of messages (expected: %d, actual:%d)", maxMessages, queueProps.TotalMessageCount))
	}

	// now attempt to receive messages, while constant detaching is happening
	receiver, err := client.NewReceiverForQueue(queueName, &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModePeekLock,
	})
	sc.PanicOnError("failed to create receiver", err)
	go func() {
		// keep updating the definition of the queue, which will cause the service to detach us.
		err := shared.ConstantlyUpdateQueue(sc.Context, adminClient, queueName, 30*time.Second)

		if errors.Is(err, context.Canceled) {
			return
		}

		sc.PanicOnError("constantly updated queue failed", err)
	}()

InfiniteLoop:
	for stats.Completed != maxMessages {
		select {
		case <-sc.Done():
			break InfiniteLoop
		default:
		}

		messages, err := receiver.ReceiveMessages(sc.Context, 1000, nil)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Application is done, cancelling...")
			break InfiniteLoop
		}

		if err != nil {
			sc.LogIfFailed("receive failed, continuing", err, stats)
			continue
		}

		stats.AddReceived(int32(len(messages)))

		wg := sync.WaitGroup{}
		wg.Add(len(messages))

		for _, m := range messages {
			go func(m *azservicebus.ReceivedMessage) {
				defer wg.Done()

				if err := receiver.CompleteMessage(sc.Context, m, nil); err != nil {
					sc.LogIfFailed("Failed completing message", err, stats)
					return
				}

				stats.AddCompleted(int32(1))
			}(m)
		}

		wg.Wait()
	}

	if stats.Completed != maxMessages {
		sc.PanicOnError("constantdetach failed", fmt.Errorf("not all messages received (got %d, wanted %d)", stats.Received, maxMessages))
	} else {
		log.Printf("All messages received")
	}
}

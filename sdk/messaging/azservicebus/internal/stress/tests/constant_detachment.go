// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func ConstantDetachment(remainingArgs []string) {
	sc := shared.MustCreateStressContext("ConstantDetachment", nil)
	defer sc.End()

	queueName := fmt.Sprintf("detach-tester-%s", sc.Nano)

	shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
	sc.PanicOnError("create a sender", err)

	// this number isn't too special, but it gives us long enough so that
	// we are guaranteed that these detaches will interfere with our receiving.
	// Bug filed for this particular test, pertaining to settlement counting:
	//   https://github.com/Azure/azure-sdk-for-go/issues/17945)
	const maxMessages = 20000

	shared.MustGenerateMessages(sc, sender, maxMessages, 1024)

	// We'll give it a little time for the messages to show up in the runtime properties.
	// Just a simple sanity check that we don't have a bug in our message generator and that
	// they all exist on the remote entity.
	time.Sleep(10 * time.Second)

	adminClient, err := admin.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("creating admin client to check properties", err)

	queueProps, err := adminClient.GetQueueRuntimeProperties(sc.Context, queueName, nil)
	sc.PanicOnError("failed to get queue runtime properties", err)

	if queueProps.TotalMessageCount != maxMessages {
		sc.PanicOnError("", fmt.Errorf("incorrect number of messages (expected: %d, actual:%d)", maxMessages, queueProps.TotalMessageCount))
	}

	// now attempt to receive messages, while constant detaching is happening
	receiver, err := shared.NewTrackingReceiverForQueue(sc.TC, client, queueName, &azservicebus.ReceiverOptions{
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

	var completed int64
	var received int

InfiniteLoop:
	for atomic.LoadInt64(&completed) != maxMessages {
		select {
		case <-sc.Context.Done():
			break InfiniteLoop
		default:
		}

		messages, err := receiver.ReceiveMessages(sc.Context, 1000, nil)

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Application is done, cancelling...")
			break InfiniteLoop
		}

		if err != nil {
			sc.LogIfFailed("receive failed, continuing", err)
			continue
		}

		received += len(messages)

		wg := sync.WaitGroup{}
		wg.Add(len(messages))

		for _, m := range messages {
			go func(m *azservicebus.ReceivedMessage) {
				defer wg.Done()

				if err := receiver.CompleteMessage(sc.Context, m, nil); err != nil {
					return
				}

				atomic.AddInt64(&completed, 1)
			}(m)
		}

		wg.Wait()
	}

	if completed != maxMessages {
		sc.PanicOnError("constantdetach failed", fmt.Errorf("not all messages received (got %d, wanted %d)", received, maxMessages))
	} else {
		log.Printf("All messages received")
	}
}

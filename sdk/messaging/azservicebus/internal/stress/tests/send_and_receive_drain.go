// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func SendAndReceiveDrain(remainingArgs []string) {
	// set a long lock duration to make it obvious when a message is being lost in our
	// internal buffer or somewhere along the way.
	// This mimics the scenario mentioned in this issue filed by a customer:
	// https://github.com/Azure/azure-sdk-for-go/issues/17853
	const lockDuration = "PT5M"

	const numToSend = 2000
	const msgPadding = 4096

	sc := shared.MustCreateStressContext("SendAndReceiveDrainTest", nil)
	defer sc.End()

	queueName := strings.ToLower(fmt.Sprintf("queue-%X", time.Now().UnixNano()))
	sc.Start(queueName, map[string]string{
		"LockDuration":   lockDuration,
		"NumToSend":      fmt.Sprintf("%d", numToSend),
		"MessagePadding": fmt.Sprintf("%d", msgPadding),
	})

	log.Printf("Creating queue")

	adminClient := shared.MustCreateAutoDeletingQueue(sc, queueName, &admin.QueueProperties{
		LockDuration: to.Ptr(lockDuration),
	})

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
	sc.PanicOnError("failed to create sender", err)

	receiver, err := shared.NewTrackingReceiverForQueue(sc.TC, client, queueName, nil)
	sc.PanicOnError("Failed to create receiver", err)

	for i := 0; i < 1000; i++ {
		log.Printf("=====> Round [%d] <====", i)

		shared.MustGenerateMessages(sc, sender, numToSend, msgPadding)

		receivedIds := sync.Map{}
		var totalCompleted int64

		for totalCompleted < numToSend {
			log.Printf("Receiving messages [%d/%d]...", totalCompleted, numToSend)
			ctx, cancel := context.WithTimeout(sc.Context, time.Minute)
			defer cancel()

			messages, err := receiver.ReceiveMessages(ctx, numToSend+100, nil)

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				// this is bad - it means we didn't get _any_ messages within an entire
				// minute and might indicate that we're hitting the customer bug.

				log.Printf("Exceeded the timeout, trying one more time real fast")

				// let's see if there is some other momentary issue happening here by doing a quick receive again.
				ctx, cancel := context.WithTimeout(sc.Context, time.Minute)
				defer cancel()
				messages, err = receiver.ReceiveMessages(ctx, numToSend+100, nil)
				sc.PanicOnError("Exceeded a minute while waiting for messages", err)
			}

			log.Printf("Got %d messages, completing...", len(messages))

			wg := sync.WaitGroup{}

			for _, m := range messages {
				wg.Add(1)
				go func(m *azservicebus.ReceivedMessage) {
					defer wg.Done()

					if len(m.Body) != msgPadding {
						sc.PanicOnError("Body length issue", fmt.Errorf("Invalid body length - expected %d, got %d", msgPadding, len(m.Body)))
					}

					if err := receiver.CompleteMessage(sc.Context, m, nil); err != nil {
						sc.PanicOnError("Failed to complete message", err)
					}

					atomic.AddInt64(&totalCompleted, 1)

					num := int(m.ApplicationProperties["Number"].(int64))

					if _, exists := receivedIds.LoadOrStore(num, true); exists {
						sc.Failf("Duplicate message received with Number %d", num)
					}
				}(m)
			}

			wg.Wait()
		}

		log.Printf("[end] Receiving messages (all received)")

		// some sanity checking
		rtp, err := adminClient.GetQueueRuntimeProperties(context.Background(), queueName, nil)
		sc.PanicOnError("Failed to get runtime propeties for queue", err)

		if rtp.ActiveMessageCount != 0 {
			sc.PanicOnError(fmt.Sprintf("No messages should be active in the queue, but actually still had %d", rtp.ActiveMessageCount), errors.New("Messages still left in queue"))
		}
	}
}

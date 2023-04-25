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
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func FiniteSendAndReceiveTest(remainingArgs []string) {
	sc := shared.MustCreateStressContext("FiniteSendAndReceiveTest", nil)

	sc.TrackEvent("Start")
	defer sc.End()

	queueName := strings.ToLower(fmt.Sprintf("queue-%X", time.Now().UnixNano()))

	log.Printf("Creating queue")

	lockDuration := "PT5M"

	shared.MustCreateAutoDeletingQueue(sc, queueName, &admin.QueueProperties{
		LockDuration: &lockDuration,
	})

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("failed to create sender", err)

	// 50000 isn't particularly special, but it does give us a decent # of receives
	// so we get a decent view into our performance.
	const messageLimit = 50000

	log.Printf("Sending %d messages (all messages will be sent before receiving begins)", messageLimit)
	stats := sc.NewStat("finite")
	shared.MustGenerateMessages(sc, sender, messageLimit, 100, stats)

	log.Printf("Starting receiving...")

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	sc.PanicOnError("Failed to create receiver", err)

	for stats.Received < messageLimit {
		log.Printf("[start] Receiving messages...")
		messages, err := receiver.ReceiveMessages(context.Background(), 1000, nil)
		log.Printf("[done] Receiving messages... %v, %v", len(messages), err)
		sc.PanicOnError("failed to create client", err)

		wg := sync.WaitGroup{}

		log.Printf("About to complete %d messages", len(messages))
		time.Sleep(10 * time.Second)

		for _, msg := range messages {
			wg.Add(1)

			go func(msg *azservicebus.ReceivedMessage) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
				defer cancel()

				err := receiver.CompleteMessage(ctx, msg, nil)

				var sbErr *azservicebus.Error

				if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeLockLost {
					stats.AddError("lock lost", err)
					return
				}

				sc.PanicOnError("failed to complete message", err)
				stats.AddReceived(1)
			}(msg)
		}

		wg.Wait()
	}
}

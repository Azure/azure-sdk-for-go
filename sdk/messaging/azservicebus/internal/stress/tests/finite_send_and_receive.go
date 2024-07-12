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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func FiniteSendAndReceiveTest(remainingArgs []string) {
	const lockDuration = "PT5M"

	// 50000 isn't particularly special, but it does give us a decent # of receives
	// so we get a decent view into our performance.
	const messageLimit = 50000

	sc := shared.MustCreateStressContext("FiniteSendAndReceiveTest", nil)
	defer sc.End()

	queueName := strings.ToLower(fmt.Sprintf("queue-%X", time.Now().UnixNano()))
	sc.Start(queueName, map[string]string{
		"LockDuration": lockDuration,
		"MessageLimit": fmt.Sprintf("%d", messageLimit),
	})

	log.Printf("Creating queue")

	shared.MustCreateAutoDeletingQueue(sc, queueName, &admin.QueueProperties{
		LockDuration: to.Ptr(lockDuration),
	})

	client, err := azservicebus.NewClient(sc.Endpoint, sc.Cred, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := shared.NewTrackingSender(sc.TC, client, queueName, nil)
	sc.PanicOnError("failed to create sender", err)

	log.Printf("Sending %d messages (all messages will be sent before receiving begins)", messageLimit)
	shared.MustGenerateMessages(sc, sender, messageLimit, 100)

	log.Printf("Starting receiving...")

	receiver, err := shared.NewTrackingReceiverForQueue(sc.TC, client, queueName, nil)
	sc.PanicOnError("Failed to create receiver", err)

	var received int64

	for received < messageLimit {
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
					return
				}

				sc.PanicOnError("failed to complete message", err)
			}(msg)
		}

		wg.Wait()
	}
}

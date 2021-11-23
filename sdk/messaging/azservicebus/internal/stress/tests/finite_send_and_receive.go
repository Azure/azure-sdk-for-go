// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func FiniteSendAndReceiveTest(remainingArgs []string) {
	sc := shared.MustCreateStressContext("FiniteSendAndReceiveTest")

	sc.TrackEvent("Start")
	defer sc.TrackEvent("End")

	queueName := strings.ToLower(fmt.Sprintf("queue-%X", time.Now().UnixNano()))

	log.Printf("Creating queue")
	shared.MustCreateAutoDeletingQueue(sc, queueName)

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create client", err)

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("failed to create sender", err)
	const messageLimit = 50000

	shared.MustGenerateMessages(sc, sender, messageLimit, 100, sc.NewStat("sender"))

	log.Printf("Sending %d messages (all messages will be sent before receiving begins)", messageLimit)
	log.Printf("Starting receiving...")

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	sc.PanicOnError("Failed to create receiver", err)

	completions := make(chan struct{}, 100)

	receiverStats := sc.NewStat("receiver")

	for receiverStats.Received == messageLimit {
		log.Printf("[start] Receiving messages...")
		messages, err := receiver.ReceiveMessages(context.Background(), 100, nil)
		log.Printf("[done] Receiving messages... %v, %v", len(messages), err)
		sc.PanicOnError("failed to create client", err)

		wg := sync.WaitGroup{}

		for _, msg := range messages {
			wg.Add(1)
			go func(msg *azservicebus.ReceivedMessage) {
				completions <- struct{}{}
				defer wg.Done()
				defer func() { <-completions }()

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				err := receiver.CompleteMessage(ctx, msg)
				sc.PanicOnError("failed to complete message", err)
			}(msg)
		}

		wg.Wait()

		receiverStats.AddReceived(int32(len(messages)))
	}
}

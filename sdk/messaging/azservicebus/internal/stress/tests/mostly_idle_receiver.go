// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

// MostlyIdleReceiver tests that if there are long idle periods that our connection continues to work and receive messages.
func MostlyIdleReceiver(remainingArgs []string) {
	sc := shared.MustCreateStressContext("MostlyIdleReceiver", nil)
	defer sc.End()

	// we'll try several levels of "idleness", with different connections to make sure they don't
	// affect each other.

	durations := []time.Duration{
		time.Second,
		30 * time.Second,
		time.Minute,
		15 * time.Minute,
		30 * time.Minute,
		time.Hour,
		2 * time.Hour,
		3 * time.Hour,
		24 * time.Hour,
		36 * time.Hour,
		2 * 24 * time.Hour,
		3 * 24 * time.Hour,
		4 * 24 * time.Hour,
	}

	wg := sync.WaitGroup{}

	log.Println("Running tests for wait times of:")

	for i, dur := range durations {
		log.Printf("%d  %s", i+1, dur)
	}

	for _, duration := range durations {
		stats := sc.NewStat(fmt.Sprintf("Duration: %s", duration))

		wg.Add(1)

		go func(stats *shared.Stats, duration time.Duration) {
			defer wg.Done()

			queueName := fmt.Sprintf("mostly-idle-receiver-%s-%s", sc.Nano, duration)
			shared.MustCreateAutoDeletingQueue(sc, queueName, nil)

			client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
			sc.PanicOnError("failed to create client", err)

			defer func() {
				err = client.Close(context.Background())
				sc.LogIfFailed("failed to close client", err, stats)
			}()

			receiver, err := client.NewReceiverForQueue(queueName, nil)
			sc.PanicOnError("failed to create receiver", err)

			sender, err := client.NewSender(queueName, nil)
			sc.PanicOnError("failed to create sender", err)

			time.AfterFunc(duration, func() {
				log.Printf("Sending message for duration %s", duration)
				err := sender.SendMessage(sc.Context, &azservicebus.Message{
					Body: []byte(fmt.Sprintf("Message for %s", duration)),
				}, nil)
				sc.PanicOnError(fmt.Sprintf("failed sending message for duration %s", duration), err)
			})

			log.Printf("Waiting for message to arrive, after duration %s", duration)
			messages, err := receiver.ReceiveMessages(sc.Context, 1, nil)
			sc.PanicOnError(fmt.Sprintf("failed receiving messages for duration %s", duration), err)

			log.Printf("Received %d messages", len(messages))
			stats.AddReceived(int32(len(messages)))

			for _, msg := range messages {
				err := receiver.CompleteMessage(sc.Context, msg, nil)
				sc.PanicOnError("failed to complete message", err)
			}
		}(stats, duration)
	}

	wg.Wait()
}

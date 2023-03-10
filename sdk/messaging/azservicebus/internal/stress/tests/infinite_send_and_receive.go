// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func InfiniteSendAndReceiveRun(remainingArgs []string) {
	sc := shared.MustCreateStressContext("InfiniteSendAndReceiveRun")
	defer sc.End()

	topicName := fmt.Sprintf("topic-%s", sc.Nano)
	sc.Start(topicName, nil)

	stats := sc.NewStat("infinite")

	cleanup := shared.MustCreateSubscriptions(sc, topicName, []string{"batch"}, nil)
	defer cleanup()

	time.AfterFunc(5*24*time.Hour, func() {
		sc.End()
	})

	go func() {
		runBatchReceiver(sc, topicName, "batch", stats)
	}()

	go func() {
		continuallySend(sc, topicName)
	}()

	<-sc.Context.Done()
}

func runBatchReceiver(sc *shared.StressContext, topicName string, subscriptionName string, stats *shared.Stats) {
	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create a client", err)

	receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)

	if err != nil {
		log.Fatalf("Failed to create receiver: %s", err.Error())
	}

	for {
		messages, err := receiver.ReceiveMessages(sc.Context, 20, nil)

		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				break
			}

			sc.PanicOnError("failed when receiving", err)
		}

		stats.AddReceived(int32(len(messages)))

		for _, msg := range messages {
			go func(msg *azservicebus.ReceivedMessage) {
				err := receiver.CompleteMessage(sc.Context, msg, nil)
				sc.LogIfFailed("complete failed", err, stats)

				if err == nil {
					stats.AddCompleted(1)
				}
			}(msg)
		}
	}
}

func continuallySend(sc *shared.StressContext, queueName string) {
	ctx, cancel := context.WithTimeout(sc.Context, 5*24*time.Hour)
	defer cancel()

	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create a connection string", err)

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("failed to create sender", err)

	defer sender.Close(ctx)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	senderStats := sc.NewStat("sender")

	for t := range ticker.C {
		err := sender.SendMessage(ctx, &azservicebus.Message{
			Body: []byte(fmt.Sprintf("hello world: %s", t.String())),
		}, nil)

		atomic.AddInt32(&senderStats.Sent, 1)
		sc.TrackMetric(string(MetricNameMessageSent), 1)

		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				log.Printf("Test complete, stopping sender loop")
				break
			}

			sc.LogIfFailed("failed to send message", err, senderStats)
			break
		}
	}
}

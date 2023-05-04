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
	sc := shared.MustCreateStressContext("InfiniteSendAndReceiveRun", &shared.StressContextOptions{
		Duration: 5 * 24 * time.Hour,
	})
	defer sc.End()

	topicName := fmt.Sprintf("topic-%s", sc.Nano)
	sc.Start(topicName, nil)

	stats := sc.NewStat("infinite")

	cleanup := shared.MustCreateSubscriptions(sc, topicName, []string{"batch"}, nil)
	defer cleanup()

	go runBatchReceiver(sc, topicName, "batch", stats)
	go continuallySend(sc, topicName)

	<-sc.Context.Done()
}

func runBatchReceiver(sc *shared.StressContext, topicName string, subscriptionName string, stats *shared.Stats) {
	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create a client", err)

	defer client.Close(context.Background())

	receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)
	sc.PanicOnError("failed to create receiver", err)

	defer receiver.Close(context.Background())

	for {
		messages, err := receiver.ReceiveMessages(sc.Context, 20, nil)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Printf("Test deadline reached, receiver is closing")
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
	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create a connection string", err)

	sender, err := client.NewSender(queueName, nil)
	sc.PanicOnError("failed to create sender", err)

	defer sender.Close(context.Background())

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	senderStats := sc.NewStat("sender")

	for t := range ticker.C {
		err := sender.SendMessage(sc.Context, &azservicebus.Message{
			Body: []byte(fmt.Sprintf("hello world: %s", t.String())),
		}, nil)

		atomic.AddInt32(&senderStats.Sent, 1)
		sc.TrackMetric(string(MetricNameMessageSent), 1)

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				log.Printf("Test deadline reached, stopping sender loop")
				break
			}

			sc.LogIfFailed("failed to send message", err, senderStats)
			break
		}
	}
}

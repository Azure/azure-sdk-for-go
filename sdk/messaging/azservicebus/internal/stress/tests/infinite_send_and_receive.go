// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tests

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// test metrics
const (
	MetricMessageSent = "MessageSent"
)

func InfiniteSendAndReceiveRun(remainingArgs []string) {
	sc := shared.MustCreateStressContext("InfiniteSendAndReceiveRun")

	topicName := fmt.Sprintf("topic-%s", sc.Nano)

	startEvent := appinsights.NewEventTelemetry("Start")
	startEvent.Properties["Topic"] = topicName
	sc.Track(startEvent)

	cleanup := shared.MustCreateSubscriptions(sc, topicName, []string{"batch"})
	defer cleanup()

	for i := 0; i < 4; i++ {
		go func(i int) {
			// give it a bunch of iterations if it should fail.
			for attempt := 0; attempt < 10000; attempt++ {
				runBatchReceiver(sc, fmt.Sprintf("%d-%d", i, attempt), topicName, "batch")
			}
		}(i)
	}

	go func() {
		for {
			continuallySend(sc, topicName)
		}
	}()

	ch := make(chan struct{})
	<-ch
}

func runBatchReceiver(sc *shared.StressContext, id string, topicName string, subscriptionName string) {
	client, err := azservicebus.NewClientFromConnectionString(sc.ConnectionString, nil)
	sc.PanicOnError("failed to create a client", err)

	receiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, nil)

	if err != nil {
		log.Fatalf("Failed to create receiver: %s", err.Error())
	}

	stats := sc.NewStat("receiver[" + id + "]")

	for {
		messages, err := receiver.ReceiveMessages(sc.Context, 20, nil)

		if err != nil {
			sc.LogIfFailed("failed to receive messages", err, stats)
			continue
		}

		stats.AddReceived(int32(len(messages)))

		for _, msg := range messages {
			go func(msg *azservicebus.ReceivedMessage) {
				err := receiver.CompleteMessage(sc.Context, msg)
				sc.LogIfFailed("complete failed", err, stats)
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
		})

		atomic.AddInt32(&senderStats.Sent, 1)
		sc.TrackMetric(MetricMessageSent, 1)

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

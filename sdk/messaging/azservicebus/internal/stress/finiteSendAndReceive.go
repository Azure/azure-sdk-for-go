// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func finiteSendAndReceiveTest(cs string, telemetryClient appinsights.TelemetryClient) {
	telemetryClient.TrackEvent("Start")
	defer telemetryClient.TrackEvent("End")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	stats := &stats{}

	// create a new queue
	adminClient, err := admin.NewClientFromConnectionString(cs, nil)

	if err != nil {
		trackError(stats, telemetryClient, "failed to create adminclient", err)
		panic(err)
	}

	queueName := strings.ToLower(fmt.Sprintf("queue-%X", time.Now().UnixNano()))

	log.Printf("Creating queue")

	_, err = adminClient.CreateQueue(context.Background(), queueName, nil, nil)

	if err != nil {
		trackError(stats, telemetryClient, fmt.Sprintf("failed to create queue %s", queueName), err)
		panic(err)
	}

	defer func() {
		_, err = adminClient.DeleteQueue(context.Background(), queueName, nil)

		if err != nil {
			trackError(stats, telemetryClient, fmt.Sprintf("failed to create queue %s", queueName), err)
		}
	}()

	client, err := azservicebus.NewClientFromConnectionString(cs, nil)

	if err != nil {
		trackError(stats, telemetryClient, "failed to create client", err)
		panic(err)
	}

	sender, err := client.NewSender(queueName, nil)

	if err != nil {
		trackError(stats, telemetryClient, "failed to create client", err)
		panic(err)
	}

	startStatsTicker(ctx, "finitesr", stats, 5*time.Second)

	const messageLimit = 10000

	log.Printf("Sending %d messages", messageLimit)

	for i := 0; i < messageLimit; i++ {
		err := sender.SendMessage(context.Background(), &azservicebus.Message{
			Body: []byte(fmt.Sprintf("Message %d", i)),
		})

		if err != nil {
			trackError(stats, telemetryClient, "failed to create client", err)
			panic(err)
		}

		atomic.AddInt32(&stats.Sent, 1)
	}

	log.Printf("Starting receiving...")

	receiver, err := client.NewReceiverForQueue(queueName, nil)

	var all []*azservicebus.ReceivedMessage

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		messages, err := receiver.ReceiveMessages(context.Background(), 100, nil)

		if err != nil {
			trackError(stats, telemetryClient, "failed to create client", err)
			panic(err)
		}

		for _, msg := range messages {
			if err := receiver.CompleteMessage(ctx, msg); err != nil {
				trackError(stats, telemetryClient, "failed to create client", err)
				panic(err)
			}
		}

		atomic.AddInt32(&stats.Received, int32(len(messages)))
		all = append(all, messages...)

		if len(all) == messageLimit {
			log.Printf("All messages received!")
			break
		}
	}
}

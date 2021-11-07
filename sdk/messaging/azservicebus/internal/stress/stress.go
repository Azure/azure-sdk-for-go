// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/devigned/tab"
	"github.com/joho/godotenv"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// Simple query to view some of the stats reported by this stress test
//
// customMetrics | where customDimensions["TestRunId"] == "169C8700A767E314"
// | project timestamp, name, valueCount
// | summarize Sum=sum(valueCount) by bin(timestamp, 30s), name
// | render timechart

type stats struct {
	Sent     int32
	Received int32
	Errors   int32
}

var receiverStats stats
var senderStats stats

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Failed to load .env file: %s", err.Error())
	}

	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	aiKey := os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY")

	if cs == "" || aiKey == "" {
		log.Fatalf("APPINSIGHTS_INSTRUMENTATIONKEY and SERVICEBUS_CONNECTION_STRING must be defined in the environment")
	}

	config := appinsights.NewTelemetryConfiguration(aiKey)
	config.MaxBatchInterval = 5 * time.Second
	telemetryClient := appinsights.NewTelemetryClientFromConfig(config)

	if len(os.Args) == 1 {
		runBasicSendAndReceiveTest(cs, telemetryClient)
	}

	switch os.Args[1] {
	case "finiteSendAndReceive":
		finiteSendAndReceiveTest(cs, telemetryClient)
	default:
		log.Printf("Unknown test: %s", os.Args[0])
	}
}

func startStatsTicker(ctx context.Context, prefix string, stats *stats, interval time.Duration) {
	go func(ctx context.Context) {
		ticker := time.NewTicker(interval)

	TickerLoop:
		for range ticker.C {
			select {
			case <-ctx.Done():
				ticker.Stop()
				break TickerLoop
			default:
			}

			log.Printf("[%s] Received: (r:%d), Sent: %d, Errors: %d",
				prefix,
				atomic.LoadInt32(&stats.Received),
				atomic.LoadInt32(&stats.Sent),
				atomic.LoadInt32(&stats.Errors),
			)
		}
	}(ctx)
}

func runBasicSendAndReceiveTest(cs string, telemetryClient appinsights.TelemetryClient) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			log.Printf("Received: (r:%d), Sent: %d, Errors: (r:%d,s:%d)",
				atomic.LoadInt32(&receiverStats.Received),
				atomic.LoadInt32(&senderStats.Sent),
				atomic.LoadInt32(&receiverStats.Errors),
				atomic.LoadInt32(&senderStats.Errors))
		}
	}()

	nanoSeconds := time.Now().UnixNano()
	topicName := fmt.Sprintf("topic-%X", nanoSeconds)

	telemetryClient.Context().CommonProperties = map[string]string{
		"Test":      "SendAndReceive",
		"TestRunId": fmt.Sprintf("%X", nanoSeconds),
	}

	log.Printf("Common properties: %+v", telemetryClient.Context().CommonProperties)

	defer func() {
		log.Printf("Flushing remaining telemetry")
		<-telemetryClient.Channel().Close()
		log.Printf("Flushed")
	}()

	startEvent := appinsights.NewEventTelemetry("Start")
	startEvent.Properties = map[string]string{
		"Topic": topicName,
	}

	telemetryClient.Track(startEvent)

	cleanup, err := createSubscriptions(telemetryClient, cs, topicName, []string{"batch"})
	defer cleanup()

	if err != nil {
		log.Printf("Failed to create topic and subscriptions: %s", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 24*5*time.Hour)
	defer cancel()

	serviceBusClient, err := azservicebus.NewClientFromConnectionString(cs, nil)
	if err != nil {
		trackError(nil, telemetryClient, "Failed to create service bus client", err)
		return
	}

	defer func() {
		if err := serviceBusClient.Close(ctx); err != nil {
			trackError(nil, telemetryClient, "Error when closing client", err)
		}
	}()

	tab.Register(&utils.StderrTracer{
		Include: map[string]bool{
			// internal.SpanProcessorClose: true,
			// internal.SpanProcessorLoop:  true,

			//internal.SpanProcessorMessage: true,
			tracing.SpanRecover:        true,
			tracing.SpanNegotiateClaim: true,
			tracing.SpanRecoverClient:  true,
			tracing.SpanRecoverLink:    true,
		},
	})

	go func() {
		for {
			runBatchReceiver(ctx, serviceBusClient, topicName, "batch", telemetryClient)
		}
	}()

	go func() {
		for {
			continuallySend(ctx, serviceBusClient, topicName, telemetryClient)
		}
	}()

	ch := make(chan struct{})
	<-ch
}

func runBatchReceiver(ctx context.Context, serviceBusClient *azservicebus.Client, topicName string, subscriptionName string, telemetryClient appinsights.TelemetryClient) {
	receiver, err := serviceBusClient.NewReceiverForSubscription(topicName, subscriptionName, nil)

	if err != nil {
		log.Fatalf("Failed to create receiver: %s", err.Error())
	}

	for {
		messages, err := receiver.ReceiveMessages(ctx, 20, nil)

		if err != nil {
			trackError(&receiverStats, telemetryClient, "receive batch failure", err)
		}

		atomic.AddInt32(&receiverStats.Received, int32(len(messages)))

		for _, msg := range messages {
			go func(msg *azservicebus.ReceivedMessage) {
				if err := receiver.CompleteMessage(ctx, msg); err != nil {
					trackError(&receiverStats, telemetryClient, "complete failed", err)
				}
			}(msg)
		}
	}
}

func continuallySend(ctx context.Context, client *azservicebus.Client, queueName string, telemetryClient appinsights.TelemetryClient) {
	sender, err := client.NewSender(queueName, nil)

	if err != nil {
		trackError(&senderStats, telemetryClient, "SenderCreate", err)
		return
	}

	defer sender.Close(ctx)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for t := range ticker.C {
		err := sender.SendMessage(ctx, &azservicebus.Message{
			Body: []byte(fmt.Sprintf("hello world: %s", t.String())),
		})

		atomic.AddInt32(&senderStats.Sent, 1)
		telemetryClient.TrackMetric("MessageSent", 1)

		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				log.Printf("Test complete, stopping sender loop")
				break
			}

			trackError(&senderStats, telemetryClient, "SendMessage", err)
			break
		}
	}
}

func createSubscriptions(telemetryClient appinsights.TelemetryClient, connectionString string, topicName string, subscriptionNames []string) (func(), error) {
	log.Printf("[BEGIN] Creating topic %s", topicName)
	defer log.Printf("[END] Creating topic %s", topicName)

	ac, err := admin.NewClientFromConnectionString(connectionString, nil)

	if err != nil {
		trackError(nil, telemetryClient, "Failed to create a topic manager", err)
		return nil, err
	}

	if _, err := ac.CreateTopic(context.Background(), topicName, nil, nil); err != nil {
		trackError(nil, telemetryClient, "Failed to create topic", err)
		return nil, err
	}

	for _, name := range subscriptionNames {
		if _, err := ac.CreateSubscription(context.Background(), topicName, name, nil, nil); err != nil {
			trackError(nil, telemetryClient, "Failed to create subscription manager", err)
		}
	}

	return func() {
		if _, err := ac.DeleteTopic(context.Background(), topicName, nil); err != nil {
			trackError(nil, telemetryClient, fmt.Sprintf("Failed to delete topic %s", topicName), err)
		}
	}, nil
}

func trackError(stats *stats, telemetryClient appinsights.TelemetryClient, message string, err error) {
	log.Printf("Exception: %s: %s", message, err.Error())

	if stats != nil {
		atomic.AddInt32(&stats.Errors, 1)
	}

	et := appinsights.NewExceptionTelemetry(err)
	et.Properties["Reason"] = message
	telemetryClient.Track(et)
}

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
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
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

var processorStats stats
var receiverStats stats
var senderStats stats

func main() {
	runBasicSendAndReceiveTest()
}

func runBasicSendAndReceiveTest() {
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

	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			log.Printf("Received: (p:%d,r:%d), Sent: %d, Errors: (p:%d,r:%d,s:%d)",
				atomic.LoadInt32(&processorStats.Received),
				atomic.LoadInt32(&receiverStats.Received),
				atomic.LoadInt32(&senderStats.Sent),
				atomic.LoadInt32(&processorStats.Errors),
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

	cleanup, err := createSubscriptions(telemetryClient, cs, topicName, []string{"processor", "batch"})
	defer cleanup()

	if err != nil {
		log.Printf("Failed to create topic and subscriptions: %s", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 24*5*time.Hour)
	defer cancel()

	serviceBusClient, err := azservicebus.NewClientWithConnectionString(cs, nil)
	if err != nil {
		trackException(nil, telemetryClient, "Failed to create service bus client", err)
		return
	}

	defer func() {
		if err := serviceBusClient.Close(ctx); err != nil {
			trackException(nil, telemetryClient, "Error when closing client", err)
		}
	}()

	tab.Register(&utils.StderrTracer{
		Include: map[string]bool{
			internal.SpanProcessorClose: true,
			internal.SpanProcessorLoop:  true,
			//internal.SpanProcessorMessage: true,
			internal.SpanRecover:        true,
			internal.SpanNegotiateClaim: true,
			internal.SpanRecoverClient:  true,
			internal.SpanRecoverLink:    true,
		},
	})

	runProcessorTest := true

	if runProcessorTest {
		go func() {
			for {
				runProcessor(ctx, serviceBusClient, topicName, "processor", telemetryClient)
			}
		}()
	} else {
		go func() {
			for {
				runBatchReceiver(ctx, serviceBusClient, topicName, "batch", telemetryClient)
			}
		}()
	}

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
			trackException(&receiverStats, telemetryClient, "receive batch failure", err)
		}

		atomic.AddInt32(&receiverStats.Received, int32(len(messages)))

		for _, msg := range messages {
			go func(msg *azservicebus.ReceivedMessage) {
				if err := receiver.CompleteMessage(ctx, msg); err != nil {
					trackException(&receiverStats, telemetryClient, "complete failed", err)
				}
			}(msg)
		}
	}
}

func runProcessor(ctx context.Context, client *azservicebus.Client, topicName string, subscriptionName string, telemetryClient appinsights.TelemetryClient) {
	log.Printf("Starting processor...")
	processor, err := client.NewProcessorForSubscription(
		topicName, subscriptionName,
		&azservicebus.ProcessorOptions{MaxConcurrentCalls: 10})

	if err != nil {
		trackException(&processorStats, telemetryClient, "Failed when creating processor", err)
		return
	}

	err = processor.Start(ctx, func(msg *azservicebus.ReceivedMessage) error {
		atomic.AddInt32(&processorStats.Received, 1)
		telemetryClient.TrackMetric("MessageReceived", 1)
		return nil
	}, func(err error) {
		log.Printf("Exception in processor: %s", err.Error())
		trackException(&processorStats, telemetryClient, "Processor.HandleError", err)
	})

	if err != nil {
		log.Printf("Exception when starting processor: %s", err.Error())
		trackException(&processorStats, telemetryClient, "Processor.Start", err)
		return
	}

	<-ctx.Done()

	telemetryClient.TrackEvent("ProcessorStopped")
	log.Print("Processor was stopped!")
}

func continuallySend(ctx context.Context, client *azservicebus.Client, queueName string, telemetryClient appinsights.TelemetryClient) {
	sender, err := client.NewSender(queueName)

	if err != nil {
		trackException(&senderStats, telemetryClient, "SenderCreate", err)
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

			trackException(&senderStats, telemetryClient, "SendMessage", err)
			break
		}
	}
}

func createSubscriptions(telemetryClient appinsights.TelemetryClient, connectionString string, topicName string, subscriptionNames []string) (func(), error) {
	log.Printf("[BEGIN] Creating topic %s", topicName)
	defer log.Printf("[END] Creating topic %s", topicName)

	tm, err := internal.NewTopicManagerWithConnectionString(connectionString)

	if err != nil {
		trackException(nil, telemetryClient, "Failed to create a topic manager", err)
		return nil, err
	}

	if _, err := tm.Put(context.TODO(), topicName); err != nil {
		trackException(nil, telemetryClient, "Failed to create topic", err)
		return nil, err
	}

	sm, err := internal.NewSubscriptionManagerForConnectionString(topicName, connectionString)

	if err != nil {
		trackException(nil, telemetryClient, "Failed to create subscription manager", err)
		return nil, err
	}

	for _, name := range subscriptionNames {
		if _, err := sm.Put(context.Background(), name); err != nil {
			trackException(nil, telemetryClient, "Failed to create subscription manager", err)
		}
	}

	return func() {
		if err := tm.Delete(context.TODO(), topicName); err != nil {
			trackException(nil, telemetryClient, fmt.Sprintf("Failed to delete topic %s", topicName), err)
		}
	}, nil
}

func trackException(stats *stats, telemetryClient appinsights.TelemetryClient, message string, err error) {
	log.Printf("Exception: %s: %s", message, err.Error())

	if stats != nil {
		atomic.AddInt32(&stats.Errors, 1)
	}

	et := appinsights.NewExceptionTelemetry(err)
	et.Properties["Reason"] = message
	telemetryClient.Track(et)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/joho/godotenv"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// Simple query to view some of the stats reported by this stress test
//
// customMetrics | where customDimensions["TestRunId"] == "169C8700A767E314"
// | project timestamp, name, valueCount
// | summarize Sum=sum(valueCount) by bin(timestamp, 30s), name
// | render timechart

var messagesSent int32
var messagesReceived int32
var exceptions int32

func main() {
	runBasicSendAndReceiveTest()
}

func runBasicSendAndReceiveTest() {
	godotenv.Load()

	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	aiKey := os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY")

	if cs == "" || aiKey == "" {
		log.Fatalf("APPINSIGHTS_INSTRUMENTATIONKEY and SERVICEBUS_CONNECTION_STRING must be defined in the environment")
	}

	config := appinsights.NewTelemetryConfiguration(aiKey)
	config.MaxBatchInterval = time.Second * 5
	telemetryClient := appinsights.NewTelemetryClientFromConfig(config)

	go func() {
		ticker := time.NewTicker(time.Second * 5)

		for range ticker.C {
			log.Printf("Received: %d, Sent: %d, Exceptions: %d", atomic.LoadInt32(&messagesReceived), atomic.LoadInt32(&messagesSent), atomic.LoadInt32(&exceptions))
		}
	}()

	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

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
		"Queue": queueName,
	}

	telemetryClient.Track(startEvent)

	cleanupQueue := createQueue(telemetryClient, cs, queueName)
	defer cleanupQueue()

	serviceBusClient, err := azservicebus.NewClient(azservicebus.WithConnectionString(cs))
	if err != nil {
		trackException(telemetryClient, "Failed to create service bus client", err)
		return
	}

	defer func() {
		if err := serviceBusClient.Close(context.TODO()); err != nil {
			trackException(telemetryClient, "Error when closing client", err)
		}
	}()

	go func() {
		for {
			runProcessor(serviceBusClient, queueName, telemetryClient)
		}
	}()

	go func() {
		for {
			continuallySend(serviceBusClient, queueName, telemetryClient)
		}
	}()

	ch := make(chan struct{})
	<-ch
}

func runProcessor(client *azservicebus.Client, queueName string, telemetryClient appinsights.TelemetryClient) {
	log.Printf("Starting processor...")
	processor, err := client.NewProcessor(azservicebus.ProcessorWithQueue(queueName), azservicebus.ProcessorWithMaxConcurrentCalls(10))

	if err != nil {
		trackException(telemetryClient, "Failed when creating processor", err)
		return
	}

	err = processor.Start(func(msg *azservicebus.ReceivedMessage) error {
		atomic.AddInt32(&messagesReceived, 1)
		telemetryClient.TrackMetric("MessageReceived", 1)
		return nil
	}, func(err error) {
		atomic.AddInt32(&exceptions, 1)
		log.Printf("Exception in processor: %s", err.Error())
		trackException(telemetryClient, "Processor.HandleError", err)
	})

	if err != nil {
		atomic.AddInt32(&exceptions, 1)
		log.Printf("Exception when starting processor: %s", err.Error())
		trackException(telemetryClient, "Processor.Start", err)
		return
	}

	<-processor.Done()

	telemetryClient.TrackEvent("ProcessorStopped")
	log.Print("Processor was stopped!")
}

func continuallySend(client *azservicebus.Client, queueName string, telemetryClient appinsights.TelemetryClient) {
	sender, err := client.NewSender(queueName)

	if err != nil {
		trackException(telemetryClient, "SenderCreate", err)
		return
	}

	defer sender.Close(context.TODO())

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	for t := range ticker.C {
		err := sender.SendMessage(context.Background(), &azservicebus.Message{
			Body: []byte(fmt.Sprintf("hello world: %s", t.String())),
		})

		atomic.AddInt32(&messagesSent, 1)
		telemetryClient.TrackMetric("MessageSent", 1)

		if err != nil {
			trackException(telemetryClient, "SendMessage", err)
			break
		}
	}
}

func createQueue(telemetryClient appinsights.TelemetryClient, connectionString string, queueName string) func() {
	log.Printf("[BEGIN] Creating queue %s", queueName)
	defer log.Printf("[END] Creating queue %s", queueName)
	ns, err := internal.NewNamespace(internal.NamespaceWithConnectionString(connectionString))

	if err != nil {
		trackException(telemetryClient, "Failed to create namespace client", err)
		return nil
	}

	qm := ns.NewQueueManager()

	_, err = qm.Put(context.TODO(), queueName)

	if err != nil {
		trackException(telemetryClient, "Failed to create queue", err)
		return nil
	}

	return func() {
		if err := qm.Delete(context.TODO(), queueName); err != nil {
			trackException(telemetryClient, "Failed to delete queue", err)
		}
	}
}

func trackException(telemetryClient appinsights.TelemetryClient, message string, err error) {
	log.Printf("Exception: %s: %s", message, err.Error())
	et := appinsights.NewExceptionTelemetry(err)
	et.Properties["ExtraMessage"] = message
	telemetryClient.Track(et)
}

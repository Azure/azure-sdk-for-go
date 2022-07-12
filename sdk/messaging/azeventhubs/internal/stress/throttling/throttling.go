// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/stress"
	"github.com/devigned/tab"
	"github.com/joho/godotenv"
)

var MaxBatches = 50
var SenderMaxRetryCount = 10

const TestIdProperty = "testId"

func main() {
	envFile := "../../../.env"
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Failed to load .env file from %s: %s", envFile, err.Error())
	}

	cs := os.Getenv("EVENTHUB_CONNECTION_STRING")

	hub, err := eventhub.NewHubFromConnectionString(cs, eventhub.HubWithSenderMaxRetryCount(SenderMaxRetryCount))

	if err != nil {
		log.Fatalf("Failed to create hub: %s", err.Error())
	}

	partitions := getPartitionCounts(context.Background(), hub)

	// Generate some large batches of messages and send them in parallel.
	// The Go SDK is fast enough that this will cause a 1TU instance to throttle
	// us, allowing you to see how our code reacts to it.
	tab.Register(&stress.StderrTracer{NoOpTracer: &tab.NoOpTracer{}})
	messageCount := sendMessages(hub)

	log.Printf("Sending complete, last expected ID = %d", messageCount)

	endSequenceNumbers := getPartitionCounts(context.Background(), hub)

	for partitionID, partition := range endSequenceNumbers {
		startSequenceNumber := partitions[partitionID].LastSequenceNumber

		log.Printf("[%s] diff: %d", partitionID, partition.LastSequenceNumber-startSequenceNumber)
	}

	// now receive and check all the messages, make sure everything arrived.
	verifyMessages(context.TODO(), hub, partitions, messageCount)
}

func sendMessages(hub *eventhub.Hub) int64 {
	var batches []eventhub.BatchIterator
	nextTestId := int64(0)

	log.Printf("Creating event batches")

	for i := 0; i < MaxBatches; i++ {
		batches = append(batches, createEventBatch(&nextTestId))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	wg := &sync.WaitGroup{}

	log.Printf("Sending event batches")

	var totalBatches int64 = 0

	for i, batch := range batches {
		wg.Add(1)

		go func(idx int, batch eventhub.BatchIterator) {
			err := hub.SendBatch(ctx, batch)

			if err != nil {
				log.Fatalf("ERROR sending batch: %s", err.Error())
			}

			wg.Done()
			atomic.AddInt64(&totalBatches, 1)
			log.Printf("[%d/%d] sent...", totalBatches, len(batches))
		}(i, batch)
	}

	wg.Wait()

	return nextTestId
}

func createEventBatch(testId *int64) eventhub.BatchIterator {
	var events []*eventhub.Event
	var data = [1024]byte{1}

	// simple minimum
	batchSize := 880

	for i := 0; i < batchSize; i++ {
		events = append(events, &eventhub.Event{
			Data: data[:],
			Properties: map[string]interface{}{
				TestIdProperty: *testId,
			},
		})

		*testId++
	}

	return eventhub.NewEventBatchIterator(events...)
}

func getPartitionCounts(ctx context.Context, hub *eventhub.Hub) map[string]*eventhub.HubPartitionRuntimeInformation {
	partitions := map[string]*eventhub.HubPartitionRuntimeInformation{}

	runtimeInfo, err := hub.GetRuntimeInformation(ctx)

	if err != nil {
		log.Fatalf("Failed to get runtime information from hub: %s", err.Error())
	}

	for _, partitionId := range runtimeInfo.PartitionIDs {
		partInfo, err := hub.GetPartitionInformation(ctx, partitionId)

		if err != nil {
			log.Fatalf("Failed to get partition info for partition ID %s: %s", partitionId, err.Error())
		}

		partitions[partitionId] = partInfo
	}

	return partitions
}

func verifyMessages(ctx context.Context, hub *eventhub.Hub, partitions map[string]*eventhub.HubPartitionRuntimeInformation, expectedMessages int64) {
	after := time.After(time.Minute * 5)

	receiverCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	messagesCh := make(chan int64, expectedMessages+10)

	for partitionID, partition := range partitions {
		go func(partitionID string, partition *eventhub.HubPartitionRuntimeInformation) {
			_, _ = hub.Receive(receiverCtx, partitionID, func(ctx context.Context, event *eventhub.Event) error {
				messagesCh <- event.Properties[TestIdProperty].(int64)

				return nil
			}, eventhub.ReceiveWithStartingOffset(partition.LastEnqueuedOffset))
		}(partitionID, partition)
	}

	log.Printf("Waiting for 5 minutes _or_ for %d unique messages to arrive", expectedMessages)

	mu := &sync.Mutex{}
	messagesReceived := map[int64]int64{}

	for i := 0; i < 5; i++ {
		go func() {
			for testId := range messagesCh {
				mu.Lock()
				messagesReceived[testId]++
				hits := messagesReceived[testId]
				length := len(messagesReceived)
				mu.Unlock()

				if hits > 1 {
					// we're getting duplicates
					log.Printf("Duplicate message with testId property %d", testId)
				}

				if int64(length) >= expectedMessages {
					log.Printf("Unique messages received: %d/%d", length, expectedMessages)
					cancel()
				}

				if length > 0 && length%1000 == 0 {
					log.Printf("Received %d messages", length)
				}
			}
		}()
	}

	select {
	case <-after:
		log.Printf("Timed out, didn't receive all messages")
		break
	case <-receiverCtx.Done():
		log.Printf("All messages received!")
		break
	}
}

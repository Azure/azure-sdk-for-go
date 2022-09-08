// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
)

const endProperty = "End"
const numProperty = "Number"
const partitionProperty = "Partition"

const messageLimit = 10000
const extraBytes = 1024

func FiniteSendAndReceiveTest(ctx context.Context) error {
	testData, err := newStressTestData("finite")

	if err != nil {
		return err
	}

	defer testData.Close()

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testData.ConnectionString, testData.HubName, nil)

	if err != nil {
		return err
	}

	defer func() { _ = producerClient.Close(context.Background()) }()

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testData.ConnectionString, testData.HubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return err
	}

	defer func() {
		err := consumerClient.Close(context.Background())

		if err != nil {
			panic(err)
		}
	}()

	ehProps, err := producerClient.GetEventHubProperties(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		err := producerClient.Close(context.Background())

		if err != nil {
			panic(err)
		}
	}()

	wg := sync.WaitGroup{}

	for _, tmpPartitionID := range ehProps.PartitionIDs {
		wg.Add(1)

		go func(partitionID string) {
			defer wg.Done()

			log.Printf("[p:%s] Starting to send messages to partition", partitionID)

			sp, err := sendEventsToPartition(context.Background(), producerClient, partitionID, messageLimit, extraBytes)

			if err != nil {
				log.Fatalf("Failed to send events to partition %s: %s", partitionID, err)
			}

			log.Printf("[p:%s] Done sending messages to partition", partitionID)

			log.Printf("[p:%s] Starting to receive messages from partition", partitionID)
			defer log.Printf("[p:%s] Done receiving messages from partition", partitionID)

			if err := consumeEventsFromPartition(context.Background(), consumerClient, sp, partitionID, messageLimit); err != nil {
				log.Fatalf("[p:%s] Failed to receive all events from partition: %s", partitionID, err)
			}
		}(tmpPartitionID)
	}

	wg.Wait()

	log.Printf("SUCCESS!")
	return nil
}

func consumeEventsFromPartition(ctx context.Context, consumerClient *azeventhubs.ConsumerClient, startPosition azeventhubs.StartPosition, partitionID string, numMessages int) error {
	log.Printf("[%s] Starting to consume events", partitionID)

	// read in 10 second chunks. If we ever end a 10 second chunk with no messages
	// then we've probably just failed.

	partitionClient, err := consumerClient.NewPartitionClient(partitionID, &azeventhubs.NewPartitionClientOptions{
		StartPosition: startPosition,
	})

	if err != nil {
		panic(err)
	}

	defer func() {
		err := partitionClient.Close(context.Background())

		if err != nil {
			panic(err)
		}
	}()

	sequenceNumbers := map[int64]int64{}

	numEmptyBatches := 0
	count := 0

	for {
		done, err := func() (bool, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			events, err := partitionClient.ReceiveEvents(ctx, 1000, nil)

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return false, fmt.Errorf("deadline exceeded, no messages arrived in 10 seconds. Got %d/%d messages", count, numMessages)
				}

				return false, err
			}

			count += len(events)

			log.Printf("[p:%s] Got %d/%d messages", partitionID, count, numMessages)

			if len(events) == 0 {
				numEmptyBatches++
			} else {
				numEmptyBatches = 0
			}

			if numEmptyBatches > 5 {
				// this is a lot of empty batches. During some of these tests, because of the sheer amount of activity, we can
				// get throttled. Giving about 50+ seconds would hopefully be enough to overcome that.
				return false, fmt.Errorf("%d empty batches in a row, link seems to be stuck. Got %d/%d messages", numEmptyBatches, count, numMessages)
			}

			for _, event := range events {
				_, exists := event.Properties[partitionProperty]

				if !exists {
					panic(fmt.Errorf("[%s] invalid message (seq: %d)- missing %s property: %#v", partitionID, event.SequenceNumber, partitionProperty, event.Properties))
				}

				eventPartitionID := event.Properties[partitionProperty].(string)
				num := event.Properties[numProperty].(int64)

				if eventPartitionID != partitionID {
					return false, fmt.Errorf("message with ID %s, with num %d, had a partition %s but we only were reading from partition %s", *event.MessageID, num, eventPartitionID, partitionID)
				}

				val, exists := sequenceNumbers[event.SequenceNumber]

				if !exists {
					sequenceNumbers[event.SequenceNumber] = num
				} else {
					return false, fmt.Errorf("duplicate message with sequence number %d was received. First occurrence with num %d, second with num %d", event.SequenceNumber, val, num)
				}

				expectedMessageCount, exists := event.Properties[endProperty].(int64)

				if exists {
					if len(sequenceNumbers) != int(expectedMessageCount) {
						return false, fmt.Errorf("end event was received but our count is off: expected %d, got %d", expectedMessageCount, len(sequenceNumbers))
					} else {
						return true, nil
					}
				}
			}

			return false, nil
		}()

		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}

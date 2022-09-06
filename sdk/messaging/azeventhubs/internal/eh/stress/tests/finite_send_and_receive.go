// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/joho/godotenv"
)

const endProperty = "End"
const numProperty = "Number"
const partitionProperty = "Partition"

const messageLimit = 100000
const extraBytes = 1024

func FiniteSendAndReceiveTest() error {
	azlog.SetEvents(azeventhubs.EventAuth, azeventhubs.EventConn, azeventhubs.EventConsumer)
	azlog.SetListener(func(e azlog.Event, s string) {
		log.Printf("[%s] %s", e, s)
	})

	defer func() {
		err := recover()

		if err != nil {
			log.Printf("FATAL ERROR: %s", err)
		}
	}()

	cs, hubName, err := loadEnvironment()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		fmt.Printf("Usage: stress\n")
		os.Exit(1)
	}

	allPartitions, err := getPartitions(cs, hubName)

	if err != nil {
		log.Printf("ERROR: failed to get partition IDs for test: %s", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	for _, partition := range allPartitions {
		wg.Add(1)
		go func(partition azeventhubs.PartitionProperties) {
			defer wg.Done()

			log.Printf("[p:%s] Starting to send messages to partition", partition.PartitionID)
			defer log.Printf("[p:%s] Done sending messages to partition", partition.PartitionID)

			err := sendEventsToPartition(cs, hubName, partition.PartitionID, messageLimit, extraBytes)

			if err != nil {
				log.Fatalf("Failed to send events to partition %s: %s", partition.PartitionID, err)
			}

			log.Printf("[p:%s] Starting to receive messages from partition", partition.PartitionID)
			defer log.Printf("[p:%s] Done receiving messages from partition", partition.PartitionID)

			if err := consumeEventsFromPartition(cs, hubName, partition, messageLimit); err != nil {
				log.Fatalf("[p:%s] Failed to receive all events from partition, started at %d: %s", partition.PartitionID, partition.LastEnqueuedSequenceNumber, err)
			}
		}(partition)
	}

	wg.Wait()

	log.Printf("SUCCESS!")
	return nil
}

func sendEventsToPartition(cs string, hubName string, partitionID string, messageLimit int, numExtraBytes int) error {
	log.Printf("Sending %d messages to partition ID %s, with messages of size %db", messageLimit, partitionID, numExtraBytes)

	extraBytes := make([]byte, numExtraBytes)

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(cs, hubName, nil)

	if err != nil {
		return err
	}

	batch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
		PartitionID: &partitionID,
	})

	if err != nil {
		return err
	}

	defer func() {
		if err := producerClient.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < messageLimit; i++ {
		ed := &azeventhubs.EventData{
			Body: extraBytes,
			ApplicationProperties: map[string]interface{}{
				numProperty:       i,
				partitionProperty: partitionID,
			},
		}

		if i == (messageLimit - 1) {
			ed.ApplicationProperties[endProperty] = messageLimit
		}

		err := batch.AddEventData(ed, nil)

		if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
			if batch.NumMessages() == 0 {
				return errors.New("single event was too large to fit into batch")
			}

			if err := producerClient.SendEventBatch(context.Background(), batch, nil); err != nil {
				return err
			}

			tempBatch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
				PartitionID: &partitionID,
			})

			if err != nil {
				return err
			}

			batch = tempBatch
			i-- // retry adding the same message
		} else if err != nil {
			return err
		}
	}

	if batch.NumMessages() > 0 {
		if err := producerClient.SendEventBatch(context.Background(), batch, nil); err != nil {
			return err
		}
	}

	return nil
}

func consumeEventsFromPartition(cs string, hubName string, partProps azeventhubs.PartitionProperties, numMessages int) error {
	log.Printf("Starting to consume events from %s, partitionID: %s, startingSequence: %d", hubName, partProps.PartitionID, partProps.LastEnqueuedSequenceNumber)

	startPosition := azeventhubs.StartPosition{
		//SequenceNumber: &firstPartition.LastEnqueuedSequenceNumber,
		SequenceNumber: to.Ptr(partProps.LastEnqueuedSequenceNumber),
		Inclusive:      false,
	}

	if partProps.IsEmpty {
		startPosition = azeventhubs.StartPosition{
			Earliest:  to.Ptr(true),
			Inclusive: true,
		}
	}

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(cs, hubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err := consumerClient.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	// read in 10 second chunks. If we ever end a 10 second chunk with no messages
	// then we've probably just failed.

	subscription, err := consumerClient.NewPartitionClient(partProps.PartitionID, &azeventhubs.NewPartitionClientOptions{
		StartPosition: startPosition,
	})

	if err != nil {
		panic(err)
	}

	sequenceNumbers := map[int64]int64{}

	numEmptyBatches := 0
	count := 0

	for {
		done, err := func() (bool, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			events, err := subscription.ReceiveEvents(ctx, 1000, nil)

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return false, fmt.Errorf("deadline exceeded, no messages arrived in 10 seconds. Got %d/%d messages", count, numMessages)
				}

				return false, err
			}

			count += len(events)

			log.Printf("[p:%s] Got %d/%d messages", partProps.PartitionID, count, numMessages)

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
				eventPartitionID := event.ApplicationProperties[partitionProperty].(string)
				num := event.ApplicationProperties[numProperty].(int64)

				if eventPartitionID != partProps.PartitionID {
					return false, fmt.Errorf("message with ID %s, with num %d, had a partition %s but we only were reading from partition %s", *event.MessageID, num, eventPartitionID, partProps.PartitionID)
				}

				val, exists := sequenceNumbers[event.SequenceNumber]

				if !exists {
					sequenceNumbers[event.SequenceNumber] = num
				} else {
					return false, fmt.Errorf("duplicate message with sequence number %d was received. First occurrence with num %d, second with num %d", event.SequenceNumber, val, num)
				}

				expectedMessageCount, exists := event.ApplicationProperties[endProperty].(int64)

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

func getPartitions(cs string, hubName string) (allPartitionProps []azeventhubs.PartitionProperties, err error) {
	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(cs, hubName, nil)

	if err != nil {
		return nil, err
	}

	defer func() { _ = producerClient.Close(context.Background()) }()

	props, err := producerClient.GetEventHubProperties(context.Background(), nil)

	if err != nil {
		return nil, err
	}

	for _, partitionID := range props.PartitionIDs {
		partProps, err := producerClient.GetPartitionProperties(context.Background(), partitionID, nil)

		if err != nil {
			return nil, err
		}

		allPartitionProps = append(allPartitionProps, partProps)
	}

	return allPartitionProps, nil
}

func loadEnvironment() (string, string, error) {
	envFilePath := ".env"

	if os.Getenv("ENV_FILE") != "" {
		envFilePath = os.Getenv("ENV_FILE")
	}

	if err := godotenv.Load(envFilePath); err != nil {
		return "", "", err
	}

	cs := os.Getenv("EVENTHUB_CONNECTION_STRING")
	hubName := os.Getenv("EVENTHUB_NAME")

	if cs == "" || hubName == "" {
		return "", "", fmt.Errorf("environment variables EVENTHUB_CONNECTION_STRING and EVENTHUB_NAME needs to be set")
	}

	return cs, hubName, nil
}

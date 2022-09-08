// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/blob"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

const inifiniteProcessorTestDuration = 5 * 24 * time.Hour

func FiniteProcessorTest(ctx context.Context) error {
	testData, err := newStressTestData("finiteprocessor")

	if err != nil {
		return err
	}

	defer testData.Close()

	blobNameBytes := make([]byte, 10)
	_, err = rand.Read(blobNameBytes)

	if err != nil {
		return err
	}

	containerName := "procstress-" + hex.EncodeToString(blobNameBytes)

	containerClient, err := blob.NewContainerClientFromConnectionString(testData.StorageConnectionString, containerName, nil)

	if err != nil {
		return err
	}

	_, err = containerClient.Create(ctx, nil)

	if err != nil {
		return err
	}

	checkpointStore, err := checkpoints.NewBlobStoreFromConnectionString(testData.StorageConnectionString, containerName, nil)

	if err != nil {
		return err
	}

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testData.ConnectionString, testData.HubName, azeventhubs.DefaultConsumerGroup, nil)

	if err != nil {
		return err
	}

	defer consumerClient.Close(ctx)

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testData.ConnectionString, testData.HubName, nil)

	if err != nil {
		return err
	}

	defer producerClient.Close(ctx)

	err = initCheckpointStore(ctx, producerClient, azeventhubs.CheckpointStoreAddress{
		ConsumerGroup:           azeventhubs.DefaultConsumerGroup,
		FullyQualifiedNamespace: testData.Namespace,
		EventHubName:            testData.HubName,
	}, checkpointStore)

	if err != nil {
		return err
	}

	const numMessages = 10000

	// start producing messages. We'll send out bursts, so we can validate that things are
	// progressing the way we expect.
	totalSent, err := sendEventsToAllPartitions(ctx, producerClient, numMessages, 500, testData.TC)

	if err != nil {
		return err
	}

	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, &azeventhubs.NewProcessorOptions{
		LoadBalancingStrategy: azeventhubs.ProcessorStrategyGreedy,
	})

	if err != nil {
		return err
	}

	processorCtx, cancelProcessor := context.WithTimeout(ctx, inifiniteProcessorTestDuration)
	defer cancelProcessor()

	var total int64 = 0

	go func() {
		for {
			partitionClient := processor.NextPartitionClient(ctx)

			if partitionClient == nil {
				return
			}

			go func() {
				defer partitionClient.Close(ctx)

				if err := processEvents(testData.TC, partitionClient, numMessages); err != nil {
					testData.TC.TrackException(fmt.Errorf("failed when processing events in ProcessorPartitionClient: %w", err))
				}

				if atomic.AddInt64(&total, numMessages) == int64(totalSent) {
					log.Printf("All messages received across all partitions!")
					cancelProcessor()
				}
			}()
		}
	}()

	if err := processor.Run(processorCtx); err != nil {
		return err
	}

	return nil
}

func processEvents(tc appinsights.TelemetryClient, partitionClient *azeventhubs.ProcessorPartitionClient, numMessages int) error {
	// TODO: calculate lag.
	log.Printf("[START] processing partition %s", partitionClient.PartitionID())
	defer log.Printf("[END] processing partition %s", partitionClient.PartitionID())

	total := 0

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		events, err := partitionClient.ReceiveEvents(ctx, 1000, nil)
		cancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		if len(events) > 0 {
			total += len(events)

			tc.TrackMetric("Received", float64(len(events)))

			err := partitionClient.UpdateCheckpoint(context.Background(), events[len(events)-1])

			if err != nil {
				return err
			}

			tc.TrackMetric("SequenceNumber", float64(events[len(events)-1].SequenceNumber))

			if total == numMessages {
				log.Printf("Finished receiving messages from %s", partitionClient.PartitionID())
				return nil
			}
		}
	}
}

func sendEventsToAllPartitions(ctx context.Context, producerClient *azeventhubs.ProducerClient, count int, extraBytes int, tc appinsights.TelemetryClient) (int, error) {
	ehProps, err := producerClient.GetEventHubProperties(ctx, nil)

	if err != nil {
		return 0, err
	}

	sendCtx, cancelSend := context.WithCancel(ctx)
	defer cancelSend()

	wg := sync.WaitGroup{}

	for tmpI, tmpPid := range ehProps.PartitionIDs {
		wg.Add(1)

		go func(i int, pid string) {
			defer wg.Done()

			_, err := sendEventsToPartition(sendCtx, producerClient, pid, count, extraBytes)

			if err != nil {
				tc.TrackException(fmt.Errorf("failed when sending %d events to partition %s: %w", count, pid, err))
				cancelSend()
			}

		}(tmpI, tmpPid)
	}

	wg.Wait()

	return len(ehProps.PartitionIDs) * count, ctx.Err()
}

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"context"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/blob"
)

func InfiniteProcessorTest(ctx context.Context) error {
	testData, err := newStressTestData("infinite")

	if err != nil {
		return err
	}

	defer testData.Close()

	blobNameBytes := make([]byte, 10)
	_, err = rand.Read(blobNameBytes)

	if err != nil {
		return err
	}

	containerName := "infprocstress-" + hex.EncodeToString(blobNameBytes)

	if err := testInitialize(ctx, testData, containerName); err != nil {
		return err
	}

	partitionToReceivedSequenceNums := map[string]trackingData{}

	{
		producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testData.ConnectionString, testData.HubName, nil)

		if err != nil {
			return err
		}

		defer func() {
			if err := producerClient.Close(context.Background()); err != nil {
				panic(err)
			}
		}()

		// pre-populate the map - from here nobody will modify the map, only
		// read it to get to the values in the `consumerData`. And those fields
		// should only be used atomically.
		ehProps, err := producerClient.GetEventHubProperties(ctx, nil)

		if err != nil {
			return err
		}

		for _, pid := range ehProps.PartitionIDs {
			partitionToReceivedSequenceNums[pid] = newTrackingData(pid)
		}

		// start producing goroutines
		for _, tempPartitionID := range ehProps.PartitionIDs {
			trackingData := partitionToReceivedSequenceNums[tempPartitionID]

			go func(partitionID string) {
				err := continuallySendEvents(ctx, producerClient, partitionID, trackingData.SentCount)

				if err != nil {
					panic(err)
				}
			}(tempPartitionID)
		}
	}

	allProcessorsDone := sync.WaitGroup{}

	for i := 0; i < 2; i++ {
		allProcessorsDone.Add(1)

		// start dispatching partitions, based on ownership
		go func(i int) {
			id := string([]byte{byte(i) + 'A'})

			checkpointStore, err := checkpoints.NewBlobStoreFromConnectionString(testData.StorageConnectionString, containerName, nil)

			if err != nil {
				panic(err)
			}

			consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testData.ConnectionString, testData.HubName, azeventhubs.DefaultConsumerGroup, nil)

			if err != nil {
				panic(err)
			}

			defer consumerClient.Close(ctx)

			processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, &azeventhubs.NewProcessorOptions{
				LoadBalancingStrategy: azeventhubs.ProcessorStrategyGreedy,
			})

			if err != nil {
				panic(err)
			}

			processorCtx, cancelProcessor := context.WithTimeout(ctx, inifiniteProcessorTestDuration)
			defer cancelProcessor()

			go func() {
				defer allProcessorsDone.Done()

				if err := processor.Run(processorCtx); err != nil {
					panic(err)
				}
			}()

			for {
				partitionClient := processor.NextPartitionClient(ctx)

				if partitionClient == nil {
					return
				}

				trackingData := partitionToReceivedSequenceNums[partitionClient.PartitionID()]

				go func() {
					defer partitionClient.Close(ctx)

					for {
						receiveCtx, cancelReceive := context.WithTimeout(ctx, 5*time.Second)
						events, err := partitionClient.ReceiveEvents(receiveCtx, 10000, nil)
						cancelReceive()

						if err != nil {
							if isCancelError(err) {
								return
							}

							var ehErr *azeventhubs.Error
							if errors.As(err, &ehErr) && ehErr.Code == azeventhubs.CodeOwnershipLost {
								log.Printf("[%s:%s] Lost ownership, stopping", id, partitionClient.PartitionID())
								return
							}

							panic(err)
						}

						for _, e := range events {
							trackingData.Inc(e.SequenceNumber)
						}

						// how far behind are we?
						stats := trackingData.Stats()

						testData.TC.TrackMetric("Received", float64(len(events)))
						testData.TC.TrackMetric("Lag", float64(stats.Sent-stats.Received))
						log.Printf("[%s:%s] just received %10d events, %10d behind, total: %10d", id, partitionClient.PartitionID(), len(events), stats.Sent-stats.Received, stats.Received)
					}
				}()
			}
		}(i)
	}

	allProcessorsDone.Wait()
	return nil
}

func continuallySendEvents(ctx context.Context, producerClient *azeventhubs.ProducerClient, partitionID string, sentCount *int64) error {
	const batchSize = 500

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		batch, err := producerClient.NewEventDataBatch(ctx, &azeventhubs.NewEventDataBatchOptions{
			PartitionID: &partitionID,
		})

		if err != nil {
			if isCancelError(err) {
				return nil
			}

			return err
		}

		base := atomic.LoadInt64(sentCount)

		for i := 0; i < batchSize; i++ {
			err := batch.AddEventData(&azeventhubs.EventData{
				Properties: map[string]any{
					"Index": base + int64(i),
				},
			}, nil)

			if err != nil {
				return err
			}
		}

		// add this _before_ we send otherwise we can end up in a funny situation
		// where the receiver gets the results before our send call in our local
		// process actually finishes, so it looks like they received messages
		// that were never sent. So we'll just increment before and things'll be fine.
		atomic.AddInt64(sentCount, batchSize)

		if err := producerClient.SendEventBatch(ctx, batch, nil); err != nil {
			atomic.AddInt64(sentCount, -batchSize)

			if isCancelError(err) {
				return nil
			}

			return err
		}

		time.Sleep(time.Second)
	}
}

func isCancelError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func testInitialize(ctx context.Context, testData *stressTestData, containerName string) error {
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

	return err
}

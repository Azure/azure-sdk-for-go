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

	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, &azeventhubs.NewProcessorOptions{
		LoadBalancingStrategy: azeventhubs.ProcessorStrategyGreedy,
	})

	if err != nil {
		return err
	}

	processorCtx, cancelProcessor := context.WithTimeout(ctx, inifiniteProcessorTestDuration)
	defer cancelProcessor()

	go func() {
		for {
			partitionClient := processor.NextPartitionClient(ctx)

			if partitionClient == nil {
				return
			}

			var numOutstanding int64
			var total int64

			receivedSequenceNumbers := map[int64]bool{}

			go func() {
				err := continuallySendEvents(ctx, producerClient, partitionClient.PartitionID(), &numOutstanding)

				if err != nil {
					panic(err)
				}
			}()

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

						panic(err)
					}

					for _, e := range events {
						if receivedSequenceNumbers[e.SequenceNumber] {
							panic(fmt.Errorf("got the same sequence number (%d) more than once for partition %s", e.SequenceNumber, partitionClient.PartitionID()))
						}
					}

					// how far behind are we?
					remaining := atomic.AddInt64(&numOutstanding, -int64(len(events)))
					total += int64(len(events))

					testData.TC.TrackMetric("Received", float64(len(events)))
					testData.TC.TrackMetric("Lag", float64(remaining))
					log.Printf("[%s] just received %d events, %d behind, total: %d", partitionClient.PartitionID(), len(events), remaining, total)
				}
			}()
		}
	}()

	if err := processor.Run(processorCtx); err != nil {
		return err
	}

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
		// where the receiver gets the results before our send call finishes and
		// we see negative numbers.
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

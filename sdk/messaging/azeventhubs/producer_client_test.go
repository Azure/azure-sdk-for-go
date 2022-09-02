// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestNewProducerClient_GetHubAndPartitionProperties(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producer, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	hubProps, err := producer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, hubProps.PartitionIDs)

	wg := sync.WaitGroup{}

	for _, partitionID := range hubProps.PartitionIDs {
		wg.Add(1)

		go func(pid string) {
			defer wg.Done()

			t.Run(fmt.Sprintf("Partition%s", pid), func(t *testing.T) {
				sendAndReceiveToPartitionTest(t, testParams.ConnectionString, testParams.EventHubName, pid)
			})
		}(partitionID)
	}

	wg.Wait()
}

func TestNewProducerClient_GetEventHubsProperties(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producer, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	props, err := producer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, props)
	require.NotEmpty(t, props.PartitionIDs)

	for _, pid := range props.PartitionIDs {
		props, err := producer.GetPartitionProperties(context.Background(), pid, nil)

		require.NoError(t, err)
		require.NotEmpty(t, props)

		require.Equal(t, pid, props.PartitionID)
	}
}

func TestNewProducerClient_SendToAny(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)
	partitions := getPartitions(t, testParams)

	require.NotNil(t, partitions)

	producer, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	batch, err := producer.NewEventDataBatch(context.Background(), nil)
	require.NoError(t, err)

	err = batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	err = producer.SendEventBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	for _, partProps := range partitions {
		wg.Add(1)

		go func(partProps azeventhubs.PartitionProperties) {
			defer wg.Done()

			consumer, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
			require.NoError(t, err)

			defer func() {
				err := consumer.Close(context.Background())
				require.NoError(t, err)
			}()

			partClient, err := consumer.NewPartitionClient(partProps.PartitionID, &azeventhubs.NewPartitionClientOptions{
				StartPosition: getStartPosition(partProps),
			})
			require.NoError(t, err)

			defer func() {
				err := partClient.Close(context.Background())
				require.NoError(t, err)
			}()

			events, err := partClient.ReceiveEvents(ctx, 1, nil)

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					return
				}

				require.NoError(t, err)
			}

			cancel()
			require.Equal(t, "hello world", string(events[0].Body))
		}(partProps)
	}

	wg.Wait()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}

func getPartitions(t *testing.T, testParams test.ConnectionParamsForTest) []azeventhubs.PartitionProperties {
	producer, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	hubProps, err := producer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	var partitions []azeventhubs.PartitionProperties

	for _, partitionID := range hubProps.PartitionIDs {
		partProps, err := producer.GetPartitionProperties(context.Background(), partitionID, nil)
		require.NoError(t, err)

		partitions = append(partitions, partProps)
	}

	sort.Slice(partitions, func(i, j int) bool {
		return partitions[i].PartitionID < partitions[j].PartitionID
	})

	return partitions
}

func sendAndReceiveToPartitionTest(t *testing.T, cs string, eventHubName string, partitionID string) {
	producer, err := azeventhubs.NewProducerClientFromConnectionString(cs, eventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	partProps, err := producer.GetPartitionProperties(context.Background(), partitionID, &azeventhubs.GetPartitionPropertiesOptions{})
	require.NoError(t, err)

	consumer, err := azeventhubs.NewConsumerClientFromConnectionString(cs, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
	require.NoError(t, err)

	defer func() {
		err := consumer.Close(context.Background())
		require.NoError(t, err)
	}()

	batch, err := producer.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
		PartitionID: &partitionID,
	})
	require.NoError(t, err)

	runID := time.Now().UnixNano()
	var expectedBodies []string

	for i := 0; i < 200; i++ {
		msg := fmt.Sprintf("%05d", i)

		err = batch.AddEventData(&azeventhubs.EventData{
			Body: []byte(msg),
			ApplicationProperties: map[string]any{
				"PartitionID": partitionID,
				"RunID":       runID,
			},
		}, nil)
		require.NoError(t, err)

		expectedBodies = append(expectedBodies, msg)
	}

	err = producer.SendEventBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	// give us 60 seconds to receive all 100 messages we sent in the batch
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var actualBodies []string

	subscription, err := consumer.NewPartitionClient(partitionID, &azeventhubs.NewPartitionClientOptions{
		StartPosition: getStartPosition(partProps),
	})
	require.NoError(t, err)

	for {
		events, err := subscription.ReceiveEvents(ctx, 100, nil)
		require.NoError(t, err)

		for _, event := range events {
			actualBodies = append(actualBodies, string(event.Body))

			require.Equal(t, partitionID, event.ApplicationProperties["PartitionID"], "No messages from other partitions")
			require.Equal(t, runID, event.ApplicationProperties["RunID"], "No messages from older runs")
		}

		if len(actualBodies) == len(expectedBodies) {
			break
		}
	}

	sort.Strings(actualBodies)
	require.Equal(t, expectedBodies, actualBodies)
}

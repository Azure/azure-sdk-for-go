// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
	"github.com/stretchr/testify/require"
)

func TestConsumerClient_DefaultAzureCredential(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	dac, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	t.Run("EventHubProperties and PartitionProperties", func(t *testing.T) {
		consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, dac, nil)
		require.NoError(t, err)

		defer func() {
			err := consumerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, dac, nil)
		require.NoError(t, err)

		defer func() {
			err := producerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		consumerProps, err := consumerClient.GetEventHubProperties(context.Background(), nil)
		require.NoError(t, err)

		producerProps, err := producerClient.GetEventHubProperties(context.Background(), nil)
		require.NoError(t, err)

		require.Equal(t, consumerProps, producerProps)

		producerPartProps, err := producerClient.GetPartitionProperties(context.Background(), consumerProps.PartitionIDs[0], nil)
		require.NoError(t, err)

		consumerPartProps, err := consumerClient.GetPartitionProperties(context.Background(), consumerProps.PartitionIDs[0], nil)
		require.NoError(t, err)

		require.Equal(t, producerPartProps, consumerPartProps)
	})

	t.Run("send and receive", func(t *testing.T) {
		producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, dac, nil)
		require.NoError(t, err)

		defer func() {
			err := producerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		firstPartition, err := producerClient.GetPartitionProperties(context.Background(), "0", nil)
		require.NoError(t, err)

		consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, dac, nil)
		require.NoError(t, err)

		defer func() {
			err := consumerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		eventDataBatch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
			PartitionID: to.Ptr(firstPartition.PartitionID),
		})
		require.NoError(t, err)

		err = eventDataBatch.AddEventData(&azeventhubs.EventData{
			Body: []byte("hello"),
		}, nil)
		require.NoError(t, err)

		err = producerClient.SendEventBatch(context.Background(), eventDataBatch, nil)
		require.NoError(t, err)

		subscription, err := consumerClient.NewPartitionClient(firstPartition.PartitionID, &azeventhubs.NewPartitionClientOptions{
			StartPosition: getStartPosition(firstPartition),
		})
		require.NoError(t, err)
		require.NotNil(t, subscription)

		defer func() {
			err := subscription.Close(context.Background())
			require.NoError(t, err)
		}()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		events, err := subscription.ReceiveEvents(ctx, 1, nil)
		require.NoError(t, err)

		require.Equal(t, "hello", string(events[0].Body))

		consumerPart, err := consumerClient.GetPartitionProperties(context.Background(), firstPartition.PartitionID, nil)
		require.NoError(t, err)
		producerPart, err := producerClient.GetPartitionProperties(context.Background(), firstPartition.PartitionID, nil)
		require.NoError(t, err)

		require.Equal(t, firstPartition.LastEnqueuedSequenceNumber+1, consumerPart.LastEnqueuedSequenceNumber)
		require.Equal(t, consumerPart, producerPart)
	})

	t.Run("EventHubProperties and PartitionProperties after send", func(t *testing.T) {
		consumerClient, err := azeventhubs.NewConsumerClient(testParams.EventHubNamespace, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, dac, nil)
		require.NoError(t, err)

		defer func() {
			err := consumerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		producerClient, err := azeventhubs.NewProducerClient(testParams.EventHubNamespace, testParams.EventHubName, dac, nil)
		require.NoError(t, err)

		defer func() {
			err := producerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		consumerProps, err := consumerClient.GetEventHubProperties(context.Background(), nil)
		require.NoError(t, err)

		producerProps, err := producerClient.GetEventHubProperties(context.Background(), nil)
		require.NoError(t, err)

		require.Equal(t, consumerProps, producerProps)

		producerPartProps, err := producerClient.GetPartitionProperties(context.Background(), consumerProps.PartitionIDs[0], nil)
		require.NoError(t, err)

		consumerPartProps, err := consumerClient.GetPartitionProperties(context.Background(), consumerProps.PartitionIDs[0], nil)
		require.NoError(t, err)

		require.Equal(t, producerPartProps, consumerPartProps)
	})
}

func TestConsumerClient_GetHubAndPartitionProperties(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	consumer, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
	require.NoError(t, err)

	defer func() {
		err := consumer.Close(context.Background())
		require.NoError(t, err)
	}()

	hubProps, err := consumer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, hubProps.PartitionIDs)

	for _, partitionID := range hubProps.PartitionIDs {
		props, err := consumer.GetPartitionProperties(context.Background(), partitionID, nil)
		require.NoError(t, err)

		require.Equal(t, partitionID, props.PartitionID)
	}
}

func TestConsumerClient_Concurrent_NoEpoch(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	partitions := mustSendEventsToAllPartitions(t, []*azeventhubs.EventData{
		{Body: []byte("hello world")},
	})

	const simultaneousClients = 5 // max you can have with a single consumer group for a single partition

	for i := 0; i < simultaneousClients; i++ {
		client, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, "$Default", nil)
		require.NoError(t, err)

		// We want all the clients open while this for loop is going.
		defer func() {
			err := client.Close(context.Background())
			require.NoError(t, err)
		}()

		partitionClient, err := client.NewPartitionClient(partitions[0].PartitionID, &azeventhubs.NewPartitionClientOptions{
			StartPosition: getStartPosition(partitions[0]),
		})
		require.NoError(t, err)

		defer func() {
			err := partitionClient.Close(context.Background())
			require.NoError(t, err)
		}()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		events, err := partitionClient.ReceiveEvents(ctx, 1, nil)
		require.NoError(t, err)

		require.Equal(t, 1, len(events))
	}
}

func TestConsumerClient_SameEpoch_StealsLink(t *testing.T) {
	partitions := mustSendEventsToAllPartitions(t, []*azeventhubs.EventData{
		{Body: []byte("hello world 1")},
	})

	ownerLevel := int64(2)

	origPartClient, cleanup := newPartitionClientForTest(t, partitions[0].PartitionID, azeventhubs.NewPartitionClientOptions{
		StartPosition: getStartPosition(partitions[0]),
		OwnerLevel:    &ownerLevel,
	})
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// open up a link, with an owner level of 2
	events, err := origPartClient.ReceiveEvents(ctx, 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	// link with owner level of 2 is alive, so now we'll steal it.

	thiefPartClient, cleanup := newPartitionClientForTest(t, partitions[0].PartitionID, azeventhubs.NewPartitionClientOptions{
		StartPosition: getStartPosition(partitions[0]),
		OwnerLevel:    &ownerLevel,
	})
	defer cleanup()

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	events, err = thiefPartClient.ReceiveEvents(ctx, 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	// the link has been stolen at this point - 'stealerPartClient' owns the link since it's last-in-wins.

	// using the original link reports that it was stolen
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	events, err = origPartClient.ReceiveEvents(ctx, 1, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "amqp:link:stolen")
	require.Empty(t, events)
}

func TestConsumerClient_LowerEpochsAreRejected(t *testing.T) {
	partitions := mustSendEventsToAllPartitions(t, []*azeventhubs.EventData{
		{Body: []byte("hello world 1")},
		{Body: []byte("hello world 2")},
	})

	highestOwnerLevel := int64(2)

	origPartClient, cleanup := newPartitionClientForTest(t, partitions[0].PartitionID, azeventhubs.NewPartitionClientOptions{
		StartPosition: getStartPosition(partitions[0]),
		OwnerLevel:    &highestOwnerLevel,
	})
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	events, err := origPartClient.ReceiveEvents(ctx, 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	lowerOwnerLevels := []*int64{
		nil, // no owner level
		to.Ptr(highestOwnerLevel - 1),
	}

	for _, ownerLevel := range lowerOwnerLevels {
		origPartClient, cleanup := newPartitionClientForTest(t, partitions[0].PartitionID, azeventhubs.NewPartitionClientOptions{
			StartPosition: getStartPosition(partitions[0]),
			OwnerLevel:    ownerLevel,
		})
		defer cleanup()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		events, err := origPartClient.ReceiveEvents(ctx, 1, nil)
		require.Error(t, err)
		// The typical error message is like this:
		//  At least one receiver for the endpoint is created with epoch of '2', and so non-epoch receiver is not allowed.
		//  Either reconnect with a higher epoch, or make sure all epoch receivers are closed or disconnected.
		require.Contains(t, err.Error(), "amqp:link:stolen")
		require.Empty(t, events)
	}

	// and the original client is unaffected
	events, err = origPartClient.ReceiveEvents(ctx, 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, events)
}

func newPartitionClientForTest(t *testing.T, partitionID string, subscribeOptions azeventhubs.NewPartitionClientOptions) (*azeventhubs.PartitionClient, func()) {
	testParams := test.GetConnectionParamsForTest(t)

	origClient, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, "$Default", &azeventhubs.ConsumerClientOptions{
		// Today we treat the link stolen error as retryable. I've filed an issue to look at making this fatal
		// instead since it's likely to be a configuration/runtime issue where the user has two consumers
		//  starting up with the same ownerlevel. Having them fight with retries is probably undesirable.
		RetryOptions: azeventhubs.RetryOptions{
			MaxRetries: -1,
		},
	})
	require.NoError(t, err)

	partClient, err := origClient.NewPartitionClient(partitionID, &subscribeOptions)
	require.NoError(t, err)

	return partClient, func() {
		err := partClient.Close(context.Background())
		require.NoError(t, err)

		err = origClient.Close(context.Background())
		require.NoError(t, err)
	}
}

func TestConsumerClient_StartPositions(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producerClient.Close(context.Background())
		require.NoError(t, err)
	}()

	batch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
		PartitionID: to.Ptr("0"),
	})
	require.NoError(t, err)

	require.NoError(t, batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("message 1"),
	}, nil))
	require.NoError(t, batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("message 2"),
	}, nil))

	origPartProps, err := producerClient.GetPartitionProperties(context.Background(), "0", nil)
	require.NoError(t, err)

	// introduce a little gap between any messages that are already in the eventhub and our new ones we're sending.
	// (this adds some peace of mind or the test below that uses the enqueued time for a filter)
	time.Sleep(time.Second)

	err = producerClient.SendEventBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	t.Run("offset", func(t *testing.T) {
		consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
		require.NoError(t, err)

		defer func() {
			err := consumerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		subscription, err := consumerClient.NewPartitionClient("0", &azeventhubs.NewPartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				Offset: &origPartProps.LastEnqueuedOffset,
			},
		})
		require.NoError(t, err)

		defer func() {
			err := subscription.Close(context.Background())
			require.NoError(t, err)
		}()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		events, err := subscription.ReceiveEvents(ctx, 2, nil)
		require.NoError(t, err)
		require.Equal(t, []string{"message 1", "message 2"}, getSortedBodies(events))
	})

	t.Run("enqueuedTime", func(t *testing.T) {
		consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
		require.NoError(t, err)

		defer func() {
			err := consumerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		subscription, err := consumerClient.NewPartitionClient("0", &azeventhubs.NewPartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				EnqueuedTime: &origPartProps.LastEnqueuedOn,
			},
		})
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		events, err := subscription.ReceiveEvents(ctx, 2, nil)
		require.NoError(t, err)
		require.Equal(t, []string{"message 1", "message 2"}, getSortedBodies(events))
	})

	t.Run("earliest", func(t *testing.T) {
		consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
		require.NoError(t, err)

		defer func() {
			err := consumerClient.Close(context.Background())
			require.NoError(t, err)
		}()

		subscription, err := consumerClient.NewPartitionClient("0", &azeventhubs.NewPartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				Earliest: to.Ptr(true),
			},
		})
		require.NoError(t, err)
		defer func() {
			err := subscription.Close(context.Background())
			require.NoError(t, err)
		}()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// we know there are _at_ two events but it's okay if they're just any events.
		events, err := subscription.ReceiveEvents(ctx, 2, nil)
		require.NoError(t, err)
		require.Equal(t, 2, len(events))
	})
}

func TestConsumerClient_StartPosition_Latest(t *testing.T) {
	testParams := test.GetConnectionParamsForTest(t)

	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, azeventhubs.DefaultConsumerGroup, nil)
	require.NoError(t, err)

	defer func() {
		err := consumerClient.Close(context.Background())
		require.NoError(t, err)
	}()

	// warm up the AMQP connection underneath. The link will be created when I start doing the receive.
	_, err = consumerClient.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	latestEventsCh := make(chan []*azeventhubs.ReceivedEventData, 1)

	go func() {
		subscription, err := consumerClient.NewPartitionClient("0", &azeventhubs.NewPartitionClientOptions{
			StartPosition: azeventhubs.StartPosition{
				Latest: to.Ptr(true),
			},
		})
		require.NoError(t, err)

		defer func() {
			err := subscription.Close(context.Background())
			require.NoError(t, err)
		}()

		events, err := subscription.ReceiveEvents(context.Background(), 2, nil)
		require.NoError(t, err)
		latestEventsCh <- events
	}()

	// give the consumer link time to spin up and start listening on the partition
	time.Sleep(5 * time.Second)

	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producerClient.Close(context.Background())
		require.NoError(t, err)
	}()

	batch, err := producerClient.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
		PartitionID: to.Ptr("0"),
	})
	require.NoError(t, err)

	require.NoError(t, batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("latest test: message 1"),
	}, nil))
	require.NoError(t, batch.AddEventData(&azeventhubs.EventData{
		Body: []byte("latest test: message 2"),
	}, nil))

	err = producerClient.SendEventBatch(context.Background(), batch, nil)
	require.NoError(t, err)

	select {
	case events := <-latestEventsCh:
		require.Equal(t, []string{"latest test: message 1", "latest test: message 2"}, getSortedBodies(events))
	case <-time.After(time.Minute):
		require.Fail(t, "Timed out waiting for events to arrrive")
	}
}

// mustSendEventsToAllPartitions sends the event given in evt to each partition in the
// eventHub, returning the sequence number just before the new message.
//
// This is useful for tests that need to work with a hub that might already have messages, and need
// to start from a particular sequence number to avoid them.
//
// NOTE: the message that's passed in does get altered so don't count on it being unchanged after calling
// this function. Each message gets an additional property (DestPartitionID), set to the parttion ID that
// we sent it to.
func mustSendEventsToAllPartitions(t *testing.T, events []*azeventhubs.EventData) []azeventhubs.PartitionProperties {
	testParams := test.GetConnectionParamsForTest(t)
	producer, err := azeventhubs.NewProducerClientFromConnectionString(testParams.ConnectionString, testParams.EventHubName, nil)
	require.NoError(t, err)

	defer func() {
		err := producer.Close(context.Background())
		require.NoError(t, err)
	}()

	hubProps, err := producer.GetEventHubProperties(context.Background(), nil)
	require.NoError(t, err)

	partitionsCh := make(chan azeventhubs.PartitionProperties, len(hubProps.PartitionIDs))

	wg := sync.WaitGroup{}
	wg.Add(len(hubProps.PartitionIDs))

	for _, partitionID := range hubProps.PartitionIDs {
		go func(partitionID string) {
			defer wg.Done()

			partProps, err := producer.GetPartitionProperties(context.Background(), partitionID, nil)
			require.NoError(t, err)
			partitionsCh <- partProps

			// send the message to the partition.
			batch, err := producer.NewEventDataBatch(context.Background(), &azeventhubs.NewEventDataBatchOptions{
				PartitionID: &partitionID,
			})
			require.NoError(t, err)

			for _, event := range events {
				if event.Properties == nil {
					event.Properties = map[string]any{}
				}

				event.Properties["DestPartitionID"] = partitionID

				err = batch.AddEventData(event, nil)
				require.NoError(t, err)
			}

			err = producer.SendEventBatch(context.Background(), batch, nil)
			require.NoError(t, err)
		}(partitionID)
	}

	wg.Wait()
	close(partitionsCh)

	var partitions []azeventhubs.PartitionProperties

	for p := range partitionsCh {
		partitions = append(partitions, p)
	}

	return partitions
}

func getStartPosition(props azeventhubs.PartitionProperties) azeventhubs.StartPosition {
	if props.IsEmpty {
		return azeventhubs.StartPosition{
			Earliest:  to.Ptr(true),
			Inclusive: true,
		}
	}

	return azeventhubs.StartPosition{
		SequenceNumber: to.Ptr(props.LastEnqueuedSequenceNumber),
		Inclusive:      false,
	}
}

func getSortedBodies(events []*azeventhubs.ReceivedEventData) []string {
	var bodies []string

	for _, e := range events {
		bodies = append(bodies, string(e.Body))
	}

	return bodies
}
